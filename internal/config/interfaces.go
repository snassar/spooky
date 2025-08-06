package config

import (
	spookyconfigtypes "spooky/internal/config/types"
)

// ConfigManager defines the main interface for configuration management
type ConfigManager interface {
	// Core configuration operations
	LoadConfig() error
	ReloadConfig() error
	GetConfig() *spookyconfigtypes.Config
	ValidateConfig() error

	// Value access operations
	GetValue(path string) (interface{}, spookyconfigtypes.ConfigSource, error)
	SetValue(path string, value interface{}, source spookyconfigtypes.ConfigSource) error
	GetString(path string) (string, spookyconfigtypes.ConfigSource, error)
	GetInt(path string) (int, spookyconfigtypes.ConfigSource, error)
	GetBool(path string) (bool, spookyconfigtypes.ConfigSource, error)
	GetStringSlice(path string) ([]string, spookyconfigtypes.ConfigSource, error)

	// Project configuration
	LoadProjectConfig(projectPath string) error
	GetProjectConfig() *spookyconfigtypes.ProjectConfig

	// CLI integration
	ApplyCLIFlags(flags map[string]interface{}) error

	// Utility operations
	GetConfigurationPath() string
	CreateDefaultConfig() error
	Close() error
}

// Loader defines the interface for configuration loading
type Loader interface {
	LoadGlobalConfig() (*spookyconfigtypes.GlobalConfig, error)
	LoadProjectConfig(projectPath string) (*spookyconfigtypes.ProjectConfig, error)
	LoadFromFile(path string) (interface{}, error)
	LoadFromEnvironment() (map[string]interface{}, error)
}

// Validator defines the interface for configuration validation
type Validator interface {
	ValidateGlobalConfig(config *spookyconfigtypes.GlobalConfig) error
	ValidateProjectConfig(config *spookyconfigtypes.ProjectConfig) error
	ValidateConfigFile(path string) error
	ValidateAgainstSchema(config interface{}, schemaName string) error
}

// EnvironmentManager defines the interface for environment variable management
type EnvironmentManager interface {
	GetEnvironmentVariable(name string) (interface{}, error)
	GetAllEnvironmentVariables() (map[string]interface{}, error)
	ValidateEnvironmentVariable(name, value string) error
	GetEnvironmentVariableDescription(name string) string
	GetEnvironmentVariableDefault(name string) interface{}
}
