package performance

import (
	spookyactionstypes "spooky/internal/actions/types"
)

// PerformanceManager defines the interface for action performance operations
type PerformanceManager interface {
	// Core performance operations
	OptimizeAction(action *spookyactionstypes.Action) error
	OptimizeActionCollection(collection *spookyactionstypes.ActionCollection) error
	GetPerformanceMetrics(action *spookyactionstypes.Action) (*spookyactionstypes.PerformanceMetrics, error)

	// Optimizer management
	CreateOptimizer(action *spookyactionstypes.Action) (ActionOptimizer, error)
	GetOptimizer(action *spookyactionstypes.Action) (ActionOptimizer, error)

	// Metrics management
	GetMetrics(action *spookyactionstypes.Action) (*spookyactionstypes.PerformanceMetrics, error)
	ListMetrics() ([]*spookyactionstypes.PerformanceMetrics, error)
	ClearMetrics(action *spookyactionstypes.Action) error

	// Configuration
	SetOptimizationLevel(level spookyactionstypes.OptimizationLevel)
	SetResourceLimits(limits *spookyactionstypes.ResourceLimits)
}

// ActionOptimizer defines the interface for action performance optimization
type ActionOptimizer interface {
	// Core optimization operations
	Optimize(action *spookyactionstypes.Action) error
	GetMetrics(action *spookyactionstypes.Action) (*spookyactionstypes.PerformanceMetrics, error)

	// Optimization management
	SetLevel(level spookyactionstypes.OptimizationLevel)
	GetLevel() spookyactionstypes.OptimizationLevel

	// Configuration
	SetResourceLimits(limits *spookyactionstypes.ResourceLimits)
	SetOptimizationTarget(target spookyactionstypes.OptimizationTarget)
}

// PerformanceAnalyzer defines the interface for performance analysis
type PerformanceAnalyzer interface {
	// Analysis operations
	AnalyzeAction(action *spookyactionstypes.Action) (*PerformanceAnalysis, error)
	AnalyzeCollection(collection *spookyactionstypes.ActionCollection) (*PerformanceAnalysis, error)
	CompareActions(action1, action2 *spookyactionstypes.Action) (*PerformanceComparison, error)

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
