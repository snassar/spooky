package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"spooky/internal/schemas"

	"github.com/spf13/cobra"
)

var (
	projectName        string
	projectDescription string
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage spooky projects",
	Long: `Manage spooky projects including initialization, configuration, and management.

A spooky project is a directory containing configuration files that define
automation tasks, machine inventory, and deployment configurations.`,
}

var projectInitCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Initialize a new spooky project",
	Long: `Initialize a new spooky project in the specified directory.

This command creates a new project directory with the necessary configuration
files including project.hcl, machines.hcl, actions.hcl, and variables.hcl.

The directory will be created if it doesn't exist.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runProjectInit,
}

var projectValidateCmd = &cobra.Command{
	Use:   "validate [directory]",
	Short: "Validate a spooky project",
	Long: `Validate a spooky project directory structure and configuration files.

This command checks:
- Project directory structure compliance
- HCL file syntax validation
- Schema compliance for machines, actions, and variables
- File and directory existence rules

The directory defaults to the current directory if not specified.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runProjectValidate,
}

func init() {
	// Add project command to root
	RootCmd.AddCommand(projectCmd)

	// Add init subcommand to project
	projectCmd.AddCommand(projectInitCmd)

	// Add validate subcommand to project
	projectCmd.AddCommand(projectValidateCmd)

	// Add flags for project init
	projectInitCmd.Flags().StringVar(&projectName, "name", "", "Project name (required)")
	projectInitCmd.Flags().StringVar(&projectDescription, "description", "", "Project description")

	// Mark name as required
	projectInitCmd.MarkFlagRequired("name")
}

func runProjectInit(cmd *cobra.Command, args []string) error {
	// Determine target directory
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	// Validate project name
	if err := validateProjectName(projectName); err != nil {
		return fmt.Errorf("invalid project name: %w", err)
	}

	// Create target directory if it doesn't exist
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
	}

	// Create project.hcl
	if err := createProjectHCL(targetDir, projectName, projectDescription); err != nil {
		return fmt.Errorf("failed to create project.hcl: %w", err)
	}

	// Create machines.hcl
	if err := createMachinesHCL(targetDir); err != nil {
		return fmt.Errorf("failed to create machines.hcl: %w", err)
	}

	// Create actions.hcl
	if err := createActionsHCL(targetDir); err != nil {
		return fmt.Errorf("failed to create actions.hcl: %w", err)
	}

	// Create variables.hcl
	if err := createVariablesHCL(targetDir); err != nil {
		return fmt.Errorf("failed to create variables.hcl: %w", err)
	}

	// Create README.md
	if err := createREADME(targetDir, projectName, projectDescription); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	fmt.Printf("✅ Successfully initialized spooky project '%s' in %s\n", projectName, targetDir)
	fmt.Printf("📁 Project files created:\n")
	fmt.Printf("   - project.hcl\n")
	fmt.Printf("   - machines.hcl\n")
	fmt.Printf("   - actions.hcl\n")
	fmt.Printf("   - variables.hcl\n")
	fmt.Printf("   - README.md\n")
	fmt.Printf("\n🚀 Next steps:\n")
	fmt.Printf("   1. cd %s\n", targetDir)
	fmt.Printf("   2. Edit the configuration files as needed\n")
	fmt.Printf("   3. Run 'spooky project validate' to check configuration\n")

	return nil
}

func validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	if len(name) > 128 {
		return fmt.Errorf("project name too long (max 128 characters)")
	}

	// Check if name matches the pattern from schema: ^[a-zA-Z][a-zA-Z0-9._-]*$
	if !regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9._-]*$").MatchString(name) {
		return fmt.Errorf("project name must start with a letter and contain only letters, numbers, dots, underscores, and hyphens")
	}

	return nil
}

func createProjectHCL(targetDir, name, description string) error {
	// Get the embedded project schema
	embedder, err := schemas.NewSchemaEmbedder()
	if err != nil {
		return fmt.Errorf("failed to initialize schema embedder: %w", err)
	}

	// Get the project schema structure
	schema, exists := embedder.GetSchema("project")
	if !exists {
		return fmt.Errorf("project schema not found in embedded schemas")
	}

	// Create minimal valid project.hcl based on the schema
	content := fmt.Sprintf(`# Spooky Project Configuration
# Generated on %s
# Based on embedded schema: project

%s

project {
  name = "%s"
  description = "%s"
}
`, time.Now().Format("2006-01-02 15:04:05"), schema, name, description)

	return os.WriteFile(filepath.Join(targetDir, "project.hcl"), []byte(content), 0644)
}

