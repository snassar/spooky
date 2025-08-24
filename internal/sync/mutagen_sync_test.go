package sync

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMutagenSync(t *testing.T) {
	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "spooky-mutagen-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test 1: New file sync (simple copy)
	t.Run("NewFile", func(t *testing.T) {
		sourcePath := filepath.Join(tempDir, "new_source.txt")
		targetPath := filepath.Join(tempDir, "new_target.txt")

		// Create source file
		sourceData := []byte("Hello, this is a new file for testing Mutagen sync!")
		if err := os.WriteFile(sourcePath, sourceData, 0644); err != nil {
			t.Fatalf("Failed to create source file: %v", err)
		}

		// Test sync
		options := DefaultSyncOptions()
		options.Verbose = true
		result, err := SyncFile(sourcePath, targetPath, options)

		if err != nil {
			t.Fatalf("SyncFile failed: %v", err)
		}

		if !result.Success {
			t.Fatalf("Sync operation failed: %v", result.Error)
		}

		// Verify target file was created
		targetData, err := os.ReadFile(targetPath)
		if err != nil {
			t.Fatalf("Failed to read target file: %v", err)
		}

		if string(targetData) != string(sourceData) {
			t.Errorf("Target file content mismatch. Expected: %q, Got: %q", sourceData, targetData)
		}

		t.Logf("New file sync successful: %d bytes transferred", result.BytesTransferred)
	})

	// Test 2: Identical files
	t.Run("IdenticalFiles", func(t *testing.T) {
		sourcePath := filepath.Join(tempDir, "identical_source.txt")
		targetPath := filepath.Join(tempDir, "identical_target.txt")

		// Create identical files
		data := []byte("This content is identical in both files.")
		if err := os.WriteFile(sourcePath, data, 0644); err != nil {
			t.Fatalf("Failed to create source file: %v", err)
		}
		if err := os.WriteFile(targetPath, data, 0644); err != nil {
			t.Fatalf("Failed to create target file: %v", err)
		}

		// Test sync
		options := DefaultSyncOptions()
		options.Verbose = true
		result, err := SyncFile(sourcePath, targetPath, options)

		if err != nil {
			t.Fatalf("SyncFile failed: %v", err)
		}

		if !result.Success {
			t.Fatalf("Sync operation failed: %v", result.Error)
		}

		if result.BytesTransferred != 0 {
			t.Errorf("Expected 0 bytes transferred for identical files, got %d", result.BytesTransferred)
		}

		t.Logf("Identical files test successful: no transfer needed")
	})

	// Test 3: Similar files (efficient delta sync)
	t.Run("SimilarFiles", func(t *testing.T) {
		sourcePath := filepath.Join(tempDir, "similar_source.txt")
		targetPath := filepath.Join(tempDir, "similar_target.txt")

		// Create source file
		sourceData := []byte("This is the source file with some content that will be modified slightly.")
		if err := os.WriteFile(sourcePath, sourceData, 0644); err != nil {
			t.Fatalf("Failed to create source file: %v", err)
		}

		// Create target file with similar but different content
		targetData := []byte("This is the target file with some content that will be modified differently.")
		if err := os.WriteFile(targetPath, targetData, 0644); err != nil {
			t.Fatalf("Failed to create target file: %v", err)
		}

		// Test sync
		options := DefaultSyncOptions()
		options.Verbose = true
		result, err := SyncFile(sourcePath, targetPath, options)

		if err != nil {
			t.Fatalf("SyncFile failed: %v", err)
		}

		if !result.Success {
			t.Fatalf("Sync operation failed: %v", result.Error)
		}

		// Verify target file has been updated to match source
		finalData, err := os.ReadFile(targetPath)
		if err != nil {
			t.Fatalf("Failed to read target file: %v", err)
		}

		if string(finalData) != string(sourceData) {
			t.Errorf("Sync failed - target file content mismatch.\nExpected: %q\nGot: %q", sourceData, finalData)
		}

		t.Logf("Similar files sync successful: %d bytes transferred, %d bytes saved, %d operations",
			result.BytesTransferred, result.BytesSaved, result.Operations)
	})

	// Test 4: Large file differences
	t.Run("LargeFiles", func(t *testing.T) {
		sourcePath := filepath.Join(tempDir, "large_source.txt")
		targetPath := filepath.Join(tempDir, "large_target.txt")

		// Create larger files for more realistic testing
		sourceData := make([]byte, 10*1024) // 10KB
		for i := range sourceData {
			sourceData[i] = byte(i % 256)
		}
		// Modify the source slightly
		copy(sourceData[5000:5010], []byte("DIFFERENT!"))

		if err := os.WriteFile(sourcePath, sourceData, 0644); err != nil {
			t.Fatalf("Failed to create source file: %v", err)
		}

		// Create target with mostly same content but some differences
		targetData := make([]byte, 10*1024)
		for i := range targetData {
			targetData[i] = byte(i % 256)
		}
		// Different modification in target
		copy(targetData[5000:5010], []byte("original!!"))

		if err := os.WriteFile(targetPath, targetData, 0644); err != nil {
			t.Fatalf("Failed to create target file: %v", err)
		}

		// Test sync
		options := DefaultSyncOptions()
		options.Verbose = true
		result, err := SyncFile(sourcePath, targetPath, options)

		if err != nil {
			t.Fatalf("SyncFile failed: %v", err)
		}

		if !result.Success {
			t.Fatalf("Sync operation failed: %v", result.Error)
		}

		// Verify target file has been updated to match source
		finalData, err := os.ReadFile(targetPath)
		if err != nil {
			t.Fatalf("Failed to read target file: %v", err)
		}

		if len(finalData) != len(sourceData) {
			t.Errorf("Target file size mismatch. Expected: %d, Got: %d", len(sourceData), len(finalData))
		}

		// Check that the files are now identical
		for i := range sourceData {
			if finalData[i] != sourceData[i] {
				t.Errorf("File content mismatch at position %d. Expected: %d, Got: %d", i, sourceData[i], finalData[i])
				break
			}
		}

		// Check efficiency - should transfer much less than the full file size
		efficiency := float64(result.BytesTransferred) / float64(len(sourceData))
		t.Logf("Large files sync successful: %d bytes transferred (%.1f%% of file), %d bytes saved, %d operations",
			result.BytesTransferred, efficiency*100, result.BytesSaved, result.Operations)

		if efficiency > 0.5 {
			t.Logf("Warning: Transfer efficiency was %.1f%% - might indicate suboptimal rsync performance", efficiency*100)
		}
	})

	// Test 5: Dry run
	t.Run("DryRun", func(t *testing.T) {
		sourcePath := filepath.Join(tempDir, "dryrun_source.txt")
		targetPath := filepath.Join(tempDir, "dryrun_target.txt")

		// Create source file
		sourceData := []byte("This is a dry run test file.")
		if err := os.WriteFile(sourcePath, sourceData, 0644); err != nil {
			t.Fatalf("Failed to create source file: %v", err)
		}

		// Test dry run
		options := DefaultSyncOptions()
		options.DryRun = true
		options.Verbose = true
		result, err := SyncFile(sourcePath, targetPath, options)

		if err != nil {
			t.Fatalf("SyncFile failed: %v", err)
		}

		if !result.Success {
			t.Fatalf("Dry run operation failed: %v", result.Error)
		}

		// Verify target file was NOT created
		if _, err := os.Stat(targetPath); !os.IsNotExist(err) {
			t.Errorf("Target file should not exist after dry run")
		}

		t.Logf("Dry run successful: would transfer %d bytes", result.BytesTransferred)
	})
}

func TestMutagenSyncEngine(t *testing.T) {
	// Test engine creation
	engine := NewMutagenSyncEngine()
	if engine == nil {
		t.Fatal("Failed to create Mutagen sync engine")
	}

	if engine.engine == nil {
		t.Fatal("Mutagen rsync engine not initialized")
	}

	t.Log("Mutagen sync engine created successfully")
}
