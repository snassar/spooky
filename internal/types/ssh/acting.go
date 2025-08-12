// Package ssh provides SSH acting types for the spooky codebase.
// This package defines the data structures for SSH acting operations and session management.
package ssh

import (
	"time"

	spookytypescommon "spooky/internal/types/common"
)

// ActingSession represents an SSH acting session for remote operations
type ActingSession struct {
	spookytypescommon.CompleteEntity

	// Session details
	SessionID  string      `json:"session_id" hcl:"session_id"`
	Connection *Connection `json:"connection" hcl:"connection"`
	Client     *Client     `json:"client" hcl:"client"`

	// Session state
	Status    ActingSessionStatus `json:"status" hcl:"status"`
	StartedAt time.Time           `json:"started_at" hcl:"started_at"`
	EndedAt   *time.Time          `json:"ended_at,omitempty" hcl:"ended_at,optional"`

	// Session configuration
	Environment map[string]string `json:"environment,omitempty" hcl:"environment,optional"`
	WorkingDir  string            `json:"working_dir,omitempty" hcl:"working_dir,optional"`
	Pty         *PtyConfig        `json:"pty,omitempty" hcl:"pty,optional"`

	// Session metrics
	CommandsExecuted int           `json:"commands_executed" hcl:"commands_executed"`
	TotalExecTime    time.Duration `json:"total_exec_time" hcl:"total_exec_time"`
	AverageExecTime  time.Duration `json:"average_exec_time" hcl:"average_exec_time"`

	// Session metadata
	UserAgent    string `json:"user_agent,omitempty" hcl:"user_agent,optional"`
	TerminalType string `json:"terminal_type,omitempty" hcl:"terminal_type,optional"`

	// Acting session specific
	ActingMode ActingMode `json:"acting_mode" hcl:"acting_mode" default:"interactive"`
	BatchSize  int        `json:"batch_size" hcl:"batch_size" default:"1"`
	Parallel   bool       `json:"parallel" hcl:"parallel" default:"false"`
}

// ActingSessionStatus represents the status of an SSH acting session
type ActingSessionStatus string

const (
	ActingSessionStatusCreated   ActingSessionStatus = "created"
	ActingSessionStatusStarting  ActingSessionStatus = "starting"
	ActingSessionStatusActive    ActingSessionStatus = "active"
	ActingSessionStatusExecuting ActingSessionStatus = "executing"
	ActingSessionStatusCompleted ActingSessionStatus = "completed"
	ActingSessionStatusFailed    ActingSessionStatus = "failed"
	ActingSessionStatusClosed    ActingSessionStatus = "closed"
)

// ActingMode represents the mode of SSH acting
type ActingMode string

const (
	ActingModeInteractive ActingMode = "interactive"
	ActingModeBatch       ActingMode = "batch"
	ActingModeScript      ActingMode = "script"
	ActingModeCommand     ActingMode = "command"
)

// ActingCommand represents a command to be executed via SSH acting
type ActingCommand struct {
	spookytypescommon.CompleteEntity

	// Command details
	Command     string            `json:"command" hcl:"command"`
	Args        []string          `json:"args,omitempty" hcl:"args,optional"`
	WorkingDir  string            `json:"working_dir,omitempty" hcl:"working_dir,optional"`
	Environment map[string]string `json:"environment,omitempty" hcl:"environment,optional"`

	// Command configuration
	Timeout       time.Duration `json:"timeout,omitempty" hcl:"timeout,optional"`
	Pty           *PtyConfig    `json:"pty,omitempty" hcl:"pty,optional"`
	Stdin         string        `json:"stdin,omitempty" hcl:"stdin,optional"`
	CaptureOutput bool          `json:"capture_output" hcl:"capture_output" default:"true"`

	// Command metadata
	Priority    int       `json:"priority" hcl:"priority" default:"0"`
	ScheduledAt time.Time `json:"scheduled_at" hcl:"scheduled_at"`
	Tags        []string  `json:"tags,omitempty" hcl:"tags,optional"`

	// Security settings
	AllowUnsafe bool `json:"allow_unsafe" hcl:"allow_unsafe" default:"false"`
	RequireSudo bool `json:"require_sudo" hcl:"require_sudo" default:"false"`

	// Acting specific
	ActingMode ActingMode `json:"acting_mode" hcl:"acting_mode" default:"command"`
	RetryCount int        `json:"retry_count" hcl:"retry_count" default:"0"`
	MaxRetries int        `json:"max_retries" hcl:"max_retries" default:"3"`
}

