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
	helpers  *SchemaHelpers

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
		helpers:         NewSchemaHelpers(logger),
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
	return e.helpers.CheckCompatibility(schema1, schema2)
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
	return e.helpers.ValidateMigration(schema, migration)
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
	return e.helpers.CheckForUpdates(schema, e.registry.List())
}

// Helper methods - using shared helpers to avoid duplication
