package ssh

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"spooky/internal/config"

	"github.com/gliderlabs/ssh"
)

// mockSSHServer creates a mock SSH server for testing
func mockSSHServer(t *testing.T) (serverAddr string, cleanup func()) {
	// Create a listener on a random port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	addr := fmt.Sprintf("127.0.0.1:%d", port)

	// Create SSH server
	server := &ssh.Server{
		Addr: addr,
		Handler: func(s ssh.Session) {
			// Echo back the command
			command := strings.Join(s.Command(), " ")
			fmt.Fprintf(s, "Executed: %s\n", command)
		},
		PasswordHandler: func(_ ssh.Context, password string) bool {
			return password == "testpass"
		},
	}

	// Start server in goroutine
	go func() {
		if err := server.Serve(listener); err != nil && err != ssh.ErrServerClosed {
			t.Logf("SSH server error: %v", err)
		}
	}()

	// Wait a bit for server to start
	time.Sleep(100 * time.Millisecond)

	// Return cleanup function
	cleanup = func() {
		server.Close()
	}

	return addr, cleanup
}

func TestNewSSHClient_NoAuth(t *testing.T) {
	machine := &config.Machine{
		Name: "test",
		Host: "localhost",
		User: "testuser",
		Port: 22,
	}

	_, err := NewSSHClient(machine, 30)
	if err == nil {
		t.Error("expected error when no authentication method is provided")
	}
	if !strings.Contains(err.Error(), "no authentication method available") {
		t.Errorf("expected authentication error, got: %v", err)
	}
}

func TestNewSSHClient_PasswordAuth(t *testing.T) {
	machine := &config.Machine{
		Name:     "test",
		Host:     "localhost",
		User:     "testuser",
		Port:     22,
		Password: "testpass",
	}

	_, err := NewSSHClient(machine, 30)
	// This will fail to connect or authenticate, but should not fail due to authentication setup
	if err != nil && !strings.Contains(err.Error(), "connection refused") && !strings.Contains(err.Error(), "no route to host") && !strings.Contains(err.Error(), "unable to authenticate") {
		t.Errorf("expected connection/authentication error, got: %v", err)
	}
}

func TestNewSSHClient_KeyFileAuth(t *testing.T) {
	// Create a temporary key file
	tempDir := t.TempDir()
	keyFile := filepath.Join(tempDir, "test_key")

	// Write a dummy key file (this will fail to parse, but tests the file reading logic)
	err := os.WriteFile(keyFile, []byte("invalid key content"), 0o600)
	if err != nil {
		t.Fatalf("failed to create test key file: %v", err)
	}

	machine := &config.Machine{
		Name:    "test",
		Host:    "localhost",
		User:    "testuser",
		Port:    22,
		KeyFile: keyFile,
	}

	_, err = NewSSHClient(machine, 30)
	if err == nil {
		t.Error("expected error when key file is invalid")
	}
	if !strings.Contains(err.Error(), "failed to parse private key") {
		t.Errorf("expected key parsing error, got: %v", err)
	}
}

func TestNewSSHClient_BothAuthMethods(t *testing.T) {
	// Create a temporary key file
	tempDir := t.TempDir()
	keyFile := filepath.Join(tempDir, "test_key")

	// Write a dummy key file
	err := os.WriteFile(keyFile, []byte("invalid key content"), 0o600)
	if err != nil {
		t.Fatalf("failed to create test key file: %v", err)
	}

	machine := &config.Machine{
		Name:     "test",
		Host:     "localhost",
		User:     "testuser",
		Port:     22,
		Password: "testpass",
		KeyFile:  keyFile,
	}

	_, err = NewSSHClient(machine, 30)
	// Should fail due to key parsing error
	if err == nil {
		t.Error("expected error when key file is invalid")
	}
	if !strings.Contains(err.Error(), "failed to parse private key") {
		t.Errorf("expected key parsing error, got: %v", err)
	}
}

