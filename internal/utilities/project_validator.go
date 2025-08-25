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
	fmt.Printf("DEBUG: ValidateProject called for targetDir: %s\n", targetDir)

	// Validate project structure
	fmt.Printf("DEBUG: Calling validateProjectStructure\n")
	pv.validateProjectStructure(targetDir)

	// Validate individual files
	fmt.Printf("DEBUG: Calling validateProjectFiles\n")
	pv.validateProjectFiles(targetDir)

	// Set final validity
	pv.result.IsValid = len(pv.result.Errors) == 0
	fmt.Printf("DEBUG: Final validation result - IsValid: %v, Errors: %d\n", pv.result.IsValid, len(pv.result.Errors))

	return pv.result
}

// validateProjectStructure validates the project directory structure using schema-driven validation
func (pv *ProjectValidator) validateProjectStructure(targetDir string) {
	// Scan the actual project directory structure
	projectDir := pv.scanProjectDirectory(targetDir)

	// Validate against ProjectDirectoryV1 schema rules
	pv.validateProjectDirectoryAgainstSchema(projectDir, targetDir)

	// Additional validation for alternative configurations (either file or directory)
	pv.validateAlternativeConfigurations(targetDir)
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

// validateProjectDirectoryAgainstSchema validates the scanned directory structure against ProjectDirectoryV1 schema rules
func (pv *ProjectValidator) validateProjectDirectoryAgainstSchema(projectDir map[string]interface{}, targetDir string) {
	// Get the scanned files and directories
	files, _ := projectDir["files"].([]string)
	directories, _ := projectDir["directories"].([]string)

	// Validate required files
	pv.validateRequiredFiles(files, targetDir)

	// Validate optional files and directories
	pv.validateOptionalFiles(files, targetDir)
	pv.validateOptionalDirectories(directories, targetDir)
}

// validateRequiredFiles validates that all required files exist
func (pv *ProjectValidator) validateRequiredFiles(files []string, targetDir string) {
	requiredFiles := []string{"project.hcl"}

	for _, requiredFile := range requiredFiles {
		found := false
		for _, file := range files {
			if file == requiredFile {
				found = true
				break
			}
		}

		if !found {
			pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
				Message: fmt.Sprintf("required file missing: %s", requiredFile),
				File:    targetDir,
				Context: "Project structure validation",
			})
		}
	}
}

// validateOptionalFiles validates optional files and their patterns
func (pv *ProjectValidator) validateOptionalFiles(files []string, targetDir string) {
	// Optional files that can be validated for content patterns
	optionalFiles := map[string]string{
		"machines.hcl":  "machines {",
		"actions.hcl":   "actions {",
		"variables.hcl": "variables {",
		"README.md":     "# ",
	}

	for _, file := range files {
		if pattern, exists := optionalFiles[file]; exists {
			filePath := filepath.Join(targetDir, file)
			pv.validateFileContent(filePath, pattern, file)
		}
	}
}

// validateOptionalDirectories validates optional directories
func (pv *ProjectValidator) validateOptionalDirectories(directories []string, targetDir string) {
	// Optional directories that can be validated
	optionalDirs := []string{"machines", "actions", "variables", "templates", "files"}

	for _, dir := range directories {
		found := false
		for _, optionalDir := range optionalDirs {
			if dir == optionalDir {
				found = true
				break
			}
		}

		if !found {
			// Unknown directory - could be a warning
			pv.result.Warnings = append(pv.result.Warnings, schemas.ValidationWarning{
				Message: fmt.Sprintf("unknown directory: %s", dir),
				File:    targetDir,
				Context: "Project structure validation",
			})
		}
	}
}

// validateFileContent validates file content against expected patterns
func (pv *ProjectValidator) validateFileContent(filePath, pattern, fileName string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		// File exists but can't be read - this is a warning, not an error
		pv.result.Warnings = append(pv.result.Warnings, schemas.ValidationWarning{
			Message: fmt.Sprintf("cannot read file %s: %v", fileName, err),
			File:    filePath,
			Context: "Project structure validation",
		})
		return
	}

	// Check if content matches expected pattern
	if !strings.Contains(string(content), pattern) {
		pv.result.Warnings = append(pv.result.Warnings, schemas.ValidationWarning{
			Message: fmt.Sprintf("file %s does not contain expected pattern: %s", fileName, pattern),
			File:    filePath,
			Context: "Project structure validation",
		})
	}
}

// validateAlternativeConfigurations validates that at least one configuration option exists for each type
func (pv *ProjectValidator) validateAlternativeConfigurations(targetDir string) {
	// Check for machines configuration (either file or directory)
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

	// Check for actions configuration (either file or directory)
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

	// Check for variables configuration (either file or directory)
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

	// Skip validation if file doesn't exist (structure validation already caught this)
	if _, err := os.Stat(projectHCLPath); os.IsNotExist(err) {
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
	fmt.Printf("DEBUG: validateActions called for projectPath: %s\n", projectPath)

	actionsHCLPath := filepath.Join(projectPath, "actions.hcl")
	actionsDirPath := filepath.Join(projectPath, "actions")

	fmt.Printf("DEBUG: Checking for actionsHCLPath: %s\n", actionsHCLPath)
	fmt.Printf("DEBUG: Checking for actionsDirPath: %s\n", actionsDirPath)

	// Always validate actions.hcl if it exists
	if _, err := os.Stat(actionsHCLPath); err == nil {
		fmt.Printf("DEBUG: actions.hcl exists, calling validateActionsFile\n")
		pv.validateActionsFile(actionsHCLPath)
	} else {
		fmt.Printf("DEBUG: actions.hcl does not exist: %v\n", err)
	}

	// Always validate actions/ directory if it exists
	if _, err := os.Stat(actionsDirPath); err == nil {
		fmt.Printf("DEBUG: actions/ directory exists, calling validateActionsDirectory\n")
		pv.validateActionsDirectory(actionsDirPath)
	} else {
		fmt.Printf("DEBUG: actions/ directory does not exist: %v\n", err)
	}
}

// validateActionsFile validates a single actions.hcl file
func (pv *ProjectValidator) validateActionsFile(filePath string) {
	fmt.Printf("DEBUG: Validating actions file: %s\n", filePath)

	validator := schemas.NewStructValidator()
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("DEBUG: Failed to read file: %v\n", err)
		pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
			Message: fmt.Sprintf("failed to read %s: %v", filepath.Base(filePath), err),
			File:    filePath,
			Context: "Actions validation",
		})
		return
	}

	fmt.Printf("DEBUG: File content length: %d bytes\n", len(content))
	fmt.Printf("DEBUG: File content preview: %s\n", string(content[:min(100, len(content))]))

	result, err := validator.ValidateHCLContent("actions", string(content))
	if err != nil {
		fmt.Printf("DEBUG: Schema validation error: %v\n", err)
		pv.result.Errors = append(pv.result.Errors, schemas.ValidationError{
			Message: fmt.Sprintf("schema validation failed: %v", err),
			File:    filePath,
			Context: "Actions validation",
		})
		return
	}

	fmt.Printf("DEBUG: Schema validation result - IsValid: %v, Errors: %d, Warnings: %d\n",
		result.IsValid, len(result.Errors), len(result.Warnings))

	if !result.IsValid {
		for _, schemaErr := range result.Errors {
			fmt.Printf("DEBUG: Schema error: %s - %s\n", schemaErr.Field, schemaErr.Message)
			schemaErr.File = filePath
			schemaErr.Context = "Schema validation"
			pv.result.Errors = append(pv.result.Errors, schemaErr)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
