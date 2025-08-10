package interfaces

import (
	"context"
	"time"

	spookytypesactions "spooky/internal/types/actions"
)

// ActionManager defines the main interface for action management and execution
type ActionManager interface {
	// Core action operations
	LoadActions(projectPath string) (*spookytypesactions.ActionCollection, error)
	GetAction(name string) (*spookytypesactions.Action, error)
	ListActions() ([]*spookytypesactions.Action, error)
	AddAction(name string, action *spookytypesactions.Action) error
	RemoveAction(name string) error

	// Acting operations
	ExecuteAction(ctx context.Context, action *spookytypesactions.Action, context *spookytypesactions.ActionContext) (*spookytypesactions.ActingSession, error)
	ExecuteActionCollection(ctx context.Context, collection *spookytypesactions.ActionCollection, context *spookytypesactions.ActionContext) (*spookytypesactions.ActingSession, error)
	PrepareAction(action *spookytypesactions.Action, context *spookytypesactions.ActionContext) error

	// Planning operations
	PlanAction(action *spookytypesactions.Action, context *spookytypesactions.ActionContext) (*spookytypesactions.ActionPlan, error)
	PlanActionCollection(collection *spookytypesactions.ActionCollection, context *spookytypesactions.ActionContext) (*spookytypesactions.ActionPlan, error)
	ValidatePlan(plan *spookytypesactions.ActionPlan) error

	// Validation operations
	ValidateAction(action *spookytypesactions.Action) error
	ValidateActionCollection(collection *spookytypesactions.ActionCollection) error
	ValidateActionContext(context *spookytypesactions.ActionContext) error

	// Merging operations
	MergeActions(actions ...*spookytypesactions.Action) (*spookytypesactions.ActionCollection, error)
	MergeWithPolicy(existing, new *spookytypesactions.ActionCollection, policy spookytypesactions.MergePolicy) (*spookytypesactions.ActionCollection, error)

	// Performance operations
	OptimizeAction(action *spookytypesactions.Action) error
	OptimizeActionCollection(collection *spookytypesactions.ActionCollection) error
	GetPerformanceMetrics(action *spookytypesactions.Action) (*spookytypesactions.PerformanceMetrics, error)

	// Configuration
	SetDefaultTimeout(timeout time.Duration)
	SetDefaultParallel(parallel bool)
	RegisterCustomValidator(name string, validator ActionValidator)

	// Utility operations
	Close() error
}

// ActionValidator defines the interface for action validation
type ActionValidator interface {
	// Core validation operations
	Validate(action *spookytypesactions.Action) error
	ValidateCollection(collection *spookytypesactions.ActionCollection) error
	ValidateContext(context *spookytypesactions.ActionContext) error
}

// ActingManager defines the interface for action execution operations
type ActingManager interface {
	// Core acting operations
	ExecuteAction(ctx context.Context, action *spookytypesactions.Action, context *spookytypesactions.ActionContext) (*spookytypesactions.ActingSession, error)
	ExecuteActionCollection(ctx context.Context, collection *spookytypesactions.ActionCollection, context *spookytypesactions.ActionContext) (*spookytypesactions.ActingSession, error)
	PrepareAction(action *spookytypesactions.Action, context *spookytypesactions.ActionContext) error

	// Actor management
	CreateActor(action *spookytypesactions.Action, context *spookytypesactions.ActionContext) (Actor, error)
	GetActor(action *spookytypesactions.Action) (Actor, error)

	// Session management
	GetSession(sessionID string) (*spookytypesactions.ActingSession, error)
	ListSessions() ([]*spookytypesactions.ActingSession, error)
	CancelSession(sessionID string) error

	// Configuration
	SetDefaultTimeout(timeout time.Duration)
	SetDefaultParallel(parallel bool)
	SetMaxConcurrent(maxConcurrent int)
}

// Actor defines the interface for individual action execution
type Actor interface {
	// Core acting operations
	Execute(ctx context.Context, context *spookytypesactions.ActionContext) (*spookytypesactions.ActingResult, error)
	Prepare(context *spookytypesactions.ActionContext) error
	Cancel() error

	// State management
	GetState() spookytypesactions.ActingState
	GetProgress() float64
	GetStatus() spookytypesactions.ActingStatus

	// Configuration
	SetTimeout(timeout time.Duration)
	SetParallel(parallel bool)
}

