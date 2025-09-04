package sync

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"

	"spooky/internal/logging"
	"spooky/internal/utilities"

	"log/slog"

	"github.com/mutagen-io/mutagen/pkg/synchronization/rsync"
)

// MutagenSyncEngine wraps Mutagen's rsync engine for efficient file synchronization
type MutagenSyncEngine struct {
	engine *rsync.Engine
}

// NewMutagenSyncEngine creates a new sync engine using Mutagen's rsync implementation.
// Returns a properly initialized MutagenSyncEngine ready for file synchronization operations.
func NewMutagenSyncEngine() *MutagenSyncEngine {
	return &MutagenSyncEngine{
		engine: rsync.NewEngine(),
	}
}

// File efficiently synchronizes a source file to a target location using Mutagen's rsync algorithm.
// It performs intelligent delta synchronization to minimize data transfer.
//
// Parameters:
//   - sourcePath: Path to the source file to synchronize
//   - targetPath: Path where the file should be synchronized to
//   - options: Synchronization options (nil will use defaults)
//
// Returns:
//   - *FileSyncResult: Detailed result of the synchronization operation
//   - error: Any error that occurred during synchronization
//
// The function handles various scenarios:
//   - New files (target doesn't exist)
//   - Identical files (no sync needed)
//   - Different files (delta sync using rsync algorithm)
func (m *MutagenSyncEngine) File(sourcePath, targetPath string, options *Options) (*FileSyncResult, error) {
	// Input validation
	if sourcePath == "" {
		return nil, fmt.Errorf("source path cannot be empty")
	}
	if targetPath == "" {
		return nil, fmt.Errorf("target path cannot be empty")
	}
	if sourcePath == targetPath {
		return nil, fmt.Errorf("source and target paths cannot be identical")
	}

	// Use default options if none provided
	if options == nil {
		options = DefaultOptions()
	}

	// Initialize result structure
	result := &FileSyncResult{
		SourcePath: sourcePath,
		TargetPath: targetPath,
	}

	// Validate source file exists and is accessible
	sourceInfo, err := m.validateSourceFile(sourcePath, result)
	if err != nil {
		return result, fmt.Errorf("source file validation failed: %w", err)
	}

	// Check target file status
	targetExists, targetInfo, err := m.checkTargetFile(targetPath, result)
	if err != nil {
		return result, fmt.Errorf("target file check failed: %w", err)
	}

	// Handle new file case (target doesn't exist)
	if !targetExists {
		return m.handleNewFile(sourcePath, targetPath, options, result, sourceInfo)
	}

	// Check if files are identical (no sync needed)
	if m.filesAreIdentical(sourcePath, targetPath, result, sourceInfo, targetInfo, options) {
		return result, nil
	}

	// Files are different, perform delta synchronization
	return m.syncDifferentFiles(sourcePath, targetPath, options, result, sourceInfo)
}

// validateSourceFile validates that the source file exists and is accessible.
// It performs comprehensive checks to ensure the source file is ready for synchronization.
//
// Parameters:
//   - sourcePath: Path to the source file to validate
//   - result: FileSyncResult to update with any errors
//
// Returns:
//   - os.FileInfo: File information if validation succeeds
//   - error: Validation error if the file cannot be accessed
func (m *MutagenSyncEngine) validateSourceFile(sourcePath string, result *FileSyncResult) (os.FileInfo, error) {
	// Check if source path is a directory
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		if os.IsNotExist(err) {
			result.Error = fmt.Errorf("source file does not exist: %s", sourcePath)
		} else if os.IsPermission(err) {
			result.Error = fmt.Errorf("permission denied accessing source file: %s", sourcePath)
		} else {
			result.Error = fmt.Errorf("failed to access source file %s: %w", sourcePath, err)
		}
		return nil, result.Error
	}

	// Ensure it's a regular file, not a directory
	if sourceInfo.IsDir() {
		result.Error = fmt.Errorf("source path is a directory, not a file: %s", sourcePath)
		return nil, result.Error
	}

	// Check if file is readable
	file, err := os.Open(sourcePath)
	if err != nil {
		result.Error = fmt.Errorf("source file is not readable: %s: %w", sourcePath, err)
		return nil, result.Error
	}
	if closeErr := file.Close(); closeErr != nil {
		// Log the error but don't fail the function since we've already read the data
		// This is a best-effort cleanup
		logger := logging.GetGlobalLogger()
		logger.Warn("failed to close file during cleanup", slog.String("error", closeErr.Error()))
	}

	return sourceInfo, nil
}

