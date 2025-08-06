package embedded

import (
	"fmt"
	"spooky/internal/schemas/types"
)

// SchemaDefinitions holds embedded schema definitions
type SchemaDefinitions struct {
	schemas map[types.SchemaType]string
}

// NewSchemaDefinitions creates a new schema definitions instance
func NewSchemaDefinitions() *SchemaDefinitions {
	return &SchemaDefinitions{
		schemas: make(map[types.SchemaType]string),
	}
}

// GetSchema retrieves a schema by type
func (sd *SchemaDefinitions) GetSchema(schemaType types.SchemaType) (string, error) {
	// This would load from the embedded files
	// For now, return a placeholder
	content, exists := sd.schemas[schemaType]
	if !exists {
		return "", fmt.Errorf("schema %s not found", schemaType)
	}
	return content, nil
}

// LoadAllSchemas loads all embedded schemas
func (sd *SchemaDefinitions) LoadAllSchemas() error {
	// This would scan the embedded files and load all schemas
	// For now, just initialize with some placeholder schemas
	sd.schemas[types.SchemaType("machines")] = "// Machines schema placeholder"
	sd.schemas[types.SchemaType("variables")] = "// Variables schema placeholder"
	sd.schemas[types.SchemaType("actions")] = "// Actions schema placeholder"
	return nil
}

// ListSchemas returns all available schema types
func (sd *SchemaDefinitions) ListSchemas() []types.SchemaType {
	var schemas []types.SchemaType
	for schemaType := range sd.schemas {
		schemas = append(schemas, schemaType)
	}
	return schemas
}
