package utilities

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"spooky/internal/schemas"
)

// Validation context constants to avoid repeated string literals
const (
	ValidationContextSchema    = "Schema validation"
	ValidationContextProject   = "Project configuration validation"
	ValidationContextMachines  = "Machines validation"
	ValidationContextActions   = "Actions validation"
	ValidationContextVariables = "Variables validation"
)

// Resource type constants to avoid repeated string literals
const (
	ResourceTypeMachines  = "machines"
	ResourceTypeActions   = "actions"
	ResourceTypeVariables = "variables"
)

// ProjectValidator validates spooky projects with comprehensive error collection
type ProjectValidator struct {
	result *schemas.ValidationResult
}

// NewProjectValidator creates a new project validator
func NewProjectValidator() *ProjectValidator {
	return &ProjectValidator{
		result: &schemas.ValidationResult{
			IsValid:  true,
			Errors:   []schemas.ValidationError{},
			Warnings: []schemas.ValidationWarning{},
		},
	}
}

// ValidateProject performs comprehensive project validation
func (pv *ProjectValidator) ValidateProject(targetDir string) *schemas.ValidationResult {
	// Validate project structure
	pv.validateProjectStructure(targetDir)

	// Validate individual files
	pv.validateProjectFiles(targetDir)

	// Set final validity
	pv.result.IsValid = len(pv.result.Errors) == 0

	return pv.result
}

// validateProjectStructure validates the project directory structure using schema-driven validation
func (pv *ProjectValidator) validateProjectStructure(targetDir string) {
	// Scan the actual project directory structure
	projectDir := pv.scanProjectDirectory(targetDir)

	// Convert to schema format and validate using struct validator
	pv.validateProjectDirectoryWithSchema(projectDir, targetDir)
}

// scanProjectDirectory scans the actual project directory and returns filesystem structure
func (pv *ProjectValidator) scanProjectDirectory(targetDir string) map[string]interface{} {
	projectDir := map[string]interface{}{
		"name":        filepath.Base(targetDir),
		"files":       []string{},
		"directories": []string{},
	}

	// Scan directory for files and subdirectories
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return projectDir
	}

	for _, entry := range entries {
		name := entry.Name()
		isDir := entry.IsDir()

		if isDir {
			projectDir["directories"] = append(projectDir["directories"].([]string), name)
		} else {
			projectDir["files"] = append(projectDir["files"].([]string), name)
		}
	}

	return projectDir
}

// validateProjectDirectoryWithSchema validates project directory structure using schema-driven validation
func (pv *ProjectValidator) validateProjectDirectoryWithSchema(projectDir map[string]interface{}, targetDir string) {
	// Convert scanned directory structure to schema format
	schemaData := map[string]interface{}{
		"project_directory": map[string]interface{}{
			"name":        projectDir["name"],
			"files":       pv.convertFilesToSchema(projectDir["files"].([]string)),
			"directories": pv.convertDirectoriesToSchema(projectDir["directories"].([]string)),
		},
	}

	// Validate using simplified validator
	validator := schemas.NewSimpleValidator()
	result := validator.ValidateData("project_directory", schemaData)

	// Add schema validation errors to our result
	if !result.IsValid {
		for _, schemaErr := range result.Errors {
			schemaErr.File = targetDir
			schemaErr.Context = "Project directory structure validation"
			pv.result.Errors = append(pv.result.Errors, schemaErr)
		}
	}

	// Add schema validation warnings to our result
	for _, schemaWarning := range result.Warnings {
		schemaWarning.File = targetDir
		schemaWarning.Context = "Project directory structure validation"
		pv.result.Warnings = append(pv.result.Warnings, schemaWarning)
	}

	// Additional validation for alternative configurations (either file or directory)
	pv.validateAlternativeConfigurations(targetDir)
}

