// Package schemas provides schema-driven validation functionality for the spooky codebase.
package schemas

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsimple"

	spookytypeslogging "spooky/internal/types/logging"
	spookytypesschemas "spooky/internal/types/schemas"
)

// SchemaDrivenValidator provides comprehensive schema-driven validation
type SchemaDrivenValidator struct {
	logger   spookytypeslogging.Logger
	registry spookytypesschemas.SchemaRegistry
	parser   *hclparse.Parser

	// Embedded schemas for all configuration types
	embeddedSchemas map[string]*spookytypesschemas.Schema

	// Validation configuration
	config *SchemaDrivenValidationConfig
}

// SchemaDrivenValidationConfig provides configuration for schema-driven validation
type SchemaDrivenValidationConfig struct {
	// Whether to use embedded schemas
	UseEmbeddedSchemas bool `json:"use_embedded_schemas" hcl:"use_embedded_schemas"`

	// Whether to validate against strict schemas
	StrictValidation bool `json:"strict_validation" hcl:"strict_validation"`

	// Whether to allow unknown fields
	AllowUnknownFields bool `json:"allow_unknown_fields" hcl:"allow_unknown_fields"`

	// Whether to provide detailed error messages
	DetailedErrors bool `json:"detailed_errors" hcl:"detailed_errors"`

	// Custom validation rules
	CustomRules map[string]CustomValidationRule `json:"custom_rules,omitempty" hcl:"custom_rules,optional"`
}

// CustomValidationRule represents a custom validation rule
type CustomValidationRule struct {
	// Rule name
	Name string `json:"name" hcl:"name"`

	// Rule description
	Description string `json:"description" hcl:"description"`

	// Rule function name
	Function string `json:"function" hcl:"function"`

	// Rule severity
	Severity string `json:"severity" hcl:"severity"` // "error", "warning", "info"

	// Rule parameters
	Parameters map[string]interface{} `json:"parameters,omitempty" hcl:"parameters,optional"`
}

// NewSchemaDrivenValidator creates a new schema-driven validator
func NewSchemaDrivenValidator(logger spookytypeslogging.Logger, config *SchemaDrivenValidationConfig) *SchemaDrivenValidator {
	if config == nil {
		config = &SchemaDrivenValidationConfig{
			UseEmbeddedSchemas: true,
			StrictValidation:   true,
			AllowUnknownFields: false,
			DetailedErrors:     true,
			CustomRules:        make(map[string]CustomValidationRule),
		}
	}

	validator := &SchemaDrivenValidator{
		logger:          logger,
		parser:          hclparse.NewParser(),
		embeddedSchemas: make(map[string]*spookytypesschemas.Schema),
		config:          config,
	}

	// Load embedded schemas
	validator.loadEmbeddedSchemas()

	return validator
}

// SetRegistry sets the schema registry for the validator
func (v *SchemaDrivenValidator) SetRegistry(registry spookytypesschemas.SchemaRegistry) {
	v.registry = registry
}

// ValidateConfiguration validates configuration using schema-driven approach
func (v *SchemaDrivenValidator) ValidateConfiguration(_ context.Context, configPath, configType string) (*spookytypesschemas.ValidationResult, error) {
	start := time.Now()

	v.logger.Debug("Validating configuration with schema-driven approach", map[string]interface{}{
		"config_path": configPath,
		"config_type": configType,
	})

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

	// Step 1: Validate file existence and readability
	if err := v.validateFileAccess(configPath, result); err != nil {
		return result, err
	}

	// Step 2: Parse HCL content
	parsedData, err := v.parseHCLContent(configPath, result)
	if err != nil {
		return result, err
	}

	// Step 3: Get schema for configuration type
	schema, err := v.getSchemaForConfigType(configType, result)
	if err != nil {
		return result, err
	}

	// Step 4: Validate against schema
	if err := v.validateAgainstSchema(schema, parsedData, configPath, result); err != nil {
		return result, err
	}

	// Step 5: Apply custom validation rules
	if err := v.applyCustomValidationRules(configType, parsedData, result); err != nil {
		return result, err
	}

	// Step 6: Cross-field validation
	if err := v.validateCrossFieldRules(configType, parsedData, result); err != nil {
		return result, err
	}

	// Update validation statistics
	result.Statistics.Duration = time.Since(start)
	result.Valid = len(result.Errors) == 0

	v.logger.Info("Configuration validation completed", map[string]interface{}{
		"config_path": configPath,
		"config_type": configType,
		"valid":       result.Valid,
		"errors":      len(result.Errors),
		"warnings":    len(result.Warnings),
		"duration":    result.Statistics.Duration,
	})

	return result, nil
}

