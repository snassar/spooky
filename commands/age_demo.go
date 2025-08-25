package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"spooky/internal/encryption"

	"github.com/spf13/cobra"
)

var ageDemoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Demonstrate age encryption system",
	Long: `Demonstrate the age encryption system with key loading and encryption/decryption.

This command shows how the age encryption system works, including:
- Loading age identities and recipients
- Encrypting and decrypting data
- Structured encryption format
- HCL configuration format

Note: This demo requires the age CLI to be installed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("=== Age Encryption Demo ===")
		fmt.Println()

		// Step 1: Check if age CLI is available
		fmt.Println("Step 1: Checking age CLI availability")
		if !commandExists("age-keygen") {
			fmt.Println("❌ age-keygen command not found")
			fmt.Println("   Please install age: https://github.com/FiloSottile/age")
			fmt.Println("   On Linux/macOS: brew install age")
			fmt.Println("   On Ubuntu/Debian: sudo apt install age")
			return nil
		}
		fmt.Println("✅ age-keygen is available")
		fmt.Println()

		// Step 2: Create temporary directory for demo keys
		fmt.Println("Step 2: Creating demo keys")
		tempDir, err := os.MkdirTemp("", "age-demo-*")
		if err != nil {
			return fmt.Errorf("failed to create temp directory: %w", err)
		}
		defer os.RemoveAll(tempDir)

		// Generate demo keys
		if err := generateDemoKeys(tempDir); err != nil {
			fmt.Printf("❌ Failed to generate demo keys: %v\n", err)
			fmt.Println("   This demo requires the age CLI to be installed and working")
			return nil
		}

		identityPath := filepath.Join(tempDir, "identity.txt")
		recipientsPath := filepath.Join(tempDir, "recipients.txt")

		fmt.Printf("✅ Generated demo keys in: %s\n", tempDir)
		fmt.Printf("   Identity: %s\n", identityPath)
		fmt.Printf("   Recipients: %s\n", recipientsPath)
		fmt.Println()

		// Step 3: Load keys and test encryption
		fmt.Println("Step 3: Testing age encryption")
		ae, err := encryption.NewAgeEncryption(identityPath, recipientsPath)
		if err != nil {
			return fmt.Errorf("failed to create age encryption: %w", err)
		}

		fmt.Printf("✅ Loaded %d identities and %d recipients\n",
			ae.GetIdentitiesCount(), ae.GetRecipientsCount())

		// Validate configuration
		if err := ae.ValidateConfiguration(); err != nil {
			return fmt.Errorf("configuration validation failed: %w", err)
		}
		fmt.Println("✅ Configuration is valid")
		fmt.Println()

		// Step 4: Demonstrate encryption/decryption
		fmt.Println("Step 4: Encryption/Decryption Demo")

		testData := "my-super-secret-password-123"
		fmt.Printf("Original data: %s\n", testData)

		// Encrypt
		encrypted, err := ae.Encrypt(testData)
		if err != nil {
			return fmt.Errorf("encryption failed: %w", err)
		}

		fmt.Printf("✅ Encrypted successfully\n")
		fmt.Printf("Encrypted data (first 100 chars): %s...\n", truncateString(encrypted, 100))

		// Verify it looks encrypted
		if !ae.IsEncrypted(encrypted) {
			return fmt.Errorf("encrypted data doesn't appear to be age-encrypted")
		}
		fmt.Println("✅ Data appears to be age-encrypted")

		// Decrypt
		decrypted, err := ae.Decrypt(encrypted)
		if err != nil {
			return fmt.Errorf("decryption failed: %w", err)
		}

		// Verify decryption
		if decrypted != testData {
			return fmt.Errorf("decryption verification failed: expected '%s', got '%s'", testData, decrypted)
		}

		fmt.Printf("✅ Decrypted successfully\n")
		fmt.Printf("Decrypted data: %s\n", decrypted)
		fmt.Println()

		// Step 5: Demonstrate structured encryption format
		fmt.Println("Step 5: Structured Encryption Format")

		// Create structured encrypted value
		structuredValue := &encryption.EncryptedValueV1{
			Data:        extractBase64FromArmored(encrypted),
			Format:      "base64",
			Algorithm:   "age",
			Version:     "v1",
			EncryptedAt: "2024-01-15T10:30:00Z",
		}

		fmt.Printf("Structured encrypted value:\n")
		fmt.Printf("  Data: %s...\n", truncateString(structuredValue.Data, 50))
		fmt.Printf("  Format: %s\n", structuredValue.Format)
		fmt.Printf("  Algorithm: %s\n", structuredValue.Algorithm)
		fmt.Printf("  Version: %s\n", structuredValue.Version)
		fmt.Printf("  EncryptedAt: %s\n", structuredValue.EncryptedAt)

		// Verify no headers/footers in data
		if strings.Contains(structuredValue.Data, "-----BEGIN AGE ENCRYPTED FILE-----") {
			fmt.Println("  ❌ Data contains AGE headers (should not)")
		} else {
			fmt.Println("  ✅ Data does not contain AGE headers (correct)")
		}

		if strings.Contains(structuredValue.Data, "-----END AGE ENCRYPTED FILE-----") {
			fmt.Println("  ❌ Data contains AGE footers (should not)")
		} else {
			fmt.Println("  ✅ Data does not contain AGE footers (correct)")
		}
		fmt.Println()

		// Step 6: Show HCL configuration format
		fmt.Println("Step 6: HCL Configuration Format")
		fmt.Println("Example HCL configuration with structured encryption:")
		fmt.Println("```hcl")
		fmt.Println("variables {")
		fmt.Println("  variable \"database_password\" {")
		fmt.Println("    description = \"Database password\"")
		fmt.Println("    encrypted_value = {")
		fmt.Printf("      data = \"%s...\"\n", truncateString(structuredValue.Data, 50))
		fmt.Printf("      format = \"%s\"\n", structuredValue.Format)
		fmt.Printf("      algorithm = \"%s\"\n", structuredValue.Algorithm)
		fmt.Printf("      version = \"%s\"\n", structuredValue.Version)
		fmt.Printf("      encrypted_at = \"%s\"\n", structuredValue.EncryptedAt)
		fmt.Println("    }")
		fmt.Println("    encrypted = false  # No longer needed")
		fmt.Println("    sensitive = true   # Automatically set")
		fmt.Println("  }")
		fmt.Println("}")
		fmt.Println("```")
		fmt.Println()

		fmt.Println("=== Age Encryption Demo Complete ===")
		fmt.Println()
		fmt.Println("Key Features Demonstrated:")
		fmt.Println("1. ✅ Age key loading (identity and recipients)")
		fmt.Println("2. ✅ Encryption and decryption")
		fmt.Println("3. ✅ Configuration validation")
		fmt.Println("4. ✅ Structured encryption format")
		fmt.Println("5. ✅ HCL configuration format")
		fmt.Println("6. ✅ No AGE headers/footers in structured data")

		return nil
	},
}

// Helper functions

// generateDemoKeys generates demo age keys
func generateDemoKeys(tempDir string) error {
	// Generate identity
	identityPath := filepath.Join(tempDir, "identity.txt")
	identityCmd := fmt.Sprintf("age-keygen -o %s", identityPath)
	if err := runCommand(identityCmd); err != nil {
		return fmt.Errorf("failed to generate identity: %w", err)
	}

	// Extract recipient from identity
	recipientsPath := filepath.Join(tempDir, "recipients.txt")
	recipientCmd := fmt.Sprintf("age-keygen -y %s > %s", identityPath, recipientsPath)
	if err := runCommand(recipientCmd); err != nil {
		return fmt.Errorf("failed to extract recipient: %w", err)
	}

	return nil
}

// extractBase64FromArmored extracts the base64 content from armored AGE format
func extractBase64FromArmored(armored string) string {
	lines := strings.Split(armored, "\n")
	var base64Content []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "-----BEGIN AGE ENCRYPTED FILE-----" &&
			line != "-----END AGE ENCRYPTED FILE-----" &&
			line != "" {
			base64Content = append(base64Content, line)
		}
	}

	return strings.Join(base64Content, "")
}

func init() {
	ageCmd.AddCommand(ageDemoCmd)
}
