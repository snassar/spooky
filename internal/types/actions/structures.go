package spookytypesactions

import (
	spookytypescommon "spooky/internal/types/common"
	"time"
)

// ActionCollection represents a collection of actions
type ActionCollection struct {
	spookytypescommon.CompleteEntity

	// Collection metadata
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`

	// Actions
	Actions map[string]*Action `json:"actions"`

	// Collection settings
	DefaultTimeout  int  `json:"default_timeout"`
	DefaultParallel bool `json:"default_parallel"`
	DefaultRetries  int  `json:"default_retries"`

	// Metadata
	Tags     []string          `json:"tags,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ActionPlan represents a plan for running actions
type ActionPlan struct {
	spookytypescommon.TimestampedEntity

	// Plan identification
	PlanID      string `json:"plan_id"`
	PlanName    string `json:"plan_name"`
	Description string `json:"description"`

	// Actions to run
	Actions []*Action `json:"actions"`

	// Execution order
	ExecutionOrder [][]string `json:"execution_order"`

	// Dependencies
	Dependencies map[string][]string `json:"dependencies"`

	// Plan settings
	Parallel     bool `json:"parallel"`
	AllowFailure bool `json:"allow_failure"`
	DryRun       bool `json:"dry_run"`

	// Validation
	Validated bool     `json:"validated"`
	Errors    []string `json:"errors,omitempty"`
	Warnings  []string `json:"warnings,omitempty"`
}

// ActionDependency represents a dependency between actions
type ActionDependency struct {
	spookytypescommon.TimestampedEntity

	// Dependency identification
	FromAction string `json:"from_action"`
	ToAction   string `json:"to_action"`

	// Dependency type
	Type string `json:"type"` // "required", "optional", "conditional"

	// Condition (for conditional dependencies)
	Condition string `json:"condition,omitempty"`

	// Metadata
	Description string            `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// ActionExecution represents the execution of a single action
type ActionExecution struct {
	spookytypescommon.TimestampedEntity

	// Execution identification
	ExecutionID string `json:"execution_id"`
	ActionName  string `json:"action_name"`

	// Execution context
	Context *ActionExecutionContext `json:"context"`

	// Execution state
	Status    string     `json:"status"`
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`

	// Results
	Results []ActingResult `json:"results,omitempty"`

	// Error handling
	Error      error `json:"error,omitempty"`
	RetryCount int   `json:"retry_count"`
	MaxRetries int   `json:"max_retries"`
}

// ActionValidation represents validation results for actions
type ActionValidation struct {
	spookytypescommon.TimestampedEntity

	// Validation identification
	ValidationID string `json:"validation_id"`
	ActionName   string `json:"action_name"`

	// Validation results
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`

	// Validation details
	SchemaValid      bool     `json:"schema_valid"`
	SchemaErrors     []string `json:"schema_errors,omitempty"`
	DependencyValid  bool     `json:"dependency_valid"`
	DependencyErrors []string `json:"dependency_errors,omitempty"`
	ResourceValid    bool     `json:"resource_valid"`
	ResourceErrors   []string `json:"resource_errors,omitempty"`

	// Metadata
	ValidationTime *time.Time        `json:"validation_time,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// ActionMetrics represents metrics for action execution
type ActionMetrics struct {
	spookytypescommon.TimestampedEntity

	// Metrics identification
	ActionName string `json:"action_name"`

	// Execution metrics
	TotalExecutions      int     `json:"total_executions"`
	SuccessfulExecutions int     `json:"successful_executions"`
	FailedExecutions     int     `json:"failed_executions"`
	SuccessRate          float64 `json:"success_rate"`

	// Timing metrics
	AverageDuration time.Duration `json:"average_duration"`
	MinDuration     time.Duration `json:"min_duration"`
	MaxDuration     time.Duration `json:"max_duration"`

	// Resource metrics
	AverageMemoryUsage int64   `json:"average_memory_usage"`
	AverageCPUUsage    float64 `json:"average_cpu_usage"`

	// Retry metrics
	TotalRetries   int     `json:"total_retries"`
	AverageRetries float64 `json:"average_retries"`

	// Machine metrics
	MachinesExecuted int `json:"machines_executed"`
	MachinesFailed   int `json:"machines_failed"`

	// Metadata
	LastExecution *time.Time        `json:"last_execution,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}
