package sync

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"spooky/internal/logging"
	"syscall"

	"github.com/pkg/errors"
)

// SyncMode defines the synchronization directionality and conflict resolution
type SyncMode string

const (
	// SyncModeOneWayReplica: Exact one-way replication (alpha → beta)
	// Beta becomes exact copy of alpha, overwriting any local changes
	SyncModeOneWayReplica SyncMode = "one-way-replica"

	// SyncModeOneWaySafe: Safe one-way sync (alpha → beta)
	// Changes only propagate from alpha to beta, conflicts are preserved
	SyncModeOneWaySafe SyncMode = "one-way-safe"

	// SyncModeTwoWaySafe: Bidirectional sync with conflict detection
	// Both endpoints can modify, conflicts are detected and preserved
	SyncModeTwoWaySafe SyncMode = "two-way-safe"

	// SyncModeTwoWayResolved: Bidirectional sync with alpha winning conflicts
	// Both endpoints can modify, alpha always wins conflicts
	SyncModeTwoWayResolved SyncMode = "two-way-resolved"
)

// SyncOptions configures file synchronization behavior
type SyncOptions struct {
	BlockLength   int32    // Block size for rsync algorithm
	CreateBackup  bool     // Create backup before overwriting
	PreservePerms bool     // Preserve file permissions
	PreserveOwner bool     // Preserve file owner
	PreserveGroup bool     // Preserve file group
	DryRun        bool     // Show what would be done without doing it
	Verbose       bool     // Show detailed progress
	SyncMode      SyncMode // Synchronization mode and directionality
}

// DefaultSyncOptions returns default synchronization options
func DefaultSyncOptions() *SyncOptions {
	return &SyncOptions{
		BlockLength:   DefaultBlockLength,
		CreateBackup:  true,
		PreservePerms: true,
		PreserveOwner: false,
		PreserveGroup: false,
		DryRun:        false,
		Verbose:       false,
		SyncMode:      SyncModeOneWayReplica, // Default to exact replication for deployments
	}
}

// FileSyncResult contains the result of a file synchronization operation
type FileSyncResult struct {
	SourcePath       string
	TargetPath       string
	BytesTransferred int64
	BytesSaved       int64
	Operations       int
	Success          bool
	Error            error
	Conflicts        []string // List of conflicts detected (for bidirectional modes)
}

// SyncFile efficiently synchronizes a source file to a target location
// using Mutagen's proven rsync algorithm to minimize data transfer
func SyncFile(sourcePath, targetPath string, options *SyncOptions) (*FileSyncResult, error) {
	// Use Mutagen's rsync implementation directly
	engine := NewMutagenSyncEngine()
	return engine.SyncFile(sourcePath, targetPath, options)
}

// SyncDirectory synchronizes entire directories with support for different modes
func SyncDirectory(sourcePath, targetPath string, options *SyncOptions) (*FileSyncResult, error) {
	if options == nil {
		options = DefaultSyncOptions()
	}

	result := &FileSyncResult{
		SourcePath: sourcePath,
		TargetPath: targetPath,
	}

	// Validate source directory exists
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		result.Error = fmt.Errorf("source directory not found: %v", err)
		return result, result.Error
	}
	if !sourceInfo.IsDir() {
		result.Error = fmt.Errorf("source path is not a directory: %s", sourcePath)
		return result, result.Error
	}

	// Create target directory if it doesn't exist
	if err := os.MkdirAll(targetPath, 0o755); err != nil {
		result.Error = fmt.Errorf("failed to create target directory: %v", err)
		return result, result.Error
	}

	// Check target directory exists
	targetInfo, err := os.Stat(targetPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to stat target directory: %v", err)
		return result, result.Error
	}
	if !targetInfo.IsDir() {
		result.Error = fmt.Errorf("target path is not a directory: %s", targetPath)
		return result, result.Error
	}

	// Perform directory synchronization based on mode
	switch options.SyncMode {
	case SyncModeOneWayReplica:
		return syncOneWayReplica(sourcePath, targetPath, options, result)
	case SyncModeOneWaySafe:
		return syncOneWaySafe(sourcePath, targetPath, options, result)
	case SyncModeTwoWaySafe:
		return syncTwoWaySafe(sourcePath, targetPath, options, result)
	case SyncModeTwoWayResolved:
		return syncTwoWayResolved(sourcePath, targetPath, options, result)
	default:
		result.Error = fmt.Errorf("unknown sync mode: %s", options.SyncMode)
		return result, result.Error
	}
}

