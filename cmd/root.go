// Package cmd provides command implementations for spooky CLI.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version information
var (
	Version   = "0.1.0"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "spooky",
	Short: "A powerful automation and orchestration tool",
	Long: `spooky is a powerful automation and orchestration tool built with Go.

It provides declarative configuration, parallel execution capabilities, 
and intelligent fact-driven decision making for heterogeneous environments.

Examples:
  spooky project init my-project
  spooky project validate my-project
  spooky --version
  spooky --help`,
	Version: Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Set version template
	RootCmd.SetVersionTemplate(fmt.Sprintf(`spooky version %s
Build time: %s
Git commit: %s
`, Version, BuildTime, GitCommit))
}
