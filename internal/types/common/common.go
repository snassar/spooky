// Package common provides foundational types used across multiple domains in the spooky codebase.
// These types serve as building blocks for domain-specific types and ensure consistency
// in data structures throughout the system.
package common

import (
	"time"
)

// TimestampedEntity provides creation and update tracking for entities
type TimestampedEntity struct {
	CreatedAt time.Time `json:"created_at" hcl:"created_at"`
	UpdatedAt time.Time `json:"updated_at" hcl:"updated_at"`
}

// NamedEntity provides name and description fields for entities
type NamedEntity struct {
	Name        string `json:"name" hcl:"name"`
	Description string `json:"description,omitempty" hcl:"description,optional"`
}

// MetadataEntity provides tags and metadata support for entities
type MetadataEntity struct {
	Tags        []string          `json:"tags,omitempty" hcl:"tags,optional"`
	Metadata    map[string]string `json:"metadata,omitempty" hcl:"metadata,optional"`
	Labels      map[string]string `json:"labels,omitempty" hcl:"labels,optional"`
	Annotations map[string]string `json:"annotations,omitempty" hcl:"annotations,optional"`
}

// ValidationEntity provides validation state and results for entities
type ValidationEntity struct {
	ValidatedAt time.Time         `json:"validated_at,omitempty" hcl:"validated_at,optional"`
	Valid       bool              `json:"valid" hcl:"valid"`
	Errors      []ValidationError `json:"errors,omitempty" hcl:"errors,optional"`
	Warnings    []ValidationError `json:"warnings,omitempty" hcl:"warnings,optional"`
}

// ValidationError represents a validation error with context
type ValidationError struct {
	Field   string      `json:"field" hcl:"field"`
	Message string      `json:"message" hcl:"message"`
	Code    string      `json:"code,omitempty" hcl:"code,optional"`
	Value   interface{} `json:"value,omitempty" hcl:"value,optional"`
}

// StatusEntity provides operational status tracking for entities
type StatusEntity struct {
	Status    string    `json:"status" hcl:"status"`
	StatusAt  time.Time `json:"status_at" hcl:"status_at"`
	StatusMsg string    `json:"status_msg,omitempty" hcl:"status_msg,optional"`
}

// CompleteEntity combines all common entity types for comprehensive entities
type CompleteEntity struct {
	TimestampedEntity
	NamedEntity
	MetadataEntity
	ValidationEntity
	StatusEntity
}

// ErrorDetails provides structured error information with context and stack traces
type ErrorDetails struct {
	Code        string                 `json:"code" hcl:"code"`
	Message     string                 `json:"message" hcl:"message"`
	Context     map[string]interface{} `json:"context,omitempty" hcl:"context,optional"`
	Stack       []string               `json:"stack,omitempty" hcl:"stack,optional"`
	Recoverable bool                   `json:"recoverable" hcl:"recoverable"`
	Timestamp   time.Time              `json:"timestamp" hcl:"timestamp"`
}

// ExportOptions provides configuration for data export operations
type ExportOptions struct {
	Format      string            `json:"format" hcl:"format"`
	Compression string            `json:"compression,omitempty" hcl:"compression,optional"`
	Encryption  EncryptionOptions `json:"encryption,omitempty" hcl:"encryption,optional"`
	Filter      map[string]string `json:"filter,omitempty" hcl:"filter,optional"`
	Include     []string          `json:"include,omitempty" hcl:"include,optional"`
	Exclude     []string          `json:"exclude,omitempty" hcl:"exclude,optional"`
}

// EncryptionOptions provides encryption configuration
type EncryptionOptions struct {
	Enabled    bool     `json:"enabled" hcl:"enabled"`
	Algorithm  string   `json:"algorithm" hcl:"algorithm"`
	Recipients []string `json:"recipients,omitempty" hcl:"recipients,optional"`
	KeyFile    string   `json:"key_file,omitempty" hcl:"key_file,optional"`
}

// ImportOptions provides configuration for data import operations
type ImportOptions struct {
	Format      string            `json:"format" hcl:"format"`
	Compression string            `json:"compression,omitempty" hcl:"compression,optional"`
	Encryption  EncryptionOptions `json:"encryption,omitempty" hcl:"encryption,optional"`
	Validate    bool              `json:"validate" hcl:"validate"`
	Overwrite   bool              `json:"overwrite" hcl:"overwrite"`
	DryRun      bool              `json:"dry_run" hcl:"dry_run"`
}

// Query provides standardized query structure for data operations
type Query struct {
	Filters    map[string]interface{} `json:"filters,omitempty" hcl:"filters,optional"`
	Sort       []SortField            `json:"sort,omitempty" hcl:"sort,optional"`
	Pagination Pagination             `json:"pagination,omitempty" hcl:"pagination,optional"`
	Fields     []string               `json:"fields,omitempty" hcl:"fields,optional"`
}

// SortField represents a sort field with direction
type SortField struct {
	Field     string `json:"field" hcl:"field"`
	Direction string `json:"direction" hcl:"direction"` // "asc" or "desc"
}

// Result provides standardized result structure for data operations
type Result struct {
	Data       interface{}            `json:"data" hcl:"data"`
	Total      int64                  `json:"total" hcl:"total"`
	Pagination Pagination             `json:"pagination,omitempty" hcl:"pagination,optional"`
	Metadata   map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// TimeRange provides time-based range for queries
type TimeRange struct {
	Start time.Time `json:"start" hcl:"start"`
	End   time.Time `json:"end" hcl:"end"`
}

// Pagination provides pagination information for queries
type Pagination struct {
	Page     int   `json:"page" hcl:"page"`
	PageSize int   `json:"page_size" hcl:"page_size"`
	Total    int64 `json:"total" hcl:"total"`
}

// EncryptionMetadata provides cross-domain encryption information and recipient management
type EncryptionMetadata struct {
	Algorithm  string                 `json:"algorithm" hcl:"algorithm"`
	Recipients []string               `json:"recipients" hcl:"recipients"`
	KeyID      string                 `json:"key_id,omitempty" hcl:"key_id,optional"`
	Metadata   map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
	CreatedAt  time.Time              `json:"created_at" hcl:"created_at"`
}
