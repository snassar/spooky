// Package schemas provides schema validation and management functionality for the spooky codebase.
package schemas

import (
	"fmt"
	"sort"
	"strings"

	spookytypeslogging "spooky/internal/types/logging"
	spookytypesschemas "spooky/internal/types/schemas"
)

// Registry provides comprehensive schema registry functionality
type Registry struct {
	logger  spookytypeslogging.Logger
	helpers *SchemaHelpers
	schemas map[string]*spookytypesschemas.Schema
	byType  map[string][]*spookytypesschemas.Schema
}

// NewRegistry creates a new schema registry instance
func NewRegistry(logger spookytypeslogging.Logger) *Registry {
	return &Registry{
		logger:  logger,
		helpers: NewSchemaHelpers(logger),
		schemas: make(map[string]*spookytypesschemas.Schema),
		byType:  make(map[string][]*spookytypesschemas.Schema),
	}
}

// Register registers a new schema
func (r *Registry) Register(schema *spookytypesschemas.Schema) error {
	if schema == nil {
		return fmt.Errorf("schema cannot be nil")
	}

	if schema.Name == "" {
		return fmt.Errorf("schema name cannot be empty")
	}

	// Check if schema already exists
	if existing, exists := r.schemas[schema.Name]; exists {
		r.logger.Warn("Schema already registered, updating", map[string]interface{}{
			"schema_name": schema.Name,
			"old_version": existing.Version,
			"new_version": schema.Version,
		})
	}

	// Register schema
	r.schemas[schema.Name] = schema

	// Add to type index
	if schema.Type != "" {
		r.byType[schema.Type] = append(r.byType[schema.Type], schema)
	}

	r.logger.Debug("Schema registered", map[string]interface{}{
		"schema_name": schema.Name,
		"schema_type": schema.Type,
		"version":     schema.Version,
	})

	return nil
}

// Get returns a schema by name and type
func (r *Registry) Get(name, schemaType string) (*spookytypesschemas.Schema, bool) {
	schema, exists := r.schemas[name]
	if !exists {
		return nil, false
	}

	// Check type if specified
	if schemaType != "" && schema.Type != schemaType {
		return nil, false
	}

	return schema, true
}

// List returns all registered schemas
func (r *Registry) List() []*spookytypesschemas.Schema {
	var schemas []*spookytypesschemas.Schema
	for _, schema := range r.schemas {
		schemas = append(schemas, schema)
	}

	// Sort by name for consistent ordering
	sort.Slice(schemas, func(i, j int) bool {
		return schemas[i].Name < schemas[j].Name
	})

	return schemas
}

// ListByType returns schemas by type
func (r *Registry) ListByType(schemaType string) []*spookytypesschemas.Schema {
	schemas, exists := r.byType[schemaType]
	if !exists {
		return []*spookytypesschemas.Schema{}
	}

	// Sort by name for consistent ordering
	sort.Slice(schemas, func(i, j int) bool {
		return schemas[i].Name < schemas[j].Name
	})

	return schemas
}

// ValidateData validates data against a schema by name and type
func (r *Registry) ValidateData(schemaName, schemaType string, data interface{}) (*spookytypesschemas.ValidationResult, error) {
	schema, exists := r.Get(schemaName, schemaType)
	if !exists {
		return nil, fmt.Errorf("schema '%s' of type '%s' not found", schemaName, schemaType)
	}

	// Create validator and validate
	validator := NewValidator(r.logger)
	return validator.Validate(schema, data)
}

// CheckCompatibility checks compatibility between schemas
func (r *Registry) CheckCompatibility(schema1, schema2 *spookytypesschemas.Schema) (*spookytypesschemas.CompatibilityResult, error) {
	return r.helpers.CheckCompatibility(schema1, schema2)
}

// GetMigrationPath gets migration path between schema versions
func (r *Registry) GetMigrationPath(fromVersion, toVersion string) ([]*spookytypesschemas.SchemaMigration, error) {
	// Find schemas with the specified versions
	var fromSchemas, toSchemas []*spookytypesschemas.Schema

	for _, schema := range r.schemas {
		if schema.Version == fromVersion {
			fromSchemas = append(fromSchemas, schema)
		}
		if schema.Version == toVersion {
			toSchemas = append(toSchemas, schema)
		}
	}

	if len(fromSchemas) == 0 {
		return nil, fmt.Errorf("no schemas found with version %s", fromVersion)
	}

	if len(toSchemas) == 0 {
		return nil, fmt.Errorf("no schemas found with version %s", toVersion)
	}

	// Find migration path
	var migrations []*spookytypesschemas.SchemaMigration

	for _, fromSchema := range fromSchemas {
		if fromSchema.Evolution != nil {
			for _, migration := range fromSchema.Evolution.Migrations {
				if migration.FromVersion == fromVersion && migration.ToVersion == toVersion {
					migrations = append(migrations, &migration)
				}
			}
		}
	}

	if len(migrations) == 0 {
		return nil, fmt.Errorf("no migration path found from version %s to %s", fromVersion, toVersion)
	}

	return migrations, nil
}

