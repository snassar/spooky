// Package schemas provides schema validation and management functionality for the spooky codebase.
package schemas

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"

	spookytypeslogging "spooky/internal/types/logging"
	spookytypesschemas "spooky/internal/types/schemas"
)

// Validator provides comprehensive schema validation functionality
type Validator struct {
	logger   spookytypeslogging.Logger
	schemas  map[string]*spookytypesschemas.Schema
	registry spookytypesschemas.SchemaRegistry
}

// NewValidator creates a new schema validator instance
func NewValidator(logger spookytypeslogging.Logger) *Validator {
	return &Validator{
		logger:  logger,
		schemas: make(map[string]*spookytypesschemas.Schema),
	}
}

// SetRegistry sets the schema registry for the validator
func (v *Validator) SetRegistry(registry spookytypesschemas.SchemaRegistry) {
	v.registry = registry
}

// LoadSchemas loads all schema files from the schemas directory
func (v *Validator) LoadSchemas(schemasDir string) error {
	v.logger.Info("Loading schemas from directory", map[string]interface{}{
		"schemas_dir": schemasDir,
	})

	// Walk through schemas directory
	err := filepath.Walk(schemasDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process schema files
		if !strings.HasSuffix(info.Name(), ".schema.hcl") {
			return nil
		}

		// Load schema from file
		schema, err := v.loadSchemaFromFile(path)
		if err != nil {
			v.logger.Warn("Failed to load schema", map[string]interface{}{
				"file":  path,
				"error": err.Error(),
			})
			return nil // Continue with other schemas
		}

		// Store schema by name
		schemaName := strings.TrimSuffix(info.Name(), ".schema.hcl")
		v.schemas[schemaName] = schema

		// Register with registry if available
		if v.registry != nil {
			if err := v.registry.Register(schema); err != nil {
				v.logger.Warn("Failed to register schema", map[string]interface{}{
					"schema_name": schemaName,
					"error":       err.Error(),
				})
			}
		}

		v.logger.Debug("Loaded schema", map[string]interface{}{
			"schema_name": schemaName,
			"file":        path,
		})

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk schemas directory: %w", err)
	}

	v.logger.Info("Schemas loaded successfully", map[string]interface{}{
		"schemas_dir": schemasDir,
		"count":       len(v.schemas),
		"schemas":     v.getSchemaNames(),
	})

	return nil
}

// ValidateFile validates an HCL file against its corresponding schema
func (v *Validator) ValidateFile(filePath string, schemaName string) (*spookytypesschemas.ValidationResult, error) {
	v.logger.Debug("Validating file against schema", map[string]interface{}{
		"file":        filePath,
		"schema_name": schemaName,
	})

	// Get schema
	schema, exists := v.schemas[schemaName]
	if !exists {
		return nil, fmt.Errorf("schema '%s' not found", schemaName)
	}

	// Read file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Parse HCL content
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(data, filePath)
	if diags.HasErrors() {
		// Create validation result with HCL parsing errors
		result := &spookytypesschemas.ValidationResult{
			Valid:       false,
			ValidatedAt: time.Now(),
			Errors:      []spookytypesschemas.SchemaError{},
			Warnings:    []spookytypesschemas.SchemaError{},
			Info:        []spookytypesschemas.SchemaError{},
		}

		// Convert HCL diagnostics to schema errors
		for _, diag := range diags {
			error := spookytypesschemas.SchemaError{
				Code:        "hcl_parse_error",
				Message:     diag.Summary,
				SchemaName:  schemaName,
				SchemaType:  schema.Type,
				Severity:    "error",
				Recoverable: false,
				Timestamp:   time.Now(),
			}
			if diag.Detail != "" {
				error.Message += ": " + diag.Detail
			}

			// Add location information
			if diag.Subject != nil {
				error.SetLocation(filePath, diag.Subject.Start.Line, diag.Subject.Start.Column)
			}

			// Add suggestions
			error.AddSuggestion("Check HCL syntax (brackets, quotes, commas)")
			error.AddSuggestion("Validate required fields are present")
			error.AddSuggestion("Check for typos in field names")

			result.Errors = append(result.Errors, error)
		}

		return result, nil
	}

	// Validate against schema
	result := v.validateAgainstSchema(file, schema, filePath)

	v.logger.Debug("File validation completed", map[string]interface{}{
		"file":          filePath,
		"schema_name":   schemaName,
		"valid":         result.Valid,
		"error_count":   len(result.Errors),
		"warning_count": len(result.Warnings),
	})

	return result, nil
}