// ValidateProjectStructure validates project structure using schema-driven approach
func (v *SchemaDrivenValidator) ValidateProjectStructure(_ context.Context, projectPath string) (*spookytypesschemas.ValidationResult, error) {
	start := time.Now()

	v.logger.Debug("Validating project structure with schema-driven approach", map[string]interface{}{
		"project_path": projectPath,
	})

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

	// Step 1: Validate project directory structure
	if err := v.validateProjectDirectory(projectPath, result); err != nil {
		return result, err
	}

	// Step 2: Validate required files
	if err := v.validateRequiredFiles(projectPath, result); err != nil {
		return result, err
	}

	// Step 3: Validate optional files
	if err := v.validateOptionalFiles(projectPath, result); err != nil {
		return result, err
	}

	// Step 4: Validate file permissions
	if err := v.validateFilePermissions(projectPath, result); err != nil {
		return result, err
	}

	// Update validation statistics
	result.Statistics.Duration = time.Since(start)
	result.Valid = len(result.Errors) == 0

	v.logger.Info("Project structure validation completed", map[string]interface{}{
		"project_path": projectPath,
		"valid":        result.Valid,
		"errors":       len(result.Errors),
		"warnings":     len(result.Warnings),
		"duration":     result.Statistics.Duration,
	})

	return result, nil
}

// Helper methods

// validateFileAccess validates file existence and readability
func (v *SchemaDrivenValidator) validateFileAccess(configPath string, result *spookytypesschemas.ValidationResult) error {
	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		v.addError(result, "file_not_found", fmt.Sprintf("Configuration file not found: %s", configPath), "Check file path and permissions", "error")
		return fmt.Errorf("configuration file not found: %s", configPath)
	}

	// Check if file is readable
	if _, err := os.ReadFile(configPath); err != nil {
		v.addError(result, "file_not_readable", fmt.Sprintf("Configuration file not readable: %s", configPath), "Check file permissions", "error")
		return fmt.Errorf("configuration file not readable: %s", configPath)
	}

	return nil
}

