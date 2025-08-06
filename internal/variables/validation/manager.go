package validation

import (
	"context"
	"fmt"
	"strings"

	"spooky/internal/logging"
	"spooky/internal/schemas"
	"spooky/internal/variables/types"
)

// Manager implements ValidationManager interface
type Manager struct {
	config           *types.ValidationConfig
	customValidators map[string]VariableValidator
	logger           logging.Logger
	errors           []types.ValidationError
}

// NewManager creates a new validation manager
func NewManager(config *types.ValidationConfig, logger logging.Logger) *Manager {
	return &Manager{
		config:           config,
		customValidators: make(map[string]VariableValidator),
		logger:           logger,
		errors:           make([]types.ValidationError, 0),
	}
}

// ValidateVariable validates a single variable
func (m *Manager) ValidateVariable(ctx context.Context, variable *types.Variable) (*types.ValidationResult, error) {
	result := &types.ValidationResult{
		Valid:    true,
		Errors:   make([]types.ValidationError, 0),
		Warnings: make([]types.ValidationWarning, 0),
	}

	// 1. Basic validation
	if err := m.validateBasic(variable); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, types.ValidationError{
			Field:   "basic",
			Message: err.Error(),
		})
	}

	// 2. Schema validation
	if err := m.validateAgainstSchema(variable, "variables-hcl"); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, types.ValidationError{
			Field:   "schema",
			Message: err.Error(),
		})
	}

	// 3. Custom validation
	for name, validator := range m.customValidators {
		if err := validator.Validate(ctx, variable); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, types.ValidationError{
				Field:   name,
				Message: err.Error(),
			})
		}
	}

	// 4. Store errors
	m.errors = append(m.errors, result.Errors...)

	return result, nil
}

// ValidateCollection validates a collection of variables
func (m *Manager) ValidateCollection(ctx context.Context, collection *types.VariableCollection) (*types.ValidationResult, error) {
	result := &types.ValidationResult{
		Valid:    true,
		Errors:   make([]types.ValidationError, 0),
		Warnings: make([]types.ValidationWarning, 0),
	}

	// Validate each variable in collection
	for _, variable := range collection.Variables {
		varResult, err := m.ValidateVariable(ctx, variable)
		if err != nil {
			return nil, fmt.Errorf("failed to validate variable %s: %w", variable.Name, err)
		}

		if !varResult.Valid {
			result.Valid = false
			result.Errors = append(result.Errors, varResult.Errors...)
		}

		result.Warnings = append(result.Warnings, varResult.Warnings...)
	}

	// Validate collection-level constraints
	if err := m.validateCollectionConstraints(collection); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, types.ValidationError{
			Field:   "collection",
			Message: err.Error(),
		})
	}

	return result, nil
}

// ValidateContext validates a variable context
func (m *Manager) ValidateContext(ctx context.Context, context *types.VariableContext) (*types.ValidationResult, error) {
	// Convert context variables to collection for validation
	variables := make([]*types.Variable, 0, len(context.Variables))
	for _, variable := range context.Variables {
		variables = append(variables, variable)
	}

	collection := &types.VariableCollection{
		Variables: variables,
	}

	return m.ValidateCollection(ctx, collection)
}

// ValidateAgainstSchema validates a variable against a schema
func (m *Manager) ValidateAgainstSchema(variable *types.Variable, schemaName string) (*types.ValidationResult, error) {
	result := &types.ValidationResult{
		Valid:    true,
		Errors:   make([]types.ValidationError, 0),
		Warnings: make([]types.ValidationWarning, 0),
	}

	// Load schema if not already loaded
	validator := schemas.NewSchemaValidator()
	if err := validator.LoadSchema(schemas.SchemaTypeVariablesHCL); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, types.ValidationError{
			Field:   "schema",
			Message: fmt.Sprintf("failed to load schema: %v", err),
		})
		return result, nil
	}

	// Validate variable against schema
	// Implementation depends on schema validation interface
	// For now, return success
	return result, nil
}

