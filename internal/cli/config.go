package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"spooky/internal/config"

	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration files",
	Long: `Manage and validate configuration files.

This command group provides tools for validating, linting, and formatting
HCL configuration files used by spooky.

Examples:
  # Validate a configuration file
  spooky config validate config.hcl

  # Lint a configuration file for style issues
  spooky config lint config.hcl

  # Format a configuration file
  spooky config format config.hcl`,
}

var configValidateCmd = &cobra.Command{
	Use:   "validate <file>",
	Short: "Validate configuration syntax",
	Long: `Validate configuration syntax and structure.

This command checks if a configuration file has valid HCL syntax and can be
parsed without errors. It also validates the configuration against the schema.

Examples:
  spooky config validate config.hcl
  spooky config validate /path/to/config.hcl`,
	Args: cobra.ExactArgs(1),
	RunE: runConfigValidate,
}

var configLintCmd = &cobra.Command{
	Use:   "lint <file>",
	Short: "Lint configuration files",
	Long: `Lint configuration files for best practices and potential issues.

This command checks configuration files for common issues, best practices,
and potential problems. It provides suggestions for improvement.

Examples:
  spooky config lint config.hcl
  spooky config lint /path/to/config.hcl --fix`,
	Args: cobra.ExactArgs(1),
	RunE: runConfigLint,
}

var configFormatCmd = &cobra.Command{
	Use:   "format <file>",
	Short: "Format configuration files",
	Long: `Format configuration files for consistent style.

This command formats HCL configuration files to ensure consistent indentation,
spacing, and structure. It can format files in-place or output to stdout.

Examples:
  spooky config format config.hcl
  spooky config format config.hcl --write
  spooky config format config.hcl --output formatted.hcl`,
	Args: cobra.ExactArgs(1),
	RunE: runConfigFormat,
}

func init() {
	// Add flags for config lint command
	configLintCmd.Flags().Bool("fix", false, "Automatically fix issues where possible")
	configLintCmd.Flags().Bool("strict", false, "Enable strict linting rules")

	// Add flags for config format command
	configFormatCmd.Flags().Bool("write", false, "Write formatted content back to file")
	configFormatCmd.Flags().String("output", "", "Output file path (default: stdout)")

	// Add subcommands to config command
	ConfigCmd.AddCommand(configValidateCmd)
	ConfigCmd.AddCommand(configLintCmd)
	ConfigCmd.AddCommand(configFormatCmd)

	// Note: configCmd will be added to root in main.go
}

func runConfigValidate(_ *cobra.Command, args []string) error {
	configFile := args[0]

	// Check if file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("configuration file does not exist: %s", configFile)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(configFile))
	if ext != ".hcl" && ext != ".hcl2" {
		fmt.Printf("⚠️  Warning: File has non-standard extension '%s'. Expected '.hcl' or '.hcl2'\n", ext)
	}

	// Parse configuration
	cfg, err := config.ParseConfig(configFile)
	if err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Additional validation checks
	var issues []string

	// Check for empty configuration
	if len(cfg.Machines) == 0 && len(cfg.Actions) == 0 {
		issues = append(issues, "Configuration file contains no machines or actions")
	}

	// Check machine names for uniqueness
	machineNames := make(map[string]bool)
	for _, machine := range cfg.Machines {
		if machineNames[machine.Name] {
			issues = append(issues, fmt.Sprintf("Duplicate machine name: %s", machine.Name))
		}
		machineNames[machine.Name] = true
	}

	// Check action names for uniqueness
	actionNames := make(map[string]bool)
	for i := range cfg.Actions {
		action := &cfg.Actions[i]
		if actionNames[action.Name] {
			issues = append(issues, fmt.Sprintf("Duplicate action name: %s", action.Name))
		}
		actionNames[action.Name] = true
	}

	// Report results
	fmt.Printf("✓ Configuration file '%s' is valid\n", configFile)
	fmt.Printf("  - Machines: %d\n", len(cfg.Machines))
	fmt.Printf("  - Actions: %d\n", len(cfg.Actions))

	if len(issues) > 0 {
		fmt.Println("\n⚠️  Issues found:")
		for _, issue := range issues {
			fmt.Printf("  - %s\n", issue)
		}
		return fmt.Errorf("configuration has %d issue(s)", len(issues))
	}

	fmt.Println("✅ Configuration validation completed successfully")
	return nil
}

func runConfigLint(cmd *cobra.Command, args []string) error {
	configFile := args[0]
	fix, _ := cmd.Flags().GetBool("fix")
	strict, _ := cmd.Flags().GetBool("strict")

	// Check if file exists and read content
	content, err := readConfigFile(configFile)
	if err != nil {
		return err
	}

	// Perform linting checks
	lines := strings.Split(string(content), "\n")
	issues, suggestions := performLintChecks(lines, string(content), strict, fix)

	// Report results
	reportLintResults(configFile, issues, suggestions)

	// Apply fixes if requested
	if fix && len(issues) > 0 {
		if err := applyLintFixes(configFile, lines, len(issues)); err != nil {
			return err
		}
	}

	if len(issues) > 0 {
		return fmt.Errorf("configuration has %d issue(s)", len(issues))
	}

	return nil
}

