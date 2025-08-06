package templates

import (
	"spooky/internal/logging"
	"spooky/internal/templates/engine"
	"spooky/internal/templates/functions"
	"spooky/internal/templates/secrets"
	"spooky/internal/templates/types"
	"spooky/internal/templates/validation"
)

// NewTemplateManager creates a new template manager with all dependencies
func NewTemplateManager(config *types.Config, logger logging.Logger) (TemplateManager, error) {
	// Create sub-managers
	engineManager := engine.NewManager(config.EngineConfig, logger)
	functionsManager := functions.NewManager(config.FunctionsConfig, logger)
	validationManager := validation.NewManager(config.ValidationConfig, logger)
	secretsManager := secrets.NewManager(config.SecretsConfig, logger)

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
func NewDefaultTemplateManager(logger logging.Logger) (TemplateManager, error) {
	config := &types.Config{
		DefaultTimeout:  30,          // 30 seconds default timeout
		MaxTemplateSize: 1024 * 1024, // 1MB max template size
		EngineConfig: &types.EngineConfig{
			StrictMode:       false,
			MaxExecutionTime: 30,
		},
		FunctionsConfig: &types.FunctionsConfig{
			BuiltinFunctions: true,
			MaxTemplateSize:  1024 * 1024,
			MaxNestingDepth:  10,
			FunctionTimeout:  5,
		},
		ValidationConfig: &types.ValidationConfig{
			StrictValidation: false,
		},
		SecretsConfig: &types.SecretsConfig{
			Enabled:             false,
			EncryptionAlgorithm: "age",
		},
	}

	return NewTemplateManager(config, logger)
}
