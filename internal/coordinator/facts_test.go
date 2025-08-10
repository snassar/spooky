package coordinator

import (
	"testing"
	"time"

	spookyfactstypes "spooky/internal/types/facts"
	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookyloggingtypes "spooky/internal/types/logging"

	"github.com/stretchr/testify/assert"
)

func TestNewCoordinatorFactsIntegration(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorFactsIntegration(nil, logger)

	assert.NotNil(t, integration)
	assert.NotNil(t, integration.logger)
}

func TestLoadFacts(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorFactsIntegration(nil, logger)

	machineNames := []string{"machine1", "machine2"}
	context, err := integration.LoadFacts(machineNames)

	// Should handle nil facts manager gracefully
	assert.NoError(t, err)
	assert.NotNil(t, context)
	assert.NotNil(t, context.MachineFacts)
}

func TestValidateFacts(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorFactsIntegration(nil, logger)

	context := &spookyinterfaces.FactsContext{
		MachineFacts: make(map[string]*spookyfactstypes.FactCollection),
		GlobalFacts:  &spookyfactstypes.FactCollection{Facts: make(map[string]*spookyfactstypes.Fact)},
		ProjectFacts: &spookyfactstypes.FactCollection{Facts: make(map[string]*spookyfactstypes.Fact)},
	}

	err := integration.ValidateFacts(context)
	assert.NoError(t, err)
}

func TestOptimizeFactsGathering(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorFactsIntegration(nil, logger)

	machineNames := []string{"machine1", "machine2", "machine3", "machine4", "machine5"}
	parallel := 10

	optimalParallel, optimalTimeout, err := integration.OptimizeFactsGathering(machineNames, parallel)

	assert.NoError(t, err)
	assert.Greater(t, optimalParallel, 0)
	assert.LessOrEqual(t, optimalParallel, parallel)
	assert.Greater(t, optimalTimeout, time.Duration(0))
}

func TestGetFactsForMachine(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorFactsIntegration(nil, logger)

	collection, err := integration.GetFactsForMachine("test-machine")

	// Should handle nil facts manager gracefully
	assert.NoError(t, err)
	assert.NotNil(t, collection)
}

func TestCacheFacts(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorFactsIntegration(nil, logger)

	context := &spookyinterfaces.FactsContext{
		MachineFacts: make(map[string]*spookyfactstypes.FactCollection),
		GlobalFacts:  &spookyfactstypes.FactCollection{Facts: make(map[string]*spookyfactstypes.Fact)},
		ProjectFacts: &spookyfactstypes.FactCollection{Facts: make(map[string]*spookyfactstypes.Fact)},
		CacheKey:     "test-cache-key",
	}

	err := integration.CacheFacts(context)
	assert.NoError(t, err)
}

func TestFactsIntegrationPerformance(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorFactsIntegration(nil, logger)

	// Test performance requirement: context loading < 100ms
	machineNames := []string{"machine1", "machine2", "machine3", "machine4", "machine5"}

	start := time.Now()
	context, err := integration.LoadFacts(machineNames)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.NotNil(t, context)
	assert.Less(t, duration, 100*time.Millisecond, "Facts loading should complete within 100ms")
}

func TestFactsIntegrationConcurrent(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorFactsIntegration(nil, logger)

	// Test concurrent operations
	const numOperations = 5
	results := make(chan error, numOperations)

	for i := 0; i < numOperations; i++ {
		go func() {
			machineNames := []string{"machine1", "machine2"}
			_, err := integration.LoadFacts(machineNames)
			results <- err
		}()
	}

	// Collect results
	for i := 0; i < numOperations; i++ {
		err := <-results
		assert.NoError(t, err)
	}
}
