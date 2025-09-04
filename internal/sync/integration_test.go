//go:build integration

package sync

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	testhelpers "spooky/internal/testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileSyncIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tests := []struct {
		name        string
		description string
		testFunc    func(t *testing.T)
	}{
		{
			name:        "TestFileSyncLocalToRemote",
			description: "Test file synchronization from local to remote",
			testFunc:    testFileSyncLocalToRemote,
		},
		{
			name:        "TestFileSyncRemoteToLocal",
			description: "Test file synchronization from remote to local",
			testFunc:    testFileSyncRemoteToLocal,
		},
		{
			name:        "TestFileSyncConflictResolution",
			description: "Test file sync conflict resolution",
			testFunc:    testFileSyncConflictResolution,
		},
		{
			name:        "TestFileSyncLargeFiles",
			description: "Test synchronization of large files",
			testFunc:    testFileSyncLargeFiles,
		},
		{
			name:        "TestFileSyncDirectoryStructure",
			description: "Test synchronization of directory structures",
			testFunc:    testFileSyncDirectoryStructure,
		},
		{
			name:        "TestFileSyncErrorHandling",
			description: "Test file sync error handling",
			testFunc:    testFileSyncErrorHandling,
		},
		{
			name:        "TestFileSyncConcurrentOperations",
			description: "Test concurrent file sync operations",
			testFunc:    testFileSyncConcurrentOperations,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Running integration test: %s", tt.description)
			tt.testFunc(t)
		})
	}
}

func testFileSyncLocalToRemote(t *testing.T) {
	// Create test file system
	testFS := testhelpers.NewTestFileSystem(t)
	defer testFS.Cleanup()

	// Create test files
	testContent1 := []byte("test content 1")
	testContent2 := []byte("test content 2")

	err := testFS.CreateFile("file1.txt", testContent1)
	require.NoError(t, err)

	err = testFS.CreateFile("subdir/file2.txt", testContent2)
	require.NoError(t, err)

	// Create target directory
	targetDir, err := os.MkdirTemp("", "sync-target-*")
	require.NoError(t, err)
	defer os.RemoveAll(targetDir)

	// Test local file sync
	options := &Options{
		BlockLength:   1024,
		CreateBackup:  false,
		PreservePerms: true,
		PreserveOwner: false,
		PreserveGroup: false,
	}

	result, err := File(testFS.RootDir(), targetDir, options)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, int(result.FilesProcessed))

	// Verify files were synced
	targetFile1 := filepath.Join(targetDir, "file1.txt")
	targetFile2 := filepath.Join(targetDir, "subdir", "file2.txt")

	assert.True(t, testFS.FileExists("file1.txt"))
	assert.True(t, testFS.FileExists("subdir/file2.txt"))

	// Verify content
	content1, err := os.ReadFile(targetFile1)
	require.NoError(t, err)
	assert.Equal(t, testContent1, content1)

	content2, err := os.ReadFile(targetFile2)
	require.NoError(t, err)
	assert.Equal(t, testContent2, content2)
}

func testFileSyncRemoteToLocal(t *testing.T) {
	// This test would require a mock SSH server
	// For now, we'll test the local sync functionality

	// Create source directory with files
	sourceDir, err := os.MkdirTemp("", "sync-source-*")
	require.NoError(t, err)
	defer os.RemoveAll(sourceDir)

	// Create test files
	testContent := []byte("remote content")
	testFile := filepath.Join(sourceDir, "remote_file.txt")
	err = os.WriteFile(testFile, testContent, 0644)
	require.NoError(t, err)

	// Create target directory
	targetDir, err := os.MkdirTemp("", "sync-target-*")
	require.NoError(t, err)
	defer os.RemoveAll(targetDir)

	// Test sync
	options := &Options{
		BlockLength:   1024,
		CreateBackup:  false,
		PreservePerms: true,
	}

	result, err := File(sourceDir, targetDir, options)
	require.NoError(t, err)
	assert.NotNil(t, result)

	// Verify file was synced
	targetFile := filepath.Join(targetDir, "remote_file.txt")
	assert.FileExists(t, targetFile)

	content, err := os.ReadFile(targetFile)
	require.NoError(t, err)
	assert.Equal(t, testContent, content)
}

