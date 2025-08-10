package interfaces

import (
	spookytypesconfig "spooky/internal/types/config"
)

// ConfigManager defines the main interface for configuration management
type ConfigManager interface {
	// Core configuration operations
	LoadConfig() error
	ReloadConfig() error
	GetConfig() *spookytypesconfig.Config
	ValidateConfig() error

	// Value access operations
	GetValue(path string) (interface{}, spookytypesconfig.ConfigSource, error)
	SetValue(path string, value interface{}, source spookytypesconfig.ConfigSource) error
	GetString(path string) (string, spookytypesconfig.ConfigSource, error)
	GetInt(path string) (int, spookytypesconfig.ConfigSource, error)
	GetBool(path string) (bool, spookytypesconfig.ConfigSource, error)
	GetStringSlice(path string) ([]string, spookytypesconfig.ConfigSource, error)

	// Project configuration
	LoadProjectConfig(projectPath string) error
	GetProjectConfig() *spookytypesconfig.ProjectConfig

	// CLI integration
	ApplyCLIFlags(flags map[string]interface{}) error

	// Utility operations
	GetConfigurationPath() string
	CreateDefaultConfig() error
	Close() error
}

// Loader defines the interface for configuration loading
type Loader interface {
	LoadGlobalConfig() (*spookytypesconfig.GlobalConfig, error)
	LoadProjectConfig(projectPath string) (*spookytypesconfig.ProjectConfig, error)
	LoadFromFile(path string) (interface{}, error)
	LoadFromEnvironment() (map[string]interface{}, error)
}

// Validator defines the interface for configuration validation
type Validator interface {
	ValidateGlobalConfig(config *spookytypesconfig.GlobalConfig) error
	ValidateProjectConfig(config *spookytypesconfig.ProjectConfig) error
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

// LoadingManager defines the interface for configuration loading operations
type LoadingManager interface {
	// Core loading operations
	LoadGlobalConfig() (*spookytypesconfig.GlobalConfig, error)
	LoadProjectConfig(projectPath string) (*spookytypesconfig.ProjectConfig, error)
	LoadFromFile(path string) (interface{}, error)
	LoadFromEnvironment() (map[string]interface{}, error)

	// Configuration
	SetConfigPath(path string) error
	SetDefaultConfig(defaultConfig *spookytypesconfig.GlobalConfig) error
	EnableAutoReload(enabled bool) error

	// Utility operations
	GetConfigurationPath() string
	ValidateConfigFile(path string) error
	Close() error
}

// EnvironmentValidator defines the interface for environment validation
type EnvironmentValidator interface {
	ValidateVariable(name, value string) error
	ValidateVariableType(name, value, expectedType string) error
	ValidateVariableFormat(name, value, format string) error
}
