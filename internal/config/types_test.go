package config

import (
	"strings"
	"testing"
)

func TestServer_Validate(t *testing.T) {
	tests := []struct {
		name    string
		server  Server
		wantErr bool
	}{
		{
			name: "valid server with password",
			server: Server{
				Name:     "test-server",
				Host:     "localhost",
				User:     "testuser",
				Password: "testpass",
			},
			wantErr: false,
		},
		{
			name: "valid server with key file",
			server: Server{
				Name:    "test-server",
				Host:    "localhost",
				User:    "testuser",
				KeyFile: "/path/to/key",
			},
			wantErr: false,
		},
		{
			name: "valid server with both password and key file",
			server: Server{
				Name:     "test-server",
				Host:     "localhost",
				User:     "testuser",
				Password: "testpass",
				KeyFile:  "/path/to/key",
			},
			wantErr: false,
		},
		{
			name: "invalid server with no authentication",
			server: Server{
				Name: "test-server",
				Host: "localhost",
				User: "testuser",
				// No password or key file
			},
			wantErr: true,
		},
		{
			name: "invalid server with empty password and key file",
			server: Server{
				Name:     "test-server",
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
			err := v.ValidateServer(&tt.server)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				// Check that the error message contains the server name
				if !strings.Contains(err.Error(), tt.server.Name) {
					t.Errorf("ValidateServer() error message should contain server name '%s', got: %s", tt.server.Name, err.Error())
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
			name: "valid action with servers and command",
			action: Action{
				Name:    "test-action",
				Command: "echo test",
				Servers: []string{"server1", "server2"},
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
		Servers: []Server{
			{
				Name:     "test-server",
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

	if len(config.Servers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(config.Servers))
	}

	if len(config.Actions) != 1 {
		t.Errorf("Expected 1 action, got %d", len(config.Actions))
	}

	if config.Servers[0].Name != "test-server" {
		t.Errorf("Expected server name 'test-server', got '%s'", config.Servers[0].Name)
	}

	if config.Actions[0].Name != "test-action" {
		t.Errorf("Expected action name 'test-action', got '%s'", config.Actions[0].Name)
	}
}

func TestServer_DefaultValues(t *testing.T) {
	server := Server{
		Name: "test-server",
		Host: "localhost",
		User: "testuser",
	}

	// Test default values
	if server.Port != 0 {
		t.Errorf("Expected default port 0, got %d", server.Port)
	}

	if server.Password != "" {
		t.Errorf("Expected default empty password, got '%s'", server.Password)
	}

	if server.KeyFile != "" {
		t.Errorf("Expected default empty key file, got '%s'", server.KeyFile)
	}

	// Tags map is nil by default in Go
	if server.Tags != nil {
		t.Errorf("Expected nil tags map by default, got %v", server.Tags)
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
	if action.Servers != nil {
		t.Errorf("Expected nil servers slice by default, got %v", action.Servers)
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
