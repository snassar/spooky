package schemas

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/hcl/v2"
)

// ConvertValidationError converts between different ValidationError types
func ConvertValidationError(err ValidationError) *ValidationError {
	return &err
}

// ConvertConfigValidationError converts go-playground/validator errors to unified format
func ConvertConfigValidationError(err validator.FieldError) *ValidationError {
	return &ValidationError{
		Field:    err.Field(),
		Message:  formatConfigValidationError(err),
		Value:    fmt.Sprintf("%v", err.Value()),
		Severity: "error",
	}
}

// ConvertFactsValidationError converts facts validation errors to unified format
// This is a placeholder since facts validation is now handled in the facts package
func ConvertFactsValidationError(err error) *ValidationError {
	return &ValidationError{
		Field:    "facts",
		Message:  err.Error(),
		Severity: "error",
	}
}

// ConvertHCLDiagnostic converts HCL diagnostics to unified format
func ConvertHCLDiagnostic(diag hcl.Diagnostic) *ValidationError {
	        validationError := &ValidationError{
		Message:  diag.Summary,
		Severity: "error",
	}

	// Extract line and column information if available
	if diag.Subject != nil {
		validationError.Line = diag.Subject.Start.Line
		validationError.Column = diag.Subject.Start.Column
	}

	// Extract field information from context
	        if diag.Subject != nil && diag.Subject.Filename != "" {
		validationError.File = diag.Subject.Filename
	}

	return validationError
}

// ConvertValidationResult converts between different ValidationResult types
func ConvertValidationResult(result *ValidationResult) *ValidationResult {
	return result
}

// formatConfigValidationError formats go-playground/validator errors
func formatConfigValidationError(e validator.FieldError) string {
	// Handle special cases for min validation
	if e.Tag() == "min" {
		return formatMinValidation(e)
	}

	// Use map for other validation tags
	errorMessages := map[string]string{
		"required":       fmt.Sprintf("%s is required", e.Field()),
		"max":            fmt.Sprintf("%s must be at most %s", e.Field(), e.Param()),
		"machine_auth":   fmt.Sprintf("either password or key_file must be specified for machine %s", e.Param()),
		"action_exec":    fmt.Sprintf("either command or script must be specified for action %s (but not both)", e.Param()),
		"unique_machine": fmt.Sprintf("duplicate machine name: %s", e.Param()),
		"unique_action":  fmt.Sprintf("duplicate action name: %s", e.Param()),
		"valid_port":     fmt.Sprintf("port must be between 1 and 65535 for machine %s", e.Param()),
		"valid_timeout":  fmt.Sprintf("timeout must be between 1 and 3600 seconds for action %s", e.Param()),
		"valid_machines": fmt.Sprintf("machine reference '%s' in action '%s' does not exist", e.Value(), e.Param()),
		"sshkeyfile":     fmt.Sprintf("SSH key file '%s' does not exist or is not readable for machine %s", e.Value(), e.Param()),
		"scriptfile":     fmt.Sprintf("script file '%s' does not exist or is not executable for action %s", e.Value(), e.Param()),
	}

	if message, exists := errorMessages[e.Tag()]; exists {
		return message
	}

	return fmt.Sprintf("%s failed validation: %s", e.Field(), e.Tag())
}

// formatMinValidation handles the complex min validation logic
func formatMinValidation(e validator.FieldError) string {
	if e.Field() == "Machines" {
		return "at least one machine must be defined"
	}

	// Handle numeric fields differently
	if e.Field() == "Port" || e.Field() == "Timeout" {
		return fmt.Sprintf("%s must be at least %s", e.Field(), e.Param())
	}

	return fmt.Sprintf("%s must be at least %s", e.Field(), e.Param())
}

// ConvertValidatorErrors converts a slice of validator errors to unified format
func ConvertValidatorErrors(err error) []ValidationError {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errors []ValidationError
		for _, e := range validationErrors {
			errors = append(errors, *ConvertConfigValidationError(e))
		}
		return errors
	}
	return nil
}

// ConvertFactsValidationResult converts facts validation results to unified format
// This is a placeholder since facts validation is now handled in the facts package
func ConvertFactsValidationResult(_ interface{}) *ValidationResult {
	return &ValidationResult{
		Valid: true,
	}
}

// AddFileInfo adds file information to validation errors
func AddFileInfo(errors []ValidationError, filePath string) []ValidationError {
	for i := range errors {
		if errors[i].File == "" {
			errors[i].File = filePath
		}
	}
	return errors
}

// AddSchemaInfo adds schema information to validation results
func AddSchemaInfo(result *ValidationResult, schemaName string) *ValidationResult {
	if result != nil {
		result.Schema = schemaName
	}
	return result
}
