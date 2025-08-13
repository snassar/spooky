package spookytypesactions

import (
	"time"

	spookytypescommon "spooky/internal/types/common"
	spookytypesmachines "spooky/internal/types/machines"
)

// ActionCollection represents a collection of actions
type ActionCollection struct {
	spookytypescommon.CompleteEntity
	Name        string            `json:"name" hcl:"name"`
	Description string            `json:"description" hcl:"description"`
	Actions     []*Action         `json:"actions" hcl:"actions"`
	Tags        []string          `json:"tags" hcl:"tags,optional"`
	Metadata    map[string]string `json:"metadata" hcl:"metadata,optional"`
	Config      *CollectionConfig `json:"config" hcl:"config,optional"`
}

// CollectionConfig represents configuration for action collections
type CollectionConfig struct {
	Parallel       bool `json:"parallel" hcl:"parallel,optional"`
	MaxConcurrent  int  `json:"max_concurrent" hcl:"max_concurrent,optional"`
	Timeout        int  `json:"timeout" hcl:"timeout,optional"`
	AllowFailures  bool `json:"allow_failures" hcl:"allow_failures,optional"`
	StopOnFailure  bool `json:"stop_on_failure" hcl:"stop_on_failure,optional"`
	ValidateBefore bool `json:"validate_before" hcl:"validate_before,optional"`
}

// ActionPlan represents a plan for running actions
type ActionPlan struct {
	spookytypescommon.CompleteEntity
	Name           string                        `json:"name" hcl:"name"`
	Description    string                        `json:"description" hcl:"description"`
	Actions        []*Action                     `json:"actions" hcl:"actions"`
	Machines       []spookytypesmachines.Machine `json:"machines" hcl:"machines"`
	RunOrder       [][]string                    `json:"run_order" hcl:"run_order"`
	Dependencies   map[string][]string           `json:"dependencies" hcl:"dependencies"`
	Parallel       bool                          `json:"parallel" hcl:"parallel,optional"`
	MaxConcurrent  int                           `json:"max_concurrent" hcl:"max_concurrent,optional"`
	Timeout        int                           `json:"timeout" hcl:"timeout,optional"`
	AllowFailures  bool                          `json:"allow_failures" hcl:"allow_failures,optional"`
	StopOnFailure  bool                          `json:"stop_on_failure" hcl:"stop_on_failure,optional"`
	ValidateBefore bool                          `json:"validate_before" hcl:"validate_before,optional"`
	EstimatedTime  time.Duration                 `json:"estimated_time" hcl:"estimated_time,optional"`
	RiskLevel      string                        `json:"risk_level" hcl:"risk_level,optional"`
	RollbackPlan   *ActionPlan                   `json:"rollback_plan" hcl:"rollback_plan,optional"`
}

// ActionDependency represents a dependency between actions
type ActionDependency struct {
	spookytypescommon.CompleteEntity
	FromAction    string `json:"from_action" hcl:"from_action"`
	ToAction      string `json:"to_action" hcl:"to_action"`
	Type          string `json:"type" hcl:"type"` // "required", "optional", "conditional"
	Condition     string `json:"condition" hcl:"condition,optional"`
	Description   string `json:"description" hcl:"description,optional"`
	Timeout       int    `json:"timeout" hcl:"timeout,optional"`
	RetryCount    int    `json:"retry_count" hcl:"retry_count,optional"`
	RetryDelay    int    `json:"retry_delay" hcl:"retry_delay,optional"`
	FailureAction string `json:"failure_action" hcl:"failure_action,optional"` // "continue", "stop", "rollback"
}

// ActionRun represents the running of a single action
type ActionRun struct {
	spookytypescommon.CompleteEntity
	RunID         string                      `json:"run_id" hcl:"run_id"`
	Action        *Action                     `json:"action" hcl:"action"`
	Machine       spookytypesmachines.Machine `json:"machine" hcl:"machine"`
	Status        string                      `json:"status" hcl:"status"` // "pending", "running", "completed", "failed", "cancelled"
	StartTime     *time.Time                  `json:"start_time" hcl:"start_time,optional"`
	EndTime       *time.Time                  `json:"end_time" hcl:"end_time,optional"`
	Duration      time.Duration               `json:"duration" hcl:"duration,optional"`
	ExitCode      int                         `json:"exit_code" hcl:"exit_code,optional"`
	Stdout        string                      `json:"stdout" hcl:"stdout,optional"`
	Stderr        string                      `json:"stderr" hcl:"stderr,optional"`
	Error         string                      `json:"error" hcl:"error,optional"`
	RetryCount    int                         `json:"retry_count" hcl:"retry_count,optional"`
	MaxRetries    int                         `json:"max_retries" hcl:"max_retries,optional"`
	RetryDelay    time.Duration               `json:"retry_delay" hcl:"retry_delay,optional"`
	ResourceUsage *ResourceUsage              `json:"resource_usage" hcl:"resource_usage,optional"`
	Metadata      map[string]string           `json:"metadata" hcl:"metadata,optional"`
}

