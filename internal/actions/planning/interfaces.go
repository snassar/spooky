package planning

import (
	"time"

	spookyactionstypes "spooky/internal/actions/types"
)

// PlanningManager defines the interface for action planning operations
type PlanningManager interface {
	// Core planning operations
	PlanAction(action *spookyactionstypes.Action, context *spookyactionstypes.ActionContext) (*spookyactionstypes.ActionPlan, error)
	PlanActionCollection(collection *spookyactionstypes.ActionCollection, context *spookyactionstypes.ActionContext) (*spookyactionstypes.ActionPlan, error)
	ValidatePlan(plan *spookyactionstypes.ActionPlan) error

	// Planner management
	CreatePlanner(action *spookyactionstypes.Action) (Planner, error)
	GetPlanner(action *spookyactionstypes.Action) (Planner, error)

	// Plan management
	GetPlan(planID string) (*spookyactionstypes.ActionPlan, error)
	ListPlans() ([]*spookyactionstypes.ActionPlan, error)
	DeletePlan(planID string) error

	// Configuration
	SetDefaultStrategy(strategy spookyactionstypes.PlanningStrategy)
	SetDefaultOptimization(optimization spookyactionstypes.PlanningOptimization)
}

// Planner defines the interface for action planning
type Planner interface {
	// Core planning operations
	Plan(context *spookyactionstypes.ActionContext) (*spookyactionstypes.ActionPlan, error)
	Validate(plan *spookyactionstypes.ActionPlan) error
	Optimize(plan *spookyactionstypes.ActionPlan) error

	// Strategy management
	SetStrategy(strategy spookyactionstypes.PlanningStrategy)
	GetStrategy() spookyactionstypes.PlanningStrategy

	// Configuration
	SetOptimization(optimization spookyactionstypes.PlanningOptimization)
	SetConstraints(constraints []spookyactionstypes.PlanningConstraint)
}

// PlanGenerator defines the interface for generating action plans
type PlanGenerator interface {
	// Plan generation
	GeneratePlan(action *spookyactionstypes.Action, context *spookyactionstypes.ActionContext) (*spookyactionstypes.ActionPlan, error)
	GenerateCollectionPlan(collection *spookyactionstypes.ActionCollection, context *spookyactionstypes.ActionContext) (*spookyactionstypes.ActionPlan, error)

	// Plan customization
	CustomizePlan(plan *spookyactionstypes.ActionPlan, options PlanOptions) error
	ApplyConstraints(plan *spookyactionstypes.ActionPlan, constraints []spookyactionstypes.PlanningConstraint) error

	// Plan validation
	ValidatePlanStructure(plan *spookyactionstypes.ActionPlan) error
	ValidatePlanDependencies(plan *spookyactionstypes.ActionPlan) error
	ValidatePlanResources(plan *spookyactionstypes.ActionPlan) error
}

// PlanOptimizer defines the interface for optimizing action plans
type PlanOptimizer interface {
	// Plan optimization
	OptimizePlan(plan *spookyactionstypes.ActionPlan, target OptimizationTarget) error
	OptimizeForSpeed(plan *spookyactionstypes.ActionPlan) error
	OptimizeForResources(plan *spookyactionstypes.ActionPlan) error
	OptimizeForReliability(plan *spookyactionstypes.ActionPlan) error

	// Optimization analysis
	AnalyzePlan(plan *spookyactionstypes.ActionPlan) (*PlanAnalysis, error)
	GetOptimizationSuggestions(plan *spookyactionstypes.ActionPlan) ([]OptimizationSuggestion, error)

	// Configuration
	SetOptimizationLevel(level spookyactionstypes.OptimizationLevel)
	SetResourceLimits(limits *spookyactionstypes.ResourceLimits)
}

// PlanValidator defines the interface for validating action plans
type PlanValidator interface {
	// Plan validation
	ValidatePlan(plan *spookyactionstypes.ActionPlan) error
	ValidatePlanStep(step *spookyactionstypes.PlanStep) error
	ValidatePlanDependencies(plan *spookyactionstypes.ActionPlan) error

	// Validation rules
	AddValidationRule(rule PlanValidationRule) error
	RemoveValidationRule(name string) error
	GetValidationRules() ([]PlanValidationRule, error)

	// Configuration
	SetStrictMode(strict bool)
	SetValidationLevel(level spookyactionstypes.ValidationLevel)
}

