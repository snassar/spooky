package interfaces

import (
	spookyschemas "spooky/internal/schemas"
	spookytypesschemas "spooky/internal/types/schemas"
)

// SchemaValidator interface defines the core validation functionality
type SchemaValidator interface {
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

// SchemaManager interface defines the schema manager functionality
type SchemaManager interface {
	Configure(config *spookytypesschemas.Config) error
	GetValidator() SchemaValidator
	GetLoader() SchemaLoader
	GetReporter() ValidationReporter
	ValidateFile(filePath, schemaName string) *spookyschemas.ValidationResult
	ValidateProject(projectPath string) *spookyschemas.ValidationResult
	Close() error
}

// SchemaLoaderAlt interface defines schema loading strategies
type SchemaLoaderAlt interface {
	LoadSchema(schemaType spookyschemas.SchemaType) (*spookyschemas.Schema, error)
	LoadAllSchemas() (map[spookyschemas.SchemaType]*spookyschemas.Schema, error)
	GetName() string
}

// SchemaReporter interface defines validation reporting strategies
type SchemaReporter interface {
	Report(result *spookyschemas.ValidationResult) error
	GetName() string
}

// SchemaValidatorAlt interface defines validation strategies
type SchemaValidatorAlt interface {
	Validate(data interface{}, schema *spookytypesschemas.Schema) *spookytypesschemas.ValidationResult
	GetName() string
}
