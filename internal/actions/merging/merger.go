package merging

import (
	"fmt"
	"time"

	"spooky/internal/actions/types"
	"spooky/internal/logging"
)

// ActionMergerImpl implements the ActionMerger interface
type ActionMergerImpl struct {
	policy             types.MergePolicy
	conflictResolution types.ConflictResolution
	mergeStrategy      types.MergeStrategy
	logger             logging.Logger
}

// NewActionMerger creates a new ActionMerger
func NewActionMerger(policy types.MergePolicy, logger logging.Logger) ActionMerger {
	return &ActionMergerImpl{
		policy:             policy,
		conflictResolution: types.ConflictResolutionFirst,
		mergeStrategy:      types.MergeStrategyAppend,
		logger:             logger,
	}
}

// Merge merges multiple actions into a collection
func (m *ActionMergerImpl) Merge(actions ...*types.Action) (*types.ActionCollection, error) {
	if len(actions) == 0 {
		return nil, fmt.Errorf("no actions provided for merging")
	}

	m.logger.Info("Merging actions", logging.Int("actions_count", len(actions)))

	// Create a new collection
	collection := &types.ActionCollection{
		Actions:   make([]*types.Action, 0, len(actions)),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata:  make(map[string]interface{}),
	}

	// Add actions to collection based on strategy
	switch m.mergeStrategy {
	case types.MergeStrategyAppend:
		// Simply append all actions
		collection.Actions = append(collection.Actions, actions...)

	case types.MergeStrategyReplace:
		// Replace with new actions
		collection.Actions = append(collection.Actions, actions...)

	case types.MergeStrategyMerge:
		// Merge actions by name, resolving conflicts
		actionMap := make(map[string]*types.Action)
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

	m.logger.Info("Successfully merged actions", logging.Int("actions_count", len(collection.Actions)))
	return collection, nil
}

// MergeWithPolicy merges action collections with a specific policy
func (m *ActionMergerImpl) MergeWithPolicy(existing, new *types.ActionCollection, policy types.MergePolicy) (*types.ActionCollection, error) {
	if existing == nil {
		return nil, fmt.Errorf("existing collection cannot be nil")
	}

	if new == nil {
		return nil, fmt.Errorf("new collection cannot be nil")
	}

	m.logger.Info("Merging collections with policy",
		logging.String("policy", policy.PolicyName),
		logging.Int("existing_count", len(existing.Actions)),
		logging.Int("new_count", len(new.Actions)))

	// Create a new collection
	collection := &types.ActionCollection{
		Actions:   make([]*types.Action, 0, len(existing.Actions)+len(new.Actions)),
		CreatedAt: existing.CreatedAt,
		UpdatedAt: time.Now(),
		Metadata:  make(map[string]interface{}),
	}

	// Merge based on policy strategy
	switch policy.Strategy {
	case types.MergeStrategyAppend:
		// Append new actions to existing
		collection.Actions = append(collection.Actions, existing.Actions...)
		collection.Actions = append(collection.Actions, new.Actions...)

	case types.MergeStrategyReplace:
		// Replace existing with new
		collection.Actions = append(collection.Actions, new.Actions...)

	case types.MergeStrategyMerge:
		// Merge by name, resolving conflicts
		actionMap := make(map[string]*types.Action)

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
		logging.String("policy", policy.PolicyName),
		logging.Int("result_count", len(collection.Actions)))

	return collection, nil
}

// SetPolicy sets the merge policy
func (m *ActionMergerImpl) SetPolicy(policy types.MergePolicy) {
	m.policy = policy
	m.logger.Debug("Set merge policy", logging.String("policy", policy.PolicyName))
}

// GetPolicy gets the current merge policy
func (m *ActionMergerImpl) GetPolicy() types.MergePolicy {
	return m.policy
}

// SetConflictResolution sets the conflict resolution strategy
func (m *ActionMergerImpl) SetConflictResolution(resolution types.ConflictResolution) {
	m.conflictResolution = resolution
	m.logger.Debug("Set conflict resolution", logging.String("resolution", string(resolution)))
}

// SetMergeStrategy sets the merge strategy
func (m *ActionMergerImpl) SetMergeStrategy(strategy types.MergeStrategy) {
	m.mergeStrategy = strategy
	m.logger.Debug("Set merge strategy", logging.String("strategy", string(strategy)))
}

// resolveConflict resolves conflicts between two actions
func (m *ActionMergerImpl) resolveConflict(existing, new *types.Action) *types.Action {
	switch m.conflictResolution {
	case types.ConflictResolutionFirst:
		// Keep the first (existing) action
		return existing

	case types.ConflictResolutionLast:
		// Keep the last (new) action
		return new

	case types.ConflictResolutionMerge:
		// Merge the actions (simplified - just take the new one)
		return new

	case types.ConflictResolutionSkip:
		// Skip the conflict (return existing)
		return existing

	case types.ConflictResolutionError:
		// This should be handled by the caller
		// For now, return the existing action
		return existing

	default:
		// Default to keeping the first
		return existing
	}
}
