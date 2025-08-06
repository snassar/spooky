package types

import (
	"time"
)

// ExportOptions defines options for encrypted export
type ExportOptions struct {
	EncryptSensitive   bool
	EncryptionKey      string
	CompressData       bool
	IncludeMetadata    bool
	FilterByMachineID  []string
	FilterByTags       map[string]string
	FilterByTimeRange  *TimeRange
	ExportFormat       string // "json", "hcl"
	IncludeSystemFacts bool
	IncludeCustomFacts bool
	Decrypt            bool // Decrypt encrypted facts during export
}

// TimeRange defines a time range for filtering
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// FactQuery defines query parameters for fact collections
type FactQuery struct {
	MachineName   string            // Query by machine name
	MachineID     string            // Query by machine ID
	Tags          map[string]string // Query by tags
	OS            string            // Query by OS
	Environment   string            // Query by environment
	Limit         int               // Limit results
	SearchQuery   string            // Text search in facts
	SearchField   string            // Specific field to search
	UpdatedBefore *time.Time        // Filter by collection time
	UpdatedAfter  *time.Time        // Filter by collection time
}

// StorageType defines the type of storage backend
type StorageType string

const (
	StorageTypeBadger StorageType = "badger"
	StorageTypeJSON   StorageType = "json"
	StorageTypeHCL    StorageType = "hcl"
)

// StorageOptions defines options for storage initialization
type StorageOptions struct {
	Type StorageType
	Path string
}
