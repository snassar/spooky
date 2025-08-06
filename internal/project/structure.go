package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ProjectStructureEngine manages project directory structure
type ProjectStructureEngine struct {
	validator *ProjectValidator
}

// NewProjectStructureEngine creates a new project structure engine
func NewProjectStructureEngine() *ProjectStructureEngine {
	return &ProjectStructureEngine{
		validator: NewProjectValidator(),
	}
}

// ValidateProjectStructure validates project directory structure
func (pse *ProjectStructureEngine) ValidateProjectStructure(projectPath string) *ValidationResult {
	result := &ValidationResult{
		Valid: true,
	}

	// Check if project path exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "path",
			Message:  "Project path does not exist",
			Value:    projectPath,
			Severity: "error",
		})
		return result
	}

	// Validate required files
	if err := pse.validateRequiredFiles(projectPath, result); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, *err)
	}

	// Validate optional directories
	pse.validateOptionalDirectories(projectPath, result)

	// Validate cross-file rules
	pse.validateCrossFileRules(projectPath, result)

	return result
}

// CreateProjectStructure creates a new project directory structure
func (pse *ProjectStructureEngine) CreateProjectStructure(projectPath string, project *Project) error {
	// Create project directory if it doesn't exist
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create required directories
	requiredDirs := []string{"facts.db"}
	for _, dir := range requiredDirs {
		dirPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create required directory %s: %w", dir, err)
		}
	}

	// Create optional directories based on project structure
	if project.Structure != nil {
		optionalDirs := []string{
			project.Structure.TemplatesDir,
			project.Structure.DataDir,
			project.Structure.ScriptsDir,
			project.Structure.LogsDir,
			project.Structure.BackupsDir,
		}

		for _, dir := range optionalDirs {
			if dir != "" {
				dirPath := filepath.Join(projectPath, dir)
				if err := os.MkdirAll(dirPath, 0755); err != nil {
					return fmt.Errorf("failed to create optional directory %s: %w", dir, err)
				}
			}
		}
	}

	// Create project.hcl file if it doesn't exist
	projectHCLPath := filepath.Join(projectPath, "project.hcl")
	if _, err := os.Stat(projectHCLPath); os.IsNotExist(err) {
		if err := pse.createProjectHCL(projectHCLPath, project); err != nil {
			return fmt.Errorf("failed to create project.hcl: %w", err)
		}
	}

	// Create README.md if it doesn't exist
	readmePath := filepath.Join(projectPath, "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		if err := pse.createREADME(readmePath, project); err != nil {
			return fmt.Errorf("failed to create README.md: %w", err)
		}
	}

	return nil
}

// MaintainProjectStructure maintains existing project structure
func (pse *ProjectStructureEngine) MaintainProjectStructure(projectPath string, project *Project) error {
	// Validate existing structure
	result := pse.ValidateProjectStructure(projectPath)
	if !result.Valid {
		return fmt.Errorf("project structure validation failed: %v", result.Errors)
	}

	// Create missing optional directories
	if project.Structure != nil {
		optionalDirs := []string{
			project.Structure.TemplatesDir,
			project.Structure.DataDir,
			project.Structure.ScriptsDir,
			project.Structure.LogsDir,
			project.Structure.BackupsDir,
		}

		for _, dir := range optionalDirs {
			if dir != "" {
				dirPath := filepath.Join(projectPath, dir)
				if _, err := os.Stat(dirPath); os.IsNotExist(err) {
					if err := os.MkdirAll(dirPath, 0755); err != nil {
						return fmt.Errorf("failed to create missing directory %s: %w", dir, err)
					}
				}
			}
		}
	}

	return nil
}

// GetProjectStructure returns information about project structure
func (pse *ProjectStructureEngine) GetProjectStructure(projectPath string) (*ProjectStructureInfo, error) {
	info := &ProjectStructureInfo{
		Path: projectPath,
	}

	// Check required files
	requiredFiles := []string{"project.hcl"}
	for _, file := range requiredFiles {
		filePath := filepath.Join(projectPath, file)
		if _, err := os.Stat(filePath); err == nil {
			info.RequiredFiles = append(info.RequiredFiles, file)
		}
	}

	// Check optional directories
	optionalDirs := []string{"templates", "data", "scripts", "logs", "backups", "machines", "actions", "variables"}
	for _, dir := range optionalDirs {
		dirPath := filepath.Join(projectPath, dir)
		if _, err := os.Stat(dirPath); err == nil {
			info.OptionalDirectories = append(info.OptionalDirectories, dir)
		}
	}

	// Check facts database
	factsDBPath := filepath.Join(projectPath, "facts.db")
	if _, err := os.Stat(factsDBPath); err == nil {
		info.FactsDatabaseExists = true
	}

	return info, nil
}

