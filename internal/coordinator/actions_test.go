package coordinator

import (
	"testing"
	"time"

	spookyactionstypes "spooky/internal/actions/types"
	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookyloggingtypes "spooky/internal/logging/types"

	"github.com/stretchr/testify/assert"
)

func TestNewCoordinatorActionsIntegration(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorActionsIntegration(nil, logger)

	assert.NotNil(t, integration)
	assert.NotNil(t, integration.logger)
}

func TestLoadActions(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorActionsIntegration(nil, logger)

	projectPath := "/test/project"
	context, err := integration.LoadActions(projectPath)

	// Should handle nil actions manager gracefully
	assert.NoError(t, err)
	assert.NotNil(t, context)
	assert.Equal(t, projectPath, context.ProjectPath)
}

func TestValidateAction(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorActionsIntegration(nil, logger)

	action := &spookyactionstypes.Action{
		Name:    "test-action",
		Command: "echo 'test'",
	}

	execContext := &spookyinterfaces.ActionExecutionContext{
		BaseContext: spookyinterfaces.BaseContext{
			ProjectPath: "/test/project",
		},
		Action: action,
	}

	err := integration.ValidateAction(action, execContext)
	assert.NoError(t, err)
}

func TestPrepareActionForExecution(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorActionsIntegration(nil, logger)

	action := &spookyactionstypes.Action{
		Name:    "test-action",
		Command: "echo 'test'",
	}

	execContext := &spookyinterfaces.ActionExecutionContext{
		BaseContext: spookyinterfaces.BaseContext{
			ProjectPath: "/test/project",
		},
		Action: action,
	}

	err := integration.PrepareActionForExecution(action, execContext)
	assert.NoError(t, err)
}

func TestExecuteAction(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorActionsIntegration(nil, logger)

	action := &spookyactionstypes.Action{
		Name:    "test-action",
		Command: "echo 'test'",
	}

	execContext := &spookyinterfaces.ActionExecutionContext{
		BaseContext: spookyinterfaces.BaseContext{
			ProjectPath: "/test/project",
		},
		Action: action,
	}

	err := integration.ExecuteAction(action, execContext)
	// Should handle nil actions manager gracefully
	assert.NoError(t, err)
}

func TestActionsIntegrationPerformance(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorActionsIntegration(nil, logger)

	// Test performance requirement: action validation < 50ms
	action := &spookyactionstypes.Action{
		Name:    "test-action",
		Command: "echo 'test'",
	}

	execContext := &spookyinterfaces.ActionExecutionContext{
		BaseContext: spookyinterfaces.BaseContext{
			ProjectPath: "/test/project",
		},
		Action: action,
	}

	start := time.Now()
	err := integration.ValidateAction(action, execContext)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.Less(t, duration, 50*time.Millisecond, "Action validation should complete within 50ms")
}

func TestActionsIntegrationConcurrent(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorActionsIntegration(nil, logger)

	// Test concurrent operations
	const numOperations = 10
	results := make(chan error, numOperations)

	for i := 0; i < numOperations; i++ {
		go func() {
			action := &spookyactionstypes.Action{
				Name:    "test-action",
				Command: "echo 'test'",
			}

			execContext := &spookyinterfaces.ActionExecutionContext{
				BaseContext: spookyinterfaces.BaseContext{
					ProjectPath: "/test/project",
				},
				Action: action,
			}

			err := integration.ValidateAction(action, execContext)
			results <- err
		}()
	}

	// Collect results
	for i := 0; i < numOperations; i++ {
		err := <-results
		assert.NoError(t, err)
	}
}
