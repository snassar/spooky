// Package integration provides the central IntegrationManager for coordinating all system integrations.
package integration

import (
	"fmt"
	spookyactions "spooky/internal/actions"
	spookyconfig "spooky/internal/config"
	spookyfacts "spooky/internal/facts"
	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookymachines "spooky/internal/machines"
	spookyschemas "spooky/internal/schemas"
	spookyssh "spooky/internal/ssh"
	spookytypeslogging "spooky/internal/types/logging"
	spookyvariables "spooky/internal/variables"
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
	// Create log manager for facts
	logManager := spookylogging.NewLogManager()
	factsLogger := logManager.GetLogger("facts")

	// Create facts components
	collector := spookyfacts.NewSystemFactCollector()
	manager := spookyfacts.NewManager(collector, nil, factsLogger)

	// Create facts integration
	integration := spookyfacts.NewIntegration(manager)

	f.logger.Info("Facts integration created successfully")
	return integration
}

// createActionsIntegration creates the actions integration
func (f *Factory) createActionsIntegration() spookyinterfaces.ActionsIntegration {
	// Create log manager for actions
	logManager := spookylogging.NewLogManager()
	actionsLoggerInterface := logManager.GetLogger("actions")

	// Cast the interface to the concrete type that actions manager expects
	actionsLoggerPtr, ok := actionsLoggerInterface.(*spookylogging.Logger)
	if !ok {
		f.logger.Error("Failed to cast logger to concrete type", fmt.Errorf("logger type assertion failed"))
		return nil
	}

	// Create SSH manager with the interface logger
	sshManager := spookyssh.NewManager(actionsLoggerInterface)

	// Create schema validator with the interface logger
	schemaValidator := spookyschemas.NewValidator(actionsLoggerInterface)

	// Load schemas from the schemas directory
	schemasDir := "internal/schemas/schemas"
	if err := schemaValidator.LoadSchemas(schemasDir); err != nil {
		f.logger.Error("Failed to load schemas", err, map[string]interface{}{
			"schemas_dir": schemasDir,
		})
		return nil
	}

	// Create actions integration with the concrete logger type
	// Dereference the pointer to get the concrete type
	integration := spookyactions.NewIntegration(*actionsLoggerPtr, nil, sshManager, schemaValidator)

	f.logger.Info("Actions integration created successfully")
	return integration
}

// createVariablesIntegration creates the variables integration
func (f *Factory) createVariablesIntegration() spookyinterfaces.VariablesIntegration {
	// Create log manager for variables
	logManager := spookylogging.NewLogManager()
	variablesLoggerInterface := logManager.GetLogger("variables")

	// Create variables loader and validator
	loader := spookyvariables.NewLoader(variablesLoggerInterface)
	validator := spookyvariables.NewValidator(variablesLoggerInterface)

	// Create variables integration
	integration := spookyvariables.NewIntegration(variablesLoggerInterface, loader, validator)

	f.logger.Info("Variables integration created successfully")
	return integration
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
	// Create log manager for machines
	logManager := spookylogging.NewLogManager()
	machinesLogger := logManager.GetLogger("machines")

	// Create machines components
	validator := spookymachines.NewValidator(machinesLogger)
	loader := spookymachines.NewLoader(machinesLogger)
	manager := spookymachines.NewManager(machinesLogger, loader, validator)

	f.logger.Info("Machines integration created successfully")
	return manager
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
	// Create log manager for config
	logManager := spookylogging.NewLogManager()
	configLoggerInterface := logManager.GetLogger("config")

	// Create config integration
	integration := spookyconfig.NewIntegration(configLoggerInterface)

	f.logger.Info("Config integration created successfully")
	return integration
}
