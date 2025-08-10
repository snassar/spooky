package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CreateTemplatesCommand creates the main templates command
func CreateTemplatesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "templates",
		Short: "Manage templates",
		Long:  "Manage and render templates",
	}

	// Add subcommands
	cmd.AddCommand(createTemplatesListCommand())
	cmd.AddCommand(createTemplatesRenderCommand())
	cmd.AddCommand(createTemplatesValidateCommand())

	return cmd
}

// createTemplatesListCommand creates the templates list subcommand
func createTemplatesListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List templates",
		Long:  "List all templates in a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement properly with correct types
			return fmt.Errorf("not implemented - interface mismatches need to be resolved")
		},
	}
	return cmd
}

// createTemplatesRenderCommand creates the templates render subcommand
func createTemplatesRenderCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "render",
		Short: "Render template",
		Long:  "Render a template with data",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement properly with correct types
			return fmt.Errorf("not implemented - interface mismatches need to be resolved")
		},
	}
	return cmd
}

// createTemplatesValidateCommand creates the templates validate subcommand
func createTemplatesValidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate templates",
		Long:  "Validate template configurations",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement properly with correct types
			return fmt.Errorf("not implemented - interface mismatches need to be resolved")
		},
	}
	return cmd
}
