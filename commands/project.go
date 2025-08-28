package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"spooky/internal/encryption"
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
	// Generate project configuration directly from Go structs
	content := schemas.GenerateProjectConfigFromStructs(name, description)
	return os.WriteFile(filepath.Join(targetDir, "project.hcl"), []byte(content), 0o644)
}

func createMachinesHCL(targetDir string) error {
	// Generate machines configuration with the new simplified format
	content := `machines {
  # Example machine using IP address as identifier
  machine "192.168.1.100" {
    port = 22
    user = "root"
    authentication "password" {
      password {
        value = "your_password_here"
        encrypted = false
      }
    }
  }

  # Example machine using FQDN as identifier
  # machine "web-server.example.com" {
  #   port = 22
  #   user = "admin"
  #   authentication "publickey" {
  #     public_key_path = "~/.ssh/id_rsa"
  #   }
  # }

  # Example machine using IPv6 address as identifier
  # machine "2001:db8::1" {
  #   port = 22
  #   user = "root"
  #   authentication "password" {
  #     password {
  #       value = "your_password_here"
  #       encrypted = false
  #     }
  #   }
  # }
}`
	return os.WriteFile(filepath.Join(targetDir, "machines.hcl"), []byte(content), 0o644)
}

func createActionsHCL(targetDir string) error {
	// Generate actions configuration directly from Go structs
	content := schemas.GenerateActionsConfigFromStructs()
	return os.WriteFile(filepath.Join(targetDir, "actions.hcl"), []byte(content), 0o644)
}

func createVariablesHCL(targetDir string) error {
	// Generate variables configuration directly from Go structs
	content := schemas.GenerateVariablesConfigFromStructs()

	return os.WriteFile(filepath.Join(targetDir, "variables.hcl"), []byte(content), 0o644)
}

// generateProjectConfigFromSchema creates a project configuration based on schema understanding
func generateProjectConfigFromSchema(name, description string, schemaData map[string]interface{}) string {
	var content strings.Builder

	// Generate header from schema metadata if available
	if metadata, ok := schemaData["metadata"].(map[string]interface{}); ok {
		if version, hasVersion := metadata["version"]; hasVersion {
			content.WriteString(fmt.Sprintf("# Spooky Project Configuration Schema v%v\n", version))
		}
		if desc, hasDesc := metadata["description"]; hasDesc {
			content.WriteString(fmt.Sprintf("# %v\n", desc))
		}
		content.WriteString("# Generated based on embedded schema\n\n")
	} else {
		content.WriteString("# Spooky Project Configuration\n# Generated based on embedded schema\n\n")
	}

	// Generate configuration blocks entirely from schema structure
	for blockName, blockData := range schemaData {
		if blockName == "metadata" {
			continue // Skip metadata block, already handled above
		}

		// Generate block header
		content.WriteString(fmt.Sprintf("%s {\n", blockName))

		// Handle special case for project block - add required name and description
		if blockName == "project" {
			content.WriteString(fmt.Sprintf("  name = \"%s\"\n", name))
			content.WriteString(fmt.Sprintf("  description = \"%s\"\n", description))
		}

		// Generate fields from schema understanding
		if blockInfo, ok := blockData.(map[string]interface{}); ok {
			for fieldName, fieldData := range blockInfo {
				// Skip name and description if we already added them for project block
				if blockName == "project" && (fieldName == "name" || fieldName == "description") {
					continue
				}

				if fieldInfo, ok := fieldData.(map[string]interface{}); ok {
					// Check if field is required
					if required, hasRequired := fieldInfo["required"]; hasRequired && required == true {
						// For required fields, add placeholder with description
						if desc, hasDesc := fieldInfo["description"]; hasDesc {
							content.WriteString(fmt.Sprintf("  # %s = <required>  # %v\n", fieldName, desc))
						} else {
							content.WriteString(fmt.Sprintf("  # %s = <required>\n", fieldName))
						}
					} else {
						// For optional fields, add commented default or placeholder
						if defaultValue, hasDefault := fieldInfo["default"]; hasDefault {
							if desc, hasDesc := fieldInfo["description"]; hasDesc {
								content.WriteString(fmt.Sprintf("  # %s = %v  # %v\n", fieldName, defaultValue, desc))
							} else {
								content.WriteString(fmt.Sprintf("  # %s = %v\n", fieldName, defaultValue))
							}
						} else {
							if desc, hasDesc := fieldInfo["description"]; hasDesc {
								content.WriteString(fmt.Sprintf("  # %s = <value>  # %v\n", fieldName, desc))
							} else {
								content.WriteString(fmt.Sprintf("  # %s = <value>\n", fieldName))
							}
						}
					}
				}
			}
		}

		content.WriteString("}\n\n")
	}

	return strings.TrimSpace(content.String())
}

