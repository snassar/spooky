package coordinator

import (
	"testing"
	"time"

	"spooky/internal/interfaces"
	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookyloggingtypes "spooky/internal/logging/types"

	"github.com/stretchr/testify/assert"
)

func TestNewCoordinatorMachinesIntegration(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorMachinesIntegration(nil, logger)

	assert.NotNil(t, integration)
	assert.NotNil(t, integration.logger)
}

func TestLoadMachines(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorMachinesIntegration(nil, logger)

	projectPath := "/test/project"
	context, err := integration.LoadMachines(projectPath)

	// Should handle nil machines manager gracefully
	assert.NoError(t, err)
	assert.NotNil(t, context)
	assert.Equal(t, projectPath, context.ProjectPath)
}

func TestValidateMachines(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorMachinesIntegration(nil, logger)

	context := &spookyinterfaces.MachinesContext{
		Machines: make(map[string]*spookyinterfaces.Machine),
	}

	err := integration.ValidateMachines(context)
	assert.NoError(t, err)
}

func TestConnectToMachine(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorMachinesIntegration(nil, logger)

	connectionContext := &interfaces.ConnectionContext{
		Timeout: 30,
	}

	err := integration.ConnectToMachine("test-machine", connectionContext)
	// Should handle nil machines manager gracefully
	assert.NoError(t, err)
}

func TestPingMachine(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorMachinesIntegration(nil, logger)

	err := integration.PingMachine("test-machine")
	// Should handle nil machines manager gracefully
	assert.NoError(t, err)
}

func TestMachinesIntegrationPerformance(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorMachinesIntegration(nil, logger)

	// Test performance requirement: machine loading < 100ms
	projectPath := "/test/project"

	start := time.Now()
	context, err := integration.LoadMachines(projectPath)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.NotNil(t, context)
	assert.Less(t, duration, 100*time.Millisecond, "Machine loading should complete within 100ms")
}

func TestMachinesIntegrationConcurrent(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorMachinesIntegration(nil, logger)

	// Test concurrent operations
	const numOperations = 10
	results := make(chan error, numOperations)

	for i := 0; i < numOperations; i++ {
		go func() {
			projectPath := "/test/project"
			_, err := integration.LoadMachines(projectPath)
			results <- err
		}()
	}

	// Collect results
	for i := 0; i < numOperations; i++ {
		err := <-results
		assert.NoError(t, err)
	}
}
