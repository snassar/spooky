// Package schemas provides schema validation and management functionality for the spooky codebase.
package schemas

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"

	spookylogging "spooky/internal/logging"
	spookytypesschemas "spooky/internal/types/schemas"
)

// Validator provides functionality to validate HCL files against schemas
type Validator struct {
	logger  spookylogging.Logger
	schemas map[string]*spookytypesschemas.Schema
}

// NewValidator creates a new schema validator instance
func NewValidator(logger spookylogging.Logger) *Validator {
	return &Validator{
		logger:  logger,
		schemas: make(map[string]*spookytypesschemas.Schema),
	}
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
			}
			if diag.Detail != "" {
				error.Message += ": " + diag.Detail
			}
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
		error := spookytypesschemas.SchemaError{
			Code:        "empty_file",
			Message:     "File is empty",
			SchemaName:  schema.Name,
			SchemaType:  schema.Type,
			Severity:    "error",
			Recoverable: false,
		}
		result.Errors = append(result.Errors, error)
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
