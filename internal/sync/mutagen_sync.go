package sync

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mutagen-io/mutagen/pkg/synchronization/rsync"
)

// MutagenSyncEngine wraps Mutagen's rsync engine for efficient file synchronization
type MutagenSyncEngine struct {
	engine *rsync.Engine
}

// NewMutagenSyncEngine creates a new sync engine using Mutagen's rsync implementation
func NewMutagenSyncEngine() *MutagenSyncEngine {
	return &MutagenSyncEngine{
		engine: rsync.NewEngine(),
	}
}

// SyncFile efficiently synchronizes a source file to a target location using Mutagen's rsync
func (m *MutagenSyncEngine) SyncFile(sourcePath, targetPath string, options *SyncOptions) (*FileSyncResult, error) {
	if options == nil {
		options = DefaultSyncOptions()
	}

	result := &FileSyncResult{
		SourcePath: sourcePath,
		TargetPath: targetPath,
	}

	// Validate source file
	sourceInfo, err := m.validateSourceFile(sourcePath, result)
	if err != nil {
		return result, err
	}

	// Check target file status
	targetExists, targetInfo, err := m.checkTargetFile(targetPath, result)
	if err != nil {
		return result, err
	}

	// Handle new file case
	if !targetExists {
		return m.handleNewFile(sourcePath, targetPath, options, result, sourceInfo)
	}

	// Check if files are identical
	if m.filesAreIdentical(sourcePath, targetPath, result, sourceInfo, targetInfo, options) {
		return result, nil
	}

	// Files are different, sync them
	return m.syncDifferentFiles(sourcePath, targetPath, options, result, sourceInfo)
}

// validateSourceFile validates that the source file exists and is accessible
func (m *MutagenSyncEngine) validateSourceFile(sourcePath string, result *FileSyncResult) (os.FileInfo, error) {
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		result.Error = fmt.Errorf("source file not found: %v", err)
		return nil, result.Error
	}
	return sourceInfo, nil
}

// checkTargetFile checks the status of the target file
func (m *MutagenSyncEngine) checkTargetFile(targetPath string, result *FileSyncResult) (bool, os.FileInfo, error) {
	targetInfo, err := os.Stat(targetPath)
	if err != nil && !os.IsNotExist(err) {
		result.Error = fmt.Errorf("failed to stat target file: %v", err)
		return false, nil, result.Error
	}
	return !os.IsNotExist(err), targetInfo, nil
}

// handleNewFile handles the case where the target file doesn't exist
func (m *MutagenSyncEngine) handleNewFile(sourcePath, targetPath string, options *SyncOptions, result *FileSyncResult, sourceInfo os.FileInfo) (*FileSyncResult, error) {
	if options.DryRun {
		return m.handleDryRunNewFile(sourcePath, targetPath, options, result, sourceInfo)
	}

	// Create target directory if it doesn't exist
	if err := m.createTargetDirectory(targetPath); err != nil {
		result.Error = err
		return result, err
	}

	// Copy the new file
	if err := m.copyNewFile(sourcePath, targetPath, options, result, sourceInfo); err != nil {
		result.Error = err
		return result, err
	}

	return result, nil
}

// handleDryRunNewFile handles dry run for new files
func (m *MutagenSyncEngine) handleDryRunNewFile(sourcePath, targetPath string, options *SyncOptions, result *FileSyncResult, sourceInfo os.FileInfo) (*FileSyncResult, error) {
	if options.Verbose {
		fmt.Printf("Would copy %s to %s (new file)\n", sourcePath, targetPath)
	}
	result.BytesTransferred = sourceInfo.Size()
	result.Operations = 1
	result.Success = true
	return result, nil
}

// createTargetDirectory creates the target directory if it doesn't exist
func (m *MutagenSyncEngine) createTargetDirectory(targetPath string) error {
	targetDir := filepath.Dir(targetPath)
	return os.MkdirAll(targetDir, 0o755)
}

// copyNewFile copies a new file to the target location
func (m *MutagenSyncEngine) copyNewFile(sourcePath, targetPath string, options *SyncOptions, result *FileSyncResult, sourceInfo os.FileInfo) error {
	if err := copyFile(sourcePath, targetPath, options); err != nil {
		return fmt.Errorf("failed to copy file: %v", err)
	}

	result.BytesTransferred = sourceInfo.Size()
	result.Operations = 1
	result.Success = true
	return nil
}

// filesAreIdentical checks if source and target files are identical
func (m *MutagenSyncEngine) filesAreIdentical(sourcePath, targetPath string, result *FileSyncResult, sourceInfo, targetInfo os.FileInfo, options *SyncOptions) bool {
	// Quick size check first
	if sourceInfo.Size() != targetInfo.Size() {
		return false
	}

	// Compare checksums
	sourceChecksum, err := FileChecksum(sourcePath)
	if err != nil {
		result.Error = fmt.Errorf("failed to compute source checksum: %v", err)
		return false
	}

	targetChecksum, err := FileChecksum(targetPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to compute target checksum: %v", err)
		return false
	}

	if bytes.Equal(sourceChecksum, targetChecksum) {
		if options.Verbose {
			fmt.Printf("Files are identical, no sync needed\n")
		}
		result.Success = true
		return true
	}

	return false
}

// syncDifferentFiles handles syncing files that are different
func (m *MutagenSyncEngine) syncDifferentFiles(sourcePath, targetPath string, options *SyncOptions, result *FileSyncResult, sourceInfo os.FileInfo) (*FileSyncResult, error) {
	if options.DryRun {
		return m.handleDryRunSync(sourcePath, targetPath, options, result, sourceInfo)
	}

	// Create backup if requested
	if err := m.createBackupIfNeeded(targetPath, options, result); err != nil {
		return result, err
	}

	// Perform the actual sync using Mutagen's rsync algorithm
	return m.performMutagenSync(sourcePath, targetPath, options, result, sourceInfo)
}

