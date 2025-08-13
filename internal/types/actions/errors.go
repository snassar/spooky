package spookytypesactions

import (
	"fmt"
	"time"

	spookytypescommon "spooky/internal/types/common"
)

// ActionError represents an error that occurred during action execution
type ActionError struct {
	spookytypescommon.ErrorDetails

	// Action identification
	ActionName string `json:"action_name"`
	ActionType string `json:"action_type"`

	// Error context
	ErrorType    string `json:"error_type"`
	ErrorMessage string `json:"error_message"`

	// Execution context
	MachineName string     `json:"machine_name,omitempty"`
	StartTime   *time.Time `json:"start_time,omitempty"`
	EndTime     *time.Time `json:"end_time,omitempty"`

	// Error details
	ExitCode   int    `json:"exit_code,omitempty"`
	Stdout     string `json:"stdout,omitempty"`
	Stderr     string `json:"stderr,omitempty"`
	RetryCount int    `json:"retry_count"`
	MaxRetries int    `json:"max_retries"`
}

// Error returns the error message
func (e *ActionError) Error() string {
	if e.MachineName != "" {
		return fmt.Sprintf("action '%s' failed on machine '%s': %s", e.ActionName, e.MachineName, e.ErrorMessage)
	}
	return fmt.Sprintf("action '%s' failed: %s", e.ActionName, e.ErrorMessage)
}

// Unwrap returns the underlying error
func (e *ActionError) Unwrap() error {
	return nil
}

// NewActionError creates a new action error
func NewActionError(actionName, actionType, errorType, errorMessage string) *ActionError {
	return &ActionError{
		ErrorDetails: spookytypescommon.ErrorDetails{
			Timestamp: time.Now(),
		},
		ActionName:   actionName,
		ActionType:   actionType,
		ErrorType:    errorType,
		ErrorMessage: errorMessage,
	}
}

// ActingError represents an error that occurred during action acting
type ActingError struct {
	spookytypescommon.ErrorDetails

	// Acting identification
	SessionID  string `json:"session_id"`
	ActionName string `json:"action_name"`

	// Error context
	ErrorType    string `json:"error_type"`
	ErrorMessage string `json:"error_message"`

	// Acting context
	MachineName string     `json:"machine_name,omitempty"`
	StartTime   *time.Time `json:"start_time,omitempty"`
	EndTime     *time.Time `json:"end_time,omitempty"`

	// Error details
	ExitCode   int    `json:"exit_code,omitempty"`
	Stdout     string `json:"stdout,omitempty"`
	Stderr     string `json:"stderr,omitempty"`
	RetryCount int    `json:"retry_count"`
	MaxRetries int    `json:"max_retries"`
}

// Error returns the error message
func (e *ActingError) Error() string {
	if e.MachineName != "" {
		return fmt.Sprintf("acting session '%s' action '%s' failed on machine '%s': %s", e.SessionID, e.ActionName, e.MachineName, e.ErrorMessage)
	}
	return fmt.Sprintf("acting session '%s' action '%s' failed: %s", e.SessionID, e.ActionName, e.ErrorMessage)
}

// Unwrap returns the underlying error
func (e *ActingError) Unwrap() error {
	return nil
}

// NewActingError creates a new acting error
func NewActingError(sessionID, actionName, errorType, errorMessage string) *ActingError {
	return &ActingError{
		ErrorDetails: spookytypescommon.ErrorDetails{
			Timestamp: time.Now(),
		},
		SessionID:    sessionID,
		ActionName:   actionName,
		ErrorType:    errorType,
		ErrorMessage: errorMessage,
	}
}

// PlanningError represents an error that occurred during action planning
type PlanningError struct {
	spookytypescommon.ErrorDetails

	// Planning identification
	PlanID   string `json:"plan_id"`
	PlanName string `json:"plan_name"`

	// Error context
	ErrorType    string `json:"error_type"`
	ErrorMessage string `json:"error_message"`

	// Planning context
	ActionName string `json:"action_name,omitempty"`
	Step       string `json:"step,omitempty"`

	// Error details
	ValidationErrors []string `json:"validation_errors,omitempty"`
	DependencyErrors []string `json:"dependency_errors,omitempty"`
	ResourceErrors   []string `json:"resource_errors,omitempty"`
}

