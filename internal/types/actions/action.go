package spookytypesactions

import (
	spookytypescommon "spooky/internal/types/common"
)

// Action represents a single action to be run on machines
type Action struct {
	spookytypescommon.CompleteEntity

	// Action identification
	Name        string `json:"name" hcl:"name"`
	Description string `json:"description" hcl:"description,optional"`
	Type        string `json:"type" hcl:"type"` // "command", "script", "template_deploy", "file_copy", "service_control"

	// Action targeting
	Machines []string `json:"machines" hcl:"machines,optional"`
	Tags     []string `json:"tags" hcl:"tags,optional"`

	// Action dependencies
	Dependencies []string `json:"dependencies" hcl:"dependencies,optional"`

	// Action run configuration
	Parallel      bool `json:"parallel" hcl:"parallel,optional"`
	MaxConcurrent int  `json:"max_concurrent" hcl:"max_concurrent,optional"`
	Timeout       int  `json:"timeout" hcl:"timeout,optional"`
	Retries       int  `json:"retries" hcl:"retries,optional"`
	RetryDelay    int  `json:"retry_delay" hcl:"retry_delay,optional"`
	AllowFailure  bool `json:"allow_failure" hcl:"allow_failure,optional"`
	StopOnFailure bool `json:"stop_on_failure" hcl:"stop_on_failure,optional"`

	// Action configuration based on type
	CommandString  string                `json:"command_string" hcl:"command,optional"`
	Command        *CommandConfig        `json:"command" hcl:"command_config,optional"`
	Script         *ScriptConfig         `json:"script" hcl:"script,optional"`
	Template       *TemplateConfig       `json:"template" hcl:"template,optional"`
	FileCopy       *FileCopyConfig       `json:"file_copy" hcl:"file_copy,optional"`
	ServiceControl *ServiceControlConfig `json:"service_control" hcl:"service_control,optional"`

	// Resource limits
	ResourceLimits *ResourceLimits `json:"resource_limits" hcl:"resource_limits,optional"`

	// Environment and variables
	Environment map[string]string `json:"environment" hcl:"environment,optional"`
	Variables   map[string]string `json:"variables" hcl:"variables,optional"`
	WorkingDir  string            `json:"working_dir" hcl:"working_dir,optional"`

	// Security and permissions
	Sudo  bool   `json:"sudo" hcl:"sudo,optional"`
	User  string `json:"user" hcl:"user,optional"`
	Group string `json:"group" hcl:"group,optional"`

	// Validation and safety
	ValidateBefore bool `json:"validate_before" hcl:"validate_before,optional"`
	DryRun         bool `json:"dry_run" hcl:"dry_run,optional"`

	// Metadata
	Metadata map[string]string `json:"metadata" hcl:"metadata,optional"`
}

// CommandConfig represents configuration for command actions
type CommandConfig struct {
	Command     string            `json:"command" hcl:"command"`
	Args        []string          `json:"args" hcl:"args,optional"`
	Environment map[string]string `json:"environment" hcl:"environment,optional"`
	WorkingDir  string            `json:"working_dir" hcl:"working_dir,optional"`
	Shell       string            `json:"shell" hcl:"shell,optional"`
	Timeout     int               `json:"timeout" hcl:"timeout,optional"`
}

// ScriptConfig represents configuration for script actions
type ScriptConfig struct {
	Script      string            `json:"script" hcl:"script"`
	Args        []string          `json:"args" hcl:"args,optional"`
	Environment map[string]string `json:"environment" hcl:"environment,optional"`
	WorkingDir  string            `json:"working_dir" hcl:"working_dir,optional"`
	Shell       string            `json:"shell" hcl:"shell,optional"`
	Timeout     int               `json:"timeout" hcl:"timeout,optional"`
	Validate    bool              `json:"validate" hcl:"validate,optional"`
}

