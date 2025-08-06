package types

import (
	spookyconfigtypes "spooky/internal/config/types"
	"time"
)

// Action represents an action to be executed on machines
type Action struct {
	// Basic action information
	Name        string `hcl:"name,label" validate:"required"`
	Description string `hcl:"description,optional"`
	Type        string `hcl:"type,optional" validate:"omitempty,oneof=command script template_deploy template_evaluate template_validate template_cleanup"`

	// Execution configuration
	Command  string          `hcl:"command,optional"`
	Script   string          `hcl:"script,optional"`
	Template *TemplateConfig `hcl:"template,block"`
	Machines []string        `hcl:"machines,optional" validate:"omitempty,dive,required"`
	Tags     []string        `hcl:"tags,optional" validate:"omitempty,dive,required"`
	Timeout  int             `hcl:"timeout,optional" validate:"omitempty,min=1,max=3600"`
	Parallel bool            `hcl:"parallel,optional"`

	// Extended properties for retries and error handling
	Retries    int `hcl:"retries,optional" validate:"omitempty,min=0,max=10"`
	RetryDelay int `hcl:"retry_delay,optional" validate:"omitempty,min=1,max=300"`

	// Dependencies and execution control
	Dependencies []string `hcl:"dependencies,optional" validate:"omitempty,dive,required"`

	// Environment and execution context
	Environment      map[string]string `hcl:"environment,optional" validate:"omitempty,dive,keys,required,endkeys,required"`
	WorkingDirectory string            `hcl:"working_directory,optional"`
	User             string            `hcl:"user,optional"`
	Sudo             bool              `hcl:"sudo,optional"`
	DryRun           bool              `hcl:"dry_run,optional"`

	// Action metadata and organization
	Category string            `hcl:"category,optional"`
	Priority int               `hcl:"priority,optional" validate:"omitempty,min=1,max=10"`
	Critical bool              `hcl:"critical,optional"`
	Metadata map[string]string `hcl:"metadata,optional" validate:"omitempty,dive,keys,required,endkeys,required"`

	// Security and validation
	ValidateBeforeRun bool `hcl:"validate_before_run,optional"`
	AllowFailure      bool `hcl:"allow_failure,optional"`

	// Performance and resource management
	MaxConcurrent  int                   `hcl:"max_concurrent,optional" validate:"omitempty,min=1,max=100"`
	ResourceLimits *ActionResourceLimits `hcl:"resource_limits,block"`

	// State tracking
	State *ActionState `hcl:"state,block"`

	// Source information
	SourceFile string `hcl:"source_file,optional"`
	SourceLine int    `hcl:"source_line,optional"`
}

// ActionState represents the current state of an action
type ActionState struct {
	Status         ActionStatus  `hcl:"status,optional"`
	StartedAt      *time.Time    `hcl:"started_at,optional"`
	CompletedAt    *time.Time    `hcl:"completed_at,optional"`
	Duration       time.Duration `hcl:"duration,optional"`
	ExecutionCount int           `hcl:"execution_count,optional"`
	SuccessCount   int           `hcl:"success_count,optional"`
	FailureCount   int           `hcl:"failure_count,optional"`
	LastError      error         `hcl:"last_error,optional"`
	Progress       float64       `hcl:"progress,optional"`
}

// ActionStatus represents the status of an action
type ActionStatus string

const (
	ActionStatusPending   ActionStatus = "pending"
	ActionStatusRunning   ActionStatus = "running"
	ActionStatusCompleted ActionStatus = "completed"
	ActionStatusFailed    ActionStatus = "failed"
	ActionStatusCancelled ActionStatus = "cancelled"
	ActionStatusSkipped   ActionStatus = "skipped"
)

// ActionCollection represents a collection of actions
type ActionCollection struct {
	Actions   []*Action              `hcl:"actions,block"`
	Metadata  map[string]interface{} `hcl:"metadata,optional"`
	CreatedAt time.Time              `hcl:"created_at,optional"`
	UpdatedAt time.Time              `hcl:"updated_at,optional"`
}

// TemplateConfig represents template configuration for actions
type TemplateConfig struct {
	Source      string `hcl:"source" validate:"required"`
	Destination string `hcl:"destination" validate:"required"`
	Validate    bool   `hcl:"validate,optional"`
	Backup      bool   `hcl:"backup,optional"`
	Permissions string `hcl:"permissions,optional" validate:"omitempty,regexp=^[0-7]{3,4}$"`
	Owner       string `hcl:"owner,optional"`
	Group       string `hcl:"group,optional"`
}