// checkTargetFile checks the status of the target file and determines if it exists.
// It handles various error conditions gracefully and provides detailed error information.
//
// Parameters:
//   - targetPath: Path to the target file to check
//   - result: FileSyncResult to update with any errors
//
// Returns:
//   - bool: true if target file exists, false otherwise
//   - os.FileInfo: File information if target exists, nil otherwise
//   - error: Any error that occurred during the check
func (m *MutagenSyncEngine) checkTargetFile(targetPath string, result *FileSyncResult) (bool, os.FileInfo, error) {
	targetInfo, err := os.Stat(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Target doesn't exist - this is normal for new files
			return false, nil, nil
		} else if os.IsPermission(err) {
			result.Error = fmt.Errorf("permission denied accessing target path: %s", targetPath)
			return false, nil, result.Error
		} else {
			result.Error = fmt.Errorf("failed to check target file status: %s: %w", targetPath, err)
			return false, nil, result.Error
		}
	}

	// Target exists - ensure it's a regular file if it's not a directory
	if !targetInfo.IsDir() {
		// Check if target file is writable
		file, err := os.OpenFile(targetPath, os.O_WRONLY, 0)
		if err != nil {
			result.Error = fmt.Errorf("target file is not writable: %s: %w", targetPath, err)
			return true, targetInfo, result.Error
		}
		if closeErr := file.Close(); closeErr != nil {
			// Log the error but don't fail the function since we've already read the data
			// This is a best-effort cleanup
			logger := logging.GetGlobalLogger()
			logger.Warn("failed to close file during cleanup", slog.String("error", closeErr.Error()))
		}
	}

	return true, targetInfo, nil
}

// handleNewFile handles the case where the target file doesn't exist.
// It creates the target directory if needed and copies the source file to the target location.
//
// Parameters:
//   - sourcePath: Path to the source file
//   - targetPath: Path where the file should be copied
//   - options: Synchronization options
//   - result: FileSyncResult to update with operation results
//   - sourceInfo: File information for the source file
//
// Returns:
//   - *FileSyncResult: Updated result with operation details
//   - error: Any error that occurred during the operation
func (m *MutagenSyncEngine) handleNewFile(sourcePath, targetPath string, options *Options, result *FileSyncResult, sourceInfo os.FileInfo) (*FileSyncResult, error) {
	if options.DryRun {
		return m.handleDryRunNewFile(sourcePath, targetPath, options, result, sourceInfo)
	}

	// Create target directory if it doesn't exist
	if err := m.createTargetDirectory(targetPath); err != nil {
		result.Error = fmt.Errorf("failed to create target directory: %w", err)
		return result, result.Error
	}

	// Copy the new file
	if err := m.copyNewFile(sourcePath, targetPath, options, result, sourceInfo); err != nil {
		result.Error = fmt.Errorf("failed to copy new file: %w", err)
		return result, result.Error
	}

	return result, nil
}

// handleDryRunNewFile handles dry run for new files
func (m *MutagenSyncEngine) handleDryRunNewFile(sourcePath, targetPath string, options *Options, result *FileSyncResult, sourceInfo os.FileInfo) (*FileSyncResult, error) {
	if options.Verbose {
		logger := logging.GetGlobalLogger()
		logger.Info("would copy file (new file)",
			slog.String("source", sourcePath),
			slog.String("target", targetPath))
	}
	result.BytesTransferred = sourceInfo.Size()
	result.Operations = 1
	result.Success = true
	return result, nil
}

// createTargetDirectory creates the target directory if it doesn't exist.
// It ensures the parent directory structure is created with appropriate permissions.
//
// Parameters:
//   - targetPath: Path to the target file (directory will be created for its parent)
//
// Returns:
//   - error: Any error that occurred during directory creation
func (m *MutagenSyncEngine) createTargetDirectory(targetPath string) error {
	targetDir := filepath.Dir(targetPath)
	if targetDir == "." || targetDir == "" {
		// No directory creation needed for current directory
		return nil
	}

	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("failed to create target directory %s: %w", targetDir, err)
	}

	return nil
}

