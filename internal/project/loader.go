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
func (l *Loader) LoadProjectConfig(_ context.Context, projectPath string) (*spookytypes.ProjectConfig, error) {
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
func (l *Loader) parseProjectBlock(block *hcl.Block, _ string) (*spookytypes.ProjectConfig, error) {
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
	var metadata *spookytypesproject.Metadata
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
	var settings *spookytypesproject.Settings
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
func (l *Loader) parseMetadataBlock(block *hcl.Block) (*spookytypesproject.Metadata, error) {
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

	metadata := &spookytypesproject.Metadata{
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

// getSettingsBlockSchema returns the HCL schema for settings blocks
func getSettingsBlockSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "parallel_workers", Required: false},
			{Name: "timeout_seconds", Required: false},
			{Name: "log_level", Required: false},
			{Name: "default_dry_run", Required: false},
			{Name: "validate_before_run", Required: false},
			{Name: "max_retries", Required: false},
			{Name: "retry_delay_seconds", Required: false},
		},
	}
}

// parseBlockContent parses the content of a settings block
func parseBlockContent(block *hcl.Block) (*hcl.BodyContent, error) {
	schema := getSettingsBlockSchema()
	content, diags := block.Body.Content(schema)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse settings block: %v", diags)
	}
	return content, nil
}

// parseIntegerAttribute parses an integer attribute from HCL content
func parseIntegerAttribute(content *hcl.BodyContent, attrName string) (int, error) {
	attr, exists := content.Attributes[attrName]
	if !exists {
		return 0, nil // Return default value if attribute doesn't exist
	}

	value, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return 0, fmt.Errorf("failed to parse %s: %v", attrName, diags)
	}

	intValue, _ := value.AsBigFloat().Int64()
	return int(intValue), nil
}

// parseStringAttribute parses a string attribute from HCL content
func parseStringAttribute(content *hcl.BodyContent, attrName string) (string, error) {
	attr, exists := content.Attributes[attrName]
	if !exists {
		return "", nil // Return default value if attribute doesn't exist
	}

	value, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return "", fmt.Errorf("failed to parse %s: %v", attrName, diags)
	}

	return value.AsString(), nil
}

// parseBooleanAttribute parses a boolean attribute from HCL content
func parseBooleanAttribute(content *hcl.BodyContent, attrName string) (bool, error) {
	attr, exists := content.Attributes[attrName]
	if !exists {
		return false, nil // Return default value if attribute doesn't exist
	}

	value, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return false, fmt.Errorf("failed to parse %s: %v", attrName, diags)
	}

	return value.True(), nil
}

// populateSettings populates settings with parsed values from HCL content
func populateSettings(settings *spookytypesproject.Settings, content *hcl.BodyContent) error {
	// Parse integer attributes
	if parallelWorkers, err := parseIntegerAttribute(content, "parallel_workers"); err != nil {
		return err
	} else {
		settings.ParallelWorkers = parallelWorkers
	}

	if timeoutSeconds, err := parseIntegerAttribute(content, "timeout_seconds"); err != nil {
		return err
	} else {
		settings.TimeoutSeconds = timeoutSeconds
	}

	if maxRetries, err := parseIntegerAttribute(content, "max_retries"); err != nil {
		return err
	} else {
		settings.MaxRetries = maxRetries
	}

	if retryDelaySeconds, err := parseIntegerAttribute(content, "retry_delay_seconds"); err != nil {
		return err
	} else {
		settings.RetryDelaySeconds = retryDelaySeconds
	}

	// Parse string attributes
	if logLevel, err := parseStringAttribute(content, "log_level"); err != nil {
		return err
	} else {
		settings.LogLevel = logLevel
	}

	// Parse boolean attributes
	if defaultDryRun, err := parseBooleanAttribute(content, "default_dry_run"); err != nil {
		return err
	} else {
		settings.DefaultDryRun = defaultDryRun
	}

	if validateBeforeRun, err := parseBooleanAttribute(content, "validate_before_run"); err != nil {
		return err
	} else {
		settings.ValidateBeforeRun = validateBeforeRun
	}

	return nil
}

// parseSettingsBlock parses a settings block from HCL
func (l *Loader) parseSettingsBlock(block *hcl.Block) (*spookytypesproject.Settings, error) {
	// Parse block content
	content, err := parseBlockContent(block)
	if err != nil {
		return nil, err
	}

	// Create default settings
	settings := l.createDefaultSettings()

	// Populate settings with parsed values
	if err := populateSettings(settings, content); err != nil {
		return nil, err
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
func (l *Loader) createDefaultMetadata() *spookytypesproject.Metadata {
	return &spookytypesproject.Metadata{
		Version:    "1.0.0",
		Author:     "spooky-user",
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}
}

// createDefaultSettings creates default project settings
func (l *Loader) createDefaultSettings() *spookytypesproject.Settings {
	return &spookytypesproject.Settings{
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
func (l *Loader) LoadProjectMetadata(_ context.Context, projectPath string) (*spookytypes.ProjectMetadata, error) {
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
