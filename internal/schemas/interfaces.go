package schemas

import (
	"spooky/internal/schemas/types"
)

// Validator interface defines the core validation functionality
type Validator interface {
	LoadSchema(schemaType types.SchemaType) error
	LoadAllSchemas() error
	ValidateFile(filePath, schemaName string) *types.ValidationResult
	ValidateData(data interface{}, schemaName string) error
	ValidateProject(projectPath string) *types.ValidationResult
	ValidateFacts(factsPath string, format string) *types.ValidationResult
	ValidateProjectDirectory(projectPath string) *types.ValidationResult
	GetSchema(schemaType types.SchemaType) (*types.Schema, error)
	ListLoadedSchemas() []types.SchemaType
	GetValidationErrors() []types.ValidationError
}

// SchemaLoader interface defines schema loading strategies
type SchemaLoader interface {
	LoadSchema(schemaType types.SchemaType) (*types.Schema, error)
	LoadAllSchemas() (map[types.SchemaType]*types.Schema, error)
	GetName() string
}

// ValidationReporter interface defines validation reporting strategies
type ValidationReporter interface {
	Report(result *types.ValidationResult) error
	GetName() string
}

// Manager interface defines the schema manager functionality
type Manager interface {
	Configure(config *types.Config) error
	GetValidator() Validator
	GetLoader() SchemaLoader
	GetReporter() ValidationReporter
	ValidateFile(filePath, schemaName string) *types.ValidationResult
	ValidateProject(projectPath string) *types.ValidationResult
	Close() error
}
