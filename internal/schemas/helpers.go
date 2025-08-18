package schemas

import (
	"fmt"
	"time"

	spookytypeslogging "spooky/internal/types/logging"
	spookytypesschemas "spooky/internal/types/schemas"
)

// SchemaHelpers provides common functionality shared across schema components
type SchemaHelpers struct {
	logger spookytypeslogging.Logger
}

// NewSchemaHelpers creates a new SchemaHelpers instance
func NewSchemaHelpers(logger spookytypeslogging.Logger) *SchemaHelpers {
	return &SchemaHelpers{
		logger: logger,
	}
}

// ValidateMigration validates a migration before applying it
func (h *SchemaHelpers) ValidateMigration(schema *spookytypesschemas.Schema, migration *spookytypesschemas.SchemaMigration) (*spookytypesschemas.ValidationResult, error) {
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
		schemaError := spookytypesschemas.NewSchemaError(schema.Name, schema.Type,
			fmt.Sprintf("Migration from version %s cannot be applied to schema version %s",
				migration.FromVersion, schema.Version))
		schemaError.Severity = ValidationError
		result.Errors = append(result.Errors, *schemaError)
		result.Valid = false
	}

	// Validate migration prerequisites
	if migration.Validation != nil {
		for _, preValidation := range migration.Validation.PreValidation {
			h.logger.Debug("Running pre-validation", map[string]interface{}{
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

// CheckForUpdates checks for available schema updates
func (h *SchemaHelpers) CheckForUpdates(schema *spookytypesschemas.Schema, registeredSchemas []*spookytypesschemas.Schema) ([]spookytypesschemas.SchemaUpdate, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema cannot be nil")
	}

	var updates []spookytypesschemas.SchemaUpdate

	// Check for newer versions
	for _, registeredSchema := range registeredSchemas {
		if registeredSchema.Type == schema.Type && registeredSchema.Name == schema.Name {
			if registeredSchema.Version != schema.Version {
				update := spookytypesschemas.SchemaUpdate{
					Type:        h.determineUpdateType(schema.Version, registeredSchema.Version),
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

// CheckCompatibility checks compatibility between two schema versions
func (h *SchemaHelpers) CheckCompatibility(schema1, schema2 *spookytypesschemas.Schema) (*spookytypesschemas.CompatibilityResult, error) {
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
	h.checkFieldCompatibility(schema1, schema2, result)

	// Check validation rule compatibility
	h.checkValidationCompatibility(schema1, schema2, result)

	// Check evolution compatibility
	h.checkEvolutionCompatibility(schema1, schema2, result)

	// Determine if migration is required
	result.MigrationRequired = len(result.Issues) > 0

	return result, nil
}

// checkFieldCompatibility checks field compatibility between schemas
func (h *SchemaHelpers) checkFieldCompatibility(schema1, schema2 *spookytypesschemas.Schema, result *spookytypesschemas.CompatibilityResult) {
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
			if h.fieldValidationChanged(field1, field2) {
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
func (h *SchemaHelpers) checkValidationCompatibility(schema1, schema2 *spookytypesschemas.Schema, result *spookytypesschemas.CompatibilityResult) {
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

	for i := range schema1.Validation.Rules {
		rule := &schema1.Validation.Rules[i]
		rules1[rule.Name] = *rule
	}

	for i := range schema2.Validation.Rules {
		rule := &schema2.Validation.Rules[i]
		rules2[rule.Name] = *rule
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
				Type:        "info",
				Description: fmt.Sprintf("Validation rule '%s' was added", ruleName),
				Severity:    "info",
				Resolution:  fmt.Sprintf("Validation rule '%s' is now enforced", ruleName),
			}
			result.Issues = append(result.Issues, issue)
		}
	}
}

// checkEvolutionCompatibility checks evolution compatibility
func (h *SchemaHelpers) checkEvolutionCompatibility(schema1, schema2 *spookytypesschemas.Schema, result *spookytypesschemas.CompatibilityResult) {
	_ = schema1 // schema1 is not used in this implementation
	// Check breaking changes
	if schema2.Evolution != nil {
		for _, breakingChange := range schema2.Evolution.BreakingChanges {
			issue := spookytypesschemas.CompatibilityIssue{
				Type:        "breaking",
				Description: breakingChange.Description,
				Severity:    "error",
				Resolution:  breakingChange.Mitigation,
			}
			result.Issues = append(result.Issues, issue)
			result.Compatible = false
		}
	}

	// Check deprecations
	if schema2.Evolution != nil {
		for _, deprecation := range schema2.Evolution.Deprecations {
			issue := spookytypesschemas.CompatibilityIssue{
				Type:        "deprecation",
				Description: deprecation.Reason,
				Severity:    "warning",
				Resolution:  deprecation.Replacement,
			}
			result.Issues = append(result.Issues, issue)
		}
	}
}

// fieldValidationChanged checks if field validation changed
func (h *SchemaHelpers) fieldValidationChanged(field1, field2 *spookytypesschemas.FieldValidation) bool {
	if field1.Required != field2.Required {
		return true
	}
	if field1.Type != field2.Type {
		return true
	}
	// Add more detailed comparison as needed
	return false
}

// determineUpdateType determines the type of update
func (h *SchemaHelpers) determineUpdateType(currentVersion, newVersion string) string {
	// Simple version comparison - in a real implementation, you'd use proper semver parsing
	if currentVersion < newVersion {
		return "upgrade"
	}
	return "downgrade"
}
