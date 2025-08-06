package types

import (
	"time"
)

// Variable represents a single variable
type Variable struct {
	Name         string                 `hcl:"name,label" json:"name"`
	Type         string                 `hcl:"type" json:"type"`
	Value        interface{}            `hcl:"value,optional" json:"value,omitempty"`
	Default      interface{}            `hcl:"default,optional" json:"default,omitempty"`
	Description  string                 `hcl:"description,optional" json:"description,omitempty"`
	Required     bool                   `hcl:"required,optional" json:"required,omitempty"`
	Sensitive    bool                   `hcl:"sensitive,optional" json:"sensitive,omitempty"`
	Encrypted    bool                   `hcl:"encrypted,optional" json:"encrypted,omitempty"`
	Scope        string                 `hcl:"scope,optional" json:"scope,omitempty"`
	Dependencies []string               `hcl:"dependencies,optional" json:"dependencies,omitempty"`
	Validation   *VariableValidation    `hcl:"validation,optional" json:"validation,omitempty"`
	Constraints  *VariableConstraints   `hcl:"constraints,optional" json:"constraints,omitempty"`
	Metadata     map[string]interface{} `hcl:"metadata,optional" json:"metadata,omitempty"`
	Resolved     bool                   `hcl:"resolved,optional" json:"resolved,omitempty"`
	CreatedAt    time.Time              `hcl:"created_at,optional" json:"created_at,omitempty"`
	UpdatedAt    time.Time              `hcl:"updated_at,optional" json:"updated_at,omitempty"`
}

// VariableValidation represents validation rules for a variable
type VariableValidation struct {
	Condition      string `hcl:"condition" json:"condition"`
	ErrorMessage   string `hcl:"error_message" json:"error_message"`
	WarningMessage string `hcl:"warning_message,optional" json:"warning_message,omitempty"`
}

// VariableConstraints represents type-specific constraints
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
	Variables []*Variable            `hcl:"variables" json:"variables"`
	Path      string                 `hcl:"path,optional" json:"path,omitempty"`
	Metadata  map[string]interface{} `hcl:"metadata,optional" json:"metadata,omitempty"`
}

// VariableContext represents a context for variable resolution
type VariableContext struct {
	Variables   map[string]*Variable `hcl:"variables" json:"variables"`
	ProjectPath string               `hcl:"project_path,optional" json:"project_path,omitempty"`
	Environment map[string]string    `hcl:"environment,optional" json:"environment,omitempty"`
}

// Configuration types
type Config struct {
	LoadingConfig      *LoadingConfig      `hcl:"loading,optional" json:"loading,omitempty"`
	ResolutionConfig   *ResolutionConfig   `hcl:"resolution,optional" json:"resolution,omitempty"`
	ValidationConfig   *ValidationConfig   `hcl:"validation,optional" json:"validation,omitempty"`
	ImportExportConfig *ImportExportConfig `hcl:"import_export,optional" json:"import_export,omitempty"`
}

type LoadingConfig struct {
	DefaultEncoding   string   `hcl:"default_encoding,optional" json:"default_encoding,omitempty"`
	MaxFileSize       int64    `hcl:"max_file_size,optional" json:"max_file_size,omitempty"`
	AllowedExtensions []string `hcl:"allowed_extensions,optional" json:"allowed_extensions,omitempty"`
}

type ResolutionConfig struct {
	MaxRecursionDepth int                    `hcl:"max_recursion_depth,optional" json:"max_recursion_depth,omitempty"`
	DefaultValues     map[string]interface{} `hcl:"default_values,optional" json:"default_values,omitempty"`
	StrictMode        bool                   `hcl:"strict_mode,optional" json:"strict_mode,omitempty"`
}

type ValidationConfig struct {
	ValidationRules     *ValidationRules `hcl:"validation_rules,optional" json:"validation_rules,omitempty"`
	StrictValidation    bool             `hcl:"strict_validation,optional" json:"strict_validation,omitempty"`
	MaxValidationErrors int              `hcl:"max_validation_errors,optional" json:"max_validation_errors,omitempty"`
}

type ImportExportConfig struct {
	ExportOptions *ExportOptions `hcl:"export_options,optional" json:"export_options,omitempty"`
	ImportOptions *ImportOptions `hcl:"import_options,optional" json:"import_options,omitempty"`
	DefaultFormat ExportFormat   `hcl:"default_format,optional" json:"default_format,omitempty"`
}

// Import/export configuration types
type HCLOptions struct {
	IncludeMetadata bool `hcl:"include_metadata,optional" json:"include_metadata,omitempty"`
	PrettyPrint     bool `hcl:"pretty_print,optional" json:"pretty_print,omitempty"`
	ValidateSchema  bool `hcl:"validate_schema,optional" json:"validate_schema,omitempty"`
}

type JSONOptions struct {
	IncludeMetadata bool `hcl:"include_metadata,optional" json:"include_metadata,omitempty"`
	PrettyPrint     bool `hcl:"pretty_print,optional" json:"pretty_print,omitempty"`
	IndentSize      int  `hcl:"indent_size,optional" json:"indent_size,omitempty"`
}

// Export and import types
type ExportFormat string
type ImportFormat string

const (
	ExportFormatHCL  ExportFormat = "hcl"
	ExportFormatJSON ExportFormat = "json"

	ImportFormatHCL  ImportFormat = "hcl"
	ImportFormatJSON ImportFormat = "json"
)

type ExportOptions struct {
	IncludeMetadata bool `hcl:"include_metadata,optional" json:"include_metadata,omitempty"`
	PrettyPrint     bool `hcl:"pretty_print,optional" json:"pretty_print,omitempty"`
}

type ImportOptions struct {
	MergePolicy string `hcl:"merge_policy,optional" json:"merge_policy,omitempty"`
	Overwrite   bool   `hcl:"overwrite,optional" json:"overwrite,omitempty"`
}

// Validation types
type ValidationResult struct {
	Valid    bool                `hcl:"valid" json:"valid"`
	Errors   []ValidationError   `hcl:"errors,optional" json:"errors,omitempty"`
	Warnings []ValidationWarning `hcl:"warnings,optional" json:"warnings,omitempty"`
}

type ValidationError struct {
	Field   string `hcl:"field" json:"field"`
	Message string `hcl:"message" json:"message"`
}

type ValidationWarning struct {
	Field   string `hcl:"field" json:"field"`
	Message string `hcl:"message" json:"message"`
}

type ValidationRules struct {
	RequiredFields  []string `hcl:"required_fields,optional" json:"required_fields,omitempty"`
	ForbiddenFields []string `hcl:"forbidden_fields,optional" json:"forbidden_fields,omitempty"`
}

// Dependency types
type DependencyGraph struct {
	Nodes map[string]*DependencyNode `hcl:"nodes" json:"nodes"`
	Edges map[string][]string        `hcl:"edges" json:"edges"`
}

type DependencyNode struct {
	Name         string   `hcl:"name" json:"name"`
	Dependencies []string `hcl:"dependencies" json:"dependencies"`
	Resolved     bool     `hcl:"resolved" json:"resolved"`
}
