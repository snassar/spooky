package merging

import (
	"fmt"
	"time"

	spookytypes "spooky/internal/types"
	spookylogging "spooky/internal/logging"
)

// ActionMergerImpl implements the ActionMerger interface
type ActionMergerImpl struct {
	policy             spookyactionstypes.MergePolicy
	conflictResolution spookyactionstypes.ConflictResolution
	mergeStrategy      spookyactionstypes.MergeStrategy
	logger             spookylogging.Logger
}

// NewActionMerger creates a new ActionMerger
func NewActionMerger(policy spookyactionstypes.MergePolicy, logger spookylogging.Logger) ActionMerger {
	return &ActionMergerImpl{
		policy:             policy,
		conflictResolution: spookyactionstypes.ConflictResolutionFirst,
		mergeStrategy:      spookyactionstypes.MergeStrategyAppend,
		logger:             logger,
	}
}

// Merge merges multiple actions into a collection
func (m *ActionMergerImpl) Merge(actions ...*spookyactionstypes.Action) (*spookyactionstypes.ActionCollection, error) {
	if len(actions) == 0 {
		return nil, fmt.Errorf("no actions provided for merging")
	}

	m.logger.Info("Merging actions", spookylogging.Int("actions_count", len(actions)))

	// Create a new collection
	collection := &spookyactionstypes.ActionCollection{
		Actions:   make([]*spookyactionstypes.Action, 0, len(actions)),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata:  make(map[string]interface{}),
	}

	// Add actions to collection based on strategy
	switch m.mergeStrategy {
	case spookyactionstypes.MergeStrategyAppend:
		// Simply append all actions
		collection.Actions = append(collection.Actions, actions...)

	case spookyactionstypes.MergeStrategyReplace:
		// Replace with new actions
		collection.Actions = append(collection.Actions, actions...)

	case spookyactionstypes.MergeStrategyMerge:
		// Merge actions by name, resolving conflicts
		actionMap := make(map[string]*spookyactionstypes.Action)
		for _, action := range actions {
			if existing, exists := actionMap[action.Name]; exists {
				// Conflict resolution
				resolvedAction := m.resolveConflict(existing, action)
				actionMap[action.Name] = resolvedAction
			} else {
				actionMap[action.Name] = action
			}
		}

		// Convert map back to slice
		for _, action := range actionMap {
			collection.Actions = append(collection.Actions, action)
		}
	}

	m.logger.Info("Successfully merged actions", spookylogging.Int("actions_count", len(collection.Actions)))
	return collection, nil
}

// MergeWithPolicy merges action collections with a specific policy
func (m *ActionMergerImpl) MergeWithPolicy(existing, new *spookyactionstypes.ActionCollection, policy spookyactionstypes.MergePolicy) (*spookyactionstypes.ActionCollection, error) {
	if existing == nil {
		return nil, fmt.Errorf("existing collection cannot be nil")
	}

	if new == nil {
		return nil, fmt.Errorf("new collection cannot be nil")
	}

	m.logger.Info("Merging collections with policy",
		spookylogging.String("policy", policy.PolicyName),
		spookylogging.Int("existing_count", len(existing.Actions)),
		spookylogging.Int("new_count", len(new.Actions)))

	// Create a new collection
	collection := &spookyactionstypes.ActionCollection{
		Actions:   make([]*spookyactionstypes.Action, 0, len(existing.Actions)+len(new.Actions)),
		CreatedAt: existing.CreatedAt,
		UpdatedAt: time.Now(),
		Metadata:  make(map[string]interface{}),
	}

	// Merge based on policy strategy
	switch policy.Strategy {
	case spookyactionstypes.MergeStrategyAppend:
		// Append new actions to existing
		collection.Actions = append(collection.Actions, existing.Actions...)
		collection.Actions = append(collection.Actions, new.Actions...)

	case spookyactionstypes.MergeStrategyReplace:
		// Replace existing with new
		collection.Actions = append(collection.Actions, new.Actions...)

	case spookyactionstypes.MergeStrategyMerge:
		// Merge by name, resolving conflicts
		actionMap := make(map[string]*spookyactionstypes.Action)

		// Add existing actions
		for _, action := range existing.Actions {
			actionMap[action.Name] = action
		}

		// Merge new actions
		for _, action := range new.Actions {
			if existingAction, exists := actionMap[action.Name]; exists {
				// Conflict resolution
				resolvedAction := m.resolveConflict(existingAction, action)
				actionMap[action.Name] = resolvedAction
			} else {
				actionMap[action.Name] = action
			}
		}

		// Convert map back to slice
		for _, action := range actionMap {
			collection.Actions = append(collection.Actions, action)
		}
	}

	m.logger.Info("Successfully merged collections with policy",
		spookylogging.String("policy", policy.PolicyName),
		spookylogging.Int("result_count", len(collection.Actions)))

	return collection, nil
}

// SetPolicy sets the merge policy
func (m *ActionMergerImpl) SetPolicy(policy spookyactionstypes.MergePolicy) {
	m.policy = policy
	m.logger.Debug("Set merge policy", spookylogging.String("policy", policy.PolicyName))
}

// GetPolicy gets the current merge policy
func (m *ActionMergerImpl) GetPolicy() spookyactionstypes.MergePolicy {
	return m.policy
}

// SetConflictResolution sets the conflict resolution strategy
func (m *ActionMergerImpl) SetConflictResolution(resolution spookyactionstypes.ConflictResolution) {
	m.conflictResolution = resolution
	m.logger.Debug("Set conflict resolution", spookylogging.String("resolution", string(resolution)))
}

// SetMergeStrategy sets the merge strategy
func (m *ActionMergerImpl) SetMergeStrategy(strategy spookyactionstypes.MergeStrategy) {
	m.mergeStrategy = strategy
	m.logger.Debug("Set merge strategy", spookylogging.String("strategy", string(strategy)))
}

// resolveConflict resolves conflicts between two actions
func (m *ActionMergerImpl) resolveConflict(existing, new *spookyactionstypes.Action) *spookyactionstypes.Action {
	switch m.conflictResolution {
	case spookyactionstypes.ConflictResolutionFirst:
		// Keep the first (existing) action
		return existing

	case spookyactionstypes.ConflictResolutionLast:
		// Keep the last (new) action
		return new

	case spookyactionstypes.ConflictResolutionMerge:
		// Merge the actions (simplified - just take the new one)
		return new

	case spookyactionstypes.ConflictResolutionSkip:
		// Skip the conflict (return existing)
		return existing

	case spookyactionstypes.ConflictResolutionError:
		// This should be handled by the caller
		// For now, return the existing action
		return existing

	default:
		// Default to keeping the first
		return existing
	}
}
