package actions

import (
	"context"
	"time"

	"spooky/internal/actions/types"
)

// ActionManager defines the main interface for action management and execution
type ActionManager interface {
	// Core action operations
	LoadActions(projectPath string) (*types.ActionCollection, error)
	GetAction(name string) (*types.Action, error)
	ListActions() ([]*types.Action, error)
	AddAction(name string, action *types.Action) error
	RemoveAction(name string) error

	// Acting operations
	ExecuteAction(ctx context.Context, action *types.Action, context *types.ActionContext) (*types.ActingSession, error)
	ExecuteActionCollection(ctx context.Context, collection *types.ActionCollection, context *types.ActionContext) (*types.ActingSession, error)
	PrepareAction(action *types.Action, context *types.ActionContext) error

	// Planning operations
	PlanAction(action *types.Action, context *types.ActionContext) (*types.ActionPlan, error)
	PlanActionCollection(collection *types.ActionCollection, context *types.ActionContext) (*types.ActionPlan, error)
	ValidatePlan(plan *types.ActionPlan) error

	// Validation operations
	ValidateAction(action *types.Action) error
	ValidateActionCollection(collection *types.ActionCollection) error
	ValidateActionContext(context *types.ActionContext) error

	// Merging operations
	MergeActions(actions ...*types.Action) (*types.ActionCollection, error)
	MergeWithPolicy(existing, new *types.ActionCollection, policy types.MergePolicy) (*types.ActionCollection, error)

	// Performance operations
	OptimizeAction(action *types.Action) error
	OptimizeActionCollection(collection *types.ActionCollection) error
	GetPerformanceMetrics(action *types.Action) (*types.PerformanceMetrics, error)

	// Configuration
	SetDefaultTimeout(timeout time.Duration)
	SetDefaultParallel(parallel bool)
	RegisterCustomValidator(name string, validator ActionValidator)

	// Utility operations
	Close() error
}

// ActionValidator defines the interface for action validation
type ActionValidator interface {
	// Core validation operations
	Validate(action *types.Action) error
	ValidateCollection(collection *types.ActionCollection) error
	ValidateContext(context *types.ActionContext) error
}
