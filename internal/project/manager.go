// Package project provides project management functionality for spooky.
package project

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	spookyinterfaces "spooky/internal/interfaces"
	spookyschemas "spooky/internal/schemas"
	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesproject "spooky/internal/types/project"
	spookytypesschemas "spooky/internal/types/schemas"
)

// Manager implements the ProjectManager interface
type Manager struct {
	logger    spookytypeslogging.Logger
	validator spookyinterfaces.ProjectValidator
	loader    spookyinterfaces.ProjectLoader
}

// NewManager creates a new ProjectManager instance
func NewManager(
	logger spookytypeslogging.Logger,
	validator spookyinterfaces.ProjectValidator,
	loader spookyinterfaces.ProjectLoader,
) spookyinterfaces.ProjectManager {
	return &Manager{
		logger:    logger,
		validator: validator,
		loader:    loader,
	}
}

// Initialize initializes a new project using schema-driven approach
func (m *Manager) Initialize(_ context.Context, projectPath string) (*spookytypes.Project, error) {
	m.logger.Info("Initializing new project with schema-driven approach", map[string]interface{}{
		"project_path": projectPath,
	})

	// Resolve absolute path
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve project path: %w", err)
	}

	// Check if project directory already exists
	if _, err := os.Stat(absPath); err == nil {
		return nil, fmt.Errorf("project directory already exists: %s", absPath)
	}

	// Create project directory
	if err := os.MkdirAll(absPath, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create project directory: %w", err)
	}

	// Parse project directory schema
	schemaParser := spookyschemas.NewProjectDirectorySchemaParser(m.logger)
	schema, err := schemaParser.ParseProjectDirectorySchema()
	if err != nil {
		return nil, fmt.Errorf("failed to parse project directory schema: %w", err)
	}

	// Create project structure based on schema
	if err := m.createProjectFromSchema(absPath, schema); err != nil {
		return nil, fmt.Errorf("failed to create project from schema: %w", err)
	}

	// Create default project configuration
	projectName := filepath.Base(absPath)
	project := &spookytypes.Project{
		Path: absPath,
		Config: &spookytypesproject.Config{
			Name:        projectName,
			Description: "A spooky project for automation and orchestration",
			Metadata: &spookytypesproject.Metadata{
				Version: "1.0.0",
				Author:  "spooky-user",
			},
			Settings: &spookytypesproject.Settings{
				ParallelWorkers:   10,
				TimeoutSeconds:    300,
				LogLevel:          "info",
				DefaultDryRun:     false,
				ValidateBeforeRun: true,
				MaxRetries:        3,
				RetryDelaySeconds: 5,
			},
		},
	}

	// Create project.hcl file
	if err := m.createProjectHCL(project); err != nil {
		return nil, fmt.Errorf("failed to create project.hcl: %w", err)
	}

	// Create README.md file
	if err := m.createREADME(project); err != nil {
		return nil, fmt.Errorf("failed to create README.md: %w", err)
	}

	m.logger.Info("Project initialized successfully with schema-driven approach", map[string]interface{}{
		"project_path": absPath,
		"project_name": projectName,
	})

	return project, nil
}

