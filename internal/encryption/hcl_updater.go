package encryption

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"spooky/internal/logging"
	"spooky/internal/utilities"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/pkg/errors"
)

// Resource type constants to avoid repeated string literals
const (
	ResourceTypeVariables = "variables"
	ResourceTypeVariable  = "variable"
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
		if block.Type == ResourceTypeVariables {
			blockModified, err := hu.processVariablesBlock(block, &contentStr)
			if err != nil {
				return errors.Wrapf(err, "failed to process variables block in %s", filePath)
			}
			modified = modified || blockModified
		}
	}

	// If no modifications were made, return early
	if !modified {
		logger := logging.GetGlobalLogger()
		logger.Debug("no variables marked for encryption found",
			slog.String("file", filePath))
		return nil
	}

	// Write the modified content back to the file
	if err := utilities.WriteFile(filePath, contentStr); err != nil {
		return errors.Wrapf(err, "failed to write modified file: %s", filePath)
	}

	logger := logging.GetGlobalLogger()
	logger.Info("successfully encrypted variables",
		slog.String("file", filePath))
	return nil
}

// processVariablesBlock processes a variables block and updates the content string
func (hu *HCLUpdater) processVariablesBlock(block *hclsyntax.Block, contentStr *string) (bool, error) {
	modified := false

	// Process each variable block
	for _, variableBlock := range block.Body.Blocks {
		if variableBlock.Type == ResourceTypeVariable {
			blockModified, err := hu.processVariableBlock(variableBlock, contentStr)
			if err != nil {
				return false, err
			}
			modified = modified || blockModified
		}
	}

	return modified, nil
}

// isVariableEncrypted checks if a variable block is marked for encryption
func (hu *HCLUpdater) isVariableEncrypted(block *hclsyntax.Block) bool {
	for _, attr := range block.Body.Attributes {
		if attr.Name == "encrypted" {
			if val, ok := attr.Expr.(*hclsyntax.LiteralValueExpr); ok {
				if val.Val.Type().IsPrimitiveType() && val.Val.True() {
					return true
				}
			}
		}
	}
	return false
}

// findValueAttribute finds the value attribute in a variable block
func (hu *HCLUpdater) findValueAttribute(block *hclsyntax.Block) (*hclsyntax.Attribute, error) {
	for _, attr := range block.Body.Attributes {
		if attr.Name == "value" {
			return attr, nil
		}
	}
	return nil, errors.Errorf("variable %s is marked as encrypted but has no value attribute", block.Labels[0])
}

// findValueStart finds the start position of the actual value in the content string
func (hu *HCLUpdater) findValueStart(contentStr string, start int) int {
	for i := start; i < len(contentStr); i++ {
		if contentStr[i] == '=' {
			// Skip equals sign and whitespace
			for j := i + 1; j < len(contentStr); j++ {
				if contentStr[j] != ' ' && contentStr[j] != '\t' {
					return j
				}
			}
			break
		}
	}
	return start
}

// findValueEnd finds the end position of the value in the content string
func (hu *HCLUpdater) findValueEnd(contentStr string, valueStart int) int {
	for i := valueStart; i < len(contentStr); i++ {
		if contentStr[i] == '\n' || contentStr[i] == '}' {
			return i
		}
	}
	return len(contentStr)
}

// replaceContentValue replaces the value in the content string with the encrypted value
func (hu *HCLUpdater) replaceContentValue(contentStr *string, valueAttr *hclsyntax.Attribute, encryptedValue string) {
	valueRange := valueAttr.Range()
	start := valueRange.Start.Byte

	// Extract the part before the value
	before := (*contentStr)[:start]

	// Find the actual value positions
	valueStart := hu.findValueStart(*contentStr, start)
	valueEnd := hu.findValueEnd(*contentStr, valueStart)

	// Extract the part after the value
	after := (*contentStr)[valueEnd:]

	// Create the new value with proper quoting
	newValue := fmt.Sprintf(`"%s"`, strings.ReplaceAll(encryptedValue, `"`, `\"`))

	// Reconstruct the content
	*contentStr = before + newValue + after
}

// processVariableBlock processes a single variable block and updates the content string
func (hu *HCLUpdater) processVariableBlock(block *hclsyntax.Block, contentStr *string) (bool, error) {
	// Check if this variable is marked for encryption
	if !hu.isVariableEncrypted(block) {
		return false, nil
	}

	// Find the value attribute
	valueAttr, err := hu.findValueAttribute(block)
	if err != nil {
		return false, err
	}

	// Get the current value
	currentValue, err := hu.extractStringValue(valueAttr.Expr)
	if err != nil {
		return false, errors.Wrapf(err, "failed to extract value for variable %s", block.Labels[0])
	}

	// Check if already encrypted
	if hu.ageEncryption.IsEncrypted(currentValue) {
		logger := logging.GetGlobalLogger()
		logger.Info("variable already encrypted, skipping", slog.String("variable", block.Labels[0]))
		return false, nil
	}

	// Encrypt the value
	encryptedValue, err := hu.ageEncryption.Encrypt(currentValue)
	if err != nil {
		return false, errors.Wrapf(err, "failed to encrypt value for variable %s", block.Labels[0])
	}

	// Update the value in the content string
	hu.replaceContentValue(contentStr, valueAttr, encryptedValue)

	logger := logging.GetGlobalLogger()
	logger.Info("encrypted variable",
		slog.String("variable", block.Labels[0]),
		slog.String("encrypted_value", encryptedValue))

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
			logger := logging.GetGlobalLogger()
			logger.Warn("failed to process file",
				slog.String("file", filePath),
				slog.String("error", err.Error()))
		}
	}

	return nil
}
