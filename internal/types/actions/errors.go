package spookytypesactions

import (
	"fmt"
	"time"

	spookytypescommon "spooky/internal/types/common"
)

// ActionError represents an error that occurred during action running
type ActionError struct {
	spookytypescommon.ErrorDetails
	ActionName  string `json:"action_name"`
	ActionType  string `json:"action_type"`
	MachineName string `json:"machine_name,omitempty"`
	ExitCode    int    `json:"exit_code,omitempty"`
	Stdout      string `json:"stdout,omitempty"`
	Stderr      string `json:"stderr,omitempty"`
}

// NewActionError creates a new ActionError
func NewActionError(actionName, actionType, message string) *ActionError {
	return &ActionError{
		ErrorDetails: spookytypescommon.ErrorDetails{
			Message:   message,
			Code:      "ACTION_ERROR",
			Timestamp: time.Now(),
		},
		ActionName: actionName,
		ActionType: actionType,
	}
}

// Error returns the error message
func (e *ActionError) Error() string {
	return fmt.Sprintf("action error [%s/%s]: %s", e.ActionName, e.ActionType, e.Message)
}

// Unwrap returns the underlying error
func (e *ActionError) Unwrap() error {
	return nil
}

// ActingError represents an error that occurred during acting
type ActingError struct {
	spookytypescommon.ErrorDetails
	SessionID   string `json:"session_id"`
	ActionName  string `json:"action_name"`
	MachineName string `json:"machine_name,omitempty"`
	Status      string `json:"status"`
}

// NewActingError creates a new ActingError
func NewActingError(sessionID, actionName, message string) *ActingError {
	return &ActingError{
		ErrorDetails: spookytypescommon.ErrorDetails{
			Message:   message,
			Code:      "ACTING_ERROR",
			Timestamp: time.Now(),
		},
		SessionID:  sessionID,
		ActionName: actionName,
	}
}

// Error returns the error message
func (e *ActingError) Error() string {
	return fmt.Sprintf("acting error [%s/%s]: %s", e.SessionID, e.ActionName, e.Message)
}

// Unwrap returns the underlying error
func (e *ActingError) Unwrap() error {
	return nil
}

// PlanningError represents an error that occurred during action planning
type PlanningError struct {
	spookytypescommon.ErrorDetails
	PlanName   string   `json:"plan_name"`
	ActionName string   `json:"action_name,omitempty"`
	Issues     []string `json:"issues,omitempty"`
}

// NewPlanningError creates a new PlanningError
func NewPlanningError(planName, message string) *PlanningError {
	return &PlanningError{
		ErrorDetails: spookytypescommon.ErrorDetails{
			Message:   message,
			Code:      "PLANNING_ERROR",
			Timestamp: time.Now(),
		},
		PlanName: planName,
	}
}

// Error returns the error message
func (e *PlanningError) Error() string {
	return fmt.Sprintf("planning error [%s]: %s", e.PlanName, e.Message)
}

// Unwrap returns the underlying error
func (e *PlanningError) Unwrap() error {
	return nil
}

// ValidationError represents a validation error for actions
type ValidationError struct {
	spookytypescommon.ErrorDetails
	ActionName string `json:"action_name,omitempty"`
	Field      string `json:"field,omitempty"`
	Value      string `json:"value,omitempty"`
	Rule       string `json:"rule,omitempty"`
}

// NewValidationError creates a new ValidationError
func NewValidationError(actionName, field, message string) *ValidationError {
	return &ValidationError{
		ErrorDetails: spookytypescommon.ErrorDetails{
			Message:   message,
			Code:      "VALIDATION_ERROR",
			Timestamp: time.Now(),
		},
		ActionName: actionName,
		Field:      field,
	}
}

// Error returns the error message
func (e *ValidationError) Error() string {
	if e.ActionName != "" && e.Field != "" {
		return fmt.Sprintf("validation error [%s.%s]: %s", e.ActionName, e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// Unwrap returns the underlying error
func (e *ValidationError) Unwrap() error {
	return nil
}

// DependencyError represents an error with action dependencies
type DependencyError struct {
	spookytypescommon.ErrorDetails
	ActionName     string     `json:"action_name"`
	Dependencies   []string   `json:"dependencies"`
	MissingActions []string   `json:"missing_actions,omitempty"`
	CircularDeps   [][]string `json:"circular_deps,omitempty"`
}

// NewDependencyError creates a new DependencyError
func NewDependencyError(actionName, message string) *DependencyError {
	return &DependencyError{
		ErrorDetails: spookytypescommon.ErrorDetails{
			Message:   message,
			Code:      "DEPENDENCY_ERROR",
			Timestamp: time.Now(),
		},
		ActionName: actionName,
	}
}

// Error returns the error message
func (e *DependencyError) Error() string {
	return fmt.Sprintf("dependency error [%s]: %s", e.ActionName, e.Message)
}

// Unwrap returns the underlying error
func (e *DependencyError) Unwrap() error {
	return nil
}

// ConfigurationError represents an error with action configuration
type ConfigurationError struct {
	spookytypescommon.ErrorDetails
	ActionName string `json:"action_name,omitempty"`
	Field      string `json:"field,omitempty"`
	ConfigType string `json:"config_type,omitempty"`
}

// NewConfigurationError creates a new ConfigurationError
func NewConfigurationError(actionName, field, message string) *ConfigurationError {
	return &ConfigurationError{
		ErrorDetails: spookytypescommon.ErrorDetails{
			Message:   message,
			Code:      "CONFIGURATION_ERROR",
			Timestamp: time.Now(),
		},
		ActionName: actionName,
		Field:      field,
	}
}

// Error returns the error message
func (e *ConfigurationError) Error() string {
	if e.ActionName != "" && e.Field != "" {
		return fmt.Sprintf("configuration error [%s.%s]: %s", e.ActionName, e.Field, e.Message)
	}
	return fmt.Sprintf("configuration error: %s", e.Message)
}

// Unwrap returns the underlying error
func (e *ConfigurationError) Unwrap() error {
	return nil
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
