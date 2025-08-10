package commands

import (
	"fmt"
	"os"
	"path/filepath"

	spookycoordinator "spooky/internal/coordinator"
	spookylogging "spooky/internal/logging"
	spookyloggingtypes "spooky/internal/types/logging"

	"github.com/spf13/cobra"
)

// CreateTemplatesCommand creates the templates command with all subcommands
func CreateTemplatesCommand() *cobra.Command {
	templatesCmd := &cobra.Command{
		Use:   "templates",
		Short: "Manage templates",
		Long:  "Manage and render templates",
	}

	// Add subcommands
	templatesCmd.AddCommand(createTemplatesRenderCommand())
	templatesCmd.AddCommand(createTemplatesValidateCommand())
	templatesCmd.AddCommand(createTemplatesListCommand())

	return templatesCmd
}

// createTemplatesRenderCommand creates the templates render subcommand
func createTemplatesRenderCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "render",
		Short: "Render templates",
		Long:  "Render templates with data",
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
			if _, err := os.Stat(projectPath); os.IsNotExist(err) {
				return fmt.Errorf("project directory does not exist: %s", projectPath)
			}

			// Validate template path
			fullTemplatePath := filepath.Join(projectPath, "templates", templatePath)
			if _, err := os.Stat(fullTemplatePath); os.IsNotExist(err) {
				return fmt.Errorf("template file does not exist: %s", fullTemplatePath)
			}

			// Create logger
			logger := spookylogging.NewLogger(spookyloggingtypes.Config{
				Level:  spookyloggingtypes.InfoLevel,
				Format: "text",
				Output: "stdout",
			})

			// Create coordinator manager
			coord, err := spookycoordinator.NewCoordinatorManagerFromProject(projectPath, logger)
			if err != nil {
				return fmt.Errorf("failed to create coordinator: %w", err)
			}

			// Load templates context
			templatesContext, err := coord.Templates().LoadTemplates(projectPath)
			if err != nil {
				return fmt.Errorf("failed to load templates: %w", err)
			}

			// Load variables context if data file provided
			if dataFile != "" {
				_, err = coord.Variables().LoadVariables(projectPath)
				if err != nil {
					return fmt.Errorf("failed to load variables: %w", err)
				}
			}

			// Render template
			fmt.Printf("Rendering template: %s\n", templatePath)

			// Get template
			template, err := coord.Templates().GetTemplate(templatePath, templatesContext)
			if err != nil {
				return fmt.Errorf("failed to get template: %w", err)
			}

			// Render template
			rendered, err := coord.Templates().RenderTemplate(template, templatesContext)
			if err != nil {
				return fmt.Errorf("failed to render template: %w", err)
			}

			// Handle output
			if dryRun {
				fmt.Println("DRY RUN: Template would be rendered as:")
				fmt.Println("---")
				fmt.Println(rendered)
				fmt.Println("---")
				return nil
			}

			if preview {
				fmt.Println("PREVIEW: Template rendered as:")
				fmt.Println("---")
				fmt.Println(rendered)
				fmt.Println("---")
				return nil
			}

			// Write to output file
			if outputFile == "" {
				// Use template name with .rendered extension
				ext := filepath.Ext(templatePath)
				base := templatePath[:len(templatePath)-len(ext)]
				outputFile = filepath.Join(projectPath, fmt.Sprintf("%s.rendered%s", base, ext))
			}

			// Ensure output directory exists
			outputDir := filepath.Dir(outputFile)
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}

			// Write rendered content
			if err := os.WriteFile(outputFile, []byte(rendered), 0644); err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}

			fmt.Printf("✅ Template rendered successfully to: %s\n", outputFile)
			return nil
		},
	}

	// Add flags
	cmd.Flags().StringP("template", "t", "", "Template file path")
	cmd.Flags().StringP("data", "d", "", "Data file path")
	cmd.Flags().StringP("output", "o", "", "Output file path")
	cmd.Flags().BoolP("dry-run", "n", false, "Dry run mode")
	cmd.Flags().BoolP("preview", "p", false, "Preview mode")

	return cmd
}

