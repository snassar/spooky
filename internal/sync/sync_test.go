package sync

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSyncFile_InputValidation(t *testing.T) {
	tests := []struct {
		name        string
		sourcePath  string
		targetPath  string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty source path",
			sourcePath:  "",
			targetPath:  "target.txt",
			expectError: true,
			errorMsg:    "source path cannot be empty",
		},
		{
			name:        "empty target path",
			sourcePath:  "source.txt",
			targetPath:  "",
			expectError: true,
			errorMsg:    "target path cannot be empty",
		},
		{
			name:        "identical paths",
			sourcePath:  "same.txt",
			targetPath:  "same.txt",
			expectError: true,
			errorMsg:    "source and target paths cannot be identical",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SyncFile(tt.sourcePath, tt.targetPath, nil)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("expected error message %q, got %q", tt.errorMsg, err.Error())
				}
				if result != nil {
					t.Errorf("expected nil result on error, got %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestSyncDirectory_InputValidation(t *testing.T) {
	tests := []struct {
		name        string
		sourcePath  string
		targetPath  string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty source path",
			sourcePath:  "",
			targetPath:  "target",
			expectError: true,
			errorMsg:    "source path cannot be empty",
		},
		{
			name:        "empty target path",
			sourcePath:  "source",
			targetPath:  "",
			expectError: true,
			errorMsg:    "target path cannot be empty",
		},
		{
			name:        "identical paths",
			sourcePath:  "same",
			targetPath:  "same",
			expectError: true,
			errorMsg:    "source and target paths cannot be identical",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SyncDirectory(tt.sourcePath, tt.targetPath, nil)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("expected error message %q, got %q", tt.errorMsg, err.Error())
				}
				if result != nil {
					t.Errorf("expected nil result on error, got %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestSyncDirectory_NonExistentSource(t *testing.T) {
	tempDir := t.TempDir()
	nonExistentDir := filepath.Join(tempDir, "nonexistent")
	targetDir := filepath.Join(tempDir, "target")

	result, err := SyncDirectory(nonExistentDir, targetDir, nil)
	if err == nil {
		t.Errorf("expected error for non-existent source directory")
	}
	if result.Error == nil {
		t.Errorf("expected result.Error to be set")
	}
}

func TestSyncDirectory_FileAsSource(t *testing.T) {
	tempDir := t.TempDir()
	sourceFile := filepath.Join(tempDir, "source.txt")
	targetDir := filepath.Join(tempDir, "target")

	// Create a file instead of directory
	if err := os.WriteFile(sourceFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	result, err := SyncDirectory(sourceFile, targetDir, nil)
	if err == nil {
		t.Errorf("expected error for file as source directory")
	}
	if result.Error == nil {
		t.Errorf("expected result.Error to be set")
	}
}

func TestDefaultSyncOptions(t *testing.T) {
	options := DefaultSyncOptions()
	if options == nil {
		t.Fatal("DefaultSyncOptions returned nil")
	}

	// Check default values
	if options.BlockLength != DefaultBlockLength {
		t.Errorf("expected BlockLength %d, got %d", DefaultBlockLength, options.BlockLength)
	}
	if !options.CreateBackup {
		t.Error("expected CreateBackup to be true")
	}
	if !options.PreservePerms {
		t.Error("expected PreservePerms to be true")
	}
	if options.PreserveOwner {
		t.Error("expected PreserveOwner to be false")
	}
	if options.PreserveGroup {
		t.Error("expected PreserveGroup to be false")
	}
	if options.DryRun {
		t.Error("expected DryRun to be false")
	}
	if options.Verbose {
		t.Error("expected Verbose to be false")
	}
	if options.SyncMode != SyncModeOneWayReplica {
		t.Errorf("expected SyncMode %s, got %s", SyncModeOneWayReplica, options.SyncMode)
	}
}

func TestSyncModes(t *testing.T) {
	// Test that all sync modes are properly defined
	expectedModes := []SyncMode{
		SyncModeOneWayReplica,
		SyncModeOneWaySafe,
		SyncModeTwoWaySafe,
		SyncModeTwoWayResolved,
	}

	for _, mode := range expectedModes {
		if mode == "" {
			t.Errorf("sync mode should not be empty")
		}
	}
}
