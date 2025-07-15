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

	// Verify directory name ends with a 3-digit number
	dirName := filepath.Base(outputDir)
	if len(dirName) != 11 { // YYYYMMDD + 3 digits
		t.Errorf("directory name should be 11 characters, got: %s", dirName)
	}
	if !strings.HasPrefix(dirName, expectedDate) {
		t.Errorf("directory name should start with date %s, got: %s", expectedDate, dirName)
	}
}

func TestWriteKeyFiles_FilePermissions(t *testing.T) {
	testPrivateKey := []byte("test private key")
	testPublicKey := []byte("test public key")
	testPassword := "testpass"

	outputDir, err := WriteKeyFiles(testPrivateKey, testPublicKey, testPassword)
	if err != nil {
		t.Fatalf("failed to write key files: %v", err)
	}

	// Clean up after test
	defer os.RemoveAll("generated_keys")

	// Check private key file permissions (should be 0600)
	privateKeyPath := filepath.Join(outputDir, "id_ed25519")
	privateKeyInfo, err := os.Stat(privateKeyPath)
	if err != nil {
		t.Errorf("failed to stat private key file: %v", err)
	}
	expectedMode := os.FileMode(0600)
	if privateKeyInfo.Mode() != expectedMode {
		t.Errorf("private key file should have mode %v, got: %v", expectedMode, privateKeyInfo.Mode())
	}

	// Check public key file permissions (should be 0644)
	publicKeyPath := filepath.Join(outputDir, "id_ed25519.pub")
	publicKeyInfo, err := os.Stat(publicKeyPath)
	if err != nil {
		t.Errorf("failed to stat public key file: %v", err)
	}
	expectedMode = os.FileMode(0644)
	if publicKeyInfo.Mode() != expectedMode {
		t.Errorf("public key file should have mode %v, got: %v", expectedMode, publicKeyInfo.Mode())
	}

	// Check password file permissions (should be 0600)
	passwordPath := filepath.Join(outputDir, "ed25519_password.txt")
	passwordInfo, err := os.Stat(passwordPath)
	if err != nil {
		t.Errorf("failed to stat password file: %v", err)
	}
	expectedMode = os.FileMode(0600)
	if passwordInfo.Mode() != expectedMode {
		t.Errorf("password file should have mode %v, got: %v", expectedMode, passwordInfo.Mode())
	}

	// Check combined file permissions (should be 0600)
	combinedPath := filepath.Join(outputDir, "ed25519_keys.txt")
	combinedInfo, err := os.Stat(combinedPath)
	if err != nil {
		t.Errorf("failed to stat combined file: %v", err)
	}
	expectedMode = os.FileMode(0600)
	if combinedInfo.Mode() != expectedMode {
		t.Errorf("combined file should have mode %v, got: %v", expectedMode, combinedInfo.Mode())
	}
}

func TestWriteKeyFiles_MultipleDirectories(t *testing.T) {
	// Test that multiple calls create different directories
	testPrivateKey := []byte("test private key")
	testPublicKey := []byte("test public key")
	testPassword := "testpass"

	// Clean up any existing directories
	os.RemoveAll("generated_keys")

	// Create first directory
	outputDir1, err := WriteKeyFiles(testPrivateKey, testPublicKey, testPassword)
	if err != nil {
		t.Fatalf("failed to write first key files: %v", err)
	}

	// Create second directory
	outputDir2, err := WriteKeyFiles(testPrivateKey, testPublicKey, testPassword)
	if err != nil {
		t.Fatalf("failed to write second key files: %v", err)
	}

	// Clean up after test
	defer os.RemoveAll("generated_keys")

	// Verify directories are different
	if outputDir1 == outputDir2 {
		t.Error("multiple calls should create different directories")
	}

	// Verify both directories exist
	if _, err := os.Stat(outputDir1); os.IsNotExist(err) {
		t.Errorf("first directory %s was not created", outputDir1)
	}
	if _, err := os.Stat(outputDir2); os.IsNotExist(err) {
		t.Errorf("second directory %s was not created", outputDir2)
	}

	// Verify both directories contain the expected files
	for _, dir := range []string{outputDir1, outputDir2} {
		expectedFiles := []string{"id_ed25519", "id_ed25519.pub", "ed25519_password.txt", "ed25519_keys.txt"}
		for _, filename := range expectedFiles {
			filepath := filepath.Join(dir, filename)
			if _, err := os.Stat(filepath); os.IsNotExist(err) {
				t.Errorf("expected file %s was not created in %s", filepath, dir)
			}
		}
	}
}