// ActingExecutor defines the interface for executing actions on machines
type ActingExecutor interface {
	// Core execution operations
	ExecuteCommand(ctx context.Context, command string, context *spookytypesactions.ActionContext) (*spookytypesactions.ActingResult, error)
	ExecuteScript(ctx context.Context, script string, context *spookytypesactions.ActionContext) (*spookytypesactions.ActingResult, error)
	ExecuteTemplate(ctx context.Context, template *spookytypesactions.TemplateConfig, context *spookytypesactions.ActionContext) (*spookytypesactions.ActingResult, error)

	// Machine management
	GetMachine(machineID string) (Machine, error)
	ListMachines() ([]Machine, error)
	ConnectMachine(machineID string) error
	DisconnectMachine(machineID string) error

	// Configuration
	SetConnectionTimeout(timeout time.Duration)
	SetCommandTimeout(timeout time.Duration)
	SetRetryAttempts(attempts int)
	SetRetryDelay(delay time.Duration)
}

// Machine defines the interface for machine operations
type Machine interface {
	// Machine information
	GetID() string
	GetName() string
	GetHost() string
	GetUser() string
	GetPort() int

	// Connection management
	Connect() error
	Disconnect() error
	IsConnected() bool

	// Execution operations
	ExecuteCommand(ctx context.Context, command string) (*spookytypesactions.ActingResult, error)
	ExecuteScript(ctx context.Context, script string) (*spookytypesactions.ActingResult, error)
	UploadFile(ctx context.Context, localPath, remotePath string) error
	DownloadFile(ctx context.Context, remotePath, localPath string) error

	// Configuration
	SetTimeout(timeout time.Duration)
	SetSudo(sudo bool)
	SetWorkingDirectory(dir string)
	SetEnvironment(env map[string]string)
}

// ActingSessionManager defines the interface for managing acting sessions
type ActingSessionManager interface {
	// Session management
	CreateSession(actionName string) (*spookytypesactions.ActingSession, error)
	GetSession(sessionID string) (*spookytypesactions.ActingSession, error)
	ListSessions() ([]*spookytypesactions.ActingSession, error)
	UpdateSession(session *spookytypesactions.ActingSession) error
	DeleteSession(sessionID string) error

	// Session state management
	StartSession(sessionID string) error
	CompleteSession(sessionID string) error
	FailSession(sessionID string, err error) error
	CancelSession(sessionID string) error

	// Session cleanup
	CleanupExpiredSessions(maxAge time.Duration) error
	CleanupAllSessions() error
}

// ActingResultProcessor defines the interface for processing acting results
type ActingResultProcessor interface {
	// Result processing
	ProcessResult(result *spookytypesactions.ActingResult) error
	ProcessResults(results []*spookytypesactions.ActingResult) error

	// Result aggregation
	AggregateResults(results []*spookytypesactions.ActingResult) (*spookytypesactions.ActingSession, error)
	CalculateSuccessRate(results []*spookytypesactions.ActingResult) float64

	// Result validation
	ValidateResult(result *spookytypesactions.ActingResult) error
	ValidateResults(results []*spookytypesactions.ActingResult) error

	// Result transformation
	TransformResult(result *spookytypesactions.ActingResult, format string) (interface{}, error)
	TransformResults(results []*spookytypesactions.ActingResult, format string) (interface{}, error)
}

// ActingProgressTracker defines the interface for tracking acting progress
type ActingProgressTracker interface {
	// Progress tracking
	StartTracking(sessionID string, totalSteps int) error
	UpdateProgress(sessionID string, completedSteps int) error
	GetProgress(sessionID string) (float64, error)
	CompleteTracking(sessionID string) error

	// Progress reporting
	GetProgressReport(sessionID string) (*ProgressReport, error)
	ListProgressReports() ([]*ProgressReport, error)

	// Progress cleanup
	CleanupProgress(sessionID string) error
	CleanupAllProgress() error
}