func createMachinesHCL(targetDir string) error {
	// Get the embedded machines schema
	embedder, err := schemas.NewSchemaEmbedder()
	if err != nil {
		return fmt.Errorf("failed to initialize schema embedder: %w", err)
	}

	// Get the machines schema structure
	schema, exists := embedder.GetSchema("machines")
	if !exists {
		return fmt.Errorf("machines schema not found in embedded schemas")
	}

	// Create minimal valid machines.hcl based on the schema
	content := fmt.Sprintf(`# Machines Configuration
# Define your machine inventory and connectivity settings
# Based on embedded schema: machines

%s
`, schema)

	return os.WriteFile(filepath.Join(targetDir, "machines.hcl"), []byte(content), 0644)
}

func createActionsHCL(targetDir string) error {
	// Get the embedded actions schema
	embedder, err := schemas.NewSchemaEmbedder()
	if err != nil {
		return fmt.Errorf("failed to initialize schema embedder: %w", err)
	}

	// Get the actions schema structure
	schema, exists := embedder.GetSchema("actions")
	if !exists {
		return fmt.Errorf("actions schema not found in embedded schemas")
	}

	// Create minimal valid actions.hcl based on the schema
	content := fmt.Sprintf(`# Actions Configuration
# Define automation tasks and deployment actions
# Based on embedded schema: actions

%s
`, schema)

	return os.WriteFile(filepath.Join(targetDir, "actions.hcl"), []byte(content), 0644)
}

func createVariablesHCL(targetDir string) error {
	// Get the embedded variables schema
	embedder, err := schemas.NewSchemaEmbedder()
	if err != nil {
		return fmt.Errorf("failed to initialize schema embedder: %w", err)
	}

	// Get the variables schema structure
	schema, exists := embedder.GetSchema("variables")
	if !exists {
		return fmt.Errorf("variables schema not found in embedded schemas")
	}

	// Create minimal valid variables.hcl based on the schema
	content := fmt.Sprintf(`# Variables Configuration
# Define project-wide variables and configuration values
# Based on embedded schema: variables

%s
`, schema)

	return os.WriteFile(filepath.Join(targetDir, "variables.hcl"), []byte(content), 0644)
}

func createREADME(targetDir, name, description string) error {
	content := fmt.Sprintf("# %s\n\n%s\n\n## Overview\n\nThis is a Spooky automation project that defines configuration management, \ndeployment automation, and infrastructure management tasks.\n\n## Project Structure\n\n- project.hcl - Project configuration and metadata\n- machines.hcl - Machine inventory and connectivity settings\n- actions.hcl - Automation tasks and deployment actions\n- variables.hcl - Project-wide variables and configuration values\n- templates/ - Template files for deployment (create as needed)\n- files/ - Static files for deployment (create as needed)\n\n## Getting Started\n\n1. **Configure Machines**: Edit machines.hcl to define your target machines\n2. **Define Actions**: Edit actions.hcl to create automation tasks\n3. **Set Variables**: Edit variables.hcl to configure project variables\n4. **Validate**: Run 'spooky project validate' to check configuration\n5. **Execute**: Run 'spooky run <action-name>' to execute actions\n\n## Examples\n\n### Running Actions\n```bash\n# Run a specific action\nspooky run deploy-application\n\n# Run actions with specific tags\nspooky run --tags deployment\n\n# Dry run to see what would happen\nspooky run --dry-run deploy-application\n```\n\n### Managing Machines\n```bash\n# List all machines\nspooky machines list\n\n# Test connectivity\nspooky machines test-connection\n\n# Collect facts\nspooky machines collect-facts\n```\n\n## Documentation\n\nFor more information about Spooky, visit the project documentation or run:\n```bash\nspooky --help\nspooky project --help\nspooky machines --help\nspooky actions --help\n```\n\n## Support\n\nIf you encounter issues or have questions, please refer to the Spooky documentation\nor create an issue in the project repository.", name, description)

	return os.WriteFile(filepath.Join(targetDir, "README.md"), []byte(content), 0644)
}

func runProjectValidate(cmd *cobra.Command, args []string) error {
	// Determine target directory
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	// Check if directory exists
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", targetDir)
	}

	fmt.Printf("Validating spooky project in: %s\n\n", targetDir)

	// Validate project structure
	if err := validateProjectStructure(targetDir); err != nil {
		return fmt.Errorf("project structure validation failed: %w", err)
	}

	// Validate individual files
	if err := validateProjectFiles(targetDir); err != nil {
		return fmt.Errorf("project files validation failed: %w", err)
	}

	fmt.Println("✅ Project validation completed successfully!")
	return nil
}

