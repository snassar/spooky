// Package schemas provides schema validation and management functionality for the spooky codebase.
package schemas

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"

	spookytypes "spooky/internal/types"
	spookytypescommon "spooky/internal/types/common"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesschemas "spooky/internal/types/schemas"
)

// Manager provides comprehensive schema management functionality
type Manager struct {
	logger            spookytypeslogging.Logger
	helpers           *SchemaHelpers
	registry          *Registry
	validator         *Validator
	enhancedValidator *EnhancedValidator
	evolutionManager  *EvolutionManager

	// Metadata validation
	metadataSchema *spookytypesschemas.Schema
	parser         *hclparse.Parser
}

// NewManager creates a new schema manager instance
func NewManager(logger spookytypeslogging.Logger) *Manager {
	registry := NewRegistry(logger)
	validator := NewValidator(logger)
	validator.SetRegistry(registry)

	enhancedValidator := NewEnhancedValidator(logger, nil)
	enhancedValidator.SetRegistry(registry)

	evolutionManager := NewEvolutionManager(logger, registry)

	manager := &Manager{
		logger:            logger,
		helpers:           NewSchemaHelpers(logger),
		registry:          registry,
		validator:         validator,
		enhancedValidator: enhancedValidator,
		evolutionManager:  evolutionManager,
		parser:            hclparse.NewParser(),
	}

	// Set the manager as the validator's manager for metadata validation
	validator.SetManager(manager)

	return manager
}

// Load loads a schema from a file
func (m *Manager) Load(filePath string) (*spookytypesschemas.Schema, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("schema file does not exist: %s", filePath)
	}

	// Read file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	// Create schema from file content
	schema := &spookytypesschemas.Schema{
		Version:     "1.0",
		Type:        "hcl",
		Name:        filepath.Base(filePath),
		Description: fmt.Sprintf("Schema loaded from %s", filePath),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Content:     string(data),
		Metadata:    make(map[string]interface{}),
	}

	// Parse HCL content to validate it
	if err := m.validateSchemaContent(data, filePath); err != nil {
		return nil, fmt.Errorf("failed to parse schema file: %w", err)
	}

	// Validate schema metadata if this is not the metadata schema itself
	if !strings.HasSuffix(filepath.Base(filePath), "schema-metadata.schema.hcl") {
		if err := m.validateSchemaMetadata(data, filePath); err != nil {
			return nil, fmt.Errorf("schema metadata validation failed: %w", err)
		}
	}

	return schema, nil
}

// LoadFromString loads a schema from a string
func (m *Manager) LoadFromString(content string) (*spookytypesschemas.Schema, error) {
	if content == "" {
		return nil, fmt.Errorf("content cannot be empty")
	}

	schema := &spookytypesschemas.Schema{
		Version:     "1.0",
		Type:        "hcl",
		Name:        "string-schema",
		Description: "Schema loaded from string content",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Content:     content,
		Metadata:    make(map[string]interface{}),
	}

	// Validate content
	if err := m.validateSchemaContent([]byte(content), "string"); err != nil {
		return nil, fmt.Errorf("failed to parse schema content: %w", err)
	}

	// Validate schema metadata
	if err := m.validateSchemaMetadata([]byte(content), "string"); err != nil {
		return nil, fmt.Errorf("schema metadata validation failed: %w", err)
	}

	return schema, nil
}

// LoadFromBytes loads a schema from bytes
func (m *Manager) LoadFromBytes(data []byte) (*spookytypesschemas.Schema, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	schema := &spookytypesschemas.Schema{
		Version:     "1.0",
		Type:        "hcl",
		Name:        "bytes-schema",
		Description: "Schema loaded from byte data",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Content:     string(data),
		Metadata:    make(map[string]interface{}),
	}

	// Validate content
	if err := m.validateSchemaContent(data, "bytes"); err != nil {
		return nil, fmt.Errorf("failed to parse schema data: %w", err)
	}

	// Validate schema metadata
	if err := m.validateSchemaMetadata(data, "bytes"); err != nil {
		return nil, fmt.Errorf("schema metadata validation failed: %w", err)
	}

	return schema, nil
}

// LoadEmbedded loads an embedded schema
func (m *Manager) LoadEmbedded(name string) (*spookytypesschemas.Schema, error) {
	if name == "" {
		return nil, fmt.Errorf("schema name cannot be empty")
	}

	// Look for embedded schema in the schemas directory
	schemasDir := filepath.Join("internal", "schemas", "schemas")
	schemaPath := filepath.Join(schemasDir, name+".schema.hcl")

	return m.Load(schemaPath)
}

// LoadEmbeddedSchema loads an embedded schema (interface compatibility)
func (m *Manager) LoadEmbeddedSchema(_ context.Context, schemaName string) (*spookytypes.Schema, error) {
	schema, err := m.LoadEmbedded(schemaName)
	if err != nil {
		return nil, err
	}

	// Convert from spookytypesschemas.Schema to spookytypes.Schema
	return &spookytypes.Schema{
		Name:    schema.Name,
		Content: schema.Content,
		// Add other fields as needed
	}, nil
}

