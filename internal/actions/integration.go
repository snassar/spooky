package spookyactions

import (
	"context"

	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookyschemas "spooky/internal/schemas"
	spookytypes "spooky/internal/types"
	spookytypesactions "spooky/internal/types/actions"
)

// Integration implements the ActionsIntegration interface
type Integration struct {
	manager *Manager
}

// NewIntegration creates a new actions integration
func NewIntegration(
	logger spookylogging.Logger,
	validator spookyinterfaces.ActionValidator,
	sshManager spookyinterfaces.SSHManager,
	schemaValidator *spookyschemas.Validator,
) spookyinterfaces.ActionsIntegration {
	manager := NewManager(logger, validator, sshManager, schemaValidator)
	if managerImpl, ok := manager.(*Manager); ok {
		return &Integration{
			manager: managerImpl,
		}
	}
	// Fallback: create a new manager directly
	return &Integration{
		manager: &Manager{
			logger:          logger,
			validator:       validator,
			sshManager:      sshManager,
			schemaValidator: schemaValidator,
		},
	}
}

// LoadActions loads actions from the given source
func (i *Integration) LoadActions(ctx context.Context, source string) ([]spookytypes.Action, error) {
	return i.manager.LoadActions(ctx, source)
}

// ValidateActions validates actions
func (i *Integration) ValidateActions(ctx context.Context, actions []spookytypes.Action) (*spookytypes.ValidationResult, error) {
	return i.manager.ValidateActions(ctx, actions)
}

// RunActions runs actions on the given machines
func (i *Integration) RunActions(ctx context.Context, actions []spookytypes.Action, machines []spookytypes.Machine) ([]spookytypes.ActingResult, error) {
	return i.manager.RunActions(ctx, actions, machines)
}

// CreateActionPlan creates an execution plan for the given actions and machines
// This method is exposed for CLI use
func (i *Integration) CreateActionPlan(ctx context.Context, actions []spookytypes.Action, machines []spookytypes.Machine) (*spookytypesactions.ActionPlan, error) {
	return i.manager.createActionPlan(ctx, actions, machines)
}
