package config

import (
	spookyconfigenvironment "spooky/internal/config/environment"
	spookyconfigloading "spooky/internal/config/loading"
	spookyconfigtypes "spooky/internal/types/config"
	spookyconfigvalidation "spooky/internal/config/validation"
	spookylogging "spooky/internal/logging"
)

// NewManagerWithConfig creates a new config manager with custom configuration
func NewManagerWithConfig(
	config *spookyconfigtypes.Config,
	logger spookylogging.Logger,
) *Manager {
	// Create sub-managers with default configs if not provided
	loadingConfig := &spookyconfigtypes.LoadingConfig{
		ConfigPath: "",
		AutoReload: true,
	}
	if config != nil {
		loadingConfig = &spookyconfigtypes.LoadingConfig{
			ConfigPath: "",
			AutoReload: true,
		}
	}

	validationConfig := &spookyconfigtypes.ValidationConfig{
		StrictValidation: false,
	}

	environmentConfig := &spookyconfigtypes.EnvironmentConfig{
		ValidateVariables: true,
	}

	loadingManager := spookyconfigloading.NewManager(loadingConfig, logger)
	validationManager := spookyconfigvalidation.NewManager(validationConfig, logger)
	environmentManager := spookyconfigenvironment.NewManager(environmentConfig, logger)

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
	config *spookyconfigtypes.Config,
	loadingManager spookyconfigloading.LoadingManager,
	validationManager spookyconfigvalidation.ValidationManager,
	environmentManager spookyconfigenvironment.EnvironmentManager,
	logger spookylogging.Logger,
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
func NewDefaultManager(logger spookylogging.Logger) *Manager {
	// Create default config
	config := &spookyconfigtypes.Config{
		GlobalConfig: &spookyconfigtypes.GlobalConfig{
			LogLevel: "info",
			Quiet:    false,
			Verbose:  false,
		},
		Environment: make(map[string]interface{}),
		Source:      spookyconfigtypes.SourceDefault,
	}

	return NewManagerWithConfig(config, logger)
}
