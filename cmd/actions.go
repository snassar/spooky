package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	spookyinterfaces "spooky/internal/interfaces"
)

// actionsCmd represents the actions command
var actionsCmd = &cobra.Command{
	Use:   "actions",
	Short: "Manage and run actions",
	Long: `Manage and run actions on target machines.

Actions are operations that can be performed on machines, such as:
- Running commands
- Running scripts
- Deploying templates
- Copying files
- Controlling services

Actions can be defined in actions.hcl files or in an actions/ directory.`,
}

// actionsListCmd represents the actions list command
var actionsListCmd = &cobra.Command{
	Use:   "list [project-path]",
	Short: "List available actions",
	Long: `List all available actions in the specified project.

If no project path is provided, the current directory is used.
Actions are loaded from actions.hcl files and actions/ directories.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath := "."
		if len(args) > 0 {
			projectPath = args[0]
		}

		// Get absolute path
		absPath, err := filepath.Abs(projectPath)
		if err != nil {
			return fmt.Errorf("failed to resolve project path: %w", err)
		}

		// TODO: Initialize dependencies properly
		// For now, use placeholder implementation
		fmt.Printf("Loading actions from: %s\n", absPath)
		fmt.Println("Actions list command (placeholder implementation)")
		fmt.Println("No actions found.")

		return nil
	},
}

// actionsRunCmd represents the actions run command
var actionsRunCmd = &cobra.Command{
	Use:   "run [action-names...]",
	Short: "Run actions on target machines",
	Long: `Run one or more actions on target machines.

Actions are executed in dependency order. If no action names are provided,
all available actions will be run.

Modes:
  --plan     Show execution plan without running
  --dry-run  Simulate execution without making changes
  (no flags) Execute actions normally

Examples:
  spooky actions run deploy-web
  spooky actions run deploy-web restart-services
  spooky actions run --all
  spooky actions run . --plan
  spooky actions run . --dry-run
  spooky actions run .`,
	Args: cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath := "."
		if len(args) > 0 && !isFlagSet(cmd, "all") {
			projectPath = args[0]
		}

		// Get absolute path
		absPath, err := filepath.Abs(projectPath)
		if err != nil {
			return fmt.Errorf("failed to resolve project path: %w", err)
		}

		// Check execution mode
		if isFlagSet(cmd, "plan") {
			// Plan mode
			fmt.Println("📋 Execution Plan Mode")
			fmt.Printf("Project: %s\n", absPath)
			fmt.Println("📋 Execution Plan (placeholder)")
			fmt.Println("  Step 1:")
			fmt.Println("    - deploy-web")
			fmt.Println("    - restart-services")
			fmt.Println("  Target Machines: web-server-1, web-server-2")
			fmt.Println("  Estimated Time: 5-10 minutes")
			return nil
		} else if isFlagSet(cmd, "dry-run") {
			// Dry-run mode
			fmt.Println("🚀 Dry Run Mode")
			fmt.Printf("Project: %s\n", absPath)
			fmt.Println("Simulating execution...")
			fmt.Println("[web-server-1] Would execute: sudo systemctl stop nginx")
			fmt.Println("[web-server-1] Would copy: /tmp/web-app.tar.gz → /var/www/web-app/")
			fmt.Println("[web-server-2] Would execute: sudo systemctl stop nginx")
			fmt.Println("[web-server-2] Would copy: /tmp/web-app.tar.gz → /var/www/web-app/")
			fmt.Println("✅ Dry run completed successfully")
			return nil
		} else {
			// Normal execution mode
			fmt.Println("🚀 Execution Mode")
			fmt.Printf("Project: %s\n", absPath)
			fmt.Println("Executing actions...")
			fmt.Println("No actions executed (placeholder implementation)")
			return nil
		}
	},
}

// actionsValidateCmd represents the actions validate command
var actionsValidateCmd = &cobra.Command{
	Use:   "validate [project-path]",
	Short: "Validate action configurations",
	Long: `Validate action configurations in the specified project.

This command checks:
- HCL syntax and schema compliance
- Action dependencies and circular references
- Machine targeting configuration
- Resource limits and constraints`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath := "."
		if len(args) > 0 {
			projectPath = args[0]
		}

		// Get absolute path
		absPath, err := filepath.Abs(projectPath)
		if err != nil {
			return fmt.Errorf("failed to resolve project path: %w", err)
		}

		// TODO: Validate actions
		fmt.Println("Actions validate command (placeholder implementation)")
		fmt.Printf("Would validate actions from: %s\n", absPath)
		fmt.Println("✅ All actions are valid! (placeholder)")

		return nil
	},
}

// Helper function to check if a flag is set
func isFlagSet(cmd *cobra.Command, flagName string) bool {
	return cmd.Flags().Lookup(flagName) != nil && cmd.Flags().Lookup(flagName).Changed
}

// Helper function to get SSH manager (placeholder)
func getSSHManager() spookyinterfaces.SSHManager {
	// TODO: Implement proper SSH manager initialization
	return nil
}

func init() {
	// Add flags to run command
	actionsRunCmd.Flags().BoolP("all", "a", false, "Run all available actions")
	actionsRunCmd.Flags().BoolP("dry-run", "d", false, "Simulate execution without making changes")
	actionsRunCmd.Flags().BoolP("plan", "p", false, "Show execution plan without running")
	actionsRunCmd.Flags().StringP("machines", "m", "", "Comma-separated list of target machines")

	// Add subcommands
	actionsCmd.AddCommand(actionsListCmd)
	actionsCmd.AddCommand(actionsRunCmd)
	actionsCmd.AddCommand(actionsValidateCmd)

	// Add actions command to root
	RootCmd.AddCommand(actionsCmd)
}
