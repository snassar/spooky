package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "spooky",
	Short: "Automation and configuration management tool",
	Long: `spooky is a Go-based automation and configuration management tool 
that focuses on creating self-contained binaries with embedded 
configuration schemas and validation rules.

It provides automation capabilities similar to Ansible, with embedded
schemas for validation and configuration management.`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show spooky version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("spooky v0.1.0")
		fmt.Println("Automation and configuration management tool")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
