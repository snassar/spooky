// Package schemas provides schema evolution and migration functionality for the spooky codebase.
package schemas

import (
	"fmt"
	"sort"
	"strings"
	"time"

	spookytypeslogging "spooky/internal/types/logging"
	spookytypesschemas "spooky/internal/types/schemas"
)

// EvolutionManager provides comprehensive schema evolution and migration functionality
type EvolutionManager struct {
	logger   spookytypeslogging.Logger
	registry spookytypesschemas.SchemaRegistry

	// Evolution tracking
	versionHistory  map[string][]spookytypesschemas.SchemaVersion
	migrationPaths  map[string][]spookytypesschemas.SchemaMigration
	breakingChanges map[string][]spookytypesschemas.BreakingChange
	deprecations    map[string][]spookytypesschemas.SchemaDeprecation
}

// NewEvolutionManager creates a new evolution manager instance
func NewEvolutionManager(logger spookytypeslogging.Logger, registry spookytypesschemas.SchemaRegistry) *EvolutionManager {
	return &EvolutionManager{
		logger:          logger,
		registry:        registry,
		versionHistory:  make(map[string][]spookytypesschemas.SchemaVersion),
		migrationPaths:  make(map[string][]spookytypesschemas.SchemaMigration),
		breakingChanges: make(map[string][]spookytypesschemas.BreakingChange),
		deprecations:    make(map[string][]spookytypesschemas.SchemaDeprecation),
	}
}

// TrackSchemaEvolution tracks schema evolution for a schema
func (e *EvolutionManager) TrackSchemaEvolution(schema *spookytypesschemas.Schema) error {
	if schema == nil {
		return fmt.Errorf("schema cannot be nil")
	}

	schemaKey := fmt.Sprintf("%s:%s", schema.Type, schema.Name)

	// Track version history
	if schema.Evolution != nil && len(schema.Evolution.Versions) > 0 {
		e.versionHistory[schemaKey] = schema.Evolution.Versions
	}

	// Track migration paths
	if schema.Evolution != nil && len(schema.Evolution.Migrations) > 0 {
		e.migrationPaths[schemaKey] = schema.Evolution.Migrations
	}

	// Track breaking changes
	if schema.Evolution != nil && len(schema.Evolution.BreakingChanges) > 0 {
		e.breakingChanges[schemaKey] = schema.Evolution.BreakingChanges
	}

	// Track deprecations
	if schema.Evolution != nil && len(schema.Evolution.Deprecations) > 0 {
		e.deprecations[schemaKey] = schema.Evolution.Deprecations
	}

	e.logger.Debug("Schema evolution tracked", map[string]interface{}{
		"schema_key":   schemaKey,
		"versions":     len(schema.Evolution.Versions),
		"migrations":   len(schema.Evolution.Migrations),
		"breaking":     len(schema.Evolution.BreakingChanges),
		"deprecations": len(schema.Evolution.Deprecations),
	})

	return nil
}

// CheckCompatibility checks compatibility between two schema versions
func (e *EvolutionManager) CheckCompatibility(schema1, schema2 *spookytypesschemas.Schema) (*spookytypesschemas.CompatibilityResult, error) {
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
	e.checkFieldCompatibility(schema1, schema2, result)

	// Check validation rule compatibility
	e.checkValidationCompatibility(schema1, schema2, result)

	// Check evolution compatibility
	e.checkEvolutionCompatibility(schema1, schema2, result)

	// Determine if migration is required
	result.MigrationRequired = len(result.Issues) > 0

	return result, nil
}