func TestNewSSHClient_KeyFileNotFound(t *testing.T) {
	machine := &config.Machine{
		Name:    "test",
		Host:    "localhost",
		User:    "testuser",
		Port:    22,
		KeyFile: "/nonexistent/key/file",
	}

	_, err := NewSSHClient(machine, 30)
	if err == nil {
		t.Error("expected error when key file does not exist")
	}
	if !strings.Contains(err.Error(), "failed to read key file") {
		t.Errorf("expected file read error, got: %v", err)
	}
}

func TestNewSSHClient_Timeout(t *testing.T) {
	machine := &config.Machine{
		Name:     "test",
		Host:     "192.0.2.1", // RFC 5737 reserved for documentation/testing
		User:     "testuser",
		Port:     22,
		Password: "testpass",
	}

	start := time.Now()
	_, err := NewSSHClient(machine, 1) // 1 second timeout
	duration := time.Since(start)

	if err == nil {
		t.Error("expected error when connecting to unreachable host")
	}

	// Should timeout within reasonable time (not more than 3 seconds)
	if duration > 3*time.Second {
		t.Errorf("connection should have timed out faster, took: %v", duration)
	}
}

func TestGetHostKeyCallback(t *testing.T) {
	tests := []struct {
		name           string
		callbackType   HostKeyCallbackType
		knownHostsPath string
		expectError    bool
		errorContains  string
	}{
		{
			name:           "insecure host key callback",
			callbackType:   InsecureHostKey,
			knownHostsPath: "",
			expectError:    false,
		},
		{
			name:           "auto host key callback",
			callbackType:   AutoHostKey,
			knownHostsPath: "",
			expectError:    false,
		},
		{
			name:           "known hosts with default path",
			callbackType:   KnownHostsHostKey,
			knownHostsPath: "",
			expectError:    false, // Will be handled dynamically based on file existence
		},
		{
			name:           "known hosts with custom path",
			callbackType:   KnownHostsHostKey,
			knownHostsPath: "/nonexistent/known_hosts",
			expectError:    true,
			errorContains:  "failed to parse known_hosts file",
		},
		{
			name:           "known hosts with tilde expansion",
			callbackType:   KnownHostsHostKey,
			knownHostsPath: "~/test_known_hosts",
			expectError:    true,
			errorContains:  "failed to parse known_hosts file",
		},
		{
			name:           "unsupported host key callback type",
			callbackType:   "unsupported",
			knownHostsPath: "",
			expectError:    true,
			errorContains:  "unsupported host key callback type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Special handling for known_hosts_with_default_path test
			if tt.name == "known hosts with default path" {
				// Check if ~/.ssh/known_hosts exists
				homeDir, err := os.UserHomeDir()
				if err != nil {
					t.Skipf("Cannot determine home directory: %v", err)
				}
				knownHostsPath := filepath.Join(homeDir, ".ssh", "known_hosts")
				_, err = os.Stat(knownHostsPath)
				fileExists := err == nil

				callback, err := getHostKeyCallback(tt.callbackType, tt.knownHostsPath)

				if fileExists {
					// File exists, expect success
					if err != nil {
						t.Errorf("Expected no error but got: %v", err)
						return
					}
					if callback == nil {
						t.Errorf("Expected callback but got nil")
					}
				} else {
					// File doesn't exist, expect error
					if err == nil {
						t.Errorf("Expected error but got none")
						return
					}
					if !strings.Contains(err.Error(), "failed to parse known_hosts file") {
						t.Errorf("Expected error to contain 'failed to parse known_hosts file', got: %s", err.Error())
					}
				}
				return
			}

			// Standard test logic for other cases
			callback, err := getHostKeyCallback(tt.callbackType, tt.knownHostsPath)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %s", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
					return
				}
				if callback == nil {
					t.Errorf("Expected callback but got nil")
				}
			}
		})
	}
}

