package actions

import (
	"context"
	"time"

	spookyactionstypes "spooky/internal/actions/types"
)

// ActionManager defines the main interface for action management and execution
type ActionManager interface {
	// Core action operations
	LoadActions(projectPath string) (*spookyactionstypes.ActionCollection, error)
	GetAction(name string) (*spookyactionstypes.Action, error)
	ListActions() ([]*spookyactionstypes.Action, error)
	AddAction(name string, action *spookyactionstypes.Action) error
	RemoveAction(name string) error

	// Acting operations
	ExecuteAction(ctx context.Context, action *spookyactionstypes.Action, context *spookyactionstypes.ActionContext) (*spookyactionstypes.ActingSession, error)
	ExecuteActionCollection(ctx context.Context, collection *spookyactionstypes.ActionCollection, context *spookyactionstypes.ActionContext) (*spookyactionstypes.ActingSession, error)
	PrepareAction(action *spookyactionstypes.Action, context *spookyactionstypes.ActionContext) error

	// Planning operations
	PlanAction(action *spookyactionstypes.Action, context *spookyactionstypes.ActionContext) (*spookyactionstypes.ActionPlan, error)
	PlanActionCollection(collection *spookyactionstypes.ActionCollection, context *spookyactionstypes.ActionContext) (*spookyactionstypes.ActionPlan, error)
	ValidatePlan(plan *spookyactionstypes.ActionPlan) error

	// Validation operations
	ValidateAction(action *spookyactionstypes.Action) error
	ValidateActionCollection(collection *spookyactionstypes.ActionCollection) error
	ValidateActionContext(context *spookyactionstypes.ActionContext) error

	// Merging operations
	MergeActions(actions ...*spookyactionstypes.Action) (*spookyactionstypes.ActionCollection, error)
	MergeWithPolicy(existing, new *spookyactionstypes.ActionCollection, policy spookyactionstypes.MergePolicy) (*spookyactionstypes.ActionCollection, error)

	// Performance operations
	OptimizeAction(action *spookyactionstypes.Action) error
	OptimizeActionCollection(collection *spookyactionstypes.ActionCollection) error
	GetPerformanceMetrics(action *spookyactionstypes.Action) (*spookyactionstypes.PerformanceMetrics, error)

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
	Validate(action *spookyactionstypes.Action) error
	ValidateCollection(collection *spookyactionstypes.ActionCollection) error
	ValidateContext(context *spookyactionstypes.ActionContext) error
}