// generateVariablesConfigFromSchema creates a variables configuration based on schema understanding
func generateVariablesConfigFromSchema(schemaData map[string]interface{}) string {
	var content strings.Builder

	// Generate header from schema metadata if available
	if metadata, ok := schemaData["metadata"].(map[string]interface{}); ok {
		if version, hasVersion := metadata["version"]; hasVersion {
			content.WriteString(fmt.Sprintf("# Variables Configuration Schema v%v\n", version))
		}
		if desc, hasDesc := metadata["description"]; hasDesc {
			content.WriteString(fmt.Sprintf("# %v\n", desc))
		}
		content.WriteString("# Generated based on embedded schema\n\n")
	} else {
		content.WriteString("# Variables Configuration\n# Generated based on embedded schema\n\n")
	}

	// Generate configuration blocks entirely from schema structure
	for blockName, blockData := range schemaData {
		if blockName == "metadata" {
			continue // Skip metadata block, already handled above
		}

		// Generate block header
		content.WriteString(fmt.Sprintf("%s \"example-%s\" {\n", blockName, blockName))

		// Generate fields from schema understanding
		if blockInfo, ok := blockData.(map[string]interface{}); ok {
			for fieldName, fieldData := range blockInfo {
				if fieldInfo, ok := fieldData.(map[string]interface{}); ok {
					// Check if field is required
					if required, hasRequired := fieldInfo["required"]; hasRequired && required == true {
						// For required fields, add placeholder with description
						if desc, hasDesc := fieldInfo["description"]; hasDesc {
							content.WriteString(fmt.Sprintf("  # %s = <required>  # %v\n", fieldName, desc))
						} else {
							content.WriteString(fmt.Sprintf("  # %s = <required>\n", fieldName))
						}
					} else {
						// For optional fields, add commented default or placeholder
						if defaultValue, hasDefault := fieldInfo["default"]; hasDefault {
							if desc, hasDesc := fieldInfo["description"]; hasDesc {
								content.WriteString(fmt.Sprintf("  # %s = %v  # %v\n", fieldName, defaultValue, desc))
							} else {
								content.WriteString(fmt.Sprintf("  # %s = %v\n", fieldName, defaultValue))
							}
						} else {
							if desc, hasDesc := fieldInfo["description"]; hasDesc {
								content.WriteString(fmt.Sprintf("  # %s = <value>  # %v\n", fieldName, desc))
							} else {
								content.WriteString(fmt.Sprintf("  # %s = <value>\n", fieldName))
							}
						}
					}
				}
			}
		}

		content.WriteString("}\n\n")
	}

	// Add usage guidance
	content.WriteString("# Usage:\n")
	content.WriteString("# 1. Replace 'example-*' with your actual variable names\n")
	content.WriteString("# 2. Fill in required fields marked with <required>\n")
	content.WriteString("# 3. Uncomment and configure optional fields as needed\n")
	content.WriteString("# 4. Add more variables by copying and modifying the examples\n")

	return strings.TrimSpace(content.String())
}

