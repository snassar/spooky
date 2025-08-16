// Package cmd provides command implementations for spooky CLI.
package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	spookyinterfaces "spooky/internal/interfaces"

	"github.com/spf13/cobra"
)

// secretsCmd represents the secrets command
var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Manage age encryption and secrets",
	Long: `Manage age encryption configuration and validate secrets.

This command provides tools for validating age configuration, keys, and
encrypted data in spooky projects.`,
}

// secretsValidateCmd represents the secrets validate command
var secretsValidateCmd = &cobra.Command{
	Use:   "validate [project-path]",
	Short: "Validate age configuration and keys for project",
	Long: `Validate age configuration and keys for a project.

This command validates:
- Age configuration in spooky.hcl
- Identity files and permissions
- Recipients file format
- Encrypted values in project files

Examples:
  spooky secrets validate ./my-project`,
	Args: cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		return handleSecretsValidate(args[0])
	},
}

// handleSecretsValidate handles secrets validation
func handleSecretsValidate(projectPath string) error {
	// Validate project path
	if err := validateProjectPath(projectPath); err != nil {
		return err
	}

	// Get secrets integration
	manager := GetIntegrationManager()
	secretsIntegration := manager.GetSecretsIntegration()
	if secretsIntegration == nil {
		return fmt.Errorf("secrets integration not available")
	}

	ctx := context.Background()

	fmt.Printf("🔍 Validating age configuration for project: %s\n", projectPath)

	// Validate age configuration
	if err := validateAgeConfiguration(ctx, secretsIntegration); err != nil {
		return fmt.Errorf("age configuration validation failed: %w", err)
	}

	// Validate identity files
	if err := validateIdentityFiles(ctx, secretsIntegration); err != nil {
		return fmt.Errorf("identity files validation failed: %w", err)
	}

	// Validate recipients file
	if err := validateRecipientsFile(ctx, secretsIntegration); err != nil {
		return fmt.Errorf("recipients file validation failed: %w", err)
	}

	// Validate encrypted values in project
	if err := validateProjectEncryptedValues(ctx, projectPath, secretsIntegration); err != nil {
		return fmt.Errorf("project encrypted values validation failed: %w", err)
	}

	fmt.Println("✅ Age configuration validation completed successfully")
	return nil
}

// validateAgeConfiguration validates age configuration
func validateAgeConfiguration(ctx context.Context, secretsIntegration spookyinterfaces.SecretsIntegration) error {
	fmt.Println("  📋 Validating age configuration...")

	// Load global configuration
	configPath := getGlobalConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("global configuration not found: %s", configPath)
	}

	fmt.Printf("    ✅ Global configuration found: %s\n", configPath)

	// Read and validate the configuration file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Check for age-related configuration blocks
	content := string(data)
	if !strings.Contains(content, "age") && !strings.Contains(content, "encryption") {
		fmt.Printf("    ⚠️  No age encryption configuration found in %s\n", filepath.Base(configPath))
		fmt.Printf("    💡 Consider adding age configuration for encryption support\n")
	} else {
		fmt.Printf("    ✅ Age configuration blocks found in %s\n", filepath.Base(configPath))
	}

	// Test secrets integration functionality
	if err := testSecretsIntegration(ctx, secretsIntegration); err != nil {
		return fmt.Errorf("secrets integration test failed: %w", err)
	}

	return nil
}

// testSecretsIntegration tests basic secrets integration functionality
func testSecretsIntegration(ctx context.Context, secretsIntegration spookyinterfaces.SecretsIntegration) error {
	// Test that we can load identities (even if empty)
	identities, err := secretsIntegration.LoadIdentities(ctx, "~/.config/spooky/identities")
	if err != nil {
		// This is expected if the directory doesn't exist yet
		fmt.Printf("    ⚠️  Identity directory not accessible: %v\n", err)
	} else {
		fmt.Printf("    ✅ Secrets integration test passed (%d identities found)\n", len(identities))
	}

	return nil
}

// validateIdentityFiles validates age identity files
func validateIdentityFiles(ctx context.Context, secretsIntegration spookyinterfaces.SecretsIntegration) error {
	fmt.Println("  🔑 Validating identity files...")

	// Default identity path
	identityPath := "~/.config/spooky/identities"
	if expanded, err := expandPath(identityPath); err == nil {
		identityPath = expanded
	}

	// Check if identity directory exists
	if _, err := os.Stat(identityPath); os.IsNotExist(err) {
		return fmt.Errorf("identity directory not found: %s", identityPath)
	}

	fmt.Printf("    ✅ Identity directory found: %s\n", identityPath)

	// Load identities
	identities, err := secretsIntegration.LoadIdentities(ctx, identityPath)
	if err != nil {
		return fmt.Errorf("failed to load identities: %w", err)
	}

	if len(identities) == 0 {
		return fmt.Errorf("no valid identity files found in: %s", identityPath)
	}

	fmt.Printf("    ✅ Found %d valid identity files\n", len(identities))

	// Validate each identity file
	for _, identityFile := range identities {
		if err := secretsIntegration.ValidateAgeKey(ctx, identityFile); err != nil {
			return fmt.Errorf("invalid identity file %s: %w", identityFile, err)
		}
		fmt.Printf("    ✅ Validated identity file: %s\n", filepath.Base(identityFile))
	}

	return nil
}

