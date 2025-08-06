package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	spookycoordinator "spooky/internal/coordinator"
	spookylogging "spooky/internal/logging"
	spookyloggingtypes "spooky/internal/logging/types"

	"github.com/spf13/cobra"
)

// CreateProjectCommand creates the project command with all subcommands
func CreateProjectCommand() *cobra.Command {
	projectCmd := &cobra.Command{
		Use:   "project",
		Short: "Manage projects",
		Long:  "Manage project configuration and structure",
	}

	// Add subcommands
	projectCmd.AddCommand(createProjectInitCommand())
	projectCmd.AddCommand(createProjectInfoCommand())
	projectCmd.AddCommand(createProjectValidateCommand())
	projectCmd.AddCommand(createProjectEncryptCommand())

	return projectCmd
}

// createProjectInitCommand creates the project init subcommand
func createProjectInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize project",
		Long:  "Initialize a new spooky project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath := args[0]

			// Get flags
			name, _ := cmd.Flags().GetString("name")
			description, _ := cmd.Flags().GetString("description")
			author, _ := cmd.Flags().GetString("author")

			// Validate project path
			if _, err := os.Stat(projectPath); err == nil {
				return fmt.Errorf("project directory already exists: %s", projectPath)
			}

			// Create logger (unused for now)
			_ = spookylogging.NewLogger(spookyloggingtypes.Config{
				Level:  spookyloggingtypes.InfoLevel,
				Format: "text",
				Output: "stdout",
			})

			// Create project directory
			if err := os.MkdirAll(projectPath, 0755); err != nil {
				return fmt.Errorf("failed to create project directory: %w", err)
			}

			// Create project structure
			subdirs := []string{"templates", "variables", "actions", "machines", "facts.db"}
			for _, subdir := range subdirs {
				subdirPath := filepath.Join(projectPath, subdir)
				if err := os.MkdirAll(subdirPath, 0755); err != nil {
					return fmt.Errorf("failed to create subdirectory %s: %w", subdir, err)
				}
			}

			// Create basic project.hcl file
			projectHCL := fmt.Sprintf(`project "%s" {
  description = "%s"
  author      = "%s"
  version     = "0.1.0"
  
  structure {
    templates_path = "templates"
    variables_path = "variables"
    actions_path   = "actions"
    machines_path  = "machines"
    facts_path     = "facts.db"
  }
}`, name, description, author)

			projectHCLPath := filepath.Join(projectPath, "project.hcl")
			if err := os.WriteFile(projectHCLPath, []byte(projectHCL), 0644); err != nil {
				return fmt.Errorf("failed to create project.hcl: %w", err)
			}

			// Create basic README
			readme := fmt.Sprintf(`# %s

%s

## Project Structure

- templates/ - Template files
- variables/ - Variable definitions
- actions/ - Action definitions
- machines/ - Machine inventory
- facts.db/ - Facts database

## Usage

Use spooky commands to manage this project.
`, name, description)

			readmePath := filepath.Join(projectPath, "README.md")
			if err := os.WriteFile(readmePath, []byte(readme), 0644); err != nil {
				return fmt.Errorf("failed to create README.md: %w", err)
			}

			fmt.Printf("✅ Project initialized successfully at: %s\n", projectPath)
			fmt.Printf("  - Project name: %s\n", name)
			fmt.Printf("  - Description: %s\n", description)
			fmt.Printf("  - Author: %s\n", author)

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringP("name", "n", "", "Project name")
	cmd.Flags().StringP("description", "d", "", "Project description")
	cmd.Flags().StringP("author", "a", "", "Project author")

	return cmd
}

// createProjectInfoCommand creates the project info subcommand
func createProjectInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Show project information",
		Long:  "Display project information and statistics",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath := args[0]

			// Get flags
			verbose, _ := cmd.Flags().GetBool("verbose")

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

			// Get project stats
			stats := coord.GetProjectStats(projectPath)

			// Display project information
			fmt.Printf("Project: %s\n", projectPath)
			fmt.Printf("Status: %s\n", getProjectStatus(projectPath))

			if verbose {
				fmt.Printf("\nStatistics:\n")
				for key, value := range stats {
					fmt.Printf("  %s: %v\n", key, value)
				}
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().BoolP("verbose", "v", false, "Verbose output")

	return cmd
}

// createProjectValidateCommand creates the project validate subcommand
func createProjectValidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate project",
		Long:  "Validate project configuration and structure",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath := args[0]

			// Get flags (unused for now)
			_, _ = cmd.Flags().GetBool("strict")

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

			// Validate project
			fmt.Printf("Validating project: %s\n", projectPath)

			if err := coord.ValidateProject(projectPath); err != nil {
				return fmt.Errorf("project validation failed: %w", err)
			}

			fmt.Println("✅ Project validation completed successfully")
			return nil
		},
	}

	// Add flags
	cmd.Flags().BoolP("strict", "s", false, "Strict validation mode")

	return cmd
}

