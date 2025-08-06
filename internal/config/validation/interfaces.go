package validation

import (
	spookyconfigtypes "spooky/internal/config/types"
)

// ValidationManager defines the interface for configuration validation
type ValidationManager interface {
	// Core validation operations
	ValidateGlobalConfig(config *spookyconfigtypes.GlobalConfig) error
	ValidateProjectConfig(config *spookyconfigtypes.ProjectConfig) error
	ValidateConfigFile(path string) error
	ValidateAgainstSchema(config interface{}, schemaName string) error

	// Configuration
	SetValidationRules(rules *spookyconfigtypes.ValidationRules) error
	EnableStrictValidation(strict bool) error

	// Utility operations
	GetValidationErrors() []spookyconfigtypes.ValidationError
	ClearValidationErrors() error
	Close() error
}

// ConfigValidator defines the interface for specific validators
type ConfigValidator interface {
	Validate(config interface{}) error
	GetName() string
	GetDescription() string
}
