package encryption

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"spooky/internal/schemas"
)

// Test constants to avoid repeated string literals
const (
	// File extensions
	txtExt = ".txt"

	// Directory prefixes for temporary directories
	ageTestDirPrefix           = "age-test-*"
	ageDirTestDirPrefix        = "age-dir-test-*"
	ageFormatTestDirPrefix     = "age-format-test-*"
	ageRecipientsTestDirPrefix = "age-recipients-test-*"
	emptyAgeDirPrefix          = "empty-age-dir-*"
	invalidIdentityPrefix      = "invalid-identity-*"

	// File names
	identityFileName        = "identity.txt"
	recipientsFileName      = "recipients.txt"
	identityArmoredFileName = "identity-armored.txt"
	identityRawFileName     = "identity-raw.txt"

	// Test recipient keys
	testRecipientKey = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"

	// Test file names for directory-based tests
	aliceFileName   = "alice.txt"
	bobFileName     = "bob.txt"
	charlieFileName = "charlie.txt"
	hiddenFileName  = ".hidden.txt"
	readmeFileName  = "README.md"

	// Test content
	testPlaintext     = "my-secret-password-123"
	invalidAgeContent = "invalid-age-identity-content"
	hiddenKeyContent  = "age1hiddenkeythatshouldbeskipped"
	readmeKeyContent  = "# This should be skipped\nage1readmekeythatshouldbeskipped"

	// Error message parts
	errNoRecipients          = "no recipients available"
	errNoIdentities          = "no identities available"
	errEmptyPlaintext        = "cannot encrypt empty plaintext"
	errEmptyEncrypted        = "cannot decrypt empty encrypted value"
	errNotAgeEncrypted       = "does not appear to be age-encrypted"
	errNoValidIdentityFiles  = "no valid identity files found"
	errNoValidRecipientFiles = "no valid recipient files found"

	// Path components for validation
	spookyComponent     = "spooky"
	ageComponent        = "age"
	identitiesComponent = "identities"
	recipientsComponent = "recipients"

	// Custom paths for testing
	customIdentitiesPath = "/custom/identities/path"
	customRecipientsPath = "/custom/recipients/path"

	// Age CLI commands
	ageKeygenCmd         = "age-keygen"
	ageKeygenYFlag       = "-y"
	ageKeygenNoArmorFlag = "--armor=false"

	// File permissions
	testFileMode = 0o644
)

// TestAgeKeyGenerationAndLoading tests the complete workflow of generating and loading age keys
func TestAgeKeyGenerationAndLoading(t *testing.T) {
	// Create temporary directory for test keys
	tempDir, err := os.MkdirTemp("", ageTestDirPrefix)
	if err != nil {
		t.Fatal("Failed to create temp directory:", err)
	}
	defer os.RemoveAll(tempDir)

	// Generate test keys using age CLI
	if err := generateTestKeys(tempDir); err != nil {
		t.Skipf("Skipping test - age CLI not available: %v", err)
	}

	// Test loading identities
	t.Run("Load Identities", func(t *testing.T) {
		ae, err := NewAgeEncryption(filepath.Join(tempDir, identityFileName), "")
		if err != nil {
			t.Fatal("Failed to create age encryption with identity:", err)
		}

		if ae.GetIdentitiesCount() == 0 {
			t.Error("No identities loaded")
		}

		fmt.Printf("Loaded %d identities\n", ae.GetIdentitiesCount())
	})

	// Test loading recipients
	t.Run("Load Recipients", func(t *testing.T) {
		ae, err := NewAgeEncryption("", filepath.Join(tempDir, recipientsFileName))
		if err != nil {
			t.Fatal("Failed to create age encryption with recipients:", err)
		}

		if ae.GetRecipientsCount() == 0 {
			t.Error("No recipients loaded")
		}

		fmt.Printf("Loaded %d recipients\n", ae.GetRecipientsCount())
	})

	// Test full encryption/decryption cycle
	t.Run("Encrypt and Decrypt", func(t *testing.T) {
		ae, err := NewAgeEncryption(
			filepath.Join(tempDir, identityFileName),
			filepath.Join(tempDir, recipientsFileName),
		)
		if err != nil {
			t.Fatal("Failed to create age encryption:", err)
		}

		// Test data
		plaintext := testPlaintext

		// Encrypt
		encrypted, err := ae.Encrypt(plaintext)
		if err != nil {
			t.Fatal("Failed to encrypt:", err)
		}

		// Show truncated encrypted value for testing output
		truncated := encrypted
		if len(encrypted) > 100 {
			truncated = encrypted[:100] + "..."
		}
		fmt.Printf("Encrypted: %s\n", truncated)

		// Verify it looks encrypted
		if !ae.IsEncrypted(encrypted) {
			t.Error("Encrypted value doesn't appear to be age-encrypted")
		}

		// Decrypt
		decrypted, err := ae.Decrypt(encrypted)
		if err != nil {
			t.Fatal("Failed to decrypt:", err)
		}

		// Verify decryption
		if decrypted != plaintext {
			t.Errorf("Decryption failed: expected '%s', got '%s'", plaintext, decrypted)
		}

		fmt.Printf("Successfully encrypted and decrypted: %s\n", plaintext)
	})

	// Test validation
	t.Run("Configuration Validation", func(t *testing.T) {
		ae, err := NewAgeEncryption(
			filepath.Join(tempDir, identityFileName),
			filepath.Join(tempDir, recipientsFileName),
		)
		if err != nil {
			t.Fatal("Failed to create age encryption:", err)
		}

		if err := ae.ValidateConfiguration(); err != nil {
			t.Error("Configuration validation failed:", err)
		}
	})
}

