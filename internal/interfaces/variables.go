package interfaces

import (
	"context"
	"io"

	spookytypes "spooky/internal/types"
)

// VariableManager defines the main interface for variable management
type VariableManager interface {
	// Core operations
	LoadVariables(ctx context.Context, path string) (*spookytypes.VariableCollection, error)
	GetVariable(ctx context.Context, name string) (*spookytypes.Variable, error)
	SetVariable(ctx context.Context, variable *spookytypes.Variable) error
	DeleteVariable(ctx context.Context, name string) error
	ListVariables(ctx context.Context) ([]*spookytypes.Variable, error)

	// Context operations
	CreateContext(ctx context.Context, variables []*spookytypes.Variable) (*spookytypes.VariableContext, error)
	ResolveContext(ctx context.Context, context *spookytypes.VariableContext) error

	// Import/export operations
	ExportVariables(ctx context.Context, format spookytypes.VariableExportFormat, path string) error
	ImportVariables(ctx context.Context, format spookytypes.VariableImportFormat, path string) error

	// Validation operations
	ValidateVariables(ctx context.Context, variables []*spookytypes.Variable) (*spookytypes.VariableValidationResult, error)
	ValidateContext(ctx context.Context, context *spookytypes.VariableContext) (*spookytypes.VariableValidationResult, error)

	// Coordinator integration methods
	LoadVariablesForProject(projectPath string) (*spookytypes.VariableCollection, error)
	ResolveVariablesForContext(context *spookytypes.VariableContext) error
	ValidateVariablesForProject(projectPath string) (*spookytypes.VariableValidationResult, error)
	ExportVariablesForProject(projectPath string, format spookytypes.VariableExportFormat, outputPath string) error
}

// VariableLoaderCore defines the interface for loading variables from different sources
type VariableLoaderCore interface {
	LoadFromFile(ctx context.Context, path string) ([]*spookytypes.Variable, error)
	LoadFromDirectory(ctx context.Context, dirPath string) ([]*spookytypes.Variable, error)
	LoadFromHCL(ctx context.Context, content []byte) ([]*spookytypes.Variable, error)
	LoadFromJSON(ctx context.Context, content []byte) ([]*spookytypes.Variable, error)
}

// VariableResolverCore defines the interface for resolving variable dependencies
type VariableResolverCore interface {
	ResolveVariable(ctx context.Context, variable *spookytypes.Variable, context *spookytypes.VariableContext) error
	ResolveDependencies(ctx context.Context, variables []*spookytypes.Variable) error
	ValidateDependencies(ctx context.Context, variables []*spookytypes.Variable) error
}

// VariableValidatorCore defines the interface for variable validation
type VariableValidatorCore interface {
	ValidateVariable(ctx context.Context, variable *spookytypes.Variable) (*spookytypes.VariableValidationResult, error)
	ValidateCollection(ctx context.Context, collection *spookytypes.VariableCollection) (*spookytypes.VariableValidationResult, error)
	ValidateContext(ctx context.Context, context *spookytypes.VariableContext) (*spookytypes.VariableValidationResult, error)
}

// ImportExportManager defines the interface for import/export operations
type ImportExportManager interface {
	ExportToHCL(ctx context.Context, variables []*spookytypes.Variable, path string) error
	ExportToJSON(ctx context.Context, variables []*spookytypes.Variable, path string) error
	ImportFromHCL(ctx context.Context, path string) ([]*spookytypes.Variable, error)
	ImportFromJSON(ctx context.Context, path string) ([]*spookytypes.Variable, error)
}

// VariableLoadingManager defines the interface for variable loading operations
type VariableLoadingManager interface {
	// Core loading operations
	LoadFromFile(ctx context.Context, path string) ([]*spookytypes.Variable, error)
	LoadFromDirectory(ctx context.Context, dirPath string) ([]*spookytypes.Variable, error)
	LoadFromHCL(ctx context.Context, content []byte) ([]*spookytypes.Variable, error)
	LoadFromJSON(ctx context.Context, content []byte) ([]*spookytypes.Variable, error)

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
	Load(ctx context.Context, source interface{}) ([]*spookytypes.Variable, error)
	GetName() string
	GetSupportedExtensions() []string
	ValidateSchema(content []byte) error
}

// ResolutionManager defines the interface for variable resolution operations
type ResolutionManager interface {
	// Core resolution operations
	ResolveVariable(ctx context.Context, variable *spookytypes.Variable, context *spookytypes.VariableContext) error
	ResolveDependencies(ctx context.Context, variables []*spookytypes.Variable) error
	ResolveContext(ctx context.Context, context *spookytypes.VariableContext) error

	// Dependency management
	ValidateDependencies(ctx context.Context, variables []*spookytypes.Variable) error
	DetectCircularDependencies(variables []*spookytypes.Variable) error
	GetDependencyGraph(variables []*spookytypes.Variable) (*spookytypes.DependencyGraph, error)

	// Configuration
	SetMaxRecursionDepth(depth int) error
	SetDefaultValues(defaults map[string]interface{}) error
	EnableStrictMode(strict bool) error

	// Utility operations
	GetUnresolvedVariables(variables []*spookytypes.Variable) []*spookytypes.Variable
	GetResolutionOrder(variables []*spookytypes.Variable) ([]*spookytypes.Variable, error)
	Close() error
}

// VariableResolver defines the interface for specific resolution strategies
type VariableResolver interface {
	Resolve(ctx context.Context, variable *spookytypes.Variable, context *spookytypes.VariableContext) error
	GetName() string
	GetSupportedTypes() []string
}

// VariableValidationManager defines the interface for variable validation operations
type VariableValidationManager interface {
	// Core validation operations
	ValidateVariable(ctx context.Context, variable *spookytypes.Variable) (*spookytypes.VariableValidationResult, error)
	ValidateCollection(ctx context.Context, collection *spookytypes.VariableCollection) (*spookytypes.VariableValidationResult, error)
	ValidateContext(ctx context.Context, context *spookytypes.VariableContext) (*spookytypes.VariableValidationResult, error)

	// Schema validation
	ValidateAgainstSchema(variable *spookytypes.Variable, schemaName string) (*spookytypes.VariableValidationResult, error)
	ValidateCollectionAgainstSchema(collection *spookytypes.VariableCollection, schemaName string) (*spookytypes.VariableValidationResult, error)

	// Custom validation
	RegisterCustomValidator(name string, validator VariableValidator) error
	UnregisterCustomValidator(name string) error
	GetCustomValidators() []string

	// Configuration
	SetValidationRules(rules *spookytypes.VariableValidationRules) error
	EnableStrictValidation(strict bool) error
	SetMaxValidationErrors(max int) error

	// Utility operations
	GetValidationErrors() []spookytypes.VariableValidationError
	ClearValidationErrors() error
	Close() error
}

// VariableValidator defines the interface for custom validators
type VariableValidator interface {
	Validate(ctx context.Context, variable *spookytypes.Variable) error
	GetName() string
	GetDescription() string
}

// VariableExporter defines the interface for specific export formats
type VariableExporter interface {
	Export(ctx context.Context, variables []*spookytypes.Variable, writer io.Writer) error
	GetName() string
	GetExtension() string
}

// VariableImporter defines the interface for specific import formats
type VariableImporter interface {
	Import(ctx context.Context, reader io.Reader) ([]*spookytypes.Variable, error)
	GetName() string
	GetExtension() string
}
