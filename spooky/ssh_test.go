package spooky

import (
	"errors"
	"io"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"
)

// --- Mock interfaces for dependency injection ---
type fileReader interface {
	ReadFile(filename string) ([]byte, error)
}

type keyParser interface {
	ParsePrivateKey(key []byte) (ssh.Signer, error)
}

type sshDialer interface {
	Dial(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error)
}

// --- Mock implementations ---
type mockFileReader struct {
	shouldFail bool
	content    []byte
}

func (m *mockFileReader) ReadFile(filename string) ([]byte, error) {
	if m.shouldFail {
		return nil, errors.New("file read error")
	}
	return m.content, nil
}

type mockKeyParser struct {
	shouldFail bool
}

func (m *mockKeyParser) ParsePrivateKey(key []byte) (ssh.Signer, error) {
	if m.shouldFail {
		return nil, errors.New("key parse error")
	}
	// Return a mock signer
	return &mockSigner{}, nil
}

type mockSigner struct{}

func (m *mockSigner) PublicKey() ssh.PublicKey { return nil }
func (m *mockSigner) Sign(rand io.Reader, data []byte) (*ssh.Signature, error) {
	return &ssh.Signature{}, nil
}

type mockSSHDialer struct {
	shouldFail bool
}

func (m *mockSSHDialer) Dial(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error) {
	if m.shouldFail {
		return nil, errors.New("ssh dial error")
	}
	return &ssh.Client{}, nil
}

// --- Test helpers ---
func createTestServer() *Server {
	return &Server{
		Name:     "testserver",
		Host:     "localhost",
		User:     "testuser",
		Port:     22,
		Password: "testpass",
	}
}

// --- Tests for NewSSHClient ---
func TestNewSSHClient_NoAuth(t *testing.T) {
	server := &Server{
		Name: "testserver",
		Host: "localhost",
		User: "testuser",
		// No password or key file
	}

	_, err := NewSSHClient(server, 30)
	if err == nil {
		t.Error("expected error for no authentication method")
	}
	if !strings.Contains(err.Error(), "no authentication method available") {
		t.Errorf("expected authentication error, got: %v", err)
	}
}

func TestNewSSHClient_PasswordAuth(t *testing.T) {
	server := createTestServer()
	server.Password = "testpass"

	// This will fail due to real SSH connection, but we're testing the auth method setup
	_, err := NewSSHClient(server, 1) // Short timeout
	if err != nil {
		// Expected to fail due to connection, but should not fail due to auth method setup
		if strings.Contains(err.Error(), "no authentication method available") {
			t.Error("password authentication method should be available")
		}
		// Other errors (like connection timeout) are expected in unit tests
	}
}

func TestNewSSHClient_KeyFileAuth(t *testing.T) {
	server := createTestServer()
	server.Password = "" // Clear password
	server.KeyFile = "/path/to/key"

	// This will fail due to real SSH connection, but we're testing the auth method setup
	_, err := NewSSHClient(server, 1) // Short timeout
	if err != nil {
		// Expected to fail due to connection or file read, but should not fail due to auth method setup
		if strings.Contains(err.Error(), "no authentication method available") {
			t.Error("key file authentication method should be available")
		}
		// Other errors (like file not found or connection timeout) are expected in unit tests
	}
}

func TestNewSSHClient_BothAuthMethods(t *testing.T) {
	server := createTestServer()
	server.Password = "testpass"
	server.KeyFile = "/path/to/key"

	// This will fail due to real SSH connection, but we're testing the auth method setup
	_, err := NewSSHClient(server, 1) // Short timeout
	if err != nil {
		// Expected to fail due to connection, but should not fail due to auth method setup
		if strings.Contains(err.Error(), "no authentication method available") {
			t.Error("both authentication methods should be available")
		}
		// Other errors (like connection timeout) are expected in unit tests
	}
}

// --- Tests for ExecuteCommand ---
// Note: Testing ExecuteCommand with real SSH connections requires mocking
// or integration tests. For unit tests, we focus on error conditions.

// --- Tests for ExecuteScript ---
func TestExecuteScript_FileReadError(t *testing.T) {
	client := &SSHClient{
		client: &ssh.Client{},
		server: createTestServer(),
	}

	// Test with non-existent file
	_, err := client.ExecuteScript("/nonexistent/script.sh")
	if err == nil {
		t.Error("expected error for non-existent script file")
	}
	if !strings.Contains(err.Error(), "failed to read script file") {
		t.Errorf("expected file read error, got: %v", err)
	}
}

// --- Tests for indentOutput ---
func TestIndentOutput_Empty(t *testing.T) {
	result := indentOutput("")
	if result != "" {
		t.Errorf("expected empty string, got: %q", result)
	}
}