// GetMigrationPath gets the migration path between two schema versions
func (e *EvolutionManager) GetMigrationPath(fromSchema, toSchema *spookytypesschemas.Schema) ([]spookytypesschemas.SchemaMigration, error) {
	if fromSchema == nil || toSchema == nil {
		return nil, fmt.Errorf("both schemas must be provided")
	}

	schemaKey := fmt.Sprintf("%s:%s", fromSchema.Type, fromSchema.Name)
	migrations, exists := e.migrationPaths[schemaKey]
	if !exists {
		return nil, fmt.Errorf("no migration paths found for schema %s", schemaKey)
	}

	var path []spookytypesschemas.SchemaMigration
	for _, migration := range migrations {
		if migration.FromVersion == fromSchema.Version && migration.ToVersion == toSchema.Version {
			path = append(path, migration)
		}
	}

	if len(path) == 0 {
		return nil, fmt.Errorf("no migration path found from version %s to %s", fromSchema.Version, toSchema.Version)
	}

	// Sort migrations by priority or order
	sort.Slice(path, func(i, j int) bool {
		return path[i].Name < path[j].Name
	})

	return path, nil
}

// ValidateMigration validates a migration before applying it
func (e *EvolutionManager) ValidateMigration(schema *spookytypesschemas.Schema, migration *spookytypesschemas.SchemaMigration) (*spookytypesschemas.ValidationResult, error) {
	if schema == nil || migration == nil {
		return nil, fmt.Errorf("schema and migration cannot be nil")
	}

	result := &spookytypesschemas.ValidationResult{
		Valid:       true,
		ValidatedAt: time.Now(),
		Errors:      []spookytypesschemas.SchemaError{},
		Warnings:    []spookytypesschemas.SchemaError{},
		Info:        []spookytypesschemas.SchemaError{},
	}

	// Validate migration version compatibility
	if migration.FromVersion != schema.Version {
		error := spookytypesschemas.NewSchemaError(schema.Name, schema.Type,
			fmt.Sprintf("Migration from version %s cannot be applied to schema version %s",
				migration.FromVersion, schema.Version))
		error.Severity = "error"
		result.Errors = append(result.Errors, *error)
		result.Valid = false
	}

	// Validate migration prerequisites
	if migration.Validation != nil {
		for _, preValidation := range migration.Validation.PreValidation {
			e.logger.Debug("Executing pre-validation", map[string]interface{}{
				"validation": preValidation,
			})
		}
	}

	// Update valid flag based on errors
	if len(result.Errors) > 0 {
		result.Valid = false
	}

	return result, nil
}

// ApplyMigration applies a migration to a schema
func (e *EvolutionManager) ApplyMigration(schema *spookytypesschemas.Schema, migration *spookytypesschemas.SchemaMigration) error {
	if schema == nil || migration == nil {
		return fmt.Errorf("schema and migration cannot be nil")
	}

	// Validate migration before applying
	result, err := e.ValidateMigration(schema, migration)
	if err != nil {
		return fmt.Errorf("migration validation failed: %w", err)
	}

	if !result.Valid {
		return fmt.Errorf("migration validation failed: %d errors", len(result.Errors))
	}

	// Apply migration
	schema.Version = migration.ToVersion
	schema.UpdatedAt = time.Now()

	// Update schema metadata
	if schema.Metadata == nil {
		schema.Metadata = make(map[string]interface{})
	}
	schema.Metadata["last_migration"] = migration.Name
	schema.Metadata["migration_date"] = time.Now()

	e.logger.Info("Migration applied", map[string]interface{}{
		"schema_name":    schema.Name,
		"from_version":   migration.FromVersion,
		"to_version":     migration.ToVersion,
		"migration_name": migration.Name,
	})

	return nil
}

// GetBreakingChanges gets breaking changes for a schema
func (e *EvolutionManager) GetBreakingChanges(schema *spookytypesschemas.Schema) ([]spookytypesschemas.BreakingChange, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema cannot be nil")
	}

	schemaKey := fmt.Sprintf("%s:%s", schema.Type, schema.Name)
	changes, exists := e.breakingChanges[schemaKey]
	if !exists {
		return []spookytypesschemas.BreakingChange{}, nil
	}

	return changes, nil
}

