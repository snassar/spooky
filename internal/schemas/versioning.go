package schemas

import (
	"fmt"
	"strings"
	"time"
)

// SchemaVersion represents a schema version
type SchemaVersion struct {
	Version      string     `json:"version"`
	ReleasedAt   time.Time  `json:"released_at"`
	Deprecated   bool       `json:"deprecated"`
	DeprecatedAt *time.Time `json:"deprecated_at,omitempty"`
	Breaking     bool       `json:"breaking"`
	Description  string     `json:"description"`
}

// SchemaMigration represents a schema migration
type SchemaMigration struct {
	FromVersion string                 `json:"from_version"`
	ToVersion   string                 `json:"to_version"`
	Type        MigrationType          `json:"type"`
	Script      string                 `json:"script,omitempty"`
	Rules       map[string]interface{} `json:"rules,omitempty"`
	Description string                 `json:"description"`
}

// MigrationType represents the type of migration
type MigrationType string

const (
	MigrationTypeAutomatic MigrationType = "automatic"
	MigrationTypeManual    MigrationType = "manual"
	MigrationTypeBreaking  MigrationType = "breaking"
)

// SchemaVersionManager manages schema versions and migrations
type SchemaVersionManager struct {
	versions   map[string]SchemaVersion
	migrations map[string]SchemaMigration
	current    string
}

// NewSchemaVersionManager creates a new schema version manager
func NewSchemaVersionManager() *SchemaVersionManager {
	svm := &SchemaVersionManager{
		versions:   make(map[string]SchemaVersion),
		migrations: make(map[string]SchemaMigration),
		current:    "1.0.0",
	}

	// Initialize with default versions
	svm.initializeDefaultVersions()
	svm.initializeDefaultMigrations()

	return svm
}

// initializeDefaultVersions initializes the default schema versions
func (svm *SchemaVersionManager) initializeDefaultVersions() {
	// Project schemas
	svm.versions["project"] = SchemaVersion{
		Version:     "1.0.0",
		ReleasedAt:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Deprecated:  false,
		Breaking:    false,
		Description: "Initial project schema version",
	}

	svm.versions["machines"] = SchemaVersion{
		Version:     "1.0.0",
		ReleasedAt:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Deprecated:  false,
		Breaking:    false,
		Description: "Initial machines schema version",
	}

	svm.versions["actions"] = SchemaVersion{
		Version:     "1.0.0",
		ReleasedAt:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Deprecated:  false,
		Breaking:    false,
		Description: "Initial actions schema version",
	}

	svm.versions["variables"] = SchemaVersion{
		Version:     "1.0.0",
		ReleasedAt:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Deprecated:  false,
		Breaking:    false,
		Description: "Initial variables schema version",
	}

	// Facts schemas
	svm.versions["facts"] = SchemaVersion{
		Version:     "1.0.0",
		ReleasedAt:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Deprecated:  false,
		Breaking:    false,
		Description: "Initial facts schema version",
	}

	// Template schemas
	svm.versions["templates"] = SchemaVersion{
		Version:     "1.0.0",
		ReleasedAt:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Deprecated:  false,
		Breaking:    false,
		Description: "Initial templates schema version",
	}

	// Configuration schemas
	svm.versions["spooky"] = SchemaVersion{
		Version:     "1.0.0",
		ReleasedAt:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Deprecated:  false,
		Breaking:    false,
		Description: "Initial spooky configuration schema version",
	}
}

// initializeDefaultMigrations initializes the default schema migrations
func (svm *SchemaVersionManager) initializeDefaultMigrations() {
	// No migrations for initial version
	// This will be populated as schemas evolve
}

// GetSchemaVersion returns the version information for a schema
func (svm *SchemaVersionManager) GetSchemaVersion(schemaName string) (SchemaVersion, error) {
	version, exists := svm.versions[schemaName]
	if !exists {
		return SchemaVersion{}, fmt.Errorf("schema version not found: %s", schemaName)
	}
	return version, nil
}

// GetCurrentVersion returns the current schema version
func (svm *SchemaVersionManager) GetCurrentVersion() string {
	return svm.current
}

// IsSchemaDeprecated checks if a schema is deprecated
func (svm *SchemaVersionManager) IsSchemaDeprecated(schemaName string) (bool, error) {
	version, err := svm.GetSchemaVersion(schemaName)
	if err != nil {
		return false, err
	}
	return version.Deprecated, nil
}

// IsSchemaBreaking checks if a schema has breaking changes
func (svm *SchemaVersionManager) IsSchemaBreaking(schemaName string) (bool, error) {
	version, err := svm.GetSchemaVersion(schemaName)
	if err != nil {
		return false, err
	}
	return version.Breaking, nil
}

