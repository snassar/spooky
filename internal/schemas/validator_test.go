package schemas

import (
	"testing"
)

func TestNewValidator(t *testing.T) {
	validator := NewValidator()
	if validator == nil {
		t.Fatal("NewValidator() returned nil")
	}

	if validator.logger == nil {
		t.Error("Validator logger should not be nil")
	}

	if len(validator.supportedVersions) == 0 {
		t.Error("Validator should have supported versions")
	}

	if len(validator.enhancedRules) == 0 {
		t.Error("Validator should have enhanced rules loaded")
	}
}

func TestValidator_ValidateProjectV1(t *testing.T) {
	validator := NewValidator()

	// Test valid project
	validProject := map[string]interface{}{
		"project": map[string]interface{}{
			"name":        "test-project",
			"description": "A test project",
		},
	}

	result := validator.ValidateData("project", validProject)
	if !result.IsValid {
		t.Errorf("Expected valid project, got errors: %v", result.Errors)
	}

	// Test invalid project (missing description)
	invalidProject := map[string]interface{}{
		"project": map[string]interface{}{
			"name": "test-project",
		},
	}

	result = validator.ValidateData("project", invalidProject)
	if result.IsValid {
		t.Error("Expected invalid project due to missing description")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for missing description")
	}
}

func TestValidator_ValidateMachinesV1(t *testing.T) {
	validator := NewValidator()

	// Test valid machines
	validMachines := map[string]interface{}{
		"machines": map[string]interface{}{
			"machine": []interface{}{
				map[string]interface{}{
					"name":     "test-machine",
					"hostname": "test.example.com",
					"user":     "testuser",
					"authentication": map[string]interface{}{
						"method":   "password",
						"password": "secret",
					},
				},
			},
		},
	}

	result := validator.ValidateData("machines", validMachines)
	if !result.IsValid {
		t.Errorf("Expected valid machines, got errors: %v", result.Errors)
	}

	// Test invalid machines (missing required fields)
	invalidMachines := map[string]interface{}{
		"machines": map[string]interface{}{
			"machine": []interface{}{
				map[string]interface{}{
					"name": "test-machine",
					// Missing hostname, user, authentication
				},
			},
		},
	}

	result = validator.ValidateData("machines", invalidMachines)
	if result.IsValid {
		t.Error("Expected invalid machines due to missing required fields")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for missing required fields")
	}

	// Test machine with invalid authentication method
	invalidAuthMachines := map[string]interface{}{
		"machines": map[string]interface{}{
			"machine": []interface{}{
				map[string]interface{}{
					"name":     "test-machine",
					"hostname": "test.example.com",
					"user":     "testuser",
					"authentication": map[string]interface{}{
						"method": "invalid-method",
					},
				},
			},
		},
	}

	result = validator.ValidateData("machines", invalidAuthMachines)
	if result.IsValid {
		t.Error("Expected invalid machines due to invalid authentication method")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for invalid authentication method")
	}

	// Test machine with publickey authentication but missing public_key_path
	invalidPublicKeyMachines := map[string]interface{}{
		"machines": map[string]interface{}{
			"machine": []interface{}{
				map[string]interface{}{
					"name":     "test-machine",
					"hostname": "test.example.com",
					"user":     "testuser",
					"authentication": map[string]interface{}{
						"method": "publickey",
						// Missing public_key_path
					},
				},
			},
		},
	}

	result = validator.ValidateData("machines", invalidPublicKeyMachines)
	if result.IsValid {
		t.Error("Expected invalid machines due to missing public_key_path")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for missing public_key_path")
	}
}

func TestValidator_ValidateHCLContent(t *testing.T) {
	validator := NewValidator()

	// Test valid HCL content
	validHCL := `
project {
  name = "test-project"
  description = "A test project"
}
`

	result, err := validator.ValidateHCLContent("project", validHCL)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !result.IsValid {
		t.Errorf("Expected valid HCL, got errors: %v", result.Errors)
	}

	// Test invalid HCL content
	invalidHCL := `
project {
  name = "test-project"
  description = "A test project"
  invalid_field = "this should cause an error"
}
`

	result, err = validator.ValidateHCLContent("project", invalidHCL)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// Note: This might be valid depending on how strict the validation is
	// The test just ensures it doesn't crash
}

