package config

import (
	"strings"
	"testing"
)

func TestMachine_Validate(t *testing.T) {
	tests := []struct {
		name    string
		machine Machine
		wantErr bool
	}{
		{
			name: "valid machine with password",
			machine: Machine{
				Name:     "test-machine",
				Host:     "localhost",
				User:     "testuser",
				Password: "testpass",
			},
			wantErr: false,
		},
		{
			name: "valid machine with key file",
			machine: Machine{
				Name:    "test-machine",
				Host:    "localhost",
				User:    "testuser",
				KeyFile: "/path/to/key",
			},
			wantErr: false,
		},
		{
			name: "valid machine with both password and key file",
			machine: Machine{
				Name:     "test-machine",
				Host:     "localhost",
				User:     "testuser",
				Password: "testpass",
				KeyFile:  "/path/to/key",
			},
			wantErr: false,
		},
		{
			name: "invalid machine with no authentication",
			machine: Machine{
				Name: "test-machine",
				Host: "localhost",
				User: "testuser",
				// No password or key file
			},
			wantErr: true,
		},
		{
			name: "invalid machine with empty password and key file",
			machine: Machine{
				Name:     "test-machine",
				Host:     "localhost",
				User:     "testuser",
				Password: "",
				KeyFile:  "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			err := v.ValidateMachine(&tt.machine)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMachine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				// Check that the error message contains the machine name
				if !strings.Contains(err.Error(), tt.machine.Name) {
					t.Errorf("ValidateMachine() error message should contain machine name '%s', got: %s", tt.machine.Name, err.Error())
				}
			}
		})
	}
}

func TestAction_Validate(t *testing.T) {
	tests := []struct {
		name    string
		action  Action
		wantErr bool
	}{
		{
			name: "valid action with command",
			action: Action{
				Name:    "test-action",
				Command: "echo test",
			},
			wantErr: false,
		},
		{
			name: "valid action with script",
			action: Action{
				Name:   "test-action",
				Script: "/path/to/script.sh",
			},
			wantErr: false,
		},
		{
			name: "valid action with description and command",
			action: Action{
				Name:        "test-action",
				Description: "Test action description",
				Command:     "echo test",
			},
			wantErr: false,
		},
		{
			name: "valid action with tags and command",
			action: Action{
				Name:    "test-action",
				Command: "echo test",
				Tags:    []string{"production", "deploy"},
			},
			wantErr: false,
		},
		{
			name: "valid action with machines and command",
			action: Action{
				Name:     "test-action",
				Command:  "echo test",
				Machines: []string{"machine1", "machine2"},
			},
			wantErr: false,
		},
		{
			name: "valid action with timeout and command",
			action: Action{
				Name:    "test-action",
				Command: "echo test",
				Timeout: 30,
			},
			wantErr: false,
		},
		{
			name: "valid action with parallel and command",
			action: Action{
				Name:     "test-action",
				Command:  "echo test",
				Parallel: true,
			},
			wantErr: false,
		},
		{
			name: "invalid action with no command or script",
			action: Action{
				Name: "test-action",
				// No command or script
			},
			wantErr: true,
		},
		{
			name: "invalid action with empty command and script",
			action: Action{
				Name:    "test-action",
				Command: "",
				Script:  "",
			},
			wantErr: true,
		},
		{
			name: "invalid action with both command and script",
			action: Action{
				Name:    "test-action",
				Command: "echo test",
				Script:  "/path/to/script.sh",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			err := v.ValidateAction(&tt.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				// Check that the error message contains the action name
				if !strings.Contains(err.Error(), tt.action.Name) {
					t.Errorf("ValidateAction() error message should contain action name '%s', got: %s", tt.action.Name, err.Error())
				}
			}
		})
	}
}

func TestConfig_Structure(t *testing.T) {
	// Test that Config struct can be created with valid data
	config := Config{
		Machines: []Machine{
			{
				Name:     "test-machine",
				Host:     "localhost",
				User:     "testuser",
				Password: "testpass",
			},
		},
		Actions: []Action{
			{
				Name:    "test-action",
				Command: "echo test",
			},
		},
	}

	if len(config.Machines) != 1 {
		t.Errorf("Expected 1 machine, got %d", len(config.Machines))
	}

	if len(config.Actions) != 1 {
		t.Errorf("Expected 1 action, got %d", len(config.Actions))
	}

	if config.Machines[0].Name != "test-machine" {
		t.Errorf("Expected machine name 'test-machine', got '%s'", config.Machines[0].Name)
	}

	if config.Actions[0].Name != "test-action" {
		t.Errorf("Expected action name 'test-action', got '%s'", config.Actions[0].Name)
	}
}

func TestMachine_DefaultValues(t *testing.T) {
	machine := Machine{
		Name: "test-machine",
		Host: "localhost",
		User: "testuser",
	}

	// Test default values
	if machine.Port != 0 {
		t.Errorf("Expected default port 0, got %d", machine.Port)
	}

	if machine.Password != "" {
		t.Errorf("Expected default empty password, got '%s'", machine.Password)
	}

	if machine.KeyFile != "" {
		t.Errorf("Expected default empty key file, got '%s'", machine.KeyFile)
	}

	// Tags map is nil by default in Go
	if machine.Tags != nil {
		t.Errorf("Expected nil tags map by default, got %v", machine.Tags)
	}
}

func TestAction_DefaultValues(t *testing.T) {
	action := Action{
		Name: "test-action",
	}

	// Test default values
	if action.Description != "" {
		t.Errorf("Expected default empty description, got '%s'", action.Description)
	}

	if action.Command != "" {
		t.Errorf("Expected default empty command, got '%s'", action.Command)
	}

	if action.Script != "" {
		t.Errorf("Expected default empty script, got '%s'", action.Script)
	}

	// Slices are nil by default in Go
	if action.Machines != nil {
		t.Errorf("Expected nil machines slice by default, got %v", action.Machines)
	}

	if action.Tags != nil {
		t.Errorf("Expected nil tags slice by default, got %v", action.Tags)
	}

	if action.Timeout != 0 {
		t.Errorf("Expected default timeout 0, got %d", action.Timeout)
	}

	if action.Parallel {
		t.Errorf("Expected default parallel false, got true")
	}
}
