package actions

import (
	"time"
)

// ActionPlan represents a plan for action execution
type ActionPlan struct {
	// Plan identification
	PlanID     string `hcl:"plan_id,optional"`
	ActionName string `hcl:"action_name,optional"`

	// Plan state
	Status    PlanningStatus `hcl:"status,optional"`
	CreatedAt time.Time      `hcl:"created_at,optional"`
	UpdatedAt time.Time      `hcl:"updated_at,optional"`

	// Plan structure
	Steps        []*PlanStep          `hcl:"steps,block"`
	Strategy     PlanningStrategy     `hcl:"strategy,optional"`
	Optimization PlanningOptimization `hcl:"optimization,optional"`

	// Plan metadata
	Metadata    map[string]interface{} `hcl:"metadata,optional"`
	Constraints []PlanningConstraint   `hcl:"constraints,block"`
}

// PlanStep represents a step in an action plan
type PlanStep struct {
	// Step identification
	StepID    string `hcl:"step_id,optional"`
	StepName  string `hcl:"step_name,optional"`
	StepOrder int    `hcl:"step_order,optional"`

	// Step configuration
	Action   *Action       `hcl:"action,block"`
	Machines []string      `hcl:"machines,optional"`
	Parallel bool          `hcl:"parallel,optional"`
	Timeout  time.Duration `hcl:"timeout,optional"`

	// Step dependencies
	Dependencies  []string `hcl:"dependencies,optional"`
	Prerequisites []string `hcl:"prerequisites,optional"`

	// Step state
	Status    PlanStepStatus `hcl:"status,optional"`
	StartTime *time.Time     `hcl:"start_time,optional"`
	EndTime   *time.Time     `hcl:"end_time,optional"`

	// Step metadata
	Metadata map[string]interface{} `hcl:"metadata,optional"`
}

// PlanningStatus represents the status of action planning
type PlanningStatus string

const (
	PlanningStatusPending   PlanningStatus = "pending"
	PlanningStatusRunning   PlanningStatus = "running"
	PlanningStatusCompleted PlanningStatus = "completed"
	PlanningStatusFailed    PlanningStatus = "failed"
	PlanningStatusCancelled PlanningStatus = "cancelled"
)

// PlanStepStatus represents the status of a plan step
type PlanStepStatus string

const (
	PlanStepStatusPending   PlanStepStatus = "pending"
	PlanStepStatusRunning   PlanStepStatus = "running"
	PlanStepStatusCompleted PlanStepStatus = "completed"
	PlanStepStatusFailed    PlanStepStatus = "failed"
	PlanStepStatusSkipped   PlanStepStatus = "skipped"
)

// PlanningStrategy represents the strategy for action planning
type PlanningStrategy string

const (
	PlanningStrategySequential PlanningStrategy = "sequential"
	PlanningStrategyParallel   PlanningStrategy = "parallel"
	PlanningStrategyHybrid     PlanningStrategy = "hybrid"
	PlanningStrategyOptimized  PlanningStrategy = "optimized"
)

// PlanningOptimization represents the optimization level for action planning
type PlanningOptimization string

const (
	PlanningOptimizationNone     PlanningOptimization = "none"
	PlanningOptimizationBasic    PlanningOptimization = "basic"
	PlanningOptimizationAdvanced PlanningOptimization = "advanced"
	PlanningOptimizationMaximum  PlanningOptimization = "maximum"
)

// PlanningConstraint represents a constraint for action planning
type PlanningConstraint struct {
	// Constraint identification
	ConstraintID   string `hcl:"constraint_id,optional"`
	ConstraintType string `hcl:"constraint_type,optional"`
	ConstraintName string `hcl:"constraint_name,optional"`

	// Constraint configuration
	Value    interface{} `hcl:"value,optional"`
	Operator string      `hcl:"operator,optional"`
	Priority int         `hcl:"priority,optional"`

	// Constraint metadata
	Description string                 `hcl:"description,optional"`
	Metadata    map[string]interface{} `hcl:"metadata,optional"`
}

// MergePolicy represents a policy for merging actions
type MergePolicy struct {
	// Policy identification
	PolicyID   string `hcl:"policy_id,optional"`
	PolicyName string `hcl:"policy_name,optional"`

	// Policy configuration
	Strategy           MergeStrategy      `hcl:"strategy,optional"`
	ConflictResolution ConflictResolution `hcl:"conflict_resolution,optional"`
	Priority           int                `hcl:"priority,optional"`

	// Policy rules
	Rules []MergeRule `hcl:"rules,block"`

	// Policy metadata
	Description string                 `hcl:"description,optional"`
	Metadata    map[string]interface{} `hcl:"metadata,optional"`
}

// MergeStrategy represents the strategy for merging actions
type MergeStrategy string

const (
	MergeStrategyAppend  MergeStrategy = "append"
	MergeStrategyReplace MergeStrategy = "replace"
	MergeStrategyMerge   MergeStrategy = "merge"
	MergeStrategySelect  MergeStrategy = "select"
	MergeStrategyCustom  MergeStrategy = "custom"
)

// ConflictResolution represents how to resolve conflicts during merging
type ConflictResolution string

const (
	ConflictResolutionFirst ConflictResolution = "first"
	ConflictResolutionLast  ConflictResolution = "last"
	ConflictResolutionMerge ConflictResolution = "merge"
	ConflictResolutionSkip  ConflictResolution = "skip"
	ConflictResolutionError ConflictResolution = "error"
)

