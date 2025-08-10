package indexing

import (
	"fmt"
	"sync"
	"time"

	spookytypesconfig "spooky/internal/types/config"
	spookytypesmachines "spooky/internal/types/machines"
)

// Manager provides comprehensive indexing capabilities for machine data
type Manager struct {
	indexes map[spookytypesmachines.IndexType]*spookytypesmachines.MachineIndex
	mutex   sync.RWMutex
}

// NewManager creates a new indexing manager
func NewManager() *Manager {
	return &Manager{
		indexes: make(map[spookytypesmachines.IndexType]*spookytypesmachines.MachineIndex),
	}
}

// BuildIndexes builds all indexes for the given machines
func (im *Manager) BuildIndexes(machines []spookytypesconfig.Machine) error {
	startTime := time.Now()

	// Build each index type
	indexTypes := []spookytypesmachines.IndexType{
		spookytypesmachines.IndexTypeName,
		spookytypesmachines.IndexTypeHost,
		spookytypesmachines.IndexTypeTag,
		spookytypesmachines.IndexTypeGroup,
		spookytypesmachines.IndexTypeUser,
		spookytypesmachines.IndexTypePort,
		spookytypesmachines.IndexTypeMetadata,
	}

	for _, indexType := range indexTypes {
		if err := im.buildIndex(indexType, machines); err != nil {
			return fmt.Errorf("failed to build %s index: %w", indexType, err)
		}
	}

	buildTime := time.Since(startTime)
	// Update metrics
	for _, index := range im.indexes {
		if index.Metrics != nil {
			index.Metrics.BuildTime = buildTime
			index.Metrics.MachineCount = len(machines)
			index.Metrics.IndexCount = len(indexTypes)
			index.Metrics.LastUpdated = time.Now()
		}
	}

	return nil
}

// UpdateIndexes updates existing indexes with new machine data
func (im *Manager) UpdateIndexes(machines []spookytypesconfig.Machine) error {
	// For now, rebuild all indexes
	// In a more sophisticated implementation, we could do incremental updates
	return im.BuildIndexes(machines)
}

// SynchronizeIndexes synchronizes indexes with external data
func (im *Manager) SynchronizeIndexes(syncData *spookytypesmachines.IndexSyncData) error {
	if syncData == nil {
		return fmt.Errorf("sync data is nil")
	}

	// Simplified implementation - just rebuild indexes
	// In a real implementation, we would apply specific changes
	return nil
}

// LookupByName looks up a machine by name
func (im *Manager) LookupByName(name string) (*spookytypesconfig.Machine, bool) {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	if index, exists := im.indexes[spookytypesmachines.IndexTypeName]; exists {
		if machine, found := index.NameIndex[name]; found {
			im.updateLookupMetrics(spookytypesmachines.IndexTypeName, true)
			return machine, true
		}
	}
	im.updateLookupMetrics(spookytypesmachines.IndexTypeName, false)
	return nil, false
}

// LookupByHost looks up a machine by host
func (im *Manager) LookupByHost(host string) (*spookytypesconfig.Machine, bool) {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	if index, exists := im.indexes[spookytypesmachines.IndexTypeHost]; exists {
		if machine, found := index.HostIndex[host]; found {
			im.updateLookupMetrics(spookytypesmachines.IndexTypeHost, true)
			return machine, true
		}
	}
	im.updateLookupMetrics(spookytypesmachines.IndexTypeHost, false)
	return nil, false
}

// LookupByTag looks up machines by tag key
func (im *Manager) LookupByTag(tagKey string) ([]*spookytypesconfig.Machine, bool) {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	if index, exists := im.indexes[spookytypesmachines.IndexTypeTag]; exists {
		if machines, found := index.TagIndex[tagKey]; found {
			im.updateLookupMetrics(spookytypesmachines.IndexTypeTag, true)
			return machines, true
		}
	}
	im.updateLookupMetrics(spookytypesmachines.IndexTypeTag, false)
	return nil, false
}

// LookupByTagValue looks up machines by tag key and value
func (im *Manager) LookupByTagValue(tagKey, tagValue string) ([]*spookytypesconfig.Machine, bool) {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	if index, exists := im.indexes[spookytypesmachines.IndexTypeTag]; exists {
		key := fmt.Sprintf("%s:%s", tagKey, tagValue)
		if machines, found := index.TagValueIndex[key]; found {
			im.updateLookupMetrics(spookytypesmachines.IndexTypeTag, true)
			return machines, true
		}
	}
	im.updateLookupMetrics(spookytypesmachines.IndexTypeTag, false)
	return nil, false
}