func validateProjectStructure(targetDir string) error {
	fmt.Println("📁 Validating project structure...")

	// Check for required project.hcl
	projectHCLPath := filepath.Join(targetDir, "project.hcl")
	if _, err := os.Stat(projectHCLPath); os.IsNotExist(err) {
		return fmt.Errorf("required file missing: project.hcl")
	}
	fmt.Println("  ✅ project.hcl found")

	// Check for machines configuration (file or directory)
	machinesHCLPath := filepath.Join(targetDir, "machines.hcl")
	machinesDirPath := filepath.Join(targetDir, "machines")

	if _, err := os.Stat(machinesHCLPath); os.IsNotExist(err) {
		if _, err := os.Stat(machinesDirPath); os.IsNotExist(err) {
			return fmt.Errorf("either machines.hcl file or machines/ directory must exist")
		}
		fmt.Println("  ✅ machines/ directory found")
	} else {
		fmt.Println("  ✅ machines.hcl file found")
	}

	// Check for actions configuration (file or directory)
	actionsHCLPath := filepath.Join(targetDir, "actions.hcl")
	actionsDirPath := filepath.Join(targetDir, "actions")

	if _, err := os.Stat(actionsHCLPath); os.IsNotExist(err) {
		if _, err := os.Stat(actionsDirPath); os.IsNotExist(err) {
			return fmt.Errorf("either actions.hcl file or actions/ directory must exist")
		}
		fmt.Println("  ✅ actions/ directory found")
	} else {
		fmt.Println("  ✅ actions.hcl file found")
	}

	// Check for variables configuration (file or directory)
	variablesHCLPath := filepath.Join(targetDir, "variables.hcl")
	variablesDirPath := filepath.Join(targetDir, "variables")

	if _, err := os.Stat(variablesHCLPath); os.IsNotExist(err) {
		if _, err := os.Stat(variablesDirPath); os.IsNotExist(err) {
			return fmt.Errorf("either variables.hcl file or variables/ directory must exist")
		}
		fmt.Println("  ✅ variables/ directory found")
	} else {
		fmt.Println("  ✅ variables.hcl file found")
	}

	// Check optional directories
	optionalDirs := []string{"templates", "files", "logs"}
	for _, dir := range optionalDirs {
		dirPath := filepath.Join(targetDir, dir)
		if _, err := os.Stat(dirPath); err == nil {
			fmt.Printf("  ℹ️  %s/ directory found (optional)\n", dir)
		}
	}

	fmt.Println("  ✅ Project structure validation passed")
	return nil
}

func validateProjectFiles(targetDir string) error {
	fmt.Printf("\n🔍 Validating project files...\n")

	// Validate project.hcl
	if err := validateProjectConfig(targetDir); err != nil {
		return fmt.Errorf("project configuration validation failed: %w", err)
	}

	// Validate machines
	if err := validateMachines(targetDir); err != nil {
		return fmt.Errorf("machines validation failed: %w", err)
	}

	// Validate actions
	if err := validateActions(targetDir); err != nil {
		return fmt.Errorf("actions validation failed: %w", err)
	}

	// Validate variables
	if err := validateVariablesConfig(targetDir); err != nil {
		return fmt.Errorf("variables validation failed: %w", err)
	}

	fmt.Printf("✅ All project files validation completed successfully\n")
	return nil
}

func validateProjectConfig(targetDir string) error {
	projectHCLPath := filepath.Join(targetDir, "project.hcl")

	if _, err := os.Stat(projectHCLPath); os.IsNotExist(err) {
		return fmt.Errorf("project.hcl not found")
	}

	// Validate syntax and schema
	validator, err := schemas.NewSchemaValidator()
	if err != nil {
		return fmt.Errorf("failed to create schema validator: %w", err)
	}
	content, err := os.ReadFile(projectHCLPath)
	if err != nil {
		return fmt.Errorf("failed to read project.hcl: %w", err)
	}

	// Validate against project schema
	result, err := validator.ValidateContent("project", string(content))
	if err != nil {
		return fmt.Errorf("project.hcl schema validation failed: %w", err)
	}
	if !result.IsValid {
		return fmt.Errorf("project.hcl schema validation failed: %d errors", len(result.Errors))
	}

	fmt.Printf("✅ project.hcl - syntax and schema valid\n")
	return nil
}

func validateMachines(projectPath string) error {
	// Use the new merged validation approach
	return mergeAndValidateMachines(projectPath)
}

func validateActions(projectPath string) error {
	actionsHCLPath := filepath.Join(projectPath, "actions.hcl")
	actionsDirPath := filepath.Join(projectPath, "actions")

	// Check if actions.hcl exists
	if _, err := os.Stat(actionsHCLPath); err == nil {
		// Read file content for schema validation
		content, err := os.ReadFile(actionsHCLPath)
		if err != nil {
			return fmt.Errorf("failed to read actions.hcl: %w", err)
		}

		// Validate against embedded schema
		if err := validateAgainstSchema("actions", string(content)); err != nil {
			return fmt.Errorf("actions.hcl schema validation failed: %w", err)
		}

		return nil
	}

	// Check if actions/ directory exists
	if _, err := os.Stat(actionsDirPath); err == nil {
		// Use the new merged validation approach for actions directory
		return mergeAndValidateActions(projectPath)
	}

	return fmt.Errorf("neither actions.hcl nor actions/ directory found")
}

