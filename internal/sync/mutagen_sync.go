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

	// Check if source file exists
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		result.Error = fmt.Errorf("source file not found: %v", err)
		return result, result.Error
	}

	// Check if target file exists
	targetInfo, err := os.Stat(targetPath)
	if err != nil && !os.IsNotExist(err) {
		result.Error = fmt.Errorf("failed to stat target file: %v", err)
		return result, result.Error
	}

	// If target doesn't exist, just copy the file
	if os.IsNotExist(err) {
		if options.DryRun {
			if options.Verbose {
				fmt.Printf("Would copy %s to %s (new file)\n", sourcePath, targetPath)
			}
			result.BytesTransferred = sourceInfo.Size()
			result.Operations = 1
			result.Success = true
			return result, nil
		}

		// Create target directory if it doesn't exist
		targetDir := filepath.Dir(targetPath)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			result.Error = fmt.Errorf("failed to create target directory: %v", err)
			return result, result.Error
		}

		// Simple copy for new files
		if err := copyFile(sourcePath, targetPath, options); err != nil {
			result.Error = fmt.Errorf("failed to copy file: %v", err)
			return result, result.Error
		}

		result.BytesTransferred = sourceInfo.Size()
		result.Operations = 1
		result.Success = true
		return result, nil
	}

	// Files exist, check if they're identical
	if sourceInfo.Size() == targetInfo.Size() {
		sourceChecksum, err := FileChecksum(sourcePath)
		if err != nil {
			result.Error = fmt.Errorf("failed to compute source checksum: %v", err)
			return result, result.Error
		}

		targetChecksum, err := FileChecksum(targetPath)
		if err != nil {
			result.Error = fmt.Errorf("failed to compute target checksum: %v", err)
			return result, result.Error
		}

		if bytes.Equal(sourceChecksum, targetChecksum) {
			if options.Verbose {
				fmt.Printf("Files are identical, no sync needed\n")
			}
			result.Success = true
			return result, nil
		}
	}

	// Files are different, use Mutagen's rsync algorithm
	if options.DryRun {
		if options.Verbose {
			fmt.Printf("Would sync %s to %s using Mutagen rsync algorithm\n", sourcePath, targetPath)
		}
		result.BytesTransferred = sourceInfo.Size()
		result.Operations = 1
		result.Success = true
		return result, nil
	}

	// Create backup if requested
	if options.CreateBackup {
		backupPath := targetPath + ".backup"
		if err := copyFile(targetPath, backupPath, nil); err != nil {
			result.Error = fmt.Errorf("failed to create backup: %v", err)
			return result, result.Error
		}
	}

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
			copyBytes += int64(op.Count) * int64(signature.BlockSize)
			if op.Start == uint64(len(signature.Hashes)-1) {
				// Last block might be shorter
				copyBytes -= int64(signature.BlockSize - signature.LastBlockSize)
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
