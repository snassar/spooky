package environment

import (
	"fmt"
	"os"
	"strings"

	spookyconfigtypes "spooky/internal/config/types"
	spookylogging "spooky/internal/logging"
)

// Manager implements EnvironmentManager interface
type Manager struct {
	config    *spookyconfigtypes.EnvironmentConfig
	validator EnvironmentValidator
	logger    spookylogging.Logger
}

// NewManager creates a new environment manager
func NewManager(config *spookyconfigtypes.EnvironmentConfig, logger spookylogging.Logger) *Manager {
	return &Manager{
		config:    config,
		validator: NewEnvironmentValidator(),
		logger:    logger,
	}
}

// GetEnvironmentVariable gets an environment variable
func (m *Manager) GetEnvironmentVariable(name string) (interface{}, error) {
	value := os.Getenv(name)
	if value == "" {
		return nil, fmt.Errorf("environment variable not found: %s", name)
	}

	// Validate variable
	if err := m.validator.ValidateVariable(name, value); err != nil {
		return nil, fmt.Errorf("environment variable validation failed: %w", err)
	}

	return value, nil
}

// GetAllEnvironmentVariables gets all environment variables
func (m *Manager) GetAllEnvironmentVariables() (map[string]interface{}, error) {
	envVars := make(map[string]interface{})

	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			name := pair[0]
			value := pair[1]

			// Validate variable
			if err := m.validator.ValidateVariable(name, value); err != nil {
				m.logger.Warn("Invalid environment variable", spookylogging.String("name", name), spookylogging.Error(err))
				continue
			}

			envVars[name] = value
		}
	}

	return envVars, nil
}

// ValidateEnvironmentVariable validates an environment variable
func (m *Manager) ValidateEnvironmentVariable(name, value string) error {
	return m.validator.ValidateVariable(name, value)
}

// SetEnvironmentVariable sets an environment variable
func (m *Manager) SetEnvironmentVariable(name, value string) error {
	// Validate variable before setting
	if err := m.validator.ValidateVariable(name, value); err != nil {
		return fmt.Errorf("environment variable validation failed: %w", err)
	}

	return os.Setenv(name, value)
}

// UnsetEnvironmentVariable unsets an environment variable
func (m *Manager) UnsetEnvironmentVariable(name string) error {
	return os.Unsetenv(name)
}

// LoadEnvironmentFile loads environment variables from a file
func (m *Manager) LoadEnvironmentFile(path string) error {
	// This would implement loading from .env files
	// For now, just check if the file exists
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("environment file not found: %w", err)
	}

	return nil
}

// GetEnvironmentVariableDescription gets the description for an environment variable
func (m *Manager) GetEnvironmentVariableDescription(name string) string {
	// This would return descriptions for known environment variables
	descriptions := map[string]string{
		"SPOOKY_LOG_LEVEL":   "Log level for spooky (debug, info, warn, error)",
		"SPOOKY_CONFIG_PATH": "Path to spooky configuration file",
		"SPOOKY_DATA_PATH":   "Path to spooky data directory",
	}

	if desc, exists := descriptions[name]; exists {
		return desc
	}

	return "No description available"
}

// GetEnvironmentVariableDefault gets the default value for an environment variable
func (m *Manager) GetEnvironmentVariableDefault(name string) interface{} {
	// This would return default values for known environment variables
	defaults := map[string]interface{}{
		"SPOOKY_LOG_LEVEL":   "info",
		"SPOOKY_CONFIG_PATH": "",
		"SPOOKY_DATA_PATH":   "",
	}

	if def, exists := defaults[name]; exists {
		return def
	}

	return nil
}

// ListEnvironmentVariables lists all environment variables
func (m *Manager) ListEnvironmentVariables() []string {
	var vars []string
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			vars = append(vars, pair[0])
		}
	}
	return vars
}

// Close closes the environment manager
func (m *Manager) Close() error {
	// Cleanup resources if needed
	return nil
}
