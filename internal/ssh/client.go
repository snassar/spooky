package ssh

import (
	"fmt"
	"os"
	"strings"
	"time"

	"spooky/internal/config"
	"spooky/internal/logging"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// HostKeyCallbackType defines the type of host key verification to use
type HostKeyCallbackType string

const (
	// InsecureHostKey allows connections without host key verification (for testing only)
	InsecureHostKey HostKeyCallbackType = "insecure"
	// KnownHostsHostKey uses known_hosts file for verification
	KnownHostsHostKey HostKeyCallbackType = "known_hosts"
	// AutoHostKey automatically accepts new host keys (less secure but more convenient)
	AutoHostKey HostKeyCallbackType = "auto"
)

// getHostKeyCallback returns the appropriate host key callback based on the type
func getHostKeyCallback(callbackType HostKeyCallbackType, knownHostsPath string) (ssh.HostKeyCallback, error) {
	switch callbackType {
	case InsecureHostKey:
		// Note: This should only be used in testing environments
		//nolint:gosec // InsecureIgnoreHostKey is intentional for testing mode
		return ssh.InsecureIgnoreHostKey(), nil
	case KnownHostsHostKey:
		path := knownHostsPath
		if path == "" {
			path = "~/.ssh/known_hosts"
		}
		// Expand ~ to home directory
		if strings.HasPrefix(path, "~") {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("failed to get home directory: %w", err)
			}
			path = strings.Replace(path, "~", homeDir, 1)
		}
		// Use golang.org/x/crypto/ssh/knownhosts for host key verification
		hostKeyCallback, err := knownhosts.New(path)
		if err != nil {
			return nil, fmt.Errorf("failed to parse known_hosts file at %s: %w", path, err)
		}
		return hostKeyCallback, nil
	case AutoHostKey:
		// Auto-accept new host keys (less secure but more convenient)
		//nolint:gosec // InsecureIgnoreHostKey is intentional for auto-accept mode
		return ssh.InsecureIgnoreHostKey(), nil
	default:
		return nil, fmt.Errorf("unsupported host key callback type: %s", callbackType)
	}
}

// NewSSHClient creates a new SSH client connection
func NewSSHClient(server *config.Server, timeout int) (*SSHClient, error) {
	return NewSSHClientWithHostKeyCallback(server, timeout, InsecureHostKey, "")
}

