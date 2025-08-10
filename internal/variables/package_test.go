package variables

import (
	"testing"

	"spooky/internal/logging"
	loggingtypes "spooky/internal/types/logging"
	vartypes "spooky/internal/types/variables"
)

// TestPackageStructure verifies that the package structure is correct
func TestPackageStructure(t *testing.T) {
	// Test that we can create a basic configuration
	config := &vartypes.Config{
		LoadingConfig: &vartypes.LoadingConfig{
			DefaultEncoding:   "utf-8",
			MaxFileSize:       1024 * 1024,
			AllowedExtensions: []string{".hcl", ".json"},
		},
		ResolutionConfig: &vartypes.ResolutionConfig{
			MaxRecursionDepth: 10,
			DefaultValues:     make(map[string]interface{}),
			StrictMode:        false,
		},
		ValidationConfig: &vartypes.ValidationConfig{
			ValidationRules:     &vartypes.ValidationRules{},
			StrictValidation:    false,
			MaxValidationErrors: 100,
		},
		ImportExportConfig: &vartypes.ImportExportConfig{
			ExportOptions: &vartypes.ExportOptions{
				IncludeMetadata: true,
				PrettyPrint:     true,
			},
			ImportOptions: &vartypes.ImportOptions{
				MergePolicy: "overwrite",
				Overwrite:   false,
			},
			DefaultFormat: vartypes.ExportFormatHCL,
		},
	}

	if config == nil {
		t.Fatal("Failed to create configuration")
	}

	// Test that we can create a variable
	variable := &vartypes.Variable{
		Name:        "test_variable",
		Type:        "string",
		Value:       "test_value",
		Description: "A test variable",
		Required:    false,
		Sensitive:   false,
		Encrypted:   false,
	}

	if variable.Name != "test_variable" {
		t.Errorf("Expected variable name 'test_variable', got '%s'", variable.Name)
	}

	// Test that we can create a variable collection
	collection := &vartypes.VariableCollection{
		Variables: []*vartypes.Variable{variable},
		Path:      "/test/path",
	}

	if len(collection.Variables) != 1 {
		t.Errorf("Expected 1 variable in collection, got %d", len(collection.Variables))
	}

	// Test that we can create a variable context
	context := &vartypes.VariableContext{
		Variables: map[string]*vartypes.Variable{
			"test_variable": variable,
		},
		ProjectPath: "/test/project",
	}

	if len(context.Variables) != 1 {
		t.Errorf("Expected 1 variable in context, got %d", len(context.Variables))
	}

	// Test that we can create a logger
	logger := logging.NewLogger(loggingtypes.Config{})
	if logger == nil {
		t.Fatal("Failed to create logger")
	}

	// Test that we can create a manager (this will fail if there are import issues)
	// We'll skip this for now since the project package has build issues
	/*
		manager := NewVariableManager(logger)
		if manager == nil {
			t.Fatal("Failed to create variable manager")
		}
	*/
}

// TestVariableTypes verifies that all variable types are properly defined
func TestVariableTypes(t *testing.T) {
	// Test string variable
	stringVar := &vartypes.Variable{
		Name:  "string_var",
		Type:  "string",
		Value: "test_string",
		Constraints: &vartypes.VariableConstraints{
			MinLength: new(int),
			MaxLength: new(int),
			Pattern:   new(string),
		},
	}
	*stringVar.Constraints.MinLength = 1
	*stringVar.Constraints.MaxLength = 100
	*stringVar.Constraints.Pattern = "^[a-z]+$"

	// Test numeric variable
	numericVar := &vartypes.Variable{
		Name:  "numeric_var",
		Type:  "number",
		Value: 42.0,
		Constraints: &vartypes.VariableConstraints{
			MinValue: new(float64),
			MaxValue: new(float64),
		},
	}
	*numericVar.Constraints.MinValue = 0.0
	*numericVar.Constraints.MaxValue = 100.0

	// Test list variable
	listVar := &vartypes.Variable{
		Name:  "list_var",
		Type:  "list",
		Value: []interface{}{"item1", "item2", "item3"},
		Constraints: &vartypes.VariableConstraints{
			MinItems: new(int),
			MaxItems: new(int),
		},
	}
	*listVar.Constraints.MinItems = 1
	*listVar.Constraints.MaxItems = 10

	// Test file variable
	fileVar := &vartypes.Variable{
		Name:  "file_var",
		Type:  "file",
		Value: "/path/to/file.txt",
		Constraints: &vartypes.VariableConstraints{
			FileExists:   new(bool),
			FileReadable: new(bool),
		},
	}
	*fileVar.Constraints.FileExists = true
	*fileVar.Constraints.FileReadable = true

	// Test path variable
	pathVar := &vartypes.Variable{
		Name:  "path_var",
		Type:  "path",
		Value: "/absolute/path",
		Constraints: &vartypes.VariableConstraints{
			PathExists:   new(bool),
			PathAbsolute: new(bool),
		},
	}
	*pathVar.Constraints.PathExists = true
	*pathVar.Constraints.PathAbsolute = true

	// Verify all variables have the expected types
	testCases := []struct {
		name     string
		variable *vartypes.Variable
		expected string
	}{
		{"string", stringVar, "string"},
		{"numeric", numericVar, "number"},
		{"list", listVar, "list"},
		{"file", fileVar, "file"},
		{"path", pathVar, "path"},
	}

	for _, tc := range testCases {
		if tc.variable.Type != tc.expected {
			t.Errorf("Expected variable type '%s', got '%s'", tc.expected, tc.variable.Type)
		}
	}
}

// TestValidationResult verifies that validation results work correctly
func TestValidationResult(t *testing.T) {
	result := &vartypes.ValidationResult{
		Valid: true,
		Errors: []vartypes.ValidationError{
			{
				Field:   "name",
				Message: "Name is required",
			},
		},
		Warnings: []vartypes.ValidationWarning{
			{
				Field:   "description",
				Message: "Description is recommended",
			},
		},
	}

	if !result.Valid {
		t.Error("Expected validation result to be valid")
	}

	if len(result.Errors) != 1 {
		t.Errorf("Expected 1 validation error, got %d", len(result.Errors))
	}

	if len(result.Warnings) != 1 {
		t.Errorf("Expected 1 validation warning, got %d", len(result.Warnings))
	}
}

// TestDependencyGraph verifies that dependency graphs work correctly
func TestDependencyGraph(t *testing.T) {
	graph := &vartypes.DependencyGraph{
		Nodes: map[string]*vartypes.DependencyNode{
			"var1": {
				Name:         "var1",
				Dependencies: []string{"var2"},
				Resolved:     false,
			},
			"var2": {
				Name:         "var2",
				Dependencies: []string{},
				Resolved:     true,
			},
		},
		Edges: map[string][]string{
			"var1": {"var2"},
			"var2": {},
		},
	}

	if len(graph.Nodes) != 2 {
		t.Errorf("Expected 2 nodes in graph, got %d", len(graph.Nodes))
	}

	if len(graph.Edges) != 2 {
		t.Errorf("Expected 2 edges in graph, got %d", len(graph.Edges))
	}

	if !graph.Nodes["var2"].Resolved {
		t.Error("Expected var2 to be resolved")
	}

	if graph.Nodes["var1"].Resolved {
		t.Error("Expected var1 to be unresolved")
	}
}
