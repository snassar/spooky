package machines

import (
	"context"
	"fmt"
	"sync"
	"time"

	spookyconfigtypes "spooky/internal/types/config"
	spookylogging "spooky/internal/logging"
	spookymachinesconnectivity "spooky/internal/machines/connectivity"
	spookymachinesindexing "spooky/internal/machines/indexing"
	spookymachinestypes "spooky/internal/machines/types"
)

// Manager implements the MachineManager interface and coordinates all machine operations
type Manager struct {
	// Subpackage managers
	indexingManager     *spookymachinesindexing.Manager
	connectivityManager *spookymachinesconnectivity.Manager

	// Configuration
	config *spookymachinestypes.IndexManagerConfig
	logger spookylogging.Logger

	// State management
	state *spookymachinestypes.IndexManagerState
	mutex sync.RWMutex

	// Background workers
	optimizationTicker *time.Ticker
	stopChan           chan struct{}
	wg                 sync.WaitGroup
}

// NewManager creates a new machine manager with the given configuration
func NewManager(config *spookymachinestypes.IndexManagerConfig, logger spookylogging.Logger) *Manager {
	if config == nil {
		config = spookymachinestypes.DefaultIndexManagerConfig()
	}

	// Create subpackage managers
	indexingManager := spookymachinesindexing.NewManager()
	connectivityManager := spookymachinesconnectivity.NewManager(nil) // Use default options

	manager := &Manager{
		indexingManager:     indexingManager,
		connectivityManager: connectivityManager,
		config:              config,
		logger:              logger,
		state: &spookymachinestypes.IndexManagerState{
			LastBuilt:     time.Time{},
			LastOptimized: time.Time{},
		},
		stopChan: make(chan struct{}),
	}

	// Start background optimization if enabled
	if config.EnableOptimization {
		manager.startOptimizationWorker()
	}

	return manager
}

// BuildIndexes builds all indexes for the given machines
func (m *Manager) BuildIndexes(machines []spookyconfigtypes.Machine) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	startTime := time.Now()
	m.logger.Info("Building machine indexes", spookylogging.Int("machine_count", len(machines)))

	// Validate machine count
	if len(machines) > m.config.MaxIndexSize {
		return fmt.Errorf("machine count %d exceeds maximum index size %d", len(machines), m.config.MaxIndexSize)
	}

	// Build indexes using the indexing subpackage
	if err := m.indexingManager.BuildIndexes(machines); err != nil {
		m.state.LastError = err
		return fmt.Errorf("failed to build indexes: %w", err)
	}

	// Update state
	m.state.LastBuilt = time.Now()
	m.state.MachineCount = len(machines)
	m.state.IndexCount = 8 // Number of index types we support
	m.state.LastError = nil

	buildTime := time.Since(startTime)
	m.logger.Info("Machine indexes built successfully",
		spookylogging.Duration("build_time", buildTime.Milliseconds()),
		spookylogging.Int("machine_count", len(machines)))

	return nil
}

// UpdateIndexes updates existing indexes with new machine data
func (m *Manager) UpdateIndexes(machines []spookyconfigtypes.Machine) error {
	return m.indexingManager.UpdateIndexes(machines)
}

// GetState returns the current state of the index manager
func (m *Manager) GetState() *spookymachinestypes.IndexManagerState {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Create a copy to avoid race conditions
	state := *m.state
	return &state
}

// Stop stops the manager and cleans up resources
func (m *Manager) Stop() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Stop background workers
	if m.optimizationTicker != nil {
		m.optimizationTicker.Stop()
	}
	close(m.stopChan)

	// Wait for background workers to finish
	m.wg.Wait()

	m.logger.Info("Machine manager stopped")
	return nil
}

// LookupByName looks up a machine by name
func (m *Manager) LookupByName(name string) (*spookyconfigtypes.Machine, bool) {
	return m.indexingManager.LookupByName(name)
}

// LookupByHost looks up a machine by host
func (m *Manager) LookupByHost(host string) (*spookyconfigtypes.Machine, bool) {
	return m.indexingManager.LookupByHost(host)
}

// LookupByTag looks up machines by tag key
func (m *Manager) LookupByTag(tagKey string) ([]*spookyconfigtypes.Machine, bool) {
	return m.indexingManager.LookupByTag(tagKey)
}

// LookupByTagValue looks up machines by tag key and value
func (m *Manager) LookupByTagValue(tagKey, tagValue string) ([]*spookyconfigtypes.Machine, bool) {
	return m.indexingManager.LookupByTagValue(tagKey, tagValue)
}

// LookupByNetwork looks up machines by network type
func (m *Manager) LookupByNetwork(networkType string) ([]*spookyconfigtypes.Machine, bool) {
	// This functionality is not available in the current configtypes.Machine structure
	// Return empty result for now
	return nil, false
}

