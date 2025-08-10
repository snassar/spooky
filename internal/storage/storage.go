package storage

import (
	"fmt"
	"io"
	"time"

	spookystoragebadger "spooky/internal/storage/badger"
	spookyfactstypes "spooky/internal/types/facts"
)

// StorageType defines the type of storage backend
type StorageType string

const (
	StorageTypeBadger StorageType = "badger"
	StorageTypeJSON   StorageType = "json"
	StorageTypeHCL    StorageType = "hcl"
)

// StorageOptions defines configuration options for storage backends
type StorageOptions struct {
	Type           StorageType
	Path           string
	EncryptEnabled bool
	CryptoManager  interface{} // Will be properly typed when secrets package is available
	DefaultTTL     time.Duration
}

// StorageBackend defines the basic storage operations
type StorageBackend interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	Close() error
}

// ExportableStorage defines operations for data export and import
type ExportableStorage interface {
	ExportToJSON(w io.Writer) error
	ExportToHCL(w io.Writer) error
	ImportFromJSON(r io.Reader) error
	ImportFromHCL(r io.Reader) error
	ExportToJSONWithEncryption(w io.Writer, options spookyfactstypes.ExportOptions) error
	ImportFromJSONWithDecryption(r io.Reader, identityFile string) error
}

// FactStorage defines the interface for fact-specific storage operations
type FactStorage interface {
	StorageBackend
	ExportableStorage

	// Fact-specific operations
	SetFactCollection(machineID string, collection *spookyfactstypes.FactCollection) error
	GetFactCollection(machineID string) (*spookyfactstypes.FactCollection, error)
	QueryFactCollections(query *spookyfactstypes.FactQuery) ([]*spookyfactstypes.FactCollection, error)
	DeleteFactCollection(machineID string) error
	DeleteFactCollections(query *spookyfactstypes.FactQuery) (int, error)
}

// ExportOptions defines options for encrypted export
type ExportOptions = spookyfactstypes.ExportOptions

// TimeRange defines a time range for filtering
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// FactQuery defines query parameters for fact collections
type FactQuery = spookyfactstypes.FactQuery

// matchesQuery checks if a fact collection matches the query criteria
func matchesQuery(collection *spookyfactstypes.FactCollection, query *spookyfactstypes.FactQuery) bool {
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

	return true
}

// NewFactStorage creates a new fact storage instance (for writing/creating)
func NewFactStorage(opts StorageOptions) (FactStorage, error) {
	switch opts.Type {
	case StorageTypeBadger, "": // Default to BadgerDB
		return spookystoragebadger.NewBadgerFactStorage(opts.Path)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", opts.Type)
	}
}

// OpenFactStorage opens an existing fact storage instance (for reading)
// This will fail if the storage doesn't exist
func OpenFactStorage(opts StorageOptions) (FactStorage, error) {
	switch opts.Type {
	case StorageTypeBadger, "": // Default to BadgerDB
		return spookystoragebadger.NewBadgerFactStorageReadOnly(opts.Path)
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
