package testing

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

// TestSSHServer provides a mock SSH server for testing
type TestSSHServer struct {
	listener net.Listener
	config   *ssh.ServerConfig
	clients  []*ssh.ServerConn
	mu       sync.RWMutex
	port     int
}

// TestFileSystem provides a temporary file system for testing
type TestFileSystem struct {
	rootDir string
	files   map[string][]byte
	mu      sync.RWMutex
}

// TestNetwork provides network condition simulation
type TestNetwork struct {
	latency    time.Duration
	packetLoss float64
	bandwidth  int64 // bytes per second
}

// TestLogger provides test-specific logging utilities
type TestLogger struct {
	t *testing.T
}

// NewTestSSHServer creates a new mock SSH server for testing
func NewTestSSHServer() (*TestSSHServer, error) {
	// Generate a test host key
	hostKey, err := generateTestHostKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate test host key: %w", err)
	}

	config := &ssh.ServerConfig{
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			// Accept any public key for testing
			return &ssh.Permissions{}, nil
		},
		PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			// Accept any password for testing
			return &ssh.Permissions{}, nil
		},
	}
	config.AddHostKey(hostKey)

	// Find an available port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	addr := listener.Addr().(*net.TCPAddr)
	port := addr.Port

	server := &TestSSHServer{
		listener: listener,
		config:   config,
		clients:  make([]*ssh.ServerConn, 0),
		port:     port,
	}

	// Start accepting connections
	go server.acceptConnections()

	return server, nil
}

// Port returns the port the server is listening on
func (s *TestSSHServer) Port() int {
	return s.port
}

// Close shuts down the SSH server
func (s *TestSSHServer) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Close all client connections
	for _, client := range s.clients {
		if err := client.Close(); err != nil {
			// Log error but continue cleanup
			fmt.Printf("Warning: failed to close SSH client: %v\n", err)
		}
	}

	return s.listener.Close()
}

// acceptConnections handles incoming SSH connections
func (s *TestSSHServer) acceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			// Server is being closed
			return
		}

		go s.handleConnection(conn)
	}
}

// handleConnection handles a single SSH connection
func (s *TestSSHServer) handleConnection(conn net.Conn) {
	// Perform SSH handshake
	serverConn, chans, reqs, err := ssh.NewServerConn(conn, s.config)
	if err != nil {
		if closeErr := conn.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close connection after handshake error: %v\n", closeErr)
		}
		return
	}

	s.mu.Lock()
	s.clients = append(s.clients, serverConn)
	s.mu.Unlock()

	// Handle global requests
	go ssh.DiscardRequests(reqs)

	// Handle channels
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			if err := newChannel.Reject(ssh.UnknownChannelType, "unknown channel type"); err != nil {
				fmt.Printf("Warning: failed to reject unknown channel type: %v\n", err)
			}
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			continue
		}

		go s.handleChannel(channel, requests)
	}
}

// handleChannel handles SSH channel requests
func (s *TestSSHServer) handleChannel(channel ssh.Channel, requests <-chan *ssh.Request) {
	defer func() {
		if err := channel.Close(); err != nil {
			fmt.Printf("Warning: failed to close SSH channel: %v\n", err)
		}
	}()

	for req := range requests {
		switch req.Type {
		case "exec":
			// Handle command execution
			if len(req.Payload) > 4 {
				command := string(req.Payload[4:])
				s.executeCommand(channel, command)
			}
			if err := req.Reply(true, nil); err != nil {
				fmt.Printf("Warning: failed to reply to exec request: %v\n", err)
			}
		case "shell":
			// Handle shell request
			if err := req.Reply(true, nil); err != nil {
				fmt.Printf("Warning: failed to reply to shell request: %v\n", err)
			}
		default:
			if err := req.Reply(false, nil); err != nil {
				fmt.Printf("Warning: failed to reply to unknown request: %v\n", err)
			}
		}
	}
}

// executeCommand executes a command and returns the result
func (s *TestSSHServer) executeCommand(channel ssh.Channel, command string) {
	// Simple command simulation for testing
	output := s.simulateCommand(command)
	if _, err := channel.Write([]byte(output)); err != nil {
		fmt.Printf("Warning: failed to write command output: %v\n", err)
	}
	if _, err := channel.SendRequest("exit-status", false, []byte{0, 0, 0, 0}); err != nil {
		fmt.Printf("Warning: failed to send exit status: %v\n", err)
	}
}

// simulateCommand simulates command execution
func (s *TestSSHServer) simulateCommand(command string) string {
	// Simple command simulation
	switch {
	case strings.HasPrefix(command, "echo "):
		return strings.TrimPrefix(command, "echo ") + "\n"
	case command == "pwd":
		return "/home/test\n"
	case command == "whoami":
		return "testuser\n"
	case strings.HasPrefix(command, "ls "):
		return "file1.txt\nfile2.txt\n"
	case command == "ls":
		return "file1.txt\nfile2.txt\n"
	default:
		return fmt.Sprintf("Command executed: %s\n", command)
	}
}

