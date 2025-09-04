package schemas

import "github.com/pkg/errors"

// Schema-specific error types
var (
	ErrSchemaNotFound         = errors.New("schema not found")
	ErrSchemaInvalid          = errors.New("schema is invalid")
	ErrSchemaVersionMismatch  = errors.New("schema version mismatch")
	ErrSchemaParseFailed      = errors.New("failed to parse schema")
	ErrSchemaEmbedFailed      = errors.New("failed to embed schema")
	ErrSchemaLoadFailed       = errors.New("failed to load schema")
	ErrValidationRuleNotFound = errors.New("validation rule not found")
	ErrTestDataNotFound       = errors.New("test data not found")
	ErrMetadataNotFound       = errors.New("metadata not found")
)

// SchemaNotFoundError provides structured error information
type SchemaNotFoundError struct {
	SchemaName string
	Available  []string
}

func (e *SchemaNotFoundError) Error() string {
	return errors.Errorf("schema '%s' not found", e.SchemaName).Error()
}

func (e *SchemaNotFoundError) Unwrap() error {
	return ErrSchemaNotFound
}

// SchemaValidationError provides detailed validation error information
type SchemaValidationError struct {
	SchemaName string
	Field      string
	Value      interface{}
	Rule       string
	Message    string
}

func (e *SchemaValidationError) Error() string {
	return errors.Errorf("schema '%s' validation failed for field '%s': %s",
		e.SchemaName, e.Field, e.Message).Error()
}

func (e *SchemaValidationError) Unwrap() error {
	return ErrSchemaInvalid
}

// SchemaVersionError provides version mismatch information
type SchemaVersionError struct {
	SchemaName    string
	Expected      string
	Actual        string
	Compatibility []string
}

func (e *SchemaVersionError) Error() string {
	return errors.Errorf("schema '%s' version mismatch: expected %s, got %s",
		e.SchemaName, e.Expected, e.Actual).Error()
}

func (e *SchemaVersionError) Unwrap() error {
	return ErrSchemaVersionMismatch
}

// SchemaParseError provides parsing error information
type SchemaParseError struct {
	SchemaName string
	FilePath   string
	Line       int
	Column     int
	Message    string
}

func (e *SchemaParseError) Error() string {
	if e.Line > 0 {
		return errors.Errorf("failed to parse schema '%s' at %s:%d:%d: %s",
			e.SchemaName, e.FilePath, e.Line, e.Column, e.Message).Error()
	}
	return errors.Errorf("failed to parse schema '%s' from %s: %s",
		e.SchemaName, e.FilePath, e.Message).Error()
}

func (e *SchemaParseError) Unwrap() error {
	return ErrSchemaParseFailed
}

// NewSchemaNotFoundError creates a new schema not found error.
func NewSchemaNotFoundError(schemaName string, available []string) error {
	return errors.WithStack(&SchemaNotFoundError{
		SchemaName: schemaName,
		Available:  available,
	})
}

// NewSchemaValidationError creates a new schema validation error.
func NewSchemaValidationError(schemaName, field string, value interface{}, rule, message string) error {
	return errors.WithStack(&SchemaValidationError{
		SchemaName: schemaName,
		Field:      field,
		Value:      value,
		Rule:       rule,
		Message:    message,
	})
}

// NewSchemaVersionError creates a new schema version error.
func NewSchemaVersionError(schemaName, expected, actual string, compatibility []string) error {
	return errors.WithStack(&SchemaVersionError{
		SchemaName:    schemaName,
		Expected:      expected,
		Actual:        actual,
		Compatibility: compatibility,
	})
}

// NewSchemaParseError creates a new schema parse error.
func NewSchemaParseError(schemaName, filePath, message string, line, column int) error {
	return errors.WithStack(&SchemaParseError{
		SchemaName: schemaName,
		FilePath:   filePath,
		Line:       line,
		Column:     column,
		Message:    message,
	})
}
