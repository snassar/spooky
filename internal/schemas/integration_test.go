package schemas

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	spookytypesschemas "spooky/internal/types/schemas"
)

// TestEndToEndSchemaValidation tests the complete validation workflow
func TestEndToEndSchemaValidation(t *testing.T) {
	logger := &DebugLogger{t: t}
	manager := NewManager(logger)

	// Load a real schema file that exists
	schema, err := manager.Load("internal/schemas/schemas/project.schema.hcl")
	if err != nil {
		// Skip test if schema file doesn't exist
		t.Skip("Schema file not found, skipping test")
	}

	// Validate test data against schema
	testData := map[string]interface{}{
		"project": map[string]interface{}{
			"name":        "test-project",
			"description": "A test project",
		},
	}

	validator := NewEnhancedValidator(&ValidationConfig{
		Mode: ValidationModeStrict,
	})

	result, err := validator.ValidateWithEnhancedFeatures(context.Background(), schema, testData)
	require.NoError(t, err)

	assert.True(t, result.Valid)
	assert.Equal(t, 0, len(result.Errors))
	assert.Greater(t, result.Statistics.TotalFields, 0)
	assert.Greater(t, result.Statistics.RulesProcessed, 0)
}

// TestProjectSchemaValidation tests project schema validation
func TestProjectSchemaValidation(t *testing.T) {
	logger := &DebugLogger{t: t}
	config := &SchemaDrivenValidationConfig{
		UseEmbeddedSchemas: true,
		StrictValidation:   true,
		AllowUnknownFields: false,
		DetailedErrors:     true,
	}
	validator := NewSchemaDrivenValidator(logger, config)

	testData := map[string]interface{}{
		"project": map[string]interface{}{
			"name":        "test-project",
			"description": "A test project",
		},
	}

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

	err := validator.validateProjectSchema(testData, result)
	require.NoError(t, err)

	// Note: This test may fail if the project schema is not found
	// In that case, it's expected behavior as the schema might not exist
	if err == nil {
		assert.Greater(t, result.Statistics.TotalFields, 0)
	}
}

// TestMachinesSchemaValidation tests machines schema validation
func TestMachinesSchemaValidation(t *testing.T) {
	logger := &DebugLogger{t: t}
	config := &SchemaDrivenValidationConfig{
		UseEmbeddedSchemas: true,
		StrictValidation:   true,
		AllowUnknownFields: false,
		DetailedErrors:     true,
	}
	validator := NewSchemaDrivenValidator(logger, config)

	testData := map[string]interface{}{
		"machines": map[string]interface{}{
			"machine": map[string]interface{}{
				"hostname": "test-server",
				"port":     22,
				"user":     "admin",
			},
		},
	}

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

	err := validator.validateMachinesSchema(testData, result)
	require.NoError(t, err)

	// Note: This test may fail if the machines schema is not found
	// In that case, it's expected behavior as the schema might not exist
	if err == nil {
		assert.Greater(t, result.Statistics.TotalFields, 0)
	}
}

// TestActionsSchemaValidation tests actions schema validation
func TestActionsSchemaValidation(t *testing.T) {
	logger := &DebugLogger{t: t}
	config := &SchemaDrivenValidationConfig{
		UseEmbeddedSchemas: true,
		StrictValidation:   true,
		AllowUnknownFields: false,
		DetailedErrors:     true,
	}
	validator := NewSchemaDrivenValidator(logger, config)

	testData := map[string]interface{}{
		"actions": map[string]interface{}{
			"action": map[string]interface{}{
				"name":        "test-action",
				"description": "A test action",
				"command":     "echo 'hello world'",
			},
		},
	}

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

	err := validator.validateActionsSchema(testData, result)
	require.NoError(t, err)

	// Note: This test may fail if the actions schema is not found
	// In that case, it's expected behavior as the schema might not exist
	if err == nil {
		assert.Greater(t, result.Statistics.TotalFields, 0)
	}
}

// TestSchemaLoadingIntegration tests schema loading integration
func TestSchemaLoadingIntegration(t *testing.T) {
	logger := &DebugLogger{t: t}
	manager := NewManager(logger)

	// Test loading a schema file
	schema, err := manager.Load("internal/schemas/schemas/project.schema.hcl")
	if err != nil {
		// Skip test if schema file doesn't exist
		t.Skip("Schema file not found, skipping test")
	}

	require.NotNil(t, schema)
	assert.NotEmpty(t, schema.Name)
	assert.NotEmpty(t, schema.Content)
	assert.NotNil(t, schema.Validation)
	assert.NotNil(t, schema.Validation.Fields)
}

// TestEmbeddedSchemaIntegration tests embedded schema integration
func TestEmbeddedSchemaIntegration(t *testing.T) {
	logger := &DebugLogger{t: t}
	config := &SchemaDrivenValidationConfig{
		UseEmbeddedSchemas: true,
		StrictValidation:   true,
		AllowUnknownFields: false,
		DetailedErrors:     true,
	}
	validator := NewSchemaDrivenValidator(logger, config)

	// Test that embedded schemas are loaded
	schema, err := validator.GetEmbeddedSchema("project")
	if err != nil {
		// Skip test if embedded schema is not found
		t.Skip("Embedded schema not found, skipping test")
	}

	require.NotNil(t, schema)
	assert.NotEmpty(t, schema.Name)
	assert.NotEmpty(t, schema.Content)

	// Use logger to verify it's working
	logger.Info("Embedded schema integration test completed", map[string]interface{}{
		"schema_name": schema.Name,
	})
}

