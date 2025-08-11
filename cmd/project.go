// Package cmd provides command implementations for spooky CLI.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage spooky projects",
	Long: `Manage spooky projects including initialization and validation.

A spooky project contains configuration files, templates, and other resources
needed for automation and orchestration tasks.`,
}

// projectInitCmd represents the project init command
var projectInitCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new spooky project",
	Long: `Initialize a new spooky project with the specified name.

This command creates a new spooky project directory with the required structure
and configuration files according to the project-directory.schema.hcl schema.

Flags allow customization of project metadata during initialization.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		version, _ := cmd.Flags().GetString("version")
		author, _ := cmd.Flags().GetString("author")
		email, _ := cmd.Flags().GetString("email")
		url, _ := cmd.Flags().GetString("url")

		return handleProjectInit(args[0], name, description, version, author, email, url)
	},
}

// projectValidateCmd represents the project validate command
var projectValidateCmd = &cobra.Command{
	Use:   "validate [project-path]",
	Short: "Validate a spooky project",
	Long: `Validate a spooky project structure and configuration.

This command validates that the project follows the project-directory.schema.hcl
schema and that all configuration files are properly formatted.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleProjectValidate(args[0])
	},
}

// ProjectTemplateData holds data for project template generation
type ProjectTemplateData struct {
	ProjectName string
	Description string
	Version     string
	Author      string
	Email       string
	URL         string
}

