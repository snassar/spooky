package commands

import (
	"spooky/internal/logging"

	"github.com/spf13/cobra"
)

// Global flags
var (
	configFile string
	quiet      bool
	verbose    bool
	logLevel   string
)

// RootCmd represents the root command of the spooky CLI.
var RootCmd = &cobra.Command{
	Use:   "spooky",
	Short: "Automation and configuration management tool",
	Long: `spooky is a Go-based automation and configuration management tool 
that focuses on creating self-contained binaries with embedded 
configuration schemas and validation rules.

It provides automation capabilities similar to Ansible, with embedded
schemas for validation and configuration management.`,
	Aliases: []string{"spooky"},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return RootCmd.Execute()
}

// GetConfigFile returns the config file path, considering the global flag
func GetConfigFile() string {
	return configFile
}

// IsCustomConfig returns true if a custom config file was specified
func IsCustomConfig() bool {
	return configFile != ""
}

// GetLogLevel returns the effective log level based on flags
func GetLogLevel() string {
	if quiet {
		return logging.LevelError
	}
	if verbose {
		return logging.LevelDebug
	}
	if logLevel != "" {
		return logLevel
	}
	return logging.LevelInfo // default
}

// IsQuiet returns true if quiet mode is enabled
func IsQuiet() bool {
	return quiet
}

// IsVerbose returns true if verbose mode is enabled
func IsVerbose() bool {
	return verbose
}

func init() {
	// Global flags
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "path to custom configuration file (overrides default spooky.hcl and embedded default)")
	RootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress all logging output except errors")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging (debug level)")
	RootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "", "set logging level (debug, info, warn, error, fatal)")

	// Mark flags as mutually exclusive
	RootCmd.MarkFlagsMutuallyExclusive("quiet", "verbose", "log-level")

	// Add alias commands
	addAliasCommands()
}

// addAliasCommands adds alias commands for common operations
func addAliasCommands() {
	// Alias for "spooky run" -> "spooky actions run"
	runAliasCmd := &cobra.Command{
		Use:   "run",
		Short: "Alias for 'spooky actions run'",
		RunE:  runAction,
	}

	// Add flags to run alias (same as runActionCmd)
	runAliasCmd.Flags().StringSliceVarP(&actionTargets, "targets", "t", nil, "override action targets (comma-separated list)")
	runAliasCmd.Flags().BoolVarP(&actionDryRun, "dry-run", "n", false, "show what would be executed without actually running")
	runAliasCmd.Flags().IntVar(&actionTimeout, "timeout", 0, "override action timeout in seconds")
	runAliasCmd.Flags().IntVar(&actionRetries, "retries", -1, "override action retry count (-1 to use action default)")

	// Alias for "spooky ping" -> "spooky machines ping"
	pingAliasCmd := &cobra.Command{
		Use:   "ping",
		Short: "Alias for 'spooky machines ping'",
		RunE:  pingMachines,
	}

	// Add flags to ping alias (same as pingCmd)
	pingAliasCmd.Flags().StringSliceVarP(&pingTargets, "targets", "t", nil, "specific machines to ping (comma-separated list)")
	pingAliasCmd.Flags().IntVar(&pingTimeout, "timeout", 30, "ping timeout in seconds")
	pingAliasCmd.Flags().BoolVar(&pingAuth, "auth", false, "attempt authentication during ping")
	pingAliasCmd.Flags().StringVar(&pingOutput, "output", "text", "output format (text, json)")

	// Add alias commands to root
	RootCmd.AddCommand(runAliasCmd)
	RootCmd.AddCommand(pingAliasCmd)
}
