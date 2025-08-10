package interfaces

import (
	spookytypescli "spooky/internal/types/cli"

	"github.com/spf13/cobra"
)

// CLIManager defines the main interface for CLI operations
type CLIManager interface {
	// Core CLI operations
	InitializeCommands() error
	ExecuteCommand(args []string) error
	GetRootCommand() *cobra.Command

	// Command management
	RegisterCommand(command *spookytypescli.Command) error
	UnregisterCommand(name string) error
	GetCommand(name string) (*spookytypescli.Command, error)
	ListCommands() []*spookytypescli.Command

	// Configuration
	SetGlobalFlags(flags map[string]interface{}) error
	SetCommandFlags(commandName string, flags map[string]interface{}) error
	EnableCompletion(enabled bool) error

	// Utility operations
	GenerateCompletion(shell string) (string, error)
	ShowHelp(commandName string) (string, error)
	Close() error
}

// CommandsManager defines the interface for command management
type CommandsManager interface {
	RegisterCommand(command *spookytypescli.Command) error
	UnregisterCommand(name string) error
	GetCommand(name string) (*spookytypescli.Command, error)
	ListCommands() []*spookytypescli.Command
	InitializeCommands() error
}

// CompletionManager defines the interface for completion generation
type CompletionManager interface {
	GenerateCompletion(shell string) (string, error)
	GenerateCompletionFile(shell, outputPath string) error
	GetSupportedShells() []string
}

// HelpManager defines the interface for help rendering
type HelpManager interface {
	ShowHelp(commandName string) (string, error)
	ShowUsage(commandName string) (string, error)
	ShowExamples(commandName string) (string, error)
}

// FlagsManager defines the interface for flag management
type FlagsManager interface {
	SetGlobalFlags(flags map[string]interface{}) error
	SetCommandFlags(commandName string, flags map[string]interface{}) error
	GetGlobalFlags() map[string]interface{}
	GetCommandFlags(commandName string) map[string]interface{}
}

// CommandBuilder defines the interface for command building
type CommandBuilder interface {
	BuildCommand(command *spookytypescli.Command) (*cobra.Command, error)
	BuildSubcommands(parent *cobra.Command, subcommands []*spookytypescli.Command) error
	ValidateCommandStructure(command *spookytypescli.Command) error
}

// CommandExecutor defines the interface for command execution
type CommandExecutor interface {
	ExecuteCommand(command *spookytypescli.Command, args []string) error
	ExecuteSubcommand(parent *spookytypescli.Command, subcommand string, args []string) error
	ValidateExecution(command *spookytypescli.Command, args []string) error
}

// CompletionGenerator defines the interface for completion generation
type CompletionGenerator interface {
	GenerateBashCompletion() (string, error)
	GenerateZshCompletion() (string, error)
	GenerateFishCompletion() (string, error)
	GeneratePowerShellCompletion() (string, error)
}

// FlagsParser defines the interface for flag parsing
type FlagsParser interface {
	ParseFlags(cmd interface{}) error
	ValidateFlags(flags map[string]interface{}) error
	GetFlagValue(name string) interface{}
}

// HelpRenderer defines the interface for help rendering
type HelpRenderer interface {
	RenderHelp(commandName string) (string, error)
	RenderUsage(commandName string) (string, error)
	RenderExamples(commandName string) (string, error)
}
