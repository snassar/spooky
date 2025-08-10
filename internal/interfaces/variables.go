package interfaces

import (
	"context"
	"io"

	spookytypesvariables "spooky/internal/types/variables"
)

// VariableManager defines the main interface for variable management
type VariableManager interface {
	// Core operations
	LoadVariables(ctx context.Context, path string) (*spookytypesvariables.VariableCollection, error)
	GetVariable(ctx context.Context, name string) (*spookytypesvariables.Variable, error)
	SetVariable(ctx context.Context, variable *spookytypesvariables.Variable) error
	DeleteVariable(ctx context.Context, name string) error
	ListVariables(ctx context.Context) ([]*spookytypesvariables.Variable, error)

	// Context operations
	CreateContext(ctx context.Context, variables []*spookytypesvariables.Variable) (*spookytypesvariables.VariableContext, error)
	ResolveContext(ctx context.Context, context *spookytypesvariables.VariableContext) error

	// Import/export operations
	ExportVariables(ctx context.Context, format spookytypesvariables.ExportFormat, path string) error
	ImportVariables(ctx context.Context, format spookytypesvariables.ImportFormat, path string) error

	// Validation operations
	ValidateVariables(ctx context.Context, variables []*spookytypesvariables.Variable) (*spookytypesvariables.ValidationResult, error)
	ValidateContext(ctx context.Context, context *spookytypesvariables.VariableContext) (*spookytypesvariables.ValidationResult, error)

	// Coordinator integration methods
	LoadVariablesForProject(projectPath string) (*spookytypesvariables.VariableCollection, error)
	ResolveVariablesForContext(context *spookytypesvariables.VariableContext) error
	ValidateVariablesForProject(projectPath string) (*spookytypesvariables.ValidationResult, error)
	ExportVariablesForProject(projectPath string, format spookytypesvariables.ExportFormat, outputPath string) error
}

// VariableLoader defines the interface for loading variables from different sources
type VariableLoader interface {
	LoadFromFile(ctx context.Context, path string) ([]*spookytypesvariables.Variable, error)
	LoadFromDirectory(ctx context.Context, dirPath string) ([]*spookytypesvariables.Variable, error)
	LoadFromHCL(ctx context.Context, content []byte) ([]*spookytypesvariables.Variable, error)
	LoadFromJSON(ctx context.Context, content []byte) ([]*spookytypesvariables.Variable, error)
}

// VariableResolver defines the interface for resolving variable dependencies
type VariableResolver interface {
	ResolveVariable(ctx context.Context, variable *spookytypesvariables.Variable, context *spookytypesvariables.VariableContext) error
	ResolveDependencies(ctx context.Context, variables []*spookytypesvariables.Variable) error
	ValidateDependencies(ctx context.Context, variables []*spookytypesvariables.Variable) error
}

// VariableValidator defines the interface for variable validation
type VariableValidator interface {
	ValidateVariable(ctx context.Context, variable *spookytypesvariables.Variable) (*spookytypesvariables.ValidationResult, error)
	ValidateCollection(ctx context.Context, collection *spookytypesvariables.VariableCollection) (*spookytypesvariables.ValidationResult, error)
	ValidateContext(ctx context.Context, context *spookytypesvariables.VariableContext) (*spookytypesvariables.ValidationResult, error)
}

// ImportExportManager defines the interface for import/export operations
type ImportExportManager interface {
	ExportToHCL(ctx context.Context, variables []*spookytypesvariables.Variable, path string) error
	ExportToJSON(ctx context.Context, variables []*spookytypesvariables.Variable, path string) error
	ImportFromHCL(ctx context.Context, path string) ([]*spookytypesvariables.Variable, error)
	ImportFromJSON(ctx context.Context, path string) ([]*spookytypesvariables.Variable, error)
}

