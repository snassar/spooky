package encryption

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewHCLUpdater(t *testing.T) {
	// Test with nil encryption (should still work)
	updater := NewHCLUpdater(nil)
	if updater == nil {
		t.Fatal("NewHCLUpdater returned nil")
	}
	if updater.ageEncryption != nil {
		t.Error("expected ageEncryption to be nil when passed nil")
	}

	// Test with actual encryption instance
	ageEncryption, err := NewAgeEncryption("", "")
	if err != nil {
		t.Skipf("skipping test due to age encryption error: %v", err)
	}

	updater = NewHCLUpdater(ageEncryption)
	if updater == nil {
		t.Fatal("NewHCLUpdater returned nil")
	}
	if updater.ageEncryption == nil {
		t.Error("expected ageEncryption to be set")
	}
}

func TestHCLUpdater_UpdateFile_InputValidation(t *testing.T) {
	updater := NewHCLUpdater(nil)

	tests := []struct {
		name        string
		filePath    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty file path",
			filePath:    "",
			expectError: true,
			errorMsg:    "file path cannot be empty",
		},
		{
			name:        "non-existent file",
			filePath:    "nonexistent.hcl",
			expectError: true,
			errorMsg:    "file does not exist: nonexistent.hcl",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := updater.UpdateFile(tt.filePath)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("expected error message %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestHCLUpdater_UpdateFile_ValidHCL(t *testing.T) {
	updater := NewHCLUpdater(nil)

	// Create a temporary HCL file
	tempDir := t.TempDir()
	hclFile := filepath.Join(tempDir, "test.hcl")

	// Write a simple HCL file
	hclContent := `
variables {
  variable "test" {
    value = "test value"
    encrypted = false
  }
}
`
	if err := os.WriteFile(hclFile, []byte(hclContent), 0644); err != nil {
		t.Fatalf("failed to create test HCL file: %v", err)
	}

	// Test updating the file (should succeed even with nil encryption)
	err := updater.UpdateFile(hclFile)
	if err != nil {
		t.Errorf("unexpected error updating valid HCL file: %v", err)
	}
}

func TestHCLUpdater_UpdateFile_InvalidHCL(t *testing.T) {
	updater := NewHCLUpdater(nil)

	// Create a temporary HCL file with invalid syntax
	tempDir := t.TempDir()
	hclFile := filepath.Join(tempDir, "invalid.hcl")

	// Write invalid HCL content
	invalidContent := `
variables {
  variable "test" {
    value = "test value"
    encrypted = false
  // Missing closing brace
`
	if err := os.WriteFile(hclFile, []byte(invalidContent), 0644); err != nil {
		t.Fatalf("failed to create test HCL file: %v", err)
	}

	// Test updating the file (should fail due to invalid HCL)
	err := updater.UpdateFile(hclFile)
	if err == nil {
		t.Errorf("expected error for invalid HCL file")
	}
}

func TestHCLUpdater_UpdateDirectory(t *testing.T) {
	updater := NewHCLUpdater(nil)

	// Create a temporary directory with HCL files
	tempDir := t.TempDir()

	// Create some HCL files
	hclFiles := []string{"test1.hcl", "test2.hcl", "not_hcl.txt"}
	for _, filename := range hclFiles {
		filePath := filepath.Join(tempDir, filename)
		content := `
variables {
  variable "test" {
    value = "test value"
    encrypted = false
  }
}
`
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file %s: %v", filename, err)
		}
	}

	// Test updating directory
	err := updater.UpdateDirectory(tempDir)
	if err != nil {
		t.Errorf("unexpected error updating directory: %v", err)
	}
}

func TestHCLUpdater_UpdateDirectory_NonExistent(t *testing.T) {
	updater := NewHCLUpdater(nil)

	// Test with non-existent directory
	err := updater.UpdateDirectory("nonexistent_directory")
	if err == nil {
		t.Errorf("expected error for non-existent directory")
	}
}
