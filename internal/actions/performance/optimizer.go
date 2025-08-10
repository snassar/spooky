package performance

import (
	"fmt"
	"time"

	spookytypes "spooky/internal/types"
	spookylogging "spooky/internal/logging"
)

// ActionOptimizerImpl implements the ActionOptimizer interface
type ActionOptimizerImpl struct {
	action             *spookyactionstypes.Action
	level              spookyactionstypes.OptimizationLevel
	resourceLimits     *spookyactionstypes.ResourceLimits
	optimizationTarget spookyactionstypes.OptimizationTarget
	logger             spookylogging.Logger
}

// NewActionOptimizer creates a new ActionOptimizer
func NewActionOptimizer(action *spookyactionstypes.Action, logger spookylogging.Logger) ActionOptimizer {
	return &ActionOptimizerImpl{
		action:             action,
		level:              spookyactionstypes.OptimizationLevelNone,
		resourceLimits:     &spookyactionstypes.ResourceLimits{},
		optimizationTarget: spookyactionstypes.OptimizationTargetSpeed,
		logger:             logger,
	}
}

// Optimize optimizes an action
func (o *ActionOptimizerImpl) Optimize(action *spookyactionstypes.Action) error {
	if action == nil {
		return fmt.Errorf("action cannot be nil")
	}

	o.logger.Info("Optimizing action", spookylogging.String("action", action.Name))

	// Apply optimizations based on level
	switch o.level {
	case spookyactionstypes.OptimizationLevelNone:
		// No optimization
		o.logger.Debug("No optimization applied", spookylogging.String("action", action.Name))

	case spookyactionstypes.OptimizationLevelBasic:
		// Basic optimization
		if err := o.applyBasicOptimizations(action); err != nil {
			return fmt.Errorf("basic optimization failed: %w", err)
		}

	case spookyactionstypes.OptimizationLevelAdvanced:
		// Advanced optimization
		if err := o.applyAdvancedOptimizations(action); err != nil {
			return fmt.Errorf("advanced optimization failed: %w", err)
		}

	case spookyactionstypes.OptimizationLevelMaximum:
		// Maximum optimization
		if err := o.applyMaximumOptimizations(action); err != nil {
			return fmt.Errorf("maximum optimization failed: %w", err)
		}
	}

	o.logger.Info("Successfully optimized action", spookylogging.String("action", action.Name))
	return nil
}

// GetMetrics gets performance metrics for an action
func (o *ActionOptimizerImpl) GetMetrics(action *spookyactionstypes.Action) (*spookyactionstypes.PerformanceMetrics, error) {
	if action == nil {
		return nil, fmt.Errorf("action cannot be nil")
	}

	o.logger.Info("Getting metrics for action", spookylogging.String("action", action.Name))

	// Create basic metrics
	metrics := &spookyactionstypes.PerformanceMetrics{
		ActionName:     action.Name,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ExecutionCount: 0,
		SuccessCount:   0,
		FailureCount:   0,
		SuccessRate:    0.0,
		Metadata:       make(map[string]interface{}),
	}

	o.logger.Info("Successfully retrieved metrics for action", spookylogging.String("action", action.Name))
	return metrics, nil
}

// SetLevel sets the optimization level
func (o *ActionOptimizerImpl) SetLevel(level spookyactionstypes.OptimizationLevel) {
	o.level = level
	o.logger.Debug("Set optimization level", spookylogging.String("level", string(level)))
}

// GetLevel gets the current optimization level
func (o *ActionOptimizerImpl) GetLevel() spookyactionstypes.OptimizationLevel {
	return o.level
}

// SetResourceLimits sets the resource limits
func (o *ActionOptimizerImpl) SetResourceLimits(limits *spookyactionstypes.ResourceLimits) {
	o.resourceLimits = limits
	o.logger.Debug("Set resource limits")
}

// SetOptimizationTarget sets the optimization target
func (o *ActionOptimizerImpl) SetOptimizationTarget(target spookyactionstypes.OptimizationTarget) {
	o.optimizationTarget = target
	o.logger.Debug("Set optimization target", spookylogging.String("target", string(target)))
}

// applyBasicOptimizations applies basic optimizations to an action
func (o *ActionOptimizerImpl) applyBasicOptimizations(action *spookyactionstypes.Action) error {
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
func (o *ActionOptimizerImpl) applyAdvancedOptimizations(action *spookyactionstypes.Action) error {
	// Apply basic optimizations first
	if err := o.applyBasicOptimizations(action); err != nil {
		return err
	}

	// Optimize based on target
	switch o.optimizationTarget {
	case spookyactionstypes.OptimizationTargetSpeed:
		// Optimize for speed
		action.Parallel = true
		if action.MaxConcurrent < 20 {
			action.MaxConcurrent = 20
		}

	case spookyactionstypes.OptimizationTargetMemory:
		// Optimize for memory usage
		if action.MaxConcurrent > 5 {
			action.MaxConcurrent = 5
		}
		action.Parallel = false

	case spookyactionstypes.OptimizationTargetCPU:
		// Optimize for CPU usage
		if action.MaxConcurrent > 3 {
			action.MaxConcurrent = 3
		}
		action.Parallel = false

	case spookyactionstypes.OptimizationTargetNetwork:
		// Optimize for network usage
		if action.MaxConcurrent > 2 {
			action.MaxConcurrent = 2
		}
		action.Parallel = false

	case spookyactionstypes.OptimizationTargetBalanced:
		// Optimize for balanced usage
		if action.MaxConcurrent == 0 {
			action.MaxConcurrent = 10
		}
		action.Parallel = true
	}

	return nil
}

// applyMaximumOptimizations applies maximum optimizations to an action
func (o *ActionOptimizerImpl) applyMaximumOptimizations(action *spookyactionstypes.Action) error {
	// Apply advanced optimizations first
	if err := o.applyAdvancedOptimizations(action); err != nil {
		return err
	}

	// Maximum optimizations
	action.Parallel = true
	action.MaxConcurrent = 50

	// Set resource limits if not specified
	if action.ResourceLimits == nil {
		action.ResourceLimits = &spookyactionstypes.ActionResourceLimits{
			MemoryMB:   1024,
			CPUPercent: 50,
			DiskMB:     100,
		}
	}

	return nil
}
