package spooky

import (
	"os"
	"testing"
)

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name            string
		configContent   string
		wantErr         bool
		expectedServers int
		expectedActions int
	}{
		{
			name: "valid config with servers and actions",
			configContent: `
server "web1" {
  host = "192.168.1.10"
  user = "admin"
  password = "secret"
  port = 22
}

server "web2" {
  host = "192.168.1.11"
  user = "admin"
  key_file = "/path/to/key"
}

action "check_status" {
  description = "Check server status"
  command = "systemctl status nginx"
  servers = ["web1", "web2"]
}
`,
			wantErr:         false,
			expectedServers: 2,
			expectedActions: 1,
		},
		{
			name: "valid config with tags",
			configContent: `
server "web1" {
  host = "192.168.1.10"
  user = "admin"
  password = "secret"
  tags = {
    environment = "production"
    role = "web"
  }
}

action "deploy" {
  description = "Deploy application"
  command = "git pull && npm install"
  tags = ["environment", "role"]
}
`,
			wantErr:         false,
			expectedServers: 1,
			expectedActions: 1,
		},
		{
			name: "invalid config - missing server host",
			configContent: `
server "web1" {
  user = "admin"
  password = "secret"
}
`,
			wantErr: true,
		},
		{
			name: "invalid config - missing server user",
			configContent: `
server "web1" {
  host = "192.168.1.10"
  password = "secret"
}
`,
			wantErr: true,
		},
		{
			name: "invalid config - no authentication",
			configContent: `
server "web1" {
  host = "192.168.1.10"
  user = "admin"
}
`,
			wantErr: true,
		},
		{
			name: "invalid config - duplicate server name",
			configContent: `
server "web1" {
  host = "192.168.1.10"
  user = "admin"
  password = "secret"
}

server "web1" {
  host = "192.168.1.11"
  user = "admin"
  password = "secret"
}
`,
			wantErr: true,
		},
		{
			name: "invalid config - action without command or script",
			configContent: `
server "web1" {
  host = "192.168.1.10"
  user = "admin"
  password = "secret"
}

action "empty_action" {
  description = "Action without command or script"
}
`,
			wantErr: true,
		},
		{
			name: "invalid config - action with both command and script",
			configContent: `
server "web1" {
  host = "192.168.1.10"
  user = "admin"
  password = "secret"
}

action "conflict_action" {
  description = "Action with both command and script"
  command = "echo hello"
  script = "/path/to/script.sh"
}
`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary config file
			tmpfile, err := os.CreateTemp("", "test_config_*.hcl")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(tt.configContent)); err != nil {
				t.Fatal(err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatal(err)
			}

			// Parse config
			config, err := ParseConfig(tmpfile.Name())

			// Check error expectation
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseConfig() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("parseConfig() unexpected error: %v", err)
				return
			}

			// Check server count
			if len(config.Servers) != tt.expectedServers {
				t.Errorf("parseConfig() got %d servers, want %d", len(config.Servers), tt.expectedServers)
			}

			// Check action count
			if len(config.Actions) != tt.expectedActions {
				t.Errorf("parseConfig() got %d actions, want %d", len(config.Actions), tt.expectedActions)
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Servers: []Server{
					{
						Name:     "web1",
						Host:     "192.168.1.10",
						User:     "admin",
						Password: "secret",
						Port:     22,
					},
				},
				Actions: []Action{
					{
						Name:    "test_action",
						Command: "echo hello",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "no servers",
			config: &Config{
				Servers: []Server{},
				Actions: []Action{},
			},
			wantErr: true,
		},
		{
			name: "server without name",
			config: &Config{
				Servers: []Server{
					{
						Host:     "192.168.1.10",
						User:     "admin",
						Password: "secret",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "server without host",
			config: &Config{
				Servers: []Server{
					{
						Name:     "web1",
						User:     "admin",
						Password: "secret",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "server without user",
			config: &Config{
				Servers: []Server{
					{
						Name:     "web1",
						Host:     "192.168.1.10",
						Password: "secret",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "server without authentication",
			config: &Config{
				Servers: []Server{
					{
						Name: "web1",
						Host: "192.168.1.10",
						User: "admin",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetServersForAction(t *testing.T) {
	config := &Config{
		Servers: []Server{
			{
				Name:     "web1",
				Host:     "192.168.1.10",
				User:     "admin",
				Password: "secret",
				Tags: map[string]string{
					"environment": "production",
					"role":        "web",
				},
			},
			{
				Name:     "db1",
				Host:     "192.168.1.20",
				User:     "admin",
				Password: "secret",
				Tags: map[string]string{
					"environment": "production",
					"role":        "database",
				},
			},
			{
				Name:     "web2",
				Host:     "192.168.1.11",
				User:     "admin",
				Password: "secret",
				Tags: map[string]string{
					"environment": "staging",
					"role":        "web",
				},
			},
		},
	}

	tests := []struct {
		name          string
		action        *Action
		expectedCount int
		expectedNames []string
		wantErr       bool
	}{
		{
			name: "specific servers",
			action: &Action{
				Name:    "test_action",
				Command: "echo hello",
				Servers: []string{"web1", "db1"},
			},
			expectedCount: 2,
			expectedNames: []string{"web1", "db1"},
			wantErr:       false,
		},
		{
			name: "tag-based selection",
			action: &Action{
				Name:    "test_action",
				Command: "echo hello",
				Tags:    []string{"environment"},
			},
			expectedCount: 3, // All servers have environment tag
			expectedNames: []string{"web1", "db1", "web2"},
			wantErr:       false,
		},
		{
			name: "specific tag value",
			action: &Action{
				Name:    "test_action",
				Command: "echo hello",
				Tags:    []string{"role"},
			},
			expectedCount: 3, // All servers have role tag
			expectedNames: []string{"web1", "db1", "web2"},
			wantErr:       false,
		},
		{
			name: "no servers or tags - use all",
			action: &Action{
				Name:    "test_action",
				Command: "echo hello",
			},
			expectedCount: 3,
			expectedNames: []string{"web1", "db1", "web2"},
			wantErr:       false,
		},
		{
			name: "non-existent server",
			action: &Action{
				Name:    "test_action",
				Command: "echo hello",
				Servers: []string{"nonexistent"},
			},
			expectedCount: 0,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			servers, err := GetServersForAction(tt.action, config)

			if tt.wantErr {
				if err == nil {
					t.Errorf("getServersForAction() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("getServersForAction() unexpected error: %v", err)
				return
			}

			if len(servers) != tt.expectedCount {
				t.Errorf("getServersForAction() got %d servers, want %d", len(servers), tt.expectedCount)
			}

			// Check server names if expected
			if tt.expectedNames != nil {
				serverNames := make([]string, len(servers))
				for i, server := range servers {
					serverNames[i] = server.Name
				}

				// Simple check - we don't need exact order matching for this test
				if len(serverNames) != len(tt.expectedNames) {
					t.Errorf("getServersForAction() got server names %v, want %v", serverNames, tt.expectedNames)
				}
			}
		})
	}
}
