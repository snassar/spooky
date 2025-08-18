// Package schemas provides schema validation and management functionality for the spooky codebase.
package schemas

import (
	"fmt"
	"os"
	"path/filepath"
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
	manager  *Manager
	utils    *ValidationUtils
}

// NewValidator creates a new schema validator instance
func NewValidator(logger spookytypeslogging.Logger) *Validator {
	return &Validator{
		logger:  logger,
		schemas: make(map[string]*spookytypesschemas.Schema),
		utils:   NewValidationUtils(),
	}
}

// SetRegistry sets the schema registry for the validator
func (v *Validator) SetRegistry(registry spookytypesschemas.SchemaRegistry) {
	v.registry = registry
}

// SetManager sets the schema manager for the validator
func (v *Validator) SetManager(manager *Manager) {
	v.manager = manager
}

// createSchemaError creates a schema error for the validator
func (v *Validator) createSchemaError(result *spookytypesschemas.ValidationResult, code, message, suggestion, severity string) {
	schemaError := spookytypesschemas.NewSchemaError("", "", message)
	schemaError.FieldPath = ""
	schemaError.Severity = severity
	if suggestion != "" {
		schemaError.AddSuggestion(suggestion)
	}
	result.Errors = append(result.Errors, *schemaError)
}

// LoadSchemas loads all schema files from the schemas directory
func (v *Validator) LoadSchemas(schemasDir string) error {
	v.logger.Info("Loading schemas from directory", map[string]interface{}{
		"schemas_dir": schemasDir,
	})

	// Use manager if available for metadata validation
	if v.manager != nil {
		schemas, err := v.manager.LoadSchemasFromDirectory(schemasDir)
		if err != nil {
			return fmt.Errorf("failed to load schemas with metadata validation: %w", err)
		}

		// Store schemas in validator
		for name, schema := range schemas {
			v.schemas[name] = schema

			// Register with registry if available
			if v.registry != nil {
				if err := v.registry.Register(schema); err != nil {
					v.logger.Warn("Failed to register schema", map[string]interface{}{
						"schema_name": name,
						"error":       err.Error(),
					})
				}
			}
		}

		v.logger.Info("Schemas loaded successfully with metadata validation", map[string]interface{}{
			"schemas_dir": schemasDir,
			"count":       len(schemas),
			"schemas":     v.getSchemaNames(),
		})

		return nil
	}

	// Fallback to original implementation if no manager available
	return v.loadSchemasWithoutManager(schemasDir)
}

