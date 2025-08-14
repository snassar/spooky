// Package schemas provides enhanced schema validation and management functionality for the spooky codebase.
package schemas

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsimple"

	spookytypeslogging "spooky/internal/types/logging"
	spookytypesschemas "spooky/internal/types/schemas"
)

// EnhancedValidator provides comprehensive schema validation with advanced features
type EnhancedValidator struct {
	logger   spookytypeslogging.Logger
	registry spookytypesschemas.SchemaRegistry
	parser   *hclparse.Parser

	// Validation configuration
	config *ValidationConfig

	// Custom validation functions
	customValidators map[string]CustomValidatorFunc
}

// ValidationConfig provides configuration for enhanced validation
type ValidationConfig struct {
	// Validation mode
	Mode ValidationMode `json:"mode" hcl:"mode"`

	// Error handling configuration
	ErrorHandling *ErrorHandlingConfig `json:"error_handling" hcl:"error_handling"`

	// Schema evolution configuration
	Evolution *EvolutionConfig `json:"evolution" hcl:"evolution"`
}

// ValidationMode represents the validation mode
type ValidationMode string

const (
	ValidationModeStrict     ValidationMode = "strict"
	ValidationModeLenient    ValidationMode = "lenient"
	ValidationModePermissive ValidationMode = "permissive"
)

// ErrorHandlingConfig provides error handling configuration
type ErrorHandlingConfig struct {
	// Whether to stop on first error
	StopOnFirstError bool `json:"stop_on_first_error" hcl:"stop_on_first_error"`

	// Maximum number of errors to collect
	MaxErrors int `json:"max_errors" hcl:"max_errors"`

	// Whether to include warnings
	IncludeWarnings bool `json:"include_warnings" hcl:"include_warnings"`

	// Whether to include context in errors
	IncludeContext bool `json:"include_context" hcl:"include_context"`

	// Whether to include suggestions in errors
	IncludeSuggestions bool `json:"include_suggestions" hcl:"include_suggestions"`
}

// EvolutionConfig provides schema evolution configuration
type EvolutionConfig struct {
	// Whether to enable evolution tracking
	EnableTracking bool `json:"enable_tracking" hcl:"enable_tracking"`

	// Whether to allow deprecated fields
	AllowDeprecated bool `json:"allow_deprecated" hcl:"allow_deprecated"`

	// Whether to warn about deprecated fields
	WarnDeprecated bool `json:"warn_deprecated" hcl:"warn_deprecated"`

	// Whether to allow breaking changes
	AllowBreaking bool `json:"allow_breaking" hcl:"allow_breaking"`
}

// CustomValidatorFunc represents a custom validation function
type CustomValidatorFunc func(ctx context.Context, data interface{}, config map[string]interface{}) (*spookytypesschemas.ValidationResult, error)

// NewEnhancedValidator creates a new enhanced validator instance
func NewEnhancedValidator(logger spookytypeslogging.Logger, config *ValidationConfig) *EnhancedValidator {
	if config == nil {
		config = &ValidationConfig{
			Mode: ValidationModeStrict,
			ErrorHandling: &ErrorHandlingConfig{
				StopOnFirstError:   false,
				MaxErrors:          100,
				IncludeWarnings:    true,
				IncludeContext:     true,
				IncludeSuggestions: true,
			},
			Evolution: &EvolutionConfig{
				EnableTracking:  true,
				AllowDeprecated: true,
				WarnDeprecated:  true,
				AllowBreaking:   false,
			},
		}
	}

	return &EnhancedValidator{
		logger:           logger,
		parser:           hclparse.NewParser(),
		config:           config,
		customValidators: make(map[string]CustomValidatorFunc),
	}
}

// SetRegistry sets the schema registry for the validator
func (v *EnhancedValidator) SetRegistry(registry spookytypesschemas.SchemaRegistry) {
	v.registry = registry
}

