package commands

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"spooky/internal/encryption"
	"spooky/internal/logging"
	"spooky/internal/schemas"
	"spooky/internal/utilities"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
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

var projectEncryptCmd = &cobra.Command{
	Use:   "encrypt [directory]",
	Short: "Encrypt sensitive values in a spooky project",
	Long: `Encrypt sensitive values in a spooky project using age encryption.

This command:
- Finds variables marked with encrypted = true
- Encrypts their plaintext values using age encryption
- Updates the HCL files with encrypted values
- Requires age identities and recipients to be configured

The directory defaults to the current directory if not specified.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runProjectEncrypt,
}

var projectConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current configuration information",
	Long: `Show information about the current spooky configuration.

This command displays:
- Which configuration file is being used
- Configuration source (custom, user, or embedded default)
- Age encryption settings
- Configuration file details`,
	RunE: runProjectConfig,
}

func init() {
	// Add project command to root
	RootCmd.AddCommand(projectCmd)

	// Add init subcommand to project
	projectCmd.AddCommand(projectInitCmd)

	// Add validate subcommand to project
	projectCmd.AddCommand(projectValidateCmd)

	// Add encrypt subcommand to project
	projectCmd.AddCommand(projectEncryptCmd)

	// Add config subcommand to project
	projectCmd.AddCommand(projectConfigCmd)

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
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
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

	logger := logging.GetGlobalLogger()
	logger.Info("✅ Successfully initialized spooky project",
		slog.String("project_name", projectName),
		slog.String("target_directory", targetDir))

	logger.Info("📁 Project files created",
		slog.String("project_hcl", "project.hcl"),
		slog.String("machines_hcl", "machines.hcl"),
		slog.String("actions_hcl", "actions.hcl"),
		slog.String("variables_hcl", "variables.hcl"),
		slog.String("readme", "README.md"))

	logger.Info("🚀 Next steps",
		slog.String("step1", fmt.Sprintf("cd %s", targetDir)),
		slog.String("step2", "Edit the configuration files as needed"),
		slog.String("step3", "Run 'spooky project validate' to check configuration"))

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

// createProjectConfigFile creates a project configuration file with the specified content and filename
func createProjectConfigFile(targetDir, filename, content string) error {
	return os.WriteFile(filepath.Join(targetDir, filename), []byte(content), 0o644)
}

func createProjectHCL(targetDir, name, description string) error {
	// Generate project configuration directly from Go structs
	content := schemas.GenerateProjectConfigFromStructs(name, description)
	return createProjectConfigFile(targetDir, "project.hcl", content)
}

func createMachinesHCL(targetDir string) error {
	// Generate machines configuration directly from Go structs
	content := schemas.GenerateMachinesConfigFromStructs()
	return createProjectConfigFile(targetDir, "machines.hcl", content)
}

func createActionsHCL(targetDir string) error {
	// Generate actions configuration directly from Go structs
	content := schemas.GenerateActionsConfigFromStructs()
	return createProjectConfigFile(targetDir, "actions.hcl", content)
}

func createVariablesHCL(targetDir string) error {
	// Generate variables configuration directly from Go structs
	content := schemas.GenerateVariablesConfigFromStructs()
	return createProjectConfigFile(targetDir, "variables.hcl", content)
}

func createREADME(targetDir, name, description string) error {
	content := fmt.Sprintf("# %s\n\n%s\n\n## Overview\n\nThis is a Spooky automation project that defines configuration management, \ndeployment automation, and infrastructure management tasks.\n\n## Project Structure\n\n- project.hcl - Project configuration and metadata\n- machines.hcl - Machine inventory and connectivity settings\n- actions.hcl - Automation tasks and deployment actions\n- variables.hcl - Project-wide variables and configuration values\n- templates/ - Template files for deployment (create as needed)\n- files/ - Static files for deployment (create as needed)\n\n## Getting Started\n\n1. **Configure Machines**: Edit machines.hcl to define your target machines\n2. **Define Actions**: Edit actions.hcl to create automation tasks\n3. **Set Variables**: Edit variables.hcl to configure project variables\n4. **Validate**: Run 'spooky project validate' to check configuration\n5. **Execute**: Run 'spooky run <action-name>' to execute actions\n\n## Examples\n\n### Running Actions\n```bash\n# Run a specific action\nspooky run deploy-application\n\n# Run actions with specific tags\nspooky run --tags deployment\n\n# Dry run to see what would happen\nspooky run --dry-run deploy-application\n```\n\n### Managing Machines\n```bash\n# List all machines\nspooky machines list\n\n# Test connectivity\nspooky machines test-connection\n\n# Collect facts\nspooky machines collect-facts\n```\n\n## Documentation\n\nFor more information about Spooky, visit the project documentation or run:\n```bash\nspooky --help\nspooky project --help\nspooky machines --help\nspooky actions --help\n```\n\n## Support\n\nIf you encounter issues or have questions, please refer to the Spooky documentation\nor create an issue in the project repository.", name, description)

	return os.WriteFile(filepath.Join(targetDir, "README.md"), []byte(content), 0o644)
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

	// Use the enhanced project validator with schema-driven validation
	validator := utilities.NewProjectValidator()
	result := validator.ValidateProject(targetDir)

	// Display validation result
	if result.IsValid {
		// Try to get project name for display, but don't fail if it's missing
		projectName, _ := getProjectName(targetDir)
		if projectName == "" {
			projectName = "unknown"
		}
		logger := logging.GetGlobalLogger()
		logger.Info("Project validation successful", slog.String("project_name", projectName))
	} else {
		logger := logging.GetGlobalLogger()
		logger.Error("project validation failed",
			slog.Int("error_count", len(result.Errors)))
		for _, err := range result.Errors {
			logger.Error("validation error",
				slog.String("message", err.Message))
		}
		return fmt.Errorf("project validation failed with %d errors", len(result.Errors))
	}

	return nil
}

