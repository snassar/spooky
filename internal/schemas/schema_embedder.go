package schemas

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

//go:embed schemafiles
var embeddedSchemas embed.FS

//go:embed testdata
var embeddedTestData embed.FS

// SchemaEmbedder provides functionality to embed and access HCL schemas
type SchemaEmbedder struct {
	schemas  map[string]string
	rules    map[string]string
	metadata map[string]string
	testdata map[string]string
}

// NewSchemaEmbedder creates a new schema embedder instance
func NewSchemaEmbedder() (*SchemaEmbedder, error) {
	embedder := &SchemaEmbedder{
		schemas:  make(map[string]string),
		rules:    make(map[string]string),
		metadata: make(map[string]string),
		testdata: make(map[string]string),
	}

	if err := embedder.loadEmbeddedSchemas(); err != nil {
		return nil, errors.Wrap(err, "failed to load embedded schemas")
	}

	if err := embedder.loadEmbeddedTestData(); err != nil {
		return nil, errors.Wrap(err, "failed to load embedded test data")
	}

	return embedder, nil
}

// loadEmbeddedSchemas loads all embedded schema files into memory
func (se *SchemaEmbedder) loadEmbeddedSchemas() error {
	return fs.WalkDir(embeddedSchemas, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Only process .hcl files
		if !strings.HasSuffix(path, ".hcl") {
			return nil
		}

		// Read the file content
		content, err := embeddedSchemas.ReadFile(path)
		if err != nil {
			return errors.Wrapf(err, "failed to read embedded file %s", path)
		}

		// Categorize the file based on its path
		dir := filepath.Dir(path)
		baseName := filepath.Base(path)
		key := strings.TrimSuffix(baseName, ".hcl")

		switch {
		case strings.Contains(dir, "structure"):
			se.schemas[key] = string(content)
		case strings.Contains(dir, "validation"):
			// For validation rules, remove the "-rules" suffix for easier access
			cleanKey := strings.TrimSuffix(key, "-rules")
			se.rules[cleanKey] = string(content)
		case strings.Contains(dir, "metadata"):
			se.metadata[key] = string(content)
		}

		return nil
	})
}

// loadEmbeddedTestData loads all embedded test data files into memory
func (se *SchemaEmbedder) loadEmbeddedTestData() error {
	return fs.WalkDir(embeddedTestData, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Only process .hcl and .tmpl files
		if !strings.HasSuffix(path, ".hcl") && !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		// Read the file content
		content, err := embeddedTestData.ReadFile(path)
		if err != nil {
			return errors.Wrapf(err, "failed to read embedded test file %s", path)
		}

		// Use the filename as the key (handle both .hcl and .tmpl extensions)
		baseName := filepath.Base(path)
		key := strings.TrimSuffix(baseName, ".hcl")
		key = strings.TrimSuffix(key, ".tmpl")
		se.testdata[key] = string(content)

		return nil
	})
}

// GetSchema retrieves a schema by name
func (se *SchemaEmbedder) GetSchema(name string) (string, bool) {
	content, exists := se.schemas[name]
	return content, exists
}

// GetValidationRules retrieves validation rules by name
func (se *SchemaEmbedder) GetValidationRules(name string) (string, bool) {
	content, exists := se.rules[name]
	return content, exists
}

// GetMetadata retrieves metadata by name
func (se *SchemaEmbedder) GetMetadata(name string) (string, bool) {
	content, exists := se.metadata[name]
	return content, exists
}

// GetTestData retrieves test data by name
func (se *SchemaEmbedder) GetTestData(name string) (string, bool) {
	content, exists := se.testdata[name]
	return content, exists
}

// ListSchemas returns a list of all available schema names
func (se *SchemaEmbedder) ListSchemas() []string {
	var schemas []string
	for name := range se.schemas {
		schemas = append(schemas, name)
	}
	return schemas
}

// ListValidationRules returns a list of all available validation rule names
func (se *SchemaEmbedder) ListValidationRules() []string {
	var rules []string
	for name := range se.rules {
		rules = append(rules, name)
	}
	return rules
}

// ListMetadata returns a list of all available metadata names
func (se *SchemaEmbedder) ListMetadata() []string {
	var metadata []string
	for name := range se.metadata {
		metadata = append(metadata, name)
	}
	return metadata
}

