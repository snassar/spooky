// Package config provides configuration loading and management for spooky.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	spookyschemas "spooky/internal/schemas"
	spookytypesschemas "spooky/internal/types/schemas"
)

// AutoSetupConfig ensures spooky configuration directory exists and is properly configured
// This function is called before any spooky command run (except --version and --help)
func AutoSetupConfig() error {
	// Determine OS and get appropriate config directory
	configDir, err := getConfigDirectory()
	if err != nil {
		return fmt.Errorf("failed to determine config directory: %w", err)
	}

	// Check if config directory exists
	configExists, err := configDirectoryExists(configDir)
	if err != nil {
		return fmt.Errorf("failed to check config directory: %w", err)
	}

	if !configExists {
		// Create config directory and default files
		if err := createConfigDirectory(configDir); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
	} else {
		// Check if required config files exist, create them if missing
		if err := ensureConfigFiles(configDir); err != nil {
			return fmt.Errorf("config file setup failed: %w", err)
		}

		// Validate existing config files
		if err := validateConfigFiles(configDir); err != nil {
			return fmt.Errorf("config validation failed: %w", err)
		}
	}

	return nil
}

// ensureConfigFiles ensures that required config files exist, creating them if missing
func ensureConfigFiles(configDir string) error {
	// Check if spooky.hcl exists
	spookyConfigPath := filepath.Join(configDir, "spooky.hcl")
	if _, err := os.Stat(spookyConfigPath); os.IsNotExist(err) {
		if err := createDefaultSpookyConfig(configDir); err != nil {
			return fmt.Errorf("failed to create default spooky.hcl: %w", err)
		}
	}

	return nil
}

// getConfigDirectory returns the appropriate config directory for the current OS
func getConfigDirectory() (string, error) {
	switch runtime.GOOS {
	case "linux", "freebsd", "openbsd", "netbsd", "dragonfly":
		// Use XDG Base Directory Specification
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigHome == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("failed to get user home directory: %w", err)
			}
			xdgConfigHome = filepath.Join(homeDir, ".config")
		}
		return filepath.Join(xdgConfigHome, "spooky"), nil

	case "darwin":
		// macOS: ~/Library/Application Support/spooky
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		return filepath.Join(homeDir, "Library", "Application Support", "spooky"), nil

	case "windows":
		// Windows: %APPDATA%\spooky
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("APPDATA environment variable not set")
		}
		return filepath.Join(appData, "spooky"), nil

	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// configDirectoryExists checks if the spooky config directory exists
func configDirectoryExists(configDir string) (bool, error) {
	info, err := os.Stat(configDir)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

// createConfigDirectory creates the spooky config directory and default configuration files
func createConfigDirectory(configDir string) error {
	// Create the config directory
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", configDir, err)
	}

	// Create default spooky.hcl
	if err := createDefaultSpookyConfig(configDir); err != nil {
		return fmt.Errorf("failed to create default spooky.hcl: %w", err)
	}

	return nil
}

// createDefaultSpookyConfig creates a comprehensive spooky.hcl file using schema-driven generation
func createDefaultSpookyConfig(configDir string) error {
	// Load the spooky schema
	schemaManager := spookyschemas.NewManager(nil) // No logger needed for this operation
	spookySchema, err := schemaManager.Load("internal/schemas/schemas/structure/spooky.hcl")
	if err != nil {
		return fmt.Errorf("failed to load spooky schema: %w", err)
	}

	// Generate HCL content from schema
	hclContent, err := generateSpookyHCLFromSchema(spookySchema)
	if err != nil {
		return fmt.Errorf("failed to generate spooky.hcl from schema: %w", err)
	}

	// Write the generated content to file
	configPath := filepath.Join(configDir, "spooky.hcl")
	if err := os.WriteFile(configPath, []byte(hclContent), 0o600); err != nil {
		return fmt.Errorf("failed to write spooky.hcl file: %w", err)
	}

	return nil
}

// generateSpookyHCLFromSchema generates HCL content based on the spooky schema
func generateSpookyHCLFromSchema(schema *spookytypesschemas.Schema) (string, error) {
	var content strings.Builder

	// Add header comment
	content.WriteString("# Spooky CLI Configuration\n")
	content.WriteString("# This file contains global configuration for the spooky CLI tool\n")
	content.WriteString("# Generated from schema: " + schema.Name + "\n\n")

	// Parse the schema content to understand the structure
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL([]byte(schema.Content), schema.Name)
	if diags.HasErrors() {
		return "", fmt.Errorf("failed to parse schema content: %s", diags.Error())
	}

	// Start spooky block
	content.WriteString("spooky {\n")

	// Parse the spooky block from the schema
	bodyContent, diags := file.Body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "spooky"},
		},
	})
	if diags.HasErrors() {
		return "", fmt.Errorf("failed to parse spooky block from schema: %s", diags.Error())
	}

	if len(bodyContent.Blocks) > 0 {
		spookyBlock := bodyContent.Blocks[0]
		if err := generateBlockFromSchema(spookyBlock, &content, 2); err != nil {
			return "", fmt.Errorf("failed to generate spooky block: %w", err)
		}
	}

	// Close spooky block
	content.WriteString("}\n")

	return content.String(), nil
}