// validateAgainstSchema validates HCL content against embedded schemas
func validateAgainstSchema(schemaName, content string) error {
	// Use the simplified validator for essential schema validation
	validator := schemas.NewSimpleValidator()

	result, err := validator.ValidateHCLContent(schemaName, content)
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

// validateMergedContent validates merged HCL content for collisions and schema compliance
// This unified function replaces the duplicate validateMergedMachines, validateMergedVariables, and validateMergedActions functions
func validateMergedContent(contents []string, fileNames []string, contentType string) error {
	// Parse HCL and extract resource information using proper HCL parsing
	resourceMap := make(map[string]interface{})

	for i, content := range contents {
		fileName := fileNames[i]

		// Use simplified parsing for collision detection
		validator := schemas.NewSimpleValidator()

		// Parse the content using the simplified validator
		parsedData, err := validator.ParseHCLContent(content)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", fileName, err)
		}

		// Extract resource names from parsed data using unified function
		resources, err2 := extractResourceNamesFromParsedData(parsedData, fileName, contentType, getResourceType(contentType))
		if err2 != nil {
			return fmt.Errorf("failed to extract %s from %s: %w", contentType, fileName, err2)
		}

		// Check for duplicates and conflicts
		for _, resource := range resources {
			var name string
			var resourceFileName string

			// Extract name and fileName based on resource type
			switch r := resource.(type) {
			case machineInfo:
				name = r.name
				resourceFileName = r.fileName
			case variableInfo:
				name = r.name
				resourceFileName = r.fileName
			case actionInfo:
				name = r.name
				resourceFileName = r.fileName
			default:
				return fmt.Errorf("unsupported resource type: %T", resource)
			}

			if existing, exists := resourceMap[name]; exists {
				var existingFileName string
				switch e := existing.(type) {
				case machineInfo:
					existingFileName = e.fileName
				case variableInfo:
					existingFileName = e.fileName
				case actionInfo:
					existingFileName = e.fileName
				}

				return fmt.Errorf("%s name collision: '%s' defined in both %s and %s",
					contentType, name, existingFileName, resourceFileName)
			}
			resourceMap[name] = resource
		}
	}

	// Now validate the merged content against the schema
	var mergedContent string
	switch contentType {
	case "machines":
		mergedContent = mergeHCLContent(contents, "machine", "machines")
	case "variables":
		mergedContent = mergeHCLContent(contents, "variable", "variables")
	case "actions":
		mergedContent = mergeHCLContent(contents, "action", "actions")
	}

	if err := validateAgainstSchema(contentType, mergedContent); err != nil {
		return fmt.Errorf("merged %s schema validation failed: %w", contentType, err)
	}

	return nil
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

// mergeHCLContent merges multiple HCL files into a single configuration
// This unified function replaces the duplicate mergeActionsContent, mergeVariablesContent, and mergeMachinesContent functions
func mergeHCLContent(contents []string, blockType, blockName string) string {
	merged := fmt.Sprintf(`# Merged %s Configuration
# Generated from multiple files

metadata {
  version = "1"
  description = "Merged %s configuration"
}

%s {
`, strings.Title(blockType), blockType, blockName)

	for _, content := range contents {
		blocks := extractHCLBlocks(content, blockType)
		for _, block := range blocks {
			cleaned := strings.TrimSpace(block)
			merged += "\n  " + cleaned + "\n"
		}
	}
	merged += "}\n"
	return merged
}

// extractHCLBlocks extracts HCL blocks of the specified type using proper HCL parsing
// This unified function replaces the duplicate extractMachineBlocks, extractActionBlocks, and extractVariableBlocks functions
func extractHCLBlocks(content, blockType string) []string {
	var blocks []string
	file, diags := hclsyntax.ParseConfig([]byte(content), "content.hcl", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return blocks
	}
	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: blockType, LabelNames: []string{"name"}},
		},
	}
	bodyContent, diags := file.Body.Content(schema)
	if diags.HasErrors() {
		return blocks
	}
	for _, block := range bodyContent.Blocks {
		startPos := block.DefRange.Start
		endPos := block.DefRange.End
		if startPos.Byte < len(content) && endPos.Byte <= len(content) {
			blockContent := content[startPos.Byte:endPos.Byte]
			blocks = append(blocks, blockContent)
		}
	}
	return blocks
}

