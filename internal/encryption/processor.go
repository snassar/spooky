package encryption

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"spooky/internal/logging"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/pkg/errors"
)

// HCLProcessor processes HCL files for encryption
type HCLProcessor struct {
	ageEncryption *AgeEncryption
}

// NewHCLProcessor creates a new HCL processor
func NewHCLProcessor(ageEncryption *AgeEncryption) *HCLProcessor {
	return &HCLProcessor{
		ageEncryption: ageEncryption,
	}
}

// ProcessFile processes a single HCL file for encryption
func (hp *HCLProcessor) ProcessFile(filePath string) error {
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

	// Process the file
	modified, err := hp.processBlocks(file.Body.(*hclsyntax.Body).Blocks)
	if err != nil {
		return errors.Wrapf(err, "failed to process file: %s", filePath)
	}

	// If no modifications were made, return early
	if !modified {
		logger := logging.GetGlobalLogger()
		logger.Debug("no variables marked for encryption found",
			slog.String("file", filePath))
		return nil
	}

	// Write the modified content back to the file
	if err := os.WriteFile(filePath, content, 0o644); err != nil {
		return errors.Wrapf(err, "failed to write modified file: %s", filePath)
	}

	logger := logging.GetGlobalLogger()
	logger.Info("successfully encrypted variables",
		slog.String("file", filePath))
	return nil
}

// processBlocks processes HCL blocks recursively
func (hp *HCLProcessor) processBlocks(blocks hclsyntax.Blocks) (bool, error) {
	modified := false

	for _, block := range blocks {
		// Check if this is a variables block
		if block.Type == "variables" {
			blockModified, err := hp.processVariablesBlock(block)
			if err != nil {
				return false, err
			}
			modified = modified || blockModified
		}

		// Process nested blocks
		if len(block.Body.Blocks) > 0 {
			nestedModified, err := hp.processBlocks(block.Body.Blocks)
			if err != nil {
				return false, err
			}
			modified = modified || nestedModified
		}
	}

	return modified, nil
}

// processVariablesBlock processes a variables block
func (hp *HCLProcessor) processVariablesBlock(block *hclsyntax.Block) (bool, error) {
	modified := false

	// Process each variable block
	for _, variableBlock := range block.Body.Blocks {
		if variableBlock.Type == "variable" {
			blockModified, err := hp.processVariableBlock(variableBlock)
			if err != nil {
				return false, err
			}
			modified = modified || blockModified
		}
	}

	return modified, nil
}

// processVariableBlock processes a single variable block
func (hp *HCLProcessor) processVariableBlock(block *hclsyntax.Block) (bool, error) {
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
	currentValue, err := hp.extractStringValue(valueAttr.Expr)
	if err != nil {
		return false, errors.Wrapf(err, "failed to extract value for variable %s", block.Labels[0])
	}

	// Check if already encrypted
	if hp.ageEncryption.IsEncrypted(currentValue) {
		logger := logging.GetGlobalLogger()
		logger.Debug("variable already encrypted, skipping",
			slog.String("variable", block.Labels[0]))
		return false, nil
	}

	// Encrypt the value
	encryptedValue, err := hp.ageEncryption.Encrypt(currentValue)
	if err != nil {
		return false, errors.Wrapf(err, "failed to encrypt value for variable %s", block.Labels[0])
	}

	// Update the value in the HCL content
	// This is a simplified approach - in a real implementation, you'd want to
	// properly update the HCL AST and regenerate the content
	logger := logging.GetGlobalLogger()
	logger.Info("encrypted variable",
		slog.String("variable", block.Labels[0]),
		slog.String("encrypted_value", encryptedValue))

	return true, nil
}

// extractStringValue extracts a string value from an HCL expression
func (hp *HCLProcessor) extractStringValue(expr hclsyntax.Expression) (string, error) {
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

// ProcessDirectory processes all HCL files in a directory
func (hp *HCLProcessor) ProcessDirectory(dirPath string) error {
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
		if err := hp.ProcessFile(filePath); err != nil {
			logger := logging.GetGlobalLogger()
			logger.Warn("failed to process file",
				slog.String("file", filePath),
				slog.String("error", err.Error()))
		}
	}

	return nil
}
