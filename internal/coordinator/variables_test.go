package coordinator

import (
	"testing"
	"time"

	spookyfactstypes "spooky/internal/facts/types"
	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookyloggingtypes "spooky/internal/logging/types"

	"github.com/stretchr/testify/assert"
)

func TestNewCoordinatorVariablesIntegration(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorVariablesIntegration(nil, logger)

	assert.NotNil(t, integration)
	assert.NotNil(t, integration.logger)
}

func TestLoadVariables(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorVariablesIntegration(nil, logger)

	projectPath := "/test/project"
	context, err := integration.LoadVariables(projectPath)

	// Should handle nil variables manager gracefully
	assert.NoError(t, err)
	assert.NotNil(t, context)
	assert.Equal(t, projectPath, context.ProjectPath)
}

func TestResolveVariables(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorVariablesIntegration(nil, logger)

	variablesContext := &spookyinterfaces.VariablesContext{
		ResolvedVariables: make(map[string]interface{}),
		VariableContext:   make(map[string]interface{}),
		ResolutionContext: make(map[string]interface{}),
	}

	factsContext := &spookyinterfaces.FactsContext{
		MachineFacts: make(map[string]*spookyfactstypes.FactCollection),
	}

	err := integration.ResolveVariables(variablesContext, factsContext)
	assert.NoError(t, err)
}

func TestValidateVariables(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorVariablesIntegration(nil, logger)

	context := &spookyinterfaces.VariablesContext{
		ResolvedVariables: make(map[string]interface{}),
		VariableContext:   make(map[string]interface{}),
		ResolutionContext: make(map[string]interface{}),
	}

	err := integration.ValidateVariables(context)
	assert.NoError(t, err)
}

func TestSubstituteVariables(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorVariablesIntegration(nil, logger)

	template := "Hello {{name}}!"
	variablesContext := &spookyinterfaces.VariablesContext{
		ResolvedVariables: map[string]interface{}{
			"name": "World",
		},
	}

	result, err := integration.SubstituteVariables(template, variablesContext)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestVariablesIntegrationPerformance(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorVariablesIntegration(nil, logger)

	// Test performance requirement: variable resolution < 50ms
	variablesContext := &spookyinterfaces.VariablesContext{
		ResolvedVariables: make(map[string]interface{}),
		VariableContext:   make(map[string]interface{}),
		ResolutionContext: make(map[string]interface{}),
	}

	factsContext := &spookyinterfaces.FactsContext{
		MachineFacts: make(map[string]*spookyfactstypes.FactCollection),
	}

	start := time.Now()
	err := integration.ResolveVariables(variablesContext, factsContext)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.Less(t, duration, 50*time.Millisecond, "Variable resolution should complete within 50ms")
}

func TestVariablesIntegrationConcurrent(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorVariablesIntegration(nil, logger)

	// Test concurrent operations
	const numOperations = 10
	results := make(chan error, numOperations)

	for i := 0; i < numOperations; i++ {
		go func() {
			variablesContext := &spookyinterfaces.VariablesContext{
				ResolvedVariables: make(map[string]interface{}),
				VariableContext:   make(map[string]interface{}),
				ResolutionContext: make(map[string]interface{}),
			}

			factsContext := &spookyinterfaces.FactsContext{
				MachineFacts: make(map[string]*spookyfactstypes.FactCollection),
			}

			err := integration.ResolveVariables(variablesContext, factsContext)
			results <- err
		}()
	}

	// Collect results
	for i := 0; i < numOperations; i++ {
		err := <-results
		assert.NoError(t, err)
	}
}
