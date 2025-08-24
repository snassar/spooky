package utilities

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"spooky/internal/schemas"
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

// validateProjectStructure validates the project directory structure
func (pv *ProjectValidator) validateProjectStructure(targetDir string) {
	// Check for required project.hcl
	projectHCLPath := filepath.Join(targetDir, "project.hcl")
	if _, err := os.Stat(projectHCLPath); os.IsNotExist(err) {
		pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
			Message: "required file missing: project.hcl",
			File:    targetDir,
			Context: "Project structure validation",
		})
		return // Can't continue without project.hcl
	}

	// Check for machines configuration (file or directory)
	machinesHCLPath := filepath.Join(targetDir, "machines.hcl")
	machinesDirPath := filepath.Join(targetDir, "machines")

	if _, err := os.Stat(machinesHCLPath); os.IsNotExist(err) {
		if _, err := os.Stat(machinesDirPath); os.IsNotExist(err) {
			pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
				Message: "either machines.hcl file or machines/ directory must exist",
				File:    targetDir,
				Context: "Project structure validation",
			})
		}
	}

	// Check for actions configuration (file or directory)
	actionsHCLPath := filepath.Join(targetDir, "actions.hcl")
	actionsDirPath := filepath.Join(targetDir, "actions")

	if _, err := os.Stat(actionsHCLPath); os.IsNotExist(err) {
		if _, err := os.Stat(actionsDirPath); os.IsNotExist(err) {
			pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
				Message: "either actions.hcl file or actions/ directory must exist",
				File:    targetDir,
				Context: "Project structure validation",
			})
		}
	}

	// Check for variables configuration (file or directory)
	variablesHCLPath := filepath.Join(targetDir, "variables.hcl")
	variablesDirPath := filepath.Join(targetDir, "variables")

	if _, err := os.Stat(variablesHCLPath); os.IsNotExist(err) {
		if _, err := os.Stat(variablesDirPath); os.IsNotExist(err) {
			pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
				Message: "either variables.hcl file or variables/ directory must exist",
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

	if _, err := os.Stat(projectHCLPath); os.IsNotExist(err) {
		pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
			Message: "project.hcl not found",
			File:    projectHCLPath,
			Context: "Project configuration validation",
		})
		return
	}

	// Validate syntax and schema using struct validator
	validator := schemas.NewStructValidator()
	content, err := os.ReadFile(projectHCLPath)
	if err != nil {
		pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
			Message: fmt.Sprintf("failed to read project.hcl: %v", err),
			File:    projectHCLPath,
			Context: "Project configuration validation",
		})
		return
	}

	// Validate against project schema
	result, err := validator.ValidateHCLContent("project", string(content))
	if err != nil {
		pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
			Message: fmt.Sprintf("project.hcl schema validation failed: %v", err),
			File:    projectHCLPath,
			Context: "Project configuration validation",
		})
		return
	}

	if !result.IsValid {
		// Add schema validation errors with file context
		for _, schemaErr := range result.Errors {
			schemaErr.File = projectHCLPath
			schemaErr.Context = "Schema validation"
			pv.result.Errors = append(pv.result.Errors, schemaErr)
		}

		for _, schemaWarning := range result.Warnings {
			schemaWarning.File = projectHCLPath
			schemaWarning.Context = "Schema validation"
			pv.result.Warnings = append(pv.result.Warnings, schemaWarning)
		}
	}
}

// validateMachines validates machine configurations
func (pv *ProjectValidator) validateMachines(projectPath string) {
	machinesHCLPath := filepath.Join(projectPath, "machines.hcl")
	machinesDirPath := filepath.Join(projectPath, "machines")

	// Always validate machines.hcl if it exists
	if _, err := os.Stat(machinesHCLPath); err == nil {
		pv.validateMachinesFile(machinesHCLPath)
	}

	// Always validate machines/ directory if it exists
	if _, err := os.Stat(machinesDirPath); err == nil {
		pv.validateMachinesDirectory(machinesDirPath)
	}
}

// validateMachinesFile validates a single machines.hcl file
func (pv *ProjectValidator) validateMachinesFile(filePath string) {
	validator := schemas.NewStructValidator()
	content, err := os.ReadFile(filePath)
	if err != nil {
		pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
			Message: fmt.Sprintf("failed to read %s: %v", filepath.Base(filePath), err),
			File:    filePath,
			Context: "Machines validation",
		})
		return
	}

	result, err := validator.ValidateHCLContent("machines", string(content))
	if err != nil {
		pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
			Message: fmt.Sprintf("schema validation failed: %v", err),
			File:    filePath,
			Context: "Machines validation",
		})
		return
	}

	if !result.IsValid {
		for _, schemaErr := range result.Errors {
			schemaErr.File = filePath
			schemaErr.Context = "Schema validation"
			pv.result.Errors = append(pv.result.Errors, schemaErr)
		}
	}
}

