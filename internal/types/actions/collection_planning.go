package actions

import (
	"time"
)

// CollectionPlan represents a planned execution of action collections
type CollectionPlan struct {
	ID               string                  `json:"id"`
	Name             string                  `json:"name"`
	Description      string                  `json:"description"`
	Collections      []*PlannedCollection    `json:"collections"`
	Dependencies     []*CollectionDependency `json:"dependencies"`
	ExecutionOrder   []string                `json:"execution_order"`
	ParallelGroups   [][]string              `json:"parallel_groups"`
	EstimatedTime    time.Duration           `json:"estimated_time"`
	RollbackPlan     *RollbackPlan           `json:"rollback_plan"`
	Validation       *CollectionValidation   `json:"validation"`
	Metadata         map[string]interface{}  `json:"metadata"`
	FactRequirements *FactRequirements       `json:"fact_requirements"`
}

// FactRequirements defines what facts are needed for collection planning
type FactRequirements struct {
	SystemFacts     []string              `json:"system_facts"`     // Required system facts
	CustomFacts     []string              `json:"custom_facts"`     // Required custom facts
	FactSources     []FactSource          `json:"fact_sources"`     // Required fact sources
	ValidationRules []*FactValidationRule `json:"validation_rules"` // Fact validation rules
}

// FactValidationRule defines a rule for validating facts
type FactValidationRule struct {
	FactKey       string        `json:"fact_key"`
	Required      bool          `json:"required"`
	Type          string        `json:"type"`           // string, int, bool, etc.
	Pattern       string        `json:"pattern"`        // regex pattern for validation
	MinValue      interface{}   `json:"min_value"`      // minimum value for numeric types
	MaxValue      interface{}   `json:"max_value"`      // maximum value for numeric types
	AllowedValues []interface{} `json:"allowed_values"` // allowed values for enum types
}

// PlannedCollection represents a collection in the execution plan
type PlannedCollection struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	Actions          []*PlannedAction `json:"actions"`
	Dependencies     []string         `json:"dependencies"`
	Dependents       []string         `json:"dependents"`
	ExecutionPhase   int              `json:"execution_phase"`
	ParallelGroup    int              `json:"parallel_group"`
	EstimatedTime    time.Duration    `json:"estimated_time"`
	Priority         int              `json:"priority"`
	CriticalPath     bool             `json:"critical_path"`
	RollbackActions  []*PlannedAction `json:"rollback_actions"`
	FactDependencies []string         `json:"fact_dependencies"` // Facts this collection depends on
}

// PlannedAction represents an action in the execution plan
type PlannedAction struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	Type             string                 `json:"type"`
	Machines         []string               `json:"machines"`
	Dependencies     []string               `json:"dependencies"`
	EstimatedTime    time.Duration          `json:"estimated_time"`
	Priority         int                    `json:"priority"`
	RetryPolicy      *RetryPolicy           `json:"retry_policy"`
	RollbackAction   *PlannedAction         `json:"rollback_action"`
	TemplateVars     map[string]interface{} `json:"template_vars"`     // Template variables from facts
	FactRequirements []string               `json:"fact_requirements"` // Facts this action needs
}

// CollectionDependency represents a dependency between collections
type CollectionDependency struct {
	FromCollection string         `json:"from_collection"`
	ToCollection   string         `json:"to_collection"`
	Type           DependencyType `json:"type"`
	Condition      string         `json:"condition"`
	Required       bool           `json:"required"`
	FactBased      bool           `json:"fact_based"`     // Whether dependency is based on facts
	FactCondition  string         `json:"fact_condition"` // Fact-based condition expression
}

// DependencyType represents the type of dependency
type DependencyType string

const (
	DependencyTypeBefore      DependencyType = "before"
	DependencyTypeAfter       DependencyType = "after"
	DependencyTypeConditional DependencyType = "conditional"
	DependencyTypeParallel    DependencyType = "parallel"
	DependencyTypeFactBased   DependencyType = "fact_based"
)

// RollbackPlan represents a plan for rolling back failed collections
type RollbackPlan struct {
	RollbackOrder   []string                    `json:"rollback_order"`
	RollbackActions map[string][]*PlannedAction `json:"rollback_actions"`
	Dependencies    []*CollectionDependency     `json:"dependencies"`
	EstimatedTime   time.Duration               `json:"estimated_time"`
}

// CollectionValidation represents validation results for collections
type CollectionValidation struct {
	Valid          bool                  `json:"valid"`
	Errors         []string              `json:"errors"`
	Warnings       []string              `json:"warnings"`
	Dependencies   []string              `json:"dependencies"`
	CircularDeps   []string              `json:"circular_dependencies"`
	MissingActions []string              `json:"missing_actions"`
	InvalidConfigs []string              `json:"invalid_configs"`
	FactValidation *FactValidationResult `json:"fact_validation"`
}

// FactValidationResult represents the result of fact validation
type FactValidationResult struct {
	Valid        bool                  `json:"valid"`
	MissingFacts []string              `json:"missing_facts"`
	InvalidFacts []string              `json:"invalid_facts"`
	FactSources  map[string]FactSource `json:"fact_sources"`
}

// RetryPolicy represents retry configuration for actions
type RetryPolicy struct {
	MaxAttempts     int           `json:"max_attempts"`
	Delay           time.Duration `json:"delay"`
	BackoffFactor   float64       `json:"backoff_factor"`
	MaxDelay        time.Duration `json:"max_delay"`
	RetryableErrors []string      `json:"retryable_errors"`
}

// PlanningOptions represents options for collection planning
type PlanningOptions struct {
	PlanName       string                 `json:"plan_name"`
	Description    string                 `json:"description"`
	Parallel       bool                   `json:"parallel"`
	Optimize       bool                   `json:"optimize"`
	MaxParallel    int                    `json:"max_parallel"`
	Timeout        time.Duration          `json:"timeout"`
	ValidateFacts  bool                   `json:"validate_facts"`
	CreateRollback bool                   `json:"create_rollback"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// FactSource represents the source of facts
type FactSource string

const (
	SourceSystem FactSource = "system"
	SourceHCL    FactSource = "hcl"
	SourceJSON   FactSource = "json"
	SourceSSH    FactSource = "ssh"
	SourceLocal  FactSource = "local"
	SourceCustom FactSource = "custom"
)

// ActionFactDependencies represents fact dependencies for an action
type ActionFactDependencies struct {
	SystemFacts     []string              `json:"system_facts"`
	CustomFacts     []string              `json:"custom_facts"`
	Sources         []FactSource          `json:"sources"`
	ValidationRules []*FactValidationRule `json:"validation_rules"`
}

// DependencyNode represents a node in the dependency graph
type DependencyNode struct {
	ID               string
	Collection       *ActionCollection
	Dependencies     []string
	Dependents       []string
	Level            int
	Visited          bool
	InProgress       bool
	FactDependencies []string
}