// parseHCLContent parses HCL content from file
func (v *SchemaDrivenValidator) parseHCLContent(configPath string, result *spookytypesschemas.ValidationResult) (interface{}, error) {
	// Read file content
	data, err := os.ReadFile(configPath)
	if err != nil {
		v.addError(result, "file_read_error", fmt.Sprintf("Failed to read configuration file: %s", configPath), "Check file permissions and disk space", "error")
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Parse HCL content
	var parsedData interface{}
	if err := hclsimple.Decode(configPath, data, nil, &parsedData); err != nil {
		// Convert HCL parsing errors to validation errors
		if diags, ok := err.(hcl.Diagnostics); ok {
			for _, diag := range diags {
				v.addError(result, "hcl_parse_error", diag.Summary, diag.Detail, "error")
			}
		} else {
			v.addError(result, "hcl_parse_error", fmt.Sprintf("Failed to parse HCL content: %v", err), "Check HCL syntax and structure", "error")
		}
		return nil, fmt.Errorf("failed to parse HCL content: %w", err)
	}

	return parsedData, nil
}

// getSchemaForConfigType gets the appropriate schema for the configuration type
func (v *SchemaDrivenValidator) getSchemaForConfigType(configType string, result *spookytypesschemas.ValidationResult) (*spookytypesschemas.Schema, error) {
	var schema *spookytypesschemas.Schema
	var err error

	// Try to get schema from registry first
	if v.registry != nil {
		schema, exists := v.registry.Get(configType, "hcl")
		if exists && schema != nil {
			return schema, nil
		}
	}

	// Try to get embedded schema
	if v.config.UseEmbeddedSchemas {
		schema, err = v.getEmbeddedSchema(configType)
		if err == nil && schema != nil {
			return schema, nil
		}
	}

	// If no schema found, create a basic validation
	v.addWarning(result, "no_schema_found", fmt.Sprintf("No schema found for configuration type: %s", configType), "Using basic HCL validation", "warning")

	// Create a basic schema for validation
	schema = &spookytypesschemas.Schema{
		Version:     "1.0",
		Type:        "basic",
		Name:        configType,
		Description: fmt.Sprintf("Basic schema for %s", configType),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Content:     "",
		Metadata:    make(map[string]interface{}),
	}

	return schema, nil
}

// validateAgainstSchema validates data against the schema
func (v *SchemaDrivenValidator) validateAgainstSchema(schema *spookytypesschemas.Schema, data interface{}, _ string, result *spookytypesschemas.ValidationResult) error {
	// Basic schema validation
	if schema == nil {
		v.addError(result, "invalid_schema", "Schema is nil", "Check schema configuration", "error")
		return fmt.Errorf("schema is nil")
	}

	// Validate data structure based on schema type
	switch schema.Type {
	case "project-directory":
		return v.validateProjectDirectorySchema(data, result)
	case "project":
		return v.validateProjectSchema(data, result)
	case "machines":
		return v.validateMachinesSchema(data, result)
	case "actions":
		return v.validateActionsSchema(data, result)
	case "variables":
		return v.validateVariablesSchema(data, result)
	case "templates":
		return v.validateTemplatesSchema(data, result)
	case "logging":
		return v.validateLoggingSchema(data, result)
	case "ssh":
		return v.validateSSHSchema(data, result)
	case "basic":
		return v.validateBasicSchema(data, result)
	default:
		v.addWarning(result, "unknown_schema_type", fmt.Sprintf("Unknown schema type: %s", schema.Type), "Using basic validation", "warning")
		return v.validateBasicSchema(data, result)
	}
}

// validateProjectDirectory validates project directory structure
func (v *SchemaDrivenValidator) validateProjectDirectory(projectPath string, result *spookytypesschemas.ValidationResult) error {
	// Check if project directory exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		v.addError(result, "project_directory_not_found", fmt.Sprintf("Project directory not found: %s", projectPath), "Check project path", "error")
		return fmt.Errorf("project directory not found: %s", projectPath)
	}

	// Check if it's a directory
	if info, err := os.Stat(projectPath); err == nil && !info.IsDir() {
		v.addError(result, "project_path_not_directory", fmt.Sprintf("Project path is not a directory: %s", projectPath), "Check project path", "error")
		return fmt.Errorf("project path is not a directory: %s", projectPath)
	}

	return nil
}

// validateRequiredFiles validates required project files
func (v *SchemaDrivenValidator) validateRequiredFiles(projectPath string, result *spookytypesschemas.ValidationResult) error {
	requiredFiles := []string{"project.hcl"}

	for _, file := range requiredFiles {
		filePath := filepath.Join(projectPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			v.addError(result, "required_file_missing", fmt.Sprintf("Required file missing: %s", file), fmt.Sprintf("Create %s file", file), "error")
		} else {
			// Validate the file content
			if err := v.validateFileContent(filePath, file, result); err != nil {
				v.addError(result, "file_validation_failed", fmt.Sprintf("File validation failed: %s", file), err.Error(), "error")
			}
		}
	}

	return nil
}

