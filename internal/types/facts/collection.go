// Package facts provides fact data structures and types.
package facts

import (
	"context"
	"time"
)

// FactCollection represents a collection of facts for a machine
type FactCollection struct {
	// Machine ID (32-character hex string from /etc/machine-id)
	MachineID string `json:"machine_id" hcl:"machine_id"`

	// Collection timestamp
	CollectedAt time.Time `json:"collected_at" hcl:"collected_at"`

	// Collection of facts for this machine
	Facts *Facts `json:"facts" hcl:"facts"`

	// Metadata about the collection
	Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// FactStorage provides storage operations for fact collections
type FactStorage interface {
	// Store stores facts for a machine
	Store(ctx context.Context, machineID string, facts *FactCollection) error

	// Get retrieves facts for a machine
	Get(ctx context.Context, machineID string) (*FactCollection, error)

	// List lists all machine IDs with stored facts
	List(ctx context.Context) ([]string, error)

	// Clear removes all facts from storage
	Clear(ctx context.Context) error

	// GetStats returns storage statistics for debugging
	GetStats() (map[string]interface{}, error)
}

// FactCollector collects facts from a machine
type FactCollector interface {
	// Collect collects facts from the given machine
	Collect(ctx context.Context, machine interface{}) (*FactCollection, error)

	// GetName returns the collector name
	GetName() string
}

// FactManager manages fact collection and storage
type FactManager interface {
	// CollectFacts collects facts from the given machine
	CollectFacts(ctx context.Context, machine interface{}) (*FactCollection, error)

	// StoreFacts stores facts for a machine
	StoreFacts(ctx context.Context, machineID string, facts *FactCollection) error

	// GetFacts retrieves facts for a machine
	GetFacts(ctx context.Context, machineID string) (*FactCollection, error)

	// ListFacts lists all machines with stored facts
	ListFacts(ctx context.Context) ([]string, error)

	// ClearFacts removes all facts from storage
	ClearFacts(ctx context.Context) error

	// ValidateFacts validates facts against schema
	ValidateFacts(ctx context.Context, facts *FactCollection) (interface{}, error)

	// ExportFacts exports facts to the given format
	ExportFacts(ctx context.Context, machineIDs []string, format string, outputPath string) error

	// ImportFacts imports facts from the given format
	ImportFacts(ctx context.Context, format string, inputPath string) error
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

// FactExportOptions provides options for fact export
type FactExportOptions struct {
	// Export format (json, hcl)
	Format string

	// Output file path
	OutputPath string

	// Include metadata
	IncludeMetadata bool

	// Include timestamps
	IncludeTimestamps bool

	// Filter by machine IDs
	MachineIDs []string

	// Filter by fact types
	FactTypes []string
}
