package sync

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"spooky/internal/logging"
	"spooky/internal/schemas"
	"spooky/internal/ssh"
	"spooky/internal/utilities"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
)

// RemoteSyncEngine handles remote file synchronization using SSH and Mutagen's rsync engine
type RemoteSyncEngine struct {
	sshManager    *ssh.Manager
	mutagenEngine *MutagenSyncEngine
}

// NewRemoteSyncEngine creates a new remote sync engine
func NewRemoteSyncEngine(sshManager *ssh.Manager) *RemoteSyncEngine {
	return &RemoteSyncEngine{
		sshManager:    sshManager,
		mutagenEngine: NewMutagenSyncEngine(),
	}
}

// RemoteOptions configures remote synchronization behavior
type RemoteOptions struct {
	*Options
	Machine         *schemas.MachinesMachineV1
	ProgressReport  func(progress *Progress)
	ConflictResolve ConflictResolution
	SyncDelete      bool // Delete files in destination that don't exist in source
}

// Progress represents synchronization progress with thread-safe operations
type Progress struct {
	mu               sync.RWMutex
	CurrentFile      string
	FilesProcessed   int64
	TotalFiles       int64
	BytesTransferred int64
	BytesSaved       int64
	CurrentOperation string
	Percentage       float64
}

// ConflictResolution defines how to handle conflicts
type ConflictResolution string

const (
	// ConflictResolutionSkip skips conflicted files during sync.
	ConflictResolutionSkip ConflictResolution = "skip"
	// ConflictResolutionBackup creates backups of conflicted files during sync.
	ConflictResolutionBackup ConflictResolution = "backup"
	// ConflictResolutionOverwrite overwrites conflicted files during sync.
	ConflictResolutionOverwrite ConflictResolution = "overwrite"
	// ConflictResolutionPrompt prompts user for conflict resolution.
	ConflictResolutionPrompt ConflictResolution = "prompt"
)

// FileProcessingJob represents a file to be processed by a worker goroutine
type FileProcessingJob struct {
	LocalPath  string
	RemotePath string
	IsDir      bool
}

// ConcurrentSyncResult holds thread-safe results from concurrent operations
type ConcurrentSyncResult struct {
	mu               sync.Mutex
	BytesTransferred int64
	BytesSaved       int64
	Operations       int
	Errors           []error
}

// AddResult safely adds a file sync result to the concurrent result
func (c *ConcurrentSyncResult) AddResult(result *FileSyncResult) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.BytesTransferred += result.BytesTransferred
	c.BytesSaved += result.BytesSaved
	c.Operations += result.Operations
}

// AddError safely adds an error to the concurrent result
func (c *ConcurrentSyncResult) AddError(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Errors = append(c.Errors, err)
}

// GetErrors returns a copy of all errors
func (c *ConcurrentSyncResult) GetErrors() []error {
	c.mu.Lock()
	defer c.mu.Unlock()

	errors := make([]error, len(c.Errors))
	copy(errors, c.Errors)
	return errors
}

// IncrementFilesProcessed atomically increments the files processed counter
func (p *Progress) IncrementFilesProcessed() {
	atomic.AddInt64(&p.FilesProcessed, 1)
}

// SetCurrentFile safely sets the current file being processed
func (p *Progress) SetCurrentFile(filename string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.CurrentFile = filename
}

// SetCurrentOperation safely sets the current operation
func (p *Progress) SetCurrentOperation(operation string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.CurrentOperation = operation
}

// UpdatePercentage safely updates the percentage
func (p *Progress) UpdatePercentage() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.TotalFiles > 0 {
		p.Percentage = float64(p.FilesProcessed) / float64(p.TotalFiles) * 100
	}
}

// GetProgressSnapshot returns a thread-safe snapshot of current progress
func (p *Progress) GetProgressSnapshot() Progress {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return Progress{
		CurrentFile:      p.CurrentFile,
		FilesProcessed:   atomic.LoadInt64(&p.FilesProcessed),
		TotalFiles:       p.TotalFiles,
		BytesTransferred: atomic.LoadInt64(&p.BytesTransferred),
		BytesSaved:       atomic.LoadInt64(&p.BytesSaved),
		CurrentOperation: p.CurrentOperation,
		Percentage:       p.Percentage,
	}
}

// RemoteSyncResult contains the result of remote synchronization
type RemoteSyncResult struct {
	*FileSyncResult
	RemoteMachine     string
	Conflicts         []string
	ResolvedConflicts []string
	Progress          *Progress
}

