// Package integration provides the central IntegrationManager for coordinating all system integrations.
package integration

import (
	spookyinterfaces "spooky/internal/interfaces"
	spookytypeslogging "spooky/internal/types/logging"
)

// Factory creates and configures IntegrationManager instances
type Factory struct {
	logger spookytypeslogging.Logger
}

// NewFactory creates a new IntegrationManager factory
func NewFactory(logger spookytypeslogging.Logger) *Factory {
	return &Factory{
		logger: logger,
	}
}

// CreateIntegrationManager creates a new IntegrationManager with all required integrations
func (f *Factory) CreateIntegrationManager() spookyinterfaces.IntegrationManager {
	// Create individual integrations
	factsIntegration := f.createFactsIntegration()
	actionsIntegration := f.createActionsIntegration()
	variablesIntegration := f.createVariablesIntegration()
	templatesIntegration := f.createTemplatesIntegration()
	machinesIntegration := f.createMachinesIntegration()
	secretsIntegration := f.createSecretsIntegration()
	configIntegration := f.createConfigIntegration()

	// Create the IntegrationManager
	manager := NewManager(
		f.logger,
		factsIntegration,
		actionsIntegration,
		variablesIntegration,
		templatesIntegration,
		machinesIntegration,
		secretsIntegration,
		configIntegration,
	)

	f.logger.Info("IntegrationManager created successfully", map[string]interface{}{
		"facts_available":     factsIntegration != nil,
		"actions_available":   actionsIntegration != nil,
		"variables_available": variablesIntegration != nil,
		"templates_available": templatesIntegration != nil,
		"machines_available":  machinesIntegration != nil,
		"secrets_available":   secretsIntegration != nil,
		"config_available":    configIntegration != nil,
	})

	return manager
}

// createFactsIntegration creates the facts integration
func (f *Factory) createFactsIntegration() spookyinterfaces.FactsIntegration {
	// For now, return nil to indicate facts integration is not yet implemented
	// This will be implemented when the facts system is fully integrated
	f.logger.Warn("Facts integration not yet implemented")
	return nil
}

// createActionsIntegration creates the actions integration
func (f *Factory) createActionsIntegration() spookyinterfaces.ActionsIntegration {
	// For now, return nil to indicate actions integration is not yet implemented
	// This will be implemented when the actions system is fully integrated
	f.logger.Warn("Actions integration not yet implemented")
	return nil
}

// createVariablesIntegration creates the variables integration
func (f *Factory) createVariablesIntegration() spookyinterfaces.VariablesIntegration {
	// For now, return nil to indicate variables integration is not yet implemented
	// This will be implemented when the variables system is fully integrated
	f.logger.Warn("Variables integration not yet implemented")
	return nil
}

// createTemplatesIntegration creates the templates integration
func (f *Factory) createTemplatesIntegration() spookyinterfaces.TemplatesIntegration {
	// For now, return nil to indicate templates integration is not yet implemented
	// This will be implemented when the templates system is fully integrated
	f.logger.Warn("Templates integration not yet implemented")
	return nil
}

// createMachinesIntegration creates the machines integration
func (f *Factory) createMachinesIntegration() spookyinterfaces.MachinesIntegration {
	// For now, return nil to indicate machines integration is not yet implemented
	// This will be implemented when the machines system is fully integrated
	f.logger.Warn("Machines integration not yet implemented")
	return nil
}

// createSecretsIntegration creates the secrets integration
func (f *Factory) createSecretsIntegration() spookyinterfaces.SecretsIntegration {
	// For now, return nil to indicate secrets integration is not yet implemented
	// This will be implemented when the secrets system is fully integrated
	f.logger.Warn("Secrets integration not yet implemented")
	return nil
}

// createConfigIntegration creates the config integration
func (f *Factory) createConfigIntegration() spookyinterfaces.ConfigIntegration {
	// For now, return nil to indicate config integration is not yet implemented
	// This will be implemented when the config system is fully integrated
	f.logger.Warn("Config integration not yet implemented")
	return nil
}
