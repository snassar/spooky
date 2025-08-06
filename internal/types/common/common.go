package common

import (
	"time"
)

// EncryptionMetadata contains information about encryption
// Used across secrets, facts, variables, and other encrypted data
type EncryptionMetadata struct {
	EncryptedAt       string   `json:"encrypted_at"`
	EncryptionVersion string   `json:"encryption_version"`
	Recipients        []string `json:"recipients"`
	Algorithm         string   `json:"algorithm,omitempty"`
}

// TimestampedEntity provides common timestamp fields
// Used across all entities that need creation/update tracking
type TimestampedEntity struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NamedEntity provides common name and description fields
// Used across actions, templates, variables, etc.
type NamedEntity struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// MetadataEntity provides common metadata fields
// Used across all entities that need additional metadata
type MetadataEntity struct {
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Tags     []string               `json:"tags,omitempty"`
}

// ValidationEntity provides common validation fields
// Used across configurable entities that need validation
type ValidationEntity struct {
	Validated   bool      `json:"validated"`
	ValidatedAt time.Time `json:"validated_at,omitempty"`
	Errors      []string  `json:"errors,omitempty"`
	Warnings    []string  `json:"warnings,omitempty"`
}

// StatusEntity provides common status tracking
// Used across entities that have operational status
type StatusEntity struct {
	Status   string    `json:"status"`
	StatusAt time.Time `json:"status_at"`
	Message  string    `json:"message,omitempty"`
	Details  string    `json:"details,omitempty"`
}

// CompleteEntity combines all common entity types
// Used for entities that need all common functionality
type CompleteEntity struct {
	NamedEntity
	TimestampedEntity
	MetadataEntity
	ValidationEntity
	StatusEntity
}

// ErrorDetails provides common error information
// Used across all error types for consistent error reporting
type ErrorDetails struct {
	Code        string                 `json:"code"`
	Message     string                 `json:"message"`
	Operation   string                 `json:"operation"`
	Context     map[string]interface{} `json:"context,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	Stack       []string               `json:"stack,omitempty"`
	Recoverable bool                   `json:"recoverable"`
}

// ExportOptions provides common export configuration
// Used across facts, variables, machines, etc. for data export
type ExportOptions struct {
	Format      string            `json:"format"`       // "json", "hcl", "yaml"
	Compress    bool              `json:"compress"`     // Whether to compress output
	IncludeMeta bool              `json:"include_meta"` // Include metadata
	Filter      map[string]string `json:"filter"`       // Filter criteria
	SortBy      string            `json:"sort_by"`      // Sort field
	SortOrder   string            `json:"sort_order"`   // "asc" or "desc"
	Limit       int               `json:"limit"`        // Limit results
	Offset      int               `json:"offset"`       // Offset for pagination
	Encrypt     bool              `json:"encrypt"`      // Whether to encrypt sensitive data
	Recipients  []string          `json:"recipients"`   // Encryption recipients
}

// ImportOptions provides common import configuration
// Used across facts, variables, machines, etc. for data import
type ImportOptions struct {
	Format     string            `json:"format"`      // "json", "hcl", "yaml"
	Validate   bool              `json:"validate"`    // Validate imported data
	Overwrite  bool              `json:"overwrite"`   // Overwrite existing data
	Merge      bool              `json:"merge"`       // Merge with existing data
	Decrypt    bool              `json:"decrypt"`     // Decrypt encrypted data
	Identity   string            `json:"identity"`    // Decryption identity
	Transform  map[string]string `json:"transform"`   // Field transformations
	SkipErrors bool              `json:"skip_errors"` // Skip import errors
	DryRun     bool              `json:"dry_run"`     // Preview import without applying
}

// TimeRange provides common time range filtering
// Used across facts, logs, actions, etc. for time-based queries
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Pagination provides common pagination support
// Used across all list operations for consistent pagination
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
	Pages    int `json:"pages"`
}

// Query provides common query parameters
// Used across all search and filter operations
type Query struct {
	Search     string            `json:"search"`     // Text search
	Filters    map[string]string `json:"filters"`    // Field filters
	SortBy     string            `json:"sort_by"`    // Sort field
	SortOrder  string            `json:"sort_order"` // Sort direction
	TimeRange  *TimeRange        `json:"time_range"` // Time filtering
	Pagination *Pagination       `json:"pagination"` // Pagination
	Limit      int               `json:"limit"`      // Result limit
	Offset     int               `json:"offset"`     // Result offset
}

// Result provides common result structure
// Used across all operations that return data with metadata
type Result struct {
	Data      interface{} `json:"data"`
	Count     int         `json:"count"`
	Total     int         `json:"total"`
	Page      int         `json:"page"`
	PageSize  int         `json:"page_size"`
	Pages     int         `json:"pages"`
	Query     *Query      `json:"query,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	Duration  string      `json:"duration"`
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Errors    []string    `json:"errors,omitempty"`
	Warnings  []string    `json:"warnings,omitempty"`
}
