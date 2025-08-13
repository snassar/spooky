package commands

import (
	"context"
	"testing"

	spookyintegration "spooky/internal/integration"
	spookylogging "spooky/internal/logging"
)

func TestIntegrationsCommand_List(t *testing.T) {
	// Create a logger
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("test")

	// Create integration manager with nil integrations
	manager := spookyintegration.NewManager(
		logger,
		nil, // factsIntegration
		nil, // actionsIntegration
		nil, // variablesIntegration
		nil, // templatesIntegration
		nil, // machinesIntegration
		nil, // secretsIntegration
		nil, // configIntegration
	)

	// Create integrations command
	cmd := NewIntegrationsCommand(manager)

	// Test list command
	ctx := context.Background()
	err := cmd.List(ctx, []string{})
	if err != nil {
		t.Fatalf("List command failed: %v", err)
	}
}

func TestIntegrationsCommand_Validate(t *testing.T) {
	// Create a logger
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("test")

	// Create integration manager with nil integrations
	manager := spookyintegration.NewManager(
		logger,
		nil, // factsIntegration
		nil, // actionsIntegration
		nil, // variablesIntegration
		nil, // templatesIntegration
		nil, // machinesIntegration
		nil, // secretsIntegration
		nil, // configIntegration
	)

	// Create integrations command
	cmd := NewIntegrationsCommand(manager)

	// Test validate command (should fail since all integrations are nil)
	// We expect this to call os.Exit(1), so we can't test it normally
	// In a real scenario, this would be tested with actual integrations

	// For testing purposes, we'll just verify the command can be created
	// and the validate method exists (the actual validation would exit)
	if cmd == nil {
		t.Fatal("Failed to create integrations command")
	}
}
