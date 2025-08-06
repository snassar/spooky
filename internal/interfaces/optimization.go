package interfaces

import (
	"fmt"
	"runtime"
	"time"
)

// OptimizationConfig provides optimization settings
type OptimizationConfig struct {
	MaxParallelWorkers int
	DefaultTimeout     time.Duration
	CacheEnabled       bool
	CompressionEnabled bool
}

// DefaultOptimizationConfig returns default optimization settings
func DefaultOptimizationConfig() *OptimizationConfig {
	return &OptimizationConfig{
		MaxParallelWorkers: runtime.NumCPU() * 2,
		DefaultTimeout:     30 * time.Second,
		CacheEnabled:       true,
		CompressionEnabled: true,
	}
}

// CalculateOptimalParallel calculates optimal parallel workers
func CalculateOptimalParallel(requestedParallel, targetCount int, config *OptimizationConfig) int {
	if requestedParallel > 0 {
		return minInt(requestedParallel, config.MaxParallelWorkers)
	}

	// Auto-calculate based on target count
	optimal := minInt(targetCount, config.MaxParallelWorkers)
	return maxInt(optimal, 2) // Minimum of 2 workers
}

// CalculateOptimalTimeout calculates optimal timeout
func CalculateOptimalTimeout(targetCount, parallel int, config *OptimizationConfig) time.Duration {
	baseTimeout := config.DefaultTimeout
	perTarget := 2 * time.Second // 2 seconds per target

	calculated := baseTimeout + (time.Duration(targetCount) * perTarget / time.Duration(parallel))
	if calculated > 5*time.Minute {
		return 5 * time.Minute // Cap at 5 minutes
	}
	return calculated
}

// EstimateTargetCount estimates the number of targets for optimization
func EstimateTargetCount(_, machine, tags string) int {
	// Default estimation - can be enhanced with actual project analysis
	if machine != "" {
		return 1
	}

	if tags != "" {
		return 5 // Estimate 5 machines per tag
	}

	return 10 // Default estimate for full project
}

// OptimizeCacheKey generates an optimized cache key
func OptimizeCacheKey(projectPath, operation string, params map[string]interface{}) string {
	// Simple cache key generation - can be enhanced
	return fmt.Sprintf("%s:%s:%v", projectPath, operation, params)
}

// OptimizeBatchSize calculates optimal batch size for operations
func OptimizeBatchSize(totalCount, parallel int) int {
	if totalCount <= parallel {
		return 1
	}

	batchSize := totalCount / parallel
	if batchSize < 1 {
		return 1
	}

	return minInt(batchSize, 100) // Cap at 100 items per batch
}

// Helper functions for min/max
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
