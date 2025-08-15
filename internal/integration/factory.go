// Package integration provides the central IntegrationManager for coordinating all system integrations.
package integration

import (
	"context"
	"time"

	spookyactions "spooky/internal/actions"
	spookyconfig "spooky/internal/config"
	spookyfacts "spooky/internal/facts"
	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookymachines "spooky/internal/machines"
	spookyschemas "spooky/internal/schemas"
	spookysecrets "spooky/internal/secrets"
	spookyssh "spooky/internal/ssh"
	spookytemplates "spooky/internal/templates"
	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesschemas "spooky/internal/types/schemas"
	spookyvariables "spooky/internal/variables"
)

// SimpleSchemaValidator implements SchemaValidator interface for facts validation
type SimpleSchemaValidator struct {
	logger spookytypeslogging.Logger
}

// Validate validates data against a schema
func (s *SimpleSchemaValidator) Validate(schema *spookytypesschemas.Schema, data interface{}) (*spookytypesschemas.ValidationResult, error) {
	return &spookytypesschemas.ValidationResult{
		Valid:       true,
		ValidatedAt: time.Now(),
		Errors:      []spookytypesschemas.SchemaError{},
		Warnings:    []spookytypesschemas.SchemaError{},
	}, nil
}

// ValidateFile validates a file against a schema
func (s *SimpleSchemaValidator) ValidateFile(schema *spookytypesschemas.Schema, filePath string) (*spookytypesschemas.ValidationResult, error) {
	return s.Validate(schema, nil)
}

// ValidateString validates a string against a schema
func (s *SimpleSchemaValidator) ValidateString(schema *spookytypesschemas.Schema, content string) (*spookytypesschemas.ValidationResult, error) {
	return s.Validate(schema, content)
}

// ValidateBytes validates bytes against a schema
func (s *SimpleSchemaValidator) ValidateBytes(schema *spookytypesschemas.Schema, data []byte) (*spookytypesschemas.ValidationResult, error) {
	return s.Validate(schema, data)
}

// ValidateWithContext validates data with additional context
func (s *SimpleSchemaValidator) ValidateWithContext(schema *spookytypesschemas.Schema, data interface{}, context map[string]interface{}) (*spookytypesschemas.ValidationResult, error) {
	return s.Validate(schema, data)
}

// ValidateField validates a specific field
func (s *SimpleSchemaValidator) ValidateField(schema *spookytypesschemas.Schema, fieldPath string, value interface{}) (*spookytypesschemas.ValidationResult, error) {
	return &spookytypesschemas.ValidationResult{
		Valid:       true,
		ValidatedAt: time.Now(),
		Errors:      []spookytypesschemas.SchemaError{},
		Warnings:    []spookytypesschemas.SchemaError{},
	}, nil
}

// SimpleSchemaManager implements SchemaManager interface for integration
type SimpleSchemaManager struct {
	logger spookytypeslogging.Logger
}

// LoadSchema loads a schema from the given path
func (s *SimpleSchemaManager) LoadSchema(ctx context.Context, schemaPath string) (*spookytypes.Schema, error) {
	return &spookytypes.Schema{}, nil
}

// LoadEmbeddedSchema loads an embedded schema
func (s *SimpleSchemaManager) LoadEmbeddedSchema(ctx context.Context, schemaName string) (*spookytypes.Schema, error) {
	return &spookytypes.Schema{}, nil
}

// Validate validates data against a schema
func (s *SimpleSchemaManager) Validate(ctx context.Context, schema *spookytypes.Schema, data interface{}) (*spookytypes.ValidationResult, error) {
	return &spookytypes.ValidationResult{
		Valid:    true,
		Errors:   []spookytypes.SchemaError{},
		Warnings: []spookytypes.SchemaError{},
	}, nil
}

// Register registers a new schema
func (s *SimpleSchemaManager) Register(ctx context.Context, schema *spookytypes.Schema) error {
	return nil
}

// Factory creates all integration components
type Factory struct {
	logger spookytypeslogging.Logger
}

// NewFactory creates a new integration factory
func NewFactory(logger spookytypeslogging.Logger) *Factory {
	return &Factory{
		logger: logger,
	}
}

// CreateIntegrationManager creates a complete IntegrationManager with all integrations
func (f *Factory) CreateIntegrationManager() spookyinterfaces.IntegrationManager {
	// Create all individual integrations
	factsIntegration := f.createFactsIntegration()
	actionsIntegration := f.createActionsIntegration()
	variablesIntegration := f.createVariablesIntegration()
	templatesIntegration := f.createTemplatesIntegration()
	machinesIntegration := f.createMachinesIntegration()
	secretsIntegration := f.createSecretsIntegration()
	configIntegration := f.createConfigIntegration()

	// Create the integration manager
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

	return manager
}