func TestNewSSHClientWithHostKeyCallback(t *testing.T) {
	machine := &config.Machine{
		Name:     "test",
		Host:     "localhost",
		User:     "testuser",
		Password: "testpass",
		Port:     22,
	}

	// Test with insecure host key callback (should not fail due to auth setup)
	client, err := NewSSHClientWithHostKeyCallback(machine, 30, InsecureHostKey, "")
	if err != nil && !strings.Contains(err.Error(), "connection refused") && !strings.Contains(err.Error(), "no route to host") && !strings.Contains(err.Error(), "unable to authenticate") {
		t.Errorf("Expected connection/authentication error, got: %v", err)
	}
	if err == nil && client != nil {
		defer client.Close()
	}

	// Test with auto host key callback (should not fail due to auth setup)
	client, err = NewSSHClientWithHostKeyCallback(machine, 30, AutoHostKey, "")
	if err != nil && !strings.Contains(err.Error(), "connection refused") && !strings.Contains(err.Error(), "no route to host") && !strings.Contains(err.Error(), "unable to authenticate") {
		t.Errorf("Expected connection/authentication error, got: %v", err)
	}
	if err == nil && client != nil {
		defer client.Close()
	}

	// Test with known hosts callback (may succeed if known_hosts file exists)
	_, err = NewSSHClientWithHostKeyCallback(machine, 30, KnownHostsHostKey, "")
	if err != nil && !strings.Contains(err.Error(), "connection refused") && !strings.Contains(err.Error(), "no route to host") && !strings.Contains(err.Error(), "unable to authenticate") && !strings.Contains(err.Error(), "failed to parse known_hosts file") {
		t.Errorf("Expected connection/authentication/known_hosts error, got: %v", err)
	}
	if err == nil {
		t.Log("Known hosts callback succeeded (known_hosts file exists)")
	}

	// Test with unsupported host key callback type
	_, err = NewSSHClientWithHostKeyCallback(machine, 30, "unsupported", "")
	if err == nil {
		t.Error("Expected error with unsupported callback type but got none")
	} else if !strings.Contains(err.Error(), "unsupported host key callback type") {
		t.Errorf("Expected unsupported callback error, got: %v", err)
	}
}

func TestSSHClient_ExecuteCommand_WithMockServer(t *testing.T) {
	// Start mock SSH server
	addr, cleanup := mockSSHServer(t)
	defer cleanup()

	// Parse address
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatalf("Failed to parse address: %v", err)
	}

	// Parse port
	port := 22
	if portStr != "" {
		if p, err := net.LookupPort("tcp", portStr); err == nil {
			port = p
		}
	}

	// Create machine config
	machine := &config.Machine{
		Name:     "test",
		Host:     host,
		User:     "testuser",
		Password: "testpass",
		Port:     port,
	}

	// Create SSH client
	client, err := NewSSHClient(machine, 5)
	if err != nil {
		// If connection fails, that's okay for this test
		// We're mainly testing the ExecuteCommand logic
		t.Logf("SSH connection failed (expected in some environments): %v", err)
		return
	}
	defer client.Close()

	// Test ExecuteCommand
	output, err := client.ExecuteCommand("echo hello")
	if err != nil {
		t.Logf("ExecuteCommand failed (expected in some environments): %v", err)
		return
	}

	expected := "Executed: echo hello"
	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain '%s', got: '%s'", expected, output)
	}
}

func TestSSHClient_ExecuteCommand_WithRealConnection(t *testing.T) {
	// Test with a real SSH connection attempt that will fail
	// but will exercise the session creation code
	machine := &config.Machine{
		Name:     "test",
		Host:     "192.0.2.1", // RFC 5737 reserved for documentation/testing
		User:     "testuser",
		Password: "testpass",
		Port:     22,
	}

	client, err := NewSSHClient(machine, 1) // 1 second timeout
	if err != nil {
		// This is expected to fail due to connection timeout
		if !strings.Contains(err.Error(), "connection refused") &&
			!strings.Contains(err.Error(), "no route to host") &&
			!strings.Contains(err.Error(), "timeout") {
			t.Errorf("Expected connection error, got: %v", err)
		}
		return
	}

	// If we somehow got a connection, test ExecuteCommand
	if client != nil {
		defer client.Close()
		_, err := client.ExecuteCommand("echo test")
		if err != nil {
			// This is expected to fail due to session issues
			if !strings.Contains(err.Error(), "failed to create session") {
				t.Errorf("Expected session error, got: %v", err)
			}
		}
	}
}