// handleDryRunSync handles dry run for file syncing
func (m *MutagenSyncEngine) handleDryRunSync(sourcePath, targetPath string, options *SyncOptions, result *FileSyncResult, sourceInfo os.FileInfo) (*FileSyncResult, error) {
	if options.Verbose {
		fmt.Printf("Would sync %s to %s using Mutagen rsync algorithm\n", sourcePath, targetPath)
	}
	result.BytesTransferred = sourceInfo.Size()
	result.Operations = 1
	result.Success = true
	return result, nil
}

// createBackupIfNeeded creates a backup of the target file if requested
func (m *MutagenSyncEngine) createBackupIfNeeded(targetPath string, options *SyncOptions, result *FileSyncResult) error {
	if !options.CreateBackup {
		return nil
	}

	backupPath := targetPath + ".backup"
	if err := copyFile(targetPath, backupPath, nil); err != nil {
		return fmt.Errorf("failed to create backup: %v", err)
	}

	if options.Verbose {
		fmt.Printf("Created backup: %s\n", backupPath)
	}

	return nil
}

// performMutagenSync performs the actual file synchronization using Mutagen's rsync algorithm
func (m *MutagenSyncEngine) performMutagenSync(sourcePath, targetPath string, options *SyncOptions, result *FileSyncResult, sourceInfo os.FileInfo) (*FileSyncResult, error) {
	if options.Verbose {
		fmt.Printf("Using Mutagen rsync to sync files...\n")
	}

	// Open source file for reading
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		result.Error = fmt.Errorf("failed to open source file: %v", err)
		return result, result.Error
	}
	defer sourceFile.Close()

	// Open target file for reading
	targetFile, err := os.Open(targetPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to open target file: %v", err)
		return result, result.Error
	}
	defer targetFile.Close()

	// Generate signature of target file (base)
	signature, err := m.engine.Signature(targetFile, uint64(options.BlockLength))
	if err != nil {
		result.Error = fmt.Errorf("failed to generate target signature: %v", err)
		return result, result.Error
	}

	if options.Verbose {
		fmt.Printf("Generated signature with %d blocks\n", len(signature.Hashes))
	}

	// Create delta operations by comparing source with target signature
	var operations []*rsync.Operation
	var literalBytes, copyBytes int64

	// Collect operations
	transmit := func(op *rsync.Operation) error {
		operations = append(operations, &rsync.Operation{
			Data:  append([]byte(nil), op.Data...), // Copy data
			Start: op.Start,
			Count: op.Count,
		})

		// Track statistics
		if len(op.Data) > 0 {
			literalBytes += int64(len(op.Data))
		} else {
			// Check for potential overflow before conversion
			if op.Count > 0 && signature.BlockSize > 0 {
				if uint64(op.Count) <= (1<<63-1)/uint64(signature.BlockSize) {
					copyBytes += int64(op.Count) * int64(signature.BlockSize)
				} else {
					// Handle overflow case
					copyBytes = 1<<63 - 1
				}
			}
			if op.Start == uint64(len(signature.Hashes)-1) {
				// Last block might be shorter
				if signature.BlockSize >= signature.LastBlockSize {
					copyBytes -= int64(signature.BlockSize - signature.LastBlockSize)
				}
			}
		}
		return nil
	}

	// Reset source file for deltification
	if _, err := sourceFile.Seek(0, io.SeekStart); err != nil {
		result.Error = fmt.Errorf("failed to seek source file: %v", err)
		return result, result.Error
	}

	// Generate delta operations
	if err := m.engine.Deltify(sourceFile, signature, rsync.DefaultMaximumDataOperationSize, transmit); err != nil {
		result.Error = fmt.Errorf("failed to create delta: %v", err)
		return result, result.Error
	}

	result.BytesTransferred = literalBytes
	result.BytesSaved = copyBytes
	result.Operations = len(operations)

	if options.Verbose {
		fmt.Printf("Delta: %d literal bytes, %d copy bytes, %d operations\n",
			literalBytes, copyBytes, result.Operations)
	}

	// Apply delta to create new target file
	tempPath := targetPath + ".tmp"
	outputFile, err := os.Create(tempPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to create temp file: %v", err)
		return result, result.Error
	}
	defer outputFile.Close()

	// Reset target file for patching
	if _, err := targetFile.Seek(0, io.SeekStart); err != nil {
		result.Error = fmt.Errorf("failed to seek target file: %v", err)
		return result, result.Error
	}

	// Apply all operations to reconstruct the source file
	for _, op := range operations {
		if err := m.engine.Patch(outputFile, targetFile, signature, op); err != nil {
			result.Error = fmt.Errorf("failed to apply operation: %v", err)
			return result, result.Error
		}
	}

	outputFile.Close()

	// Replace target file with new version
	if err := os.Rename(tempPath, targetPath); err != nil {
		result.Error = fmt.Errorf("failed to replace target file: %v", err)
		return result, result.Error
	}

	// Preserve file attributes if requested
	if options.PreservePerms || options.PreserveOwner || options.PreserveGroup {
		if err := preserveAttributes(sourcePath, targetPath, options); err != nil {
			result.Error = fmt.Errorf("failed to preserve attributes: %v", err)
			return result, result.Error
		}
	}

	result.Success = true
	return result, nil
}
