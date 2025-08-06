package commands

import (
	"spooky/internal/cli/types"

	"github.com/spf13/cobra"
)

// CommandsManager defines the interface for command management
type CommandsManager interface {
	// Core command operations
	RegisterCommand(command *types.Command) error
	UnregisterCommand(name string) error
	GetCommand(name string) (*types.Command, error)
	ListCommands() []*types.Command
	InitializeCommands() error

	// Command creation
	CreateActionsCommand() *cobra.Command
	CreateFactsCommand() *cobra.Command
	CreateMachinesCommand() *cobra.Command
	CreateProjectCommand() *cobra.Command
	CreateTemplatesCommand() *cobra.Command
	CreateVariablesCommand() *cobra.Command

	// Configuration
	SetCommandFlags(commandName string, flags map[string]interface{}) error
	SetCommandExamples(commandName string, examples []string) error

	// Utility operations
	ValidateCommand(command *types.Command) error
	Close() error
}

// CommandBuilder defines the interface for command building
type CommandBuilder interface {
	BuildCommand(command *types.Command) (*cobra.Command, error)
	BuildSubcommands(parent *cobra.Command, subcommands []*types.Command) error
	ValidateCommandStructure(command *types.Command) error
}

// CommandExecutor defines the interface for command execution
type CommandExecutor interface {
	ExecuteCommand(command *types.Command, args []string) error
	ExecuteSubcommand(parent *types.Command, subcommand string, args []string) error
	ValidateExecution(command *types.Command, args []string) error
}