// PlanScheduler defines the interface for scheduling action plans
type PlanScheduler interface {
	// Plan scheduling
	SchedulePlan(plan *spookyactionstypes.ActionPlan, schedule ScheduleOptions) error
	ReschedulePlan(planID string, schedule ScheduleOptions) error
	CancelScheduledPlan(planID string) error

	// Schedule management
	GetScheduledPlans() ([]*ScheduledPlan, error)
	GetScheduledPlan(planID string) (*ScheduledPlan, error)

	// Schedule validation
	ValidateSchedule(schedule ScheduleOptions) error
	CheckScheduleConflicts(schedule ScheduleOptions) ([]ScheduleConflict, error)
}

// PlanOptions represents options for plan customization
type PlanOptions struct {
	Strategy      spookyactionstypes.PlanningStrategy     `json:"strategy"`
	Optimization  spookyactionstypes.PlanningOptimization `json:"optimization"`
	Constraints   []spookyactionstypes.PlanningConstraint `json:"constraints"`
	Timeout       time.Duration                           `json:"timeout"`
	Parallel      bool                                    `json:"parallel"`
	MaxConcurrent int                                     `json:"max_concurrent"`
}

// OptimizationTarget represents the target for plan optimization
type OptimizationTarget string

const (
	OptimizationTargetSpeed       OptimizationTarget = "speed"
	OptimizationTargetResources   OptimizationTarget = "resources"
	OptimizationTargetReliability OptimizationTarget = "reliability"
	OptimizationTargetCost        OptimizationTarget = "cost"
	OptimizationTargetBalanced    OptimizationTarget = "balanced"
)

// PlanAnalysis represents analysis results for a plan
type PlanAnalysis struct {
	PlanID            string                 `json:"plan_id"`
	TotalSteps        int                    `json:"total_steps"`
	EstimatedDuration time.Duration          `json:"estimated_duration"`
	ResourceUsage     *ResourceUsage         `json:"resource_usage"`
	Dependencies      []DependencyInfo       `json:"dependencies"`
	Bottlenecks       []BottleneckInfo       `json:"bottlenecks"`
	OptimizationScore float64                `json:"optimization_score"`
	Metadata          map[string]interface{} `json:"metadata"`
}

// ResourceUsage represents resource usage information
type ResourceUsage struct {
	MemoryMB    int     `json:"memory_mb"`
	CPUPercent  float64 `json:"cpu_percent"`
	DiskMB      int     `json:"disk_mb"`
	NetworkMBPS int     `json:"network_mbps"`
}

// DependencyInfo represents dependency information
type DependencyInfo struct {
	StepID       string        `json:"step_id"`
	Dependencies []string      `json:"dependencies"`
	Dependents   []string      `json:"dependents"`
	CriticalPath bool          `json:"critical_path"`
	WaitTime     time.Duration `json:"wait_time"`
}

// BottleneckInfo represents bottleneck information
type BottleneckInfo struct {
	StepID      string        `json:"step_id"`
	Type        string        `json:"type"`
	Description string        `json:"description"`
	Impact      float64       `json:"impact"`
	Duration    time.Duration `json:"duration"`
}

// OptimizationSuggestion represents an optimization suggestion
type OptimizationSuggestion struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Impact      float64  `json:"impact"`
	Difficulty  string   `json:"difficulty"`
	Steps       []string `json:"steps"`
}

// ScheduleOptions represents options for plan scheduling
type ScheduleOptions struct {
	StartTime    time.Time     `json:"start_time"`
	EndTime      time.Time     `json:"end_time"`
	Priority     int           `json:"priority"`
	RetryCount   int           `json:"retry_count"`
	RetryDelay   time.Duration `json:"retry_delay"`
	MaxDuration  time.Duration `json:"max_duration"`
	AllowOverlap bool          `json:"allow_overlap"`
}

// ScheduledPlan represents a scheduled plan
type ScheduledPlan struct {
	PlanID    string                         `json:"plan_id"`
	Plan      *spookyactionstypes.ActionPlan `json:"plan"`
	Schedule  ScheduleOptions                `json:"schedule"`
	Status    string                         `json:"status"`
	CreatedAt time.Time                      `json:"created_at"`
	UpdatedAt time.Time                      `json:"updated_at"`
}

// ScheduleConflict represents a schedule conflict
type ScheduleConflict struct {
	PlanID1      string `json:"plan_id_1"`
	PlanID2      string `json:"plan_id_2"`
	ConflictType string `json:"conflict_type"`
	Description  string `json:"description"`
	Severity     string `json:"severity"`
	Resolution   string `json:"resolution"`
}

// PlanValidationRule defines a custom plan validation rule
type PlanValidationRule interface {
	Name() string
	Validate(plan *spookyactionstypes.ActionPlan) error
	ValidateStep(step *spookyactionstypes.PlanStep) error
	GetDescription() string
}
