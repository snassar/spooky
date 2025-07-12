package main

import (
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

// --- Test generatePassword ---
func TestGeneratePassword(t *testing.T) {
	// Test that password is generated with correct length
	password := generatePassword()
	if len(password) != 25 {
		t.Errorf("expected password length 25, got %d", len(password))
	}

	// Test that password contains only valid characters
	const validChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for _, char := range password {
		if !strings.ContainsRune(validChars, char) {
			t.Errorf("password contains invalid character: %c", char)
		}
	}

	// Test that multiple passwords are different (randomness)
	password1 := generatePassword()
	password2 := generatePassword()
	if password1 == password2 {
		t.Error("generated passwords should be different")
	}
}

// --- Test generateEd25519KeyPair ---
func TestGenerateEd25519KeyPair(t *testing.T) {
	// Test successful key generation
	privateKeyPEM, publicKeyBytes, err := generateEd25519KeyPair("testpassword")
	if err != nil {
		t.Fatalf("failed to generate key pair: %v", err)
	}

	// Test private key PEM format
	if len(privateKeyPEM) == 0 {
		t.Error("private key PEM should not be empty")
	}

	// Verify PEM block structure
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		t.Error("failed to decode private key PEM")
	}
	if block.Type != "OPENSSH PRIVATE KEY" {
		t.Errorf("expected PEM type 'OPENSSH PRIVATE KEY', got '%s'", block.Type)
	}

	// Test public key format
	if len(publicKeyBytes) == 0 {
		t.Error("public key should not be empty")
	}

	// Verify public key starts with ssh-ed25519
	publicKeyStr := string(publicKeyBytes)
	if !strings.HasPrefix(publicKeyStr, "ssh-ed25519 ") {
		t.Errorf("public key should start with 'ssh-ed25519 ', got: %s", publicKeyStr[:20])
	}

	// Test that public key can be parsed
	_, _, _, _, err = ssh.ParseAuthorizedKey(publicKeyBytes)
	if err != nil {
		t.Errorf("failed to parse generated public key: %v", err)
	}

	// Test key data structure (should be 64 bytes: 32 private + 32 public)
	if len(block.Bytes) != 64 {
		t.Errorf("expected key data length 64, got %d", len(block.Bytes))
	}
}

func TestGenerateEd25519KeyPair_Consistency(t *testing.T) {
	// Test that multiple key pairs are different
	privateKey1, publicKey1, err := generateEd25519KeyPair("password1")
	if err != nil {
		t.Fatalf("failed to generate first key pair: %v", err)
	}

	privateKey2, publicKey2, err := generateEd25519KeyPair("password2")
	if err != nil {
		t.Fatalf("failed to generate second key pair: %v", err)
	}

	// Keys should be different
	if string(privateKey1) == string(privateKey2) {
		t.Error("generated private keys should be different")
	}
	if string(publicKey1) == string(publicKey2) {
		t.Error("generated public keys should be different")
	}
}

// --- Test writeKeyFiles ---
func TestWriteKeyFiles(t *testing.T) {
	// Create test data
	testPrivateKey := []byte("test private key content")
	testPublicKey := []byte("ssh-ed25519 test public key")
	testPassword := "testpassword123"

	// Test successful file writing
	outputDir, err := writeKeyFiles(testPrivateKey, testPublicKey, testPassword)
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

	outputDir, err := writeKeyFiles(testPrivateKey, testPublicKey, testPassword)
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

	outputDir, err := writeKeyFiles(testPrivateKey, testPublicKey, testPassword)
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

// --- Integration test for full workflow ---
func TestGenerateAndWriteKeys_Integration(t *testing.T) {
	// Test the full workflow: generate keys and write files
	password := generatePassword()

	privateKeyPEM, publicKey, err := generateEd25519KeyPair(password)
	if err != nil {
		t.Fatalf("failed to generate key pair: %v", err)
	}

	outputDir, err := writeKeyFiles(privateKeyPEM, publicKey, password)
	if err != nil {
		t.Fatalf("failed to write key files: %v", err)
	}

	// Clean up after test
	defer os.RemoveAll("generated_keys")

	// Verify the generated files can be read back
	privateKeyPath := filepath.Join(outputDir, "id_ed25519")
	if _, err := os.ReadFile(privateKeyPath); err != nil {
		t.Errorf("failed to read back private key file: %v", err)
	}

	publicKeyPath := filepath.Join(outputDir, "id_ed25519.pub")
	if _, err := os.ReadFile(publicKeyPath); err != nil {
		t.Errorf("failed to read back public key file: %v", err)
	}
}
