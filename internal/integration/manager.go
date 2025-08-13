// Package integration provides the central IntegrationManager for coordinating all system integrations.
package integration

import (
	"context"
	"fmt"
	"sync"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesschemas "spooky/internal/types/schemas"
)

// Manager implements the IntegrationManager interface and provides system-wide coordination
type Manager struct {
	logger spookytypeslogging.Logger

	// Integration components
	factsIntegration     spookyinterfaces.FactsIntegration
	actionsIntegration   spookyinterfaces.ActionsIntegration
	variablesIntegration spookyinterfaces.VariablesIntegration
	templatesIntegration spookyinterfaces.TemplatesIntegration
	machinesIntegration  spookyinterfaces.MachinesIntegration
	secretsIntegration   spookyinterfaces.SecretsIntegration
	configIntegration    spookyinterfaces.ConfigIntegration

	// Thread safety
	mu sync.RWMutex

	// Health tracking
	healthStatus map[string]bool
	healthMu     sync.RWMutex
}

// NewManager creates a new IntegrationManager with all required integrations
func NewManager(
	logger spookytypeslogging.Logger,
	factsIntegration spookyinterfaces.FactsIntegration,
	actionsIntegration spookyinterfaces.ActionsIntegration,
	variablesIntegration spookyinterfaces.VariablesIntegration,
	templatesIntegration spookyinterfaces.TemplatesIntegration,
	machinesIntegration spookyinterfaces.MachinesIntegration,
	secretsIntegration spookyinterfaces.SecretsIntegration,
	configIntegration spookyinterfaces.ConfigIntegration,
) spookyinterfaces.IntegrationManager {
	manager := &Manager{
		logger:               logger,
		factsIntegration:     factsIntegration,
		actionsIntegration:   actionsIntegration,
		variablesIntegration: variablesIntegration,
		templatesIntegration: templatesIntegration,
		machinesIntegration:  machinesIntegration,
		secretsIntegration:   secretsIntegration,
		configIntegration:    configIntegration,
		healthStatus:         make(map[string]bool),
	}

	// Initialize health status
	manager.initializeHealthStatus()

	return manager
}

// GetFactsIntegration returns the facts integration
func (m *Manager) GetFactsIntegration() spookyinterfaces.FactsIntegration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.factsIntegration
}

// GetActionsIntegration returns the actions integration
func (m *Manager) GetActionsIntegration() spookyinterfaces.ActionsIntegration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.actionsIntegration
}

// GetVariablesIntegration returns the variables integration
func (m *Manager) GetVariablesIntegration() spookyinterfaces.VariablesIntegration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.variablesIntegration
}

// GetTemplatesIntegration returns the templates integration
func (m *Manager) GetTemplatesIntegration() spookyinterfaces.TemplatesIntegration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.templatesIntegration
}

// GetMachinesIntegration returns the machines integration
func (m *Manager) GetMachinesIntegration() spookyinterfaces.MachinesIntegration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.machinesIntegration
}

// GetSecretsIntegration returns the secrets integration
func (m *Manager) GetSecretsIntegration() spookyinterfaces.SecretsIntegration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.secretsIntegration
}

// GetConfigIntegration returns the config integration
func (m *Manager) GetConfigIntegration() spookyinterfaces.ConfigIntegration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.configIntegration
}

// InitializeHealthStatus initializes the health status for all integrations
func (m *Manager) initializeHealthStatus() {
	m.healthMu.Lock()
	defer m.healthMu.Unlock()

	// Set initial health status based on integration availability
	m.healthStatus["facts"] = m.factsIntegration != nil
	m.healthStatus["actions"] = m.actionsIntegration != nil
	m.healthStatus["variables"] = m.variablesIntegration != nil
	m.healthStatus["templates"] = m.templatesIntegration != nil
	m.healthStatus["machines"] = m.machinesIntegration != nil
	m.healthStatus["secrets"] = m.secretsIntegration != nil
	m.healthStatus["config"] = m.configIntegration != nil

	m.logger.Info("IntegrationManager initialized", map[string]interface{}{
		"facts_available":     m.healthStatus["facts"],
		"actions_available":   m.healthStatus["actions"],
		"variables_available": m.healthStatus["variables"],
		"templates_available": m.healthStatus["templates"],
		"machines_available":  m.healthStatus["machines"],
		"secrets_available":   m.healthStatus["secrets"],
		"config_available":    m.healthStatus["config"],
	})
}

// ValidateSystemHealth validates the health of all integrations
func (m *Manager) ValidateSystemHealth(ctx context.Context) (*spookytypesschemas.ValidationResult, error) {
	m.healthMu.RLock()
	defer m.healthMu.RUnlock()

	var errors []spookytypes.SchemaError
	var warnings []spookytypes.SchemaError

	// Check each integration's health
	for integration, healthy := range m.healthStatus {
		if !healthy {
			errors = append(errors, spookytypesschemas.SchemaError{
				Message: fmt.Sprintf("integration %s is not available", integration),
			})
		}
	}

	// Additional health checks can be added here
	// For example, testing connectivity, validating configurations, etc.

	valid := len(errors) == 0

	m.logger.Info("System health validation completed", map[string]interface{}{
		"valid":    valid,
		"errors":   len(errors),
		"warnings": len(warnings),
	})

	return &spookytypesschemas.ValidationResult{
		Valid:    valid,
		Errors:   errors,
		Warnings: warnings,
	}, nil
}

// GetHealthStatus returns the current health status of all integrations
func (m *Manager) GetHealthStatus() map[string]bool {
	m.healthMu.RLock()
	defer m.healthMu.RUnlock()

	// Return a copy to prevent external modification
	status := make(map[string]bool)
	for k, v := range m.healthStatus {
		status[k] = v
	}

	return status
}

// UpdateHealthStatus updates the health status of a specific integration
func (m *Manager) UpdateHealthStatus(integration string, healthy bool) {
	m.healthMu.Lock()
	defer m.healthMu.Unlock()

	m.healthStatus[integration] = healthy

	m.logger.Info("Integration health status updated", map[string]interface{}{
		"integration": integration,
		"healthy":     healthy,
	})
}

// CoordinatedOperation performs a coordinated operation across multiple integrations
func (m *Manager) CoordinatedOperation(ctx context.Context, operation func() error) error {
	m.logger.Info("Starting coordinated operation")

	// Validate system health before operation
	healthResult, err := m.ValidateSystemHealth(ctx)
	if err != nil {
		return fmt.Errorf("failed to validate system health: %w", err)
	}

	if !healthResult.Valid {
		return fmt.Errorf("system health validation failed: %v", healthResult.Errors)
	}

	// Perform the coordinated operation
	if err := operation(); err != nil {
		m.logger.Error("Coordinated operation failed", err, map[string]interface{}{
			"operation": "coordinated",
		})
		return fmt.Errorf("coordinated operation failed: %w", err)
	}

	m.logger.Info("Coordinated operation completed successfully")
	return nil
}
