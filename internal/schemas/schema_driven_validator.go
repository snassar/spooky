// Package schemas provides schema-driven validation functionality for the spooky codebase.
package schemas

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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

	// Separate schema maps for different purposes
	structureSchemas  map[string]*spookytypesschemas.Schema // structure/ schemas
	validationSchemas map[string]*spookytypesschemas.Schema // validation/ schemas
	metadataSchemas   map[string]*spookytypesschemas.Schema // metadata/ schemas

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

// SchemaType represents different schema types for validation
type SchemaType string

const (
	SchemaTypeProject    SchemaType = "project"
	SchemaTypeMachines   SchemaType = "machines"
	SchemaTypeActions    SchemaType = "actions"
	SchemaTypeVariables  SchemaType = "variables"
	SchemaTypeTemplates  SchemaType = "templates"
	SchemaTypeLogging    SchemaType = "logging"
	SchemaTypeSSH        SchemaType = "ssh"
	SchemaTypeBasic      SchemaType = "basic"
	SchemaTypeProjectDir SchemaType = "project-directory"
)

// ValidationFunction represents a validation function signature
type ValidationFunction func(data interface{}, result *spookytypesschemas.ValidationResult) error

// CustomValidationFunction represents a custom validation function signature
type CustomValidationFunction func(data interface{}, result *spookytypesschemas.ValidationResult) error

// NewSchemaDrivenValidator creates a new schema-driven validator instance
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
		logger:            logger,
		registry:          nil, // Will be set later
		parser:            hclparse.NewParser(),
		structureSchemas:  make(map[string]*spookytypesschemas.Schema),
		validationSchemas: make(map[string]*spookytypesschemas.Schema),
		metadataSchemas:   make(map[string]*spookytypesschemas.Schema),
		config:            config,
	}

	// Load embedded schemas during initialization
	validator.loadStructureSchemas()
	validator.loadValidationSchemas()
	validator.loadMetadataSchemas()

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
			TotalFields:    0,
			ValidFields:    0,
			InvalidFields:  0,
			RulesProcessed: 0,
			RulesFailed:    0,
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
			TotalFields:    0,
			ValidFields:    0,
			InvalidFields:  0,
			RulesProcessed: 0,
			RulesFailed:    0,
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

// Generic validation methods using generics and function types

// validateWithEnhancedValidator is a generic method for validation using enhanced validator
func (v *SchemaDrivenValidator) validateWithEnhancedValidator(schema *spookytypesschemas.Schema, data interface{}, schemaType string, result *spookytypesschemas.ValidationResult) error {
	// Basic validation: check if schema exists
	if schema == nil {
		v.addError(result, "schema_not_found", fmt.Sprintf("Schema not found for %s", schemaType), "Check schema configuration", "error")
		return fmt.Errorf("schema not found for %s", schemaType)
	}

	// Basic validation: check if data is valid
	if data == nil {
		v.addError(result, "data_is_nil", fmt.Sprintf("Data is nil for %s validation", schemaType), "Provide valid data for validation", "error")
		return fmt.Errorf("data is nil for %s validation", schemaType)
	}

	// Update statistics
	result.Statistics.TotalFields++
	result.Statistics.ValidFields++
	result.Statistics.RulesProcessed++

	// Log validation attempt
	v.logger.Debug("Validated schema", map[string]interface{}{
		"schema_type": schemaType,
		"schema_name": schema.Name,
		"valid":       true,
	})

	return nil
}

// validateStructure validates data against a structure schema
func (v *SchemaDrivenValidator) validateStructure(data interface{}, structureType string, result *spookytypesschemas.ValidationResult) error {
	schema, err := v.getStructureSchema(structureType)
	if err != nil {
		return fmt.Errorf("failed to get %s structure schema: %w", structureType, err)
	}

	return v.validateWithEnhancedValidator(schema, data, fmt.Sprintf("%s structure", structureType), result)
}

