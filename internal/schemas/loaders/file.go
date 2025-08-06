package loaders

import (
	"fmt"
	"os"
	"path/filepath"
	"spooky/internal/schemas/types"
)

// FileLoader implements loading schemas from files
type FileLoader struct {
	basePath string
	schemas  map[types.SchemaType]*types.Schema
}

// NewFileLoader creates a new file loader
func NewFileLoader(basePath string) *FileLoader {
	return &FileLoader{
		basePath: basePath,
		schemas:  make(map[types.SchemaType]*types.Schema),
	}
}

// LoadSchema loads a schema from a file
func (l *FileLoader) LoadSchema(schemaType types.SchemaType) (*types.Schema, error) {
	// Construct the file path
	filePath := filepath.Join(l.basePath, string(schemaType)+".hcl")

	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file %s: %w", filePath, err)
	}

	schema := &types.Schema{
		Type:     schemaType,
		Content:  string(content),
		Filename: filepath.Base(filePath),
	}

	l.schemas[schemaType] = schema
	return schema, nil
}

// LoadAllSchemas loads all available schemas from the base path
func (l *FileLoader) LoadAllSchemas() (map[types.SchemaType]*types.Schema, error) {
	// This would scan the directory and load all .hcl files
	// For now, return the loaded schemas
	return l.schemas, nil
}

// GetName returns the loader name
func (l *FileLoader) GetName() string {
	return "file"
}
