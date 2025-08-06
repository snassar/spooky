package commands

import (
	"fmt"

	spookyclitypes "spooky/internal/cli/types"
)

// Executor implements CommandExecutor interface
type Executor struct{}

// NewExecutor creates a new command executor
func NewExecutor() *Executor {
	return &Executor{}
}

// ExecuteCommand executes a command
func (e *Executor) ExecuteCommand(command *spookyclitypes.Command, args []string) error {
	if command == nil {
		return fmt.Errorf("command cannot be nil")
	}

	// Validate execution
	if err := e.ValidateExecution(command, args); err != nil {
		return fmt.Errorf("execution validation failed: %w", err)
	}

	// TODO: Implement actual command execution logic
	// This would typically involve calling the coordinator to execute the command

	return nil
}

// ExecuteSubcommand executes a subcommand
func (e *Executor) ExecuteSubcommand(parent *spookyclitypes.Command, subcommand string, args []string) error {
	if parent == nil {
		return fmt.Errorf("parent command cannot be nil")
	}

	if subcommand == "" {
		return fmt.Errorf("subcommand name cannot be empty")
	}

	// Find the subcommand
	var targetCommand *spookyclitypes.Command
	for _, cmd := range parent.Subcommands {
		if cmd.Name == subcommand {
			targetCommand = cmd
			break
		}
	}

	if targetCommand == nil {
		return fmt.Errorf("subcommand not found: %s", subcommand)
	}

	// Execute the subcommand
	return e.ExecuteCommand(targetCommand, args)
}

// ValidateExecution validates command execution
func (e *Executor) ValidateExecution(command *spookyclitypes.Command, _ []string) error {
	if command == nil {
		return fmt.Errorf("command cannot be nil")
	}

	// Basic validation - can be extended based on specific requirements
	if command.Name == "" {
		return fmt.Errorf("command name cannot be empty")
	}

	if command.Use == "" {
		return fmt.Errorf("command use cannot be empty")
	}

	// TODO: Add more specific validation logic based on command type and arguments

	return nil
}