// GetDeprecations gets deprecations for a schema
func (e *EvolutionManager) GetDeprecations(schema *spookytypesschemas.Schema) ([]spookytypesschemas.SchemaDeprecation, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema cannot be nil")
	}

	schemaKey := fmt.Sprintf("%s:%s", schema.Type, schema.Name)
	deprecations, exists := e.deprecations[schemaKey]
	if !exists {
		return []spookytypesschemas.SchemaDeprecation{}, nil
	}

	return deprecations, nil
}

// GenerateMigrationScript generates a migration script
func (e *EvolutionManager) GenerateMigrationScript(fromSchema, toSchema *spookytypesschemas.Schema) (string, error) {
	if fromSchema == nil || toSchema == nil {
		return "", fmt.Errorf("both schemas must be provided")
	}

	// Get migration path
	migrations, err := e.GetMigrationPath(fromSchema, toSchema)
	if err != nil {
		return "", err
	}

	var script strings.Builder
	script.WriteString(fmt.Sprintf("# Migration script from %s to %s\n", fromSchema.Version, toSchema.Version))
	script.WriteString(fmt.Sprintf("# Generated on %s\n\n", time.Now().Format(time.RFC3339)))

	for i, migration := range migrations {
		script.WriteString(fmt.Sprintf("# Step %d: %s\n", i+1, migration.Name))
		if migration.Description != "" {
			script.WriteString(fmt.Sprintf("# %s\n", migration.Description))
		}

		if migration.Script != "" {
			script.WriteString(migration.Script)
			script.WriteString("\n\n")
		}
	}

	return script.String(), nil
}

// CheckForUpdates checks for available schema updates
func (e *EvolutionManager) CheckForUpdates(schema *spookytypesschemas.Schema) ([]spookytypesschemas.SchemaUpdate, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema cannot be nil")
	}

	var updates []spookytypesschemas.SchemaUpdate

	// Check for newer versions
	for _, registeredSchema := range e.registry.List() {
		if registeredSchema.Type == schema.Type && registeredSchema.Name == schema.Name {
			if registeredSchema.Version != schema.Version {
				update := spookytypesschemas.SchemaUpdate{
					Type:        e.determineUpdateType(schema.Version, registeredSchema.Version),
					Description: fmt.Sprintf("New version available: %s", registeredSchema.Version),
					Version:     registeredSchema.Version,
					Date:        registeredSchema.UpdatedAt,
					Priority:    "medium",
				}
				updates = append(updates, update)
			}
		}
	}

	return updates, nil
}

// Helper methods

// checkFieldCompatibility checks field compatibility between schemas
func (e *EvolutionManager) checkFieldCompatibility(schema1, schema2 *spookytypesschemas.Schema, result *spookytypesschemas.CompatibilityResult) {
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
			if e.fieldValidationChanged(field1, field2) {
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
func (e *EvolutionManager) checkValidationCompatibility(schema1, schema2 *spookytypesschemas.Schema, result *spookytypesschemas.CompatibilityResult) {
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
func (e *EvolutionManager) checkEvolutionCompatibility(schema1, schema2 *spookytypesschemas.Schema, result *spookytypesschemas.CompatibilityResult) {
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
func (e *EvolutionManager) fieldValidationChanged(field1, field2 *spookytypesschemas.FieldValidation) bool {
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
		if e.constraintsChanged(field1.Constraints, field2.Constraints) {
			return true
		}
	}

	return false
}

// constraintsChanged checks if field constraints changed
func (e *EvolutionManager) constraintsChanged(constraints1, constraints2 *spookytypesschemas.FieldConstraints) bool {
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

// determineUpdateType determines the type of update
func (e *EvolutionManager) determineUpdateType(fromVersion, toVersion string) string {
	// This is a simplified implementation
	// In a real implementation, you would parse semantic versions
	if fromVersion == toVersion {
		return "none"
	}

	return "patch"
}
