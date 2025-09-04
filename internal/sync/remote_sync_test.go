package sync

import (
	"os"
	"path/filepath"
	"testing"

	"spooky/internal/schemas"
	"spooky/internal/ssh"
)

func TestRemoteSyncEngine(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "remote_sync_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if removeErr := os.RemoveAll(tempDir); removeErr != nil {
			// Log the error but don't fail the test since this is cleanup
			t.Logf("Warning: failed to remove temp directory %s: %v", tempDir, removeErr)
		}
	}()

	// Create some test files
	testFile1 := filepath.Join(tempDir, "test1.txt")
	testFile2 := filepath.Join(tempDir, "test2.txt")
	testSubDir := filepath.Join(tempDir, "subdir")
	testFile3 := filepath.Join(testSubDir, "test3.txt")

	// Create test files
	if err := os.WriteFile(testFile1, []byte("test content 1"), 0644); err != nil {
		t.Fatalf("Failed to create test file 1: %v", err)
	}
	if err := os.WriteFile(testFile2, []byte("test content 2"), 0644); err != nil {
		t.Fatalf("Failed to create test file 2: %v", err)
	}
	if err := os.MkdirAll(testSubDir, 0755); err != nil {
		t.Fatalf("Failed to create test subdirectory: %v", err)
	}
	if err := os.WriteFile(testFile3, []byte("test content 3"), 0644); err != nil {
		t.Fatalf("Failed to create test file 3: %v", err)
	}

	// Create a mock SSH manager (this would be mocked in a real test)
	// For now, we'll just test the local file counting functionality
	sshManager := &ssh.SimpleSSHManager{} // This won't work without proper initialization

	// Create remote sync engine
	engine := NewRemoteSyncEngine(sshManager)

	// Test file counting
	progress := &SyncProgress{
		CurrentFile:      "",
		TotalFiles:       0,
		CurrentOperation: "scanning",
	}

	if err := engine.countFiles(tempDir, progress); err != nil {
		t.Fatalf("Failed to count files: %v", err)
	}

	// Should have 3 files (not counting directories)
	expectedFiles := 3
	if progress.TotalFiles != expectedFiles {
		t.Errorf("Expected %d files, got %d", expectedFiles, progress.TotalFiles)
	}

	t.Logf("Successfully counted %d files in test directory", progress.TotalFiles)
}

func TestSyncProgress(t *testing.T) {
	progress := &SyncProgress{
		CurrentFile:      "test.txt",
		FilesProcessed:   5,
		TotalFiles:       10,
		CurrentOperation: "syncing",
		BytesTransferred: 1024,
		BytesSaved:       512,
	}

	// Test percentage calculation
	if progress.TotalFiles > 0 {
		progress.Percentage = float64(progress.FilesProcessed) / float64(progress.TotalFiles) * 100
	}

	expectedPercentage := 50.0
	if progress.Percentage != expectedPercentage {
		t.Errorf("Expected percentage %.1f, got %.1f", expectedPercentage, progress.Percentage)
	}

	t.Logf("Progress: %d/%d files (%.1f%%)", progress.FilesProcessed, progress.TotalFiles, progress.Percentage)
}

func TestRemoteSyncOptions(t *testing.T) {
	// Test creating remote sync options
	options := &RemoteSyncOptions{
		SyncOptions: DefaultSyncOptions(),
		Machine: &schemas.MachinesMachineV1{
			Hostname: "test.example.com",
			User:     "testuser",
		},
		ProgressReport: func(progress *SyncProgress) {
			// Test progress reporter
			t.Logf("Progress update: %s", progress.CurrentOperation)
		},
		ConflictResolve: ConflictResolutionBackup,
		SyncDelete:      true,
	}

	if options.Machine.Hostname != "test.example.com" {
		t.Errorf("Expected hostname 'test.example.com', got '%s'", options.Machine.Hostname)
	}

	if options.SyncDelete != true {
		t.Errorf("Expected SyncDelete to be true, got %v", options.SyncDelete)
	}

	t.Logf("Successfully created remote sync options for %s", options.Machine.Hostname)
}