// ValidateWithEnhancedFeatures validates data with enhanced features
func (v *EnhancedValidator) ValidateWithEnhancedFeatures(ctx context.Context, schema *spookytypesschemas.Schema, data interface{}) (*spookytypesschemas.ValidationResult, error) {
	start := time.Now()

	// Create validation result
	result := &spookytypesschemas.ValidationResult{
		Valid:       true,
		ValidatedAt: time.Now(),
		Errors:      []spookytypesschemas.SchemaError{},
		Warnings:    []spookytypesschemas.SchemaError{},
		Info:        []spookytypesschemas.SchemaError{},
		Details:     make(map[string]interface{}),
		Statistics: &spookytypesschemas.ValidationStatistics{
			TotalFields:   0,
			ValidFields:   0,
			InvalidFields: 0,
			RulesExecuted: 0,
			RulesFailed:   0,
		},
	}

	// Validate schema configuration
	if err := v.validateSchemaConfiguration(schema, result); err != nil {
		return result, err
	}

	// Check for early termination
	if v.shouldStopValidation(result) {
		return result, nil
	}

	// Validate data structure
	if err := v.validateDataStructure(schema, data, result); err != nil {
		return result, err
	}

	// Check for early termination
	if v.shouldStopValidation(result) {
		return result, nil
	}

	// Parse HCL if data is HCL
	var parsedData interface{}
	if hclData, ok := data.([]byte); ok {
		parsed, err := v.parseHCLData(hclData, "data")
		if err != nil {
			v.addError(result, "hcl_parse_error", "Failed to parse HCL data", err.Error(), "error")
			return result, nil
		}
		parsedData = parsed
	} else if hclString, ok := data.(string); ok {
		parsed, err := v.parseHCLData([]byte(hclString), "data")
		if err != nil {
			v.addError(result, "hcl_parse_error", "Failed to parse HCL string", err.Error(), "error")
			return result, nil
		}
		parsedData = parsed
	} else {
		parsedData = data
	}

	// Validate field constraints
	if err := v.validateFieldConstraints(schema, parsedData, result); err != nil {
		return result, err
	}

	// Check for early termination
	if v.shouldStopValidation(result) {
		return result, nil
	}

	// Validate cross-field rules
	if err := v.validateCrossFieldRules(schema, parsedData, result); err != nil {
		return result, err
	}

	// Check for early termination
	if v.shouldStopValidation(result) {
		return result, nil
	}

	// Validate custom rules
	if err := v.validateCustomRules(ctx, schema, parsedData, result); err != nil {
		return result, err
	}

	// Check for early termination
	if v.shouldStopValidation(result) {
		return result, nil
	}

	// Validate schema evolution
	if err := v.validateSchemaEvolution(schema, parsedData, result); err != nil {
		return result, err
	}

	// Update statistics
	result.Statistics.Duration = time.Since(start)
	result.Statistics.ValidFields = result.Statistics.TotalFields - result.Statistics.InvalidFields

	// Update valid flag based on errors
	if len(result.Errors) > 0 {
		result.Valid = false
	}

	// Generate recommendations
	v.generateEnhancedRecommendations(schema, result)

	return result, nil
}

// ValidateFileWithEnhancedFeatures validates a file with enhanced features
func (v *EnhancedValidator) ValidateFileWithEnhancedFeatures(ctx context.Context, schema *spookytypesschemas.Schema, filePath string) (*spookytypesschemas.ValidationResult, error) {
	// Read file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	result, err := v.ValidateWithEnhancedFeatures(ctx, schema, data)
	if err != nil {
		return result, err
	}

	// Add file context
	result.Details["file_path"] = filePath
	result.Details["file_size"] = len(data)
	result.Details["file_modified"] = v.getFileModTime(filePath)

	return result, nil
}

// RegisterCustomValidator registers a custom validation function
func (v *EnhancedValidator) RegisterCustomValidator(name string, validator CustomValidatorFunc) error {
	if name == "" {
		return fmt.Errorf("validator name cannot be empty")
	}

	if validator == nil {
		return fmt.Errorf("validator function cannot be nil")
	}

	v.customValidators[name] = validator

	v.logger.Debug("Custom validator registered", map[string]interface{}{
		"validator_name": name,
	})

	return nil
}

// Helper methods

// validateSchemaConfiguration validates the schema configuration itself
func (v *EnhancedValidator) validateSchemaConfiguration(schema *spookytypesschemas.Schema, result *spookytypesschemas.ValidationResult) error {
	if schema == nil {
		v.addError(result, "schema_nil", "Schema cannot be nil", "", "error")
		return fmt.Errorf("schema cannot be nil")
	}

	// Validate required schema fields
	if schema.Name == "" {
		v.addError(result, "schema_name_missing", "Schema name is required", "", "error")
	}

	if schema.Type == "" {
		v.addError(result, "schema_type_missing", "Schema type is required", "", "error")
	}

	if schema.Version == "" {
		v.addWarning(result, "schema_version_missing", "Schema version is recommended", "", "warning")
	}

	return nil
}

