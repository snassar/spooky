package sync

import (
	"os"
	"path/filepath"
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
	// TODO: Implement directory synchronization with mode support
	// This would handle the broader use cases like:
	// - Database replication
	// - Configuration management
	// - Backup systems
	// - Disaster recovery
	return nil, nil
}

// copyFile performs a simple file copy
func copyFile(src, dst string, options *SyncOptions) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	// Create target directory if it doesn't exist
	targetDir := filepath.Dir(dst)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = destination.ReadFrom(source)
	return err
}

// preserveAttributes preserves file attributes from source to target
func preserveAttributes(sourcePath, targetPath string, options *SyncOptions) error {
	// This is a simplified implementation
	// In a full implementation, you would use syscall to preserve
	// permissions, owner, group, timestamps, etc.
	return nil
}
