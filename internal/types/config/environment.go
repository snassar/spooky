package config

// EnvironmentVariable represents an environment variable
type EnvironmentVariable struct {
	Name         string      `hcl:"name"`
	Value        string      `hcl:"value"`
	Description  string      `hcl:"description,optional"`
	DefaultValue interface{} `hcl:"default_value,optional"`
	Required     bool        `hcl:"required,optional"`
	Type         string      `hcl:"type,optional" validate:"omitempty,oneof=string int bool"`
	Validation   string      `hcl:"validation,optional"`
}

// EnvironmentSettings represents environment configuration
type EnvironmentSettings struct {
	Variables map[string]EnvironmentVariable `hcl:"variables,optional"`
	Files     []string                       `hcl:"files,optional"`
	Prefix    string                         `hcl:"prefix,optional"`
	Validate  bool                           `hcl:"validate,optional"`
}

// EnvironmentSource represents the source of an environment variable
type EnvironmentSource string

const (
	EnvSourceSystem  EnvironmentSource = "system"
	EnvSourceFile    EnvironmentSource = "file"
	EnvSourceCLI     EnvironmentSource = "cli"
	EnvSourceDefault EnvironmentSource = "default"
)

// EnvironmentValue represents an environment variable value with its source
type EnvironmentValue struct {
	Value  interface{}       `hcl:"value"`
	Source EnvironmentSource `hcl:"source"`
	File   string            `hcl:"file,optional"`
	Line   int               `hcl:"line,optional"`
}

// EnvironmentValidationRule represents a validation rule for environment variables
type EnvironmentValidationRule struct {
	Name          string   `hcl:"name"`
	Pattern       string   `hcl:"pattern,optional"`
	MinLength     int      `hcl:"min_length,optional"`
	MaxLength     int      `hcl:"max_length,optional"`
	MinValue      int      `hcl:"min_value,optional"`
	MaxValue      int      `hcl:"max_value,optional"`
	AllowedValues []string `hcl:"allowed_values,optional"`
	Required      bool     `hcl:"required,optional"`
}

// EnvironmentValidationResult represents the result of environment validation
type EnvironmentValidationResult struct {
	Valid    bool               `hcl:"valid"`
	Errors   []EnvironmentError `hcl:"errors,optional"`
	Warnings []EnvironmentError `hcl:"warnings,optional"`
}
