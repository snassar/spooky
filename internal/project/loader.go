// Package project provides project management functionality for spooky.
package project

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"

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

// LoadProjectConfig loads project configuration from project.hcl
func (l *Loader) LoadProjectConfig(ctx context.Context, projectPath string) (*spookytypes.ProjectConfig, error) {
	l.logger.Info("Loading project configuration", map[string]interface{}{
		"project_path": projectPath,
	})

	// Look for project.hcl file
	projectHCLPath := filepath.Join(projectPath, "project.hcl")

	// Check if project.hcl exists
	if _, err := os.Stat(projectHCLPath); os.IsNotExist(err) {
		// Fallback to basic config if no project.hcl exists
		projectName := filepath.Base(projectPath)
		return l.createDefaultProjectConfig(projectName), nil
	}

	// Parse HCL file
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(projectHCLPath)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse project.hcl: %v", diags)
	}

	// Extract project block
	content, diags := file.Body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type: "project",
			},
		},
	})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse project block: %v", diags)
	}

	if len(content.Blocks) == 0 {
		// No project block found, use default config
		projectName := filepath.Base(projectPath)
		return l.createDefaultProjectConfig(projectName), nil
	}

	// Parse the first project block
	projectBlock := content.Blocks[0]
	return l.parseProjectBlock(projectBlock, projectPath)
}

// parseProjectBlock parses a project block from HCL
func (l *Loader) parseProjectBlock(block *hcl.Block, projectPath string) (*spookytypes.ProjectConfig, error) {
	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "name", Required: true},
			{Name: "description", Required: false},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "metadata"},
			{Type: "settings"},
		},
	})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse project block attributes: %v", diags)
	}

	// Parse name attribute
	nameAttr := content.Attributes["name"]
	name, diags := nameAttr.Expr.Value(nil)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse project name: %v", diags)
	}

	// Parse description attribute (optional)
	var description string
	if descAttr, exists := content.Attributes["description"]; exists {
		desc, diags := descAttr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse project description: %v", diags)
		}
		description = desc.AsString()
	}

	// Parse metadata block
	var metadata *spookytypesproject.ProjectMetadata
	for _, block := range content.Blocks {
		if block.Type == "metadata" {
			var err error
			metadata, err = l.parseMetadataBlock(block)
			if err != nil {
				return nil, fmt.Errorf("failed to parse metadata block: %v", err)
			}
			break
		}
	}

	// Parse settings block
	var settings *spookytypesproject.ProjectSettings
	for _, block := range content.Blocks {
		if block.Type == "settings" {
			var err error
			settings, err = l.parseSettingsBlock(block)
			if err != nil {
				return nil, fmt.Errorf("failed to parse settings block: %v", err)
			}
			break
		}
	}

	// Use defaults if not specified
	if metadata == nil {
		metadata = l.createDefaultMetadata()
	}
	if settings == nil {
		settings = l.createDefaultSettings()
	}

	return &spookytypes.ProjectConfig{
		Name:        name.AsString(),
		Description: description,
		Metadata:    metadata,
		Settings:    settings,
	}, nil
}

// parseMetadataBlock parses a metadata block from HCL
func (l *Loader) parseMetadataBlock(block *hcl.Block) (*spookytypesproject.ProjectMetadata, error) {
	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "version", Required: false},
			{Name: "author", Required: false},
			{Name: "tags", Required: false},
		},
	})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse metadata block: %v", diags)
	}

	metadata := &spookytypesproject.ProjectMetadata{
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	// Parse version
	if versionAttr, exists := content.Attributes["version"]; exists {
		version, diags := versionAttr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse version: %v", diags)
		}
		metadata.Version = version.AsString()
	}

	// Parse author
	if authorAttr, exists := content.Attributes["author"]; exists {
		author, diags := authorAttr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse author: %v", diags)
		}
		metadata.Author = author.AsString()
	}

	// Parse tags
	if tagsAttr, exists := content.Attributes["tags"]; exists {
		tags, diags := tagsAttr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse tags: %v", diags)
		}
		if tags.Type().IsListType() {
			for _, tag := range tags.AsValueSlice() {
				metadata.Tags = append(metadata.Tags, tag.AsString())
			}
		}
	}

	return metadata, nil
}

