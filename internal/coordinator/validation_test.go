package coordinator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateProjectPath(t *testing.T) {
	// Test valid project path
	err := ValidateProjectPath("/valid/project/path")
	assert.NoError(t, err)

	// Test empty project path
	err = ValidateProjectPath("")
	assert.Error(t, err)
}

func TestValidateMachineNames(t *testing.T) {
	// Test valid machine names
	machineNames := []string{"machine1", "machine2", "machine3"}
	err := ValidateMachineNames(machineNames)
	assert.NoError(t, err)

	// Test empty machine names
	err = ValidateMachineNames([]string{})
	assert.Error(t, err)

	// Test invalid machine names
	invalidNames := []string{"", "machine-1", "machine 1"}
	err = ValidateMachineNames(invalidNames)
	assert.Error(t, err)
}

func TestValidateParallelWorkers(t *testing.T) {
	// Test valid parallel workers
	err := ValidateParallelWorkers(5)
	assert.NoError(t, err)

	// Test invalid parallel workers
	err = ValidateParallelWorkers(-1)
	assert.Error(t, err)
}

func TestValidateTimeout(t *testing.T) {
	// Test valid timeout
	err := ValidateTimeout(30 * time.Second)
	assert.NoError(t, err)

	// Test invalid timeout
	err = ValidateTimeout(0)
	assert.Error(t, err)

	err = ValidateTimeout(-1 * time.Second)
	assert.Error(t, err)
}

func TestValidateActionData(t *testing.T) {
	// Test valid action
	action := "test-action"
	err := ValidateAction(action)
	assert.NoError(t, err)

	// Test invalid action
	err = ValidateAction(nil)
	assert.Error(t, err)
}

func TestValidateTemplateData(t *testing.T) {
	// Test valid template
	template := "test-template"
	err := ValidateTemplate(template)
	assert.NoError(t, err)

	// Test invalid template
	err = ValidateTemplate(nil)
	assert.Error(t, err)
}
