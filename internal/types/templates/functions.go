package templates

import (
	"time"
)

// FunctionsConfig represents template functions configuration
// Aligns with template-functions.schema.hcl schema
type FunctionsConfig struct {
	AllowedFunctions   []string      `hcl:"allowed_functions,optional" json:"allowed_functions"`
	RestrictedPatterns []string      `hcl:"restricted_patterns,optional" json:"restricted_patterns"`
	MaxTemplateSize    int64         `hcl:"max_template_size,optional" json:"max_template_size"`
	MaxNestingDepth    int           `hcl:"max_nesting_depth,optional" json:"max_nesting_depth"`
	BuiltinFunctions   bool          `hcl:"builtin_functions,optional" json:"builtin_functions"`
	FunctionTimeout    time.Duration `hcl:"function_timeout,optional" json:"function_timeout"`
}

// EngineConfig represents template engine configuration
type EngineConfig struct {
	Delimiters       []string      `hcl:"delimiters,optional" json:"delimiters"`
	MaxExecutionTime time.Duration `hcl:"max_execution_time,optional" json:"max_execution_time"`
	StrictMode       bool          `hcl:"strict_mode,optional" json:"strict_mode"`
}

// ValidationConfig represents template validation configuration
type ValidationConfig struct {
	ValidationRules  *ValidationRules `hcl:"validation_rules,optional" json:"validation_rules"`
	StrictValidation bool             `hcl:"strict_validation,optional" json:"strict_validation"`
}

// SecretsConfig represents template secrets configuration
type SecretsConfig struct {
	Enabled             bool   `hcl:"enabled,optional" json:"enabled"`
	EncryptionKey       string `hcl:"encryption_key,optional" json:"encryption_key"`
	EncryptionAlgorithm string `hcl:"encryption_algorithm,optional" json:"encryption_algorithm"`
}

// Config represents the main template configuration
type Config struct {
	EngineConfig     *EngineConfig     `hcl:"engine,optional" json:"engine"`
	FunctionsConfig  *FunctionsConfig  `hcl:"functions,optional" json:"functions"`
	ValidationConfig *ValidationConfig `hcl:"validation,optional" json:"validation"`
	SecretsConfig    *SecretsConfig    `hcl:"secrets,optional" json:"secrets"`
	DefaultTimeout   time.Duration     `hcl:"default_timeout,optional" json:"default_timeout"`
	MaxTemplateSize  int64             `hcl:"max_template_size,optional" json:"max_template_size"`
}

// ValidationRules represents validation rules
type ValidationRules struct {
	RequiredFields  []string `hcl:"required_fields,optional" json:"required_fields"`
	ForbiddenFields []string `hcl:"forbidden_fields,optional" json:"forbidden_fields"`
}