// validateDataStructure validates the basic data structure
func (v *EnhancedValidator) validateDataStructure(schema *spookytypesschemas.Schema, data interface{}, result *spookytypesschemas.ValidationResult) error {
	if data == nil {
		v.addError(result, "data_nil", "Data cannot be nil", "", "error")
		return fmt.Errorf("data cannot be nil")
	}

	return nil
}

// validateFieldConstraints validates field constraints for a schema
func (v *EnhancedValidator) validateFieldConstraints(schema *spookytypesschemas.Schema, data interface{}, result *spookytypesschemas.ValidationResult) error {
	if schema.Validation == nil || schema.Validation.Fields == nil {
		return nil
	}

	for fieldPath, fieldValidation := range schema.Validation.Fields {
		result.Statistics.TotalFields++

		// Extract field value from data
		fieldValue := v.extractFieldValue(data, fieldPath)

		if err := v.validateFieldValue(fieldValidation, fieldValue, fieldPath, result); err != nil {
			result.Statistics.InvalidFields++
			result.Statistics.RulesFailed++
		} else {
			result.Statistics.ValidFields++
		}
		result.Statistics.RulesExecuted++

		// Check for early termination
		if v.shouldStopValidation(result) {
			return nil
		}
	}

	return nil
}

// validateFieldValue validates a single field value
func (v *EnhancedValidator) validateFieldValue(fieldValidation *spookytypesschemas.FieldValidation, value interface{}, fieldPath string, result *spookytypesschemas.ValidationResult) error {
	// Check required field
	if fieldValidation.Required && value == nil {
		v.addError(result, "required_field_missing", fmt.Sprintf("Required field '%s' is missing", fieldPath), fmt.Sprintf("Add the required field '%s'", fieldPath), "error")
		return fmt.Errorf("required field missing")
	}

	// Skip validation if value is nil and not required
	if value == nil {
		return nil
	}

	// Validate constraints
	if fieldValidation.Constraints != nil {
		if err := v.validateFieldConstraintsValue(fieldValidation.Constraints, value, fieldPath, result); err != nil {
			return err
		}
	}

	return nil
}

// validateFieldConstraintsValue validates field constraints
func (v *EnhancedValidator) validateFieldConstraintsValue(constraints *spookytypesschemas.FieldConstraints, value interface{}, fieldPath string, result *spookytypesschemas.ValidationResult) error {
	// String constraints
	if strValue, ok := value.(string); ok {
		if constraints.MinLength != nil && len(strValue) < *constraints.MinLength {
			v.addError(result, "string_too_short", fmt.Sprintf("Field '%s' length %d is less than minimum %d", fieldPath, len(strValue), *constraints.MinLength), fmt.Sprintf("Increase the length of field '%s' to at least %d characters", fieldPath, *constraints.MinLength), "error")
		}

		if constraints.MaxLength != nil && len(strValue) > *constraints.MaxLength {
			v.addError(result, "string_too_long", fmt.Sprintf("Field '%s' length %d exceeds maximum %d", fieldPath, len(strValue), *constraints.MaxLength), fmt.Sprintf("Reduce the length of field '%s' to at most %d characters", fieldPath, *constraints.MaxLength), "error")
		}

		if constraints.Pattern != nil {
			matched, err := regexp.MatchString(*constraints.Pattern, strValue)
			if err != nil {
				v.addError(result, "invalid_regex", fmt.Sprintf("Invalid regex pattern for field '%s': %v", fieldPath, err), "", "error")
			} else if !matched {
				v.addError(result, "pattern_mismatch", fmt.Sprintf("Field '%s' value '%s' does not match pattern '%s'", fieldPath, strValue, *constraints.Pattern), fmt.Sprintf("Ensure field '%s' matches the required pattern", fieldPath), "error")
			}
		}

		if constraints.Format != nil {
			if err := v.validateFormat(*constraints.Format, strValue, fieldPath, result); err != nil {
				return err
			}
		}
	}

	// Numeric constraints
	if numValue, ok := v.toFloat64(value); ok {
		if constraints.Min != nil && numValue < *constraints.Min {
			v.addError(result, "number_too_small", fmt.Sprintf("Field '%s' value %f is less than minimum %f", fieldPath, numValue, *constraints.Min), fmt.Sprintf("Increase the value of field '%s' to at least %f", fieldPath, *constraints.Min), "error")
		}

		if constraints.Max != nil && numValue > *constraints.Max {
			v.addError(result, "number_too_large", fmt.Sprintf("Field '%s' value %f exceeds maximum %f", fieldPath, numValue, *constraints.Max), fmt.Sprintf("Reduce the value of field '%s' to at most %f", fieldPath, *constraints.Max), "error")
		}
	}

	// Enum constraints
	if constraints.Enum != nil && len(constraints.Enum) > 0 {
		found := false
		for _, enumValue := range constraints.Enum {
			if value == enumValue {
				found = true
				break
			}
		}
		if !found {
			v.addError(result, "invalid_enum_value", fmt.Sprintf("Field '%s' value '%v' is not in allowed enum values", fieldPath, value), fmt.Sprintf("Use one of the allowed values for field '%s': %v", fieldPath, constraints.Enum), "error")
		}
	}

	return nil
}

