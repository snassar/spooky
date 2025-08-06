package project

import (
	"time"
)

// Project represents a spooky project with metadata, configuration, and structure
type Project struct {
	// Project metadata
	Name        string          `hcl:"name"`
	Description string          `hcl:"description,optional"`
	Version     string          `hcl:"version,optional"`
	Author      string          `hcl:"author,optional"`
	Contact     *ProjectContact `hcl:"contact,block,optional"`

	// Project configuration
	Environment string   `hcl:"environment,optional"`
	Region      string   `hcl:"region,optional"`
	Tags        []string `hcl:"tags,optional"`

	// Project structure
	Structure *ProjectStructure `hcl:"structure,block,optional"`

	// Project isolation and security
	Isolation *ProjectIsolation `hcl:"isolation,block,optional"`

	// Project dependencies
	Dependencies *ProjectDependencies `hcl:"dependencies,block,optional"`

	// Project execution settings
	Execution *ProjectExecution `hcl:"execution,block,optional"`

	// Project metadata
	Metadata map[string]string `hcl:"metadata,optional"`

	// Internal fields
	Path      string    `hcl:"-"`
	CreatedAt time.Time `hcl:"-"`
	UpdatedAt time.Time `hcl:"-"`
}

// ProjectContact represents project contact information
type ProjectContact struct {
	Email string `hcl:"email,optional"`
	URL   string `hcl:"url,optional"`
}

// ProjectStructure represents project directory structure configuration
type ProjectStructure struct {
	TemplatesDir string `hcl:"templates_dir,optional"`
	DataDir      string `hcl:"data_dir,optional"`
	ScriptsDir   string `hcl:"scripts_dir,optional"`
	LogsDir      string `hcl:"logs_dir,optional"`
	BackupsDir   string `hcl:"backups_dir,optional"`
}

// ProjectIsolation represents project isolation and security settings
type ProjectIsolation struct {
	Enabled         bool     `hcl:"enabled,optional"`
	FactsScope      string   `hcl:"facts_scope,optional"`
	VariablesScope  string   `hcl:"variables_scope,optional"`
	MachineAccess   string   `hcl:"machine_access,optional"`
	AllowedMachines []string `hcl:"allowed_machines,optional"`
	AllowedTags     []string `hcl:"allowed_tags,optional"`
}

// ProjectDependencies represents project dependencies and imports
type ProjectDependencies struct {
	Imports         []string `hcl:"imports,optional"`
	SharedVariables []string `hcl:"shared_variables,optional"`
	SharedFacts     []string `hcl:"shared_facts,optional"`
}

// ProjectExecution represents project execution configuration
type ProjectExecution struct {
	DefaultTimeout        int  `hcl:"default_timeout,optional"`
	MaxParallel           int  `hcl:"max_parallel,optional"`
	DryRunDefault         bool `hcl:"dry_run_default,optional"`
	ValidateBeforeExecute bool `hcl:"validate_before_execute,optional"`
	BackupBeforeChanges   bool `hcl:"backup_before_changes,optional"`
}

// ProjectMetadata represents project metadata for filtering and organization
type ProjectMetadata struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Version     string            `json:"version,omitempty"`
	Environment string            `json:"environment,omitempty"`
	Region      string            `json:"region,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Path        string            `json:"path"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// ProjectConfiguration represents project-specific configuration
type ProjectConfiguration struct {
	// Project-specific configuration overrides
	FactsStorage  string `hcl:"facts_storage,optional"`
	VariablesPath string `hcl:"variables_path,optional"`
	MachinesPath  string `hcl:"machines_path,optional"`
	ActionsPath   string `hcl:"actions_path,optional"`
	TemplatesPath string `hcl:"templates_path,optional"`

	// Project-specific settings
	Settings map[string]interface{} `hcl:"settings,optional"`
}

// ProjectWrapper wraps Project for HCL parsing
type ProjectWrapper struct {
	Project *Project `hcl:"project,block"`
}

// Default values for project configuration
const (
	DefaultProjectEnvironment = "development"
	DefaultTimeout            = 300
	DefaultMaxParallel        = 10
	DefaultTemplatesDir       = "templates"
	DefaultDataDir            = "data"
	DefaultScriptsDir         = "scripts"
	DefaultLogsDir            = "logs"
	DefaultBackupsDir         = "backups"
	DefaultFactsScope         = "project"
	DefaultVariablesScope     = "project"
	DefaultMachineAccess      = "all"
)