// SyncRemoteDirectory synchronizes a local directory with a remote directory using rsync over SSH
func (r *RemoteSyncEngine) SyncRemoteDirectory(ctx context.Context, localPath, remotePath string, options *RemoteOptions) (*RemoteSyncResult, error) {
	if options == nil {
		options = &RemoteOptions{
			Options: DefaultOptions(),
		}
	}

	result := &RemoteSyncResult{
		FileSyncResult: &FileSyncResult{
			SourcePath: localPath,
			TargetPath: remotePath,
		},
		RemoteMachine: options.Machine.Hostname,
	}

	// Validate local path
	localInfo, err := os.Stat(localPath)
	if err != nil {
		result.Error = fmt.Errorf("local path not found: %v", err)
		return result, result.Error
	}

	if !localInfo.IsDir() {
		result.Error = fmt.Errorf("local path is not a directory: %s", localPath)
		return result, result.Error
	}

	// Initialize progress tracking
	progress := &Progress{
		CurrentFile:      "",
		TotalFiles:       0,
		CurrentOperation: "scanning",
	}
	result.Progress = progress

	// Count total files for progress reporting
	if err := r.countFiles(localPath, progress); err != nil {
		result.Error = fmt.Errorf("failed to count files: %v", err)
		return result, result.Error
	}

	// Create remote directory if it doesn't exist
	if err := r.createRemoteDirectory(ctx, options.Machine, remotePath); err != nil {
		result.Error = fmt.Errorf("failed to create remote directory: %v", err)
		return result, result.Error
	}

	// Perform synchronization based on mode
	switch options.Mode {
	case ModeOneWayReplica:
		err = r.syncOneWayReplicaRemote(ctx, localPath, remotePath, options, result, progress)
	case ModeOneWaySafe:
		err = r.syncOneWaySafeRemote(ctx, localPath, remotePath, options, result, progress)
	case ModeTwoWaySafe:
		err = r.syncTwoWaySafeRemote(ctx, localPath, remotePath, options, result, progress)
	case ModeTwoWayResolved:
		err = r.syncTwoWayResolvedRemote(ctx, localPath, remotePath, options, result, progress)
	default:
		result.Error = fmt.Errorf("unknown sync mode: %s", options.Mode)
		return result, result.Error
	}

	if err != nil {
		result.Error = err
		return result, err
	}

	return result, nil
}

// countFiles counts the total number of files to be synchronized
func (r *RemoteSyncEngine) countFiles(path string, progress *Progress) error {
	return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			progress.TotalFiles++
		}
		return nil
	})
}

// createRemoteDirectory creates the remote directory if it doesn't exist
func (r *RemoteSyncEngine) createRemoteDirectory(ctx context.Context, machine *schemas.MachinesMachineV1, remotePath string) error {
	command := fmt.Sprintf("mkdir -p %s", remotePath)

	result, err := r.sshManager.RunCommandOnMachine(ctx, machine, command)
	if err != nil {
		return fmt.Errorf("failed to create remote directory: %v", err)
	}

	if result.ExitCode != 0 {
		return fmt.Errorf("failed to create remote directory, exit code: %d", result.ExitCode)
	}

	return nil
}

// getOptimalConcurrency determines the optimal number of concurrent workers
func (r *RemoteSyncEngine) getOptimalConcurrency(options *RemoteOptions) int {
	if options.MaxConcurrency > 0 {
		return options.MaxConcurrency
	}

	// Auto-detect based on CPU cores and system resources
	// For I/O bound operations like file sync, we can use more workers than CPU cores
	cpuCores := runtime.NumCPU()
	optimal := cpuCores * 2 // Start with 2x CPU cores for I/O bound work

	// Cap at reasonable maximum to avoid overwhelming the system
	if optimal > 10 {
		optimal = 10
	}

	// Minimum of 1 worker
	if optimal < 1 {
		optimal = 1
	}

	return optimal
}

// processFileConcurrently handles individual file processing in a worker goroutine
func (r *RemoteSyncEngine) processFileConcurrently(ctx context.Context, job FileProcessingJob, options *RemoteOptions, progress *Progress, concurrentResult *ConcurrentSyncResult) {
	if job.IsDir {
		// Handle directory creation
		if err := r.createRemoteDirectory(ctx, options.Machine, job.RemotePath); err != nil {
			concurrentResult.AddError(fmt.Errorf("failed to create remote directory %s: %v", job.RemotePath, err))
			return
		}

		if options.Verbose {
			logger := logging.GetGlobalLogger()
			logger.Info("created remote directory", slog.String("path", job.RemotePath))
		}
		return
	}

	// Update progress for file processing
	progress.SetCurrentFile(filepath.Base(job.LocalPath))
	progress.SetCurrentOperation("syncing")

	// Handle file synchronization
	fileResult, err := r.syncFileToRemote(ctx, job.LocalPath, job.RemotePath, options, progress)
	if err != nil {
		concurrentResult.AddError(fmt.Errorf("failed to sync file %s: %v", job.LocalPath, err))
		return
	}

	// Update progress atomically
	progress.IncrementFilesProcessed()
	progress.UpdatePercentage()

	concurrentResult.AddResult(fileResult)
}

