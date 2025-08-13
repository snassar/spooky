package integration

import (
	"context"
	"testing"

	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
)

func TestNewManager(t *testing.T) {
	// Create a logger
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("test")

	// Create manager with nil integrations (for testing)
	manager := NewManager(
		logger,
		nil, // factsIntegration
		nil, // actionsIntegration
		nil, // variablesIntegration
		nil, // templatesIntegration
		nil, // machinesIntegration
		nil, // secretsIntegration
		nil, // configIntegration
	)

	if manager == nil {
		t.Fatal("NewManager returned nil")
	}

	// Verify it implements the interface
	var _ spookyinterfaces.IntegrationManager = manager
}

func TestManager_GetIntegrations(t *testing.T) {
	// Create a logger
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("test")

	// Create manager with nil integrations
	manager := NewManager(
		logger,
		nil, // factsIntegration
		nil, // actionsIntegration
		nil, // variablesIntegration
		nil, // templatesIntegration
		nil, // machinesIntegration
		nil, // secretsIntegration
		nil, // configIntegration
	)

	// Test getting integrations
	if manager.GetFactsIntegration() != nil {
		t.Error("Expected nil facts integration")
	}

	if manager.GetActionsIntegration() != nil {
		t.Error("Expected nil actions integration")
	}

	if manager.GetVariablesIntegration() != nil {
		t.Error("Expected nil variables integration")
	}

	if manager.GetTemplatesIntegration() != nil {
		t.Error("Expected nil templates integration")
	}

	if manager.GetMachinesIntegration() != nil {
		t.Error("Expected nil machines integration")
	}

	if manager.GetSecretsIntegration() != nil {
		t.Error("Expected nil secrets integration")
	}

	if manager.GetConfigIntegration() != nil {
		t.Error("Expected nil config integration")
	}
}

func TestManager_InternalMethods(t *testing.T) {
	// Create a logger
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("test")

	// Create manager with nil integrations
	manager := NewManager(
		logger,
		nil, // factsIntegration
		nil, // actionsIntegration
		nil, // variablesIntegration
		nil, // templatesIntegration
		nil, // machinesIntegration
		nil, // secretsIntegration
		nil, // configIntegration
	)

	// Cast to concrete type to test internal methods
	managerImpl, ok := manager.(*Manager)
	if !ok {
		t.Fatal("Failed to cast to concrete Manager type")
	}

	// Test system health validation
	ctx := context.Background()
	result, err := managerImpl.ValidateSystemHealth(ctx)
	if err != nil {
		t.Fatalf("ValidateSystemHealth failed: %v", err)
	}

	// Since all integrations are nil, health should be invalid
	if result.Valid {
		t.Error("Expected system health to be invalid when all integrations are nil")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected validation errors when all integrations are nil")
	}

	// Test getting health status
	status := managerImpl.GetHealthStatus()

	// Check that all integrations are marked as unhealthy
	expectedUnhealthy := []string{"facts", "actions", "variables", "templates", "machines", "secrets", "config"}
	for _, integration := range expectedUnhealthy {
		if status[integration] {
			t.Errorf("Expected integration %s to be unhealthy", integration)
		}
	}

	// Test updating health status
	managerImpl.UpdateHealthStatus("facts", true)

	// Check that the status was updated
	status = managerImpl.GetHealthStatus()
	if !status["facts"] {
		t.Error("Expected facts integration to be marked as healthy")
	}

	// Test coordinated operation with unhealthy system
	err = managerImpl.CoordinatedOperation(ctx, func() error {
		return nil
	})

	// Should fail because system health is invalid
	if err == nil {
		t.Error("Expected coordinated operation to fail with unhealthy system")
	}

	// Mark all integrations as healthy
	managerImpl.UpdateHealthStatus("facts", true)
	managerImpl.UpdateHealthStatus("actions", true)
	managerImpl.UpdateHealthStatus("variables", true)
	managerImpl.UpdateHealthStatus("templates", true)
	managerImpl.UpdateHealthStatus("machines", true)
	managerImpl.UpdateHealthStatus("secrets", true)
	managerImpl.UpdateHealthStatus("config", true)

	// Test coordinated operation with healthy system
	err = managerImpl.CoordinatedOperation(ctx, func() error {
		return nil
	})

	if err != nil {
		t.Errorf("Expected coordinated operation to succeed with healthy system: %v", err)
	}
}
