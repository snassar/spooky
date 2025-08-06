package performance

import (
	"fmt"
	"time"

	"spooky/internal/actions/types"
	"spooky/internal/logging"
)

// ActionOptimizerImpl implements the ActionOptimizer interface
type ActionOptimizerImpl struct {
	action             *types.Action
	level              types.OptimizationLevel
	resourceLimits     *types.ResourceLimits
	optimizationTarget types.OptimizationTarget
	logger             logging.Logger
}

// NewActionOptimizer creates a new ActionOptimizer
func NewActionOptimizer(action *types.Action, logger logging.Logger) ActionOptimizer {
	return &ActionOptimizerImpl{
		action:             action,
		level:              types.OptimizationLevelNone,
		resourceLimits:     &types.ResourceLimits{},
		optimizationTarget: types.OptimizationTargetSpeed,
		logger:             logger,
	}
}

// Optimize optimizes an action
func (o *ActionOptimizerImpl) Optimize(action *types.Action) error {
	if action == nil {
		return fmt.Errorf("action cannot be nil")
	}

	o.logger.Info("Optimizing action", logging.String("action", action.Name))

	// Apply optimizations based on level
	switch o.level {
	case types.OptimizationLevelNone:
		// No optimization
		o.logger.Debug("No optimization applied", logging.String("action", action.Name))

	case types.OptimizationLevelBasic:
		// Basic optimization
		if err := o.applyBasicOptimizations(action); err != nil {
			return fmt.Errorf("basic optimization failed: %w", err)
		}

	case types.OptimizationLevelAdvanced:
		// Advanced optimization
		if err := o.applyAdvancedOptimizations(action); err != nil {
			return fmt.Errorf("advanced optimization failed: %w", err)
		}

	case types.OptimizationLevelMaximum:
		// Maximum optimization
		if err := o.applyMaximumOptimizations(action); err != nil {
			return fmt.Errorf("maximum optimization failed: %w", err)
		}
	}

	o.logger.Info("Successfully optimized action", logging.String("action", action.Name))
	return nil
}

// GetMetrics gets performance metrics for an action
func (o *ActionOptimizerImpl) GetMetrics(action *types.Action) (*types.PerformanceMetrics, error) {
	if action == nil {
		return nil, fmt.Errorf("action cannot be nil")
	}

	o.logger.Info("Getting metrics for action", logging.String("action", action.Name))

	// Create basic metrics
	metrics := &types.PerformanceMetrics{
		ActionName:     action.Name,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ExecutionCount: 0,
		SuccessCount:   0,
		FailureCount:   0,
		SuccessRate:    0.0,
		Metadata:       make(map[string]interface{}),
	}

	o.logger.Info("Successfully retrieved metrics for action", logging.String("action", action.Name))
	return metrics, nil
}

// SetLevel sets the optimization level
func (o *ActionOptimizerImpl) SetLevel(level types.OptimizationLevel) {
	o.level = level
	o.logger.Debug("Set optimization level", logging.String("level", string(level)))
}

// GetLevel gets the current optimization level
func (o *ActionOptimizerImpl) GetLevel() types.OptimizationLevel {
	return o.level
}

// SetResourceLimits sets the resource limits
func (o *ActionOptimizerImpl) SetResourceLimits(limits *types.ResourceLimits) {
	o.resourceLimits = limits
	o.logger.Debug("Set resource limits")
}

// SetOptimizationTarget sets the optimization target
func (o *ActionOptimizerImpl) SetOptimizationTarget(target types.OptimizationTarget) {
	o.optimizationTarget = target
	o.logger.Debug("Set optimization target", logging.String("target", string(target)))
}

// applyBasicOptimizations applies basic optimizations to an action
func (o *ActionOptimizerImpl) applyBasicOptimizations(action *types.Action) error {
	// Set reasonable defaults if not specified
	if action.Timeout == 0 {
		action.Timeout = 300 // 5 minutes default
	}

	if action.MaxConcurrent == 0 {
		action.MaxConcurrent = 10 // Default concurrent limit
	}

	// Enable parallel execution for non-critical actions
	if !action.Critical && !action.Parallel {
		action.Parallel = true
	}

	return nil
}

// applyAdvancedOptimizations applies advanced optimizations to an action
func (o *ActionOptimizerImpl) applyAdvancedOptimizations(action *types.Action) error {
	// Apply basic optimizations first
	if err := o.applyBasicOptimizations(action); err != nil {
		return err
	}

	// Optimize based on target
	switch o.optimizationTarget {
	case types.OptimizationTargetSpeed:
		// Optimize for speed
		action.Parallel = true
		if action.MaxConcurrent < 20 {
			action.MaxConcurrent = 20
		}

	case types.OptimizationTargetMemory:
		// Optimize for memory usage
		if action.MaxConcurrent > 5 {
			action.MaxConcurrent = 5
		}
		action.Parallel = false

	case types.OptimizationTargetCPU:
		// Optimize for CPU usage
		if action.MaxConcurrent > 3 {
			action.MaxConcurrent = 3
		}
		action.Parallel = false

	case types.OptimizationTargetNetwork:
		// Optimize for network usage
		if action.MaxConcurrent > 2 {
			action.MaxConcurrent = 2
		}
		action.Parallel = false

	case types.OptimizationTargetBalanced:
		// Optimize for balanced usage
		if action.MaxConcurrent == 0 {
			action.MaxConcurrent = 10
		}
		action.Parallel = true
	}

	return nil
}

// applyMaximumOptimizations applies maximum optimizations to an action
func (o *ActionOptimizerImpl) applyMaximumOptimizations(action *types.Action) error {
	// Apply advanced optimizations first
	if err := o.applyAdvancedOptimizations(action); err != nil {
		return err
	}

	// Maximum optimizations
	action.Parallel = true
	action.MaxConcurrent = 50

	// Set resource limits if not specified
	if action.ResourceLimits == nil {
		action.ResourceLimits = &types.ActionResourceLimits{
			MemoryMB:   1024,
			CPUPercent: 50,
			DiskMB:     100,
		}
	}

	return nil
}
