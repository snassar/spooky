package types

import (
	"time"
)

// Template represents a template with metadata and configuration
// Aligns with template-metadata.hcl schema
type Template struct {
	// Template metadata (from template-metadata.hcl)
	Name        string             `hcl:"name,label" json:"name"`
	Source      string             `hcl:"source" json:"source"`
	Destination string             `hcl:"destination,optional" json:"destination,omitempty"`
	Variables   []TemplateVariable `hcl:"variables,optional" json:"variables,omitempty"`
	Metadata    *TemplateMetadata  `hcl:"metadata,optional" json:"metadata,omitempty"`

	// Additional fields for internal use
	Content   string    `hcl:"content,optional" json:"content,omitempty"`
	Path      string    `hcl:"path,optional" json:"path,omitempty"`
	CreatedAt time.Time `hcl:"created_at,optional" json:"created_at,omitempty"`
	UpdatedAt time.Time `hcl:"updated_at,optional" json:"updated_at,omitempty"`
}

// TemplateVariable represents a variable required by a template
// Aligns with template-metadata.hcl variables definition
type TemplateVariable struct {
	Name        string      `hcl:"name" json:"name"`
	Type        string      `hcl:"type" json:"type"` // string, number, boolean, list, object
	Required    bool        `hcl:"required,optional" json:"required,omitempty"`
	Default     interface{} `hcl:"default,optional" json:"default,omitempty"`
	Description string      `hcl:"description,optional" json:"description,omitempty"`
}

// TemplateMetadata represents template metadata
// Aligns with template-metadata.hcl metadata definition
type TemplateMetadata struct {
	Description  string   `hcl:"description,optional" json:"description,omitempty"`
	Author       string   `hcl:"author,optional" json:"author,omitempty"`
	Version      string   `hcl:"version,optional" json:"version,omitempty"`
	Tags         []string `hcl:"tags,optional" json:"tags,omitempty"`
	OutputFormat string   `hcl:"output_format,optional" json:"output_format,omitempty"` // config, script, documentation, data, other
}
