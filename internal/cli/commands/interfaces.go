package commands

import (
	spookyclitypes "spooky/internal/cli/types"

	"github.com/spf13/cobra"
)

// CommandsManager defines the interface for command management
type CommandsManager interface {
	// Core command operations
	RegisterCommand(command *spookyclitypes.Command) error
	UnregisterCommand(name string) error
	GetCommand(name string) (*spookyclitypes.Command, error)
	ListCommands() []*spookyclitypes.Command
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
	ValidateCommand(command *spookyclitypes.Command) error
	Close() error
}

// CommandBuilder defines the interface for command building
type CommandBuilder interface {
	BuildCommand(command *spookyclitypes.Command) (*cobra.Command, error)
	BuildSubcommands(parent *cobra.Command, subcommands []*spookyclitypes.Command) error
	ValidateCommandStructure(command *spookyclitypes.Command) error
}

// CommandExecutor defines the interface for command execution
type CommandExecutor interface {
	ExecuteCommand(command *spookyclitypes.Command, args []string) error
	ExecuteSubcommand(parent *spookyclitypes.Command, subcommand string, args []string) error
	ValidateExecution(command *spookyclitypes.Command, args []string) error
}