// ValidateDirectory validates all HCL files in a directory against their schemas
func (v *Validator) ValidateDirectory(dirPath string, schemaMapping map[string]string) (map[string]*spookytypesschemas.ValidationResult, error) {
	v.logger.Info("Validating directory", map[string]interface{}{
		"dir_path": dirPath,
	})

	results := make(map[string]*spookytypesschemas.ValidationResult)

	// Walk through directory
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process HCL files
		if !strings.HasSuffix(info.Name(), ".hcl") {
			return nil
		}

		// Determine schema name for this file
		schemaName := v.determineSchemaName(path, schemaMapping)
		if schemaName == "" {
			v.logger.Debug("No schema mapping found for file", map[string]interface{}{
				"file": path,
			})
			return nil
		}

		// Validate file
		result, err := v.ValidateFile(path, schemaName)
		if err != nil {
			v.logger.Warn("Failed to validate file", map[string]interface{}{
				"file":        path,
				"schema_name": schemaName,
				"error":       err.Error(),
			})
			return nil // Continue with other files
		}

		results[path] = result
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	v.logger.Info("Directory validation completed", map[string]interface{}{
		"dir_path":        dirPath,
		"files_validated": len(results),
	})

	return results, nil
}

// Validate validates data against a schema
func (v *Validator) Validate(schema *spookytypesschemas.Schema, data interface{}) (*spookytypesschemas.ValidationResult, error) {
	start := time.Now()

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

	// Validate data structure
	if err := v.validateDataStructure(schema, data, result); err != nil {
		return result, err
	}

	// Validate field constraints
	if err := v.validateSchemaFieldConstraints(schema, data, result); err != nil {
		return result, err
	}

	// Validate cross-field rules
	if err := v.validateCrossFieldRules(schema, data, result); err != nil {
		return result, err
	}

	// Validate custom rules
	if err := v.validateCustomRules(schema, data, result); err != nil {
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
	v.generateRecommendations(schema, result)

	return result, nil
}

// ValidateFileAgainstSchema validates a file against a schema
func (v *Validator) ValidateFileAgainstSchema(schema *spookytypesschemas.Schema, filePath string) (*spookytypesschemas.ValidationResult, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	result, err := v.Validate(schema, data)
	if err != nil {
		return result, err
	}

	// Add file context
	result.Details["file_path"] = filePath
	result.Details["file_size"] = len(data)

	return result, nil
}

// ValidateString validates a string against a schema
func (v *Validator) ValidateString(schema *spookytypesschemas.Schema, content string) (*spookytypesschemas.ValidationResult, error) {
	return v.Validate(schema, content)
}

// ValidateBytes validates bytes against a schema
func (v *Validator) ValidateBytes(schema *spookytypesschemas.Schema, data []byte) (*spookytypesschemas.ValidationResult, error) {
	return v.Validate(schema, data)
}

// ValidateWithContext validates data with additional context
func (v *Validator) ValidateWithContext(schema *spookytypesschemas.Schema, data interface{}, context map[string]interface{}) (*spookytypesschemas.ValidationResult, error) {
	result, err := v.Validate(schema, data)
	if err != nil {
		return result, err
	}

	// Add context to result
	for key, value := range context {
		result.Details[key] = value
	}

	return result, nil
}

// ValidateField validates a specific field
func (v *Validator) ValidateField(schema *spookytypesschemas.Schema, fieldPath string, value interface{}) (*spookytypesschemas.ValidationResult, error) {
	result := &spookytypesschemas.ValidationResult{
		Valid:       true,
		ValidatedAt: time.Now(),
		Errors:      []spookytypesschemas.SchemaError{},
		Warnings:    []spookytypesschemas.SchemaError{},
		Info:        []spookytypesschemas.SchemaError{},
		Details:     make(map[string]interface{}),
	}

	// Find field validation rules
	if schema.Validation != nil && schema.Validation.Fields != nil {
		if fieldValidation, exists := schema.Validation.Fields[fieldPath]; exists {
			if err := v.validateFieldValue(fieldValidation, value, fieldPath, result); err != nil {
				return result, err
			}
		}
	}

	// Update valid flag based on errors
	if len(result.Errors) > 0 {
		result.Valid = false
	}

	return result, nil
}

// validateSchemaConfiguration validates the schema configuration itself
func (v *Validator) validateSchemaConfiguration(schema *spookytypesschemas.Schema, result *spookytypesschemas.ValidationResult) error {
	if schema == nil {
		error := spookytypesschemas.NewSchemaError("", "", "Schema cannot be nil")
		error.Severity = "error"
		result.Errors = append(result.Errors, *error)
		return fmt.Errorf("schema cannot be nil")
	}

	// Validate required schema fields
	if schema.Name == "" {
		error := spookytypesschemas.NewSchemaError(schema.Name, schema.Type, "Schema name is required")
		error.Severity = "error"
		result.Errors = append(result.Errors, *error)
	}

	if schema.Type == "" {
		error := spookytypesschemas.NewSchemaError(schema.Name, schema.Type, "Schema type is required")
		error.Severity = "error"
		result.Errors = append(result.Errors, *error)
	}

	if schema.Version == "" {
		error := spookytypesschemas.NewSchemaError(schema.Name, schema.Type, "Schema version is required")
		error.Severity = "warning"
		result.Warnings = append(result.Warnings, *error)
	}

	return nil
}

// validateDataStructure validates the basic data structure
func (v *Validator) validateDataStructure(schema *spookytypesschemas.Schema, data interface{}, result *spookytypesschemas.ValidationResult) error {
	if data == nil {
		error := spookytypesschemas.NewSchemaError(schema.Name, schema.Type, "Data cannot be nil")
		error.Severity = "error"
		result.Errors = append(result.Errors, *error)
		return fmt.Errorf("data cannot be nil")
	}

	// Basic structure validation based on schema type
	switch schema.Type {
	case "hcl":
		if _, ok := data.([]byte); !ok {
			if _, ok := data.(string); !ok {
				error := spookytypesschemas.NewSchemaError(schema.Name, schema.Type, "HCL schema expects byte array or string data")
				error.Severity = "error"
				result.Errors = append(result.Errors, *error)
			}
		}
	}

	return nil
}

// validateSchemaFieldConstraints validates field constraints for a schema
func (v *Validator) validateSchemaFieldConstraints(schema *spookytypesschemas.Schema, data interface{}, result *spookytypesschemas.ValidationResult) error {
	if schema.Validation == nil || schema.Validation.Fields == nil {
		return nil
	}

	for fieldPath, fieldValidation := range schema.Validation.Fields {
		result.Statistics.TotalFields++

		// Extract field value from data (simplified - would need proper HCL parsing)
		fieldValue := v.extractFieldValue(data, fieldPath)

		if err := v.validateFieldValue(fieldValidation, fieldValue, fieldPath, result); err != nil {
			result.Statistics.InvalidFields++
			result.Statistics.RulesFailed++
		} else {
			result.Statistics.ValidFields++
		}
		result.Statistics.RulesExecuted++
	}

	return nil
}

// validateFieldValue validates a single field value
func (v *Validator) validateFieldValue(fieldValidation *spookytypesschemas.FieldValidation, value interface{}, fieldPath string, result *spookytypesschemas.ValidationResult) error {
	// Check required field
	if fieldValidation.Required && value == nil {
		error := spookytypesschemas.NewSchemaError("", "", fmt.Sprintf("Required field '%s' is missing", fieldPath))
		error.FieldPath = fieldPath
		error.Severity = "error"
		error.AddSuggestion(fmt.Sprintf("Add the required field '%s'", fieldPath))
		result.Errors = append(result.Errors, *error)
		return fmt.Errorf("required field missing")
	}

	// Skip validation if value is nil and not required
	if value == nil {
		return nil
	}

	// Validate constraints
	if fieldValidation.Constraints != nil {
		if err := v.validateFieldConstraints(fieldValidation.Constraints, value, fieldPath, result); err != nil {
			return err
		}
	}

	return nil
}

// validateFieldConstraints validates field constraints
func (v *Validator) validateFieldConstraints(constraints *spookytypesschemas.FieldConstraints, value interface{}, fieldPath string, result *spookytypesschemas.ValidationResult) error {
	// String constraints
	if strValue, ok := value.(string); ok {
		if constraints.MinLength != nil && len(strValue) < *constraints.MinLength {
			error := spookytypesschemas.NewSchemaError("", "", fmt.Sprintf("Field '%s' length %d is less than minimum %d", fieldPath, len(strValue), *constraints.MinLength))
			error.FieldPath = fieldPath
			error.Value = strValue
			error.Severity = "error"
			error.AddSuggestion(fmt.Sprintf("Increase the length of field '%s' to at least %d characters", fieldPath, *constraints.MinLength))
			result.Errors = append(result.Errors, *error)
		}

		if constraints.MaxLength != nil && len(strValue) > *constraints.MaxLength {
			error := spookytypesschemas.NewSchemaError("", "", fmt.Sprintf("Field '%s' length %d exceeds maximum %d", fieldPath, len(strValue), *constraints.MaxLength))
			error.FieldPath = fieldPath
			error.Value = strValue
			error.Severity = "error"
			error.AddSuggestion(fmt.Sprintf("Reduce the length of field '%s' to at most %d characters", fieldPath, *constraints.MaxLength))
			result.Errors = append(result.Errors, *error)
		}

		if constraints.Pattern != nil {
			matched, err := regexp.MatchString(*constraints.Pattern, strValue)
			if err != nil {
				error := spookytypesschemas.NewSchemaError("", "", fmt.Sprintf("Invalid regex pattern for field '%s': %v", fieldPath, err))
				error.FieldPath = fieldPath
				error.Severity = "error"
				result.Errors = append(result.Errors, *error)
			} else if !matched {
				error := spookytypesschemas.NewSchemaError("", "", fmt.Sprintf("Field '%s' value '%s' does not match pattern '%s'", fieldPath, strValue, *constraints.Pattern))
				error.FieldPath = fieldPath
				error.Value = strValue
				error.Severity = "error"
				error.AddSuggestion(fmt.Sprintf("Ensure field '%s' matches the required pattern", fieldPath))
				result.Errors = append(result.Errors, *error)
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
			error := spookytypesschemas.NewSchemaError("", "", fmt.Sprintf("Field '%s' value %f is less than minimum %f", fieldPath, numValue, *constraints.Min))
			error.FieldPath = fieldPath
			error.Value = numValue
			error.Severity = "error"
			error.AddSuggestion(fmt.Sprintf("Increase the value of field '%s' to at least %f", fieldPath, *constraints.Min))
			result.Errors = append(result.Errors, *error)
		}

		if constraints.Max != nil && numValue > *constraints.Max {
			error := spookytypesschemas.NewSchemaError("", "", fmt.Sprintf("Field '%s' value %f exceeds maximum %f", fieldPath, numValue, *constraints.Max))
			error.FieldPath = fieldPath
			error.Value = numValue
			error.Severity = "error"
			error.AddSuggestion(fmt.Sprintf("Reduce the value of field '%s' to at most %f", fieldPath, *constraints.Max))
			result.Errors = append(result.Errors, *error)
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
			error := spookytypesschemas.NewSchemaError("", "", fmt.Sprintf("Field '%s' value '%v' is not in allowed enum values", fieldPath, value))
			error.FieldPath = fieldPath
			error.Value = value
			error.Severity = "error"
			error.AddSuggestion(fmt.Sprintf("Use one of the allowed values for field '%s': %v", fieldPath, constraints.Enum))
			result.Errors = append(result.Errors, *error)
		}
	}

	return nil
}

// validateFormat validates format constraints
func (v *Validator) validateFormat(format, value, fieldPath string, result *spookytypesschemas.ValidationResult) error {
	switch format {
	case "email":
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(value) {
			error := spookytypesschemas.NewSchemaError("", "", fmt.Sprintf("Field '%s' value '%s' is not a valid email address", fieldPath, value))
			error.FieldPath = fieldPath
			error.Value = value
			error.Severity = "error"
			error.AddSuggestion("Enter a valid email address (e.g., user@example.com)")
			result.Errors = append(result.Errors, *error)
		}
	case "uri":
		if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
			error := spookytypesschemas.NewSchemaError("", "", fmt.Sprintf("Field '%s' value '%s' is not a valid URI", fieldPath, value))
			error.FieldPath = fieldPath
			error.Value = value
			error.Severity = "error"
			error.AddSuggestion("Enter a valid URI starting with http:// or https://")
			result.Errors = append(result.Errors, *error)
		}
	}
	return nil
}

// validateCrossFieldRules validates cross-field validation rules
func (v *Validator) validateCrossFieldRules(schema *spookytypesschemas.Schema, data interface{}, result *spookytypesschemas.ValidationResult) error {
	if schema.Validation == nil || len(schema.Validation.CrossFieldValidations) == 0 {
		return nil
	}

	for _, crossFieldValidation := range schema.Validation.CrossFieldValidations {
		result.Statistics.RulesExecuted++

		// Extract field values (simplified - would need proper HCL parsing)
		fieldValues := make(map[string]interface{})
		for _, field := range crossFieldValidation.Fields {
			fieldValues[field] = v.extractFieldValue(data, field)
		}

		// Evaluate condition (simplified - would need expression evaluation)
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
	}

	return nil
}

// validateCustomRules validates custom validation rules
func (v *Validator) validateCustomRules(schema *spookytypesschemas.Schema, data interface{}, result *spookytypesschemas.ValidationResult) error {
	if schema.Validation == nil || len(schema.Validation.CustomValidators) == 0 {
		return nil
	}

	for _, customValidator := range schema.Validation.CustomValidators {
		result.Statistics.RulesExecuted++

		// Execute custom validator (simplified - would need plugin system)
		if err := v.executeCustomValidator(customValidator, data); err != nil {
			error := spookytypesschemas.NewSchemaError("", "", fmt.Sprintf("Custom validation failed: %v", err))
			error.Severity = "error"
			error.AddContext("validator", customValidator)
			result.Errors = append(result.Errors, *error)
			result.Statistics.RulesFailed++
		}
	}

	return nil
}

// generateRecommendations generates validation recommendations
func (v *Validator) generateRecommendations(schema *spookytypesschemas.Schema, result *spookytypesschemas.ValidationResult) {
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

// loadSchemaFromFile loads a schema from an HCL file
func (v *Validator) loadSchemaFromFile(filePath string) (*spookytypesschemas.Schema, error) {
	// Read file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	// Create a basic schema from the file content
	schema := &spookytypesschemas.Schema{
		Version:     "1.0",
		Type:        "hcl",
		Name:        filepath.Base(filePath),
		Description: fmt.Sprintf("Schema loaded from %s", filePath),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Content:     string(data),
		Metadata:    make(map[string]interface{}),
	}

	// Parse HCL content to validate it
	parser := hclparse.NewParser()
	_, diags := parser.ParseHCL(data, filePath)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse schema file: %s", diags.Error())
	}

	return schema, nil
}

// validateAgainstSchema validates an HCL file against a schema
func (v *Validator) validateAgainstSchema(file *hcl.File, schema *spookytypesschemas.Schema, filePath string) *spookytypesschemas.ValidationResult {
	result := &spookytypesschemas.ValidationResult{
		Valid:       true,
		ValidatedAt: time.Now(),
		Errors:      []spookytypesschemas.SchemaError{},
		Warnings:    []spookytypesschemas.SchemaError{},
		Info:        []spookytypesschemas.SchemaError{},
		Details:     make(map[string]interface{}),
	}

	// Basic validation: check if the file has content
	if len(file.Bytes) == 0 {
		error := spookytypesschemas.NewSchemaError(schema.Name, schema.Type, "File is empty")
		error.Severity = "error"
		error.SetLocation(filePath, 1, 1)
		result.Errors = append(result.Errors, *error)
		result.Valid = false
	}

	// Add validation details
	result.Details["file_path"] = filePath
	result.Details["schema_name"] = schema.Name
	result.Details["schema_type"] = schema.Type
	result.Details["file_size"] = len(file.Bytes)

	// Update valid flag based on errors
	if len(result.Errors) > 0 {
		result.Valid = false
	}

	return result
}

// determineSchemaName determines the schema name for a file based on mapping
func (v *Validator) determineSchemaName(filePath string, schemaMapping map[string]string) string {
	fileName := filepath.Base(filePath)

	// Check direct file name mapping
	if schemaName, exists := schemaMapping[fileName]; exists {
		return schemaName
	}

	// Check file extension mapping
	ext := filepath.Ext(fileName)
	if schemaName, exists := schemaMapping[ext]; exists {
		return schemaName
	}

	// Default mappings based on file name patterns
	if strings.Contains(fileName, "machines") {
		return "machines"
	}
	if strings.Contains(fileName, "variables") {
		return "variables"
	}
	if strings.Contains(fileName, "project") {
		return "project"
	}
	if strings.Contains(fileName, "actions") {
		return "actions"
	}
	if strings.Contains(fileName, "facts") {
		return "facts"
	}
	if strings.Contains(fileName, "templates") {
		return "templates"
	}

	return ""
}

// getSchemaNames returns a list of loaded schema names
func (v *Validator) getSchemaNames() []string {
	var names []string
	for name := range v.schemas {
		names = append(names, name)
	}
	return names
}

// GetSchema returns a schema by name
func (v *Validator) GetSchema(name string) (*spookytypesschemas.Schema, bool) {
	schema, exists := v.schemas[name]
	return schema, exists
}

// ListSchemas returns all loaded schema names
func (v *Validator) ListSchemas() []string {
	return v.getSchemaNames()
}

// Helper methods for validation

// extractFieldValue extracts a field value from data (simplified)
func (v *Validator) extractFieldValue(data interface{}, fieldPath string) interface{} {
	// This is a simplified implementation
	// In a real implementation, this would parse HCL and extract the field value
	return nil
}

// toFloat64 converts a value to float64
func (v *Validator) toFloat64(value interface{}) (float64, bool) {
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

// evaluateCrossFieldCondition evaluates a cross-field condition (simplified)
func (v *Validator) evaluateCrossFieldCondition(condition string, fieldValues map[string]interface{}) bool {
	// This is a simplified implementation
	// In a real implementation, this would evaluate expressions
	return true
}

// executeCustomValidator executes a custom validator (simplified)
func (v *Validator) executeCustomValidator(validator string, data interface{}) error {
	// This is a simplified implementation
	// In a real implementation, this would execute custom validation logic
	return nil
}
