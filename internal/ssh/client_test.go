package ssh

import (
	"errors"
	"io"
	"strings"
	"testing"

	"spooky/internal/config"

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
func createTestServer() *config.Server {
	return &config.Server{
		Name:     "testserver",
		Host:     "localhost",
		User:     "testuser",
		Port:     22,
		Password: "testpass",
	}
}

// --- Tests for NewSSHClient ---
func TestNewSSHClient_NoAuth(t *testing.T) {
	server := &config.Server{
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
		Client: &ssh.Client{},
		Server: createTestServer(),
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