// generateMachinesConfigFromSchema creates a machines configuration based on schema understanding
func generateMachinesConfigFromSchema(schemaData map[string]interface{}) string {
	var content strings.Builder

	// Generate header from schema metadata if available
	if metadata, ok := schemaData["metadata"].(map[string]interface{}); ok {
		if version, hasVersion := metadata["version"]; hasVersion {
			content.WriteString(fmt.Sprintf("# Machines Configuration Schema v%v\n", version))
		}
		if desc, hasDesc := metadata["description"]; hasDesc {
			content.WriteString(fmt.Sprintf("# %v\n", desc))
		}
		content.WriteString("# Generated based on embedded schema\n\n")
	} else {
		content.WriteString("# Machines Configuration\n# Generated based on embedded schema\n\n")
	}

	// Generate configuration blocks entirely from schema structure
	for blockName, blockData := range schemaData {
		if blockName == "metadata" {
			continue // Skip metadata block, already handled above
		}

		// Generate block header
		content.WriteString(fmt.Sprintf("%s \"example-%s\" {\n", blockName, blockName))

		// Generate fields from schema understanding
		if blockInfo, ok := blockData.(map[string]interface{}); ok {
			for fieldName, fieldData := range blockInfo {
				if fieldInfo, ok := fieldData.(map[string]interface{}); ok {
					// Check if field is required
					if required, hasRequired := fieldInfo["required"]; hasRequired && required == true {
						// For required fields, add placeholder with description
						if desc, hasDesc := fieldInfo["description"]; hasDesc {
							content.WriteString(fmt.Sprintf("  # %s = <required>  # %v\n", fieldName, desc))
						} else {
							content.WriteString(fmt.Sprintf("  # %s = <required>\n", fieldName))
						}
					} else {
						// For optional fields, add commented default or placeholder
						if defaultValue, hasDefault := fieldInfo["default"]; hasDefault {
							if desc, hasDesc := fieldInfo["description"]; hasDesc {
								content.WriteString(fmt.Sprintf("  # %s = %v  # %v\n", fieldName, defaultValue, desc))
							} else {
								content.WriteString(fmt.Sprintf("  # %s = %v\n", fieldName, defaultValue))
							}
						} else {
							if desc, hasDesc := fieldInfo["description"]; hasDesc {
								content.WriteString(fmt.Sprintf("  # %s = <value>  # %v\n", fieldName, desc))
							} else {
								content.WriteString(fmt.Sprintf("  # %s = <value>\n", fieldName))
							}
						}
					}
				}
			}
		}

		content.WriteString("}\n\n")
	}

	// Add usage guidance
	content.WriteString("# Usage:\n")
	content.WriteString("# 1. Replace 'example-*' with your actual machine/group names\n")
	content.WriteString("# 2. Fill in required fields marked with <required>\n")
	content.WriteString("# 3. Uncomment and configure optional fields as needed\n")
	content.WriteString("# 4. Add more machines/groups by copying and modifying the examples\n")

	return strings.TrimSpace(content.String())
}

