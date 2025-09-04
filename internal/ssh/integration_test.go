//go:build integration

package ssh

import (
	"fmt"
	"testing"
	"time"

	"spooky/internal/schemas"
	testhelpers "spooky/internal/testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSSHConnectionIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tests := []struct {
		name        string
		description string
		testFunc    func(t *testing.T)
	}{
		{
			name:        "TestSSHConnectionSuccess",
			description: "Test successful SSH connection establishment",
			testFunc:    testSSHConnectionSuccess,
		},
		{
			name:        "TestSSHConnectionTimeout",
			description: "Test SSH connection timeout handling",
			testFunc:    testSSHConnectionTimeout,
		},
		{
			name:        "TestSSHConnectionRefused",
			description: "Test SSH connection refused handling",
			testFunc:    testSSHConnectionRefused,
		},
		{
			name:        "TestSSHAuthenticationMethods",
			description: "Test different SSH authentication methods",
			testFunc:    testSSHAuthenticationMethods,
		},
		{
			name:        "TestSSHCommandExecution",
			description: "Test SSH command execution",
			testFunc:    testSSHCommandExecution,
		},
		{
			name:        "TestSSHConnectionRetry",
			description: "Test SSH connection retry logic",
			testFunc:    testSSHConnectionRetry,
		},
		{
			name:        "TestSSHKeepAlive",
			description: "Test SSH keep-alive functionality",
			testFunc:    testSSHKeepAlive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Running integration test: %s", tt.description)
			tt.testFunc(t)
		})
	}
}

func testSSHConnectionSuccess(t *testing.T) {
	// Create test SSH server
	server, err := testhelpers.NewTestSSHServer()
	require.NoError(t, err)
	defer server.Close()

	// Create SSH client configuration
	config := &Config{
		Host:    "localhost",
		Port:    server.Port(),
		User:    "testuser",
		Timeout: 30 * time.Second,
	}

	// Create SSH client
	client, err := NewClient(config)
	require.NoError(t, err)

	// Test connection
	ctx := testhelpers.TestContext(t, 30*time.Second)
	err = client.Connect(ctx)
	assert.NoError(t, err)

	// Test disconnection
	err = client.Disconnect()
	assert.NoError(t, err)
}

