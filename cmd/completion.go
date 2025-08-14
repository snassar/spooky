// Package cmd provides command implementations for spooky CLI.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for spooky.

This command generates completion scripts for bash, zsh, and fish shells.
The generated scripts provide tab completion for all spooky commands and flags.

Examples:
  spooky completion bash > ~/.local/share/bash-completion/completions/spooky
  spooky completion zsh > ~/.zsh/completions/_spooky
  spooky completion fish > ~/.config/fish/completions/spooky.fish`,
	ValidArgs: []string{"bash", "zsh", "fish"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		shell := args[0]

		switch shell {
		case "bash":
			return cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			return cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			return cmd.Root().GenFishCompletion(os.Stdout, true)
		default:
			return cmd.Help()
		}
	},
}