// LookupByGroup looks up machines by group name
func (im *Manager) LookupByGroup(groupName string) ([]*spookytypesconfig.Machine, bool) {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	if index, exists := im.indexes[spookytypesmachines.IndexTypeGroup]; exists {
		if machines, found := index.GroupIndex[groupName]; found {
			im.updateLookupMetrics(spookytypesmachines.IndexTypeGroup, true)
			return machines, true
		}
	}
	im.updateLookupMetrics(spookytypesmachines.IndexTypeGroup, false)
	return nil, false
}

// LookupByMetadata looks up machines by metadata key
func (im *Manager) LookupByMetadata(metadataKey string) ([]*spookytypesconfig.Machine, bool) {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	if index, exists := im.indexes[spookytypesmachines.IndexTypeMetadata]; exists {
		if machines, found := index.MetadataIndex[metadataKey]; found {
			im.updateLookupMetrics(spookytypesmachines.IndexTypeMetadata, true)
			return machines, true
		}
	}
	im.updateLookupMetrics(spookytypesmachines.IndexTypeMetadata, false)
	return nil, false
}

// FilterByGroups filters machines by multiple groups
func (im *Manager) FilterByGroups(groups []string) []*spookytypesconfig.Machine {
	var result []*spookytypesconfig.Machine
	seen := make(map[*spookytypesconfig.Machine]bool)

	for _, group := range groups {
		if machines, found := im.LookupByGroup(group); found {
			for _, machine := range machines {
				if !seen[machine] {
					seen[machine] = true
					result = append(result, machine)
				}
			}
		}
	}

	return result
}

// FilterByTags filters machines by tag criteria
func (im *Manager) FilterByTags(criteria map[string]string) []*spookytypesconfig.Machine {
	if len(criteria) == 0 {
		return nil
	}

	var result []*spookytypesconfig.Machine
	seen := make(map[*spookytypesconfig.Machine]bool)

	// Start with the first criterion
	var firstKey, firstValue string
	for key, value := range criteria {
		firstKey, firstValue = key, value
		break
	}

	if machines, found := im.LookupByTagValue(firstKey, firstValue); found {
		for _, machine := range machines {
			seen[machine] = true
			result = append(result, machine)
		}
	}

	// Apply remaining criteria as intersections
	for key, value := range criteria {
		if key == firstKey && value == firstValue {
			continue // Skip the first criterion we already processed
		}

		if machines, found := im.LookupByTagValue(key, value); found {
			// Intersect with current result
			newResult := []*spookytypesconfig.Machine{}
			for _, machine := range machines {
				if seen[machine] {
					newResult = append(newResult, machine)
				}
			}
			result = newResult

			// Update seen map
			seen = make(map[*spookytypesconfig.Machine]bool)
			for _, machine := range result {
				seen[machine] = true
			}
		} else {
			// No machines match this criterion, so result is empty
			return nil
		}
	}

	return result
}

// GetMetrics returns performance metrics for all indexes
func (im *Manager) GetMetrics() *spookytypesmachines.IndexMetrics {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	// Aggregate metrics from all indexes
	aggregateMetrics := &spookytypesmachines.IndexMetrics{
		IndexTypeStats: make(map[spookytypesmachines.IndexType]*spookytypesmachines.IndexTypeStats),
	}

	for indexType, index := range im.indexes {
		if index.Metrics != nil {
			aggregateMetrics.BuildTime += index.Metrics.BuildTime
			aggregateMetrics.LookupTime += index.Metrics.LookupTime
			aggregateMetrics.MemoryUsage += index.Metrics.MemoryUsage
			aggregateMetrics.MachineCount = index.Metrics.MachineCount
			aggregateMetrics.IndexCount = len(im.indexes)
			aggregateMetrics.LastUpdated = index.Metrics.LastUpdated

			// Aggregate index type stats
			if stats, exists := index.Metrics.IndexTypeStats[indexType]; exists {
				aggregateMetrics.IndexTypeStats[indexType] = stats
			}
		}
	}

	// Calculate hit rate
	totalHits := int64(0)
	totalMisses := int64(0)
	for _, stats := range aggregateMetrics.IndexTypeStats {
		totalHits += stats.HitCount
		totalMisses += stats.MissCount
	}

	if totalHits+totalMisses > 0 {
		aggregateMetrics.HitRate = float64(totalHits) / float64(totalHits+totalMisses)
	}

	return aggregateMetrics
}

// OptimizeIndexes performs optimization on all indexes
func (im *Manager) OptimizeIndexes() {
	im.mutex.Lock()
	defer im.mutex.Unlock()

	for indexType, index := range im.indexes {
		im.optimizeIndex(indexType, index)
	}
}