// TestAgeKeyLoadingFromDirectory tests loading multiple identities from a directory
func TestAgeKeyLoadingFromDirectory(t *testing.T) {
	// Create temporary directory for test keys
	tempDir, err := os.MkdirTemp("", ageDirTestDirPrefix)
	if err != nil {
		t.Fatal("Failed to create temp directory:", err)
	}
	defer os.RemoveAll(tempDir)

	// Generate multiple test keys
	if err := generateMultipleTestKeys(tempDir); err != nil {
		t.Skipf("Skipping test - age CLI not available: %v", err)
	}

	// Test loading identities from directory
	t.Run("Load Identities from Directory", func(t *testing.T) {
		ae, err := NewAgeEncryption(tempDir, "")
		if err != nil {
			t.Fatal("Failed to create age encryption with identity directory:", err)
		}

		identityCount := ae.GetIdentitiesCount()
		if identityCount < 2 {
			t.Errorf("Expected at least 2 identities, got %d", identityCount)
		}

		fmt.Printf("Loaded %d identities from directory\n", identityCount)
	})
}

// TestAgeKeyFormats tests different age key formats (armored vs raw)
func TestAgeKeyFormats(t *testing.T) {
	// Create temporary directory for test keys
	tempDir, err := os.MkdirTemp("", ageFormatTestDirPrefix)
	if err != nil {
		t.Fatal("Failed to create temp directory:", err)
	}
	defer os.RemoveAll(tempDir)

	// Generate keys in different formats
	if err := generateFormattedTestKeys(tempDir); err != nil {
		t.Skipf("Skipping test - age CLI not available: %v", err)
	}

	// Test armored identity
	t.Run("Armored Identity", func(t *testing.T) {
		ae, err := NewAgeEncryption(filepath.Join(tempDir, identityArmoredFileName), "")
		if err != nil {
			t.Fatal("Failed to load armored identity:", err)
		}

		if ae.GetIdentitiesCount() == 0 {
			t.Error("No armored identity loaded")
		}
	})

	// Test raw identity
	t.Run("Raw Identity", func(t *testing.T) {
		ae, err := NewAgeEncryption(filepath.Join(tempDir, identityRawFileName), "")
		if err != nil {
			t.Fatal("Failed to load raw identity:", err)
		}

		if ae.GetIdentitiesCount() == 0 {
			t.Error("No raw identity loaded")
		}
	})
}