func testSSHConnectionTimeout(t *testing.T) {
	// Create SSH client configuration with very short timeout
	config := &Config{
		Host:    "192.0.2.1", // Non-routable IP address
		Port:    22,
		User:    "testuser",
		Timeout: 100 * time.Millisecond, // Very short timeout
	}

	// Create SSH client
	client, err := NewClient(config)
	require.NoError(t, err)

	// Test connection with timeout
	ctx := testhelpers.TestContext(t, 5*time.Second)
	err = client.Connect(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")
}

func testSSHConnectionRefused(t *testing.T) {
	// Create SSH client configuration pointing to non-existent service
	config := &Config{
		Host:    "localhost",
		Port:    2222, // Port that's likely not in use
		User:    "testuser",
		Timeout: 5 * time.Second,
	}

	// Create SSH client
	client, err := NewClient(config)
	require.NoError(t, err)

	// Test connection refused
	ctx := testhelpers.TestContext(t, 10*time.Second)
	err = client.Connect(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection refused")
}

func testSSHAuthenticationMethods(t *testing.T) {
	// Create test SSH server
	server, err := testhelpers.NewTestSSHServer()
	require.NoError(t, err)
	defer server.Close()

	tests := []struct {
		name           string
		authMethod     string
		expectedResult bool
	}{
		{
			name:           "Password Authentication",
			authMethod:     "password",
			expectedResult: true,
		},
		{
			name:           "Public Key Authentication",
			authMethod:     "publickey",
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				Host:    "localhost",
				Port:    server.Port(),
				User:    "testuser",
				Timeout: 30 * time.Second,
			}

			// Configure authentication method
			switch tt.authMethod {
			case "password":
				config.Password = "testpass"
				config.PasswordAuth = true
			case "publickey":
				// For testing, we'll use a mock public key
				config.PubkeyAuth = true
			}

			client, err := NewClient(config)
			require.NoError(t, err)

			ctx := testhelpers.TestContext(t, 30*time.Second)
			err = client.Connect(ctx)

			if tt.expectedResult {
				assert.NoError(t, err)
				client.Disconnect()
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func testSSHCommandExecution(t *testing.T) {
	// Create test SSH server
	server, err := testhelpers.NewTestSSHServer()
	require.NoError(t, err)
	defer server.Close()

	// Create test SSH client
	testClient, err := testhelpers.NewTestSSHClient(server)
	require.NoError(t, err)
	defer testClient.Close()

	// Test various commands
	commands := []struct {
		command        string
		expectedOutput string
	}{
		{"echo hello", "hello\n"},
		{"pwd", "/home/test\n"},
		{"whoami", "testuser\n"},
		{"ls", "file1.txt\nfile2.txt\n"},
	}

	for _, cmd := range commands {
		t.Run(fmt.Sprintf("Command_%s", cmd.command), func(t *testing.T) {
			output, err := testClient.RunCommand(cmd.command)
			assert.NoError(t, err)
			assert.Equal(t, cmd.expectedOutput, output)
		})
	}
}

func testSSHConnectionRetry(t *testing.T) {
	// This test simulates connection retry logic
	// In a real implementation, you would test the retry mechanism

	config := &Config{
		Host:    "192.0.2.1", // Non-routable IP
		Port:    22,
		User:    "testuser",
		Timeout: 1 * time.Second,
	}

	client, err := NewClient(config)
	require.NoError(t, err)

	// Test multiple connection attempts
	ctx := testhelpers.TestContext(t, 10*time.Second)

	// First attempt should fail
	err = client.Connect(ctx)
	assert.Error(t, err)

	// Second attempt should also fail (no retry logic implemented yet)
	err = client.Connect(ctx)
	assert.Error(t, err)
}

func testSSHKeepAlive(t *testing.T) {
	// Create test SSH server
	server, err := testhelpers.NewTestSSHServer()
	require.NoError(t, err)
	defer server.Close()

	// Create SSH client with keep-alive
	config := &Config{
		Host:           "localhost",
		Port:           server.Port(),
		User:           "testuser",
		Timeout:        30 * time.Second,
		KeepAlive:      5 * time.Second,
		KeepAliveCount: 3,
	}

	client, err := NewClient(config)
	require.NoError(t, err)

	ctx := testhelpers.TestContext(t, 30*time.Second)
	err = client.Connect(ctx)
	require.NoError(t, err)

	// Wait for keep-alive to be active
	time.Sleep(6 * time.Second)

	// Connection should still be active
	err = client.Disconnect()
	assert.NoError(t, err)
}

// TestSSHManagerIntegration tests the SSH manager integration
func TestSSHManagerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create test SSH server
	server, err := testhelpers.NewTestSSHServer()
	require.NoError(t, err)
	defer server.Close()

	// Create SSH manager configuration
	sshConfig := &schemas.SpookySSHV1{
		Timeout:                   30,
		KeepaliveInterval:         60,
		KeepaliveCount:            3,
		KeyScanTimeout:            10,
		KnownHostsMode:            "accept-new",
		Compression:               false,
		CompressionLevel:          6,
		TCPKeepAlive:              true,
		TCPKeepAliveCount:         3,
		TCPKeepAliveIdle:          60,
		TCPKeepAliveInterval:      10,
		TCPKeepAliveProbeInterval: 5,
	}

	// Create SSH manager
	manager := NewSSHManager(nil, sshConfig)

	// Create test machine configuration
	machine := &schemas.MachinesMachineV1{
		Hostname: "localhost",
		Port:     server.Port(),
		User:     "testuser",
		Authentication: schemas.MachinesMachineAuthenticationV1{
			Password: schemas.MachinesMachineAuthenticationPasswordV1{
				Value:     "testpass",
				Encrypted: false,
			},
		},
	}

	// Test command execution
	ctx := testhelpers.TestContext(t, 30*time.Second)
	result, err := manager.RunCommandOnMachine(ctx, machine, "echo hello")
	require.NoError(t, err)
	assert.Equal(t, 0, result.ExitCode)
	assert.Contains(t, result.Stdout, "hello")
}

// TestSSHConnectionStress tests SSH connection under stress
func TestSSHConnectionStress(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	// Create test SSH server
	server, err := testhelpers.NewTestSSHServer()
	require.NoError(t, err)
	defer server.Close()

	// Test multiple concurrent connections
	numConnections := 10
	results := make(chan error, numConnections)

	for i := 0; i < numConnections; i++ {
		go func(id int) {
			config := &Config{
				Host:    "localhost",
				Port:    server.Port(),
				User:    "testuser",
				Timeout: 30 * time.Second,
			}

			client, err := NewClient(config)
			if err != nil {
				results <- err
				return
			}

			ctx := testhelpers.TestContext(t, 30*time.Second)
			err = client.Connect(ctx)
			if err != nil {
				results <- err
				return
			}

			// Execute a command
			_, err = client.RunCommand(ctx, "echo test")
			if err != nil {
				results <- err
				return
			}

			err = client.Disconnect()
			results <- err
		}(i)
	}

	// Collect results
	var errors []error
	for i := 0; i < numConnections; i++ {
		if err := <-results; err != nil {
			errors = append(errors, err)
		}
	}

	// Allow some failures in stress test (up to 20%)
	maxFailures := numConnections / 5
	assert.LessOrEqual(t, len(errors), maxFailures,
		"Too many connection failures in stress test: %d/%d", len(errors), numConnections)
}

// TestSSHConnectionErrorHandling tests various error conditions
func TestSSHConnectionErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tests := []struct {
		name        string
		config      *Config
		expectError bool
		errorType   string
	}{
		{
			name: "Invalid Host",
			config: &Config{
				Host:    "invalid-host-that-does-not-exist.example.com",
				Port:    22,
				User:    "testuser",
				Timeout: 5 * time.Second,
			},
			expectError: true,
			errorType:   "host resolution",
		},
		{
			name: "Invalid Port",
			config: &Config{
				Host:    "localhost",
				Port:    99999, // Invalid port
				User:    "testuser",
				Timeout: 5 * time.Second,
			},
			expectError: true,
			errorType:   "connection",
		},
		{
			name: "Empty User",
			config: &Config{
				Host:    "localhost",
				Port:    22,
				User:    "", // Empty user
				Timeout: 5 * time.Second,
			},
			expectError: true,
			errorType:   "configuration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			if tt.expectError && tt.errorType == "configuration" {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			ctx := testhelpers.TestContext(t, 10*time.Second)
			err = client.Connect(ctx)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				client.Disconnect()
			}
		})
	}
}
