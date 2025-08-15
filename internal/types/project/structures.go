// Package project provides types for project management in the spooky codebase.
// These types define the structure for spooky projects and their configuration.
package project

import (
	"time"
)

// Project represents a spooky project with its configuration and metadata
type Project struct {
	// Basic entity fields
	Name        string    `json:"name" hcl:"name"`
	Description string    `json:"description,omitempty" hcl:"description,optional"`
	CreatedAt   time.Time `json:"created_at" hcl:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" hcl:"updated_at"`

	// Project path on the filesystem
	Path string `json:"path" hcl:"path"`

	// Project configuration
	Config *Config `json:"config" hcl:"config"`

	// Project metadata
	Metadata *Metadata `json:"metadata" hcl:"metadata"`

	// Project settings
	Settings *Settings `json:"settings" hcl:"settings"`

	// Project validation state
	Validated   bool      `json:"validated" hcl:"validated"`
	ValidatedAt time.Time `json:"validated_at,omitempty" hcl:"validated_at,optional"`
}

// Config represents the project.hcl configuration
type Config struct {
	// Project name
	Name string `json:"name" hcl:"name"`

	// Project description
	Description string `json:"description,omitempty" hcl:"description,optional"`

	// Project metadata
	Metadata *Metadata `json:"metadata,omitempty" hcl:"metadata,optional"`

	// Project settings
	Settings *Settings `json:"settings,omitempty" hcl:"settings,optional"`
}

// Metadata provides metadata for the project
type Metadata struct {
	// Project version
	Version string `json:"version,omitempty" hcl:"version,optional"`

	// Project author
	Author string `json:"author,omitempty" hcl:"author,optional"`

	// Project tags
	Tags []string `json:"tags,omitempty" hcl:"tags,optional"`

	// Project URL
	URL string `json:"url,omitempty" hcl:"url,optional"`

	// Project license
	License string `json:"license,omitempty" hcl:"license,optional"`

	// Project creation date
	CreatedAt time.Time `json:"created_at,omitempty" hcl:"created_at,optional"`

	// Project last modified date
	ModifiedAt time.Time `json:"modified_at,omitempty" hcl:"modified_at,optional"`
}

// Settings provides configuration settings for the project
type Settings struct {
	// Parallel workers for operations
	ParallelWorkers int `json:"parallel_workers,omitempty" hcl:"parallel_workers,optional"`

	// Timeout for operations in seconds
	TimeoutSeconds int `json:"timeout_seconds,omitempty" hcl:"timeout_seconds,optional"`

	// Log level for project operations
	LogLevel string `json:"log_level,omitempty" hcl:"log_level,optional"`

	// Whether to enable dry-run mode by default
	DefaultDryRun bool `json:"default_dry_run,omitempty" hcl:"default_dry_run,optional"`

	// Whether to enable verbose output by default
	DefaultVerbose bool `json:"default_verbose,omitempty" hcl:"default_verbose,optional"`

	// Default output format
	DefaultFormat string `json:"default_format,omitempty" hcl:"default_format,optional"`

	// Whether to validate before operations
	ValidateBeforeRun bool `json:"validate_before_run,omitempty" hcl:"validate_before_run,optional"`

	// Whether to allow failures without stopping
	AllowFailures bool `json:"allow_failures,omitempty" hcl:"allow_failures,optional"`

	// Maximum retry attempts
	MaxRetries int `json:"max_retries,omitempty" hcl:"max_retries,optional"`

	// Retry delay in seconds
	RetryDelaySeconds int `json:"retry_delay_seconds,omitempty" hcl:"retry_delay_seconds,optional"`
}

// Directory represents the project directory structure
type Directory struct {
	// Project root path
	Root string `json:"root" hcl:"root"`

	// Required directories
	Directories []string `json:"directories" hcl:"directories"`

	// Required files
	Files []string `json:"files" hcl:"files"`

	// Optional directories
	OptionalDirectories []string `json:"optional_directories,omitempty" hcl:"optional_directories,optional"`

	// Optional files
	OptionalFiles []string `json:"optional_files,omitempty" hcl:"optional_files,optional"`
}

