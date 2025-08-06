package interfaces

import (
	spookyactionstypes "spooky/internal/actions/types"
	spookyconfigtypes "spooky/internal/config/types"
	spookyfactstypes "spooky/internal/facts/types"
	spookytemplatestypes "spooky/internal/templates/types"
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
	GlobalConfig  *spookyconfigtypes.GlobalConfig
	ProjectConfig *spookyconfigtypes.ProjectConfig
	Environment   map[string]interface{}
}

// FactsContext provides facts data for integrations
type FactsContext struct {
	BaseContext
	MachineFacts map[string]*spookyfactstypes.FactCollection
	GlobalFacts  *spookyfactstypes.FactCollection
	ProjectFacts *spookyfactstypes.FactCollection
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
	Templates     map[string]*spookytemplatestypes.Template
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
	Action           *spookyactionstypes.Action
	Decrypt          bool
}
