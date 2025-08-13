// Package facts provides fact collection, storage, and management functionality.
package facts

import (
	"context"
	"time"

	spookytypes "spooky/internal/types"
	spookytypesfacts "spooky/internal/types/facts"
)

// FactCollection represents a collection of facts for a machine
type FactCollection struct {
	// Machine ID (32-character hex string from /etc/machine-id)
	MachineID string `json:"machine_id" hcl:"machine_id"`

	// Collection timestamp
	CollectedAt time.Time `json:"collected_at" hcl:"collected_at"`

	// Collection of facts for this machine
	Facts *spookytypesfacts.Facts `json:"facts" hcl:"facts"`

	// Metadata about the collection
	Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// FactCollector collects facts from a machine
type FactCollector interface {
	// Collect collects facts from the given machine
	Collect(ctx context.Context, machine *spookytypes.Machine) (*FactCollection, error)

	// GetName returns the collector name
	GetName() string
}

// FactManager manages fact collection and export
type FactManager interface {
	// CollectFacts collects facts from the given machine
	CollectFacts(ctx context.Context, machine *spookytypes.Machine) (*FactCollection, error)

	// GetFacts retrieves facts for a specific machine (collects on demand)
	GetFacts(ctx context.Context, machineID string) (*FactCollection, error)

	// ValidateFacts validates facts against schema
	ValidateFacts(ctx context.Context, facts *FactCollection) (*spookytypes.ValidationResult, error)

	// ExportFacts exports facts to the given format
	ExportFacts(ctx context.Context, machineIDs []string, format string, outputPath string) error

	// GetStorageStats returns storage statistics for debugging
	GetStorageStats() (map[string]interface{}, error)

	// ClearFacts removes all facts from memory
	ClearFacts(ctx context.Context) error
}

// FactCollectionOptions provides options for fact collection
type FactCollectionOptions struct {
	// Timeout for fact collection
	Timeout time.Duration

	// Parallel workers for collection
	ParallelWorkers int

	// Retry attempts
	RetryAttempts int

	// Retry delay
	RetryDelay time.Duration

	// Include system facts
	IncludeSystem bool

	// Include enhanced facts
	IncludeEnhanced bool

	// Include application facts
	IncludeApplications bool

	// Include custom facts
	IncludeCustom bool
}

// DefaultFactCollectionOptions returns default fact collection options
func DefaultFactCollectionOptions() *FactCollectionOptions {
	return &FactCollectionOptions{
		Timeout:             30 * time.Second,
		ParallelWorkers:     4,
		RetryAttempts:       3,
		RetryDelay:          5 * time.Second,
		IncludeSystem:       true,
		IncludeEnhanced:     true,
		IncludeApplications: false,
		IncludeCustom:       false,
	}
}

// FactExportOptions provides options for fact export
type FactExportOptions struct {
	// Export format (json, hcl)
	Format string

	// Output path
	OutputPath string

	// Include metadata
	IncludeMetadata bool

	// Pretty print
	PrettyPrint bool

	// Compress output
	Compress bool
}

// DefaultFactExportOptions returns default fact export options
func DefaultFactExportOptions() *FactExportOptions {
	return &FactExportOptions{
		Format:          "json",
		OutputPath:      "",
		IncludeMetadata: true,
		PrettyPrint:     true,
		Compress:        false,
	}
}

// FactImportOptions provides options for fact import
type FactImportOptions struct {
	// Import format (json, hcl)
	Format string

	// Input path
	InputPath string

	// Overwrite existing facts
	Overwrite bool

	// Validate on import
	Validate bool
}

// DefaultFactImportOptions returns default fact import options
func DefaultFactImportOptions() *FactImportOptions {
	return &FactImportOptions{
		Format:    "json",
		InputPath: "",
		Overwrite: false,
		Validate:  true,
	}
}
