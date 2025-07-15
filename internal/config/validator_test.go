package config

import (
	"testing"
)

func TestValidateConfig(t *testing.T) {
	testCases := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1", Password: "pass1"},
				},
				Actions: []Action{
					{Name: "action1", Command: "cmd1"},
				},
			},
			wantErr: false,
		},
		{
			name: "no servers",
			config: &Config{
				Servers: []Server{},
				Actions: []Action{
					{Name: "action1", Command: "cmd1"},
				},
			},
			wantErr: true,
		},
		{
			name: "server without name",
			config: &Config{
				Servers: []Server{
					{Host: "host1", User: "user1", Password: "pass1"},
				},
			},
			wantErr: true,
		},
		{
			name: "server without host",
			config: &Config{
				Servers: []Server{
					{Name: "server1", User: "user1", Password: "pass1"},
				},
			},
			wantErr: true,
		},
		{
			name: "server without user",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", Password: "pass1"},
				},
			},
			wantErr: true,
		},
		{
			name: "server without authentication",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1"},
				},
			},
			wantErr: true,
		},
		{
			name: "duplicate server names",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1", Password: "pass1"},
					{Name: "server1", Host: "host2", User: "user2", Password: "pass2"},
				},
			},
			wantErr: true,
		},
		{
			name: "action without name",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1", Password: "pass1"},
				},
				Actions: []Action{
					{Command: "cmd1"},
				},
			},
			wantErr: true,
		},
		{
			name: "action without command or script",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1", Password: "pass1"},
				},
				Actions: []Action{
					{Name: "action1"},
				},
			},
			wantErr: true,
		},
		{
			name: "action with both command and script",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1", Password: "pass1"},
				},
				Actions: []Action{
					{Name: "action1", Command: "cmd1", Script: "script1"},
				},
			},
			wantErr: true,
		},
		{
			name: "duplicate action names",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1", Password: "pass1"},
				},
				Actions: []Action{
					{Name: "action1", Command: "cmd1"},
					{Name: "action1", Script: "script1"},
				},
			},
			wantErr: true,
		},
		{
			name: "server with key file only",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1", KeyFile: "/path/to/key"},
				},
			},
			wantErr: false,
		},
		{
			name: "server with password and key file",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1", Password: "pass1", KeyFile: "/path/to/key"},
				},
			},
			wantErr: false,
		},
		{
			name: "action with script only",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1", Password: "pass1"},
				},
				Actions: []Action{
					{Name: "action1", Script: "script1"},
				},
			},
			wantErr: false,
		},
		{
			name: "server with default port",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1", Password: "pass1", Port: 0}, // Will default to 22
				},
			},
			wantErr: false,
		},
		{
			name: "server with custom port",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1", Password: "pass1", Port: 2222},
				},
			},
			wantErr: false,
		},
		{
			name: "action with timeout",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1", Password: "pass1"},
				},
				Actions: []Action{
					{Name: "action1", Command: "cmd1", Timeout: 60},
				},
			},
			wantErr: false,
		},
		{
			name: "action with description",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1", Password: "pass1"},
				},
				Actions: []Action{
					{Name: "action1", Command: "cmd1", Description: "Test action"},
				},
			},
			wantErr: false,
		},
		{
			name: "server with tags",
			config: &Config{
				Servers: []Server{
					{
						Name:     "server1",
						Host:     "host1",
						User:     "user1",
						Password: "pass1",
						Tags:     map[string]string{"env": "prod", "region": "us-west"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "action with servers list",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1", Password: "pass1"},
					{Name: "server2", Host: "host2", User: "user2", Password: "pass2"},
				},
				Actions: []Action{
					{Name: "action1", Command: "cmd1", Servers: []string{"server1"}},
				},
			},
			wantErr: false,
		},
		{
			name: "action with tags",
			config: &Config{
				Servers: []Server{
					{
						Name:     "server1",
						Host:     "host1",
						User:     "user1",
						Password: "pass1",
						Tags:     map[string]string{"env": "prod"},
					},
				},
				Actions: []Action{
					{Name: "action1", Command: "cmd1", Tags: []string{"env"}},
				},
			},
			wantErr: false,
		},
		{
			name: "action with parallel flag",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1", Password: "pass1"},
				},
				Actions: []Action{
					{Name: "action1", Command: "cmd1", Parallel: true},
				},
			},
			wantErr: false,
		},
		{
			name: "empty actions list",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1", Password: "pass1"},
				},
				Actions: []Action{},
			},
			wantErr: false, // Actions are optional
		},
		{
			name: "nil actions list",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1", Password: "pass1"},
				},
				Actions: nil,
			},
			wantErr: false, // Actions are optional
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateConfig(tc.config)
			if tc.wantErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateConfig_ErrorMessages(t *testing.T) {
	testCases := []struct {
		name          string
		config        *Config
		expectedError string
	}{
		{
			name: "no servers error message",
			config: &Config{
				Servers: []Server{},
			},
			expectedError: "at least one server must be defined",
		},
		{
			name: "server without name error message",
			config: &Config{
				Servers: []Server{
					{Host: "host1", User: "user1", Password: "pass1"},
				},
			},
			expectedError: "server name cannot be empty",
		},
		{
			name: "server without host error message",
			config: &Config{
				Servers: []Server{
					{Name: "server1", User: "user1", Password: "pass1"},
				},
			},
			expectedError: "server host cannot be empty for server server1",
		},
		{
			name: "server without user error message",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", Password: "pass1"},
				},
			},
			expectedError: "server user cannot be empty for server server1",
		},
		{
			name: "server without authentication error message",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1"},
				},
			},
			expectedError: "either password or key_file must be specified for server server1",
		},
		{
			name: "duplicate server names error message",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1", Password: "pass1"},
					{Name: "server1", Host: "host2", User: "user2", Password: "pass2"},
				},
			},
			expectedError: "duplicate server name: server1",
		},
		{
			name: "action without name error message",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1", Password: "pass1"},
				},
				Actions: []Action{
					{Command: "cmd1"},
				},
			},
			expectedError: "action name cannot be empty",
		},
		{
			name: "action without command or script error message",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1", Password: "pass1"},
				},
				Actions: []Action{
					{Name: "action1"},
				},
			},
			expectedError: "either command or script must be specified for action action1",
		},
		{
			name: "action with both command and script error message",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1", Password: "pass1"},
				},
				Actions: []Action{
					{Name: "action1", Command: "cmd1", Script: "script1"},
				},
			},
			expectedError: "cannot specify both command and script for action action1",
		},
		{
			name: "duplicate action names error message",
			config: &Config{
				Servers: []Server{
					{Name: "server1", Host: "host1", User: "user1", Password: "pass1"},
				},
				Actions: []Action{
					{Name: "action1", Command: "cmd1"},
					{Name: "action1", Script: "script1"},
				},
			},
			expectedError: "duplicate action name: action1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateConfig(tc.config)
			if err == nil {
				t.Error("expected error but got none")
				return
			}
			if !contains(err.Error(), tc.expectedError) {
				t.Errorf("expected error to contain '%s', got: %v", tc.expectedError, err)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
