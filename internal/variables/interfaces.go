package variables

import (
	"context"
	spookyvariablestypes "spooky/internal/variables/types"
)

// VariableManager defines the main interface for variable management
type VariableManager interface {
	// Core operations
	LoadVariables(ctx context.Context, path string) (*spookyvariablestypes.VariableCollection, error)
	GetVariable(ctx context.Context, name string) (*spookyvariablestypes.Variable, error)
	SetVariable(ctx context.Context, variable *spookyvariablestypes.Variable) error
	DeleteVariable(ctx context.Context, name string) error
	ListVariables(ctx context.Context) ([]*spookyvariablestypes.Variable, error)

	// Context operations
	CreateContext(ctx context.Context, variables []*spookyvariablestypes.Variable) (*spookyvariablestypes.VariableContext, error)
	ResolveContext(ctx context.Context, context *spookyvariablestypes.VariableContext) error

	// Import/export operations
	ExportVariables(ctx context.Context, format spookyvariablestypes.ExportFormat, path string) error
	ImportVariables(ctx context.Context, format spookyvariablestypes.ImportFormat, path string) error

	// Validation operations
	ValidateVariables(ctx context.Context, variables []*spookyvariablestypes.Variable) (*spookyvariablestypes.ValidationResult, error)
	ValidateContext(ctx context.Context, context *spookyvariablestypes.VariableContext) (*spookyvariablestypes.ValidationResult, error)

	// Coordinator integration methods
	LoadVariablesForProject(projectPath string) (*spookyvariablestypes.VariableCollection, error)
	ResolveVariablesForContext(context *spookyvariablestypes.VariableContext) error
	ValidateVariablesForProject(projectPath string) (*spookyvariablestypes.ValidationResult, error)
	ExportVariablesForProject(projectPath string, format spookyvariablestypes.ExportFormat, outputPath string) error
}

// VariableLoader defines the interface for loading variables from different sources
type VariableLoader interface {
	LoadFromFile(ctx context.Context, path string) ([]*spookyvariablestypes.Variable, error)
	LoadFromDirectory(ctx context.Context, dirPath string) ([]*spookyvariablestypes.Variable, error)
	LoadFromHCL(ctx context.Context, content []byte) ([]*spookyvariablestypes.Variable, error)
	LoadFromJSON(ctx context.Context, content []byte) ([]*spookyvariablestypes.Variable, error)
}

// VariableResolver defines the interface for resolving variable dependencies
type VariableResolver interface {
	ResolveVariable(ctx context.Context, variable *spookyvariablestypes.Variable, context *spookyvariablestypes.VariableContext) error
	ResolveDependencies(ctx context.Context, variables []*spookyvariablestypes.Variable) error
	ValidateDependencies(ctx context.Context, variables []*spookyvariablestypes.Variable) error
}

// VariableValidator defines the interface for variable validation
type VariableValidator interface {
	ValidateVariable(ctx context.Context, variable *spookyvariablestypes.Variable) (*spookyvariablestypes.ValidationResult, error)
	ValidateCollection(ctx context.Context, collection *spookyvariablestypes.VariableCollection) (*spookyvariablestypes.ValidationResult, error)
	ValidateContext(ctx context.Context, context *spookyvariablestypes.VariableContext) (*spookyvariablestypes.ValidationResult, error)
}

// ImportExportManager defines the interface for import/export operations
type ImportExportManager interface {
	ExportToHCL(ctx context.Context, variables []*spookyvariablestypes.Variable, path string) error
	ExportToJSON(ctx context.Context, variables []*spookyvariablestypes.Variable, path string) error
	ImportFromHCL(ctx context.Context, path string) ([]*spookyvariablestypes.Variable, error)
	ImportFromJSON(ctx context.Context, path string) ([]*spookyvariablestypes.Variable, error)
}
