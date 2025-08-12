// Package project provides project management functionality for spooky.
package project

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesproject "spooky/internal/types/project"
)

// Loader implements the ProjectLoader interface
type Loader struct {
	logger spookytypeslogging.Logger
}

// NewLoader creates a new ProjectLoader instance
func NewLoader(logger spookytypeslogging.Logger) spookyinterfaces.ProjectLoader {
	return &Loader{
		logger: logger,
	}
}

// LoadProject loads a project from disk
func (l *Loader) LoadProject(ctx context.Context, projectPath string) (*spookytypes.Project, error) {
	l.logger.Info("Loading project from disk", map[string]interface{}{
		"project_path": projectPath,
	})

	// Resolve absolute path
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve project path: %w", err)
	}

	// Check if project directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("project directory does not exist: %s", absPath)
	}

	// Load project configuration
	config, err := l.LoadProjectConfig(ctx, projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load project config: %w", err)
	}

	// Load project metadata
	metadata, err := l.LoadProjectMetadata(ctx, projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load project metadata: %w", err)
	}

	// Create project instance
	project := &spookytypes.Project{
		Path:      absPath,
		Config:    config,
		Metadata:  metadata,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	l.logger.Info("Project loaded successfully", map[string]interface{}{
		"project_path": absPath,
		"project_name": config.Name,
	})

	return project, nil
}

// LoadProjectConfig loads project configuration
func (l *Loader) LoadProjectConfig(ctx context.Context, projectPath string) (*spookytypes.ProjectConfig, error) {
	l.logger.Info("Loading project configuration", map[string]interface{}{
		"project_path": projectPath,
	})

	// For now, create a basic config since we don't have HCL parsing implemented yet
	// In a full implementation, this would parse the project.hcl file
	projectName := filepath.Base(projectPath)

	config := &spookytypes.ProjectConfig{
		Name:        projectName,
		Description: "A spooky project for automation and orchestration",
		Metadata: &spookytypesproject.ProjectMetadata{
			Version: "1.0.0",
			Author:  "spooky-user",
		},
		Settings: &spookytypesproject.ProjectSettings{
			ParallelWorkers:   10,
			TimeoutSeconds:    300,
			LogLevel:          "info",
			DefaultDryRun:     false,
			ValidateBeforeRun: true,
			MaxRetries:        3,
			RetryDelaySeconds: 5,
		},
	}

	l.logger.Info("Project configuration loaded", map[string]interface{}{
		"project_name": config.Name,
	})

	return config, nil
}

// LoadProjectMetadata loads project metadata
func (l *Loader) LoadProjectMetadata(ctx context.Context, projectPath string) (*spookytypes.ProjectMetadata, error) {
	l.logger.Info("Loading project metadata", map[string]interface{}{
		"project_path": projectPath,
	})

	// For now, create basic metadata since we don't have full parsing implemented yet
	// In a full implementation, this would extract metadata from project.hcl and other sources
	metadata := &spookytypesproject.ProjectMetadata{
		Version:    "1.0.0",
		Author:     "spooky-user",
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	l.logger.Info("Project metadata loaded", map[string]interface{}{
		"project_path": projectPath,
		"version":      metadata.Version,
		"author":       metadata.Author,
	})

	return metadata, nil
}
