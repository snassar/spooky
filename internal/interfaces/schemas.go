package interfaces

import (
	spookyschemas "spooky/internal/schemas"
	spookytypesschemas "spooky/internal/types/schemas"
)

// Validator interface defines the core validation functionality
type Validator interface {
	LoadSchema(schemaType spookyschemas.SchemaType) error
	LoadAllSchemas() error
	ValidateFile(filePath, schemaName string) *spookyschemas.ValidationResult
	ValidateData(data interface{}, schemaName string) error
	ValidateProject(projectPath string) *spookyschemas.ValidationResult
	ValidateFacts(factsPath string, format string) *spookyschemas.ValidationResult
	ValidateProjectDirectory(projectPath string) *spookyschemas.ValidationResult
	GetSchema(schemaType spookyschemas.SchemaType) (*spookyschemas.Schema, error)
	ListLoadedSchemas() []spookyschemas.SchemaType
	GetValidationErrors() []spookyschemas.ValidationError
}

// SchemaLoader interface defines schema loading strategies
type SchemaLoader interface {
	LoadSchema(schemaType spookyschemas.SchemaType) (*spookyschemas.Schema, error)
	LoadAllSchemas() (map[spookyschemas.SchemaType]*spookyschemas.Schema, error)
	GetName() string
}

// ValidationReporter interface defines validation reporting strategies
type ValidationReporter interface {
	Report(result *spookyschemas.ValidationResult) error
	GetName() string
}

// Manager interface defines the schema manager functionality
type Manager interface {
	Configure(config *spookyschemas.Config) error
	GetValidator() Validator
	GetLoader() SchemaLoader
	GetReporter() ValidationReporter
	ValidateFile(filePath, schemaName string) *spookyschemas.ValidationResult
	ValidateProject(projectPath string) *spookyschemas.ValidationResult
	Close() error
}

// Loader interface defines schema loading strategies
type Loader interface {
	LoadSchema(schemaType spookyschemas.SchemaType) (*spookyschemas.Schema, error)
	LoadAllSchemas() (map[spookyschemas.SchemaType]*spookyschemas.Schema, error)
	GetName() string
}

// Reporter interface defines validation reporting strategies
type Reporter interface {
	Report(result *spookyschemas.ValidationResult) error
	GetName() string
}

// Validator interface defines validation strategies
type SchemaValidator interface {
	Validate(data interface{}, schema *spookytypesschemas.Schema) *spookytypesschemas.ValidationResult
	GetName() string
}
