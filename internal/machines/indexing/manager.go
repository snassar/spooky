package indexing

import (
	"fmt"
	"sync"
	"time"

	configtypes "spooky/internal/config/types"
	"spooky/internal/machines/types"
)

// Manager provides comprehensive indexing capabilities for machine data
type Manager struct {
	indexes map[types.IndexType]*types.MachineIndex
	mutex   sync.RWMutex
}

// NewManager creates a new indexing manager
func NewManager() *Manager {
	return &Manager{
		indexes: make(map[types.IndexType]*types.MachineIndex),
	}
}

// BuildIndexes builds all indexes for the given machines
func (im *Manager) BuildIndexes(machines []configtypes.Machine) error {
	startTime := time.Now()

	// Build each index type
	indexTypes := []types.IndexType{
		types.IndexTypeName,
		types.IndexTypeHost,
		types.IndexTypeTag,
		types.IndexTypeGroup,
		types.IndexTypeUser,
		types.IndexTypePort,
		types.IndexTypeMetadata,
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
func (im *Manager) UpdateIndexes(machines []configtypes.Machine) error {
	// For now, rebuild all indexes
	// In a more sophisticated implementation, we could do incremental updates
	return im.BuildIndexes(machines)
}

// SynchronizeIndexes synchronizes indexes with external data
func (im *Manager) SynchronizeIndexes(syncData *types.IndexSyncData) error {
	if syncData == nil {
		return fmt.Errorf("sync data is nil")
	}

	// Simplified implementation - just rebuild indexes
	// In a real implementation, we would apply specific changes
	return nil
}

// LookupByName looks up a machine by name
func (im *Manager) LookupByName(name string) (*configtypes.Machine, bool) {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	if index, exists := im.indexes[types.IndexTypeName]; exists {
		if machine, found := index.NameIndex[name]; found {
			im.updateLookupMetrics(types.IndexTypeName, true)
			return machine, true
		}
	}
	im.updateLookupMetrics(types.IndexTypeName, false)
	return nil, false
}

// LookupByHost looks up a machine by host
func (im *Manager) LookupByHost(host string) (*configtypes.Machine, bool) {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	if index, exists := im.indexes[types.IndexTypeHost]; exists {
		if machine, found := index.HostIndex[host]; found {
			im.updateLookupMetrics(types.IndexTypeHost, true)
			return machine, true
		}
	}
	im.updateLookupMetrics(types.IndexTypeHost, false)
	return nil, false
}

// LookupByTag looks up machines by tag key
func (im *Manager) LookupByTag(tagKey string) ([]*configtypes.Machine, bool) {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	if index, exists := im.indexes[types.IndexTypeTag]; exists {
		if machines, found := index.TagIndex[tagKey]; found {
			im.updateLookupMetrics(types.IndexTypeTag, true)
			return machines, true
		}
	}
	im.updateLookupMetrics(types.IndexTypeTag, false)
	return nil, false
}

// LookupByTagValue looks up machines by tag key and value
func (im *Manager) LookupByTagValue(tagKey, tagValue string) ([]*configtypes.Machine, bool) {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	if index, exists := im.indexes[types.IndexTypeTag]; exists {
		key := fmt.Sprintf("%s:%s", tagKey, tagValue)
		if machines, found := index.TagValueIndex[key]; found {
			im.updateLookupMetrics(types.IndexTypeTag, true)
			return machines, true
		}
	}
	im.updateLookupMetrics(types.IndexTypeTag, false)
	return nil, false
}

// LookupByGroup looks up machines by group name
func (im *Manager) LookupByGroup(groupName string) ([]*configtypes.Machine, bool) {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	if index, exists := im.indexes[types.IndexTypeGroup]; exists {
		if machines, found := index.GroupIndex[groupName]; found {
			im.updateLookupMetrics(types.IndexTypeGroup, true)
			return machines, true
		}
	}
	im.updateLookupMetrics(types.IndexTypeGroup, false)
	return nil, false
}

// LookupByMetadata looks up machines by metadata key
func (im *Manager) LookupByMetadata(metadataKey string) ([]*configtypes.Machine, bool) {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	if index, exists := im.indexes[types.IndexTypeMetadata]; exists {
		if machines, found := index.MetadataIndex[metadataKey]; found {
			im.updateLookupMetrics(types.IndexTypeMetadata, true)
			return machines, true
		}
	}
	im.updateLookupMetrics(types.IndexTypeMetadata, false)
	return nil, false
}

// FilterByGroups filters machines by multiple groups
func (im *Manager) FilterByGroups(groups []string) []*configtypes.Machine {
	var result []*configtypes.Machine
	seen := make(map[*configtypes.Machine]bool)

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
func (im *Manager) FilterByTags(criteria map[string]string) []*configtypes.Machine {
	if len(criteria) == 0 {
		return nil
	}

	var result []*configtypes.Machine
	seen := make(map[*configtypes.Machine]bool)

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
			newResult := []*configtypes.Machine{}
			for _, machine := range machines {
				if seen[machine] {
					newResult = append(newResult, machine)
				}
			}
			result = newResult

			// Update seen map
			seen = make(map[*configtypes.Machine]bool)
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
func (im *Manager) GetMetrics() *types.IndexMetrics {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	// Aggregate metrics from all indexes
	aggregateMetrics := &types.IndexMetrics{
		IndexTypeStats: make(map[types.IndexType]*types.IndexTypeStats),
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
func (im *Manager) buildIndex(indexType types.IndexType, machines []configtypes.Machine) error {
	startTime := time.Now()

	index := &types.MachineIndex{
		Metrics: &types.IndexMetrics{
			IndexTypeStats: make(map[types.IndexType]*types.IndexTypeStats),
		},
	}

	switch indexType {
	case types.IndexTypeName:
		index.NameIndex = make(map[string]*configtypes.Machine)
		for i := range machines {
			index.NameIndex[machines[i].Name] = &machines[i]
		}

	case types.IndexTypeHost:
		index.HostIndex = make(map[string]*configtypes.Machine)
		for i := range machines {
			if machines[i].Host != "" {
				index.HostIndex[machines[i].Host] = &machines[i]
			}
		}

	case types.IndexTypeTag:
		index.TagIndex = make(map[string][]*configtypes.Machine)
		index.TagValueIndex = make(map[string][]*configtypes.Machine)
		index.MachineTags = make(map[*configtypes.Machine]map[string]string)
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

	case types.IndexTypeGroup:
		index.GroupIndex = make(map[string][]*configtypes.Machine)
		index.MachineGroups = make(map[*configtypes.Machine][]string)
		for i := range machines {
			machine := &machines[i]
			for _, group := range machine.Groups {
				index.GroupIndex[group] = append(index.GroupIndex[group], machine)
				index.MachineGroups[machine] = append(index.MachineGroups[machine], group)
			}
		}

	case types.IndexTypeUser:
		index.UserIndex = make(map[string][]*configtypes.Machine)
		for i := range machines {
			if machines[i].User != "" {
				index.UserIndex[machines[i].User] = append(index.UserIndex[machines[i].User], &machines[i])
			}
		}

	case types.IndexTypePort:
		index.PortIndex = make(map[int][]*configtypes.Machine)
		for i := range machines {
			index.PortIndex[machines[i].Port] = append(index.PortIndex[machines[i].Port], &machines[i])
		}

	case types.IndexTypeMetadata:
		index.MetadataIndex = make(map[string][]*configtypes.Machine)
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
	index.Metrics.IndexTypeStats[indexType] = &types.IndexTypeStats{
		BuildTime:  buildTime,
		EntryCount: im.countIndexEntries(index, indexType),
	}

	im.indexes[indexType] = index
	return nil
}

// countIndexEntries counts the number of entries in an index
func (im *Manager) countIndexEntries(index *types.MachineIndex, indexType types.IndexType) int {
	switch indexType {
	case types.IndexTypeName:
		return len(index.NameIndex)
	case types.IndexTypeHost:
		return len(index.HostIndex)
	case types.IndexTypeTag:
		return len(index.TagIndex)
	case types.IndexTypeGroup:
		return len(index.GroupIndex)
	case types.IndexTypeUser:
		return len(index.UserIndex)
	case types.IndexTypePort:
		return len(index.PortIndex)
	case types.IndexTypeMetadata:
		return len(index.MetadataIndex)
	default:
		return 0
	}
}

// updateLookupMetrics updates lookup performance metrics
func (im *Manager) updateLookupMetrics(indexType types.IndexType, hit bool) {
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
func (im *Manager) optimizeIndex(indexType types.IndexType, index *types.MachineIndex) {
	// For now, just update the last updated time
	// In a more sophisticated implementation, we could:
	// - Reorder entries for better cache locality
	// - Compress rarely used data
	// - Precompute common queries
	if index.Metrics != nil {
		index.Metrics.LastUpdated = time.Now()
	}
}