// generateBlockFromSchema generates HCL content for a block based on its schema definition
func generateBlockFromSchema(block *hcl.Block, content *strings.Builder, indent int) error {
	indentStr := strings.Repeat("  ", indent)

	// Parse the block content to understand its structure
	blockContent, diags := block.Body.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "*", Required: false}, // Allow any attribute
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "*"}, // Allow any blocks
		},
	})
	if diags.HasErrors() {
		return fmt.Errorf("failed to parse block content: %s", diags.Error())
	}

	// Generate attributes with their default values
	for _, attr := range blockContent.Attributes {
		if err := generateAttributeFromSchema(attr, content, indentStr); err != nil {
			return fmt.Errorf("failed to generate attribute %s: %w", attr.Name, err)
		}
	}

	// Generate nested blocks
	for _, nestedBlock := range blockContent.Blocks {
		if err := generateNestedBlockFromSchema(nestedBlock, content, indent); err != nil {
			return fmt.Errorf("failed to generate nested block %s: %w", nestedBlock.Type, err)
		}
	}

	return nil
}

// generateAttributeFromSchema generates HCL content for an attribute based on its schema definition
func generateAttributeFromSchema(attr *hcl.Attribute, content *strings.Builder, indentStr string) error {
	// Try to parse the attribute as a field validation block
	if block, ok := attr.Expr.(*hclsyntax.ObjectConsExpr); ok {
		return generateFieldValidationFromBlock(attr.Name, block, content, indentStr)
	}

	// Try to parse as a simple value
	if val, diags := attr.Expr.Value(nil); !diags.HasErrors() {
		// Extract default value and type information
		defaultValue := extractDefaultValue(val)
		content.WriteString(fmt.Sprintf("%s%s = %s\n", indentStr, attr.Name, defaultValue))
		return nil
	}

	// If we can't parse it, just write the attribute name with a placeholder
	content.WriteString(fmt.Sprintf("%s%s = \"<value>\"\n", indentStr, attr.Name))
	return nil
}

// generateFieldValidationFromBlock generates HCL content from a field validation block
func generateFieldValidationFromBlock(fieldName string, block *hclsyntax.ObjectConsExpr, content *strings.Builder, indentStr string) error {
	// Extract type, default value, and description from the validation block
	var fieldType, defaultValue, description string

	for _, item := range block.Items {
		if key, ok := item.KeyExpr.(*hclsyntax.LiteralValueExpr); ok {
			keyStr := key.Val.AsString()

			if val, diags := item.ValueExpr.Value(nil); !diags.HasErrors() {
				switch keyStr {
				case "type":
					fieldType = val.AsString()
				case "default":
					defaultValue = formatValue(val)
				case "description":
					description = val.AsString()
				}
			}
		}
	}

	// Add comment with description if available
	if description != "" {
		content.WriteString(fmt.Sprintf("%s# %s\n", indentStr, description))
	}

	// Generate the field with appropriate default value
	if defaultValue != "" {
		content.WriteString(fmt.Sprintf("%s%s = %s\n", indentStr, fieldName, defaultValue))
	} else {
		// Generate appropriate default based on type
		defaultValue = generateDefaultForType(fieldType)
		content.WriteString(fmt.Sprintf("%s%s = %s\n", indentStr, fieldName, defaultValue))
	}

	return nil
}

// generateNestedBlockFromSchema generates HCL content for a nested block
func generateNestedBlockFromSchema(block *hcl.Block, content *strings.Builder, parentIndent int) error {
	indentStr := strings.Repeat("  ", parentIndent)

	// Add comment for the block
	content.WriteString(fmt.Sprintf("%s# %s configuration\n", indentStr, cases.Title(language.English).String(block.Type)))
	content.WriteString(fmt.Sprintf("%s%s {\n", indentStr, block.Type))

	// Generate the block content
	if err := generateBlockFromSchema(block, content, parentIndent+1); err != nil {
		return err
	}

	content.WriteString(fmt.Sprintf("%s}\n\n", indentStr))
	return nil
}

