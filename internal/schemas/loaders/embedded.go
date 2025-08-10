package loaders

import (
	"fmt"
	"spooky/internal/types/schemas"
)

// EmbeddedLoader implements loading schemas from embedded resources
type EmbeddedLoader struct {
	schemas map[types.SchemaType]*types.Schema
}

// NewEmbeddedLoader creates a new embedded loader
func NewEmbeddedLoader() *EmbeddedLoader {
	return &EmbeddedLoader{
		schemas: make(map[types.SchemaType]*types.Schema),
	}
}

// LoadSchema loads a schema from embedded resources
func (l *EmbeddedLoader) LoadSchema(schemaType types.SchemaType) (*types.Schema, error) {
	// This would integrate with the existing embedded schema system
	// For now, return a placeholder
	schema := &types.Schema{
		Type:     schemaType,
		Content:  fmt.Sprintf("// Schema for %s", schemaType),
		Filename: string(schemaType) + ".hcl",
	}

	l.schemas[schemaType] = schema
	return schema, nil
}

// LoadAllSchemas loads all available schemas
func (l *EmbeddedLoader) LoadAllSchemas() (map[types.SchemaType]*types.Schema, error) {
	// This would load all available schemas
	// For now, return the loaded schemas
	return l.schemas, nil
}

// GetName returns the loader name
func (l *EmbeddedLoader) GetName() string {
	return "embedded"
}