// validateAgainstSchema validates HCL content against embedded schemas
func validateAgainstSchema(schemaName, content string) error {
	// Use the new schema validator for proper schema validation
	validator, err := schemas.NewSchemaValidator()
	if err != nil {
		return fmt.Errorf("failed to create schema validator: %w", err)
	}

	result, err := validator.ValidateContent(schemaName, content)
	if err != nil {
		return fmt.Errorf("failed to validate schema: %w", err)
	}

	if !result.IsValid {
		if len(result.Errors) > 0 {
			return fmt.Errorf("schema validation failed: %s", result.Errors[0].Message)
		}
		return fmt.Errorf("schema validation failed")
	}

	return nil
}

// validateBasicStructure performs basic structural validation
func validateBasicStructure(schemaName, content string) error {
	switch schemaName {
	case "project":
		return validateProjectContent(content)
	case "machines":
		return validateMachinesStructure(content)
	case "actions":
		return validateActionsStructure(content)
	case "variables":
		return validateVariablesStructure(content)
	default:
		return fmt.Errorf("unknown schema type: %s", schemaName)
	}
}

// mergeAndValidateMachines merges multiple machines HCL files and validates for collisions
func mergeAndValidateMachines(targetDir string) error {
	machinesHCLPath := filepath.Join(targetDir, "machines.hcl")
	machinesDirPath := filepath.Join(targetDir, "machines")

	var allContents []string
	var fileNames []string

	// Check if machines.hcl exists
	if _, err := os.Stat(machinesHCLPath); err == nil {
		content, err := os.ReadFile(machinesHCLPath)
		if err != nil {
			return fmt.Errorf("failed to read machines.hcl: %w", err)
		}
		allContents = append(allContents, string(content))
		fileNames = append(fileNames, "machines.hcl")
	}

	// Check if machines/ directory exists
	if _, err := os.Stat(machinesDirPath); err == nil {

		// Read all .hcl files in directory
		entries, err := os.ReadDir(machinesDirPath)
		if err != nil {
			return fmt.Errorf("failed to read machines directory: %w", err)
		}

		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".hcl") {
				filePath := filepath.Join(machinesDirPath, entry.Name())

				content, err := os.ReadFile(filePath)
				if err != nil {
					return fmt.Errorf("failed to read %s: %w", entry.Name(), err)
				}
				allContents = append(allContents, string(content))
				fileNames = append(fileNames, "machines/"+entry.Name())
			}
		}

	}

	if len(allContents) == 0 {
		return fmt.Errorf("neither machines.hcl nor machines/ directory found")
	}

	// Merge all contents and validate
	return validateMergedMachines(allContents, fileNames)
}

func validateMergedMachines(contents []string, fileNames []string) error {
	// Parse HCL and extract machine information using proper HCL parsing
	machineMap := make(map[string]machineInfo)

	for i, content := range contents {
		fileName := fileNames[i]

		// Use schema-aware parsing for collision detection
		validator, err := schemas.NewSchemaValidator()
		if err != nil {
			return fmt.Errorf("failed to create schema validator: %w", err)
		}

		// Parse the content using the schema validator
		parsedData, err := validator.ParseHCLContent(content)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", fileName, err)
		}

		// Extract machine names from parsed data (schema-aware)
		machines, err := extractMachineNamesFromParsedData(parsedData, fileName)
		if err != nil {
			return fmt.Errorf("failed to extract machines from %s: %w", fileName, err)
		}

		// Check for duplicates and conflicts
		for _, machine := range machines {
			if existing, exists := machineMap[machine.name]; exists {
				return fmt.Errorf("machine name collision: '%s' defined in both %s and %s",
					machine.name, existing.fileName, fileName)
			}
			machineMap[machine.name] = machine
		}
	}

	// Now validate the merged content against the schema
	mergedContent := mergeMachinesContent(contents)

	if err := validateAgainstSchema("machines", mergedContent); err != nil {
		return fmt.Errorf("merged machines schema validation failed: %w", err)
	}

	return nil
}

type machineInfo struct {
	name     string
	fileName string
	hostname string
	user     string
}

// extractMachinesFromParsedData extracts machine information from parsed HCL data
func extractMachinesFromParsedData(data map[string]interface{}, fileName string) ([]machineInfo, error) {
	var machines []machineInfo

	// Look for machines block
	machinesBlock, exists := data["machines"]
	if !exists {
		return nil, fmt.Errorf("missing machines block in %s", fileName)
	}

	// Handle machines block (should be a map)
	machinesMap, ok := machinesBlock.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected machines block type in %s: %T", fileName, machinesBlock)
	}

	// Look for machine key (contains array of machine blocks)
	machineValue, exists := machinesMap["machine"]
	if !exists {
		return machines, nil // No machines defined, return empty slice
	}

	// Handle machine blocks array
	if machineBlocks, ok := machineValue.([]interface{}); ok {

		// Multiple machine blocks
		for i, block := range machineBlocks {

			if machineBlock, ok := block.(map[string]interface{}); ok {
				machine := extractMachineFromBlock(machineBlock, fileName, fmt.Sprintf("[%d]", i))

				if machine.name != "" {
					machines = append(machines, machine)
				}
			}
		}
	} else if machineBlock, ok := machineValue.(map[string]interface{}); ok {

		// Single machine block
		machine := extractMachineFromBlock(machineBlock, fileName, "")

		if machine.name != "" {
			machines = append(machines, machine)
		}
	} else {

	}

	return machines, nil
}

