package encryption

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// TransparentDecryptor provides transparent decryption for template context
type TransparentDecryptor struct {
	ageEncryption *AgeEncryption
}

// NewTransparentDecryptor creates a new transparent decryptor
func NewTransparentDecryptor(ageEncryption *AgeEncryption) *TransparentDecryptor {
	return &TransparentDecryptor{
		ageEncryption: ageEncryption,
	}
}

// DecryptVariable decrypts a variable value if it's encrypted
func (td *TransparentDecryptor) DecryptVariable(variable interface{}) (interface{}, error) {
	// Handle different variable types
	switch v := variable.(type) {
	case map[string]interface{}:
		return td.decryptMapVariable(v)
	case string:
		// Check if it's a simple encrypted string (legacy format)
		if td.ageEncryption.IsEncrypted(v) {
			return td.ageEncryption.Decrypt(v)
		}
		return v, nil
	default:
		return v, nil
	}
}

// decryptMapVariable decrypts a variable stored as a map (structured format)
func (td *TransparentDecryptor) decryptMapVariable(variable map[string]interface{}) (interface{}, error) {
	// Check if this variable has an encrypted_value field
	if encryptedValue, exists := variable["encrypted_value"]; exists {
		return td.decryptEncryptedValue(encryptedValue)
	}

	// Check if this variable has a value field that might be encrypted
	if value, exists := variable["value"]; exists {
		if valueStr, ok := value.(string); ok {
			if td.ageEncryption.IsEncrypted(valueStr) {
				return td.ageEncryption.Decrypt(valueStr)
			}
		}
		return value, nil
	}

	// No encryption, return as-is
	return variable, nil
}

// decryptEncryptedValue decrypts a structured encrypted value
func (td *TransparentDecryptor) decryptEncryptedValue(encryptedValue interface{}) (string, error) {
	// Convert to map if it's not already
	var encryptedMap map[string]interface{}
	switch ev := encryptedValue.(type) {
	case map[string]interface{}:
		encryptedMap = ev
	default:
		return "", errors.New("encrypted_value must be a map - decryption cannot proceed")
	}

	// Extract the data field
	data, exists := encryptedMap["data"]
	if !exists {
		return "", errors.New("encrypted_value missing required 'data' field - decryption cannot proceed")
	}

	dataStr, ok := data.(string)
	if !ok {
		return "", errors.New("encrypted_value.data must be a string - decryption cannot proceed")
	}

	// Security-critical: Validate data before decryption
	if dataStr == "" {
		return "", errors.New("encrypted_value.data is empty - decryption cannot proceed")
	}

	// Get the format (default to base64)
	format := "base64"
	if formatVal, exists := encryptedMap["format"]; exists {
		if formatStr, ok := formatVal.(string); ok {
			format = formatStr
		}
	}

	// Reconstruct the full armored format if needed
	var encryptedData string
	switch format {
	case "base64":
		// Reconstruct the full armored format
		encryptedData = fmt.Sprintf("-----BEGIN AGE ENCRYPTED FILE-----\n%s\n-----END AGE ENCRYPTED FILE-----", dataStr)
	case "armored":
		// Already in armored format
		encryptedData = dataStr
	case "compact":
		// Remove line breaks and reconstruct
		cleanData := strings.ReplaceAll(dataStr, "\n", "")
		encryptedData = fmt.Sprintf("-----BEGIN AGE ENCRYPTED FILE-----%s-----END AGE ENCRYPTED FILE-----", cleanData)
	default:
		return "", errors.Errorf("unsupported encryption format: %s - decryption cannot proceed", format)
	}

	// Security-critical: Validate reconstructed data
	if encryptedData == "" {
		return "", errors.New("failed to reconstruct encrypted data - decryption cannot proceed")
	}

	// Decrypt using age encryption
	decrypted, err := td.ageEncryption.Decrypt(encryptedData)
	if err != nil {
		return "", errors.Wrap(err, "failed to decrypt value - decryption operation failed")
	}

	// Security-critical: Validate decrypted result
	if decrypted == "" {
		return "", errors.New("decryption produced empty result - this may indicate a security issue or corrupted data")
	}

	return decrypted, nil
}

