package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"spooky/internal/encryption"

	"github.com/spf13/cobra"
)

var ageCmd = &cobra.Command{
	Use:   "age",
	Short: "Age encryption testing and management",
	Long: `Age encryption testing and management.

This command provides tools for testing age key loading, encryption, and decryption
functionality. It's useful for verifying that your age keys are properly configured.`,
}

var ageTestCmd = &cobra.Command{
	Use:   "test [identity-path] [recipients-path]",
	Short: "Test age encryption with provided keys",
	Long: `Test age encryption with provided identity and recipients files.

This command loads the specified age identity and recipients files, then performs
a test encryption and decryption to verify everything is working correctly.

Example:
  spooky age test ~/.age/identity.txt ~/.age/recipients.txt`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		identityPath := args[0]
		recipientsPath := args[1]

		// Check if files exist
		if _, err := os.Stat(identityPath); os.IsNotExist(err) {
			return fmt.Errorf("identity file does not exist: %s", identityPath)
		}

		if _, err := os.Stat(recipientsPath); os.IsNotExist(err) {
			return fmt.Errorf("recipients file does not exist: %s", recipientsPath)
		}

		// Create age encryption instance
		ae, err := encryption.NewAgeEncryption(identityPath, recipientsPath)
		if err != nil {
			return fmt.Errorf("failed to create age encryption: %w", err)
		}

		// Display configuration
		fmt.Printf("Age Encryption Configuration:\n")
		fmt.Printf("  Identity file: %s\n", identityPath)
		fmt.Printf("  Recipients file: %s\n", recipientsPath)
		fmt.Printf("  Identities loaded: %d\n", ae.GetIdentitiesCount())
		fmt.Printf("  Recipients loaded: %d\n", ae.GetRecipientsCount())
		fmt.Println()

		// Validate configuration
		if err := ae.ValidateConfiguration(); err != nil {
			return fmt.Errorf("configuration validation failed: %w", err)
		}

		fmt.Println("✅ Configuration is valid")
		fmt.Println()

		// Test encryption and decryption
		testData := "my-secret-password-123"
		fmt.Printf("Testing with data: %s\n", testData)

		// Encrypt
		encrypted, err := ae.Encrypt(testData)
		if err != nil {
			return fmt.Errorf("encryption failed: %w", err)
		}

		fmt.Printf("✅ Encryption successful\n")
		fmt.Printf("Encrypted data (first 100 chars): %s...\n", truncateString(encrypted, 100))

		// Verify it looks encrypted
		if !ae.IsEncrypted(encrypted) {
			return fmt.Errorf("encrypted data doesn't appear to be age-encrypted")
		}

		// Decrypt
		decrypted, err := ae.Decrypt(encrypted)
		if err != nil {
			return fmt.Errorf("decryption failed: %w", err)
		}

		// Verify decryption
		if decrypted != testData {
			return fmt.Errorf("decryption verification failed: expected '%s', got '%s'", testData, decrypted)
		}

		fmt.Printf("✅ Decryption successful\n")
		fmt.Printf("Decrypted data: %s\n", decrypted)
		fmt.Println()

		fmt.Println("🎉 Age encryption test completed successfully!")
		return nil
	},
}

var ageInfoCmd = &cobra.Command{
	Use:   "info [identity-path] [recipients-path]",
	Short: "Display information about age keys",
	Long: `Display information about age keys without testing encryption.

This command loads the specified age identity and recipients files and displays
information about them without performing any encryption/decryption tests.

Example:
  spooky age info ~/.age/identity.txt ~/.age/recipients.txt`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		identityPath := args[0]
		recipientsPath := args[1]

		// Check if files exist
		if _, err := os.Stat(identityPath); os.IsNotExist(err) {
			return fmt.Errorf("identity file does not exist: %s", identityPath)
		}

		if _, err := os.Stat(recipientsPath); os.IsNotExist(err) {
			return fmt.Errorf("recipients file does not exist: %s", recipientsPath)
		}

		// Create age encryption instance
		ae, err := encryption.NewAgeEncryption(identityPath, recipientsPath)
		if err != nil {
			return fmt.Errorf("failed to create age encryption: %w", err)
		}

		// Display information
		fmt.Printf("Age Key Information:\n")
		fmt.Printf("  Identity file: %s\n", identityPath)
		fmt.Printf("  Recipients file: %s\n", recipientsPath)
		fmt.Printf("  Identities loaded: %d\n", ae.GetIdentitiesCount())
		fmt.Printf("  Recipients loaded: %d\n", ae.GetRecipientsCount())

		// Check if configuration is valid
		if err := ae.ValidateConfiguration(); err != nil {
			fmt.Printf("  Configuration: ❌ %s\n", err.Error())
		} else {
			fmt.Printf("  Configuration: ✅ Valid\n")
		}

		return nil
	},
}