// LoadSchema loads a schema from the given path (interface compatibility)
func (m *Manager) LoadSchema(_ context.Context, schemaPath string) (*spookytypes.Schema, error) {
	schema, err := m.Load(schemaPath)
	if err != nil {
		return nil, err
	}

	// Convert from spookytypesschemas.Schema to spookytypes.Schema
	return &spookytypes.Schema{
		Name:    schema.Name,
		Content: schema.Content,
		// Add other fields as needed
	}, nil
}

// LoadWithValidation loads a schema with validation
func (m *Manager) LoadWithValidation(filePath string) (*spookytypesschemas.Schema, *spookytypesschemas.ValidationResult, error) {
	schema, err := m.Load(filePath)
	if err != nil {
		return nil, nil, err
	}

	// Validate the schema itself
	result, err := m.validator.Validate(schema, schema.Content)
	if err != nil {
		return schema, nil, err
	}

	return schema, result, nil
}

// LoadMultiple loads multiple schemas
func (m *Manager) LoadMultiple(filePaths []string) (map[string]*spookytypesschemas.Schema, error) {
	schemas := make(map[string]*spookytypesschemas.Schema)

	for _, filePath := range filePaths {
		schema, err := m.Load(filePath)
		if err != nil {
			m.logger.Warn("Failed to load schema", map[string]interface{}{
				"file":  filePath,
				"error": err.Error(),
			})
			continue
		}

		schemaName := strings.TrimSuffix(filepath.Base(filePath), ".schema.hcl")
		schemas[schemaName] = schema
	}

	return schemas, nil
}

// LoadSchemasFromDirectory loads all schema files from a directory with metadata validation
func (m *Manager) LoadSchemasFromDirectory(dirPath string) (map[string]*spookytypesschemas.Schema, error) {
	if dirPath == "" {
		return nil, fmt.Errorf("directory path cannot be empty")
	}

	// Check if directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("schema directory does not exist: %s", dirPath)
	}

	var schemaFiles []string

	// Walk through directory to find schema files
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process schema files
		if strings.HasSuffix(info.Name(), ".schema.hcl") {
			schemaFiles = append(schemaFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk schema directory: %w", err)
	}

	m.logger.Info("Found schema files", map[string]interface{}{
		"directory": dirPath,
		"count":     len(schemaFiles),
		"files":     schemaFiles,
	})

	// Load all schemas with metadata validation
	return m.LoadMultiple(schemaFiles)
}

// Validate validates data against a schema
func (m *Manager) Validate(_ context.Context, schema *spookytypes.Schema, data interface{}) (*spookytypes.ValidationResult, error) {
	// Convert from spookytypes.Schema to spookytypesschemas.Schema
	schemasSchema := &spookytypesschemas.Schema{
		Name:    schema.Name,
		Content: schema.Content,
		// Add other fields as needed
	}

	result, err := m.validator.Validate(schemasSchema, data)
	if err != nil {
		return nil, err
	}

	// Convert from spookytypesschemas.ValidationResult to spookytypes.ValidationResult
	return &spookytypes.ValidationResult{
		Valid:       result.Valid,
		ValidatedAt: result.ValidatedAt,
		Errors:      result.Errors,
		Warnings:    result.Warnings,
		Info:        result.Info,
	}, nil
}

// ValidateWithEnhancedFeatures validates data with enhanced features
func (m *Manager) ValidateWithEnhancedFeatures(_ context.Context, schema *spookytypesschemas.Schema, data interface{}) (*spookytypesschemas.ValidationResult, error) {
	return m.enhancedValidator.ValidateWithEnhancedFeatures(context.Background(), schema, data)
}

// ValidateFile validates a file against a schema
func (m *Manager) ValidateFile(schema *spookytypesschemas.Schema, filePath string) (*spookytypesschemas.ValidationResult, error) {
	return m.validator.ValidateFileAgainstSchema(schema, filePath)
}

// ValidateString validates a string against a schema
func (m *Manager) ValidateString(schema *spookytypesschemas.Schema, content string) (*spookytypesschemas.ValidationResult, error) {
	return m.validator.ValidateString(schema, content)
}

// ValidateBytes validates bytes against a schema
func (m *Manager) ValidateBytes(schema *spookytypesschemas.Schema, data []byte) (*spookytypesschemas.ValidationResult, error) {
	return m.validator.ValidateBytes(schema, data)
}

// ValidateWithContext validates data with additional context
func (m *Manager) ValidateWithContext(schema *spookytypesschemas.Schema, data interface{}, context map[string]interface{}) (*spookytypesschemas.ValidationResult, error) {
	return m.validator.ValidateWithContext(schema, data, context)
}

// ValidateField validates a specific field
func (m *Manager) ValidateField(schema *spookytypesschemas.Schema, fieldPath string, value interface{}) (*spookytypesschemas.ValidationResult, error) {
	return m.validator.ValidateField(schema, fieldPath, value)
}

// Register registers a new schema
func (m *Manager) Register(_ context.Context, schema *spookytypes.Schema) error {
	// Convert from spookytypes.Schema to spookytypesschemas.Schema
	schemasSchema := &spookytypesschemas.Schema{
		Name:    schema.Name,
		Content: schema.Content,
		// Add other fields as needed
	}

	return m.registry.Register(schemasSchema)
}

