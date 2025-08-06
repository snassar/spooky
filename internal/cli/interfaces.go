package cli

import (
	spookyclitypes "spooky/internal/cli/types"

	"github.com/spf13/cobra"
)

// CLIManager defines the main interface for CLI operations
type CLIManager interface {
	// Core CLI operations
	InitializeCommands() error
	ExecuteCommand(args []string) error
	GetRootCommand() *cobra.Command

	// Command management
	RegisterCommand(command *spookyclitypes.Command) error
	UnregisterCommand(name string) error
	GetCommand(name string) (*spookyclitypes.Command, error)
	ListCommands() []*spookyclitypes.Command

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
	RegisterCommand(command *spookyclitypes.Command) error
	UnregisterCommand(name string) error
	GetCommand(name string) (*spookyclitypes.Command, error)
	ListCommands() []*spookyclitypes.Command
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