// createProjectFromSchema creates project structure based on the schema
func (m *Manager) createProjectFromSchema(projectPath string, schema *spookyschemas.ProjectDirectorySchema) error {
	m.logger.Debug("Creating project structure from schema", map[string]interface{}{
		"project_path": projectPath,
	})

	// Create required directories from schema
	for _, dir := range schema.GetRequiredDirectories() {
		dirPath := filepath.Join(projectPath, dir.Name)
		if err := os.MkdirAll(dirPath, 0o755); err != nil {
			return fmt.Errorf("failed to create required directory %s: %w", dir.Name, err)
		}
		m.logger.Debug("Created required directory", map[string]interface{}{
			"directory": dir.Name,
			"path":      dirPath,
		})
	}

	// Create optional directories from schema (with user choice)
	for _, dir := range schema.GetOptionalDirectories() {
		if schema.ShouldCreateOptionalDirectory(dir) {
			dirPath := filepath.Join(projectPath, dir.Name)
			if err := os.MkdirAll(dirPath, 0o755); err != nil {
				return fmt.Errorf("failed to create optional directory %s: %w", dir.Name, err)
			}
			m.logger.Debug("Created optional directory", map[string]interface{}{
				"directory":   dir.Name,
				"path":        dirPath,
				"description": dir.Description,
			})
		}
	}

	// Create required files from schema
	for _, file := range schema.GetRequiredFiles() {
		if err := m.createFileFromSchema(projectPath, file); err != nil {
			return fmt.Errorf("failed to create required file %s: %w", file.Name, err)
		}
	}

	// Create optional files from schema (with user choice)
	for _, file := range schema.GetOptionalFiles() {
		if schema.ShouldCreateOptionalFile(file) {
			if err := m.createFileFromSchema(projectPath, file); err != nil {
				return fmt.Errorf("failed to create optional file %s: %w", file.Name, err)
			}
		}
	}

	m.logger.Info("Project structure created from schema successfully", map[string]interface{}{
		"project_path": projectPath,
	})

	return nil
}

// createFileFromSchema creates a file based on schema definition
func (m *Manager) createFileFromSchema(projectPath string, file spookyschemas.SchemaFile) error {
	filePath := filepath.Join(projectPath, file.Name)

	// Skip if file already exists
	if _, err := os.Stat(filePath); err == nil {
		m.logger.Debug("File already exists, skipping", map[string]interface{}{
			"file": file.Name,
			"path": filePath,
		})
		return nil
	}

	// Create file based on type
	switch file.Name {
	case "README.md":
		// README.md will be created by createREADME method
		return nil
	case "project.hcl":
		// project.hcl will be created by createProjectHCL method
		return nil
	case "machines.hcl":
		return m.createMachinesHCL(filePath)
	case "actions.hcl":
		return m.createActionsHCL(filePath)
	case "variables.hcl":
		return m.createVariablesHCL(filePath)
	case "recipients.txt":
		return m.createRecipientsTXT(filePath)
	default:
		// Create empty file for other types
		if err := os.WriteFile(filePath, []byte(""), 0o600); err != nil {
			return fmt.Errorf("failed to create file %s: %w", file.Name, err)
		}
		m.logger.Debug("Created file from schema", map[string]interface{}{
			"file":        file.Name,
			"path":        filePath,
			"description": file.Description,
		})
		return nil
	}
}

// createMachinesHCL creates a default machines.hcl file
func (m *Manager) createMachinesHCL(filePath string) error {
	content := `# Machine inventory for spooky project
# Define your machines here

machines {
  # Example machine definition
  # machine "web-server" {
  #   hostname = "web.example.com"
  #   port = 22
  #   user = "admin"
  #   
  #   authentication {
  #     method = "ssh_key"
  #     key_path = "~/.ssh/id_rsa"
  #   }
  #   
  #   tags = ["web", "production"]
  # }
}
`
	return os.WriteFile(filePath, []byte(content), 0o600)
}

// createActionsHCL creates a default actions.hcl file
func (m *Manager) createActionsHCL(filePath string) error {
	content := `# Actions for spooky project
# Define your actions here

actions {
  # Example action definition
  # action "deploy-web" {
  #   description = "Deploy web application"
  #   
  #   machines = ["web-server"]
  #   parallel = true
  #   
  #   template {
  #     source = "templates/deploy.sh.tmpl"
  #     destination = "/tmp/deploy.sh"
  #     permissions = "0755"
  #   }
  #   
  #   command = "/tmp/deploy.sh"
  # }
}
`
	return os.WriteFile(filePath, []byte(content), 0o600)
}