func testFileSyncConflictResolution(t *testing.T) {
	// Create source and target directories
	sourceDir, err := os.MkdirTemp("", "sync-source-*")
	require.NoError(t, err)
	defer os.RemoveAll(sourceDir)

	targetDir, err := os.MkdirTemp("", "sync-target-*")
	require.NoError(t, err)
	defer os.RemoveAll(targetDir)

	// Create conflicting files
	sourceContent := []byte("source content")
	targetContent := []byte("target content")

	sourceFile := filepath.Join(sourceDir, "conflict.txt")
	targetFile := filepath.Join(targetDir, "conflict.txt")

	err = os.WriteFile(sourceFile, sourceContent, 0644)
	require.NoError(t, err)

	err = os.WriteFile(targetFile, targetContent, 0644)
	require.NoError(t, err)

	// Test sync with conflict resolution
	options := &Options{
		BlockLength:   1024,
		CreateBackup:  true, // Enable backup for conflicts
		PreservePerms: true,
	}

	result, err := File(sourceDir, targetDir, options)
	require.NoError(t, err)
	assert.NotNil(t, result)

	// Verify conflict was resolved (source should win by default)
	content, err := os.ReadFile(targetFile)
	require.NoError(t, err)
	assert.Equal(t, sourceContent, content)

	// Verify backup was created
	backupFile := targetFile + ".backup"
	assert.FileExists(t, backupFile)

	backupContent, err := os.ReadFile(backupFile)
	require.NoError(t, err)
	assert.Equal(t, targetContent, backupContent)
}

func testFileSyncLargeFiles(t *testing.T) {
	// Create test file system
	testFS := testhelpers.NewTestFileSystem(t)
	defer testFS.Cleanup()

	// Create a large file (1MB)
	largeContent := make([]byte, 1024*1024) // 1MB
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}

	err := testFS.CreateFile("large_file.bin", largeContent)
	require.NoError(t, err)

	// Create target directory
	targetDir, err := os.MkdirTemp("", "sync-target-*")
	require.NoError(t, err)
	defer os.RemoveAll(targetDir)

	// Test sync with large file
	options := &Options{
		BlockLength:   4096, // 4KB blocks
		CreateBackup:  false,
		PreservePerms: true,
	}

	start := time.Now()
	result, err := File(testFS.RootDir(), targetDir, options)
	duration := time.Since(start)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, int(result.FilesProcessed))

	// Verify large file was synced correctly
	targetFile := filepath.Join(targetDir, "large_file.bin")
	assert.FileExists(t, targetFile)

	content, err := os.ReadFile(targetFile)
	require.NoError(t, err)
	assert.Equal(t, largeContent, content)

	// Log performance
	t.Logf("Synced 1MB file in %v", duration)
}

func testFileSyncDirectoryStructure(t *testing.T) {
	// Create test file system with complex directory structure
	testFS := testhelpers.NewTestFileSystem(t)
	defer testFS.Cleanup()

	// Create nested directory structure
	directories := []string{
		"level1",
		"level1/level2",
		"level1/level2/level3",
		"another_dir",
		"another_dir/subdir",
	}

	for _, dir := range directories {
		err := testFS.CreateFile(dir+"/.gitkeep", []byte(""))
		require.NoError(t, err)
	}

	// Create files at different levels
	files := map[string][]byte{
		"root_file.txt":           []byte("root content"),
		"level1/file1.txt":        []byte("level1 content"),
		"level1/level2/file2.txt": []byte("level2 content"),
		"another_dir/file3.txt":   []byte("another content"),
	}

	for path, content := range files {
		err := testFS.CreateFile(path, content)
		require.NoError(t, err)
	}

	// Create target directory
	targetDir, err := os.MkdirTemp("", "sync-target-*")
	require.NoError(t, err)
	defer os.RemoveAll(targetDir)

	// Test directory sync
	options := &Options{
		BlockLength:   1024,
		CreateBackup:  false,
		PreservePerms: true,
	}

	result, err := Directory(testFS.RootDir(), targetDir, ModeOneWayReplica, options)
	require.NoError(t, err)
	assert.NotNil(t, result)

	// Verify directory structure was preserved
	for _, dir := range directories {
		targetPath := filepath.Join(targetDir, dir)
		assert.DirExists(t, targetPath)
	}

	// Verify all files were synced
	for path, expectedContent := range files {
		targetFile := filepath.Join(targetDir, path)
		assert.FileExists(t, targetFile)

		content, err := os.ReadFile(targetFile)
		require.NoError(t, err)
		assert.Equal(t, expectedContent, content)
	}
}

func testFileSyncErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		sourcePath  string
		targetPath  string
		expectError bool
		errorType   string
	}{
		{
			name:        "Non-existent source",
			sourcePath:  "/non/existent/path",
			targetPath:  "/tmp/target",
			expectError: true,
			errorType:   "source not found",
		},
		{
			name:        "Invalid target path",
			sourcePath:  "/tmp",
			targetPath:  "/invalid/target/path/that/cannot/be/created",
			expectError: true,
			errorType:   "target creation",
		},
		{
			name:        "Identical paths",
			sourcePath:  "/tmp",
			targetPath:  "/tmp",
			expectError: true,
			errorType:   "identical paths",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := &Options{
				BlockLength:   1024,
				CreateBackup:  false,
				PreservePerms: true,
			}

			result, err := File(tt.sourcePath, tt.targetPath, options)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func testFileSyncConcurrentOperations(t *testing.T) {
	// Create test file system
	testFS := testhelpers.NewTestFileSystem(t)
	defer testFS.Cleanup()

	// Create multiple test files
	numFiles := 10
	for i := 0; i < numFiles; i++ {
		content := []byte(fmt.Sprintf("content %d", i))
		filename := fmt.Sprintf("file_%d.txt", i)
		err := testFS.CreateFile(filename, content)
		require.NoError(t, err)
	}

	// Create target directories for concurrent operations
	numTargets := 3
	targetDirs := make([]string, numTargets)
	for i := 0; i < numTargets; i++ {
		targetDir, err := os.MkdirTemp("", fmt.Sprintf("sync-target-%d-*", i))
		require.NoError(t, err)
		targetDirs[i] = targetDir
		defer os.RemoveAll(targetDir)
	}

	// Run concurrent sync operations
	results := make(chan error, numTargets)

	for i := 0; i < numTargets; i++ {
		go func(targetDir string) {
			options := &Options{
				BlockLength:   1024,
				CreateBackup:  false,
				PreservePerms: true,
			}

			_, err := File(testFS.RootDir(), targetDir, options)
			results <- err
		}(targetDirs[i])
	}

	// Collect results
	var errors []error
	for i := 0; i < numTargets; i++ {
		if err := <-results; err != nil {
			errors = append(errors, err)
		}
	}

	// All operations should succeed
	assert.Empty(t, errors, "Concurrent sync operations should not fail")

	// Verify all target directories have the same files
	for _, targetDir := range targetDirs {
		files, err := testFS.ListFiles()
		require.NoError(t, err)

		for _, file := range files {
			targetFile := filepath.Join(targetDir, file)
			assert.FileExists(t, targetFile)
		}
	}
}

// TestEndToEndSyncIntegration tests complete end-to-end sync scenarios
func TestEndToEndSyncIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping end-to-end integration test in short mode")
	}

	t.Run("CompleteSyncWorkflow", func(t *testing.T) {
		// Create source directory with various file types
		sourceDir, err := os.MkdirTemp("", "sync-source-*")
		require.NoError(t, err)
		defer os.RemoveAll(sourceDir)

		// Create test files of different types
		testFiles := map[string][]byte{
			"text.txt":          []byte("This is a text file"),
			"binary.bin":        {0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD},
			"empty.txt":         []byte(""),
			"large.txt":         make([]byte, 10000), // 10KB file
			"subdir/nested.txt": []byte("nested content"),
		}

		for path, content := range testFiles {
			fullPath := filepath.Join(sourceDir, path)
			dir := filepath.Dir(fullPath)

			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatalf("Failed to create directory: %v", err)
			}

			if err := os.WriteFile(fullPath, content, 0644); err != nil {
				t.Fatalf("Failed to create file: %v", err)
			}
		}

		// Create target directory
		targetDir, err := os.MkdirTemp("", "sync-target-*")
		require.NoError(t, err)
		defer os.RemoveAll(targetDir)

		// Perform complete sync
		options := &Options{
			BlockLength:   1024,
			CreateBackup:  true,
			PreservePerms: true,
			PreserveOwner: false,
			PreserveGroup: false,
		}

		result, err := Directory(sourceDir, targetDir, ModeOneWayReplica, options)
		require.NoError(t, err)
		assert.NotNil(t, result)

		// Verify complete sync
		for path, expectedContent := range testFiles {
			targetFile := filepath.Join(targetDir, path)
			assert.FileExists(t, targetFile)

			content, err := os.ReadFile(targetFile)
			require.NoError(t, err)
			assert.Equal(t, expectedContent, content)
		}

		// Verify directory structure
		assert.DirExists(t, filepath.Join(targetDir, "subdir"))

		t.Logf("Successfully synced %d files", len(testFiles))
	})
}