// validateRequiredFiles validates required project files
func (pse *ProjectStructureEngine) validateRequiredFiles(projectPath string, result *ValidationResult) *ValidationError {
	// Check project.hcl
	projectHCLPath := filepath.Join(projectPath, "project.hcl")
	if _, err := os.Stat(projectHCLPath); os.IsNotExist(err) {
		return &ValidationError{
			Field:    "project.hcl",
			Message:  "Project configuration file not found",
			Value:    projectHCLPath,
			Severity: "error",
		}
	}

	return nil
}

// validateOptionalDirectories validates optional project directories
func (pse *ProjectStructureEngine) validateOptionalDirectories(projectPath string, result *ValidationResult) {
	optionalDirs := []string{"templates", "data", "scripts", "logs", "backups", "machines", "actions", "variables"}

	for _, dir := range optionalDirs {
		dirPath := filepath.Join(projectPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			result.Warnings = append(result.Warnings, ValidationError{
				Field:    dir,
				Message:  fmt.Sprintf("Optional directory '%s' not found", dir),
				Value:    dirPath,
				Severity: "warning",
			})
		}
	}
}

// validateCrossFileRules validates cross-file validation rules
func (pse *ProjectStructureEngine) validateCrossFileRules(projectPath string, result *ValidationResult) {
	// Check for machines file or directory
	machinesHCLPath := filepath.Join(projectPath, "machines.hcl")
	machinesDirPath := filepath.Join(projectPath, "machines")

	machinesHCLExists := false
	machinesDirExists := false

	if _, err := os.Stat(machinesHCLPath); err == nil {
		machinesHCLExists = true
	}
	if _, err := os.Stat(machinesDirPath); err == nil {
		machinesDirExists = true
	}

	if !machinesHCLExists && !machinesDirExists {
		result.Warnings = append(result.Warnings, ValidationError{
			Field:    "machines",
			Message:  "Neither machines.hcl file nor machines/ directory found",
			Severity: "warning",
		})
	}

	// Check for actions file or directory
	actionsHCLPath := filepath.Join(projectPath, "actions.hcl")
	actionsDirPath := filepath.Join(projectPath, "actions")

	actionsHCLExists := false
	actionsDirExists := false

	if _, err := os.Stat(actionsHCLPath); err == nil {
		actionsHCLExists = true
	}
	if _, err := os.Stat(actionsDirPath); err == nil {
		actionsDirExists = true
	}

	if !actionsHCLExists && !actionsDirExists {
		result.Warnings = append(result.Warnings, ValidationError{
			Field:    "actions",
			Message:  "Neither actions.hcl file nor actions/ directory found",
			Severity: "warning",
		})
	}

	// Check for variables file or directory
	variablesHCLPath := filepath.Join(projectPath, "variables.hcl")
	variablesDirPath := filepath.Join(projectPath, "variables")

	variablesHCLExists := false
	variablesDirExists := false

	if _, err := os.Stat(variablesHCLPath); err == nil {
		variablesHCLExists = true
	}
	if _, err := os.Stat(variablesDirPath); err == nil {
		variablesDirExists = true
	}

	if !variablesHCLExists && !variablesDirExists {
		result.Warnings = append(result.Warnings, ValidationError{
			Field:    "variables",
			Message:  "Neither variables.hcl file nor variables/ directory found",
			Severity: "warning",
		})
	}
}