// copyNewFile copies a new file to the target location
func (m *MutagenSyncEngine) copyNewFile(sourcePath, targetPath string, options *Options, result *FileSyncResult, sourceInfo os.FileInfo) error {
	if err := copyFile(sourcePath, targetPath, options); err != nil {
		return fmt.Errorf("failed to copy file: %v", err)
	}

	result.BytesTransferred = sourceInfo.Size()
	result.Operations = 1
	result.Success = true
	return nil
}

// filesAreIdentical checks if source and target files are identical
func (m *MutagenSyncEngine) filesAreIdentical(sourcePath, targetPath string, result *FileSyncResult, sourceInfo, targetInfo os.FileInfo, options *Options) bool {
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
			logger := logging.GetGlobalLogger()
			logger.Info("files are identical, no sync needed")
		}
		result.Success = true
		return true
	}

	return false
}

// syncDifferentFiles handles syncing files that are different
func (m *MutagenSyncEngine) syncDifferentFiles(sourcePath, targetPath string, options *Options, result *FileSyncResult, sourceInfo os.FileInfo) (*FileSyncResult, error) {
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
func (m *MutagenSyncEngine) handleDryRunSync(sourcePath, targetPath string, options *Options, result *FileSyncResult, sourceInfo os.FileInfo) (*FileSyncResult, error) {
	if options.Verbose {
		logger := logging.GetGlobalLogger()
		logger.Info("would sync file using Mutagen rsync algorithm",
			slog.String("source", sourcePath),
			slog.String("target", targetPath))
	}
	result.BytesTransferred = sourceInfo.Size()
	result.Operations = 1
	result.Success = true
	return result, nil
}

// createBackupIfNeeded creates a backup of the target file if requested
func (m *MutagenSyncEngine) createBackupIfNeeded(targetPath string, options *Options, result *FileSyncResult) error {
	if !options.CreateBackup {
		return nil
	}

	backupPath := targetPath + ".backup"
	if err := copyFile(targetPath, backupPath, nil); err != nil {
		return fmt.Errorf("failed to create backup: %v", err)
	}

	if options.Verbose {
		logger := logging.GetGlobalLogger()
		logger.Info("created backup", slog.String("path", backupPath))
	}

	return nil
}

// performMutagenSync performs the actual file synchronization using Mutagen's rsync algorithm
func (m *MutagenSyncEngine) performMutagenSync(sourcePath, targetPath string, options *Options, result *FileSyncResult, sourceInfo os.FileInfo) (*FileSyncResult, error) {
	if options.Verbose {
		logger := logging.GetGlobalLogger()
		logger.Info("using Mutagen rsync to sync files")
	}

	// Open files
	sourceFile, targetFile, err := m.openSyncFiles(sourcePath, targetPath, result)
	if err != nil {
		return result, err
	}
	defer func() {
		if closeErr := sourceFile.Close(); closeErr != nil {
			// Log the error but don't fail the function since we've already processed the data
			// This is a best-effort cleanup
			logger := logging.GetGlobalLogger()
			logger.Warn("failed to close source file during cleanup", slog.String("error", closeErr.Error()))
		}
	}()
	defer func() {
		if closeErr := targetFile.Close(); closeErr != nil {
			// Log the error but don't fail the function since we've already processed the data
			// This is a best-effort cleanup
			logger := logging.GetGlobalLogger()
			logger.Warn("failed to close target file during cleanup", slog.String("error", closeErr.Error()))
		}
	}()

	// Generate signature
	signature, err := m.generateSignature(targetFile, options, result)
	if err != nil {
		return result, err
	}

	// Create delta operations
	operations, literalBytes, copyBytes, err := m.createDeltaOperations(sourceFile, signature, result)
	if err != nil {
		return result, err
	}

	// Update result statistics
	m.updateResultStatistics(result, literalBytes, copyBytes, operations, options)

	// Apply delta to create new target file
	if err := m.applyDeltaOperations(targetFile, operations, signature, targetPath, result); err != nil {
		return result, err
	}

	// Preserve file attributes if requested
	if err := m.preserveFileAttributes(sourcePath, targetPath, options, result); err != nil {
		return result, err
	}

	result.Success = true
	return result, nil
}

// openSyncFiles opens source and target files for synchronization
func (m *MutagenSyncEngine) openSyncFiles(sourcePath, targetPath string, result *FileSyncResult) (*os.File, *os.File, error) {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		result.Error = fmt.Errorf("failed to open source file: %v", err)
		return nil, nil, result.Error
	}

	targetFile, err := os.Open(targetPath)
	if err != nil {
		if closeErr := sourceFile.Close(); closeErr != nil {
			// Log the error but don't fail the function since we're already in an error state
			// This is a best-effort cleanup
			logger := logging.GetGlobalLogger()
			logger.Warn("failed to close source file during error cleanup", slog.String("error", closeErr.Error()))
		}
		result.Error = fmt.Errorf("failed to open target file: %v", err)
		return nil, nil, result.Error
	}

	return sourceFile, targetFile, nil
}

