package completion

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// Generator implements completion generation for different shells
type Generator struct {
	rootCommand *cobra.Command
}

// NewGenerator creates a new completion generator
func NewGenerator(rootCommand *cobra.Command) *Generator {
	return &Generator{
		rootCommand: rootCommand,
	}
}

// GenerateBashCompletion generates bash completion script
func (g *Generator) GenerateBashCompletion() (string, error) {
	if g.rootCommand == nil {
		return "", fmt.Errorf("root command is nil")
	}

	var buf bytes.Buffer
	err := g.rootCommand.GenBashCompletion(&buf)
	if err != nil {
		return "", fmt.Errorf("failed to generate bash completion: %w", err)
	}

	return buf.String(), nil
}

// GenerateZshCompletion generates zsh completion script
func (g *Generator) GenerateZshCompletion() (string, error) {
	if g.rootCommand == nil {
		return "", fmt.Errorf("root command is nil")
	}

	var buf bytes.Buffer
	err := g.rootCommand.GenZshCompletion(&buf)
	if err != nil {
		return "", fmt.Errorf("failed to generate zsh completion: %w", err)
	}

	return buf.String(), nil
}

// GenerateFishCompletion generates fish completion script
func (g *Generator) GenerateFishCompletion() (string, error) {
	if g.rootCommand == nil {
		return "", fmt.Errorf("root command is nil")
	}

	var buf bytes.Buffer
	err := g.rootCommand.GenFishCompletion(&buf, true)
	if err != nil {
		return "", fmt.Errorf("failed to generate fish completion: %w", err)
	}

	return buf.String(), nil
}

// GeneratePowerShellCompletion generates PowerShell completion script
func (g *Generator) GeneratePowerShellCompletion() (string, error) {
	if g.rootCommand == nil {
		return "", fmt.Errorf("root command is nil")
	}

	var buf bytes.Buffer
	err := g.rootCommand.GenPowerShellCompletion(&buf)
	if err != nil {
		return "", fmt.Errorf("failed to generate PowerShell completion: %w", err)
	}

	return buf.String(), nil
}

// GenerateCompletionForShell generates completion for a specific shell
func (g *Generator) GenerateCompletionForShell(shell string) (string, error) {
	switch strings.ToLower(shell) {
	case "bash":
		return g.GenerateBashCompletion()
	case "zsh":
		return g.GenerateZshCompletion()
	case "fish":
		return g.GenerateFishCompletion()
	case "powershell":
		return g.GeneratePowerShellCompletion()
	default:
		return "", fmt.Errorf("unsupported shell: %s", shell)
	}
}

// GetSupportedShells returns the list of supported shells
func (g *Generator) GetSupportedShells() []string {
	return []string{"bash", "zsh", "fish", "powershell"}
}

// ValidateShell validates if a shell is supported
func (g *Generator) ValidateShell(shell string) error {
	supportedShells := g.GetSupportedShells()
	for _, supported := range supportedShells {
		if strings.ToLower(shell) == supported {
			return nil
		}
	}
	return fmt.Errorf("unsupported shell: %s. Supported shells: %v", shell, supportedShells)
}
