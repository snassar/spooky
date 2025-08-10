package interfaces

import (
	spookytypesactions "spooky/internal/types/actions"
	spookytypesconfig "spooky/internal/types/config"
	spookytypesfacts "spooky/internal/types/facts"
	spookytypestemplates "spooky/internal/types/templates"
	spookytypesvariables "spooky/internal/types/variables"
)

// IntegrationManager coordinates all system integrations
type IntegrationManager interface {
	// System integrations
	Facts() FactsIntegration
	Actions() ActionsIntegration
	Variables() VariablesIntegration
	Templates() TemplatesIntegration
	Machines() MachinesIntegration
	Secrets() SecretsIntegration

	// Context management
	LoadContextForAction(action *spookytypesconfig.Action, projectPath string, machineNames []string) (*ActionExecutionContext, error)
	ValidateActionWithAllSystems(action *spookytypesconfig.Action, context *ActionExecutionContext) error
	PrepareActionForExecution(action *spookytypesconfig.Action, context *ActionExecutionContext) error
	ExecuteAction(action *spookytypesconfig.Action, context *ActionExecutionContext) error

	// Health and status
	ValidateIntegrationHealth() error
	GetIntegrationStats() map[string]interface{}
	ClearAllCaches() error

	// Project operations
	LoadProjectContext(projectPath string) (*ProjectContext, error)
	ValidateProject(projectPath string) error
	GetProjectStats(projectPath string) map[string]interface{}
}

// FactsIntegration defines the interface for facts system integration
type FactsIntegration interface {
	CollectAllFacts(server string) (*spookytypesfacts.FactCollection, error)
	GetFact(server, key string) (*spookytypesfacts.Fact, error)
	ValidateFacts(collection *spookytypesfacts.FactCollection) error
}

// ActionsIntegration defines the interface for actions system integration
type ActionsIntegration interface {
	LoadActions(projectPath string) (*spookytypesactions.ActionCollection, error)
	GetAction(name string) (*spookytypesactions.Action, error)
	ValidateAction(action *spookytypesactions.Action) error
}

// VariablesIntegration defines the interface for variables system integration
type VariablesIntegration interface {
	LoadVariables(projectPath string) ([]*spookytypesvariables.Variable, error)
	GetVariable(name string) (*spookytypesvariables.Variable, error)
	ValidateVariables(variables []*spookytypesvariables.Variable) error
}

// TemplatesIntegration defines the interface for templates system integration
type TemplatesIntegration interface {
	LoadTemplates(projectPath string) ([]*spookytypestemplates.Template, error)
	GetTemplate(name string) (*spookytypestemplates.Template, error)
	ValidateTemplate(template *spookytypestemplates.Template) error
}

// MachinesIntegration defines the interface for machines system integration
type MachinesIntegration interface {
	LoadMachines(projectPath string) ([]*spookytypesconfig.Machine, error)
	GetMachine(name string) (*spookytypesconfig.Machine, error)
	ValidateMachine(machine *spookytypesconfig.Machine) error
}

// ConfigIntegration defines the interface for configuration system integration
type ConfigIntegration interface {
	LoadConfig(projectPath string) (*ConfigContext, error)
	ValidateConfig(projectPath string) error
	GetConfigValue(path string) (interface{}, error)
	SetConfigValue(path string, value interface{}) error
	GetConfigString(path string) (string, error)
	GetConfigInt(path string) (int, error)
	GetConfigBool(path string) (bool, error)
	GetConfigStringSlice(path string) ([]string, error)
	ApplyCLIFlags(flags map[string]interface{}) error
	GetConfigurationPath() string
	CreateDefaultConfig() error
	ReloadConfig() error
	GetProjectConfig() interface{}
}

// ProjectContext provides project-wide context
type ProjectContext struct {
	BaseContext
	FactsContext     *FactsContext
	VariablesContext *VariablesContext
	TemplatesContext *TemplatesContext
	MachinesContext  *MachinesContext
	ActionsContext   *ActionsContext
}