// ProjectInfo represents comprehensive project information
type ProjectInfo struct {
	Project  *Project              `json:"project"`
	Metadata *ProjectMetadata      `json:"metadata"`
	Config   *ProjectConfiguration `json:"config,omitempty"`
	Stats    *ProjectStats         `json:"stats,omitempty"`
	Status   *ProjectStatus        `json:"status,omitempty"`
}

// ProjectStats represents project statistics
type ProjectStats struct {
	FactsCount     int       `json:"facts_count"`
	VariablesCount int       `json:"variables_count"`
	MachinesCount  int       `json:"machines_count"`
	ActionsCount   int       `json:"actions_count"`
	TemplatesCount int       `json:"templates_count"`
	LastUpdated    time.Time `json:"last_updated"`
}

// ProjectStatus represents project validation status
type ProjectStatus struct {
	Valid         bool      `json:"valid"`
	Errors        []string  `json:"errors,omitempty"`
	Warnings      []string  `json:"warnings,omitempty"`
	LastValidated time.Time `json:"last_validated,omitempty"`
}

// SetDefaults sets default values for the project
func (p *Project) SetDefaults() {
	// Set default environment if not specified
	if p.Environment == "" {
		p.Environment = "development"
	}

	// Set default structure if not specified
	if p.Structure == nil {
		p.Structure = &ProjectStructure{
			TemplatesDir: "templates",
			DataDir:      "data",
			ScriptsDir:   "scripts",
			LogsDir:      "logs",
			BackupsDir:   "backups",
		}
	}

	// Set default isolation if not specified
	if p.Isolation == nil {
		p.Isolation = &ProjectIsolation{
			Enabled:         true,
			FactsScope:      "project",
			VariablesScope:  "project",
			MachineAccess:   "all",
			AllowedMachines: []string{},
			AllowedTags:     []string{},
		}
	} else {
		// Set defaults for missing isolation fields while preserving existing values
		if p.Isolation.FactsScope == "" {
			p.Isolation.FactsScope = "project"
		}
		if p.Isolation.VariablesScope == "" {
			p.Isolation.VariablesScope = "project"
		}
		if p.Isolation.MachineAccess == "" {
			p.Isolation.MachineAccess = "all"
		}
		// Initialize arrays if they're nil (but don't override existing arrays)
		if p.Isolation.AllowedMachines == nil {
			p.Isolation.AllowedMachines = []string{}
		}
		if p.Isolation.AllowedTags == nil {
			p.Isolation.AllowedTags = []string{}
		}
	}

	// Set default execution if not specified
	if p.Execution == nil {
		p.Execution = &ProjectExecution{
			DefaultTimeout:        300,
			MaxParallel:           10,
			DryRunDefault:         false,
			ValidateBeforeExecute: true,
			BackupBeforeChanges:   false,
		}
	}

	// Set default metadata if not specified
	if p.Metadata == nil {
		p.Metadata = make(map[string]string)
	}

	// Set timestamps if not set
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	p.UpdatedAt = time.Now()
}

// GetMetadata returns project metadata
func (p *Project) GetMetadata() *ProjectMetadata {
	return &ProjectMetadata{
		Name:        p.Name,
		Description: p.Description,
		Version:     p.Version,
		Environment: p.Environment,
		Region:      p.Region,
		Tags:        p.Tags,
		Path:        p.Path,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Metadata:    p.Metadata,
	}
}

// GetInfo returns comprehensive project information
func (p *Project) GetInfo() *ProjectInfo {
	return &ProjectInfo{
		Project:  p,
		Metadata: p.GetMetadata(),
		Config:   &ProjectConfiguration{},
		Stats:    &ProjectStats{},
		Status:   &ProjectStatus{},
	}
}

// IsValid returns true if the project has a valid name
func (p *Project) IsValid() bool {
	return p.Name != "" && p.Path != ""
}

// HasTag returns true if the project has the specified tag
func (p *Project) HasTag(tag string) bool {
	for _, t := range p.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// HasAnyTag returns true if the project has any of the specified tags
func (p *Project) HasAnyTag(tags []string) bool {
	for _, tag := range tags {
		if p.HasTag(tag) {
			return true
		}
	}
	return false
}

// HasAllTags returns true if the project has all of the specified tags
func (p *Project) HasAllTags(tags []string) bool {
	for _, tag := range tags {
		if !p.HasTag(tag) {
			return false
		}
	}
	return true
}

// MatchesEnvironment returns true if the project matches the specified environment
func (p *Project) MatchesEnvironment(environment string) bool {
	return p.Environment == environment
}

// MatchesRegion returns true if the project matches the specified region
func (p *Project) MatchesRegion(region string) bool {
	return p.Region == region
}
