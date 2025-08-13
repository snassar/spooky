package commands

import (
	"context"
	"fmt"
	"os"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
)

// IntegrationsCommand handles integration management
type IntegrationsCommand struct {
	integrationManager spookyinterfaces.IntegrationManager
}

// NewIntegrationsCommand creates a new integrations command
func NewIntegrationsCommand(manager spookyinterfaces.IntegrationManager) *IntegrationsCommand {
	return &IntegrationsCommand{
		integrationManager: manager,
	}
}

// List lists all available integrations and their status
func (c *IntegrationsCommand) List(ctx context.Context, args []string) error {
	// Get health status from the integration manager
	managerImpl, ok := c.integrationManager.(interface {
		GetHealthStatus() map[string]bool
	})
	if !ok {
		return fmt.Errorf("integration manager does not support health status")
	}

	status := managerImpl.GetHealthStatus()

	fmt.Println("Available Integrations:")
	fmt.Println("=======================")

	integrations := []string{
		"facts", "actions", "variables", "templates",
		"machines", "secrets", "config",
	}

	for _, integration := range integrations {
		healthy := status[integration]
		statusStr := "❌ unavailable"
		if healthy {
			statusStr = "✅ available"
		}
		fmt.Printf("%-12s %s\n", integration, statusStr)
	}

	return nil
}

// Validate validates all integrations are working
func (c *IntegrationsCommand) Validate(ctx context.Context, args []string) error {
	// Get health status from the integration manager
	managerImpl, ok := c.integrationManager.(interface {
		GetHealthStatus() map[string]bool
		ValidateSystemHealth(context.Context) (*spookytypes.ValidationResult, error)
	})
	if !ok {
		return fmt.Errorf("integration manager does not support validation")
	}

	// Validate system health
	result, err := managerImpl.ValidateSystemHealth(ctx)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if result.Valid {
		fmt.Println("✅ All integrations are working correctly")
		return nil
	}

	fmt.Println("❌ Integration validation failed:")
	for _, error := range result.Errors {
		fmt.Printf("  - %s\n", error.Message)
	}

	// Exit with error code if validation fails
	os.Exit(1)
	return nil
}