// GetMigrationPath returns the migration path from one version to another
func (svm *SchemaVersionManager) GetMigrationPath(schemaName, fromVersion, toVersion string) ([]SchemaMigration, error) {
	var migrations []SchemaMigration

	// First, check for explicit migrations
	migrationKey := fmt.Sprintf("%s:%s:%s", schemaName, fromVersion, toVersion)
	if migration, exists := svm.migrations[migrationKey]; exists {
		migrations = append(migrations, migration)
		return migrations, nil
	}

	// If no explicit migration exists, generate automatic migration path
	// This is a simplified version - in a real implementation, this would be more sophisticated
	if fromVersion != toVersion {
		// Generate automatic migration
		automaticMigration := SchemaMigration{
			FromVersion: fromVersion,
			ToVersion:   toVersion,
			Type:        MigrationTypeAutomatic,
			Description: fmt.Sprintf("Automatic migration from %s to %s", fromVersion, toVersion),
			Rules: map[string]interface{}{
				"auto_generated": true,
				"schema_name":    schemaName,
			},
		}
		migrations = append(migrations, automaticMigration)
	}

	return migrations, nil
}

// ValidateSchemaCompatibility validates if a schema is compatible with the current version
func (svm *SchemaVersionManager) ValidateSchemaCompatibility(schemaName, version string) error {
	schemaVersion, err := svm.GetSchemaVersion(schemaName)
	if err != nil {
		return fmt.Errorf("schema not found: %s", err)
	}

	// Check if the schema is deprecated
	if schemaVersion.Deprecated {
		deprecatedAt := schemaVersion.DeprecatedAt
		if deprecatedAt != nil && time.Now().After(*deprecatedAt) {
			return fmt.Errorf("schema %s version %s is deprecated since %s", schemaName, version, deprecatedAt.Format("2006-01-02"))
		}
	}

	// Check if the schema has breaking changes
	if schemaVersion.Breaking {
		return fmt.Errorf("schema %s version %s contains breaking changes", schemaName, version)
	}

	return nil
}

// AddSchemaVersion adds a new schema version
func (svm *SchemaVersionManager) AddSchemaVersion(schemaName string, version SchemaVersion) error {
	svm.versions[schemaName] = version
	return nil
}

// AddMigration adds a new schema migration
func (svm *SchemaVersionManager) AddMigration(schemaName, fromVersion, toVersion string, migration SchemaMigration) error {
	migrationKey := fmt.Sprintf("%s:%s:%s", schemaName, fromVersion, toVersion)
	svm.migrations[migrationKey] = migration
	return nil
}

// ListSchemaVersions returns all schema versions
func (svm *SchemaVersionManager) ListSchemaVersions() map[string]SchemaVersion {
	return svm.versions
}

// ListMigrations returns all migrations for a schema
func (svm *SchemaVersionManager) ListMigrations(schemaName string) []SchemaMigration {
	var migrations []SchemaMigration
	for key, migration := range svm.migrations {
		if strings.HasPrefix(key, schemaName+":") {
			migrations = append(migrations, migration)
		}
	}
	return migrations
}

// UpdateCurrentVersion updates the current schema version
func (svm *SchemaVersionManager) UpdateCurrentVersion(version string) {
	svm.current = version
}

// GetSchemaVersionHistory returns the version history for a schema
func (svm *SchemaVersionManager) GetSchemaVersionHistory(schemaName string) ([]SchemaVersion, error) {
	// For now, return just the current version
	// In a real implementation, this would return the full history
	version, err := svm.GetSchemaVersion(schemaName)
	if err != nil {
		return nil, err
	}
	return []SchemaVersion{version}, nil
}

// ValidateMigration validates if a migration is valid
func (svm *SchemaVersionManager) ValidateMigration(migration SchemaMigration) error {
	if migration.FromVersion == "" {
		return fmt.Errorf("migration must specify from_version")
	}
	if migration.ToVersion == "" {
		return fmt.Errorf("migration must specify to_version")
	}
	if migration.Type == "" {
		return fmt.Errorf("migration must specify type")
	}
	if migration.Description == "" {
		return fmt.Errorf("migration must specify description")
	}
	return nil
}

// GetDeprecatedSchemas returns all deprecated schemas
func (svm *SchemaVersionManager) GetDeprecatedSchemas() []string {
	var deprecated []string
	for name, version := range svm.versions {
		if version.Deprecated {
			deprecated = append(deprecated, name)
		}
	}
	return deprecated
}

// GetBreakingSchemas returns all schemas with breaking changes
func (svm *SchemaVersionManager) GetBreakingSchemas() []string {
	var breaking []string
	for name, version := range svm.versions {
		if version.Breaking {
			breaking = append(breaking, name)
		}
	}
	return breaking
}

// SuggestMigration suggests migration for a schema
func (svm *SchemaVersionManager) SuggestMigration(schemaName, currentVersion string) ([]SchemaMigration, error) {
	// Get the target version (current schema version)
	targetVersion := svm.GetCurrentVersion()

	// If already at target version, no migration needed
	if currentVersion == targetVersion {
		return []SchemaMigration{}, nil
	}

	// Get migration path
	migrations, err := svm.GetMigrationPath(schemaName, currentVersion, targetVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get migration path: %w", err)
	}

	return migrations, nil
}

// GetMigrationSuggestions returns migration suggestions for all schemas
func (svm *SchemaVersionManager) GetMigrationSuggestions() map[string][]SchemaMigration {
	suggestions := make(map[string][]SchemaMigration)

	// Check all schemas for migration suggestions
	for schemaName, version := range svm.versions {
		if migrations, err := svm.SuggestMigration(schemaName, version.Version); err == nil && len(migrations) > 0 {
			suggestions[schemaName] = migrations
		}
	}

	return suggestions
}
