package templates

import (
	"context"
	"fmt"
	"testing"
	"time"

	spookylogging "spooky/internal/logging"
	spookyschemas "spooky/internal/schemas"
	spookytypes "spooky/internal/types"
	spookytypesschemas "spooky/internal/types/schemas"
	spookytypestemplates "spooky/internal/types/templates"

	"github.com/stretchr/testify/assert"
)

func TestValidateTemplateWithSchema(t *testing.T) {
	// Create a logger
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("test")

	// Create a template validator with schema validation
	validator := &TemplateValidator{
		logger: logger,
		schemaValidator: &MockSchemaValidator{
			shouldValidate: true,
		},
		schemaManager: nil,
	}

	// Create a test template
	template := &spookytypes.Template{
		ID:            "test-template",
		SourcePath:    "templates/test.tmpl",
		Type:          "template",
		Scope:         "project",
		SecurityLevel: "standard",
		Engine:        "go-template",
		Content:       "Hello {{.Name}}!",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Test validation
	result, err := validator.ValidateTemplateComprehensive(context.Background(), template)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if !result.Valid {
		t.Errorf("Expected template to be valid, but got errors: %v", result.Errors)
	}
}

func TestValidateTemplateWithoutSchemaValidator(t *testing.T) {
	// Create a logger
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("test")

	// Create a template validator without schema validation
	validator := &TemplateValidator{
		logger: logger,
		// No schema validator set
	}

	// Create a test template
	template := &spookytypes.Template{
		ID:            "test-template",
		SourcePath:    "templates/test.tmpl",
		Type:          "template",
		Scope:         "project",
		SecurityLevel: "standard",
		Engine:        "go-template",
		Content:       "Hello {{.Name}}!",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Test validation
	result, err := validator.ValidateTemplateComprehensive(context.Background(), template)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	// Should return error about missing schema validator
	if result.Valid {
		t.Error("Expected template to be invalid due to missing schema validator")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected validation errors due to missing schema validator")
	}

	// Check for the specific error message
	found := false
	for _, err := range result.Errors {
		if err.Message == "schema validator not configured" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected error message about missing schema validator")
	}
}

func TestEnhancedTemplateMetadata(t *testing.T) {
	// Create enhanced metadata
	metadata := &spookytypestemplates.TemplateMetadata{
		Name:         "test-template",
		Description:  "A test template for validation",
		Author:       "test-author",
		Version:      "1.0.0",
		Tags:         []string{"test", "template", "validation"},
		Category:     "testing",
		Subcategory:  "unit-tests",
		Priority:     1,
		Keywords:     []string{"test", "validation", "template"},
		Dependencies: []string{"base-template"},
		Compatibility: map[string]string{
			"spooky": ">=1.0.0",
		},
		UsageCount: 5,
		LastUsed:   time.Now(),
	}

	// Verify all fields are set correctly
	assert.Equal(t, "test-template", metadata.Name)
	assert.Equal(t, "A test template for validation", metadata.Description)
	assert.Equal(t, "test-author", metadata.Author)
	assert.Equal(t, "1.0.0", metadata.Version)
	assert.Equal(t, []string{"test", "template", "validation"}, metadata.Tags)
	assert.Equal(t, "testing", metadata.Category)
	assert.Equal(t, "unit-tests", metadata.Subcategory)
	assert.Equal(t, 1, metadata.Priority)
	assert.Equal(t, []string{"test", "validation", "template"}, metadata.Keywords)
	assert.Equal(t, []string{"base-template"}, metadata.Dependencies)
	assert.Equal(t, map[string]string{"spooky": ">=1.0.0"}, metadata.Compatibility)
	assert.Equal(t, 5, metadata.UsageCount)
	assert.False(t, metadata.LastUsed.IsZero())
}

func TestEnhancedMetadataIndexer(t *testing.T) {
	// Create enhanced indexer
	indexer := NewEnhancedMetadataIndexer()

	// Create test metadata
	metadata1 := &spookytypestemplates.TemplateMetadata{
		Name:        "web-deployment",
		Description: "Web application deployment template",
		Tags:        []string{"web", "deployment", "kubernetes"},
		Category:    "deployment",
		Keywords:    []string{"web", "app", "deploy"},
	}

	metadata2 := &spookytypestemplates.TemplateMetadata{
		Name:        "database-config",
		Description: "Database configuration template",
		Tags:        []string{"database", "config", "postgres"},
		Category:    "configuration",
		Keywords:    []string{"database", "postgres", "config"},
	}

	// Index metadata
	err := indexer.IndexMetadata(metadata1)
	assert.NoError(t, err)

	err = indexer.IndexMetadata(metadata2)
	assert.NoError(t, err)

	// Test search by name
	results, err := indexer.Search("web", &SearchFilters{})
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "web-deployment", results[0].Metadata.Name)
	assert.Greater(t, results[0].Score, 0.0)

	// Test search by tags
	results, err = indexer.Search("kubernetes", &SearchFilters{})
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "web-deployment", results[0].Metadata.Name)

	// Test search by category
	results, err = indexer.Search("deployment", &SearchFilters{Category: "deployment"})
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "web-deployment", results[0].Metadata.Name)

	// Test search by keywords
	results, err = indexer.Search("postgres", &SearchFilters{})
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "database-config", results[0].Metadata.Name)

	// Test full-text search
	results, err = indexer.Search("application", &SearchFilters{})
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "web-deployment", results[0].Metadata.Name)
}

