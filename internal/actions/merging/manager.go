package merging

import (
	"fmt"
	"sync"

	"spooky/internal/actions/types"
	"spooky/internal/logging"
)

// Manager implements the MergingManager interface
type Manager struct {
	// Configuration
	defaultPolicy      types.MergePolicy
	conflictResolution types.ConflictResolution

	// State
	mergers  map[string]ActionMerger
	policies map[string]types.MergePolicy
	logger   logging.Logger
	mu       sync.RWMutex
}

// NewManager creates a new MergingManager
func NewManager(logger logging.Logger) *Manager {
	return &Manager{
		defaultPolicy:      types.MergePolicy{PolicyName: "default"},
		conflictResolution: types.ConflictResolutionFirst,
		mergers:            make(map[string]ActionMerger),
		policies:           make(map[string]types.MergePolicy),
		logger:             logger,
	}
}

// MergeActions merges multiple actions into a collection
func (m *Manager) MergeActions(actions ...*types.Action) (*types.ActionCollection, error) {
	if len(actions) == 0 {
		return nil, fmt.Errorf("no actions provided for merging")
	}

	m.logger.Info("Merging actions", logging.Int("actions_count", len(actions)))

	// Create a merger with default policy
	merger, err := m.CreateMerger(m.defaultPolicy)
	if err != nil {
		return nil, fmt.Errorf("failed to create merger: %w", err)
	}

	// Merge the actions
	collection, err := merger.Merge(actions...)
	if err != nil {
		return nil, fmt.Errorf("failed to merge actions: %w", err)
	}

	m.logger.Info("Successfully merged actions", logging.Int("actions_count", len(collection.Actions)))
	return collection, nil
}

// MergeWithPolicy merges action collections with a specific policy
func (m *Manager) MergeWithPolicy(existing, new *types.ActionCollection, policy types.MergePolicy) (*types.ActionCollection, error) {
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

	// Create a merger with the specified policy
	merger, err := m.CreateMerger(policy)
	if err != nil {
		return nil, fmt.Errorf("failed to create merger: %w", err)
	}

	// Merge the collections
	collection, err := merger.MergeWithPolicy(existing, new, policy)
	if err != nil {
		return nil, fmt.Errorf("failed to merge collections: %w", err)
	}

	m.logger.Info("Successfully merged collections with policy",
		logging.String("policy", policy.PolicyName),
		logging.Int("result_count", len(collection.Actions)))

	return collection, nil
}

// CreateMerger creates a new merger with a specific policy
func (m *Manager) CreateMerger(policy types.MergePolicy) (ActionMerger, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if merger already exists for this policy
	policyKey := policy.PolicyName
	if merger, exists := m.mergers[policyKey]; exists {
		return merger, nil
	}

	// Create new merger
	merger := NewActionMerger(policy, m.logger)
	m.mergers[policyKey] = merger

	m.logger.Debug("Created merger for policy", logging.String("policy", policy.PolicyName))
	return merger, nil
}

// GetMerger gets an existing merger for a policy
func (m *Manager) GetMerger(policy types.MergePolicy) (ActionMerger, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	policyKey := policy.PolicyName
	merger, exists := m.mergers[policyKey]
	if !exists {
		return nil, fmt.Errorf("merger not found for policy %s", policy.PolicyName)
	}

	return merger, nil
}

// GetMergePolicy gets a merge policy by name
func (m *Manager) GetMergePolicy(name string) (types.MergePolicy, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	policy, exists := m.policies[name]
	if !exists {
		return types.MergePolicy{}, fmt.Errorf("merge policy not found: %s", name)
	}

	return policy, nil
}

// ListMergePolicies lists all available merge policies
func (m *Manager) ListMergePolicies() ([]types.MergePolicy, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	policies := make([]types.MergePolicy, 0, len(m.policies))
	for _, policy := range m.policies {
		policies = append(policies, policy)
	}

	return policies, nil
}

// AddMergePolicy adds a new merge policy
func (m *Manager) AddMergePolicy(policy types.MergePolicy) error {
	if policy.PolicyName == "" {
		return fmt.Errorf("policy name cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.policies[policy.PolicyName] = policy
	m.logger.Info("Added merge policy", logging.String("policy", policy.PolicyName))
	return nil
}

// RemoveMergePolicy removes a merge policy by name
func (m *Manager) RemoveMergePolicy(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.policies[name]; !exists {
		return fmt.Errorf("merge policy not found: %s", name)
	}

	delete(m.policies, name)
	m.logger.Info("Removed merge policy", logging.String("policy", name))
	return nil
}

// SetDefaultPolicy sets the default merge policy
func (m *Manager) SetDefaultPolicy(policy types.MergePolicy) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.defaultPolicy = policy
	m.logger.Info("Set default merge policy", logging.String("policy", policy.PolicyName))
}

// SetConflictResolution sets the conflict resolution strategy
func (m *Manager) SetConflictResolution(resolution types.ConflictResolution) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.conflictResolution = resolution
	m.logger.Info("Set conflict resolution", logging.String("resolution", string(resolution)))
}
