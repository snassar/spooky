package spookytypesactions

import (
	"time"

	spookytypescommon "spooky/internal/types/common"
	spookytypesmachines "spooky/internal/types/machines"
	spookytypesvariables "spooky/internal/types/variables"
)

// ActionExecutionContext represents the context for action execution
type ActionExecutionContext struct {
	spookytypescommon.TimestampedEntity

	// Project context
	ProjectPath string `json:"project_path"`
	ProjectName string `json:"project_name"`

	// Action context
	ActionName string `json:"action_name"`
	ActionType string `json:"action_type"`

	// Machine context
	TargetMachines []spookytypesmachines.Machine `json:"target_machines"`
	MachineNames   []string                      `json:"machine_names"`

	// Variable context
	Variables map[string]*spookytypesvariables.Variable `json:"variables"`
	Resolved  map[string]interface{}                    `json:"resolved"`

	// Execution settings
	Parallel     bool          `json:"parallel"`
	Timeout      time.Duration `json:"timeout"`
	Retries      int           `json:"retries"`
	RetryDelay   time.Duration `json:"retry_delay"`
	DryRun       bool          `json:"dry_run"`
	AllowFailure bool          `json:"allow_failure"`

	// Environment
	Environment      map[string]string `json:"environment"`
	WorkingDirectory string            `json:"working_directory"`
	User             string            `json:"user"`
	Sudo             bool              `json:"sudo"`

	// Resource limits
	MaxConcurrent  int             `json:"max_concurrent"`
	ResourceLimits *ResourceLimits `json:"resource_limits,omitempty"`

	// Execution state
	StartTime *time.Time     `json:"start_time,omitempty"`
	EndTime   *time.Time     `json:"end_time,omitempty"`
	Status    string         `json:"status"`
	Error     error          `json:"error,omitempty"`
	Results   []ActingResult `json:"results,omitempty"`
}

// ActingSession represents a session for running actions
type ActingSession struct {
	spookytypescommon.TimestampedEntity

	// Session identification
	SessionID   string `json:"session_id"`
	ProjectPath string `json:"project_path"`

	// Actions to run
	Actions []Action `json:"actions"`

	// Target machines
	Machines []spookytypesmachines.Machine `json:"machines"`

	// Execution context
	Context *ActionExecutionContext `json:"context"`

	// Session state
	Status    string     `json:"status"`
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Progress  float64    `json:"progress"`
	Completed int        `json:"completed"`
	Total     int        `json:"total"`
	Failed    int        `json:"failed"`
	Skipped   int        `json:"skipped"`

	// Results
	Results []ActingResult `json:"results,omitempty"`
	Errors  []error        `json:"errors,omitempty"`
}

// ActingResult represents the result of running an action
type ActingResult struct {
	spookytypescommon.TimestampedEntity

	// Action identification
	ActionName string `json:"action_name"`
	ActionType string `json:"action_type"`

	// Machine identification
	MachineName string `json:"machine_name"`
	MachineHost string `json:"machine_host"`

	// Execution details
	StartTime *time.Time    `json:"start_time,omitempty"`
	EndTime   *time.Time    `json:"end_time,omitempty"`
	Duration  time.Duration `json:"duration"`
	Status    string        `json:"status"`
	ExitCode  int           `json:"exit_code"`

	// Output
	Stdout   string `json:"stdout,omitempty"`
	Stderr   string `json:"stderr,omitempty"`
	Combined string `json:"combined,omitempty"`

	// Error information
	Error        error  `json:"error,omitempty"`
	ErrorType    string `json:"error_type,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`

	// Metadata
	RetryCount int               `json:"retry_count"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// SessionStatus represents the status of an acting session
type SessionStatus string

const (
	SessionStatusPending   SessionStatus = "pending"
	SessionStatusRunning   SessionStatus = "running"
	SessionStatusCompleted SessionStatus = "completed"
	SessionStatusFailed    SessionStatus = "failed"
	SessionStatusCancelled SessionStatus = "cancelled"
)

// ResultStatus represents the status of an action result
type ResultStatus string

const (
	ResultStatusPending   ResultStatus = "pending"
	ResultStatusRunning   ResultStatus = "running"
	ResultStatusCompleted ResultStatus = "completed"
	ResultStatusFailed    ResultStatus = "failed"
	ResultStatusSkipped   ResultStatus = "skipped"
	ResultStatusCancelled ResultStatus = "cancelled"
)