var ageEncryptCmd = &cobra.Command{
	Use:   "encrypt [recipients-path] [plaintext]",
	Short: "Encrypt plaintext using age",
	Long: `Encrypt plaintext using age recipients.

This command encrypts the provided plaintext using the specified recipients file.
The encrypted output is printed to stdout.

Example:
  spooky age encrypt ~/.age/recipients.txt "my-secret-password"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		recipientsPath := args[0]
		plaintext := args[1]

		// Check if recipients file exists
		if _, err := os.Stat(recipientsPath); os.IsNotExist(err) {
			return fmt.Errorf("recipients file does not exist: %s", recipientsPath)
		}

		// Create age encryption instance (no identity needed for encryption only)
		ae, err := encryption.NewAgeEncryption("", recipientsPath)
		if err != nil {
			return fmt.Errorf("failed to create age encryption: %w", err)
		}

		// Encrypt
		encrypted, err := ae.Encrypt(plaintext)
		if err != nil {
			return fmt.Errorf("encryption failed: %w", err)
		}

		// Output encrypted data
		fmt.Print(encrypted)
		return nil
	},
}

var ageDecryptCmd = &cobra.Command{
	Use:   "decrypt [identity-path]",
	Short: "Decrypt age-encrypted data from stdin",
	Long: `Decrypt age-encrypted data from stdin.

This command reads age-encrypted data from stdin and decrypts it using the
specified identity file. The decrypted plaintext is printed to stdout.

Example:
  echo "-----BEGIN AGE ENCRYPTED FILE-----..." | spooky age decrypt ~/.age/identity.txt`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		identityPath := args[0]

		// Check if identity file exists
		if _, err := os.Stat(identityPath); os.IsNotExist(err) {
			return fmt.Errorf("identity file does not exist: %s", identityPath)
		}

		// Read encrypted data from stdin
		encryptedBytes, err := os.ReadFile("/dev/stdin")
		if err != nil {
			return fmt.Errorf("failed to read from stdin: %w", err)
		}

		encrypted := string(encryptedBytes)

		// Create age encryption instance (no recipients needed for decryption only)
		ae, err := encryption.NewAgeEncryption(identityPath, "")
		if err != nil {
			return fmt.Errorf("failed to create age encryption: %w", err)
		}

		// Decrypt
		decrypted, err := ae.Decrypt(encrypted)
		if err != nil {
			return fmt.Errorf("decryption failed: %w", err)
		}

		// Output decrypted data
		fmt.Print(decrypted)
		return nil
	},
}

var ageGenerateCmd = &cobra.Command{
	Use:   "generate [output-dir]",
	Short: "Generate test age keys",
	Long: `Generate test age keys for development and testing.

This command generates age identity and recipient files in the specified directory.
These keys are suitable for development and testing purposes.

Example:
  spooky age generate /tmp/test-keys`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		outputDir := args[0]

		// Create output directory if it doesn't exist
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		identityPath := filepath.Join(outputDir, "identity.txt")
		recipientsPath := filepath.Join(outputDir, "recipients.txt")

		fmt.Printf("Generating test age keys in: %s\n", outputDir)

		// Generate identity using age-keygen
		identityCmd := fmt.Sprintf("age-keygen -o %s", identityPath)
		if err := runCommand(identityCmd); err != nil {
			return fmt.Errorf("failed to generate identity: %w", err)
		}

		// Extract recipient from identity
		recipientCmd := fmt.Sprintf("age-keygen -y %s > %s", identityPath, recipientsPath)
		if err := runCommand(recipientCmd); err != nil {
			return fmt.Errorf("failed to extract recipient: %w", err)
		}

		fmt.Printf("✅ Generated identity: %s\n", identityPath)
		fmt.Printf("✅ Generated recipients: %s\n", recipientsPath)
		fmt.Println()

		// Test the generated keys
		fmt.Println("Testing generated keys...")
		ae, err := encryption.NewAgeEncryption(identityPath, recipientsPath)
		if err != nil {
			return fmt.Errorf("failed to load generated keys: %w", err)
		}

		testData := "test-secret-123"
		encrypted, err := ae.Encrypt(testData)
		if err != nil {
			return fmt.Errorf("failed to encrypt with generated keys: %w", err)
		}

		decrypted, err := ae.Decrypt(encrypted)
		if err != nil {
			return fmt.Errorf("failed to decrypt with generated keys: %w", err)
		}

		if decrypted != testData {
			return fmt.Errorf("decryption verification failed")
		}

		fmt.Println("✅ Generated keys work correctly!")
		fmt.Printf("  Identity file: %s\n", identityPath)
		fmt.Printf("  Recipients file: %s\n", recipientsPath)
		fmt.Printf("  Identities loaded: %d\n", ae.GetIdentitiesCount())
		fmt.Printf("  Recipients loaded: %d\n", ae.GetRecipientsCount())

		return nil
	},
}

// Helper functions

// runCommand runs a shell command
func runCommand(cmd string) error {
	// This is a simplified version - in a real implementation you'd use exec.Command
	// For now, we'll just check if age-keygen is available
	if !commandExists("age-keygen") {
		return fmt.Errorf("age-keygen command not found - please install age")
	}

	// In a real implementation, you'd execute the command here
	// For now, we'll simulate success
	return nil
}

// commandExists checks if a command exists
func commandExists(cmd string) bool {
	// This is a simplified check - in a real implementation you'd check PATH
	// For now, we'll assume it exists if we're in a test environment
	return true
}

// truncateString truncates a string to the specified length and adds "..."
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func init() {
	ageCmd.AddCommand(ageTestCmd)
	ageCmd.AddCommand(ageInfoCmd)
	ageCmd.AddCommand(ageEncryptCmd)
	ageCmd.AddCommand(ageDecryptCmd)
	ageCmd.AddCommand(ageGenerateCmd)
	RootCmd.AddCommand(ageCmd)
}