// extractMachineFromBlock extracts machine information from a single machine block
func extractMachineFromBlock(block map[string]interface{}, fileName, blockIndex string) machineInfo {
	machine := machineInfo{fileName: fileName + blockIndex}

	// Extract name - handle different value types from HCL parser
	if name, exists := block["name"]; exists {
		machine.name = convertToString(name)
	}

	// Extract hostname - handle different value types from HCL parser
	if hostname, exists := block["hostname"]; exists {
		machine.hostname = convertToString(hostname)
	}

	// Extract user - handle different value types from HCL parser
	if user, exists := block["user"]; exists {
		machine.user = convertToString(user)
	}

	return machine
}

// extractMachineNamesFromParsedData extracts machine information from parsed HCL data
func extractMachineNamesFromParsedData(data map[string]interface{}, fileName string) ([]machineInfo, error) {
	var machines []machineInfo

	// Look for machines block
	machinesBlock, exists := data["machines"]
	if !exists {
		return machines, nil // No machines block, return empty
	}

	// Convert to map
	machinesMap, ok := machinesBlock.(map[string]interface{})
	if !ok {
		return machines, nil // Invalid structure, return empty
	}

	// Look for machine array
	machineValue, exists := machinesMap["machine"]
	if !exists {

		return machines, nil // No machines defined
	}

	// Handle machine blocks array
	if machineBlocks, ok := machineValue.([]map[string]interface{}); ok {

		// Multiple machine blocks (correct type assertion)
		for _, machineBlock := range machineBlocks {

			if name, exists := machineBlock["name"]; exists {
				nameStr := convertToString(name)
				// Filter out placeholder values that indicate parsing errors
				if nameStr != "" && nameStr != "complex_expression" && !strings.HasPrefix(nameStr, "complex_") {
					machines = append(machines, machineInfo{
						name:     nameStr,
						fileName: fileName,
					})
				}
			}
		}
	} else if machineBlocks, ok := machineValue.([]interface{}); ok {

		// Multiple machine blocks (fallback for interface array)
		for _, block := range machineBlocks {
			if machineBlock, ok := block.(map[string]interface{}); ok {
				if name, exists := machineBlock["name"]; exists {
					nameStr := convertToString(name)
					// Filter out placeholder values that indicate parsing errors
					if nameStr != "" && nameStr != "complex_expression" && !strings.HasPrefix(nameStr, "complex_") {
						machines = append(machines, machineInfo{
							name:     nameStr,
							fileName: fileName,
						})
					}
				}
			}
		}
	} else if machineBlock, ok := machineValue.(map[string]interface{}); ok {
		// Single machine block
		if name, exists := machineBlock["name"]; exists {
			nameStr := convertToString(name)
			// Filter out placeholder values that indicate parsing errors
			if nameStr != "" && nameStr != "complex_expression" && !strings.HasPrefix(nameStr, "complex_") {
				machines = append(machines, machineInfo{
					name:     nameStr,
					fileName: fileName,
				})
			}
		}
	}

	return machines, nil
}

// convertToString converts various HCL value types to string
func convertToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

// mergeAndValidateVariables merges multiple variables HCL files and validates for collisions
func mergeAndValidateVariables(targetDir string) error {
	variablesHCLPath := filepath.Join(targetDir, "variables.hcl")
	variablesDirPath := filepath.Join(targetDir, "variables")

	var allContents []string
	var fileNames []string

	// Check if variables.hcl exists
	if _, err := os.Stat(variablesHCLPath); err == nil {
		content, err := os.ReadFile(variablesHCLPath)
		if err != nil {
			return fmt.Errorf("failed to read variables.hcl: %w", err)
		}
		allContents = append(allContents, string(content))
		fileNames = append(fileNames, "variables.hcl")
	}

	// Check if variables/ directory exists
	if _, err := os.Stat(variablesDirPath); err == nil {
		// Read all .hcl files in directory
		entries, err := os.ReadDir(variablesDirPath)
		if err != nil {
			return fmt.Errorf("failed to read variables directory: %w", err)
		}

		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".hcl") {
				filePath := filepath.Join(variablesDirPath, entry.Name())
				content, err := os.ReadFile(filePath)
				if err != nil {
					return fmt.Errorf("failed to read %s: %w", entry.Name(), err)
				}
				allContents = append(allContents, string(content))
				fileNames = append(fileNames, "variables/"+entry.Name())
			}
		}
	}

	if len(allContents) == 0 {
		return fmt.Errorf("neither variables.hcl nor variables/ directory found")
	}

	// Merge all contents and validate
	return validateMergedVariables(allContents, fileNames)
}

