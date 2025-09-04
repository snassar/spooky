package sync

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewMutagenSyncEngine(t *testing.T) {
	engine := NewMutagenSyncEngine()
	if engine == nil {
		t.Fatal("NewMutagenSyncEngine returned nil")
	}
	if engine.engine == nil {
		t.Fatal("MutagenSyncEngine.engine is nil")
	}
}

func TestMutagenSyncEngine_SyncFile_InputValidation(t *testing.T) {
	engine := NewMutagenSyncEngine()

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
			result, err := engine.SyncFile(tt.sourcePath, tt.targetPath, nil)
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

func TestMutagenSyncEngine_validateSourceFile(t *testing.T) {
	engine := NewMutagenSyncEngine()

	// Create a temporary file for testing
	tempDir := t.TempDir()
	sourceFile := filepath.Join(tempDir, "source.txt")
	if err := os.WriteFile(sourceFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	tests := []struct {
		name        string
		sourcePath  string
		expectError bool
	}{
		{
			name:        "valid file",
			sourcePath:  sourceFile,
			expectError: false,
		},
		{
			name:        "non-existent file",
			sourcePath:  filepath.Join(tempDir, "nonexistent.txt"),
			expectError: true,
		},
		{
			name:        "directory instead of file",
			sourcePath:  tempDir,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &FileSyncResult{}
			fileInfo, err := engine.validateSourceFile(tt.sourcePath, result)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if fileInfo != nil {
					t.Errorf("expected nil fileInfo on error, got %v", fileInfo)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if fileInfo == nil {
					t.Errorf("expected fileInfo but got nil")
				}
			}
		})
	}
}

func TestMutagenSyncEngine_checkTargetFile(t *testing.T) {
	engine := NewMutagenSyncEngine()

	// Create a temporary directory for testing
	tempDir := t.TempDir()
	targetFile := filepath.Join(tempDir, "target.txt")

	// Create a test file
	if err := os.WriteFile(targetFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	tests := []struct {
		name         string
		targetPath   string
		expectExists bool
		expectError  bool
	}{
		{
			name:         "existing file",
			targetPath:   targetFile,
			expectExists: true,
			expectError:  false,
		},
		{
			name:         "non-existent file",
			targetPath:   filepath.Join(tempDir, "nonexistent.txt"),
			expectExists: false,
			expectError:  false,
		},
		{
			name:         "directory",
			targetPath:   tempDir,
			expectExists: true,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &FileSyncResult{}
			exists, fileInfo, err := engine.checkTargetFile(tt.targetPath, result)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if exists != tt.expectExists {
				t.Errorf("expected exists=%v, got %v", tt.expectExists, exists)
			}

			if tt.expectExists && fileInfo == nil {
				t.Errorf("expected fileInfo but got nil")
			}
		})
	}
}

func TestMutagenSyncEngine_createTargetDirectory(t *testing.T) {
	engine := NewMutagenSyncEngine()

	tempDir := t.TempDir()

	tests := []struct {
		name        string
		targetPath  string
		expectError bool
	}{
		{
			name:        "create nested directory",
			targetPath:  filepath.Join(tempDir, "nested", "deep", "file.txt"),
			expectError: false,
		},
		{
			name:        "current directory",
			targetPath:  "file.txt",
			expectError: false,
		},
		{
			name:        "empty directory",
			targetPath:  "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := engine.createTargetDirectory(tt.targetPath)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