// parseSettingsBlock parses a settings block from HCL
func (l *Loader) parseSettingsBlock(block *hcl.Block) (*spookytypesproject.ProjectSettings, error) {
	content, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "parallel_workers", Required: false},
			{Name: "timeout_seconds", Required: false},
			{Name: "log_level", Required: false},
			{Name: "default_dry_run", Required: false},
			{Name: "validate_before_run", Required: false},
			{Name: "max_retries", Required: false},
			{Name: "retry_delay_seconds", Required: false},
		},
	})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse settings block: %v", diags)
	}

	settings := l.createDefaultSettings()

	// Parse parallel_workers
	if workersAttr, exists := content.Attributes["parallel_workers"]; exists {
		workers, diags := workersAttr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse parallel_workers: %v", diags)
		}
		workersInt, _ := workers.AsBigFloat().Int64()
		settings.ParallelWorkers = int(workersInt)
	}

	// Parse timeout_seconds
	if timeoutAttr, exists := content.Attributes["timeout_seconds"]; exists {
		timeout, diags := timeoutAttr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse timeout_seconds: %v", diags)
		}
		timeoutInt, _ := timeout.AsBigFloat().Int64()
		settings.TimeoutSeconds = int(timeoutInt)
	}

	// Parse log_level
	if logLevelAttr, exists := content.Attributes["log_level"]; exists {
		logLevel, diags := logLevelAttr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse log_level: %v", diags)
		}
		settings.LogLevel = logLevel.AsString()
	}

	// Parse default_dry_run
	if dryRunAttr, exists := content.Attributes["default_dry_run"]; exists {
		dryRun, diags := dryRunAttr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse default_dry_run: %v", diags)
		}
		settings.DefaultDryRun = dryRun.True()
	}

	// Parse validate_before_run
	if validateAttr, exists := content.Attributes["validate_before_run"]; exists {
		validate, diags := validateAttr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse validate_before_run: %v", diags)
		}
		settings.ValidateBeforeRun = validate.True()
	}

	// Parse max_retries
	if retriesAttr, exists := content.Attributes["max_retries"]; exists {
		retries, diags := retriesAttr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse max_retries: %v", diags)
		}
		retriesInt, _ := retries.AsBigFloat().Int64()
		settings.MaxRetries = int(retriesInt)
	}

	// Parse retry_delay_seconds
	if delayAttr, exists := content.Attributes["retry_delay_seconds"]; exists {
		delay, diags := delayAttr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to parse retry_delay_seconds: %v", diags)
		}
		delayInt, _ := delay.AsBigFloat().Int64()
		settings.RetryDelaySeconds = int(delayInt)
	}

	return settings, nil
}

// createDefaultProjectConfig creates a default project configuration
func (l *Loader) createDefaultProjectConfig(projectName string) *spookytypes.ProjectConfig {
	return &spookytypes.ProjectConfig{
		Name:        projectName,
		Description: "A spooky project for automation and orchestration",
		Metadata:    l.createDefaultMetadata(),
		Settings:    l.createDefaultSettings(),
	}
}

// createDefaultMetadata creates default project metadata
func (l *Loader) createDefaultMetadata() *spookytypesproject.ProjectMetadata {
	return &spookytypesproject.ProjectMetadata{
		Version:    "1.0.0",
		Author:     "spooky-user",
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}
}

// createDefaultSettings creates default project settings
func (l *Loader) createDefaultSettings() *spookytypesproject.ProjectSettings {
	return &spookytypesproject.ProjectSettings{
		ParallelWorkers:   10,
		TimeoutSeconds:    300,
		LogLevel:          "info",
		DefaultDryRun:     false,
		ValidateBeforeRun: true,
		MaxRetries:        3,
		RetryDelaySeconds: 5,
	}
}

// LoadProjectMetadata loads project metadata from project.hcl and other sources
func (l *Loader) LoadProjectMetadata(ctx context.Context, projectPath string) (*spookytypes.ProjectMetadata, error) {
	l.logger.Info("Loading project metadata", map[string]interface{}{
		"project_path": projectPath,
	})

	// Try to load metadata from project.hcl first
	projectHCLPath := filepath.Join(projectPath, "project.hcl")
	if _, err := os.Stat(projectHCLPath); err == nil {
		// Parse project.hcl to extract metadata
		parser := hclparse.NewParser()
		file, diags := parser.ParseHCLFile(projectHCLPath)
		if !diags.HasErrors() {
			content, diags := file.Body.Content(&hcl.BodySchema{
				Blocks: []hcl.BlockHeaderSchema{
					{Type: "project"},
				},
			})
			if !diags.HasErrors() && len(content.Blocks) > 0 {
				projectBlock := content.Blocks[0]
				projectContent, diags := projectBlock.Body.Content(&hcl.BodySchema{
					Blocks: []hcl.BlockHeaderSchema{
						{Type: "metadata"},
					},
				})
				if !diags.HasErrors() && len(projectContent.Blocks) > 0 {
					metadata, err := l.parseMetadataBlock(projectContent.Blocks[0])
					if err == nil {
						l.logger.Info("Project metadata loaded from project.hcl", map[string]interface{}{
							"project_path": projectPath,
							"version":      metadata.Version,
							"author":       metadata.Author,
						})
						return metadata, nil
					}
				}
			}
		}
	}

	// Fallback to default metadata
	metadata := l.createDefaultMetadata()
	l.logger.Info("Project metadata loaded (default)", map[string]interface{}{
		"project_path": projectPath,
		"version":      metadata.Version,
		"author":       metadata.Author,
	})

	return metadata, nil
}
