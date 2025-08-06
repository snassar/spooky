package merging

import (
	"spooky/internal/actions/types"
)

// MergingManager defines the interface for action merging operations
type MergingManager interface {
	// Core merging operations
	MergeActions(actions ...*types.Action) (*types.ActionCollection, error)
	MergeWithPolicy(existing, new *types.ActionCollection, policy types.MergePolicy) (*types.ActionCollection, error)

	// Merger management
	CreateMerger(policy types.MergePolicy) (ActionMerger, error)
	GetMerger(policy types.MergePolicy) (ActionMerger, error)

	// Policy management
	GetMergePolicy(name string) (types.MergePolicy, error)
	ListMergePolicies() ([]types.MergePolicy, error)
	AddMergePolicy(policy types.MergePolicy) error
	RemoveMergePolicy(name string) error

	// Configuration
	SetDefaultPolicy(policy types.MergePolicy)
	SetConflictResolution(resolution types.ConflictResolution)
}

// ActionMerger defines the interface for action merging
type ActionMerger interface {
	// Core merging operations
	Merge(actions ...*types.Action) (*types.ActionCollection, error)
	MergeWithPolicy(existing, new *types.ActionCollection, policy types.MergePolicy) (*types.ActionCollection, error)

	// Policy management
	SetPolicy(policy types.MergePolicy)
	GetPolicy() types.MergePolicy

	// Configuration
	SetConflictResolution(resolution types.ConflictResolution)
	SetMergeStrategy(strategy types.MergeStrategy)
}

// MergeResult represents the result of a merge operation
type MergeResult struct {
	Success       bool                   `json:"success"`
	MergedActions []*types.Action        `json:"merged_actions"`
	Conflicts     []MergeConflict        `json:"conflicts"`
	Warnings      []MergeWarning         `json:"warnings"`
	Details       map[string]interface{} `json:"details"`
}

// MergeConflict represents a merge conflict
type MergeConflict struct {
	ActionName  string      `json:"action_name"`
	Field       string      `json:"field"`
	Value1      interface{} `json:"value_1"`
	Value2      interface{} `json:"value_2"`
	Resolution  string      `json:"resolution"`
	Description string      `json:"description"`
}

// MergeWarning represents a merge warning
type MergeWarning struct {
	ActionName string `json:"action_name"`
	Field      string `json:"field"`
	Message    string `json:"message"`
	Severity   string `json:"severity"`
}