// extractDefaultValue extracts a default value from a cty.Value
func extractDefaultValue(val cty.Value) string {
	if val.IsNull() {
		return "null"
	}

	switch {
	case val.Type() == cty.String:
		return fmt.Sprintf("%q", val.AsString())
	case val.Type() == cty.Number:
		return val.AsBigFloat().String()
	case val.Type() == cty.Bool:
		return fmt.Sprintf("%t", val.True())
	case val.Type().IsListType():
		// Handle list values
		if val.LengthInt() == 0 {
			return "[]"
		}
		var items []string
		for it := val.ElementIterator(); it.Next(); {
			_, itemVal := it.Element()
			items = append(items, extractDefaultValue(itemVal))
		}
		return fmt.Sprintf("[%s]", strings.Join(items, ", "))
	case val.Type().IsMapType():
		// Handle map values
		if val.LengthInt() == 0 {
			return "{}"
		}
		var items []string
		for it := val.ElementIterator(); it.Next(); {
			key, itemVal := it.Element()
			items = append(items, fmt.Sprintf("%q = %s", key.AsString(), extractDefaultValue(itemVal)))
		}
		return fmt.Sprintf("{\n    %s\n  }", strings.Join(items, "\n    "))
	default:
		return fmt.Sprintf("%v", val)
	}
}

// formatValue formats a cty.Value for HCL output
func formatValue(val cty.Value) string {
	return extractDefaultValue(val)
}

// generateDefaultForType generates an appropriate default value for a given type
func generateDefaultForType(fieldType string) string {
	switch fieldType {
	case "string":
		return `""`
	case "integer":
		return "0"
	case "number":
		return "0.0"
	case "boolean":
		return "false"
	case "list":
		return "[]"
	case "map":
		return "{}"
	case "object":
		return "{}"
	default:
		return `""`
	}
}

// validateConfigFiles validates that the existing config files are valid HCL
func validateConfigFiles(configDir string) error {
	// Check if spooky.hcl exists and is valid
	spookyConfigPath := filepath.Join(configDir, "spooky.hcl")
	if err := validateHCLFile(spookyConfigPath, "spooky.hcl"); err != nil {
		return err
	}

	return nil
}

// validateHCLFile validates that a file exists and contains valid HCL
func validateHCLFile(filePath, fileName string) error {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("required configuration file %s does not exist at %s", fileName, filePath)
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", fileName, err)
	}

	// Basic HCL validation - check for common syntax errors
	if err := validateHCLSyntax(content, fileName); err != nil {
		return fmt.Errorf("invalid HCL syntax in %s: %w", fileName, err)
	}

	return nil
}

// validateHCLSyntax performs basic HCL syntax validation
func validateHCLSyntax(content []byte, _ string) error {
	contentStr := string(content)

	// Check for balanced braces first (most critical)
	if !hasBalancedBraces(contentStr) {
		return fmt.Errorf("unbalanced braces in HCL content")
	}

	// Check for basic HCL structure
	if !strings.Contains(contentStr, "{") {
		return fmt.Errorf("missing required HCL block structure (no opening braces)")
	}

	// Check for basic assignment syntax
	lines := strings.Split(contentStr, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for basic assignment or block syntax
		if !isValidHCLSyntax(line) {
			return fmt.Errorf("invalid HCL syntax at line %d: %s", i+1, line)
		}
	}

	return nil
}

// hasBalancedBraces checks if the content has balanced braces
func hasBalancedBraces(content string) bool {
	stack := 0
	for _, char := range content {
		switch char {
		case '{':
			stack++
		case '}':
			stack--
			if stack < 0 {
				return false
			}
		}
	}
	return stack == 0
}

// isValidHCLSyntax checks if a line has valid HCL syntax
func isValidHCLSyntax(line string) bool {
	// Skip comments and empty lines
	if strings.HasPrefix(line, "#") || line == "" {
		return true
	}

	// Check for assignment syntax: key = value
	if strings.Contains(line, "=") {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key != "" && value != "" {
				return true
			}
		}
	}

	// Check for block syntax: block_name "name" {
	if strings.Contains(line, "{") {
		// Basic block validation
		return true
	}

	// Check for block closing
	if strings.TrimSpace(line) == "}" {
		return true
	}

	// Allow lines that are just whitespace or comments
	if strings.TrimSpace(line) == "" {
		return true
	}

	return false
}

// GetConfigDirectory returns the config directory path (for external use)
func GetConfigDirectory() (string, error) {
	return getConfigDirectory()
}
