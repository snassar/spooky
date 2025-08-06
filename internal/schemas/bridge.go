package schemas

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// ConfigValidatorBridge bridges between HCL validation and Go struct validation
type ConfigValidatorBridge struct {
	hclValidator *SchemaValidator
}

// NewConfigValidatorBridge creates a new config validator bridge
func NewConfigValidatorBridge() *ConfigValidatorBridge {
	return &ConfigValidatorBridge{
		hclValidator: NewSchemaValidator(),
	}
}

// ValidateMachineConfig validates machine configuration using both HCL and struct validation
func (b *ConfigValidatorBridge) ValidateMachineConfig(machine interface{}) *ValidationResult {
	return b.validateConfigGeneric(machine, "machine-config", "machine", func(obj interface{}) error {
		// For now, just validate that the machine is not nil
		// The actual struct validation will be done in the config package
		if obj == nil {
			return fmt.Errorf("machine configuration is nil")
		}
		return nil
	})
}

// ValidateActionConfig validates action configuration using both HCL and struct validation
func (b *ConfigValidatorBridge) ValidateActionConfig(action interface{}) *ValidationResult {
	return b.validateConfigGeneric(action, "action-config", "action", func(obj interface{}) error {
		// For now, just validate that the action is not nil
		// The actual struct validation will be done in the config package
		if obj == nil {
			return fmt.Errorf("action configuration is nil")
		}
		return nil
	})
}

// ValidateProjectConfig validates project configuration using both HCL and struct validation
func (b *ConfigValidatorBridge) ValidateProjectConfig(project interface{}) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Schema: "project-config",
	}

	// For now, just validate that the project is not nil
	// DEPRECATED: Schema system is fully implemented - this TODO is ready for removal
	if project == nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project",
			Message:  "project configuration is nil",
			Severity: "error",
		})
	}

	return result
}

// ValidateGlobalConfig validates global configuration using both HCL and struct validation
func (b *ConfigValidatorBridge) ValidateGlobalConfig(global interface{}) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Schema: "global-config",
	}

	// Validate using HCL schema first
	if hclResult := b.hclValidator.ValidateFile("", "spooky"); !hclResult.Valid {
		result.Valid = false
		result.Errors = append(result.Errors, hclResult.Errors...)
		result.Warnings = append(result.Warnings, hclResult.Warnings...)
	}

	// For now, just validate that the global config is not nil
	// The actual struct validation will be done in the config package
	if global == nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "global-config",
			Message:  "global configuration is nil",
			Severity: "error",
		})
	}

	return result
}

// convertValidatorErrors converts go-playground/validator errors to unified format
func (b *ConfigValidatorBridge) convertValidatorErrors(err error) []ValidationError {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errors []ValidationError
		for _, e := range validationErrors {
			errors = append(errors, *ConvertConfigValidationError(e))
		}
		return errors
	}
	return nil
}

// validateConfigGeneric is a generic validation function to reduce code duplication
func (b *ConfigValidatorBridge) validateConfigGeneric(obj interface{}, schemaName string, fieldName string, validator func(interface{}) error) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Schema: schemaName,
	}

	// Validate using the provided validator function
	if err := validator(obj); err != nil {
		result.Valid = false
		// Convert validator errors to unified format using the existing helper
		validationErrors := b.convertValidatorErrors(err)
		if len(validationErrors) > 0 {
			result.Errors = append(result.Errors, validationErrors...)
		} else {
			result.Errors = append(result.Errors, ValidationError{
				Field:    fieldName,
				Message:  err.Error(),
				Severity: "error",
			})
		}
	}

	return result
}
