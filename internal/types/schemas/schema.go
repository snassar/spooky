// Package schemas provides types for schema validation and management in the spooky codebase.
// These types define the structure for HCL schemas and validation operations.
package schemas

import (
	"time"
)

// Schema represents a HCL schema definition
type Schema struct {
	// Schema metadata
	Version     string    `json:"version" hcl:"version"`
	Type        string    `json:"type" hcl:"type"`
	Name        string    `json:"name" hcl:"name"`
	Description string    `json:"description" hcl:"description"`
	CreatedAt   time.Time `json:"created_at" hcl:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" hcl:"updated_at"`

	// Schema content
	Content string `json:"content" hcl:"content"`

	// Schema validation rules
	Validation *SchemaValidation `json:"validation,omitempty" hcl:"validation,optional"`

	// Schema evolution and migration
	Evolution *SchemaEvolution `json:"evolution,omitempty" hcl:"evolution,optional"`

	// Schema metadata
	Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`

	// Schema compatibility
	Compatibility *SchemaCompatibility `json:"compatibility,omitempty" hcl:"compatibility,optional"`
}

// SchemaValidation represents schema validation configuration
type SchemaValidation struct {
	// Whether validation is enabled
	Enabled bool `json:"enabled" hcl:"enabled"`

	// Validation mode (strict, lenient, permissive)
	Mode string `json:"mode" hcl:"mode"`

	// Validation rules
	Rules []ValidationRule `json:"rules,omitempty" hcl:"rules,optional"`

	// Custom validation functions
	CustomValidators []string `json:"custom_validators,omitempty" hcl:"custom_validators,optional"`

	// Validation error handling
	ErrorHandling *ValidationErrorHandling `json:"error_handling,omitempty" hcl:"error_handling,optional"`

	// Field validation rules
	Fields map[string]*FieldValidation `json:"fields,omitempty" hcl:"fields,optional"`

	// Cross-field validation rules
	CrossFieldValidations []CrossFieldValidation `json:"cross_field_validations,omitempty" hcl:"cross_field_validations,optional"`
}

// FieldValidation represents validation rules for a specific field
type FieldValidation struct {
	// Field type
	Type string `json:"type" hcl:"type"`

	// Whether field is required
	Required bool `json:"required" hcl:"required"`

	// Default value
	Default interface{} `json:"default,omitempty" hcl:"default,optional"`

	// Validation constraints
	Constraints *FieldConstraints `json:"constraints,omitempty" hcl:"constraints,optional"`

	// Custom validation rules
	CustomRules []string `json:"custom_rules,omitempty" hcl:"custom_rules,optional"`

	// Field description
	Description string `json:"description,omitempty" hcl:"description,optional"`

	// Field examples
	Examples []interface{} `json:"examples,omitempty" hcl:"examples,optional"`
}

// FieldConstraints represents validation constraints for a field
type FieldConstraints struct {
	// String constraints
	MinLength *int    `json:"min_length,omitempty" hcl:"min_length,optional"`
	MaxLength *int    `json:"max_length,omitempty" hcl:"max_length,optional"`
	Pattern   *string `json:"pattern,omitempty" hcl:"pattern,optional"`
	Format    *string `json:"format,omitempty" hcl:"format,optional"`

	// Numeric constraints
	Min *float64 `json:"min,omitempty" hcl:"min,optional"`
	Max *float64 `json:"max,omitempty" hcl:"max,optional"`

	// Array constraints
	MinItems *int `json:"min_items,omitempty" hcl:"min_items,optional"`
	MaxItems *int `json:"max_items,omitempty" hcl:"max_items,optional"`

	// Enum values
	Enum []interface{} `json:"enum,omitempty" hcl:"enum,optional"`

	// Custom constraint expression
	Expression *string `json:"expression,omitempty" hcl:"expression,optional"`
}

// CrossFieldValidation represents validation rules that depend on multiple fields
type CrossFieldValidation struct {
	// Validation name
	Name string `json:"name" hcl:"name"`

	// Fields involved in validation
	Fields []string `json:"fields" hcl:"fields"`

	// Validation condition
	Condition string `json:"condition" hcl:"condition"`

	// Error message
	Message string `json:"message" hcl:"message"`

	// Validation severity
	Severity string `json:"severity" hcl:"severity"` // "error", "warning", "info"
}

// ValidationRule represents a validation rule
type ValidationRule struct {
	// Rule name
	Name string `json:"name" hcl:"name"`

	// Rule type
	Type string `json:"type" hcl:"type"`

	// Rule pattern (for regex rules)
	Pattern string `json:"pattern,omitempty" hcl:"pattern,optional"`

	// Rule condition (for conditional rules)
	Condition string `json:"condition,omitempty" hcl:"condition,optional"`

	// Rule message
	Message string `json:"message" hcl:"message"`

	// Rule severity
	Severity string `json:"severity" hcl:"severity"` // "error", "warning", "info"

	// Rule parameters
	Parameters map[string]interface{} `json:"parameters,omitempty" hcl:"parameters,optional"`

	// Rule dependencies
	Dependencies []string `json:"dependencies,omitempty" hcl:"dependencies,optional"`

	// Rule processing order
	Priority int `json:"priority,omitempty" hcl:"priority,optional"`
}

// ValidationErrorHandling represents validation error handling configuration
type ValidationErrorHandling struct {
	// Whether to stop on first error
	StopOnFirstError bool `json:"stop_on_first_error" hcl:"stop_on_first_error"`

	// Maximum number of errors to collect
	MaxErrors int `json:"max_errors" hcl:"max_errors"`

	// Whether to include warnings
	IncludeWarnings bool `json:"include_warnings" hcl:"include_warnings"`

	// Whether to include context in errors
	IncludeContext bool `json:"include_context" hcl:"include_context"`

	// Whether to include suggestions in errors
	IncludeSuggestions bool `json:"include_suggestions" hcl:"include_suggestions"`

	// Error grouping strategy
	GroupErrors bool `json:"group_errors" hcl:"group_errors"`
}

// SchemaEvolution represents schema evolution and migration support
type SchemaEvolution struct {
	// Schema version history
	Versions []SchemaVersion `json:"versions,omitempty" hcl:"versions,optional"`

	// Migration scripts
	Migrations []SchemaMigration `json:"migrations,omitempty" hcl:"migrations,optional"`

	// Deprecation notices
	Deprecations []SchemaDeprecation `json:"deprecations,omitempty" hcl:"deprecations,optional"`

	// Breaking changes
	BreakingChanges []BreakingChange `json:"breaking_changes,omitempty" hcl:"breaking_changes,optional"`

	// Evolution policy
	Policy *EvolutionPolicy `json:"policy,omitempty" hcl:"policy,optional"`
}

// SchemaVersion represents a schema version
type SchemaVersion struct {
	// Version number
	Version string `json:"version" hcl:"version"`

	// Release date
	ReleasedAt time.Time `json:"released_at" hcl:"released_at"`

	// Version description
	Description string `json:"description,omitempty" hcl:"description,optional"`

	// Changes in this version
	Changes []VersionChange `json:"changes,omitempty" hcl:"changes,optional"`

	// Compatibility information
	Compatibility *VersionCompatibility `json:"compatibility,omitempty" hcl:"compatibility,optional"`
}

// VersionChange represents a change in a schema version
type VersionChange struct {
	// Change type
	Type string `json:"type" hcl:"type"` // "added", "removed", "modified", "deprecated"

	// Change description
	Description string `json:"description" hcl:"description"`

	// Affected fields
	Fields []string `json:"fields,omitempty" hcl:"fields,optional"`

	// Migration guidance
	MigrationGuidance string `json:"migration_guidance,omitempty" hcl:"migration_guidance,optional"`

	// Breaking change indicator
	Breaking bool `json:"breaking,omitempty" hcl:"breaking,optional"`
}

// VersionCompatibility represents version compatibility information
type VersionCompatibility struct {
	// Compatible versions
	CompatibleWith []string `json:"compatible_with,omitempty" hcl:"compatible_with,optional"`

	// Incompatible versions
	IncompatibleWith []string `json:"incompatible_with,omitempty" hcl:"incompatible_with,optional"`

	// Migration required
	MigrationRequired bool `json:"migration_required,omitempty" hcl:"migration_required,optional"`
}

// SchemaMigration represents a schema migration
type SchemaMigration struct {
	// Migration name
	Name string `json:"name" hcl:"name"`

	// Source version
	FromVersion string `json:"from_version" hcl:"from_version"`

	// Target version
	ToVersion string `json:"to_version" hcl:"to_version"`

	// Migration script
	Script string `json:"script,omitempty" hcl:"script,optional"`

	// Migration function
	Function string `json:"function,omitempty" hcl:"function,optional"`

	// Migration description
	Description string `json:"description,omitempty" hcl:"description,optional"`

	// Migration validation
	Validation *MigrationValidation `json:"validation,omitempty" hcl:"validation,optional"`
}

// MigrationValidation represents migration validation rules
type MigrationValidation struct {
	// Pre-migration validation
	PreValidation []string `json:"pre_validation,omitempty" hcl:"pre_validation,optional"`

	// Post-migration validation
	PostValidation []string `json:"post_validation,omitempty" hcl:"post_validation,optional"`

	// Rollback validation
	RollbackValidation []string `json:"rollback_validation,omitempty" hcl:"rollback_validation,optional"`
}

// SchemaDeprecation represents a schema deprecation
type SchemaDeprecation struct {
	// Deprecated field or feature
	Field string `json:"field" hcl:"field"`

	// Deprecation date
	DeprecatedAt time.Time `json:"deprecated_at" hcl:"deprecated_at"`

	// Removal date
	RemovedAt *time.Time `json:"removed_at,omitempty" hcl:"removed_at,optional"`

	// Deprecation reason
	Reason string `json:"reason,omitempty" hcl:"reason,optional"`

	// Replacement field
	Replacement string `json:"replacement,omitempty" hcl:"replacement,optional"`

	// Migration guidance
	MigrationGuidance string `json:"migration_guidance,omitempty" hcl:"migration_guidance,optional"`
}

// BreakingChange represents a breaking change
type BreakingChange struct {
	// Changed field or feature
	Field string `json:"field" hcl:"field"`

	// Change description
	Description string `json:"description" hcl:"description"`

	// Impact level
	Impact string `json:"impact" hcl:"impact"` // "high", "medium", "low"

	// Mitigation strategy
	Mitigation string `json:"mitigation,omitempty" hcl:"mitigation,optional"`

	// Affected versions
	AffectedVersions []string `json:"affected_versions,omitempty" hcl:"affected_versions,optional"`
}

// EvolutionPolicy represents schema evolution policy
type EvolutionPolicy struct {
	// Versioning strategy
	VersioningStrategy string `json:"versioning_strategy" hcl:"versioning_strategy"`

	// Breaking change policy
	BreakingChangePolicy string `json:"breaking_change_policy" hcl:"breaking_change_policy"`

	// Deprecation policy
	DeprecationPolicy string `json:"deprecation_policy" hcl:"deprecation_policy"`

	// Migration policy
	MigrationPolicy string `json:"migration_policy" hcl:"migration_policy"`

	// Support duration
	SupportDuration time.Duration `json:"support_duration,omitempty" hcl:"support_duration,optional"`
}

// SchemaCompatibility represents schema compatibility information
type SchemaCompatibility struct {
	// Compatible schema types
	CompatibleTypes []string `json:"compatible_types,omitempty" hcl:"compatible_types,optional"`

	// Compatible versions
	CompatibleVersions []string `json:"compatible_versions,omitempty" hcl:"compatible_versions,optional"`

	// Compatibility matrix
	CompatibilityMatrix map[string][]string `json:"compatibility_matrix,omitempty" hcl:"compatibility_matrix,optional"`

	// Forward compatibility
	ForwardCompatible bool `json:"forward_compatible,omitempty" hcl:"forward_compatible,optional"`

	// Backward compatibility
	BackwardCompatible bool `json:"backward_compatible,omitempty" hcl:"backward_compatible,optional"`
}

// SchemaError represents a schema-related error
type SchemaError struct {
	// Error details
	Code        string                 `json:"code" hcl:"code"`
	Message     string                 `json:"message" hcl:"message"`
	Context     map[string]interface{} `json:"context,omitempty" hcl:"context,optional"`
	Stack       []string               `json:"stack,omitempty" hcl:"stack,optional"`
	Recoverable bool                   `json:"recoverable" hcl:"recoverable"`

	// Schema information
	SchemaName string `json:"schema_name" hcl:"schema_name"`
	SchemaType string `json:"schema_type" hcl:"schema_type"`

	// Validation information
	FieldPath string      `json:"field_path,omitempty" hcl:"field_path,optional"`
	Value     interface{} `json:"value,omitempty" hcl:"value,optional"`

	// Error severity
	Severity string `json:"severity" hcl:"severity"` // "error", "warning", "info"

	// Error suggestions
	Suggestions []string `json:"suggestions,omitempty" hcl:"suggestions,optional"`

	// Error documentation
	Documentation string `json:"documentation,omitempty" hcl:"documentation,optional"`

	// Error timestamp
	Timestamp time.Time `json:"timestamp" hcl:"timestamp"`

	// Error location
	Location *ErrorLocation `json:"location,omitempty" hcl:"location,optional"`

	// Related errors
	RelatedErrors []*SchemaError `json:"related_errors,omitempty" hcl:"related_errors,optional"`
}

// ErrorLocation represents the location of an error
type ErrorLocation struct {
	// File path
	FilePath string `json:"file_path,omitempty" hcl:"file_path,optional"`

	// Line number
	Line int `json:"line,omitempty" hcl:"line,optional"`

	// Column number
	Column int `json:"column,omitempty" hcl:"column,optional"`

	// Byte offset
	ByteOffset int `json:"byte_offset,omitempty" hcl:"byte_offset,optional"`

	// Character offset
	CharOffset int `json:"char_offset,omitempty" hcl:"char_offset,optional"`
}

// NewSchemaError creates a new schema error
func NewSchemaError(schemaName, schemaType, message string) *SchemaError {
	return &SchemaError{
		Code:        "schema_error",
		Message:     message,
		Recoverable: true,
		SchemaName:  schemaName,
		SchemaType:  schemaType,
		Severity:    "error",
		Timestamp:   time.Now(),
	}
}

// Error implements the error interface
func (e *SchemaError) Error() string {
	return e.Message
}

// Unwrap returns the underlying error
func (e *SchemaError) Unwrap() error {
	return nil
}

// AddSuggestion adds a suggestion to the error
func (e *SchemaError) AddSuggestion(suggestion string) {
	e.Suggestions = append(e.Suggestions, suggestion)
}

// AddContext adds context information to the error
func (e *SchemaError) AddContext(key string, value interface{}) {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
}

// SetLocation sets the error location
func (e *SchemaError) SetLocation(filePath string, line, column int) {
	e.Location = &ErrorLocation{
		FilePath: filePath,
		Line:     line,
		Column:   column,
	}
}

// SchemaRegistry manages available schemas
type SchemaRegistry interface {
	// Register registers a new schema
	Register(schema *Schema) error

	// Get returns a schema by name and type
	Get(name, schemaType string) (*Schema, bool)

	// List returns all registered schemas
	List() []*Schema

	// ListByType returns schemas by type
	ListByType(schemaType string) []*Schema

	// ValidateData validates data against a schema by name and type
	ValidateData(schemaName, schemaType string, data interface{}) (*ValidationResult, error)

	// CheckCompatibility checks compatibility between schemas
	CheckCompatibility(schema1, schema2 *Schema) (*CompatibilityResult, error)

	// GetMigrationPath gets migration path between schema versions
	GetMigrationPath(fromVersion, toVersion string) ([]*SchemaMigration, error)
}

// CompatibilityResult represents the result of a compatibility check
type CompatibilityResult struct {
	// Whether schemas are compatible
	Compatible bool `json:"compatible" hcl:"compatible"`

	// Compatibility issues
	Issues []CompatibilityIssue `json:"issues,omitempty" hcl:"issues,optional"`

	// Migration requirements
	MigrationRequired bool `json:"migration_required,omitempty" hcl:"migration_required,optional"`

	// Migration path
	MigrationPath []*SchemaMigration `json:"migration_path,omitempty" hcl:"migration_path,optional"`
}

// CompatibilityIssue represents a compatibility issue
type CompatibilityIssue struct {
	// Issue type
	Type string `json:"type" hcl:"type"` // "breaking", "deprecation", "warning"

	// Issue description
	Description string `json:"description" hcl:"description"`

	// Affected fields
	Fields []string `json:"fields,omitempty" hcl:"fields,optional"`

	// Issue severity
	Severity string `json:"severity" hcl:"severity"` // "error", "warning", "info"

	// Resolution guidance
	Resolution string `json:"resolution,omitempty" hcl:"resolution,optional"`
}

// ValidationResult represents the result of a validation operation
type ValidationResult struct {
	// Whether the validation passed
	Valid bool `json:"valid" hcl:"valid"`

	// Validation timestamp
	ValidatedAt time.Time `json:"validated_at" hcl:"validated_at"`

	// Validation errors
	Errors []SchemaError `json:"errors,omitempty" hcl:"errors,optional"`

	// Validation warnings
	Warnings []SchemaError `json:"warnings,omitempty" hcl:"warnings,optional"`

	// Validation info messages
	Info []SchemaError `json:"info,omitempty" hcl:"info,optional"`

	// Validation details
	Details map[string]interface{} `json:"details,omitempty" hcl:"details,optional"`

	// Validation statistics
	Statistics *ValidationStatistics `json:"statistics,omitempty" hcl:"statistics,optional"`

	// Validation recommendations
	Recommendations []string `json:"recommendations,omitempty" hcl:"recommendations,optional"`
}

// ValidationStatistics represents validation statistics
type ValidationStatistics struct {
	// Total fields validated
	TotalFields int `json:"total_fields" hcl:"total_fields"`

	// Valid fields
	ValidFields int `json:"valid_fields" hcl:"valid_fields"`

	// Invalid fields
	InvalidFields int `json:"invalid_fields" hcl:"invalid_fields"`

	// Validation duration
	Duration time.Duration `json:"duration,omitempty" hcl:"duration,optional"`

	// Rules processed
	RulesProcessed int `json:"rules_processed" hcl:"rules_processed"`

	// Rules failed
	RulesFailed int `json:"rules_failed" hcl:"rules_failed"`
}

// SchemaValidator provides schema validation functionality
type SchemaValidator interface {
	// Validate validates data against a schema
	Validate(schema *Schema, data interface{}) (*ValidationResult, error)

	// ValidateFile validates a file against a schema
	ValidateFile(schema *Schema, filePath string) (*ValidationResult, error)

	// ValidateString validates a string against a schema
	ValidateString(schema *Schema, content string) (*ValidationResult, error)

	// ValidateBytes validates bytes against a schema
	ValidateBytes(schema *Schema, data []byte) (*ValidationResult, error)

	// ValidateWithContext validates data with additional context
	ValidateWithContext(schema *Schema, data interface{}, context map[string]interface{}) (*ValidationResult, error)

	// ValidateField validates a specific field
	ValidateField(schema *Schema, fieldPath string, value interface{}) (*ValidationResult, error)
}

// SchemaLoader provides schema loading functionality
type SchemaLoader interface {
	// Load loads a schema from a file
	Load(filePath string) (*Schema, error)

	// LoadFromString loads a schema from a string
	LoadFromString(content string) (*Schema, error)

	// LoadFromBytes loads a schema from bytes
	LoadFromBytes(data []byte) (*Schema, error)

	// LoadEmbedded loads an embedded schema
	LoadEmbedded(name string) (*Schema, error)

	// LoadWithValidation loads a schema with validation
	LoadWithValidation(filePath string) (*Schema, *ValidationResult, error)

	// LoadMultiple loads multiple schemas
	LoadMultiple(filePaths []string) (map[string]*Schema, error)
}

// SchemaManager provides comprehensive schema management functionality
type SchemaManager interface {
	// SchemaRegistry provides schema registration and retrieval
	SchemaRegistry

	// SchemaValidator provides schema validation
	SchemaValidator

	// SchemaLoader provides schema loading
	SchemaLoader

	// Schema evolution and migration
	CheckForUpdates(schema *Schema) ([]SchemaUpdate, error)
	ApplyMigration(schema *Schema, migration *SchemaMigration) error
	ValidateMigration(schema *Schema, migration *SchemaMigration) (*ValidationResult, error)

	// Schema analysis
	AnalyzeSchema(schema *Schema) (*SchemaAnalysis, error)
	CompareSchemas(schema1, schema2 *Schema) (*SchemaComparison, error)

	// Schema documentation
	GenerateDocumentation(schema *Schema) (string, error)
	GenerateExamples(schema *Schema) ([]string, error)
}

// SchemaUpdate represents a schema update
type SchemaUpdate struct {
	// Update type
	Type string `json:"type" hcl:"type"` // "patch", "minor", "major"

	// Update description
	Description string `json:"description" hcl:"description"`

	// Update version
	Version string `json:"version" hcl:"version"`

	// Update date
	Date time.Time `json:"date" hcl:"date"`

	// Update changes
	Changes []VersionChange `json:"changes,omitempty" hcl:"changes,optional"`

	// Update priority
	Priority string `json:"priority" hcl:"priority"` // "low", "medium", "high", "critical"
}

// SchemaAnalysis represents schema analysis results
type SchemaAnalysis struct {
	// Schema complexity
	Complexity *SchemaComplexity `json:"complexity,omitempty" hcl:"complexity,optional"`

	// Schema coverage
	Coverage *SchemaCoverage `json:"coverage,omitempty" hcl:"coverage,optional"`

	// Schema quality
	Quality *SchemaQuality `json:"quality,omitempty" hcl:"quality,optional"`

	// Schema recommendations
	Recommendations []string `json:"recommendations,omitempty" hcl:"recommendations,optional"`
}

// SchemaComplexity represents schema complexity metrics
type SchemaComplexity struct {
	// Number of fields
	FieldCount int `json:"field_count" hcl:"field_count"`

	// Number of validation rules
	RuleCount int `json:"rule_count" hcl:"rule_count"`

	// Nesting depth
	MaxDepth int `json:"max_depth" hcl:"max_depth"`

	// Cyclomatic complexity
	CyclomaticComplexity int `json:"cyclomatic_complexity" hcl:"cyclomatic_complexity"`

	// Complexity score
	Score float64 `json:"score" hcl:"score"`
}

// SchemaCoverage represents schema coverage metrics
type SchemaCoverage struct {
	// Field coverage percentage
	FieldCoverage float64 `json:"field_coverage" hcl:"field_coverage"`

	// Validation coverage percentage
	ValidationCoverage float64 `json:"validation_coverage" hcl:"validation_coverage"`

	// Documentation coverage percentage
	DocumentationCoverage float64 `json:"documentation_coverage" hcl:"documentation_coverage"`

	// Example coverage percentage
	ExampleCoverage float64 `json:"example_coverage" hcl:"example_coverage"`
}

// SchemaQuality represents schema quality metrics
type SchemaQuality struct {
	// Overall quality score
	Score float64 `json:"score" hcl:"score"`

	// Quality dimensions
	Dimensions map[string]float64 `json:"dimensions,omitempty" hcl:"dimensions,optional"`

	// Quality issues
	Issues []string `json:"issues,omitempty" hcl:"issues,optional"`

	// Quality improvements
	Improvements []string `json:"improvements,omitempty" hcl:"improvements,optional"`
}

// SchemaComparison represents schema comparison results
type SchemaComparison struct {
	// Comparison summary
	Summary string `json:"summary" hcl:"summary"`

	// Added fields
	Added []string `json:"added,omitempty" hcl:"added,optional"`

	// Removed fields
	Removed []string `json:"removed,omitempty" hcl:"removed,optional"`

	// Modified fields
	Modified []string `json:"modified,omitempty" hcl:"modified,optional"`

	// Breaking changes
	BreakingChanges []string `json:"breaking_changes,omitempty" hcl:"breaking_changes,optional"`

	// Compatibility status
	Compatible bool `json:"compatible" hcl:"compatible"`
}
