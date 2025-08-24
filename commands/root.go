package commands

import (
	"github.com/spf13/cobra"
)

// Global flags
var (
	configFile string
	quiet      bool
	verbose    bool
	logLevel   string
)

var RootCmd = &cobra.Command{
	Use:   "spooky",
	Short: "Automation and configuration management tool",
	Long: `spooky is a Go-based automation and configuration management tool 
that focuses on creating self-contained binaries with embedded 
configuration schemas and validation rules.

It provides automation capabilities similar to Ansible, with embedded
schemas for validation and configuration management.`,
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
		return "error"
	}
	if verbose {
		return "debug"
	}
	if logLevel != "" {
		return logLevel
	}
	return "info" // default
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
}