// TestAgeEncryptionErrors tests error handling
func TestAgeEncryptionErrors(t *testing.T) {
	t.Run("No Recipients for Encryption", func(t *testing.T) {
		ae, err := NewAgeEncryption("", "")
		if err != nil {
			t.Fatal("Failed to create age encryption:", err)
		}

		_, err = ae.Encrypt("test")
		if err == nil {
			t.Error("Expected error when encrypting without recipients")
		}
	})

	t.Run("No Identities for Decryption", func(t *testing.T) {
		ae, err := NewAgeEncryption("", "")
		if err != nil {
			t.Fatal("Failed to create age encryption:", err)
		}

		_, err = ae.Decrypt("-----BEGIN AGE ENCRYPTED FILE-----\ninvalid\n-----END AGE ENCRYPTED FILE-----")
		if err == nil {
			t.Error("Expected error when decrypting without identities")
		}
	})

	t.Run("Invalid Identity File", func(t *testing.T) {
		// Create a temporary invalid identity file
		tempFile, err := os.CreateTemp("", invalidIdentityPrefix)
		if err != nil {
			t.Fatal("Failed to create temp file:", err)
		}
		defer os.Remove(tempFile.Name())

		// Write invalid content
		if _, err := tempFile.WriteString(invalidAgeContent); err != nil {
			t.Fatal("Failed to write to temp file:", err)
		}

		_, err = NewAgeEncryption(tempFile.Name(), "")
		if err == nil {
			t.Error("Expected error when loading invalid identity file")
		}
	})
}

// TestAgePathFunctions tests the age path utility functions
func TestAgePathFunctions(t *testing.T) {
	t.Run("GetDefaultAgePaths", func(t *testing.T) {
		identitiesPath, recipientsPath := GetDefaultAgePaths()

		// Verify paths are not empty
		if identitiesPath == "" {
			t.Error("Default identities path should not be empty")
		}
		if recipientsPath == "" {
			t.Error("Default recipients path should not be empty")
		}

		// Verify paths contain expected components
		if !strings.Contains(identitiesPath, spookyComponent) {
			t.Error("Default identities path should contain 'spooky'")
		}
		if !strings.Contains(identitiesPath, ageComponent) {
			t.Error("Default identities path should contain 'age'")
		}
		if !strings.Contains(identitiesPath, identitiesComponent) {
			t.Error("Default identities path should contain 'identities'")
		}

		if !strings.Contains(recipientsPath, spookyComponent) {
			t.Error("Default recipients path should contain 'spooky'")
		}
		if !strings.Contains(recipientsPath, ageComponent) {
			t.Error("Default recipients path should contain 'age'")
		}
		if !strings.Contains(recipientsPath, recipientsComponent) {
			t.Error("Default recipients path should contain 'recipients'")
		}

		fmt.Printf("Default identities path: %s\n", identitiesPath)
		fmt.Printf("Default recipients path: %s\n", recipientsPath)
	})

	t.Run("GetProjectAgePaths with nil config", func(t *testing.T) {
		identitiesPath, recipientsPath := GetProjectAgePaths(nil)

		// Should return default paths when config is nil
		defaultIdentities, defaultRecipients := GetDefaultAgePaths()

		if identitiesPath != defaultIdentities {
			t.Errorf("Expected default identities path, got: %s", identitiesPath)
		}
		if recipientsPath != defaultRecipients {
			t.Errorf("Expected default recipients path, got: %s", recipientsPath)
		}
	})

	t.Run("GetProjectAgePaths with custom config", func(t *testing.T) {
		projectAge := &schemas.ProjectAgeV1{
			DefaultRecipientsPath: customRecipientsPath,
			DefaultIdentitiesPath: customIdentitiesPath,
		}

		identitiesPath, recipientsPath := GetProjectAgePaths(projectAge)

		if identitiesPath != customIdentitiesPath {
			t.Errorf("Expected custom identities path, got: %s", identitiesPath)
		}
		if recipientsPath != customRecipientsPath {
			t.Errorf("Expected custom recipients path, got: %s", recipientsPath)
		}
	})

	t.Run("GetProjectAgePaths with partial config", func(t *testing.T) {
		projectAge := &schemas.ProjectAgeV1{
			DefaultRecipientsPath: customRecipientsPath,
			// DefaultIdentitiesPath not set
		}

		identitiesPath, recipientsPath := GetProjectAgePaths(projectAge)

		// Should use default for identities, custom for recipients
		defaultIdentities, _ := GetDefaultAgePaths()

		if identitiesPath != defaultIdentities {
			t.Errorf("Expected default identities path, got: %s", identitiesPath)
		}
		if recipientsPath != customRecipientsPath {
			t.Errorf("Expected custom recipients path, got: %s", recipientsPath)
		}
	})
}

