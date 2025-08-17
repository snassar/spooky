// Package schemas provides schema parsing and validation functionality
package schemas

import (
	spookytypeslogging "spooky/internal/types/logging"
)

// ProjectDirectorySchema represents the parsed project directory schema
type ProjectDirectorySchema struct {
	RequiredFiles   []SchemaFile
	OptionalFiles   []SchemaFile
	RequiredDirs    []SchemaDirectory
	OptionalDirs    []SchemaDirectory
	ValidationRules []string
}

// SchemaFile represents a file definition in the schema
type SchemaFile struct {
	Name        string
	Type        string
	Required    bool
	Description string
	Validate    string
	Pattern     string
}

// SchemaDirectory represents a directory definition in the schema
type SchemaDirectory struct {
	Name        string
	Type        string
	Required    bool
	Description string
	Validate    string
	Pattern     string
}

// ProjectDirectorySchemaParser parses project directory schemas
type ProjectDirectorySchemaParser struct {
	logger spookytypeslogging.Logger
}

// NewProjectDirectorySchemaParser creates a new schema parser
func NewProjectDirectorySchemaParser(logger spookytypeslogging.Logger) *ProjectDirectorySchemaParser {
	return &ProjectDirectorySchemaParser{
		logger: logger,
	}
}

// ParseProjectDirectorySchema parses the project directory schema from the embedded schema file
func (p *ProjectDirectorySchemaParser) ParseProjectDirectorySchema() (*ProjectDirectorySchema, error) {
	p.logger.Debug("Parsing project directory schema")

	// For now, we'll hardcode the schema structure based on the HCL file
	// In a full implementation, this would parse the actual HCL schema file
	schema := &ProjectDirectorySchema{
		RequiredFiles: []SchemaFile{
			{
				Name:        "project.hcl",
				Type:        "file",
				Required:    true,
				Description: "Main project configuration file",
				Validate:    "hcl_project_config",
				Pattern:     "project \"[a-zA-Z0-9_-]+\" {",
			},
		},
		OptionalFiles: []SchemaFile{
			{
				Name:        "machines.hcl",
				Type:        "file",
				Required:    false,
				Description: "Machine inventory definitions",
				Validate:    "hcl_machines_config",
				Pattern:     "machines {",
			},
			{
				Name:        "actions.hcl",
				Type:        "file",
				Required:    false,
				Description: "Main actions file",
				Validate:    "hcl_actions_config",
				Pattern:     "actions {",
			},
			{
				Name:        "variables.hcl",
				Type:        "file",
				Required:    false,
				Description: "Main variables file",
				Validate:    "hcl_variables_config",
				Pattern:     "variables {",
			},
			{
				Name:        "README.md",
				Type:        "file",
				Required:    false,
				Description: "Project documentation",
				Pattern:     "# .*",
			},
			{
				Name:        "recipients.txt",
				Type:        "file",
				Required:    false,
				Description: "Project-specific age recipients",
				Pattern:     "age1[a-zA-Z0-9]+",
			},
		},
		RequiredDirs: []SchemaDirectory{
			// No required directories according to schema
		},
		OptionalDirs: []SchemaDirectory{
			{
				Name:        "machines",
				Type:        "directory",
				Required:    false,
				Description: "Machine inventory files directory",
				Validate:    "hcl_machines_files",
				Pattern:     ".*\\.hcl$",
			},
			{
				Name:        "actions",
				Type:        "directory",
				Required:    false,
				Description: "Organized action files",
				Validate:    "hcl_actions_files",
				Pattern:     ".*\\.hcl$",
			},
			{
				Name:        "variables",
				Type:        "directory",
				Required:    false,
				Description: "Variables files directory",
				Validate:    "hcl_variables_files",
				Pattern:     ".*\\.hcl$",
			},
			{
				Name:        "templates",
				Type:        "directory",
				Required:    false,
				Description: "Template files for dynamic content",
				Validate:    "directory_exists",
			},
			{
				Name:        "files",
				Type:        "directory",
				Required:    false,
				Description: "Static files to be deployed",
				Validate:    "directory_exists",
			},
			{
				Name:        "logs",
				Type:        "directory",
				Required:    false,
				Description: "Log files directory",
				Validate:    "directory_exists",
			},
		},
		ValidationRules: []string{
			"machines_file_or_directory_exists",
			"actions_file_or_directory_exists",
			"variables_file_or_directory_exists",
			"no_circular_references",
			"logging_file_output_requires_logs_directory",
			"logging_file_path_validation",
		},
	}

	p.logger.Debug("Project directory schema parsed successfully", map[string]interface{}{
		"required_files": len(schema.RequiredFiles),
		"optional_files": len(schema.OptionalFiles),
		"required_dirs":  len(schema.RequiredDirs),
		"optional_dirs":  len(schema.OptionalDirs),
	})

	return schema, nil
}

// GetRequiredDirectories returns all required directories from the schema
func (s *ProjectDirectorySchema) GetRequiredDirectories() []SchemaDirectory {
	return s.RequiredDirs
}

// GetOptionalDirectories returns all optional directories from the schema
func (s *ProjectDirectorySchema) GetOptionalDirectories() []SchemaDirectory {
	return s.OptionalDirs
}

// GetRequiredFiles returns all required files from the schema
func (s *ProjectDirectorySchema) GetRequiredFiles() []SchemaFile {
	return s.RequiredFiles
}

// GetOptionalFiles returns all optional files from the schema
func (s *ProjectDirectorySchema) GetOptionalFiles() []SchemaFile {
	return s.OptionalFiles
}

// ShouldCreateOptionalDirectory determines if an optional directory should be created
// This can be enhanced with user preferences or configuration
func (s *ProjectDirectorySchema) ShouldCreateOptionalDirectory(dir SchemaDirectory) bool {
	// For now, create commonly useful directories by default
	commonlyUseful := []string{"files", "logs"}
	for _, useful := range commonlyUseful {
		if dir.Name == useful {
			return true
		}
	}
	return false
}

// ShouldCreateOptionalFile determines if an optional file should be created
// This can be enhanced with user preferences or configuration
func (s *ProjectDirectorySchema) ShouldCreateOptionalFile(file SchemaFile) bool {
	// For now, create README.md by default
	if file.Name == "README.md" {
		return true
	}
	return false
}
