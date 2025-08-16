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

// SaveVariables saves variables to the given destination
func (i *Integration) SaveVariables(ctx context.Context, variables map[string]*spookytypes.Variable, destination string) error {
	return i.manager.SaveVariables(ctx, variables, destination)
}

// EncryptVariables encrypts all variables that have encrypted=true
func (i *Integration) EncryptVariables(ctx context.Context, projectPath string, secretsIntegration spookyinterfaces.SecretsIntegration, recipients []string, dryRun bool) error {
	return i.manager.EncryptVariables(ctx, projectPath, secretsIntegration, recipients, dryRun)
}

// DecryptVariables decrypts age-encrypted values in variables for debugging
func (i *Integration) DecryptVariables(ctx context.Context, variables map[string]*spookytypes.Variable, secretsIntegration spookyinterfaces.SecretsIntegration, identityPath string) error {
	return i.manager.DecryptVariables(ctx, variables, secretsIntegration, identityPath)
}
