package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "spooky-schemas",
	Short: "Manage and inspect embedded schemas",
	Long: `spooky-schemas is a utility tool for working with embedded schemas,
validation rules, and test data in the spooky project.

It provides functionality to list, inspect, and validate schema definitions.`,
}

var schemasCmd = &cobra.Command{
	Use:   "schemas",
	Short: "Show embedded schemas summary",
	Long:  `Display a comprehensive summary of all embedded schemas, validation rules, and test data.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("⚠️  Schema embedder has been replaced with struct-based schemas")
		fmt.Println("   Use 'spooky project validate' to validate configuration files")
		fmt.Println("   Use 'spooky project init' to generate configuration templates")
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all embedded schemas and test data",
	Long:  `List all available embedded schemas, validation rules, and test data files.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("⚠️  Schema embedder has been replaced with struct-based schemas")
		fmt.Println("   Use 'spooky project validate' to validate configuration files")
		fmt.Println("   Use 'spooky project init' to generate configuration templates")
	},
}

var showCmd = &cobra.Command{
	Use:   "show [schema|rule|test] [name]",
	Short: "Show content of a specific schema, rule, or test data",
	Long:  `Display the content of a specific embedded schema, validation rule, or test data file.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("⚠️  Schema embedder has been replaced with struct-based schemas")
		fmt.Println("   Use 'spooky project validate' to validate configuration files")
		fmt.Println("   Use 'spooky project init' to generate configuration templates")
	},
}

func init() {
	rootCmd.AddCommand(schemasCmd)
	schemasCmd.AddCommand(listCmd)
	schemasCmd.AddCommand(showCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