// createVariablesHCL creates a default variables.hcl file
func (m *Manager) createVariablesHCL(filePath string) error {
	content := `# Variables for spooky project
# Define your variables here

variables {
  # Example variable definitions
  # app_version = "1.0.0"
  # environment = "production"
  # backup_retention_days = 30
}
`
	return os.WriteFile(filePath, []byte(content), 0o600)
}

// createRecipientsTXT creates a default recipients.txt file
func (m *Manager) createRecipientsTXT(filePath string) error {
	content := `# Age recipients for this project
# Add your age public keys here (one per line)
# Example: age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p
`
	return os.WriteFile(filePath, []byte(content), 0o600)
}

// Load loads a project from the given path
func (m *Manager) Load(ctx context.Context, projectPath string) (*spookytypes.Project, error) {
	m.logger.Info("Loading project", map[string]interface{}{
		"project_path": projectPath,
	})

	// Use the loader to load the project
	project, err := m.loader.LoadProject(ctx, projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load project: %w", err)
	}

	m.logger.Info("Project loaded successfully", map[string]interface{}{
		"project_path": projectPath,
		"project_name": project.Config.Name,
	})

	return project, nil
}

// Validate validates a project
func (m *Manager) Validate(ctx context.Context, project *spookytypes.Project) (*spookytypesschemas.ValidationResult, error) {
	m.logger.Info("Validating project", map[string]interface{}{
		"project_path": project.Path,
		"project_name": project.Config.Name,
	})

	// Use the validator to validate the project
	result, err := m.validator.ValidateProject(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("failed to validate project: %w", err)
	}

	if result.Valid {
		m.logger.Info("Project validation passed", map[string]interface{}{
			"project_path": project.Path,
		})
	} else {
		m.logger.Warn("Project validation failed", map[string]interface{}{
			"project_path": project.Path,
			"errors":       len(result.Errors),
			"warnings":     len(result.Warnings),
		})
	}

	return result, nil
}

// Save saves a project to disk
func (m *Manager) Save(_ context.Context, project *spookytypes.Project) error {
	m.logger.Info("Saving project", map[string]interface{}{
		"project_path": project.Path,
		"project_name": project.Config.Name,
	})

	// Create project.hcl file
	if err := m.createProjectHCL(project); err != nil {
		return fmt.Errorf("failed to save project.hcl: %w", err)
	}

	m.logger.Info("Project saved successfully", map[string]interface{}{
		"project_path": project.Path,
	})

	return nil
}

// Delete deletes a project
func (m *Manager) Delete(_ context.Context, projectPath string) error {
	m.logger.Info("Deleting project", map[string]interface{}{
		"project_path": projectPath,
	})

	// Resolve absolute path
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return fmt.Errorf("failed to resolve project path: %w", err)
	}

	// Check if project directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("project directory does not exist: %s", absPath)
	}

	// Remove the project directory
	if err := os.RemoveAll(absPath); err != nil {
		return fmt.Errorf("failed to delete project directory: %w", err)
	}

	m.logger.Info("Project deleted successfully", map[string]interface{}{
		"project_path": absPath,
	})

	return nil
}

// createProjectHCL creates the project.hcl file using schema-driven generation
func (m *Manager) createProjectHCL(project *spookytypes.Project) error {
	m.logger.Debug("Creating project.hcl using schema-driven generation", map[string]interface{}{
		"project_path": project.Path,
	})

	// Load the project schema
	schemaManager := spookyschemas.NewManager(m.logger)
	projectSchema, err := schemaManager.Load("internal/schemas/schemas/structure/project.hcl")
	if err != nil {
		return fmt.Errorf("failed to load project schema: %w", err)
	}

	// Generate HCL content from schema and project data
	hclContent, err := m.generateProjectHCLFromSchema(projectSchema, project)
	if err != nil {
		return fmt.Errorf("failed to generate project.hcl from schema: %w", err)
	}

	// Write the generated content to file
	filePath := filepath.Join(project.Path, "project.hcl")
	if err := os.WriteFile(filePath, []byte(hclContent), 0o644); err != nil {
		return fmt.Errorf("failed to write project.hcl file: %w", err)
	}

	m.logger.Debug("Generated project.hcl from schema", map[string]interface{}{
		"file_path":      filePath,
		"content_length": len(hclContent),
	})

	return nil
}

