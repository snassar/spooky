package commands

import (
	"fmt"

	spookytypescli "spooky/internal/types/cli"
)

// Actor implements CommandActor interface
type Actor struct{}

// NewActor creates a new command actor
func NewActor() *Actor {
	return &Actor{}
}

// RunCommand runs a command
func (a *Actor) RunCommand(command *spookytypescli.Command, args []string) error {
	if command == nil {
		return fmt.Errorf("command cannot be nil")
	}

	// Validate command
	if err := a.ValidateCommand(command, args); err != nil {
		return fmt.Errorf("command validation failed: %w", err)
	}

	// TODO: Implement actual command running logic
	// This would typically involve calling the coordinator to run the command

	return nil
}

// RunSubcommand runs a subcommand
func (a *Actor) RunSubcommand(parent *spookytypescli.Command, subcommand string, args []string) error {
	if parent == nil {
		return fmt.Errorf("parent command cannot be nil")
	}

	if subcommand == "" {
		return fmt.Errorf("subcommand name cannot be empty")
	}

	// Find the subcommand
	var targetCommand *spookytypescli.Command
	for _, cmd := range parent.Subcommands {
		if cmd.Name == subcommand {
			targetCommand = cmd
			break
		}
	}

	if targetCommand == nil {
		return fmt.Errorf("subcommand not found: %s", subcommand)
	}

	// Run the subcommand
	return a.RunCommand(targetCommand, args)
}

// ValidateCommand validates command
func (a *Actor) ValidateCommand(command *spookytypescli.Command, _ []string) error {
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
