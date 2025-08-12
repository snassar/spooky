// Package variables provides variable type definitions for the spooky codebase.
package variables

import (
	"time"

	spookytypescommon "spooky/internal/types/common"
	spookytypesschemas "spooky/internal/types/schemas"
)

// Variable represents a single variable definition
type Variable struct {
	spookytypescommon.CompleteEntity

	// Core variable properties
	Name        string        `hcl:"name" json:"name"`
	Type        VariableType  `hcl:"type" json:"type"`
	Description string        `hcl:"description,optional" json:"description,omitempty"`
	Default     interface{}   `hcl:"default,optional" json:"default,omitempty"`
	Required    bool          `hcl:"required,optional" json:"required,omitempty"`
	Sensitive   bool          `hcl:"sensitive,optional" json:"sensitive,omitempty"`
	Encrypted   bool          `hcl:"encrypted,optional" json:"encrypted,omitempty"`
	Scope       VariableScope `hcl:"scope,optional" json:"scope,omitempty"`

	// Dependencies and validation
	Dependencies []string             `hcl:"dependencies,optional" json:"dependencies,omitempty"`
	Validation   *VariableValidation  `hcl:"validation,block" json:"validation,omitempty"`
	Constraints  *VariableConstraints `hcl:"constraints,block" json:"constraints,omitempty"`

	// Metadata and source information
	SourceFile string                 `hcl:"-" json:"source_file,omitempty"`
	SourceLine int                    `hcl:"-" json:"source_line,omitempty"`
	Metadata   map[string]interface{} `hcl:"metadata,optional" json:"metadata,omitempty"`

	// Runtime state
	ResolvedValue   interface{} `hcl:"-" json:"resolved_value,omitempty"`
	IsResolved      bool        `hcl:"-" json:"is_resolved,omitempty"`
	ResolutionError string      `hcl:"-" json:"resolution_error,omitempty"`
}

// VariableType represents the type of a variable
type VariableType string

const (
	VariableTypeString   VariableType = "string"
	VariableTypeNumber   VariableType = "number"
	VariableTypeFloat    VariableType = "float"
	VariableTypeBool     VariableType = "bool"
	VariableTypeList     VariableType = "list"
	VariableTypeMap      VariableType = "map"
	VariableTypeObject   VariableType = "object"
	VariableTypeDuration VariableType = "duration"
	VariableTypeIP       VariableType = "ip"
	VariableTypeCIDR     VariableType = "cidr"
	VariableTypePath     VariableType = "path"
	VariableTypeFile     VariableType = "file"
	VariableTypeSecret   VariableType = "secret"
)

// VariableScope represents the scope of a variable
type VariableScope string

const (
	VariableScopeProject   VariableScope = "project"
	VariableScopeGlobal    VariableScope = "global"
	VariableScopeInherited VariableScope = "inherited"
)

// VariableValidation represents validation rules for a variable
type VariableValidation struct {
	Condition      string `hcl:"condition" json:"condition"`
	ErrorMessage   string `hcl:"error_message" json:"error_message"`
	WarningMessage string `hcl:"warning_message,optional" json:"warning_message,omitempty"`
}

// VariableConstraints represents type-specific constraints for a variable
type VariableConstraints struct {
	// String constraints
	MinLength *int    `hcl:"min_length,optional" json:"min_length,omitempty"`
	MaxLength *int    `hcl:"max_length,optional" json:"max_length,omitempty"`
	Pattern   *string `hcl:"pattern,optional" json:"pattern,omitempty"`

	// Numeric constraints
	MinValue *float64 `hcl:"min_value,optional" json:"min_value,omitempty"`
	MaxValue *float64 `hcl:"max_value,optional" json:"max_value,omitempty"`

	// List constraints
	MinItems *int `hcl:"min_items,optional" json:"min_items,omitempty"`
	MaxItems *int `hcl:"max_items,optional" json:"max_items,omitempty"`

	// File constraints
	FileExists   *bool   `hcl:"file_exists,optional" json:"file_exists,omitempty"`
	FileReadable *bool   `hcl:"file_readable,optional" json:"file_readable,omitempty"`
	FileSizeMax  *string `hcl:"file_size_max,optional" json:"file_size_max,omitempty"`

	// Path constraints
	PathExists   *bool `hcl:"path_exists,optional" json:"path_exists,omitempty"`
	PathAbsolute *bool `hcl:"path_absolute,optional" json:"path_absolute,omitempty"`
	PathRelative *bool `hcl:"path_relative,optional" json:"path_relative,omitempty"`
}

