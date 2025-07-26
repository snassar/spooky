package facts

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationError(t *testing.T) {
	// Test ValidationError struct and Error method
	validationError := &ValidationError{
		Field:   "test_field",
		Message: "test message",
		Value:   "test_value",
	}

	errorString := validationError.Error()
	assert.Contains(t, errorString, "validation error in test_field")
	assert.Contains(t, errorString, "test message")
}

func TestValidationResult(t *testing.T) {
	// Test ValidationResult struct
	result := &ValidationResult{
		Valid: true,
		Errors: []*ValidationError{
			{
				Field:   "field1",
				Message: "error1",
			},
		},
		Warnings: []*ValidationError{
			{
				Field:   "field2",
				Message: "warning1",
			},
		},
	}

	assert.True(t, result.Valid)
	assert.Len(t, result.Errors, 1)
	assert.Len(t, result.Warnings, 1)
	assert.Equal(t, "field1", result.Errors[0].Field)
	assert.Equal(t, "error1", result.Errors[0].Message)
	assert.Equal(t, "field2", result.Warnings[0].Field)
	assert.Equal(t, "warning1", result.Warnings[0].Message)
}

func TestValidateCustomFacts(t *testing.T) {
	// Test valid custom facts
	validCustomFacts := map[string]*CustomFacts{
		"test-server": {
			Custom: map[string]interface{}{
				"os": map[string]interface{}{
					"name":    "linux",
					"version": "20.04",
				},
				"hardware": map[string]interface{}{
					"cpu":    "Intel i7",
					"memory": "16GB",
				},
			},
			Overrides: map[string]interface{}{
				"network": map[string]interface{}{
					"hostname": "custom-hostname",
				},
			},
			Source: "test-source",
		},
	}

	result := ValidateCustomFacts(validCustomFacts)
	assert.True(t, result.Valid)
	assert.Len(t, result.Errors, 0)

	// Test custom facts with invalid server ID
	invalidCustomFacts := map[string]*CustomFacts{
		"": { // Empty server ID
			Custom: map[string]interface{}{
				"os": map[string]interface{}{
					"name": "linux",
				},
			},
		},
	}

	result = ValidateCustomFacts(invalidCustomFacts)
	assert.False(t, result.Valid)
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "server_id", result.Errors[0].Field)
	assert.Contains(t, result.Errors[0].Message, "server ID cannot be empty")

	// Test custom facts with invalid server ID characters
	invalidCustomFacts = map[string]*CustomFacts{
		"invalid@server": { // Invalid characters
			Custom: map[string]interface{}{
				"os": map[string]interface{}{
					"name": "linux",
				},
			},
		},
	}

	result = ValidateCustomFacts(invalidCustomFacts)
	assert.False(t, result.Valid)
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "server_id", result.Errors[0].Field)
	assert.Contains(t, result.Errors[0].Message, "server ID contains invalid characters")

	// Test custom facts with empty category
	invalidCustomFacts = map[string]*CustomFacts{
		"test-server": {
			Custom: map[string]interface{}{
				"": map[string]interface{}{ // Empty category
					"name": "linux",
				},
			},
		},
	}

	result = ValidateCustomFacts(invalidCustomFacts)
	assert.False(t, result.Valid)
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "custom", result.Errors[0].Field)
	assert.Contains(t, result.Errors[0].Message, "category name cannot be empty")

	// Test custom facts with empty fact key
	invalidCustomFacts = map[string]*CustomFacts{
		"test-server": {
			Custom: map[string]interface{}{
				"os": map[string]interface{}{
					"": "linux", // Empty fact key
				},
			},
		},
	}

	result = ValidateCustomFacts(invalidCustomFacts)
	assert.False(t, result.Valid)
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "custom", result.Errors[0].Field)
	assert.Contains(t, result.Errors[0].Message, "fact key cannot be empty")

	// Test custom facts with nil fact value
	invalidCustomFacts = map[string]*CustomFacts{
		"test-server": {
			Custom: map[string]interface{}{
				"os": map[string]interface{}{
					"name": nil, // Nil fact value
				},
			},
		},
	}

	result = ValidateCustomFacts(invalidCustomFacts)
	assert.False(t, result.Valid)
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "custom", result.Errors[0].Field)
	assert.Contains(t, result.Errors[0].Message, "fact value cannot be nil")

	// Test custom facts with invalid category type
	invalidCustomFacts = map[string]*CustomFacts{
		"test-server": {
			Custom: map[string]interface{}{
				"os": "not-an-object", // Not an object
			},
		},
	}

	result = ValidateCustomFacts(invalidCustomFacts)
	assert.False(t, result.Valid)
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "custom", result.Errors[0].Field)
	assert.Contains(t, result.Errors[0].Message, "category os must be an object")
}

