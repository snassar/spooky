package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"spooky/internal/utilities"
)

var rootCmd = &cobra.Command{
	Use:   "spooky-hcl",
	Short: "Validate and manipulate HCL files",
	Long: `spooky-hcl is a utility tool for working with HCL (HashiCorp Configuration Language) files.

It provides functionality to validate HCL syntax, check file integrity, and perform
various HCL-related operations.`,
}

var validateCmd = &cobra.Command{
	Use:   "validate [path]",
	Short: "Validate HCL files or directories",
	Long: `Validate HCL files or directories. Automatically detects whether the input is a file or directory.

Examples:
  spooky-hcl validate project.hcl          # Validate single file
  spooky-hcl validate ./schemas/           # Validate all .hcl files in directory
  spooky-hcl validate < input.hcl          # Validate content from stdin`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		quick, _ := cmd.Flags().GetBool("quick")

		// Handle stdin input
		if len(args) == 0 {
			validateStdin(quick)
			return
		}

		path := args[0]
		validatePath(path, quick)
	},
}

func validateStdin(quick bool) {
	// Read from stdin
	scanner := bufio.NewScanner(os.Stdin)
	var content string
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", errors.WithStack(err))
		os.Exit(1)
	}

	if quick {
		// Quick validation
		if utilities.IsValidHCL(content) {
			fmt.Println("✅ Valid HCL content")
		} else {
			fmt.Fprintf(os.Stderr, "❌ Invalid HCL content\n")
			os.Exit(1)
		}
	} else {
		// Detailed validation
		validator := utilities.NewHCLValidator()
		result, err := validator.ValidateContent(content, "stdin")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating content: %v\n", errors.WithStack(err))
			os.Exit(1)
		}

		fmt.Print(utilities.FormatValidationResult(result, "stdin"))
		if !result.IsValid {
			os.Exit(1)
		}
	}
}

func validatePath(path string, quick bool) {
	// Check if path exists
	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error accessing path: %v\n", errors.WithStack(err))
		os.Exit(1)
	}

	if info.IsDir() {
		validateDirectory(path, quick)
	} else {
		validateFile(path, quick)
	}
}

func validateFile(filePath string, quick bool) {
	if quick {
		// Quick validation
		isValid, err := utilities.ValidateHCLFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", errors.WithStack(err))
			os.Exit(1)
		}

		if isValid {
			fmt.Printf("✅ Valid HCL file: %s\n", filePath)
		} else {
			fmt.Fprintf(os.Stderr, "❌ Invalid HCL file: %s\n", filePath)
			os.Exit(1)
		}
	} else {
		// Detailed validation
		validator := utilities.NewHCLValidator()
		result, err := validator.ValidateFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating file: %v\n", errors.WithStack(err))
			os.Exit(1)
		}

		fmt.Print(utilities.FormatValidationResult(result, filePath))
		if !result.IsValid {
			os.Exit(1)
		}
	}
}

func validateDirectory(dirPath string, quick bool) {
	validator := utilities.NewHCLValidator()
	results, err := validator.ValidateDirectory(dirPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error validating directory: %v\n", errors.WithStack(err))
		os.Exit(1)
	}

	if quick {
		// Quick validation - just show summary
		validCount := 0
		invalidCount := 0

		for filePath, result := range results {
			if result.IsValid {
				validCount++
			} else {
				invalidCount++
				fmt.Printf("❌ %s\n", filePath)
			}
		}

		fmt.Printf("📁 Directory: %s\n", dirPath)
		fmt.Printf("✅ Valid: %d, ❌ Invalid: %d, 📊 Total: %d\n", validCount, invalidCount, len(results))

		if invalidCount > 0 {
			os.Exit(1)
		}
	} else {
		// Detailed validation
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

		if invalidCount > 0 {
			os.Exit(1)
		}
	}
}

func init() {
	validateCmd.Flags().BoolP("quick", "q", false, "Quick validation (exit code only)")
	rootCmd.AddCommand(validateCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
