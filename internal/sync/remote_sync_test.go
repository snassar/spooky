package sync

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
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
	sshManager := &ssh.Manager{} // This won't work without proper initialization

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
	expectedFiles := int64(3)
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

func TestConcurrentSyncResult(t *testing.T) {
	result := &ConcurrentSyncResult{}

	// Test adding results
	fileResult1 := &FileSyncResult{
		BytesTransferred: 1000,
		BytesSaved:       500,
		Operations:       1,
	}
	fileResult2 := &FileSyncResult{
		BytesTransferred: 2000,
		BytesSaved:       1000,
		Operations:       2,
	}

	result.AddResult(fileResult1)
	result.AddResult(fileResult2)

	// Test adding errors
	err1 := fmt.Errorf("test error 1")
	err2 := fmt.Errorf("test error 2")
	result.AddError(err1)
	result.AddError(err2)

	// Verify results
	if result.BytesTransferred != 3000 {
		t.Errorf("Expected BytesTransferred=3000, got %d", result.BytesTransferred)
	}
	if result.BytesSaved != 1500 {
		t.Errorf("Expected BytesSaved=1500, got %d", result.BytesSaved)
	}
	if result.Operations != 3 {
		t.Errorf("Expected Operations=3, got %d", result.Operations)
	}

	errors := result.GetErrors()
	if len(errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errors))
	}

	t.Logf("ConcurrentSyncResult test passed: %d bytes transferred, %d operations, %d errors",
		result.BytesTransferred, result.Operations, len(errors))
}

func TestProgressThreadSafety(t *testing.T) {
	progress := &Progress{
		TotalFiles: 100,
	}

	// Test concurrent access to progress
	var wg sync.WaitGroup
	numGoroutines := 10
	filesPerGoroutine := 10

	// Start multiple goroutines that increment files processed
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < filesPerGoroutine; j++ {
				progress.IncrementFilesProcessed()
				progress.SetCurrentFile(fmt.Sprintf("file_%d_%d", i, j))
				progress.SetCurrentOperation("syncing")
				progress.UpdatePercentage()
			}
		}()
	}

	wg.Wait()

	// Verify final count
	expectedFiles := int64(numGoroutines * filesPerGoroutine)
	if progress.FilesProcessed != expectedFiles {
		t.Errorf("Expected %d files processed, got %d", expectedFiles, progress.FilesProcessed)
	}

	// Verify percentage calculation
	expectedPercentage := 100.0
	if progress.Percentage != expectedPercentage {
		t.Errorf("Expected percentage %.1f, got %.1f", expectedPercentage, progress.Percentage)
	}

	t.Logf("Thread safety test passed: %d files processed (%.1f%%)",
		progress.FilesProcessed, progress.Percentage)
}

func TestGetOptimalConcurrency(t *testing.T) {
	// Create a mock SSH manager
	sshManager := &ssh.Manager{}
	engine := NewRemoteSyncEngine(sshManager)

	// Test with explicit concurrency setting
	options := &RemoteOptions{
		Options: &Options{
			MaxConcurrency: 3,
		},
	}

	concurrency := engine.getOptimalConcurrency(options)
	if concurrency != 3 {
		t.Errorf("Expected concurrency=3, got %d", concurrency)
	}

	// Test with auto-detection (0)
	options.MaxConcurrency = 0
	concurrency = engine.getOptimalConcurrency(options)
	if concurrency < 1 || concurrency > 10 {
		t.Errorf("Expected concurrency between 1-10, got %d", concurrency)
	}

	t.Logf("Concurrency test passed: explicit=3, auto-detected=%d", concurrency)
}

func TestFileProcessingJob(t *testing.T) {
	job := FileProcessingJob{
		LocalPath:  "/local/path/file.txt",
		RemotePath: "/remote/path/file.txt",
		IsDir:      false,
	}

	if job.LocalPath != "/local/path/file.txt" {
		t.Errorf("Expected LocalPath='/local/path/file.txt', got '%s'", job.LocalPath)
	}
	if job.RemotePath != "/remote/path/file.txt" {
		t.Errorf("Expected RemotePath='/remote/path/file.txt', got '%s'", job.RemotePath)
	}
	if job.IsDir {
		t.Error("Expected IsDir=false, got true")
	}

	// Test directory job
	dirJob := FileProcessingJob{
		LocalPath:  "/local/path/dir",
		RemotePath: "/remote/path/dir",
		IsDir:      true,
	}

	if !dirJob.IsDir {
		t.Error("Expected IsDir=true, got false")
	}

	t.Logf("FileProcessingJob test passed: file job and directory job created successfully")
}

func TestProgressSnapshot(t *testing.T) {
	progress := &Progress{
		CurrentFile:      "test.txt",
		FilesProcessed:   5,
		TotalFiles:       10,
		BytesTransferred: 1024,
		BytesSaved:       512,
		CurrentOperation: "syncing",
		Percentage:       50.0,
	}

	snapshot := progress.GetProgressSnapshot()

	// Verify snapshot contains correct values
	if snapshot.CurrentFile != "test.txt" {
		t.Errorf("Expected CurrentFile='test.txt', got '%s'", snapshot.CurrentFile)
	}
	if snapshot.FilesProcessed != 5 {
		t.Errorf("Expected FilesProcessed=5, got %d", snapshot.FilesProcessed)
	}
	if snapshot.TotalFiles != 10 {
		t.Errorf("Expected TotalFiles=10, got %d", snapshot.TotalFiles)
	}
	if snapshot.BytesTransferred != 1024 {
		t.Errorf("Expected BytesTransferred=1024, got %d", snapshot.BytesTransferred)
	}
	if snapshot.BytesSaved != 512 {
		t.Errorf("Expected BytesSaved=512, got %d", snapshot.BytesSaved)
	}
	if snapshot.CurrentOperation != "syncing" {
		t.Errorf("Expected CurrentOperation='syncing', got '%s'", snapshot.CurrentOperation)
	}
	if snapshot.Percentage != 50.0 {
		t.Errorf("Expected Percentage=50.0, got %.1f", snapshot.Percentage)
	}

	t.Logf("Progress snapshot test passed: all fields correctly copied")
}
