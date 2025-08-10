package ssh

import (
	"time"
)

// SSHAction represents an SSH action to act
type SSHAction struct {
	Name        string            `hcl:"name"`
	Command     string            `hcl:"command"`
	Script      string            `hcl:"script,optional"`
	Timeout     time.Duration     `hcl:"timeout,optional"`
	Environment map[string]string `hcl:"environment,optional"`
	WorkingDir  string            `hcl:"working_dir,optional"`
}

// ActionResult represents the result of an action acting
type ActionResult struct {
	Action    *SSHAction    `hcl:"action"`
	Success   bool          `hcl:"success"`
	ExitCode  int           `hcl:"exit_code"`
	Stdout    string        `hcl:"stdout"`
	Stderr    string        `hcl:"stderr"`
	Duration  time.Duration `hcl:"duration"`
	Error     string        `hcl:"error,optional"`
	Timestamp time.Time     `hcl:"timestamp"`
}

// TemplateAction represents a template-based action
type TemplateAction struct {
	Name       string                 `hcl:"name"`
	Template   string                 `hcl:"template"`
	Data       map[string]interface{} `hcl:"data,optional"`
	OutputFile string                 `hcl:"output_file,optional"`
	Timeout    time.Duration          `hcl:"timeout,optional"`
}

// ExecuteConfig represents acting configuration
type ExecuteConfig struct {
	Timeout     time.Duration `hcl:"timeout,optional"`
	Parallel    bool          `hcl:"parallel,optional"`
	MaxParallel int           `hcl:"max_parallel,optional"`
	Retries     int           `hcl:"retries,optional"`
	RetryDelay  time.Duration `hcl:"retry_delay,optional"`
}

// ActingConfig represents acting configuration
type ActingConfig struct {
	DefaultTimeout time.Duration `hcl:"default_timeout,optional"`
	MaxParallel    int           `hcl:"max_parallel,optional"`
	EnableRetries  bool          `hcl:"enable_retries,optional"`
	MaxRetries     int           `hcl:"max_retries,optional"`
	RetryDelay     time.Duration `hcl:"retry_delay,optional"`
}
