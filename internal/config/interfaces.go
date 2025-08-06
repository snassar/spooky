package config

import (
	"spooky/internal/config/types"
)

// ConfigManager defines the main interface for configuration management
type ConfigManager interface {
	// Core configuration operations
	LoadConfig() error
	ReloadConfig() error
	GetConfig() *types.Config
	ValidateConfig() error

	// Value access operations
	GetValue(path string) (interface{}, types.ConfigSource, error)
	SetValue(path string, value interface{}, source types.ConfigSource) error
	GetString(path string) (string, types.ConfigSource, error)
	GetInt(path string) (int, types.ConfigSource, error)
	GetBool(path string) (bool, types.ConfigSource, error)
	GetStringSlice(path string) ([]string, types.ConfigSource, error)

	// Project configuration
	LoadProjectConfig(projectPath string) error
	GetProjectConfig() *types.ProjectConfig

	// CLI integration
	ApplyCLIFlags(flags map[string]interface{}) error

	// Utility operations
	GetConfigurationPath() string
	CreateDefaultConfig() error
	Close() error
}

// Loader defines the interface for configuration loading
type Loader interface {
	LoadGlobalConfig() (*types.GlobalConfig, error)
	LoadProjectConfig(projectPath string) (*types.ProjectConfig, error)
	LoadFromFile(path string) (interface{}, error)
	LoadFromEnvironment() (map[string]interface{}, error)
}

// Validator defines the interface for configuration validation
type Validator interface {
	ValidateGlobalConfig(config *types.GlobalConfig) error
	ValidateProjectConfig(config *types.ProjectConfig) error
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