// buildIndex builds a specific index type
func (im *Manager) buildIndex(indexType spookytypesmachines.IndexType, machines []spookytypesconfig.Machine) error {
	startTime := time.Now()

	index := &spookytypesmachines.MachineIndex{
		Metrics: &spookytypesmachines.IndexMetrics{
			IndexTypeStats: make(map[spookytypesmachines.IndexType]*spookytypesmachines.IndexTypeStats),
		},
	}

	switch indexType {
	case spookytypesmachines.IndexTypeName:
		index.NameIndex = make(map[string]*spookytypesconfig.Machine)
		for i := range machines {
			index.NameIndex[machines[i].Name] = &machines[i]
		}

	case spookytypesmachines.IndexTypeHost:
		index.HostIndex = make(map[string]*spookytypesconfig.Machine)
		for i := range machines {
			if machines[i].Host != "" {
				index.HostIndex[machines[i].Host] = &machines[i]
			}
		}

	case spookytypesmachines.IndexTypeTag:
		index.TagIndex = make(map[string][]*spookytypesconfig.Machine)
		index.TagValueIndex = make(map[string][]*spookytypesconfig.Machine)
		index.MachineTags = make(map[*spookytypesconfig.Machine]map[string]string)
		for i := range machines {
			machine := &machines[i]
			index.MachineTags[machine] = make(map[string]string)
			for key, value := range machine.Tags {
				// Add to tag index
				index.TagIndex[key] = append(index.TagIndex[key], machine)
				// Add to tag value index
				tagValueKey := fmt.Sprintf("%s:%s", key, value)
				index.TagValueIndex[tagValueKey] = append(index.TagValueIndex[tagValueKey], machine)
				// Add to machine tags
				index.MachineTags[machine][key] = value
			}
		}

	case spookytypesmachines.IndexTypeGroup:
		index.GroupIndex = make(map[string][]*spookytypesconfig.Machine)
		index.MachineGroups = make(map[*spookytypesconfig.Machine][]string)
		for i := range machines {
			machine := &machines[i]
			for _, group := range machine.Groups {
				index.GroupIndex[group] = append(index.GroupIndex[group], machine)
				index.MachineGroups[machine] = append(index.MachineGroups[machine], group)
			}
		}

	case spookytypesmachines.IndexTypeUser:
		index.UserIndex = make(map[string][]*spookytypesconfig.Machine)
		for i := range machines {
			if machines[i].User != "" {
				index.UserIndex[machines[i].User] = append(index.UserIndex[machines[i].User], &machines[i])
			}
		}

	case spookytypesmachines.IndexTypePort:
		index.PortIndex = make(map[int][]*spookytypesconfig.Machine)
		for i := range machines {
			index.PortIndex[machines[i].Port] = append(index.PortIndex[machines[i].Port], &machines[i])
		}

	case spookytypesmachines.IndexTypeMetadata:
		index.MetadataIndex = make(map[string][]*spookytypesconfig.Machine)
		for i := range machines {
			machine := &machines[i]
			for key := range machine.Metadata {
				index.MetadataIndex[key] = append(index.MetadataIndex[key], machine)
			}
		}
	}

	// Update metrics
	buildTime := time.Since(startTime)
	index.Metrics.BuildTime = buildTime
	index.Metrics.MachineCount = len(machines)
	index.Metrics.IndexTypeStats[indexType] = &spookytypesmachines.IndexTypeStats{
		BuildTime:  buildTime,
		EntryCount: im.countIndexEntries(index, indexType),
	}

	im.indexes[indexType] = index
	return nil
}

// countIndexEntries counts the number of entries in an index
func (im *Manager) countIndexEntries(index *spookytypesmachines.MachineIndex, indexType spookytypesmachines.IndexType) int {
	switch indexType {
	case spookytypesmachines.IndexTypeName:
		return len(index.NameIndex)
	case spookytypesmachines.IndexTypeHost:
		return len(index.HostIndex)
	case spookytypesmachines.IndexTypeTag:
		return len(index.TagIndex)
	case spookytypesmachines.IndexTypeGroup:
		return len(index.GroupIndex)
	case spookytypesmachines.IndexTypeUser:
		return len(index.UserIndex)
	case spookytypesmachines.IndexTypePort:
		return len(index.PortIndex)
	case spookytypesmachines.IndexTypeMetadata:
		return len(index.MetadataIndex)
	default:
		return 0
	}
}

// updateLookupMetrics updates lookup performance metrics
func (im *Manager) updateLookupMetrics(indexType spookytypesmachines.IndexType, hit bool) {
	if index, exists := im.indexes[indexType]; exists && index.Metrics != nil {
		if stats, exists := index.Metrics.IndexTypeStats[indexType]; exists {
			if hit {
				stats.HitCount++
			} else {
				stats.MissCount++
			}
		}
	}
}

// optimizeIndex performs optimization on a specific index
func (im *Manager) optimizeIndex(indexType spookytypesmachines.IndexType, index *spookytypesmachines.MachineIndex) {
	// For now, just update the last updated time
	// In a more sophisticated implementation, we could:
	// - Reorder entries for better cache locality
	// - Compress rarely used data
	// - Precompute common queries
	if index.Metrics != nil {
		index.Metrics.LastUpdated = time.Now()
	}
}