// LoadingManager defines the interface for variable loading operations
type LoadingManager interface {
	// Core loading operations
	LoadFromFile(ctx context.Context, path string) ([]*spookytypesvariables.Variable, error)
	LoadFromDirectory(ctx context.Context, dirPath string) ([]*spookytypesvariables.Variable, error)
	LoadFromHCL(ctx context.Context, content []byte) ([]*spookytypesvariables.Variable, error)
	LoadFromJSON(ctx context.Context, content []byte) ([]*spookytypesvariables.Variable, error)

	// Schema validation
	ValidateFileAgainstSchema(path string, schemaType string) error
	ValidateContentAgainstSchema(content []byte, schemaType string) error

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
	Load(ctx context.Context, source interface{}) ([]*spookytypesvariables.Variable, error)
	GetName() string
	GetSupportedExtensions() []string
	ValidateSchema(content []byte) error
}

// ResolutionManager defines the interface for variable resolution operations
type ResolutionManager interface {
	// Core resolution operations
	ResolveVariable(ctx context.Context, variable *spookytypesvariables.Variable, context *spookytypesvariables.VariableContext) error
	ResolveDependencies(ctx context.Context, variables []*spookytypesvariables.Variable) error
	ResolveContext(ctx context.Context, context *spookytypesvariables.VariableContext) error

	// Dependency management
	ValidateDependencies(ctx context.Context, variables []*spookytypesvariables.Variable) error
	DetectCircularDependencies(variables []*spookytypesvariables.Variable) error
	GetDependencyGraph(variables []*spookytypesvariables.Variable) (*spookytypesvariables.DependencyGraph, error)

	// Configuration
	SetMaxRecursionDepth(depth int) error
	SetDefaultValues(defaults map[string]interface{}) error
	EnableStrictMode(strict bool) error

	// Utility operations
	GetUnresolvedVariables(variables []*spookytypesvariables.Variable) []*spookytypesvariables.Variable
	GetResolutionOrder(variables []*spookytypesvariables.Variable) ([]*spookytypesvariables.Variable, error)
	Close() error
}

// VariableResolver defines the interface for specific resolution strategies
type VariableResolver interface {
	Resolve(ctx context.Context, variable *spookytypesvariables.Variable, context *spookytypesvariables.VariableContext) error
	GetName() string
	GetSupportedTypes() []string
}

// ValidationManager defines the interface for variable validation operations
type ValidationManager interface {
	// Core validation operations
	ValidateVariable(ctx context.Context, variable *spookytypesvariables.Variable) (*spookytypesvariables.ValidationResult, error)
	ValidateCollection(ctx context.Context, collection *spookytypesvariables.VariableCollection) (*spookytypesvariables.ValidationResult, error)
	ValidateContext(ctx context.Context, context *spookytypesvariables.VariableContext) (*spookytypesvariables.ValidationResult, error)

	// Schema validation
	ValidateAgainstSchema(variable *spookytypesvariables.Variable, schemaName string) (*spookytypesvariables.ValidationResult, error)
	ValidateCollectionAgainstSchema(collection *spookytypesvariables.VariableCollection, schemaName string) (*spookytypesvariables.ValidationResult, error)

	// Custom validation
	RegisterCustomValidator(name string, validator VariableValidator) error
	UnregisterCustomValidator(name string) error
	GetCustomValidators() []string

	// Configuration
	SetValidationRules(rules *spookytypesvariables.ValidationRules) error
	EnableStrictValidation(strict bool) error
	SetMaxValidationErrors(max int) error

	// Utility operations
	GetValidationErrors() []spookytypesvariables.ValidationError
	ClearValidationErrors() error
	Close() error
}

// VariableValidator defines the interface for custom validators
type VariableValidator interface {
	Validate(ctx context.Context, variable *spookytypesvariables.Variable) error
	GetName() string
	GetDescription() string
}

// VariableExporter defines the interface for specific export formats
type VariableExporter interface {
	Export(ctx context.Context, variables []*spookytypesvariables.Variable, writer io.Writer) error
	GetName() string
	GetExtension() string
}

// VariableImporter defines the interface for specific import formats
type VariableImporter interface {
	Import(ctx context.Context, reader io.Reader) ([]*spookytypesvariables.Variable, error)
	GetName() string
	GetExtension() string
}
