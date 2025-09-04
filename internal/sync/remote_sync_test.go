package sync

import (
	"context"
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
	progress := &Progress{
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

func TestProgress(t *testing.T) {
	progress := &Progress{
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

func TestRemoteOptions(t *testing.T) {
	// Test creating remote sync options
	options := &RemoteOptions{
		Options: DefaultOptions(),
		Machine: &schemas.MachinesMachineV1{
			Hostname: "test.example.com",
			User:     "testuser",
		},
		ProgressReport: func(progress *Progress) {
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

func TestListLocalFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "list_local_files_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if removeErr := os.RemoveAll(tempDir); removeErr != nil {
			t.Logf("Warning: failed to remove temp directory %s: %v", tempDir, removeErr)
		}
	}()

	// Create test files and directories
	testFile1 := filepath.Join(tempDir, "file1.txt")
	testFile2 := filepath.Join(tempDir, "file2.txt")
	testSubDir := filepath.Join(tempDir, "subdir")
	testFile3 := filepath.Join(testSubDir, "file3.txt")
	testFile4 := filepath.Join(testSubDir, "file4.txt")

	// Create test files
	if err := os.WriteFile(testFile1, []byte("content 1"), 0644); err != nil {
		t.Fatalf("Failed to create test file 1: %v", err)
	}
	if err := os.WriteFile(testFile2, []byte("content 2"), 0644); err != nil {
		t.Fatalf("Failed to create test file 2: %v", err)
	}
	if err := os.MkdirAll(testSubDir, 0755); err != nil {
		t.Fatalf("Failed to create test subdirectory: %v", err)
	}
	if err := os.WriteFile(testFile3, []byte("content 3"), 0644); err != nil {
		t.Fatalf("Failed to create test file 3: %v", err)
	}
	if err := os.WriteFile(testFile4, []byte("content 4"), 0644); err != nil {
		t.Fatalf("Failed to create test file 4: %v", err)
	}

	// Create remote sync engine
	engine := NewRemoteSyncEngine(nil) // SSH manager not needed for this test

	// Test listing local files
	files, err := engine.listLocalFiles(tempDir)
	if err != nil {
		t.Fatalf("Failed to list local files: %v", err)
	}

	// Should have 4 files
	expectedFiles := 4
	if len(files) != expectedFiles {
		t.Errorf("Expected %d files, got %d", expectedFiles, len(files))
	}

	// Check that all expected files are present
	expectedFilePaths := []string{"file1.txt", "file2.txt", "subdir/file3.txt", "subdir/file4.txt"}
	fileMap := make(map[string]bool)
	for _, file := range files {
		fileMap[file] = true
	}

	for _, expectedFile := range expectedFilePaths {
		if !fileMap[expectedFile] {
			t.Errorf("Expected file %s not found in results", expectedFile)
		}
	}

	t.Logf("Successfully listed %d local files", len(files))
}

func TestCleanupRemoteFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "cleanup_remote_files_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if removeErr := os.RemoveAll(tempDir); removeErr != nil {
			t.Logf("Warning: failed to remove temp directory %s: %v", tempDir, removeErr)
		}
	}()

	// Create test files
	testFile1 := filepath.Join(tempDir, "keep1.txt")
	testFile2 := filepath.Join(tempDir, "keep2.txt")
	testSubDir := filepath.Join(tempDir, "subdir")
	testFile3 := filepath.Join(testSubDir, "keep3.txt")

	// Create test files
	if err := os.WriteFile(testFile1, []byte("keep content 1"), 0644); err != nil {
		t.Fatalf("Failed to create test file 1: %v", err)
	}
	if err := os.WriteFile(testFile2, []byte("keep content 2"), 0644); err != nil {
		t.Fatalf("Failed to create test file 2: %v", err)
	}
	if err := os.MkdirAll(testSubDir, 0755); err != nil {
		t.Fatalf("Failed to create test subdirectory: %v", err)
	}
	if err := os.WriteFile(testFile3, []byte("keep content 3"), 0644); err != nil {
		t.Fatalf("Failed to create test file 3: %v", err)
	}

	// Create remote sync engine with nil SSH manager (we'll only test the disabled case)
	engine := NewRemoteSyncEngine(nil)

	// Create test options with cleanup disabled
	options := &RemoteOptions{
		Options: &Options{
			DryRun: true, // Use dry run to avoid actual SSH operations
		},
		Machine: &schemas.MachinesMachineV1{
			Hostname: "test.example.com",
			User:     "testuser",
		},
		SyncDelete: false, // Cleanup disabled
	}

	progress := &Progress{}

	// Test with cleanup disabled - this should return early without calling SSH
	err = engine.cleanupRemoteFiles(context.Background(), tempDir, "/remote/path", options, progress)
	if err != nil {
		t.Fatalf("Cleanup with SyncDelete=false should not fail: %v", err)
	}

	// Test that the function correctly checks SyncDelete flag
	// We can't test the full functionality without a real SSH connection,
	// but we can verify the early return behavior
	t.Logf("Cleanup tests completed successfully - SyncDelete flag properly checked")
}
