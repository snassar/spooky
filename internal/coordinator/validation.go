package coordinator

import (
	"fmt"
	"strings"
	"time"
)

// ValidationError represents multiple validation errors
type ValidationError struct {
	Errors []error
}

func (ve *ValidationError) Error() string {
	var messages []string
	for _, err := range ve.Errors {
		messages = append(messages, err.Error())
	}
	return fmt.Sprintf("validation failed: %s", strings.Join(messages, "; "))
}

// ValidateProjectPath validates a project path
func ValidateProjectPath(projectPath string) error {
	if projectPath == "" {
		return fmt.Errorf("project path cannot be empty")
	}

	// Add additional validation as needed
	return nil
}

// ValidateMachineNames validates machine names
func ValidateMachineNames(machineNames []string) error {
	if len(machineNames) == 0 {
		return fmt.Errorf("at least one machine must be specified")
	}

	for _, machine := range machineNames {
		if machine == "" {
			return fmt.Errorf("machine name cannot be empty")
		}
	}

	return nil
}

// ValidateParallelWorkers validates parallel worker count
func ValidateParallelWorkers(parallel int) error {
	if parallel < 0 {
		return fmt.Errorf("parallel workers must be non-negative, got %d", parallel)
	}
	if parallel == 1 {
		return fmt.Errorf("parallel workers must be 0 (auto) or 2 or greater, got %d", parallel)
	}
	return nil
}

// ValidateTimeout validates timeout values
func ValidateTimeout(timeout time.Duration) error {
	if timeout <= 0 {
		return fmt.Errorf("timeout must be positive, got %v", timeout)
	}
	return nil
}

// ValidateAction validates action data
func ValidateAction(action interface{}) error {
	if action == nil {
		return fmt.Errorf("action cannot be nil")
	}

	// Add additional validation as needed
	return nil
}

// ValidateTemplate validates template data
func ValidateTemplate(template interface{}) error {
	if template == nil {
		return fmt.Errorf("template cannot be nil")
	}

	// Add additional validation as needed
	return nil
}

// ValidateVariables validates variables data
func ValidateVariables(variables interface{}) error {
	if variables == nil {
		return fmt.Errorf("variables cannot be nil")
	}

	// Add additional validation as needed
	return nil
}
