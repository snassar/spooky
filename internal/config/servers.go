package config

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// IndexMetrics tracks performance metrics for the index system
type IndexMetrics struct {
	BuildTime   time.Duration
	LookupTime  time.Duration
	MemoryUsage int64
	HitRate     float64
	ServerCount int
	TagCount    int
	LastUpdated time.Time
}

// TagIndex maps tag keys to servers for O(1) lookup
type TagIndex map[string][]*Server

// ServerTagIndex maps servers to their tags for reverse lookups
type ServerTagIndex map[*Server]map[string]string

// CompositeIndex provides multi-level indexing for enterprise-scale deployments
type CompositeIndex struct {
	TagIndex       TagIndex
	ServerTagIndex ServerTagIndex
	TagCount       map[string]int // Popularity tracking for optimization
	Metrics        *IndexMetrics  // Performance metrics
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
func buildEnterpriseIndex(servers []Server) *CompositeIndex {
	startTime := time.Now()

	tagIndex := make(TagIndex)
	serverTagIndex := make(ServerTagIndex)
	tagCount := make(map[string]int)
	uniqueTagNames := make(map[string]struct{})

	for i := range servers {
		server := &servers[i]
		serverTagIndex[server] = make(map[string]string)

		for tagName, tagValue := range server.Tags {
			if tagValue != "" {
				key := fmt.Sprintf("%s=%s", tagName, tagValue)
				tagIndex[key] = append(tagIndex[key], server)
				serverTagIndex[server][tagName] = tagValue
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
	for _, servers := range tagIndex {
		memoryUsage += int64(len(servers) * 8) // Pointer size
	}

	metrics := &IndexMetrics{
		BuildTime:   buildTime,
		ServerCount: len(servers),
		TagCount:    len(tagCount),
		MemoryUsage: memoryUsage,
		LastUpdated: time.Now(),
	}

	return &CompositeIndex{
		TagIndex:       tagIndex,
		ServerTagIndex: serverTagIndex,
		TagCount:       tagCount,
		Metrics:        metrics,
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

// GetServersForActionLarge provides optimized lookup for enterprise-scale deployments
func GetServersForActionLarge(config *Config, action *Action, index *CompositeIndex) ([]*Server, error) {
	startTime := time.Now()

	if len(action.Servers) > 0 {
		servers, err := findServersByName(config.Servers, action.Servers)
		updateLookupMetrics(index, time.Since(startTime))
		return servers, err
	}

	if len(action.Tags) > 0 {
		// Use map for O(1) deduplication with pre-allocated capacity
		serverMap := make(map[*Server]bool, len(action.Tags)*100)

		// Sort tags by popularity for better cache locality
		sortedTags := sortTagsByPopularity(action.Tags, index.TagCount)

		// For multiple tags, we need intersection (servers that match ALL tags)
		firstTag := true
		for _, tag := range sortedTags {
			if servers, exists := index.TagIndex[tag]; exists {
				if firstTag {
					// First tag: add all matching servers
					for _, server := range servers {
						serverMap[server] = true
					}
					firstTag = false
				} else {
					// Subsequent tags: keep only servers that match this tag too
					currentServers := make(map[*Server]bool)
					for _, server := range servers {
						if serverMap[server] {
							currentServers[server] = true
						}
					}
					serverMap = currentServers
				}
			} else {
				// Tag doesn't exist, no servers match
				serverMap = make(map[*Server]bool)
				break
			}
		}

		// Convert to slice with pre-allocated capacity
		targetServers := make([]*Server, 0, len(serverMap))
		for server := range serverMap {
			targetServers = append(targetServers, server)
		}

		updateLookupMetrics(index, time.Since(startTime))
		return targetServers, nil
	}

	servers := getAllServers(config.Servers)
	updateLookupMetrics(index, time.Since(startTime))
	return servers, nil
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
		ic.index = buildEnterpriseIndex(config.Servers)
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
	hash := fmt.Sprintf("%d-%d", len(config.Servers), len(config.Actions))
	for _, server := range config.Servers {
		hash += fmt.Sprintf("-%s-%s", server.Name, server.Host)
	}
	return hash
}

// findServersByName finds servers by their names (helper for optimized lookup)
func findServersByName(servers []Server, serverNames []string) ([]*Server, error) {
	serverMap := make(map[string]*Server)
	for i := range servers {
		serverMap[servers[i].Name] = &servers[i]
	}

	targetServers := make([]*Server, 0, len(serverNames))
	for _, serverName := range serverNames {
		if server, exists := serverMap[serverName]; exists {
			targetServers = append(targetServers, server)
		} else {
			return nil, fmt.Errorf("server '%s' not found", serverName)
		}
	}
	return targetServers, nil
}

// getAllServers returns all servers (helper for optimized lookup)
func getAllServers(servers []Server) []*Server {
	targetServers := make([]*Server, len(servers))
	for i := range servers {
		targetServers[i] = &servers[i]
	}
	return targetServers
}

// GetServersForAction returns the list of servers that should execute an action
func GetServersForAction(action *Action, config *Config) ([]*Server, error) {
	var targetServers []*Server

	// If specific servers are specified, use those
	if len(action.Servers) > 0 {
		serverMap := make(map[string]*Server)
		for i := range config.Servers {
			serverMap[config.Servers[i].Name] = &config.Servers[i]
		}

		for _, serverName := range action.Servers {
			if server, exists := serverMap[serverName]; exists {
				targetServers = append(targetServers, server)
			} else {
				return nil, fmt.Errorf("server '%s' not found for action '%s'", serverName, action.Name)
			}
		}
		return targetServers, nil
	}

	// If tags are specified, find servers with matching tags
	if len(action.Tags) > 0 {
		for i := range config.Servers {
			server := &config.Servers[i]
			for _, tag := range action.Tags {
				if value, exists := server.Tags[tag]; exists && value != "" {
					targetServers = append(targetServers, server)
					break
				}
			}
		}
		return targetServers, nil
	}

	// If no servers or tags specified, use all servers
	for i := range config.Servers {
		targetServers = append(targetServers, &config.Servers[i])
	}

	return targetServers, nil
}

// GetIndexMetrics returns current performance metrics
func (ic *IndexCache) GetIndexMetrics() *IndexMetrics {
	ic.mutex.RLock()
	defer ic.mutex.RUnlock()

	if ic.metrics == nil {
		return &IndexMetrics{}
	}

	return ic.metrics
}
