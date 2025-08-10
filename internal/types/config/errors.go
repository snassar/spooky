package config

import (
	"fmt"
)

// ConfigErrorDetails represents additional configuration error details
type ConfigErrorDetails struct {
	Path   string `hcl:"path,optional"`
	Source string `hcl:"source,optional"`
}

// ValidationErrorDetails represents additional validation error details
type ValidationErrorDetails struct {
	Value interface{} `hcl:"value,optional"`
	Rule  string      `hcl:"rule,optional"`
}

// LoadingError represents a configuration loading error
type LoadingError struct {
	File    string `hcl:"file"`
	Message string `hcl:"message"`
	Details string `hcl:"details,optional"`
}

// Error implements the error interface
func (e *LoadingError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("failed to load config file '%s': %s (%s)", e.File, e.Message, e.Details)
	}
	return fmt.Sprintf("failed to load config file '%s': %s", e.File, e.Message)
}

// ParsingError represents a configuration parsing error
type ParsingError struct {
	File    string `hcl:"file"`
	Line    int    `hcl:"line,optional"`
	Column  int    `hcl:"column,optional"`
	Message string `hcl:"message"`
	Context string `hcl:"context,optional"`
}

// Error implements the error interface
func (e *ParsingError) Error() string {
	if e.Line > 0 && e.Column > 0 {
		return fmt.Sprintf("parsing error in '%s' at line %d, column %d: %s", e.File, e.Line, e.Column, e.Message)
	}
	return fmt.Sprintf("parsing error in '%s': %s", e.File, e.Message)
}

// SchemaError represents a schema validation error
type SchemaError struct {
	Schema  string `hcl:"schema"`
	Field   string `hcl:"field,optional"`
	Message string `hcl:"message"`
	Details string `hcl:"details,optional"`
}

// Error implements the error interface
func (e *SchemaError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("schema validation error for '%s' field '%s': %s", e.Schema, e.Field, e.Message)
	}
	return fmt.Sprintf("schema validation error for '%s': %s", e.Schema, e.Message)
}

// EnvironmentError represents an environment variable error
type EnvironmentError struct {
	Variable string `hcl:"variable"`
	Message  string `hcl:"message"`
	Value    string `hcl:"value,optional"`
}

// Error implements the error interface
func (e *EnvironmentError) Error() string {
	if e.Value != "" {
		return fmt.Sprintf("environment variable error for '%s' (value: %s): %s", e.Variable, e.Value, e.Message)
	}
	return fmt.Sprintf("environment variable error for '%s': %s", e.Variable, e.Message)
}

// Error types for specific operations
var (
	ErrConfigNotLoaded     = fmt.Errorf("configuration not loaded")
	ErrConfigPathInvalid   = fmt.Errorf("invalid configuration path")
	ErrConfigFileNotFound  = fmt.Errorf("configuration file not found")
	ErrConfigFileInvalid   = fmt.Errorf("configuration file is invalid")
	ErrConfigValueNotFound = fmt.Errorf("configuration value not found")
	ErrConfigValueInvalid  = fmt.Errorf("configuration value is invalid")
)

// ConfigValidator defines the interface for configuration validators
type ConfigValidator interface {
	Validate(config interface{}) error
	GetName() string
	GetDescription() string
}

// EnvironmentValidator defines the interface for environment variable validators
type EnvironmentValidator interface {
	ValidateVariable(name, value string) error
	GetName() string
	GetDescription() string
}

// ConfigParser defines the interface for configuration parsers
type ConfigParser interface {
	Parse(data []byte) (interface{}, error)
	ParseFile(path string) (interface{}, error)
	GetName() string
	GetSupportedExtensions() []string
}

// XDGManager defines the interface for XDG base directory management
type XDGManager interface {
	GetConfigHome() string
	GetDataHome() string
	GetCacheHome() string
	GetRuntimeDir() string
	GetConfigDirs() []string
	GetDataDirs() []string
}
