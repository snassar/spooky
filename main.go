package main

import (
	"fmt"
	"os"

	"spooky/internal/cli"
	"spooky/internal/logging"

	"github.com/spf13/cobra"
)

func main() {
	// coverage-ignore: main function is entry point, tested via integration tests
	var rootCmd = &cobra.Command{
		Use:   "spooky",
		Short: "Spooky is a server configuration and automation tool",
		Long: `Spooky is a powerful server configuration and automation tool that allows you to:
- Connect to multiple servers via SSH
- Execute commands and scripts from HCL2 configuration files
- Manage server operations in a declarative way
- Support for parallel execution and error handling
- Collect and manage server facts
- Use templates for dynamic configuration`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Configure logger based on global flags
			config := cli.GetGlobalConfig()
			logging.ConfigureLogger(config.LogLevel, "json", config.LogFile, config.Quiet, config.Verbose)

			logger := logging.GetLogger()
			logger.Info("Starting spooky application",
				logging.String("version", "1.0.0"), // TODO: Add version from build info
			)
		},
	}

	// Add global flags
	cli.AddGlobalFlags(rootCmd)

	// Initialize CLI commands
	cli.InitCommands()

	// Add subcommands
	rootCmd.AddCommand(cli.ExecuteCmd)
	rootCmd.AddCommand(cli.ValidateCmd)
	rootCmd.AddCommand(cli.ListCmd)
	rootCmd.AddCommand(cli.FactsCmd)

	if err := rootCmd.Execute(); err != nil {
		// Configure logger for error output if not already configured
		config := cli.GetGlobalConfig()
		logging.ConfigureLogger(config.LogLevel, "json", config.LogFile, config.Quiet, config.Verbose)

		logger := logging.GetLogger()
		logger.Error("Application execution failed", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Configure logger for success message if not already configured
	config := cli.GetGlobalConfig()
	logging.ConfigureLogger(config.LogLevel, "json", config.LogFile, config.Quiet, config.Verbose)

	logger := logging.GetLogger()
	logger.Info("Application completed successfully")
}
