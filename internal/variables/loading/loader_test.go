package loading

import (
	"context"
	"os"
	"testing"

	"spooky/internal/schemas"
)

func TestHCLVariableLoader_ParseHCLToVariables(t *testing.T) {
	validator := schemas.NewSchemaValidator()
	loader := NewHCLVariableLoader(validator)

	// Test HCL content
	hclContent := []byte(`
variables {
  variable "test_string" {
    type        = "string"
    description = "A test string variable"
    default     = "default_value"
  }

  variable "test_number" {
    type        = "number"
    description = "A test number variable"
    default     = 42
    required    = true
  }

  variable "test_complex" {
    type        = "object"
    description = "A complex variable with validation"
    
    validation {
      condition     = "length(var.test_complex.key1) > 0"
      error_message = "key1 cannot be empty"
    }
    
    constraints {
      min_length = 1
      max_length = 100
    }
  }
}
`)

	variables, err := loader.parseHCLToVariables(hclContent)
	if err != nil {
		t.Fatalf("Failed to parse HCL: %v", err)
	}

	if len(variables) != 3 {
		t.Fatalf("Expected 3 variables, got %d", len(variables))
	}

	// Check first variable
	if variables[0].Name != "test_string" {
		t.Errorf("Expected name 'test_string', got '%s'", variables[0].Name)
	}
	if variables[0].Type != "string" {
		t.Errorf("Expected type 'string', got '%s'", variables[0].Type)
	}
	if variables[0].Description != "A test string variable" {
		t.Errorf("Expected description 'A test string variable', got '%s'", variables[0].Description)
	}

	// Check second variable
	if variables[1].Name != "test_number" {
		t.Errorf("Expected name 'test_number', got '%s'", variables[1].Name)
	}
	if variables[1].Type != "number" {
		t.Errorf("Expected type 'number', got '%s'", variables[1].Type)
	}
	if !variables[1].Required {
		t.Error("Expected required to be true")
	}

	// Check third variable
	if variables[2].Name != "test_complex" {
		t.Errorf("Expected name 'test_complex', got '%s'", variables[2].Name)
	}
	if variables[2].Type != "object" {
		t.Errorf("Expected type 'object', got '%s'", variables[2].Type)
	}
	if variables[2].Validation == nil {
		t.Error("Expected validation to be set")
	} else {
		if variables[2].Validation.Condition != "length(var.test_complex.key1) > 0" {
			t.Errorf("Expected condition 'length(var.test_complex.key1) > 0', got '%s'", variables[2].Validation.Condition)
		}
	}
	if variables[2].Constraints == nil {
		t.Error("Expected constraints to be set")
	}
}

func TestJSONVariableLoader_ParseJSONToVariables(t *testing.T) {
	validator := schemas.NewSchemaValidator()
	loader := NewJSONVariableLoader(validator)

	// Test JSON content in export format
	jsonData := map[string]interface{}{
		"metadata": map[string]interface{}{
			"version": "1.0",
			"format":  "json",
		},
		"variables": map[string]interface{}{
			"test_string": map[string]interface{}{
				"name":        "test_string",
				"type":        "string",
				"value":       "test_value",
				"description": "A test string variable",
			},
			"test_number": map[string]interface{}{
				"name":     "test_number",
				"type":     "number",
				"value":    float64(42),
				"required": true,
			},
		},
	}

	variables, err := loader.parseJSONToVariables(jsonData)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if len(variables) != 2 {
		t.Fatalf("Expected 2 variables, got %d", len(variables))
	}

	// Check first variable
	if variables[0].Name != "test_string" {
		t.Errorf("Expected name 'test_string', got '%s'", variables[0].Name)
	}
	if variables[0].Type != "string" {
		t.Errorf("Expected type 'string', got '%s'", variables[0].Type)
	}
	if variables[0].Description != "A test string variable" {
		t.Errorf("Expected description 'A test string variable', got '%s'", variables[0].Description)
	}

	// Check second variable
	if variables[1].Name != "test_number" {
		t.Errorf("Expected name 'test_number', got '%s'", variables[1].Name)
	}
	if variables[1].Type != "number" {
		t.Errorf("Expected type 'number', got '%s'", variables[1].Type)
	}
	if !variables[1].Required {
		t.Error("Expected required to be true")
	}
}

func TestHCLVariableLoader_Load(t *testing.T) {
	validator := schemas.NewSchemaValidator()
	loader := NewHCLVariableLoader(validator)

	// Create a temporary file with HCL content
	tmpFile := createTempHCLFile(t, `
variables {
  variable "test_var" {
    type        = "string"
    description = "Test variable"
    default     = "test_value"
  }
}`)
	defer tmpFile.Close()

	variables, err := loader.Load(context.Background(), tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load variables: %v", err)
	}

	if len(variables) != 1 {
		t.Fatalf("Expected 1 variable, got %d", len(variables))
	}

	if variables[0].Name != "test_var" {
		t.Errorf("Expected name 'test_var', got '%s'", variables[0].Name)
	}
}

func TestJSONVariableLoader_Load(t *testing.T) {
	validator := schemas.NewSchemaValidator()
	loader := NewJSONVariableLoader(validator)

	// Create a temporary file with JSON content
	tmpFile := createTempJSONFile(t, `{
  "metadata": {
    "version": "1.0",
    "format": "json"
  },
  "variables": {
    "test_var": {
      "name": "test_var",
      "type": "string",
      "value": "test_value",
      "description": "Test variable"
    }
  }
}`)
	defer tmpFile.Close()

	variables, err := loader.Load(context.Background(), tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load variables: %v", err)
	}

	if len(variables) != 1 {
		t.Fatalf("Expected 1 variable, got %d", len(variables))
	}

	if variables[0].Name != "test_var" {
		t.Errorf("Expected name 'test_var', got '%s'", variables[0].Name)
	}
}

// Helper functions for creating temporary files
func createTempHCLFile(t *testing.T, content string) *os.File {
	tmpFile, err := os.CreateTemp("", "test-*.hcl")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Reopen for reading
	tmpFile, err = os.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to reopen temp file: %v", err)
	}

	return tmpFile
}

func createTempJSONFile(t *testing.T, content string) *os.File {
	tmpFile, err := os.CreateTemp("", "test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Reopen for reading
	tmpFile, err = os.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to reopen temp file: %v", err)
	}

	return tmpFile
}
