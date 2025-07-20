package facts

import (
	"fmt"
	"sync"
	"time"

	"spooky/internal/ssh"
)

// Manager coordinates fact collection from multiple sources
type Manager struct {
	sshClient      *ssh.SSHClient
	sshCollector   *SSHCollector
	localCollector *LocalCollector
	hclCollector   *HCLCollector
	tofuCollector  *OpenTofuCollector

	// Cache for collected facts
	cache      map[string]*FactCollection
	cacheMutex sync.RWMutex

	// Configuration
	defaultTTL time.Duration
}

// NewManager creates a new fact collection manager
func NewManager(sshClient *ssh.SSHClient) *Manager {
	return &Manager{
		sshClient:      sshClient,
		sshCollector:   NewSSHCollector(sshClient),
		localCollector: NewLocalCollector(),
		hclCollector:   NewHCLCollector(),
		tofuCollector:  NewOpenTofuCollector(),
		cache:          make(map[string]*FactCollection),
		defaultTTL:     DefaultTTL,
	}
}

// CollectAllFacts collects facts from all sources for a server
func (m *Manager) CollectAllFacts(server string) (*FactCollection, error) {
	// Check cache first
	if cached := m.getCachedFacts(server); cached != nil {
		return cached, nil
	}

	// Collect from all sources
	var collections []*FactCollection
	var errors []error

	// SSH collection (if server is remote)
	if server != "local" {
		if collection, err := m.sshCollector.Collect(server); err == nil {
			collections = append(collections, collection)
		} else {
			errors = append(errors, fmt.Errorf("SSH collection failed: %w", err))
		}
	}

	// Local collection
	if collection, err := m.localCollector.Collect(server); err == nil {
		collections = append(collections, collection)
	} else {
		errors = append(errors, fmt.Errorf("local collection failed: %w", err))
	}

	// HCL collection
	if collection, err := m.hclCollector.Collect(server); err == nil {
		collections = append(collections, collection)
	} else {
		errors = append(errors, fmt.Errorf("HCL collection failed: %w", err))
	}

	// OpenTofu collection
	if collection, err := m.tofuCollector.Collect(server); err == nil {
		collections = append(collections, collection)
	} else {
		errors = append(errors, fmt.Errorf("OpenTofu collection failed: %w", err))
	}

	// Merge all collections
	if len(collections) == 0 {
		return nil, fmt.Errorf("no facts collected from any source: %v", errors)
	}

	merged := m.mergeCollections(collections)

	// Cache the result
	m.cacheFacts(server, merged)

	return merged, nil
}

// CollectSpecificFacts collects only the specified facts
func (m *Manager) CollectSpecificFacts(server string, keys []string) (*FactCollection, error) {
	// Check cache first for specific keys
	if cached := m.getCachedFacts(server); cached != nil {
		if filtered := m.getFilteredCachedFacts(cached, keys); filtered != nil {
			return filtered, nil
		}
	}

	// Collect from appropriate sources based on keys
	collections, errors := m.collectFromSources(server, keys)

	if len(collections) == 0 {
		return nil, fmt.Errorf("no facts collected from any source: %v", errors)
	}

	merged := m.mergeCollections(collections)

	// Cache the result
	m.cacheFacts(server, merged)

	return merged, nil
}

// getFilteredCachedFacts returns filtered facts from cache if all are available and not expired
func (m *Manager) getFilteredCachedFacts(cached *FactCollection, keys []string) *FactCollection {
	// Check if all requested keys are in cache and not expired
	for _, key := range keys {
		if fact, exists := cached.Facts[key]; !exists || m.isExpired(fact) {
			return nil
		}
	}

	// Return only requested facts
	filtered := &FactCollection{
		Server:    cached.Server,
		Timestamp: cached.Timestamp,
		Facts:     make(map[string]*Fact),
	}
	for _, key := range keys {
		if fact, exists := cached.Facts[key]; exists {
			filtered.Facts[key] = fact
		}
	}
	return filtered
}

