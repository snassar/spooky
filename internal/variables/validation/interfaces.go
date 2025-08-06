package validation

import (
	"context"
	"spooky/internal/variables/types"
)

// ValidationManager defines the interface for variable validation operations
type ValidationManager interface {
	// Core validation operations
	ValidateVariable(ctx context.Context, variable *types.Variable) (*types.ValidationResult, error)
	ValidateCollection(ctx context.Context, collection *types.VariableCollection) (*types.ValidationResult, error)
	ValidateContext(ctx context.Context, context *types.VariableContext) (*types.ValidationResult, error)

	// Schema validation
	ValidateAgainstSchema(variable *types.Variable, schemaName string) (*types.ValidationResult, error)
	ValidateCollectionAgainstSchema(collection *types.VariableCollection, schemaName string) (*types.ValidationResult, error)

	// Custom validation
	RegisterCustomValidator(name string, validator VariableValidator) error
	UnregisterCustomValidator(name string) error
	GetCustomValidators() []string

	// Configuration
	SetValidationRules(rules *types.ValidationRules) error
	EnableStrictValidation(strict bool) error
	SetMaxValidationErrors(max int) error

	// Utility operations
	GetValidationErrors() []types.ValidationError
	ClearValidationErrors() error
	Close() error
}

// VariableValidator defines the interface for custom validators
type VariableValidator interface {
	Validate(ctx context.Context, variable *types.Variable) error
	GetName() string
	GetDescription() string
}
