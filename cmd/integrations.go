// Package cmd provides command implementations for spooky CLI.
package cmd

import (
	"context"
	"fmt"

	spookycli "spooky/internal/cli/commands"
	spookyintegration "spooky/internal/integration"
	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookytypesconfig "spooky/internal/types/config"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypessecrets "spooky/internal/types/secrets"

	"github.com/spf13/cobra"
)

// Global instances for integration dependency injection
var (
	integrationManager spookyinterfaces.IntegrationManager
	integrationsLogger spookytypeslogging.Logger
)

// GetIntegrationManager returns the global integration manager instance
func GetIntegrationManager() spookyinterfaces.IntegrationManager {
	return integrationManager
}

// InitializeIntegrationsDependencies initializes integration-related dependencies
func InitializeIntegrationsDependencies() error {
	// Create log manager for integrations component
	logManager := spookylogging.NewLogManager()
	integrationsLogger = logManager.GetLogger("integrations")

	// Create default configuration for now
	config := &spookytypesconfig.Config{
		// Add default age config
		Age: &spookytypessecrets.AgeConfig{
			Identities: "~/.config/spooky/identities",
			Recipients: "~/.config/spooky/recipients.txt",
		},
	}

	// Create integration factory and manager
	factory := spookyintegration.NewFactory(integrationsLogger, config)
	integrationManager = factory.CreateIntegrationManager()

	return nil
}

// setupEncryption prepares the encryption environment
func setupEncryption() (spookyinterfaces.SecretsIntegration, error) {
	manager := GetIntegrationManager()
	if manager == nil {
		return nil, fmt.Errorf("integration manager not initialized")
	}

	secretsIntegration := manager.GetSecretsIntegration()
	if secretsIntegration == nil {
		return nil, fmt.Errorf("secrets integration not available")
	}

	return secretsIntegration, nil
}

// loadRecipients loads recipients from the global config
func loadRecipients(secretsIntegration spookyinterfaces.SecretsIntegration) ([]string, error) {
	recipients, err := secretsIntegration.LoadRecipients(context.Background(), "~/.config/spooky/recipients.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to load recipients: %w", err)
	}

	if len(recipients) == 0 {
		return nil, fmt.Errorf("no recipients found in ~/.config/spooky/recipients.txt")
	}

	return recipients, nil
}

// getIdentityPath returns the path to the identity file for decryption
func getIdentityPath() string {
	// Default identity path
	identityPath := "~/.config/spooky/identities"

	// Check if identity directory exists
	if expanded, err := expandPath(identityPath); err == nil {
		identityPath = expanded
	}

	// For now, return the first identity file found
	// In a real implementation, this would be more sophisticated
	return identityPath
}

// handleEncryptionOperation performs encryption for a specific integration type
func handleEncryptionOperation(projectPath string, dryRun bool, operationType string, encryptFunc func(context.Context, string, spookyinterfaces.SecretsIntegration, []string, bool) error) error {
	secretsIntegration, err := setupEncryption()
	if err != nil {
		return err
	}

	fmt.Printf("%s encryption for %s (dry-run: %t)\n", operationType, projectPath, dryRun)

	recipients, err := loadRecipients(secretsIntegration)
	if err != nil {
		return err
	}

	fmt.Printf("Loaded %d recipients for encryption\n", len(recipients))

	if err := encryptFunc(context.Background(), projectPath, secretsIntegration, recipients, dryRun); err != nil {
		return fmt.Errorf("failed to encrypt %s: %w", operationType, err)
	}

	if dryRun {
		fmt.Println("Dry run completed - no changes made")
	} else {
		fmt.Printf("%s encryption completed successfully\n", operationType)
	}

	return nil
}

// integrationsCmd represents the integrations command
var integrationsCmd = &cobra.Command{
	Use:   "integrations",
	Short: "Manage system integrations",
	Long: `Manage and validate system integrations.

This command provides tools to check the status of all system integrations
including facts, actions, variables, templates, machines, secrets, and configuration.`,
}

// integrationsListCmd represents the integrations list command
var integrationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available integrations",
	Long: `List all available integrations and their current status.

This command shows which integrations are available and working correctly.`,
	RunE: func(_ *cobra.Command, args []string) error {
		// Initialize dependencies if not already done
		if integrationManager == nil {
			if err := InitializeIntegrationsDependencies(); err != nil {
				return err
			}
		}

		// Create integrations command
		integrationsCmd := spookycli.NewIntegrationsCommand(integrationManager)

		// Run list command
		ctx := context.Background()
		return integrationsCmd.List(ctx, args)
	},
}

// integrationsValidateCmd represents the integrations validate command
var integrationsValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate all integrations",
	Long: `Validate that all integrations are working correctly.

This command performs comprehensive validation of all system integrations
and reports any issues found.`,
	RunE: func(_ *cobra.Command, args []string) error {
		// Initialize dependencies if not already done
		if integrationManager == nil {
			if err := InitializeIntegrationsDependencies(); err != nil {
				return err
			}
		}

		// Create integrations command
		integrationsCmd := spookycli.NewIntegrationsCommand(integrationManager)

		// Run validate command
		ctx := context.Background()
		return integrationsCmd.Validate(ctx, args)
	},
}

func init() {
	// Add subcommands to integrations command
	integrationsCmd.AddCommand(integrationsListCmd)
	integrationsCmd.AddCommand(integrationsValidateCmd)

	// Add integrations command to root
	RootCmd.AddCommand(integrationsCmd)
}
