package interfaces

import (
	"spooky/internal/actions/types"
	configtypes "spooky/internal/config/types"
	factstypes "spooky/internal/facts/types"
	templatesTypes "spooky/internal/templates/types"
	"time"
)

// BaseContext provides base context for all integrations
type BaseContext struct {
	ProjectPath string
	Timestamp   time.Time
	Metadata    map[string]interface{}
}

// ConfigContext provides configuration data for integrations
type ConfigContext struct {
	BaseContext
	GlobalConfig  *configtypes.GlobalConfig
	ProjectConfig *configtypes.ProjectConfig
	Environment   map[string]interface{}
}

// FactsContext provides facts data for integrations
type FactsContext struct {
	BaseContext
	MachineFacts map[string]*factstypes.FactCollection
	GlobalFacts  *factstypes.FactCollection
	ProjectFacts *factstypes.FactCollection
	CacheKey     string
}

// VariablesContext provides variables data for integrations
type VariablesContext struct {
	BaseContext
	ResolvedVariables map[string]interface{}
	VariableContext   map[string]interface{}
	ResolutionContext map[string]interface{}
}

// TemplatesContext provides templates data for integrations
type TemplatesContext struct {
	BaseContext
	Templates     map[string]*templatesTypes.Template
	RenderedCache map[string]string
	Functions     map[string]interface{}
}

// ActionExecutionContext provides all context for action execution
type ActionExecutionContext struct {
	BaseContext
	FactsContext     *FactsContext
	VariablesContext *VariablesContext
	TemplatesContext *TemplatesContext
	MachineNames     []string
	Action           *types.Action
	Decrypt          bool
}
