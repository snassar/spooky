package config

import (
	spookyconfigenvironment "spooky/internal/config/environment"
	spookyconfigloading "spooky/internal/config/loading"
	spookyconfigvalidation "spooky/internal/config/validation"
	spookyconfigtypes "spooky/internal/types/config"
	spookytypeslogging "spooky/internal/types/logging"
)

// NewManagerWithConfig creates a new config manager with custom configuration
func NewManagerWithConfig(
	config *spookyconfigtypes.Config,
	logger spookytypeslogging.Logger,
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
// Note: This function accepts concrete types but the Manager internally implements
// the ConfigManager interface, maintaining interface-based architecture
func NewManagerWithDependencies(
	config *spookyconfigtypes.Config,
	loadingManager *spookyconfigloading.Manager,
	validationManager *spookyconfigvalidation.Manager,
	environmentManager *spookyconfigenvironment.Manager,
	logger spookytypeslogging.Logger,
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
func NewDefaultManager(logger spookytypeslogging.Logger) *Manager {
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