// generateProjectHCLFromSchema generates HCL content based on the project schema
func (m *Manager) generateProjectHCLFromSchema(schema *spookytypesschemas.Schema, project *spookytypes.Project) (string, error) {
	var content strings.Builder

	// Start project block
	content.WriteString("project {\n")

	// Add required fields first
	if project.Config.Name != "" {
		content.WriteString(fmt.Sprintf("  name = %q\n", project.Config.Name))
	}

	// Add optional fields based on schema structure
	if project.Config.Description != "" {
		content.WriteString(fmt.Sprintf("  description = %q\n", project.Config.Description))
	}

	// Add metadata fields as direct attributes (not in metadata block)
	if project.Config.Metadata != nil {
		if project.Config.Metadata.Version != "" {
			content.WriteString(fmt.Sprintf("  version = %q\n", project.Config.Metadata.Version))
		}
		if project.Config.Metadata.Author != "" {
			content.WriteString(fmt.Sprintf("  author = %q\n", project.Config.Metadata.Author))
		}
		if project.Config.Metadata.Email != "" {
			content.WriteString(fmt.Sprintf("  email = %q\n", project.Config.Metadata.Email))
		}
		if project.Config.Metadata.URL != "" {
			content.WriteString(fmt.Sprintf("  url = %q\n", project.Config.Metadata.URL))
		}
	}

	// Add settings block if settings exist
	if project.Config.Settings != nil {
		content.WriteString("  settings {\n")

		if project.Config.Settings.ParallelWorkers > 0 {
			content.WriteString(fmt.Sprintf("    parallel_workers = %d\n", project.Config.Settings.ParallelWorkers))
		}
		if project.Config.Settings.TimeoutSeconds > 0 {
			content.WriteString(fmt.Sprintf("    timeout_seconds = %d\n", project.Config.Settings.TimeoutSeconds))
		}
		if project.Config.Settings.LogLevel != "" {
			content.WriteString(fmt.Sprintf("    log_level = %q\n", project.Config.Settings.LogLevel))
		}
		if project.Config.Settings.DefaultDryRun {
			content.WriteString("    default_dry_run = true\n")
		}
		if project.Config.Settings.ValidateBeforeRun {
			content.WriteString("    validate_before_run = true\n")
		}
		if project.Config.Settings.MaxRetries > 0 {
			content.WriteString(fmt.Sprintf("    max_retries = %d\n", project.Config.Settings.MaxRetries))
		}
		if project.Config.Settings.RetryDelaySeconds > 0 {
			content.WriteString(fmt.Sprintf("    retry_delay_seconds = %d\n", project.Config.Settings.RetryDelaySeconds))
		}

		content.WriteString("  }\n")
	}

	// Close project block
	content.WriteString("}\n")

	return content.String(), nil
}

// createREADME creates a README.md file for the project
func (m *Manager) createREADME(project *spookytypes.Project) error {
	const readmeTemplate = `# {{.Name}}

{{.Description}}

## Project Structure

This project follows the spooky project-directory.schema.hcl structure:

- ` + "`" + `project.hcl` + "`" + ` - Main project configuration
- ` + "`" + `files/` + "`" + ` - Static files for deployment
- ` + "`" + `logs/` + "`" + ` - Log files directory

## Configuration

Edit ` + "`" + `project.hcl` + "`" + ` to customize your project settings, run parameters, and facts collection configuration.

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

	file, err := os.Create(filepath.Join(project.Path, "README.md"))
	if err != nil {
		return fmt.Errorf("failed to create README.md file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, project.Config); err != nil {
		return fmt.Errorf("failed to run README template: %w", err)
	}

	return nil
}
