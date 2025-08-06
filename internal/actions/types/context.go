package types

import (
	"time"

	configtypes "spooky/internal/config/types"
	factstypes "spooky/internal/facts/types"
	vartypes "spooky/internal/variables/types"
)

// ActionContext represents the context for action execution
type ActionContext struct {
	// Project information
	ProjectPath string `hcl:"project_path,optional"`
	ProjectName string `hcl:"project_name,optional"`

	// Execution context
	Facts     *factstypes.FactCollection `hcl:"facts,block"`
	Variables *vartypes.VariableContext  `hcl:"variables,block"`
	Machines  []configtypes.Machine      `hcl:"machines,block"`

	// Execution settings
	Timeout       time.Duration `hcl:"timeout,optional"`
	Parallel      bool          `hcl:"parallel,optional"`
	DryRun        bool          `hcl:"dry_run,optional"`
	MaxConcurrent int           `hcl:"max_concurrent,optional"`

	// Environment
	Environment      map[string]string `hcl:"environment,optional"`
	WorkingDirectory string            `hcl:"working_directory,optional"`
	User             string            `hcl:"user,optional"`
	Sudo             bool              `hcl:"sudo,optional"`

	// Metadata
	SessionID string    `hcl:"session_id,optional"`
	CreatedAt time.Time `hcl:"created_at,optional"`
	UpdatedAt time.Time `hcl:"updated_at,optional"`

	// Custom context data
	CustomData map[string]interface{} `hcl:"custom_data,optional"`
}

// ActingContext represents the context for action acting (execution)
type ActingContext struct {
	// Base action context
	*ActionContext

	// Acting-specific settings
	ActorID     string      `hcl:"actor_id,optional"`
	ActorType   string      `hcl:"actor_type,optional"`
	ActingState ActingState `hcl:"acting_state,optional"`

	// Progress tracking
	Progress  float64    `hcl:"progress,optional"`
	StartTime *time.Time `hcl:"start_time,optional"`
	EndTime   *time.Time `hcl:"end_time,optional"`

	// Acting metadata
	ActingSessionID string `hcl:"acting_session_id,optional"`
	ActingAttempt   int    `hcl:"acting_attempt,optional"`
	ActingRetries   int    `hcl:"acting_retries,optional"`
}

// ActingState represents the state of action acting
type ActingState string

const (
	ActingStatePending   ActingState = "pending"
	ActingStateRunning   ActingState = "running"
	ActingStateCompleted ActingState = "completed"
	ActingStateFailed    ActingState = "failed"
	ActingStateCancelled ActingState = "cancelled"
	ActingStateSkipped   ActingState = "skipped"
)

// ActingStatus represents the status of action acting
type ActingStatus string

const (
	ActingStatusPending   ActingStatus = "pending"
	ActingStatusRunning   ActingStatus = "running"
	ActingStatusCompleted ActingStatus = "completed"
	ActingStatusFailed    ActingStatus = "failed"
	ActingStatusCancelled ActingStatus = "cancelled"
	ActingStatusSkipped   ActingStatus = "skipped"
)

// ActingSession represents a session for action acting
type ActingSession struct {
	// Session identification
	SessionID  string `hcl:"session_id,optional"`
	ActionName string `hcl:"action_name,optional"`

	// Session state
	Status    ActingStatus  `hcl:"status,optional"`
	StartTime *time.Time    `hcl:"start_time,optional"`
	EndTime   *time.Time    `hcl:"end_time,optional"`
	Duration  time.Duration `hcl:"duration,optional"`

	// Session results
	Results []*ActingResult `hcl:"results,block"`
	Error   error           `hcl:"error,optional"`

	// Session metadata
	CreatedAt time.Time              `hcl:"created_at,optional"`
	UpdatedAt time.Time              `hcl:"updated_at,optional"`
	Metadata  map[string]interface{} `hcl:"metadata,optional"`
}

// ActingResult represents the result of an action acting
type ActingResult struct {
	// Result identification
	ResultID   string `hcl:"result_id,optional"`
	ActionName string `hcl:"action_name,optional"`
	MachineID  string `hcl:"machine_id,optional"`

	// Result state
	Status    ActingStatus  `hcl:"status,optional"`
	StartTime *time.Time    `hcl:"start_time,optional"`
	EndTime   *time.Time    `hcl:"end_time,optional"`
	Duration  time.Duration `hcl:"duration,optional"`

	// Result data
	Output   string `hcl:"output,optional"`
	Error    string `hcl:"error,optional"`
	ExitCode int    `hcl:"exit_code,optional"`

	// Result metadata
	CreatedAt time.Time              `hcl:"created_at,optional"`
	Metadata  map[string]interface{} `hcl:"metadata,optional"`
}

// NewActionContext creates a new ActionContext
func NewActionContext(projectPath string) *ActionContext {
	return &ActionContext{
		ProjectPath: projectPath,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Environment: make(map[string]string),
		CustomData:  make(map[string]interface{}),
	}
}

// NewActingContext creates a new ActingContext
func NewActingContext(actionContext *ActionContext) *ActingContext {
	return &ActingContext{
		ActionContext: actionContext,
		ActingState:   ActingStatePending,
	}
}

// NewActingSession creates a new ActingSession
func NewActingSession(actionName string) *ActingSession {
	now := time.Now()
	return &ActingSession{
		ActionName: actionName,
		Status:     ActingStatusPending,
		StartTime:  &now,
		CreatedAt:  now,
		UpdatedAt:  now,
		Results:    make([]*ActingResult, 0),
		Metadata:   make(map[string]interface{}),
	}
}

// NewActingResult creates a new ActingResult
func NewActingResult(actionName, machineID string) *ActingResult {
	now := time.Now()
	return &ActingResult{
		ActionName: actionName,
		MachineID:  machineID,
		Status:     ActingStatusPending,
		StartTime:  &now,
		CreatedAt:  now,
		Metadata:   make(map[string]interface{}),
	}
}