// validateFormat validates format constraints
func (v *EnhancedValidator) validateFormat(format, value, fieldPath string, result *spookytypesschemas.ValidationResult) error {
	switch format {
	case "email":
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(value) {
			v.addError(result, "invalid_email", fmt.Sprintf("Field '%s' value '%s' is not a valid email address", fieldPath, value), "Enter a valid email address (e.g., user@example.com)", "error")
		}
	case "uri":
		if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
			v.addError(result, "invalid_uri", fmt.Sprintf("Field '%s' value '%s' is not a valid URI", fieldPath, value), "Enter a valid URI starting with http:// or https://", "error")
		}
	case "ipv4":
		ipv4Regex := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
		if !ipv4Regex.MatchString(value) {
			v.addError(result, "invalid_ipv4", fmt.Sprintf("Field '%s' value '%s' is not a valid IPv4 address", fieldPath, value), "Enter a valid IPv4 address (e.g., 192.168.1.1)", "error")
		}
	}
	return nil
}

// validateCrossFieldRules validates cross-field validation rules
func (v *EnhancedValidator) validateCrossFieldRules(schema *spookytypesschemas.Schema, data interface{}, result *spookytypesschemas.ValidationResult) error {
	if schema.Validation == nil || len(schema.Validation.CrossFieldValidations) == 0 {
		return nil
	}

	for _, crossFieldValidation := range schema.Validation.CrossFieldValidations {
		result.Statistics.RulesExecuted++

		// Extract field values
		fieldValues := make(map[string]interface{})
		for _, field := range crossFieldValidation.Fields {
			fieldValues[field] = v.extractFieldValue(data, field)
		}

		// Evaluate condition
		if !v.evaluateCrossFieldCondition(crossFieldValidation.Condition, fieldValues) {
			error := spookytypesschemas.NewSchemaError("", "", crossFieldValidation.Message)
			error.Severity = crossFieldValidation.Severity
			error.AddContext("fields", crossFieldValidation.Fields)
			error.AddContext("condition", crossFieldValidation.Condition)

			switch crossFieldValidation.Severity {
			case "error":
				result.Errors = append(result.Errors, *error)
			case "warning":
				result.Warnings = append(result.Warnings, *error)
			case "info":
				result.Info = append(result.Info, *error)
			}

			result.Statistics.RulesFailed++
		}

		// Check for early termination
		if v.shouldStopValidation(result) {
			return nil
		}
	}

	return nil
}

// validateCustomRules validates custom validation rules
func (v *EnhancedValidator) validateCustomRules(ctx context.Context, schema *spookytypesschemas.Schema, data interface{}, result *spookytypesschemas.ValidationResult) error {
	if schema.Validation == nil || len(schema.Validation.CustomValidators) == 0 {
		return nil
	}

	for _, customValidator := range schema.Validation.CustomValidators {
		result.Statistics.RulesExecuted++

		// Execute custom validator
		if validator, exists := v.customValidators[customValidator]; exists {
			if customResult, err := validator(ctx, data, nil); err != nil {
				error := spookytypesschemas.NewSchemaError("", "", fmt.Sprintf("Custom validation failed: %v", err))
				error.Severity = "error"
				error.AddContext("validator", customValidator)
				result.Errors = append(result.Errors, *error)
				result.Statistics.RulesFailed++
			} else if customResult != nil {
				// Merge custom validation results
				result.Errors = append(result.Errors, customResult.Errors...)
				result.Warnings = append(result.Warnings, customResult.Warnings...)
				result.Info = append(result.Info, customResult.Info...)
			}
		} else {
			v.addWarning(result, "custom_validator_not_found", fmt.Sprintf("Custom validator '%s' not found", customValidator), "Register the custom validator or remove the reference", "warning")
		}

		// Check for early termination
		if v.shouldStopValidation(result) {
			return nil
		}
	}

	return nil
}