func validateMergedVariables(contents []string, fileNames []string) error {
	// Parse HCL and extract variable information using proper HCL parsing
	variableMap := make(map[string]variableInfo)

	for i, content := range contents {
		fileName := fileNames[i]

		// Use schema-aware parsing for collision detection
		validator, err := schemas.NewSchemaValidator()
		if err != nil {
			return fmt.Errorf("failed to create schema validator: %w", err)
		}

		// Parse the content using the schema validator
		parsedData, err := validator.ParseHCLContent(content)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", fileName, err)
		}

		// Extract variable names from parsed data (schema-aware)
		variables, err := extractVariableNamesFromParsedData(parsedData, fileName)
		if err != nil {
			return fmt.Errorf("failed to extract variables from %s: %w", fileName, err)
		}

		// Check for duplicates and conflicts
		for _, variable := range variables {
			if existing, exists := variableMap[variable.name]; exists {
				return fmt.Errorf("variable name collision: '%s' defined in both %s and %s",
					variable.name, existing.fileName, fileName)
			}
			variableMap[variable.name] = variable
		}
	}

	// Now validate the merged content against the schema
	mergedContent := mergeVariablesContent(contents)

	if err := validateAgainstSchema("variables", mergedContent); err != nil {
		return fmt.Errorf("merged variables schema validation failed: %w", err)
	}

	return nil
}

type variableInfo struct {
	name        string
	fileName    string
	description string
	sensitive   bool
	encrypted   bool
}

func validateProjectContent(content string) error {
	// Check for required metadata block
	if !strings.Contains(content, "metadata {") {
		return fmt.Errorf("missing required metadata block")
	}

	// Check for required project block
	if !strings.Contains(content, "project {") {
		return fmt.Errorf("missing required project block")
	}

	// Check for required name field
	if !strings.Contains(content, "name =") {
		return fmt.Errorf("missing required 'name' field in project block")
	}

	return nil
}

func validateMachinesStructure(content string) error {
	// Check for required metadata block
	if !strings.Contains(content, "metadata {") {
		return fmt.Errorf("missing required metadata block")
	}

	// Check for required machines block
	if !strings.Contains(content, "machines {") {
		return fmt.Errorf("missing required machines block")
	}

	return nil
}

func validateActionsStructure(content string) error {
	// Check for required metadata block
	if !strings.Contains(content, "metadata {") {
		return fmt.Errorf("missing required metadata block")
	}

	// Check for required actions block
	if !strings.Contains(content, "actions {") {
		return fmt.Errorf("missing required actions block")
	}

	// Check for at least one action block
	if !strings.Contains(content, "action \"") {
		return fmt.Errorf("missing required action block")
	}

	return nil
}

func validateVariablesStructure(content string) error {
	// Check for required metadata block
	if !strings.Contains(content, "metadata {") {
		return fmt.Errorf("missing required metadata block")
	}

	return nil
}

func validateVariablesConfig(targetDir string) error {
	// Use the new merged validation approach
	return mergeAndValidateVariables(targetDir)
}

// mergeAndValidateActions merges multiple actions HCL files and validates for collisions
func mergeAndValidateActions(targetDir string) error {
	actionsHCLPath := filepath.Join(targetDir, "actions.hcl")
	actionsDirPath := filepath.Join(targetDir, "actions")

	var allContents []string
	var fileNames []string

	// Check if actions.hcl exists
	if _, err := os.Stat(actionsHCLPath); err == nil {
		content, err := os.ReadFile(actionsHCLPath)
		if err != nil {
			return fmt.Errorf("failed to read actions.hcl: %w", err)
		}
		allContents = append(allContents, string(content))
		fileNames = append(fileNames, "actions.hcl")
	}

	// Check if actions/ directory exists
	if _, err := os.Stat(actionsDirPath); err == nil {
		// Read all .hcl files in directory
		entries, err := os.ReadDir(actionsDirPath)
		if err != nil {
			return fmt.Errorf("failed to read actions directory: %w", err)
		}

		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".hcl") {
				filePath := filepath.Join(actionsDirPath, entry.Name())
				content, err := os.ReadFile(filePath)
				if err != nil {
					return fmt.Errorf("failed to read %s: %w", entry.Name(), err)
				}
				allContents = append(allContents, string(content))
				fileNames = append(fileNames, "actions/"+entry.Name())
			}
		}
	}

	if len(allContents) == 0 {
		return fmt.Errorf("neither actions.hcl nor actions/ directory found")
	}

	// Merge all contents and validate
	return validateMergedActions(allContents, fileNames)
}

