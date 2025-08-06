package types

import (
	"time"
)

// Command represents a CLI command
type Command struct {
	Name        string                 `hcl:"name"`
	Use         string                 `hcl:"use"`
	Short       string                 `hcl:"short,optional"`
	Long        string                 `hcl:"long,optional"`
	Aliases     []string               `hcl:"aliases,optional"`
	Examples    []string               `hcl:"examples,optional"`
	Subcommands []*Command             `hcl:"subcommands,optional"`
	Flags       map[string]interface{} `hcl:"flags,optional"`
	CreatedAt   time.Time              `hcl:"created_at,optional"`
}

// Config represents the main CLI configuration
type Config struct {
	CommandsConfig   *CommandsConfig   `hcl:"commands,optional"`
	CompletionConfig *CompletionConfig `hcl:"completion,optional"`
	HelpConfig       *HelpConfig       `hcl:"help,optional"`
	FlagsConfig      *FlagsConfig      `hcl:"flags,optional"`
	EnableCompletion bool              `hcl:"enable_completion,optional"`
	EnableHelp       bool              `hcl:"enable_help,optional"`

	// Logging configuration
	LogLevel string `hcl:"log_level,optional"`
	LogFile  string `hcl:"log_file,optional"`
	Quiet    bool   `hcl:"quiet,optional"`
	Verbose  bool   `hcl:"verbose,optional"`
}

// CommandsConfig represents commands configuration
type CommandsConfig struct {
	AutoInitialize   bool `hcl:"auto_initialize,optional"`
	ValidateCommands bool `hcl:"validate_commands,optional"`
}

// CompletionConfig represents completion configuration
type CompletionConfig struct {
	EnabledShells []string `hcl:"enabled_shells,optional"`
	OutputPath    string   `hcl:"output_path,optional"`
}

// HelpConfig represents help configuration
type HelpConfig struct {
	EnableExamples bool `hcl:"enable_examples,optional"`
	EnableUsage    bool `hcl:"enable_usage,optional"`
}

// FlagsConfig represents flags configuration
type FlagsConfig struct {
	GlobalFlags  map[string]interface{}            `hcl:"global_flags,optional"`
	CommandFlags map[string]map[string]interface{} `hcl:"command_flags,optional"`
}