// VariableCollection represents a collection of variables
type VariableCollection struct {
	spookytypescommon.CompleteEntity

	Variables map[string]*Variable `hcl:"-" json:"variables"`
	Source    string               `hcl:"-" json:"source"`
	Count     int                  `hcl:"-" json:"count"`
}

// VariableResolutionConfig represents configuration for variable resolution
type VariableResolutionConfig struct {
	AllowSelfReference bool   `hcl:"allow_self_reference,optional" json:"allow_self_reference,omitempty"`
	AllowCircularDeps  bool   `hcl:"allow_circular_deps,optional" json:"allow_circular_deps,omitempty"`
	MaxResolutionDepth int    `hcl:"max_resolution_depth,optional" json:"max_resolution_depth,omitempty"`
	FailOnMissing      bool   `hcl:"fail_on_missing,optional" json:"fail_on_missing,omitempty"`
	UseEnvironment     bool   `hcl:"use_environment,optional" json:"use_environment,omitempty"`
	EnvironmentPrefix  string `hcl:"environment_prefix,optional" json:"environment_prefix,omitempty"`
}

// VariableSecurityConfig represents security configuration for variables
type VariableSecurityConfig struct {
	EncryptionEnabled bool   `hcl:"encryption_enabled,optional" json:"encryption_enabled,omitempty"`
	EncryptionMethod  string `hcl:"encryption_method,optional" json:"encryption_method,omitempty"`
	KeyFile           string `hcl:"key_file,optional" json:"key_file,omitempty"`
	KeyID             string `hcl:"key_id,optional" json:"key_id,omitempty"`
	MaskSensitive     bool   `hcl:"mask_sensitive,optional" json:"mask_sensitive,omitempty"`
	AuditTrail        bool   `hcl:"audit_trail,optional" json:"audit_trail,omitempty"`
}

// VariableContext represents the context for variable resolution
type VariableContext struct {
	ProjectPath string                 `json:"project_path"`
	Environment map[string]interface{} `json:"environment"`
	Facts       map[string]interface{} `json:"facts"`
	Machines    map[string]interface{} `json:"machines"`
	UserData    map[string]interface{} `json:"user_data"`
	Timestamp   time.Time              `json:"timestamp"`
}

// VariableResolutionResult represents the result of variable resolution
type VariableResolutionResult struct {
	Variables map[string]*Variable   `json:"variables"`
	Resolved  map[string]interface{} `json:"resolved"`
	Errors    []VariableError        `json:"errors,omitempty"`
	Warnings  []VariableWarning      `json:"warnings,omitempty"`
	Duration  time.Duration          `json:"duration"`
}

// VariableError represents an error during variable processing
type VariableError struct {
	spookytypescommon.ErrorDetails

	VariableName string                 `json:"variable_name"`
	ErrorType    VariableErrorType      `json:"error_type"`
	Context      map[string]interface{} `json:"context,omitempty"`
}

// VariableErrorType represents the type of variable error
type VariableErrorType string

const (
	VariableErrorTypeValidation   VariableErrorType = "validation"
	VariableErrorTypeResolution   VariableErrorType = "resolution"
	VariableErrorTypeDependency   VariableErrorType = "dependency"
	VariableErrorTypeCircular     VariableErrorType = "circular"
	VariableErrorTypeMissing      VariableErrorType = "missing"
	VariableErrorTypeTypeMismatch VariableErrorType = "type_mismatch"
	VariableErrorTypeConstraint   VariableErrorType = "constraint"
	VariableErrorTypeSecurity     VariableErrorType = "security"
)

// VariableWarning represents a warning during variable processing
type VariableWarning struct {
	spookytypescommon.ErrorDetails

	VariableName string                 `json:"variable_name"`
	WarningType  VariableWarningType    `json:"warning_type"`
	Context      map[string]interface{} `json:"context,omitempty"`
}

