package keygen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWriteKeyFiles(t *testing.T) {
	// Create test data
	testPrivateKey := []byte("test private key content")
	testPublicKey := []byte("ssh-ed25519 test public key")
	testPassword := "testpassword123"

	// Test successful file writing
	outputDir, err := WriteKeyFiles(testPrivateKey, testPublicKey, testPassword)
	if err != nil {
		t.Fatalf("failed to write key files: %v", err)
	}

	// Clean up after test
	defer os.RemoveAll("generated_keys")

	// Verify output directory exists
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Errorf("output directory %s was not created", outputDir)
	}

	// Verify all expected files exist
	expectedFiles := []string{
		"id_ed25519",
		"id_ed25519.pub",
		"ed25519_password.txt",
		"ed25519_keys.txt",
	}

	for _, filename := range expectedFiles {
		filepath := filepath.Join(outputDir, filename)
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			t.Errorf("expected file %s was not created", filepath)
		}
	}

	// Verify private key file content
	privateKeyPath := filepath.Join(outputDir, "id_ed25519")
	privateKeyContent, err := os.ReadFile(privateKeyPath)
	if err != nil {
		t.Errorf("failed to read private key file: %v", err)
	}
	if string(privateKeyContent) != string(testPrivateKey) {
		t.Error("private key file content does not match")
	}

	// Verify public key file content
	publicKeyPath := filepath.Join(outputDir, "id_ed25519.pub")
	publicKeyContent, err := os.ReadFile(publicKeyPath)
	if err != nil {
		t.Errorf("failed to read public key file: %v", err)
	}
	if string(publicKeyContent) != string(testPublicKey) {
		t.Error("public key file content does not match")
	}

	// Verify password file content
	passwordPath := filepath.Join(outputDir, "ed25519_password.txt")
	passwordContent, err := os.ReadFile(passwordPath)
	if err != nil {
		t.Errorf("failed to read password file: %v", err)
	}
	expectedPasswordContent := fmt.Sprintf("Generated password for Ed25519 key: %s\n", testPassword)
	if string(passwordContent) != expectedPasswordContent {
		t.Errorf("password file content does not match. Expected: %s, Got: %s", expectedPasswordContent, string(passwordContent))
	}

	// Verify combined file content
	combinedPath := filepath.Join(outputDir, "ed25519_keys.txt")
	combinedContent, err := os.ReadFile(combinedPath)
	if err != nil {
		t.Errorf("failed to read combined file: %v", err)
	}
	combinedStr := string(combinedContent)
	if !strings.Contains(combinedStr, "=== ED25519 SSH Key Pair ===") {
		t.Error("combined file should contain header")
	}
	if !strings.Contains(combinedStr, string(testPrivateKey)) {
		t.Error("combined file should contain private key")
	}
	if !strings.Contains(combinedStr, string(testPublicKey)) {
		t.Error("combined file should contain public key")
	}
	if !strings.Contains(combinedStr, testPassword) {
		t.Error("combined file should contain password")
	}
}

func TestWriteKeyFiles_DirectoryNaming(t *testing.T) {
	// Test that directory naming follows expected pattern
	testPrivateKey := []byte("test private key")
	testPublicKey := []byte("test public key")
	testPassword := "testpass"

	outputDir, err := WriteKeyFiles(testPrivateKey, testPublicKey, testPassword)
	if err != nil {
		t.Fatalf("failed to write key files: %v", err)
	}

	// Clean up after test
	defer os.RemoveAll("generated_keys")

	// Verify directory name format: generated_keys/YYYYMMDDNNN
	expectedDate := time.Now().Format("20060102")
	if !strings.Contains(outputDir, "generated_keys/") {
		t.Error("output directory should be in generated_keys/")
	}
	if !strings.Contains(outputDir, expectedDate) {
		t.Error("output directory should contain current date")
	}
}

func TestWriteKeyFiles_FilePermissions(t *testing.T) {
	// Test that files have correct permissions
	testPrivateKey := []byte("test private key")
	testPublicKey := []byte("test public key")
	testPassword := "testpass"

	outputDir, err := WriteKeyFiles(testPrivateKey, testPublicKey, testPassword)
	if err != nil {
		t.Fatalf("failed to write key files: %v", err)
	}

	// Clean up after test
	defer os.RemoveAll("generated_keys")

	privateKeyPath := filepath.Join(outputDir, "id_ed25519")
	publicKeyPath := filepath.Join(outputDir, "id_ed25519.pub")

	if os.PathSeparator == '\\' {
		// Windows: Check if file is writable (should be)
		privateFile, err := os.OpenFile(privateKeyPath, os.O_WRONLY, 0)
		if err != nil {
			t.Errorf("private key file should be writable on Windows: %v", err)
		} else {
			privateFile.Close()
		}
		publicFile, err := os.OpenFile(publicKeyPath, os.O_WRONLY, 0)
		if err != nil {
			t.Errorf("public key file should be writable on Windows: %v", err)
		} else {
			publicFile.Close()
		}
	} else {
		// Unix: Check file modes
		if info, err := os.Stat(privateKeyPath); err == nil {
			mode := info.Mode()
			if mode&0777 != 0600 {
				t.Errorf("private key file should have permissions 0600, got %o", mode&0777)
			}
		}
		if info, err := os.Stat(publicKeyPath); err == nil {
			mode := info.Mode()
			if mode&0777 != 0644 {
				t.Errorf("public key file should have permissions 0644, got %o", mode&0777)
			}
		}
	}
}

func TestWriteKeyPair(t *testing.T) {
	// Create test key pair
	keyPair := &KeyPair{
		PrivateKey: []byte("test private key content"),
		PublicKey:  []byte("ssh-ed25519 test public key"),
		Password:   "testpassword123",
	}

	// Test successful key pair writing
	keyFiles, err := WriteKeyPair(keyPair)
	if err != nil {
		t.Fatalf("failed to write key pair: %v", err)
	}

	// Clean up after test
	defer os.RemoveAll("generated_keys")

	// Verify KeyFiles struct is populated
	if keyFiles == nil {
		t.Fatal("keyFiles should not be nil")
	}

	if keyFiles.OutputDir == "" {
		t.Error("OutputDir should not be empty")
	}

	if keyFiles.PrivateKeyPath == "" {
		t.Error("PrivateKeyPath should not be empty")
	}

	if keyFiles.PublicKeyPath == "" {
		t.Error("PublicKeyPath should not be empty")
	}

	if keyFiles.PasswordPath == "" {
		t.Error("PasswordPath should not be empty")
	}

	if keyFiles.CombinedPath == "" {
		t.Error("CombinedPath should not be empty")
	}

	// Verify files exist
	if _, err := os.Stat(keyFiles.PrivateKeyPath); os.IsNotExist(err) {
		t.Errorf("private key file %s was not created", keyFiles.PrivateKeyPath)
	}

	if _, err := os.Stat(keyFiles.PublicKeyPath); os.IsNotExist(err) {
		t.Errorf("public key file %s was not created", keyFiles.PublicKeyPath)
	}

	if _, err := os.Stat(keyFiles.PasswordPath); os.IsNotExist(err) {
		t.Errorf("password file %s was not created", keyFiles.PasswordPath)
	}

	if _, err := os.Stat(keyFiles.CombinedPath); os.IsNotExist(err) {
		t.Errorf("combined file %s was not created", keyFiles.CombinedPath)
	}
}
