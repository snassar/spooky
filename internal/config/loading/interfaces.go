package loading

import (
	spookyconfigtypes "spooky/internal/config/types"
)

// LoadingManager defines the interface for configuration loading operations
type LoadingManager interface {
	// Core loading operations
	LoadGlobalConfig() (*spookyconfigtypes.GlobalConfig, error)
	LoadProjectConfig(projectPath string) (*spookyconfigtypes.ProjectConfig, error)
	LoadFromFile(path string) (interface{}, error)
	LoadFromEnvironment() (map[string]interface{}, error)

	// Configuration
	SetConfigPath(path string) error
	SetDefaultConfig(defaultConfig *spookyconfigtypes.GlobalConfig) error
	EnableAutoReload(enabled bool) error

	// Utility operations
	GetConfigPath() string
	ValidateConfigPath(path string) error
	Close() error
}

// ConfigParser defines the interface for configuration parsing
type ConfigParser interface {
	ParseHCL(content []byte) (interface{}, error)
	ParseJSON(content []byte) (interface{}, error)
	ValidateFormat(content []byte, format string) error
}

// XDGManager defines the interface for XDG directory handling
type XDGManager interface {
	GetConfigHome() string
	GetConfigDirs() []string
	GetDataHome() string
	GetDataDirs() []string
	GetCacheHome() string
	GetRuntimeDir() string
	CreateConfigDirectories() error
}