// syncOneWayReplicaRemote performs exact one-way replication (local → remote) with parallel processing
func (r *RemoteSyncEngine) syncOneWayReplicaRemote(ctx context.Context, localPath, remotePath string, options *RemoteOptions, result *RemoteSyncResult, progress *Progress) error {
	logger := logging.GetGlobalLogger()
	logger.Info("performing one-way replica sync to remote with parallel processing",
		slog.String("source", localPath),
		slog.String("target", remotePath),
		slog.String("machine", options.Machine.Hostname))

	// Determine optimal concurrency level
	maxWorkers := r.getOptimalConcurrency(options)
	logger.Info("using parallel processing", slog.Int("workers", maxWorkers))

	// Create channels for job distribution and results
	jobChan := make(chan FileProcessingJob, maxWorkers*2) // Buffered channel
	concurrentResult := &ConcurrentSyncResult{}
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobChan {
				// Process each file/directory concurrently
				r.processFileConcurrently(ctx, job, options, progress, concurrentResult)
			}
		}(i)
	}

	// Walk directory and send jobs to workers
	walkErr := make(chan error, 1)
	go func() {
		defer close(jobChan)
		walkErr <- filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return errors.Wrapf(err, "failed to walk local path %s", path)
			}

			// Calculate relative path from local source
			relPath, err := filepath.Rel(localPath, path)
			if err != nil {
				return fmt.Errorf("failed to calculate relative path: %v", err)
			}

			// Skip root directory
			if relPath == "." {
				return nil
			}

			remoteFile := filepath.Join(remotePath, relPath)

			// Send job to worker
			job := FileProcessingJob{
				LocalPath:  path,
				RemotePath: remoteFile,
				IsDir:      info.IsDir(),
			}
			jobChan <- job

			return nil
		})
	}()

	// Wait for directory walk to complete
	if err := <-walkErr; err != nil {
		// Close job channel to stop workers
		close(jobChan)
		wg.Wait()
		result.Error = errors.Wrap(err, "failed to walk directory")
		return result.Error
	}

	// Wait for all workers to complete
	wg.Wait()

	// Check for errors from concurrent operations
	if errors := concurrentResult.GetErrors(); len(errors) > 0 {
		// Log all errors but continue with cleanup
		for _, err := range errors {
			logger.Error("concurrent sync error", slog.String("error", err.Error()))
		}
		result.Error = fmt.Errorf("encountered %d errors during parallel sync", len(errors))
	}

	// Update result with aggregated data
	result.BytesTransferred += concurrentResult.BytesTransferred
	result.BytesSaved += concurrentResult.BytesSaved
	result.Operations += concurrentResult.Operations

	// Clean up files on remote that don't exist locally if delete is enabled
	if options.SyncDelete {
		if err := r.cleanupRemoteFiles(ctx, localPath, remotePath, options, progress); err != nil {
			result.Error = errors.Wrap(err, "failed to cleanup remote files")
			return result.Error
		}
	}

	result.Success = true
	return nil
}

// syncFileToRemote synchronizes a single file to the remote machine using Mutagen's rsync engine
func (r *RemoteSyncEngine) syncFileToRemote(ctx context.Context, localPath, remotePath string, options *RemoteOptions, progress *Progress) (*FileSyncResult, error) {
	result := &FileSyncResult{
		SourcePath: localPath,
		TargetPath: remotePath,
	}

	// Update progress and setup cleanup
	tempRemotePath := r.setupProgressAndCleanup(localPath, options, progress)
	defer r.cleanupTempFile(tempRemotePath)

	// Check if remote file exists and download it if needed
	remoteExists, err := r.downloadRemoteFileIfExists(ctx, options.Machine, remotePath, tempRemotePath)
	if err != nil {
		result.Error = fmt.Errorf("failed to check/download remote file: %v", err)
		return result, result.Error
	}

	// Handle conflict detection in safe mode
	if remoteExists && options.Mode == ModeOneWaySafe {
		if conflict, err := r.checkForConflicts(localPath, tempRemotePath, remotePath, options); err != nil {
			result.Error = err
			return result, result.Error
		} else if conflict {
			result.Conflicts = append(result.Conflicts, remotePath)
			return result, nil
		}
	}

	// Perform file synchronization
	if err := r.performFileSync(localPath, tempRemotePath, remoteExists, options, result); err != nil {
		result.Error = err
		return result, result.Error
	}

	// Upload the synced file to the remote machine
	if err := r.uploadFileToRemote(ctx, tempRemotePath, remotePath, options); err != nil {
		result.Error = fmt.Errorf("failed to upload synced file: %v", err)
		return result, result.Error
	}

	// Update progress and mark success
	r.updateProgress(progress, options)
	result.Success = true
	return result, nil
}