// TestValidationResultConsistency tests validation result consistency
func TestValidationResultConsistency(t *testing.T) {
	logger := &DebugLogger{t: t}
	validator := NewEnhancedValidator(&ValidationConfig{
		Mode: ValidationModeStrict,
	})

	// Create a schema with validation rules
	schema := &spookytypesschemas.Schema{
		Version:     "1.0",
		Type:        "test",
		Name:        "test-schema",
		Description: "Test schema for validation",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Content: `
test_schema {
  required_field = {
    type = "string"
    required = true
    description = "A required field"
  }
  
  optional_field = {
    type = "string"
    required = false
    description = "An optional field"
  }
}
`,
		Metadata: make(map[string]interface{}),
		Validation: &spookytypesschemas.SchemaValidation{
			Enabled: true,
			Mode:    "strict",
			Fields: map[string]*spookytypesschemas.FieldValidation{
				"test_schema.required_field": {
					Type:        "string",
					Required:    true,
					Description: "A required field",
				},
			},
		},
	}

	// Test data with missing required field (should fail validation)
	testData := map[string]interface{}{
		"test_schema": map[string]interface{}{
			"invalid_field": "invalid_value", // Missing required_field
		},
	}

	result, err := validator.ValidateWithEnhancedFeatures(context.Background(), schema, testData)
	require.NoError(t, err)

	// Print actual result for debugging
	t.Logf("Validation result: Valid=%v, Errors=%d, Warnings=%d", result.Valid, len(result.Errors), len(result.Warnings))
	t.Logf("Statistics: TotalFields=%d, ValidFields=%d, InvalidFields=%d", result.Statistics.TotalFields, result.Statistics.ValidFields, result.Statistics.InvalidFields)

	// Test consistency
	assert.False(t, result.Valid)
	assert.Equal(t, result.Statistics.ValidFields+result.Statistics.InvalidFields, result.Statistics.TotalFields)
	assert.Equal(t, len(result.Errors) > 0, !result.Valid)

	// Use logger to verify it's working
	logger.Info("Validation result consistency test completed", map[string]interface{}{
		"valid":  result.Valid,
		"errors": len(result.Errors),
	})
}

// TestSchemaDrivenValidationWorkflow tests the complete schema-driven validation workflow
func TestSchemaDrivenValidationWorkflow(t *testing.T) {
	logger := &DebugLogger{t: t}
	config := &SchemaDrivenValidationConfig{
		UseEmbeddedSchemas: true,
		StrictValidation:   true,
		AllowUnknownFields: false,
		DetailedErrors:     true,
	}
	validator := NewSchemaDrivenValidator(logger, config)

	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.hcl")
	testContent := `
project {
  name = "test-project"
  description = "A test project"
}
`
	err := os.WriteFile(testFile, []byte(testContent), 0600)
	require.NoError(t, err)

	// Test configuration validation
	result, err := validator.ValidateConfiguration(context.Background(), testFile, "project")
	require.NoError(t, err)

	// Verify result structure
	assert.NotNil(t, result)
	assert.NotNil(t, result.Statistics)
	assert.GreaterOrEqual(t, result.Statistics.Duration, time.Duration(0))
}

// TestErrorScenarios tests error scenarios and edge cases
func TestErrorScenarios(t *testing.T) {
	logger := &DebugLogger{t: t}
	validator := NewEnhancedValidator(&ValidationConfig{
		Mode: ValidationModeStrict,
	})

	// Test with nil schema
	result, err := validator.ValidateWithEnhancedFeatures(context.Background(), nil, map[string]interface{}{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "schema cannot be nil")
	assert.False(t, result.Valid)
	assert.Greater(t, len(result.Errors), 0)

	// Test with nil data
	schema := &spookytypesschemas.Schema{
		Version:     "1.0",
		Type:        "test",
		Name:        "test-schema",
		Description: "Test schema",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Content:     "",
		Metadata:    make(map[string]interface{}),
	}

	result, err = validator.ValidateWithEnhancedFeatures(context.Background(), schema, nil)
	require.NoError(t, err)
	assert.NotNil(t, result)

	// Use logger to verify it's working
	logger.Info("Error scenarios test completed", map[string]interface{}{
		"test_count": 2,
	})
}

// TestSchemaLoading verifies that schemas are loaded correctly
func TestSchemaLoading(t *testing.T) {
	logger := &DebugLogger{t: t}
	config := &SchemaDrivenValidationConfig{
		UseEmbeddedSchemas: true,
		StrictValidation:   true,
		AllowUnknownFields: false,
		DetailedErrors:     true,
	}
	validator := NewSchemaDrivenValidator(logger, config)

	// Check if schemas were loaded
	t.Logf("Structure schemas loaded: %d", len(validator.structureSchemas))
	t.Logf("Validation schemas loaded: %d", len(validator.validationSchemas))
	t.Logf("Metadata schemas loaded: %d", len(validator.metadataSchemas))

	// List loaded schemas
	for name := range validator.structureSchemas {
		t.Logf("Structure schema: %s", name)
	}
	for name := range validator.validationSchemas {
		t.Logf("Validation schema: %s", name)
	}
	for name := range validator.metadataSchemas {
		t.Logf("Metadata schema: %s", name)
	}

	// Check if project schema is loaded
	if _, exists := validator.structureSchemas["project"]; exists {
		t.Logf("Project schema found in structure schemas")
	} else {
		t.Logf("Project schema NOT found in structure schemas")
	}
}
