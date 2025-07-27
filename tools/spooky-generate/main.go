package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	generateType  string
	actionsCount  int
	machinesCount int
	outputPath    string
	projectName   string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "spooky-generate",
		Short: "Generate spooky configuration files for testing",
		Long: `Generate spooky configuration files for testing purposes.
		
This tool generates valid HCL configuration files that can be used in spooky projects.
It can generate both actions and inventory files with intertwined tags for realistic testing.`,
		RunE: runGenerate,
	}

	// Add flags
	rootCmd.Flags().StringVarP(&generateType, "generate", "g", "all", "What to generate: all, inventory-only")
	rootCmd.Flags().IntVarP(&actionsCount, "actions", "a", 25, "Number of actions to generate")
	rootCmd.Flags().IntVarP(&machinesCount, "machines", "m", 10, "Number of machines to generate")
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output directory for spooky project (default: ./spooky-project)")
	rootCmd.Flags().StringVarP(&projectName, "name", "n", "generated-project", "Name for the spooky project")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runGenerate(cmd *cobra.Command, args []string) error {
	generator := NewGenerator()

	switch generateType {
	case "all":
		return generator.GenerateProject()
	case "inventory-only":
		return generator.GenerateInventoryOnly()
	default:
		return fmt.Errorf("invalid generate type: %s. Use 'all' or 'inventory-only'", generateType)
	}
}
