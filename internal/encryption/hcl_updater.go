package encryption

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/pkg/errors"
)

// HCLUpdater updates HCL files with encrypted values
type HCLUpdater struct {
	ageEncryption *AgeEncryption
}

// NewHCLUpdater creates a new HCL updater
func NewHCLUpdater(ageEncryption *AgeEncryption) *HCLUpdater {
	return &HCLUpdater{
		ageEncryption: ageEncryption,
	}
}

// UpdateFile updates a single HCL file with encrypted values
func (hu *HCLUpdater) UpdateFile(filePath string) error {
	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read file: %s", filePath)
	}

	// Parse HCL
	file, diags := hclsyntax.ParseConfig(content, filepath.Base(filePath), hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return errors.Wrapf(err, "failed to parse HCL file: %s", filePath)
	}

	// Convert content to string for manipulation
	contentStr := string(content)
	modified := false

	// Process variables blocks
	for _, block := range file.Body.(*hclsyntax.Body).Blocks {
		if block.Type == "variables" {
			blockModified, err := hu.processVariablesBlock(block, &contentStr)
			if err != nil {
				return errors.Wrapf(err, "failed to process variables block in %s", filePath)
			}
			modified = modified || blockModified
		}
	}

	// If no modifications were made, return early
	if !modified {
		fmt.Printf("No variables marked for encryption found in %s\n", filePath)
		return nil
	}

	// Write the modified content back to the file
	if err := os.WriteFile(filePath, []byte(contentStr), 0644); err != nil {
		return errors.Wrapf(err, "failed to write modified file: %s", filePath)
	}

	fmt.Printf("Successfully encrypted variables in %s\n", filePath)
	return nil
}

// processVariablesBlock processes a variables block and updates the content string
func (hu *HCLUpdater) processVariablesBlock(block *hclsyntax.Block, contentStr *string) (bool, error) {
	modified := false

	// Process each variable block
	for _, variableBlock := range block.Body.Blocks {
		if variableBlock.Type == "variable" {
			blockModified, err := hu.processVariableBlock(variableBlock, contentStr)
			if err != nil {
				return false, err
			}
			modified = modified || blockModified
		}
	}

	return modified, nil
}

// processVariableBlock processes a single variable block and updates the content string
func (hu *HCLUpdater) processVariableBlock(block *hclsyntax.Block, contentStr *string) (bool, error) {
	// Check if this variable is marked for encryption
	encrypted := false
	for _, attr := range block.Body.Attributes {
		if attr.Name == "encrypted" {
			if val, ok := attr.Expr.(*hclsyntax.LiteralValueExpr); ok {
				if val.Val.Type().IsPrimitiveType() && val.Val.True() {
					encrypted = true
				}
			}
		}
	}

	if !encrypted {
		return false, nil
	}

	// Find the value attribute
	var valueAttr *hclsyntax.Attribute
	for _, attr := range block.Body.Attributes {
		if attr.Name == "value" {
			valueAttr = attr
			break
		}
	}

	if valueAttr == nil {
		return false, errors.Errorf("variable %s is marked as encrypted but has no value attribute", block.Labels[0])
	}

	// Get the current value
	currentValue, err := hu.extractStringValue(valueAttr.Expr)
	if err != nil {
		return false, errors.Wrapf(err, "failed to extract value for variable %s", block.Labels[0])
	}

	// Check if already encrypted
	if hu.ageEncryption.IsEncrypted(currentValue) {
		fmt.Printf("Variable %s is already encrypted, skipping\n", block.Labels[0])
		return false, nil
	}

	// Encrypt the value
	encryptedValue, err := hu.ageEncryption.Encrypt(currentValue)
	if err != nil {
		return false, errors.Wrapf(err, "failed to encrypt value for variable %s", block.Labels[0])
	}

	// Update the value in the content string
	// Find the range of the value attribute and replace it
	valueRange := valueAttr.Range()
	start := valueRange.Start.Byte
	end := valueRange.End.Byte

	// Extract the part before the value
	before := (*contentStr)[:start]

	// Find the actual value part (after the equals sign and whitespace)
	valueStart := start
	for i := start; i < len(*contentStr); i++ {
		if (*contentStr)[i] == '=' {
			// Skip equals sign and whitespace
			for j := i + 1; j < len(*contentStr); j++ {
				if (*contentStr)[j] != ' ' && (*contentStr)[j] != '\t' {
					valueStart = j
					break
				}
			}
			break
		}
	}

	// Find the end of the value (before the next newline or closing brace)
	valueEnd := end
	for i := valueStart; i < len(*contentStr); i++ {
		if (*contentStr)[i] == '\n' || (*contentStr)[i] == '}' {
			valueEnd = i
			break
		}
	}

	// Extract the part after the value
	after := (*contentStr)[valueEnd:]

	// Create the new value with proper quoting
	newValue := fmt.Sprintf(`"%s"`, strings.ReplaceAll(encryptedValue, `"`, `\"`))

	// Reconstruct the content
	*contentStr = before + newValue + after

	fmt.Printf("Encrypted variable %s\n", block.Labels[0])
	fmt.Printf("Original value: %s\n", currentValue)
	fmt.Printf("Encrypted value: %s\n", encryptedValue)

	return true, nil
}

// extractStringValue extracts a string value from an HCL expression
func (hu *HCLUpdater) extractStringValue(expr hclsyntax.Expression) (string, error) {
	switch v := expr.(type) {
	case *hclsyntax.LiteralValueExpr:
		if v.Val.Type().IsPrimitiveType() {
			return v.Val.AsString(), nil
		}
		return "", errors.New("expression is not a string literal")
	case *hclsyntax.TemplateExpr:
		// Handle template expressions (interpolated strings)
		parts := make([]string, len(v.Parts))
		for i, part := range v.Parts {
			if lit, ok := part.(*hclsyntax.LiteralValueExpr); ok {
				parts[i] = lit.Val.AsString()
			} else {
				return "", errors.New("template expression contains non-literal parts")
			}
		}
		return strings.Join(parts, ""), nil
	default:
		return "", errors.New("unsupported expression type")
	}
}

// UpdateDirectory updates all HCL files in a directory
func (hu *HCLUpdater) UpdateDirectory(dirPath string) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return errors.Wrapf(err, "failed to read directory: %s", dirPath)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(entry.Name(), ".hcl") {
			continue
		}

		filePath := filepath.Join(dirPath, entry.Name())
		if err := hu.UpdateFile(filePath); err != nil {
			fmt.Printf("Warning: failed to process %s: %v\n", filePath, err)
		}
	}

	return nil
}