// syncOneWayReplica performs exact one-way replication (source → target)
// Target becomes exact copy of source, overwriting any local changes
func syncOneWayReplica(sourcePath, targetPath string, options *SyncOptions, result *FileSyncResult) (*FileSyncResult, error) {
	if options.Verbose {
		logger := logging.GetGlobalLogger()
		logger.Info("performing one-way replica sync",
			slog.String("source", sourcePath),
			slog.String("target", targetPath))
	}

	// Walk source directory and sync all files
	err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "failed to walk source path %s", path)
		}

		// Calculate relative path from source
		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return fmt.Errorf("failed to calculate relative path: %v", err)
		}

		// Skip root directory
		if relPath == "." {
			return nil
		}

		targetFile := filepath.Join(targetPath, relPath)

		if info.IsDir() {
			// Create directory in target
			if err := os.MkdirAll(targetFile, info.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", targetFile, err)
			}
			if options.Verbose {
				logger := logging.GetGlobalLogger()
				logger.Info("created directory", slog.String("path", targetFile))
			}
		} else {
			// Sync file using existing SyncFile function
			fileResult, err := SyncFile(path, targetFile, options)
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
		return result, result.Error
	}

	// Remove files in target that don't exist in source (cleanup)
	err = filepath.Walk(targetPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "failed to walk target path %s", path)
		}

		// Skip root directory
		if path == targetPath {
			return nil
		}

		relPath, err := filepath.Rel(targetPath, path)
		if err != nil {
			return fmt.Errorf("failed to calculate relative path: %v", err)
		}

		sourceFile := filepath.Join(sourcePath, relPath)
		if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
			if options.DryRun {
				if options.Verbose {
					logger := logging.GetGlobalLogger()
					logger.Info("would remove file", slog.String("path", path))
				}
			} else {
				if info.IsDir() {
					if err := os.RemoveAll(path); err != nil {
						return fmt.Errorf("failed to remove directory %s: %v", path, err)
					}
				} else {
					if err := os.Remove(path); err != nil {
						return fmt.Errorf("failed to remove file %s: %v", path, err)
					}
				}
				if options.Verbose {
					logger := logging.GetGlobalLogger()
					logger.Info("removed file", slog.String("path", path))
				}
			}
		}

		return nil
	})

	if err != nil {
		result.Error = errors.Wrap(err, "failed to cleanup target directory")
		return result, result.Error
	}

	result.Success = true
	return result, nil
}

// syncOneWaySafe performs safe one-way sync (source → target)
// Changes only propagate from source to target, conflicts are preserved
func syncOneWaySafe(sourcePath, targetPath string, options *SyncOptions, result *FileSyncResult) (*FileSyncResult, error) {
	if options.Verbose {
		logger := logging.GetGlobalLogger()
		logger.Info("performing one-way safe sync",
			slog.String("source", sourcePath),
			slog.String("target", targetPath))
	}

	// Walk source directory and sync files, but preserve conflicts
	err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "failed to walk source path %s", path)
		}

		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return fmt.Errorf("failed to calculate relative path: %v", err)
		}

		if relPath == "." {
			return nil
		}

		targetFile := filepath.Join(targetPath, relPath)

		if info.IsDir() {
			if err := os.MkdirAll(targetFile, info.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", targetFile, err)
			}
		} else {
			// Check if target file exists and is different
			if targetInfo, err := os.Stat(targetFile); err == nil {
				// Target exists, check for conflicts
				if targetInfo.ModTime().After(info.ModTime()) {
					// Target is newer, preserve it (conflict)
					if options.Verbose {
						logger := logging.GetGlobalLogger()
						logger.Info("conflict detected, preserving target", slog.String("path", targetFile))
					}
					result.Conflicts = append(result.Conflicts, targetFile)
					return nil
				}
			}

			// No conflict, sync the file
			fileResult, err := SyncFile(path, targetFile, options)
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
		result.Error = err
		return result, err
	}

	result.Success = true
	return result, nil
}

// combineSyncResults combines results from two sync operations
func combineSyncResults(source, target *FileSyncResult) *FileSyncResult {
	return &FileSyncResult{
		BytesTransferred: source.BytesTransferred + target.BytesTransferred,
		BytesSaved:       source.BytesSaved + target.BytesSaved,
		Operations:       source.Operations + target.Operations,
		Conflicts:        append(source.Conflicts, target.Conflicts...),
		Success:          true,
	}
}

// applyCombinedResults applies combined results to the target result struct
func applyCombinedResults(result *FileSyncResult, combined *FileSyncResult) {
	result.BytesTransferred = combined.BytesTransferred
	result.BytesSaved = combined.BytesSaved
	result.Operations = combined.Operations
	result.Conflicts = combined.Conflicts
	result.Success = combined.Success
}