// ProgressReport represents a progress report for acting operations
type ProgressReport struct {
	SessionID      string    `json:"session_id"`
	ActionName     string    `json:"action_name"`
	TotalSteps     int       `json:"total_steps"`
	CompletedSteps int       `json:"completed_steps"`
	Progress       float64   `json:"progress"`
	StartTime      time.Time `json:"start_time"`
	LastUpdate     time.Time `json:"last_update"`
	EstimatedEnd   time.Time `json:"estimated_end,omitempty"`
	Status         string    `json:"status"`
}

// MergingManager defines the interface for action merging operations
type MergingManager interface {
	// Core merging operations
	MergeActions(actions ...*spookytypesactions.Action) (*spookytypesactions.ActionCollection, error)
	MergeWithPolicy(existing, new *spookytypesactions.ActionCollection, policy spookytypesactions.MergePolicy) (*spookytypesactions.ActionCollection, error)

	// Merger management
	CreateMerger(policy spookytypesactions.MergePolicy) (ActionMerger, error)
	GetMerger(policy spookytypesactions.MergePolicy) (ActionMerger, error)

	// Policy management
	GetMergePolicy(name string) (spookytypesactions.MergePolicy, error)
	ListMergePolicies() ([]spookytypesactions.MergePolicy, error)
	AddMergePolicy(policy spookytypesactions.MergePolicy) error
	RemoveMergePolicy(name string) error

	// Configuration
	SetDefaultPolicy(policy spookytypesactions.MergePolicy)
	SetConflictResolution(resolution spookytypesactions.ConflictResolution)
}

// ActionMerger defines the interface for action merging
type ActionMerger interface {
	// Core merging operations
	Merge(actions ...*spookytypesactions.Action) (*spookytypesactions.ActionCollection, error)
	MergeWithPolicy(existing, new *spookytypesactions.ActionCollection, policy spookytypesactions.MergePolicy) (*spookytypesactions.ActionCollection, error)

	// Policy management
	SetPolicy(policy spookytypesactions.MergePolicy)
	GetPolicy() spookytypesactions.MergePolicy

	// Configuration
	SetConflictResolution(resolution spookytypesactions.ConflictResolution)
	SetMergeStrategy(strategy spookytypesactions.MergeStrategy)
}

// MergeResult represents the result of a merge operation
type MergeResult struct {
	Success       bool                         `json:"success"`
	MergedActions []*spookytypesactions.Action `json:"merged_actions"`
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

// PerformanceManager defines the interface for action performance operations
type PerformanceManager interface {
	// Core performance operations
	OptimizeAction(action *spookytypesactions.Action) error
	OptimizeActionCollection(collection *spookytypesactions.ActionCollection) error
	GetPerformanceMetrics(action *spookytypesactions.Action) (*spookytypesactions.PerformanceMetrics, error)

	// Optimizer management
	CreateOptimizer(action *spookytypesactions.Action) (ActionOptimizer, error)
	GetOptimizer(action *spookytypesactions.Action) (ActionOptimizer, error)

	// Metrics management
	GetMetrics(action *spookytypesactions.Action) (*spookytypesactions.PerformanceMetrics, error)
	ListMetrics() ([]*spookytypesactions.PerformanceMetrics, error)
	ClearMetrics(action *spookytypesactions.Action) error

	// Configuration
	SetOptimizationLevel(level spookytypesactions.OptimizationLevel)
	SetResourceLimits(limits *spookytypesactions.ResourceLimits)
}

// ActionOptimizer defines the interface for action performance optimization
type ActionOptimizer interface {
	// Core optimization operations
	Optimize(action *spookytypesactions.Action) error
	GetMetrics(action *spookytypesactions.Action) (*spookytypesactions.PerformanceMetrics, error)

	// Optimization management
	SetLevel(level spookytypesactions.OptimizationLevel)
	GetLevel() spookytypesactions.OptimizationLevel

	// Configuration
	SetResourceLimits(limits *spookytypesactions.ResourceLimits)
	SetOptimizationTarget(target spookytypesactions.OptimizationTarget)
}

// PerformanceAnalyzer defines the interface for performance analysis
type PerformanceAnalyzer interface {
	// Analysis operations
	AnalyzeAction(action *spookytypesactions.Action) (*PerformanceAnalysis, error)
	AnalyzeCollection(collection *spookytypesactions.ActionCollection) (*PerformanceAnalysis, error)
	CompareActions(action1, action2 *spookytypesactions.Action) (*PerformanceComparison, error)

	// Analysis configuration
	SetAnalysisDepth(depth int)
	SetIncludeMetrics(include bool)
}

// PerformanceAnalysis represents performance analysis results
type PerformanceAnalysis struct {
	ActionName      string                 `json:"action_name"`
	ExecutionTime   float64                `json:"execution_time"`
	MemoryUsage     int64                  `json:"memory_usage"`
	CPUUsage        float64                `json:"cpu_usage"`
	NetworkUsage    int64                  `json:"network_usage"`
	Bottlenecks     []Bottleneck           `json:"bottlenecks"`
	Optimizations   []Optimization         `json:"optimizations"`
	Recommendations []Recommendation       `json:"recommendations"`
	Details         map[string]interface{} `json:"details"`
}

// PerformanceComparison represents performance comparison results
type PerformanceComparison struct {
	Action1Name       string                 `json:"action1_name"`
	Action2Name       string                 `json:"action2_name"`
	ExecutionTimeDiff float64                `json:"execution_time_diff"`
	MemoryUsageDiff   int64                  `json:"memory_usage_diff"`
	CPUUsageDiff      float64                `json:"cpu_usage_diff"`
	NetworkUsageDiff  int64                  `json:"network_usage_diff"`
	Winner            string                 `json:"winner"`
	Details           map[string]interface{} `json:"details"`
}

// Bottleneck represents a performance bottleneck
type Bottleneck struct {
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Impact      float64 `json:"impact"`
	Location    string  `json:"location"`
	Severity    string  `json:"severity"`
}

// Optimization represents a performance optimization
type Optimization struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Impact      float64  `json:"impact"`
	Difficulty  string   `json:"difficulty"`
	Steps       []string `json:"steps"`
}