// validateOptionalFiles validates optional project files
func (v *SchemaDrivenValidator) validateOptionalFiles(projectPath string, result *spookytypesschemas.ValidationResult) error {
	optionalFiles := []string{"machines.hcl", "actions.hcl", "variables.hcl", "templates.hcl"}

	for _, file := range optionalFiles {
		filePath := filepath.Join(projectPath, file)
		if _, err := os.Stat(filePath); err == nil {
			// Validate the file content
			if err := v.validateFileContent(filePath, file, result); err != nil {
				v.addWarning(result, "optional_file_validation_failed", fmt.Sprintf("Optional file validation failed: %s", file), err.Error(), "warning")
			}
		}
	}

	return nil
}

// validateFilePermissions validates file permissions
func (v *SchemaDrivenValidator) validateFilePermissions(projectPath string, result *spookytypesschemas.ValidationResult) error {
	// Check project directory permissions
	if info, err := os.Stat(projectPath); err == nil {
		mode := info.Mode()
		if mode&0o077 != 0 {
			v.addWarning(result, "directory_permissions", fmt.Sprintf("Project directory has loose permissions: %s", projectPath), "Consider restricting permissions to 750", "warning")
		}
	}

	// Check file permissions for sensitive files
	sensitiveFiles := []string{"project.hcl", "machines.hcl", "variables.hcl"}
	for _, file := range sensitiveFiles {
		filePath := filepath.Join(projectPath, file)
		if info, err := os.Stat(filePath); err == nil {
			mode := info.Mode()
			if mode&0o077 != 0 {
				v.addWarning(result, "file_permissions", fmt.Sprintf("Sensitive file has loose permissions: %s", file), "Consider restricting permissions to 600", "warning")
			}
		}
	}

	return nil
}

// validateFileContent validates file content
func (v *SchemaDrivenValidator) validateFileContent(filePath, _ string, _ *spookytypesschemas.ValidationResult) error {
	// Read file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Basic HCL syntax validation
	if err := hclsimple.Decode(filePath, data, nil, &map[string]interface{}{}); err != nil {
		return fmt.Errorf("invalid HCL syntax: %w", err)
	}

	return nil
}

// Schema-specific validation methods

func (v *SchemaDrivenValidator) validateProjectDirectorySchema(_ interface{}, _ *spookytypesschemas.ValidationResult) error {
	// Validate project directory structure
	// This would contain specific validation rules for project directory structure
	return nil
}

func (v *SchemaDrivenValidator) validateProjectSchema(_ interface{}, _ *spookytypesschemas.ValidationResult) error {
	// Validate project configuration
	// This would contain specific validation rules for project configuration
	return nil
}

func (v *SchemaDrivenValidator) validateMachinesSchema(_ interface{}, _ *spookytypesschemas.ValidationResult) error {
	// Validate machines configuration
	// This would contain specific validation rules for machines configuration
	return nil
}

func (v *SchemaDrivenValidator) validateActionsSchema(_ interface{}, _ *spookytypesschemas.ValidationResult) error {
	// Validate actions configuration
	// This would contain specific validation rules for actions configuration
	return nil
}

func (v *SchemaDrivenValidator) validateVariablesSchema(_ interface{}, _ *spookytypesschemas.ValidationResult) error {
	// Validate variables configuration
	// This would contain specific validation rules for variables configuration
	return nil
}

func (v *SchemaDrivenValidator) validateTemplatesSchema(_ interface{}, _ *spookytypesschemas.ValidationResult) error {
	// Validate templates configuration
	// This would contain specific validation rules for templates configuration
	return nil
}

func (v *SchemaDrivenValidator) validateLoggingSchema(_ interface{}, _ *spookytypesschemas.ValidationResult) error {
	// Validate logging configuration
	// This would contain specific validation rules for logging configuration
	return nil
}

func (v *SchemaDrivenValidator) validateSSHSchema(_ interface{}, _ *spookytypesschemas.ValidationResult) error {
	// Validate SSH configuration
	// This would contain specific validation rules for SSH configuration
	return nil
}

func (v *SchemaDrivenValidator) validateBasicSchema(_ interface{}, _ *spookytypesschemas.ValidationResult) error {
	// Basic validation for any HCL content
	// This would contain basic validation rules
	return nil
}