// generateActionsConfigFromSchema creates an actions configuration based on schema understanding
func generateActionsConfigFromSchema(schemaData map[string]interface{}) string {
	var content strings.Builder

	// Generate header from schema metadata if available
	if metadata, ok := schemaData["metadata"].(map[string]interface{}); ok {
		if version, hasVersion := metadata["version"]; hasVersion {
			content.WriteString(fmt.Sprintf("# Actions Configuration Schema v%v\n", version))
		}
		if desc, hasDesc := metadata["description"]; hasDesc {
			content.WriteString(fmt.Sprintf("# %v\n", desc))
		}
		content.WriteString("# Generated based on embedded schema\n\n")
	} else {
		content.WriteString("# Actions Configuration\n# Generated based on embedded schema\n\n")
	}

	// Generate configuration blocks entirely from schema structure
	for blockName, blockData := range schemaData {
		if blockName == "metadata" {
			continue // Skip metadata block, already handled above
		}

		// Generate block header
		content.WriteString(fmt.Sprintf("%s \"example-%s\" {\n", blockName, blockName))

		// Generate fields from schema understanding
		if blockInfo, ok := blockData.(map[string]interface{}); ok {
			for fieldName, fieldData := range blockInfo {
				if fieldInfo, ok := fieldData.(map[string]interface{}); ok {
					// Check if field is required
					if required, hasRequired := fieldInfo["required"]; hasRequired && required == true {
						// For required fields, add placeholder with description
						if desc, hasDesc := fieldInfo["description"]; hasDesc {
							content.WriteString(fmt.Sprintf("  # %s = <required>  # %v\n", fieldName, desc))
						} else {
							content.WriteString(fmt.Sprintf("  # %s = <required>\n", fieldName))
						}
					} else {
						// For optional fields, add commented default or placeholder
						if defaultValue, hasDefault := fieldInfo["default"]; hasDefault {
							if desc, hasDesc := fieldInfo["description"]; hasDesc {
								content.WriteString(fmt.Sprintf("  # %s = %v  # %v\n", fieldName, defaultValue, desc))
							} else {
								content.WriteString(fmt.Sprintf("  # %s = %v\n", fieldName, defaultValue))
							}
						} else {
							if desc, hasDesc := fieldInfo["description"]; hasDesc {
								content.WriteString(fmt.Sprintf("  # %s = <value>  # %v\n", fieldName, desc))
							} else {
								content.WriteString(fmt.Sprintf("  # %s = <value>\n", fieldName))
							}
						}
					}
				}
			}
		}

		content.WriteString("}\n\n")
	}

	// Add usage guidance
	content.WriteString("# Usage:\n")
	content.WriteString("# 1. Replace 'example-*' with your actual action names\n")
	content.WriteString("# 2. Fill in required fields marked with <required>\n")
	content.WriteString("# 3. Uncomment and configure optional fields as needed\n")
	content.WriteString("# 4. Add more actions by copying and modifying the examples\n")

	return strings.TrimSpace(content.String())
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
		fmt.Printf("Project \"%s\" is valid\n", projectName)
	} else {
		fmt.Printf("Project validation failed with %d errors:\n", len(result.Errors))
		for _, err := range result.Errors {
			fmt.Printf("  - %s\n", err.Message)
		}
		return fmt.Errorf("project validation failed with %d errors", len(result.Errors))
	}

	return nil
}

