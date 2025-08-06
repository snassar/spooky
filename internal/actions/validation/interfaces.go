package validation

import (
	"spooky/internal/actions/types"
)

// ValidationManager defines the interface for action validation operations
type ValidationManager interface {
	// Core validation operations
	ValidateAction(action *types.Action) error
	ValidateActionCollection(collection *types.ActionCollection) error
	ValidateActionContext(context *types.ActionContext) error

	// Validator management
	CreateValidator(action *types.Action) (ActionValidator, error)
	GetValidator(action *types.Action) (ActionValidator, error)

	// Custom validation
	AddValidationRule(rule ValidationRule) error
	RemoveValidationRule(name string) error
	GetValidationRules() ([]ValidationRule, error)

	// Configuration
	SetStrictMode(strict bool)
	SetValidationLevel(level types.ValidationLevel)
}

// ActionValidator defines the interface for action validation
type ActionValidator interface {
	// Core validation operations
	Validate(action *types.Action) error
	ValidateCollection(collection *types.ActionCollection) error
	ValidateContext(context *types.ActionContext) error

	// Rule management
	AddRule(rule ValidationRule) error
	RemoveRule(name string) error
	GetRules() ([]ValidationRule, error)

	// Configuration
	SetStrictMode(strict bool)
	SetLevel(level types.ValidationLevel)
}

// ValidationRule defines a custom validation rule
type ValidationRule interface {
	Name() string
	Validate(action *types.Action) error
	ValidateCollection(collection *types.ActionCollection) error
	ValidateContext(context *types.ActionContext) error
}

// ValidationResult represents the result of a validation operation
type ValidationResult struct {
	Valid    bool                   `json:"valid"`
	Errors   []ValidationError      `json:"errors"`
	Warnings []ValidationWarning    `json:"warnings"`
	Details  map[string]interface{} `json:"details"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field    string `json:"field"`
	Message  string `json:"message"`
	Code     string `json:"code"`
	Severity string `json:"severity"`
}

// ValidationWarning represents a validation warning
type ValidationWarning struct {
	Field    string `json:"field"`
	Message  string `json:"message"`
	Code     string `json:"code"`
	Severity string `json:"severity"`
}
