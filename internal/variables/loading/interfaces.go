package loading

import (
	"context"
	"spooky/internal/schemas"
	"spooky/internal/variables/types"
)

// LoadingManager defines the interface for variable loading operations
type LoadingManager interface {
	// Core loading operations
	LoadFromFile(ctx context.Context, path string) ([]*types.Variable, error)
	LoadFromDirectory(ctx context.Context, dirPath string) ([]*types.Variable, error)
	LoadFromHCL(ctx context.Context, content []byte) ([]*types.Variable, error)
	LoadFromJSON(ctx context.Context, content []byte) ([]*types.Variable, error)

	// Schema validation
	ValidateFileAgainstSchema(path string, schemaType schemas.SchemaType) error
	ValidateContentAgainstSchema(content []byte, schemaType schemas.SchemaType) error

	// Configuration
	SetDefaultEncoding(encoding string) error
	SetMaxFileSize(maxSize int64) error
	SetAllowedExtensions(extensions []string) error

	// Utility operations
	ValidateFile(path string) error
	GetSupportedFormats() []string
	Close() error
}

// VariableLoader defines the interface for specific format loaders
type VariableLoader interface {
	Load(ctx context.Context, source interface{}) ([]*types.Variable, error)
	GetName() string
	GetSupportedExtensions() []string
	ValidateSchema(content []byte) error
}
