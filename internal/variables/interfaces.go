package variables

import (
	"context"
	"spooky/internal/variables/types"
)

// VariableManager defines the main interface for variable management
type VariableManager interface {
	// Core operations
	LoadVariables(ctx context.Context, path string) (*types.VariableCollection, error)
	GetVariable(ctx context.Context, name string) (*types.Variable, error)
	SetVariable(ctx context.Context, variable *types.Variable) error
	DeleteVariable(ctx context.Context, name string) error
	ListVariables(ctx context.Context) ([]*types.Variable, error)

	// Context operations
	CreateContext(ctx context.Context, variables []*types.Variable) (*types.VariableContext, error)
	ResolveContext(ctx context.Context, context *types.VariableContext) error

	// Import/export operations
	ExportVariables(ctx context.Context, format types.ExportFormat, path string) error
	ImportVariables(ctx context.Context, format types.ImportFormat, path string) error

	// Validation operations
	ValidateVariables(ctx context.Context, variables []*types.Variable) (*types.ValidationResult, error)
	ValidateContext(ctx context.Context, context *types.VariableContext) (*types.ValidationResult, error)

	// Coordinator integration methods
	LoadVariablesForProject(projectPath string) (*types.VariableCollection, error)
	ResolveVariablesForContext(context *types.VariableContext) error
	ValidateVariablesForProject(projectPath string) (*types.ValidationResult, error)
	ExportVariablesForProject(projectPath string, format types.ExportFormat, outputPath string) error
}

// VariableLoader defines the interface for loading variables from different sources
type VariableLoader interface {
	LoadFromFile(ctx context.Context, path string) ([]*types.Variable, error)
	LoadFromDirectory(ctx context.Context, dirPath string) ([]*types.Variable, error)
	LoadFromHCL(ctx context.Context, content []byte) ([]*types.Variable, error)
	LoadFromJSON(ctx context.Context, content []byte) ([]*types.Variable, error)
}

// VariableResolver defines the interface for resolving variable dependencies
type VariableResolver interface {
	ResolveVariable(ctx context.Context, variable *types.Variable, context *types.VariableContext) error
	ResolveDependencies(ctx context.Context, variables []*types.Variable) error
	ValidateDependencies(ctx context.Context, variables []*types.Variable) error
}

// VariableValidator defines the interface for variable validation
type VariableValidator interface {
	ValidateVariable(ctx context.Context, variable *types.Variable) (*types.ValidationResult, error)
	ValidateCollection(ctx context.Context, collection *types.VariableCollection) (*types.ValidationResult, error)
	ValidateContext(ctx context.Context, context *types.VariableContext) (*types.ValidationResult, error)
}

// ImportExportManager defines the interface for import/export operations
type ImportExportManager interface {
	ExportToHCL(ctx context.Context, variables []*types.Variable, path string) error
	ExportToJSON(ctx context.Context, variables []*types.Variable, path string) error
	ImportFromHCL(ctx context.Context, path string) ([]*types.Variable, error)
	ImportFromJSON(ctx context.Context, path string) ([]*types.Variable, error)
}
