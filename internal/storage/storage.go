package storage

import (
	"fmt"
	"time"

	spookyfactstypes "spooky/internal/types/facts"
	spookystoragebadger "spooky/internal/storage/badger"
)

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
