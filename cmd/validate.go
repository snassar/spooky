package cmd

import (
	"fmt"
	"os"

	"spooky/internal/utilities"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate HCL files and content",
	Long: `The validate command provides functionality to validate HCL files and content.
It checks for proper HCL syntax and provides detailed error reporting.`,
}

var validateFileCmd = &cobra.Command{
	Use:   "file [filepath]",
	Short: "Validate a specific HCL file",
	Long:  `Validate a specific HCL file and display detailed validation results.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		// Create validator
		validator := utilities.NewHCLValidator()

		// Validate the file
		result, err := validator.ValidateFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating file: %v\n", err)
			os.Exit(1)
		}

		// Display results
		fmt.Print(utilities.FormatValidationResult(result, filePath))
	},
}

var validateContentCmd = &cobra.Command{
	Use:   "content [content]",
	Short: "Validate HCL content from string",
	Long:  `Validate HCL content provided as a string argument.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		content := args[0]

		// Create validator
		validator := utilities.NewHCLValidator()

		// Validate the content
		result, err := validator.ValidateContent(content, "input")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating content: %v\n", err)
			os.Exit(1)
		}

		// Display results
		fmt.Print(utilities.FormatValidationResult(result, "input"))
	},
}

var validateDirectoryCmd = &cobra.Command{
	Use:   "directory [dirpath]",
	Short: "Validate all HCL files in a directory",
	Long:  `Validate all HCL files in a directory and display validation results for each.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dirPath := args[0]

		// Check if directory exists
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Directory does not exist: %s\n", dirPath)
			os.Exit(1)
		}

		// Create validator
		validator := utilities.NewHCLValidator()

		// Validate the directory
		results, err := validator.ValidateDirectory(dirPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating directory: %v\n", err)
			os.Exit(1)
		}

		// Display results
		fmt.Printf("📁 Validation Results for Directory: %s\n", dirPath)
		fmt.Printf("📊 Total Files: %d\n\n", len(results))

		validCount := 0
		invalidCount := 0

		for filePath, result := range results {
			if result.IsValid {
				validCount++
				fmt.Printf("✅ %s\n", filePath)
			} else {
				invalidCount++
				fmt.Printf("❌ %s\n", filePath)
			}
		}

		fmt.Printf("\n📈 Summary:\n")
		fmt.Printf("  Valid files: %d\n", validCount)
		fmt.Printf("  Invalid files: %d\n", invalidCount)
		fmt.Printf("  Total files: %d\n", len(results))

		// Exit with error if any files are invalid
		if invalidCount > 0 {
			os.Exit(1)
		}
	},
}

var validateQuickCmd = &cobra.Command{
	Use:   "quick [filepath]",
	Short: "Quick validation check (returns exit code only)",
	Long: `Perform a quick validation check that only returns an exit code.
Useful for scripts and automation.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		// Quick validation
		isValid, err := utilities.ValidateHCLFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if !isValid {
			fmt.Fprintf(os.Stderr, "File is not valid HCL: %s\n", filePath)
			os.Exit(1)
		}

		// File is valid, exit with success
		fmt.Printf("File is valid HCL: %s\n", filePath)
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.AddCommand(validateFileCmd)
	validateCmd.AddCommand(validateContentCmd)
	validateCmd.AddCommand(validateDirectoryCmd)
	validateCmd.AddCommand(validateQuickCmd)
}