// getKeys returns the keys from a map as a slice of strings
func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// getResourceType returns the resource type name for a given content type
func getResourceType(contentType string) string {
	switch contentType {
	case "machines":
		return "machine"
	case "variables":
		return "variable"
	case "actions":
		return "action"
	default:
		return contentType
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
	return validateMergedContent(allContents, fileNames, "machines")
}

type machineInfo struct {
	name     string
	fileName string
	hostname string
	user     string
}

type variableInfo struct {
	name        string
	fileName    string
	description string
	sensitive   bool
	encrypted   bool
}

type actionInfo struct {
	name        string
	fileName    string
	description string
	actionType  string
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

// extractResourceNamesFromParsedData extracts resource information from parsed HCL data
// This unified function replaces the duplicate extractMachineNamesFromParsedData, extractVariableNamesFromParsedData, and extractActionNamesFromParsedData functions
func extractResourceNamesFromParsedData(data map[string]interface{}, fileName, blockType, resourceType string) ([]interface{}, error) {
	var resources []interface{}

	// Look for the main block (e.g., "machines", "variables", "actions")
	mainBlock, exists := data[blockType]
	if !exists {
		return resources, nil // No main block, return empty
	}

	// Convert to map
	mainMap, ok := mainBlock.(map[string]interface{})
	if !ok {
		return resources, nil // Invalid structure, return empty
	}

	// Handle different resource structures based on type
	switch blockType {
	case "machines":
		return extractMachineResources(mainMap, fileName, resourceType)
	case "variables":
		return extractVariableResources(mainMap, fileName, resourceType)
	case "actions":
		return extractActionResources(mainMap, fileName, resourceType)
	default:
		return nil, fmt.Errorf("unsupported block type: %s", blockType)
	}
}

// extractMachineResources extracts machine resources from the main block
func extractMachineResources(mainMap map[string]interface{}, fileName, resourceType string) ([]interface{}, error) {
	var resources []interface{}

	// Look for machine array
	machineValue, exists := mainMap[resourceType]
	if !exists {
		return resources, nil // No machines defined
	}

	// Handle machine blocks array
	if machineBlocks, ok := machineValue.([]map[string]interface{}); ok {
		// Multiple machine blocks (correct type assertion)
		for _, machineBlock := range machineBlocks {
			if name, exists := machineBlock["name"]; exists {
				nameStr := convertToString(name)
				// Filter out placeholder values that indicate parsing errors
				if nameStr != "" && nameStr != "complex_expression" && !strings.HasPrefix(nameStr, "complex_") {
					resources = append(resources, machineInfo{
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
						resources = append(resources, machineInfo{
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
				resources = append(resources, machineInfo{
					name:     nameStr,
					fileName: fileName,
				})
			}
		}
	}

	return resources, nil
}

// extractVariableResources extracts variable resources from the main block
func extractVariableResources(mainMap map[string]interface{}, fileName, resourceType string) ([]interface{}, error) {
	var resources []interface{}

	// Look for variable blocks (variable { ... })
	variableBlocks, exists := mainMap[resourceType]
	if !exists {
		return resources, nil
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

					resources = append(resources, variable)
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

				resources = append(resources, variable)
			}
		}
	}

	return resources, nil
}

// extractActionResources extracts action resources from the main block
func extractActionResources(mainMap map[string]interface{}, fileName, resourceType string) ([]interface{}, error) {
	var resources []interface{}
	logger := logging.GetGlobalLogger()

	// Look for action blocks (action "name" { ... })
	for actionName, actionValue := range mainMap {
		logger.Debug("found action",
			slog.String("file", fileName),
			slog.String("action_name", actionName))

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

			resources = append(resources, action)
		} else {
			// Simple action block without metadata
			action := actionInfo{
				name:     actionName,
				fileName: fileName,
			}
			resources = append(resources, action)
		}
	}

	return resources, nil
}

// extractActionNamesFromParsedData extracts action information from parsed HCL data
func extractActionNamesFromParsedData(data map[string]interface{}, fileName string) ([]actionInfo, error) {
	var actions []actionInfo
	logger := logging.GetGlobalLogger()

	// Look for actions block
	actionsBlock, exists := data["actions"]
	if !exists {
		logger.Debug("no actions block found",
			slog.String("file", fileName))
		return actions, nil // No actions block, return empty
	}

	logger.Debug("found actions block",
		slog.String("file", fileName),
		slog.String("type", fmt.Sprintf("%T", actionsBlock)))

	// Convert to map
	actionsMap, ok := actionsBlock.(map[string]interface{})
	if !ok {
		logger.Debug("actions block is not a map",
			slog.String("file", fileName))
		return actions, nil // Invalid structure, return empty
	}

	logger.Debug("found actions map keys",
		slog.String("file", fileName),
		slog.Any("keys", getKeys(actionsMap)))

	// Look for action blocks (action "name" { ... })
	for actionName, actionValue := range actionsMap {
		logger.Debug("found action",
			slog.String("file", fileName),
			slog.String("action_name", actionName))

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

// getProjectName extracts the project name from project.hcl file
func getProjectName(targetDir string) (string, error) {
	projectHCLPath := filepath.Join(targetDir, "project.hcl")
	content, err := os.ReadFile(projectHCLPath)
	if err != nil {
		return "", nil // Return empty string instead of error for missing files
	}

	// Regex to extract project name from "project \"name\" {"
	re := regexp.MustCompile(`project\s+["']([^"']+)["']\s*{`)
	matches := re.FindStringSubmatch(string(content))
	if len(matches) < 2 {
		return "unknown", nil // Return default name instead of error
	}

	return matches[1], nil
}

func runProjectEncrypt(cmd *cobra.Command, args []string) error {
	// Determine target directory
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	// Validate project directory
	validator := utilities.NewProjectValidator()
	result := validator.ValidateProject(targetDir)
	if !result.IsValid {
		return fmt.Errorf("project directory validation failed: %s", result.Errors[0].Message)
	}

	// Get age configuration from spooky config
	configManager, err := utilities.NewConfigManager()
	if err != nil {
		return fmt.Errorf("failed to create config manager: %w", err)
	}

	// Get effective config to find age settings (respect --config flag)
	customConfigFile := GetConfigFile()
	_, err = configManager.GetEffectiveConfig(customConfigFile)
	if err != nil {
		return fmt.Errorf("failed to get effective config: %w", err)
	}

	// Parse config to find age settings
	// For now, we'll use default paths - in a real implementation,
	// you'd parse the HCL config to get the actual paths
	identitiesPath := "~/.config/spooky/age/identities"
	recipientsPath := "~/.config/spooky/age/recipients"

	// Expand home directory
	if strings.HasPrefix(identitiesPath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		identitiesPath = filepath.Join(homeDir, identitiesPath[2:])
	}

	if strings.HasPrefix(recipientsPath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		recipientsPath = filepath.Join(homeDir, recipientsPath[2:])
	}

	// Create age encryption instance
	ageEncryption, err := encryption.NewAgeEncryption(identitiesPath, recipientsPath)
	if err != nil {
		return fmt.Errorf("failed to create age encryption: %w", err)
	}

	// Validate age configuration
	if err := ageEncryption.ValidateConfiguration(); err != nil {
		return fmt.Errorf("age encryption configuration invalid: %w", err)
	}

	logger := logging.GetGlobalLogger()
	logger.Info("Age encryption configured",
		slog.Int("identities_count", ageEncryption.GetIdentitiesCount()),
		slog.Int("recipients_count", ageEncryption.GetRecipientsCount()))

	// Create HCL updater
	updater := encryption.NewHCLUpdater(ageEncryption)

	// Process the project directory
	logger.Info("Processing project directory", slog.String("target_directory", targetDir))
	if err := updater.UpdateDirectory(targetDir); err != nil {
		return fmt.Errorf("failed to process project directory: %w", err)
	}

	logger.Info("Project encryption completed successfully!")
	return nil
}

func runProjectConfig(cmd *cobra.Command, args []string) error {
	// Get custom config file from --config flag
	customConfigFile := GetConfigFile()

	logger := logging.GetGlobalLogger()
	logger.Info("=== Spooky Configuration Information ===")

	// Create config manager
	configManager, err := utilities.NewConfigManager()
	if err != nil {
		return fmt.Errorf("failed to create config manager: %w", err)
	}

	// Get effective config info
	effectiveInfo, err := configManager.GetEffectiveConfigInfo(customConfigFile)
	if err != nil {
		return fmt.Errorf("failed to get effective config info: %w", err)
	}

	// Display configuration information
	logger.Info("Configuration information",
		slog.String("source", effectiveInfo.Source),
		slog.String("config_file", effectiveInfo.ConfigFile))

	if effectiveInfo.Exists {
		logger.Info("Configuration file details",
			slog.Int64("size_bytes", effectiveInfo.Size))
		if effectiveInfo.ModTime != "" {
			logger.Info("Configuration file modified", slog.String("modified_time", effectiveInfo.ModTime))
		}
	} else {
		logger.Info("Configuration file status", slog.String("status", "Does not exist"))
	}

	// Show config priority
	logger.Info("Configuration priority information",
		slog.String("priority1", "Custom config file (--config flag)"),
		slog.String("priority2", "User config (~/.config/spooky/spooky.hcl)"),
		slog.String("priority3", "Embedded default"))

	// Show current --config flag status
	logger.Info("Configuration flag status")
	if customConfigFile != "" {
		logger.Info("Custom config file specified", slog.String("config_file", customConfigFile))
	} else {
		logger.Info("Custom config file not specified")
	}

	return nil
}