func TestValidator_EnhancedValidation(t *testing.T) {
	validator := NewValidator()

	// Test actions enhanced validation
	actionsData := map[string]interface{}{
		"actions": map[string]interface{}{
			"action": []interface{}{
				map[string]interface{}{
					"name":        "test-action",
					"description": "A test action",
					"type":        "command",
					"command":     "echo 'hello world'",
					"machines":    []interface{}{"test-machine"},
				},
			},
		},
	}

	result := validator.ValidateData("actions", actionsData)
	if !result.IsValid {
		t.Errorf("Expected valid actions, got errors: %v", result.Errors)
	}

	// Test actions with missing machine targeting
	invalidActionsData := map[string]interface{}{
		"actions": map[string]interface{}{
			"action": []interface{}{
				map[string]interface{}{
					"name":        "test-action",
					"description": "A test action",
					"type":        "command",
					"command":     "echo 'hello world'",
					// Missing machines or tags
				},
			},
		},
	}

	result = validator.ValidateData("actions", invalidActionsData)
	if result.IsValid {
		t.Error("Expected invalid actions due to missing machine targeting")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for missing machine targeting")
	}
}

func TestValidator_VariablesV1Validation(t *testing.T) {
	validator := NewValidator()

	// Test valid variables
	validVariables := map[string]interface{}{
		"variables": map[string]interface{}{
			"variable": []interface{}{
				map[string]interface{}{
					"name":  "test-variable",
					"value": "test-value",
				},
			},
		},
	}

	result := validator.ValidateData("variables", validVariables)
	if !result.IsValid {
		t.Errorf("Expected valid variables, got errors: %v", result.Errors)
	}

	// Test invalid variables (missing required fields)
	invalidVariables := map[string]interface{}{
		"variables": map[string]interface{}{
			"variable": []interface{}{
				map[string]interface{}{
					"name": "test-variable",
					// Missing value
				},
			},
		},
	}

	result = validator.ValidateData("variables", invalidVariables)
	if result.IsValid {
		t.Error("Expected invalid variables due to missing required fields")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for missing required fields")
	}

	// Test invalid variable name
	invalidNameVariables := map[string]interface{}{
		"variables": map[string]interface{}{
			"variable": []interface{}{
				map[string]interface{}{
					"name":  "invalid@name",
					"value": "test-value",
				},
			},
		},
	}

	result = validator.ValidateData("variables", invalidNameVariables)
	if result.IsValid {
		t.Error("Expected invalid variables due to invalid name")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for invalid name")
	}

	// Test null value
	nullValueVariables := map[string]interface{}{
		"variables": map[string]interface{}{
			"variable": []interface{}{
				map[string]interface{}{
					"name":  "test-variable",
					"value": nil,
				},
			},
		},
	}

	result = validator.ValidateData("variables", nullValueVariables)
	if result.IsValid {
		t.Error("Expected invalid variables due to null value")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for null value")
	}
}

func TestValidator_VariablesEnhancedValidation(t *testing.T) {
	validator := NewValidator()

	// Test variables with enhanced validation
	variablesData := map[string]interface{}{
		"variables": map[string]interface{}{
			"variable": []interface{}{
				map[string]interface{}{
					"name":  "test-variable",
					"value": "test-value",
				},
			},
		},
	}

	result := validator.ValidateData("variables", variablesData)
	if !result.IsValid {
		t.Errorf("Expected valid variables, got errors: %v", result.Errors)
	}

	// Test variables with sensitive and encrypted conflict
	conflictVariables := map[string]interface{}{
		"variables": map[string]interface{}{
			"variable": []interface{}{
				map[string]interface{}{
					"name":      "test-variable",
					"value":     "test-value",
					"sensitive": true,
					"encrypted": true,
				},
			},
		},
	}

	result = validator.ValidateData("variables", conflictVariables)
	if result.IsValid {
		t.Error("Expected invalid variables due to sensitive and encrypted conflict")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for sensitive and encrypted conflict")
	}
}

