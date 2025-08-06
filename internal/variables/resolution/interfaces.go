package resolution

import (
	"context"
	"spooky/internal/variables/types"
)

// ResolutionManager defines the interface for variable resolution operations
type ResolutionManager interface {
	// Core resolution operations
	ResolveVariable(ctx context.Context, variable *types.Variable, context *types.VariableContext) error
	ResolveDependencies(ctx context.Context, variables []*types.Variable) error
	ResolveContext(ctx context.Context, context *types.VariableContext) error

	// Dependency management
	ValidateDependencies(ctx context.Context, variables []*types.Variable) error
	DetectCircularDependencies(variables []*types.Variable) error
	GetDependencyGraph(variables []*types.Variable) (*types.DependencyGraph, error)

	// Configuration
	SetMaxRecursionDepth(depth int) error
	SetDefaultValues(defaults map[string]interface{}) error
	EnableStrictMode(strict bool) error

	// Utility operations
	GetUnresolvedVariables(variables []*types.Variable) []*types.Variable
	GetResolutionOrder(variables []*types.Variable) ([]*types.Variable, error)
	Close() error
}

// VariableResolver defines the interface for specific resolution strategies
type VariableResolver interface {
	Resolve(ctx context.Context, variable *types.Variable, context *types.VariableContext) error
	GetName() string
	GetSupportedTypes() []string
}