// setupProgressAndCleanup initializes progress tracking and returns temp file path
func (r *RemoteSyncEngine) setupProgressAndCleanup(localPath string, options *RemoteOptions, progress *Progress) string {
	progress.CurrentFile = filepath.Base(localPath)
	progress.CurrentOperation = "syncing"
	if options.ProgressReport != nil {
		options.ProgressReport(progress)
	}
	return localPath + ".remote"
}

// cleanupTempFile removes the temporary remote file
func (r *RemoteSyncEngine) cleanupTempFile(tempRemotePath string) {
	if removeErr := os.Remove(tempRemotePath); removeErr != nil {
		// Log the error but don't fail the function since this is cleanup
		// This is a best-effort cleanup
		utilities.HandleCleanupError(removeErr, "temporary_file", "remove")
	}
}

// checkForConflicts checks for file conflicts in safe mode
func (r *RemoteSyncEngine) checkForConflicts(localPath, tempRemotePath, remotePath string, options *RemoteOptions) (bool, error) {
	localInfo, err := os.Stat(localPath)
	if err != nil {
		return false, fmt.Errorf("failed to stat local file: %v", err)
	}

	remoteInfo, err := os.Stat(tempRemotePath)
	if err != nil {
		return false, fmt.Errorf("failed to stat remote file: %v", err)
	}

	if remoteInfo.ModTime().After(localInfo.ModTime()) {
		// Remote is newer, preserve it (conflict)
		if options.Verbose {
			logger := logging.GetGlobalLogger()
			logger.Info("conflict detected, preserving remote file", slog.String("path", remotePath))
		}
		return true, nil
	}

	return false, nil
}

// performFileSync handles the actual file synchronization logic
func (r *RemoteSyncEngine) performFileSync(localPath, tempRemotePath string, remoteExists bool, options *RemoteOptions, result *FileSyncResult) error {
	if remoteExists {
		return r.syncExistingFile(localPath, tempRemotePath, options, result)
	}
	return r.copyNewFile(localPath, tempRemotePath, options, result)
}

// syncExistingFile syncs an existing file using Mutagen
func (r *RemoteSyncEngine) syncExistingFile(localPath, tempRemotePath string, options *RemoteOptions, result *FileSyncResult) error {
	syncResult, err := r.mutagenEngine.File(localPath, tempRemotePath, options.Options)
	if err != nil {
		return fmt.Errorf("failed to sync file using Mutagen: %v", err)
	}

	// Copy sync statistics
	result.BytesTransferred = syncResult.BytesTransferred
	result.BytesSaved = syncResult.BytesSaved
	result.Operations = syncResult.Operations
	return nil
}

// copyNewFile copies a new file and sets up statistics
func (r *RemoteSyncEngine) copyNewFile(localPath, tempRemotePath string, options *RemoteOptions, result *FileSyncResult) error {
	if err := copyFile(localPath, tempRemotePath, options.Options); err != nil {
		return fmt.Errorf("failed to copy new file: %v", err)
	}

	// Get file size for statistics
	if info, err := os.Stat(localPath); err == nil {
		result.BytesTransferred = info.Size()
		result.Operations = 1
	}
	return nil
}

// updateProgress updates the progress counter and reports if needed
func (r *RemoteSyncEngine) updateProgress(progress *Progress, options *RemoteOptions) {
	progress.FilesProcessed++
	if options.ProgressReport != nil {
		options.ProgressReport(progress)
	}
}

// downloadRemoteFileIfExists downloads a remote file to a local path if it exists
func (r *RemoteSyncEngine) downloadRemoteFileIfExists(ctx context.Context, machine *schemas.MachinesMachineV1, remotePath, localPath string) (bool, error) {
	// Check if remote file exists
	command := fmt.Sprintf("test -f %s && echo 'exists'", remotePath)
	result, err := r.sshManager.RunCommandOnMachine(ctx, machine, command)
	if err != nil {
		return false, fmt.Errorf("failed to check remote file existence: %v", err)
	}

	if result.ExitCode != 0 || strings.TrimSpace(result.Stdout) != "exists" {
		return false, nil // File doesn't exist
	}

	// Download the file using scp-like approach
	downloadCmd := fmt.Sprintf("cat %s", remotePath)
	result, err = r.sshManager.RunCommandOnMachine(ctx, machine, downloadCmd)
	if err != nil {
		return false, fmt.Errorf("failed to download remote file: %v", err)
	}

	if result.ExitCode != 0 {
		return false, fmt.Errorf("download command failed with exit code %d", result.ExitCode)
	}

	// Write the content to local file
	if err := os.WriteFile(localPath, []byte(result.Stdout), 0644); err != nil {
		return false, fmt.Errorf("failed to write local file: %v", err)
	}

	return true, nil
}

