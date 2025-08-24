package encryption

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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

	contentStr := string(content)
	modified := false

	// Find all variable blocks with encrypted = true
	// This regex matches: variable "name" { ... value = "..." ... encrypted = true ... }
	pattern := regexp.MustCompile(`variable\s+"([^"]+)"\s*\{[^}]*value\s*=\s*"([^"]*)"[^}]*encrypted\s*=\s*true[^}]*\}`)

	matches := pattern.FindAllStringSubmatch(contentStr, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			variableName := match[1]
			originalValue := match[2]

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

			// Replace the value in the content
			oldValue := fmt.Sprintf(`value = "%s"`, originalValue)
			newValue := fmt.Sprintf(`value = "%s"`, strings.ReplaceAll(encryptedValue, `"`, `\"`))

			contentStr = strings.Replace(contentStr, oldValue, newValue, 1)
			modified = true

			fmt.Printf("Encrypted variable %s\n", variableName)
			fmt.Printf("Original value: %s\n", originalValue)
			fmt.Printf("Encrypted value: %s\n", encryptedValue)
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
