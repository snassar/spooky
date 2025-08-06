package facts

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"os"
	"spooky/internal/facts/types"
	"spooky/internal/logging"
	"spooky/internal/ssh"
)

// Manager implements the FactManager interface
type Manager struct {
	// Core components
	collector types.FactCollector
	storage   FactStorage
	logger    logging.Logger

	// Configuration
	defaultTTL time.Duration
	cache      map[string]*types.FactCollection
	cacheMu    sync.RWMutex

	// Custom collectors
	customCollectors map[string]types.FactCollector
}

// NewManager creates a new fact collection manager
func NewManager(sshClient *ssh.SSHClient, logger logging.Logger) *Manager {
	return &Manager{
		collector:        NewSSHCollector(sshClient),
		storage:          nil, // Will be configured when storage is provided
		logger:           logger,
		customCollectors: make(map[string]types.FactCollector),
		cache:            make(map[string]*types.FactCollection),
		defaultTTL:       30 * time.Minute,
	}
}

// NewManagerWithStorage creates a new fact collection manager with storage
func NewManagerWithStorage(sshClient *ssh.SSHClient, storage FactStorage, logger logging.Logger) *Manager {
	manager := NewManager(sshClient, logger)
	manager.storage = storage
	return manager
}

// ConfigureHCLCollector configures the HCL collector with a file path
func (m *Manager) ConfigureHCLCollector(filePath string) {
	m.collector = NewHCLCollector(filePath, nil, types.MergePolicyReplace)
}

// ConfigureOpenTofuCollector configures the OpenTofu collector with a state path
func (m *Manager) ConfigureOpenTofuCollector(statePath string) {
	m.collector = NewOpenTofuCollector(statePath, nil, types.MergePolicyReplace)
}

