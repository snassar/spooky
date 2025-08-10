package config

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	spookyconfigenvironment "spooky/internal/config/environment"
	spookyconfigloading "spooky/internal/config/loading"
	spookyconfigtypes "spooky/internal/types/config"
	spookyconfigvalidation "spooky/internal/config/validation"
	spookylogging "spooky/internal/logging"
)

// Manager implements ConfigManager interface
type Manager struct {
	config             *spookyconfigtypes.Config
	loadingManager     spookyconfigloading.LoadingManager
	validationManager  spookyconfigvalidation.ValidationManager
	environmentManager spookyconfigenvironment.EnvironmentManager
	logger             spookylogging.Logger
	mutex              sync.RWMutex
}

// NewManager creates a new configuration manager
func NewManager(
	config *spookyconfigtypes.Config,
	loadingManager spookyconfigloading.LoadingManager,
	validationManager spookyconfigvalidation.ValidationManager,
	environmentManager spookyconfigenvironment.EnvironmentManager,
	logger spookylogging.Logger,
) *Manager {
	return &Manager{
		config:             config,
		loadingManager:     loadingManager,
		validationManager:  validationManager,
		environmentManager: environmentManager,
		logger:             logger,
	}
}

// LoadConfig loads the configuration from all sources
func (m *Manager) LoadConfig() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 1. Load global configuration
	globalConfig, err := m.loadingManager.LoadGlobalConfig()
	if err != nil {
		return fmt.Errorf("failed to load global configuration: %w", err)
	}

	// 2. Load environment variables
	envVars, err := m.environmentManager.GetAllEnvironmentVariables()
	if err != nil {
		return fmt.Errorf("failed to load environment variables: %w", err)
	}

	// 3. Validate global configuration
	if err := m.validationManager.ValidateGlobalConfig(globalConfig); err != nil {
		return fmt.Errorf("global configuration validation failed: %w", err)
	}

	// 4. Create configuration context
	m.config = &spookyconfigtypes.Config{
		GlobalConfig: globalConfig,
		Environment:  envVars,
		Source:       spookyconfigtypes.SourceGlobal,
	}

	return nil
}

// ReloadConfig reloads the configuration from all sources
func (m *Manager) ReloadConfig() error {
	return m.LoadConfig()
}

// GetConfig returns the current configuration
func (m *Manager) GetConfig() *spookyconfigtypes.Config {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.config
}

// GetValue gets a configuration value by dot-separated path
func (m *Manager) GetValue(path string) (interface{}, spookyconfigtypes.ConfigSource, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.config == nil {
		return nil, spookyconfigtypes.SourceDefault, fmt.Errorf("no configuration loaded")
	}

	value, source, err := m.getValueByPath(m.config.GlobalConfig, path)
	if err != nil {
		return nil, spookyconfigtypes.SourceDefault, err
	}

	return value, source, nil
}

// SetValue sets a configuration value by dot-separated path
func (m *Manager) SetValue(path string, value interface{}, source spookyconfigtypes.ConfigSource) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.config == nil {
		return fmt.Errorf("no configuration loaded")
	}

	// Set the source for tracking where the value came from
	if m.config.Source == spookyconfigtypes.SourceDefault {
		m.config.Source = source
	}

	return m.setValueByPath(m.config.GlobalConfig, path, value)
}

// GetString gets a string configuration value
func (m *Manager) GetString(path string) (string, spookyconfigtypes.ConfigSource, error) {
	value, source, err := m.GetValue(path)
	if err != nil {
		return "", spookyconfigtypes.SourceDefault, err
	}

	if str, ok := value.(string); ok {
		return str, source, nil
	}

	return "", spookyconfigtypes.SourceDefault, fmt.Errorf("value at path %s is not a string", path)
}

// GetInt gets an integer configuration value
func (m *Manager) GetInt(path string) (int, spookyconfigtypes.ConfigSource, error) {
	value, source, err := m.GetValue(path)
	if err != nil {
		return 0, spookyconfigtypes.SourceDefault, err
	}

	if i, ok := value.(int); ok {
		return i, source, nil
	}

	return 0, spookyconfigtypes.SourceDefault, fmt.Errorf("value at path %s is not an integer", path)
}

// GetBool gets a boolean configuration value
func (m *Manager) GetBool(path string) (bool, spookyconfigtypes.ConfigSource, error) {
	value, source, err := m.GetValue(path)
	if err != nil {
		return false, spookyconfigtypes.SourceDefault, err
	}

	if b, ok := value.(bool); ok {
		return b, source, nil
	}

	return false, spookyconfigtypes.SourceDefault, fmt.Errorf("value at path %s is not a boolean", path)
}

