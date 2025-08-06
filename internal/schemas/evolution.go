package schemas

import (
	"fmt"
	"strings"
	"time"
)

// SchemaChange represents a change in a schema
type SchemaChange struct {
	Type        ChangeType `json:"type"`
	Field       string     `json:"field,omitempty"`
	OldValue    string     `json:"old_value,omitempty"`
	NewValue    string     `json:"new_value,omitempty"`
	Description string     `json:"description"`
	Breaking    bool       `json:"breaking"`
	Severity    string     `json:"severity"` // "low", "medium", "high", "critical"
}

// ChangeType represents the type of schema change
type ChangeType string

const (
	ChangeTypeAdded      ChangeType = "added"
	ChangeTypeRemoved    ChangeType = "removed"
	ChangeTypeModified   ChangeType = "modified"
	ChangeTypeRenamed    ChangeType = "renamed"
	ChangeTypeDeprecated ChangeType = "deprecated"
)

// SchemaEvolutionManager manages schema evolution and change detection
type SchemaEvolutionManager struct {
	versionManager *SchemaVersionManager
	changes        map[string][]SchemaChange
	notifications  []SchemaNotification
}

// SchemaNotification represents a schema update notification
type SchemaNotification struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	SchemaName string    `json:"schema_name"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
	Read       bool      `json:"read"`
	ActionURL  string    `json:"action_url,omitempty"`
	Severity   string    `json:"severity"`
}

// NewSchemaEvolutionManager creates a new schema evolution manager
func NewSchemaEvolutionManager() *SchemaEvolutionManager {
	return &SchemaEvolutionManager{
		versionManager: NewSchemaVersionManager(),
		changes:        make(map[string][]SchemaChange),
		notifications:  make([]SchemaNotification, 0),
	}
}

// DetectChanges detects changes between two schema versions
func (sem *SchemaEvolutionManager) DetectChanges(schemaName, oldVersion, newVersion string) ([]SchemaChange, error) {
	var changes []SchemaChange

	// Get the old and new schema content
	oldSchema, err := sem.getSchemaContent(schemaName, oldVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get old schema content: %w", err)
	}

	newSchema, err := sem.getSchemaContent(schemaName, newVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get new schema content: %w", err)
	}

	// Compare the schemas and detect changes
	changes = sem.compareSchemas(oldSchema, newSchema)

	// Store the changes
	sem.changes[schemaName] = changes

	return changes, nil
}

// compareSchemas compares two schema contents and returns detected changes
func (sem *SchemaEvolutionManager) compareSchemas(oldContent, newContent string) []SchemaChange {
	var changes []SchemaChange

	// Simple string-based comparison for now
	// In a real implementation, this would parse the HCL and compare structures

	oldLines := strings.Split(oldContent, "\n")
	newLines := strings.Split(newContent, "\n")

	// Detect added lines
	for i, line := range newLines {
		if i >= len(oldLines) || oldLines[i] != line {
			// This is a new or modified line
			if strings.Contains(line, "variable") || strings.Contains(line, "block") {
				changes = append(changes, SchemaChange{
					Type:        ChangeTypeAdded,
					Field:       fmt.Sprintf("line_%d", i+1),
					NewValue:    strings.TrimSpace(line),
					Description: fmt.Sprintf("Added new schema element at line %d", i+1),
					Breaking:    false,
					Severity:    "low",
				})
			}
		}
	}

	// Detect removed lines
	for i, line := range oldLines {
		if i >= len(newLines) || newLines[i] != line {
			// This line was removed or modified
			if strings.Contains(line, "variable") || strings.Contains(line, "block") {
				changes = append(changes, SchemaChange{
					Type:        ChangeTypeRemoved,
					Field:       fmt.Sprintf("line_%d", i+1),
					OldValue:    strings.TrimSpace(line),
					Description: fmt.Sprintf("Removed schema element at line %d", i+1),
					Breaking:    true,
					Severity:    "high",
				})
			}
		}
	}

	return changes
}

// getSchemaContent gets the content of a schema at a specific version
func (sem *SchemaEvolutionManager) getSchemaContent(schemaName, _ string) (string, error) {
	// For now, just get the current schema content
	// In a real implementation, this would retrieve the specific version
	schemaType := SchemaType(schemaName)
	content, err := GetSchema(schemaType)
	if err != nil {
		return "", fmt.Errorf("failed to get schema content: %w", err)
	}
	return string(content), nil
}

// CheckCompatibility checks if a schema is compatible with the current version
func (sem *SchemaEvolutionManager) CheckCompatibility(schemaName, version string) (bool, []SchemaChange, error) {
	// Get the current version
	currentVersion := sem.versionManager.GetCurrentVersion()

	// If the version is the same as current, it's compatible
	if version == currentVersion {
		return true, []SchemaChange{}, nil
	}

	// Detect changes between the specified version and current version
	changes, err := sem.DetectChanges(schemaName, version, currentVersion)
	if err != nil {
		return false, nil, fmt.Errorf("failed to detect changes: %w", err)
	}

	// If no changes detected, it's compatible
	if len(changes) == 0 {
		return true, changes, nil
	}

	// Check if any changes are breaking
	compatible := true
	for _, change := range changes {
		if change.Breaking {
			compatible = false
			break
		}
	}

	return compatible, changes, nil
}

// AutoMigrate automatically migrates a schema to the latest version
func (sem *SchemaEvolutionManager) AutoMigrate(schemaName, fromVersion string) (string, []SchemaChange, error) {
	// Get the current version
	toVersion := sem.versionManager.GetCurrentVersion()

	// Get the migration path
	migrations, err := sem.versionManager.GetMigrationPath(schemaName, fromVersion, toVersion)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get migration path: %w", err)
	}

	var allChanges []SchemaChange

	// Apply each migration
	for _, migration := range migrations {
		if migration.Type == MigrationTypeAutomatic {
			// Apply automatic migration
			changes, err := sem.applyMigration(schemaName, migration)
			if err != nil {
				return "", nil, fmt.Errorf("failed to apply migration: %w", err)
			}
			allChanges = append(allChanges, changes...)
		} else {
			// Manual migration required
			return "", nil, fmt.Errorf("manual migration required for %s: %s", schemaName, migration.Description)
		}
	}

	// Create notification about the migration
	sem.createNotification(SchemaNotification{
		Type:       "migration",
		SchemaName: schemaName,
		Message:    fmt.Sprintf("Schema %s automatically migrated from %s to %s", schemaName, fromVersion, toVersion),
		CreatedAt:  time.Now(),
		Severity:   "info",
	})

	return toVersion, allChanges, nil
}

// applyMigration applies a single migration
func (sem *SchemaEvolutionManager) applyMigration(_ string, _ SchemaMigration) ([]SchemaChange, error) {
	// For now, just return empty changes
	// In a real implementation, this would apply the migration rules
	return []SchemaChange{}, nil
}

// createNotification creates a new schema notification
func (sem *SchemaEvolutionManager) createNotification(notification SchemaNotification) {
	notification.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	notification.Read = false
	sem.notifications = append(sem.notifications, notification)
}

// GetNotifications returns all schema notifications
func (sem *SchemaEvolutionManager) GetNotifications() []SchemaNotification {
	return sem.notifications
}

// GetUnreadNotifications returns unread schema notifications
func (sem *SchemaEvolutionManager) GetUnreadNotifications() []SchemaNotification {
	var unread []SchemaNotification
	for _, notification := range sem.notifications {
		if !notification.Read {
			unread = append(unread, notification)
		}
	}
	return unread
}

// MarkNotificationRead marks a notification as read
func (sem *SchemaEvolutionManager) MarkNotificationRead(notificationID string) error {
	for i, notification := range sem.notifications {
		if notification.ID == notificationID {
			sem.notifications[i].Read = true
			return nil
		}
	}
	return fmt.Errorf("notification not found: %s", notificationID)
}

// ClearNotifications clears all notifications
func (sem *SchemaEvolutionManager) ClearNotifications() {
	sem.notifications = make([]SchemaNotification, 0)
}

// GetSchemaChanges returns all changes for a schema
func (sem *SchemaEvolutionManager) GetSchemaChanges(schemaName string) []SchemaChange {
	return sem.changes[schemaName]
}

// GetAllChanges returns all schema changes
func (sem *SchemaEvolutionManager) GetAllChanges() map[string][]SchemaChange {
	return sem.changes
}

// ValidateSchemaEvolution validates if schema evolution is valid
func (sem *SchemaEvolutionManager) ValidateSchemaEvolution(schemaName, fromVersion, toVersion string) error {
	// Check if the versions exist
	_, err := sem.versionManager.GetSchemaVersion(schemaName)
	if err != nil {
		return fmt.Errorf("schema not found: %s", err)
	}

	// Check if the migration path exists
	migrations, err := sem.versionManager.GetMigrationPath(schemaName, fromVersion, toVersion)
	if err != nil {
		return fmt.Errorf("failed to get migration path: %w", err)
	}

	if len(migrations) == 0 {
		return fmt.Errorf("no migration path found from %s to %s", fromVersion, toVersion)
	}

	// Validate each migration
	for _, migration := range migrations {
		if err := sem.versionManager.ValidateMigration(migration); err != nil {
			return fmt.Errorf("invalid migration: %w", err)
		}
	}

	return nil
}

// GetEvolutionSummary returns a summary of schema evolution
func (sem *SchemaEvolutionManager) GetEvolutionSummary() map[string]interface{} {
	summary := make(map[string]interface{})

	// Count total changes
	totalChanges := 0
	breakingChanges := 0
	for _, changes := range sem.changes {
		totalChanges += len(changes)
		for _, change := range changes {
			if change.Breaking {
				breakingChanges++
			}
		}
	}

	// Count notifications
	unreadNotifications := len(sem.GetUnreadNotifications())

	summary["total_changes"] = totalChanges
	summary["breaking_changes"] = breakingChanges
	summary["unread_notifications"] = unreadNotifications
	summary["current_version"] = sem.versionManager.GetCurrentVersion()
	summary["deprecated_schemas"] = sem.versionManager.GetDeprecatedSchemas()
	summary["breaking_schemas"] = sem.versionManager.GetBreakingSchemas()

	return summary
}

// SuggestMigration suggests a migration path for a schema
func (sem *SchemaEvolutionManager) SuggestMigration(schemaName, currentVersion string) ([]SchemaMigration, error) {
	// Get the current schema version
	toVersion := sem.versionManager.GetCurrentVersion()

	// Get the migration path
	migrations, err := sem.versionManager.GetMigrationPath(schemaName, currentVersion, toVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get migration path: %w", err)
	}

	// Create notification about the suggested migration
	sem.createNotification(SchemaNotification{
		Type:       "migration_suggestion",
		SchemaName: schemaName,
		Message:    fmt.Sprintf("Migration suggested for schema %s from %s to %s", schemaName, currentVersion, toVersion),
		CreatedAt:  time.Now(),
		Severity:   "medium",
	})

	return migrations, nil
}

// ManageSchemaEvolutionWorkflow manages the complete workflow for schema evolution
func (sem *SchemaEvolutionManager) ManageSchemaEvolutionWorkflow(schemaName, fromVersion, toVersion string) (map[string]interface{}, error) {
	workflow := make(map[string]interface{})

	// Step 1: Validate the evolution
	if err := sem.ValidateSchemaEvolution(schemaName, fromVersion, toVersion); err != nil {
		workflow["status"] = "failed"
		workflow["error"] = err.Error()
		return workflow, err
	}

	// Step 2: Detect changes
	changes, err := sem.DetectChanges(schemaName, fromVersion, toVersion)
	if err != nil {
		workflow["status"] = "failed"
		workflow["error"] = err.Error()
		return workflow, err
	}

	// Step 3: Check compatibility
	compatible, _, err := sem.CheckCompatibility(schemaName, fromVersion)
	if err != nil {
		workflow["status"] = "failed"
		workflow["error"] = err.Error()
		return workflow, err
	}

	// Step 4: Get migration path
	migrations, err := sem.versionManager.GetMigrationPath(schemaName, fromVersion, toVersion)
	if err != nil {
		workflow["status"] = "failed"
		workflow["error"] = err.Error()
		return workflow, err
	}

	// Step 5: Create workflow summary
	workflow["status"] = "ready"
	workflow["schema_name"] = schemaName
	workflow["from_version"] = fromVersion
	workflow["to_version"] = toVersion
	workflow["compatible"] = compatible
	workflow["changes"] = len(changes)
	workflow["migrations"] = len(migrations)
	workflow["breaking_changes"] = func() int {
		count := 0
		for _, change := range changes {
			if change.Breaking {
				count++
			}
		}
		return count
	}()

	// Step 6: Generate workflow steps
	var steps []map[string]interface{}

	// Step 6.1: Pre-migration validation
	steps = append(steps, map[string]interface{}{
		"step":        1,
		"action":      "validate_schema",
		"description": "Validate current schema structure",
		"status":      "pending",
	})

	// Step 6.2: Backup current schema
	steps = append(steps, map[string]interface{}{
		"step":        2,
		"action":      "backup_schema",
		"description": "Create backup of current schema",
		"status":      "pending",
	})

	// Step 6.3: Apply migrations
	for i, migration := range migrations {
		steps = append(steps, map[string]interface{}{
			"step":        3 + i,
			"action":      "apply_migration",
			"description": fmt.Sprintf("Apply migration: %s", migration.Description),
			"migration":   migration,
			"status":      "pending",
		})
	}

	// Step 6.4: Post-migration validation
	steps = append(steps, map[string]interface{}{
		"step":        3 + len(migrations),
		"action":      "validate_migrated_schema",
		"description": "Validate migrated schema structure",
		"status":      "pending",
	})

	// Step 6.5: Update schema version
	steps = append(steps, map[string]interface{}{
		"step":        4 + len(migrations),
		"action":      "update_version",
		"description": "Update schema version to target version",
		"status":      "pending",
	})

	workflow["steps"] = steps
	workflow["total_steps"] = len(steps)

	// Step 7: Create notification
	sem.createNotification(SchemaNotification{
		Type:       "evolution_workflow",
		SchemaName: schemaName,
		Message:    fmt.Sprintf("Schema evolution workflow created for %s: %s → %s", schemaName, fromVersion, toVersion),
		CreatedAt:  time.Now(),
		Severity:   "medium",
	})

	return workflow, nil
}

// ExecuteSchemaEvolutionWorkflow executes a schema evolution workflow
func (sem *SchemaEvolutionManager) ExecuteSchemaEvolutionWorkflow(workflow map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	result["workflow_id"] = fmt.Sprintf("workflow_%d", time.Now().UnixNano())
	result["started_at"] = time.Now()
	result["status"] = "in_progress"

	// Get workflow details
	schemaName := workflow["schema_name"].(string)
	fromVersion := workflow["from_version"].(string)
	toVersion := workflow["to_version"].(string)
	steps := workflow["steps"].([]map[string]interface{})

	var executedSteps []map[string]interface{}

	// Execute each step
	for i, step := range steps {
		step["started_at"] = time.Now()
		step["status"] = "executing"

		// Execute step based on action
		switch step["action"] {
		case "validate_schema":
			// Validate current schema
			if err := sem.ValidateSchemaEvolution(schemaName, fromVersion, fromVersion); err != nil {
				step["status"] = "failed"
				step["error"] = err.Error()
				step["completed_at"] = time.Now()
				executedSteps = append(executedSteps, step)
				result["status"] = "failed"
				result["failed_step"] = i + 1
				result["error"] = err.Error()
				result["executed_steps"] = executedSteps
				return result, err
			}

		case "backup_schema":
			// Create backup (simulated)
			step["backup_path"] = fmt.Sprintf("/tmp/schema_backup_%s_%s.hcl", schemaName, fromVersion)

		case "apply_migration":
			// Apply migration
			migration := step["migration"].(SchemaMigration)
			if _, err := sem.applyMigration(schemaName, migration); err != nil {
				step["status"] = "failed"
				step["error"] = err.Error()
				step["completed_at"] = time.Now()
				executedSteps = append(executedSteps, step)
				result["status"] = "failed"
				result["failed_step"] = i + 1
				result["error"] = err.Error()
				result["executed_steps"] = executedSteps
				return result, err
			}

		case "validate_migrated_schema":
			// Validate migrated schema
			if err := sem.ValidateSchemaEvolution(schemaName, toVersion, toVersion); err != nil {
				step["status"] = "failed"
				step["error"] = err.Error()
				step["completed_at"] = time.Now()
				executedSteps = append(executedSteps, step)
				result["status"] = "failed"
				result["failed_step"] = i + 1
				result["error"] = err.Error()
				result["executed_steps"] = executedSteps
				return result, err
			}

		case "update_version":
			// Update version (simulated)
			step["new_version"] = toVersion
		}

		step["status"] = "completed"
		step["completed_at"] = time.Now()
		executedSteps = append(executedSteps, step)
	}

	result["status"] = "completed"
	result["completed_at"] = time.Now()
	result["executed_steps"] = executedSteps

	// Create completion notification
	sem.createNotification(SchemaNotification{
		Type:       "evolution_completed",
		SchemaName: schemaName,
		Message:    fmt.Sprintf("Schema evolution completed for %s: %s → %s", schemaName, fromVersion, toVersion),
		CreatedAt:  time.Now(),
		Severity:   "info",
	})

	return result, nil
}