// CollectAllFacts collects facts from all sources for a server
func (m *Manager) CollectAllFacts(server string) (*types.FactCollection, error) {
	// Check cache first
	if cached := m.getCachedFacts(server); cached != nil {
		return cached, nil
	}

	// Collect from all sources
	var collections []*types.FactCollection
	var errors []error

	// SSH collection (if server is remote)
	if server != "local" {
		if collection, err := m.collector.Collect(server); err == nil {
			collections = append(collections, collection)
		} else {
			errors = append(errors, fmt.Errorf("SSH collection failed: %w", err))
		}
	}

	// Local collection
	if collection, err := m.collector.Collect(server); err == nil {
		collections = append(collections, collection)
	} else {
		errors = append(errors, fmt.Errorf("local collection failed: %w", err))
	}

	// HCL collection
	if m.collector != nil {
		if collection, err := m.collector.Collect(server); err == nil {
			collections = append(collections, collection)
		} else {
			errors = append(errors, fmt.Errorf("HCL collection failed: %w", err))
		}
	}

	// OpenTofu collection
	if m.collector != nil {
		if collection, err := m.collector.Collect(server); err == nil {
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
func (m *Manager) CollectSpecificFacts(server string, keys []string) (*types.FactCollection, error) {
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

// GetFact retrieves a single fact
func (m *Manager) GetFact(server, key string) (*types.Fact, error) {
	// Check cache first
	if cached := m.getCachedFacts(server); cached != nil {
		if fact, exists := cached.Facts[key]; exists && !m.isExpired(fact) {
			return fact, nil
		}
	}

	// Collect from appropriate source
	sources := m.determineSources([]string{key})

	for _, source := range sources {
		var fact *types.Fact
		var err error

		switch source {
		case types.SourceSSH:
			if server != "local" {
				fact, err = m.collector.GetFact(server, key)
			}
		case types.SourceLocal:
			fact, err = m.collector.GetFact(server, key)
		case types.SourceHCL:
			if m.collector != nil {
				fact, err = m.collector.GetFact(server, key)
			} else {
				err = fmt.Errorf("HCL collector not configured")
			}
		case types.SourceOpenTofu:
			if m.collector != nil {
				fact, err = m.collector.GetFact(server, key)
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

// PersistFacts persists a fact collection to storage
func (m *Manager) PersistFacts(machineID string, collection *types.FactCollection) error {
	if m.storage == nil {
		return nil // No storage configured
	}

	// Generate machine ID if not provided
	if machineID == "" {
		machineID = m.GenerateMachineID(collection)
	}

	return m.storage.SetFactCollection(machineID, collection)
}

// LoadPersistedFacts loads facts from storage for a server
func (m *Manager) LoadPersistedFacts(server string) (*types.FactCollection, error) {
	if m.storage == nil {
		return nil, fmt.Errorf("no storage configured")
	}

	// Try to find facts by server name first
	query := &types.FactQuery{
		MachineName: server,
		Limit:       1,
	}

	factCollections, err := m.storage.QueryFactCollections(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query facts: %w", err)
	}

	if len(factCollections) == 0 {
		return nil, nil // No facts found
	}

	return factCollections[0], nil
}

// QueryPersistedFacts queries facts from storage
func (m *Manager) QueryPersistedFacts(query *types.FactQuery) ([]*types.FactCollection, error) {
	if m.storage == nil {
		return nil, fmt.Errorf("no storage configured")
	}

	return m.storage.QueryFactCollections(query)
}

// DeletePersistedFacts deletes facts from storage
func (m *Manager) DeletePersistedFacts(query *types.FactQuery) (int, error) {
	if m.storage == nil {
		return 0, fmt.Errorf("no storage configured")
	}

	return m.storage.DeleteFactCollections(query)
}

// ExportFacts exports combined facts to JSON
func (m *Manager) ExportFacts(w io.Writer) error {
	if m.storage == nil {
		return fmt.Errorf("no storage configured")
	}

	return m.storage.ExportToJSON(w)
}

// ImportFacts imports facts from JSON
func (m *Manager) ImportFacts(r io.Reader) error {
	if m.storage == nil {
		return fmt.Errorf("no storage configured")
	}

	return m.storage.ImportFromJSON(r)
}

// ExportFactsWithEncryption exports facts with encryption support
func (m *Manager) ExportFactsWithEncryption(w io.Writer, opts types.ExportOptions) error {
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

// ClearCache clears the fact cache
func (m *Manager) ClearCache() {
	m.cacheMu.Lock()
	defer m.cacheMu.Unlock()
	m.cache = make(map[string]*types.FactCollection)
}

// ClearExpiredCache removes expired facts from cache
func (m *Manager) ClearExpiredCache() {
	m.cacheMu.Lock()
	defer m.cacheMu.Unlock()

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

// GetAllFacts returns all cached facts from all servers
func (m *Manager) GetAllFacts() ([]*types.Fact, error) {
	m.cacheMu.RLock()
	defer m.cacheMu.RUnlock()

	var allFacts []*types.Fact
	for _, collection := range m.cache {
		for _, fact := range collection.Facts {
			if !m.isExpired(fact) {
				allFacts = append(allFacts, fact)
			}
		}
	}

	return allFacts, nil
}

// SetDefaultTTL sets the default TTL for facts
func (m *Manager) SetDefaultTTL(ttl time.Duration) {
	m.defaultTTL = ttl
}

// RegisterCustomCollector registers a custom fact collector
func (m *Manager) RegisterCustomCollector(name string, collector types.FactCollector) {
	m.customCollectors[name] = collector
}

// ImportCustomFacts imports facts from a custom source
func (m *Manager) ImportCustomFacts(source, server string, mergePolicy types.MergePolicy) (*types.FactCollection, error) {
	var collector types.FactCollector
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
	var existingCollection *types.FactCollection
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
func (m *Manager) ImportCustomFactsWithOptions(source string, options *types.ImportOptions) error {
	// For now, implement a simplified version
	// This would need to be expanded based on the options
	_, err := m.ImportCustomFacts(source, options.Server, types.MergePolicyReplace)
	return err
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
		if fact.Value != nil {
			customFacts[key] = fact.Value
		}
	}

	return customFacts, nil
}

// GenerateMachineID generates a machine ID from fact collection
func (m *Manager) GenerateMachineID(facts *types.FactCollection) string {
	// Use machine_id fact if available
	if machineID, exists := facts.Facts["machine_id"]; exists {
		if id, ok := machineID.Value.(string); ok && id != "" {
			return id
		}
	}

	// Fallback: generate UUID from hostname + IP
	return m.generateUUIDFromFacts(facts)
}

// Close closes the storage connection
func (m *Manager) Close() error {
	if m.storage != nil {
		return m.storage.Close()
	}
	return nil
}

// GetFactCollection retrieves a fact collection from storage
func (m *Manager) GetFactCollection(machineID string) (*types.FactCollection, error) {
	if m.storage == nil {
		return nil, fmt.Errorf("no storage configured")
	}
	return m.storage.GetFactCollection(machineID)
}

// SetFactCollection sets a fact collection for a machine
func (m *Manager) SetFactCollection(machineID string, collection *types.FactCollection) error {
	if m.storage == nil {
		return fmt.Errorf("storage not configured")
	}

	// Cache the collection
	m.cacheFacts(machineID, collection)

	// Persist to storage
	return m.storage.SetFactCollection(machineID, collection)
}

// GetStorage returns the underlying storage for coordinator integration
func (m *Manager) GetStorage() FactStorage {
	return m.storage
}

// Helper methods

// getFilteredCachedFacts returns filtered facts from cache if all are available and not expired
func (m *Manager) getFilteredCachedFacts(cached *types.FactCollection, keys []string) *types.FactCollection {
	// Check if all requested keys are in cache and not expired
	for _, key := range keys {
		if fact, exists := cached.Facts[key]; !exists || m.isExpired(fact) {
			return nil
		}
	}

	// Return only requested facts
	filtered := &types.FactCollection{
		Server:    cached.Server,
		Timestamp: cached.Timestamp,
		Facts:     make(map[string]*types.Fact),
	}
	for _, key := range keys {
		if fact, exists := cached.Facts[key]; exists {
			filtered.Facts[key] = fact
		}
	}
	return filtered
}

// collectFromSources collects facts from the appropriate sources
func (m *Manager) collectFromSources(server string, keys []string) ([]*types.FactCollection, []error) {
	var collections []*types.FactCollection
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
func (m *Manager) collectFromSource(source types.FactSource, server string, keys []string) (*types.FactCollection, error) {
	switch source {
	case types.SourceSSH:
		if server != "local" {
			return m.collector.CollectSpecific(server, keys)
		}
	case types.SourceLocal:
		return m.collector.CollectSpecific(server, keys)
	case types.SourceHCL:
		if m.collector != nil {
			return m.collector.CollectSpecific(server, keys)
		}
		return nil, fmt.Errorf("HCL collector not configured")
	case types.SourceOpenTofu:
		if m.collector != nil {
			return m.collector.CollectSpecific(server, keys)
		}
		return nil, fmt.Errorf("OpenTofu collector not configured")
	}
	return nil, nil
}

// determineSources determines which sources to use based on fact keys
func (m *Manager) determineSources(keys []string) []types.FactSource {
	sources := make(map[types.FactSource]bool)

	for _, key := range keys {
		switch {
		case m.isSystemFact(key):
			sources[types.SourceSSH] = true
			sources[types.SourceLocal] = true
		case m.isOSFact(key):
			sources[types.SourceSSH] = true
			sources[types.SourceLocal] = true
		case m.isHardwareFact(key):
			sources[types.SourceSSH] = true
			sources[types.SourceLocal] = true
		case m.isNetworkFact(key):
			sources[types.SourceSSH] = true
			sources[types.SourceLocal] = true
		case m.isEnvironmentFact(key):
			sources[types.SourceSSH] = true
			sources[types.SourceLocal] = true
		case m.isHCLFact(key):
			sources[types.SourceHCL] = true
		case m.isOpenTofuFact(key):
			sources[types.SourceOpenTofu] = true
		}
	}

	result := make([]types.FactSource, 0, len(sources))
	for source := range sources {
		result = append(result, source)
	}

	return result
}

// isSystemFact checks if a fact key is a system fact
func (m *Manager) isSystemFact(key string) bool {
	systemFacts := []string{"machine_id", "hostname", "fqdn"}
	for _, fact := range systemFacts {
		if key == fact {
			return true
		}
	}
	return false
}

// isOSFact checks if a fact key is an OS fact
func (m *Manager) isOSFact(key string) bool {
	osFacts := []string{"os.name", "os.version", "os.distribution", "os.architecture", "os.kernel"}
	for _, fact := range osFacts {
		if key == fact {
			return true
		}
	}
	return false
}

// isHardwareFact checks if a fact key is a hardware fact
func (m *Manager) isHardwareFact(key string) bool {
	hardwareFacts := []string{"cpu.cores", "cpu.model", "cpu.arch", "memory.total", "memory.used", "memory.available", "disk.total", "disk.used", "disk.available"}
	for _, fact := range hardwareFacts {
		if key == fact {
			return true
		}
	}
	return false
}

// isNetworkFact checks if a fact key is a network fact
func (m *Manager) isNetworkFact(key string) bool {
	networkFacts := []string{"network.ips", "network.macs", "dns"}
	for _, fact := range networkFacts {
		if key == fact {
			return true
		}
	}
	return false
}

// isEnvironmentFact checks if a fact key is an environment fact
func (m *Manager) isEnvironmentFact(key string) bool {
	return key == "environment"
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
func (m *Manager) mergeCollections(collections []*types.FactCollection) *types.FactCollection {
	if len(collections) == 0 {
		return nil
	}

	merged := &types.FactCollection{
		Server:    collections[0].Server,
		Timestamp: time.Now(),
		Facts:     make(map[string]*types.Fact),
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
func (m *Manager) cacheFacts(server string, collection *types.FactCollection) {
	m.cacheMu.Lock()
	defer m.cacheMu.Unlock()
	m.cache[server] = collection
}

// cacheFact caches a single fact
func (m *Manager) cacheFact(server string, fact *types.Fact) {
	m.cacheMu.Lock()
	defer m.cacheMu.Unlock()

	if m.cache[server] == nil {
		m.cache[server] = &types.FactCollection{
			Server:    server,
			Timestamp: time.Now(),
			Facts:     make(map[string]*types.Fact),
		}
	}

	m.cache[server].Facts[fact.Key] = fact
}

// getCachedFacts retrieves cached facts for a server
func (m *Manager) getCachedFacts(server string) *types.FactCollection {
	m.cacheMu.RLock()
	defer m.cacheMu.RUnlock()

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
func (m *Manager) isExpired(fact *types.Fact) bool {
	if fact.TTL == 0 {
		return false // No expiration
	}
	return time.Since(fact.Timestamp) > fact.TTL
}

// generateUUIDFromFacts generates a UUID from fact data
func (m *Manager) generateUUIDFromFacts(facts *types.FactCollection) string {
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

// isHTTPURL checks if a string is an HTTPS URL (HTTP is not allowed for security)
func isHTTPURL(s string) bool {
	return len(s) > 8 && s[:8] == "https://"
}

// ImportFactsToProject imports facts from file and saves to project structure
func (m *Manager) ImportFactsToProject(ctx context.Context, importPath string, projectPath string, format string) error {
	// 1. Open import file
	file, err := os.Open(importPath)
	if err != nil {
		return fmt.Errorf("failed to open import file: %w", err)
	}
	defer file.Close()

	// 2. Import facts based on format
	switch format {
	case "json":
		// Use existing JSON import
		if err := m.storage.ImportFromJSON(file); err != nil {
			return fmt.Errorf("failed to import facts from JSON: %w", err)
		}
		return nil
	case "hcl":
		// Use new HCL import
		if err := m.storage.ImportFromHCL(file); err != nil {
			return fmt.Errorf("failed to import facts from HCL: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unsupported import format: %s", format)
	}
}

// ExportFactsFromProject exports facts from project to file
func (m *Manager) ExportFactsFromProject(ctx context.Context, projectPath string, exportPath string, format string) error {
	// 1. Create output file
	file, err := os.Create(exportPath)
	if err != nil {
		return fmt.Errorf("failed to create export file: %w", err)
	}
	defer file.Close()

	// 2. Export based on format
	switch format {
	case "json":
		return m.ExportFacts(file)
	case "hcl":
		// Create exporter and export to HCL
		exporter := NewExporter(m.storage)
		query := &types.FactQuery{}
		return exporter.ExportToHCL(file, query)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}
