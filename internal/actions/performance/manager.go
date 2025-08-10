package performance

import (
	"fmt"
	"sync"

	spookytypes "spooky/internal/types"
	spookylogging "spooky/internal/logging"
)

// Manager implements the PerformanceManager interface
type Manager struct {
	// Configuration
	optimizationLevel spookyactionstypes.OptimizationLevel
	resourceLimits    *spookyactionstypes.ResourceLimits

	// State
	optimizers map[string]ActionOptimizer
	metrics    map[string]*spookyactionstypes.PerformanceMetrics
	logger     spookylogging.Logger
	mu         sync.RWMutex
}

// NewManager creates a new PerformanceManager
func NewManager(logger spookylogging.Logger) *Manager {
	return &Manager{
		optimizationLevel: spookyactionstypes.OptimizationLevelNone,
		resourceLimits:    &spookyactionstypes.ResourceLimits{},
		optimizers:        make(map[string]ActionOptimizer),
		metrics:           make(map[string]*spookyactionstypes.PerformanceMetrics),
		logger:            logger,
	}
}

// OptimizeAction optimizes a single action
func (m *Manager) OptimizeAction(action *spookyactionstypes.Action) error {
	if action == nil {
		return fmt.Errorf("action cannot be nil")
	}

	m.logger.Info("Optimizing action", spookylogging.String("action", action.Name))

	// Create an optimizer for this action
	optimizer, err := m.CreateOptimizer(action)
	if err != nil {
		return fmt.Errorf("failed to create optimizer for action %s: %w", action.Name, err)
	}

	// Optimize the action
	if err := optimizer.Optimize(action); err != nil {
		return fmt.Errorf("failed to optimize action %s: %w", action.Name, err)
	}

	m.logger.Info("Successfully optimized action", spookylogging.String("action", action.Name))
	return nil
}

// OptimizeActionCollection optimizes a collection of actions
func (m *Manager) OptimizeActionCollection(collection *spookyactionstypes.ActionCollection) error {
	if collection == nil {
		return fmt.Errorf("action collection cannot be nil")
	}

	m.logger.Info("Optimizing action collection", spookylogging.Int("actions_count", len(collection.Actions)))

	// Optimize each action in the collection
	for _, action := range collection.Actions {
		if err := m.OptimizeAction(action); err != nil {
			return fmt.Errorf("failed to optimize action %s in collection: %w", action.Name, err)
		}
	}

	m.logger.Info("Successfully optimized action collection", spookylogging.Int("actions_count", len(collection.Actions)))
	return nil
}

// GetPerformanceMetrics gets performance metrics for an action
func (m *Manager) GetPerformanceMetrics(action *spookyactionstypes.Action) (*spookyactionstypes.PerformanceMetrics, error) {
	if action == nil {
		return nil, fmt.Errorf("action cannot be nil")
	}

	m.logger.Info("Getting performance metrics", spookylogging.String("action", action.Name))

	// Create an optimizer for this action
	optimizer, err := m.CreateOptimizer(action)
	if err != nil {
		return nil, fmt.Errorf("failed to create optimizer for action %s: %w", action.Name, err)
	}

	// Get metrics from the optimizer
	metrics, err := optimizer.GetMetrics(action)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics for action %s: %w", action.Name, err)
	}

	m.logger.Info("Successfully retrieved performance metrics", spookylogging.String("action", action.Name))
	return metrics, nil
}

// CreateOptimizer creates a new optimizer for an action
func (m *Manager) CreateOptimizer(action *spookyactionstypes.Action) (ActionOptimizer, error) {
	if action == nil {
		return nil, fmt.Errorf("action cannot be nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if optimizer already exists
	if optimizer, exists := m.optimizers[action.Name]; exists {
		return optimizer, nil
	}

	// Create new optimizer
	optimizer := NewActionOptimizer(action, m.logger)
	m.optimizers[action.Name] = optimizer

	m.logger.Debug("Created optimizer for action", spookylogging.String("action", action.Name))
	return optimizer, nil
}

// GetOptimizer gets an existing optimizer for an action
func (m *Manager) GetOptimizer(action *spookyactionstypes.Action) (ActionOptimizer, error) {
	if action == nil {
		return nil, fmt.Errorf("action cannot be nil")
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	optimizer, exists := m.optimizers[action.Name]
	if !exists {
		return nil, fmt.Errorf("optimizer not found for action %s", action.Name)
	}

	return optimizer, nil
}

// GetMetrics gets performance metrics for an action
func (m *Manager) GetMetrics(action *spookyactionstypes.Action) (*spookyactionstypes.PerformanceMetrics, error) {
	if action == nil {
		return nil, fmt.Errorf("action cannot be nil")
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics, exists := m.metrics[action.Name]
	if !exists {
		return nil, fmt.Errorf("metrics not found for action %s", action.Name)
	}

	return metrics, nil
}

// ListMetrics lists all performance metrics
func (m *Manager) ListMetrics() ([]*spookyactionstypes.PerformanceMetrics, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics := make([]*spookyactionstypes.PerformanceMetrics, 0, len(m.metrics))
	for _, metric := range m.metrics {
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// ClearMetrics clears performance metrics for an action
func (m *Manager) ClearMetrics(action *spookyactionstypes.Action) error {
	if action == nil {
		return fmt.Errorf("action cannot be nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.metrics, action.Name)
	m.logger.Info("Cleared metrics for action", spookylogging.String("action", action.Name))
	return nil
}

// SetOptimizationLevel sets the optimization level
func (m *Manager) SetOptimizationLevel(level spookyactionstypes.OptimizationLevel) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.optimizationLevel = level
	m.logger.Info("Set optimization level", spookylogging.String("level", string(level)))
}

// SetResourceLimits sets the resource limits
func (m *Manager) SetResourceLimits(limits *spookyactionstypes.ResourceLimits) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.resourceLimits = limits
	m.logger.Info("Set resource limits")
}
