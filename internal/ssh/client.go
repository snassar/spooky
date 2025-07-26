package ssh

import (
	"bytes"
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
		// Check if the known_hosts file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to parse known_hosts file at %s: file does not exist", path)
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

// NewSSHClient creates a new SSH client for the given machine
func NewSSHClient(machine *config.Machine, timeout int) (*SSHClient, error) {
	return NewSSHClientWithHostKeyCallback(machine, timeout, InsecureHostKey, "")
}

// NewSSHClientWithHostKeyCallback creates a new SSH client with custom host key verification
func NewSSHClientWithHostKeyCallback(machine *config.Machine, timeout int, hostKeyType HostKeyCallbackType, knownHostsPath string) (*SSHClient, error) {
	if machine == nil {
		return nil, fmt.Errorf("machine configuration cannot be nil")
	}

	// Validate timeout value
	if timeout <= 0 {
		return nil, fmt.Errorf("timeout must be positive, got %d", timeout)
	}

	logger := logging.GetLogger()

	logger.Info("Creating SSH client",
		logging.Server(machine.Name),
		logging.Host(machine.Host),
		logging.Port(machine.Port),
		logging.String("user", machine.User),
		logging.Int("timeout_seconds", timeout),
		logging.String("host_key_type", string(hostKeyType)),
	)

	var authMethods []ssh.AuthMethod

	// Add password authentication if provided
	if machine.Password != "" {
		logger.Debug("Adding password authentication",
			logging.Server(machine.Name),
		)
		authMethods = append(authMethods, ssh.Password(machine.Password))
	}

	// Add key-based authentication if provided
	if machine.KeyFile != "" {
		logger.Debug("Adding key-based authentication",
			logging.Server(machine.Name),
			logging.String("key_file", machine.KeyFile),
		)

		key, err := os.ReadFile(machine.KeyFile)
		if err != nil {
			logger.Error("Failed to read SSH key file", err,
				logging.Server(machine.Name),
				logging.String("key_file", machine.KeyFile),
			)
			return nil, fmt.Errorf("failed to read key file %s: %w", machine.KeyFile, err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			logger.Error("Failed to parse SSH private key", err,
				logging.Server(machine.Name),
				logging.String("key_file", machine.KeyFile),
			)
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if len(authMethods) == 0 {
		logger.Error("No authentication method available", fmt.Errorf("no auth methods"),
			logging.Server(machine.Name),
		)
		return nil, fmt.Errorf("no authentication method available for server %s", machine.Name)
	}

	// Get host key callback
	hostKeyCallback, err := getHostKeyCallback(hostKeyType, knownHostsPath)
	if err != nil {
		logger.Error("Failed to create host key callback", err,
			logging.Server(machine.Name),
			logging.String("host_key_type", string(hostKeyType)),
		)
		return nil, fmt.Errorf("failed to create host key callback: %w", err)
	}

	// SSH client configuration
	sshConfig := &ssh.ClientConfig{
		User:            machine.User,
		Auth:            authMethods,
		HostKeyCallback: hostKeyCallback,
		Timeout:         time.Duration(timeout) * time.Second,
	}

	// Connect to the server
	startTime := time.Now()
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", machine.Host, machine.Port), sshConfig)
	if err != nil {
		logger.Error("Failed to establish SSH connection", err,
			logging.Server(machine.Name),
			logging.Host(machine.Host),
			logging.Port(machine.Port),
			logging.String("user", machine.User),
			logging.Duration("duration_ms", time.Since(startTime).Milliseconds()),
		)
		return nil, fmt.Errorf("failed to connect to %s@%s:%d: %w", machine.User, machine.Host, machine.Port, err)
	}

	logger.Info("SSH connection established successfully",
		logging.Server(machine.Name),
		logging.Host(machine.Host),
		logging.Port(machine.Port),
		logging.String("user", machine.User),
		logging.Duration("duration_ms", time.Since(startTime).Milliseconds()),
		logging.Int("auth_methods", len(authMethods)),
	)

	return &SSHClient{
		client: client,
		config: machine,
	}, nil
}

// Connect establishes a connection to the machine
func (c *SSHClient) Connect() error {
	if c.client == nil {
		return fmt.Errorf("SSH client not initialized")
	}
	return nil
}

// Close closes the SSH connection
func (c *SSHClient) Close() error {
	if c.client == nil {
		return nil
	}
	return c.client.Close()
}

// GetMachine returns the machine configuration
func (c *SSHClient) GetMachine() *config.Machine {
	return c.config
}

// ExecuteCommand executes a command on the remote server
func (c *SSHClient) ExecuteCommand(command string) (string, error) {
	logger := logging.GetLogger()

	if c.client == nil {
		logger.Error("No SSH connection available", fmt.Errorf("client is nil"),
			logging.Server(c.config.Name),
		)
		return "", fmt.Errorf("failed to create session: no SSH connection exists (Client is nil)")
	}

	logger.Debug("Creating SSH session",
		logging.Server(c.config.Name),
		logging.String("command_length", fmt.Sprintf("%d chars", len(command))),
	)

	session, err := c.client.NewSession()
	if err != nil {
		logger.Error("Failed to create SSH session", err,
			logging.Server(c.config.Name),
		)
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	if err := session.Run(command); err != nil {
		logger.Error("Command execution failed", err,
			logging.Server(c.config.Name),
			logging.String("command_length", fmt.Sprintf("%d chars", len(command))),
			logging.String("stderr", stderr.String()),
		)
		return "", fmt.Errorf("command execution failed: %w", err)
	}

	output := stdout.String()
	logger.Debug("Command executed successfully",
		logging.Server(c.config.Name),
		logging.String("command_length", fmt.Sprintf("%d chars", len(command))),
		logging.String("output_length", fmt.Sprintf("%d chars", len(output))),
	)

	return output, nil
}

// ExecuteScript executes a script file on the remote server
func (c *SSHClient) ExecuteScript(scriptPath string) (string, error) {
	logger := logging.GetLogger()

	logger.Info("Loading script file",
		logging.Server(c.config.Name),
		logging.String("script_path", scriptPath),
	)

	scriptContent, err := os.ReadFile(scriptPath)
	if err != nil {
		logger.Error("Failed to read script file", err,
			logging.Server(c.config.Name),
			logging.String("script_path", scriptPath),
		)
		return "", fmt.Errorf("failed to read script file %s: %w", scriptPath, err)
	}

	logger.Debug("Script file loaded successfully",
		logging.Server(c.config.Name),
		logging.String("script_path", scriptPath),
		logging.String("script_size", fmt.Sprintf("%d bytes", len(scriptContent))),
	)

	// Execute the script content
	return c.ExecuteCommand(string(scriptContent))
}
