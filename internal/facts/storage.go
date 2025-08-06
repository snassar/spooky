package facts

import (
	"fmt"
	"io"
	"time"

	spookybadger "spooky/internal/facts/storage/badger"
	"spooky/internal/facts/types"
)

// FactStorage defines the interface for persistent fact storage
type FactStorage interface {
	// Core fact collection operations
	SetFactCollection(machineID string, collection *types.FactCollection) error
	GetFactCollection(machineID string) (*types.FactCollection, error)
	QueryFactCollections(query *FactQuery) ([]*types.FactCollection, error)
	DeleteFactCollection(machineID string) error
	DeleteFactCollections(query *FactQuery) (int, error)

	// Import/Export operations
	ExportToJSON(w io.Writer) error
	ImportFromJSON(r io.Reader) error
	ImportFromHCL(r io.Reader) error // NEW: HCL import
	ExportToJSONWithEncryption(w io.Writer, opts ExportOptions) error
	ImportFromJSONWithDecryption(r io.Reader, identityFile string) error

	// Close the storage connection
	Close() error
}

// ExportOptions defines options for encrypted export
type ExportOptions = types.ExportOptions

// TimeRange defines a time range for filtering
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// FactQuery defines query parameters for fact collections
type FactQuery = types.FactQuery

// matchesQuery checks if a fact collection matches the query criteria
func matchesQuery(collection *types.FactCollection, query *FactQuery) bool {
	if query.MachineName != "" && collection.Server != query.MachineName {
		return false
	}

	if query.MachineID != "" {
		if machineID, exists := collection.Facts["machine_id"]; exists {
			if id, ok := machineID.Value.(string); ok && id != query.MachineID {
				return false
			}
		} else {
			return false
		}
	}

	if query.OS != "" {
		// Check both possible OS fact keys
		osFact, exists := collection.Facts["system.os.name"]
		if !exists {
			osFact, exists = collection.Facts["os.name"]
		}
		if exists {
			if os, ok := osFact.Value.(string); ok && os != query.OS {
				return false
			}
		} else {
			return false
		}
	}

	// Check tag matching
	if len(query.Tags) > 0 {
		for tagKey, tagValue := range query.Tags {
			factKey := "tags." + tagKey
			if tagFact, exists := collection.Facts[factKey]; exists {
				if tag, ok := tagFact.Value.(string); ok && tag != tagValue {
					return false
				}
			} else {
				return false
			}
		}
	}

	// Check environment matching
	if query.Environment != "" {
		if envFact, exists := collection.Facts["tags.environment"]; exists {
			if env, ok := envFact.Value.(string); ok && env != query.Environment {
				return false
			}
		} else {
			return false
		}
	}

	if query.UpdatedBefore != nil && collection.Timestamp.After(*query.UpdatedBefore) {
		return false
	}

	if query.UpdatedAfter != nil && collection.Timestamp.Before(*query.UpdatedAfter) {
		return false
	}

	// TODO: Implement text search functionality
	return true
}

// StorageType defines the type of storage backend
type StorageType string

const (
	StorageTypeBadger StorageType = "badger"
	StorageTypeJSON   StorageType = "json"
	StorageTypeHCL    StorageType = "hcl"
)

// StorageOptions defines configuration for fact storage
type StorageOptions struct {
	Type StorageType
	Path string
}

// NewFactStorage creates a new fact storage instance (for writing/creating)
func NewFactStorage(opts StorageOptions) (FactStorage, error) {
	switch opts.Type {
	case StorageTypeJSON:
		return NewJSONFactStorage(opts.Path)
	case StorageTypeHCL:
		return NewHCLFactStorage(opts.Path)
	case StorageTypeBadger, "": // Default to BadgerDB
		return spookybadger.NewBadgerFactStorage(opts.Path)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", opts.Type)
	}
}

// OpenFactStorage opens an existing fact storage instance (for reading)
// This will fail if the storage doesn't exist
func OpenFactStorage(opts StorageOptions) (FactStorage, error) {
	switch opts.Type {
	case StorageTypeJSON:
		return NewJSONFactStorage(opts.Path)
	case StorageTypeHCL:
		return NewHCLFactStorage(opts.Path)
	case StorageTypeBadger, "": // Default to BadgerDB
		return spookybadger.NewBadgerFactStorageReadOnly(opts.Path)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", opts.Type)
	}
}

// NewFactStorageReadOnly creates a new fact storage instance for read-only operations
// This will fail if the database doesn't exist
// Deprecated: Use OpenFactStorage instead
func NewFactStorageReadOnly(opts StorageOptions) (FactStorage, error) {
	return OpenFactStorage(opts)
}
