package interfaces

import "spooky/internal/config/types"

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
	LoadContextForAction(action *types.Action, projectPath string, machineNames []string) (*ActionExecutionContext, error)
	ValidateActionWithAllSystems(action *types.Action, context *ActionExecutionContext) error
	PrepareActionForExecution(action *types.Action, context *ActionExecutionContext) error
	ExecuteAction(action *types.Action, context *ActionExecutionContext) error

	// Health and status
	ValidateIntegrationHealth() error
	GetIntegrationStats() map[string]interface{}
	ClearAllCaches() error

	// Project operations
	LoadProjectContext(projectPath string) (*ProjectContext, error)
	ValidateProject(projectPath string) error
	GetProjectStats(projectPath string) map[string]interface{}
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