// uploadFileToRemote uploads a local file to the remote machine
func (r *RemoteSyncEngine) uploadFileToRemote(ctx context.Context, localPath, remotePath string, options *RemoteOptions) error {
	// Read the local file content
	content, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("failed to read local file: %v", err)
	}

	// Escape the content for shell command
	escapedContent := strings.ReplaceAll(string(content), "'", "'\"'\"'")

	// Upload using echo and redirection (simple but effective for small files)
	uploadCmd := fmt.Sprintf("echo '%s' > %s", escapedContent, remotePath)

	result, err := r.sshManager.RunCommandOnMachine(ctx, options.Machine, uploadCmd)
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	if result.ExitCode != 0 {
		return fmt.Errorf("upload command failed with exit code %d", result.ExitCode)
	}

	// Set file permissions if requested
	if options.PreservePerms {
		if info, err := os.Stat(localPath); err == nil {
			chmodCmd := fmt.Sprintf("chmod %o %s", info.Mode().Perm(), remotePath)
			result, err := r.sshManager.RunCommandOnMachine(ctx, options.Machine, chmodCmd)
			if err != nil {
				logger := logging.GetGlobalLogger()
				logger.Warn("failed to set remote file permissions", slog.String("error", err.Error()))
			} else if result.ExitCode != 0 {
				logger := logging.GetGlobalLogger()
				logger.Warn("chmod command failed", slog.Int("exit_code", result.ExitCode))
			}
		}
	}

	return nil
}

// listLocalFiles returns a list of all files in the local directory
func (r *RemoteSyncEngine) listLocalFiles(localPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// Convert to relative path from localPath
			relPath, err := filepath.Rel(localPath, path)
			if err != nil {
				return err
			}
			files = append(files, relPath)
		}
		return nil
	})

	return files, err
}

// removeRemoteFile safely removes a file from the remote system
func (r *RemoteSyncEngine) removeRemoteFile(ctx context.Context, machine *schemas.MachinesMachineV1, remoteFilePath string) error {
	// Use rm command to remove the file
	command := fmt.Sprintf("rm -f %s", remoteFilePath)
	result, err := r.sshManager.RunCommandOnMachine(ctx, machine, command)
	if err != nil {
		return fmt.Errorf("failed to remove remote file %s: %v", remoteFilePath, err)
	}

	if result.ExitCode != 0 {
		return fmt.Errorf("failed to remove remote file %s, exit code: %d", remoteFilePath, result.ExitCode)
	}

	return nil
}

// cleanupRemoteFiles removes files on remote that don't exist locally
func (r *RemoteSyncEngine) cleanupRemoteFiles(ctx context.Context, localPath, remotePath string, options *RemoteOptions, progress *Progress) error {
	logger := logging.GetGlobalLogger()

	// Check if cleanup is enabled
	if !options.SyncDelete {
		logger.Debug("cleanup disabled, skipping remote file cleanup")
		return nil
	}

	// Update progress
	if progress != nil {
		progress.CurrentOperation = "cleanup"
	}

	// Get list of local files
	localFiles, err := r.listLocalFiles(localPath)
	if err != nil {
		return fmt.Errorf("failed to list local files: %v", err)
	}

	// Create a map for faster lookup
	localFileMap := make(map[string]bool)
	for _, file := range localFiles {
		localFileMap[file] = true
	}

	// Get list of remote files
	remoteFiles, err := r.listRemoteFiles(ctx, options.Machine, remotePath)
	if err != nil {
		return fmt.Errorf("failed to list remote files: %v", err)
	}

	// Find files to remove (exist on remote but not locally)
	var filesToRemove []string
	for _, remoteFile := range remoteFiles {
		// Convert remote file path to relative path from remotePath
		relPath, err := filepath.Rel(remotePath, remoteFile)
		if err != nil {
			logger.Warn("failed to get relative path for remote file", slog.String("file", remoteFile), slog.Any("error", err))
			continue
		}

		// Check if this file exists locally
		if !localFileMap[relPath] {
			filesToRemove = append(filesToRemove, remoteFile)
		}
	}

	// Remove files that don't exist locally
	removedCount := 0
	for _, fileToRemove := range filesToRemove {
		if options.DryRun {
			logger.Info("would remove remote file", slog.String("file", fileToRemove))
			removedCount++
		} else {
			err := r.removeRemoteFile(ctx, options.Machine, fileToRemove)
			if err != nil {
				logger.Error("failed to remove remote file", slog.String("file", fileToRemove), slog.Any("error", err))
				// Continue with other files even if one fails
				continue
			}
			logger.Debug("removed remote file", slog.String("file", fileToRemove))
			removedCount++
		}

		// Update progress
		if progress != nil {
			progress.FilesProcessed++
		}
	}

	logger.Info("remote file cleanup completed",
		slog.Int("files_removed", removedCount),
		slog.Int("total_remote_files", len(remoteFiles)),
		slog.Int("total_local_files", len(localFiles)))

	return nil
}