func TestWriteKeyFiles_EmptyContent(t *testing.T) {
	// Test with empty content
	outputDir, err := WriteKeyFiles([]byte{}, []byte{}, "")
	if err != nil {
		t.Fatalf("failed to write key files with empty content: %v", err)
	}

	// Clean up after test
	defer os.RemoveAll("generated_keys")

	// Verify files were created (even with empty content)
	expectedFiles := []string{"id_ed25519", "id_ed25519.pub", "ed25519_password.txt", "ed25519_keys.txt"}
	for _, filename := range expectedFiles {
		filepath := filepath.Join(outputDir, filename)
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			t.Errorf("expected file %s was not created", filepath)
		}
	}

	// Verify password file contains the empty password
	passwordPath := filepath.Join(outputDir, "ed25519_password.txt")
	passwordContent, err := os.ReadFile(passwordPath)
	if err != nil {
		t.Errorf("failed to read password file: %v", err)
	}
	expectedPasswordContent := "Generated password for Ed25519 key: \n"
	if string(passwordContent) != expectedPasswordContent {
		t.Errorf("password file content does not match. Expected: %s, Got: %s", expectedPasswordContent, string(passwordContent))
	}
}

func TestWriteKeyFiles_LargeContent(t *testing.T) {
	// Test with large content
	largePrivateKey := make([]byte, 10000)
	for i := range largePrivateKey {
		largePrivateKey[i] = byte(i % 256)
	}

	largePublicKey := make([]byte, 5000)
	for i := range largePublicKey {
		largePublicKey[i] = byte(i % 256)
	}

	largePassword := strings.Repeat("a", 1000)

	outputDir, err := WriteKeyFiles(largePrivateKey, largePublicKey, largePassword)
	if err != nil {
		t.Fatalf("failed to write key files with large content: %v", err)
	}

	// Clean up after test
	defer os.RemoveAll("generated_keys")

	// Verify files were created and contain the large content
	privateKeyPath := filepath.Join(outputDir, "id_ed25519")
	privateKeyContent, err := os.ReadFile(privateKeyPath)
	if err != nil {
		t.Errorf("failed to read private key file: %v", err)
	}
	if len(privateKeyContent) != len(largePrivateKey) {
		t.Errorf("private key file size mismatch. Expected: %d, Got: %d", len(largePrivateKey), len(privateKeyContent))
	}

	publicKeyPath := filepath.Join(outputDir, "id_ed25519.pub")
	publicKeyContent, err := os.ReadFile(publicKeyPath)
	if err != nil {
		t.Errorf("failed to read public key file: %v", err)
	}
	if len(publicKeyContent) != len(largePublicKey) {
		t.Errorf("public key file size mismatch. Expected: %d, Got: %d", len(largePublicKey), len(publicKeyContent))
	}

	passwordPath := filepath.Join(outputDir, "ed25519_password.txt")
	passwordContent, err := os.ReadFile(passwordPath)
	if err != nil {
		t.Errorf("failed to read password file: %v", err)
	}
	expectedPasswordContent := fmt.Sprintf("Generated password for Ed25519 key: %s\n", largePassword)
	if string(passwordContent) != expectedPasswordContent {
		t.Error("password file content does not match large password")
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

	// Verify file paths are correct
	if !strings.HasSuffix(keyFiles.PrivateKeyPath, "id_ed25519") {
		t.Errorf("private key path should end with 'id_ed25519', got: %s", keyFiles.PrivateKeyPath)
	}

	if !strings.HasSuffix(keyFiles.PublicKeyPath, "id_ed25519.pub") {
		t.Errorf("public key path should end with 'id_ed25519.pub', got: %s", keyFiles.PublicKeyPath)
	}

	if !strings.HasSuffix(keyFiles.PasswordPath, "ed25519_password.txt") {
		t.Errorf("password path should end with 'ed25519_password.txt', got: %s", keyFiles.PasswordPath)
	}

	if !strings.HasSuffix(keyFiles.CombinedPath, "ed25519_keys.txt") {
		t.Errorf("combined path should end with 'ed25519_keys.txt', got: %s", keyFiles.CombinedPath)
	}

	// Verify all paths are in the same directory
	dir := filepath.Dir(keyFiles.PrivateKeyPath)
	if filepath.Dir(keyFiles.PublicKeyPath) != dir {
		t.Error("all files should be in the same directory")
	}
	if filepath.Dir(keyFiles.PasswordPath) != dir {
		t.Error("all files should be in the same directory")
	}
	if filepath.Dir(keyFiles.CombinedPath) != dir {
		t.Error("all files should be in the same directory")
	}
	if keyFiles.OutputDir != dir {
		t.Error("OutputDir should match the directory of the files")
	}
}

func TestWriteKeyPair_EmptyKeyPair(t *testing.T) {
	// Test with empty key pair
	keyPair := &KeyPair{
		PrivateKey: []byte{},
		PublicKey:  []byte{},
		Password:   "",
	}

	keyFiles, err := WriteKeyPair(keyPair)
	if err != nil {
		t.Fatalf("failed to write empty key pair: %v", err)
	}

	// Clean up after test
	defer os.RemoveAll("generated_keys")

	// Verify files were created
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