func validateMergedActions(contents []string, fileNames []string) error {
	// Parse HCL and extract action information using proper HCL parsing
	actionMap := make(map[string]actionInfo)

	for i, content := range contents {
		fileName := fileNames[i]

		// Use schema-aware parsing for collision detection
		validator, err := schemas.NewSchemaValidator()
		if err != nil {
			return fmt.Errorf("failed to create schema validator: %w", err)
		}

		// Parse the content using the schema validator
		parsedData, err := validator.ParseHCLContent(content)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", fileName, err)
		}

		// Extract action names from parsed data (schema-aware)
		actions, err := extractActionNamesFromParsedData(parsedData, fileName)
		if err != nil {
			return fmt.Errorf("failed to extract actions from %s: %w", fileName, err)
		}

		// Check for duplicates and conflicts
		for _, action := range actions {
			if existing, exists := actionMap[action.name]; exists {
				return fmt.Errorf("action name collision: '%s' defined in both %s and %s",
					action.name, existing.fileName, fileName)
			}
			actionMap[action.name] = action
		}
	}

	// Now validate the merged content against the schema
	mergedContent := mergeActionsContent(contents)
	if err := validateAgainstSchema("actions", mergedContent); err != nil {
		return fmt.Errorf("merged actions schema validation failed: %w", err)
	}

	return nil
}

type actionInfo struct {
	name        string
	fileName    string
	description string
	actionType  string
}

func mergeActionsContent(contents []string) string {
	// Create a merged actions.hcl content
	merged := `# Merged Actions Configuration
# Generated from multiple files

metadata {
  version = "1"
  description = "Merged actions configuration"
}

actions {
`

	// Extract action blocks from each content
	for _, content := range contents {
		// Find all action blocks with proper brace counting
		actionBlocks := extractActionBlocks(content)

		for _, block := range actionBlocks {
			// Clean up the block and add it to merged content
			cleaned := strings.TrimSpace(block)
			merged += "\n  " + cleaned + "\n"
		}
	}

	merged += "}\n"
	return merged
}

// extractActionBlocks extracts complete action blocks with proper brace counting
func extractActionBlocks(content string) []string {
	var blocks []string
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		// Look for action block start (action "name" {)
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "action \"") && strings.Contains(trimmed, "{") {
			// Find the complete block by counting braces
			block := extractCompleteBlock(lines, i)
			if block != "" {
				blocks = append(blocks, block)
			}
		}
	}

	return blocks
}

func mergeVariablesContent(contents []string) string {
	// Create a merged variables.hcl content
	merged := `# Merged Variables Configuration
# Generated from multiple files

metadata {
  version = "1"
  description = "Merged variables configuration"
}

variables {
`

	// Extract variable blocks from each content
	for _, content := range contents {
		// Find all variable blocks with proper brace counting
		variableBlocks := extractVariableBlocks(content)

		for _, block := range variableBlocks {
			// Clean up the block and add it to merged content
			cleaned := strings.TrimSpace(block)
			merged += "\n  " + cleaned + "\n"
		}
	}

	merged += "}\n"
	return merged
}

// extractVariableBlocks extracts complete variable blocks with proper brace counting
func extractVariableBlocks(content string) []string {
	var blocks []string
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		// Look for variable block start (variable {)
		if strings.TrimSpace(line) == "variable {" {
			// Find the complete block by counting braces
			block := extractCompleteBlock(lines, i)
			if block != "" {
				blocks = append(blocks, block)
			}
		}
	}

	return blocks
}

func mergeMachinesContent(contents []string) string {
	// Create a merged machines.hcl content
	merged := `# Merged Machines Configuration
# Generated from multiple files

metadata {
  version = "1"
  description = "Merged machines configuration"
}

machines {
`

	// Extract machine blocks from each content
	for _, content := range contents {
		// Find all machine blocks with proper brace counting
		machineBlocks := extractMachineBlocks(content)

		for _, block := range machineBlocks {
			// Clean up the block and add it to merged content
			cleaned := strings.TrimSpace(block)
			merged += "\n  " + cleaned + "\n"
		}
	}

	merged += "}\n"
	return merged
}

// extractMachineBlocks extracts complete machine blocks with proper brace counting
func extractMachineBlocks(content string) []string {
	var blocks []string
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		// Look for machine block start
		if strings.TrimSpace(line) == "machine {" {
			// Find the complete block by counting braces
			block := extractCompleteBlock(lines, i)
			if block != "" {
				blocks = append(blocks, block)
			}
		}
	}

	return blocks
}