// collectFromSources collects facts from the appropriate sources
func (m *Manager) collectFromSources(server string, keys []string) ([]*FactCollection, []error) {
	var collections []*FactCollection
	var errors []error

	// Determine which sources to use based on fact keys
	sources := m.determineSources(keys)

	for _, source := range sources {
		collection, err := m.collectFromSource(source, server, keys)
		if err == nil && collection != nil {
			collections = append(collections, collection)
		} else if err != nil {
			errors = append(errors, fmt.Errorf("%s collection failed: %w", source, err))
		}
	}

	return collections, errors
}

// collectFromSource collects facts from a specific source
func (m *Manager) collectFromSource(source FactSource, server string, keys []string) (*FactCollection, error) {
	switch source {
	case SourceSSH:
		if server != "local" {
			return m.sshCollector.CollectSpecific(server, keys)
		}
	case SourceLocal:
		return m.localCollector.CollectSpecific(server, keys)
	case SourceHCL:
		return m.hclCollector.CollectSpecific(server, keys)
	case SourceOpenTofu:
		return m.tofuCollector.CollectSpecific(server, keys)
	}
	return nil, nil
}

// GetFact retrieves a single fact
func (m *Manager) GetFact(server, key string) (*Fact, error) {
	// Check cache first
	if cached := m.getCachedFacts(server); cached != nil {
		if fact, exists := cached.Facts[key]; exists && !m.isExpired(fact) {
			return fact, nil
		}
	}

	// Collect from appropriate source
	sources := m.determineSources([]string{key})

	for _, source := range sources {
		var fact *Fact
		var err error

		switch source {
		case SourceSSH:
			if server != "local" {
				fact, err = m.sshCollector.GetFact(server, key)
			}
		case SourceLocal:
			fact, err = m.localCollector.GetFact(server, key)
		case SourceHCL:
			fact, err = m.hclCollector.GetFact(server, key)
		case SourceOpenTofu:
			fact, err = m.tofuCollector.GetFact(server, key)
		}

		if err == nil && fact != nil {
			// Cache the fact
			m.cacheFact(server, fact)
			return fact, nil
		}
	}

	return nil, fmt.Errorf("fact %s not found for server %s", key, server)
}

// ClearCache clears the fact cache
func (m *Manager) ClearCache() {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()
	m.cache = make(map[string]*FactCollection)
}

// ClearExpiredCache removes expired facts from cache
func (m *Manager) ClearExpiredCache() {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	for server, collection := range m.cache {
		expiredKeys := []string{}
		for key, fact := range collection.Facts {
			if m.isExpired(fact) {
				expiredKeys = append(expiredKeys, key)
			}
		}

		for _, key := range expiredKeys {
			delete(collection.Facts, key)
		}

		if len(collection.Facts) == 0 {
			delete(m.cache, server)
		}
	}
}

// SetDefaultTTL sets the default TTL for facts
func (m *Manager) SetDefaultTTL(ttl time.Duration) {
	m.defaultTTL = ttl
}

// Helper methods

// determineSources determines which sources to use based on fact keys
func (m *Manager) determineSources(keys []string) []FactSource {
	sources := make(map[FactSource]bool)

	for _, key := range keys {
		switch {
		case m.isSystemFact(key):
			sources[SourceSSH] = true
			sources[SourceLocal] = true
		case m.isOSFact(key):
			sources[SourceSSH] = true
			sources[SourceLocal] = true
		case m.isHardwareFact(key):
			sources[SourceSSH] = true
			sources[SourceLocal] = true
		case m.isNetworkFact(key):
			sources[SourceSSH] = true
			sources[SourceLocal] = true
		case m.isEnvironmentFact(key):
			sources[SourceSSH] = true
			sources[SourceLocal] = true
		case m.isHCLFact(key):
			sources[SourceHCL] = true
		case m.isOpenTofuFact(key):
			sources[SourceOpenTofu] = true
		}
	}

	result := make([]FactSource, 0, len(sources))
	for source := range sources {
		result = append(result, source)
	}

	return result
}