// MockSchemaValidator implements SchemaValidator for testing
type MockSchemaValidator struct {
	shouldValidate bool
}

// MockSchemaManager embeds the real schema manager for testing
type MockSchemaManager struct {
	*spookyschemas.Manager
	shouldLoad bool
}

func (m *MockSchemaValidator) Validate(_ *spookytypesschemas.Schema, _ interface{}) (*spookytypesschemas.ValidationResult, error) {
	if m.shouldValidate {
		return &spookytypesschemas.ValidationResult{
			Valid:       true,
			ValidatedAt: time.Now(),
			Errors:      []spookytypesschemas.SchemaError{},
			Warnings:    []spookytypesschemas.SchemaError{},
			Info:        []spookytypesschemas.SchemaError{},
			Details:     make(map[string]interface{}),
		}, nil
	}

	return &spookytypesschemas.ValidationResult{
		Valid:       false,
		ValidatedAt: time.Now(),
		Errors: []spookytypesschemas.SchemaError{
			{
				Message: "mock validation failed",
			},
		},
		Warnings: []spookytypesschemas.SchemaError{},
		Info:     []spookytypesschemas.SchemaError{},
		Details:  make(map[string]interface{}),
	}, nil
}

func (m *MockSchemaValidator) ValidateFile(_ *spookytypesschemas.Schema, _ string) (*spookytypesschemas.ValidationResult, error) {
	return m.Validate(nil, nil)
}

func (m *MockSchemaValidator) ValidateString(_ *spookytypesschemas.Schema, _ string) (*spookytypesschemas.ValidationResult, error) {
	return m.Validate(nil, nil)
}

func (m *MockSchemaValidator) ValidateBytes(_ *spookytypesschemas.Schema, _ []byte) (*spookytypesschemas.ValidationResult, error) {
	return m.Validate(nil, nil)
}

func (m *MockSchemaValidator) ValidateWithContext(_ *spookytypesschemas.Schema, _ interface{}, _ map[string]interface{}) (*spookytypesschemas.ValidationResult, error) {
	return m.Validate(nil, nil)
}

func (m *MockSchemaValidator) ValidateField(_ *spookytypesschemas.Schema, _ string, _ interface{}) (*spookytypesschemas.ValidationResult, error) {
	return m.Validate(nil, nil)
}

func (m *MockSchemaManager) Load(_ string) (*spookytypesschemas.Schema, error) {
	if m.shouldLoad {
		return &spookytypesschemas.Schema{
			Name:    "mock-schema",
			Version: "0.20250817.0",
		}, nil
	}
	return nil, fmt.Errorf("mock schema load failed")
}
