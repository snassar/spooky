package performance

import (
	"fmt"
	"sync"

	"spooky/internal/actions/types"
	"spooky/internal/logging"
)

// Manager implements the PerformanceManager interface
type Manager struct {
	// Configuration
	optimizationLevel types.OptimizationLevel
	resourceLimits    *types.ResourceLimits

	// State
	optimizers map[string]ActionOptimizer
	metrics    map[string]*types.PerformanceMetrics
	logger     logging.Logger
	mu         sync.RWMutex
}

// NewManager creates a new PerformanceManager
func NewManager(logger logging.Logger) *Manager {
	return &Manager{
		optimizationLevel: types.OptimizationLevelNone,
		resourceLimits:    &types.ResourceLimits{},
		optimizers:        make(map[string]ActionOptimizer),
		metrics:           make(map[string]*types.PerformanceMetrics),
		logger:            logger,
	}
}

// OptimizeAction optimizes a single action
func (m *Manager) OptimizeAction(action *types.Action) error {
	if action == nil {
		return fmt.Errorf("action cannot be nil")
	}

	m.logger.Info("Optimizing action", logging.String("action", action.Name))

	// Create an optimizer for this action
	optimizer, err := m.CreateOptimizer(action)
	if err != nil {
		return fmt.Errorf("failed to create optimizer for action %s: %w", action.Name, err)
	}

	// Optimize the action
	if err := optimizer.Optimize(action); err != nil {
		return fmt.Errorf("failed to optimize action %s: %w", action.Name, err)
	}

	m.logger.Info("Successfully optimized action", logging.String("action", action.Name))
	return nil
}

// OptimizeActionCollection optimizes a collection of actions
func (m *Manager) OptimizeActionCollection(collection *types.ActionCollection) error {
	if collection == nil {
		return fmt.Errorf("action collection cannot be nil")
	}

	m.logger.Info("Optimizing action collection", logging.Int("actions_count", len(collection.Actions)))

	// Optimize each action in the collection
	for _, action := range collection.Actions {
		if err := m.OptimizeAction(action); err != nil {
			return fmt.Errorf("failed to optimize action %s in collection: %w", action.Name, err)
		}
	}

	m.logger.Info("Successfully optimized action collection", logging.Int("actions_count", len(collection.Actions)))
	return nil
}

// GetPerformanceMetrics gets performance metrics for an action
func (m *Manager) GetPerformanceMetrics(action *types.Action) (*types.PerformanceMetrics, error) {
	if action == nil {
		return nil, fmt.Errorf("action cannot be nil")
	}

	m.logger.Info("Getting performance metrics", logging.String("action", action.Name))

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

	m.logger.Info("Successfully retrieved performance metrics", logging.String("action", action.Name))
	return metrics, nil
}

// CreateOptimizer creates a new optimizer for an action
func (m *Manager) CreateOptimizer(action *types.Action) (ActionOptimizer, error) {
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

	m.logger.Debug("Created optimizer for action", logging.String("action", action.Name))
	return optimizer, nil
}

// GetOptimizer gets an existing optimizer for an action
func (m *Manager) GetOptimizer(action *types.Action) (ActionOptimizer, error) {
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
func (m *Manager) GetMetrics(action *types.Action) (*types.PerformanceMetrics, error) {
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
func (m *Manager) ListMetrics() ([]*types.PerformanceMetrics, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics := make([]*types.PerformanceMetrics, 0, len(m.metrics))
	for _, metric := range m.metrics {
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// ClearMetrics clears performance metrics for an action
func (m *Manager) ClearMetrics(action *types.Action) error {
	if action == nil {
		return fmt.Errorf("action cannot be nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.metrics, action.Name)
	m.logger.Info("Cleared metrics for action", logging.String("action", action.Name))
	return nil
}

// SetOptimizationLevel sets the optimization level
func (m *Manager) SetOptimizationLevel(level types.OptimizationLevel) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.optimizationLevel = level
	m.logger.Info("Set optimization level", logging.String("level", string(level)))
}

// SetResourceLimits sets the resource limits
func (m *Manager) SetResourceLimits(limits *types.ResourceLimits) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.resourceLimits = limits
	m.logger.Info("Set resource limits")
}