// NewTestFileSystem creates a new temporary file system for testing
func NewTestFileSystem(t *testing.T) *TestFileSystem {
	rootDir, err := os.MkdirTemp("", "test-filesystem-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	return &TestFileSystem{
		rootDir: rootDir,
		files:   make(map[string][]byte),
	}
}

// RootDir returns the root directory of the test file system
func (fs *TestFileSystem) RootDir() string {
	return fs.rootDir
}

// CreateFile creates a file in the test file system
func (fs *TestFileSystem) CreateFile(path string, content []byte) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fullPath := filepath.Join(fs.rootDir, path)
	dir := filepath.Dir(fullPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fs.files[path] = content
	return nil
}

// ReadFile reads a file from the test file system
func (fs *TestFileSystem) ReadFile(path string) ([]byte, error) {
	fullPath := filepath.Join(fs.rootDir, path)
	return os.ReadFile(fullPath)
}

// FileExists checks if a file exists in the test file system
func (fs *TestFileSystem) FileExists(path string) bool {
	fullPath := filepath.Join(fs.rootDir, path)
	_, err := os.Stat(fullPath)
	return !os.IsNotExist(err)
}

// ListFiles lists all files in the test file system
func (fs *TestFileSystem) ListFiles() ([]string, error) {
	var files []string
	err := filepath.Walk(fs.rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(fs.rootDir, path)
			if err != nil {
				return err
			}
			files = append(files, relPath)
		}
		return nil
	})
	return files, err
}

// Cleanup removes the test file system
func (fs *TestFileSystem) Cleanup() error {
	return os.RemoveAll(fs.rootDir)
}

// NewTestNetwork creates a new network condition simulator
func NewTestNetwork(latency time.Duration, packetLoss float64, bandwidth int64) *TestNetwork {
	return &TestNetwork{
		latency:    latency,
		packetLoss: packetLoss,
		bandwidth:  bandwidth,
	}
}

// SimulateLatency simulates network latency
func (n *TestNetwork) SimulateLatency() {
	if n.latency > 0 {
		time.Sleep(n.latency)
	}
}

// NewTestLogger creates a new test logger
func NewTestLogger(t *testing.T) *TestLogger {
	return &TestLogger{t: t}
}

// Log logs a message with the test context
func (l *TestLogger) Log(message string, args ...interface{}) {
	l.t.Logf(message, args...)
}

// Logf logs a formatted message with the test context
func (l *TestLogger) Logf(format string, args ...interface{}) {
	l.t.Logf(format, args...)
}

// generateTestHostKey generates a test host key for the SSH server
func generateTestHostKey() (ssh.Signer, error) {
	// For testing purposes, we'll create a simple RSA key
	// In a real implementation, you might want to use a pre-generated key
	// or generate one with proper entropy
	return nil, fmt.Errorf("host key generation not implemented for testing")
}

// TestSSHClient provides a test SSH client that connects to the mock server
type TestSSHClient struct {
	server *TestSSHServer
	client *ssh.Client
}

// NewTestSSHClient creates a new test SSH client
func NewTestSSHClient(server *TestSSHServer) (*TestSSHClient, error) {
	config := &ssh.ClientConfig{
		User: "testuser",
		Auth: []ssh.AuthMethod{
			ssh.Password("testpass"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For testing only
	}

	addr := fmt.Sprintf("localhost:%d", server.Port())
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH server: %w", err)
	}

	return &TestSSHClient{
		server: server,
		client: client,
	}, nil
}

// RunCommand executes a command on the test SSH server
func (c *TestSSHClient) RunCommand(command string) (string, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			fmt.Printf("Warning: failed to close SSH session: %v\n", err)
		}
	}()

	output, err := session.CombinedOutput(command)
	return string(output), err
}

// Close closes the SSH client connection
func (c *TestSSHClient) Close() error {
	return c.client.Close()
}

// TestContext provides a test context with timeout
func TestContext(t *testing.T, timeout time.Duration) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(cancel)
	return ctx
}

// AssertNoError asserts that an error is nil
func AssertNoError(t *testing.T, err error, message string) {
	if err != nil {
		t.Fatalf("%s: %v", message, err)
	}
}

// AssertError asserts that an error is not nil
func AssertError(t *testing.T, err error, message string) {
	if err == nil {
		t.Fatalf("%s: expected error but got none", message)
	}
}

// AssertEqual asserts that two values are equal
func AssertEqual(t *testing.T, expected, actual interface{}, message string) {
	if expected != actual {
		t.Fatalf("%s: expected %v, got %v", message, expected, actual)
	}
}

// AssertNotNil asserts that a value is not nil
func AssertNotNil(t *testing.T, value interface{}, message string) {
	if value == nil {
		t.Fatalf("%s: expected non-nil value", message)
	}
}
