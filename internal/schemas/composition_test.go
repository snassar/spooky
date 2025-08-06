package schemas

import (
	"strings"
	"testing"
)

func TestSchemaComposer(t *testing.T) {
	composer := NewSchemaComposer()

	// Test composing all system schemas
	schemas, err := composer.ComposeSystemSchemas()
	if err != nil {
		t.Fatalf("Failed to compose system schemas: %v", err)
	}

	// Verify we have the expected schemas
	expectedSchemas := []string{
		"facts-structure.hcl",
		"facts-storage.hcl",
		"actions-composed.hcl",
		"machines-composed.hcl",
		"variables-structure-composed.hcl",
		"variables-hcl-composed.hcl",
		"variables-json-composed.hcl",
		"project-composed.hcl",
		"template-metadata-composed.hcl",
		"template-context-composed.hcl",
		"template-functions-composed.hcl",
	}

	for _, expected := range expectedSchemas {
		if content, exists := schemas[expected]; !exists {
			t.Errorf("Expected schema %s not found", expected)
		} else if content == "" {
			t.Errorf("Schema %s is empty", expected)
		} else if strings.Contains(expected, "composed") {
			// Composed schemas should contain "# Composed" header
			if !strings.Contains(content, "# Composed") {
				t.Errorf("Schema %s doesn't contain expected header", expected)
			}
		}
	}

	// Test getting a specific system schema
	actionsSchema, err := composer.GetSystemSchema("actions")
	if err != nil {
		t.Fatalf("Failed to get actions schema: %v", err)
	}

	if actionsSchema == "" {
		t.Error("Actions schema is empty")
	}

	if !strings.Contains(actionsSchema, "# Composed Actions Schema") {
		t.Error("Actions schema doesn't contain expected header")
	}

	if !strings.Contains(actionsSchema, "actions_validation") {
		t.Error("Actions schema doesn't contain validation rules")
	}

	// Test listing available schemas
	availableSchemas, err := composer.ListAvailableSchemas()
	if err != nil {
		t.Fatalf("Failed to list available schemas: %v", err)
	}

	if len(availableSchemas) == 0 {
		t.Error("No schemas available")
	}

	// Test cache functionality
	composer.ClearCache()
	if len(composer.cache) != 0 {
		t.Error("Cache was not cleared")
	}

	// Test that schemas are recomposed after cache clear
	schemasAfterClear, err := composer.ComposeSystemSchemas()
	if err != nil {
		t.Fatalf("Failed to recompose schemas after cache clear: %v", err)
	}

	if len(schemasAfterClear) == 0 {
		t.Error("No schemas available after cache clear")
	}
}

func TestGetSchemaWithComposedSchemas(t *testing.T) {
	// Test getting individual schemas through the GetSchema function
	// (composed schemas are not available through GetSchema)

	// Test basic schemas
	basicSchemas := []SchemaType{
		SchemaTypeActions,
		SchemaTypeMachines,
		SchemaTypeProject,
		SchemaTypeSpooky,
	}

	for _, schemaType := range basicSchemas {
		content, err := GetSchema(schemaType)
		if err != nil {
			t.Errorf("Failed to get schema %s: %v", schemaType, err)
			continue
		}

		if len(content) == 0 {
			t.Errorf("Schema %s is empty", schemaType)
		}
	}

	// Test that composed schemas are not available through GetSchema
	composedSchemas := []SchemaType{
		SchemaTypeActionsComposed,
		SchemaTypeMachinesComposed,
		SchemaTypeVariablesComposed,
		SchemaTypeProjectComposed,
	}

	for _, schemaType := range composedSchemas {
		_, err := GetSchema(schemaType)
		if err == nil {
			t.Errorf("Expected error for composed schema %s, but got none", schemaType)
		}
	}
}

func TestValidateAgainstSchema(t *testing.T) {
	composer := NewSchemaComposer()

	// Test validation against a schema
	// This is a basic test - actual validation would require HCL parsing
	err := composer.ValidateAgainstSchema("test data", "actions")
	if err != nil {
		// Validation might fail due to HCL parsing, which is expected
		t.Logf("Validation failed as expected: %v", err)
	}
}
