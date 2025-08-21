package commands

import (
	"github.com/spf13/cobra"
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
