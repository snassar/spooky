package interfaces

// VariablesIntegration defines the interface for variables system integration
type VariablesIntegration interface {
	// LoadVariables loads variables from the project
	LoadVariables(projectPath string) (*VariablesContext, error)

	// ResolveVariables resolves variables using facts context
	ResolveVariables(variables *VariablesContext, facts *FactsContext) error

	// ValidateVariables validates variables data
	ValidateVariables(variables *VariablesContext) error

	// SubstituteVariables substitutes variables in a template string
	SubstituteVariables(template string, variables *VariablesContext) (string, error)

	// GetVariable gets a specific variable by name
	GetVariable(name string, context *VariablesContext) (interface{}, error)

	// SetVariable sets a variable value
	SetVariable(name string, value interface{}, context *VariablesContext) error

	// ListVariables lists all available variables
	ListVariables(context *VariablesContext) (map[string]interface{}, error)
}
