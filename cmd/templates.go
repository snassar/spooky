package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"

	"github.com/spf13/cobra"
)

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Manage templates",
	Long:  `Manage template loading, rendering, validation, and discovery.`,
}

var templatesRenderCmd = &cobra.Command{
	Use:   "render [project] [template]",
	Short: "Render a template",
	Long:  `Render a template with the given data and output to a file.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath := args[0]
		templatePath := args[1]

		// Get flags
		dataFile, _ := cmd.Flags().GetString("data")
		outputFile, _ := cmd.Flags().GetString("output")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		preview, _ := cmd.Flags().GetBool("preview")

		// Validate project path
		if err := validateProjectPath(projectPath); err != nil {
			return fmt.Errorf("invalid project path: %w", err)
		}

		// Get templates integration
		integration := GetIntegrationManager().GetTemplatesIntegration()
		if integration == nil {
			return fmt.Errorf("templates integration not available")
		}

		// Load data file if provided
		var data map[string]interface{}
		var err error
		if dataFile != "" {
			data, err = loadDataFile(dataFile)
			if err != nil {
				return fmt.Errorf("failed to load data file: %w", err)
			}
		}

		// Load template
		ctx := context.Background()
		template, err := integration.LoadTemplate(ctx, templatePath)
		if err != nil {
			return fmt.Errorf("failed to load template: %w", err)
		}

		// Render template
		result, err := integration.RenderTemplate(ctx, template, data)
		if err != nil {
			return fmt.Errorf("failed to render template: %w", err)
		}

		// Handle output
		if dryRun || preview {
			fmt.Println("=== Template Rendering Result ===")
			fmt.Println(result)
			return nil
		}

		if outputFile != "" {
			// Write to output file
			if err := os.WriteFile(outputFile, []byte(result), 0644); err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}
			fmt.Printf("Template rendered successfully to: %s\n", outputFile)
		} else {
			// Write to stdout
			fmt.Print(result)
		}

		return nil
	},
}

var templatesValidateCmd = &cobra.Command{
	Use:   "validate [project]",
	Short: "Validate templates",
	Long:  `Validate templates in the project for syntax and security.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath := args[0]

		// Get flags
		templatePath, _ := cmd.Flags().GetString("template")

		// Validate project path
		if err := validateProjectPath(projectPath); err != nil {
			return fmt.Errorf("invalid project path: %w", err)
		}

		// Get templates integration
		integration := GetIntegrationManager().GetTemplatesIntegration()
		if integration == nil {
			return fmt.Errorf("templates integration not available")
		}

		ctx := context.Background()

		if templatePath != "" {
			// Validate specific template
			template, err := integration.LoadTemplate(ctx, templatePath)
			if err != nil {
				return fmt.Errorf("failed to load template: %w", err)
			}

			result, err := integration.ValidateTemplate(ctx, template)
			if err != nil {
				return fmt.Errorf("failed to validate template: %w", err)
			}

			if result.Valid {
				fmt.Printf("✅ Template %s is valid\n", templatePath)
			} else {
				fmt.Printf("❌ Template %s has validation errors:\n", templatePath)
				for _, err := range result.Errors {
					fmt.Printf("  - %s\n", err.Message)
				}
				return fmt.Errorf("template validation failed")
			}
		} else {
			// Validate all templates in project
			templatesDir := filepath.Join(projectPath, "templates")
			if err := validateAllTemplates(ctx, integration, templatesDir); err != nil {
				return fmt.Errorf("failed to validate templates: %w", err)
			}
		}

		return nil
	},
}

