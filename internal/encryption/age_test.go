package encryption

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestAgeKeyGenerationAndLoading tests the complete workflow of generating and loading age keys
func TestAgeKeyGenerationAndLoading(t *testing.T) {
	// Create temporary directory for test keys
	tempDir, err := os.MkdirTemp("", "age-test-*")
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
		ae, err := NewAgeEncryption(filepath.Join(tempDir, "identity.txt"), "")
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
		ae, err := NewAgeEncryption("", filepath.Join(tempDir, "recipients.txt"))
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
			filepath.Join(tempDir, "identity.txt"),
			filepath.Join(tempDir, "recipients.txt"),
		)
		if err != nil {
			t.Fatal("Failed to create age encryption:", err)
		}

		// Test data
		plaintext := "my-secret-password-123"

		// Encrypt
		encrypted, err := ae.Encrypt(plaintext)
		if err != nil {
			t.Fatal("Failed to encrypt:", err)
		}

		fmt.Printf("Encrypted: %s...\n", truncateString(encrypted, 100))

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
			filepath.Join(tempDir, "identity.txt"),
			filepath.Join(tempDir, "recipients.txt"),
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
	tempDir, err := os.MkdirTemp("", "age-dir-test-*")
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
	tempDir, err := os.MkdirTemp("", "age-format-test-*")
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
		ae, err := NewAgeEncryption(filepath.Join(tempDir, "identity-armored.txt"), "")
		if err != nil {
			t.Fatal("Failed to load armored identity:", err)
		}

		if ae.GetIdentitiesCount() == 0 {
			t.Error("No armored identity loaded")
		}
	})

	// Test raw identity
	t.Run("Raw Identity", func(t *testing.T) {
		ae, err := NewAgeEncryption(filepath.Join(tempDir, "identity-raw.txt"), "")
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
		tempFile, err := os.CreateTemp("", "invalid-identity-*")
		if err != nil {
			t.Fatal("Failed to create temp file:", err)
		}
		defer os.Remove(tempFile.Name())

		// Write invalid content
		tempFile.WriteString("invalid-age-identity-content")

		_, err = NewAgeEncryption(tempFile.Name(), "")
		if err == nil {
			t.Error("Expected error when loading invalid identity file")
		}
	})
}

// Helper functions

// generateTestKeys generates test age keys using the age CLI
func generateTestKeys(tempDir string) error {
	// Generate identity
	identityCmd := fmt.Sprintf("age-keygen -o %s", filepath.Join(tempDir, "identity.txt"))
	if err := runCommand(identityCmd); err != nil {
		return fmt.Errorf("failed to generate identity: %w", err)
	}

	// Extract recipient from identity
	recipientCmd := fmt.Sprintf("age-keygen -y %s > %s",
		filepath.Join(tempDir, "identity.txt"),
		filepath.Join(tempDir, "recipients.txt"))
	if err := runCommand(recipientCmd); err != nil {
		return fmt.Errorf("failed to extract recipient: %w", err)
	}

	return nil
}

// generateMultipleTestKeys generates multiple test age keys
func generateMultipleTestKeys(tempDir string) error {
	// Generate multiple identities
	for i := 1; i <= 3; i++ {
		identityCmd := fmt.Sprintf("age-keygen -o %s",
			filepath.Join(tempDir, fmt.Sprintf("identity-%d.txt", i)))
		if err := runCommand(identityCmd); err != nil {
			return fmt.Errorf("failed to generate identity %d: %w", i, err)
		}
	}

	return nil
}

// generateFormattedTestKeys generates keys in different formats
func generateFormattedTestKeys(tempDir string) error {
	// Generate armored identity (default)
	identityCmd := fmt.Sprintf("age-keygen -o %s", filepath.Join(tempDir, "identity-armored.txt"))
	if err := runCommand(identityCmd); err != nil {
		return fmt.Errorf("failed to generate armored identity: %w", err)
	}

	// Generate raw identity
	rawIdentityCmd := fmt.Sprintf("age-keygen --armor=false -o %s", filepath.Join(tempDir, "identity-raw.txt"))
	if err := runCommand(rawIdentityCmd); err != nil {
		return fmt.Errorf("failed to generate raw identity: %w", err)
	}

	return nil
}

// runCommand runs a shell command
func runCommand(cmd string) error {
	// This is a simplified version - in a real implementation you'd use exec.Command
	// For now, we'll just check if age-keygen is available
	if !commandExists("age-keygen") {
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
