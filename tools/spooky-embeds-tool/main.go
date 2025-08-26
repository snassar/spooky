package main

import (
	"fmt"
	"os"

	"spooky/internal/utilities"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "spooky-embeds",
	Short: "Embedded files management tool for spooky",
	Long: `spooky-embeds demonstrates how spooky handles embedded files:
- Default configuration files
- HCL validation of embedded files
- File embedding and retrieval`,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all embedded files",
	Run: func(cmd *cobra.Command, args []string) {
		embedder, err := utilities.NewFileEmbedder()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file embedder: %v\n", err)
			os.Exit(1)
		}

		embedder.PrintFileSummary()
	},
}

var showCmd = &cobra.Command{
	Use:   "show [file-name]",
	Short: "Show content of an embedded file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fileName := args[0]

		embedder, err := utilities.NewFileEmbedder()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file embedder: %v\n", err)
			os.Exit(1)
		}

		content, exists := embedder.GetFile(fileName)
		if !exists {
			fmt.Fprintf(os.Stderr, "File '%s' not found\n", fileName)
			fmt.Printf("Available files:\n")
			for _, name := range embedder.ListFiles() {
				fmt.Printf("  - %s\n", name)
			}
			os.Exit(1)
		}

		fmt.Printf("📄 Content of '%s':\n", fileName)
		fmt.Printf("Size: %d bytes\n\n", len(content))
		fmt.Println(content)
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate [file-name]",
	Short: "Validate HCL syntax of an embedded file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fileName := args[0]

		embedder, err := utilities.NewFileEmbedder()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file embedder: %v\n", err)
			os.Exit(1)
		}

		content, exists := embedder.GetFile(fileName)
		if !exists {
			fmt.Fprintf(os.Stderr, "File '%s' not found\n", fileName)
			os.Exit(1)
		}

		// Validate HCL
		validator := utilities.NewHCLValidator()
		result, err := validator.ValidateContent(content, fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating file: %v\n", err)
			os.Exit(1)
		}

		if result.IsValid {
			fmt.Printf("✅ File '%s' is valid HCL\n", fileName)
			fmt.Printf("📊 File Size: %d bytes\n", result.FileSize)
			fmt.Printf("📦 Block Count: %d\n", result.BlockCount)
		} else {
			fmt.Printf("❌ File '%s' has HCL errors:\n", fileName)
			for _, err := range result.Errors {
				fmt.Printf("  - %s\n", err.Message)
			}
			os.Exit(1)
		}
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Demonstrate default configuration functionality",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("⚙️  Default Configuration Demo:")

		// Get embedded default config
		embedder, err := utilities.NewFileEmbedder()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file embedder: %v\n", err)
			os.Exit(1)
		}

		defaultConfig, err := embedder.GetDefaultConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting default config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("📄 Default Configuration:\n")
		fmt.Printf("Size: %d bytes\n\n", len(defaultConfig))
		fmt.Println(defaultConfig)

		// Validate the config
		fmt.Println("\n🔍 Validating default configuration...")
		validator := utilities.NewHCLValidator()
		result, err := validator.ValidateContent(defaultConfig, "default-config.hcl")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating config: %v\n", err)
			os.Exit(1)
		}

		if result.IsValid {
			fmt.Println("✅ Default configuration is valid HCL!")
		} else {
			fmt.Println("❌ Default configuration has errors:")
			for _, err := range result.Errors {
				fmt.Printf("  - %s\n", err.Message)
			}
			os.Exit(1)
		}

		// Test config manager integration
		fmt.Println("\n🔧 Testing config manager integration...")
		configManager, err := utilities.NewConfigManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config manager: %v\n", err)
			os.Exit(1)
		}

		if !configManager.ConfigExists() {
			fmt.Println("📝 Creating default configuration file...")
			if err := configManager.CreateDefaultConfig(); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating default config: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ Default configuration created at: %s\n", configManager.GetConfigPath())
		} else {
			fmt.Printf("✅ Configuration already exists at: %s\n", configManager.GetConfigPath())
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(configCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
