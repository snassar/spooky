package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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
			if err := os.WriteFile(outputFile, []byte(result), 0o600); err != nil {
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

func loadDataFile(_ string) (map[string]interface{}, error) {
	// For now, return empty data
	// This will be enhanced to load from JSON/HCL files
	return make(map[string]interface{}), nil
}

func validateAllTemplates(ctx context.Context, integration spookyinterfaces.TemplatesIntegration, templatesDir string) error {
	// Check if templates directory exists
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		fmt.Printf("No templates directory found: %s\n", templatesDir)
		return nil
	}

	fmt.Printf("Validating templates in: %s\n", templatesDir)

	// Scan for template files
	templateFiles, err := filepath.Glob(filepath.Join(templatesDir, "*.tmpl"))
	if err != nil {
		return fmt.Errorf("failed to scan templates directory: %w", err)
	}

	if len(templateFiles) == 0 {
		fmt.Println("No template files found")
		return nil
	}

	var totalErrors int
	var totalWarnings int

	// Validate each template file
	for _, templateFile := range templateFiles {
		fmt.Printf("Validating: %s\n", filepath.Base(templateFile))

		// Load template
		template, err := integration.LoadTemplate(ctx, templateFile)
		if err != nil {
			fmt.Printf("  ❌ Failed to load template: %v\n", err)
			totalErrors++
			continue
		}

		// Validate template using schema validation
		result, err := integration.ValidateTemplate(ctx, template)
		if err != nil {
			fmt.Printf("  ❌ Validation failed: %v\n", err)
			totalErrors++
			continue
		}

		if result.Valid {
			fmt.Printf("  ✅ Template is valid")
			if len(result.Warnings) > 0 {
				fmt.Printf(" (with %d warnings)", len(result.Warnings))
			}
			fmt.Println()
		} else {
			fmt.Printf("  ❌ Template has %d errors", len(result.Errors))
			if len(result.Warnings) > 0 {
				fmt.Printf(" and %d warnings", len(result.Warnings))
			}
			fmt.Println()

			// Print errors
			for i := range result.Errors {
				fmt.Printf("    - Error: %s\n", result.Errors[i].Message)
			}

			// Print warnings
			for i := range result.Warnings {
				fmt.Printf("    - Warning: %s\n", result.Warnings[i].Message)
			}

			totalErrors += len(result.Errors)
			totalWarnings += len(result.Warnings)
		}
	}

	// Print summary
	fmt.Printf("\nValidation Summary:\n")
	fmt.Printf("  Templates processed: %d\n", len(templateFiles))
	fmt.Printf("  Total errors: %d\n", totalErrors)
	fmt.Printf("  Total warnings: %d\n", totalWarnings)

	if totalErrors > 0 {
		return fmt.Errorf("validation failed with %d errors", totalErrors)
	}

	return nil
}

func listTemplates(projectPath string) ([]*spookytypes.Template, error) {
	templatesDir := filepath.Join(projectPath, "templates")

	// Check if templates directory exists
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		return []*spookytypes.Template{}, nil
	}

	// Scan for template files
	templateFiles, err := filepath.Glob(filepath.Join(templatesDir, "*.tmpl"))
	if err != nil {
		return nil, fmt.Errorf("failed to scan templates directory: %w", err)
	}

	var templates []*spookytypes.Template

	// Load each template
	for _, templateFile := range templateFiles {
		// Create a basic template object for now
		// In a full implementation, this would load the actual template content
		template := &spookytypes.Template{
			ID:         filepath.Base(templateFile),
			SourcePath: templateFile,
			Type:       "template",
			Scope:      "project",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		templates = append(templates, template)
	}

	return templates, nil
}

func searchTemplates(projectPath, query string, tags []string, category string) ([]*spookytypes.Template, error) {
	// Get templates integration
	integration := GetIntegrationManager().GetTemplatesIntegration()
	if integration == nil {
		return nil, fmt.Errorf("templates integration not available")
	}

	// Load all templates in the project
	allTemplates, err := loadAllTemplates(integration, projectPath)
	if err != nil {
		return nil, err
	}

	// Filter templates based on search criteria
	return filterTemplates(allTemplates, query, tags, category), nil
}

// loadAllTemplates loads all templates from the project
func loadAllTemplates(integration spookyinterfaces.TemplatesIntegration, projectPath string) ([]*spookytypes.Template, error) {
	templatesDir := filepath.Join(projectPath, "templates")
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		return []*spookytypes.Template{}, nil
	}

	// Scan for template files
	templateFiles, err := filepath.Glob(filepath.Join(templatesDir, "*.tmpl"))
	if err != nil {
		return nil, fmt.Errorf("failed to scan templates directory: %w", err)
	}

	var allTemplates []*spookytypes.Template

	// Load each template
	for _, templateFile := range templateFiles {
		template, err := integration.LoadTemplate(context.Background(), templateFile)
		if err != nil {
			// Skip templates that can't be loaded
			continue
		}
		allTemplates = append(allTemplates, template)
	}

	return allTemplates, nil
}

// filterTemplates filters templates based on search criteria
func filterTemplates(templates []*spookytypes.Template, query string, tags []string, category string) []*spookytypes.Template {
	var results []*spookytypes.Template
	queryLower := strings.ToLower(query)

	for _, template := range templates {
		if matchesTemplate(template, queryLower, tags, category) {
			results = append(results, template)
		}
	}

	return results
}

// matchesTemplate checks if a template matches the search criteria
func matchesTemplate(template *spookytypes.Template, queryLower string, tags []string, category string) bool {
	// Search in template name
	if strings.Contains(strings.ToLower(template.ID), queryLower) {
		return true
	}

	// Search in metadata
	if template.Metadata != nil {
		// Search in metadata name
		if template.Metadata.Name != "" && strings.Contains(strings.ToLower(template.Metadata.Name), queryLower) {
			return true
		}

		// Search in metadata description
		if template.Metadata.Description != "" && strings.Contains(strings.ToLower(template.Metadata.Description), queryLower) {
			return true
		}

		// Check tag filter
		if len(tags) > 0 && !hasMatchingTag(template.Metadata.Tags, tags) {
			return false
		}

		// Check category filter
		if category != "" && template.Metadata.Category != category {
			return false
		}
	}

	return false
}

// hasMatchingTag checks if template has any of the required tags
func hasMatchingTag(templateTags, searchTags []string) bool {
	if templateTags == nil {
		return false
	}

	for _, searchTag := range searchTags {
		for _, templateTag := range templateTags {
			if searchTag == templateTag {
				return true
			}
		}
	}

	return false
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