// syncOneWaySafeRemote performs safe one-way sync (local → remote)
func (r *RemoteSyncEngine) syncOneWaySafeRemote(ctx context.Context, localPath, remotePath string, options *RemoteOptions, result *RemoteSyncResult, progress *Progress) error {
	// Similar to one-way replica but with conflict detection
	// For now, delegate to the replica implementation
	return r.syncOneWayReplicaRemote(ctx, localPath, remotePath, options, result, progress)
}

// syncTwoWaySafeRemote performs bidirectional sync with conflict detection
func (r *RemoteSyncEngine) syncTwoWaySafeRemote(ctx context.Context, localPath, remotePath string, options *RemoteOptions, result *RemoteSyncResult, progress *Progress) error {
	logger := logging.GetGlobalLogger()
	logger.Info("performing two-way safe sync",
		slog.String("local", localPath),
		slog.String("remote", remotePath),
		slog.String("machine", options.Machine.Hostname))

	// Step 1: Sync local → remote (one-way safe)
	logger.Info("syncing local to remote (safe mode)")
	localToRemoteResult := &RemoteSyncResult{
		FileSyncResult: &FileSyncResult{
			SourcePath: localPath,
			TargetPath: remotePath,
		},
		RemoteMachine: options.Machine.Hostname,
		Progress:      progress,
	}

	if err := r.syncOneWaySafeRemote(ctx, localPath, remotePath, options, localToRemoteResult, progress); err != nil {
		return fmt.Errorf("failed to sync local to remote: %v", err)
	}

	// Step 2: Sync remote → local (one-way safe)
	logger.Info("syncing remote to local (safe mode)")
	remoteToLocalResult := &RemoteSyncResult{
		FileSyncResult: &FileSyncResult{
			SourcePath: remotePath,
			TargetPath: localPath,
		},
		RemoteMachine: options.Machine.Hostname,
		Progress:      progress,
	}

	if err := r.syncRemoteToLocalSafe(ctx, remotePath, localPath, options, remoteToLocalResult, progress); err != nil {
		return fmt.Errorf("failed to sync remote to local: %v", err)
	}

	// Step 3: Combine results and detect conflicts
	result.BytesTransferred = localToRemoteResult.BytesTransferred + remoteToLocalResult.BytesTransferred
	result.BytesSaved = localToRemoteResult.BytesSaved + remoteToLocalResult.BytesSaved
	result.Operations = localToRemoteResult.Operations + remoteToLocalResult.Operations
	result.Conflicts = append(localToRemoteResult.Conflicts, remoteToLocalResult.Conflicts...)
	result.Success = localToRemoteResult.Success && remoteToLocalResult.Success

	if len(result.Conflicts) > 0 {
		logger.Warn("conflicts detected during two-way safe sync",
			slog.Int("conflict_count", len(result.Conflicts)),
			slog.String("conflicts", fmt.Sprintf("%v", result.Conflicts)))
	}

	return nil
}