// ResourceUsage represents resource usage during action running
type ResourceUsage struct {
	spookytypescommon.CompleteEntity
	CPUPercent    float64 `json:"cpu_percent" hcl:"cpu_percent,optional"`
	MemoryMB      int     `json:"memory_mb" hcl:"memory_mb,optional"`
	DiskMB        int     `json:"disk_mb" hcl:"disk_mb,optional"`
	NetworkMB     int     `json:"network_mb" hcl:"network_mb,optional"`
	ProcessCount  int     `json:"process_count" hcl:"process_count,optional"`
	OpenFiles     int     `json:"open_files" hcl:"open_files,optional"`
	MaxMemoryMB   int     `json:"max_memory_mb" hcl:"max_memory_mb,optional"`
	MaxCPUPercent float64 `json:"max_cpu_percent" hcl:"max_cpu_percent,optional"`
}

// ActionMetrics represents metrics for action running
type ActionMetrics struct {
	spookytypescommon.CompleteEntity
	ActionName        string             `json:"action_name" hcl:"action_name"`
	TotalRuns         int                `json:"total_runs" hcl:"total_runs"`
	SuccessfulRuns    int                `json:"successful_runs" hcl:"successful_runs"`
	FailedRuns        int                `json:"failed_runs" hcl:"failed_runs"`
	SuccessRate       float64            `json:"success_rate" hcl:"success_rate"`
	AverageDuration   time.Duration      `json:"average_duration" hcl:"average_duration,optional"`
	MinDuration       time.Duration      `json:"min_duration" hcl:"min_duration,optional"`
	MaxDuration       time.Duration      `json:"max_duration" hcl:"max_duration,optional"`
	TotalDuration     time.Duration      `json:"total_duration" hcl:"total_duration,optional"`
	MachinesRun       int                `json:"machines_run" hcl:"machines_run"`
	LastRun           *time.Time         `json:"last_run" hcl:"last_run,optional"`
	LastSuccess       *time.Time         `json:"last_success" hcl:"last_success,optional"`
	LastFailure       *time.Time         `json:"last_failure" hcl:"last_failure,optional"`
	ConcurrentRuns    int                `json:"concurrent_runs" hcl:"concurrent_runs,optional"`
	MaxConcurrentRuns int                `json:"max_concurrent_runs" hcl:"max_concurrent_runs,optional"`
	ResourceUsage     *ResourceUsage     `json:"resource_usage" hcl:"resource_usage,optional"`
	ErrorBreakdown    map[string]int     `json:"error_breakdown" hcl:"error_breakdown,optional"`
	PerformanceTrends []PerformancePoint `json:"performance_trends" hcl:"performance_trends,optional"`
}

// PerformancePoint represents a performance measurement point
type PerformancePoint struct {
	spookytypescommon.CompleteEntity
	Timestamp      time.Time      `json:"timestamp" hcl:"timestamp"`
	Duration       time.Duration  `json:"duration" hcl:"duration"`
	Success        bool           `json:"success" hcl:"success"`
	ResourceUsage  *ResourceUsage `json:"resource_usage" hcl:"resource_usage,optional"`
	MachineCount   int            `json:"machine_count" hcl:"machine_count"`
	ConcurrentRuns int            `json:"concurrent_runs" hcl:"concurrent_runs"`
}

// ActionValidation represents validation results for actions
type ActionValidation struct {
	spookytypescommon.CompleteEntity
	ActionName     string              `json:"action_name" hcl:"action_name"`
	Valid          bool                `json:"valid" hcl:"valid"`
	Errors         []ValidationError   `json:"errors" hcl:"errors,optional"`
	Warnings       []ValidationWarning `json:"warnings" hcl:"warnings,optional"`
	ValidatedAt    time.Time           `json:"validated_at" hcl:"validated_at"`
	ValidatedBy    string              `json:"validated_by" hcl:"validated_by,optional"`
	SchemaVersion  string              `json:"schema_version" hcl:"schema_version,optional"`
	ValidationTime time.Duration       `json:"validation_time" hcl:"validation_time,optional"`
}

// ValidationWarning represents a validation warning
type ValidationWarning struct {
	spookytypescommon.CompleteEntity
	Field       string   `json:"field" hcl:"field"`
	Message     string   `json:"message" hcl:"message"`
	Code        string   `json:"code" hcl:"code"`
	Line        int      `json:"line" hcl:"line,optional"`
	Column      int      `json:"column" hcl:"column,optional"`
	Context     string   `json:"context" hcl:"context,optional"`
	Suggestions []string `json:"suggestions" hcl:"suggestions,optional"`
}
