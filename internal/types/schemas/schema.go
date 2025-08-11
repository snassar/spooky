// Package schemas provides types for schema validation and management in the spooky codebase.
// These types define the structure for HCL schemas and validation operations.
package schemas

import (
	"time"
)

// Schema represents a HCL schema definition
type Schema struct {
	// Schema metadata
	Version     string    `json:"version" hcl:"version"`
	Type        string    `json:"type" hcl:"type"`
	Name        string    `json:"name" hcl:"name"`
	Description string    `json:"description" hcl:"description"`
	CreatedAt   time.Time `json:"created_at" hcl:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" hcl:"updated_at"`

	// Schema content
	Content string `json:"content" hcl:"content"`

	// Schema validation rules
	Validation *SchemaValidation `json:"validation,omitempty" hcl:"validation,optional"`

	// Schema metadata
	Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// SchemaValidation represents schema validation configuration
type SchemaValidation struct {
	// Whether validation is enabled
	Enabled bool `json:"enabled" hcl:"enabled"`

	// Validation mode (strict, lenient)
	Mode string `json:"mode" hcl:"mode"`

	// Validation rules
	Rules []ValidationRule `json:"rules,omitempty" hcl:"rules,optional"`

	// Custom validation functions
	CustomValidators []string `json:"custom_validators,omitempty" hcl:"custom_validators,optional"`

	// Validation error handling
	ErrorHandling *ValidationErrorHandling `json:"error_handling,omitempty" hcl:"error_handling,optional"`
}

// ValidationRule represents a validation rule
type ValidationRule struct {
	// Rule name
	Name string `json:"name" hcl:"name"`

	// Rule type
	Type string `json:"type" hcl:"type"`

	// Rule pattern (for regex rules)
	Pattern string `json:"pattern,omitempty" hcl:"pattern,optional"`

	// Rule condition (for conditional rules)
	Condition string `json:"condition,omitempty" hcl:"condition,optional"`

	// Rule message
	Message string `json:"message" hcl:"message"`

	// Rule severity
	Severity string `json:"severity" hcl:"severity"` // "error", "warning", "info"

	// Rule parameters
	Parameters map[string]interface{} `json:"parameters,omitempty" hcl:"parameters,optional"`
}

// ValidationErrorHandling represents validation error handling configuration
type ValidationErrorHandling struct {
	// Whether to stop on first error
	StopOnFirstError bool `json:"stop_on_first_error" hcl:"stop_on_first_error"`

	// Maximum number of errors to collect
	MaxErrors int `json:"max_errors" hcl:"max_errors"`

	// Whether to include warnings
	IncludeWarnings bool `json:"include_warnings" hcl:"include_warnings"`

	// Whether to include context in errors
	IncludeContext bool `json:"include_context" hcl:"include_context"`
}

// SchemaError represents a schema-related error
type SchemaError struct {
	// Error details
	Code        string                 `json:"code" hcl:"code"`
	Message     string                 `json:"message" hcl:"message"`
	Context     map[string]interface{} `json:"context,omitempty" hcl:"context,optional"`
	Stack       []string               `json:"stack,omitempty" hcl:"stack,optional"`
	Recoverable bool                   `json:"recoverable" hcl:"recoverable"`

	// Schema information
	SchemaName string `json:"schema_name" hcl:"schema_name"`
	SchemaType string `json:"schema_type" hcl:"schema_type"`

	// Validation information
	FieldPath string      `json:"field_path,omitempty" hcl:"field_path,optional"`
	Value     interface{} `json:"value,omitempty" hcl:"value,optional"`

	// Error severity
	Severity string `json:"severity" hcl:"severity"` // "error", "warning", "info"
}

// NewSchemaError creates a new schema error
func NewSchemaError(schemaName, schemaType, message string) *SchemaError {
	return &SchemaError{
		Code:        "schema_error",
		Message:     message,
		Recoverable: true,
		SchemaName:  schemaName,
		SchemaType:  schemaType,
		Severity:    "error",
	}
}

// Error implements the error interface
func (e *SchemaError) Error() string {
	return e.Message
}

// Unwrap returns the underlying error
func (e *SchemaError) Unwrap() error {
	return nil
}

// SchemaRegistry manages available schemas
type SchemaRegistry interface {
	// Register registers a new schema
	Register(schema *Schema) error

	// Get returns a schema by name and type
	Get(name, schemaType string) (*Schema, bool)

	// List returns all registered schemas
	List() []*Schema

	// ListByType returns schemas by type
	ListByType(schemaType string) []*Schema

	// Validate validates data against a schema
	Validate(schemaName, schemaType string, data interface{}) (*ValidationResult, error)
}

// ValidationResult represents the result of a validation operation
type ValidationResult struct {
	// Whether the validation passed
	Valid bool `json:"valid" hcl:"valid"`

	// Validation timestamp
	ValidatedAt time.Time `json:"validated_at" hcl:"validated_at"`

	// Validation errors
	Errors []SchemaError `json:"errors,omitempty" hcl:"errors,optional"`

	// Validation warnings
	Warnings []SchemaError `json:"warnings,omitempty" hcl:"warnings,optional"`

	// Validation info messages
	Info []SchemaError `json:"info,omitempty" hcl:"info,optional"`

	// Validation details
	Details map[string]interface{} `json:"details,omitempty" hcl:"details,optional"`
}

// SchemaValidator provides schema validation functionality
type SchemaValidator interface {
	// Validate validates data against a schema
	Validate(schema *Schema, data interface{}) (*ValidationResult, error)

	// ValidateFile validates a file against a schema
	ValidateFile(schema *Schema, filePath string) (*ValidationResult, error)

	// ValidateString validates a string against a schema
	ValidateString(schema *Schema, content string) (*ValidationResult, error)

	// ValidateBytes validates bytes against a schema
	ValidateBytes(schema *Schema, data []byte) (*ValidationResult, error)
}

// SchemaLoader provides schema loading functionality
type SchemaLoader interface {
	// Load loads a schema from a file
	Load(filePath string) (*Schema, error)

	// LoadFromString loads a schema from a string
	LoadFromString(content string) (*Schema, error)

	// LoadFromBytes loads a schema from bytes
	LoadFromBytes(data []byte) (*Schema, error)

	// LoadEmbedded loads an embedded schema
	LoadEmbedded(name string) (*Schema, error)
}
