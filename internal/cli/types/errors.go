package types

import (
	"fmt"
)

// CLIError represents a CLI-specific error
type CLIError struct {
	Code    string
	Message string
	Err     error
}

// Error implements the error interface
func (e *CLIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *CLIError) Unwrap() error {
	return e.Err
}

// NewCLIError creates a new CLI error
func NewCLIError(code, message string, err error) *CLIError {
	return &CLIError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common error codes
const (
	ErrCodeCommandNotFound  = "COMMAND_NOT_FOUND"
	ErrCodeInvalidArguments = "INVALID_ARGUMENTS"
	ErrCodeCommandExecution = "COMMAND_EXECUTION"
	ErrCodeConfiguration    = "CONFIGURATION"
	ErrCodeValidation       = "VALIDATION"
	ErrCodeCoordinator      = "COORDINATOR"
)