// createProjectHCL creates a basic project.hcl file
func (pse *ProjectStructureEngine) createProjectHCL(path string, project *Project) error {
	content := fmt.Sprintf("project \"%s\" {\n", project.Name)

	if project.Description != "" {
		content += fmt.Sprintf("  description = \"%s\"\n", project.Description)
	}

	if project.Version != "" {
		content += fmt.Sprintf("  version = \"%s\"\n", project.Version)
	}

	if project.Environment != "" {
		content += fmt.Sprintf("  environment = \"%s\"\n", project.Environment)
	}

	// Add file references
	content += "  inventory_file = \"machines.hcl\"\n"
	content += "  actions_file = \"actions.hcl\"\n"

	// Add project settings
	if project.Execution != nil {
		if project.Execution.DefaultTimeout > 0 {
			content += fmt.Sprintf("  default_timeout = %d\n", project.Execution.DefaultTimeout)
		}
		content += fmt.Sprintf("  default_parallel = %t\n", project.Execution.MaxParallel > 0)
		content += fmt.Sprintf("  dry_run_default = %t\n", project.Execution.DryRunDefault)
		content += fmt.Sprintf("  validate_before_execute = %t\n", project.Execution.ValidateBeforeExecute)
		content += fmt.Sprintf("  backup_before_changes = %t\n", project.Execution.BackupBeforeChanges)
	} else {
		content += "  default_parallel = true\n"
	}

	// Add tags
	if len(project.Tags) > 0 {
		content += "  tags = {\n"
		for _, tag := range project.Tags {
			// Parse tag in format "key=value" or just "key"
			if idx := strings.Index(tag, "="); idx != -1 {
				key := tag[:idx]
				value := tag[idx+1:]
				content += fmt.Sprintf("    %s = \"%s\"\n", key, value)
			} else {
				content += fmt.Sprintf("    %s = \"true\"\n", tag)
			}
		}
		content += "  }\n"
	}

	// Add storage configuration
	content += "  \n  storage {\n"
	content += "    type = \"badgerdb\"\n"
	content += "    path = \"facts.db\"\n"
	content += "  }\n"

	// Add logging configuration
	content += "  \n  logging {\n"
	content += "    level = \"info\"\n"
	content += "    format = \"json\"\n"
	content += "    output = \"logs/spooky.log\"\n"
	content += "  }\n"

	// Add SSH configuration
	content += "  \n  ssh {\n"
	content += "    timeout = 30\n"
	content += "    keepalive_interval = 30\n"
	content += "    keepalive_count = 3\n"
	content += "    key_scan_timeout = 10\n"
	content += "    known_hosts_strict = false\n"
	content += "    connection_pool_size = 10\n"
	content += "  }\n"

	// Add isolation configuration
	if project.Isolation != nil {
		content += "  \n  isolation {\n"
		content += fmt.Sprintf("    enabled = %t\n", project.Isolation.Enabled)
		if project.Isolation.FactsScope != "" {
			content += fmt.Sprintf("    facts_scope = \"%s\"\n", project.Isolation.FactsScope)
		}
		if project.Isolation.VariablesScope != "" {
			content += fmt.Sprintf("    variables_scope = \"%s\"\n", project.Isolation.VariablesScope)
		}
		if project.Isolation.MachineAccess != "" {
			content += fmt.Sprintf("    machine_access = \"%s\"\n", project.Isolation.MachineAccess)
		}
		if len(project.Isolation.AllowedMachines) > 0 {
			content += "    allowed_machines = [\n"
			for _, machine := range project.Isolation.AllowedMachines {
				content += fmt.Sprintf("      \"%s\",\n", machine)
			}
			content += "    ]\n"
		}
		if len(project.Isolation.AllowedTags) > 0 {
			content += "    allowed_tags = [\n"
			for _, tag := range project.Isolation.AllowedTags {
				content += fmt.Sprintf("      \"%s\",\n", tag)
			}
			content += "    ]\n"
		}
		content += "  }\n"
	}

	content += "}\n"

	return os.WriteFile(path, []byte(content), 0644)
}

// createREADME creates a basic README.md file
func (pse *ProjectStructureEngine) createREADME(path string, project *Project) error {
	content := fmt.Sprintf("# %s\n\n", project.Name)
	content += fmt.Sprintf("%s\n\n", project.Description)
	content += "## Project Information\n\n"
	content += fmt.Sprintf("- **Name**: %s\n", project.Name)
	content += fmt.Sprintf("- **Version**: %s\n", project.Version)
	content += fmt.Sprintf("- **Environment**: %s\n", project.Environment)
	content += fmt.Sprintf("- **Description**: %s\n\n", project.Description)
	content += "## Project Structure\n\n"
	content += "This project follows the standard spooky project structure:\n\n"
	content += "- project.hcl - Project configuration\n"
	content += "- facts.db/ - Facts database\n"
	content += "- templates/ - Template files\n"
	content += "- data/ - Data files\n"
	content += "- scripts/ - Script files\n"
	content += "- logs/ - Log files\n"
	content += "- backups/ - Backup files\n"
	content += "- machines.hcl or machines/ - Machine inventory\n"
	content += "- actions.hcl or actions/ - Action definitions\n"
	content += "- variables.hcl or variables/ - Variable definitions\n\n"
	content += "## Usage\n\n"
	content += "Use spooky commands to manage this project:\n\n"
	content += "```bash\n"
	content += "# Validate project\n"
	content += fmt.Sprintf("spooky project validate %s\n\n", project.Name)
	content += "# Show project information\n"
	content += fmt.Sprintf("spooky project show %s\n\n", project.Name)
	content += "# Run actions\n"
	content += fmt.Sprintf("spooky actions run %s\n", project.Name)
	content += "```\n"

	return os.WriteFile(path, []byte(content), 0644)
}

// ProjectStructureInfo represents project structure information
type ProjectStructureInfo struct {
	Path                string   `json:"path"`
	RequiredFiles       []string `json:"required_files"`
	OptionalDirectories []string `json:"optional_directories"`
	FactsDatabaseExists bool     `json:"facts_database_exists"`
}
