package commands

import (
	"github.com/spf13/cobra"
)

// Global flags
var (
	configFile string
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

func init() {
	// Global flags
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "path to custom configuration file (overrides default spooky.hcl and embedded default)")
}
