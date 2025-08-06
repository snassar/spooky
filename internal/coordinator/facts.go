package coordinator

import (
	"fmt"
	"io"
	"time"

	"spooky/internal/facts"
	"spooky/internal/facts/types"
	"spooky/internal/interfaces"
	"spooky/internal/logging"
)

// CoordinatorFactsIntegration implements facts system integration
type CoordinatorFactsIntegration struct {
	factsManager facts.FactManager
	logger       logging.Logger
}

// NewCoordinatorFactsIntegration creates a new facts integration
func NewCoordinatorFactsIntegration(factsManager facts.FactManager, logger logging.Logger) *CoordinatorFactsIntegration {
	return &CoordinatorFactsIntegration{
		factsManager: factsManager,
		logger:       logger,
	}
}

// LoadFacts loads facts for the specified machines
func (fi *CoordinatorFactsIntegration) LoadFacts(machineNames []string) (*interfaces.FactsContext, error) {
	context := &interfaces.FactsContext{
		BaseContext: interfaces.BaseContext{
			Timestamp: time.Now(),
		},
		MachineFacts: make(map[string]*types.FactCollection),
	}

	// Load machine-specific facts
	if fi.factsManager != nil {
		for _, machine := range machineNames {
			facts, err := fi.factsManager.GetFactCollection(machine)
			if err != nil {
				fi.logger.Warn("Failed to load facts for machine",
					logging.String("machine", machine),
					logging.Error(err))
				continue
			}
			context.MachineFacts[machine] = facts
		}

		// Load global facts if available
		if globalFacts, err := fi.factsManager.LoadPersistedFacts("global"); err == nil {
			context.GlobalFacts = globalFacts
		}

		// Load project facts if available
		if projectFacts, err := fi.factsManager.LoadPersistedFacts("project"); err == nil {
			context.ProjectFacts = projectFacts
		}
	}

	// Generate cache key
	context.CacheKey = fi.generateCacheKey(machineNames)

	return context, nil
}

// CollectFacts collects facts from the specified machines
func (fi *CoordinatorFactsIntegration) CollectFacts(machineNames []string) (*interfaces.FactsContext, error) {
	// Validate inputs
	if err := interfaces.ValidateMachineNames(machineNames); err != nil {
		return nil, err
	}

	fi.logger.Info("Collecting facts from machines", logging.StringSlice("machines", machineNames))

	// Create facts context
	context := &interfaces.FactsContext{
		MachineFacts: make(map[string]*types.FactCollection),
		GlobalFacts:  nil,
		ProjectFacts: nil,
	}

	// Collect facts from each machine
	if fi.factsManager != nil {
		for _, machine := range machineNames {
			fi.logger.Info("Collecting facts from machine", logging.String("machine", machine))

			// Collect facts using the facts manager
			facts, err := fi.factsManager.CollectAllFacts(machine)
			if err != nil {
				fi.logger.Warn("Failed to collect facts for machine",
					logging.String("machine", machine),
					logging.Error(err))
				continue
			}

			context.MachineFacts[machine] = facts
			fi.logger.Info("Successfully collected facts from machine",
				logging.String("machine", machine),
				logging.Int("facts_count", len(facts.Facts)))
		}
	}

	// Generate cache key
	context.CacheKey = fi.generateCacheKey(machineNames)

	fi.logger.Info("Completed facts collection",
		logging.Int("machines_processed", len(context.MachineFacts)),
		logging.Int("total_machines", len(machineNames)))

	return context, nil
}

// ValidateFacts validates facts data integrity
func (fi *CoordinatorFactsIntegration) ValidateFacts(factsContext *interfaces.FactsContext) error {
	if factsContext == nil {
		return fmt.Errorf("facts context cannot be nil")
	}

	// Validate machine facts
	for machine, factCollection := range factsContext.MachineFacts {
		if factCollection == nil {
			return fmt.Errorf("fact collection for machine '%s' is nil", machine)
		}
		if factCollection.Facts == nil {
			return fmt.Errorf("facts for machine '%s' are nil", machine)
		}
	}

	// Validate global facts if present
	if factsContext.GlobalFacts != nil && factsContext.GlobalFacts.Facts == nil {
		return fmt.Errorf("global facts are nil")
	}

	// Validate project facts if present
	if factsContext.ProjectFacts != nil && factsContext.ProjectFacts.Facts == nil {
		return fmt.Errorf("project facts are nil")
	}

	return nil
}

// OptimizeFactsGathering optimizes facts gathering performance
func (fi *CoordinatorFactsIntegration) OptimizeFactsGathering(machineNames []string, parallel int) (int, time.Duration, error) {
	// Validate inputs
	if err := interfaces.ValidateMachineNames(machineNames); err != nil {
		return 0, 0, err
	}

	if err := interfaces.ValidateParallelWorkers(parallel); err != nil {
		return 0, 0, err
	}

	// Calculate optimal parallel workers
	optimalParallel := interfaces.CalculateOptimalParallel(parallel, len(machineNames), interfaces.DefaultOptimizationConfig())

	// Calculate optimal timeout
	optimalTimeout := interfaces.CalculateOptimalTimeout(len(machineNames), optimalParallel, interfaces.DefaultOptimizationConfig())

	return optimalParallel, optimalTimeout, nil
}

// GetFactsForMachine gets facts for a specific machine
func (fi *CoordinatorFactsIntegration) GetFactsForMachine(machine string) (*types.FactCollection, error) {
	if fi.factsManager == nil {
		return &types.FactCollection{}, nil
	}
	return fi.factsManager.GetFactCollection(machine)
}