func TestValidateCustomFactsWithOverrides(t *testing.T) {
	// Test custom facts with valid overrides
	validCustomFacts := map[string]*CustomFacts{
		"test-server": {
			Custom: map[string]interface{}{
				"os": map[string]interface{}{
					"name": "linux",
				},
			},
			Overrides: map[string]interface{}{
				"network": map[string]interface{}{
					"hostname": "custom-hostname",
				},
			},
		},
	}

	result := ValidateCustomFacts(validCustomFacts)
	assert.True(t, result.Valid)
	assert.Len(t, result.Errors, 0)

	// Test custom facts with invalid override category
	invalidCustomFacts := map[string]*CustomFacts{
		"test-server": {
			Custom: map[string]interface{}{
				"os": map[string]interface{}{
					"name": "linux",
				},
			},
			Overrides: map[string]interface{}{
				"": map[string]interface{}{ // Empty override category
					"hostname": "custom-hostname",
				},
			},
		},
	}

	result = ValidateCustomFacts(invalidCustomFacts)
	assert.False(t, result.Valid)
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "overrides", result.Errors[0].Field)
	assert.Contains(t, result.Errors[0].Message, "override category cannot be empty")

	// Test custom facts with invalid override key
	invalidCustomFacts = map[string]*CustomFacts{
		"test-server": {
			Custom: map[string]interface{}{
				"os": map[string]interface{}{
					"name": "linux",
				},
			},
			Overrides: map[string]interface{}{
				"network": map[string]interface{}{
					"": "custom-hostname", // Empty override key
				},
			},
		},
	}

	result = ValidateCustomFacts(invalidCustomFacts)
	assert.False(t, result.Valid)
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "overrides", result.Errors[0].Field)
	assert.Contains(t, result.Errors[0].Message, "override key cannot be empty")

	// Test custom facts with nil override value
	invalidCustomFacts = map[string]*CustomFacts{
		"test-server": {
			Custom: map[string]interface{}{
				"os": map[string]interface{}{
					"name": "linux",
				},
			},
			Overrides: map[string]interface{}{
				"network": map[string]interface{}{
					"hostname": nil, // Nil override value
				},
			},
		},
	}

	result = ValidateCustomFacts(invalidCustomFacts)
	assert.False(t, result.Valid)
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "overrides", result.Errors[0].Field)
	assert.Contains(t, result.Errors[0].Message, "override value cannot be nil")
}

func TestValidateWithTestDataFromExamples(t *testing.T) {
	// Test validation with data patterns from examples/testing
	examplesDir := "../../examples/testing"

	// Test with valid project data patterns
	validProjectPath := examplesDir + "/test-valid-project"
	if _, err := os.Stat(validProjectPath); os.IsNotExist(err) {
		t.Skip("Test data directory not found, skipping test")
	}

	// Test custom facts similar to test data
	customFacts := map[string]*CustomFacts{
		"example-server": {
			Custom: map[string]interface{}{
				"os": map[string]interface{}{
					"name":    "linux",
					"version": "20.04",
				},
				"hardware": map[string]interface{}{
					"cpu":    "Intel(R) Core(TM) i7-8700K CPU @ 3.70GHz",
					"memory": "16GB",
					"disk":   "500GB SSD",
				},
				"network": map[string]interface{}{
					"hostname": "example-server",
					"ip":       "192.168.1.100",
				},
			},
			Source: "test-valid-project",
		},
	}

	result := ValidateCustomFacts(customFacts)
	assert.True(t, result.Valid)
	assert.Len(t, result.Errors, 0)
}