// Get returns a schema by name and type
func (m *Manager) Get(name, schemaType string) (*spookytypesschemas.Schema, bool) {
	return m.registry.Get(name, schemaType)
}

// List returns all registered schemas
func (m *Manager) List() []*spookytypesschemas.Schema {
	return m.registry.List()
}

// ListByType returns schemas by type
func (m *Manager) ListByType(schemaType string) []*spookytypesschemas.Schema {
	return m.registry.ListByType(schemaType)
}

// ValidateData validates data against a schema by name and type
func (m *Manager) ValidateData(schemaName, schemaType string, data interface{}) (*spookytypesschemas.ValidationResult, error) {
	return m.registry.ValidateData(schemaName, schemaType, data)
}

// CheckCompatibility checks compatibility between schemas
func (m *Manager) CheckCompatibility(schema1, schema2 *spookytypesschemas.Schema) (*spookytypesschemas.CompatibilityResult, error) {
	return m.evolutionManager.CheckCompatibility(schema1, schema2)
}

// TrackSchemaEvolution tracks schema evolution
func (m *Manager) TrackSchemaEvolution(schema *spookytypesschemas.Schema) error {
	return m.evolutionManager.TrackSchemaEvolution(schema)
}

// GetBreakingChanges gets breaking changes for a schema
func (m *Manager) GetBreakingChanges(schema *spookytypesschemas.Schema) ([]spookytypesschemas.BreakingChange, error) {
	return m.evolutionManager.GetBreakingChanges(schema)
}

// GetDeprecations gets deprecations for a schema
func (m *Manager) GetDeprecations(schema *spookytypesschemas.Schema) ([]spookytypesschemas.SchemaDeprecation, error) {
	return m.evolutionManager.GetDeprecations(schema)
}

// GenerateMigrationScript generates a migration script
func (m *Manager) GenerateMigrationScript(fromSchema, toSchema *spookytypesschemas.Schema) (string, error) {
	return m.evolutionManager.GenerateMigrationScript(fromSchema, toSchema)
}

// LoadSchemas loads schemas from a directory
func (m *Manager) LoadSchemas(schemasDir string) error {
	return m.validator.LoadSchemas(schemasDir)
}

// GetSchemaStatistics returns schema statistics
func (m *Manager) GetSchemaStatistics() map[string]interface{} {
	return m.registry.GetSchemaStatistics()
}

// GetMigrationPath gets migration path between schema versions
func (m *Manager) GetMigrationPath(fromVersion, toVersion string) ([]*spookytypesschemas.SchemaMigration, error) {
	return m.registry.GetMigrationPath(fromVersion, toVersion)
}

// GetMigrationPathBetweenSchemas gets migration path between two schemas
func (m *Manager) GetMigrationPathBetweenSchemas(fromSchema, toSchema *spookytypesschemas.Schema) ([]spookytypesschemas.SchemaMigration, error) {
	return m.evolutionManager.GetMigrationPath(fromSchema, toSchema)
}

// CheckForUpdates checks for schema updates
func (m *Manager) CheckForUpdates(schema *spookytypesschemas.Schema) ([]spookytypesschemas.SchemaUpdate, error) {
	return m.helpers.CheckForUpdates(schema, m.registry.List())
}

// ApplyMigration applies a migration to a schema
func (m *Manager) ApplyMigration(schema *spookytypesschemas.Schema, migration *spookytypesschemas.SchemaMigration) error {
	if schema == nil || migration == nil {
		return fmt.Errorf("schema and migration cannot be nil")
	}

	// Validate migration before applying
	result, err := m.ValidateMigration(schema, migration)
	if err != nil {
		return fmt.Errorf("migration validation failed: %w", err)
	}

	if !result.Valid {
		return fmt.Errorf("migration validation failed: %d errors", len(result.Errors))
	}

	// Apply migration (simplified - would need actual migration logic)
	schema.Version = migration.ToVersion
	schema.UpdatedAt = time.Now()

	m.logger.Info("Migration applied", map[string]interface{}{
		"schema_name":    schema.Name,
		"from_version":   migration.FromVersion,
		"to_version":     migration.ToVersion,
		"migration_name": migration.Name,
	})

	return nil
}

// ValidateMigration validates a migration
func (m *Manager) ValidateMigration(schema *spookytypesschemas.Schema, migration *spookytypesschemas.SchemaMigration) (*spookytypesschemas.ValidationResult, error) {
	return m.helpers.ValidateMigration(schema, migration)
}