// generateSignature generates the signature of the target file
func (m *MutagenSyncEngine) generateSignature(targetFile *os.File, options *Options, result *FileSyncResult) (*rsync.Signature, error) {
	// Convert int32 to uint64 safely
	if options.BlockLength < 0 {
		result.Error = fmt.Errorf("block length cannot be negative")
		return nil, result.Error
	}
	blockLength := uint64(options.BlockLength)

	signature, err := m.engine.Signature(targetFile, blockLength)
	if err != nil {
		result.Error = fmt.Errorf("failed to generate target signature: %v", err)
		return nil, result.Error
	}

	if options.Verbose {
		logger := logging.GetGlobalLogger()
		logger.Info("generated signature", slog.Int("blocks", len(signature.Hashes)))
	}

	return signature, nil
}

// createDeltaOperations creates delta operations and tracks statistics
func (m *MutagenSyncEngine) createDeltaOperations(sourceFile *os.File, signature *rsync.Signature, result *FileSyncResult) ([]*rsync.Operation, int64, int64, error) {
	var operations []*rsync.Operation
	var literalBytes, copyBytes int64

	// Create transmit function for collecting operations and statistics
	transmit := func(op *rsync.Operation) error {
		operations = append(operations, &rsync.Operation{
			Data:  append([]byte(nil), op.Data...), // Copy data
			Start: op.Start,
			Count: op.Count,
		})

		// Track operation statistics
		literalBytes, copyBytes = m.trackOperationStats(op, signature, literalBytes, copyBytes)
		return nil
	}

	// Reset source file for deltification
	if _, err := sourceFile.Seek(0, io.SeekStart); err != nil {
		result.Error = fmt.Errorf("failed to seek source file: %v", err)
		return nil, 0, 0, result.Error
	}

	// Generate delta operations
	if err := m.engine.Deltify(sourceFile, signature, rsync.DefaultMaximumDataOperationSize, transmit); err != nil {
		result.Error = fmt.Errorf("failed to create delta: %v", err)
		return nil, 0, 0, result.Error
	}

	return operations, literalBytes, copyBytes, nil
}

// trackOperationStats tracks statistics for a single operation
func (m *MutagenSyncEngine) trackOperationStats(op *rsync.Operation, signature *rsync.Signature, literalBytes, copyBytes int64) (int64, int64) {
	if len(op.Data) > 0 {
		// Track literal bytes
		literalBytes = m.trackLiteralBytes(op, literalBytes)
	} else {
		// Track copy bytes
		copyBytes = m.trackCopyBytes(op, signature, copyBytes)
		// Adjust for last block if needed
		copyBytes = m.adjustLastBlockSize(op, signature, copyBytes)
	}
	return literalBytes, copyBytes
}

// trackLiteralBytes safely tracks literal bytes with overflow protection
func (m *MutagenSyncEngine) trackLiteralBytes(op *rsync.Operation, literalBytes int64) int64 {
	if dataLen, err := utilities.SafeInt64(len(op.Data)); err == nil {
		if newLiteralBytes, err := utilities.SafeAddInt64(literalBytes, dataLen); err == nil {
			return newLiteralBytes
		}
	}
	// Handle overflow by capping at maximum safe value
	return math.MaxInt64
}

