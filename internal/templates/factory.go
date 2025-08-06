package templates

import (
	spookylogging "spooky/internal/logging"
	spookytemplatesengine "spooky/internal/templates/engine"
	spookytemplatesfunctions "spooky/internal/templates/functions"
	spookytemplatessecrets "spooky/internal/templates/secrets"
	spookytemplatestypes "spooky/internal/templates/types"
	spookytemplatesvalidation "spooky/internal/templates/validation"
)

// NewTemplateManager creates a new template manager with all dependencies
func NewTemplateManager(config *spookytemplatestypes.Config, logger spookylogging.Logger) (TemplateManager, error) {
	// Create sub-managers
	engineManager := spookytemplatesengine.NewManager(config.EngineConfig, logger)
	functionsManager := spookytemplatesfunctions.NewManager(config.FunctionsConfig, logger)
	validationManager := spookytemplatesvalidation.NewManager(config.ValidationConfig, logger)
	secretsManager := spookytemplatessecrets.NewManager(config.SecretsConfig, logger)

	// Create main manager
	manager := NewManager(
		config,
		engineManager,
		functionsManager,
		validationManager,
		secretsManager,
		logger,
	)

	return manager, nil
}

// NewDefaultTemplateManager creates a template manager with default configuration
func NewDefaultTemplateManager(logger spookylogging.Logger) (TemplateManager, error) {
	config := &spookytemplatestypes.Config{
		DefaultTimeout:  30,          // 30 seconds default timeout
		MaxTemplateSize: 1024 * 1024, // 1MB max template size
		EngineConfig: &spookytemplatestypes.EngineConfig{
			StrictMode:       false,
			MaxExecutionTime: 30,
		},
		FunctionsConfig: &spookytemplatestypes.FunctionsConfig{
			BuiltinFunctions: true,
			MaxTemplateSize:  1024 * 1024,
			MaxNestingDepth:  10,
			FunctionTimeout:  5,
		},
		ValidationConfig: &spookytemplatestypes.ValidationConfig{
			StrictValidation: false,
		},
		SecretsConfig: &spookytemplatestypes.SecretsConfig{
			Enabled:             false,
			EncryptionAlgorithm: "age",
		},
	}

	return NewTemplateManager(config, logger)
}
