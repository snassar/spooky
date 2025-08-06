package project

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"spooky/internal/config/types"
	"spooky/internal/logging"
)

// ProjectManager coordinates all project system components
type ProjectManager struct {
	validator    *ProjectValidator
	structure    *ProjectStructureEngine
	identity     *ProjectIdentityManager
	configEngine *ProjectConfigurationEngine
	isolation    *ProjectIsolationEngine
	dependencies *ProjectDependenciesEngine
	portability  *ProjectPortabilityEngine

	// Integration managers - now handled by coordinator
}

// NewProjectManager creates a new project manager
func NewProjectManager() *ProjectManager {
	return &ProjectManager{
		validator:    NewProjectValidator(),
		structure:    NewProjectStructureEngine(),
		identity:     NewProjectIdentityManager(),
		configEngine: NewProjectConfigurationEngine(),
		isolation:    NewProjectIsolationEngine(),
		dependencies: NewProjectDependenciesEngine(),
		portability:  NewProjectPortabilityEngine(),
	}
}

// LoadProject loads a project from the specified path
func (pm *ProjectManager) LoadProject(projectPath string) (*Project, error) {
	logger := logging.GetLogger()

	// Validate project path using schema validation
	pathResult := pm.structure.ValidateProjectStructure(projectPath)
	if !pathResult.Valid {
		return nil, fmt.Errorf("project path validation failed: %v", pathResult.Errors)
	}

	// Load project.hcl file
	projectHCLPath := filepath.Join(projectPath, "project.hcl")
	if _, err := os.Stat(projectHCLPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("project.hcl file not found at %s", projectHCLPath)
	}

	// Parse project.hcl file using HCL parser with schema validation
	logger.Info("Loading project configuration", logging.String("path", projectHCLPath))

	// For now, create a basic project structure
	// In a real implementation, this would use the config manager to parse project.hcl
	project := &Project{
		Path:      projectPath,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set defaults for any missing fields
	project.SetDefaults()

	// Validate project using schema validation
	validationResult := pm.validator.ValidateProject(project)
	if !validationResult.Valid {
		return nil, fmt.Errorf("project validation failed: %v", validationResult.Errors)
	}

	logger.Info("Project loaded successfully",
		logging.String("name", project.Name),
		logging.String("path", project.Path))

	// Integration managers are now handled by the coordinator package

	return project, nil
}

// ApplyProjectIsolation applies isolation settings to a project
func (pm *ProjectManager) ApplyProjectIsolation(project *Project) (*IsolationContext, error) {
	return pm.isolation.ApplyIsolation(project)
}

// ValidateProjectIsolation validates project isolation configuration
func (pm *ProjectManager) ValidateProjectIsolation(project *Project) *IsolationResult {
	return pm.isolation.ValidateIsolation(project)
}

// GetProjectIsolationSummary returns a summary of project isolation settings
func (pm *ProjectManager) GetProjectIsolationSummary(project *Project) map[string]interface{} {
	return pm.isolation.GetIsolationSummary(project)
}

// CheckProjectMachineAccess checks if a machine is accessible for a project
func (pm *ProjectManager) CheckProjectMachineAccess(project *Project, machineName string, machineTags []string) bool {
	return pm.isolation.CheckMachineAccess(project, machineName, machineTags)
}

// convertProjectConfig converts config.ProjectConfig to project.Project
func (pm *ProjectManager) convertProjectConfig(config *types.ProjectConfig, projectPath string) *Project {
	project := &Project{
		Path:      projectPath,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Convert basic fields
	if config.Name != "" {
		project.Name = config.Name
	}
	if config.Description != "" {
		project.Description = config.Description
	}
	if config.Version != "" {
		project.Version = config.Version
	}
	if config.Environment != "" {
		project.Environment = config.Environment
	}

	// Convert tags from map[string]string to []string
	if len(config.Tags) > 0 {
		project.Tags = make([]string, 0, len(config.Tags))
		for key := range config.Tags {
			// Use just the key as the tag, since the validation expects simple tags
			project.Tags = append(project.Tags, key)
		}
	}

	// Convert execution configuration from basic fields
	if project.Execution == nil {
		project.Execution = &ProjectExecution{}
	}
	// Set fields from config.ProjectConfig
	project.Execution.DefaultTimeout = config.DefaultTimeout
	project.Execution.DryRunDefault = config.DryRunDefault
	project.Execution.ValidateBeforeExecute = config.ValidateBeforeExecute
	project.Execution.BackupBeforeChanges = config.BackupBeforeChanges

	// Set MaxParallel based on DefaultParallel (since config uses bool, project uses int)
	if config.DefaultParallel {
		project.Execution.MaxParallel = 10
	} else {
		project.Execution.MaxParallel = 1
	}

	// Convert isolation configuration
	if config.Isolation != nil {
		project.Isolation = &ProjectIsolation{
			Enabled:         config.Isolation.Enabled,
			FactsScope:      config.Isolation.FactsScope,
			VariablesScope:  config.Isolation.VariablesScope,
			MachineAccess:   config.Isolation.MachineAccess,
			AllowedMachines: config.Isolation.AllowedMachines,
			AllowedTags:     config.Isolation.AllowedTags,
		}
	}

	// Note: The config.ProjectConfig structure is simpler than project.Project
	// Additional fields like Author, Region, Contact, Structure, Dependencies,
	// and detailed Execution settings are not available in the
	// current config.ProjectConfig structure. These would need to be added
	// to the config system or handled separately.

	// Set defaults for missing fields
	project.SetDefaults()

	return project
}

// CreateProject creates a new project at the specified path
func (pm *ProjectManager) CreateProject(projectPath string, projectConfig *Project) error {
	// Generate project name if not provided
	if projectConfig.Name == "" {
		name, err := pm.identity.GenerateProjectName(projectPath)
		if err != nil {
			return fmt.Errorf("failed to generate project name: %w", err)
		}
		projectConfig.Name = name
	}

	// Validate project name
	nameResult := pm.identity.ValidateProjectName(projectConfig.Name, projectPath)
	if !nameResult.Valid {
		return fmt.Errorf("project name validation failed: %v", nameResult.Errors)
	}

	// Set project path
	projectConfig.Path = projectPath
	projectConfig.CreatedAt = time.Now()
	projectConfig.UpdatedAt = time.Now()

	// Set defaults
	projectConfig.SetDefaults()

	// Validate project configuration
	validationResult := pm.validator.ValidateProject(projectConfig)
	if !validationResult.Valid {
		return fmt.Errorf("project configuration validation failed: %v", validationResult.Errors)
	}

	// Create project structure
	if err := pm.structure.CreateProjectStructure(projectPath, projectConfig); err != nil {
		return fmt.Errorf("failed to create project structure: %w", err)
	}

	return nil
}

// ValidateProject validates a project configuration
func (pm *ProjectManager) ValidateProject(projectPath string) *ValidationResult {
	result := &ValidationResult{
		Valid: true,
	}

	// Validate project structure
	structureResult := pm.structure.ValidateProjectStructure(projectPath)
	if !structureResult.Valid {
		result.Valid = false
		result.Errors = append(result.Errors, structureResult.Errors...)
	}
	result.Warnings = append(result.Warnings, structureResult.Warnings...)

	// Load and validate project configuration
	project, err := pm.LoadProject(projectPath)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:    "project",
			Message:  err.Error(),
			Severity: "error",
		})
		return result
	}

	// Validate project configuration
	configResult := pm.validator.ValidateProject(project)
	if !configResult.Valid {
		result.Valid = false
		result.Errors = append(result.Errors, configResult.Errors...)
	}
	result.Warnings = append(result.Warnings, configResult.Warnings...)

	// Validate project identity
	identityResult := pm.identity.ValidateProjectName(project.Name, projectPath)
	if !identityResult.Valid {
		result.Valid = false
		result.Errors = append(result.Errors, identityResult.Errors...)
	}
	result.Warnings = append(result.Warnings, identityResult.Warnings...)

	// Validate project portability
	portabilityResult := pm.identity.ValidateProjectPortability(project)
	result.Warnings = append(result.Warnings, portabilityResult.Warnings...)

	return result
}

