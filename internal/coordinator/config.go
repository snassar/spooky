package coordinator

import (
	"fmt"
	"time"

	spookyconfig "spooky/internal/config"
	spookyconfigtypes "spooky/internal/types/config"
	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
)

// CoordinatorConfigIntegration implements config system integration
type CoordinatorConfigIntegration struct {
	configManager spookyconfig.ConfigManager
	logger        spookylogging.Logger
}

// NewCoordinatorConfigIntegration creates a new config integration
func NewCoordinatorConfigIntegration(configManager spookyconfig.ConfigManager, logger spookylogging.Logger) *CoordinatorConfigIntegration {
	return &CoordinatorConfigIntegration{
		configManager: configManager,
		logger:        logger,
	}
}

// LoadConfig loads configuration for the project
func (ci *CoordinatorConfigIntegration) LoadConfig(projectPath string) (*spookyinterfaces.ConfigContext, error) {
	// Load global configuration
	if err := ci.configManager.LoadConfig(); err != nil {
		return nil, fmt.Errorf("failed to load global configuration: %w", err)
	}

	// Load project configuration if provided
	if projectPath != "" {
		if err := ci.configManager.LoadProjectConfig(projectPath); err != nil {
			return nil, fmt.Errorf("failed to load project configuration: %w", err)
		}
	}

	// Get the loaded configuration
	config := ci.configManager.GetConfig()
	if config == nil {
		return nil, fmt.Errorf("no configuration loaded")
	}

	// Create config context
	configContext := &spookyinterfaces.ConfigContext{
		BaseContext: spookyinterfaces.BaseContext{
			ProjectPath: projectPath,
			Timestamp:   time.Now(),
		},
		GlobalConfig:  config.GlobalConfig,
		ProjectConfig: config.ProjectConfig,
		Environment:   config.Environment,
	}

	ci.logger.Info("Loaded configuration",
		spookylogging.String("project", projectPath),
		spookylogging.String("source", string(config.Source)))

	return configContext, nil
}

// ValidateConfig validates configuration for the project
func (ci *CoordinatorConfigIntegration) ValidateConfig(projectPath string) error {
	// Validate global configuration
	if err := ci.configManager.ValidateConfig(); err != nil {
		return fmt.Errorf("global configuration validation failed: %w", err)
	}

	// Validate project configuration if provided
	if projectPath != "" {
		if err := ci.configManager.LoadProjectConfig(projectPath); err != nil {
			return fmt.Errorf("project configuration validation failed: %w", err)
		}
	}

	ci.logger.Info("Validated configuration",
		spookylogging.String("project", projectPath))

	return nil
}

// GetConfigValue gets a configuration value by path
func (ci *CoordinatorConfigIntegration) GetConfigValue(path string) (interface{}, error) {
	value, _, err := ci.configManager.GetValue(path)
	return value, err
}

// SetConfigValue sets a configuration value by path
func (ci *CoordinatorConfigIntegration) SetConfigValue(path string, value interface{}) error {
	return ci.configManager.SetValue(path, value, spookyconfigtypes.SourceCLI)
}

// GetConfigString gets a string configuration value
func (ci *CoordinatorConfigIntegration) GetConfigString(path string) (string, error) {
	value, _, err := ci.configManager.GetString(path)
	return value, err
}

// GetConfigInt gets an integer configuration value
func (ci *CoordinatorConfigIntegration) GetConfigInt(path string) (int, error) {
	value, _, err := ci.configManager.GetInt(path)
	return value, err
}

// GetConfigBool gets a boolean configuration value
func (ci *CoordinatorConfigIntegration) GetConfigBool(path string) (bool, error) {
	value, _, err := ci.configManager.GetBool(path)
	return value, err
}

// GetConfigStringSlice gets a string slice configuration value
func (ci *CoordinatorConfigIntegration) GetConfigStringSlice(path string) ([]string, error) {
	value, _, err := ci.configManager.GetStringSlice(path)
	return value, err
}

// ApplyCLIFlags applies CLI flags to configuration
func (ci *CoordinatorConfigIntegration) ApplyCLIFlags(flags map[string]interface{}) error {
	return ci.configManager.ApplyCLIFlags(flags)
}

// GetConfigurationPath gets the configuration file path
func (ci *CoordinatorConfigIntegration) GetConfigurationPath() string {
	return ci.configManager.GetConfigurationPath()
}

// CreateDefaultConfig creates default configuration
func (ci *CoordinatorConfigIntegration) CreateDefaultConfig() error {
	return ci.configManager.CreateDefaultConfig()
}

// ReloadConfig reloads the configuration
func (ci *CoordinatorConfigIntegration) ReloadConfig() error {
	return ci.configManager.ReloadConfig()
}

// GetProjectConfig gets the project configuration
func (ci *CoordinatorConfigIntegration) GetProjectConfig() interface{} {
	return ci.configManager.GetProjectConfig()
}