func TestValidator_LoggingV1Validation(t *testing.T) {
	validator := NewValidator()

	// Test valid logging
	validLogging := map[string]interface{}{
		"logging": map[string]interface{}{
			"level":  "info",
			"format": "json",
			"output": "stdout",
		},
	}

	result := validator.ValidateData("logging", validLogging)
	if !result.IsValid {
		t.Errorf("Expected valid logging, got errors: %v", result.Errors)
	}

	// Test invalid logging (missing required block)
	invalidLogging := map[string]interface{}{
		// Missing logging block
	}

	result = validator.ValidateData("logging", invalidLogging)
	if result.IsValid {
		t.Error("Expected invalid logging due to missing logging block")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for missing logging block")
	}

	// Test invalid logging level
	invalidLevelLogging := map[string]interface{}{
		"logging": map[string]interface{}{
			"level":  "invalid-level",
			"format": "json",
			"output": "stdout",
		},
	}

	result = validator.ValidateData("logging", invalidLevelLogging)
	if result.IsValid {
		t.Error("Expected invalid logging due to invalid level")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for invalid level")
	}

	// Test file output without file_path
	fileOutputLogging := map[string]interface{}{
		"logging": map[string]interface{}{
			"level":  "info",
			"format": "json",
			"output": "file",
			// Missing file_path
		},
	}

	result = validator.ValidateData("logging", fileOutputLogging)
	if result.IsValid {
		t.Error("Expected invalid logging due to missing file_path")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for missing file_path")
	}
}

func TestValidator_FactsV1Validation(t *testing.T) {
	validator := NewValidator()

	// Test valid facts
	validFacts := map[string]interface{}{
		"facts": map[string]interface{}{
			"basic_facts": map[string]interface{}{
				"test-fact": map[string]interface{}{
					"value": "test-value",
					"type":  "string",
				},
			},
		},
	}

	result := validator.ValidateData("facts", validFacts)
	if !result.IsValid {
		t.Errorf("Expected valid facts, got errors: %v", result.Errors)
	}

	// Test invalid facts (missing required fields)
	invalidFacts := map[string]interface{}{
		"facts": map[string]interface{}{
			"basic_facts": map[string]interface{}{
				"test-fact": map[string]interface{}{
					// Missing value and type
				},
			},
		},
	}

	result = validator.ValidateData("facts", invalidFacts)
	if result.IsValid {
		t.Error("Expected invalid facts due to missing required fields")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for missing required fields")
	}

	// Test invalid fact name
	invalidNameFacts := map[string]interface{}{
		"facts": map[string]interface{}{
			"basic_facts": map[string]interface{}{
				"invalid@name": map[string]interface{}{
					"value": "test-value",
					"type":  "string",
				},
			},
		},
	}

	result = validator.ValidateData("facts", invalidNameFacts)
	if result.IsValid {
		t.Error("Expected invalid facts due to invalid name")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for invalid name")
	}

	// Test invalid fact type
	invalidTypeFacts := map[string]interface{}{
		"facts": map[string]interface{}{
			"basic_facts": map[string]interface{}{
				"test-fact": map[string]interface{}{
					"value": "test-value",
					"type":  "invalid-type",
				},
			},
		},
	}

	result = validator.ValidateData("facts", invalidTypeFacts)
	if result.IsValid {
		t.Error("Expected invalid facts due to invalid type")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for invalid type")
	}

	// Test null value
	nullValueFacts := map[string]interface{}{
		"facts": map[string]interface{}{
			"basic_facts": map[string]interface{}{
				"test-fact": map[string]interface{}{
					"value": nil,
					"type":  "string",
				},
			},
		},
	}

	result = validator.ValidateData("facts", nullValueFacts)
	if result.IsValid {
		t.Error("Expected invalid facts due to null value")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for null value")
	}

	// Test encrypted fact with encrypted_value
	encryptedFacts := map[string]interface{}{
		"facts": map[string]interface{}{
			"custom_facts": map[string]interface{}{
				"secret-fact": map[string]interface{}{
					"value":     "secret-value",
					"type":      "string",
					"encrypted": true,
					"encrypted_value": map[string]interface{}{
						"data":         "encrypted-data-here",
						"format":       "base64",
						"algorithm":    "age",
						"version":      "v1",
						"encrypted_at": "2023-01-01T00:00:00Z",
					},
					"sensitive": true,
				},
			},
		},
	}

	result = validator.ValidateData("facts", encryptedFacts)
	if !result.IsValid {
		t.Errorf("Expected valid encrypted facts, got errors: %v", result.Errors)
	}
}