// ListTestData returns a list of all available test data names
func (se *SchemaEmbedder) ListTestData() []string {
	var testdata []string
	for name := range se.testdata {
		testdata = append(testdata, name)
	}
	return testdata
}

// GetAllSchemas returns all schemas as a map
func (se *SchemaEmbedder) GetAllSchemas() map[string]string {
	result := make(map[string]string)
	for k, v := range se.schemas {
		result[k] = v
	}
	return result
}

// GetAllValidationRules returns all validation rules as a map
func (se *SchemaEmbedder) GetAllValidationRules() map[string]string {
	result := make(map[string]string)
	for k, v := range se.rules {
		result[k] = v
	}
	return result
}

// GetAllMetadata returns all metadata as a map
func (se *SchemaEmbedder) GetAllMetadata() map[string]string {
	result := make(map[string]string)
	for k, v := range se.metadata {
		result[k] = v
	}
	return result
}

// GetAllTestData returns all test data as a map
func (se *SchemaEmbedder) GetAllTestData() map[string]string {
	result := make(map[string]string)
	for k, v := range se.testdata {
		result[k] = v
	}
	return result
}

// GetSchemaWithRules returns both schema and validation rules for a given name
func (se *SchemaEmbedder) GetSchemaWithRules(name string) (schema string, rules string, err error) {
	schemaContent, schemaExists := se.GetSchema(name)
	if !schemaExists {
		available := se.ListSchemas()
		return "", "", NewSchemaNotFoundError(name, available)
	}

	rulesContent, rulesExist := se.GetValidationRules(name)
	if !rulesExist {
		return schemaContent, "", nil // Schema exists but no rules
	}

	return schemaContent, rulesContent, nil
}

// ValidateSchemaExists checks if a schema exists
func (se *SchemaEmbedder) ValidateSchemaExists(name string) bool {
	_, exists := se.schemas[name]
	return exists
}

// GetSchemaInfo returns information about a schema
func (se *SchemaEmbedder) GetSchemaInfo(name string) (*SchemaInfo, error) {
	schema, exists := se.GetSchema(name)
	if !exists {
		available := se.ListSchemas()
		return nil, NewSchemaNotFoundError(name, available)
	}

	info := &SchemaInfo{
		Name:    name,
		Content: schema,
		Size:    len(schema),
	}

	// Check if validation rules exist
	if rules, rulesExist := se.GetValidationRules(name); rulesExist {
		info.HasValidationRules = true
		info.ValidationRulesSize = len(rules)
	}

	// Check if metadata exists
	if metadata, metadataExist := se.GetMetadata(name); metadataExist {
		info.HasMetadata = true
		info.MetadataSize = len(metadata)
	}

	return info, nil
}

// SchemaInfo contains information about a schema
type SchemaInfo struct {
	Name                string
	Content             string
	Size                int
	HasValidationRules  bool
	ValidationRulesSize int
	HasMetadata         bool
	MetadataSize        int
}

// PrintSchemaSummary prints a summary of all embedded schemas
func (se *SchemaEmbedder) PrintSchemaSummary() {
	fmt.Println("=== Embedded Schemas Summary ===")

	fmt.Printf("\n📁 Structure Schemas (%d):\n", len(se.schemas))
	for _, name := range se.ListSchemas() {
		info, _ := se.GetSchemaInfo(name)
		fmt.Printf("  - %s (%d bytes)\n", name, info.Size)
	}

	fmt.Printf("\n✅ Validation Rules (%d):\n", len(se.rules))
	for _, name := range se.ListValidationRules() {
		rules, _ := se.GetValidationRules(name)
		fmt.Printf("  - %s (%d bytes)\n", name, len(rules))
	}

	fmt.Printf("\n📋 Metadata (%d):\n", len(se.metadata))
	for _, name := range se.ListMetadata() {
		metadata, _ := se.GetMetadata(name)
		fmt.Printf("  - %s (%d bytes)\n", name, len(metadata))
	}

	fmt.Printf("\n🧪 Test Data (%d):\n", len(se.testdata))
	for _, name := range se.ListTestData() {
		testdata, _ := se.GetTestData(name)
		fmt.Printf("  - %s (%d bytes)\n", name, len(testdata))
	}

	fmt.Printf("\n📊 Total: %d files embedded\n", len(se.schemas)+len(se.rules)+len(se.metadata)+len(se.testdata))
}
