package cli

import (
	"fmt"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypescli "spooky/internal/types/cli"

	"github.com/spf13/cobra"
)

// Manager implements CLIManager interface
type Manager struct {
	config            *spookytypescli.Config
	commandsManager   spookyinterfaces.CommandsManager
	completionManager spookyinterfaces.CompletionManager
	helpManager       spookyinterfaces.HelpManager
	flagsManager      spookyinterfaces.FlagsManager
	rootCommand       *cobra.Command
	logger            spookyinterfaces.Logger
}

// NewManager creates a new CLI manager
func NewManager(
	config *spookytypescli.Config,
	commandsManager spookyinterfaces.CommandsManager,
	completionManager spookyinterfaces.CompletionManager,
	helpManager spookyinterfaces.HelpManager,
	flagsManager spookyinterfaces.FlagsManager,
	logger spookyinterfaces.Logger,
) *Manager {
	return &Manager{
		config:            config,
		commandsManager:   commandsManager,
		completionManager: completionManager,
		helpManager:       helpManager,
		flagsManager:      flagsManager,
		logger:            logger,
	}
}

// InitializeCommands initializes all CLI commands
func (m *Manager) InitializeCommands() error {
	// TODO: Implement properly with correct types
	return fmt.Errorf("not implemented - interface mismatches need to be resolved")
}

// ExecuteCommand executes a command with arguments
func (m *Manager) ExecuteCommand(args []string) error {
	// TODO: Implement properly with correct types
	return fmt.Errorf("not implemented - interface mismatches need to be resolved")
}

// GetRootCommand returns the root command
func (m *Manager) GetRootCommand() *cobra.Command {
	// TODO: Implement properly with correct types
	return nil
}

// RegisterCommand registers a new command
func (m *Manager) RegisterCommand(command *spookytypescli.Command) error {
	return m.commandsManager.RegisterCommand(command)
}

// UnregisterCommand unregisters a command
func (m *Manager) UnregisterCommand(name string) error {
	return m.commandsManager.UnregisterCommand(name)
}

// GetCommand gets a command by name
func (m *Manager) GetCommand(name string) (*spookytypescli.Command, error) {
	return m.commandsManager.GetCommand(name)
}

// ListCommands lists all commands
func (m *Manager) ListCommands() []*spookytypescli.Command {
	return m.commandsManager.ListCommands()
}

// SetGlobalFlags sets global flags
func (m *Manager) SetGlobalFlags(flags map[string]interface{}) error {
	return m.flagsManager.SetGlobalFlags(flags)
}

// SetCommandFlags sets flags for a specific command
func (m *Manager) SetCommandFlags(commandName string, flags map[string]interface{}) error {
	return m.flagsManager.SetCommandFlags(commandName, flags)
}

// EnableCompletion enables command completion
func (m *Manager) EnableCompletion(enabled bool) error {
	// TODO: Implement properly with correct types
	return fmt.Errorf("not implemented - interface mismatches need to be resolved")
}

// GenerateCompletion generates completion for a shell
func (m *Manager) GenerateCompletion(shell string) (string, error) {
	return m.completionManager.GenerateCompletion(shell)
}

// ShowHelp shows help for a command
func (m *Manager) ShowHelp(commandName string) (string, error) {
	return m.helpManager.ShowHelp(commandName)
}

// Close closes the CLI manager
func (m *Manager) Close() error {
	// TODO: Implement properly with correct types
	return fmt.Errorf("not implemented - interface mismatches need to be resolved")
}
