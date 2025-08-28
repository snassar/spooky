package encryption

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/pkg/errors"
)

// SimpleHCLUpdater updates HCL files with encrypted values using string replacement
type SimpleHCLUpdater struct {
	ageEncryption *AgeEncryption
}

// NewSimpleHCLUpdater creates a new simple HCL updater
func NewSimpleHCLUpdater(ageEncryption *AgeEncryption) *SimpleHCLUpdater {
	return &SimpleHCLUpdater{
		ageEncryption: ageEncryption,
	}
}

// UpdateFile updates a single HCL file with encrypted values
func (su *SimpleHCLUpdater) UpdateFile(filePath string) error {
	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read file: %s", filePath)
	}

	// Parse HCL content
	file, diags := hclsyntax.ParseConfig(content, filepath.Base(filePath), hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return errors.Wrapf(err, "failed to parse HCL file: %s", filePath)
	}

	// Define schema for variable blocks
	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
			},
		},
	}

	// Extract variable blocks
	bodyContent, diags := file.Body.Content(schema)
	if diags.HasErrors() {
		return errors.Wrapf(err, "failed to decode variable blocks: %s", filePath)
	}

	modified := false
	contentStr := string(content)

	// Process each variable block
	for _, block := range bodyContent.Blocks {
		variableName := block.Labels[0]

		// Define schema for variable attributes
		varSchema := &hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{Name: "value", Required: false},
				{Name: "encrypted", Required: false},
			},
		}

		// Extract variable attributes
		varContent, diags := block.Body.Content(varSchema)
		if diags.HasErrors() {
			continue // Skip this variable if we can't parse it
		}

		// Check if encrypted flag is set to true
		if encryptedAttr, exists := varContent.Attributes["encrypted"]; exists {
			var encrypted bool
			if diags := gohcl.DecodeExpression(encryptedAttr.Expr, nil, &encrypted); diags.HasErrors() {
				continue // Skip if we can't decode the encrypted flag
			}

			if !encrypted {
				continue // Skip if not marked for encryption
			}
		} else {
			continue // Skip if no encrypted flag
		}

		// Get the value to encrypt
		if valueAttr, exists := varContent.Attributes["value"]; exists {
			var originalValue string
			if diags := gohcl.DecodeExpression(valueAttr.Expr, nil, &originalValue); diags.HasErrors() {
				continue // Skip if we can't decode the value
			}

			// Check if already encrypted
			if su.ageEncryption.IsEncrypted(originalValue) {
				fmt.Printf("Variable %s is already encrypted, skipping\n", variableName)
				continue
			}

			// Encrypt the value
			encryptedValue, err := su.ageEncryption.Encrypt(originalValue)
			if err != nil {
				return errors.Wrapf(err, "failed to encrypt value for variable %s", variableName)
			}

			// Replace the value in the content using the original position
			startPos := valueAttr.Expr.Range().Start
			endPos := valueAttr.Expr.Range().End

			if startPos.Byte < len(content) && endPos.Byte <= len(content) {
				newValue := fmt.Sprintf(`"%s"`, strings.ReplaceAll(encryptedValue, `"`, `\"`))

				contentStr = contentStr[:startPos.Byte] + newValue + contentStr[endPos.Byte:]
				modified = true

				fmt.Printf("Encrypted variable %s\n", variableName)
				fmt.Printf("Original value: %s\n", originalValue)
				fmt.Printf("Encrypted value: %s\n", encryptedValue)
			}
		}
	}

	// If no modifications were made, return early
	if !modified {
		fmt.Printf("No variables marked for encryption found in %s\n", filePath)
		return nil
	}

	// Write the modified content back to the file
	if err := os.WriteFile(filePath, []byte(contentStr), 0o644); err != nil {
		return errors.Wrapf(err, "failed to write modified file: %s", filePath)
	}

	fmt.Printf("Successfully encrypted variables in %s\n", filePath)
	return nil
}

// UpdateDirectory updates all HCL files in a directory
func (su *SimpleHCLUpdater) UpdateDirectory(dirPath string) error {
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
		if err := su.UpdateFile(filePath); err != nil {
			fmt.Printf("Warning: failed to process %s: %v\n", filePath, err)
		}
	}

	return nil
}
