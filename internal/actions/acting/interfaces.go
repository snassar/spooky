package acting

import (
	"context"
	"time"

	spookyactionstypes "spooky/internal/actions/types"
)

// ActingManager defines the interface for action execution operations
type ActingManager interface {
	// Core acting operations
	ExecuteAction(ctx context.Context, action *spookyactionstypes.Action, context *spookyactionstypes.ActionContext) (*spookyactionstypes.ActingSession, error)
	ExecuteActionCollection(ctx context.Context, collection *spookyactionstypes.ActionCollection, context *spookyactionstypes.ActionContext) (*spookyactionstypes.ActingSession, error)
	PrepareAction(action *spookyactionstypes.Action, context *spookyactionstypes.ActionContext) error

	// Actor management
	CreateActor(action *spookyactionstypes.Action, context *spookyactionstypes.ActionContext) (Actor, error)
	GetActor(action *spookyactionstypes.Action) (Actor, error)

	// Session management
	GetSession(sessionID string) (*spookyactionstypes.ActingSession, error)
	ListSessions() ([]*spookyactionstypes.ActingSession, error)
	CancelSession(sessionID string) error

	// Configuration
	SetDefaultTimeout(timeout time.Duration)
	SetDefaultParallel(parallel bool)
	SetMaxConcurrent(maxConcurrent int)
}

// Actor defines the interface for individual action execution
type Actor interface {
	// Core acting operations
	Execute(ctx context.Context, context *spookyactionstypes.ActionContext) (*spookyactionstypes.ActingResult, error)
	Prepare(context *spookyactionstypes.ActionContext) error
	Cancel() error

	// State management
	GetState() spookyactionstypes.ActingState
	GetProgress() float64
	GetStatus() spookyactionstypes.ActingStatus

	// Configuration
	SetTimeout(timeout time.Duration)
	SetParallel(parallel bool)
}

// ActingExecutor defines the interface for executing actions on machines
type ActingExecutor interface {
	// Core execution operations
	ExecuteCommand(ctx context.Context, command string, context *spookyactionstypes.ActionContext) (*spookyactionstypes.ActingResult, error)
	ExecuteScript(ctx context.Context, script string, context *spookyactionstypes.ActionContext) (*spookyactionstypes.ActingResult, error)
	ExecuteTemplate(ctx context.Context, template *spookyactionstypes.TemplateConfig, context *spookyactionstypes.ActionContext) (*spookyactionstypes.ActingResult, error)

	// Machine management
	GetMachine(machineID string) (Machine, error)
	ListMachines() ([]Machine, error)
	ConnectMachine(machineID string) error
	DisconnectMachine(machineID string) error

	// Configuration
	SetConnectionTimeout(timeout time.Duration)
	SetCommandTimeout(timeout time.Duration)
	SetRetryAttempts(attempts int)
	SetRetryDelay(delay time.Duration)
}

// Machine defines the interface for machine operations
type Machine interface {
	// Machine information
	GetID() string
	GetName() string
	GetHost() string
	GetUser() string
	GetPort() int

	// Connection management
	Connect() error
	Disconnect() error
	IsConnected() bool

	// Execution operations
	ExecuteCommand(ctx context.Context, command string) (*spookyactionstypes.ActingResult, error)
	ExecuteScript(ctx context.Context, script string) (*spookyactionstypes.ActingResult, error)
	UploadFile(ctx context.Context, localPath, remotePath string) error
	DownloadFile(ctx context.Context, remotePath, localPath string) error

	// Configuration
	SetTimeout(timeout time.Duration)
	SetSudo(sudo bool)
	SetWorkingDirectory(dir string)
	SetEnvironment(env map[string]string)
}

// ActingSessionManager defines the interface for managing acting sessions
type ActingSessionManager interface {
	// Session management
	CreateSession(actionName string) (*spookyactionstypes.ActingSession, error)
	GetSession(sessionID string) (*spookyactionstypes.ActingSession, error)
	ListSessions() ([]*spookyactionstypes.ActingSession, error)
	UpdateSession(session *spookyactionstypes.ActingSession) error
	DeleteSession(sessionID string) error

	// Session state management
	StartSession(sessionID string) error
	CompleteSession(sessionID string) error
	FailSession(sessionID string, err error) error
	CancelSession(sessionID string) error

	// Session cleanup
	CleanupExpiredSessions(maxAge time.Duration) error
	CleanupAllSessions() error
}

// ActingResultProcessor defines the interface for processing acting results
type ActingResultProcessor interface {
	// Result processing
	ProcessResult(result *spookyactionstypes.ActingResult) error
	ProcessResults(results []*spookyactionstypes.ActingResult) error

	// Result aggregation
	AggregateResults(results []*spookyactionstypes.ActingResult) (*spookyactionstypes.ActingSession, error)
	CalculateSuccessRate(results []*spookyactionstypes.ActingResult) float64

	// Result validation
	ValidateResult(result *spookyactionstypes.ActingResult) error
	ValidateResults(results []*spookyactionstypes.ActingResult) error

	// Result transformation
	TransformResult(result *spookyactionstypes.ActingResult, format string) (interface{}, error)
	TransformResults(results []*spookyactionstypes.ActingResult, format string) (interface{}, error)
}

// ActingProgressTracker defines the interface for tracking acting progress
type ActingProgressTracker interface {
	// Progress tracking
	StartTracking(sessionID string, totalSteps int) error
	UpdateProgress(sessionID string, completedSteps int) error
	GetProgress(sessionID string) (float64, error)
	CompleteTracking(sessionID string) error

	// Progress reporting
	GetProgressReport(sessionID string) (*ProgressReport, error)
	ListProgressReports() ([]*ProgressReport, error)

	// Progress cleanup
	CleanupProgress(sessionID string) error
	CleanupAllProgress() error
}

// ProgressReport represents a progress report for acting operations
type ProgressReport struct {
	SessionID      string    `json:"session_id"`
	ActionName     string    `json:"action_name"`
	TotalSteps     int       `json:"total_steps"`
	CompletedSteps int       `json:"completed_steps"`
	Progress       float64   `json:"progress"`
	StartTime      time.Time `json:"start_time"`
	LastUpdate     time.Time `json:"last_update"`
	EstimatedEnd   time.Time `json:"estimated_end,omitempty"`
	Status         string    `json:"status"`
}
