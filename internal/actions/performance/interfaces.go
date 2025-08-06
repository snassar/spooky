package performance

import (
	"spooky/internal/actions/types"
)

// PerformanceManager defines the interface for action performance operations
type PerformanceManager interface {
	// Core performance operations
	OptimizeAction(action *types.Action) error
	OptimizeActionCollection(collection *types.ActionCollection) error
	GetPerformanceMetrics(action *types.Action) (*types.PerformanceMetrics, error)

	// Optimizer management
	CreateOptimizer(action *types.Action) (ActionOptimizer, error)
	GetOptimizer(action *types.Action) (ActionOptimizer, error)

	// Metrics management
	GetMetrics(action *types.Action) (*types.PerformanceMetrics, error)
	ListMetrics() ([]*types.PerformanceMetrics, error)
	ClearMetrics(action *types.Action) error

	// Configuration
	SetOptimizationLevel(level types.OptimizationLevel)
	SetResourceLimits(limits *types.ResourceLimits)
}

// ActionOptimizer defines the interface for action performance optimization
type ActionOptimizer interface {
	// Core optimization operations
	Optimize(action *types.Action) error
	GetMetrics(action *types.Action) (*types.PerformanceMetrics, error)

	// Optimization management
	SetLevel(level types.OptimizationLevel)
	GetLevel() types.OptimizationLevel

	// Configuration
	SetResourceLimits(limits *types.ResourceLimits)
	SetOptimizationTarget(target types.OptimizationTarget)
}

// PerformanceAnalyzer defines the interface for performance analysis
type PerformanceAnalyzer interface {
	// Analysis operations
	AnalyzeAction(action *types.Action) (*PerformanceAnalysis, error)
	AnalyzeCollection(collection *types.ActionCollection) (*PerformanceAnalysis, error)
	CompareActions(action1, action2 *types.Action) (*PerformanceComparison, error)

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
