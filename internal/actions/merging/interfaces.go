package merging

import (
	spookyactionstypes "spooky/internal/actions/types"
)

// MergingManager defines the interface for action merging operations
type MergingManager interface {
	// Core merging operations
	MergeActions(actions ...*spookyactionstypes.Action) (*spookyactionstypes.ActionCollection, error)
	MergeWithPolicy(existing, new *spookyactionstypes.ActionCollection, policy spookyactionstypes.MergePolicy) (*spookyactionstypes.ActionCollection, error)

	// Merger management
	CreateMerger(policy spookyactionstypes.MergePolicy) (ActionMerger, error)
	GetMerger(policy spookyactionstypes.MergePolicy) (ActionMerger, error)

	// Policy management
	GetMergePolicy(name string) (spookyactionstypes.MergePolicy, error)
	ListMergePolicies() ([]spookyactionstypes.MergePolicy, error)
	AddMergePolicy(policy spookyactionstypes.MergePolicy) error
	RemoveMergePolicy(name string) error

	// Configuration
	SetDefaultPolicy(policy spookyactionstypes.MergePolicy)
	SetConflictResolution(resolution spookyactionstypes.ConflictResolution)
}

// ActionMerger defines the interface for action merging
type ActionMerger interface {
	// Core merging operations
	Merge(actions ...*spookyactionstypes.Action) (*spookyactionstypes.ActionCollection, error)
	MergeWithPolicy(existing, new *spookyactionstypes.ActionCollection, policy spookyactionstypes.MergePolicy) (*spookyactionstypes.ActionCollection, error)

	// Policy management
	SetPolicy(policy spookyactionstypes.MergePolicy)
	GetPolicy() spookyactionstypes.MergePolicy

	// Configuration
	SetConflictResolution(resolution spookyactionstypes.ConflictResolution)
	SetMergeStrategy(strategy spookyactionstypes.MergeStrategy)
}

// MergeResult represents the result of a merge operation
type MergeResult struct {
	Success       bool                         `json:"success"`
	MergedActions []*spookyactionstypes.Action `json:"merged_actions"`
	Conflicts     []MergeConflict              `json:"conflicts"`
	Warnings      []MergeWarning               `json:"warnings"`
	Details       map[string]interface{}       `json:"details"`
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
