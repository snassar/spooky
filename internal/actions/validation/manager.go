package validation

import (
	"fmt"
	"sync"

	spookytypes "spooky/internal/types"
	spookylogging "spooky/internal/logging"
)

// Manager implements the ValidationManager interface
type Manager struct {
	// Configuration
	strictMode      bool
	validationLevel spookyactionstypes.ValidationLevel

	// State
	validators map[string]ActionValidator
	rules      []ValidationRule
	logger     spookylogging.Logger
	mu         sync.RWMutex
}

// NewManager creates a new ValidationManager
func NewManager(logger spookylogging.Logger) *Manager {
	return &Manager{
		strictMode:      false,
		validationLevel: spookyactionstypes.ValidationLevelBasic,
		validators:      make(map[string]ActionValidator),
		rules:           make([]ValidationRule, 0),
		logger:          logger,
	}
}

// ValidateAction validates a single action
func (m *Manager) ValidateAction(action *spookyactionstypes.Action) error {
	if action == nil {
		return fmt.Errorf("action cannot be nil")
	}

	m.logger.Info("Validating action", spookylogging.String("action", action.Name))

	// Create a validator for this action
	validator, err := m.CreateValidator(action)
	if err != nil {
		return fmt.Errorf("failed to create validator for action %s: %w", action.Name, err)
	}

	// Validate the action
	if err := validator.Validate(action); err != nil {
		return fmt.Errorf("action validation failed: %w", err)
	}

	m.logger.Info("Action validation successful", spookylogging.String("action", action.Name))
	return nil
}

// ValidateActionCollection validates a collection of actions
func (m *Manager) ValidateActionCollection(collection *spookyactionstypes.ActionCollection) error {
	if collection == nil {
		return fmt.Errorf("action collection cannot be nil")
	}

	m.logger.Info("Validating action collection", spookylogging.Int("actions_count", len(collection.Actions)))

	// Validate each action in the collection
	for _, action := range collection.Actions {
		if err := m.ValidateAction(action); err != nil {
			return fmt.Errorf("failed to validate action %s in collection: %w", action.Name, err)
		}
	}

	m.logger.Info("Action collection validation successful", spookylogging.Int("actions_count", len(collection.Actions)))
	return nil
}

// ValidateActionContext validates an action context
func (m *Manager) ValidateActionContext(context *spookyactionstypes.ActionContext) error {
	if context == nil {
		return fmt.Errorf("action context cannot be nil")
	}

	m.logger.Info("Validating action context")

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

	// Validate variables if specified
	if context.Variables != nil {
		// Variables context validation is handled by the variables package
		// We just ensure it's not nil
	}

	m.logger.Info("Action context validation successful")
	return nil
}

// CreateValidator creates a new validator for an action
func (m *Manager) CreateValidator(action *spookyactionstypes.Action) (ActionValidator, error) {
	if action == nil {
		return nil, fmt.Errorf("action cannot be nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if validator already exists
	if validator, exists := m.validators[action.Name]; exists {
		return validator, nil
	}

	// Create new validator
	validator := NewActionValidator(action, m.logger)
	m.validators[action.Name] = validator

	m.logger.Debug("Created validator for action", spookylogging.String("action", action.Name))
	return validator, nil
}

// GetValidator gets an existing validator for an action
func (m *Manager) GetValidator(action *spookyactionstypes.Action) (ActionValidator, error) {
	if action == nil {
		return nil, fmt.Errorf("action cannot be nil")
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	validator, exists := m.validators[action.Name]
	if !exists {
		return nil, fmt.Errorf("validator not found for action %s", action.Name)
	}

	return validator, nil
}

// AddValidationRule adds a custom validation rule
func (m *Manager) AddValidationRule(rule ValidationRule) error {
	if rule == nil {
		return fmt.Errorf("validation rule cannot be nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.rules = append(m.rules, rule)
	m.logger.Info("Added validation rule", spookylogging.String("rule", rule.Name()))
	return nil
}

// RemoveValidationRule removes a validation rule by name
func (m *Manager) RemoveValidationRule(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, rule := range m.rules {
		if rule.Name() == name {
			m.rules = append(m.rules[:i], m.rules[i+1:]...)
			m.logger.Info("Removed validation rule", spookylogging.String("rule", name))
			return nil
		}
	}

	return fmt.Errorf("validation rule not found: %s", name)
}

// GetValidationRules returns all validation rules
func (m *Manager) GetValidationRules() ([]ValidationRule, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	rules := make([]ValidationRule, len(m.rules))
	copy(rules, m.rules)
	return rules, nil
}

// SetStrictMode sets the strict validation mode
func (m *Manager) SetStrictMode(strict bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.strictMode = strict
	m.logger.Info("Set strict validation mode", spookylogging.Bool("strict", strict))
}

// SetValidationLevel sets the validation level
func (m *Manager) SetValidationLevel(level spookyactionstypes.ValidationLevel) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.validationLevel = level
	m.logger.Info("Set validation level", spookylogging.String("level", string(level)))
}