// NewSSHClientWithHostKeyCallback creates a new SSH client with custom host key verification
func NewSSHClientWithHostKeyCallback(server *config.Server, timeout int, hostKeyType HostKeyCallbackType, knownHostsPath string) (*SSHClient, error) {
	logger := logging.GetLogger()

	logger.Info("Creating SSH client",
		logging.Server(server.Name),
		logging.Host(server.Host),
		logging.Port(server.Port),
		logging.String("user", server.User),
		logging.Int("timeout_seconds", timeout),
		logging.String("host_key_type", string(hostKeyType)),
	)

	var authMethods []ssh.AuthMethod

	// Add password authentication if provided
	if server.Password != "" {
		logger.Debug("Adding password authentication",
			logging.Server(server.Name),
		)
		authMethods = append(authMethods, ssh.Password(server.Password))
	}

	// Add key-based authentication if provided
	if server.KeyFile != "" {
		logger.Debug("Adding key-based authentication",
			logging.Server(server.Name),
			logging.String("key_file", server.KeyFile),
		)

		key, err := os.ReadFile(server.KeyFile)
		if err != nil {
			logger.Error("Failed to read SSH key file", err,
				logging.Server(server.Name),
				logging.String("key_file", server.KeyFile),
			)
			return nil, fmt.Errorf("failed to read key file %s: %w", server.KeyFile, err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			logger.Error("Failed to parse SSH private key", err,
				logging.Server(server.Name),
				logging.String("key_file", server.KeyFile),
			)
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if len(authMethods) == 0 {
		logger.Error("No authentication method available", fmt.Errorf("no auth methods"),
			logging.Server(server.Name),
		)
		return nil, fmt.Errorf("no authentication method available for server %s", server.Name)
	}

	// Get host key callback
	hostKeyCallback, err := getHostKeyCallback(hostKeyType, knownHostsPath)
	if err != nil {
		logger.Error("Failed to create host key callback", err,
			logging.Server(server.Name),
			logging.String("host_key_type", string(hostKeyType)),
		)
		return nil, fmt.Errorf("failed to create host key callback: %w", err)
	}

	// SSH client configuration
	sshConfig := &ssh.ClientConfig{
		User:            server.User,
		Auth:            authMethods,
		HostKeyCallback: hostKeyCallback,
		Timeout:         time.Duration(timeout) * time.Second,
	}

	// Connect to the server
	startTime := time.Now()
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", server.Host, server.Port), sshConfig)
	if err != nil {
		logger.Error("Failed to establish SSH connection", err,
			logging.Server(server.Name),
			logging.Host(server.Host),
			logging.Port(server.Port),
			logging.String("user", server.User),
			logging.Duration("duration_ms", time.Since(startTime).Milliseconds()),
		)
		return nil, fmt.Errorf("failed to connect to %s@%s:%d: %w", server.User, server.Host, server.Port, err)
	}

	logger.Info("SSH connection established successfully",
		logging.Server(server.Name),
		logging.Host(server.Host),
		logging.Port(server.Port),
		logging.String("user", server.User),
		logging.Duration("duration_ms", time.Since(startTime).Milliseconds()),
		logging.Int("auth_methods", len(authMethods)),
	)

	return &SSHClient{
		Client: client,
		Server: server,
	}, nil
}

// Close closes the SSH connection
func (s *SSHClient) Close() error {
	if s.Client == nil {
		return nil
	}
	return s.Client.Close()
}

// ExecuteCommand executes a command on the remote server
func (s *SSHClient) ExecuteCommand(command string) (string, error) {
	logger := logging.GetLogger()

	if s.Client == nil {
		logger.Error("No SSH connection available", fmt.Errorf("client is nil"),
			logging.Server(s.Server.Name),
		)
		return "", fmt.Errorf("failed to create session: no SSH connection exists (Client is nil)")
	}

	logger.Debug("Creating SSH session",
		logging.Server(s.Server.Name),
		logging.String("command_length", fmt.Sprintf("%d chars", len(command))),
	)

	session, err := s.Client.NewSession()
	if err != nil {
		logger.Error("Failed to create SSH session", err,
			logging.Server(s.Server.Name),
		)
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Capture stdout and stderr
	var stdout, stderr strings.Builder
	session.Stdout = &stdout
	session.Stderr = &stderr

	// Execute the command
	startTime := time.Now()
	err = session.Run(command)
	if err != nil {
		logger.Error("Command execution failed", err,
			logging.Server(s.Server.Name),
			logging.String("command_length", fmt.Sprintf("%d chars", len(command))),
			logging.String("stderr", stderr.String()),
			logging.Duration("duration_ms", time.Since(startTime).Milliseconds()),
		)
		return "", fmt.Errorf("command execution failed: %s, stderr: %s", err, stderr.String())
	}

	output := stdout.String()
	logger.Debug("Command executed successfully",
		logging.Server(s.Server.Name),
		logging.String("command_length", fmt.Sprintf("%d chars", len(command))),
		logging.String("output_length", fmt.Sprintf("%d chars", len(output))),
		logging.Duration("duration_ms", time.Since(startTime).Milliseconds()),
	)

	return output, nil
}

// ExecuteScript executes a script file on the remote server
func (s *SSHClient) ExecuteScript(scriptPath string) (string, error) {
	logger := logging.GetLogger()

	logger.Info("Loading script file",
		logging.Server(s.Server.Name),
		logging.String("script_path", scriptPath),
	)

	// Read the script file
	scriptContent, err := os.ReadFile(scriptPath)
	if err != nil {
		logger.Error("Failed to read script file", err,
			logging.Server(s.Server.Name),
			logging.String("script_path", scriptPath),
		)
		return "", fmt.Errorf("failed to read script file %s: %w", scriptPath, err)
	}

	logger.Debug("Script file loaded successfully",
		logging.Server(s.Server.Name),
		logging.String("script_path", scriptPath),
		logging.String("script_size", fmt.Sprintf("%d bytes", len(scriptContent))),
	)

	// Execute the script content
	return s.ExecuteCommand(string(scriptContent))
}
