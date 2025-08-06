package flags

import (
	"fmt"
	"sync"

	spookyclitypes "spooky/internal/cli/types"
	spookylogging "spooky/internal/logging"
)

// Manager implements FlagsManager interface
type Manager struct {
	config       *spookyclitypes.FlagsConfig
	globalFlags  map[string]interface{}
	commandFlags map[string]map[string]interface{}
	parser       FlagsParser
	logger       spookylogging.Logger
	mutex        sync.RWMutex
}

// NewManager creates a new flags manager
func NewManager(
	config *spookyclitypes.FlagsConfig,
	parser FlagsParser,
	logger spookylogging.Logger,
) *Manager {
	return &Manager{
		config:       config,
		globalFlags:  make(map[string]interface{}),
		commandFlags: make(map[string]map[string]interface{}),
		parser:       parser,
		logger:       logger,
	}
}

// SetGlobalFlags sets global flags
func (m *Manager) SetGlobalFlags(flags map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Validate flags
	if err := m.parser.ValidateFlags(flags); err != nil {
		return fmt.Errorf("invalid global flags: %w", err)
	}

	// Set flags
	for name, value := range flags {
		m.globalFlags[name] = value
	}

	m.logger.Info("Global flags set", spookylogging.Int("count", len(flags)))
	return nil
}

// SetCommandFlags sets flags for a specific command
func (m *Manager) SetCommandFlags(commandName string, flags map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if commandName == "" {
		return fmt.Errorf("command name cannot be empty")
	}

	// Validate flags
	if err := m.parser.ValidateFlags(flags); err != nil {
		return fmt.Errorf("invalid command flags for %s: %w", commandName, err)
	}

	// Initialize command flags map if it doesn't exist
	if m.commandFlags[commandName] == nil {
		m.commandFlags[commandName] = make(map[string]interface{})
	}

	// Set flags
	for name, value := range flags {
		m.commandFlags[commandName][name] = value
	}

	m.logger.Info("Command flags set",
		spookylogging.String("command", commandName),
		spookylogging.Int("count", len(flags)))
	return nil
}

// GetGlobalFlags gets global flags
func (m *Manager) GetGlobalFlags() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	flags := make(map[string]interface{})
	for name, value := range m.globalFlags {
		flags[name] = value
	}

	return flags
}

// GetCommandFlags gets flags for a specific command
func (m *Manager) GetCommandFlags(commandName string) map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if commandName == "" {
		return nil
	}

	flags := make(map[string]interface{})
	if commandFlags, exists := m.commandFlags[commandName]; exists {
		for name, value := range commandFlags {
			flags[name] = value
		}
	}

	return flags
}