// syncTwoWaySafe performs bidirectional sync with conflict detection
// Both endpoints can modify, conflicts are detected and preserved
func syncTwoWaySafe(sourcePath, targetPath string, options *SyncOptions, result *FileSyncResult) (*FileSyncResult, error) {
	if options.Verbose {
		logger := logging.GetGlobalLogger()
		logger.Info("performing two-way safe sync",
			slog.String("source", sourcePath),
			slog.String("target", targetPath))
	}

	// First, sync source → target (one-way safe)
	sourceToTarget, err := syncOneWaySafe(sourcePath, targetPath, options, &FileSyncResult{
		SourcePath: sourcePath,
		TargetPath: targetPath,
	})
	if err != nil {
		result.Error = err
		return result, err
	}

	// Then, sync target → source (one-way safe)
	targetToSource, err := syncOneWaySafe(targetPath, sourcePath, options, &FileSyncResult{
		SourcePath: targetPath,
		TargetPath: sourcePath,
	})
	if err != nil {
		result.Error = err
		return result, err
	}

	// Combine results using helper function
	combined := combineSyncResults(sourceToTarget, targetToSource)
	applyCombinedResults(result, combined)

	return result, nil
}

// syncTwoWayResolved performs bidirectional sync with source winning conflicts
// Both endpoints can modify, source always wins conflicts
func syncTwoWayResolved(sourcePath, targetPath string, options *SyncOptions, result *FileSyncResult) (*FileSyncResult, error) {
	if options.Verbose {
		logger := logging.GetGlobalLogger()
		logger.Info("performing two-way resolved sync (source wins)",
			slog.String("source", sourcePath),
			slog.String("target", targetPath))
	}

	// Use one-way replica for source → target (source always wins)
	sourceToTarget, err := syncOneWayReplica(sourcePath, targetPath, options, &FileSyncResult{
		SourcePath: sourcePath,
		TargetPath: targetPath,
	})
	if err != nil {
		result.Error = err
		return result, err
	}

	// Use one-way safe for target → source (preserve source changes)
	targetToSource, err := syncOneWaySafe(targetPath, sourcePath, options, &FileSyncResult{
		SourcePath: targetPath,
		TargetPath: sourcePath,
	})
	if err != nil {
		result.Error = err
		return result, err
	}

	// Combine results using helper function
	combined := combineSyncResults(sourceToTarget, targetToSource)
	applyCombinedResults(result, combined)

	return result, nil
}

// copyFile performs a simple file copy
func copyFile(src, dst string, options *SyncOptions) error {
	source, err := os.Open(src)
	if err != nil {
		return errors.Wrapf(err, "failed to open source file %s", src)
	}
	defer source.Close()

	// Create target directory if it doesn't exist
	targetDir := filepath.Dir(dst)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return errors.Wrapf(err, "failed to create target directory %s", targetDir)
	}

	destination, err := os.Create(dst)
	if err != nil {
		return errors.Wrapf(err, "failed to create destination file %s", dst)
	}
	defer destination.Close()

	_, err = destination.ReadFrom(source)
	return errors.Wrapf(err, "failed to copy data from %s to %s", src, dst)
}

// preserveAttributes preserves file attributes from source to target
func preserveAttributes(sourcePath, targetPath string, options *SyncOptions) error {
	if !options.PreservePerms && !options.PreserveOwner && !options.PreserveGroup {
		return nil // No attributes to preserve
	}

	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return errors.Wrapf(err, "failed to stat source file %s", sourcePath)
	}

	if options.PreservePerms {
		if err := os.Chmod(targetPath, sourceInfo.Mode()); err != nil {
			return errors.Wrapf(err, "failed to preserve permissions for %s", targetPath)
		}
	}

	if options.PreserveOwner || options.PreserveGroup {
		stat, ok := sourceInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return errors.New("failed to get system-specific file info")
		}

		if options.PreserveOwner && options.PreserveGroup {
			if err := os.Chown(targetPath, int(stat.Uid), int(stat.Gid)); err != nil {
				return errors.Wrapf(err, "failed to preserve owner and group for %s", targetPath)
			}
		} else if options.PreserveOwner {
			// Get current group of target
			targetInfo, err := os.Stat(targetPath)
			if err != nil {
				return errors.Wrapf(err, "failed to stat target file %s", targetPath)
			}
			targetStat, ok := targetInfo.Sys().(*syscall.Stat_t)
			if !ok {
				return errors.New("failed to get target system-specific file info")
			}
			if err := os.Chown(targetPath, int(stat.Uid), int(targetStat.Gid)); err != nil {
				return errors.Wrapf(err, "failed to preserve owner for %s", targetPath)
			}
		} else if options.PreserveGroup {
			// Get current owner of target
			targetInfo, err := os.Stat(targetPath)
			if err != nil {
				return errors.Wrapf(err, "failed to stat target file %s", targetPath)
			}
			targetStat, ok := targetInfo.Sys().(*syscall.Stat_t)
			if !ok {
				return errors.New("failed to get target system-specific file info")
			}
			if err := os.Chown(targetPath, int(targetStat.Uid), int(stat.Gid)); err != nil {
				return errors.Wrapf(err, "failed to preserve group for %s", targetPath)
			}
		}
	}

	return nil
}
