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
	schemas map[string]*spookytypesschemas.Schema
	byType  map[string][]*spookytypesschemas.Schema
}

// NewRegistry creates a new schema registry instance
func NewRegistry(logger spookytypeslogging.Logger) *Registry {
	return &Registry{
		logger:  logger,
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
	if schema1 == nil || schema2 == nil {
		return nil, fmt.Errorf("both schemas must be provided")
	}

	result := &spookytypesschemas.CompatibilityResult{
		Compatible: true,
		Issues:     []spookytypesschemas.CompatibilityIssue{},
	}

	// Check basic compatibility
	if schema1.Type != schema2.Type {
		issue := spookytypesschemas.CompatibilityIssue{
			Type:        "breaking",
			Description: fmt.Sprintf("Schema types differ: %s vs %s", schema1.Type, schema2.Type),
			Severity:    "error",
			Resolution:  "Schemas must be of the same type to be compatible",
		}
		result.Issues = append(result.Issues, issue)
		result.Compatible = false
	}

	// Check version compatibility
	if schema1.Version != schema2.Version {
		issue := spookytypesschemas.CompatibilityIssue{
			Type:        "warning",
			Description: fmt.Sprintf("Schema versions differ: %s vs %s", schema1.Version, schema2.Version),
			Severity:    "warning",
			Resolution:  "Consider migrating to the same version for full compatibility",
		}
		result.Issues = append(result.Issues, issue)
	}

	// Check field compatibility
	r.checkFieldCompatibility(schema1, schema2, result)

	// Check validation rule compatibility
	r.checkValidationCompatibility(schema1, schema2, result)

	// Check evolution compatibility
	r.checkEvolutionCompatibility(schema1, schema2, result)

	// Determine if migration is required
	result.MigrationRequired = len(result.Issues) > 0

	return result, nil
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

// checkFieldCompatibility checks field compatibility between schemas
func (r *Registry) checkFieldCompatibility(schema1, schema2 *spookytypesschemas.Schema, result *spookytypesschemas.CompatibilityResult) {
	if schema1.Validation == nil || schema2.Validation == nil {
		return
	}

	fields1 := schema1.Validation.Fields
	fields2 := schema2.Validation.Fields

	if fields1 == nil || fields2 == nil {
		return
	}

	// Check for removed fields
	for fieldName := range fields1 {
		if _, exists := fields2[fieldName]; !exists {
			issue := spookytypesschemas.CompatibilityIssue{
				Type:        "breaking",
				Description: fmt.Sprintf("Field '%s' was removed in schema2", fieldName),
				Fields:      []string{fieldName},
				Severity:    "error",
				Resolution:  fmt.Sprintf("Field '%s' is no longer available", fieldName),
			}
			result.Issues = append(result.Issues, issue)
			result.Compatible = false
		}
	}

	// Check for added fields
	for fieldName := range fields2 {
		if _, exists := fields1[fieldName]; !exists {
			issue := spookytypesschemas.CompatibilityIssue{
				Type:        "info",
				Description: fmt.Sprintf("Field '%s' was added in schema2", fieldName),
				Fields:      []string{fieldName},
				Severity:    "info",
				Resolution:  fmt.Sprintf("Field '%s' is now available", fieldName),
			}
			result.Issues = append(result.Issues, issue)
		}
	}

	// Check for modified fields
	for fieldName, field1 := range fields1 {
		if field2, exists := fields2[fieldName]; exists {
			if r.fieldValidationChanged(field1, field2) {
				issue := spookytypesschemas.CompatibilityIssue{
					Type:        "warning",
					Description: fmt.Sprintf("Field '%s' validation rules changed", fieldName),
					Fields:      []string{fieldName},
					Severity:    "warning",
					Resolution:  fmt.Sprintf("Review validation rules for field '%s'", fieldName),
				}
				result.Issues = append(result.Issues, issue)
			}
		}
	}
}

// checkValidationCompatibility checks validation rule compatibility
func (r *Registry) checkValidationCompatibility(schema1, schema2 *spookytypesschemas.Schema, result *spookytypesschemas.CompatibilityResult) {
	if schema1.Validation == nil || schema2.Validation == nil {
		return
	}

	// Check validation mode changes
	if schema1.Validation.Mode != schema2.Validation.Mode {
		issue := spookytypesschemas.CompatibilityIssue{
			Type:        "warning",
			Description: fmt.Sprintf("Validation mode changed from %s to %s", schema1.Validation.Mode, schema2.Validation.Mode),
			Severity:    "warning",
			Resolution:  "Review validation behavior changes",
		}
		result.Issues = append(result.Issues, issue)
	}

	// Check validation rule changes
	rules1 := make(map[string]spookytypesschemas.ValidationRule)
	rules2 := make(map[string]spookytypesschemas.ValidationRule)

	for _, rule := range schema1.Validation.Rules {
		rules1[rule.Name] = rule
	}

	for _, rule := range schema2.Validation.Rules {
		rules2[rule.Name] = rule
	}

	// Check for removed rules
	for ruleName := range rules1 {
		if _, exists := rules2[ruleName]; !exists {
			issue := spookytypesschemas.CompatibilityIssue{
				Type:        "info",
				Description: fmt.Sprintf("Validation rule '%s' was removed", ruleName),
				Severity:    "info",
				Resolution:  fmt.Sprintf("Validation rule '%s' is no longer enforced", ruleName),
			}
			result.Issues = append(result.Issues, issue)
		}
	}

	// Check for added rules
	for ruleName := range rules2 {
		if _, exists := rules1[ruleName]; !exists {
			issue := spookytypesschemas.CompatibilityIssue{
				Type:        "warning",
				Description: fmt.Sprintf("Validation rule '%s' was added", ruleName),
				Severity:    "warning",
				Resolution:  fmt.Sprintf("New validation rule '%s' may affect existing data", ruleName),
			}
			result.Issues = append(result.Issues, issue)
		}
	}
}

// checkEvolutionCompatibility checks evolution compatibility
func (r *Registry) checkEvolutionCompatibility(schema1, schema2 *spookytypesschemas.Schema, result *spookytypesschemas.CompatibilityResult) {
	if schema1.Evolution == nil || schema2.Evolution == nil {
		return
	}

	// Check for breaking changes
	for _, breakingChange := range schema2.Evolution.BreakingChanges {
		issue := spookytypesschemas.CompatibilityIssue{
			Type:        "breaking",
			Description: breakingChange.Description,
			Fields:      []string{breakingChange.Field},
			Severity:    "error",
			Resolution:  breakingChange.Mitigation,
		}
		result.Issues = append(result.Issues, issue)
		result.Compatible = false
	}

	// Check for deprecations
	for _, deprecation := range schema2.Evolution.Deprecations {
		issue := spookytypesschemas.CompatibilityIssue{
			Type:        "deprecation",
			Description: fmt.Sprintf("Field '%s' is deprecated: %s", deprecation.Field, deprecation.Reason),
			Fields:      []string{deprecation.Field},
			Severity:    "warning",
			Resolution:  fmt.Sprintf("Use '%s' instead of '%s'", deprecation.Replacement, deprecation.Field),
		}
		result.Issues = append(result.Issues, issue)
	}
}

// fieldValidationChanged checks if field validation rules changed
func (r *Registry) fieldValidationChanged(field1, field2 *spookytypesschemas.FieldValidation) bool {
	// Check basic field properties
	if field1.Required != field2.Required {
		return true
	}

	if field1.Type != field2.Type {
		return true
	}

	// Check constraints
	if field1.Constraints == nil && field2.Constraints != nil {
		return true
	}

	if field1.Constraints != nil && field2.Constraints == nil {
		return true
	}

	if field1.Constraints != nil && field2.Constraints != nil {
		if r.constraintsChanged(field1.Constraints, field2.Constraints) {
			return true
		}
	}

	return false
}

// constraintsChanged checks if field constraints changed
func (r *Registry) constraintsChanged(constraints1, constraints2 *spookytypesschemas.FieldConstraints) bool {
	// Check string constraints
	if constraints1.MinLength != constraints2.MinLength {
		return true
	}

	if constraints1.MaxLength != constraints2.MaxLength {
		return true
	}

	if constraints1.Pattern != constraints2.Pattern {
		return true
	}

	if constraints1.Format != constraints2.Format {
		return true
	}

	// Check numeric constraints
	if constraints1.Min != constraints2.Min {
		return true
	}

	if constraints1.Max != constraints2.Max {
		return true
	}

	// Check array constraints
	if constraints1.MinItems != constraints2.MinItems {
		return true
	}

	if constraints1.MaxItems != constraints2.MaxItems {
		return true
	}

	// Check enum values
	if len(constraints1.Enum) != len(constraints2.Enum) {
		return true
	}

	for i, enumValue := range constraints1.Enum {
		if i >= len(constraints2.Enum) || enumValue != constraints2.Enum[i] {
			return true
		}
	}

	return false
}

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
