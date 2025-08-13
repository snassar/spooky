package spookytypesactions

import (
	spookytypescommon "spooky/internal/types/common"
)

// Action represents a single action definition
type Action struct {
	spookytypescommon.CompleteEntity

	// Basic action properties
	Name        string `hcl:"name,label" json:"name"`
	Description string `hcl:"description" json:"description"`
	Type        string `hcl:"type" json:"type"`

	// Action execution configuration
	Command        string            `hcl:"command,optional" json:"command,omitempty"`
	Script         string            `hcl:"script,optional" json:"script,omitempty"`
	Variables      map[string]string `hcl:"variables,optional" json:"variables,omitempty"`
	Template       *TemplateConfig   `hcl:"template,block" json:"template,omitempty"`
	FileCopy       *FileCopyConfig   `hcl:"file_copy,block" json:"file_copy,omitempty"`
	ServiceControl *ServiceConfig    `hcl:"service_control,block" json:"service_control,omitempty"`

	// Machine targeting
	Machines []string `hcl:"machines,optional" json:"machines,omitempty"`
	Tags     []string `hcl:"tags,optional" json:"tags,omitempty"`

	// Execution settings
	Timeout          int               `hcl:"timeout,optional" json:"timeout,omitempty"`
	Parallel         bool              `hcl:"parallel,optional" json:"parallel,omitempty"`
	Retries          int               `hcl:"retries,optional" json:"retries,omitempty"`
	RetryDelay       int               `hcl:"retry_delay,optional" json:"retry_delay,omitempty"`
	Dependencies     []string          `hcl:"dependencies,optional" json:"dependencies,omitempty"`
	Environment      map[string]string `hcl:"environment,optional" json:"environment,omitempty"`
	WorkingDirectory string            `hcl:"working_directory,optional" json:"working_directory,omitempty"`
	User             string            `hcl:"user,optional" json:"user,omitempty"`
	Sudo             bool              `hcl:"sudo,optional" json:"sudo,omitempty"`
	DryRun           bool              `hcl:"dry_run,optional" json:"dry_run,omitempty"`

	// Metadata and organization
	Category string            `hcl:"category,optional" json:"category,omitempty"`
	Priority int               `hcl:"priority,optional" json:"priority,omitempty"`
	Critical bool              `hcl:"critical,optional" json:"critical,omitempty"`
	Metadata map[string]string `hcl:"metadata,optional" json:"metadata,omitempty"`

	// Security and validation
	ValidateBeforeRun bool `hcl:"validate_before_run,optional" json:"validate_before_run,omitempty"`
	AllowFailure      bool `hcl:"allow_failure,optional" json:"allow_failure,omitempty"`

	// Performance and resource limits
	MaxConcurrent  int             `hcl:"max_concurrent,optional" json:"max_concurrent,omitempty"`
	ResourceLimits *ResourceLimits `hcl:"resource_limits,block" json:"resource_limits,omitempty"`
}

// TemplateConfig represents template deployment configuration
type TemplateConfig struct {
	Source      string `hcl:"source" json:"source"`
	Destination string `hcl:"destination" json:"destination"`
	Validate    bool   `hcl:"validate,optional" json:"validate,omitempty"`
	Backup      bool   `hcl:"backup,optional" json:"backup,omitempty"`
	Permissions string `hcl:"permissions,optional" json:"permissions,omitempty"`
	Owner       string `hcl:"owner,optional" json:"owner,omitempty"`
	Group       string `hcl:"group,optional" json:"group,omitempty"`
}

// FileCopyConfig represents file copy configuration
type FileCopyConfig struct {
	Source      string `hcl:"source" json:"source"`
	Destination string `hcl:"destination" json:"destination"`
	Backup      bool   `hcl:"backup,optional" json:"backup,omitempty"`
	Permissions string `hcl:"permissions,optional" json:"permissions,omitempty"`
	Owner       string `hcl:"owner,optional" json:"owner,omitempty"`
	Group       string `hcl:"group,optional" json:"group,omitempty"`
}

// ServiceConfig represents service control configuration
type ServiceConfig struct {
	Service       string `hcl:"service" json:"service"`
	Action        string `hcl:"action" json:"action"`
	Systemd       bool   `hcl:"systemd,optional" json:"systemd,omitempty"`
	Timeout       int    `hcl:"timeout,optional" json:"timeout,omitempty"`
	WaitForStatus string `hcl:"wait_for_status,optional" json:"wait_for_status,omitempty"`
	WaitTimeout   int    `hcl:"wait_timeout,optional" json:"wait_timeout,omitempty"`
}

// ResourceLimits represents resource limits for action execution
type ResourceLimits struct {
	MemoryMB   int `hcl:"memory_mb,optional" json:"memory_mb,omitempty"`
	CPUPercent int `hcl:"cpu_percent,optional" json:"cpu_percent,omitempty"`
	DiskMB     int `hcl:"disk_mb,optional" json:"disk_mb,omitempty"`
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