// trackCopyBytes safely tracks copy bytes with overflow protection
func (m *MutagenSyncEngine) trackCopyBytes(op *rsync.Operation, signature *rsync.Signature, copyBytes int64) int64 {
	if op.Count == 0 || signature.BlockSize == 0 {
		return copyBytes
	}

	// Convert uint32 to int64 safely with bounds checking
	if op.Count > math.MaxInt64 {
		return math.MaxInt64
	}
	if signature.BlockSize > math.MaxInt64 {
		return math.MaxInt64
	}
	count64 := int64(op.Count)
	blockSize64 := int64(signature.BlockSize)

	// Use safe multiplication
	multiplied, err := utilities.SafeMultiplyInt64(count64, blockSize64)
	if err != nil {
		return math.MaxInt64
	}

	// Add to existing copy bytes
	if newCopyBytes, err := utilities.SafeAddInt64(copyBytes, multiplied); err == nil {
		return newCopyBytes
	}

	// Handle overflow by capping at maximum safe value
	return math.MaxInt64
}

// adjustLastBlockSize adjusts copy bytes for the last block which might be shorter
func (m *MutagenSyncEngine) adjustLastBlockSize(op *rsync.Operation, signature *rsync.Signature, copyBytes int64) int64 {
	// Convert len(signature.Hashes) to uint64 safely (int fits in uint64)
	hashCount := uint64(len(signature.Hashes))

	if op.Start != hashCount-1 {
		return copyBytes
	}

	if signature.BlockSize < signature.LastBlockSize {
		return copyBytes
	}

	// Safe subtraction since we're dealing with positive values
	// Convert uint32 to int64 safely with bounds checking
	blockDiff := signature.BlockSize - signature.LastBlockSize
	if blockDiff > math.MaxInt64 {
		return copyBytes
	}
	adjustment := int64(blockDiff)

	if newCopyBytes, err := utilities.SafeSubtractInt64(copyBytes, adjustment); err == nil {
		return newCopyBytes
	}

	// Handle underflow by setting to 0
	return 0
}

// updateResultStatistics updates the result with operation statistics
func (m *MutagenSyncEngine) updateResultStatistics(result *FileSyncResult, literalBytes, copyBytes int64, operations []*rsync.Operation, options *Options) {
	result.BytesTransferred = literalBytes
	result.BytesSaved = copyBytes
	result.Operations = len(operations)

	if options.Verbose {
		logger := logging.GetGlobalLogger()
		logger.Info("delta generated",
			slog.Int64("literal_bytes", literalBytes),
			slog.Int64("copy_bytes", copyBytes),
			slog.Int("operations", result.Operations))
	}
}

// applyDeltaOperations applies delta operations to create the new target file
func (m *MutagenSyncEngine) applyDeltaOperations(targetFile *os.File, operations []*rsync.Operation, signature *rsync.Signature, targetPath string, result *FileSyncResult) error {
	// Create temporary file
	tempPath := targetPath + ".tmp"
	outputFile, err := os.Create(tempPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to create temp file: %v", err)
		return result.Error
	}
	defer func() {
		if closeErr := outputFile.Close(); closeErr != nil {
			// Log the error but don't fail the function since we've already written the data
			// This is a best-effort cleanup
			logger := logging.GetGlobalLogger()
			logger.Warn("failed to close output file during cleanup", slog.String("error", closeErr.Error()))
		}
	}()

	// Reset target file for patching
	if _, err := targetFile.Seek(0, io.SeekStart); err != nil {
		result.Error = fmt.Errorf("failed to seek target file: %v", err)
		return result.Error
	}

	// Apply all operations to reconstruct the source file
	for _, op := range operations {
		if err := m.engine.Patch(outputFile, targetFile, signature, op); err != nil {
			result.Error = fmt.Errorf("failed to apply operation: %v", err)
			return result.Error
		}
	}

	if closeErr := outputFile.Close(); closeErr != nil {
		// Log the error but don't fail the function since we've already written the data
		// This is a best-effort cleanup
		logger := logging.GetGlobalLogger()
		logger.Warn("failed to close output file during cleanup", slog.String("error", closeErr.Error()))
	}

	// Replace target file with new version
	if err := os.Rename(tempPath, targetPath); err != nil {
		result.Error = fmt.Errorf("failed to replace target file: %v", err)
		return result.Error
	}

	return nil
}

// preserveFileAttributes preserves file attributes if requested
func (m *MutagenSyncEngine) preserveFileAttributes(sourcePath, targetPath string, options *Options, result *FileSyncResult) error {
	if !options.PreservePerms && !options.PreserveOwner && !options.PreserveGroup {
		return nil
	}

	if err := preserveAttributes(sourcePath, targetPath, options); err != nil {
		result.Error = fmt.Errorf("failed to preserve attributes: %v", err)
		return result.Error
	}

	return nil
}
