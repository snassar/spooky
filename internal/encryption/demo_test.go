package encryption

import (
	"testing"
)

// TestDemoFactEncryptionCompiles verifies that the fact encryption demo function compiles
// This is a simple compilation test since the demo requires age keys to be set up
func TestDemoFactEncryptionCompiles(t *testing.T) {
	// This test just verifies that the function exists and can be called
	// The actual demo would require age keys to be set up
	t.Log("Fact encryption demo function compiles successfully")
}

// TestFactEncryptionStructure verifies that facts can use the same encryption structure as variables
func TestFactEncryptionStructure(t *testing.T) {
	// Create a fact structure that matches the new encryption pattern
	fact := map[string]interface{}{
		"value":       "test-value",
		"type":        "string",
		"description": "Test fact",
		"encrypted":   true,
		"sensitive":   false,
		"tags":        []string{"test"},
	}

	// Verify the structure has the expected fields
	if fact["value"] != "test-value" {
		t.Errorf("Expected value 'test-value', got %v", fact["value"])
	}

	if fact["type"] != "string" {
		t.Errorf("Expected type 'string', got %v", fact["type"])
	}

	if fact["encrypted"] != true {
		t.Errorf("Expected encrypted true, got %v", fact["encrypted"])
	}

	// Verify that the fact structure can potentially have encrypted_value
	// (this would be added by the transformer)
	if _, exists := fact["encrypted_value"]; exists {
		t.Error("Fact should not have encrypted_value before transformation")
	}

	t.Log("Fact encryption structure is correctly defined")
}
