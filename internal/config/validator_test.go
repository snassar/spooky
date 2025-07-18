package config

import (
	"strings"
	"testing"
)

func TestNewValidator(t *testing.T) {
	v := NewValidator()
	if v == nil {
		t.Fatal("NewValidator() returned nil")
	}
	if v.validate == nil {
		t.Fatal("validator instance is nil")
	}
}

func TestValidator_ValidateServer(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name    string
		server  Server
		wantErr bool
		errMsg  string
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
			},
			wantErr: true,
			errMsg:  "either password or key_file must be specified for server test-server",
		},
		{
			name: "invalid server with empty name",
			server: Server{
				Name:     "",
				Host:     "localhost",
				User:     "testuser",
				Password: "testpass",
			},
			wantErr: true,
			errMsg:  "Name is required",
		},
		{
			name: "invalid server with empty host",
			server: Server{
				Name:     "test-server",
				Host:     "",
				User:     "testuser",
				Password: "testpass",
			},
			wantErr: true,
			errMsg:  "Host is required",
		},
		{
			name: "invalid server with empty user",
			server: Server{
				Name:     "test-server",
				Host:     "localhost",
				User:     "",
				Password: "testpass",
			},
			wantErr: true,
			errMsg:  "User is required",
		},

		{
			name: "valid server with valid port",
			server: Server{
				Name:     "test-server",
				Host:     "localhost",
				User:     "testuser",
				Password: "testpass",
				Port:     22,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateServer(&tt.server)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateServer() error message = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestValidator_ValidateAction(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name    string
		action  Action
		wantErr bool
		errMsg  string
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
			name: "invalid action with no command or script",
			action: Action{
				Name: "test-action",
			},
			wantErr: true,
			errMsg:  "either command or script must be specified for action test-action",
		},
		{
			name: "invalid action with both command and script",
			action: Action{
				Name:    "test-action",
				Command: "echo test",
				Script:  "/path/to/script.sh",
			},
			wantErr: true,
			errMsg:  "either command or script must be specified for action test-action",
		},
		{
			name: "invalid action with empty name",
			action: Action{
				Name:    "",
				Command: "echo test",
			},
			wantErr: true,
			errMsg:  "Name is required",
		},

		{
			name: "valid action with valid timeout",
			action: Action{
				Name:    "test-action",
				Command: "echo test",
				Timeout: 30,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateAction(&tt.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateAction() error message = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestValidator_ValidateConfig(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: Config{
				Servers: []Server{
					{
						Name:     "server1",
						Host:     "localhost",
						User:     "testuser",
						Password: "testpass",
					},
				},
				Actions: []Action{
					{
						Name:    "action1",
						Command: "echo test",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid config with no servers",
			config: Config{
				Servers: []Server{},
				Actions: []Action{
					{
						Name:    "action1",
						Command: "echo test",
					},
				},
			},
			wantErr: true,
			errMsg:  "at least one server must be defined",
		},
		{
			name: "invalid config with duplicate server names",
			config: Config{
				Servers: []Server{
					{
						Name:     "server1",
						Host:     "localhost",
						User:     "testuser",
						Password: "testpass",
					},
					{
						Name:     "server1",
						Host:     "localhost2",
						User:     "testuser2",
						Password: "testpass2",
					},
				},
				Actions: []Action{
					{
						Name:    "action1",
						Command: "echo test",
					},
				},
			},
			wantErr: true,
			errMsg:  "duplicate server name: server1",
		},
		{
			name: "invalid config with duplicate action names",
			config: Config{
				Servers: []Server{
					{
						Name:     "server1",
						Host:     "localhost",
						User:     "testuser",
						Password: "testpass",
					},
				},
				Actions: []Action{
					{
						Name:    "action1",
						Command: "echo test",
					},
					{
						Name:   "action1",
						Script: "/path/to/script.sh",
					},
				},
			},
			wantErr: true,
			errMsg:  "duplicate action name: action1",
		},
		{
			name: "invalid config with invalid server reference",
			config: Config{
				Servers: []Server{
					{
						Name:     "server1",
						Host:     "localhost",
						User:     "testuser",
						Password: "testpass",
					},
				},
				Actions: []Action{
					{
						Name:    "action1",
						Command: "echo test",
						Servers: []string{"server1", "nonexistent"},
					},
				},
			},
			wantErr: true,
			errMsg:  "server reference 'nonexistent' in action 'action1' does not exist",
		},
		{
			name: "valid config with valid server references",
			config: Config{
				Servers: []Server{
					{
						Name:     "server1",
						Host:     "localhost",
						User:     "testuser",
						Password: "testpass",
					},
					{
						Name:     "server2",
						Host:     "localhost2",
						User:     "testuser2",
						Password: "testpass2",
					},
				},
				Actions: []Action{
					{
						Name:    "action1",
						Command: "echo test",
						Servers: []string{"server1", "server2"},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.validateConfig(&tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateConfig() error message = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestValidator_StructLevelValidation(t *testing.T) {
	v := NewValidator()

	// Test server struct-level validation
	t.Run("Server struct-level validation", func(t *testing.T) {
		// Test authentication validation
		server := Server{
			Name: "test-server",
			Host: "localhost",
			User: "testuser",
			// No password or key file
		}
		err := v.ValidateServer(&server)
		if err == nil {
			t.Error("Expected error for server without authentication")
		}
		if !strings.Contains(err.Error(), "either password or key_file must be specified") {
			t.Errorf("Expected authentication error, got: %v", err)
		}

	})

	// Test action struct-level validation
	t.Run("Action struct-level validation", func(t *testing.T) {
		// Test execution validation
		action := Action{
			Name: "test-action",
			// No command or script
		}
		err := v.ValidateAction(&action)
		if err == nil {
			t.Error("Expected error for action without command or script")
		}
		if !strings.Contains(err.Error(), "either command or script must be specified") {
			t.Errorf("Expected execution error, got: %v", err)
		}

	})
}
