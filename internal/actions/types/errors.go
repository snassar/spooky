package types

import (
	"fmt"
	"time"
)

// ActionError represents an error that occurred during action operations
type ActionError struct {
	ActionName string                 `json:"action_name"`
	ErrorType  string                 `json:"error_type"`
	Message    string                 `json:"message"`
	Timestamp  time.Time              `json:"timestamp"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

// Error returns the error message
func (e *ActionError) Error() string {
	if e.ActionName != "" {
		return fmt.Sprintf("action '%s': %s", e.ActionName, e.Message)
	}
	return e.Message
}

// ActingError represents an error that occurred during action acting
type ActingError struct {
	ActionName string                 `json:"action_name"`
	MachineID  string                 `json:"machine_id,omitempty"`
	SessionID  string                 `json:"session_id,omitempty"`
	ErrorType  string                 `json:"error_type"`
	Message    string                 `json:"message"`
	Timestamp  time.Time              `json:"timestamp"`
	ExitCode   int                    `json:"exit_code,omitempty"`
	Output     string                 `json:"output,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

// Error returns the error message
func (e *ActingError) Error() string {
	if e.MachineID != "" {
		return fmt.Sprintf("action '%s' on machine '%s': %s", e.ActionName, e.MachineID, e.Message)
	}
	if e.ActionName != "" {
		return fmt.Sprintf("action '%s': %s", e.ActionName, e.Message)
	}
	return e.Message
}

// PlanningError represents an error that occurred during action planning
type PlanningError struct {
	ActionName string                 `json:"action_name"`
	PlanID     string                 `json:"plan_id,omitempty"`
	ErrorType  string                 `json:"error_type"`
	Message    string                 `json:"message"`
	Timestamp  time.Time              `json:"timestamp"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

// Error returns the error message
func (e *PlanningError) Error() string {
	if e.ActionName != "" {
		return fmt.Sprintf("planning action '%s': %s", e.ActionName, e.Message)
	}
	return e.Message
}

// ValidationError represents an error that occurred during action validation
type ValidationError struct {
	ActionName string                 `json:"action_name"`
	Field      string                 `json:"field,omitempty"`
	ErrorType  string                 `json:"error_type"`
	Message    string                 `json:"message"`
	Timestamp  time.Time              `json:"timestamp"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

// Error returns the error message
func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error in action '%s' field '%s': %s", e.ActionName, e.Field, e.Message)
	}
	if e.ActionName != "" {
		return fmt.Sprintf("validation error in action '%s': %s", e.ActionName, e.Message)
	}
	return e.Message
}

// MergingError represents an error that occurred during action merging
type MergingError struct {
	ActionName string                 `json:"action_name"`
	Policy     string                 `json:"policy,omitempty"`
	ErrorType  string                 `json:"error_type"`
	Message    string                 `json:"message"`
	Timestamp  time.Time              `json:"timestamp"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

// Error returns the error message
func (e *MergingError) Error() string {
	if e.Policy != "" {
		return fmt.Sprintf("merging error for action '%s' with policy '%s': %s", e.ActionName, e.Policy, e.Message)
	}
	if e.ActionName != "" {
		return fmt.Sprintf("merging error for action '%s': %s", e.ActionName, e.Message)
	}
	return e.Message
}

// PerformanceError represents an error that occurred during performance optimization
type PerformanceError struct {
	ActionName string                 `json:"action_name"`
	ErrorType  string                 `json:"error_type"`
	Message    string                 `json:"message"`
	Timestamp  time.Time              `json:"timestamp"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

// Error returns the error message
func (e *PerformanceError) Error() string {
	if e.ActionName != "" {
		return fmt.Sprintf("performance error for action '%s': %s", e.ActionName, e.Message)
	}
	return e.Message
}

// Error types
const (
	ErrorTypeValidation  = "validation"
	ErrorTypeActing      = "acting"
	ErrorTypePlanning    = "planning"
	ErrorTypeMerging     = "merging"
	ErrorTypePerformance = "performance"
	ErrorTypeTimeout     = "timeout"
	ErrorTypeDependency  = "dependency"
	ErrorTypeResource    = "resource"
	ErrorTypePermission  = "permission"
	ErrorTypeNetwork     = "network"
	ErrorTypeInternal    = "internal"
)

// NewActionError creates a new ActionError
func NewActionError(actionName, errorType, message string) *ActionError {
	return &ActionError{
		ActionName: actionName,
		ErrorType:  errorType,
		Message:    message,
		Timestamp:  time.Now(),
		Details:    make(map[string]interface{}),
	}
}

// NewActingError creates a new ActingError
func NewActingError(actionName, machineID, errorType, message string) *ActingError {
	return &ActingError{
		ActionName: actionName,
		MachineID:  machineID,
		ErrorType:  errorType,
		Message:    message,
		Timestamp:  time.Now(),
		Details:    make(map[string]interface{}),
	}
}

// NewPlanningError creates a new PlanningError
func NewPlanningError(actionName, errorType, message string) *PlanningError {
	return &PlanningError{
		ActionName: actionName,
		ErrorType:  errorType,
		Message:    message,
		Timestamp:  time.Now(),
		Details:    make(map[string]interface{}),
	}
}

// NewValidationError creates a new ValidationError
func NewValidationError(actionName, field, errorType, message string) *ValidationError {
	return &ValidationError{
		ActionName: actionName,
		Field:      field,
		ErrorType:  errorType,
		Message:    message,
		Timestamp:  time.Now(),
		Details:    make(map[string]interface{}),
	}
}

// NewMergingError creates a new MergingError
func NewMergingError(actionName, policy, errorType, message string) *MergingError {
	return &MergingError{
		ActionName: actionName,
		Policy:     policy,
		ErrorType:  errorType,
		Message:    message,
		Timestamp:  time.Now(),
		Details:    make(map[string]interface{}),
	}
}

// NewPerformanceError creates a new PerformanceError
func NewPerformanceError(actionName, errorType, message string) *PerformanceError {
	return &PerformanceError{
		ActionName: actionName,
		ErrorType:  errorType,
		Message:    message,
		Timestamp:  time.Now(),
		Details:    make(map[string]interface{}),
	}
}
