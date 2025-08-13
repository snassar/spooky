// Package templates provides template-related type definitions.
package templates

import (
	"context"
	"time"
)

// Template represents a template with metadata and configuration
type Template struct {
	// Template identification
	ID string `json:"id" hcl:"id"`

	// Template source location
	SourcePath string `json:"source_path" hcl:"source_path"`

	// Template destination
	DestinationPath string `json:"destination_path,omitempty" hcl:"destination_path,optional"`

	// Template type classification
	Type string `json:"type" hcl:"type"`

	// Template scope
	Scope string `json:"scope" hcl:"scope"`

	// Template security level
	SecurityLevel string `json:"security_level" hcl:"security_level"`

	// Template rendering engine
	Engine string `json:"engine" hcl:"engine"`

	// Template variables
	Variables map[string]interface{} `json:"variables,omitempty" hcl:"variables,optional"`

	// Template context data
	ContextData map[string]interface{} `json:"context_data,omitempty" hcl:"context_data,optional"`

	// Template functions and restrictions
	Functions map[string]interface{} `json:"functions,omitempty" hcl:"functions,optional"`

	// Template metadata
	Metadata *TemplateMetadata `json:"metadata,omitempty" hcl:"metadata,optional"`

	// Template content
	Content string `json:"content,omitempty" hcl:"content,optional"`

	// Creation timestamp
	CreatedAt time.Time `json:"created_at" hcl:"created_at"`

	// Last update timestamp
	UpdatedAt time.Time `json:"updated_at" hcl:"updated_at"`
}

// TemplateMetadata represents template metadata and information
type TemplateMetadata struct {
	Name        string   `json:"name,omitempty" hcl:"name,optional"`
	Description string   `json:"description,omitempty" hcl:"description,optional"`
	Author      string   `json:"author,omitempty" hcl:"author,optional"`
	Version     string   `json:"version,omitempty" hcl:"version,optional"`
	Tags        []string `json:"tags,omitempty" hcl:"tags,optional"`
	License     string   `json:"license,omitempty" hcl:"license,optional"`
}

// TemplateContext represents the context available to templates during rendering
type TemplateContext struct {
	// Project information
	Project map[string]interface{} `json:"project,omitempty"`

	// Machine facts
	Facts map[string]interface{} `json:"facts,omitempty"`

	// Inventory information
	Machines []map[string]interface{} `json:"machines,omitempty"`

	// Environment variables
	Environment map[string]string `json:"environment,omitempty"`

	// Custom data
	CustomData map[string]interface{} `json:"custom_data,omitempty"`

	// Variables
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// TemplateEngine represents a template rendering engine
type TemplateEngine interface {
	// RenderTemplate renders a template with the given context
	RenderTemplate(ctx context.Context, template *Template, context *TemplateContext) (string, error)

	// ValidateTemplate validates a template
	ValidateTemplate(ctx context.Context, template *Template) error

	// GetTemplateFunctions returns available template functions
	GetTemplateFunctions() map[string]interface{}
}

// TemplateManager manages template operations
type TemplateManager interface {
	// LoadTemplate loads a template from the given path
	LoadTemplate(ctx context.Context, templatePath string) (*Template, error)

	// SaveTemplate saves a template to the given path
	SaveTemplate(ctx context.Context, template *Template, templatePath string) error

	// RenderTemplate renders a template with the given context
	RenderTemplate(ctx context.Context, template *Template, context *TemplateContext) (string, error)

	// ValidateTemplate validates a template
	ValidateTemplate(ctx context.Context, template *Template) error

	// ListTemplates lists all templates in the given directory
	ListTemplates(ctx context.Context, directory string) ([]*Template, error)
}

// TemplateError represents template-related errors
type TemplateError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
	Line    int    `json:"line,omitempty"`
	Column  int    `json:"column,omitempty"`
}

// Error returns the error message
func (e *TemplateError) Error() string {
	return e.Message
}

// TemplateValidationResult represents template validation results
type TemplateValidationResult struct {
	Valid    bool             `json:"valid"`
	Errors   []*TemplateError `json:"errors,omitempty"`
	Warnings []*TemplateError `json:"warnings,omitempty"`
}
