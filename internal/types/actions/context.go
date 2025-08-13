package spookytypesactions

import (
	"context"
	"time"

	spookytypescommon "spooky/internal/types/common"
	spookytypesmachines "spooky/internal/types/machines"
)

// ActionRunContext represents the context for action running
type ActionRunContext struct {
	spookytypescommon.CompleteEntity

	// Context identification
	ContextID string    `json:"context_id" hcl:"context_id"`
	SessionID string    `json:"session_id" hcl:"session_id"`
	CreatedAt time.Time `json:"created_at" hcl:"created_at"`

	// Project context
	ProjectPath string `json:"project_path" hcl:"project_path"`
	ProjectName string `json:"project_name" hcl:"project_name"`

	// Action context
	ActionName string `json:"action_name" hcl:"action_name"`
	ActionType string `json:"action_type" hcl:"action_type"`

	// Machine context
	MachineName string `json:"machine_name" hcl:"machine_name"`
	MachineHost string `json:"machine_host" hcl:"machine_host"`

	// Running context
	UserID      string            `json:"user_id" hcl:"user_id"`
	Environment map[string]string `json:"environment" hcl:"environment,optional"`
	Variables   map[string]string `json:"variables" hcl:"variables,optional"`

	// Configuration
	Timeout    time.Duration `json:"timeout" hcl:"timeout,optional"`
	RetryCount int           `json:"retry_count" hcl:"retry_count,optional"`
	MaxRetries int           `json:"max_retries" hcl:"max_retries,optional"`
	Parallel   bool          `json:"parallel" hcl:"parallel,optional"`

	// State
	Status    string     `json:"status" hcl:"status"` // "pending", "running", "completed", "failed", "cancelled"
	StartTime *time.Time `json:"start_time" hcl:"start_time,optional"`
	EndTime   *time.Time `json:"end_time" hcl:"end_time,optional"`

	// Results
	ExitCode int    `json:"exit_code" hcl:"exit_code,optional"`
	Stdout   string `json:"stdout" hcl:"stdout,optional"`
	Stderr   string `json:"stderr" hcl:"stderr,optional"`
	Error    string `json:"error" hcl:"error,optional"`

	// Metadata
	Tags     []string          `json:"tags" hcl:"tags,optional"`
	Metadata map[string]string `json:"metadata" hcl:"metadata,optional"`

	// Go context
	Context context.Context `json:"-" hcl:"-"`
}

// ActingSession represents a session for running actions
type ActingSession struct {
	spookytypescommon.CompleteEntity

	// Session identification
	SessionID string    `json:"session_id" hcl:"session_id"`
	CreatedAt time.Time `json:"created_at" hcl:"created_at"`
	ExpiresAt time.Time `json:"expires_at" hcl:"expires_at"`

	// Session context
	UserID      string `json:"user_id" hcl:"user_id"`
	ProjectPath string `json:"project_path" hcl:"project_path"`
	ProjectName string `json:"project_name" hcl:"project_name"`

	// Session state
	Status    string     `json:"status" hcl:"status"` // "active", "completed", "failed", "cancelled"
	StartTime *time.Time `json:"start_time" hcl:"start_time,optional"`
	EndTime   *time.Time `json:"end_time" hcl:"end_time,optional"`

	// Session configuration
	Parallel      bool          `json:"parallel" hcl:"parallel,optional"`
	MaxConcurrent int           `json:"max_concurrent" hcl:"max_concurrent,optional"`
	Timeout       time.Duration `json:"timeout" hcl:"timeout,optional"`
	AllowFailures bool          `json:"allow_failures" hcl:"allow_failures,optional"`

	// Session results
	TotalActions     int     `json:"total_actions" hcl:"total_actions"`
	CompletedActions int     `json:"completed_actions" hcl:"completed_actions"`
	FailedActions    int     `json:"failed_actions" hcl:"failed_actions"`
	SuccessRate      float64 `json:"success_rate" hcl:"success_rate"`

	// Machine inventory access
	MachineInventory []spookytypesmachines.Machine           `json:"machine_inventory" hcl:"machine_inventory,optional"`
	MachineCache     map[string]*spookytypesmachines.Machine `json:"machine_cache" hcl:"machine_cache,optional"`

	// Session metadata
	Tags     []string          `json:"tags" hcl:"tags,optional"`
	Metadata map[string]string `json:"metadata" hcl:"metadata,optional"`

	// Go context
	Context context.Context `json:"-" hcl:"-"`
}

// ActingResult represents the result of running an action
type ActingResult struct {
	spookytypescommon.CompleteEntity

	// Result identification
	ResultID    string `json:"result_id" hcl:"result_id"`
	SessionID   string `json:"session_id" hcl:"session_id"`
	ActionName  string `json:"action_name" hcl:"action_name"`
	MachineName string `json:"machine_name" hcl:"machine_name"`

	// Result state
	Status    string        `json:"status" hcl:"status"` // "success", "failure", "cancelled"
	StartTime time.Time     `json:"start_time" hcl:"start_time"`
	EndTime   time.Time     `json:"end_time" hcl:"end_time"`
	Duration  time.Duration `json:"duration" hcl:"duration"`

	// Result data
	ExitCode  int    `json:"exit_code" hcl:"exit_code"`
	Stdout    string `json:"stdout" hcl:"stdout,optional"`
	Stderr    string `json:"stderr" hcl:"stderr,optional"`
	Error     string `json:"error" hcl:"error,optional"`
	ErrorType string `json:"error_type" hcl:"error_type,optional"`

	// Performance data
	CPUUsage    float64 `json:"cpu_usage" hcl:"cpu_usage,optional"`
	MemoryUsage int64   `json:"memory_usage" hcl:"memory_usage,optional"`
	DiskUsage   int64   `json:"disk_usage" hcl:"disk_usage,optional"`

	// Retry information
	RetryCount int `json:"retry_count" hcl:"retry_count"`
	MaxRetries int `json:"max_retries" hcl:"max_retries"`

	// Metadata
	Tags     []string          `json:"tags" hcl:"tags,optional"`
	Metadata map[string]string `json:"metadata" hcl:"metadata,optional"`
}