// convertToSchema is a generic function that converts items to schema format
func (pv *ProjectValidator) convertToSchema(items []string, expectedItems map[string]map[string]interface{}) []map[string]interface{} {
	var schemaItems []map[string]interface{}

	// Add found items to schema
	for _, itemName := range items {
		if expectedItem, exists := expectedItems[itemName]; exists {
			schemaItems = append(schemaItems, expectedItem)
		}
	}

	return schemaItems
}

// ExpectedItem represents a configuration item (file or directory)
type ExpectedItem struct {
	Name        string
	Type        string
	Required    bool
	Description string
	Pattern     string
}

// getExpectedFiles returns the expected files configuration
func (pv *ProjectValidator) getExpectedFiles() map[string]map[string]interface{} {
	expectedFiles := []ExpectedItem{
		{
			Name:        "project.hcl",
			Type:        "file",
			Required:    true,
			Description: "Main project configuration file",
			Pattern:     "project \"[a-zA-Z0-9_-]+\" {",
		},
		{
			Name:        "machines.hcl",
			Type:        "file",
			Required:    false,
			Description: "Machine inventory definitions",
			Pattern:     "machines {",
		},
		{
			Name:        "actions.hcl",
			Type:        "file",
			Required:    false,
			Description: "Main actions file",
			Pattern:     "actions {",
		},
		{
			Name:        "variables.hcl",
			Type:        "file",
			Required:    false,
			Description: "Main variables file",
			Pattern:     "variables {",
		},
		{
			Name:        "README.md",
			Type:        "file",
			Required:    false,
			Description: "Project documentation",
			Pattern:     "# .*",
		},
	}

	return pv.convertExpectedItemsToMap(expectedFiles)
}

// getExpectedDirectories returns the expected directories configuration
func (pv *ProjectValidator) getExpectedDirectories() map[string]map[string]interface{} {
	expectedDirs := []ExpectedItem{
		{
			Name:        "machines",
			Type:        "directory",
			Required:    false,
			Description: "Machine inventory files directory",
			Pattern:     ".*\\.hcl$",
		},
		{
			Name:        "actions",
			Type:        "directory",
			Required:    false,
			Description: "Organized action files",
			Pattern:     ".*\\.hcl$",
		},
		{
			Name:        "variables",
			Type:        "directory",
			Required:    false,
			Description: "Variables files directory",
			Pattern:     ".*\\.hcl$",
		},
		{
			Name:        "templates",
			Type:        "directory",
			Required:    false,
			Description: "Template files for dynamic content",
			Pattern:     "",
		},
		{
			Name:        "files",
			Type:        "directory",
			Required:    false,
			Description: "Static files to be deployed",
			Pattern:     "",
		},
	}

	return pv.convertExpectedItemsToMap(expectedDirs)
}

// convertExpectedItemsToMap converts a slice of ExpectedItem to the expected map format
func (pv *ProjectValidator) convertExpectedItemsToMap(items []ExpectedItem) map[string]map[string]interface{} {
	result := make(map[string]map[string]interface{})
	for _, item := range items {
		result[item.Name] = map[string]interface{}{
			"name":        item.Name,
			"type":        item.Type,
			"required":    item.Required,
			"description": item.Description,
			"pattern":     item.Pattern,
		}
	}
	return result
}

// convertFilesToSchema converts file names to schema format
func (pv *ProjectValidator) convertFilesToSchema(files []string) []map[string]interface{} {
	return pv.convertToSchema(files, pv.getExpectedFiles())
}

// convertDirectoriesToSchema converts directory names to schema format
func (pv *ProjectValidator) convertDirectoriesToSchema(directories []string) []map[string]interface{} {
	return pv.convertToSchema(directories, pv.getExpectedDirectories())
}

