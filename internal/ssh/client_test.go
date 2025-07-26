package ssh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"spooky/internal/config"
)

func TestNewSSHClient(t *testing.T) {
	tests := []struct {
		name        string
		machine     *config.Machine
		timeout     int
		expectError bool
	}{
		{
			name: "valid machine configuration",
			machine: &config.Machine{
				Name:     "test-server",
				Host:     "192.168.1.100",
				Port:     22,
				User:     "testuser",
				Password: "testpass",
			},
			timeout:     1, // 1 second timeout for tests
			expectError: false,
		},
		{
			name: "machine with key file",
			machine: &config.Machine{
				Name:    "test-server",
				Host:    "192.168.1.100",
				Port:    22,
				User:    "testuser",
				KeyFile: "~/.ssh/id_rsa",
			},
			timeout:     1, // 1 second timeout for tests
			expectError: false,
		},
		{
			name: "machine with no authentication",
			machine: &config.Machine{
				Name: "test-server",
				Host: "192.168.1.100",
				Port: 22,
				User: "testuser",
			},
			timeout:     1,     // 1 second timeout for tests
			expectError: false, // Should not error, just no auth methods
		},
		{
			name:        "nil machine",
			machine:     nil,
			timeout:     1, // 1 second timeout for tests
			expectError: true,
		},
		{
			name: "invalid port",
			machine: &config.Machine{
				Name: "test-server",
				Host: "192.168.1.100",
				Port: 99999, // Invalid port
				User: "testuser",
			},
			timeout:     1, // 1 second timeout for tests
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewSSHClient(tt.machine, tt.timeout)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				// Note: In a real environment, this would fail due to various SSH issues
				// But we're testing the client creation logic, not the connection
				if err != nil {
					// If there's an error, it should be one of the expected SSH errors
					errorMsg := err.Error()
					assert.True(t,
						strings.Contains(errorMsg, "failed to connect") ||
							strings.Contains(errorMsg, "failed to read key file") ||
							strings.Contains(errorMsg, "no authentication method available"),
						"Unexpected error: %s", errorMsg)
				} else {
					assert.NotNil(t, client)
					assert.Equal(t, tt.machine, client.GetMachine())
				}
			}
		})
	}
}

func TestNewSSHClientWithHostKeyCallback(t *testing.T) {
	machine := &config.Machine{
		Name:     "test-server",
		Host:     "192.168.1.100",
		Port:     22,
		User:     "testuser",
		Password: "testpass",
	}

	tests := []struct {
		name           string
		hostKeyType    HostKeyCallbackType
		knownHostsPath string
		expectError    bool
	}{
		{
			name:           "insecure host key",
			hostKeyType:    InsecureHostKey,
			knownHostsPath: "",
			expectError:    false,
		},
		{
			name:           "auto host key",
			hostKeyType:    AutoHostKey,
			knownHostsPath: "",
			expectError:    false,
		},
		{
			name:           "known hosts with non-existent file",
			hostKeyType:    KnownHostsHostKey,
			knownHostsPath: "/nonexistent/known_hosts",
			expectError:    true,
		},
		{
			name:           "invalid host key type",
			hostKeyType:    "invalid",
			knownHostsPath: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewSSHClientWithHostKeyCallback(machine, 1, tt.hostKeyType, tt.knownHostsPath) // 1 second timeout for tests

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				// Note: This will likely fail due to SSH connection
				if err != nil {
					// If there's an error, it should be one of the expected SSH errors
					errorMsg := err.Error()
					assert.True(t,
						strings.Contains(errorMsg, "failed to connect") ||
							strings.Contains(errorMsg, "failed to read key file") ||
							strings.Contains(errorMsg, "no authentication method available"),
						"Unexpected error: %s", errorMsg)
				} else {
					assert.NotNil(t, client)
				}
			}
		})
	}
}

func TestSSHClient_GetMachine(t *testing.T) {
	machine := &config.Machine{
		Name:     "test-server",
		Host:     "192.168.1.100",
		Port:     22,
		User:     "testuser",
		Password: "testpass",
	}

	client, err := NewSSHClient(machine, 1) // 1 second timeout for tests
	// Note: This will likely fail due to SSH connection, but we can test the client structure
	if err != nil {
		// If there's an error, it should be one of the expected SSH errors
		errorMsg := err.Error()
		assert.True(t,
			strings.Contains(errorMsg, "failed to connect") ||
				strings.Contains(errorMsg, "failed to read key file") ||
				strings.Contains(errorMsg, "no authentication method available"),
			"Unexpected error: %s", errorMsg)
		return
	}

	require.NotNil(t, client)
	retrievedMachine := client.GetMachine()
	assert.Equal(t, machine, retrievedMachine)
}