// Helper methods - using shared helpers to avoid duplication

// GetSchemaInfo returns detailed information about a schema
func (r *Registry) GetSchemaInfo(schemaName string) (map[string]interface{}, error) {
	schema, exists := r.Get(schemaName, "")
	if !exists {
		return nil, fmt.Errorf("schema '%s' not found", schemaName)
	}

	info := map[string]interface{}{
		"name":        schema.Name,
		"type":        schema.Type,
		"version":     schema.Version,
		"description": schema.Description,
		"created_at":  schema.CreatedAt,
		"updated_at":  schema.UpdatedAt,
	}

	// Add validation information
	if schema.Validation != nil {
		validationInfo := map[string]interface{}{
			"enabled": schema.Validation.Enabled,
			"mode":    schema.Validation.Mode,
		}

		if schema.Validation.Fields != nil {
			validationInfo["field_count"] = len(schema.Validation.Fields)
		}

		if schema.Validation.Rules != nil {
			validationInfo["rule_count"] = len(schema.Validation.Rules)
		}

		info["validation"] = validationInfo
	}

	// Add evolution information
	if schema.Evolution != nil {
		evolutionInfo := map[string]interface{}{}

		if schema.Evolution.Versions != nil {
			evolutionInfo["version_count"] = len(schema.Evolution.Versions)
		}

		if schema.Evolution.Migrations != nil {
			evolutionInfo["migration_count"] = len(schema.Evolution.Migrations)
		}

		if schema.Evolution.Deprecations != nil {
			evolutionInfo["deprecation_count"] = len(schema.Evolution.Deprecations)
		}

		if schema.Evolution.BreakingChanges != nil {
			evolutionInfo["breaking_change_count"] = len(schema.Evolution.BreakingChanges)
		}

		info["evolution"] = evolutionInfo
	}

	// Add compatibility information
	if schema.Compatibility != nil {
		compatibilityInfo := map[string]interface{}{
			"forward_compatible":  schema.Compatibility.ForwardCompatible,
			"backward_compatible": schema.Compatibility.BackwardCompatible,
		}

		if schema.Compatibility.CompatibleTypes != nil {
			compatibilityInfo["compatible_types"] = schema.Compatibility.CompatibleTypes
		}

		if schema.Compatibility.CompatibleVersions != nil {
			compatibilityInfo["compatible_versions"] = schema.Compatibility.CompatibleVersions
		}

		info["compatibility"] = compatibilityInfo
	}

	return info, nil
}

// SearchSchemas searches for schemas by various criteria
func (r *Registry) SearchSchemas(criteria map[string]interface{}) ([]*spookytypesschemas.Schema, error) {
	var results []*spookytypesschemas.Schema

	for _, schema := range r.schemas {
		if r.matchesCriteria(schema, criteria) {
			results = append(results, schema)
		}
	}

	// Sort results by name
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results, nil
}

// matchesCriteria checks if a schema matches the search criteria
func (r *Registry) matchesCriteria(schema *spookytypesschemas.Schema, criteria map[string]interface{}) bool {
	for key, value := range criteria {
		switch key {
		case "name":
			if name, ok := value.(string); ok {
				if !strings.Contains(strings.ToLower(schema.Name), strings.ToLower(name)) {
					return false
				}
			}
		case "type":
			if schemaType, ok := value.(string); ok {
				if schema.Type != schemaType {
					return false
				}
			}
		case "version":
			if version, ok := value.(string); ok {
				if schema.Version != version {
					return false
				}
			}
		case "description":
			if description, ok := value.(string); ok {
				if !strings.Contains(strings.ToLower(schema.Description), strings.ToLower(description)) {
					return false
				}
			}
		}
	}

	return true
}

// GetSchemaStatistics returns statistics about registered schemas
func (r *Registry) GetSchemaStatistics() map[string]interface{} {
	stats := map[string]interface{}{
		"total_schemas": len(r.schemas),
		"by_type":       make(map[string]int),
		"by_version":    make(map[string]int),
	}

	// Count by type
	for schemaType, schemas := range r.byType {
		stats["by_type"].(map[string]int)[schemaType] = len(schemas)
	}

	// Count by version
	for _, schema := range r.schemas {
		version := schema.Version
		if version == "" {
			version = "unknown"
		}
		stats["by_version"].(map[string]int)[version]++
	}

	return stats
}