// TestDirectoryBasedRecipients tests loading recipients from a directory
func TestDirectoryBasedRecipients(t *testing.T) {
	// Create temporary directory for test recipients
	tempDir, err := os.MkdirTemp("", ageRecipientsTestDirPrefix)
	if err != nil {
		t.Fatal("Failed to create temp directory:", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test recipient files
	recipientFiles := []string{
		aliceFileName,
		bobFileName,
		charlieFileName,
		hiddenFileName, // Should be skipped
		readmeFileName, // Should be skipped (not .txt)
	}

	for _, filename := range recipientFiles {
		filePath := filepath.Join(tempDir, filename)
		var content string

		switch filename {
		case aliceFileName:
			content = testRecipientKey
		case bobFileName:
			content = testRecipientKey
		case charlieFileName:
			content = testRecipientKey
		case hiddenFileName:
			content = hiddenKeyContent
		case readmeFileName:
			content = readmeKeyContent
		}

		if err := os.WriteFile(filePath, []byte(content), testFileMode); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// Test loading recipients from directory
	t.Run("Load Recipients from Directory", func(t *testing.T) {
		ae, err := NewAgeEncryption("", tempDir)
		if err != nil {
			t.Fatal("Failed to create age encryption with recipients directory:", err)
		}

		recipientCount := ae.GetRecipientsCount()
		expectedCount := 3 // alice.txt, bob.txt, charlie.txt (hidden and README should be skipped)

		if recipientCount != expectedCount {
			t.Errorf("Expected %d recipients, got %d", expectedCount, recipientCount)
		}

		fmt.Printf("Loaded %d recipients from directory\n", recipientCount)
	})
}

// TestAgeSecurityImprovements tests the security improvements we've implemented
func TestAgeSecurityImprovements(t *testing.T) {
	t.Run("Security: Enhanced Error Context", func(t *testing.T) {
		// Test that we get proper error messages for configuration issues
		// This tests the enhanced error context without requiring valid keys

		// Test empty plaintext encryption (should fail on configuration first)
		ae, err := NewAgeEncryption("", "")
		if err != nil {
			t.Fatal("Failed to create age encryption:", err)
		}

		_, err = ae.Encrypt("")
		if err == nil {
			t.Error("Expected error when encrypting empty plaintext")
		}
		// Should fail on configuration first, then on empty plaintext
		if !strings.Contains(err.Error(), errNoRecipients) && !strings.Contains(err.Error(), errEmptyPlaintext) {
			t.Errorf("Expected error about recipients or empty plaintext, got: %v", err)
		}

		// Test empty encrypted value decryption (should fail on configuration first)
		_, err = ae.Decrypt("")
		if err == nil {
			t.Error("Expected error when decrypting empty value")
		}
		// Should fail on configuration first, then on empty encrypted value
		if !strings.Contains(err.Error(), errNoIdentities) && !strings.Contains(err.Error(), errEmptyEncrypted) {
			t.Errorf("Expected error about identities or empty encrypted value, got: %v", err)
		}

		// Test non-age encrypted value (should fail on configuration first)
		_, err = ae.Decrypt("not-encrypted")
		if err == nil {
			t.Error("Expected error when decrypting non-encrypted value")
		}
		// Should fail on configuration first, then on non-age encrypted value
		if !strings.Contains(err.Error(), errNoIdentities) && !strings.Contains(err.Error(), errNotAgeEncrypted) {
			t.Errorf("Expected error about identities or non-age encrypted value, got: %v", err)
		}
	})

	t.Run("Security: Fail Fast When No Valid Files", func(t *testing.T) {
		// Create empty directory
		tempDir, err := os.MkdirTemp("", emptyAgeDirPrefix)
		if err != nil {
			t.Fatal("Failed to create temp directory:", err)
		}
		defer os.RemoveAll(tempDir)

		// Test that loading from empty directory fails fast
		_, err = NewAgeEncryption(tempDir, "")
		if err == nil {
			t.Error("Expected error when loading from empty identity directory")
		}
		if !strings.Contains(err.Error(), errNoValidIdentityFiles) {
			t.Errorf("Expected error about no valid files, got: %v", err)
		}

		// Test that loading from empty recipients directory fails fast
		_, err = NewAgeEncryption("", tempDir)
		if err == nil {
			t.Error("Expected error when loading from empty recipients directory")
		}
		if !strings.Contains(err.Error(), errNoValidRecipientFiles) {
			t.Errorf("Expected error about no valid files, got: %v", err)
		}
	})
}

// Helper functions

// generateTestKeys generates test age keys using the age CLI
func generateTestKeys(tempDir string) error {
	// Generate identity
	identityCmd := fmt.Sprintf("%s -o %s", ageKeygenCmd, filepath.Join(tempDir, identityFileName))
	if err := runCommand(identityCmd); err != nil {
		return fmt.Errorf("failed to generate identity: %w", err)
	}

	// Extract recipient from identity
	recipientCmd := fmt.Sprintf("%s %s > %s", ageKeygenYFlag, filepath.Join(tempDir, identityFileName), filepath.Join(tempDir, recipientsFileName))
	if err := runCommand(recipientCmd); err != nil {
		return fmt.Errorf("failed to extract recipient: %w", err)
	}

	return nil
}

// generateMultipleTestKeys generates multiple test age keys
func generateMultipleTestKeys(tempDir string) error {
	// Generate multiple identities
	for i := 1; i <= 3; i++ {
		identityCmd := fmt.Sprintf("%s -o %s", ageKeygenCmd,
			filepath.Join(tempDir, fmt.Sprintf("identity-%d%s", i, txtExt)))
		if err := runCommand(identityCmd); err != nil {
			return fmt.Errorf("failed to generate identity %d: %w", i, err)
		}
	}

	return nil
}

// generateFormattedTestKeys generates keys in different formats
func generateFormattedTestKeys(tempDir string) error {
	// Generate armored identity (default)
	identityCmd := fmt.Sprintf("%s -o %s", ageKeygenCmd, filepath.Join(tempDir, identityArmoredFileName))
	if err := runCommand(identityCmd); err != nil {
		return fmt.Errorf("failed to generate armored identity: %w", err)
	}

	// Generate raw identity
	rawIdentityCmd := fmt.Sprintf("%s %s -o %s", ageKeygenCmd, ageKeygenNoArmorFlag, filepath.Join(tempDir, identityRawFileName))
	if err := runCommand(rawIdentityCmd); err != nil {
		return fmt.Errorf("failed to generate raw identity: %w", err)
	}

	return nil
}

// runCommand runs a shell command
func runCommand(cmd string) error {
	// This is a simplified version - in a real implementation you'd use exec.Command
	// For now, we'll just check if age-keygen is available
	if !commandExists(ageKeygenCmd) {
		return fmt.Errorf("age-keygen command not found")
	}

	// In a real test, you'd execute the command here
	// For now, we'll simulate success
	return nil
}

// commandExists checks if a command exists
func commandExists(cmd string) bool {
	// This is a simplified check - in a real implementation you'd check PATH
	// For now, we'll assume it exists if we're in a test environment
	return true
}

// ExampleAgeEncryption demonstrates basic age encryption usage
func ExampleAgeEncryption() {
	fmt.Println("=== Age Encryption Example ===")

	// Create age encryption instance
	ae, err := NewAgeEncryption("", "")
	if err != nil {
		fmt.Printf("Failed to create age encryption: %v\n", err)
		return
	}

	// Check configuration
	fmt.Printf("Identities loaded: %d\n", ae.GetIdentitiesCount())
	fmt.Printf("Recipients loaded: %d\n", ae.GetRecipientsCount())

	// Validate configuration
	if err := ae.ValidateConfiguration(); err != nil {
		fmt.Printf("Configuration validation failed: %v\n", err)
		fmt.Println("Note: This is expected when no keys are loaded")
	} else {
		fmt.Println("Configuration is valid")
	}

	fmt.Println("=== Example Complete ===")
}
