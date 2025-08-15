// Package project provides project management functionality for spooky.
package project

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	spookyinterfaces "spooky/internal/interfaces"
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

// Initialize initializes a new project
func (m *Manager) Initialize(_ context.Context, projectPath string) (*spookytypes.Project, error) {
	m.logger.Info("Initializing new project", map[string]interface{}{
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

	// Create required directories according to project-directory.schema.hcl
	requiredDirs := []string{
		// No required directories
	}
	for _, dir := range requiredDirs {
		dirPath := filepath.Join(absPath, dir)
		if err := os.MkdirAll(dirPath, 0o755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create optional but commonly useful directories
	optionalDirs := []string{
		"files", // Optional but useful for static files
		"logs",  // Optional but useful for log files
	}
	for _, dir := range optionalDirs {
		dirPath := filepath.Join(absPath, dir)
		if err := os.MkdirAll(dirPath, 0o755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
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

	m.logger.Info("Project initialized successfully", map[string]interface{}{
		"project_path": absPath,
		"project_name": projectName,
	})

	return project, nil
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

// createProjectHCL creates the project.hcl file
func (m *Manager) createProjectHCL(project *spookytypes.Project) error {
	const projectTemplate = `project {
  name = "{{.Name}}"
  description = "{{.Description}}"
  {{- if .Metadata}}
  {{- if .Metadata.Version}}
  version = "{{.Metadata.Version}}"
  {{- end}}
  {{- if .Metadata.Author}}
  author = "{{.Metadata.Author}}"
  {{- end}}
  {{- if .Metadata.URL}}
  url = "{{.Metadata.URL}}"
  {{- end}}
  {{- end}}

  {{- if .Settings}}
  settings {
    {{- if .Settings.ParallelWorkers}}
    parallel_workers = {{.Settings.ParallelWorkers}}
    {{- end}}
    {{- if .Settings.TimeoutSeconds}}
    timeout_seconds = {{.Settings.TimeoutSeconds}}
    {{- end}}
    {{- if .Settings.LogLevel}}
    log_level = "{{.Settings.LogLevel}}"
    {{- end}}
    {{- if .Settings.DefaultDryRun}}
    default_dry_run = {{.Settings.DefaultDryRun}}
    {{- end}}
    {{- if .Settings.ValidateBeforeRun}}
    validate_before_run = {{.Settings.ValidateBeforeRun}}
    {{- end}}
    {{- if .Settings.MaxRetries}}
    max_retries = {{.Settings.MaxRetries}}
    {{- end}}
    {{- if .Settings.RetryDelaySeconds}}
    retry_delay_seconds = {{.Settings.RetryDelaySeconds}}
    {{- end}}
  }
  {{- end}}
}
`

	tmpl, err := template.New("project").Parse(projectTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse project template: %w", err)
	}

	file, err := os.Create(filepath.Join(project.Path, "project.hcl"))
	if err != nil {
		return fmt.Errorf("failed to create project.hcl file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, project.Config); err != nil {
		return fmt.Errorf("failed to run project template: %w", err)
	}

	return nil
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
