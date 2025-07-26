package facts

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"encoding/json"
	"net/http" // Added for HTTP custom facts
	"os"
	"spooky/internal/ssh"
)

// Manager coordinates fact collection from multiple sources
type Manager struct {
	sshClient      *ssh.SSHClient
	sshCollector   *SSHCollector
	localCollector *LocalCollector
	hclCollector   *HCLCollector
	tofuCollector  *OpenTofuCollector

	// Custom collectors
	customCollectors map[string]FactCollector

	// Cache for collected facts
	cache      map[string]*FactCollection
	cacheMutex sync.RWMutex

	// Storage for persistent facts
	storage FactStorage

	// Configuration
	defaultTTL time.Duration
}

// NewManager creates a new fact collection manager
func NewManager(sshClient *ssh.SSHClient) *Manager {
	return &Manager{
		sshClient:        sshClient,
		sshCollector:     NewSSHCollector(sshClient),
		localCollector:   NewLocalCollector(),
		hclCollector:     nil, // Will be configured when file path is provided
		tofuCollector:    nil, // Will be configured when state path is provided
		customCollectors: make(map[string]FactCollector),
		cache:            make(map[string]*FactCollection),
		defaultTTL:       DefaultTTL,
	}
}

// NewManagerWithStorage creates a new fact collection manager with storage
func NewManagerWithStorage(sshClient *ssh.SSHClient, storage FactStorage) *Manager {
	return &Manager{
		sshClient:        sshClient,
		sshCollector:     NewSSHCollector(sshClient),
		localCollector:   NewLocalCollector(),
		hclCollector:     nil, // Will be configured when file path is provided
		tofuCollector:    nil, // Will be configured when state path is provided
		customCollectors: make(map[string]FactCollector),
		cache:            make(map[string]*FactCollection),
		storage:          storage,
		defaultTTL:       DefaultTTL,
	}
}

// ConfigureHCLCollector configures the HCL collector with a file path
func (m *Manager) ConfigureHCLCollector(filePath string) {
	m.hclCollector = NewHCLCollector(filePath, nil, MergePolicyReplace)
}