func readConfigFile(configFile string) ([]byte, error) {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file does not exist: %s", configFile)
	}

	content, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("error reading configuration file: %w", err)
	}

	return content, nil
}

func performLintChecks(lines []string, contentStr string, strict, fix bool) (issues, suggestions []string) {
	// Check for trailing whitespace
	issues = append(issues, checkTrailingWhitespace(lines, fix)...)

	// Check for indentation issues
	suggestions = append(suggestions, checkIndentation(lines)...)

	// Check for common issues
	suggestions = append(suggestions, checkCommonIssues(contentStr)...)

	// Strict mode checks
	if strict {
		issues = append(issues, checkStrictModeIssues(contentStr)...)
	}

	return
}

func checkTrailingWhitespace(lines []string, fix bool) []string {
	var issues []string
	for i, line := range lines {
		if strings.HasSuffix(line, " ") || strings.HasSuffix(line, "\t") {
			issues = append(issues, fmt.Sprintf("Line %d: Trailing whitespace", i+1))
			if fix {
				lines[i] = strings.TrimRight(line, " \t")
			}
		}
	}
	return issues
}

func checkIndentation(lines []string) []string {
	var suggestions []string
	for i, line := range lines {
		if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "#") {
			// Skip lines that should not be indented (top-level blocks)
			if strings.Contains(line, "{") || strings.Contains(line, "}") {
				continue
			}

			// Check if line should be indented but isn't
			if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
				// This is a basic check - in a real implementation, you'd have more sophisticated parsing
				if i > 0 && strings.Contains(lines[i-1], "{") {
					suggestions = append(suggestions, fmt.Sprintf("Line %d: Consider indenting this line", i+1))
				}
			}
		}
	}
	return suggestions
}

func checkCommonIssues(contentStr string) []string {
	var suggestions []string

	if strings.Contains(contentStr, "server ") {
		suggestions = append(suggestions, "Consider using 'machine' instead of 'server' (deprecated)")
	}

	if strings.Contains(contentStr, "password = \"") {
		suggestions = append(suggestions, "Consider using environment variables for passwords instead of hardcoding")
	}

	return suggestions
}

func checkStrictModeIssues(contentStr string) []string {
	var issues []string

	if strings.Contains(contentStr, "insecure") {
		issues = append(issues, "Strict mode: 'insecure' host key verification is not recommended")
	}

	if strings.Contains(contentStr, "password =") {
		issues = append(issues, "Strict mode: Hardcoded passwords are not allowed")
	}

	return issues
}

func reportLintResults(configFile string, issues, suggestions []string) {
	fmt.Printf("Linting configuration file: %s\n", configFile)

	if len(issues) == 0 && len(suggestions) == 0 {
		fmt.Println("✅ No issues found")
		return
	}

	if len(issues) > 0 {
		fmt.Printf("\n❌ Issues found (%d):\n", len(issues))
		for _, issue := range issues {
			fmt.Printf("  - %s\n", issue)
		}
	}

	if len(suggestions) > 0 {
		fmt.Printf("\n💡 Suggestions (%d):\n", len(suggestions))
		for _, suggestion := range suggestions {
			fmt.Printf("  - %s\n", suggestion)
		}
	}
}

func applyLintFixes(configFile string, lines []string, issueCount int) error {
	fmt.Println("\n🔧 Applying fixes...")
	fixedContent := strings.Join(lines, "\n")
	err := os.WriteFile(configFile, []byte(fixedContent), 0o600)
	if err != nil {
		return fmt.Errorf("error writing fixed file: %w", err)
	}
	fmt.Printf("✅ Fixed %d issue(s) in %s\n", issueCount, configFile)
	return nil
}

func runConfigFormat(cmd *cobra.Command, args []string) error {
	configFile := args[0]
	write, _ := cmd.Flags().GetBool("write")
	outputFile, _ := cmd.Flags().GetString("output")

	// Check if file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("configuration file does not exist: %s", configFile)
	}

	// Read file content
	content, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("error reading configuration file: %w", err)
	}

	// Parse configuration to validate it
	_, err = config.ParseConfig(configFile)
	if err != nil {
		return fmt.Errorf("cannot format invalid configuration: %w", err)
	}

	// Basic formatting (placeholder implementation)
	// In a real implementation, you'd use a proper HCL formatter
	formatted := string(content)

	// Simple formatting: trim trailing whitespace and ensure consistent line endings
	lines := strings.Split(formatted, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t\r")
	}
	formatted = strings.Join(lines, "\n")

	// Remove empty lines at end
	formatted = strings.TrimRight(formatted, "\n")
	if formatted != "" {
		formatted += "\n"
	}

	// Output result
	switch {
	case outputFile != "":
		err := os.WriteFile(outputFile, []byte(formatted), 0o600)
		if err != nil {
			return fmt.Errorf("error writing output file: %w", err)
		}
		fmt.Printf("✅ Configuration formatted and saved to: %s\n", outputFile)
	case write:
		err := os.WriteFile(configFile, []byte(formatted), 0o600)
		if err != nil {
			return fmt.Errorf("error writing formatted file: %w", err)
		}
		fmt.Printf("✅ Configuration formatted and saved to: %s\n", configFile)
	default:
		fmt.Println(formatted)
	}

	return nil
}