// LookupBySubnet looks up machines by subnet
func (m *Manager) LookupBySubnet(subnet string) ([]*spookyconfigtypes.Machine, bool) {
	// This functionality is not available in the current configtypes.Machine structure
	// Return empty result for now
	return nil, false
}

// FilterByTags filters machines by tag criteria
func (m *Manager) FilterByTags(criteria map[string]string) []*spookyconfigtypes.Machine {
	return m.indexingManager.FilterByTags(criteria)
}

// GetIndexMetrics returns performance metrics for all indexes
func (m *Manager) GetIndexMetrics() *spookymachinestypes.IndexMetrics {
	return m.indexingManager.GetMetrics()
}

// GetIndexPerformance returns detailed performance statistics
func (m *Manager) GetIndexPerformance() *spookymachinestypes.IndexPerformanceStats {
	metrics := m.indexingManager.GetMetrics()

	// Convert metrics to performance stats
	stats := &spookymachinestypes.IndexPerformanceStats{
		TotalBuildTime:   metrics.BuildTime,
		TotalLookupTime:  metrics.LookupTime,
		TotalMemoryUsage: metrics.MemoryUsage,
		AverageHitRate:   metrics.HitRate,
		IndexTypeStats:   make(map[spookymachinestypes.IndexType]*spookymachinestypes.IndexTypePerformanceStats),
	}

	// Convert index type stats
	for indexType, metricStats := range metrics.IndexTypeStats {
		stats.IndexTypeStats[indexType] = &spookymachinestypes.IndexTypePerformanceStats{
			BuildTime:   metricStats.BuildTime,
			LookupTime:  metricStats.LookupTime,
			MemoryUsage: metricStats.MemoryUsage,
			HitRate:     float64(metricStats.HitCount) / float64(metricStats.HitCount+metricStats.MissCount),
			EntryCount:  metricStats.EntryCount,
		}
	}

	return stats
}

// OptimizeIndexes performs optimization on all indexes
func (m *Manager) OptimizeIndexes() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.state.IsOptimizing {
		return fmt.Errorf("optimization already in progress")
	}

	m.state.IsOptimizing = true
	defer func() {
		m.state.IsOptimizing = false
		m.state.LastOptimized = time.Now()
	}()

	m.logger.Debug("Performing index optimization")
	m.indexingManager.OptimizeIndexes()

	return nil
}

// CleanupIndexes cleans up unused indexes and frees memory
func (m *Manager) CleanupIndexes() error {
	m.logger.Debug("Cleaning up indexes")
	// The indexing manager handles cleanup internally
	return nil
}

// ValidateIndexes validates the integrity of all indexes
func (m *Manager) ValidateIndexes() error {
	m.logger.Debug("Validating indexes")
	// For now, just check if indexes exist
	metrics := m.indexingManager.GetMetrics()
	if metrics.MachineCount == 0 {
		return fmt.Errorf("no machines indexed")
	}
	return nil
}

// startOptimizationWorker starts the background optimization worker
func (m *Manager) startOptimizationWorker() {
	m.optimizationTicker = time.NewTicker(m.config.OptimizationInterval)

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		for {
			select {
			case <-m.optimizationTicker.C:
				m.performBackgroundOptimization()
			case <-m.stopChan:
				return
			}
		}
	}()
}

// performBackgroundOptimization performs background index optimization
func (m *Manager) performBackgroundOptimization() {
	m.mutex.Lock()
	if m.state.IsOptimizing {
		m.mutex.Unlock()
		return
	}
	m.state.IsOptimizing = true
	m.mutex.Unlock()

	defer func() {
		m.mutex.Lock()
		m.state.IsOptimizing = false
		m.state.LastOptimized = time.Now()
		m.mutex.Unlock()
	}()

	m.logger.Debug("Performing background index optimization")
	m.indexingManager.OptimizeIndexes()
}

// TestMachineConnectivity tests connectivity to a specific machine
func (m *Manager) TestMachineConnectivity(ctx context.Context, machine *spookyconfigtypes.Machine) []spookymachinestypes.ConnectivityTestResult {
	return m.connectivityManager.TestMachineConnectivity(ctx, machine)
}

// TestMachinesConnectivity tests connectivity to multiple machines
func (m *Manager) TestMachinesConnectivity(ctx context.Context, machines []*spookyconfigtypes.Machine) map[string][]spookymachinestypes.ConnectivityTestResult {
	return m.connectivityManager.TestMachinesConnectivity(ctx, machines)
}

// GetConnectivityTestSummary provides a summary of connectivity test results
func (m *Manager) GetConnectivityTestSummary(results map[string][]spookymachinestypes.ConnectivityTestResult) *spookymachinestypes.ConnectivityTestSummary {
	return m.connectivityManager.GetTestSummary(results)
}