func TestIndentOutput_SingleLine(t *testing.T) {
	input := "hello world"
	expected := "    hello world\n"
	result := indentOutput(input)
	if result != expected {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, result)
	}
}

func TestIndentOutput_MultipleLines(t *testing.T) {
	input := "line1\nline2\nline3"
	expected := "    line1\n    line2\n    line3\n"
	result := indentOutput(input)
	if result != expected {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, result)
	}
}

func TestIndentOutput_WithEmptyLines(t *testing.T) {
	input := "line1\n\nline3"
	expected := "    line1\n    \n    line3\n"
	result := indentOutput(input)
	if result != expected {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, result)
	}
}

func TestIndentOutput_WithTrailingNewline(t *testing.T) {
	input := "line1\nline2\n"
	expected := "    line1\n    line2\n" // The function doesn't add an extra newline
	result := indentOutput(input)
	if result != expected {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, result)
	}
}

// --- Tests for ExecuteConfig ---
func TestExecuteConfig_EmptyActions(t *testing.T) {
	config := &Config{
		Servers: []Server{
			{
				Name:     "testserver",
				Host:     "localhost",
				User:     "testuser",
				Password: "testpass",
			},
		},
		Actions: []Action{}, // Empty actions
	}

	err := ExecuteConfig(config)
	if err != nil {
		t.Errorf("expected no error for empty actions, got: %v", err)
	}
}

func TestExecuteConfig_GetServersError(t *testing.T) {
	config := &Config{
		Servers: []Server{
			{
				Name:     "testserver",
				Host:     "localhost",
				User:     "testuser",
				Password: "testpass",
			},
		},
		Actions: []Action{
			{
				Name:    "testaction",
				Command: "echo hello",
				Servers: []string{"nonexistent"}, // This will cause GetServersForAction to fail
			},
		},
	}

	err := ExecuteConfig(config)
	if err == nil {
		t.Error("expected error for non-existent server")
	}
	if !strings.Contains(err.Error(), "failed to get servers for action") {
		t.Errorf("expected server error, got: %v", err)
	}
}

// --- Tests for executeActionSequential ---
func TestExecuteActionSequential_EmptyServers(t *testing.T) {
	action := &Action{
		Name:    "testaction",
		Command: "echo hello",
	}
	servers := []*Server{} // Empty servers list

	err := executeActionSequential(action, servers)
	if err != nil {
		t.Errorf("expected no error for empty servers, got: %v", err)
	}
}

// --- Tests for executeActionParallel ---
func TestExecuteActionParallel_EmptyServers(t *testing.T) {
	action := &Action{
		Name:    "testaction",
		Command: "echo hello",
	}
	servers := []*Server{} // Empty servers list

	err := executeActionParallel(action, servers)
	if err != nil {
		t.Errorf("expected no error for empty servers, got: %v", err)
	}
}

// --- Integration-style tests (skipped in unit test environment) ---
func TestExecuteActionSequential_WithRealServer(t *testing.T) {
	t.Skip("Skipping integration test that requires real SSH server")

	action := &Action{
		Name:    "testaction",
		Command: "echo hello",
	}
	servers := []*Server{
		{
			Name:     "testserver",
			Host:     "localhost",
			User:     "testuser",
			Password: "testpass",
			Port:     22,
		},
	}

	// This will fail due to SSH connection, but we're testing the function structure
	err := executeActionSequential(action, servers)
	if err == nil {
		// If it succeeds, that's fine (unlikely in test environment)
	} else {
		// Expected to fail due to SSH connection issues
		if !strings.Contains(err.Error(), "Failed to connect") &&
			!strings.Contains(err.Error(), "ssh:") &&
			!strings.Contains(err.Error(), "dial tcp") {
			t.Errorf("unexpected error type: %v", err)
		}
	}
}

func TestExecuteActionParallel_WithRealServer(t *testing.T) {
	t.Skip("Skipping integration test that requires real SSH server")

	action := &Action{
		Name:    "testaction",
		Command: "echo hello",
	}
	servers := []*Server{
		{
			Name:     "testserver",
			Host:     "localhost",
			User:     "testuser",
			Password: "testpass",
			Port:     22,
		},
	}

	// This will fail due to SSH connection, but we're testing the function structure
	err := executeActionParallel(action, servers)
	if err == nil {
		// If it succeeds, that's fine (unlikely in test environment)
	} else {
		// Expected to fail due to SSH connection issues
		if !strings.Contains(err.Error(), "failed to connect") &&
			!strings.Contains(err.Error(), "ssh:") &&
			!strings.Contains(err.Error(), "dial tcp") {
			t.Errorf("unexpected error type: %v", err)
		}
	}
}
