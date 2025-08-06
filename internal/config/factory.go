package config

import (
	"spooky/internal/config/environment"
	"spooky/internal/config/loading"
	"spooky/internal/config/types"
	"spooky/internal/config/validation"
	"spooky/internal/logging"
)

// NewManagerWithConfig creates a new config manager with custom configuration
func NewManagerWithConfig(
	config *types.Config,
	logger logging.Logger,
) *Manager {
	// Create sub-managers with default configs if not provided
	loadingConfig := &types.LoadingConfig{
		ConfigPath: "",
		AutoReload: true,
	}
	if config != nil {
		loadingConfig = &types.LoadingConfig{
			ConfigPath: "",
			AutoReload: true,
		}
	}

	validationConfig := &types.ValidationConfig{
		StrictValidation: false,
	}

	environmentConfig := &types.EnvironmentConfig{
		ValidateVariables: true,
	}

	loadingManager := loading.NewManager(loadingConfig, logger)
	validationManager := validation.NewManager(validationConfig, logger)
	environmentManager := environment.NewManager(environmentConfig, logger)

	// Create main manager
	return NewManager(
		config,
		loadingManager,
		validationManager,
		environmentManager,
		logger,
	)
}

// NewManagerWithDependencies creates a new config manager with custom dependencies
func NewManagerWithDependencies(
	config *types.Config,
	loadingManager loading.LoadingManager,
	validationManager validation.ValidationManager,
	environmentManager environment.EnvironmentManager,
	logger logging.Logger,
) *Manager {
	return NewManager(
		config,
		loadingManager,
		validationManager,
		environmentManager,
		logger,
	)
}

// NewDefaultManager creates a new config manager with default configuration
func NewDefaultManager(logger logging.Logger) *Manager {
	// Create default config
	config := &types.Config{
		GlobalConfig: &types.GlobalConfig{
			LogLevel: "info",
			Quiet:    false,
			Verbose:  false,
		},
		Environment: make(map[string]interface{}),
		Source:      types.SourceDefault,
	}

	return NewManagerWithConfig(config, logger)
}