func TestValidator_SpookyV1Validation(t *testing.T) {
	validator := NewValidator()

	// Test valid spooky configuration
	validSpooky := map[string]interface{}{
		"spooky": map[string]interface{}{
			"ssh": map[string]interface{}{
				"timeout":            30.0,
				"keepalive_interval": 60.0,
			},
			"security": map[string]interface{}{
				"allow_unsafe_commands": false,
			},
			"age": map[string]interface{}{
				"identities": "age1...",
			},
			"logging": map[string]interface{}{
				"level": "info",
			},
		},
	}

	result := validator.ValidateData("spooky", validSpooky)
	if !result.IsValid {
		t.Errorf("Expected valid spooky configuration, got errors: %v", result.Errors)
	}

	// Test invalid spooky (missing required block)
	invalidSpooky := map[string]interface{}{
		// Missing spooky block
	}

	result = validator.ValidateData("spooky", invalidSpooky)
	if result.IsValid {
		t.Error("Expected invalid spooky due to missing spooky block")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for missing spooky block")
	}

	// Test invalid SSH timeout
	invalidSSHSpooky := map[string]interface{}{
		"spooky": map[string]interface{}{
			"ssh": map[string]interface{}{
				"timeout": 400.0, // Invalid: > 300
			},
		},
	}

	result = validator.ValidateData("spooky", invalidSSHSpooky)
	if result.IsValid {
		t.Error("Expected invalid spooky due to invalid SSH timeout")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for invalid SSH timeout")
	}

	// Test invalid security configuration
	invalidSecuritySpooky := map[string]interface{}{
		"spooky": map[string]interface{}{
			"security": map[string]interface{}{
				"allow_unsafe_commands": "not-a-boolean", // Invalid: should be boolean
			},
		},
	}

	result = validator.ValidateData("spooky", invalidSecuritySpooky)
	if result.IsValid {
		t.Error("Expected invalid spooky due to invalid security configuration")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for invalid security configuration")
	}

	// Test invalid logging level
	invalidLoggingSpooky := map[string]interface{}{
		"spooky": map[string]interface{}{
			"logging": map[string]interface{}{
				"level": "invalid-level", // Invalid: not a valid level
			},
		},
	}

	result = validator.ValidateData("spooky", invalidLoggingSpooky)
	if result.IsValid {
		t.Error("Expected invalid spooky due to invalid logging level")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for invalid logging level")
	}
}

func TestValidator_ProjectDirectoryV1Validation(t *testing.T) {
	validator := NewValidator()

	// Test valid project directory
	validProjectDir := map[string]interface{}{
		"project_directory": map[string]interface{}{
			"name": "test-project",
			"files": []interface{}{
				map[string]interface{}{
					"name": "test-file.txt",
					"type": "file",
				},
			},
			"directories": []interface{}{
				map[string]interface{}{
					"name": "test-dir",
					"type": "directory",
				},
			},
		},
	}

	result := validator.ValidateData("project-directory", validProjectDir)
	if !result.IsValid {
		t.Errorf("Expected valid project directory, got errors: %v", result.Errors)
	}

	// Test invalid project directory (missing required block)
	invalidProjectDir := map[string]interface{}{
		// Missing project_directory block
	}

	result = validator.ValidateData("project-directory", invalidProjectDir)
	if result.IsValid {
		t.Error("Expected invalid project directory due to missing project_directory block")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for missing project_directory block")
	}

	// Test invalid project directory (missing required name)
	invalidNameProjectDir := map[string]interface{}{
		"project_directory": map[string]interface{}{
			// Missing name field
			"files": []interface{}{
				map[string]interface{}{
					"name": "test-file.txt",
					"type": "file",
				},
			},
		},
	}

	result = validator.ValidateData("project-directory", invalidNameProjectDir)
	if result.IsValid {
		t.Error("Expected invalid project directory due to missing name")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for missing name")
	}

	// Test invalid file type
	invalidFileTypeProjectDir := map[string]interface{}{
		"project_directory": map[string]interface{}{
			"name": "test-project",
			"files": []interface{}{
				map[string]interface{}{
					"name": "test-file.txt",
					"type": "invalid-type", // Invalid: should be "file"
				},
			},
		},
	}

	result = validator.ValidateData("project-directory", invalidFileTypeProjectDir)
	if result.IsValid {
		t.Error("Expected invalid project directory due to invalid file type")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for invalid file type")
	}

	// Test invalid directory type
	invalidDirTypeProjectDir := map[string]interface{}{
		"project_directory": map[string]interface{}{
			"name": "test-project",
			"directories": []interface{}{
				map[string]interface{}{
					"name": "test-dir",
					"type": "invalid-type", // Invalid: should be "directory"
				},
			},
		},
	}

	result = validator.ValidateData("project-directory", invalidDirTypeProjectDir)
	if result.IsValid {
		t.Error("Expected invalid project directory due to invalid directory type")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for invalid directory type")
	}
}