// validateAgainstSchema validates HCL content against embedded schemas
func validateAgainstSchema(schemaName, content string) error {
	// Use the new unified validator for proper schema validation
	validator := schemas.NewValidator()

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
		validator := schemas.NewValidator()

		// Parse the content using the unified validator
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
		validator := schemas.NewValidator()

		// Parse the content using the unified validator
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
		validator := schemas.NewValidator()

		// Parse the content using the unified validator
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

// extractActionBlocks extracts action blocks using proper HCL parsing
func extractActionBlocks(content string) []string {
	var blocks []string

	// Parse HCL content
	file, diags := hclsyntax.ParseConfig([]byte(content), "content.hcl", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		// If parsing fails, return empty slice
		return blocks
	}

	// Define schema for action blocks
	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "action",
				LabelNames: []string{"name"},
			},
		},
	}

	// Extract action blocks
	bodyContent, diags := file.Body.Content(schema)
	if diags.HasErrors() {
		return blocks
	}

	// Convert blocks back to HCL strings by extracting the original content
	for _, block := range bodyContent.Blocks {
		// Get the range of the block in the original content
		startPos := block.DefRange.Start
		endPos := block.DefRange.End

		// Extract the block content from the original string
		if startPos.Byte < len(content) && endPos.Byte <= len(content) {
			blockContent := content[startPos.Byte:endPos.Byte]
			blocks = append(blocks, blockContent)
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

// extractVariableBlocks extracts variable blocks using proper HCL parsing
func extractVariableBlocks(content string) []string {
	var blocks []string

	// Parse HCL content
	file, diags := hclsyntax.ParseConfig([]byte(content), "content.hcl", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		// If parsing fails, return empty slice
		return blocks
	}

	// Define schema for variable blocks
	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "variable",
				LabelNames: []string{},
			},
		},
	}

	// Extract variable blocks
	bodyContent, diags := file.Body.Content(schema)
	if diags.HasErrors() {
		return blocks
	}

	// Convert blocks back to HCL strings by extracting the original content
	for _, block := range bodyContent.Blocks {
		// Get the range of the block in the original content
		startPos := block.DefRange.Start
		endPos := block.DefRange.End

		// Extract the block content from the original string
		if startPos.Byte < len(content) && endPos.Byte <= len(content) {
			blockContent := content[startPos.Byte:endPos.Byte]
			blocks = append(blocks, blockContent)
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

// extractMachineBlocks extracts machine blocks using proper HCL parsing
func extractMachineBlocks(content string) []string {
	var blocks []string

	// Parse HCL content
	file, diags := hclsyntax.ParseConfig([]byte(content), "content.hcl", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		// If parsing fails, return empty slice
		return blocks
	}

	// Define schema for machine blocks
	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "machine",
				LabelNames: []string{},
			},
		},
	}

	// Extract machine blocks
	bodyContent, diags := file.Body.Content(schema)
	if diags.HasErrors() {
		return blocks
	}

	// Convert blocks back to HCL strings by extracting the original content
	for _, block := range bodyContent.Blocks {
		// Get the range of the block in the original content
		startPos := block.DefRange.Start
		endPos := block.DefRange.End

		// Extract the block content from the original string
		if startPos.Byte < len(content) && endPos.Byte <= len(content) {
			blockContent := content[startPos.Byte:endPos.Byte]
			blocks = append(blocks, blockContent)
		}
	}

	return blocks
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

	fmt.Printf("Age encryption configured with %d identities and %d recipients\n",
		ageEncryption.GetIdentitiesCount(), ageEncryption.GetRecipientsCount())

	// Create HCL updater
	updater := encryption.NewSimpleHCLUpdater(ageEncryption)

	// Process the project directory
	fmt.Printf("Processing project directory: %s\n", targetDir)
	if err := updater.UpdateDirectory(targetDir); err != nil {
		return fmt.Errorf("failed to process project directory: %w", err)
	}

	fmt.Println("Project encryption completed successfully!")
	return nil
}

func runProjectConfig(cmd *cobra.Command, args []string) error {
	// Get custom config file from --config flag
	customConfigFile := GetConfigFile()

	fmt.Println("=== Spooky Configuration Information ===")
	fmt.Println()

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
	fmt.Printf("Configuration Source: %s\n", effectiveInfo.Source)
	fmt.Printf("Configuration File: %s\n", effectiveInfo.ConfigFile)

	if effectiveInfo.Exists {
		fmt.Printf("File Size: %d bytes\n", effectiveInfo.Size)
		if effectiveInfo.ModTime != "" {
			fmt.Printf("Last Modified: %s\n", effectiveInfo.ModTime)
		}
	} else {
		fmt.Println("File Status: Does not exist")
	}

	// Show config priority
	fmt.Println()
	fmt.Println("Configuration Priority:")
	fmt.Println("  1. Custom config file (--config flag)")
	fmt.Println("  2. User config (~/.config/spooky/spooky.hcl)")
	fmt.Println("  3. Embedded default")

	// Show current --config flag status
	fmt.Println()
	if customConfigFile != "" {
		fmt.Printf("--config flag: %s\n", customConfigFile)
	} else {
		fmt.Println("--config flag: Not specified")
	}

	return nil
}