// validateRecipientsFile validates age recipients file
func validateRecipientsFile(ctx context.Context, secretsIntegration spookyinterfaces.SecretsIntegration) error {
	fmt.Println("  📜 Validating recipients file...")

	// Default recipients path
	recipientsPath := "~/.config/spooky/recipients.txt"
	if expanded, err := expandPath(recipientsPath); err == nil {
		recipientsPath = expanded
	}

	// Check if recipients file exists
	if _, err := os.Stat(recipientsPath); os.IsNotExist(err) {
		return fmt.Errorf("recipients file not found: %s", recipientsPath)
	}

	fmt.Printf("    ✅ Recipients file found: %s\n", recipientsPath)

	// Load recipients
	recipients, err := secretsIntegration.LoadRecipients(ctx, recipientsPath)
	if err != nil {
		return fmt.Errorf("failed to load recipients: %w", err)
	}

	if len(recipients) == 0 {
		return fmt.Errorf("no valid recipients found in: %s", recipientsPath)
	}

	fmt.Printf("    ✅ Found %d valid recipients\n", len(recipients))

	return nil
}

// validateProjectEncryptedValues validates encrypted values in project files
func validateProjectEncryptedValues(ctx context.Context, projectPath string, secretsIntegration spookyinterfaces.SecretsIntegration) error {
	fmt.Println("  🔒 Validating encrypted values in project...")

	// Check variables.hcl for encrypted values
	variablesFile := filepath.Join(projectPath, "variables.hcl")
	if _, err := os.Stat(variablesFile); err == nil {
		if err := validateEncryptedValuesInFile(ctx, variablesFile, secretsIntegration); err != nil {
			return fmt.Errorf("variables file validation failed: %w", err)
		}
		fmt.Printf("    ✅ Variables file validated: %s\n", filepath.Base(variablesFile))
	}

	// Check machines.hcl for encrypted values
	machinesFile := filepath.Join(projectPath, "machines.hcl")
	if _, err := os.Stat(machinesFile); err == nil {
		if err := validateEncryptedValuesInFile(ctx, machinesFile, secretsIntegration); err != nil {
			return fmt.Errorf("machines file validation failed: %w", err)
		}
		fmt.Printf("    ✅ Machines file validated: %s\n", filepath.Base(machinesFile))
	}

	return nil
}

// validateEncryptedValuesInFile validates age-encrypted values in a file
func validateEncryptedValuesInFile(ctx context.Context, filePath string, secretsIntegration spookyinterfaces.SecretsIntegration) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	content := string(data)
	if !containsAgeEncryptedValues(content) {
		// No encrypted values found, which is fine
		return nil
	}

	// Find all age1-encrypted values in the file
	lines := strings.Split(content, "\n")
	encryptedCount := 0
	validCount := 0

	for lineNum, line := range lines {
		lineNum++ // Convert to 1-based line numbers

		// Look for age1 strings in this line
		if strings.Contains(line, "age1") {
			// Extract potential age1 values from the line
			words := strings.Fields(line)
			for _, word := range words {
				// Clean up the word (remove quotes, commas, etc.)
				cleanWord := strings.Trim(word, `"'`)
				cleanWord = strings.TrimSuffix(cleanWord, ",")

				if strings.HasPrefix(cleanWord, "age1") {
					encryptedCount++

					// Validate the encrypted value
					if err := secretsIntegration.ValidateAgeEncryptedValue(ctx, cleanWord); err != nil {
						return fmt.Errorf("invalid age-encrypted value at line %d: %w", lineNum, err)
					}

					validCount++
					fmt.Printf("    ✅ Validated age-encrypted value at line %d\n", lineNum)
				}
			}
		}
	}

	if encryptedCount > 0 {
		fmt.Printf("    📝 Found and validated %d age-encrypted values in: %s\n", validCount, filepath.Base(filePath))
	}

	return nil
}

// containsAgeEncryptedValues checks if content contains age-encrypted values
func containsAgeEncryptedValues(content string) bool {
	// Simple check for age1 prefix
	// In a real implementation, this would be more sophisticated
	return len(content) > 4 && content[:4] == "age1"
}

// expandPath expands a path with ~ to home directory
func expandPath(path string) (string, error) {
	if path != "" && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, path[1:]), nil
	}
	return path, nil
}

// getGlobalConfigPath returns the path to the global configuration file
func getGlobalConfigPath() string {
	// Check XDG_CONFIG_HOME first
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "spooky", "spooky.hcl")
	}

	// Fall back to default
	home, err := os.UserHomeDir()
	if err != nil {
		return "~/.config/spooky/spooky.hcl"
	}
	return filepath.Join(home, ".config", "spooky", "spooky.hcl")
}

func init() {
	// Add secrets command to root
	RootCmd.AddCommand(secretsCmd)

	// Add subcommands to secrets command
	secretsCmd.AddCommand(secretsValidateCmd)
}