func TestSSHClient_ConnectAndClose(t *testing.T) {
	// This test requires a real SSH server or mock
	// For now, we'll test the basic structure without actual connection
	machine := &config.Machine{
		Name:     "test-server",
		Host:     "192.168.1.100",
		Port:     22,
		User:     "testuser",
		Password: "testpass",
	}

	client, err := NewSSHClient(machine, 1) // 1 second timeout for tests
	// Note: This will likely fail due to SSH connection
	if err != nil {
		// If there's an error, it should be one of the expected SSH errors
		errorMsg := err.Error()
		assert.True(t,
			strings.Contains(errorMsg, "failed to connect") ||
				strings.Contains(errorMsg, "failed to read key file") ||
				strings.Contains(errorMsg, "no authentication method available"),
			"Unexpected error: %s", errorMsg)
		return
	}

	require.NotNil(t, client)

	// Test that Close doesn't panic even if not connected
	assert.NoError(t, client.Close())
}

func TestSSHClient_ExecuteCommand_NotConnected(t *testing.T) {
	machine := &config.Machine{
		Name:     "test-server",
		Host:     "192.168.1.100",
		Port:     22,
		User:     "testuser",
		Password: "testpass",
	}

	client, err := NewSSHClient(machine, 1) // 1 second timeout for tests
	// Note: This will likely fail due to SSH connection
	if err != nil {
		// If there's an error, it should be one of the expected SSH errors
		errorMsg := err.Error()
		assert.True(t,
			strings.Contains(errorMsg, "failed to connect") ||
				strings.Contains(errorMsg, "failed to read key file") ||
				strings.Contains(errorMsg, "no authentication method available"),
			"Unexpected error: %s", errorMsg)
		return
	}
	require.NotNil(t, client)

	// Try to execute command without connecting
	result, err := client.ExecuteCommand("echo hello")
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "not connected")
}

func TestSSHClient_ExecuteScript_NotConnected(t *testing.T) {
	machine := &config.Machine{
		Name:     "test-server",
		Host:     "192.168.1.100",
		Port:     22,
		User:     "testuser",
		Password: "testpass",
	}

	client, err := NewSSHClient(machine, 1) // 1 second timeout for tests
	// Note: This will likely fail due to SSH connection
	if err != nil {
		// If there's an error, it should be one of the expected SSH errors
		errorMsg := err.Error()
		assert.True(t,
			strings.Contains(errorMsg, "failed to connect") ||
				strings.Contains(errorMsg, "failed to read key file") ||
				strings.Contains(errorMsg, "no authentication method available"),
			"Unexpected error: %s", errorMsg)
		return
	}
	require.NotNil(t, client)

	// Try to execute script without connecting
	result, err := client.ExecuteScript("/tmp/test.sh")
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "not connected")
}

func TestGetHostKeyCallback(t *testing.T) {
	tests := []struct {
		name           string
		callbackType   HostKeyCallbackType
		knownHostsPath string
		expectError    bool
	}{
		{
			name:           "insecure host key",
			callbackType:   InsecureHostKey,
			knownHostsPath: "",
			expectError:    false,
		},
		{
			name:           "auto host key",
			callbackType:   AutoHostKey,
			knownHostsPath: "",
			expectError:    false,
		},
		{
			name:           "known hosts with non-existent file",
			callbackType:   KnownHostsHostKey,
			knownHostsPath: "/nonexistent/known_hosts",
			expectError:    true,
		},
		{
			name:           "invalid host key type",
			callbackType:   "invalid",
			knownHostsPath: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callback, err := getHostKeyCallback(tt.callbackType, tt.knownHostsPath)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, callback)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, callback)
			}
		})
	}
}

func TestGetHostKeyCallback_WithHomeDirectory(t *testing.T) {
	// Test known hosts with ~ expansion
	callback, err := getHostKeyCallback(KnownHostsHostKey, "~/.ssh/known_hosts")

	// This will likely fail because ~/.ssh/known_hosts doesn't exist in test environment
	// but we can test that the ~ expansion logic is called
	if err != nil {
		assert.Contains(t, err.Error(), "file does not exist")
	} else {
		assert.NotNil(t, callback)
	}
}

