package ssh

import (
	"fmt"
	"os"
	"strings"
	"time"

	"spooky/internal/config"

	"golang.org/x/crypto/ssh"
)

// NewSSHClient creates a new SSH client connection
func NewSSHClient(server *config.Server, timeout int) (*SSHClient, error) {
	var authMethods []ssh.AuthMethod

	// Add password authentication if provided
	if server.Password != "" {
		authMethods = append(authMethods, ssh.Password(server.Password))
	}

	// Add key-based authentication if provided
	if server.KeyFile != "" {
		key, err := os.ReadFile(server.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read key file %s: %w", server.KeyFile, err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("no authentication method available for server %s", server.Name)
	}

	// SSH client configuration
	sshConfig := &ssh.ClientConfig{
		User:            server.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(timeout) * time.Second,
	}

	// Connect to the server
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", server.Host, server.Port), sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s@%s:%d: %w", server.User, server.Host, server.Port, err)
	}

	return &SSHClient{
		Client: client,
		Server: server,
	}, nil
}

// Close closes the SSH connection
func (s *SSHClient) Close() error {
	return s.Client.Close()
}

// ExecuteCommand executes a command on the remote server
func (s *SSHClient) ExecuteCommand(command string) (string, error) {
	session, err := s.Client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Capture stdout and stderr
	var stdout, stderr strings.Builder
	session.Stdout = &stdout
	session.Stderr = &stderr

	// Execute the command
	err = session.Run(command)
	if err != nil {
		return "", fmt.Errorf("command execution failed: %s, stderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

// ExecuteScript executes a script file on the remote server
func (s *SSHClient) ExecuteScript(scriptPath string) (string, error) {
	// Read the script file
	scriptContent, err := os.ReadFile(scriptPath)
	if err != nil {
		return "", fmt.Errorf("failed to read script file %s: %w", scriptPath, err)
	}

	// Execute the script content
	return s.ExecuteCommand(string(scriptContent))
}