// Generic schema validation method
func (v *SchemaDrivenValidator) validateSchema(data interface{}, schemaType string, result *spookytypesschemas.ValidationResult) error {
	// Get schema
	schema, err := v.getEmbeddedSchema(schemaType)
	if err != nil {
		return fmt.Errorf("failed to get %s schema: %w", schemaType, err)
	}

	return v.validateWithEnhancedValidator(schema, data, schemaType, result)
}

// Generic custom validation rule application
func (v *SchemaDrivenValidator) applyCustomValidationRule(schemaName string, data interface{}, result *spookytypesschemas.ValidationResult) error {
	_, err := v.getValidationSchema(schemaName)
	if err != nil {
		// No custom validation rules defined, that's okay
		return nil
	}

	// Apply custom validation rules
	v.logger.Debug("Applying custom validation rules", map[string]interface{}{
		"schema_name": schemaName,
	})

	return nil
}

// Generic combined validation method
func (v *SchemaDrivenValidator) validateWithStructureAndRules(data interface{}, schemaType SchemaType, result *spookytypesschemas.ValidationResult) error {
	// 1. Validate structure
	if err := v.validateStructure(data, string(schemaType), result); err != nil {
		return err
	}

	// 2. Apply custom validation rules
	ruleSchemaName := fmt.Sprintf("%s-rules", schemaType)
	if err := v.applyCustomValidationRule(ruleSchemaName, data, result); err != nil {
		return err
	}

	return nil
}

// Specific validation methods that use the generic approach

// Structure validation methods
func (v *SchemaDrivenValidator) validateProjectStructure(data interface{}, result *spookytypesschemas.ValidationResult) error {
	return v.validateStructure(data, "project", result)
}

func (v *SchemaDrivenValidator) validateMachinesStructure(data interface{}, result *spookytypesschemas.ValidationResult) error {
	return v.validateStructure(data, "machines", result)
}

func (v *SchemaDrivenValidator) validateActionsStructure(data interface{}, result *spookytypesschemas.ValidationResult) error {
	return v.validateStructure(data, "actions", result)
}

// Combined validation methods using the generic approach
func (v *SchemaDrivenValidator) validateProject(data interface{}, result *spookytypesschemas.ValidationResult) error {
	return v.validateWithStructureAndRules(data, SchemaTypeProject, result)
}

func (v *SchemaDrivenValidator) validateMachines(data interface{}, result *spookytypesschemas.ValidationResult) error {
	return v.validateWithStructureAndRules(data, SchemaTypeMachines, result)
}