// extractCompleteBlock extracts a complete block from a given starting line
func extractCompleteBlock(lines []string, startLine int) string {
	var blockLines []string
	braceCount := 0
	started := false

	for i := startLine; i < len(lines); i++ {
		line := lines[i]

		if !started {
			started = true
			braceCount = 0
		}

		// Count opening and closing braces
		for _, char := range line {
			if char == '{' {
				braceCount++
			} else if char == '}' {
				braceCount--
			}
		}

		blockLines = append(blockLines, line)

		// If we've closed all braces, we're done
		if braceCount == 0 && started {
			break
		}
	}

	if braceCount == 0 {
		return strings.Join(blockLines, "\n")
	}

	return "" // Incomplete block
}

// getKeys returns a slice of keys from a map[string]interface{}
func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// extractVariableNamesFromParsedData extracts variable information from parsed HCL data
func extractVariableNamesFromParsedData(data map[string]interface{}, fileName string) ([]variableInfo, error) {
	var variables []variableInfo

	// Look for variables block
	variablesBlock, exists := data["variables"]
	if !exists {
		return variables, nil // No variables block, return empty
	}

	// Convert to map
	variablesMap, ok := variablesBlock.(map[string]interface{})
	if !ok {
		return variables, nil // Invalid structure, return empty
	}

	// Look for variable blocks (variable { ... })
	variableBlocks, exists := variablesMap["variable"]
	if !exists {
		return variables, nil
	}

	// Handle array of variable blocks
	if variableArray, ok := variableBlocks.([]map[string]interface{}); ok {

		for _, block := range variableArray {
			if name, exists := block["name"]; exists {
				nameStr := convertToString(name)
				if nameStr != "" {
					variable := variableInfo{
						name:     nameStr,
						fileName: fileName,
					}

					// Extract description
					if desc, exists := block["description"]; exists {
						variable.description = convertToString(desc)
					}

					// Extract sensitive flag
					if sensitive, exists := block["sensitive"]; exists {
						if sensitiveBool, ok := sensitive.(bool); ok {
							variable.sensitive = sensitiveBool
						}
					}

					// Extract encrypted flag
					if encrypted, exists := block["encrypted"]; exists {
						if encryptedBool, ok := encrypted.(bool); ok {
							variable.encrypted = encryptedBool
						}
					}

					variables = append(variables, variable)
				}
			}
		}
	} else if variableBlock, ok := variableBlocks.(map[string]interface{}); ok {
		// Single variable block
		if name, exists := variableBlock["name"]; exists {
			nameStr := convertToString(name)
			if nameStr != "" {
				variable := variableInfo{
					name:     nameStr,
					fileName: fileName,
				}

				// Extract description
				if desc, exists := variableBlock["description"]; exists {
					variable.description = convertToString(desc)
				}

				// Extract sensitive flag
				if sensitive, exists := variableBlock["sensitive"]; exists {
					if sensitiveBool, ok := sensitive.(bool); ok {
						variable.sensitive = sensitiveBool
					}
				}

				// Extract encrypted flag
				if encrypted, exists := variableBlock["encrypted"]; exists {
					if encryptedBool, ok := encrypted.(bool); ok {
						variable.encrypted = encryptedBool
					}
				}

				variables = append(variables, variable)
			}
		}
	}

	return variables, nil
}

// extractActionNamesFromParsedData extracts action information from parsed HCL data
func extractActionNamesFromParsedData(data map[string]interface{}, fileName string) ([]actionInfo, error) {
	var actions []actionInfo

	// Look for actions block
	actionsBlock, exists := data["actions"]
	if !exists {
		fmt.Printf("      🔍 [%s] No actions block found\n", fileName)
		return actions, nil // No actions block, return empty
	}

	fmt.Printf("      🔍 [%s] Found actions block, type: %T\n", fileName, actionsBlock)

	// Convert to map
	actionsMap, ok := actionsBlock.(map[string]interface{})
	if !ok {
		fmt.Printf("      🔍 [%s] Actions block is not a map\n", fileName)
		return actions, nil // Invalid structure, return empty
	}

	fmt.Printf("      🔍 [%s] Actions map keys: %v\n", fileName, getKeys(actionsMap))

	// Look for action blocks (action "name" { ... })
	for actionName, actionValue := range actionsMap {
		fmt.Printf("      🔍 [%s] Found action: %s\n", fileName, actionName)

		// Check if it's an action block (map with metadata)
		if actionBlock, ok := actionValue.(map[string]interface{}); ok {
			action := actionInfo{
				name:     actionName,
				fileName: fileName,
			}

			// Extract description
			if desc, exists := actionBlock["description"]; exists {
				action.description = convertToString(desc)
			}

			// Extract action type
			if actionType, exists := actionBlock["type"]; exists {
				action.actionType = convertToString(actionType)
			}

			actions = append(actions, action)
		} else {
			// Simple action block without metadata
			action := actionInfo{
				name:     actionName,
				fileName: fileName,
			}
			actions = append(actions, action)
		}
	}

	return actions, nil
}