// GetStringSlice gets a string slice configuration value
func (m *Manager) GetStringSlice(path string) ([]string, spookyconfigtypes.ConfigSource, error) {
	value, source, err := m.GetValue(path)
	if err != nil {
		return nil, spookyconfigtypes.SourceDefault, err
	}

	if slice, ok := value.([]string); ok {
		return slice, source, nil
	}

	return nil, spookyconfigtypes.SourceDefault, fmt.Errorf("value at path %s is not a string slice", path)
}

// LoadProjectConfig loads a project configuration
func (m *Manager) LoadProjectConfig(projectPath string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Load project configuration
	projectConfig, err := m.loadingManager.LoadProjectConfig(projectPath)
	if err != nil {
		return fmt.Errorf("failed to load project configuration: %w", err)
	}

	// Validate project configuration
	if err := m.validationManager.ValidateProjectConfig(projectConfig); err != nil {
		return fmt.Errorf("project configuration validation failed: %w", err)
	}

	// Set project configuration
	if m.config == nil {
		// Load global config if not already loaded
		globalConfig, err := m.loadingManager.LoadGlobalConfig()
		if err != nil {
			return fmt.Errorf("failed to load global configuration: %w", err)
		}

		envVars, err := m.environmentManager.GetAllEnvironmentVariables()
		if err != nil {
			return fmt.Errorf("failed to load environment variables: %w", err)
		}

		m.config = &spookyconfigtypes.Config{
			GlobalConfig:  globalConfig,
			ProjectConfig: projectConfig,
			Environment:   envVars,
			Source:        spookyconfigtypes.SourceProject,
		}
	} else {
		// Update existing context
		m.config.ProjectConfig = projectConfig
		m.config.Source = spookyconfigtypes.SourceProject
	}

	return nil
}

// GetProjectConfig returns the current project configuration
func (m *Manager) GetProjectConfig() *spookyconfigtypes.ProjectConfig {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.config == nil {
		return nil
	}

	return m.config.ProjectConfig
}

// ApplyCLIFlags applies CLI flags to the configuration
func (m *Manager) ApplyCLIFlags(flags map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.config == nil {
		return fmt.Errorf("no configuration loaded")
	}

	// Apply CLI overrides
	m.config.CLIFlags = flags
	m.config.Source = spookyconfigtypes.SourceCLI

	return nil
}

// ValidateConfig validates the current configuration
func (m *Manager) ValidateConfig() error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.config == nil {
		return fmt.Errorf("no configuration loaded")
	}

	if err := m.validationManager.ValidateGlobalConfig(m.config.GlobalConfig); err != nil {
		return fmt.Errorf("global configuration validation failed: %w", err)
	}

	if m.config.ProjectConfig != nil {
		if err := m.validationManager.ValidateProjectConfig(m.config.ProjectConfig); err != nil {
			return fmt.Errorf("project configuration validation failed: %w", err)
		}
	}

	return nil
}

// GetConfigurationPath returns the current configuration file path
func (m *Manager) GetConfigurationPath() string {
	return m.loadingManager.GetConfigPath()
}

// CreateDefaultConfig creates a default configuration file
func (m *Manager) CreateDefaultConfig() error {
	// This would create a default configuration file
	// For now, just return success
	return nil
}

// Close closes the configuration manager
func (m *Manager) Close() error {
	// Close all submanagers
	if err := m.loadingManager.Close(); err != nil {
		return fmt.Errorf("failed to close loading manager: %w", err)
	}

	if err := m.validationManager.Close(); err != nil {
		return fmt.Errorf("failed to close validation manager: %w", err)
	}

	if err := m.environmentManager.Close(); err != nil {
		return fmt.Errorf("failed to close environment manager: %w", err)
	}

	return nil
}