// Recommendation represents a performance recommendation
type Recommendation struct {
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Priority    string  `json:"priority"`
	Effort      string  `json:"effort"`
	Benefit     float64 `json:"benefit"`
}

// PlanningManager defines the interface for action planning operations
type PlanningManager interface {
	// Core planning operations
	PlanAction(action *spookytypesactions.Action, context *spookytypesactions.ActionContext) (*spookytypesactions.ActionPlan, error)
	PlanActionCollection(collection *spookytypesactions.ActionCollection, context *spookytypesactions.ActionContext) (*spookytypesactions.ActionPlan, error)
	ValidatePlan(plan *spookytypesactions.ActionPlan) error

	// Planner management
	CreatePlanner(action *spookytypesactions.Action) (Planner, error)
	GetPlanner(action *spookytypesactions.Action) (Planner, error)

	// Plan management
	GetPlan(planID string) (*spookytypesactions.ActionPlan, error)
	ListPlans() ([]*spookytypesactions.ActionPlan, error)
	DeletePlan(planID string) error

	// Configuration
	SetDefaultStrategy(strategy spookytypesactions.PlanningStrategy)
	SetDefaultOptimization(optimization spookytypesactions.PlanningOptimization)
}

// Planner defines the interface for action planning
type Planner interface {
	// Core planning operations
	Plan(context *spookytypesactions.ActionContext) (*spookytypesactions.ActionPlan, error)
	Validate(plan *spookytypesactions.ActionPlan) error
	Optimize(plan *spookytypesactions.ActionPlan) error

	// Strategy management
	SetStrategy(strategy spookytypesactions.PlanningStrategy)
	GetStrategy() spookytypesactions.PlanningStrategy

	// Configuration
	SetOptimization(optimization spookytypesactions.PlanningOptimization)
	SetConstraints(constraints []spookytypesactions.PlanningConstraint)
}

// PlanGenerator defines the interface for generating action plans
type PlanGenerator interface {
	// Plan generation
	GeneratePlan(action *spookytypesactions.Action, context *spookytypesactions.ActionContext) (*spookytypesactions.ActionPlan, error)
	GenerateCollectionPlan(collection *spookytypesactions.ActionCollection, context *spookytypesactions.ActionContext) (*spookytypesactions.ActionPlan, error)

	// Plan customization
	CustomizePlan(plan *spookytypesactions.ActionPlan, options PlanOptions) error
	ApplyConstraints(plan *spookytypesactions.ActionPlan, constraints []spookytypesactions.PlanningConstraint) error

	// Plan validation
	ValidatePlanStructure(plan *spookytypesactions.ActionPlan) error
	ValidatePlanDependencies(plan *spookytypesactions.ActionPlan) error
	ValidatePlanResources(plan *spookytypesactions.ActionPlan) error
}

