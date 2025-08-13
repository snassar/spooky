package variables

import (
	"context"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
)

// Integration implements the VariablesIntegration interface
type Integration struct {
	manager spookyinterfaces.VariablesIntegration
}

// NewIntegration creates a new variables integration
func NewIntegration(
	logger spookytypeslogging.Logger,
	loader spookyinterfaces.VariableLoader,
	validator spookyinterfaces.VariableValidator,
) spookyinterfaces.VariablesIntegration {
	manager := NewManager(logger, loader, validator)
	return &Integration{
		manager: manager,
	}
}

// LoadVariables loads variables from the given source
func (i *Integration) LoadVariables(ctx context.Context, source string) (map[string]*spookytypes.Variable, error) {
	return i.manager.LoadVariables(ctx, source)
}

// ResolveVariables resolves variables with the given context
func (i *Integration) ResolveVariables(ctx context.Context, variables map[string]*spookytypes.Variable, context *spookytypes.VariableContext) (*spookytypes.VariableResolutionResult, error) {
	return i.manager.ResolveVariables(ctx, variables, context)
}

// ValidateVariables validates variables
func (i *Integration) ValidateVariables(ctx context.Context, variables map[string]*spookytypes.Variable) (*spookytypes.ValidationResult, error) {
	return i.manager.ValidateVariables(ctx, variables)
}
