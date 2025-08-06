package validators

import (
	"fmt"
	spookyfactstypes "spooky/internal/facts/types"
)

// Validator defines the interface for fact validation operations
type Validator interface {
	// Validation operations
	ValidateCollection(collection *spookyfactstypes.FactCollection) error
	ValidateFact(fact *spookyfactstypes.Fact) error
	ValidateCustomFacts(facts map[string]*spookyfactstypes.CustomFacts) error

	// Schema validation
	ValidateAgainstSchema(collection *spookyfactstypes.FactCollection, schema interface{}) error

	// Custom validation
	AddValidationRule(rule ValidationRule) error
	RemoveValidationRule(name string) error
}

// ValidationRule defines a custom validation rule
type ValidationRule interface {
	Name() string
	Validate(fact *spookyfactstypes.Fact) error
	ValidateCollection(collection *spookyfactstypes.FactCollection) error
}

// Manager provides validator management functionality
type Manager struct {
	validator Validator
	rules     map[string]ValidationRule
}

// NewManager creates a new validator manager
func NewManager(validator Validator) *Manager {
	return &Manager{
		validator: validator,
		rules:     make(map[string]ValidationRule),
	}
}

// GetValidator returns the underlying validator
func (m *Manager) GetValidator() Validator {
	return m.validator
}

// SetValidator sets the underlying validator
func (m *Manager) SetValidator(validator Validator) {
	m.validator = validator
}

// AddRule adds a validation rule
func (m *Manager) AddRule(rule ValidationRule) error {
	if rule == nil {
		return fmt.Errorf("validation rule cannot be nil")
	}
	m.rules[rule.Name()] = rule
	return nil
}

// RemoveRule removes a validation rule
func (m *Manager) RemoveRule(name string) error {
	if _, exists := m.rules[name]; !exists {
		return fmt.Errorf("validation rule %s not found", name)
	}
	delete(m.rules, name)
	return nil
}

// GetRule gets a validation rule by name
func (m *Manager) GetRule(name string) (ValidationRule, bool) {
	rule, exists := m.rules[name]
	return rule, exists
}