// validateMachinesDirectory validates all .hcl files in the machines directory
func (pv *ProjectValidator) validateMachinesDirectory(dirPath string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
			Message: fmt.Sprintf("failed to read machines directory: %v", err),
			File:    dirPath,
			Context: "Machines validation",
		})
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".hcl") {
			filePath := filepath.Join(dirPath, entry.Name())
			pv.validateMachinesFile(filePath)
		}
	}
}

// validateActions validates action configurations
func (pv *ProjectValidator) validateActions(projectPath string) {
	actionsHCLPath := filepath.Join(projectPath, "actions.hcl")
	actionsDirPath := filepath.Join(projectPath, "actions")

	// Always validate actions.hcl if it exists
	if _, err := os.Stat(actionsHCLPath); err == nil {
		pv.validateActionsFile(actionsHCLPath)
	}

	// Always validate actions/ directory if it exists
	if _, err := os.Stat(actionsDirPath); err == nil {
		pv.validateActionsDirectory(actionsDirPath)
	}

	// Also check for any action-related .hcl files in the root directory
	pv.validateActionFilesInRoot(projectPath)
}

// validateActionsFile validates a single actions.hcl file
func (pv *ProjectValidator) validateActionsFile(filePath string) {
	validator := schemas.NewStructValidator()
	content, err := os.ReadFile(filePath)
	if err != nil {
		pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
			Message: fmt.Sprintf("failed to read %s: %v", filepath.Base(filePath), err),
			File:    filePath,
			Context: "Actions validation",
		})
		return
	}

	result, err := validator.ValidateHCLContent("actions", string(content))
	if err != nil {
		pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
			Message: fmt.Sprintf("schema validation failed: %v", err),
			File:    filePath,
			Context: "Actions validation",
		})
		return
	}

	if !result.IsValid {
		for _, schemaErr := range result.Errors {
			schemaErr.File = filePath
			schemaErr.Context = "Schema validation"
			pv.result.Errors = append(pv.result.Errors, schemaErr)
		}
	}
}

// validateActionsDirectory validates all .hcl files in the actions directory
func (pv *ProjectValidator) validateActionsDirectory(dirPath string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
			Message: fmt.Sprintf("failed to read actions directory: %v", err),
			File:    dirPath,
			Context: "Actions validation",
		})
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".hcl") {
			filePath := filepath.Join(dirPath, entry.Name())
			pv.validateActionsFile(filePath)
		}
	}
}

// validateActionFilesInRoot validates action-related .hcl files in the root directory
func (pv *ProjectValidator) validateActionFilesInRoot(projectPath string) {
	entries, err := os.ReadDir(projectPath)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".hcl") {
			fileName := entry.Name()
			// Check if this looks like an action-related file
			if strings.Contains(fileName, "action") || strings.Contains(fileName, "invalid") {
				filePath := filepath.Join(projectPath, fileName)
				pv.validateActionsFile(filePath)
			}
		}
	}
}

// validateVariablesConfig validates variable configurations
func (pv *ProjectValidator) validateVariablesConfig(projectPath string) {
	variablesHCLPath := filepath.Join(projectPath, "variables.hcl")
	variablesDirPath := filepath.Join(projectPath, "variables")

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
	validator := schemas.NewStructValidator()
	content, err := os.ReadFile(filePath)
	if err != nil {
		pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
			Message: fmt.Sprintf("failed to read %s: %v", filepath.Base(filePath), err),
			File:    filePath,
			Context: "Variables validation",
		})
		return
	}

	result, err := validator.ValidateHCLContent("variables", string(content))
	if err != nil {
		pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
			Message: fmt.Sprintf("schema validation failed: %v", err),
			File:    filePath,
			Context: "Variables validation",
		})
		return
	}

	if !result.IsValid {
		for _, schemaErr := range result.Errors {
			schemaErr.File = filePath
			schemaErr.Context = "Schema validation"
			pv.result.Errors = append(pv.result.Errors, schemaErr)
		}
	}
}

// validateVariablesDirectory validates all .hcl files in the variables directory
func (pv *ProjectValidator) validateVariablesDirectory(dirPath string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
			Message: fmt.Sprintf("failed to read variables directory: %v", err),
			File:    dirPath,
			Context: "Variables validation",
		})
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".hcl") {
			filePath := filepath.Join(dirPath, entry.Name())
			pv.validateVariablesFile(filePath)
		}
	}
}
