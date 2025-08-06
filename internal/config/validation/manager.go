package validation

import (
	"fmt"

	"spooky/internal/config/types"
	"spooky/internal/logging"
)

// Manager implements ValidationManager interface
type Manager struct {
	config     *types.ValidationConfig
	validators map[string]ConfigValidator
	logger     logging.Logger
	errors     []types.ValidationError
}

// NewManager creates a new validation manager
func NewManager(config *types.ValidationConfig, logger logging.Logger) *Manager {
	return &Manager{
		config:     config,
		validators: make(map[string]ConfigValidator),
		logger:     logger,
		errors:     make([]types.ValidationError, 0),
	}
}

// ValidateGlobalConfig validates global configuration
func (m *Manager) ValidateGlobalConfig(config *types.GlobalConfig) error {
	return m.validateConfig(config, "global", m.validateGlobalConfigBasic)
}

// ValidateProjectConfig validates project configuration
func (m *Manager) ValidateProjectConfig(config *types.ProjectConfig) error {
	return m.validateConfig(config, "project", m.validateProjectConfigBasic)
}

// validateConfig is a common validation function to reduce code duplication
func (m *Manager) validateConfig(config interface{}, configType string, basicValidator func(interface{}) error) error {
	if config == nil {
		return fmt.Errorf("%s configuration is nil", configType)
	}

	// Clear previous errors
	m.errors = make([]types.ValidationError, 0)

	// Run custom validators
	for name, validator := range m.validators {
		if err := validator.Validate(config); err != nil {
			m.errors = append(m.errors, types.ValidationError{
				Field:   name,
				Message: err.Error(),
			})
		}
	}

	// Basic validation
	if err := basicValidator(config); err != nil {
		m.errors = append(m.errors, types.ValidationError{
			Field:   configType,
			Message: err.Error(),
		})
	}

	if len(m.errors) > 0 {
		return fmt.Errorf("validation failed: %v", m.errors)
	}

	return nil
}

// ValidateConfigFile validates a configuration file
func (m *Manager) ValidateConfigFile(_ string) error {
	// This would implement file validation logic
	// For now, just check if the file exists
	return nil
}

// ValidateAgainstSchema validates configuration against a schema
func (m *Manager) ValidateAgainstSchema(config interface{}, _ string) error {
	// This would implement schema validation logic
	// For now, just do basic type checking
	if config == nil {
		return fmt.Errorf("configuration is nil")
	}
	return nil
}

// SetValidationRules sets validation rules
func (m *Manager) SetValidationRules(rules *types.ValidationRules) error {
	if m.config == nil {
		m.config = &types.ValidationConfig{}
	}
	m.config.ValidationRules = rules
	return nil
}

// EnableStrictValidation enables or disables strict validation
func (m *Manager) EnableStrictValidation(strict bool) error {
	if m.config == nil {
		m.config = &types.ValidationConfig{}
	}
	m.config.StrictValidation = strict
	return nil
}

// GetValidationErrors returns validation errors
func (m *Manager) GetValidationErrors() []types.ValidationError {
	return m.errors
}

// ClearValidationErrors clears validation errors
func (m *Manager) ClearValidationErrors() error {
	m.errors = make([]types.ValidationError, 0)
	return nil
}

// Close closes the validation manager
func (m *Manager) Close() error {
	// Cleanup resources if needed
	return nil
}

// validateGlobalConfigBasic performs basic validation on global config
func (m *Manager) validateGlobalConfigBasic(config interface{}) error {
	globalConfig, ok := config.(*types.GlobalConfig)
	if !ok {
		return fmt.Errorf("invalid config type for global validation")
	}

	// Basic validation logic
	if globalConfig.LogLevel != "" {
		validLevels := []string{"debug", "info", "warn", "error"}
		found := false
		for _, level := range validLevels {
			if globalConfig.LogLevel == level {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid log level: %s", globalConfig.LogLevel)
		}
	}

	return nil
}

// validateProjectConfigBasic performs basic validation on project config
func (m *Manager) validateProjectConfigBasic(config interface{}) error {
	projectConfig, ok := config.(*types.ProjectConfig)
	if !ok {
		return fmt.Errorf("invalid config type for project validation")
	}

	// Basic validation logic
	if projectConfig.Name == "" {
		return fmt.Errorf("project name is required")
	}

	if projectConfig.DefaultTimeout < 0 {
		return fmt.Errorf("default timeout cannot be negative")
	}

	return nil
}