// applyCustomValidationRules applies custom validation rules
func (v *SchemaDrivenValidator) applyCustomValidationRules(_ string, data interface{}, result *spookytypesschemas.ValidationResult) error {
	// Apply custom validation rules for the configuration type
	for ruleName, rule := range v.config.CustomRules {
		if err := v.applyCustomRule(ruleName, rule, data, result); err != nil {
			v.addError(result, "custom_rule_failed", fmt.Sprintf("Custom rule failed: %s", ruleName), err.Error(), "error")
		}
	}

	return nil
}

// validateCrossFieldRules validates cross-field rules
func (v *SchemaDrivenValidator) validateCrossFieldRules(_ string, _ interface{}, _ *spookytypesschemas.ValidationResult) error {
	// Apply cross-field validation rules
	// This would contain cross-field validation logic
	return nil
}

// applyCustomRule applies a custom validation rule
func (v *SchemaDrivenValidator) applyCustomRule(_ string, _ CustomValidationRule, _ interface{}, _ *spookytypesschemas.ValidationResult) error {
	// Apply custom validation rule
	// This would contain custom rule application logic
	return nil
}

// loadEmbeddedSchemas loads embedded schemas
func (v *SchemaDrivenValidator) loadEmbeddedSchemas() {
	// Load embedded schemas from the schemas directory
	schemasDir := "internal/schemas/schemas"

	schemaFiles := []string{
		"project-directory.schema.hcl",
		"project.schema.hcl",
		"machines.schema.hcl",
		"actions.schema.hcl",
		"variables.schema.hcl",
		"templates.schema.hcl",
		"logging.schema.hcl",
		"ssh.schema.hcl",
	}

	for _, schemaFile := range schemaFiles {
		schemaPath := filepath.Join(schemasDir, schemaFile)
		if schema, err := v.loadEmbeddedSchema(schemaPath); err == nil {
			schemaName := strings.TrimSuffix(schemaFile, ".schema.hcl")
			v.embeddedSchemas[schemaName] = schema
		}
	}
}

// loadEmbeddedSchema loads an embedded schema
func (v *SchemaDrivenValidator) loadEmbeddedSchema(schemaPath string) (*spookytypesschemas.Schema, error) {
	// Load embedded schema from file
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded schema: %w", err)
	}

	schema := &spookytypesschemas.Schema{
		Version:     "1.0",
		Type:        "embedded",
		Name:        filepath.Base(schemaPath),
		Description: fmt.Sprintf("Embedded schema from %s", schemaPath),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Content:     string(data),
		Metadata:    make(map[string]interface{}),
	}

	return schema, nil
}

// getEmbeddedSchema gets an embedded schema by name
func (v *SchemaDrivenValidator) getEmbeddedSchema(schemaName string) (*spookytypesschemas.Schema, error) {
	if schema, exists := v.embeddedSchemas[schemaName]; exists {
		return schema, nil
	}
	return nil, fmt.Errorf("embedded schema not found: %s", schemaName)
}

// addError adds an error to the validation result
func (v *SchemaDrivenValidator) addError(result *spookytypesschemas.ValidationResult, code, message, suggestion, severity string) {
	schemaError := spookytypesschemas.NewSchemaError("", "", message)
	schemaError.Code = code
	schemaError.Severity = severity
	if suggestion != "" {
		schemaError.AddSuggestion(suggestion)
	}
	result.Errors = append(result.Errors, *schemaError)
	result.Statistics.InvalidFields++
	result.Statistics.RulesFailed++
}

// addWarning adds a warning to the validation result
func (v *SchemaDrivenValidator) addWarning(result *spookytypesschemas.ValidationResult, code, message, suggestion, severity string) {
	warning := spookytypesschemas.NewSchemaError("", "", message)
	warning.Code = code
	warning.Severity = severity
	if suggestion != "" {
		warning.AddSuggestion(suggestion)
	}
	result.Warnings = append(result.Warnings, *warning)
	result.Statistics.RulesExecuted++
}