// GetProjectInfo returns comprehensive project information
func (pm *ProjectManager) GetProjectInfo(projectPath string) (*ProjectInfo, error) {
	// Load project
	project, err := pm.LoadProject(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load project: %w", err)
	}

	// Get project structure information
	_, err = pm.structure.GetProjectStructure(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get project structure: %w", err)
	}

	// Create project info
	info := project.GetInfo()

	// Add structure information
	info.Config = &ProjectConfiguration{
		FactsStorage:  "badgerdb",
		VariablesPath: "variables",
		MachinesPath:  "machines",
		ActionsPath:   "actions",
		TemplatesPath: "templates",
	}

	// Add basic stats (TODO: implement actual counting)
	info.Stats = &ProjectStats{
		FactsCount:     0,
		VariablesCount: 0,
		MachinesCount:  0,
		ActionsCount:   0,
		TemplatesCount: 0,
		LastUpdated:    project.UpdatedAt,
	}

	// Add validation status
	validationResult := pm.ValidateProject(projectPath)
	info.Status = &ProjectStatus{
		Valid:         validationResult.Valid,
		Errors:        make([]string, 0),
		Warnings:      make([]string, 0),
		LastValidated: time.Now(),
	}

	// Convert validation errors to strings
	for _, err := range validationResult.Errors {
		info.Status.Errors = append(info.Status.Errors, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}

	for _, warning := range validationResult.Warnings {
		info.Status.Warnings = append(info.Status.Warnings, fmt.Sprintf("%s: %s", warning.Field, warning.Message))
	}

	return info, nil
}

// ListProjects lists all projects in a directory
func (pm *ProjectManager) ListProjects(directory string) ([]*ProjectInfo, error) {
	var projects []*ProjectInfo

	// Read directory entries
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	// Check each directory for projects
	for _, entry := range entries {
		if entry.IsDir() {
			projectPath := filepath.Join(directory, entry.Name())
			projectHCLPath := filepath.Join(projectPath, "project.hcl")

			// Check if this directory contains a project
			if _, err := os.Stat(projectHCLPath); err == nil {
				projectInfo, err := pm.GetProjectInfo(projectPath)
				if err != nil {
					// Skip invalid projects but continue listing
					continue
				}
				projects = append(projects, projectInfo)
			}
		}
	}

	return projects, nil
}

// UpdateProject updates an existing project
func (pm *ProjectManager) UpdateProject(projectPath string, updates *Project) error {
	// Load existing project
	existingProject, err := pm.LoadProject(projectPath)
	if err != nil {
		return fmt.Errorf("failed to load existing project: %w", err)
	}

	// Apply updates
	if updates.Name != "" {
		existingProject.Name = updates.Name
	}
	if updates.Description != "" {
		existingProject.Description = updates.Description
	}
	if updates.Version != "" {
		existingProject.Version = updates.Version
	}
	if updates.Environment != "" {
		existingProject.Environment = updates.Environment
	}
	if updates.Region != "" {
		existingProject.Region = updates.Region
	}
	if updates.Tags != nil {
		existingProject.Tags = updates.Tags
	}
	if updates.Structure != nil {
		existingProject.Structure = updates.Structure
	}
	if updates.Isolation != nil {
		existingProject.Isolation = updates.Isolation
	}
	if updates.Dependencies != nil {
		existingProject.Dependencies = updates.Dependencies
	}
	if updates.Execution != nil {
		existingProject.Execution = updates.Execution
	}
	if updates.Metadata != nil {
		existingProject.Metadata = updates.Metadata
	}

	// Update timestamp
	existingProject.UpdatedAt = time.Now()

	// Validate updated project
	validationResult := pm.validator.ValidateProject(existingProject)
	if !validationResult.Valid {
		return fmt.Errorf("project validation failed after update: %v", validationResult.Errors)
	}

	// Maintain project structure
	if err := pm.structure.MaintainProjectStructure(projectPath, existingProject); err != nil {
		return fmt.Errorf("failed to maintain project structure: %w", err)
	}

	// TODO: Save updated project.hcl file
	// For now, we'll just return success

	return nil
}

// DeleteProject deletes a project (use with caution)
func (pm *ProjectManager) DeleteProject(projectPath string, force bool) error {
	// Validate project path
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return fmt.Errorf("project path does not exist: %s", projectPath)
	}

	// Check if it's actually a project
	projectHCLPath := filepath.Join(projectPath, "project.hcl")
	if _, err := os.Stat(projectHCLPath); os.IsNotExist(err) {
		return fmt.Errorf("path does not contain a valid project: %s", projectPath)
	}

	// If not forced, validate the project first
	if !force {
		validationResult := pm.ValidateProject(projectPath)
		if !validationResult.Valid {
			return fmt.Errorf("cannot delete invalid project: %v", validationResult.Errors)
		}
	}

	// Remove the project directory
	if err := os.RemoveAll(projectPath); err != nil {
		return fmt.Errorf("failed to delete project directory: %w", err)
	}

	return nil
}