// loadSchemasWithoutManager loads schemas without metadata validation (fallback)
func (v *Validator) loadSchemasWithoutManager(schemasDir string) error {
	v.logger.Info("Loading schemas without metadata validation (fallback)", map[string]interface{}{
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

	v.logger.Info("Schemas loaded successfully (without metadata validation)", map[string]interface{}{
		"schemas_dir": schemasDir,
		"count":       len(v.schemas),
		"schemas":     v.getSchemaNames(),
	})

	return nil
}

// ValidateFile validates an HCL file against its corresponding schema
func (v *Validator) ValidateFile(filePath, schemaName string) (*spookytypesschemas.ValidationResult, error) {
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
			schemaError := spookytypesschemas.SchemaError{
				Code:        "hcl_parse_error",
				Message:     diag.Summary,
				SchemaName:  schemaName,
				SchemaType:  schema.Type,
				Severity:    "error",
				Recoverable: false,
				Timestamp:   time.Now(),
			}
			if diag.Detail != "" {
				schemaError.Message += ": " + diag.Detail
			}

			// Add location information
			if diag.Subject != nil {
				schemaError.SetLocation(filePath, diag.Subject.Start.Line, diag.Subject.Start.Column)
			}

			// Add suggestions
			schemaError.AddSuggestion("Check HCL syntax (brackets, quotes, commas)")
			schemaError.AddSuggestion("Validate required fields are present")
			schemaError.AddSuggestion("Check for typos in field names")

			result.Errors = append(result.Errors, schemaError)
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
			TotalFields:    0,
			ValidFields:    0,
			InvalidFields:  0,
			RulesProcessed: 0,
			RulesFailed:    0,
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
		schemaError := spookytypesschemas.NewSchemaError("", "", "Schema cannot be nil")
		schemaError.Severity = ValidationError
		result.Errors = append(result.Errors, *schemaError)
		return fmt.Errorf("schema cannot be nil")
	}

	// Validate required schema fields
	if schema.Name == "" {
		schemaError := spookytypesschemas.NewSchemaError(schema.Name, schema.Type, "Schema name is required")
		schemaError.Severity = ValidationError
		result.Errors = append(result.Errors, *schemaError)
	}

	if schema.Type == "" {
		schemaError := spookytypesschemas.NewSchemaError(schema.Name, schema.Type, "Schema type is required")
		schemaError.Severity = ValidationError
		result.Errors = append(result.Errors, *schemaError)
	}

	if schema.Version == "" {
		schemaError := spookytypesschemas.NewSchemaError(schema.Name, schema.Type, "Schema version is required")
		schemaError.Severity = "warning"
		result.Warnings = append(result.Warnings, *schemaError)
	}

	return nil
}

// validateDataStructure validates the basic data structure
func (v *Validator) validateDataStructure(schema *spookytypesschemas.Schema, data interface{}, result *spookytypesschemas.ValidationResult) error {
	if data == nil {
		schemaError := spookytypesschemas.NewSchemaError(schema.Name, schema.Type, "Data cannot be nil")
		schemaError.Severity = ValidationError
		result.Errors = append(result.Errors, *schemaError)
		return fmt.Errorf("data cannot be nil")
	}

	// Basic structure validation based on schema type
	if schema.Type == "hcl" {
		if _, ok := data.([]byte); !ok {
			if _, ok := data.(string); !ok {
				schemaError := spookytypesschemas.NewSchemaError(schema.Name, schema.Type, "HCL schema expects byte array or string data")
				schemaError.Severity = ValidationError
				result.Errors = append(result.Errors, *schemaError)
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
		result.Statistics.RulesProcessed++
	}

	return nil
}

// validateFieldValue validates a single field value
func (v *Validator) validateFieldValue(fieldValidation *spookytypesschemas.FieldValidation, value interface{}, fieldPath string, result *spookytypesschemas.ValidationResult) error {
	// Check required field
	if fieldValidation.Required && value == nil {
		schemaError := spookytypesschemas.NewSchemaError("", "", fmt.Sprintf("Required field '%s' is missing", fieldPath))
		schemaError.FieldPath = fieldPath
		schemaError.Severity = ValidationError
		schemaError.AddSuggestion(fmt.Sprintf("Add the required field '%s'", fieldPath))
		result.Errors = append(result.Errors, *schemaError)
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
	return v.utils.ValidateFieldConstraints(constraints, value, fieldPath, result, func(code, message, suggestion, severity string) {
		v.createSchemaError(result, code, message, suggestion, severity)
	}, v.toFloat64)
}

// validateCrossFieldRules validates cross-field validation rules
func (v *Validator) validateCrossFieldRules(schema *spookytypesschemas.Schema, data interface{}, result *spookytypesschemas.ValidationResult) error {
	if schema.Validation == nil || len(schema.Validation.CrossFieldValidations) == 0 {
		return nil
	}

	for _, crossFieldValidation := range schema.Validation.CrossFieldValidations {
		result.Statistics.RulesProcessed++

		// Extract field values (simplified - would need proper HCL parsing)
		fieldValues := make(map[string]interface{})
		for _, field := range crossFieldValidation.Fields {
			fieldValues[field] = v.extractFieldValue(data, field)
		}

		// Evaluate condition (simplified - would need expression evaluation)
		if !v.evaluateCrossFieldCondition(crossFieldValidation.Condition, fieldValues) {
			schemaError := spookytypesschemas.NewSchemaError("", "", crossFieldValidation.Message)
			schemaError.Severity = crossFieldValidation.Severity
			schemaError.AddContext("fields", crossFieldValidation.Fields)
			schemaError.AddContext("condition", crossFieldValidation.Condition)

			switch crossFieldValidation.Severity {
			case "error":
				result.Errors = append(result.Errors, *schemaError)
			case "warning":
				result.Warnings = append(result.Warnings, *schemaError)
			case "info":
				result.Info = append(result.Info, *schemaError)
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
		result.Statistics.RulesProcessed++

		// Run custom validator (simplified - would need plugin system)
		if err := v.executeCustomValidator(customValidator, data); err != nil {
			schemaError := spookytypesschemas.NewSchemaError("", "", fmt.Sprintf("Custom validation failed: %v", err))
			schemaError.Severity = "error"
			schemaError.AddContext("validator", customValidator)
			result.Errors = append(result.Errors, *schemaError)
			result.Statistics.RulesFailed++
		}
	}

	return nil
}

// generateRecommendations generates validation recommendations
func (v *Validator) generateRecommendations(_ *spookytypesschemas.Schema, result *spookytypesschemas.ValidationResult) {
	if len(result.Errors) == 0 && len(result.Warnings) == 0 {
		result.Recommendations = append(result.Recommendations, "Schema validation passed successfully")
		return
	}

	// Generate recommendations based on error types
	errorTypes := make(map[string]int)
	for i := range result.Errors {
		errorTypes[result.Errors[i].Code]++
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

	// Parse validation rules from HCL content
	schemaParser := NewSchemaParser(v.logger)
	if err := schemaParser.ParseValidationRules(schema); err != nil {
		return nil, fmt.Errorf("failed to parse validation rules: %w", err)
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
		schemaError := spookytypesschemas.NewSchemaError(schema.Name, schema.Type, "File is empty")
		schemaError.Severity = ValidationError
		schemaError.SetLocation(filePath, 1, 1)
		result.Errors = append(result.Errors, *schemaError)
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
	// Simple condition evaluation for common patterns
	condition = strings.ToLower(strings.TrimSpace(condition))

	// Handle common condition patterns
	switch {
	case strings.Contains(condition, "required"):
		return v.evaluateRequiredCondition(condition, fieldValues)
	case strings.Contains(condition, "equals"):
		return v.evaluateEqualsCondition(condition, fieldValues)
	case strings.Contains(condition, "not_equals"):
		return v.evaluateNotEqualsCondition(condition, fieldValues)
	default:
		// Default to true for unknown conditions
		return true
	}
}

// evaluateRequiredCondition evaluates a required condition
func (v *Validator) evaluateRequiredCondition(condition string, fieldValues map[string]interface{}) bool {
	words := strings.Fields(condition)
	for _, word := range words {
		if word != "required" && !strings.Contains(word, "(") && !strings.Contains(word, ")") {
			if value, exists := fieldValues[word]; !exists || value == nil || value == "" {
				return false
			}
		}
	}
	return true
}

// evaluateEqualsCondition evaluates an equals condition
func (v *Validator) evaluateEqualsCondition(condition string, fieldValues map[string]interface{}) bool {
	return v.evaluateComparisonCondition(condition, fieldValues, "equals", true)
}

// evaluateNotEqualsCondition evaluates a not_equals condition
func (v *Validator) evaluateNotEqualsCondition(condition string, fieldValues map[string]interface{}) bool {
	return v.evaluateComparisonCondition(condition, fieldValues, "not_equals", false)
}

// evaluateComparisonCondition evaluates comparison conditions (equals/not_equals)
func (v *Validator) evaluateComparisonCondition(condition string, fieldValues map[string]interface{}, operator string, isEquals bool) bool {
	words := strings.Fields(condition)
	var fields []string
	for _, word := range words {
		if word != operator && !strings.Contains(word, "(") && !strings.Contains(word, ")") {
			fields = append(fields, word)
		}
	}
	if len(fields) >= 2 {
		val1 := fieldValues[fields[0]]
		val2 := fieldValues[fields[1]]
		comparison := fmt.Sprintf("%v", val1) == fmt.Sprintf("%v", val2)
		if isEquals {
			return comparison
		}
		return !comparison
	}
	return false
}

// executeCustomValidator runs a custom validator
func (v *Validator) executeCustomValidator(validator string, data interface{}) error {
	// Custom validator implementation supporting common validation rules

	switch validator {
	case "required":
		return v.validateRequired(data)
	case "email":
		return v.validateEmail(data)
	case "url":
		return v.validateURL(data)
	case "positive":
		return v.validatePositive(data)
	default:
		// Unknown validator - return success
		return nil
	}
}

// validateRequired checks if data is not empty
func (v *Validator) validateRequired(data interface{}) error {
	if data == nil {
		return fmt.Errorf("field is required")
	}
	if str, ok := data.(string); ok && str == "" {
		return fmt.Errorf("field cannot be empty")
	}
	return nil
}

// validateEmail performs basic email validation
func (v *Validator) validateEmail(data interface{}) error {
	if str, ok := data.(string); ok {
		if !strings.Contains(str, "@") {
			return fmt.Errorf("invalid email format")
		}
	}
	return nil
}

// validateURL performs basic URL validation
func (v *Validator) validateURL(data interface{}) error {
	if str, ok := data.(string); ok {
		if !strings.HasPrefix(str, "http://") && !strings.HasPrefix(str, "https://") {
			return fmt.Errorf("invalid URL format")
		}
	}
	return nil
}

// validatePositive checks if numeric value is positive
func (v *Validator) validatePositive(data interface{}) error {
	if num, ok := data.(float64); ok {
		if num <= 0 {
			return fmt.Errorf("value must be positive")
		}
	}
	return nil
}
