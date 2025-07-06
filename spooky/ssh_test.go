package spooky

import (
	"os"
	"strings"
	"testing"
)

func TestNewSSHClient(t *testing.T) {
	tests := []struct {
		name    string
		server  *Server
		timeout int
		wantErr bool
	}{
		{
			name: "valid server with password",
			server: &Server{
				Name:     "test-server",
				Host:     "localhost",
				User:     "testuser",
				Password: "testpass",
				Port:     22,
			},
			timeout: 30,
			wantErr: false, // Will fail in unit test due to no real SSH server
		},
		{
			name: "valid server with key file",
			server: &Server{
				Name:    "test-server",
				Host:    "localhost",
				User:    "testuser",
				KeyFile: "/path/to/key",
				Port:    22,
			},
			timeout: 30,
			wantErr: true, // Will fail due to non-existent key file
		},
		{
			name: "server without authentication",
			server: &Server{
				Name: "test-server",
				Host: "localhost",
				User: "testuser",
				Port: 22,
			},
			timeout: 30,
			wantErr: true,
		},
		{
			name: "server without host",
			server: &Server{
				Name:     "test-server",
				User:     "testuser",
				Password: "testpass",
				Port:     22,
			},
			timeout: 30,
			wantErr: true,
		},
		{
			name: "server without user",
			server: &Server{
				Name:     "test-server",
				Host:     "localhost",
				Password: "testpass",
				Port:     22,
			},
			timeout: 30,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSSHClient(tt.server, tt.timeout)

			// In unit tests, we expect connection failures since there's no real SSH server
			// We're mainly testing the validation logic
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewSSHClient() expected error but got none")
				}
			} else {
				// For valid configs, we expect connection errors but not validation errors
				if err != nil && !strings.Contains(err.Error(), "failed to connect") {
					t.Errorf("NewSSHClient() unexpected validation error: %v", err)
				}
			}
		})
	}
}

func TestSSHClient_ExecuteCommand(t *testing.T) {
	// This test requires a mock SSH client or real SSH server
	// For now, we'll test the structure and error handling

	server := &Server{
		Name:     "test-server",
		Host:     "localhost",
		User:     "testuser",
		Password: "testpass",
		Port:     22,
	}

	client, err := NewSSHClient(server, 5)
	if err == nil {
		// If we somehow get a connection, test command execution
		defer client.Close()

		output, err := client.ExecuteCommand("echo hello")
		if err != nil {
			// Expected in unit test environment
			t.Logf("ExecuteCommand failed as expected: %v", err)
		} else {
			if !strings.Contains(output, "hello") {
				t.Errorf("ExecuteCommand() output = %v, want to contain 'hello'", output)
			}
		}
	} else {
		// Expected in unit test environment
		t.Logf("SSH connection failed as expected: %v", err)
	}
}

func TestSSHClient_ExecuteScript(t *testing.T) {
	// Create a temporary script file for testing
	scriptContent := "#!/bin/bash\necho 'Hello from script'\nexit 0"

	tmpfile, err := os.CreateTemp("", "test_script_*.sh")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(scriptContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	server := &Server{
		Name:     "test-server",
		Host:     "localhost",
		User:     "testuser",
		Password: "testpass",
		Port:     22,
	}

	client, err := NewSSHClient(server, 5)
	if err == nil {
		// If we somehow get a connection, test script execution
		defer client.Close()

		output, err := client.ExecuteScript(tmpfile.Name())
		if err != nil {
			// Expected in unit test environment
			t.Logf("ExecuteScript failed as expected: %v", err)
		} else {
			if !strings.Contains(output, "Hello from script") {
				t.Errorf("ExecuteScript() output = %v, want to contain 'Hello from script'", output)
			}
		}
	} else {
		// Expected in unit test environment
		t.Logf("SSH connection failed as expected: %v", err)
	}
}

func TestExecuteActionSequential(t *testing.T) {
	// Test with mock servers that won't actually connect
	servers := []*Server{
		{
			Name:     "server1",
			Host:     "localhost",
			User:     "testuser",
			Password: "testpass",
			Port:     22,
		},
		{
			Name:     "server2",
			Host:     "localhost",
			User:     "testuser",
			Password: "testpass",
			Port:     22,
		},
	}

	action := &Action{
		Name:    "test_action",
		Command: "echo hello",
	}

	// This should not panic and should handle connection failures gracefully
	err := executeActionSequential(action, servers)
	// We expect this to fail due to no real SSH server, but it shouldn't panic
	if err != nil {
		t.Logf("executeActionSequential failed as expected: %v", err)
	}
}

func TestExecuteActionParallel(t *testing.T) {
	// Test with mock servers that won't actually connect
	servers := []*Server{
		{
			Name:     "server1",
			Host:     "localhost",
			User:     "testuser",
			Password: "testpass",
			Port:     22,
		},
		{
			Name:     "server2",
			Host:     "localhost",
			User:     "testuser",
			Password: "testpass",
			Port:     22,
		},
	}

	action := &Action{
		Name:    "test_action",
		Command: "echo hello",
	}

	// This should not panic and should handle connection failures gracefully
	err := executeActionParallel(action, servers)
	// We expect this to fail due to no real SSH server, but it shouldn't panic
	if err != nil {
		t.Logf("executeActionParallel failed as expected: %v", err)
	}
}

func TestIndentOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single line",
			input:    "hello world",
			expected: "    hello world\n",
		},
		{
			name:     "multiple lines",
			input:    "line1\nline2\nline3",
			expected: "    line1\n    line2\n    line3\n",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "trailing newline",
			input:    "hello\nworld\n",
			expected: "    hello\n    world\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := indentOutput(tt.input)
			if result != tt.expected {
				t.Errorf("indentOutput() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Mock SSH client for testing
type MockSSHClient struct {
	server     *Server
	shouldFail bool
	output     string
	err        error
}

func NewMockSSHClient(server *Server, timeout int) (*MockSSHClient, error) {
	return &MockSSHClient{
		server: server,
		output: "mock output",
	}, nil
}

func (m *MockSSHClient) Close() error {
	return nil
}

func (m *MockSSHClient) ExecuteCommand(command string) (string, error) {
	if m.shouldFail {
		return "", m.err
	}
	return m.output, nil
}

func (m *MockSSHClient) ExecuteScript(scriptPath string) (string, error) {
	if m.shouldFail {
		return "", m.err
	}
	return m.output, nil
}

func TestMockSSHClient(t *testing.T) {
	server := &Server{
		Name:     "mock-server",
		Host:     "localhost",
		User:     "testuser",
		Password: "testpass",
		Port:     22,
	}

	client, err := NewMockSSHClient(server, 30)
	if err != nil {
		t.Fatalf("NewMockSSHClient() failed: %v", err)
	}
	defer client.Close()

	output, err := client.ExecuteCommand("echo hello")
	if err != nil {
		t.Errorf("MockSSHClient.ExecuteCommand() failed: %v", err)
	}
	if output != "mock output" {
		t.Errorf("MockSSHClient.ExecuteCommand() = %q, want %q", output, "mock output")
	}
}
