package keygen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateSSHKeys(t *testing.T) {
	// Test the default SSH key generation workflow
	err := GenerateSSHKeys()
	if err != nil {
		t.Fatalf("failed to generate SSH keys: %v", err)
	}

	// Clean up after test
	defer os.RemoveAll("generated_keys")

	// Verify that at least one key directory was created
	entries, err := os.ReadDir("generated_keys")
	if err != nil {
		t.Fatalf("failed to read generated_keys directory: %v", err)
	}

	if len(entries) == 0 {
		t.Error("no key directories were created")
	}

	// Check that at least one directory contains the expected files
	foundValidDirectory := false
	for _, entry := range entries {
		if entry.IsDir() {
			dirPath := filepath.Join("generated_keys", entry.Name())
			expectedFiles := []string{"id_ed25519", "id_ed25519.pub", "ed25519_password.txt", "ed25519_keys.txt"}

			allFilesExist := true
			for _, filename := range expectedFiles {
				filePath := filepath.Join(dirPath, filename)
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					allFilesExist = false
					break
				}
			}

			if allFilesExist {
				foundValidDirectory = true
				break
			}
		}
	}

	if !foundValidDirectory {
		t.Error("no valid key directory with all expected files was found")
	}
}

func TestGenerateSSHKeysWithType(t *testing.T) {
	testCases := []struct {
		name    string
		keyType KeyType
		wantErr bool
	}{
		{"Ed25519 key", Ed25519KeyType, false},
		{"invalid key type", KeyType("invalid"), true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean up any existing generated_keys before each test
			os.RemoveAll("generated_keys")

			err := GenerateSSHKeysWithType(tc.keyType)

			if tc.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Clean up after test
			defer os.RemoveAll("generated_keys")

			// Verify that at least one key directory was created
			entries, err := os.ReadDir("generated_keys")
			if err != nil {
				t.Fatalf("failed to read generated_keys directory: %v", err)
			}

			if len(entries) == 0 {
				t.Error("no key directories were created")
			}

			// Check that at least one directory contains the expected files
			foundValidDirectory := false
			for _, entry := range entries {
				if entry.IsDir() {
					dirPath := filepath.Join("generated_keys", entry.Name())
					expectedFiles := []string{"id_ed25519", "id_ed25519.pub", "ed25519_password.txt", "ed25519_keys.txt"}

					allFilesExist := true
					for _, filename := range expectedFiles {
						filePath := filepath.Join(dirPath, filename)
						if _, err := os.Stat(filePath); os.IsNotExist(err) {
							allFilesExist = false
							break
						}
					}

					if allFilesExist {
						foundValidDirectory = true
						break
					}
				}
			}

			if !foundValidDirectory {
				t.Error("no valid key directory with all expected files was found")
			}
		})
	}
}

func TestGenerateSSHKeys_Integration(t *testing.T) {
	// Test the full integration workflow
	err := GenerateSSHKeys()
	if err != nil {
		t.Fatalf("failed to generate SSH keys: %v", err)
	}

	// Clean up after test
	defer os.RemoveAll("generated_keys")

	// Find the generated key directory
	entries, err := os.ReadDir("generated_keys")
	if err != nil {
		t.Fatalf("failed to read generated_keys directory: %v", err)
	}

	if len(entries) == 0 {
		t.Fatal("no key directories were created")
	}

	// Get the first directory (should be the one we just created)
	var keyDir string
	for _, entry := range entries {
		if entry.IsDir() {
			keyDir = filepath.Join("generated_keys", entry.Name())
			break
		}
	}

	if keyDir == "" {
		t.Fatal("no key directory found")
	}

	// Verify the generated files can be read back
	privateKeyPath := filepath.Join(keyDir, "id_ed25519")
	privateKeyContent, err := os.ReadFile(privateKeyPath)
	if err != nil {
		t.Errorf("failed to read back private key file: %v", err)
	}
	if len(privateKeyContent) == 0 {
		t.Error("private key file is empty")
	}

	publicKeyPath := filepath.Join(keyDir, "id_ed25519.pub")
	publicKeyContent, err := os.ReadFile(publicKeyPath)
	if err != nil {
		t.Errorf("failed to read back public key file: %v", err)
	}
	if len(publicKeyContent) == 0 {
		t.Error("public key file is empty")
	}

	// Verify public key format
	publicKeyStr := string(publicKeyContent)
	if !strings.HasPrefix(publicKeyStr, "ssh-ed25519 ") {
		t.Errorf("public key should start with 'ssh-ed25519 ', got: %s", publicKeyStr[:20])
	}

	// Verify password file exists and contains a password
	passwordPath := filepath.Join(keyDir, "ed25519_password.txt")
	passwordContent, err := os.ReadFile(passwordPath)
	if err != nil {
		t.Errorf("failed to read password file: %v", err)
	}
	if len(passwordContent) == 0 {
		t.Error("password file is empty")
	}

	// Verify combined file exists and contains all components
	combinedPath := filepath.Join(keyDir, "ed25519_keys.txt")
	combinedContent, err := os.ReadFile(combinedPath)
	if err != nil {
		t.Errorf("failed to read combined file: %v", err)
	}
	combinedStr := string(combinedContent)
	if !strings.Contains(combinedStr, "=== ED25519 SSH Key Pair ===") {
		t.Error("combined file should contain header")
	}
	if !strings.Contains(combinedStr, string(privateKeyContent)) {
		t.Error("combined file should contain private key")
	}
	if !strings.Contains(combinedStr, string(publicKeyContent)) {
		t.Error("combined file should contain public key")
	}
}