// CacheFacts caches facts for later use
func (fi *CoordinatorFactsIntegration) CacheFacts(factsContext *interfaces.FactsContext) error {
	if factsContext == nil {
		return fmt.Errorf("facts context cannot be nil")
	}

	// If facts manager is nil, just return success (no caching)
	if fi.factsManager == nil {
		return nil
	}

	// Cache machine facts
	for machine, factCollection := range factsContext.MachineFacts {
		if factCollection != nil {
			if err := fi.factsManager.SetFactCollection(machine, factCollection); err != nil {
				fi.logger.Warn("Failed to cache facts for machine",
					logging.String("machine", machine),
					logging.Error(err))
			}
		}
	}

	// Cache global facts if present
	if factsContext.GlobalFacts != nil {
		if err := fi.factsManager.PersistFacts("global", factsContext.GlobalFacts); err != nil {
			fi.logger.Warn("Failed to cache global facts", logging.Error(err))
		}
	}

	// Cache project facts if present
	if factsContext.ProjectFacts != nil {
		if err := fi.factsManager.PersistFacts("project", factsContext.ProjectFacts); err != nil {
			fi.logger.Warn("Failed to cache project facts", logging.Error(err))
		}
	}

	return nil
}

// CollectFactsForAction collects facts needed for action execution
func (fi *CoordinatorFactsIntegration) CollectFactsForAction(_ interface{}, machines []string) (*interfaces.FactsContext, error) {
	return fi.LoadFacts(machines)
}

// ValidateActionWithFacts validates an action using facts data
func (fi *CoordinatorFactsIntegration) ValidateActionWithFacts(action interface{}, factsContext *interfaces.FactsContext) error {
	// Basic validation - ensure facts context is valid
	if err := fi.ValidateFacts(factsContext); err != nil {
		return fmt.Errorf("facts validation failed: %w", err)
	}

	// Action-specific facts validation
	if action != nil {
		// Check if action requires specific facts
		// This could include:
		// - Required machine facts for action execution
		// - Required global facts for action configuration
		// - Required project facts for action context

		// For now, we'll do basic validation
		// In a real implementation, this would check action metadata for required facts

		// Example validation: ensure we have facts for all machines
		if len(factsContext.MachineFacts) == 0 {
			fi.logger.Warn("No machine facts available for action validation")
		}

		// Example validation: check for required global facts
		if factsContext.GlobalFacts == nil {
			fi.logger.Warn("No global facts available for action validation")
		}

		// Example validation: check for required project facts
		if factsContext.ProjectFacts == nil {
			fi.logger.Warn("No project facts available for action validation")
		}
	}

	return nil
}

// GetFactValue gets a specific fact value from the context
func (fi *CoordinatorFactsIntegration) GetFactValue(factKey string, factsContext *interfaces.FactsContext) (interface{}, error) {
	if factsContext == nil {
		return nil, fmt.Errorf("facts context cannot be nil")
	}

	// Try to get fact from machine facts first
	for machine, factCollection := range factsContext.MachineFacts {
		if fact, exists := factCollection.Facts[factKey]; exists {
			fi.logger.Debug("Found fact in machine facts",
				logging.String("fact_key", factKey),
				logging.String("machine", machine))
			return fact.Value, nil
		}
	}

	// Try global facts
	if factsContext.GlobalFacts != nil {
		if fact, exists := factsContext.GlobalFacts.Facts[factKey]; exists {
			fi.logger.Debug("Found fact in global facts",
				logging.String("fact_key", factKey))
			return fact.Value, nil
		}
	}

	// Try project facts
	if factsContext.ProjectFacts != nil {
		if fact, exists := factsContext.ProjectFacts.Facts[factKey]; exists {
			fi.logger.Debug("Found fact in project facts",
				logging.String("fact_key", factKey))
			return fact.Value, nil
		}
	}

	return nil, fmt.Errorf("fact '%s' not found in any context", factKey)
}

// ExportFacts exports facts to JSON format with encryption support
func (fi *CoordinatorFactsIntegration) ExportFacts(w io.Writer, opts facts.ExportOptions) error {
	if fi.factsManager == nil {
		return fmt.Errorf("facts manager is not available")
	}

	// Use the facts manager's export functionality
	return fi.factsManager.ExportFacts(w)
}

// ExportToHCL exports facts to HCL format
func (fi *CoordinatorFactsIntegration) ExportToHCL(w io.Writer, query *facts.FactQuery) error {
	if fi.factsManager == nil {
		return fmt.Errorf("facts manager is not available")
	}

	// Create an exporter using the facts manager's storage
	storage := fi.factsManager.GetStorage()
	if storage == nil {
		return fmt.Errorf("facts storage is not available")
	}

	exporter := facts.NewExporter(storage)
	return exporter.ExportToHCL(w, query)
}

// ImportFacts imports facts from JSON format
func (fi *CoordinatorFactsIntegration) ImportFacts(r io.Reader) error {
	if fi.factsManager == nil {
		return fmt.Errorf("facts manager is not available")
	}

	// Use the facts manager's import functionality
	return fi.factsManager.ImportFacts(r)
}

// ImportFromHCL imports facts from HCL format
func (fi *CoordinatorFactsIntegration) ImportFromHCL(r io.Reader) error {
	if fi.factsManager == nil {
		return fmt.Errorf("facts manager is not available")
	}

	// Get the storage from the facts manager
	storage := fi.factsManager.GetStorage()
	if storage == nil {
		return fmt.Errorf("facts storage is not available")
	}

	// Use the storage's HCL import functionality
	return storage.ImportFromHCL(r)
}

// generateCacheKey generates a cache key for the given machine names
func (fi *CoordinatorFactsIntegration) generateCacheKey(machineNames []string) string {
	// Simple cache key generation - could be enhanced with hashing
	key := fmt.Sprintf("facts_%d", len(machineNames))
	for _, machine := range machineNames {
		key += "_" + machine
	}
	return key
}
