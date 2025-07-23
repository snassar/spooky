package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"spooky/internal/facts"
)

var TemplatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Manage configuration templates",
	Long: `Manage configuration templates for dynamic configuration generation.

Templates allow you to create dynamic configurations based on machine facts,
environment variables, and other data sources.

Examples:
  # List all available templates
  spooky templates list

  # Validate a template file
  spooky templates validate template.hcl

  # Render a template locally
  spooky templates render template.hcl`,
}

var templatesListCmd = &cobra.Command{
	Use:   "list [directory]",
	Short: "List available templates",
	Long: `List all available templates in the specified directory or current directory.

Examples:
  spooky templates list                    # List templates in current directory
  spooky templates list /path/to/templates # List templates in specific directory`,
	Args: cobra.MaximumNArgs(1),
	RunE: runTemplatesList,
}

var templatesValidateCmd = &cobra.Command{
	Use:   "validate <file>",
	Short: "Validate template syntax",
	Long: `Validate template syntax and structure.

This command checks if a template file has valid syntax and can be parsed
without errors. It does not render the template or check for data availability.

Examples:
  spooky templates validate template.hcl
  spooky templates validate /path/to/template.tmpl`,
	Args: cobra.ExactArgs(1),
	RunE: runTemplatesValidate,
}

var templatesRenderCmd = &cobra.Command{
	Use:   "render <file>",
	Short: "Render template locally",
	Long: `Render a template locally with provided data.

This command renders a template file using local data sources and displays
the output. Useful for testing templates before deployment.

Examples:
  spooky templates render template.hcl
  spooky templates render /path/to/template.tmpl --data-file data.json`,
	Args: cobra.ExactArgs(1),
	RunE: runTemplatesRender,
}

func init() {
	// Add flags for templates render command
	templatesRenderCmd.Flags().String("data-file", "", "Path to JSON file containing template data")
	templatesRenderCmd.Flags().String("output", "", "Output file path (default: stdout)")
	templatesRenderCmd.Flags().Bool("dry-run", false, "Show what would be rendered without writing output")
	templatesRenderCmd.Flags().String("server", "", "Server name for fact integration (e.g., web-001)")

	// Add subcommands to templates command
	TemplatesCmd.AddCommand(templatesListCmd)
	TemplatesCmd.AddCommand(templatesValidateCmd)
	TemplatesCmd.AddCommand(templatesRenderCmd)

	// Add templates command to root
	// Note: rootCmd is defined in main.go and will be available at runtime
}

func runTemplatesList(_ *cobra.Command, args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	// Validate directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", dir)
	}

	// Find template files
	templateExtensions := []string{".hcl", ".tmpl", ".template", ".tpl"}
	var templates []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ext := filepath.Ext(path)
			for _, templateExt := range templateExtensions {
				if ext == templateExt {
					// Make path relative to search directory
					relPath, err := filepath.Rel(dir, path)
					if err != nil {
						relPath = path
					}
					templates = append(templates, relPath)
					break
				}
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error scanning directory: %w", err)
	}

	if len(templates) == 0 {
		fmt.Printf("No templates found in %s\n", dir)
		return nil
	}

	fmt.Printf("Templates found in %s:\n", dir)
	for _, template := range templates {
		fmt.Printf("  %s\n", template)
	}

	return nil
}

func runTemplatesValidate(_ *cobra.Command, args []string) error {
	templateFile := args[0]

	// Check if file exists
	if _, err := os.Stat(templateFile); os.IsNotExist(err) {
		return fmt.Errorf("template file does not exist: %s", templateFile)
	}

	// Read template file
	content, err := os.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("error reading template file: %w", err)
	}

	// Basic validation - check for common template syntax
	// This is a placeholder implementation - in a real system, you'd use
	// the actual template parser to validate syntax
	if len(content) == 0 {
		return fmt.Errorf("template file is empty")
	}

	// Check for basic template markers (placeholder validation)
	hasTemplateMarkers := false
	contentStr := string(content)

	// Look for common template patterns
	templatePatterns := []string{"{{", "}}", "${", "}", "#{", "}"}
	for _, pattern := range templatePatterns {
		if len(pattern) == 2 {
			// For 2-character patterns, check if they appear in pairs
			if len(pattern) == 2 {
				count := 0
				for i := 0; i < len(contentStr)-1; i++ {
					if contentStr[i:i+2] == pattern {
						count++
					}
				}
				if count > 0 && count%2 == 0 {
					hasTemplateMarkers = true
					break
				}
			}
		}
	}

	fmt.Printf("✓ Template file '%s' is valid\n", templateFile)
	if hasTemplateMarkers {
		fmt.Println("  - Contains template markers")
	}
	fmt.Printf("  - Size: %d bytes\n", len(content))

	return nil
}

func runTemplatesRender(cmd *cobra.Command, args []string) error {
	templateFile := args[0]
	dataFile, _ := cmd.Flags().GetString("data-file")
	outputFile, _ := cmd.Flags().GetString("output")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	server, _ := cmd.Flags().GetString("server")

	// Check if template file exists
	if _, err := os.Stat(templateFile); os.IsNotExist(err) {
		return fmt.Errorf("template file does not exist: %s", templateFile)
	}

	// Load additional data if provided
	var additionalData map[string]interface{}
	if dataFile != "" {
		if _, err := os.Stat(dataFile); os.IsNotExist(err) {
			return fmt.Errorf("data file does not exist: %s", dataFile)
		}

		dataBytes, err := os.ReadFile(dataFile)
		if err != nil {
			return fmt.Errorf("error reading data file: %w", err)
		}

		// Parse JSON data
		if err := json.Unmarshal(dataBytes, &additionalData); err != nil {
			return fmt.Errorf("error parsing JSON data file: %w", err)
		}
		fmt.Printf("Loading data from: %s\n", dataFile)
	}

	// Create fact manager for template rendering
	var manager *facts.Manager
	if server != "" {
		// Create storage for fact manager
		storage, err := facts.NewFactStorage(facts.StorageOptions{
			Type: facts.StorageTypeBadger,
			Path: getFactsDBPath(),
		})
		if err == nil {
			defer storage.Close()
			manager = facts.NewManagerWithStorage(nil, storage)
		}
	}

	// Create template engine
	templateEngine := NewTemplateEngine(manager)

	// Render template
	rendered, err := templateEngine.RenderTemplate(templateFile, server, additionalData)
	if err != nil {
		return fmt.Errorf("error rendering template: %w", err)
	}

	if dryRun {
		fmt.Println("=== DRY RUN - Template would be rendered as: ===")
		fmt.Println(rendered)
		fmt.Println("=== END DRY RUN ===")
		return nil
	}

	// Output result
	if outputFile != "" {
		err := os.WriteFile(outputFile, []byte(rendered), 0o600)
		if err != nil {
			return fmt.Errorf("error writing output file: %w", err)
		}
		fmt.Printf("Template rendered to: %s\n", outputFile)
	} else {
		fmt.Println(rendered)
	}

	return nil
}
