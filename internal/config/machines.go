package config

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// IndexMetrics tracks performance metrics for the index system
type IndexMetrics struct {
	BuildTime    time.Duration
	LookupTime   time.Duration
	MemoryUsage  int64
	HitRate      float64
	MachineCount int
	TagCount     int
	LastUpdated  time.Time
}

// TagIndex maps tag keys to machines for O(1) lookup
type TagIndex map[string][]*Machine

// MachineTagIndex maps machines to their tags for reverse lookups
type MachineTagIndex map[*Machine]map[string]string

// CompositeIndex provides multi-level indexing for enterprise-scale deployments
type CompositeIndex struct {
	TagIndex        TagIndex
	MachineTagIndex MachineTagIndex
	TagCount        map[string]int // Popularity tracking for optimization
	Metrics         *IndexMetrics  // Performance metrics
}

// IndexCache provides thread-safe caching of indexes
type IndexCache struct {
	index      *CompositeIndex
	lastBuilt  time.Time
	configHash string
	mutex      sync.RWMutex
	metrics    *IndexMetrics // Cache-level metrics
}

// buildEnterpriseIndex creates a composite index for optimal lookup performance
func buildEnterpriseIndex(machines []Machine) *CompositeIndex {
	startTime := time.Now()

	tagIndex := make(TagIndex)
	machineTagIndex := make(MachineTagIndex)
	tagCount := make(map[string]int)
	uniqueTagNames := make(map[string]struct{})

	for i := range machines {
		machine := &machines[i]
		machineTagIndex[machine] = make(map[string]string)

		for tagName, tagValue := range machine.Tags {
			if tagValue != "" {
				key := fmt.Sprintf("%s=%s", tagName, tagValue)
				tagIndex[key] = append(tagIndex[key], machine)
				machineTagIndex[machine][tagName] = tagValue
				uniqueTagNames[tagName] = struct{}{}
			}
		}
	}
	// Set tagCount to 1 for each unique tag name
	for tagName := range uniqueTagNames {
		tagCount[tagName] = 1
	}

	buildTime := time.Since(startTime)

	// Calculate memory usage (rough estimate)
	memoryUsage := int64(len(tagIndex) * 64) // Approximate size of map entries
	for _, machines := range tagIndex {
		memoryUsage += int64(len(machines) * 8) // Pointer size
	}

	metrics := &IndexMetrics{
		BuildTime:    buildTime,
		MachineCount: len(machines),
		TagCount:     len(tagCount),
		MemoryUsage:  memoryUsage,
		LastUpdated:  time.Now(),
	}

	return &CompositeIndex{
		TagIndex:        tagIndex,
		MachineTagIndex: machineTagIndex,
		TagCount:        tagCount,
		Metrics:         metrics,
	}
}

// sortTagsByPopularity sorts tags by usage frequency for better cache locality
func sortTagsByPopularity(tags []string, tagCount map[string]int) []string {
	sorted := make([]string, len(tags))
	copy(sorted, tags)

	sort.Slice(sorted, func(i, j int) bool {
		return tagCount[sorted[i]] > tagCount[sorted[j]]
	})

	return sorted
}

// GetMachinesForActionLarge provides optimized lookup for enterprise-scale deployments
func GetMachinesForActionLarge(config *Config, action *Action, index *CompositeIndex) ([]*Machine, error) {
	startTime := time.Now()

	if len(action.Machines) > 0 {
		machines, err := findMachinesByName(config.Machines, action.Machines)
		updateLookupMetrics(index, time.Since(startTime))
		return machines, err
	}

	if len(action.Tags) > 0 {
		// Use map for O(1) deduplication with pre-allocated capacity
		machineMap := make(map[*Machine]bool, len(action.Tags)*100)

		// Sort tags by popularity for better cache locality
		sortedTags := sortTagsByPopularity(action.Tags, index.TagCount)

		// For multiple tags, we need intersection (machines that match ALL tags)
		firstTag := true
		for _, tag := range sortedTags {
			if machines, exists := index.TagIndex[tag]; exists {
				if firstTag {
					// First tag: add all matching machines
					for _, machine := range machines {
						machineMap[machine] = true
					}
					firstTag = false
				} else {
					// Subsequent tags: keep only machines that match this tag too
					currentMachines := make(map[*Machine]bool)
					for _, machine := range machines {
						if machineMap[machine] {
							currentMachines[machine] = true
						}
					}
					machineMap = currentMachines
				}
			} else {
				// Tag doesn't exist, no machines match
				machineMap = make(map[*Machine]bool)
				break
			}
		}

		// Convert to slice with pre-allocated capacity
		targetMachines := make([]*Machine, 0, len(machineMap))
		for machine := range machineMap {
			targetMachines = append(targetMachines, machine)
		}

		updateLookupMetrics(index, time.Since(startTime))
		return targetMachines, nil
	}

	machines := getAllMachines(config.Machines)
	updateLookupMetrics(index, time.Since(startTime))
	return machines, nil
}