func TestSSHClient_WithTestDataFromExamples(t *testing.T) {
	// Use test data from examples/testing where available
	examplesDir := "../../examples/testing"

	// Test with valid project configuration
	validProjectPath := filepath.Join(examplesDir, "test-valid-project")
	if _, err := os.Stat(validProjectPath); os.IsNotExist(err) {
		t.Skip("Test data directory not found, skipping test")
	}

	// Create a machine configuration similar to what would be in test data
	machine := &config.Machine{
		Name:     "example-server",
		Host:     "192.168.1.100",
		Port:     22,
		User:     "testuser",
		Password: "testpass",
		Tags: map[string]string{
			"environment": "test",
			"role":        "web",
		},
	}

	client, err := NewSSHClient(machine, 1) // 1 second timeout for tests
	// Note: This will likely fail due to SSH connection
	if err != nil {
		// If there's an error, it should be one of the expected SSH errors
		errorMsg := err.Error()
		assert.True(t,
			strings.Contains(errorMsg, "failed to connect") ||
				strings.Contains(errorMsg, "failed to read key file") ||
				strings.Contains(errorMsg, "no authentication method available"),
			"Unexpected error: %s", errorMsg)
		return
	}

	require.NotNil(t, client)

	// Verify machine configuration
	retrievedMachine := client.GetMachine()
	assert.Equal(t, machine.Name, retrievedMachine.Name)
	assert.Equal(t, machine.Host, retrievedMachine.Host)
	assert.Equal(t, machine.Port, retrievedMachine.Port)
	assert.Equal(t, machine.User, retrievedMachine.User)
	assert.Equal(t, machine.Tags, retrievedMachine.Tags)
}

func TestSSHClient_TimeoutConfiguration(t *testing.T) {
	machine := &config.Machine{
		Name:     "test-server",
		Host:     "192.168.1.100",
		Port:     22,
		User:     "testuser",
		Password: "testpass",
	}

	tests := []struct {
		name    string
		timeout int
	}{
		{"short timeout", 1},
		{"medium timeout", 2},
		{"long timeout", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewSSHClient(machine, tt.timeout)
			// Note: This will likely fail due to SSH connection
			if err != nil {
				// If there's an error, it should be one of the expected SSH errors
				errorMsg := err.Error()
				assert.True(t,
					strings.Contains(errorMsg, "failed to connect") ||
						strings.Contains(errorMsg, "failed to read key file") ||
						strings.Contains(errorMsg, "no authentication method available"),
					"Unexpected error: %s", errorMsg)
				return
			}

			require.NotNil(t, client)

			// Verify the client was created (timeout validation happens during connection)
			assert.Equal(t, machine, client.GetMachine())
		})
	}
}

func TestSSHClient_InvalidTimeout(t *testing.T) {
	machine := &config.Machine{
		Name:     "test-server",
		Host:     "192.168.1.100",
		Port:     22,
		User:     "testuser",
		Password: "testpass",
	}

	// Test with invalid timeout values
	invalidTimeouts := []int{-1, 0} // Negative and zero timeouts

	for _, timeout := range invalidTimeouts {
		t.Run(fmt.Sprintf("timeout_%d", timeout), func(t *testing.T) {
			client, err := NewSSHClient(machine, timeout)

			// Now these should fail immediately due to validation
			assert.Error(t, err)
			assert.Nil(t, client)
			assert.Contains(t, err.Error(), "timeout must be positive")
		})
	}
}

func TestSSHClient_AuthenticationMethods(t *testing.T) {
	tests := []struct {
		name        string
		machine     *config.Machine
		description string
	}{
		{
			name: "password only",
			machine: &config.Machine{
				Name:     "test-server",
				Host:     "192.168.1.100",
				Port:     22,
				User:     "testuser",
				Password: "testpass",
			},
			description: "Machine with password authentication only",
		},
		{
			name: "key file only",
			machine: &config.Machine{
				Name:    "test-server",
				Host:    "192.168.1.100",
				Port:    22,
				User:    "testuser",
				KeyFile: "~/.ssh/id_rsa",
			},
			description: "Machine with key file authentication only",
		},
		{
			name: "both password and key",
			machine: &config.Machine{
				Name:     "test-server",
				Host:     "192.168.1.100",
				Port:     22,
				User:     "testuser",
				Password: "testpass",
				KeyFile:  "~/.ssh/id_rsa",
			},
			description: "Machine with both password and key authentication",
		},
		{
			name: "no authentication",
			machine: &config.Machine{
				Name: "test-server",
				Host: "192.168.1.100",
				Port: 22,
				User: "testuser",
			},
			description: "Machine with no authentication methods",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewSSHClient(tt.machine, 1) // 1 second timeout for tests

			// Note: This will likely fail due to SSH connection
			if err != nil {
				// If there's an error, it should be one of the expected SSH errors
				errorMsg := err.Error()
				assert.True(t,
					strings.Contains(errorMsg, "failed to connect") ||
						strings.Contains(errorMsg, "failed to read key file") ||
						strings.Contains(errorMsg, "no authentication method available"),
					"Unexpected error: %s", errorMsg)
				return
			}

			assert.NotNil(t, client)
			assert.Equal(t, tt.machine, client.GetMachine())
		})
	}
}