// PlanOptimizer defines the interface for optimizing action plans
type PlanOptimizer interface {
	// Plan optimization
	OptimizePlan(plan *spookytypesactions.ActionPlan, target OptimizationTarget) error
	OptimizeForSpeed(plan *spookytypesactions.ActionPlan) error
	OptimizeForResources(plan *spookytypesactions.ActionPlan) error
	OptimizeForReliability(plan *spookytypesactions.ActionPlan) error

	// Optimization analysis
	AnalyzePlan(plan *spookytypesactions.ActionPlan) (*PlanAnalysis, error)
	GetOptimizationSuggestions(plan *spookytypesactions.ActionPlan) ([]OptimizationSuggestion, error)

	// Configuration
	SetOptimizationLevel(level spookytypesactions.OptimizationLevel)
	SetResourceLimits(limits *spookytypesactions.ResourceLimits)
}

// PlanValidator defines the interface for validating action plans
type PlanValidator interface {
	// Plan validation
	ValidatePlan(plan *spookytypesactions.ActionPlan) error
	ValidatePlanStep(step *spookytypesactions.PlanStep) error
	ValidatePlanDependencies(plan *spookytypesactions.ActionPlan) error

	// Validation rules
	AddValidationRule(rule PlanValidationRule) error
	RemoveValidationRule(name string) error
	GetValidationRules() ([]PlanValidationRule, error)

	// Configuration
	SetStrictMode(strict bool)
	SetValidationLevel(level spookytypesactions.ValidationLevel)
}

// PlanScheduler defines the interface for scheduling action plans
type PlanScheduler interface {
	// Plan scheduling
	SchedulePlan(plan *spookytypesactions.ActionPlan, schedule ScheduleOptions) error
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
	Strategy      spookytypesactions.PlanningStrategy     `json:"strategy"`
	Optimization  spookytypesactions.PlanningOptimization `json:"optimization"`
	Constraints   []spookytypesactions.PlanningConstraint `json:"constraints"`
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
	Plan      *spookytypesactions.ActionPlan `json:"plan"`
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
	Validate(plan *spookytypesactions.ActionPlan) error
	ValidateStep(step *spookytypesactions.PlanStep) error
	GetDescription() string
}

// ValidationManager defines the interface for action validation operations
type ValidationManager interface {
	// Core validation operations
	ValidateAction(action *spookytypesactions.Action) error
	ValidateActionCollection(collection *spookytypesactions.ActionCollection) error
	ValidateActionContext(context *spookytypesactions.ActionContext) error

	// Validator management
	CreateValidator(action *spookytypesactions.Action) (ActionValidator, error)
	GetValidator(action *spookytypesactions.Action) (ActionValidator, error)

	// Custom validation
	AddValidationRule(rule ValidationRule) error
	RemoveValidationRule(name string) error
	GetValidationRules() ([]ValidationRule, error)

	// Configuration
	SetStrictMode(strict bool)
	SetValidationLevel(level spookytypesactions.ValidationLevel)
}

// ValidationRule defines a custom validation rule
type ValidationRule interface {
	Name() string
	Validate(action *spookytypesactions.Action) error
	ValidateCollection(collection *spookytypesactions.ActionCollection) error
	ValidateContext(context *spookytypesactions.ActionContext) error
}

// ValidationResult represents the result of a validation operation
type ValidationResult struct {
	Valid    bool                   `json:"valid"`
	Errors   []ValidationError      `json:"errors"`
	Warnings []ValidationWarning    `json:"warnings"`
	Details  map[string]interface{} `json:"details"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field    string `json:"field"`
	Message  string `json:"message"`
	Code     string `json:"code"`
	Severity string `json:"severity"`
}

// ValidationWarning represents a validation warning
type ValidationWarning struct {
	Field    string `json:"field"`
	Message  string `json:"message"`
	Code     string `json:"code"`
	Severity string `json:"severity"`
}