// ActionResourceLimits represents resource limits for actions
type ActionResourceLimits struct {
	MemoryMB   int `hcl:"memory_mb,optional" validate:"omitempty,min=1,max=32768"`
	CPUPercent int `hcl:"cpu_percent,optional" validate:"omitempty,min=1,max=100"`
	DiskMB     int `hcl:"disk_mb,optional" validate:"omitempty,min=1,max=1048576"`
}

// NewAction creates a new Action from a config Action
func NewAction(configAction *spookyconfigtypes.Action) *Action {
	if configAction == nil {
		return nil
	}

	action := &Action{
		Name:              configAction.Name,
		Description:       configAction.Description,
		Type:              configAction.Type,
		Command:           configAction.Command,
		Script:            configAction.Script,
		Machines:          configAction.Machines,
		Tags:              configAction.Tags,
		Timeout:           configAction.Timeout,
		Parallel:          configAction.Parallel,
		Retries:           configAction.Retries,
		RetryDelay:        configAction.RetryDelay,
		Dependencies:      configAction.Dependencies,
		Environment:       configAction.Environment,
		WorkingDirectory:  configAction.WorkingDirectory,
		User:              configAction.User,
		Sudo:              configAction.Sudo,
		DryRun:            configAction.DryRun,
		Category:          configAction.Category,
		Priority:          configAction.Priority,
		Critical:          configAction.Critical,
		Metadata:          configAction.Metadata,
		ValidateBeforeRun: configAction.ValidateBeforeRun,
		AllowFailure:      configAction.AllowFailure,
		MaxConcurrent:     configAction.MaxConcurrent,
	}

	// Convert template config if present
	if configAction.Template != nil {
		action.Template = &TemplateConfig{
			Source:      configAction.Template.Source,
			Destination: configAction.Template.Destination,
			Validate:    configAction.Template.Validate,
			Backup:      configAction.Template.Backup,
			Permissions: configAction.Template.Permissions,
			Owner:       configAction.Template.Owner,
			Group:       configAction.Template.Group,
		}
	}

	// Convert resource limits if present
	if configAction.ResourceLimits != nil {
		action.ResourceLimits = &ActionResourceLimits{
			MemoryMB:   configAction.ResourceLimits.MemoryMB,
			CPUPercent: configAction.ResourceLimits.CPUPercent,
			DiskMB:     configAction.ResourceLimits.DiskMB,
		}
	}

	return action
}

// ToConfigAction converts an Action to a config.Action
func (a *Action) ToConfigAction() *spookyconfigtypes.Action {
	if a == nil {
		return nil
	}

	configAction := &spookyconfigtypes.Action{
		Name:              a.Name,
		Description:       a.Description,
		Type:              a.Type,
		Command:           a.Command,
		Script:            a.Script,
		Machines:          a.Machines,
		Tags:              a.Tags,
		Timeout:           a.Timeout,
		Parallel:          a.Parallel,
		Retries:           a.Retries,
		RetryDelay:        a.RetryDelay,
		Dependencies:      a.Dependencies,
		Environment:       a.Environment,
		WorkingDirectory:  a.WorkingDirectory,
		User:              a.User,
		Sudo:              a.Sudo,
		DryRun:            a.DryRun,
		Category:          a.Category,
		Priority:          a.Priority,
		Critical:          a.Critical,
		Metadata:          a.Metadata,
		ValidateBeforeRun: a.ValidateBeforeRun,
		AllowFailure:      a.AllowFailure,
		MaxConcurrent:     a.MaxConcurrent,
	}

	// Convert template config if present
	if a.Template != nil {
		configAction.Template = &spookyconfigtypes.TemplateConfig{
			Source:      a.Template.Source,
			Destination: a.Template.Destination,
			Validate:    a.Template.Validate,
			Backup:      a.Template.Backup,
			Permissions: a.Template.Permissions,
			Owner:       a.Template.Owner,
			Group:       a.Template.Group,
		}
	}

	// Convert resource limits if present
	if a.ResourceLimits != nil {
		configAction.ResourceLimits = &spookyconfigtypes.ActionResourceLimits{
			MemoryMB:   a.ResourceLimits.MemoryMB,
			CPUPercent: a.ResourceLimits.CPUPercent,
			DiskMB:     a.ResourceLimits.DiskMB,
		}
	}

	return configAction
}