// AnalyzeSchema analyzes a schema
func (m *Manager) AnalyzeSchema(schema *spookytypesschemas.Schema) (*spookytypesschemas.SchemaAnalysis, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema cannot be nil")
	}

	analysis := &spookytypesschemas.SchemaAnalysis{
		Complexity: &spookytypesschemas.SchemaComplexity{},
		Coverage:   &spookytypesschemas.SchemaCoverage{},
		Quality:    &spookytypesschemas.SchemaQuality{},
	}

	// Analyze complexity
	if schema.Validation != nil {
		if schema.Validation.Fields != nil {
			analysis.Complexity.FieldCount = len(schema.Validation.Fields)
		}
		if schema.Validation.Rules != nil {
			analysis.Complexity.RuleCount = len(schema.Validation.Rules)
		}
	}

	// Calculate complexity score
	analysis.Complexity.Score = m.calculateComplexityScore(analysis.Complexity)

	// Analyze coverage
	analysis.Coverage = m.calculateCoverage(schema)

	// Analyze quality
	analysis.Quality = m.calculateQuality(schema)

	// Generate recommendations
	analysis.Recommendations = m.generateAnalysisRecommendations(analysis)

	return analysis, nil
}

// CompareSchemas compares two schemas
func (m *Manager) CompareSchemas(schema1, schema2 *spookytypesschemas.Schema) (*spookytypesschemas.SchemaComparison, error) {
	if schema1 == nil || schema2 == nil {
		return nil, fmt.Errorf("both schemas must be provided")
	}

	comparison := &spookytypesschemas.SchemaComparison{
		Compatible: true,
	}

	// Compare basic properties
	if schema1.Type != schema2.Type {
		comparison.Summary = "Schemas have different types"
		comparison.Compatible = false
	} else {
		comparison.Summary = "Schemas are of the same type"
	}

	// Compare fields
	if schema1.Validation != nil && schema2.Validation != nil {
		comparison.Added, comparison.Removed, comparison.Modified = m.compareFields(
			schema1.Validation.Fields, schema2.Validation.Fields)
	}

	// Check for breaking changes
	if schema2.Evolution != nil {
		for i := range schema2.Evolution.BreakingChanges {
			breakingChange := &schema2.Evolution.BreakingChanges[i]
			comparison.BreakingChanges = append(comparison.BreakingChanges, breakingChange.Field)
		}
		if len(comparison.BreakingChanges) > 0 {
			comparison.Compatible = false
		}
	}

	return comparison, nil
}

// GenerateDocumentation generates documentation for a schema
func (m *Manager) GenerateDocumentation(schema *spookytypesschemas.Schema) (string, error) {
	if schema == nil {
		return "", fmt.Errorf("schema cannot be nil")
	}

	var doc strings.Builder

	// Schema header
	doc.WriteString(fmt.Sprintf("# %s\n\n", schema.Name))
	doc.WriteString(fmt.Sprintf("**Type:** %s\n", schema.Type))
	doc.WriteString(fmt.Sprintf("**Version:** %s\n", schema.Version))
	doc.WriteString(fmt.Sprintf("**Description:** %s\n\n", schema.Description))

	// Validation information
	if schema.Validation != nil {
		doc.WriteString("## Validation\n\n")
		doc.WriteString(fmt.Sprintf("**Mode:** %s\n", schema.Validation.Mode))
		doc.WriteString(fmt.Sprintf("**Enabled:** %t\n\n", schema.Validation.Enabled))

		// Field documentation
		if schema.Validation.Fields != nil {
			doc.WriteString("### Fields\n\n")
			for fieldName, field := range schema.Validation.Fields {
				doc.WriteString(fmt.Sprintf("#### %s\n", fieldName))
				doc.WriteString(fmt.Sprintf("**Type:** %s\n", field.Type))
				doc.WriteString(fmt.Sprintf("**Required:** %t\n", field.Required))
				if field.Description != "" {
					doc.WriteString(fmt.Sprintf("**Description:** %s\n", field.Description))
				}
				doc.WriteString("\n")
			}
		}

		// Validation rules
		if len(schema.Validation.Rules) > 0 {
			doc.WriteString("### Validation Rules\n\n")
			for i := range schema.Validation.Rules {
				rule := &schema.Validation.Rules[i]
				doc.WriteString(fmt.Sprintf("#### %s\n", rule.Name))
				doc.WriteString(fmt.Sprintf("**Type:** %s\n", rule.Type))
				doc.WriteString(fmt.Sprintf("**Severity:** %s\n", rule.Severity))
				doc.WriteString(fmt.Sprintf("**Message:** %s\n\n", rule.Message))
			}
		}
	}

	// Evolution information
	if schema.Evolution != nil {
		doc.WriteString("## Evolution\n\n")

		if len(schema.Evolution.Deprecations) > 0 {
			doc.WriteString("### Deprecations\n\n")
			for i := range schema.Evolution.Deprecations {
				deprecation := &schema.Evolution.Deprecations[i]
				doc.WriteString(fmt.Sprintf("- **%s:** %s\n", deprecation.Field, deprecation.Reason))
				if deprecation.Replacement != "" {
					doc.WriteString(fmt.Sprintf("  - **Replacement:** %s\n", deprecation.Replacement))
				}
				doc.WriteString("\n")
			}
		}

		if len(schema.Evolution.BreakingChanges) > 0 {
			doc.WriteString("### Breaking Changes\n\n")
			for i := range schema.Evolution.BreakingChanges {
				breakingChange := &schema.Evolution.BreakingChanges[i]
				doc.WriteString(fmt.Sprintf("- **%s:** %s\n", breakingChange.Field, breakingChange.Description))
				doc.WriteString(fmt.Sprintf("  - **Impact:** %s\n", breakingChange.Impact))
				if breakingChange.Mitigation != "" {
					doc.WriteString(fmt.Sprintf("  - **Mitigation:** %s\n", breakingChange.Mitigation))
				}
				doc.WriteString("\n")
			}
		}
	}

	return doc.String(), nil
}