// getValueByPath gets a value from a configuration object by dot-separated path
func (m *Manager) getValueByPath(config *spookyconfigtypes.GlobalConfig, path string) (interface{}, spookyconfigtypes.ConfigSource, error) {
	parts := strings.Split(path, ".")
	if len(parts) < 2 {
		return nil, spookyconfigtypes.SourceDefault, fmt.Errorf("invalid path format: %s", path)
	}

	section := parts[0]
	field := parts[1]

	switch section {
	case "storage":
		if config.Storage == nil {
			return nil, spookyconfigtypes.SourceDefault, fmt.Errorf("storage configuration not available")
		}
		return m.getStructField(config.Storage, field), spookyconfigtypes.SourceGlobal, nil
	case "facts":
		if config.Facts == nil {
			return nil, spookyconfigtypes.SourceDefault, fmt.Errorf("facts configuration not available")
		}
		return m.getStructField(config.Facts, field), spookyconfigtypes.SourceGlobal, nil
	case "ssh":
		if config.SSH == nil {
			return nil, spookyconfigtypes.SourceDefault, fmt.Errorf("SSH configuration not available")
		}
		return m.getStructField(config.SSH, field), spookyconfigtypes.SourceGlobal, nil
	case "templates":
		if config.Templates == nil {
			return nil, spookyconfigtypes.SourceDefault, fmt.Errorf("templates configuration not available")
		}
		return m.getStructField(config.Templates, field), spookyconfigtypes.SourceGlobal, nil
	case "security":
		if config.Security == nil {
			return nil, spookyconfigtypes.SourceDefault, fmt.Errorf("security configuration not available")
		}
		return m.getStructField(config.Security, field), spookyconfigtypes.SourceGlobal, nil
	case "age":
		if config.Age == nil {
			return nil, spookyconfigtypes.SourceDefault, fmt.Errorf("age configuration not available")
		}
		return m.getStructField(config.Age, field), spookyconfigtypes.SourceGlobal, nil
	case "logging":
		if config.Logging == nil {
			return nil, spookyconfigtypes.SourceDefault, fmt.Errorf("logging configuration not available")
		}
		return m.getStructField(config.Logging, field), spookyconfigtypes.SourceGlobal, nil
	case "performance":
		if config.Performance == nil {
			return nil, spookyconfigtypes.SourceDefault, fmt.Errorf("performance configuration not available")
		}
		return m.getStructField(config.Performance, field), spookyconfigtypes.SourceGlobal, nil
	default:
		return nil, spookyconfigtypes.SourceDefault, fmt.Errorf("unknown configuration section: %s", section)
	}
}

// setValueByPath sets a value in a configuration object by dot-separated path
func (m *Manager) setValueByPath(config *spookyconfigtypes.GlobalConfig, path string, value interface{}) error {
	parts := strings.Split(path, ".")
	if len(parts) < 2 {
		return fmt.Errorf("invalid path format: %s", path)
	}

	section := parts[0]
	field := parts[1]

	switch section {
	case "storage":
		if config.Storage == nil {
			config.Storage = &spookyconfigtypes.StorageConfig{}
		}
		return m.setStructField(config.Storage, field, value)
	case "facts":
		if config.Facts == nil {
			config.Facts = &spookyconfigtypes.FactsConfig{}
		}
		return m.setStructField(config.Facts, field, value)
	case "ssh":
		if config.SSH == nil {
			config.SSH = &spookyconfigtypes.SSHConfig{}
		}
		return m.setStructField(config.SSH, field, value)
	case "templates":
		if config.Templates == nil {
			config.Templates = &spookyconfigtypes.TemplatesConfig{}
		}
		return m.setStructField(config.Templates, field, value)
	case "security":
		if config.Security == nil {
			config.Security = &spookyconfigtypes.SecurityConfig{}
		}
		return m.setStructField(config.Security, field, value)
	case "age":
		if config.Age == nil {
			config.Age = &spookyconfigtypes.AgeConfig{}
		}
		return m.setStructField(config.Age, field, value)
	case "logging":
		if config.Logging == nil {
			config.Logging = &spookyconfigtypes.LoggingConfig{}
		}
		return m.setStructField(config.Logging, field, value)
	case "performance":
		if config.Performance == nil {
			config.Performance = &spookyconfigtypes.PerformanceConfig{}
		}
		return m.setStructField(config.Performance, field, value)
	default:
		return fmt.Errorf("unknown configuration section: %s", section)
	}
}

// getStructField gets a field value from a struct using reflection
func (m *Manager) getStructField(obj interface{}, fieldName string) interface{} {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return nil
	}

	return field.Interface()
}

// setStructField sets a field value on a struct using reflection
func (m *Manager) setStructField(obj interface{}, fieldName string, value interface{}) error {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() || !field.CanSet() {
		return fmt.Errorf("field %s is not valid or cannot be set", fieldName)
	}

	// Convert value to the correct type
	val := reflect.ValueOf(value)
	if val.Type().ConvertibleTo(field.Type()) {
		field.Set(val.Convert(field.Type()))
		return nil
	}

	return fmt.Errorf("cannot convert value of type %T to field type %s", value, field.Type())
}