// syncTwoWayResolvedRemote performs bidirectional sync with local winning conflicts
func (r *RemoteSyncEngine) syncTwoWayResolvedRemote(ctx context.Context, localPath, remotePath string, options *RemoteOptions, result *RemoteSyncResult, progress *Progress) error {
	logger := logging.GetGlobalLogger()
	logger.Info("performing two-way resolved sync (local wins conflicts)",
		slog.String("local", localPath),
		slog.String("remote", remotePath),
		slog.String("machine", options.Machine.Hostname))

	// Step 1: Sync local → remote (one-way replica - local always wins)
	logger.Info("syncing local to remote (replica mode - local wins)")
	localToRemoteResult := &RemoteSyncResult{
		FileSyncResult: &FileSyncResult{
			SourcePath: localPath,
			TargetPath: remotePath,
		},
		RemoteMachine: options.Machine.Hostname,
		Progress:      progress,
	}

	if err := r.syncOneWayReplicaRemote(ctx, localPath, remotePath, options, localToRemoteResult, progress); err != nil {
		return fmt.Errorf("failed to sync local to remote: %v", err)
	}

	// Step 2: Sync remote → local (one-way safe - preserve local changes)
	logger.Info("syncing remote to local (safe mode - preserve local changes)")
	remoteToLocalResult := &RemoteSyncResult{
		FileSyncResult: &FileSyncResult{
			SourcePath: remotePath,
			TargetPath: localPath,
		},
		RemoteMachine: options.Machine.Hostname,
		Progress:      progress,
	}

	if err := r.syncRemoteToLocalSafe(ctx, remotePath, localPath, options, remoteToLocalResult, progress); err != nil {
		return fmt.Errorf("failed to sync remote to local: %v", err)
	}

	// Step 3: Combine results (conflicts are resolved by local winning)
	result.BytesTransferred = localToRemoteResult.BytesTransferred + remoteToLocalResult.BytesTransferred
	result.BytesSaved = localToRemoteResult.BytesSaved + remoteToLocalResult.BytesSaved
	result.Operations = localToRemoteResult.Operations + remoteToLocalResult.Operations
	result.Conflicts = append(localToRemoteResult.Conflicts, remoteToLocalResult.Conflicts...)
	result.Success = localToRemoteResult.Success && remoteToLocalResult.Success

	// Log resolved conflicts
	if len(result.Conflicts) > 0 {
		logger.Info("conflicts resolved during two-way sync (local wins)",
			slog.Int("resolved_count", len(result.Conflicts)),
			slog.String("resolved_conflicts", fmt.Sprintf("%v", result.Conflicts)))
		result.ResolvedConflicts = result.Conflicts
		result.Conflicts = nil // Clear conflicts since they're resolved
	}

	return nil
}

// syncRemoteToLocalSafe performs safe one-way sync from remote to local
func (r *RemoteSyncEngine) syncRemoteToLocalSafe(ctx context.Context, remotePath, localPath string, options *RemoteOptions, result *RemoteSyncResult, progress *Progress) error {
	logger := logging.GetGlobalLogger()
	logger.Info("performing one-way safe sync from remote to local",
		slog.String("source", remotePath),
		slog.String("target", localPath),
		slog.String("machine", options.Machine.Hostname))

	// Create local directory if it doesn't exist
	if err := os.MkdirAll(localPath, 0755); err != nil {
		return fmt.Errorf("failed to create local directory: %v", err)
	}

	// Get list of remote files
	remoteFiles, err := r.listRemoteFiles(ctx, options.Machine, remotePath)
	if err != nil {
		return fmt.Errorf("failed to list remote files: %v", err)
	}

	// Sync each remote file to local
	for _, remoteFile := range remoteFiles {
		relPath, err := filepath.Rel(remotePath, remoteFile)
		if err != nil {
			return fmt.Errorf("failed to calculate relative path: %v", err)
		}

		localFile := filepath.Join(localPath, relPath)

		// Create local directory if needed
		localDir := filepath.Dir(localFile)
		if err := os.MkdirAll(localDir, 0755); err != nil {
			return fmt.Errorf("failed to create local directory %s: %v", localDir, err)
		}

		// Sync file from remote to local
		fileResult, err := r.syncFileFromRemote(ctx, remoteFile, localFile, options, progress)
		if err != nil {
			return fmt.Errorf("failed to sync file %s: %v", remoteFile, err)
		}

		result.BytesTransferred += fileResult.BytesTransferred
		result.BytesSaved += fileResult.BytesSaved
		result.Operations += fileResult.Operations
		result.Conflicts = append(result.Conflicts, fileResult.Conflicts...)
	}

	// Clean up local files that don't exist on remote if delete is enabled
	if options.SyncDelete {
		if err := r.cleanupLocalFiles(ctx, remotePath, localPath, options, progress); err != nil {
			return fmt.Errorf("failed to cleanup local files: %v", err)
		}
	}

	result.Success = true
	return nil
}

// listRemoteFiles lists all files in a remote directory
func (r *RemoteSyncEngine) listRemoteFiles(ctx context.Context, machine *schemas.MachinesMachineV1, remotePath string) ([]string, error) {
	// Use find command to list all files recursively
	command := fmt.Sprintf("find %s -type f 2>/dev/null || true", remotePath)
	result, err := r.sshManager.RunCommandOnMachine(ctx, machine, command)
	if err != nil {
		return nil, fmt.Errorf("failed to list remote files: %v", err)
	}

	if result.ExitCode != 0 {
		return nil, fmt.Errorf("failed to list remote files, exit code: %d", result.ExitCode)
	}

	// Parse the output and filter out empty lines
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")
	var files []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			files = append(files, strings.TrimSpace(line))
		}
	}

	return files, nil
}

