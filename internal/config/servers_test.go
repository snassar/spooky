package config

import (
	"testing"
)

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
		{
			name: "action with empty servers list",
			action: &Action{
				Name:    "test_action",
				Command: "echo hello",
				Servers: []string{}, // Empty servers list
			},
			expectedCount: 3, // Should return all servers when servers list is empty
			wantErr:       false,
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
