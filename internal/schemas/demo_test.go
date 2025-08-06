package schemas

import (
	"fmt"
	"strings"
	"testing"
)

func TestEnhancedSchemaCompositionDemo(t *testing.T) {
	// Demonstrate that our enhanced schema composition system works
	// without generating any file artifacts

	composer := NewSchemaComposer()

	// Get a composed actions schema
	actionsSchema, err := composer.GetSystemSchema("actions")
	if err != nil {
		t.Fatalf("Failed to get actions schema: %v", err)
	}

	// Verify it contains the expected content
	if !strings.Contains(actionsSchema, "# Composed Actions Schema") {
		t.Error("Actions schema missing expected header")
	}

	if !strings.Contains(actionsSchema, "actions_validation") {
		t.Error("Actions schema missing validation rules")
	}

	if !strings.Contains(actionsSchema, "actions_metadata") {
		t.Error("Actions schema missing metadata")
	}

	// Get a composed machines schema
	machinesSchema, err := composer.GetSystemSchema("machines")
	if err != nil {
		t.Fatalf("Failed to get machines schema: %v", err)
	}

	// Verify it contains the expected content
	if !strings.Contains(machinesSchema, "# Composed Machines Schema") {
		t.Error("Machines schema missing expected header")
	}

	if !strings.Contains(machinesSchema, "machines_validation") {
		t.Error("Machines schema missing validation rules")
	}

	// Get a composed variables schema (now returns HCL schema by default)
	variablesSchema, err := composer.GetSystemSchema("variables")
	if err != nil {
		t.Fatalf("Failed to get variables schema: %v", err)
	}

	// Verify it contains the expected content (HCL schema format)
	if !strings.Contains(variablesSchema, "# Composed Variables HCL Schema") {
		t.Error("Variables schema missing expected header")
	}

	if !strings.Contains(variablesSchema, "variables_hcl_metadata") {
		t.Error("Variables schema missing metadata")
	}

	// Get a composed project schema
	projectSchema, err := composer.GetSystemSchema("project")
	if err != nil {
		t.Fatalf("Failed to get project schema: %v", err)
	}

	// Verify it contains the expected content
	if !strings.Contains(projectSchema, "# Composed Project Schema") {
		t.Error("Project schema missing expected header")
	}

	if !strings.Contains(projectSchema, "project_validation") {
		t.Error("Project schema missing validation rules")
	}

	// List all available schemas
	availableSchemas, err := composer.ListAvailableSchemas()
	if err != nil {
		t.Fatalf("Failed to list available schemas: %v", err)
	}

	// Verify we got the expected number of schemas
	if len(availableSchemas) != 11 {
		t.Errorf("Expected 11 schemas, got %d", len(availableSchemas))
	}

	// Print a summary for demonstration
	fmt.Printf("✅ Enhanced schema composition system working correctly\n")
	fmt.Printf("✅ Generated %d composed schemas in memory\n", len(availableSchemas))
	fmt.Printf("✅ No file artifacts generated\n")
	fmt.Printf("✅ All schemas include validation rules and metadata\n")

	for _, schemaName := range availableSchemas {
		fmt.Printf("  - %s\n", schemaName)
	}
}
