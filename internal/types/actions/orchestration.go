package actions

import (
	"time"
)

// OrchestrationPlan represents the complete orchestration strategy
type OrchestrationPlan struct {
	Actions              []*Action
	OrchestrationOrder   []string
	ParallelGroups       [][]string // Groups that can run in parallel
	SequentialOrder      []string   // Actions that must run sequentially
	ConcurrencyLevel     int        // Optimal concurrency for parallel groups
	EstimatedDuration    time.Duration
	ResourceRequirements map[string]interface{}
	RollbackStrategy     *RollbackStrategy
}

// RollbackStrategy defines how to handle failures
type RollbackStrategy struct {
	EnableRollback  bool
	RollbackOrder   []string // Reverse orchestration order for rollback
	CriticalActions []string // Actions that trigger rollback on failure
}

// OrchestrationResult represents the result of action orchestration
type OrchestrationResult struct {
	Status    OrchestrationStatus    `json:"status"`
	StartTime time.Time              `json:"start_time"`
	EndTime   time.Time              `json:"end_time"`
	Duration  time.Duration          `json:"duration"`
	Output    string                 `json:"output"`
	Error     string                 `json:"error"`
	ExitCode  int                    `json:"exit_code"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// OrchestrationStatus represents the status of orchestration
type OrchestrationStatus string

const (
	OrchestrationStatusPending   OrchestrationStatus = "pending"
	OrchestrationStatusRunning   OrchestrationStatus = "running"
	OrchestrationStatusCompleted OrchestrationStatus = "completed"
	OrchestrationStatusFailed    OrchestrationStatus = "failed"
	OrchestrationStatusCancelled OrchestrationStatus = "cancelled"
)
