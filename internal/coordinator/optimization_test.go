package coordinator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculateOptimalParallel(t *testing.T) {
	// Test with valid inputs
	config := &OptimizationConfig{
		MaxParallelWorkers: 20,
		DefaultTimeout:     30 * time.Second,
	}

	parallel := CalculateOptimalParallel(10, 5, config)
	assert.GreaterOrEqual(t, parallel, 2)
	assert.LessOrEqual(t, parallel, 20)

	// Test with zero requested parallel
	parallel = CalculateOptimalParallel(0, 5, config)
	assert.GreaterOrEqual(t, parallel, 2)
	assert.LessOrEqual(t, parallel, 5)
}

func TestCalculateOptimalTimeout(t *testing.T) {
	// Test with valid inputs
	config := &OptimizationConfig{
		MaxParallelWorkers: 20,
		DefaultTimeout:     30 * time.Second,
	}

	timeout := CalculateOptimalTimeout(5, 10, config)
	assert.Greater(t, timeout, 30*time.Second)
	assert.LessOrEqual(t, timeout, 5*time.Minute)

	// Test with zero target count
	timeout = CalculateOptimalTimeout(0, 10, config)
	assert.Equal(t, 30*time.Second, timeout) // Should use default
}

func TestEstimateTargetCount(t *testing.T) {
	// Test with machine specified
	count := EstimateTargetCount("", "machine1", "")
	assert.Equal(t, 1, count)

	// Test with tags specified
	count = EstimateTargetCount("", "", "web,db")
	assert.Equal(t, 5, count)

	// Test with no filters
	count = EstimateTargetCount("", "", "")
	assert.Equal(t, 10, count)
}

func TestOptimizeCacheKey(t *testing.T) {
	// Test with valid inputs
	params := map[string]interface{}{
		"machine": "test-machine",
		"tags":    []string{"web", "db"},
	}

	cacheKey := OptimizeCacheKey("/test/project", "facts-gather", params)
	assert.NotEmpty(t, cacheKey)
	assert.Contains(t, cacheKey, "/test/project")
	assert.Contains(t, cacheKey, "facts-gather")
}

func TestOptimizeBatchSize(t *testing.T) {
	// Test with valid inputs
	batchSize := OptimizeBatchSize(100, 10)
	assert.GreaterOrEqual(t, batchSize, 1)
	assert.LessOrEqual(t, batchSize, 100)

	// Test with zero total count
	batchSize = OptimizeBatchSize(0, 10)
	assert.Equal(t, 1, batchSize)

	// Test with total count less than parallel
	batchSize = OptimizeBatchSize(5, 10)
	assert.Equal(t, 1, batchSize)
}

func TestDefaultOptimizationConfig(t *testing.T) {
	config := DefaultOptimizationConfig()
	assert.NotNil(t, config)
	assert.Greater(t, config.MaxParallelWorkers, 0)
	assert.Greater(t, config.DefaultTimeout, time.Duration(0))
	assert.True(t, config.CacheEnabled)
	assert.True(t, config.CompressionEnabled)
}