// Validation represents project validation results
type Validation struct {
	// Whether the project is valid
	Valid bool `json:"valid" hcl:"valid"`

	// Validation timestamp
	ValidatedAt time.Time `json:"validated_at" hcl:"validated_at"`

	// Validation errors
	Errors []Error `json:"errors,omitempty" hcl:"errors,optional"`

	// Validation warnings
	Warnings []Error `json:"warnings,omitempty" hcl:"warnings,optional"`

	// Validation details
	Details map[string]interface{} `json:"details,omitempty" hcl:"details,optional"`
}

// Error represents a project-related error
type Error struct {
	// Error details
	Code        string                 `json:"code" hcl:"code"`
	Message     string                 `json:"message" hcl:"message"`
	Context     map[string]interface{} `json:"context,omitempty" hcl:"context,optional"`
	Stack       []string               `json:"stack,omitempty" hcl:"stack,optional"`
	Recoverable bool                   `json:"recoverable" hcl:"recoverable"`

	// Project path where the error occurred
	ProjectPath string `json:"project_path" hcl:"project_path"`

	// File path where the error occurred (if applicable)
	FilePath string `json:"file_path,omitempty" hcl:"file_path,optional"`

	// Line number where the error occurred (if applicable)
	LineNumber int `json:"line_number,omitempty" hcl:"line_number,optional"`

	// Column number where the error occurred (if applicable)
	ColumnNumber int `json:"column_number,omitempty" hcl:"column_number,optional"`

	// Error severity
	Severity string `json:"severity" hcl:"severity"` // "error", "warning", "info"
}

// NewError creates a new project error
func NewError(projectPath, message, severity string) *Error {
	return &Error{
		Code:        "project_error",
		Message:     message,
		Recoverable: severity != "error",
		ProjectPath: projectPath,
		Severity:    severity,
	}
}

// Error implements the error interface
func (e *Error) Error() string {
	return e.Message
}

// Unwrap returns the underlying error
func (e *Error) Unwrap() error {
	return nil
}

// ProjectInfo provides information about a project
type Info struct {
	// Project name
	Name string `json:"name" hcl:"name"`

	// Project description
	Description string `json:"description,omitempty" hcl:"description,optional"`

	// Project path
	Path string `json:"path" hcl:"path"`

	// Project version
	Version string `json:"version,omitempty" hcl:"version,optional"`

	// Project author
	Author string `json:"author,omitempty" hcl:"author,optional"`

	// Project tags
	Tags []string `json:"tags,omitempty" hcl:"tags,optional"`

	// Project creation date
	CreatedAt time.Time `json:"created_at,omitempty" hcl:"created_at,optional"`

	// Project last modified date
	ModifiedAt time.Time `json:"modified_at,omitempty" hcl:"modified_at,optional"`

	// Project validation status
	Valid       bool      `json:"valid" hcl:"valid"`
	ValidatedAt time.Time `json:"validated_at,omitempty" hcl:"validated_at,optional"`

	// Project statistics
	Statistics *Statistics `json:"statistics,omitempty" hcl:"statistics,optional"`
}

// Statistics provides statistics about the project
type Statistics struct {
	// Number of machines in the project
	MachineCount int `json:"machine_count" hcl:"machine_count"`

	// Number of actions in the project
	ActionCount int `json:"action_count" hcl:"action_count"`

	// Number of templates in the project
	TemplateCount int `json:"template_count" hcl:"template_count"`

	// Number of variables in the project
	VariableCount int `json:"variable_count" hcl:"variable_count"`

	// Total project size in bytes
	TotalSize int64 `json:"total_size" hcl:"total_size"`

	// Last facts collection timestamp
	LastFactsCollection time.Time `json:"last_facts_collection,omitempty" hcl:"last_facts_collection,optional"`
}
