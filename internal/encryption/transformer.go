package encryption

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Transformer handles the transformation of unencrypted values to encrypted values
type Transformer struct {
	ageEncryption *AgeEncryption
}

// NewTransformer creates a new transformer
func NewTransformer(ageEncryption *AgeEncryption) *Transformer {
	return &Transformer{
		ageEncryption: ageEncryption,
	}
}

// TransformVariable transforms a variable from unencrypted to encrypted format
func (t *Transformer) TransformVariable(variable map[string]interface{}) (map[string]interface{}, error) {
	// Check if this variable should be encrypted
	encrypted, exists := variable["encrypted"]
	if !exists {
		return variable, nil // No encryption flag, leave as-is
	}

	encryptedBool, ok := encrypted.(bool)
	if !ok {
		return variable, nil // Not a boolean, leave as-is
	}

	if !encryptedBool {
		return variable, nil // Not marked for encryption, leave as-is
	}

	// Get the value to encrypt
	value, exists := variable["value"]
	if !exists {
		return nil, errors.New("variable marked for encryption but has no value field")
	}

	valueStr, ok := value.(string)
	if !ok {
		return nil, errors.New("variable value must be a string for encryption")
	}

	// Check if already encrypted
	if t.ageEncryption.IsEncrypted(valueStr) {
		// Already encrypted, just update the flags
		result := make(map[string]interface{})
		for k, v := range variable {
			result[k] = v
		}
		result["encrypted"] = false
		result["sensitive"] = true
		return result, nil
	}

	// Create encrypted value
	encryptedValue, err := t.CreateEncryptedValue(valueStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create encrypted value")
	}

	// Create new variable map
	result := make(map[string]interface{})
	for k, v := range variable {
		if k != "value" && k != "encrypted" {
			result[k] = v
		}
	}

	// Add encrypted value and update flags
	result["encrypted_value"] = encryptedValue
	result["encrypted"] = false // No longer needed
	result["sensitive"] = true  // Automatically sensitive when encrypted

	return result, nil
}

// TransformVariablesMap transforms all variables in a map
func (t *Transformer) TransformVariablesMap(variables map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for key, value := range variables {
		if variableMap, ok := value.(map[string]interface{}); ok {
			transformedVariable, err := t.TransformVariable(variableMap)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to transform variable %s", key)
			}
			result[key] = transformedVariable
		} else {
			result[key] = value
		}
	}

	return result, nil
}

// TransformMachineAuthentication transforms authentication fields in a machine
func (t *Transformer) TransformMachineAuthentication(auth map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for key, value := range auth {
		switch key {
		case "passphrase", "password", "certificate_passphrase":
			// These fields can be encrypted
			if authField, ok := value.(map[string]interface{}); ok {
				transformedField, err := t.TransformVariable(authField)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to transform %s", key)
				}
				result[key] = transformedField
			} else {
				result[key] = value
			}
		default:
			// Other fields are not encrypted
			result[key] = value
		}
	}

	return result, nil
}

// TransformMachine transforms a single machine
func (t *Transformer) TransformMachine(machine map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for key, value := range machine {
		switch key {
		case ResourceTypeVariables:
			// Transform machine variables
			if variablesMap, ok := value.(map[string]interface{}); ok {
				transformedVariables, err := t.TransformVariablesMap(variablesMap)
				if err != nil {
					return nil, errors.Wrap(err, "failed to transform machine variables")
				}
				result[key] = transformedVariables
			} else {
				result[key] = value
			}

		case "authentication":
			// Transform authentication fields
			if authMap, ok := value.(map[string]interface{}); ok {
				transformedAuth, err := t.TransformMachineAuthentication(authMap)
				if err != nil {
					return nil, errors.Wrap(err, "failed to transform machine authentication")
				}
				result[key] = transformedAuth
			} else {
				result[key] = value
			}

		default:
			// Other fields are not transformed
			result[key] = value
		}
	}

	return result, nil
}

// TransformMachinesMap transforms all machines in a map
func (t *Transformer) TransformMachinesMap(machines map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for machineName, machineData := range machines {
		if machineMap, ok := machineData.(map[string]interface{}); ok {
			transformedMachine, err := t.TransformMachine(machineMap)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to transform machine %s", machineName)
			}
			result[machineName] = transformedMachine
		} else {
			result[machineName] = machineData
		}
	}

	return result, nil
}

// CreateEncryptedValue creates a structured encrypted value
func (t *Transformer) CreateEncryptedValue(plaintext string) (map[string]interface{}, error) {
	// Security-critical: Validate input before encryption
	if plaintext == "" {
		return nil, errors.New("cannot encrypt empty plaintext - this could lead to data loss")
	}

	// Encrypt the plaintext
	encryptedArmored, err := t.ageEncryption.Encrypt(plaintext)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt value - encryption operation failed")
	}

	// Security-critical: Validate encrypted result
	if encryptedArmored == "" {
		return nil, errors.New("encryption produced empty result - this may indicate a security issue")
	}

	// Extract the base64 content (remove headers and footers)
	lines := strings.Split(encryptedArmored, "\n")
	var base64Content []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "-----BEGIN AGE ENCRYPTED FILE-----" &&
			line != "-----END AGE ENCRYPTED FILE-----" &&
			line != "" {
			base64Content = append(base64Content, line)
		}
	}

	// Security-critical: Validate extracted content
	if len(base64Content) == 0 {
		return nil, errors.New("failed to extract encrypted content - encryption result is malformed")
	}

	data := strings.Join(base64Content, "")

	return map[string]interface{}{
		"data":         data,
		"format":       "base64",
		"algorithm":    "age",
		"version":      "v1",
		"encrypted_at": time.Now().UTC().Format(time.RFC3339),
	}, nil
}
