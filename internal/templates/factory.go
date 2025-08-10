package templates

import (
	spookyinterfaces "spooky/internal/interfaces"
	spookytemplatesengine "spooky/internal/templates/engine"
	spookytemplatesfunctions "spooky/internal/templates/functions"
	spookytemplatessecrets "spooky/internal/templates/secrets"
	spookytemplatesvalidation "spooky/internal/templates/validation"
	spookytypes "spooky/internal/types"
)

// NewTemplateManager creates a new template manager with all dependencies
func NewTemplateManager(config *spookytypes.TemplateConfig, logger spookyinterfaces.Logger) (spookyinterfaces.TemplateManager, error) {
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
func NewDefaultTemplateManager(logger spookyinterfaces.Logger) (spookyinterfaces.TemplateManager, error) {
	config := &spookytypes.TemplateConfig{
		DefaultTimeout:  30,          // 30 seconds default timeout
		MaxTemplateSize: 1024 * 1024, // 1MB max template size
		EngineConfig: &spookytypes.EngineConfig{
			StrictMode:       false,
			MaxExecutionTime: 30,
		},
		FunctionsConfig: &spookytypes.FunctionsConfig{
			BuiltinFunctions: true,
			MaxTemplateSize:  1024 * 1024,
			MaxNestingDepth:  10,
			FunctionTimeout:  5,
		},
		ValidationConfig: &spookytypes.TemplateValidationConfig{
			StrictValidation: false,
		},
		SecretsConfig: &spookytypes.TemplateSecretsConfig{
			Enabled:             false,
			EncryptionAlgorithm: "age",
		},
	}

	return NewTemplateManager(config, logger)
}