// GenerateExamples generates examples for a schema
func (m *Manager) GenerateExamples(schema *spookytypesschemas.Schema) ([]string, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema cannot be nil")
	}

	var examples []string

	// Generate basic example
	basicExample := m.generateBasicExample(schema)
	examples = append(examples, basicExample)

	// Generate field-specific examples
	if schema.Validation != nil && schema.Validation.Fields != nil {
		for fieldName, field := range schema.Validation.Fields {
			if len(field.Examples) > 0 {
				fieldExample := m.generateFieldExample(fieldName, field)
				examples = append(examples, fieldExample)
			}
		}
	}

	return examples, nil
}

// Helper methods

// validateSchemaContent validates schema content
func (m *Manager) validateSchemaContent(data []byte, _ string) error {
	// Basic validation - check if content is not empty
	if len(data) == 0 {
		return fmt.Errorf("schema content is empty")
	}

	// Check for basic HCL structure
	content := string(data)
	if !strings.Contains(content, "{") || !strings.Contains(content, "}") {
		return fmt.Errorf("schema content does not appear to be valid HCL")
	}

	return nil
}

// validateSchemaMetadata validates schema metadata against the meta-schema
func (m *Manager) validateSchemaMetadata(data []byte, source string) error {
	// Load the metadata schema if not already loaded
	if m.metadataSchema == nil {
		if err := m.loadMetadataSchema(); err != nil {
			return fmt.Errorf("failed to load metadata schema: %w", err)
		}
	}

	// Extract metadata block from the schema content
	metadataBlock, err := m.extractMetadataBlock(data, source)
	if err != nil {
		return fmt.Errorf("failed to extract metadata block: %w", err)
	}

	// If no metadata block found, that's an error
	if metadataBlock == nil {
		return fmt.Errorf("schema must contain a metadata block")
	}

	// Validate the metadata block using simple validation rules
	if err := m.validateMetadataBlock(metadataBlock); err != nil {
		return fmt.Errorf("metadata validation failed: %w", err)
	}

	m.logger.Debug("Schema metadata validation passed", map[string]interface{}{
		"source": source,
	})

	return nil
}

// validateMetadataBlock validates a metadata block using simple validation rules
func (m *Manager) validateMetadataBlock(metadata map[string]interface{}) error {
	var errors []string

	// Check required fields
	if err := m.validateRequiredFields(metadata); err != nil {
		errors = append(errors, err.Error())
	}

	// Validate schema_version format
	if err := m.validateSchemaVersion(metadata); err != nil {
		errors = append(errors, err.Error())
	}

	// Validate schema_type format
	if err := m.validateSchemaType(metadata); err != nil {
		errors = append(errors, err.Error())
	}

	// Validate last_updated format
	if err := m.validateLastUpdated(metadata); err != nil {
		errors = append(errors, err.Error())
	}

	// Validate description length
	if err := m.validateDescription(metadata); err != nil {
		errors = append(errors, err.Error())
	}

	// If there are any errors, return them
	if len(errors) > 0 {
		return fmt.Errorf("metadata validation failed: %s", strings.Join(errors, "; "))
	}

	return nil
}

// isValidScalVerFormat checks if a string is in ScalVer format
func (m *Manager) isValidScalVerFormat(version string) bool {
	return spookytypescommon.IsValidScalVerFormat(version)
}

// isValidSchemaTypeFormat checks if a string is a valid schema type format
func (m *Manager) isValidSchemaTypeFormat(schemaType string) bool {
	// Simple regex check for lowercase with hyphens
	matched, _ := regexp.MatchString(`^[a-z][a-z0-9-]*$`, schemaType)
	return matched
}

// isValidDateFormat checks if a string is in YYYY-MM-DD format
func (m *Manager) isValidDateFormat(date string) bool {
	// Simple regex check for YYYY-MM-DD format
	matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, date)
	return matched
}

// validateRequiredFields validates required fields in metadata
func (m *Manager) validateRequiredFields(metadata map[string]interface{}) error {
	requiredFields := []string{"schema_version", "schema_type", "schema_name", "last_updated", "description"}
	for _, field := range requiredFields {
		if value, exists := metadata[field]; !exists || value == nil {
			return fmt.Errorf("required field '%s' is missing", field)
		}
	}
	return nil
}

// validateSchemaVersion validates schema_version format
func (m *Manager) validateSchemaVersion(metadata map[string]interface{}) error {
	if version, exists := metadata["schema_version"]; exists && version != nil {
		if versionStr, ok := version.(string); ok {
			if !m.isValidScalVerFormat(versionStr) {
				return fmt.Errorf("schema_version must be in ScalVer format (0.YYYYMMDD.N), got: %s", versionStr)
			}
		} else {
			return fmt.Errorf("schema_version must be a string")
		}
	}
	return nil
}