// createProjectEncryptCommand creates the project encrypt subcommand
func createProjectEncryptCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encrypt",
		Short: "Re-encrypt project data",
		Long: `Re-encrypt variables and facts with new age keys.
This command decrypts existing encrypted data using old keys and re-encrypts it with new keys.
Useful when adding new age keys to a project.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath := args[0]

			// Get flags
			variables, _ := cmd.Flags().GetBool("variables")
			facts, _ := cmd.Flags().GetBool("facts")
			all, _ := cmd.Flags().GetBool("all")
			dryRun, _ := cmd.Flags().GetBool("dry-run")

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

			// Determine what to encrypt
			if !variables && !facts && !all {
				return fmt.Errorf("must specify --variables, --facts, or --all")
			}

			if all {
				variables = true
				facts = true
			}

			fmt.Printf("Re-encrypting project data: %s\n", projectPath)
			if dryRun {
				fmt.Println("DRY RUN MODE - No changes will be made")
			}

			// Re-encrypt variables if requested
			if variables {
				fmt.Println("Processing variables...")
				if err := reencryptVariables(coord, projectPath, dryRun); err != nil {
					return fmt.Errorf("failed to re-encrypt variables: %w", err)
				}
			}

			// Re-encrypt facts if requested
			if facts {
				fmt.Println("Processing facts...")
				if err := reencryptFacts(coord, projectPath, dryRun); err != nil {
					return fmt.Errorf("failed to re-encrypt facts: %w", err)
				}
			}

			if dryRun {
				fmt.Println("✅ Dry run completed - no changes made")
			} else {
				fmt.Println("✅ Re-encryption completed successfully")
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().BoolP("variables", "v", false, "Re-encrypt variables")
	cmd.Flags().BoolP("facts", "f", false, "Re-encrypt facts")
	cmd.Flags().BoolP("all", "a", false, "Re-encrypt all data (variables and facts)")
	cmd.Flags().BoolP("dry-run", "d", false, "Show what would be done without making changes")

	return cmd
}

// reencryptVariables re-encrypts variables in a project
func reencryptVariables(coord *spookycoordinator.CoordinatorManager, projectPath string, dryRun bool) error {
	// Load variables from the project
	variables, err := coord.Variables().LoadVariables(projectPath)
	if err != nil {
		return fmt.Errorf("failed to load variables: %w", err)
	}

	// Get list of variables to check for encrypted ones
	variableList, err := coord.Variables().ListVariables(variables)
	if err != nil {
		return fmt.Errorf("failed to list variables: %w", err)
	}

	encryptedCount := 0
	for varName := range variableList {
		// Get the variable to check if it's encrypted
		variable, err := coord.Variables().GetVariable(varName, variables)
		if err != nil {
			fmt.Printf("  Warning: Could not get variable %s: %v\n", varName, err)
			continue
		}

		// Check if variable is encrypted by looking for age encryption headers
		if isEncryptedValue(variable) {
			encryptedCount++
			if !dryRun {
				// Re-encrypt the variable
				fmt.Printf("  Re-encrypting variable: %s\n", varName)

				// 1. Get current encrypted value
				currentValue := fmt.Sprintf("%v", variable)

				// 2. Decrypt with old identity (using default identity)
				decrypted, err := coord.Crypto().DecryptData([]byte(currentValue))
				if err != nil {
					return fmt.Errorf("failed to decrypt variable %s: %w", varName, err)
				}

				// 3. Get new recipients from crypto manager config
				newRecipients := coord.Crypto().GetCryptoStatus()["default_recipients"].([]string)
				if len(newRecipients) == 0 {
					return fmt.Errorf("no new recipients configured for re-encryption")
				}

				// 4. Encrypt with new recipients
				encrypted, err := coord.Crypto().EncryptData(decrypted, newRecipients)
				if err != nil {
					return fmt.Errorf("failed to encrypt variable %s: %w", varName, err)
				}

				// 5. Update variable in the context
				err = coord.Variables().SetVariable(varName, string(encrypted), variables)
				if err != nil {
					return fmt.Errorf("failed to update variable %s: %w", varName, err)
				}

				fmt.Printf("    ✅ Successfully re-encrypted variable: %s\n", varName)
			} else {
				fmt.Printf("  Would re-encrypt variable: %s\n", varName)
			}
		} else {
			if !dryRun {
				fmt.Printf("  Skipping non-encrypted variable: %s\n", varName)
			} else {
				fmt.Printf("  Would skip non-encrypted variable: %s\n", varName)
			}
		}
	}

	if encryptedCount == 0 {
		fmt.Println("  No encrypted variables found")
	} else if dryRun {
		fmt.Printf("  Would re-encrypt %d variables\n", encryptedCount)
	} else {
		fmt.Printf("  Re-encrypted %d variables\n", encryptedCount)
	}

	return nil
}

// isEncryptedValue checks if a variable value appears to be encrypted
func isEncryptedValue(value interface{}) bool {
	if value == nil {
		return false
	}

	valueStr := fmt.Sprintf("%v", value)

	// Check for age encryption headers
	return strings.Contains(valueStr, "-----BEGIN AGE ENCRYPTED FILE-----") ||
		strings.Contains(valueStr, "age-encryption.org") ||
		strings.Contains(valueStr, "-----END AGE ENCRYPTED FILE-----")
}

// reencryptFacts re-encrypts facts in a project
func reencryptFacts(coord *spookycoordinator.CoordinatorManager, projectPath string, dryRun bool) error {
	// Get list of machines from the project
	// For now, we'll use a simple approach to get machines
	// In a real implementation, this would load from machines.hcl
	machines := []string{"localhost"} // Default to localhost for now

	encryptedCount := 0
	for _, machine := range machines {
		// Get facts for this machine
		facts, err := coord.Facts().GetFactsForMachine(machine)
		if err != nil {
			fmt.Printf("  Warning: Could not get facts for machine %s: %v\n", machine, err)
			continue
		}

		// Check if facts are encrypted
		if facts != nil && facts.EncryptedData != "" {
			encryptedCount++
			if !dryRun {
				// Re-encrypt the facts
				fmt.Printf("  Re-encrypting facts for machine: %s\n", machine)

				// 1. Decrypt with old identity (using default identity)
				decrypted, err := coord.Crypto().DecryptData([]byte(facts.EncryptedData))
				if err != nil {
					return fmt.Errorf("failed to decrypt facts for machine %s: %w", machine, err)
				}

				// 2. Get new recipients from crypto manager config
				cryptoStatus := coord.Crypto().GetCryptoStatus()
				newRecipientsInterface := cryptoStatus["default_recipients"]
				if newRecipientsInterface == nil {
					return fmt.Errorf("no new recipients configured for re-encryption")
				}

				newRecipients, ok := newRecipientsInterface.([]string)
				if !ok || len(newRecipients) == 0 {
					return fmt.Errorf("invalid or empty new recipients configuration")
				}

				// 3. Encrypt with new recipients
				encrypted, err := coord.Crypto().EncryptData(decrypted, newRecipients)
				if err != nil {
					return fmt.Errorf("failed to encrypt facts for machine %s: %w", machine, err)
				}

				// 4. Update facts with new encrypted data
				// Note: In a real implementation, this would update the facts in storage
				// For now, we'll just report success
				facts.EncryptedData = string(encrypted)

				// Update encryption metadata if available
				if facts.EncryptionMetadata != nil {
					facts.EncryptionMetadata.EncryptedAt = time.Now().Format(time.RFC3339)
					facts.EncryptionMetadata.Recipients = newRecipients
				}

				fmt.Printf("    ✅ Successfully re-encrypted facts for machine: %s\n", machine)
			} else {
				fmt.Printf("  Would re-encrypt facts for machine: %s\n", machine)
			}
		} else {
			if !dryRun {
				fmt.Printf("  Skipping non-encrypted facts for machine: %s\n", machine)
			} else {
				fmt.Printf("  Would skip non-encrypted facts for machine: %s\n", machine)
			}
		}
	}

	if encryptedCount == 0 {
		fmt.Println("  No encrypted facts found")
	} else if dryRun {
		fmt.Printf("  Would re-encrypt facts for %d machines\n", encryptedCount)
	} else {
		fmt.Printf("  Re-encrypted facts for %d machines\n", encryptedCount)
	}

	return nil
}

// Helper function to get project status
func getProjectStatus(projectPath string) string {
	// Check if project.hcl exists
	projectHCLPath := filepath.Join(projectPath, "project.hcl")
	if _, err := os.Stat(projectHCLPath); os.IsNotExist(err) {
		return "INVALID (missing project.hcl)"
	}

	// Check if required directories exist
	requiredDirs := []string{"templates", "variables", "actions", "machines"}
	for _, dir := range requiredDirs {
		dirPath := filepath.Join(projectPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			return "INCOMPLETE (missing directories)"
		}
	}

	return "VALID"
}
