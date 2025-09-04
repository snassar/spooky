package sync

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"spooky/internal/logging"
	"spooky/internal/schemas"
	"spooky/internal/ssh"
	"spooky/internal/utilities"
	"strings"

	"github.com/pkg/errors"
)

// RemoteSyncEngine handles remote file synchronization using SSH and Mutagen's rsync engine
type RemoteSyncEngine struct {
	sshManager    *ssh.SimpleSSHManager
	mutagenEngine *MutagenSyncEngine
}

// NewRemoteSyncEngine creates a new remote sync engine
func NewRemoteSyncEngine(sshManager *ssh.SimpleSSHManager) *RemoteSyncEngine {
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

// Progress represents synchronization progress
type Progress struct {
	CurrentFile      string
	FilesProcessed   int
	TotalFiles       int
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

// syncOneWayReplicaRemote performs exact one-way replication (local → remote)
func (r *RemoteSyncEngine) syncOneWayReplicaRemote(ctx context.Context, localPath, remotePath string, options *RemoteOptions, result *RemoteSyncResult, progress *Progress) error {
	logger := logging.GetGlobalLogger()
	logger.Info("performing one-way replica sync to remote",
		slog.String("source", localPath),
		slog.String("target", remotePath),
		slog.String("machine", options.Machine.Hostname))

	// Walk local directory and sync all files
	err := filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
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

		if info.IsDir() {
			// Create directory on remote
			if err := r.createRemoteDirectory(ctx, options.Machine, remoteFile); err != nil {
				return fmt.Errorf("failed to create remote directory %s: %v", remoteFile, err)
			}
			if options.Verbose {
				logger.Info("created remote directory", slog.String("path", remoteFile))
			}
		} else {
			// Sync file to remote
			fileResult, err := r.syncFileToRemote(ctx, path, remoteFile, options, progress)
			if err != nil {
				return fmt.Errorf("failed to sync file %s: %v", path, err)
			}

			result.BytesTransferred += fileResult.BytesTransferred
			result.BytesSaved += fileResult.BytesSaved
			result.Operations += fileResult.Operations
		}

		return nil
	})

	if err != nil {
		result.Error = errors.Wrap(err, "failed to sync directory")
		return result.Error
	}

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

// cleanupRemoteFiles removes files on remote that don't exist locally
func (r *RemoteSyncEngine) cleanupRemoteFiles(ctx context.Context, localPath, remotePath string, options *RemoteOptions, progress *Progress) error {
	// This would involve:
	// 1. Getting a list of files on remote
	// 2. Comparing with local files
	// 3. Removing files that don't exist locally

	// For now, this is a placeholder implementation
	logger := logging.GetGlobalLogger()
	logger.Info("cleanup not yet implemented")

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
	// This would involve:
	// 1. Sync local → remote (one-way safe)
	// 2. Sync remote → local (one-way safe)
	// 3. Detect and report conflicts

	// For now, this is a placeholder implementation
	logger := logging.GetGlobalLogger()
	logger.Info("two-way safe sync not yet implemented")

	return nil
}

// syncTwoWayResolvedRemote performs bidirectional sync with local winning conflicts
func (r *RemoteSyncEngine) syncTwoWayResolvedRemote(ctx context.Context, localPath, remotePath string, options *RemoteOptions, result *RemoteSyncResult, progress *Progress) error {
	// This would involve:
	// 1. Sync local → remote (one-way replica - local always wins)
	// 2. Sync remote → local (one-way safe - preserve local changes)

	// For now, this is a placeholder implementation
	logger := logging.GetGlobalLogger()
	logger.Info("two-way resolved sync not yet implemented")

	return nil
}