// validateSchemaEvolution validates schema evolution
func (v *EnhancedValidator) validateSchemaEvolution(schema *spookytypesschemas.Schema, data interface{}, result *spookytypesschemas.ValidationResult) error {
	if !v.config.Evolution.EnableTracking || schema.Evolution == nil {
		return nil
	}

	// Check for deprecated fields
	if v.config.Evolution.WarnDeprecated {
		for _, deprecation := range schema.Evolution.Deprecations {
			if v.extractFieldValue(data, deprecation.Field) != nil {
				message := fmt.Sprintf("Field '%s' is deprecated", deprecation.Field)
				if deprecation.Reason != "" {
					message += fmt.Sprintf(": %s", deprecation.Reason)
				}
				if deprecation.Replacement != "" {
					message += fmt.Sprintf(" (use '%s' instead)", deprecation.Replacement)
				}
				v.addWarning(result, "deprecated_field", message, deprecation.MigrationGuidance, "warning")
			}
		}
	}

	// Check for breaking changes
	if !v.config.Evolution.AllowBreaking {
		for _, breakingChange := range schema.Evolution.BreakingChanges {
			if v.extractFieldValue(data, breakingChange.Field) != nil {
				message := fmt.Sprintf("Breaking change detected for field '%s': %s", breakingChange.Field, breakingChange.Description)
				v.addError(result, "breaking_change", message, breakingChange.Mitigation, "error")
			}
		}
	}

	return nil
}

// generateEnhancedRecommendations generates enhanced validation recommendations
func (v *EnhancedValidator) generateEnhancedRecommendations(schema *spookytypesschemas.Schema, result *spookytypesschemas.ValidationResult) {
	if len(result.Errors) == 0 && len(result.Warnings) == 0 {
		result.Recommendations = append(result.Recommendations, "Schema validation passed successfully")
		return
	}

	// Generate recommendations based on error types
	errorTypes := make(map[string]int)
	for _, error := range result.Errors {
		errorTypes[error.Code]++
	}

	for errorCode, count := range errorTypes {
		switch errorCode {
		case "hcl_parse_error":
			result.Recommendations = append(result.Recommendations,
				fmt.Sprintf("Fix %d HCL syntax errors by checking brackets, quotes, and commas", count))
		case "required_field_missing":
			result.Recommendations = append(result.Recommendations,
				fmt.Sprintf("Add %d missing required fields", count))
		case "field_validation_failed":
			result.Recommendations = append(result.Recommendations,
				fmt.Sprintf("Fix %d field validation errors by checking constraints", count))
		case "deprecated_field":
			result.Recommendations = append(result.Recommendations,
				fmt.Sprintf("Update %d deprecated fields to their replacements", count))
		case "breaking_change":
			result.Recommendations = append(result.Recommendations,
				fmt.Sprintf("Address %d breaking changes before proceeding", count))
		}
	}

	// Add general recommendations
	if len(result.Errors) > 0 {
		result.Recommendations = append(result.Recommendations, "Review and fix all validation errors before proceeding")
	}
	if len(result.Warnings) > 0 {
		result.Recommendations = append(result.Recommendations, "Consider addressing validation warnings for better schema compliance")
	}
}

// Helper methods

// addError adds an error to the validation result
func (v *EnhancedValidator) addError(result *spookytypesschemas.ValidationResult, code, message, suggestion, severity string) {
	error := spookytypesschemas.NewSchemaError("", "", message)
	error.Code = code
	error.Severity = severity
	if suggestion != "" {
		error.AddSuggestion(suggestion)
	}
	result.Errors = append(result.Errors, *error)
}

