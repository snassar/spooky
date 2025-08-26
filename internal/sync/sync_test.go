package sync

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSyncDirectory(t *testing.T) {
	// Create temporary directories for testing
	sourceDir, err := os.MkdirTemp("", "sync_source")
	if err != nil {
		t.Fatalf("Failed to create source temp dir: %v", err)
	}
	defer os.RemoveAll(sourceDir)

	targetDir, err := os.MkdirTemp("", "sync_target")
	if err != nil {
		t.Fatalf("Failed to create target temp dir: %v", err)
	}
	defer os.RemoveAll(targetDir)

	// Create test files in source directory
	testFiles := map[string]string{
		"file1.txt":        "content1",
		"file2.txt":        "content2",
		"subdir/file3.txt": "content3",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(sourceDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
			t.Fatalf("Failed to create file %s: %v", fullPath, err)
		}
	}

	// Test one-way replica sync
	t.Run("OneWayReplica", func(t *testing.T) {
		options := DefaultSyncOptions()
		options.SyncMode = SyncModeOneWayReplica
		options.Verbose = true

		result, err := SyncDirectory(sourceDir, targetDir, options)
		if err != nil {
			t.Fatalf("SyncDirectory failed: %v", err)
		}

		if !result.Success {
			t.Fatalf("Sync failed: %v", result.Error)
		}

		// Verify files were synced
		for path, content := range testFiles {
			targetPath := filepath.Join(targetDir, path)
			data, err := os.ReadFile(targetPath)
			if err != nil {
				t.Errorf("Failed to read synced file %s: %v", targetPath, err)
			} else if string(data) != content {
				t.Errorf("File content mismatch for %s: expected %q, got %q", path, content, string(data))
			}
		}
	})

	// Test one-way safe sync
	t.Run("OneWaySafe", func(t *testing.T) {
		// Create a different target directory
		targetDir2, err := os.MkdirTemp("", "sync_target2")
		if err != nil {
			t.Fatalf("Failed to create target temp dir: %v", err)
		}
		defer os.RemoveAll(targetDir2)

		// Create a newer file in target that should be preserved
		conflictFile := filepath.Join(targetDir2, "file1.txt")
		if err := os.WriteFile(conflictFile, []byte("newer content"), 0o644); err != nil {
			t.Fatalf("Failed to create conflict file: %v", err)
		}

		options := DefaultSyncOptions()
		options.SyncMode = SyncModeOneWaySafe
		options.Verbose = true

		result, err := SyncDirectory(sourceDir, targetDir2, options)
		if err != nil {
			t.Fatalf("SyncDirectory failed: %v", err)
		}

		if !result.Success {
			t.Fatalf("Sync failed: %v", result.Error)
		}

		// Check that conflict was detected
		if len(result.Conflicts) == 0 {
			t.Error("Expected conflict to be detected")
		}

		// Verify conflict file was preserved
		data, err := os.ReadFile(conflictFile)
		if err != nil {
			t.Errorf("Failed to read conflict file: %v", err)
		} else if string(data) != "newer content" {
			t.Errorf("Conflict file was modified: expected 'newer content', got %q", string(data))
		}
	})
}

func TestSyncDirectoryValidation(t *testing.T) {
	// Test with non-existent source
	_, err := SyncDirectory("/non/existent/path", "/tmp/target", DefaultSyncOptions())
	if err == nil {
		t.Error("Expected error for non-existent source directory")
	}

	// Test with file as source
	tempFile, err := os.CreateTemp("", "test_file")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = SyncDirectory(tempFile.Name(), "/tmp/target", DefaultSyncOptions())
	if err == nil {
		t.Error("Expected error for file as source directory")
	}
}