// validateSchemaType validates schema_type format
func (m *Manager) validateSchemaType(metadata map[string]interface{}) error {
	if schemaType, exists := metadata["schema_type"]; exists && schemaType != nil {
		if typeStr, ok := schemaType.(string); ok {
			if !m.isValidSchemaTypeFormat(typeStr) {
				return fmt.Errorf("schema_type must be lowercase with hyphens only, got: %s", typeStr)
			}
		} else {
			return fmt.Errorf("schema_type must be a string")
		}
	}
	return nil
}

// validateLastUpdated validates last_updated format
func (m *Manager) validateLastUpdated(metadata map[string]interface{}) error {
	if lastUpdated, exists := metadata["last_updated"]; exists && lastUpdated != nil {
		if dateStr, ok := lastUpdated.(string); ok {
			if !m.isValidDateFormat(dateStr) {
				return fmt.Errorf("last_updated must be in YYYY-MM-DD format, got: %s", dateStr)
			}
		} else {
			return fmt.Errorf("last_updated must be a string")
		}
	}
	return nil
}

// validateDescription validates description length
func (m *Manager) validateDescription(metadata map[string]interface{}) error {
	if description, exists := metadata["description"]; exists && description != nil {
		if descStr, ok := description.(string); ok {
			if len(descStr) < 10 {
				return fmt.Errorf("description must be at least 10 characters long")
			}
			if len(descStr) > 500 {
				return fmt.Errorf("description must be at most 500 characters long")
			}
		} else {
			return fmt.Errorf("description must be a string")
		}
	}
	return nil
}

// loadMetadataSchema loads the schema-metadata.schema.hcl file
func (m *Manager) loadMetadataSchema() error {
	// Look for the metadata schema in the schemas directory
	// Try multiple possible paths
	possiblePaths := []string{
		filepath.Join("internal", "schemas", "schemas", "schema-metadata.schema.hcl"),
		filepath.Join("schemas", "schema-metadata.schema.hcl"),
		"schema-metadata.schema.hcl",
	}

	var metadataSchemaPath string
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			metadataSchemaPath = path
			break
		}
	}

	if metadataSchemaPath == "" {
		return fmt.Errorf("metadata schema file not found in any of the expected locations: %v", possiblePaths)
	}

	// Load the metadata schema without metadata validation to avoid infinite recursion
	schema, err := m.loadSchemaWithoutMetadataValidation(metadataSchemaPath)
	if err != nil {
		return fmt.Errorf("failed to load metadata schema: %w", err)
	}

	m.metadataSchema = schema
	return nil
}

// loadSchemaWithoutMetadataValidation loads a schema without validating its metadata
// This is used to load the metadata schema itself to avoid infinite recursion
func (m *Manager) loadSchemaWithoutMetadataValidation(filePath string) (*spookytypesschemas.Schema, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("schema file does not exist: %s", filePath)
	}

	// Read file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	// Create schema from file content
	schema := &spookytypesschemas.Schema{
		Version:     "1.0",
		Type:        "hcl",
		Name:        filepath.Base(filePath),
		Description: fmt.Sprintf("Schema loaded from %s", filePath),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Content:     string(data),
		Metadata:    make(map[string]interface{}),
	}

	// Parse HCL content to validate it
	if err := m.validateSchemaContent(data, filePath); err != nil {
		return nil, fmt.Errorf("failed to parse schema file: %w", err)
	}

	return schema, nil
}

// extractMetadataBlock extracts the metadata block from schema content
func (m *Manager) extractMetadataBlock(data []byte, source string) (map[string]interface{}, error) {
	// Parse HCL content
	file, diags := m.parser.ParseHCL(data, source)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL content: %s", diags.Error())
	}

	// Use a more direct approach - parse the body without schema validation
	body := file.Body

	// Try to get all blocks using a very permissive approach
	// We'll use the raw body and look for metadata blocks manually
	content, diags := body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type: "metadata",
			},
		},
		// Don't specify any attributes to be permissive
	})

	// If we get errors, it might be because of other blocks
	// Let's try a different approach - parse the raw content
	if diags.HasErrors() {
		// Try to parse with a more permissive approach
		// Look for metadata block in the raw content
		contentStr := string(data)
		if !strings.Contains(contentStr, "metadata {") {
			return nil, nil // No metadata block found
		}

		// If we found metadata block in the string, try to extract it
		// This is a fallback approach
		return m.extractMetadataFromString(contentStr)
	}

	// Look for metadata block
	for _, block := range content.Blocks {
		if block.Type == "metadata" {
			metadata, err := m.parseMetadataBlock(block)
			if err != nil {
				return nil, fmt.Errorf("failed to parse metadata block: %w", err)
			}
			return metadata, nil
		}
	}

	// No metadata block found
	return nil, nil
}

