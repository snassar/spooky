package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	spookyclitypes "spooky/internal/types/cli"
)

// Builder implements CommandBuilder interface
type Builder struct{}

// NewBuilder creates a new command builder
func NewBuilder() *Builder {
	return &Builder{}
}

// BuildCommand builds a cobra command from a command type
func (b *Builder) BuildCommand(command *spookyclitypes.Command) (*cobra.Command, error) {
	if command == nil {
		return nil, fmt.Errorf("command cannot be nil")
	}

	cobraCmd := &cobra.Command{
		Use:     command.Use,
		Short:   command.Short,
		Long:    command.Long,
		Aliases: command.Aliases,
		Example: b.buildExample(command.Examples),
	}

	// Add subcommands if any
	if len(command.Subcommands) > 0 {
		if err := b.BuildSubcommands(cobraCmd, command.Subcommands); err != nil {
			return nil, fmt.Errorf("failed to build subcommands: %w", err)
		}
	}

	return cobraCmd, nil
}

// BuildSubcommands builds subcommands for a parent command
func (b *Builder) BuildSubcommands(parent *cobra.Command, subcommands []*spookyclitypes.Command) error {
	for _, subcommand := range subcommands {
		cobraSubcommand, err := b.BuildCommand(subcommand)
		if err != nil {
			return fmt.Errorf("failed to build subcommand %s: %w", subcommand.Name, err)
		}
		parent.AddCommand(cobraSubcommand)
	}
	return nil
}

// ValidateCommandStructure validates a command structure
func (b *Builder) ValidateCommandStructure(command *spookyclitypes.Command) error {
	if command == nil {
		return fmt.Errorf("command cannot be nil")
	}

	if command.Name == "" {
		return fmt.Errorf("command name cannot be empty")
	}

	if command.Use == "" {
		return fmt.Errorf("command use cannot be empty")
	}

	// Validate subcommands recursively
	for _, subcommand := range command.Subcommands {
		if err := b.ValidateCommandStructure(subcommand); err != nil {
			return fmt.Errorf("subcommand %s validation failed: %w", subcommand.Name, err)
		}
	}

	return nil
}

// Helper method to build example string
func (b *Builder) buildExample(examples []string) string {
	if len(examples) == 0 {
		return ""
	}

	var result string
	for i, example := range examples {
		if i > 0 {
			result += "\n"
		}
		result += example
	}
	return result
}