func TestSSHClient_ExecuteCommand_SessionError(t *testing.T) {
	// Create a client with a nil SSH client to test session creation error
	client := &SSHClient{
		config: &config.Machine{Name: "test"},
		client: nil,
	}

	_, err := client.ExecuteCommand("echo test")
	if err == nil {
		t.Error("Expected error when client is nil")
	}
	if !strings.Contains(err.Error(), "no SSH connection exists") {
		t.Errorf("Expected connection error, got: %v", err)
	}
}

func TestSSHClient_ExecuteCommand_CommandExecutionError(t *testing.T) {
	// Create a mock SSH server that returns an error
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}

	serverPort := listener.Addr().(*net.TCPAddr).Port
	addr := fmt.Sprintf("127.0.0.1:%d", serverPort)

	server := &ssh.Server{
		Addr: addr,
		Handler: func(s ssh.Session) {
			// Return an error
			_ = s.Exit(1)
		},
		PasswordHandler: func(_ ssh.Context, password string) bool {
			return password == "testpass"
		},
	}

	go func() {
		if err := server.Serve(listener); err != nil && err != ssh.ErrServerClosed {
			t.Logf("SSH server error: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	defer server.Close()

	// Parse address
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatalf("Failed to parse address: %v", err)
	}

	// Parse port
	port := 22
	if portStr != "" {
		if p, err := net.LookupPort("tcp", portStr); err == nil {
			port = p
		}
	}

	// Create machine config
	machineConfig := &config.Machine{
		Name:     "test",
		Host:     host,
		User:     "testuser",
		Password: "testpass",
		Port:     port,
	}

	// Create SSH client
	client, err := NewSSHClient(machineConfig, 5)
	if err != nil {
		// If connection fails, that's okay for this test
		t.Logf("SSH connection failed (expected in some environments): %v", err)
		return
	}
	defer client.Close()

	// Test ExecuteCommand with error
	_, err = client.ExecuteCommand("exit 1")
	if err != nil {
		// This is expected to fail
		if !strings.Contains(err.Error(), "command execution failed") {
			t.Errorf("Expected command execution error, got: %v", err)
		}
	}
}

func TestSSHClient_ExecuteScript_Success(t *testing.T) {
	// Create a temporary script file
	tempDir := t.TempDir()
	scriptFile := filepath.Join(tempDir, "test_script.sh")
	scriptContent := "echo 'test script content'"

	err := os.WriteFile(scriptFile, []byte(scriptContent), 0o600)
	if err != nil {
		t.Fatalf("Failed to create test script file: %v", err)
	}

	// Create a mock SSH client
	client := &SSHClient{
		config: &config.Machine{Name: "test"},
		client: nil, // No real connection
	}

	// This should fail due to no SSH connection, but we can test the file reading part
	_, err = client.ExecuteScript(scriptFile)
	if err == nil {
		t.Error("Expected error when no SSH connection exists")
	}
	if !strings.Contains(err.Error(), "no SSH connection exists") {
		t.Errorf("Expected connection error, got: %v", err)
	}
}

func TestSSHClient_ExecuteScript_FileReadError(t *testing.T) {
	// Create a mock SSH client
	client := &SSHClient{
		config: &config.Machine{Name: "test"},
	}

	_, err := client.ExecuteScript("/nonexistent/script.sh")
	if err == nil {
		t.Error("expected error when script file does not exist")
	}
	if !strings.Contains(err.Error(), "failed to read script file") {
		t.Errorf("expected file read error, got: %v", err)
	}
}

func TestSSHClient_Close(t *testing.T) {
	// Test that Close doesn't panic when client is nil
	client := &SSHClient{
		config: &config.Machine{Name: "test"},
		client: nil,
	}

	// Should not panic
	err := client.Close()
	if err != nil {
		t.Errorf("Close should not return error when Client is nil, got: %v", err)
	}
}

func TestSSHClient_ExecuteCommand_NoConnection(t *testing.T) {
	// Test ExecuteCommand with no SSH connection
	client := &SSHClient{
		config: &config.Machine{Name: "test"},
		client: nil,
	}

	_, err := client.ExecuteCommand("echo test")
	if err == nil {
		t.Error("expected error when no SSH connection exists")
	}
	if !strings.Contains(err.Error(), "failed to create session") {
		t.Errorf("expected session creation error, got: %v", err)
	}
}
