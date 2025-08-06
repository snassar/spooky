package main

import (
	"fmt"
	"os"

	spookycli "spooky/internal/cli"
	spookylogging "spooky/internal/logging"

	"github.com/spf13/cobra"
)

// Version information - set at build time via ldflags
var (
	version = "dev"
	commit  = "unknown"
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
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			// Configure logger based on global flags
			config := spookycli.GetGlobalConfig()
			spookylogging.ConfigureLogger(config.LogLevel, "json", config.LogFile, config.Quiet, config.Verbose)

			logger := spookylogging.GetLogger()
			logger.Info("Starting spooky application",
				spookylogging.String("version", fmt.Sprintf("%s-%s", version, commit)),
			)
		},
	}

	// Add global flags
	spookycli.AddGlobalFlags(rootCmd)

	// Initialize CLI commands
	spookycli.InitCommands()

	// Add subcommands

	rootCmd.AddCommand(spookycli.InitCmd)
	rootCmd.AddCommand(spookycli.ValidateCmd)
	rootCmd.AddCommand(spookycli.ListCmd)
	rootCmd.AddCommand(spookycli.ListMachinesCmd)
	rootCmd.AddCommand(spookycli.ListActionsCmd)
	rootCmd.AddCommand(spookycli.ListTemplatesCmd)
	rootCmd.AddCommand(spookycli.ListFactsCmd)
	rootCmd.AddCommand(spookycli.GatherFactsCmd)
	rootCmd.AddCommand(spookycli.RenderTemplateCmd)
	rootCmd.AddCommand(spookycli.ValidateTemplateCmd)

	if err := rootCmd.Execute(); err != nil {
		// Configure logger for error output if not already configured
		config := spookycli.GetGlobalConfig()
		spookylogging.ConfigureLogger(config.LogLevel, "json", config.LogFile, config.Quiet, config.Verbose)

		logger := spookylogging.GetLogger()
		logger.Error("Application execution failed", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Configure logger for success message if not already configured
	config := spookycli.GetGlobalConfig()
	spookylogging.ConfigureLogger(config.LogLevel, "json", config.LogFile, config.Quiet, config.Verbose)

	logger := spookylogging.GetLogger()
	logger.Info("Application completed successfully")
}