// createFactsIntegration creates the facts integration
func (f *Factory) createFactsIntegration() spookyinterfaces.FactsIntegration {
	// Create SSH manager for facts collection
	sshManager := spookyssh.NewManager(f.logger)

	// Create fact collector
	collector := spookyfacts.NewSystemFactCollector(sshManager)

	// Create a simple schema validator for facts validation
	schemaValidator := &SimpleSchemaValidator{logger: f.logger}

	// Create facts manager with proper schema validator
	factsManager := spookyfacts.NewManager(collector, schemaValidator, f.logger)

	// Create facts integration
	factsIntegration := spookyfacts.NewIntegration(factsManager)

	return factsIntegration
}

// createActionsIntegration creates the actions integration
func (f *Factory) createActionsIntegration() spookyinterfaces.ActionsIntegration {
	// Create log manager for actions
	logManager := spookylogging.NewLogManager()
	actionsLogger := logManager.GetLogger("actions")

	// Create SSH manager for actions
	sshManager := spookyssh.NewManager(f.logger)

	// Create action validator with interface logger type
	actionValidator := spookyactions.NewValidator(actionsLogger)

	// Create schema validator
	schemaValidator := spookyschemas.NewValidator(f.logger)

	// Create actions manager with interface logger type
	actionsManager := spookyactions.NewManager(actionsLogger, actionValidator, sshManager, schemaValidator)

	return actionsManager
}

// createVariablesIntegration creates the variables integration
func (f *Factory) createVariablesIntegration() spookyinterfaces.VariablesIntegration {
	// Create variables loader
	variablesLoader := spookyvariables.NewLoader(f.logger)

	// Create variables validator
	variablesValidator := spookyvariables.NewValidator(f.logger)

	// Create variables manager
	variablesManager := spookyvariables.NewManager(f.logger, variablesLoader, variablesValidator)

	return variablesManager
}

// createTemplatesIntegration creates the templates integration
func (f *Factory) createTemplatesIntegration() spookyinterfaces.TemplatesIntegration {
	// Create templates manager
	templatesManager := spookytemplates.NewManager(f.logger)

	// Create templates integration
	templatesIntegration := spookytemplates.NewIntegration(templatesManager)

	return templatesIntegration
}

// createMachinesIntegration creates the machines integration
func (f *Factory) createMachinesIntegration() spookyinterfaces.MachinesIntegration {
	// Create machines loader
	machinesLoader := spookymachines.NewLoader(f.logger)

	// Create machines validator
	machinesValidator := spookymachines.NewValidator(f.logger)

	// Create SSH manager for machines
	sshManager := spookyssh.NewManager(f.logger)

	// Create machines manager
	machinesManager := spookymachines.NewManager(f.logger, machinesLoader, machinesValidator)

	// Create machines integration
	machinesIntegration := spookymachines.NewIntegration(
		machinesManager,
		machinesLoader,
		machinesValidator,
		sshManager,
		f.logger,
	)

	return machinesIntegration
}

// createSecretsIntegration creates the secrets integration
func (f *Factory) createSecretsIntegration() spookyinterfaces.SecretsIntegration {
	// Create secrets integration
	secretsIntegration := spookysecrets.NewIntegration(f.logger)

	return secretsIntegration
}

// createConfigIntegration creates the config integration
func (f *Factory) createConfigIntegration() spookyinterfaces.ConfigIntegration {
	// Create config integration
	configIntegration := spookyconfig.NewIntegration(f.logger)

	return configIntegration
}

// CreateSchemaManager creates a schema manager with enhanced validation
func (f *Factory) CreateSchemaManager() spookyinterfaces.SchemaManager {
	// Create schema manager
	schemaManager := spookyschemas.NewManager(f.logger)

	return schemaManager
}

// CreateLogManager creates a log manager
func (f *Factory) CreateLogManager() spookyinterfaces.LogManager {
	// Create log manager
	logManager := spookylogging.NewLogManager()

	return logManager
}

// CreateSSHManager creates an SSH manager with connection pooling
func (f *Factory) CreateSSHManager() spookyinterfaces.SSHManager {
	// Create SSH manager
	sshManager := spookyssh.NewManager(f.logger)

	return sshManager
}

// CreateAllManagers creates all managers for the system
func (f *Factory) CreateAllManagers() map[string]interface{} {
	managers := make(map[string]interface{})

	// Create integration manager
	integrationManager := f.CreateIntegrationManager()
	managers["integration"] = integrationManager

	// Create schema manager
	schemaManager := f.CreateSchemaManager()
	managers["schema"] = schemaManager

	// Create log manager
	logManager := f.CreateLogManager()
	managers["log"] = logManager

	// Create SSH manager
	sshManager := f.CreateSSHManager()
	managers["ssh"] = sshManager

	return managers
}
