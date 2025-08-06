package badger

import (
	spookyfactstypes "spooky/internal/facts/types"
)

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

	// TODO: Implement text search functionality
	return true
}
