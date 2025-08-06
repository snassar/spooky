package types

import "errors"

// Schema-related errors
var (
	ErrSchemaNotFound    = errors.New("schema not found")
	ErrInvalidSchemaType = errors.New("invalid schema type")
	ErrSchemaLoadFailed  = errors.New("failed to load schema")
	ErrValidationFailed  = errors.New("validation failed")
	ErrInvalidConfig     = errors.New("invalid configuration")
)
