package environment

// EnvironmentManager defines the interface for environment variable management
type EnvironmentManager interface {
	// Core environment operations
	GetEnvironmentVariable(name string) (interface{}, error)
	GetAllEnvironmentVariables() (map[string]interface{}, error)
	ValidateEnvironmentVariable(name, value string) error

	// Configuration
	SetEnvironmentVariable(name, value string) error
	UnsetEnvironmentVariable(name string) error
	LoadEnvironmentFile(path string) error

	// Utility operations
	GetEnvironmentVariableDescription(name string) string
	GetEnvironmentVariableDefault(name string) interface{}
	ListEnvironmentVariables() []string
	Close() error
}

// EnvironmentValidator defines the interface for environment validation
type EnvironmentValidator interface {
	ValidateVariable(name, value string) error
	ValidateVariableType(name, value, expectedType string) error
	ValidateVariableFormat(name, value, format string) error
}