// DecryptVariablesMap decrypts all encrypted variables in a map
func (td *TransparentDecryptor) DecryptVariablesMap(variables map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for key, value := range variables {
		decryptedValue, err := td.DecryptVariable(value)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to decrypt variable %s", key)
		}
		result[key] = decryptedValue
	}

	return result, nil
}

// DecryptMachineVariables decrypts machine-specific variables
func (td *TransparentDecryptor) DecryptMachineVariables(machine map[string]interface{}) (map[string]interface{}, error) {
	// Check if machine has variables
	if variables, exists := machine[ResourceTypeVariables]; exists {
		if variablesMap, ok := variables.(map[string]interface{}); ok {
			decryptedVariables, err := td.DecryptVariablesMap(variablesMap)
			if err != nil {
				return nil, err
			}

			// Create a copy of the machine map and update variables
			result := make(map[string]interface{})
			for k, v := range machine {
				result[k] = v
			}
			result[ResourceTypeVariables] = decryptedVariables
			return result, nil
		}
	}

	return machine, nil
}

// DecryptMachinesMap decrypts all encrypted values in a machines map
func (td *TransparentDecryptor) DecryptMachinesMap(machines map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for machineName, machineData := range machines {
		if machineMap, ok := machineData.(map[string]interface{}); ok {
			// Decrypt machine variables
			decryptedMachine, err := td.DecryptMachineVariables(machineMap)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to decrypt machine %s", machineName)
			}

			// Decrypt authentication fields
			if auth, exists := decryptedMachine["authentication"]; exists {
				if authMap, ok := auth.(map[string]interface{}); ok {
					decryptedAuth, err := td.decryptAuthentication(authMap)
					if err != nil {
						return nil, errors.Wrapf(err, "failed to decrypt authentication for machine %s", machineName)
					}
					decryptedMachine["authentication"] = decryptedAuth
				}
			}

			result[machineName] = decryptedMachine
		} else {
			result[machineName] = machineData
		}
	}

	return result, nil
}

// decryptAuthentication decrypts authentication fields in a machine
func (td *TransparentDecryptor) decryptAuthentication(auth map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for key, value := range auth {
		switch key {
		case "passphrase", "password", "certificate_passphrase":
			// These fields can be encrypted
			decryptedValue, err := td.DecryptVariable(value)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to decrypt %s", key)
			}
			result[key] = decryptedValue
		default:
			// Other fields are not encrypted
			result[key] = value
		}
	}

	return result, nil
}

// CreateEncryptedValue creates a structured encrypted value
func (td *TransparentDecryptor) CreateEncryptedValue(plaintext string) (*EncryptedValueV1, error) {
	// Security-critical: Validate input before encryption
	if plaintext == "" {
		return nil, errors.New("cannot encrypt empty plaintext - this could lead to data loss")
	}

	// Encrypt the plaintext
	encryptedArmored, err := td.ageEncryption.Encrypt(plaintext)
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

	return &EncryptedValueV1{
		Data:        data,
		Format:      "base64",
		Algorithm:   "age",
		Version:     "v1",
		EncryptedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// EncryptedValueV1 represents a structured encrypted value
type EncryptedValueV1 struct {
	Data        string `json:"data" required:"true" description:"The encrypted data (base64 content only, no headers/footers)"`
	Format      string `json:"format" enum:"base64,armored,compact" default:"base64" description:"Format of the encrypted data"`
	Algorithm   string `json:"algorithm" default:"age" description:"Encryption algorithm used"`
	Version     string `json:"version" default:"v1" description:"Encryption version"`
	EncryptedAt string `json:"encrypted_at" description:"ISO 8601 timestamp when the value was encrypted"`
}
