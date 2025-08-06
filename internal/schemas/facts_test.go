package schemas

import (
	"strings"
	"testing"
)

func TestFactsSchemaEmbedding(t *testing.T) {
	// Test that all facts schemas are properly embedded and accessible

	// Test global facts schemas
	globalFactsSchemas := []SchemaType{
		SchemaTypeFactsStructure,
		SchemaTypeFactsStorage,
	}

	for _, schemaType := range globalFactsSchemas {
		content, err := GetSchema(schemaType)
		if err != nil {
			t.Errorf("Failed to get schema %s: %v", schemaType, err)
			continue
		}

		if len(content) == 0 {
			t.Errorf("Schema %s is empty", schemaType)
			continue
		}

		contentStr := string(content)
		if !strings.Contains(contentStr, "# Global Facts") && !strings.Contains(contentStr, "global_facts") {
			t.Errorf("Schema %s doesn't contain expected global facts content", schemaType)
		}
	}

	// Test project facts schemas
	projectFactsSchemas := []SchemaType{
		SchemaTypeFactsStructure,
		SchemaTypeFactsStorage,
	}

	for _, schemaType := range projectFactsSchemas {
		content, err := GetSchema(schemaType)
		if err != nil {
			t.Errorf("Failed to get schema %s: %v", schemaType, err)
			continue
		}

		if len(content) == 0 {
			t.Errorf("Schema %s is empty", schemaType)
			continue
		}

		contentStr := string(content)
		if !strings.Contains(contentStr, "# Project Facts") && !strings.Contains(contentStr, "project_facts") {
			t.Errorf("Schema %s doesn't contain expected project facts content", schemaType)
		}
	}

	// Test custom facts schema
	customContent, err := GetSchema(SchemaTypeCustomFactsHCL)
	if err != nil {
		t.Errorf("Failed to get custom facts schema: %v", err)
	} else if len(customContent) == 0 {
		t.Errorf("Custom facts schema is empty")
	} else {
		contentStr := string(customContent)
		if !strings.Contains(contentStr, "# Custom Facts") && !strings.Contains(contentStr, "custom_facts") {
			t.Errorf("Custom facts schema doesn't contain expected content")
		}
	}

	// Test facts export schema
	exportContent, err := GetSchema(SchemaTypeFactsStructure)
	if err != nil {
		t.Errorf("Failed to get facts export schema: %v", err)
	} else if len(exportContent) == 0 {
		t.Errorf("Facts export schema is empty")
	} else {
		contentStr := string(exportContent)
		if !strings.Contains(contentStr, "# Facts Export") && !strings.Contains(contentStr, "facts_export") {
			t.Errorf("Facts export schema doesn't contain expected content")
		}
	}

	t.Logf("✅ All facts schemas properly embedded and accessible")
}
