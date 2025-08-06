package interfaces

import (
	"spooky/internal/actions/types"
)

// ActionsIntegration defines the interface for actions system integration
type ActionsIntegration interface {
	// LoadActions loads actions from the project
	LoadActions(projectPath string) (*ActionsContext, error)

	// ValidateAction validates an action using the execution context
	ValidateAction(action *types.Action, context *ActionExecutionContext) error

	// PrepareActionForExecution prepares an action for execution
	PrepareActionForExecution(action *types.Action, context *ActionExecutionContext) error

	// ExecuteAction executes an action with the given context
	ExecuteAction(action *types.Action, context *ActionExecutionContext) error

	// GetAction gets a specific action by name
	GetAction(name string) (*types.Action, error)

	// ListActions lists all available actions
	ListActions() ([]*types.Action, error)

	// ListActionsFromProject lists actions from a specific project
	ListActionsFromProject(projectPath string) ([]*types.Action, error)

	// AddAction adds a new action to the project
	AddAction(name string, action *types.Action) error

	// RemoveAction removes an action from the project
	RemoveAction(name string) error
}

// ActionsContext provides actions data for integrations
type ActionsContext struct {
	BaseContext
	Actions map[string]*types.Action
}
