package variables

import (
	"context"
	"testing"

	"spooky/internal/logging"
	loggingtypes "spooky/internal/logging/types"
	"spooky/internal/variables/importexport"
	"spooky/internal/variables/loading"
	"spooky/internal/variables/resolution"
	"spooky/internal/variables/types"
	"spooky/internal/variables/validation"
)

func TestNewVariableManager(t *testing.T) {
	logger := logging.NewLogger(loggingtypes.Config{})

	// Test the factory function
	manager := NewVariableManager(logger)
	if manager == nil {
		t.Fatal("NewVariableManager returned nil")
	}

	// Test that it implements the interface
	var _ VariableManager = manager
}

func TestManagerBasicOperations(t *testing.T) {
	logger := logging.NewLogger(loggingtypes.Config{})
	manager := NewVariableManager(logger)

	ctx := context.Background()

	// Test creating a variable
	variable := &types.Variable{
		Name:  "test_var",
		Value: "test_value",
		Type:  "string",
	}

	// Test setting a variable
	err := manager.SetVariable(ctx, variable)
	if err != nil {
		t.Fatalf("Failed to set variable: %v", err)
	}

	// Test getting a variable
	retrieved, err := manager.GetVariable(ctx, "test_var")
	if err != nil {
		t.Fatalf("Failed to get variable: %v", err)
	}

	if retrieved.Name != "test_var" {
		t.Errorf("Expected variable name 'test_var', got '%s'", retrieved.Name)
	}

	if retrieved.Value != "test_value" {
		t.Errorf("Expected variable value 'test_value', got '%v'", retrieved.Value)
	}

	// Test listing variables
	variables, err := manager.ListVariables(ctx)
	if err != nil {
		t.Fatalf("Failed to list variables: %v", err)
	}

	if len(variables) != 1 {
		t.Errorf("Expected 1 variable, got %d", len(variables))
	}

	// Test deleting a variable
	err = manager.DeleteVariable(ctx, "test_var")
	if err != nil {
		t.Fatalf("Failed to delete variable: %v", err)
	}

	// Verify it's deleted
	_, err = manager.GetVariable(ctx, "test_var")
	if err == nil {
		t.Error("Expected error when getting deleted variable")
	}
}

func TestManagerWithDependencies(t *testing.T) {
	logger := logging.NewLogger(loggingtypes.Config{})

	// Create default configuration
	config := &types.Config{
		LoadingConfig: &types.LoadingConfig{
			DefaultEncoding:   "utf-8",
			MaxFileSize:       1024 * 1024,
			AllowedExtensions: []string{".hcl", ".json"},
		},
		ResolutionConfig: &types.ResolutionConfig{
			MaxRecursionDepth: 10,
			DefaultValues:     make(map[string]interface{}),
			StrictMode:        false,
		},
		ValidationConfig: &types.ValidationConfig{
			ValidationRules:     &types.ValidationRules{},
			StrictValidation:    false,
			MaxValidationErrors: 100,
		},
		ImportExportConfig: &types.ImportExportConfig{
			ExportOptions: &types.ExportOptions{
				IncludeMetadata: true,
				PrettyPrint:     true,
			},
			ImportOptions: &types.ImportOptions{
				MergePolicy: "overwrite",
				Overwrite:   false,
			},
			DefaultFormat: types.ExportFormatHCL,
		},
	}

	// Create concrete implementations
	loadingManager := loading.NewManager(config.LoadingConfig, logger)
	resolutionManager := resolution.NewManager(config.ResolutionConfig, logger)
	validationManager := validation.NewManager(config.ValidationConfig, logger)
	importExportManager := importexport.NewManager(config.ImportExportConfig, logger)

	// Create manager with dependencies
	manager := NewManager(
		config,
		loadingManager,
		resolutionManager,
		validationManager,
		importExportManager,
		logger,
	)

	if manager == nil {
		t.Fatal("NewManager returned nil")
	}

	// Test that it implements the interface
	var _ VariableManager = manager
}
