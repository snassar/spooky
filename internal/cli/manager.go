package cli

import (
	"fmt"

	"spooky/internal/cli/commands"
	"spooky/internal/cli/completion"
	"spooky/internal/cli/flags"
	"spooky/internal/cli/help"
	"spooky/internal/cli/types"
	"spooky/internal/logging"

	"github.com/spf13/cobra"
)

// Manager implements CLIManager interface
type Manager struct {
	config            *types.Config
	commandsManager   commands.CommandsManager
	completionManager completion.CompletionManager
	helpManager       help.HelpManager
	flagsManager      flags.FlagsManager
	rootCommand       *cobra.Command
	logger            logging.Logger
}

// NewManager creates a new CLI manager
func NewManager(
	config *types.Config,
	commandsManager commands.CommandsManager,
	completionManager completion.CompletionManager,
	helpManager help.HelpManager,
	flagsManager flags.FlagsManager,
	logger logging.Logger,
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
	// 1. Initialize commands manager
	if err := m.commandsManager.InitializeCommands(); err != nil {
		return fmt.Errorf("failed to initialize commands: %w", err)
	}

	// 2. Create root command
	m.rootCommand = &cobra.Command{
		Use:   "spooky",
		Short: "Spooky is a server configuration and automation tool",
		Long: `Spooky is a powerful server configuration and automation tool that allows you to:
- Connect to multiple servers via SSH
- Execute commands and scripts from HCL2 configuration files
- Manage server operations in a declarative way
- Support for parallel execution and error handling
- Collect and manage server facts
- Use templates for dynamic configuration`,
	}

	// 3. Add global flags
	m.addGlobalFlags()

	// 4. Add commands
	m.addCommands()

	m.logger.Info("CLI commands initialized")
	return nil
}

// ExecuteCommand executes a command with arguments
func (m *Manager) ExecuteCommand(args []string) error {
	if m.rootCommand == nil {
		return fmt.Errorf("commands not initialized")
	}

	m.rootCommand.SetArgs(args)
	return m.rootCommand.Execute()
}

// GetRootCommand returns the root command
func (m *Manager) GetRootCommand() *cobra.Command {
	return m.rootCommand
}

// RegisterCommand registers a new command
func (m *Manager) RegisterCommand(command *types.Command) error {
	return m.commandsManager.RegisterCommand(command)
}

// UnregisterCommand unregisters a command
func (m *Manager) UnregisterCommand(name string) error {
	return m.commandsManager.UnregisterCommand(name)
}

// GetCommand gets a command by name
func (m *Manager) GetCommand(name string) (*types.Command, error) {
	return m.commandsManager.GetCommand(name)
}

// ListCommands lists all commands
func (m *Manager) ListCommands() []*types.Command {
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
	if enabled && m.rootCommand != nil {
		m.rootCommand.CompletionOptions.DisableDefaultCmd = false
	}
	return nil
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
	// Close all sub-managers
	if err := m.commandsManager.Close(); err != nil {
		return fmt.Errorf("failed to close commands manager: %w", err)
	}

	return nil
}

// Helper methods
func (m *Manager) addGlobalFlags() {
	// Add global flags to root command
	m.rootCommand.PersistentFlags().StringP("config", "c", "", "config file path")
	m.rootCommand.PersistentFlags().StringP("log-level", "l", "info", "log level")
	m.rootCommand.PersistentFlags().BoolP("quiet", "q", false, "quiet mode")
	m.rootCommand.PersistentFlags().BoolP("verbose", "v", false, "verbose mode")
}

func (m *Manager) addCommands() {
	// Add all commands to root command
	m.rootCommand.AddCommand(m.commandsManager.CreateActionsCommand())
	m.rootCommand.AddCommand(m.commandsManager.CreateFactsCommand())
	m.rootCommand.AddCommand(m.commandsManager.CreateMachinesCommand())
	m.rootCommand.AddCommand(m.commandsManager.CreateProjectCommand())
	m.rootCommand.AddCommand(m.commandsManager.CreateTemplatesCommand())
	m.rootCommand.AddCommand(m.commandsManager.CreateVariablesCommand())
}