// ValidateProjectDependencies validates project dependencies
func (pm *ProjectManager) ValidateProjectDependencies(projectPath string) *DependencyResult {
	// Load project
	project, err := pm.LoadProject(projectPath)
	if err != nil {
		return &DependencyResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:    "project",
					Message:  err.Error(),
					Severity: "error",
				},
			},
		}
	}

	// Create projects map for dependency validation
	projects := map[string]*Project{
		projectPath: project,
	}

	return pm.dependencies.ValidateDependencies(projects)
}

// GetProjectPortability validates project portability
func (pm *ProjectManager) GetProjectPortability(projectPath string) *PortabilityValidation {
	// Load project
	project, err := pm.LoadProject(projectPath)
	if err != nil {
		return &PortabilityValidation{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:    "project",
					Message:  err.Error(),
					Severity: "error",
				},
			},
		}
	}

	return pm.portability.ValidateProjectPortability(project)
}

// MigrateProject migrates a project to a different environment
func (pm *ProjectManager) MigrateProject(projectPath string, targetEnvironment string) (*Project, error) {
	// Load project
	project, err := pm.LoadProject(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load project: %w", err)
	}

	// Migrate project
	migratedProject, err := pm.portability.MigrateProject(project, targetEnvironment)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate project: %w", err)
	}

	return migratedProject, nil
}

// GetProjectConfiguration loads and resolves project configuration
func (pm *ProjectManager) GetProjectConfiguration(projectPath string) (*ConfigurationContext, error) {
	return pm.configEngine.LoadProjectConfiguration(projectPath)
}

// ApplyProjectConfiguration applies CLI overrides to project configuration
func (pm *ProjectManager) ApplyProjectConfiguration(context *ConfigurationContext, cliFlags map[string]interface{}) error {
	return pm.configEngine.ApplyCLIOverrides(context, cliFlags)
}

// GetProjectIsolation applies isolation settings to a project
func (pm *ProjectManager) GetProjectIsolation(projectPath string) (*IsolationContext, error) {
	// Load project
	project, err := pm.LoadProject(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load project: %w", err)
	}

	return pm.isolation.ApplyIsolation(project)
}