// MergeRule represents a rule for merging actions
type MergeRule struct {
	// Rule identification
	RuleID   string `hcl:"rule_id,optional"`
	RuleName string `hcl:"rule_name,optional"`

	// Rule configuration
	Field    string      `hcl:"field,optional"`
	Operator string      `hcl:"operator,optional"`
	Value    interface{} `hcl:"value,optional"`
	Action   string      `hcl:"action,optional"`

	// Rule metadata
	Description string                 `hcl:"description,optional"`
	Metadata    map[string]interface{} `hcl:"metadata,optional"`
}

// PerformanceMetrics represents performance metrics for actions
type PerformanceMetrics struct {
	// Metrics identification
	ActionName string `hcl:"action_name,optional"`
	MachineID  string `hcl:"machine_id,optional"`

	// Timing metrics
	StartTime   *time.Time    `hcl:"start_time,optional"`
	EndTime     *time.Time    `hcl:"end_time,optional"`
	Duration    time.Duration `hcl:"duration,optional"`
	AverageTime time.Duration `hcl:"average_time,optional"`
	MinTime     time.Duration `hcl:"min_time,optional"`
	MaxTime     time.Duration `hcl:"max_time,optional"`

	// Resource metrics
	MemoryUsage  int64   `hcl:"memory_usage,optional"`
	CPUUsage     float64 `hcl:"cpu_usage,optional"`
	DiskUsage    int64   `hcl:"disk_usage,optional"`
	NetworkUsage int64   `hcl:"network_usage,optional"`

	// Execution metrics
	ExecutionCount int     `hcl:"execution_count,optional"`
	SuccessCount   int     `hcl:"success_count,optional"`
	FailureCount   int     `hcl:"failure_count,optional"`
	SuccessRate    float64 `hcl:"success_rate,optional"`

	// Performance metadata
	CreatedAt time.Time              `hcl:"created_at,optional"`
	UpdatedAt time.Time              `hcl:"updated_at,optional"`
	Metadata  map[string]interface{} `hcl:"metadata,optional"`
}

// ResourceLimits represents resource limits for actions
type ResourceLimits struct {
	// Memory limits
	MemoryMB      int `hcl:"memory_mb,optional"`
	MemoryPercent int `hcl:"memory_percent,optional"`

	// CPU limits
	CPUPercent int `hcl:"cpu_percent,optional"`
	CPUCores   int `hcl:"cpu_cores,optional"`

	// Disk limits
	DiskMB      int `hcl:"disk_mb,optional"`
	DiskPercent int `hcl:"disk_percent,optional"`

	// Network limits
	NetworkMBPS int `hcl:"network_mbps,optional"`

	// Process limits
	MaxProcesses int `hcl:"max_processes,optional"`
	MaxFiles     int `hcl:"max_files,optional"`
}

// OptimizationLevel represents the level of performance optimization
type OptimizationLevel string

const (
	OptimizationLevelNone     OptimizationLevel = "none"
	OptimizationLevelBasic    OptimizationLevel = "basic"
	OptimizationLevelAdvanced OptimizationLevel = "advanced"
	OptimizationLevelMaximum  OptimizationLevel = "maximum"
)

// OptimizationTarget represents the target for performance optimization
type OptimizationTarget string

const (
	OptimizationTargetSpeed    OptimizationTarget = "speed"
	OptimizationTargetMemory   OptimizationTarget = "memory"
	OptimizationTargetCPU      OptimizationTarget = "cpu"
	OptimizationTargetNetwork  OptimizationTarget = "network"
	OptimizationTargetBalanced OptimizationTarget = "balanced"
)

// ExportFormat represents the format for exporting actions
type ExportFormat string

const (
	ExportFormatJSON ExportFormat = "json"
	ExportFormatHCL  ExportFormat = "hcl"
)

// ValidationLevel represents the level of validation
type ValidationLevel string

const (
	ValidationLevelNone     ValidationLevel = "none"
	ValidationLevelBasic    ValidationLevel = "basic"
	ValidationLevelStrict   ValidationLevel = "strict"
	ValidationLevelComplete ValidationLevel = "complete"
)

// NewActionPlan creates a new ActionPlan
func NewActionPlan(actionName string) *ActionPlan {
	now := time.Now()
	return &ActionPlan{
		ActionName:   actionName,
		Status:       PlanningStatusPending,
		CreatedAt:    now,
		UpdatedAt:    now,
		Steps:        make([]*PlanStep, 0),
		Strategy:     PlanningStrategySequential,
		Optimization: PlanningOptimizationNone,
		Metadata:     make(map[string]interface{}),
		Constraints:  make([]PlanningConstraint, 0),
	}
}

// NewPlanStep creates a new PlanStep
func NewPlanStep(stepName string, stepOrder int) *PlanStep {
	return &PlanStep{
		StepID:        stepName,
		StepName:      stepName,
		StepOrder:     stepOrder,
		Status:        PlanStepStatusPending,
		Machines:      make([]string, 0),
		Dependencies:  make([]string, 0),
		Prerequisites: make([]string, 0),
		Metadata:      make(map[string]interface{}),
	}
}

// NewMergePolicy creates a new MergePolicy
func NewMergePolicy(policyName string) *MergePolicy {
	return &MergePolicy{
		PolicyID:           policyName,
		PolicyName:         policyName,
		Strategy:           MergeStrategyAppend,
		ConflictResolution: ConflictResolutionFirst,
		Priority:           0,
		Rules:              make([]MergeRule, 0),
		Metadata:           make(map[string]interface{}),
	}
}

// NewPerformanceMetrics creates a new PerformanceMetrics
func NewPerformanceMetrics(actionName string) *PerformanceMetrics {
	now := time.Now()
	return &PerformanceMetrics{
		ActionName: actionName,
		CreatedAt:  now,
		UpdatedAt:  now,
		Metadata:   make(map[string]interface{}),
	}
}