// isSystemFact checks if a fact key is a system fact
func (m *Manager) isSystemFact(key string) bool {
	systemFacts := []string{FactMachineID, FactHostname, FactFQDN}
	for _, fact := range systemFacts {
		if key == fact {
			return true
		}
	}
	return false
}

// isOSFact checks if a fact key is an OS fact
func (m *Manager) isOSFact(key string) bool {
	osFacts := []string{FactOSName, FactOSVersion, FactOSDistro, FactOSArch, FactOSKernel}
	for _, fact := range osFacts {
		if key == fact {
			return true
		}
	}
	return false
}

// isHardwareFact checks if a fact key is a hardware fact
func (m *Manager) isHardwareFact(key string) bool {
	hardwareFacts := []string{FactCPUCores, FactCPUModel, FactCPUArch, FactMemoryTotal, FactMemoryUsed, FactMemoryAvail, FactDiskTotal, FactDiskUsed, FactDiskAvail}
	for _, fact := range hardwareFacts {
		if key == fact {
			return true
		}
	}
	return false
}

// isNetworkFact checks if a fact key is a network fact
func (m *Manager) isNetworkFact(key string) bool {
	networkFacts := []string{FactNetworkIPs, FactNetworkMACs, FactDNS}
	for _, fact := range networkFacts {
		if key == fact {
			return true
		}
	}
	return false
}

// isEnvironmentFact checks if a fact key is an environment fact
func (m *Manager) isEnvironmentFact(key string) bool {
	return key == FactEnvironment
}

// isHCLFact checks if a fact key is an HCL fact
func (m *Manager) isHCLFact(_ string) bool {
	// TODO: Implement HCL fact detection
	return false
}

// isOpenTofuFact checks if a fact key is an OpenTofu fact
func (m *Manager) isOpenTofuFact(_ string) bool {
	// TODO: Implement OpenTofu fact detection
	return false
}

// mergeCollections merges multiple fact collections
func (m *Manager) mergeCollections(collections []*FactCollection) *FactCollection {
	if len(collections) == 0 {
		return nil
	}

	merged := &FactCollection{
		Server:    collections[0].Server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*Fact),
	}

	for _, collection := range collections {
		for key, fact := range collection.Facts {
			// Use the most recent fact if there are conflicts
			if existing, exists := merged.Facts[key]; !exists || fact.Timestamp.After(existing.Timestamp) {
				merged.Facts[key] = fact
			}
		}
	}

	return merged
}

// cacheFacts caches a fact collection
func (m *Manager) cacheFacts(server string, collection *FactCollection) {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()
	m.cache[server] = collection
}

// cacheFact caches a single fact
func (m *Manager) cacheFact(server string, fact *Fact) {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	if m.cache[server] == nil {
		m.cache[server] = &FactCollection{
			Server:    server,
			Timestamp: time.Now(),
			Facts:     make(map[string]*Fact),
		}
	}

	m.cache[server].Facts[fact.Key] = fact
}

// getCachedFacts retrieves cached facts for a server
func (m *Manager) getCachedFacts(server string) *FactCollection {
	m.cacheMutex.RLock()
	defer m.cacheMutex.RUnlock()

	if collection, exists := m.cache[server]; exists {
		// Check if any facts are expired
		hasValidFacts := false
		for _, fact := range collection.Facts {
			if !m.isExpired(fact) {
				hasValidFacts = true
				break
			}
		}

		if hasValidFacts {
			return collection
		}
	}

	return nil
}

// isExpired checks if a fact has expired
func (m *Manager) isExpired(fact *Fact) bool {
	if fact.TTL == 0 {
		return false // No expiration
	}
	return time.Since(fact.Timestamp) > fact.TTL
}

// GetAllFacts returns all cached facts from all servers
func (m *Manager) GetAllFacts() ([]*Fact, error) {
	m.cacheMutex.RLock()
	defer m.cacheMutex.RUnlock()

	var allFacts []*Fact
	for _, collection := range m.cache {
		for _, fact := range collection.Facts {
			if !m.isExpired(fact) {
				allFacts = append(allFacts, fact)
			}
		}
	}

	return allFacts, nil
}
