package cmd

import (
	"fmt"
	"os"
	"strings"

	"spooky/internal/schemas"

	"github.com/spf13/cobra"
)

var schemasCmd = &cobra.Command{
	Use:   "schemas",
	Short: "Manage and inspect embedded schemas",
	Long: `The schemas command provides functionality to work with embedded schemas,
validation rules, and test data. It allows you to list, inspect, and validate
schema definitions.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize the schema embedder
		embedder, err := schemas.NewSchemaEmbedder()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing schema embedder: %v\n", err)
			os.Exit(1)
		}

		// Print schema summary
		embedder.PrintSchemaSummary()
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all embedded schemas and test data",
	Long:  `List all available embedded schemas, validation rules, and test data files.`,
	Run: func(cmd *cobra.Command, args []string) {
		embedder, err := schemas.NewSchemaEmbedder()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing schema embedder: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("📋 Embedded Schemas:")
		schemas := embedder.ListSchemas()
		for i, name := range schemas {
			fmt.Printf("  %d. %s\n", i+1, name)
		}

		fmt.Println("\n🔍 Validation Rules:")
		rules := embedder.ListValidationRules()
		for i, name := range rules {
			fmt.Printf("  %d. %s\n", i+1, name)
		}

		fmt.Println("\n🧪 Test Data:")
		testdata := embedder.ListTestData()
		for i, name := range testdata {
			fmt.Printf("  %d. %s\n", i+1, name)
		}
	},
}

var showCmd = &cobra.Command{
	Use:   "show [schema|rule|test] [name]",
	Short: "Show content of a specific schema, rule, or test data",
	Long:  `Display the content of a specific embedded schema, validation rule, or test data file.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		embedder, err := schemas.NewSchemaEmbedder()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing schema embedder: %v\n", err)
			os.Exit(1)
		}

		itemType := args[0]
		name := args[1]

		var content string
		var exists bool

		switch itemType {
		case "schema":
			content, exists = embedder.GetSchema(name)
		case "rule":
			content, exists = embedder.GetValidationRules(name)
		case "test":
			content, exists = embedder.GetTestData(name)
		default:
			fmt.Fprintf(os.Stderr, "Invalid item type: %s. Use 'schema', 'rule', or 'test'\n", itemType)
			os.Exit(1)
		}

		if !exists {
			fmt.Fprintf(os.Stderr, "%s '%s' not found\n", itemType, name)
			os.Exit(1)
		}

		fmt.Printf("📄 %s: %s\n", itemType, name)
		fmt.Println("─" + strings.Repeat("─", 50))
		fmt.Println(content)
	},
}

func init() {
	rootCmd.AddCommand(schemasCmd)
	schemasCmd.AddCommand(listCmd)
	schemasCmd.AddCommand(showCmd)
}