// syncFileFromRemote synchronizes a single file from remote to local
func (r *RemoteSyncEngine) syncFileFromRemote(ctx context.Context, remotePath, localPath string, options *RemoteOptions, progress *Progress) (*FileSyncResult, error) {
	result := &FileSyncResult{
		SourcePath: remotePath,
		TargetPath: localPath,
	}

	// Update progress
	progress.CurrentFile = filepath.Base(remotePath)
	progress.CurrentOperation = "downloading"

	// Check if local file exists
	localExists := false
	if localInfo, err := os.Stat(localPath); err == nil {
		localExists = true
		if !localInfo.IsDir() {
			// Handle conflict detection in safe mode
			if options.Mode == ModeOneWaySafe || options.Mode == ModeTwoWaySafe {
				if conflict, err := r.checkForConflictsRemoteToLocal(ctx, remotePath, localPath, options); err != nil {
					result.Error = err
					return result, result.Error
				} else if conflict {
					result.Conflicts = append(result.Conflicts, localPath)
					return result, nil
				}
			}
		}
	}

	// Download the remote file
	content, err := r.downloadRemoteFileContent(ctx, options.Machine, remotePath)
	if err != nil {
		result.Error = fmt.Errorf("failed to download remote file: %v", err)
		return result, result.Error
	}

	// Write to local file
	if err := os.WriteFile(localPath, content, 0644); err != nil {
		result.Error = fmt.Errorf("failed to write local file: %v", err)
		return result, result.Error
	}

	// Update statistics
	result.BytesTransferred = int64(len(content))
	result.Operations = 1
	if localExists {
		result.BytesSaved = 0 // No compression for simple download
	}

	// Update progress
	progress.CurrentOperation = "completed"
	result.Success = true
	return result, nil
}

// downloadRemoteFileContent downloads the content of a remote file
func (r *RemoteSyncEngine) downloadRemoteFileContent(ctx context.Context, machine *schemas.MachinesMachineV1, remotePath string) ([]byte, error) {
	command := fmt.Sprintf("cat %s", remotePath)
	result, err := r.sshManager.RunCommandOnMachine(ctx, machine, command)
	if err != nil {
		return nil, fmt.Errorf("failed to download remote file: %v", err)
	}

	if result.ExitCode != 0 {
		return nil, fmt.Errorf("download command failed with exit code %d", result.ExitCode)
	}

	return []byte(result.Stdout), nil
}

// checkForConflictsRemoteToLocal checks for conflicts when syncing from remote to local
func (r *RemoteSyncEngine) checkForConflictsRemoteToLocal(ctx context.Context, remotePath, localPath string, options *RemoteOptions) (bool, error) {
	// Download remote file content
	remoteContent, err := r.downloadRemoteFileContent(ctx, options.Machine, remotePath)
	if err != nil {
		return false, fmt.Errorf("failed to download remote file for conflict check: %v", err)
	}

	// Read local file content
	localContent, err := os.ReadFile(localPath)
	if err != nil {
		return false, fmt.Errorf("failed to read local file for conflict check: %v", err)
	}

	// Compare content
	return !bytes.Equal(remoteContent, localContent), nil
}

// cleanupLocalFiles removes local files that don't exist on remote
func (r *RemoteSyncEngine) cleanupLocalFiles(ctx context.Context, remotePath, localPath string, options *RemoteOptions, progress *Progress) error {
	// Get list of remote files
	remoteFiles, err := r.listRemoteFiles(ctx, options.Machine, remotePath)
	if err != nil {
		return fmt.Errorf("failed to list remote files for cleanup: %v", err)
	}

	// Create a set of remote files for quick lookup
	remoteFileSet := make(map[string]bool)
	for _, remoteFile := range remoteFiles {
		relPath, err := filepath.Rel(remotePath, remoteFile)
		if err != nil {
			continue
		}
		remoteFileSet[relPath] = true
	}

	// Walk local directory and remove files not in remote
	return filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(localPath, path)
		if err != nil {
			return err
		}

		// Skip if file exists on remote
		if remoteFileSet[relPath] {
			return nil
		}

		// Remove local file that doesn't exist on remote
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("failed to remove local file %s: %v", path, err)
		}

		if options.Verbose {
			logger := logging.GetGlobalLogger()
			logger.Info("removed local file not present on remote", slog.String("path", path))
		}

		return nil
	})
}