// updateLookupMetrics updates the lookup time metrics for the index
func updateLookupMetrics(index *CompositeIndex, lookupTime time.Duration) {
	if index.Metrics != nil {
		index.Metrics.LookupTime = lookupTime
		index.Metrics.LastUpdated = time.Now()
	}
}

// GetIndex provides thread-safe access to cached index
func (ic *IndexCache) GetIndex(config *Config) *CompositeIndex {
	ic.mutex.RLock()
	if ic.isValid(config) {
		defer ic.mutex.RUnlock()
		return ic.index
	}
	ic.mutex.RUnlock()

	ic.mutex.Lock()
	defer ic.mutex.Unlock()

	// Rebuild if needed
	if !ic.isValid(config) {
		ic.index = buildEnterpriseIndex(config.Machines)
		ic.lastBuilt = time.Now()
		ic.configHash = computeConfigHash(config)
		ic.metrics = ic.index.Metrics // Set cache metrics
	}
	return ic.index
}

// isValid checks if the cached index is still valid
func (ic *IndexCache) isValid(config *Config) bool {
	return ic.index != nil &&
		ic.configHash == computeConfigHash(config) &&
		time.Since(ic.lastBuilt) < 5*time.Minute // Cache for 5 minutes
}

// computeConfigHash creates a hash of the config for cache invalidation
func computeConfigHash(config *Config) string {
	// Simple hash implementation - could be improved with proper hashing
	hash := fmt.Sprintf("%d-%d", len(config.Machines), len(config.Actions))
	for _, machine := range config.Machines {
		hash += fmt.Sprintf("-%s-%s", machine.Name, machine.Host)
	}
	return hash
}

// findMachinesByName finds machines by their names (helper for optimized lookup)
func findMachinesByName(machines []Machine, machineNames []string) ([]*Machine, error) {
	machineMap := make(map[string]*Machine)
	for i := range machines {
		machineMap[machines[i].Name] = &machines[i]
	}

	targetMachines := make([]*Machine, 0, len(machineNames))
	for _, machineName := range machineNames {
		if machine, exists := machineMap[machineName]; exists {
			targetMachines = append(targetMachines, machine)
		} else {
			return nil, fmt.Errorf("machine '%s' not found", machineName)
		}
	}
	return targetMachines, nil
}

// getAllMachines returns all machines (helper for optimized lookup)
func getAllMachines(machines []Machine) []*Machine {
	targetMachines := make([]*Machine, len(machines))
	for i := range machines {
		targetMachines[i] = &machines[i]
	}
	return targetMachines
}

// GetMachinesForAction returns the list of machines that should execute an action
func GetMachinesForAction(action *Action, config *Config) ([]*Machine, error) {
	var targetMachines []*Machine

	// If specific machines are specified, use those
	if len(action.Machines) > 0 {
		machineMap := make(map[string]*Machine)
		for i := range config.Machines {
			machineMap[config.Machines[i].Name] = &config.Machines[i]
		}

		for _, machineName := range action.Machines {
			if machine, exists := machineMap[machineName]; exists {
				targetMachines = append(targetMachines, machine)
			} else {
				return nil, fmt.Errorf("machine '%s' not found", machineName)
			}
		}
		return targetMachines, nil
	}

	// If tags are specified, find machines matching those tags
	if len(action.Tags) > 0 {
		for i := range config.Machines {
			machine := &config.Machines[i]
			matchesAllTags := true

			for _, tag := range action.Tags {
				// Check if the machine has this tag
				found := false
				for machineTagKey, machineTagValue := range machine.Tags {
					if machineTagKey == tag || fmt.Sprintf("%s=%s", machineTagKey, machineTagValue) == tag {
						found = true
						break
					}
				}
				if !found {
					matchesAllTags = false
					break
				}
			}

			if matchesAllTags {
				targetMachines = append(targetMachines, machine)
			}
		}
		return targetMachines, nil
	}

	// If no machines or tags specified, return all machines
	for i := range config.Machines {
		targetMachines = append(targetMachines, &config.Machines[i])
	}

	return targetMachines, nil
}

// GetIndexMetrics returns the current metrics for the index cache
func (ic *IndexCache) GetIndexMetrics() *IndexMetrics {
	ic.mutex.RLock()
	defer ic.mutex.RUnlock()
	return ic.metrics
}
