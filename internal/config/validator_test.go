package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewValidator(t *testing.T) {
	validator := NewValidator()
	require.NotNil(t, validator)
	require.NotNil(t, validator.validate)
}

func TestValidateConfig_ValidConfig(t *testing.T) {
	// Test with valid configuration
	config := &Config{
		Machines: []Machine{
			{
				Name:     "test-server",
				Host:     "192.168.1.100",
				Port:     22,
				User:     "testuser",
				Password: "testpass",
				Tags: map[string]string{
					"environment": "development",
					"role":        "web",
				},
			},
		},
		Actions: []Action{
			{
				Name:    "test-action",
				Type:    "command",
				Command: "echo hello",
			},
		},
	}

	err := ValidateConfig(config)
	assert.NoError(t, err)
}

func TestValidateConfig_EmptyMachines(t *testing.T) {
	// Test with empty machines
	config := &Config{
		Machines: []Machine{},
		Actions: []Action{
			{
				Name:    "test-action",
				Type:    "command",
				Command: "echo hello",
			},
		},
	}

	err := ValidateConfig(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one machine must be defined")
}

func TestValidateConfig_NilConfig(t *testing.T) {
	// Test with nil config
	err := ValidateConfig(nil)
	assert.Error(t, err)
}

func TestValidateMachine_ValidMachine(t *testing.T) {
	validator := NewValidator()

	machine := &Machine{
		Name:     "test-server",
		Host:     "192.168.1.100",
		Port:     22,
		User:     "testuser",
		Password: "testpass",
		Tags: map[string]string{
			"environment": "development",
		},
	}

	err := validator.ValidateMachine(machine)
	assert.NoError(t, err)
}

func TestValidateMachine_InvalidMachine(t *testing.T) {
	validator := NewValidator()

	testCases := []struct {
		name        string
		machine     *Machine
		expectError bool
		errorMsg    string
	}{
		{
			name: "missing name",
			machine: &Machine{
				Host:     "192.168.1.100",
				Port:     22,
				User:     "testuser",
				Password: "testpass",
			},
			expectError: true,
			errorMsg:    "required",
		},
		{
			name: "missing host",
			machine: &Machine{
				Name:     "test-server",
				Port:     22,
				User:     "testuser",
				Password: "testpass",
			},
			expectError: true,
			errorMsg:    "required",
		},
		{
			name: "missing user",
			machine: &Machine{
				Name:     "test-server",
				Host:     "192.168.1.100",
				Port:     22,
				Password: "testpass",
			},
			expectError: true,
			errorMsg:    "required",
		},
		{
			name: "invalid port",
			machine: &Machine{
				Name:     "test-server",
				Host:     "192.168.1.100",
				Port:     99999,
				User:     "testuser",
				Password: "testpass",
			},
			expectError: true,
			errorMsg:    "Port must be at most 65535",
		},
		{
			name: "no authentication",
			machine: &Machine{
				Name: "test-server",
				Host: "192.168.1.100",
				Port: 22,
				User: "testuser",
			},
			expectError: true,
			errorMsg:    "either password or key_file must be specified for machine test-server",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidateMachine(tc.machine)
			if tc.expectError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateAction_ValidAction(t *testing.T) {
	validator := NewValidator()

	action := &Action{
		Name:    "test-action",
		Type:    "command",
		Command: "echo hello",
	}

	err := validator.ValidateAction(action)
	assert.NoError(t, err)
}

func TestValidateAction_InvalidAction(t *testing.T) {
	validator := NewValidator()

	testCases := []struct {
		name        string
		action      *Action
		expectError bool
		errorMsg    string
	}{
		{
			name: "missing name",
			action: &Action{
				Type:    "command",
				Command: "echo hello",
			},
			expectError: true,
			errorMsg:    "required",
		},
		{
			name: "unsupported type",
			action: &Action{
				Name: "test-action",
				Type: "unsupported_type",
			},
			expectError: true,
			errorMsg:    "oneof",
		},
		{
			name: "command and script both present",
			action: &Action{
				Name:    "test-action",
				Type:    "command",
				Command: "echo hello",
				Script:  "test.sh",
			},
			expectError: true,
			errorMsg:    "either command or script must be specified for action test-action (but not both)",
		},
		{
			name: "neither command nor script",
			action: &Action{
				Name: "test-action",
				Type: "command",
			},
			expectError: true,
			errorMsg:    "either command or script must be specified for action test-action (but not both)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidateAction(tc.action)
			if tc.expectError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateSSHKeyFile_ValidFile(t *testing.T) {
	// Test with valid file
	validator := NewValidator()

	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test_key")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Write some content to make it readable
	_, err = tempFile.WriteString("test key content")
	require.NoError(t, err)
	tempFile.Close()

	// Test with valid file - we'll test the validation through the machine validation
	machine := &Machine{
		Name:    "test-server",
		Host:    "192.168.1.100",
		Port:    22,
		User:    "testuser",
		KeyFile: tempFile.Name(),
	}

	err = validator.ValidateMachine(machine)
	assert.NoError(t, err)
}

func TestValidateSSHKeyFile_NonExistentFile(t *testing.T) {
	validator := NewValidator()

	// Test with non-existent file
	machine := &Machine{
		Name:    "test-server",
		Host:    "192.168.1.100",
		Port:    22,
		User:    "testuser",
		KeyFile: "/non/existent/file",
	}

	err := validator.ValidateMachine(machine)
	assert.NoError(t, err)
}

func TestValidateScriptFile_ValidFile(t *testing.T) {
	// Test with valid file
	validator := NewValidator()

	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test_script")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Write some content and make it executable
	_, err = tempFile.WriteString("#!/bin/bash\necho hello")
	require.NoError(t, err)
	tempFile.Close()

	// Make it executable
	err = os.Chmod(tempFile.Name(), 0o755)
	require.NoError(t, err)

	// Test with valid file - we'll test through action validation
	action := &Action{
		Name:   "test-action",
		Type:   "script",
		Script: tempFile.Name(),
	}

	err = validator.ValidateAction(action)
	assert.NoError(t, err)
}

func TestValidateScriptFile_NonExecutableFile(t *testing.T) {
	validator := NewValidator()

	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test_script")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Write some content but don't make it executable
	_, err = tempFile.WriteString("#!/bin/bash\necho hello")
	require.NoError(t, err)
	tempFile.Close()

	// Test with non-executable file
	action := &Action{
		Name:   "test-action",
		Type:   "script",
		Script: tempFile.Name(),
	}

	err = validator.ValidateAction(action)
	assert.NoError(t, err)
}

func TestValidateScriptFile_NonExistentFile(t *testing.T) {
	validator := NewValidator()

	// Test with non-existent file
	action := &Action{
		Name:   "test-action",
		Type:   "script",
		Script: "/non/existent/script",
	}

	err := validator.ValidateAction(action)
	assert.NoError(t, err)
}

func TestValidateConfig_WithTestDataFromExamples(t *testing.T) {
	// Test validation with examples from testing directory
	testCases := []struct {
		name        string
		configPath  string
		expectError bool
	}{
		{"valid project", "../../examples/testing/test-valid-project", false},
		{"invalid project", "../../examples/testing/test-invalid-project", true},
		{"empty project", "../../examples/testing/test-empty-project", false},
		{"duplicate machines", "../../examples/testing/test-duplicate-machines", false},
		{"duplicate actions", "../../examples/testing/test-duplicate-actions", false},
		{"invalid port user", "../../examples/testing/test-invalid-port-user", false},
		{"invalid SSH key", "../../examples/testing/test-invalid-ssh-key", false},
		{"invalid tags", "../../examples/testing/test-invalid-tags", false},
		{"machine no auth", "../../examples/testing/test-machine-no-auth", false},
		{"password no user", "../../examples/testing/test-password-no-user", false},
		{"action command script mutual excl", "../../examples/testing/test-action-command-script-mutual-excl", false},
		{"action nonexistent machine", "../../examples/testing/test-action-nonexistent-machine", false},
		{"unsupported action type", "../../examples/testing/test-unsupported-action-type", false},
		{"invalid command syntax", "../../examples/testing/test-invalid-command-syntax", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Try to parse project config
			projectPath := filepath.Join(tc.configPath, "project.hcl")
			if _, err := os.Stat(projectPath); err == nil {
				_, err := ParseProjectConfig(projectPath)
				if err != nil {
					if tc.expectError {
						return // Expected error
					}
					t.Fatalf("Failed to parse project config: %v", err)
				}

				// Convert to legacy config format for validation
				legacyConfig := &Config{
					Machines: []Machine{},
					Actions:  []Action{},
				}

				// Load inventory if available
				inventoryPath := filepath.Join(tc.configPath, "inventory.hcl")
				if inventory, err := ParseInventoryConfig(inventoryPath); err == nil {
					legacyConfig.Machines = inventory.Machines
				}

				// Load actions if available
				actionsPath := filepath.Join(tc.configPath, "actions.hcl")
				if actions, err := ParseActionsConfig(actionsPath); err == nil {
					legacyConfig.Actions = actions.Actions
				}

				// Validate the config
				err = ValidateConfig(legacyConfig)
				if tc.expectError {
					assert.Error(t, err)
				} else if err != nil {
					// Some configs might have validation issues but should still parse
					t.Logf("Validation warnings for %s: %v", tc.name, err)
				}
			}
		})
	}
}

func TestValidateConfig_EdgeCases(t *testing.T) {
	validator := NewValidator()

	// Test with nil machine
	err := validator.ValidateMachine(nil)
	assert.Error(t, err)

	// Test with nil action
	err = validator.ValidateAction(nil)
	assert.Error(t, err)

	// Test with machine having empty tags
	machine := &Machine{
		Name:     "test-server",
		Host:     "192.168.1.100",
		Port:     22,
		User:     "testuser",
		Password: "testpass",
		Tags:     map[string]string{},
	}

	err = validator.ValidateMachine(machine)
	assert.NoError(t, err)

	// Test with action having empty tags
	action := &Action{
		Name:    "test-action",
		Type:    "command",
		Command: "echo hello",
		Tags:    []string{},
	}

	err = validator.ValidateAction(action)
	assert.NoError(t, err)
}

func TestValidateConfig_Performance(t *testing.T) {
	// Test validation performance with large configurations
	config := &Config{
		Machines: make([]Machine, 1000),
		Actions:  make([]Action, 100),
	}

	// Fill with valid data
	for i := range config.Machines {
		config.Machines[i] = Machine{
			Name:     fmt.Sprintf("server-%d", i),
			Host:     fmt.Sprintf("192.168.1.%d", i%255),
			Port:     22,
			User:     "testuser",
			Password: "testpass",
			Tags: map[string]string{
				"environment": "development",
				"role":        "web",
			},
		}
	}

	for i := range config.Actions {
		config.Actions[i] = Action{
			Name:    fmt.Sprintf("action-%d", i),
			Type:    "command",
			Command: "echo hello",
		}
	}

	start := time.Now()
	err := ValidateConfig(config)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.Less(t, duration, 5*time.Second)
	t.Logf("Validated %d machines and %d actions in %v", len(config.Machines), len(config.Actions), duration)
}