// VariableWarningType represents the type of variable warning
type VariableWarningType string

const (
	VariableWarningTypeDeprecated   VariableWarningType = "deprecated"
	VariableWarningTypeUnused       VariableWarningType = "unused"
	VariableWarningTypeSensitive    VariableWarningType = "sensitive"
	VariableWarningTypeUnencrypted  VariableWarningType = "unencrypted"
	VariableWarningTypeDefaultValue VariableWarningType = "default_value"
	VariableWarningTypeEnvironment  VariableWarningType = "environment"
	VariableWarningTypeDependency   VariableWarningType = "dependency"
)

// VariableValidationResult represents the result of variable validation
type VariableValidationResult struct {
	spookytypesschemas.ValidationResult

	Variables map[string]*Variable `json:"variables"`
	Valid     bool                 `json:"valid"`
	Errors    []VariableError      `json:"errors,omitempty"`
	Warnings  []VariableWarning    `json:"warnings,omitempty"`
}

// VariableExportOptions represents options for variable export
type VariableExportOptions struct {
	spookytypescommon.ExportOptions

	Format        string   `json:"format"`
	IncludeValues bool     `json:"include_values"`
	IncludeMeta   bool     `json:"include_meta"`
	Variables     []string `json:"variables,omitempty"`
	Encrypt       bool     `json:"encrypt,omitempty"`
}

// VariableImportOptions represents options for variable import
type VariableImportOptions struct {
	spookytypescommon.ImportOptions

	Format    string   `json:"format"`
	Overwrite bool     `json:"overwrite"`
	Variables []string `json:"variables,omitempty"`
	Decrypt   bool     `json:"decrypt,omitempty"`
	Validate  bool     `json:"validate,omitempty"`
}

// VariableQuery represents a query for variables
type VariableQuery struct {
	spookytypescommon.Query

	Name      string        `json:"name,omitempty"`
	Type      VariableType  `json:"type,omitempty"`
	Scope     VariableScope `json:"scope,omitempty"`
	Required  *bool         `json:"required,omitempty"`
	Sensitive *bool         `json:"sensitive,omitempty"`
	Source    string        `json:"source,omitempty"`
}

// VariableResult represents the result of a variable operation
type VariableResult struct {
	spookytypescommon.Result

	Variables []*Variable   `json:"variables"`
	Count     int           `json:"count"`
	Duration  time.Duration `json:"duration"`
}

// VariableFile represents a variables file
type VariableFile struct {
	spookytypescommon.CompleteEntity

	Path      string                 `json:"path"`
	Variables map[string]*Variable   `json:"variables"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Valid     bool                   `json:"valid"`
	Errors    []VariableError        `json:"errors,omitempty"`
	Warnings  []VariableWarning      `json:"warnings,omitempty"`
}

// VariableSource represents the source of a variable
type VariableSource struct {
	Type     string `json:"type"`
	File     string `json:"file,omitempty"`
	Line     int    `json:"line,omitempty"`
	Priority int    `json:"priority"`
}

// VariableDependency represents a variable dependency
type VariableDependency struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
	Resolved bool   `json:"resolved"`
	Error    string `json:"error,omitempty"`
}

// VariableDependencyGraph represents the dependency graph of variables
type VariableDependencyGraph struct {
	Nodes  map[string]*Variable `json:"nodes"`
	Edges  map[string][]string  `json:"edges"`
	Cycles [][]string           `json:"cycles,omitempty"`
}

// VariableStatistics represents statistics about variables
type VariableStatistics struct {
	TotalCount     int `json:"total_count"`
	RequiredCount  int `json:"required_count"`
	SensitiveCount int `json:"sensitive_count"`
	EncryptedCount int `json:"encrypted_count"`
	ResolvedCount  int `json:"resolved_count"`
	ErrorCount     int `json:"error_count"`
	WarningCount   int `json:"warning_count"`

	TypeBreakdown  map[VariableType]int  `json:"type_breakdown"`
	ScopeBreakdown map[VariableScope]int `json:"scope_breakdown"`
}