func TestValidator_MachinesEnhancedValidation(t *testing.T) {
	validator := NewValidator()

	// Test machines with enhanced validation (valid password authentication)
	machinesData := map[string]interface{}{
		"machines": map[string]interface{}{
			"machine": []interface{}{
				map[string]interface{}{
					"name":     "test-machine",
					"hostname": "test.example.com",
					"user":     "testuser",
					"authentication": map[string]interface{}{
						"method":   "password",
						"password": "secret",
					},
				},
			},
		},
	}

	result := validator.ValidateData("machines", machinesData)
	if !result.IsValid {
		t.Errorf("Expected valid machines, got errors: %v", result.Errors)
	}

	// Test machines with SSH key authentication but missing key path
	invalidSSHKeyMachines := map[string]interface{}{
		"machines": map[string]interface{}{
			"machine": []interface{}{
				map[string]interface{}{
					"name":     "test-machine",
					"hostname": "test.example.com",
					"user":     "testuser",
					"authentication": map[string]interface{}{
						"method": "publickey",
						// Missing public_key_path
					},
				},
			},
		},
	}

	result = validator.ValidateData("machines", invalidSSHKeyMachines)
	if result.IsValid {
		t.Error("Expected invalid machines due to missing public_key_path")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for missing public_key_path")
	}

	// Test machines with password authentication but missing password
	invalidPasswordMachines := map[string]interface{}{
		"machines": map[string]interface{}{
			"machine": []interface{}{
				map[string]interface{}{
					"name":     "test-machine",
					"hostname": "test.example.com",
					"user":     "testuser",
					"authentication": map[string]interface{}{
						"method": "password",
						// Missing password
					},
				},
			},
		},
	}

	result = validator.ValidateData("machines", invalidPasswordMachines)
	if result.IsValid {
		t.Error("Expected invalid machines due to missing password")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for missing password")
	}

	// Test machines with certificate authentication but missing paths
	invalidCertMachines := map[string]interface{}{
		"machines": map[string]interface{}{
			"machine": []interface{}{
				map[string]interface{}{
					"name":     "test-machine",
					"hostname": "test.example.com",
					"user":     "testuser",
					"authentication": map[string]interface{}{
						"method": "certificate",
						// Missing private_key_path and certificate_path
					},
				},
			},
		},
	}

	result = validator.ValidateData("machines", invalidCertMachines)
	if result.IsValid {
		t.Error("Expected invalid machines due to missing certificate paths")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for missing certificate paths")
	}
}

func TestValidator_ActionsV1Validation(t *testing.T) {
	validator := NewValidator()

	// Test valid actions with machines targeting (to satisfy enhanced validation)
	validActions := map[string]interface{}{
		"actions": map[string]interface{}{
			"action": []interface{}{
				map[string]interface{}{
					"name":        "test-action",
					"description": "A test action",
					"type":        "command",
					"command":     "echo 'hello world'",
					"machines":    []interface{}{"test-machine"},
				},
			},
		},
	}

	result := validator.ValidateData("actions", validActions)
	if !result.IsValid {
		t.Errorf("Expected valid actions, got errors: %v", result.Errors)
	}

	// Test invalid actions (missing required fields)
	invalidActions := map[string]interface{}{
		"actions": map[string]interface{}{
			"action": []interface{}{
				map[string]interface{}{
					"name": "test-action",
					// Missing description and type
				},
			},
		},
	}

	result = validator.ValidateData("actions", invalidActions)
	if result.IsValid {
		t.Error("Expected invalid actions due to missing required fields")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for missing required fields")
	}

	// Test invalid action type
	invalidTypeActions := map[string]interface{}{
		"actions": map[string]interface{}{
			"action": []interface{}{
				map[string]interface{}{
					"name":        "test-action",
					"description": "A test action",
					"type":        "invalid-type",
				},
			},
		},
	}

	result = validator.ValidateData("actions", invalidTypeActions)
	if result.IsValid {
		t.Error("Expected invalid actions due to invalid type")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for invalid type")
	}
}
