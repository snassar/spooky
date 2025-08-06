package commands

import (
	"fmt"
	"sync"

	spookyclitypes "spooky/internal/cli/types"
	spookylogging "spooky/internal/logging"

	"github.com/spf13/cobra"
)

// Manager implements CommandsManager interface
type Manager struct {
	config   *spookyclitypes.CommandsConfig
	commands map[string]*spookyclitypes.Command
	builder  CommandBuilder
	executor CommandExecutor
	logger   spookylogging.Logger
	mutex    sync.RWMutex
}

// NewManager creates a new commands manager
func NewManager(
	config *spookyclitypes.CommandsConfig,
	builder CommandBuilder,
	executor CommandExecutor,
	logger spookylogging.Logger,
) *Manager {
	return &Manager{
		config:   config,
		commands: make(map[string]*spookyclitypes.Command),
		builder:  builder,
		executor: executor,
		logger:   logger,
	}
}

// RegisterCommand registers a new command
func (m *Manager) RegisterCommand(command *spookyclitypes.Command) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 1. Validate command
	if err := m.ValidateCommand(command); err != nil {
		return fmt.Errorf("command validation failed: %w", err)
	}

	// 2. Check for conflicts
	if _, exists := m.commands[command.Name]; exists {
		return fmt.Errorf("command already exists: %s", command.Name)
	}

	// 3. Register command
	m.commands[command.Name] = command

	m.logger.Info("Command registered", spookylogging.String("name", command.Name))
	return nil
}

// UnregisterCommand unregisters a command
func (m *Manager) UnregisterCommand(name string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.commands[name]; !exists {
		return fmt.Errorf("command not found: %s", name)
	}

	delete(m.commands, name)

	m.logger.Info("Command unregistered", spookylogging.String("name", name))
	return nil
}

// GetCommand gets a command by name
func (m *Manager) GetCommand(name string) (*spookyclitypes.Command, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	command, exists := m.commands[name]
	if !exists {
		return nil, fmt.Errorf("command not found: %s", name)
	}

	return command, nil
}

// ListCommands lists all registered commands
func (m *Manager) ListCommands() []*spookyclitypes.Command {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	commands := make([]*spookyclitypes.Command, 0, len(m.commands))
	for _, command := range m.commands {
		commands = append(commands, command)
	}

	return commands
}

// InitializeCommands initializes all commands
func (m *Manager) InitializeCommands() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Initialize all registered commands
	for _, command := range m.commands {
		if err := m.builder.ValidateCommandStructure(command); err != nil {
			return fmt.Errorf("command structure validation failed for %s: %w", command.Name, err)
		}
	}

	m.logger.Info("All commands initialized", spookylogging.Int("count", len(m.commands)))
	return nil
}

// CreateActionsCommand creates the actions command
func (m *Manager) CreateActionsCommand() *cobra.Command {
	return CreateActionsCommand()
}

// CreateFactsCommand creates the facts command
func (m *Manager) CreateFactsCommand() *cobra.Command {
	return CreateFactsCommand()
}

// CreateMachinesCommand creates the machines command
func (m *Manager) CreateMachinesCommand() *cobra.Command {
	return CreateMachinesCommand()
}

// CreateProjectCommand creates the project command
func (m *Manager) CreateProjectCommand() *cobra.Command {
	return CreateProjectCommand()
}

// CreateTemplatesCommand creates the templates command
func (m *Manager) CreateTemplatesCommand() *cobra.Command {
	return CreateTemplatesCommand()
}

// CreateVariablesCommand creates the variables command
func (m *Manager) CreateVariablesCommand() *cobra.Command {
	return CreateVariablesCommand()
}

// SetCommandFlags sets flags for a command
func (m *Manager) SetCommandFlags(commandName string, flags map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	command, exists := m.commands[commandName]
	if !exists {
		return fmt.Errorf("command not found: %s", commandName)
	}

	command.Flags = flags
	return nil
}

// SetCommandExamples sets examples for a command
func (m *Manager) SetCommandExamples(commandName string, examples []string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	command, exists := m.commands[commandName]
	if !exists {
		return fmt.Errorf("command not found: %s", commandName)
	}

	command.Examples = examples
	return nil
}

// ValidateCommand validates a command
func (m *Manager) ValidateCommand(command *spookyclitypes.Command) error {
	if command == nil {
		return fmt.Errorf("command cannot be nil")
	}

	if command.Name == "" {
		return fmt.Errorf("command name cannot be empty")
	}

	if command.Use == "" {
		return fmt.Errorf("command use cannot be empty")
	}

	return nil
}

// Close closes the commands manager
func (m *Manager) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Clear all commands
	m.commands = make(map[string]*spookyclitypes.Command)

	m.logger.Info("Commands manager closed")
	return nil
}
