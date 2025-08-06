package completion

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	spookyclitypes "spooky/internal/cli/types"
	spookylogging "spooky/internal/logging"
)

// Manager implements CompletionManager interface
type Manager struct {
	config    *spookyclitypes.CompletionConfig
	generator CompletionGenerator
	rootCmd   *cobra.Command
	logger    spookylogging.Logger
}

// NewManager creates a new completion manager
func NewManager(
	config *spookyclitypes.CompletionConfig,
	generator CompletionGenerator,
	rootCmd *cobra.Command,
	logger spookylogging.Logger,
) *Manager {
	return &Manager{
		config:    config,
		generator: generator,
		rootCmd:   rootCmd,
		logger:    logger,
	}
}

// GenerateCompletion generates completion for a shell
func (m *Manager) GenerateCompletion(shell string) (string, error) {
	if m.rootCmd == nil {
		return "", fmt.Errorf("root command not set")
	}

	var buf strings.Builder

	switch shell {
	case "bash":
		err := m.rootCmd.GenBashCompletion(&buf)
		if err != nil {
			return "", fmt.Errorf("failed to generate bash completion: %w", err)
		}
	case "zsh":
		err := m.rootCmd.GenZshCompletion(&buf)
		if err != nil {
			return "", fmt.Errorf("failed to generate zsh completion: %w", err)
		}
	case "fish":
		err := m.rootCmd.GenFishCompletion(&buf, true)
		if err != nil {
			return "", fmt.Errorf("failed to generate fish completion: %w", err)
		}
	case "powershell":
		err := m.rootCmd.GenPowerShellCompletion(&buf)
		if err != nil {
			return "", fmt.Errorf("failed to generate powershell completion: %w", err)
		}
	default:
		return "", fmt.Errorf("unsupported shell: %s", shell)
	}

	return buf.String(), nil
}

// GenerateCompletionFile generates completion file for a shell
func (m *Manager) GenerateCompletionFile(shell, outputPath string) error {
	completion, err := m.GenerateCompletion(shell)
	if err != nil {
		return fmt.Errorf("failed to generate completion: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(completion), 0o600); err != nil {
		return fmt.Errorf("failed to write completion file: %w", err)
	}

	m.logger.Info("Completion file generated",
		spookylogging.String("shell", shell),
		spookylogging.String("path", outputPath))
	return nil
}

// GetSupportedShells returns supported shells
func (m *Manager) GetSupportedShells() []string {
	return []string{"bash", "zsh", "fish", "powershell"}
}