// extractMetadataFromString extracts metadata from HCL string content
// This is a fallback method when strict parsing fails
func (m *Manager) extractMetadataFromString(content string) (map[string]interface{}, error) {
	// Simple string-based extraction for testing purposes
	// In a real implementation, this would be more sophisticated

	// Look for metadata block
	metadataStart := strings.Index(content, "metadata {")
	if metadataStart == -1 {
		return nil, nil
	}

	// Find the end of the metadata block
	braceCount := 0
	metadataEnd := -1
	for i := metadataStart; i < len(content); i++ {
		if content[i] == '{' {
			braceCount++
		} else if content[i] == '}' {
			braceCount--
			if braceCount == 0 {
				metadataEnd = i + 1
				break
			}
		}
	}

	if metadataEnd == -1 {
		return nil, fmt.Errorf("unclosed metadata block")
	}

	// Extract the metadata block content
	metadataContent := content[metadataStart:metadataEnd]

	// Parse the metadata content using simple string parsing
	metadata := make(map[string]interface{})

	// Extract key-value pairs from the metadata block
	lines := strings.Split(metadataContent, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "metadata {" || line == "}" {
			continue
		}

		// Look for key = value pattern
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				// Remove quotes if present
				if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
					value = value[1 : len(value)-1]
				}

				metadata[key] = value
			}
		}
	}

	return metadata, nil
}

// parseMetadataBlock parses a metadata block into a map
func (m *Manager) parseMetadataBlock(block *hcl.Block) (map[string]interface{}, error) {
	// Define the schema for metadata block attributes
	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "schema_version", Required: false},
			{Name: "schema_type", Required: false},
			{Name: "schema_name", Required: false},
			{Name: "last_updated", Required: false},
			{Name: "compatibility", Required: false},
			{Name: "description", Required: false},
			{Name: "scalver_format", Required: false},
		},
	})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse metadata block attributes: %s", diags.Error())
	}

	metadata := make(map[string]interface{})

	// Parse each attribute
	for name, attr := range content.Attributes {
		val, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse metadata attribute %s: %s", name, diags.Error())
		}

		// Convert cty.Value to appropriate Go type
		switch val.Type() {
		case cty.String:
			metadata[name] = val.AsString()
		case cty.Number:
			if val.IsNull() {
				metadata[name] = nil
			} else {
				// Try to convert to int first, then float
				if val.CanIterateElements() {
					// Handle arrays
					var result []string
					val.ForEachElement(func(_, value cty.Value) bool {
						if value.Type() == cty.String {
							result = append(result, value.AsString())
						}
						return false
					})
					metadata[name] = result
				} else {
					// Handle numbers
					bigFloat := val.AsBigFloat()
					if bigFloat.IsInt() {
						intVal, _ := bigFloat.Int64()
						metadata[name] = int(intVal)
					} else {
						floatVal, _ := bigFloat.Float64()
						metadata[name] = floatVal
					}
				}
			}
		case cty.Bool:
			metadata[name] = val.True()
		default:
			// For complex types, try to convert to string representation
			metadata[name] = val.AsString()
		}
	}

	return metadata, nil
}

// calculateComplexityScore calculates a complexity score
func (m *Manager) calculateComplexityScore(complexity *spookytypesschemas.SchemaComplexity) float64 {
	// Simple scoring algorithm
	score := float64(complexity.FieldCount) * 0.1
	score += float64(complexity.RuleCount) * 0.2
	score += float64(complexity.MaxDepth) * 0.3
	score += float64(complexity.CyclomaticComplexity) * 0.1

	// Normalize to 0-100 range
	if score > 100 {
		score = 100
	}

	return score
}

// calculateCoverage calculates schema coverage
func (m *Manager) calculateCoverage(schema *spookytypesschemas.Schema) *spookytypesschemas.SchemaCoverage {
	coverage := &spookytypesschemas.SchemaCoverage{}

	if schema.Validation != nil {
		// Calculate field coverage
		if schema.Validation.Fields != nil {
			coverage.FieldCoverage = 100.0 // Assuming all fields are covered
		}

		// Calculate validation coverage
		if len(schema.Validation.Rules) > 0 {
			coverage.ValidationCoverage = 100.0 // Assuming all validations are covered
		}
	}

	// Calculate documentation coverage
	if schema.Description != "" {
		coverage.DocumentationCoverage = 100.0
	}

	// Calculate example coverage
	if schema.Validation != nil && schema.Validation.Fields != nil {
		totalFields := len(schema.Validation.Fields)
		fieldsWithExamples := 0
		for _, field := range schema.Validation.Fields {
			if len(field.Examples) > 0 {
				fieldsWithExamples++
			}
		}
		if totalFields > 0 {
			coverage.ExampleCoverage = float64(fieldsWithExamples) / float64(totalFields) * 100.0
		}
	}

	return coverage
}

// calculateQuality calculates schema quality
func (m *Manager) calculateQuality(schema *spookytypesschemas.Schema) *spookytypesschemas.SchemaQuality {
	quality := &spookytypesschemas.SchemaQuality{
		Dimensions: make(map[string]float64),
	}

	// Calculate various quality dimensions
	quality.Dimensions["completeness"] = m.calculateCompleteness(schema)
	quality.Dimensions["consistency"] = m.calculateConsistency(schema)
	quality.Dimensions["clarity"] = m.calculateClarity(schema)

	// Calculate overall score
	totalScore := 0.0
	for _, score := range quality.Dimensions {
		totalScore += score
	}
	quality.Score = totalScore / float64(len(quality.Dimensions))

	return quality
}

