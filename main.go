package main

import (
	"fmt"
	"os"

	"spooky/internal/cli"

	"github.com/spf13/cobra"
)

func main() {
	// coverage-ignore: main function is entry point, tested via integration tests
	var rootCmd = &cobra.Command{
		Use:   "spooky",
		Short: "Spooky is a SSH automation tool that runs actions from HCL2 configuration files",
		Long: `Spooky is a powerful SSH automation tool that allows you to:
- Connect to multiple servers via SSH
- Execute commands and scripts from HCL2 configuration files
- Manage server operations in a declarative way
- Support for parallel execution and error handling`,
	}

	// Initialize CLI commands
	cli.InitCommands()

	// Add subcommands
	rootCmd.AddCommand(cli.ExecuteCmd)
	rootCmd.AddCommand(cli.ValidateCmd)
	rootCmd.AddCommand(cli.ListCmd)
	rootCmd.AddCommand(cli.KeygenCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