// ActingCommandResult represents the result of a command execution via SSH acting
type ActingCommandResult struct {
	spookytypescommon.CompleteEntity

	// Command details
	Command *ActingCommand `json:"command" hcl:"command"`
	Session *ActingSession `json:"session" hcl:"session"`

	// Execution results
	Success  bool   `json:"success" hcl:"success"`
	ExitCode int    `json:"exit_code" hcl:"exit_code"`
	Stdout   string `json:"stdout,omitempty" hcl:"stdout,optional"`
	Stderr   string `json:"stderr,omitempty" hcl:"stderr,optional"`
	Error    string `json:"error,omitempty" hcl:"error,optional"`

	// Execution metrics
	StartTime     time.Time     `json:"start_time" hcl:"start_time"`
	EndTime       time.Time     `json:"end_time" hcl:"end_time"`
	Duration      time.Duration `json:"duration" hcl:"duration"`
	RetryAttempts int           `json:"retry_attempts" hcl:"retry_attempts"`

	// Execution metadata
	Hostname    string            `json:"hostname,omitempty" hcl:"hostname,optional"`
	Username    string            `json:"username,omitempty" hcl:"username,optional"`
	WorkingDir  string            `json:"working_dir,omitempty" hcl:"working_dir,optional"`
	Environment map[string]string `json:"environment,omitempty" hcl:"environment,optional"`

	// Security information
	CommandHash string `json:"command_hash,omitempty" hcl:"command_hash,optional"`
	AuditTrail  string `json:"audit_trail,omitempty" hcl:"audit_trail,optional"`

	// Acting specific
	ActingMode ActingMode `json:"acting_mode" hcl:"acting_mode"`
	BatchID    string     `json:"batch_id,omitempty" hcl:"batch_id,optional"`
	SequenceID int        `json:"sequence_id" hcl:"sequence_id"`
}

// ActingBatch represents a batch of commands to be executed via SSH acting
type ActingBatch struct {
	spookytypescommon.CompleteEntity

	// Batch details
	BatchID     string `json:"batch_id" hcl:"batch_id"`
	Name        string `json:"name" hcl:"name"`
	Description string `json:"description,omitempty" hcl:"description,optional"`

	// Batch commands
	Commands []*ActingCommand `json:"commands" hcl:"commands"`

	// Batch configuration
	Parallel    bool          `json:"parallel" hcl:"parallel" default:"false"`
	MaxParallel int           `json:"max_parallel" hcl:"max_parallel" default:"1"`
	Timeout     time.Duration `json:"timeout,omitempty" hcl:"timeout,optional"`

	// Batch metadata
	Priority    int       `json:"priority" hcl:"priority" default:"0"`
	ScheduledAt time.Time `json:"scheduled_at" hcl:"scheduled_at"`
	Tags        []string  `json:"tags,omitempty" hcl:"tags,optional"`

	// Batch state
	Status    ActingBatchStatus `json:"status" hcl:"status"`
	StartedAt *time.Time        `json:"started_at,omitempty" hcl:"started_at,optional"`
	EndedAt   *time.Time        `json:"ended_at,omitempty" hcl:"ended_at,optional"`

	// Batch metrics
	TotalCommands     int           `json:"total_commands" hcl:"total_commands"`
	CompletedCommands int           `json:"completed_commands" hcl:"completed_commands"`
	FailedCommands    int           `json:"failed_commands" hcl:"failed_commands"`
	TotalDuration     time.Duration `json:"total_duration" hcl:"total_duration"`
}

// ActingBatchStatus represents the status of an SSH acting batch
type ActingBatchStatus string

const (
	ActingBatchStatusPending   ActingBatchStatus = "pending"
	ActingBatchStatusStarting  ActingBatchStatus = "starting"
	ActingBatchStatusRunning   ActingBatchStatus = "running"
	ActingBatchStatusCompleted ActingBatchStatus = "completed"
	ActingBatchStatusFailed    ActingBatchStatus = "failed"
	ActingBatchStatusCancelled ActingBatchStatus = "cancelled"
)

