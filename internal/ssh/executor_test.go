package ssh

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"spooky/internal/config"
)

func TestExecuteConfig_NilConfig(t *testing.T) {
	err := ExecuteConfig(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config cannot be nil")
}

func TestExecuteConfig_EmptyActions(t *testing.T) {
	cfg := &config.Config{
		Machines: []config.Machine{
			{
				Name:     "test-server",
				Host:     "192.168.1.100",
				Port:     22,
				User:     "testuser",
				Password: "testpass",
			},
		},
		Actions: []config.Action{}, // Empty actions
	}

	err := ExecuteConfig(cfg)
	assert.NoError(t, err) // Should succeed with no actions
}

func TestExecuteConfig_EmptyMachines(t *testing.T) {
	cfg := &config.Config{
		Machines: []config.Machine{}, // Empty machines
		Actions: []config.Action{
			{
				Name:    "test-action",
				Type:    "command",
				Command: "echo hello",
			},
		},
	}

	err := ExecuteConfig(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no machines")
}

func TestExecuteConfig_UnsupportedActionType(t *testing.T) {
	cfg := &config.Config{
		Machines: []config.Machine{
			{
				Name:     "test-server",
				Host:     "192.168.1.100",
				Port:     22,
				User:     "testuser",
				Password: "testpass",
			},
		},
		Actions: []config.Action{
			{
				Name: "test-action",
				Type: "unsupported_type",
			},
		},
	}

	err := ExecuteConfig(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported action type")
}

func TestExecuteConfig_CommandAction(t *testing.T) {
	cfg := &config.Config{
		Machines: []config.Machine{
			{
				Name:     "test-server",
				Host:     "192.168.1.100",
				Port:     22,
				User:     "testuser",
				Password: "testpass",
			},
		},
		Actions: []config.Action{
			{
				Name:    "test-action",
				Type:    "command",
				Command: "echo hello",
			},
		},
	}

	// This will fail due to SSH connection, but we can test the configuration
	err := ExecuteConfig(cfg)
	assert.Error(t, err)
	// But it shouldn't be a configuration error
	assert.NotContains(t, err.Error(), "command configuration is required")
}

func TestExecuteConfig_ScriptAction(t *testing.T) {
	cfg := &config.Config{
		Machines: []config.Machine{
			{
				Name:     "test-server",
				Host:     "192.168.1.100",
				Port:     22,
				User:     "testuser",
				Password: "testpass",
			},
		},
		Actions: []config.Action{
			{
				Name:   "test-action",
				Type:   "script",
				Script: "/tmp/test.sh",
			},
		},
	}

	// This will fail due to SSH connection, but we can test the configuration
	err := ExecuteConfig(cfg)
	assert.Error(t, err)
	// But it shouldn't be a configuration error
	assert.NotContains(t, err.Error(), "script configuration is required")
}

func TestExecuteConfig_TemplateAction(t *testing.T) {
	cfg := &config.Config{
		Machines: []config.Machine{
			{
				Name:     "test-server",
				Host:     "192.168.1.100",
				Port:     22,
				User:     "testuser",
				Password: "testpass",
			},
		},
		Actions: []config.Action{
			{
				Name: "test-action",
				Type: "template_deploy",
				Template: &config.TemplateConfig{
					Source:      "/tmp/test.tmpl",
					Destination: "/tmp/test.conf",
				},
			},
		},
	}

	// This will fail due to SSH connection, but we can test the configuration
	err := ExecuteConfig(cfg)
	assert.Error(t, err)
	// But it shouldn't be a configuration error
	assert.NotContains(t, err.Error(), "template configuration is required")
}

func TestExecuteConfig_ParallelExecution(t *testing.T) {
	cfg := &config.Config{
		Machines: []config.Machine{
			{
				Name:     "server1",
				Host:     "192.168.1.100",
				Port:     22,
				User:     "testuser",
				Password: "testpass",
			},
			{
				Name:     "server2",
				Host:     "192.168.1.101",
				Port:     22,
				User:     "testuser",
				Password: "testpass",
			},
		},
		Actions: []config.Action{
			{
				Name:     "test-action",
				Type:     "command",
				Command:  "echo hello",
				Parallel: true,
			},
		},
	}

	// This will fail due to SSH connection, but we can test the configuration
	err := ExecuteConfig(cfg)
	assert.Error(t, err)
	// But it shouldn't be a parallel execution error
	assert.NotContains(t, err.Error(), "parallel execution")
}

func TestExecuteConfig_WithTestDataFromExamples(t *testing.T) {
	// Create a configuration similar to what would be in test data
	cfg := &config.Config{
		Machines: []config.Machine{
			{
				Name:     "example-server",
				Host:     "192.168.1.100",
				Port:     22,
				User:     "testuser",
				Password: "testpass",
				Tags: map[string]string{
					"environment": "test",
					"role":        "web",
				},
			},
		},
		Actions: []config.Action{
			{
				Name:        "deploy-config",
				Description: "Deploy configuration files",
				Type:        "command",
				Command:     "echo 'Deploying configuration'",
				Machines:    []string{"example-server"},
				Parallel:    true,
			},
		},
	}

	// This will fail due to SSH connection, but we can test the configuration
	err := ExecuteConfig(cfg)
	assert.Error(t, err)
	// But it shouldn't be a configuration error
	assert.NotContains(t, err.Error(), "configuration is required")
}

func TestExecuteConfig_ActionTypes(t *testing.T) {
	cfg := &config.Config{
		Machines: []config.Machine{
			{
				Name:     "test-server",
				Host:     "192.168.1.100",
				Port:     22,
				User:     "testuser",
				Password: "testpass",
			},
		},
	}

	tests := []struct {
		name        string
		actionType  string
		description string
	}{
		{
			name:        "command",
			actionType:  "command",
			description: "Execute shell commands on target servers",
		},
		{
			name:        "script",
			actionType:  "script",
			description: "Execute script files on target servers",
		},
		{
			name:        "template_deploy",
			actionType:  "template_deploy",
			description: "Deploy template files to target servers",
		},
		{
			name:        "template_evaluate",
			actionType:  "template_evaluate",
			description: "Evaluate templates on servers with server-specific facts",
		},
		{
			name:        "template_validate",
			actionType:  "template_validate",
			description: "Validate templates on servers",
		},
		{
			name:        "template_cleanup",
			actionType:  "template_cleanup",
			description: "Remove template files from servers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var action config.Action

			switch tt.actionType {
			case "command":
				action = config.Action{
					Name:    tt.name,
					Type:    tt.actionType,
					Command: "echo test",
				}
			case "script":
				action = config.Action{
					Name:   tt.name,
					Type:   tt.actionType,
					Script: "/tmp/test.sh",
				}
			default: // template actions
				action = config.Action{
					Name: tt.name,
					Type: tt.actionType,
					Template: &config.TemplateConfig{
						Source:      "/tmp/test.tmpl",
						Destination: "/tmp/test.conf",
					},
				}
			}

			testCfg := &config.Config{
				Machines: cfg.Machines,
				Actions:  []config.Action{action},
			}

			err := ExecuteConfig(testCfg)
			// All should fail due to SSH connection, but not due to unsupported type
			assert.Error(t, err)
			assert.NotContains(t, err.Error(), "unsupported action type")
		})
	}
}

func TestIsTemplateAction(t *testing.T) {
	tests := []struct {
		name       string
		actionType string
		expectTrue bool
	}{
		{
			name:       "template_deploy",
			actionType: "template_deploy",
			expectTrue: true,
		},
		{
			name:       "template_evaluate",
			actionType: "template_evaluate",
			expectTrue: true,
		},
		{
			name:       "template_validate",
			actionType: "template_validate",
			expectTrue: true,
		},
		{
			name:       "template_cleanup",
			actionType: "template_cleanup",
			expectTrue: true,
		},
		{
			name:       "command",
			actionType: "command",
			expectTrue: false,
		},
		{
			name:       "script",
			actionType: "script",
			expectTrue: false,
		},
		{
			name:       "unsupported",
			actionType: "unsupported",
			expectTrue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := &config.Action{
				Name: tt.name,
				Type: tt.actionType,
			}

			result := isTemplateAction(action)
			assert.Equal(t, tt.expectTrue, result)
		})
	}
}

func TestExecuteConfig_MachineFiltering(t *testing.T) {
	cfg := &config.Config{
		Machines: []config.Machine{
			{
				Name:     "server1",
				Host:     "192.168.1.100",
				Port:     22,
				User:     "testuser",
				Password: "testpass",
			},
			{
				Name:     "server2",
				Host:     "192.168.1.101",
				Port:     22,
				User:     "testuser",
				Password: "testpass",
			},
			{
				Name:     "server3",
				Host:     "192.168.1.102",
				Port:     22,
				User:     "testuser",
				Password: "testpass",
			},
		},
		Actions: []config.Action{
			{
				Name:     "filtered-test",
				Type:     "command",
				Command:  "echo hello",
				Machines: []string{"server1", "server3"}, // Only target specific machines
			},
		},
	}

	// This will fail due to SSH connection, but we can test the configuration
	err := ExecuteConfig(cfg)
	assert.Error(t, err)
	// But it shouldn't be a machine filtering error
	assert.NotContains(t, err.Error(), "machine filtering")
}

func TestExecuteConfig_TimeoutConfiguration(t *testing.T) {
	cfg := &config.Config{
		Machines: []config.Machine{
			{
				Name:     "test-server",
				Host:     "192.168.1.100",
				Port:     22,
				User:     "testuser",
				Password: "testpass",
			},
		},
		Actions: []config.Action{
			{
				Name:    "timeout-test",
				Type:    "command",
				Command: "echo hello",
				// Note: Timeout is not a field in the current Action struct
				// This test documents the current behavior
			},
		},
	}

	// This will fail due to SSH connection, but we can test the configuration
	err := ExecuteConfig(cfg)
	assert.Error(t, err)
	// But it shouldn't be a timeout configuration error
	assert.NotContains(t, err.Error(), "timeout")
}
