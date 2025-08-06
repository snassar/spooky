package validation

import (
	"fmt"

	"spooky/internal/actions/types"
	"spooky/internal/logging"
)

// ActionValidatorImpl implements the ActionValidator interface
type ActionValidatorImpl struct {
	action          *types.Action
	strictMode      bool
	validationLevel types.ValidationLevel
	rules           []ValidationRule
	logger          logging.Logger
}

// NewActionValidator creates a new ActionValidator
func NewActionValidator(action *types.Action, logger logging.Logger) ActionValidator {
	return &ActionValidatorImpl{
		action:          action,
		strictMode:      false,
		validationLevel: types.ValidationLevelBasic,
		rules:           make([]ValidationRule, 0),
		logger:          logger,
	}
}

// Validate validates an action
func (v *ActionValidatorImpl) Validate(action *types.Action) error {
	if action == nil {
		return fmt.Errorf("action cannot be nil")
	}

	v.logger.Info("Validating action", logging.String("action", action.Name))

	// Basic action validation
	if action.Name == "" {
		return fmt.Errorf("action name cannot be empty")
	}

	if action.Command == "" {
		return fmt.Errorf("action command cannot be empty")
	}

	// Validate action configuration
	if err := v.validateActionConfig(action); err != nil {
		return fmt.Errorf("action configuration validation failed: %w", err)
	}

	// Run custom validation rules
	for _, rule := range v.rules {
		if err := rule.Validate(action); err != nil {
			if v.strictMode {
				return fmt.Errorf("custom rule '%s' validation failed: %w", rule.Name(), err)
			}
			v.logger.Warn("Custom validation rule failed",
				logging.String("rule", rule.Name()),
				logging.String("error", err.Error()))
		}
	}

	v.logger.Info("Action validation successful", logging.String("action", action.Name))
	return nil
}

// ValidateCollection validates an action collection
func (v *ActionValidatorImpl) ValidateCollection(collection *types.ActionCollection) error {
	if collection == nil {
		return fmt.Errorf("action collection cannot be nil")
	}

	v.logger.Info("Validating action collection", logging.Int("actions_count", len(collection.Actions)))

	// Validate each action in the collection
	for _, action := range collection.Actions {
		if err := v.Validate(action); err != nil {
			return fmt.Errorf("failed to validate action %s in collection: %w", action.Name, err)
		}
	}

	v.logger.Info("Action collection validation successful", logging.Int("actions_count", len(collection.Actions)))
	return nil
}

// ValidateContext validates an action context
func (v *ActionValidatorImpl) ValidateContext(context *types.ActionContext) error {
	if context == nil {
		return fmt.Errorf("action context cannot be nil")
	}

	v.logger.Info("Validating action context")

	// Basic context validation
	if context.ProjectPath == "" {
		return fmt.Errorf("project path cannot be empty")
	}

	// Validate machines if specified
	if len(context.Machines) > 0 {
		for i, machine := range context.Machines {
			if machine.Name == "" {
				return fmt.Errorf("machine %d: machine name cannot be empty", i)
			}
		}
	}

	v.logger.Info("Action context validation successful")
	return nil
}

// AddRule adds a custom validation rule
func (v *ActionValidatorImpl) AddRule(rule ValidationRule) error {
	if rule == nil {
		return fmt.Errorf("validation rule cannot be nil")
	}

	v.rules = append(v.rules, rule)
	v.logger.Debug("Added validation rule", logging.String("rule", rule.Name()))
	return nil
}

// RemoveRule removes a validation rule by name
func (v *ActionValidatorImpl) RemoveRule(name string) error {
	for i, rule := range v.rules {
		if rule.Name() == name {
			v.rules = append(v.rules[:i], v.rules[i+1:]...)
			v.logger.Debug("Removed validation rule", logging.String("rule", name))
			return nil
		}
	}

	return fmt.Errorf("validation rule not found: %s", name)
}

// GetRules returns all validation rules
func (v *ActionValidatorImpl) GetRules() ([]ValidationRule, error) {
	rules := make([]ValidationRule, len(v.rules))
	copy(rules, v.rules)
	return rules, nil
}

// SetStrictMode sets the strict validation mode
func (v *ActionValidatorImpl) SetStrictMode(strict bool) {
	v.strictMode = strict
	v.logger.Debug("Set strict validation mode", logging.Bool("strict", strict))
}

// SetLevel sets the validation level
func (v *ActionValidatorImpl) SetLevel(level types.ValidationLevel) {
	v.validationLevel = level
	v.logger.Debug("Set validation level", logging.String("level", string(level)))
}

// validateActionConfig validates the action configuration
func (v *ActionValidatorImpl) validateActionConfig(action *types.Action) error {
	// Validate timeout if specified
	if action.Timeout < 0 {
		return fmt.Errorf("action timeout cannot be negative")
	}

	// Validate retries if specified
	if action.Retries < 0 {
		return fmt.Errorf("action retries cannot be negative")
	}

	// Validate retry delay if specified
	if action.RetryDelay < 0 {
		return fmt.Errorf("action retry delay cannot be negative")
	}

	// Validate environment variables if specified
	if action.Environment != nil {
		for key, value := range action.Environment {
			if key == "" {
				return fmt.Errorf("environment variable key cannot be empty")
			}
			if value == "" {
				return fmt.Errorf("environment variable value cannot be empty for key: %s", key)
			}
		}
	}

	return nil
}