// ActingBatchResult represents the result of a batch execution via SSH acting
type ActingBatchResult struct {
	spookytypescommon.CompleteEntity

	// Batch details
	Batch   *ActingBatch           `json:"batch" hcl:"batch"`
	Session *ActingSession         `json:"session" hcl:"session"`
	Results []*ActingCommandResult `json:"results" hcl:"results"`

	// Batch results
	Success bool   `json:"success" hcl:"success"`
	Error   string `json:"error,omitempty" hcl:"error,optional"`

	// Batch metrics
	StartTime time.Time     `json:"start_time" hcl:"start_time"`
	EndTime   time.Time     `json:"end_time" hcl:"end_time"`
	Duration  time.Duration `json:"duration" hcl:"duration"`

	// Batch statistics
	TotalCommands      int `json:"total_commands" hcl:"total_commands"`
	SuccessfulCommands int `json:"successful_commands" hcl:"successful_commands"`
	FailedCommands     int `json:"failed_commands" hcl:"failed_commands"`
	SkippedCommands    int `json:"skipped_commands" hcl:"skipped_commands"`

	// Batch metadata
	Hostname string `json:"hostname,omitempty" hcl:"hostname,optional"`
	Username string `json:"username,omitempty" hcl:"username,optional"`

	// Security information
	AuditTrail string `json:"audit_trail,omitempty" hcl:"audit_trail,optional"`
}

// ActingScript represents a script to be executed via SSH acting
type ActingScript struct {
	spookytypescommon.CompleteEntity

	// Script details
	ScriptID string `json:"script_id" hcl:"script_id"`
	Name     string `json:"name" hcl:"name"`
	Content  string `json:"content" hcl:"content"`
	Language string `json:"language" hcl:"language" default:"bash"`

	// Script configuration
	WorkingDir  string            `json:"working_dir,omitempty" hcl:"working_dir,optional"`
	Environment map[string]string `json:"environment,omitempty" hcl:"environment,optional"`
	Arguments   []string          `json:"arguments,omitempty" hcl:"arguments,optional"`

	// Script settings
	Timeout       time.Duration `json:"timeout,omitempty" hcl:"timeout,optional"`
	Pty           *PtyConfig    `json:"pty,omitempty" hcl:"pty,optional"`
	CaptureOutput bool          `json:"capture_output" hcl:"capture_output" default:"true"`

	// Script metadata
	Priority    int       `json:"priority" hcl:"priority" default:"0"`
	ScheduledAt time.Time `json:"scheduled_at" hcl:"scheduled_at"`
	Tags        []string  `json:"tags,omitempty" hcl:"tags,optional"`

	// Security settings
	AllowUnsafe bool `json:"allow_unsafe" hcl:"allow_unsafe" default:"false"`
	RequireSudo bool `json:"require_sudo" hcl:"require_sudo" default:"false"`

	// Script validation
	IsValid         bool   `json:"is_valid" hcl:"is_valid"`
	ValidationError string `json:"validation_error,omitempty" hcl:"validation_error,optional"`
}

// ActingScriptResult represents the result of a script execution via SSH acting
type ActingScriptResult struct {
	spookytypescommon.CompleteEntity

	// Script details
	Script  *ActingScript  `json:"script" hcl:"script"`
	Session *ActingSession `json:"session" hcl:"session"`

	// Execution results
	Success  bool   `json:"success" hcl:"success"`
	ExitCode int    `json:"exit_code" hcl:"exit_code"`
	Stdout   string `json:"stdout,omitempty" hcl:"stdout,optional"`
	Stderr   string `json:"stderr,omitempty" hcl:"stderr,optional"`
	Error    string `json:"error,omitempty" hcl:"error,optional"`

	// Execution metrics
	StartTime time.Time     `json:"start_time" hcl:"start_time"`
	EndTime   time.Time     `json:"end_time" hcl:"end_time"`
	Duration  time.Duration `json:"duration" hcl:"duration"`

	// Execution metadata
	Hostname    string            `json:"hostname,omitempty" hcl:"hostname,optional"`
	Username    string            `json:"username,omitempty" hcl:"username,optional"`
	WorkingDir  string            `json:"working_dir,omitempty" hcl:"working_dir,optional"`
	Environment map[string]string `json:"environment,omitempty" hcl:"environment,optional"`

	// Security information
	ScriptHash string `json:"script_hash,omitempty" hcl:"script_hash,optional"`
	AuditTrail string `json:"audit_trail,omitempty" hcl:"audit_trail,optional"`
}
