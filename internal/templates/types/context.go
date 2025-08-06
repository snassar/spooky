package types

import (
	"spooky/internal/config/types"
	variablesTypes "spooky/internal/variables/types"
	"time"
)

// TemplateContext holds all data available to templates
// Aligns with template-context.hcl schema
type TemplateContext struct {
	// Project configuration
	Project *types.ProjectConfig `hcl:"project,optional" json:"project"`

	// Machine facts (from facts.db or JSON)
	Facts map[string]interface{} `hcl:"facts,optional" json:"facts"`

	// Inventory information
	Machines []*types.Machine `hcl:"machines,optional" json:"machines"`

	// Actions configuration
	Actions []*types.Action `hcl:"actions,optional" json:"actions"`

	// Server-specific data (when --server is specified)
	ServerFacts map[string]interface{} `hcl:"server_facts,optional" json:"server_facts"`

	// Environment variables
	Environment map[string]string `hcl:"environment,optional" json:"environment"`

	// Variables context
	Variables *variablesTypes.VariableContext `hcl:"variables,optional" json:"variables"`

	// Additional fields for internal use
	ProjectPath string                 `hcl:"project_path,optional" json:"project_path"`
	Data        map[string]interface{} `hcl:"data,optional" json:"data"`
	Functions   map[string]interface{} `hcl:"functions,optional" json:"functions"`
	CreatedAt   time.Time              `hcl:"created_at,optional" json:"created_at"`
}

// TemplateData represents the data available in templates (legacy compatibility)
type TemplateData struct {
	// System facts
	System map[string]interface{} `hcl:"system,optional" json:"system"`

	// Custom facts
	Custom map[string]interface{} `hcl:"custom,optional" json:"custom"`

	// Environment variables
	Env map[string]string `hcl:"env,optional" json:"env"`

	// Additional data
	Data map[string]interface{} `hcl:"data,optional" json:"data"`
}