// validateAlternativeConfigurations validates that at least one configuration option exists for each type
func (pv *ProjectValidator) validateAlternativeConfigurations(targetDir string) {
	// Define resource types to validate
	resourceTypes := []struct {
		name         string
		filePath     string
		dirPath      string
		resourceType string
	}{
		{
			name:         "machines",
			filePath:     "machines.hcl",
			dirPath:      ResourceTypeMachines,
			resourceType: "machines",
		},
		{
			name:         "actions",
			filePath:     "actions.hcl",
			dirPath:      ResourceTypeActions,
			resourceType: "actions",
		},
		{
			name:         "variables",
			filePath:     "variables.hcl",
			dirPath:      ResourceTypeVariables,
			resourceType: "variables",
		},
	}

	// Validate each resource type
	for _, resource := range resourceTypes {
		pv.validateResourceType(targetDir, resource.name, resource.filePath, resource.dirPath, resource.resourceType)
	}
}

// validateResourceType validates that a specific resource type has either a file or directory configuration
func (pv *ProjectValidator) validateResourceType(targetDir, name, filePath, dirPath, resourceType string) {
	hclPath := filepath.Join(targetDir, filePath)
	dirPathFull := filepath.Join(targetDir, dirPath)

	if _, err := os.Stat(hclPath); os.IsNotExist(err) {
		if _, err := os.Stat(dirPathFull); os.IsNotExist(err) {
			pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
				Message: fmt.Sprintf("either %s file or %s/ directory must exist", filePath, resourceType),
				File:    targetDir,
				Context: "Project structure validation",
			})
		}
	}
}

// validateProjectFiles validates individual project files
func (pv *ProjectValidator) validateProjectFiles(targetDir string) {
	// Validate project.hcl
	pv.validateProjectConfig(targetDir)

	// Validate machines
	pv.validateMachines(targetDir)

	// Validate actions
	pv.validateActions(targetDir)

	// Validate variables
	pv.validateVariablesConfig(targetDir)
}

// validateProjectConfig validates the project.hcl file
func (pv *ProjectValidator) validateProjectConfig(targetDir string) {
	projectHCLPath := filepath.Join(targetDir, "project.hcl")

	// Skip validation if file doesn't exist (structure validation already caught this)
	if _, err := os.Stat(projectHCLPath); os.IsNotExist(err) {
		return
	}

	// Validate syntax and schema using simplified validator
	validator := schemas.NewSimpleValidator()
	content, err := os.ReadFile(projectHCLPath)
	if err != nil {
		pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
			Message: fmt.Sprintf("failed to read project.hcl: %v", err),
			File:    projectHCLPath,
			Context: ValidationContextProject,
		})
		return
	}

	// Validate against project schema
	result, err := validator.ValidateHCLContent("project", string(content))
	if err != nil {
		pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
			Message: fmt.Sprintf("project.hcl schema validation failed: %v", err),
			File:    projectHCLPath,
			Context: ValidationContextProject,
		})
		return
	}

	if !result.IsValid {
		// Add schema validation errors with file context
		for _, schemaErr := range result.Errors {
			schemaErr.File = projectHCLPath
			schemaErr.Context = ValidationContextSchema
			pv.result.Errors = append(pv.result.Errors, schemaErr)
		}

		for _, schemaWarning := range result.Warnings {
			schemaWarning.File = projectHCLPath
			schemaWarning.Context = ValidationContextSchema
			pv.result.Warnings = append(pv.result.Warnings, schemaWarning)
		}
	}
}

// validateMachines validates machine configurations
func (pv *ProjectValidator) validateMachines(projectPath string) {
	machinesHCLPath := filepath.Join(projectPath, "machines.hcl")
	machinesDirPath := filepath.Join(projectPath, ResourceTypeMachines)

	// Always validate machines.hcl if it exists
	if _, err := os.Stat(machinesHCLPath); err == nil {
		pv.validateMachinesFile(machinesHCLPath)
	}

	// Always validate machines/ directory if it exists
	if _, err := os.Stat(machinesDirPath); err == nil {
		pv.validateMachinesDirectory(machinesDirPath)
	}
}