// ValidateCollectionAgainstSchema validates a collection against a schema
func (m *Manager) ValidateCollectionAgainstSchema(collection *types.VariableCollection, schemaName string) (*types.ValidationResult, error) {
	result := &types.ValidationResult{
		Valid:    true,
		Errors:   make([]types.ValidationError, 0),
		Warnings: make([]types.ValidationWarning, 0),
	}

	// Validate each variable against schema
	for _, variable := range collection.Variables {
		varResult, err := m.ValidateAgainstSchema(variable, schemaName)
		if err != nil {
			return nil, fmt.Errorf("failed to validate variable %s against schema: %w", variable.Name, err)
		}

		if !varResult.Valid {
			result.Valid = false
			result.Errors = append(result.Errors, varResult.Errors...)
		}

		result.Warnings = append(result.Warnings, varResult.Warnings...)
	}

	return result, nil
}

// RegisterCustomValidator registers a custom validator
func (m *Manager) RegisterCustomValidator(name string, validator VariableValidator) error {
	if name == "" {
		return fmt.Errorf("validator name cannot be empty")
	}

	if validator == nil {
		return fmt.Errorf("validator cannot be nil")
	}

	m.customValidators[name] = validator
	return nil
}

// UnregisterCustomValidator unregisters a custom validator
func (m *Manager) UnregisterCustomValidator(name string) error {
	if _, exists := m.customValidators[name]; !exists {
		return fmt.Errorf("validator %s not found", name)
	}

	delete(m.customValidators, name)
	return nil
}

// GetCustomValidators returns the list of custom validators
func (m *Manager) GetCustomValidators() []string {
	validators := make([]string, 0, len(m.customValidators))
	for name := range m.customValidators {
		validators = append(validators, name)
	}
	return validators
}

// SetValidationRules sets validation rules
func (m *Manager) SetValidationRules(rules *types.ValidationRules) error {
	if m.config == nil {
		m.config = &types.ValidationConfig{}
	}
	m.config.ValidationRules = rules
	return nil
}

// EnableStrictValidation enables or disables strict validation
func (m *Manager) EnableStrictValidation(strict bool) error {
	if m.config == nil {
		m.config = &types.ValidationConfig{}
	}
	m.config.StrictValidation = strict
	return nil
}

// SetMaxValidationErrors sets the maximum number of validation errors
func (m *Manager) SetMaxValidationErrors(max int) error {
	if m.config == nil {
		m.config = &types.ValidationConfig{}
	}
	m.config.MaxValidationErrors = max
	return nil
}

// GetValidationErrors returns all validation errors
func (m *Manager) GetValidationErrors() []types.ValidationError {
	return m.errors
}

// ClearValidationErrors clears all validation errors
func (m *Manager) ClearValidationErrors() error {
	m.errors = make([]types.ValidationError, 0)
	return nil
}

// Close closes the validation manager and releases resources
func (m *Manager) Close() error {
	// Clean up resources if needed
	return nil
}

// Helper methods
func (m *Manager) validateBasic(variable *types.Variable) error {
	if variable.Name == "" {
		return fmt.Errorf("variable name cannot be empty")
	}

	if variable.Type == "" {
		return fmt.Errorf("variable type cannot be empty")
	}

	if strings.Contains(variable.Name, " ") {
		return fmt.Errorf("variable name cannot contain spaces")
	}

	return nil
}

func (m *Manager) validateAgainstSchema(variable *types.Variable, schemaName string) error {
	// Load schema if not already loaded
	validator := schemas.NewSchemaValidator()
	if err := validator.LoadSchema(schemas.SchemaTypeVariablesHCL); err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	// Validate variable against schema
	// Implementation depends on schema validation interface
	return nil
}

func (m *Manager) validateCollectionConstraints(collection *types.VariableCollection) error {
	// Check for duplicate variable names
	names := make(map[string]bool)
	for _, variable := range collection.Variables {
		if names[variable.Name] {
			return fmt.Errorf("duplicate variable name: %s", variable.Name)
		}
		names[variable.Name] = true
	}

	return nil
}
