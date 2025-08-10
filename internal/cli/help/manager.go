package help

import (
	"fmt"
	spookyinterfaces "spooky/internal/interfaces"
	spookytypeslogging "spooky/internal/types/logging"
	"strings"

	spookyclitypes "spooky/internal/types/cli"

	"github.com/spf13/cobra"
)

// Manager implements HelpManager interface
type Manager struct {
	config   *spookyclitypes.HelpConfig
	renderer spookyinterfaces.HelpRenderer
	rootCmd  *cobra.Command
	logger   spookytypeslogging.Logger
}

// NewManager creates a new help manager
func NewManager(
	config *spookyclitypes.HelpConfig,
	renderer spookyinterfaces.HelpRenderer,
	rootCmd *cobra.Command,
	logger spookytypeslogging.Logger,
) *Manager {
	return &Manager{
		config:   config,
		renderer: renderer,
		rootCmd:  rootCmd,
		logger:   logger,
	}
}

// ShowHelp shows help for a command
func (m *Manager) ShowHelp(commandName string) (string, error) {
	if m.rootCmd == nil {
		return "", fmt.Errorf("root command not set")
	}

	if commandName == "" {
		// Show root command help
		return m.rootCmd.Long, nil
	}

	// Find the command
	cmd := m.findCommand(m.rootCmd, commandName)
	if cmd == nil {
		return "", fmt.Errorf("command not found: %s", commandName)
	}

	return cmd.Long, nil
}

// ShowUsage shows usage for a command
func (m *Manager) ShowUsage(commandName string) (string, error) {
	if m.rootCmd == nil {
		return "", fmt.Errorf("root command not set")
	}

	if commandName == "" {
		// Show root command usage
		return m.rootCmd.UsageString(), nil
	}

	// Find the command
	cmd := m.findCommand(m.rootCmd, commandName)
	if cmd == nil {
		return "", fmt.Errorf("command not found: %s", commandName)
	}

	return cmd.UsageString(), nil
}

// ShowExamples shows examples for a command
func (m *Manager) ShowExamples(commandName string) (string, error) {
	if m.rootCmd == nil {
		return "", fmt.Errorf("root command not set")
	}

	if commandName == "" {
		return "", fmt.Errorf("command name required for examples")
	}

	// Find the command
	cmd := m.findCommand(m.rootCmd, commandName)
	if cmd == nil {
		return "", fmt.Errorf("command not found: %s", commandName)
	}

	if cmd.Example == "" {
		return "No examples available for this command", nil
	}

	var buf strings.Builder
	buf.WriteString("Examples:\n")
	buf.WriteString(cmd.Example)

	return buf.String(), nil
}

// Helper method to find a command by name
func (m *Manager) findCommand(root *cobra.Command, name string) *cobra.Command {
	if root.Name() == name {
		return root
	}

	for _, cmd := range root.Commands() {
		if cmd.Name() == name {
			return cmd
		}
		// Recursively search subcommands
		if found := m.findCommand(cmd, name); found != nil {
			return found
		}
	}

	return nil
}