// ConfigureOpenTofuCollector configures the OpenTofu collector with a state path
func (m *Manager) ConfigureOpenTofuCollector(statePath string) {
	m.tofuCollector = NewOpenTofuCollector(statePath, nil, MergePolicyReplace)
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
	if m.hclCollector != nil {
		if collection, err := m.hclCollector.Collect(server); err == nil {
			collections = append(collections, collection)
		} else {
			errors = append(errors, fmt.Errorf("HCL collection failed: %w", err))
		}
	}

	// OpenTofu collection
	if m.tofuCollector != nil {
		if collection, err := m.tofuCollector.Collect(server); err == nil {
			collections = append(collections, collection)
		} else {
			errors = append(errors, fmt.Errorf("OpenTofu collection failed: %w", err))
		}
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
		if m.hclCollector != nil {
			return m.hclCollector.CollectSpecific(server, keys)
		}
		return nil, fmt.Errorf("HCL collector not configured")
	case SourceOpenTofu:
		if m.tofuCollector != nil {
			return m.tofuCollector.CollectSpecific(server, keys)
		}
		return nil, fmt.Errorf("OpenTofu collector not configured")
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
			if m.hclCollector != nil {
				fact, err = m.hclCollector.GetFact(server, key)
			} else {
				err = fmt.Errorf("HCL collector not configured")
			}
		case SourceOpenTofu:
			if m.tofuCollector != nil {
				fact, err = m.tofuCollector.GetFact(server, key)
			} else {
				err = fmt.Errorf("OpenTofu collector not configured")
			}
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
func (m *Manager) isHCLFact(key string) bool {
	// HCL facts have prefixes that match the patterns used in HCLCollector
	hclPrefixes := []string{
		"machine.", // machine.name, machine.host, machine.port, etc.
		"config.",  // config.machine_count, config.unique_tags, etc.
		"action.",  // action.name, action.description, etc.
		"hcl.",     // hcl.config, hcl.variables, etc.
	}

	for _, prefix := range hclPrefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}

// isOpenTofuFact checks if a fact key is an OpenTofu fact
func (m *Manager) isOpenTofuFact(key string) bool {
	// OpenTofu facts have prefixes that match the patterns used in OpenTofuCollector
	opentofuPrefixes := []string{
		"opentofu.", // opentofu.version, opentofu.terraform_version, etc.
	}

	for _, prefix := range opentofuPrefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
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

// Storage-related methods

// PersistFacts persists a fact collection to storage
func (m *Manager) PersistFacts(_ string, collection *FactCollection) error {
	if m.storage == nil {
		return nil // No storage configured
	}

	// Generate machine ID and convert to MachineFacts
	machineID := m.GenerateMachineID(collection)
	machineFacts := ConvertFactCollectionToMachineFacts(machineID, collection)

	return m.storage.SetMachineFacts(machineID, machineFacts)
}

// LoadPersistedFacts loads facts from storage for a server
func (m *Manager) LoadPersistedFacts(server string) (*FactCollection, error) {
	if m.storage == nil {
		return nil, fmt.Errorf("no storage configured")
	}

	// For now, we need to query by machine name since we don't have the machine ID
	query := &FactQuery{MachineName: server}
	machineFacts, err := m.storage.QueryFacts(query)
	if err != nil {
		return nil, err
	}

	if len(machineFacts) == 0 {
		return nil, fmt.Errorf("no facts found for server: %s", server)
	}

	// Use the first match
	return ConvertMachineFactsToFactCollection(machineFacts[0]), nil
}

// QueryPersistedFacts queries facts from storage
func (m *Manager) QueryPersistedFacts(query *FactQuery) ([]*FactCollection, error) {
	if m.storage == nil {
		return nil, fmt.Errorf("no storage configured")
	}

	machineFacts, err := m.storage.QueryFacts(query)
	if err != nil {
		return nil, err
	}

	var collections []*FactCollection
	for _, facts := range machineFacts {
		collection := ConvertMachineFactsToFactCollection(facts)
		collections = append(collections, collection)
	}

	return collections, nil
}

// ExportFacts exports all facts from storage to JSON
func (m *Manager) ExportFacts(w io.Writer) error {
	if m.storage == nil {
		return fmt.Errorf("no storage configured")
	}

	return m.storage.ExportToJSON(w)
}

// ImportFacts imports facts from JSON into storage
func (m *Manager) ImportFacts(r io.Reader) error {
	if m.storage == nil {
		return fmt.Errorf("no storage configured")
	}

	return m.storage.ImportFromJSON(r)
}

// CollectAndPersistFacts collects facts and persists them to storage
func (m *Manager) CollectAndPersistFacts(server string) (*FactCollection, error) {
	collection, err := m.CollectAllFacts(server)
	if err != nil {
		return nil, err
	}

	if m.storage != nil {
		// Generate machine ID from facts
		machineID := m.GenerateMachineID(collection)

		// Convert to MachineFacts and persist
		machineFacts := ConvertFactCollectionToMachineFacts(machineID, collection)
		if err := m.storage.SetMachineFacts(machineID, machineFacts); err != nil {
			return nil, fmt.Errorf("failed to persist facts: %w", err)
		}
	}

	return collection, nil
}

// GatherAndPersistFacts is an alias for CollectAndPersistFacts (new naming)
func (m *Manager) GatherAndPersistFacts(machine string) (*FactCollection, error) {
	return m.CollectAndPersistFacts(machine)
}

// GenerateMachineID generates a machine ID from fact collection
func (m *Manager) GenerateMachineID(facts *FactCollection) string {
	// Use machine_id fact if available
	if machineID, exists := facts.Facts["machine_id"]; exists {
		if id, ok := machineID.Value.(string); ok && id != "" {
			return id
		}
	}

	// Fallback: generate UUID from hostname + IP + action file
	return m.generateUUIDFromFacts(facts)
}

// generateUUIDFromFacts generates a UUID from fact data
func (m *Manager) generateUUIDFromFacts(facts *FactCollection) string {
	// Simple hash-based ID generation
	// In a real implementation, you'd use a proper UUID library
	data := facts.Server
	if hostname, exists := facts.Facts["hostname"]; exists {
		if str, ok := hostname.Value.(string); ok {
			data += str
		}
	}
	if ips, exists := facts.Facts["network.ips"]; exists {
		if ipList, ok := ips.Value.([]string); ok && len(ipList) > 0 {
			data += ipList[0]
		}
	}

	// Simple hash for now - in production use crypto/sha256
	hash := 0
	for _, char := range data {
		hash = ((hash << 5) - hash) + int(char)
		hash &= hash // Convert to 32-bit integer
	}

	return fmt.Sprintf("machine-%x", hash)
}

// GetMachineFacts retrieves machine facts from storage
func (m *Manager) GetMachineFacts(machineID string) (*MachineFacts, error) {
	if m.storage == nil {
		return nil, fmt.Errorf("no storage configured")
	}
	return m.storage.GetMachineFacts(machineID)
}

// SetMachineFacts stores machine facts to storage
func (m *Manager) SetMachineFacts(machineID string, facts *MachineFacts) error {
	if m.storage == nil {
		return fmt.Errorf("no storage configured")
	}
	return m.storage.SetMachineFacts(machineID, facts)
}

// QueryMachineFacts queries machine facts from storage
func (m *Manager) QueryMachineFacts(query *FactQuery) ([]*MachineFacts, error) {
	if m.storage == nil {
		return nil, fmt.Errorf("no storage configured")
	}
	return m.storage.QueryFacts(query)
}

// DeleteMachineFacts deletes machine facts from storage
func (m *Manager) DeleteMachineFacts(machineID string) error {
	if m.storage == nil {
		return fmt.Errorf("no storage configured")
	}
	return m.storage.DeleteMachineFacts(machineID)
}

// DeleteFacts deletes facts matching query criteria
func (m *Manager) DeleteFacts(query *FactQuery) (int, error) {
	if m.storage == nil {
		return 0, fmt.Errorf("no storage configured")
	}
	return m.storage.DeleteFacts(query)
}

// ExportFactsWithEncryption exports facts with encryption support
func (m *Manager) ExportFactsWithEncryption(w io.Writer, opts ExportOptions) error {
	if m.storage == nil {
		return fmt.Errorf("no storage configured")
	}
	return m.storage.ExportToJSONWithEncryption(w, opts)
}

// ImportFactsWithDecryption imports facts with decryption support
func (m *Manager) ImportFactsWithDecryption(r io.Reader, identityFile string) error {
	if m.storage == nil {
		return fmt.Errorf("no storage configured")
	}
	return m.storage.ImportFromJSONWithDecryption(r, identityFile)
}

// Close closes the storage connection
func (m *Manager) Close() error {
	if m.storage != nil {
		return m.storage.Close()
	}
	return nil
}

// RegisterCustomCollector registers a custom fact collector
func (m *Manager) RegisterCustomCollector(name string, collector FactCollector) {
	m.customCollectors[name] = collector
}

// ImportCustomFacts imports facts from a custom source (JSON file or HTTP endpoint)
func (m *Manager) ImportCustomFacts(source, server string, mergePolicy MergePolicy) (*FactCollection, error) {
	var collector FactCollector
	var err error

	// Determine source type and create appropriate collector
	if isHTTPURL(source) {
		collector = NewHTTPCollector(source, nil, 30*time.Second, mergePolicy)
	} else {
		// Assume local file
		collector = NewJSONCollector(source, mergePolicy)
	}

	// Collect facts from the custom source
	newCollection, err := collector.Collect(server)
	if err != nil {
		return nil, fmt.Errorf("failed to collect facts from %s: %w", source, err)
	}

	// Get existing facts if we have storage
	var existingCollection *FactCollection
	if m.storage != nil {
		existingCollection, _ = m.LoadPersistedFacts(server)
	}

	// Merge facts according to policy
	merger := NewFactMerger(mergePolicy)
	mergedCollection, err := merger.MergeCollections(existingCollection, newCollection)
	if err != nil {
		return nil, fmt.Errorf("failed to merge facts: %w", err)
	}

	// Persist merged facts if we have storage
	if m.storage != nil {
		if err := m.PersistFacts(server, mergedCollection); err != nil {
			return nil, fmt.Errorf("failed to persist merged facts: %w", err)
		}
	}

	// Update cache
	m.cacheFacts(server, mergedCollection)

	return mergedCollection, nil
}

// ImportCustomFactsWithOptions imports facts with enhanced options
func (m *Manager) ImportCustomFactsWithOptions(source string, options *ImportOptions) error {
	// Load custom facts from source
	var customFacts map[string]*CustomFacts
	var err error

	if isHTTPURL(source) {
		customFacts, err = m.loadHTTPCustomFacts(source)
	} else {
		customFacts, err = m.loadLocalCustomFacts(source)
	}

	if err != nil {
		return fmt.Errorf("failed to load custom facts: %w", err)
	}

	// Filter custom facts based on selectFacts list
	customFacts = m.filterCustomFacts(customFacts, options.SelectFacts)

	// Validate facts if requested
	if options.Validate {
		result := ValidateCustomFacts(customFacts)
		if !result.Valid {
			return fmt.Errorf("fact validation failed: %v", result.Errors)
		}
	}

	// Apply facts to storage
	for serverID, facts := range customFacts {
		// Skip if specific server is requested and this doesn't match
		if options.Server != "" && serverID != options.Server {
			continue
		}

		if options.DryRun {
			fmt.Printf("DRY RUN: Would import facts for %s\n", serverID)
			continue
		}

		existing, err := m.LoadPersistedFacts(serverID)
		if err != nil && !strings.Contains(err.Error(), "no facts found") {
			return fmt.Errorf("failed to get existing facts for %s: %w", serverID, err)
		}

		var merged *FactCollection
		if existing != nil {
			merged = m.mergeCustomFacts(existing, facts, options)
		} else {
			merged = m.convertCustomToFactCollection(facts, serverID)
		}

		if err := m.PersistFacts(serverID, merged); err != nil {
			return fmt.Errorf("failed to store merged facts for %s: %w", serverID, err)
		}

		fmt.Printf("Imported facts for %s\n", serverID)
	}

	return nil
}

// GetCustomFacts retrieves custom facts for template usage
func (m *Manager) GetCustomFacts(server string) (map[string]interface{}, error) {
	if m.storage == nil {
		return nil, fmt.Errorf("no storage configured")
	}

	// Load persisted facts
	collection, err := m.LoadPersistedFacts(server)
	if err != nil {
		return nil, err
	}

	// Extract custom facts
	customFacts := make(map[string]interface{})
	for key, fact := range collection.Facts {
		if strings.HasPrefix(key, "custom.") {
			parts := strings.Split(key, ".")
			if len(parts) >= 3 {
				category := parts[1]
				factKey := parts[2]

				if customFacts[category] == nil {
					customFacts[category] = make(map[string]interface{})
				}

				if categoryMap, ok := customFacts[category].(map[string]interface{}); ok {
					categoryMap[factKey] = fact.Value
				}
			}
		}
	}

	return customFacts, nil
}

// loadLocalCustomFacts loads custom facts from a local file
func (m *Manager) loadLocalCustomFacts(filePath string) (map[string]*CustomFacts, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var customFacts map[string]*CustomFacts
	if err := json.Unmarshal(data, &customFacts); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return customFacts, nil
}

// loadHTTPCustomFacts loads custom facts from an HTTP endpoint
func (m *Manager) loadHTTPCustomFacts(url string) (map[string]*CustomFacts, error) {
	// Enforce HTTPS-only connections for security
	if !strings.HasPrefix(url, "https://") {
		return nil, fmt.Errorf("HTTPS is required for custom facts import. HTTP URLs are not allowed for security reasons. Use https:// instead of http://")
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make HTTP request
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Add default headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "spooky-facts-collector/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status %d: %s", resp.StatusCode, resp.Status)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response as CustomFacts
	var customFacts map[string]*CustomFacts
	if err := json.Unmarshal(body, &customFacts); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Set source for all custom facts
	for _, facts := range customFacts {
		facts.Source = url
	}

	return customFacts, nil
}

// mergeCustomFacts merges custom facts with existing facts
func (m *Manager) mergeCustomFacts(existing *FactCollection, custom *CustomFacts, _ *ImportOptions) *FactCollection {
	merged := existing.Clone()

	// Merge custom facts
	if custom.Custom != nil {
		merged = m.mergeCustomFactSections(merged, custom.Custom, "custom")
	}

	// Apply overrides
	if custom.Overrides != nil {
		merged = ApplyOverrides(merged, custom.Overrides)
	}

	// Update metadata
	merged.Timestamp = time.Now()

	return merged
}

// mergeCustomFactSections merges custom fact sections
func (m *Manager) mergeCustomFactSections(collection *FactCollection, custom map[string]interface{}, prefix string) *FactCollection {
	for category, facts := range custom {
		if factsMap, ok := facts.(map[string]interface{}); ok {
			for key, value := range factsMap {
				factKey := fmt.Sprintf("%s.%s.%s", prefix, category, key)
				collection.Facts[factKey] = &Fact{
					Key:       factKey,
					Value:     value,
					Source:    string(SourceCustom),
					Server:    collection.Server,
					Timestamp: time.Now(),
					TTL:       DefaultTTL,
					Metadata:  map[string]interface{}{"category": category},
				}
			}
		}
	}

	return collection
}

// convertCustomToFactCollection converts CustomFacts to FactCollection
func (m *Manager) convertCustomToFactCollection(custom *CustomFacts, server string) *FactCollection {
	collection := &FactCollection{
		Server:    server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*Fact),
	}

	// Add custom facts
	if custom.Custom != nil {
		collection = m.mergeCustomFactSections(collection, custom.Custom, "custom")
	}

	// Add overrides
	if custom.Overrides != nil {
		collection = ApplyOverrides(collection, custom.Overrides)
	}

	return collection
}

// isHTTPURL checks if a string is an HTTPS URL (HTTP is not allowed for security)
func isHTTPURL(s string) bool {
	return len(s) > 8 && s[:8] == "https://"
}

// filterCustomFacts filters custom facts based on selectFacts list
func (m *Manager) filterCustomFacts(customFacts map[string]*CustomFacts, selectFacts []string) map[string]*CustomFacts {
	if len(selectFacts) == 0 {
		return customFacts // No filtering needed
	}

	filtered := make(map[string]*CustomFacts)

	for serverID, facts := range customFacts {
		filteredFacts := &CustomFacts{
			Custom:    make(map[string]interface{}),
			Overrides: make(map[string]interface{}),
			Source:    facts.Source,
		}

		// Filter custom facts
		if facts.Custom != nil {
			for category, categoryFacts := range facts.Custom {
				if categoryMap, ok := categoryFacts.(map[string]interface{}); ok {
					filteredCategory := make(map[string]interface{})

					for key, value := range categoryMap {
						// Check if this fact matches any of the select patterns
						factPath := fmt.Sprintf("%s.%s", category, key)
						if m.matchesSelectPattern(factPath, selectFacts) {
							filteredCategory[key] = value
						}
					}

					if len(filteredCategory) > 0 {
						filteredFacts.Custom[category] = filteredCategory
					}
				}
			}
		}

		// Filter overrides
		if facts.Overrides != nil {
			for category, categoryFacts := range facts.Overrides {
				if categoryMap, ok := categoryFacts.(map[string]interface{}); ok {
					filteredCategory := make(map[string]interface{})

					for key, value := range categoryMap {
						// Check if this fact matches any of the select patterns
						factPath := fmt.Sprintf("%s.%s", category, key)
						if m.matchesSelectPattern(factPath, selectFacts) {
							filteredCategory[key] = value
						}
					}

					if len(filteredCategory) > 0 {
						filteredFacts.Overrides[category] = filteredCategory
					}
				}
			}
		}

		// Only include server if it has any filtered facts
		if len(filteredFacts.Custom) > 0 || len(filteredFacts.Overrides) > 0 {
			filtered[serverID] = filteredFacts
		}
	}

	return filtered
}

// matchesSelectPattern checks if a fact path matches any of the select patterns
func (m *Manager) matchesSelectPattern(factPath string, selectFacts []string) bool {
	for _, pattern := range selectFacts {
		// Exact match
		if factPath == pattern {
			return true
		}

		// Category match (e.g., "application" matches "application.name", "application.version")
		if strings.HasPrefix(factPath, pattern+".") {
			return true
		}

		// Wildcard match (e.g., "*.name" matches "application.name", "environment.name")
		if strings.HasPrefix(pattern, "*.") {
			suffix := strings.TrimPrefix(pattern, "*.")
			if strings.HasSuffix(factPath, "."+suffix) {
				return true
			}
		}

		// Wildcard match (e.g., "*.port" matches "application.port", "monitoring.prometheus_port")
		if strings.HasSuffix(pattern, ".*") {
			category := strings.TrimSuffix(pattern, ".*")
			if strings.HasPrefix(factPath, category+".") {
				return true
			}
		}
	}

	return false
}