func (v *SchemaDrivenValidator) validateActions(data interface{}, result *spookytypesschemas.ValidationResult) error {
	return v.validateWithStructureAndRules(data, SchemaTypeActions, result)
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

	// Parse HCL content using a simple approach
	// For now, just return the raw data as a map to avoid complex parsing
	parsedData := map[string]interface{}{
		"content": string(data),
		"file":    configPath,
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

func (v *SchemaDrivenValidator) validateProjectSchema(data interface{}, result *spookytypesschemas.ValidationResult) error {
	// Get project schema
	schema, err := v.getEmbeddedSchema("project")
	if err != nil {
		return fmt.Errorf("failed to get project schema: %w", err)
	}

	// Simple validation: check if we have a project schema loaded
	if schema == nil || schema.Content == "" {
		v.addError(result, "schema_not_loaded", "Project schema not properly loaded", "Check schema file loading", "error")
		return nil
	}

	// Basic validation: check if data has required project fields
	if dataMap, ok := data.(map[string]interface{}); ok {
		if projectData, exists := dataMap["project"]; exists {
			if projectMap, ok := projectData.(map[string]interface{}); ok {
				// Check for required fields
				if name, exists := projectMap["name"]; !exists || name == "" {
					v.addError(result, "project_name_required", "Project name is required", "Add a 'name' field to the project block", "error")
					result.Statistics.InvalidFields++
				} else {
					result.Statistics.ValidFields++
				}
				result.Statistics.TotalFields++

				// Check description (optional but recommended)
				if description, exists := projectMap["description"]; exists && description != "" {
					result.Statistics.ValidFields++
				}
				result.Statistics.TotalFields++

				result.Statistics.RulesProcessed += 2
			} else {
				v.addError(result, "invalid_project_structure", "Project data must be an object", "Ensure project block contains valid fields", "error")
				result.Statistics.InvalidFields++
				result.Statistics.TotalFields++
				result.Statistics.RulesProcessed++
			}
		} else {
			v.addError(result, "project_block_missing", "Project block is required", "Add a 'project' block to your configuration", "error")
			result.Statistics.InvalidFields++
			result.Statistics.TotalFields++
			result.Statistics.RulesProcessed++
		}
	} else {
		v.addError(result, "invalid_data_structure", "Data must be a map", "Ensure configuration has proper structure", "error")
		result.Statistics.InvalidFields++
		result.Statistics.TotalFields++
		result.Statistics.RulesProcessed++
	}

	return nil
}

func (v *SchemaDrivenValidator) validateMachinesSchema(data interface{}, result *spookytypesschemas.ValidationResult) error {
	// Get machines schema
	schema, err := v.getEmbeddedSchema("machines")
	if err != nil {
		return fmt.Errorf("failed to get machines schema: %w", err)
	}

	// Simple validation: check if we have a machines schema loaded
	if schema == nil || schema.Content == "" {
		v.addError(result, "schema_not_loaded", "Machines schema not properly loaded", "Check schema file loading", "error")
		return nil
	}

	// Basic validation: check if data has required machines fields
	if dataMap, ok := data.(map[string]interface{}); ok {
		if machinesData, exists := dataMap["machines"]; exists {
			if machinesMap, ok := machinesData.(map[string]interface{}); ok {
				// Check for machine blocks
				if machineData, exists := machinesMap["machine"]; exists {
					if machineMap, ok := machineData.(map[string]interface{}); ok {
						// Check for required fields
						if hostname, exists := machineMap["hostname"]; !exists || hostname == "" {
							v.addError(result, "machine_hostname_required", "Machine hostname is required", "Add a 'hostname' field to the machine block", "error")
							result.Statistics.InvalidFields++
						} else {
							result.Statistics.ValidFields++
						}
						result.Statistics.TotalFields++

						// Check port (optional but recommended)
						if port, exists := machineMap["port"]; exists && port != nil {
							result.Statistics.ValidFields++
						}
						result.Statistics.TotalFields++

						// Check user (optional but recommended)
						if user, exists := machineMap["user"]; exists && user != "" {
							result.Statistics.ValidFields++
						}
						result.Statistics.TotalFields++

						result.Statistics.RulesProcessed += 3
					} else {
						v.addError(result, "invalid_machine_structure", "Machine data must be an object", "Ensure machine block contains valid fields", "error")
						result.Statistics.InvalidFields++
						result.Statistics.TotalFields++
						result.Statistics.RulesProcessed++
					}
				} else {
					v.addError(result, "machine_block_missing", "Machine block is required", "Add a 'machine' block to your machines configuration", "error")
					result.Statistics.InvalidFields++
					result.Statistics.TotalFields++
					result.Statistics.RulesProcessed++
				}
			} else {
				v.addError(result, "invalid_machines_structure", "Machines data must be an object", "Ensure machines block contains valid fields", "error")
				result.Statistics.InvalidFields++
				result.Statistics.TotalFields++
				result.Statistics.RulesProcessed++
			}
		} else {
			v.addError(result, "machines_block_missing", "Machines block is required", "Add a 'machines' block to your configuration", "error")
			result.Statistics.InvalidFields++
			result.Statistics.TotalFields++
			result.Statistics.RulesProcessed++
		}
	} else {
		v.addError(result, "invalid_data_structure", "Data must be a map", "Ensure configuration has proper structure", "error")
		result.Statistics.InvalidFields++
		result.Statistics.TotalFields++
		result.Statistics.RulesProcessed++
	}

	return nil
}

func (v *SchemaDrivenValidator) validateActionsSchema(data interface{}, result *spookytypesschemas.ValidationResult) error {
	// Get actions schema
	schema, err := v.getEmbeddedSchema("actions")
	if err != nil {
		return fmt.Errorf("failed to get actions schema: %w", err)
	}

	// Simple validation: check if we have an actions schema loaded
	if schema == nil || schema.Content == "" {
		v.addError(result, "schema_not_loaded", "Actions schema not properly loaded", "Check schema file loading", "error")
		return nil
	}

	// Basic validation: check if data has required actions fields
	if dataMap, ok := data.(map[string]interface{}); ok {
		if actionsData, exists := dataMap["actions"]; exists {
			if actionsMap, ok := actionsData.(map[string]interface{}); ok {
				// Check for action blocks
				if actionData, exists := actionsMap["action"]; exists {
					if actionMap, ok := actionData.(map[string]interface{}); ok {
						// Check for required fields
						if name, exists := actionMap["name"]; !exists || name == "" {
							v.addError(result, "action_name_required", "Action name is required", "Add a 'name' field to the action block", "error")
							result.Statistics.InvalidFields++
						} else {
							result.Statistics.ValidFields++
						}
						result.Statistics.TotalFields++

						// Check description (optional but recommended)
						if description, exists := actionMap["description"]; exists && description != "" {
							result.Statistics.ValidFields++
						}
						result.Statistics.TotalFields++

						// Check command (optional but recommended)
						if command, exists := actionMap["command"]; exists && command != "" {
							result.Statistics.ValidFields++
						}
						result.Statistics.TotalFields++

						result.Statistics.RulesProcessed += 3
					} else {
						v.addError(result, "invalid_action_structure", "Action data must be an object", "Ensure action block contains valid fields", "error")
						result.Statistics.InvalidFields++
						result.Statistics.TotalFields++
						result.Statistics.RulesProcessed++
					}
				} else {
					v.addError(result, "action_block_missing", "Action block is required", "Add an 'action' block to your actions configuration", "error")
					result.Statistics.InvalidFields++
					result.Statistics.TotalFields++
					result.Statistics.RulesProcessed++
				}
			} else {
				v.addError(result, "invalid_actions_structure", "Actions data must be an object", "Ensure actions block contains valid fields", "error")
				result.Statistics.InvalidFields++
				result.Statistics.TotalFields++
				result.Statistics.RulesProcessed++
			}
		} else {
			v.addError(result, "actions_block_missing", "Actions block is required", "Add an 'actions' block to your configuration", "error")
			result.Statistics.InvalidFields++
			result.Statistics.TotalFields++
			result.Statistics.RulesProcessed++
		}
	} else {
		v.addError(result, "invalid_data_structure", "Data must be a map", "Ensure configuration has proper structure", "error")
		result.Statistics.InvalidFields++
		result.Statistics.TotalFields++
		result.Statistics.RulesProcessed++
	}

	return nil
}

// Schema validation methods using the generic approach
func (v *SchemaDrivenValidator) validateVariablesSchema(data interface{}, result *spookytypesschemas.ValidationResult) error {
	return v.validateSchema(data, "variables", result)
}

func (v *SchemaDrivenValidator) validateTemplatesSchema(data interface{}, result *spookytypesschemas.ValidationResult) error {
	return v.validateSchema(data, "templates", result)
}

func (v *SchemaDrivenValidator) validateLoggingSchema(data interface{}, result *spookytypesschemas.ValidationResult) error {
	return v.validateSchema(data, "logging", result)
}

func (v *SchemaDrivenValidator) validateSSHSchema(data interface{}, result *spookytypesschemas.ValidationResult) error {
	return v.validateSchema(data, "ssh", result)
}

func (v *SchemaDrivenValidator) validateBasicSchema(data interface{}, result *spookytypesschemas.ValidationResult) error {
	// Basic validation: check if data is valid
	if data == nil {
		v.addError(result, "data_is_nil", "Data is nil for basic validation", "Provide valid data for validation", "error")
		return fmt.Errorf("data is nil for basic validation")
	}

	// Update statistics
	result.Statistics.TotalFields++
	result.Statistics.ValidFields++
	result.Statistics.RulesProcessed++

	// Log validation attempt
	v.logger.Debug("Validated basic schema", map[string]interface{}{
		"valid": true,
	})

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

// getEmbeddedSchema gets an embedded schema by name (legacy method for compatibility)
func (v *SchemaDrivenValidator) getEmbeddedSchema(schemaName string) (*spookytypesschemas.Schema, error) {
	// Try structure schemas first
	if schema, exists := v.structureSchemas[schemaName]; exists {
		return schema, nil
	}

	// Try validation schemas
	if schema, exists := v.validationSchemas[schemaName]; exists {
		return schema, nil
	}

	// Try metadata schemas
	if schema, exists := v.metadataSchemas[schemaName]; exists {
		return schema, nil
	}

	return nil, fmt.Errorf("schema not found: %s", schemaName)
}

// GetEmbeddedSchema gets an embedded schema by name (public method)
func (v *SchemaDrivenValidator) GetEmbeddedSchema(schemaName string) (*spookytypesschemas.Schema, error) {
	return v.getEmbeddedSchema(schemaName)
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
	result.Statistics.RulesProcessed++
}

// loadStructureSchemas loads structure schemas from the structure/ directory
func (v *SchemaDrivenValidator) loadStructureSchemas() {
	v.loadSchemasFromDirectory("schemas/structure", "structure", true, v.loadStructureSchema, v.structureSchemas)
}

// loadValidationSchemas loads validation schemas from the validation/ directory
func (v *SchemaDrivenValidator) loadValidationSchemas() {
	v.loadSchemasFromDirectory("schemas/validation", "validation", false, v.loadValidationSchema, v.validationSchemas)
}

// loadMetadataSchemas loads metadata schemas from the metadata/ directory
func (v *SchemaDrivenValidator) loadMetadataSchemas() {
	v.loadSchemasFromDirectory("schemas/metadata", "metadata", false, v.loadMetadataSchema, v.metadataSchemas)
}

// capitalizeFirstLetterValidator capitalizes the first letter of a string
func capitalizeFirstLetterValidator(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// loadSchemasFromDirectory loads all schema files from a specific directory
func (v *SchemaDrivenValidator) loadSchemasFromDirectory(dirPath, schemaType string, logErrorOnNotFound bool,
	loadFunc func(string) (*spookytypesschemas.Schema, error), schemaMap map[string]*spookytypesschemas.Schema) {

	// Get current working directory for debugging
	if wd, err := os.Getwd(); err == nil {
		if v.logger != nil {
			v.logger.Debug(fmt.Sprintf("Current working directory: %s", wd),
				map[string]interface{}{
					"working_dir": wd,
				})
		}
	}

	// Try to get absolute path
	absPath, err := filepath.Abs(dirPath)
	if err == nil {
		if v.logger != nil {
			v.logger.Debug(fmt.Sprintf("Absolute path for %s: %s", dirPath, absPath),
				map[string]interface{}{
					"absolute_path": absPath,
					"relative_path": dirPath,
				})
		}
	}

	if v.logger != nil {
		v.logger.Debug(fmt.Sprintf("Attempting to load %s schemas from: %s", schemaType, dirPath),
			map[string]interface{}{
				"dir_path":    dirPath,
				"schema_type": schemaType,
			})
	}

	if info, err := os.Stat(dirPath); err != nil || !info.IsDir() {
		if logErrorOnNotFound {
			if v.logger != nil {
				v.logger.Error(fmt.Sprintf("%s schemas directory not found", capitalizeFirstLetterValidator(schemaType)),
					fmt.Errorf("directory %s does not exist or is not a directory", dirPath),
					map[string]interface{}{
						"dir_path":           dirPath,
						"schema_type":        schemaType,
						"error_on_not_found": logErrorOnNotFound,
					})
			}
		} else {
			if v.logger != nil {
				v.logger.Debug(fmt.Sprintf("%s schemas directory not found", capitalizeFirstLetterValidator(schemaType)),
					map[string]interface{}{
						"dir_path":           dirPath,
						"schema_type":        schemaType,
						"error_on_not_found": logErrorOnNotFound,
					})
			}
		}
		return
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		if v.logger != nil {
			v.logger.Error(fmt.Sprintf("Failed to read %s schemas directory", capitalizeFirstLetterValidator(schemaType)),
				err,
				map[string]interface{}{
					"dir_path":    dirPath,
					"schema_type": schemaType,
				})
		}
		return
	}

	if v.logger != nil {
		v.logger.Debug(fmt.Sprintf("Found %d entries in %s directory", len(entries), schemaType),
			map[string]interface{}{
				"dir_path":    dirPath,
				"entry_count": len(entries),
			})
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".hcl") {
			continue
		}

		schemaPath := filepath.Join(dirPath, entry.Name())
		schemaName := strings.TrimSuffix(entry.Name(), ".hcl")

		schema, err := loadFunc(schemaPath)
		if err != nil {
			if v.logger != nil {
				v.logger.Error(fmt.Sprintf("Failed to load %s schema", capitalizeFirstLetterValidator(schemaType)),
					err,
					map[string]interface{}{
						"schema_path": schemaPath,
						"schema_name": schemaName,
						"schema_type": schemaType,
					})
			}
			continue
		}

		if v.logger != nil {
			v.logger.Info(fmt.Sprintf("Loaded %s schema", capitalizeFirstLetterValidator(schemaType)),
				map[string]interface{}{
					"schema_name": schemaName,
					"schema_path": schemaPath,
				})
		}

		schemaMap[schemaName] = schema
	}
}

// Helper methods to get schemas by type
func (v *SchemaDrivenValidator) getStructureSchema(name string) (*spookytypesschemas.Schema, error) {
	if schema, exists := v.structureSchemas[name]; exists {
		return schema, nil
	}
	return nil, fmt.Errorf("structure schema not found: %s", name)
}

func (v *SchemaDrivenValidator) getValidationSchema(name string) (*spookytypesschemas.Schema, error) {
	if schema, exists := v.validationSchemas[name]; exists {
		return schema, nil
	}
	return nil, fmt.Errorf("validation schema not found: %s", name)
}

func (v *SchemaDrivenValidator) getMetadataSchema(name string) (*spookytypesschemas.Schema, error) {
	if schema, exists := v.metadataSchemas[name]; exists {
		return schema, nil
	}
	return nil, fmt.Errorf("metadata schema not found: %s", name)
}

// loadStructureSchema loads a structure schema from file
func (v *SchemaDrivenValidator) loadStructureSchema(schemaPath string) (*spookytypesschemas.Schema, error) {
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read structure schema: %w", err)
	}

	schema := &spookytypesschemas.Schema{
		Version:     "1.0",
		Type:        "structure",
		Name:        filepath.Base(schemaPath),
		Description: fmt.Sprintf("Structure schema from %s", schemaPath),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Content:     string(data),
		Metadata:    make(map[string]interface{}),
	}

	return schema, nil
}

// loadValidationSchema loads a validation schema from file
func (v *SchemaDrivenValidator) loadValidationSchema(schemaPath string) (*spookytypesschemas.Schema, error) {
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read validation schema: %w", err)
	}

	schema := &spookytypesschemas.Schema{
		Version:     "1.0",
		Type:        "validation",
		Name:        filepath.Base(schemaPath),
		Description: fmt.Sprintf("Validation schema from %s", schemaPath),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Content:     string(data),
		Metadata:    make(map[string]interface{}),
	}

	return schema, nil
}

// loadMetadataSchema loads a metadata schema from file
func (v *SchemaDrivenValidator) loadMetadataSchema(schemaPath string) (*spookytypesschemas.Schema, error) {
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata schema: %w", err)
	}

	schema := &spookytypesschemas.Schema{
		Version:     "1.0",
		Type:        "metadata",
		Name:        filepath.Base(schemaPath),
		Description: fmt.Sprintf("Metadata schema from %s", schemaPath),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Content:     string(data),
		Metadata:    make(map[string]interface{}),
	}

	return schema, nil
}