// createTemplatesValidateCommand creates the templates validate subcommand
func createTemplatesValidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate templates",
		Long:  "Validate template syntax and structure",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath := args[0]

			// Get flags
			templatePath, _ := cmd.Flags().GetString("template")
			templates, _ := cmd.Flags().GetStringSlice("templates")

			// Validate project path
			if _, err := os.Stat(projectPath); os.IsNotExist(err) {
				return fmt.Errorf("project directory does not exist: %s", projectPath)
			}

			// Create logger
			logger := spookylogging.NewLogger(spookyloggingtypes.Config{
				Level:  spookyloggingtypes.InfoLevel,
				Format: "text",
				Output: "stdout",
			})

			// Create coordinator manager
			coord, err := spookycoordinator.NewCoordinatorManagerFromProject(projectPath, logger)
			if err != nil {
				return fmt.Errorf("failed to create coordinator: %w", err)
			}

			// Load templates context
			templatesContext, err := coord.Templates().LoadTemplates(projectPath)
			if err != nil {
				return fmt.Errorf("failed to load templates: %w", err)
			}

			// Validate templates
			var templatesToValidate []string

			if templatePath != "" {
				templatesToValidate = []string{templatePath}
			} else if len(templates) > 0 {
				templatesToValidate = templates
			} else {
				// Validate all templates
				allTemplates, err := coord.Templates().ListTemplates(templatesContext)
				if err != nil {
					return fmt.Errorf("failed to list templates: %w", err)
				}
				for _, template := range allTemplates {
					templatesToValidate = append(templatesToValidate, template.Name)
				}
			}

			if len(templatesToValidate) == 0 {
				fmt.Println("No templates found to validate")
				return nil
			}

			fmt.Printf("Validating %d templates...\n", len(templatesToValidate))

			validCount := 0
			invalidCount := 0

			for _, templateName := range templatesToValidate {
				fmt.Printf("Validating %s... ", templateName)

				template, err := coord.Templates().GetTemplate(templateName, templatesContext)
				if err != nil {
					fmt.Printf("❌ FAILED: %v\n", err)
					invalidCount++
					continue
				}

				if err := coord.Templates().ValidateTemplate(template, templatesContext); err != nil {
					fmt.Printf("❌ INVALID: %v\n", err)
					invalidCount++
				} else {
					fmt.Printf("✅ VALID\n")
					validCount++
				}
			}

			fmt.Printf("\nValidation complete: %d valid, %d invalid\n", validCount, invalidCount)

			if invalidCount > 0 {
				return fmt.Errorf("validation failed: %d templates are invalid", invalidCount)
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringP("template", "t", "", "Template file path")
	cmd.Flags().StringSliceP("templates", "T", []string{}, "Template files to validate")

	return cmd
}

// createTemplatesListCommand creates the templates list subcommand
func createTemplatesListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List templates",
		Long:  "List available templates",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath := args[0]

			// Get flags
			_, _ = cmd.Flags().GetString("filter") // Keep flag but don't use for now

			// Validate project path
			if _, err := os.Stat(projectPath); os.IsNotExist(err) {
				return fmt.Errorf("project directory does not exist: %s", projectPath)
			}

			// Create logger
			logger := spookylogging.NewLogger(spookyloggingtypes.Config{
				Level:  spookyloggingtypes.InfoLevel,
				Format: "text",
				Output: "stdout",
			})

			// Create coordinator manager
			coord, err := spookycoordinator.NewCoordinatorManagerFromProject(projectPath, logger)
			if err != nil {
				return fmt.Errorf("failed to create coordinator: %w", err)
			}

			// Load templates context
			templatesContext, err := coord.Templates().LoadTemplates(projectPath)
			if err != nil {
				return fmt.Errorf("failed to load templates: %w", err)
			}

			// List templates
			templates, err := coord.Templates().ListTemplates(templatesContext)
			if err != nil {
				return fmt.Errorf("failed to list templates: %w", err)
			}

			// Display templates
			if len(templates) == 0 {
				fmt.Println("No templates found")
				return nil
			}

			fmt.Printf("Found %d templates:\n\n", len(templates))

			for _, template := range templates {
				fmt.Printf("  %s\n", template.Name)
				if template.Metadata != nil && template.Metadata.Description != "" {
					fmt.Printf("    Description: %s\n", template.Metadata.Description)
				}
				if template.Metadata != nil && template.Metadata.OutputFormat != "" {
					fmt.Printf("    Type: %s\n", template.Metadata.OutputFormat)
				}
				fmt.Println()
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringP("filter", "f", "", "Filter expression")

	return cmd
}