// TemplateConfig represents configuration for template deployment actions
type TemplateConfig struct {
	Source       string            `json:"source" hcl:"source"`
	Destination  string            `json:"destination" hcl:"destination"`
	Permissions  string            `json:"permissions" hcl:"permissions,optional"`
	Owner        string            `json:"owner" hcl:"owner,optional"`
	Group        string            `json:"group" hcl:"group,optional"`
	Backup       bool              `json:"backup" hcl:"backup,optional"`
	BackupSuffix string            `json:"backup_suffix" hcl:"backup_suffix,optional"`
	Variables    map[string]string `json:"variables" hcl:"variables,optional"`
	Validate     bool              `json:"validate" hcl:"validate,optional"`
	Cleanup      bool              `json:"cleanup" hcl:"cleanup,optional"`
}

// FileCopyConfig represents configuration for file copy actions
type FileCopyConfig struct {
	Source      string `json:"source" hcl:"source"`
	Destination string `json:"destination" hcl:"destination"`
	Permissions string `json:"permissions" hcl:"permissions,optional"`
	Owner       string `json:"owner" hcl:"owner,optional"`
	Group       string `json:"group" hcl:"group,optional"`
	Backup      bool   `json:"backup" hcl:"backup,optional"`
	Recursive   bool   `json:"recursive" hcl:"recursive,optional"`
	Preserve    bool   `json:"preserve" hcl:"preserve,optional"`
}

// ServiceControlConfig represents configuration for service control actions
type ServiceControlConfig struct {
	Service       string `json:"service" hcl:"service"`
	Action        string `json:"action" hcl:"action"` // "start", "stop", "restart", "reload", "enable", "disable", "status"
	Systemd       bool   `json:"systemd" hcl:"systemd,optional"`
	Timeout       int    `json:"timeout" hcl:"timeout,optional"`
	WaitForStatus string `json:"wait_for_status" hcl:"wait_for_status,optional"`
	WaitTimeout   int    `json:"wait_timeout" hcl:"wait_timeout,optional"`
	CheckInterval int    `json:"check_interval" hcl:"check_interval,optional"`
	MaxRetries    int    `json:"max_retries" hcl:"max_retries,optional"`
	RetryDelay    int    `json:"retry_delay" hcl:"retry_delay,optional"`
}

// ResourceLimits represents resource limits for action running
type ResourceLimits struct {
	MemoryMB      int     `json:"memory_mb" hcl:"memory_mb,optional"`
	CPUPercent    float64 `json:"cpu_percent" hcl:"cpu_percent,optional"`
	DiskMB        int     `json:"disk_mb" hcl:"disk_mb,optional"`
	NetworkMB     int     `json:"network_mb" hcl:"network_mb,optional"`
	ProcessCount  int     `json:"process_count" hcl:"process_count,optional"`
	OpenFiles     int     `json:"open_files" hcl:"open_files,optional"`
	MaxMemoryMB   int     `json:"max_memory_mb" hcl:"max_memory_mb,optional"`
	MaxCPUPercent float64 `json:"max_cpu_percent" hcl:"max_cpu_percent,optional"`
}

// ActionType represents the type of action
type ActionType string

const (
	ActionTypeCommand          ActionType = "command"
	ActionTypeScript           ActionType = "script"
	ActionTypeTemplateDeploy   ActionType = "template_deploy"
	ActionTypeTemplateEvaluate ActionType = "template_evaluate"
	ActionTypeTemplateValidate ActionType = "template_validate"
	ActionTypeTemplateCleanup  ActionType = "template_cleanup"
	ActionTypeFileCopy         ActionType = "file_copy"
	ActionTypeServiceControl   ActionType = "service_control"
)

// ServiceAction represents the type of service action
type ServiceAction string

const (
	ServiceActionStart   ServiceAction = "start"
	ServiceActionStop    ServiceAction = "stop"
	ServiceActionRestart ServiceAction = "restart"
	ServiceActionReload  ServiceAction = "reload"
	ServiceActionEnable  ServiceAction = "enable"
	ServiceActionDisable ServiceAction = "disable"
	ServiceActionStatus  ServiceAction = "status"
)

// ServiceStatus represents the status of a service
type ServiceStatus string

const (
	ServiceStatusActive   ServiceStatus = "active"
	ServiceStatusInactive ServiceStatus = "inactive"
	ServiceStatusFailed   ServiceStatus = "failed"
	ServiceStatusAny      ServiceStatus = "any"
)
