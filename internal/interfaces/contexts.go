package interfaces

import (
	spookytypesactions "spooky/internal/types/actions"
	spookytypesconfig "spooky/internal/types/config"
	spookytypesfacts "spooky/internal/types/facts"
	spookytypestemplates "spooky/internal/types/templates"
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
	GlobalConfig  *spookytypesconfig.GlobalConfig
	ProjectConfig *spookytypesconfig.ProjectConfig
	Environment   map[string]interface{}
}

// FactsContext provides facts data for integrations
type FactsContext struct {
	BaseContext
	MachineFacts map[string]*spookytypesfacts.FactCollection
	GlobalFacts  *spookytypesfacts.FactCollection
	ProjectFacts *spookytypesfacts.FactCollection
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
	Templates     map[string]*spookytypestemplates.Template
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
	Action           *spookytypesactions.Action
	Decrypt          bool
}