// addWarning adds a warning to the validation result
func (v *EnhancedValidator) addWarning(result *spookytypesschemas.ValidationResult, code, message, suggestion, severity string) {
	warning := spookytypesschemas.NewSchemaError("", "", message)
	warning.Code = code
	warning.Severity = severity
	if suggestion != "" {
		warning.AddSuggestion(suggestion)
	}
	result.Warnings = append(result.Warnings, *warning)
}

// shouldStopValidation checks if validation should stop
func (v *EnhancedValidator) shouldStopValidation(result *spookytypesschemas.ValidationResult) bool {
	if v.config.ErrorHandling.StopOnFirstError && len(result.Errors) > 0 {
		return true
	}

	if len(result.Errors) >= v.config.ErrorHandling.MaxErrors {
		return true
	}

	return false
}

// parseHCLData parses HCL data
func (v *EnhancedValidator) parseHCLData(data []byte, source string) (interface{}, error) {
	// Use hclsimple to parse the HCL data
	var result map[string]interface{}
	err := hclsimple.Decode(source, data, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("HCL parsing failed: %w", err)
	}
	return result, nil
}

// extractFieldValue extracts a field value from data
func (v *EnhancedValidator) extractFieldValue(data interface{}, fieldPath string) interface{} {
	// Handle map[string]interface{} data structure
	if dataMap, ok := data.(map[string]interface{}); ok {
		// Simple field path extraction (supports dot notation)
		parts := strings.Split(fieldPath, ".")
		current := dataMap

		for _, part := range parts {
			if val, exists := current[part]; exists {
				if nestedMap, isMap := val.(map[string]interface{}); isMap {
					current = nestedMap
				} else {
					return val
				}
			} else {
				return nil
			}
		}
		return current
	}

	return nil
}

// toFloat64 converts a value to float64
func (v *EnhancedValidator) toFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case int32:
		return float64(v), true
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}

// evaluateCrossFieldCondition evaluates a cross-field condition
func (v *EnhancedValidator) evaluateCrossFieldCondition(condition string, fieldValues map[string]interface{}) bool {
	// Simple condition evaluation for common patterns
	condition = strings.ToLower(strings.TrimSpace(condition))

	// Handle common condition patterns
	switch {
	case strings.Contains(condition, "required"):
		// Check if required fields are present
		requiredFields := extractRequiredFields(condition)
		for _, field := range requiredFields {
			if value, exists := fieldValues[field]; !exists || value == nil || value == "" {
				return false
			}
		}
		return true

	case strings.Contains(condition, "equals"):
		// Check if two fields are equal
		fields := extractFieldNames(condition)
		if len(fields) >= 2 {
			val1 := fieldValues[fields[0]]
			val2 := fieldValues[fields[1]]
			return fmt.Sprintf("%v", val1) == fmt.Sprintf("%v", val2)
		}
		return false

	case strings.Contains(condition, "not_equals"):
		// Check if two fields are not equal
		fields := extractFieldNames(condition)
		if len(fields) >= 2 {
			val1 := fieldValues[fields[0]]
			val2 := fieldValues[fields[1]]
			return fmt.Sprintf("%v", val1) != fmt.Sprintf("%v", val2)
		}
		return false

	default:
		// Default to true for unknown conditions
		return true
	}
}

// Helper functions for condition evaluation
func extractRequiredFields(condition string) []string {
	// Simple extraction of field names from "required" conditions
	// This is a basic implementation - could be enhanced with proper parsing
	var fields []string
	if strings.Contains(condition, "(") && strings.Contains(condition, ")") {
		start := strings.Index(condition, "(")
		end := strings.LastIndex(condition, ")")
		if start < end {
			fieldList := condition[start+1 : end]
			fields = strings.Split(fieldList, ",")
			for i, field := range fields {
				fields[i] = strings.TrimSpace(field)
			}
		}
	}
	return fields
}

func extractFieldNames(condition string) []string {
	// Simple extraction of field names from conditions
	// This is a basic implementation - could be enhanced with proper parsing
	var fields []string
	words := strings.Fields(condition)
	for _, word := range words {
		if !strings.Contains(word, "(") && !strings.Contains(word, ")") &&
			word != "equals" && word != "not_equals" && word != "required" {
			fields = append(fields, word)
		}
	}
	return fields
}

// getFileModTime gets file modification time
func (v *EnhancedValidator) getFileModTime(filePath string) time.Time {
	if info, err := os.Stat(filePath); err == nil {
		return info.ModTime()
	}
	return time.Time{}
}