// calculateCompleteness calculates schema completeness
func (m *Manager) calculateCompleteness(schema *spookytypesschemas.Schema) float64 {
	score := 0.0

	if schema.Name != "" {
		score += 20
	}
	if schema.Type != "" {
		score += 20
	}
	if schema.Version != "" {
		score += 20
	}
	if schema.Description != "" {
		score += 20
	}
	if schema.Validation != nil {
		score += 20
	}

	return score
}

// calculateConsistency calculates schema consistency
func (m *Manager) calculateConsistency(_ *spookytypesschemas.Schema) float64 {
	// Simplified consistency calculation
	return 80.0 // Assume good consistency
}

// calculateClarity calculates schema clarity
func (m *Manager) calculateClarity(schema *spookytypesschemas.Schema) float64 {
	score := 0.0

	if schema.Description != "" {
		score += 40
	}

	if schema.Validation != nil && schema.Validation.Fields != nil {
		for _, field := range schema.Validation.Fields {
			if field.Description != "" {
				score += 30
				break
			}
		}
	}

	if len(schema.Metadata) > 0 {
		score += 30
	}

	return score
}

// generateAnalysisRecommendations generates recommendations based on analysis
func (m *Manager) generateAnalysisRecommendations(analysis *spookytypesschemas.SchemaAnalysis) []string {
	var recommendations []string

	// Complexity recommendations
	if analysis.Complexity.Score > 80 {
		recommendations = append(recommendations, "Consider simplifying the schema to reduce complexity")
	}

	// Coverage recommendations
	if analysis.Coverage.DocumentationCoverage < 50 {
		recommendations = append(recommendations, "Add more documentation to improve clarity")
	}

	if analysis.Coverage.ExampleCoverage < 30 {
		recommendations = append(recommendations, "Add examples for fields to improve usability")
	}

	// Quality recommendations
	if analysis.Quality.Score < 70 {
		recommendations = append(recommendations, "Review schema quality and address identified issues")
	}

	return recommendations
}

// compareFields compares fields between two schemas
func (m *Manager) compareFields(fields1, fields2 map[string]*spookytypesschemas.FieldValidation) (added, removed, modified []string) {
	// Find added fields
	for fieldName := range fields2 {
		if _, exists := fields1[fieldName]; !exists {
			added = append(added, fieldName)
		}
	}

	// Find removed fields
	for fieldName := range fields1 {
		if _, exists := fields2[fieldName]; !exists {
			removed = append(removed, fieldName)
		}
	}

	// Find modified fields
	for fieldName, field1 := range fields1 {
		if field2, exists := fields2[fieldName]; exists {
			if m.fieldValidationChanged(field1, field2) {
				modified = append(modified, fieldName)
			}
		}
	}

	return added, removed, modified
}

// fieldValidationChanged checks if field validation changed
func (m *Manager) fieldValidationChanged(field1, field2 *spookytypesschemas.FieldValidation) bool {
	if field1.Required != field2.Required {
		return true
	}
	if field1.Type != field2.Type {
		return true
	}
	// Add more detailed comparison as needed
	return false
}

// generateBasicExample generates a basic example
func (m *Manager) generateBasicExample(schema *spookytypesschemas.Schema) string {
	var example strings.Builder

	example.WriteString(fmt.Sprintf("# Example: %s\n\n", schema.Name))
	example.WriteString("```hcl\n")

	if schema.Validation != nil && schema.Validation.Fields != nil {
		for fieldName, field := range schema.Validation.Fields {
			example.WriteString(fmt.Sprintf("%s = ", fieldName))

			switch field.Type {
			case "string":
				if field.Default != nil {
					example.WriteString(fmt.Sprintf(`"%v"`, field.Default))
				} else {
					example.WriteString(`"example_value"`)
				}
			case "integer":
				if field.Default != nil {
					example.WriteString(fmt.Sprintf("%v", field.Default))
				} else {
					example.WriteString("42")
				}
			case "boolean":
				if field.Default != nil {
					example.WriteString(fmt.Sprintf("%v", field.Default))
				} else {
					example.WriteString("true")
				}
			default:
				example.WriteString(`"example_value"`)
			}
			example.WriteString("\n")
		}
	}

	example.WriteString("```\n")
	return example.String()
}

// generateFieldExample generates an example for a specific field
func (m *Manager) generateFieldExample(fieldName string, field *spookytypesschemas.FieldValidation) string {
	var example strings.Builder

	example.WriteString(fmt.Sprintf("# Example: %s field\n\n", fieldName))
	example.WriteString("```hcl\n")

	for i, fieldExample := range field.Examples {
		example.WriteString(fmt.Sprintf("# Example %d\n", i+1))
		example.WriteString(fmt.Sprintf("%s = %v\n\n", fieldName, fieldExample))
	}

	example.WriteString("```\n")
	return example.String()
}