var templatesListCmd = &cobra.Command{
	Use:   "list [project]",
	Short: "List templates",
	Long:  `List all templates in the project.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath := args[0]

		// Get flags
		format, _ := cmd.Flags().GetString("format")

		// Validate project path
		if err := validateProjectPath(projectPath); err != nil {
			return fmt.Errorf("invalid project path: %w", err)
		}

		// Get templates integration
		integration := GetIntegrationManager().GetTemplatesIntegration()
		if integration == nil {
			return fmt.Errorf("templates integration not available")
		}

		// List templates
		templates, err := listTemplates(projectPath)
		if err != nil {
			return fmt.Errorf("failed to list templates: %w", err)
		}

		// Output in requested format
		switch format {
		case "json":
			outputJSON(templates)
		case "hcl":
			outputHCL(templates)
		case "table":
			outputTable(templates)
		default:
			return fmt.Errorf("unsupported format: %s", format)
		}

		return nil
	},
}

var templatesSearchCmd = &cobra.Command{
	Use:   "search [project] [query]",
	Short: "Search templates",
	Long:  `Search templates in the project by name, description, or tags.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath := args[0]
		query := args[1]

		// Get flags
		tags, _ := cmd.Flags().GetStringSlice("tags")
		category, _ := cmd.Flags().GetString("category")

		// Validate project path
		if err := validateProjectPath(projectPath); err != nil {
			return fmt.Errorf("invalid project path: %w", err)
		}

		// Get templates integration
		integration := GetIntegrationManager().GetTemplatesIntegration()
		if integration == nil {
			return fmt.Errorf("templates integration not available")
		}

		// Search templates
		templates, err := searchTemplates(projectPath, query, tags, category)
		if err != nil {
			return fmt.Errorf("failed to search templates: %w", err)
		}

		// Output results
		if len(templates) == 0 {
			fmt.Println("No templates found matching your search criteria")
			return nil
		}

		fmt.Printf("Found %d templates matching '%s':\n", len(templates), query)
		for i, template := range templates {
			name := template.ID
			description := ""
			if template.Metadata != nil {
				if template.Metadata.Name != "" {
					name = template.Metadata.Name
				}
				description = template.Metadata.Description
			}
			fmt.Printf("%d. %s", i+1, name)
			if description != "" {
				fmt.Printf(" - %s", description)
			}
			fmt.Println()
		}

		return nil
	},
}

func init() {
	// Add flags to render command
	templatesRenderCmd.Flags().String("data", "", "Data file (JSON/HCL) for template variables")
	templatesRenderCmd.Flags().String("output", "", "Output file path")
	templatesRenderCmd.Flags().Bool("dry-run", false, "Show what would be rendered without writing files")
	templatesRenderCmd.Flags().Bool("preview", false, "Preview the rendered template")

	// Add flags to validate command
	templatesValidateCmd.Flags().String("template", "", "Specific template to validate")

	// Add flags to list command
	templatesListCmd.Flags().String("format", "table", "Output format (table, json, hcl)")

	// Add flags to search command
	templatesSearchCmd.Flags().StringSlice("tags", []string{}, "Filter by tags")
	templatesSearchCmd.Flags().String("category", "", "Filter by category")

	// Add subcommands
	templatesCmd.AddCommand(templatesRenderCmd)
	templatesCmd.AddCommand(templatesValidateCmd)
	templatesCmd.AddCommand(templatesListCmd)
	templatesCmd.AddCommand(templatesSearchCmd)

	// Add templates command to root
	RootCmd.AddCommand(templatesCmd)
}

// Helper functions

func loadDataFile(dataFile string) (map[string]interface{}, error) {
	// For now, return empty data
	// This will be enhanced to load from JSON/HCL files
	return make(map[string]interface{}), nil
}

func validateAllTemplates(ctx context.Context, integration spookyinterfaces.TemplatesIntegration, templatesDir string) error {
	// For now, just check if templates directory exists
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		fmt.Printf("No templates directory found: %s\n", templatesDir)
		return nil
	}

	fmt.Printf("Validating templates in: %s\n", templatesDir)
	// This will be enhanced to validate all templates in the directory
	return nil
}

func listTemplates(projectPath string) ([]*spookytypes.Template, error) {
	// For now, return empty list
	// This will be enhanced to scan templates directory
	return []*spookytypes.Template{}, nil
}

func searchTemplates(projectPath, query string, tags []string, category string) ([]*spookytypes.Template, error) {
	// For now, return empty list
	// This will be enhanced to search templates
	return []*spookytypes.Template{}, nil
}

func outputJSON(templates []*spookytypes.Template) {
	// For now, just print count
	fmt.Printf("Found %d templates\n", len(templates))
}

func outputHCL(templates []*spookytypes.Template) {
	// For now, just print count
	fmt.Printf("Found %d templates\n", len(templates))
}

func outputTable(templates []*spookytypes.Template) {
	// For now, just print count
	fmt.Printf("Found %d templates\n", len(templates))
}