// validateHCLFile is a generic function to validate HCL files with a given schema type
func (pv *ProjectValidator) validateHCLFile(filePath, schemaType, context string) {
	validator := schemas.NewSimpleValidator()
	content, err := os.ReadFile(filePath)
	if err != nil {
		pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
			Message: fmt.Sprintf("failed to read %s: %v", filepath.Base(filePath), err),
			File:    filePath,
			Context: context,
		})
		return
	}

	result, err := validator.ValidateHCLContent(schemaType, string(content))
	if err != nil {
		pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
			Message: fmt.Sprintf("schema validation failed: %v", err),
			File:    filePath,
			Context: context,
		})
		return
	}

	if !result.IsValid {
		for _, schemaErr := range result.Errors {
			schemaErr.File = filePath
			schemaErr.Context = context
			pv.result.Errors = append(pv.result.Errors, schemaErr)
		}
	}
}

// validateResourceFile validates a single resource file with the specified schema type and context
func (pv *ProjectValidator) validateResourceFile(filePath, resourceType, context string) {
	pv.validateHCLFile(filePath, resourceType, context)
}

// validateMachinesFile validates a single machines.hcl file
func (pv *ProjectValidator) validateMachinesFile(filePath string) {
	pv.validateResourceFile(filePath, ResourceTypeMachines, ValidationContextMachines)
}

// validateResourceDirectory validates all .hcl files in a resource directory
func (pv *ProjectValidator) validateResourceDirectory(dirPath, resourceType, context string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
			Message: fmt.Sprintf("failed to read %s directory: %v", resourceType, err),
			File:    dirPath,
			Context: context,
		})
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".hcl") {
			filePath := filepath.Join(dirPath, entry.Name())
			pv.validateResourceFile(filePath, resourceType, context)
		}
	}
}

// validateMachinesDirectory validates all .hcl files in the machines directory
func (pv *ProjectValidator) validateMachinesDirectory(dirPath string) {
	pv.validateResourceDirectory(dirPath, ResourceTypeMachines, ValidationContextMachines)
}

// validateActions validates action configurations
func (pv *ProjectValidator) validateActions(projectPath string) {
	actionsHCLPath := filepath.Join(projectPath, "actions.hcl")
	actionsDirPath := filepath.Join(projectPath, ResourceTypeActions)

	// Always validate actions.hcl if it exists
	if _, err := os.Stat(actionsHCLPath); err == nil {
		pv.validateActionsFile(actionsHCLPath)
	}

	// Always validate actions/ directory if it exists
	if _, err := os.Stat(actionsDirPath); err == nil {
		pv.validateActionsDirectory(actionsDirPath)
	}
}

// validateActionsFile validates a single actions.hcl file
func (pv *ProjectValidator) validateActionsFile(filePath string) {
	pv.validateResourceFile(filePath, ResourceTypeActions, ValidationContextActions)
}

// validateActionsDirectory validates all .hcl files in the actions directory
func (pv *ProjectValidator) validateActionsDirectory(dirPath string) {
	pv.validateResourceDirectory(dirPath, ResourceTypeActions, ValidationContextActions)
}

// validateVariablesConfig validates variable configurations
func (pv *ProjectValidator) validateVariablesConfig(projectPath string) {
	variablesHCLPath := filepath.Join(projectPath, "variables.hcl")
	variablesDirPath := filepath.Join(projectPath, ResourceTypeVariables)

	// Always validate variables.hcl if it exists
	if _, err := os.Stat(variablesHCLPath); err == nil {
		pv.validateVariablesFile(variablesHCLPath)
	}

	// Always validate variables/ directory if it exists
	if _, err := os.Stat(variablesDirPath); err == nil {
		pv.validateVariablesDirectory(variablesDirPath)
	}
}

// validateVariablesFile validates a single variables.hcl file
func (pv *ProjectValidator) validateVariablesFile(filePath string) {
	pv.validateResourceFile(filePath, ResourceTypeVariables, ValidationContextVariables)
}

// validateVariablesDirectory validates all .hcl files in the variables directory
func (pv *ProjectValidator) validateVariablesDirectory(dirPath string) {
	pv.validateResourceDirectory(dirPath, ResourceTypeVariables, ValidationContextVariables)
}