// Error returns the error message
func (e *PlanningError) Error() string {
	if e.ActionName != "" {
		return fmt.Sprintf("planning failed for plan '%s' action '%s': %s", e.PlanName, e.ActionName, e.ErrorMessage)
	}
	return fmt.Sprintf("planning failed for plan '%s': %s", e.PlanName, e.ErrorMessage)
}

// Unwrap returns the underlying error
func (e *PlanningError) Unwrap() error {
	return nil
}

// NewPlanningError creates a new planning error
func NewPlanningError(planID, planName, errorType, errorMessage string) *PlanningError {
	return &PlanningError{
		ErrorDetails: spookytypescommon.ErrorDetails{
			Timestamp: time.Now(),
		},
		PlanID:       planID,
		PlanName:     planName,
		ErrorType:    errorType,
		ErrorMessage: errorMessage,
	}
}

// ValidationError represents an error that occurred during action validation
type ValidationError struct {
	spookytypescommon.ErrorDetails

	// Validation identification
	ValidationID string `json:"validation_id"`
	ActionName   string `json:"action_name"`

	// Error context
	ErrorType    string `json:"error_type"`
	ErrorMessage string `json:"error_message"`

	// Validation context
	Field string `json:"field,omitempty"`
	Value string `json:"value,omitempty"`
	Rule  string `json:"rule,omitempty"`

	// Error details
	SchemaErrors     []string `json:"schema_errors,omitempty"`
	DependencyErrors []string `json:"dependency_errors,omitempty"`
	ResourceErrors   []string `json:"resource_errors,omitempty"`
}

// Error returns the error message
func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation failed for action '%s' field '%s': %s", e.ActionName, e.Field, e.ErrorMessage)
	}
	return fmt.Sprintf("validation failed for action '%s': %s", e.ActionName, e.ErrorMessage)
}

// Unwrap returns the underlying error
func (e *ValidationError) Unwrap() error {
	return nil
}

// NewValidationError creates a new validation error
func NewValidationError(validationID, actionName, errorType, errorMessage string) *ValidationError {
	return &ValidationError{
		ErrorDetails: spookytypescommon.ErrorDetails{
			Timestamp: time.Now(),
		},
		ValidationID: validationID,
		ActionName:   actionName,
		ErrorType:    errorType,
		ErrorMessage: errorMessage,
	}
}

// DependencyError represents an error that occurred during dependency resolution
type DependencyError struct {
	spookytypescommon.ErrorDetails

	// Dependency identification
	FromAction string `json:"from_action"`
	ToAction   string `json:"to_action"`

	// Error context
	ErrorType    string `json:"error_type"`
	ErrorMessage string `json:"error_message"`

	// Dependency context
	DependencyType string `json:"dependency_type"`
	Condition      string `json:"condition,omitempty"`

	// Error details
	CircularDependencies []string `json:"circular_dependencies,omitempty"`
	MissingDependencies  []string `json:"missing_dependencies,omitempty"`
	InvalidDependencies  []string `json:"invalid_dependencies,omitempty"`
}

// Error returns the error message
func (e *DependencyError) Error() string {
	return fmt.Sprintf("dependency error from '%s' to '%s': %s", e.FromAction, e.ToAction, e.ErrorMessage)
}

// Unwrap returns the underlying error
func (e *DependencyError) Unwrap() error {
	return nil
}

// NewDependencyError creates a new dependency error
func NewDependencyError(fromAction, toAction, errorType, errorMessage string) *DependencyError {
	return &DependencyError{
		ErrorDetails: spookytypescommon.ErrorDetails{
			Timestamp: time.Now(),
		},
		FromAction:   fromAction,
		ToAction:     toAction,
		ErrorType:    errorType,
		ErrorMessage: errorMessage,
	}
}

// Error type constants
const (
	ErrorTypeValidation     = "validation"
	ErrorTypePlanning       = "planning"
	ErrorTypeActing         = "acting"
	ErrorTypeDependency     = "dependency"
	ErrorTypeResource       = "resource"
	ErrorTypeTimeout        = "timeout"
	ErrorTypeConnection     = "connection"
	ErrorTypeAuthentication = "authentication"
	ErrorTypePermission     = "permission"
	ErrorTypeFileSystem     = "filesystem"
	ErrorTypeNetwork        = "network"
	ErrorTypeSystem         = "system"
	ErrorTypeUnknown        = "unknown"
)
