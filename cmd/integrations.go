// Package cmd provides command implementations for spooky CLI.
package cmd

import (
	"context"

	spookycli "spooky/internal/cli/commands"
	spookyintegration "spooky/internal/integration"
	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookytypeslogging "spooky/internal/types/logging"

	"github.com/spf13/cobra"
)

// Global instances for integration dependency injection
var (
	integrationManager spookyinterfaces.IntegrationManager
	integrationsLogger spookytypeslogging.Logger
)

// InitializeIntegrationsDependencies initializes integration-related dependencies
func InitializeIntegrationsDependencies() error {
	// Create log manager for integrations component
	logManager := spookylogging.NewLogManager()
	integrationsLogger = logManager.GetLogger("integrations")

	// Create integration factory and manager
	factory := spookyintegration.NewFactory(integrationsLogger)
	integrationManager = factory.CreateIntegrationManager()

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
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize dependencies if not already done
		if integrationManager == nil {
			if err := InitializeIntegrationsDependencies(); err != nil {
				return err
			}
		}

		// Create integrations command
		integrationsCmd := spookycli.NewIntegrationsCommand(integrationManager)

		// Execute list command
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
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize dependencies if not already done
		if integrationManager == nil {
			if err := InitializeIntegrationsDependencies(); err != nil {
				return err
			}
		}

		// Create integrations command
		integrationsCmd := spookycli.NewIntegrationsCommand(integrationManager)

		// Execute validate command
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