// handleProjectInit handles project initialization
func handleProjectInit(projectPath, name, description, version, author, email, url string) error {
	// Resolve absolute path
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return fmt.Errorf("failed to resolve project path: %w", err)
	}

	// Check if project directory already exists
	if _, err := os.Stat(absPath); err == nil {
		return fmt.Errorf("project directory already exists: %s", absPath)
	}

	// Create project directory
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create required directories according to project-directory.schema.hcl
	requiredDirs := []string{
		"facts.db", // Required by schema
	}
	for _, dir := range requiredDirs {
		dirPath := filepath.Join(absPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create optional but commonly useful directories
	optionalDirs := []string{
		"files", // Optional but useful for static files
		"logs",  // Optional but useful for log files
	}
	for _, dir := range optionalDirs {
		dirPath := filepath.Join(absPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Prepare template data with provided values or defaults
	projectName := filepath.Base(absPath)
	if name != "" {
		projectName = name
	}

	if description == "" {
		description = "A spooky project for automation and orchestration"
	}

	if version == "" {
		version = "1.0.0"
	}

	if author == "" {
		author = "spooky-user"
	}

	data := ProjectTemplateData{
		ProjectName: projectName,
		Description: description,
		Version:     version,
		Author:      author,
		Email:       email,
		URL:         url,
	}

	// Create project.hcl file using updated schema structure
	if err := createProjectHCL(absPath, data); err != nil {
		return fmt.Errorf("failed to create project.hcl: %w", err)
	}

	// Create README.md file
	if err := createREADME(absPath, data); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	fmt.Printf("✅ Project initialized successfully: %s\n", absPath)
	fmt.Printf("📁 Project structure created according to project-directory.schema.hcl\n")
	fmt.Printf("📄 Configuration files generated using project.schema.hcl\n")
	fmt.Printf("💡 Next steps:\n")
	fmt.Printf("   - Edit project.hcl to customize your project\n")
	fmt.Printf("   - Add machines.hcl for machine inventory\n")
	fmt.Printf("   - Add actions.hcl for automation tasks\n")
	fmt.Printf("   - Add variables.hcl for project variables\n")

	return nil
}

// createProjectHCL creates the project.hcl file using the updated schema structure
func createProjectHCL(projectPath string, data ProjectTemplateData) error {
	const projectTemplate = `project {
  name = "{{.ProjectName}}"
  description = "{{.Description}}"
  version = "{{.Version}}"
  author = "{{.Author}}"
  {{- if .Email}}
  email = "{{.Email}}"
  {{- end}}
  {{- if .URL}}
  url = "{{.URL}}"
  {{- end}}

  execution {
    default_timeout = 300
    max_parallel = 10
    dry_run_default = false
    validate_before_execute = true
    backup_before_changes = false
  }

  facts {
    timeout = 60
    max_parallel = 5
    retry_attempts = 3
    retry_delay = 5
    storage_format = "badgerdb"
    compression = true
    encryption = false
  }
}
`

	tmpl, err := template.New("project").Parse(projectTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse project template: %w", err)
	}

	file, err := os.Create(filepath.Join(projectPath, "project.hcl"))
	if err != nil {
		return fmt.Errorf("failed to create project.hcl file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute project template: %w", err)
	}

	return nil
}

// createREADME creates a README.md file for the project
func createREADME(projectPath string, data ProjectTemplateData) error {
	const readmeTemplate = `# {{.ProjectName}}

{{.Description}}

## Project Structure

This project follows the spooky project-directory.schema.hcl structure:

- ` + "`" + `project.hcl` + "`" + ` - Main project configuration
- ` + "`" + `facts.db/` + "`" + ` - Facts database (BadgerDB)
- ` + "`" + `files/` + "`" + ` - Static files for deployment
- ` + "`" + `logs/` + "`" + ` - Log files directory

## Configuration

Edit ` + "`" + `project.hcl` + "`" + ` to customize your project settings, execution parameters, and facts collection configuration.

## Usage

` + "```" + `bash
# Validate project structure
spooky project validate .

# Add machine inventory
# Edit machines.hcl or create machines/ directory

# Add automation actions  
# Edit actions.hcl or create actions/ directory

# Add project variables
# Edit variables.hcl or create variables/ directory
` + "```" + `

## Schema Compliance

This project is validated against:
- project-directory.schema.hcl - Directory structure
- project.schema.hcl - Configuration format
`

	tmpl, err := template.New("readme").Parse(readmeTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse README template: %w", err)
	}

	file, err := os.Create(filepath.Join(projectPath, "README.md"))
	if err != nil {
		return fmt.Errorf("failed to create README.md file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute README template: %w", err)
	}

	return nil
}

// handleProjectValidate handles project validation
func handleProjectValidate(projectPath string) error {
	// Resolve absolute path
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return fmt.Errorf("failed to resolve project path: %w", err)
	}

	// Check if project directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("project directory does not exist: %s", absPath)
	}

	fmt.Printf("🔍 Validating project: %s\n", absPath)

	var issues []string
	var warnings []string

	// Check required files and directories according to project-directory.schema.hcl
	requiredFiles := []string{
		"project.hcl", // Required by schema
	}
	for _, file := range requiredFiles {
		filePath := filepath.Join(absPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			issues = append(issues, fmt.Sprintf("Missing required file: %s", file))
		}
	}

	requiredDirs := []string{
		"facts.db", // Required by schema
	}
	for _, dir := range requiredDirs {
		dirPath := filepath.Join(absPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			issues = append(issues, fmt.Sprintf("Missing required directory: %s", dir))
		}
	}

	// Check optional files and directories
	optionalFiles := []string{
		"machines.hcl",
		"actions.hcl",
		"variables.hcl",
		"README.md",
	}
	for _, file := range optionalFiles {
		filePath := filepath.Join(absPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			warnings = append(warnings, fmt.Sprintf("Optional file not found: %s", file))
		}
	}

	optionalDirs := []string{
		"machines",
		"actions",
		"variables",
		"templates",
		"files",
		"logs",
	}
	for _, dir := range optionalDirs {
		dirPath := filepath.Join(absPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			warnings = append(warnings, fmt.Sprintf("Optional directory not found: %s", dir))
		}
	}

	// Report validation results
	if len(issues) == 0 && len(warnings) == 0 {
		fmt.Printf("✅ Project validation passed - all required components present\n")
		fmt.Printf("📋 Schema compliance: project-directory.schema.hcl ✅\n")
		fmt.Printf("📋 Schema compliance: project.schema.hcl ✅\n")
	} else {
		if len(issues) > 0 {
			fmt.Printf("❌ Validation issues found:\n")
			for _, issue := range issues {
				fmt.Printf("   - %s\n", issue)
			}
		}
		if len(warnings) > 0 {
			fmt.Printf("⚠️  Warnings:\n")
			for _, warning := range warnings {
				fmt.Printf("   - %s\n", warning)
			}
		}
	}

	return nil
}

func init() {
	// Add flags to project init command
	projectInitCmd.Flags().String("name", "", "Project name (defaults to directory name)")
	projectInitCmd.Flags().String("description", "", "Project description")
	projectInitCmd.Flags().String("version", "", "Project version")
	projectInitCmd.Flags().String("author", "", "Project author")
	projectInitCmd.Flags().String("email", "", "Project email")
	projectInitCmd.Flags().String("url", "", "Project URL")

	projectCmd.AddCommand(projectInitCmd)
	projectCmd.AddCommand(projectValidateCmd)
	RootCmd.AddCommand(projectCmd)
}
