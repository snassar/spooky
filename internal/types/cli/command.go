// Package cli provides types for command-line interface operations in the spooky codebase.
// These types define the structure for CLI commands, flags, and run context.
package cli

import (
	"context"
)

// Command represents a CLI command with its run logic
type Command interface {
	// Name returns the command name
	Name() string

	// Description returns the command description
	Description() string

	// Usage returns the command usage string
	Usage() string

	// Run runs the command with the given context and arguments
	Execute(ctx context.Context, args []string) error

	// Validate validates the command arguments and flags
	Validate(args []string) error
}

// CommandContext provides context for command running
type CommandContext struct {
	// Command name
	Name string `json:"name" hcl:"name"`

	// Command description
	Description string `json:"description" hcl:"description"`

	// Command usage
	Usage string `json:"usage" hcl:"usage"`

	// Command flags
	Flags *CommandFlags `json:"flags" hcl:"flags"`

	// Command arguments
	Args []string `json:"args" hcl:"args"`

	// Command metadata
	Metadata map[string]interface{} `json:"metadata,omitempty" hcl:"metadata,optional"`
}

// CommandFlags represents command flags
type CommandFlags struct {
	// Global flags
	Global map[string]interface{} `json:"global,omitempty" hcl:"global,optional"`

	// Command-specific flags
	Command map[string]interface{} `json:"command,omitempty" hcl:"command,optional"`

	// Flag validation rules
	Validation map[string]interface{} `json:"validation,omitempty" hcl:"validation,optional"`
}

// CommandError represents a command-related error
type CommandError struct {
	// Error details
	Code        string                 `json:"code" hcl:"code"`
	Message     string                 `json:"message" hcl:"message"`
	Context     map[string]interface{} `json:"context,omitempty" hcl:"context,optional"`
	Stack       []string               `json:"stack,omitempty" hcl:"stack,optional"`
	Recoverable bool                   `json:"recoverable" hcl:"recoverable"`

	// Command information
	CommandName string `json:"command_name" hcl:"command_name"`
	CommandType string `json:"command_type" hcl:"command_type"`

	// Error severity
	Severity string `json:"severity" hcl:"severity"` // "error", "warning", "info"
}

// NewCommandError creates a new command error
func NewCommandError(commandName, commandType, message string) *CommandError {
	return &CommandError{
		Code:        "command_error",
		Message:     message,
		Recoverable: true,
		CommandName: commandName,
		CommandType: commandType,
		Severity:    "error",
	}
}

// Error implements the error interface
func (e *CommandError) Error() string {
	return e.Message
}

// Unwrap returns the underlying error
func (e *CommandError) Unwrap() error {
	return nil
}

// CommandHelp provides help information for a command
type CommandHelp struct {
	Name        string   `json:"name" hcl:"name"`
	Description string   `json:"description" hcl:"description"`
	Usage       string   `json:"usage" hcl:"usage"`
	Examples    []string `json:"examples,omitempty" hcl:"examples,optional"`
	Flags       []Flag   `json:"flags,omitempty" hcl:"flags,optional"`
	Subcommands []string `json:"subcommands,omitempty" hcl:"subcommands,optional"`
}

// Flag represents a command-line flag
type Flag struct {
	Name        string      `json:"name" hcl:"name"`
	Short       string      `json:"short,omitempty" hcl:"short,optional"`
	Description string      `json:"description" hcl:"description"`
	Required    bool        `json:"required" hcl:"required"`
	Default     interface{} `json:"default,omitempty" hcl:"default,optional"`
	Type        string      `json:"type" hcl:"type"` // "string", "int", "bool", etc.
}

// CommandRegistry manages available commands
type CommandRegistry interface {
	// Register registers a new command
	Register(command Command) error

	// Get returns a command by name
	Get(name string) (Command, bool)

	// List returns all registered commands
	List() []Command

	// Run runs a command by name
	Run(ctx context.Context, name string, args []string) error
}

// CommandExecutor provides command run functionality
type CommandExecutor interface {
	// Run runs a command with the given context
	Run(ctx context.Context, command Command, args []string) error

	// Validate validates command arguments before running
	Validate(command Command, args []string) error

	// Help provides help information for a command
	Help(command Command) *CommandHelp
}

// CommandFactory creates command instances
type CommandFactory interface {
	// Create creates a new command instance
	Create(name string) (Command, error)

	// Available returns available command names
	Available() []string
}
