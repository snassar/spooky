package authentication

import (
	"fmt"
	"io/ioutil"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookysshtypes "spooky/internal/types/ssh"

	"golang.org/x/crypto/ssh"
)

// Manager implements AuthenticationEngine interface
type Manager struct {
	config *spookysshtypes.AuthenticationConfig
	logger spookyinterfaces.Logger
}

// NewManager creates a new authentication manager
func NewManager(config *spookysshtypes.AuthenticationConfig, logger spookyinterfaces.Logger) *Manager {
	return &Manager{
		config: config,
		logger: logger,
	}
}

// Authenticate authenticates an SSH connection
func (m *Manager) Authenticate(connection *spookysshtypes.SSHConnection, auth *spookysshtypes.AuthenticationConfig) error {
	if connection == nil {
		return fmt.Errorf("connection cannot be nil")
	}

	if auth == nil {
		return fmt.Errorf("authentication config cannot be nil")
	}

	m.logger.Info("Authenticating SSH connection", spookylogging.String("host", connection.Host))

	// Create SSH client config
	sshConfig := &ssh.ClientConfig{
		User:            connection.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // In production, use proper host key verification
		Timeout:         30 * time.Second,            // Default timeout
	}

	// Add authentication methods
	if auth.Password != "" {
		sshConfig.Auth = append(sshConfig.Auth, ssh.Password(auth.Password))
	}

	if auth.KeyFile != "" {
		key, err := m.loadPrivateKey(auth.KeyFile)
		if err != nil {
			return fmt.Errorf("failed to load private key: %w", err)
		}
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeys(key))
	}

	// Test connection
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", connection.Host, connection.Port), sshConfig)
	if err != nil {
		return fmt.Errorf("SSH authentication failed: %w", err)
	}
	defer client.Close()

	m.logger.Info("SSH authentication successful", spookylogging.String("host", connection.Host))
	return nil
}

// ValidateAuthentication validates authentication configuration
func (m *Manager) ValidateAuthentication(auth *spookysshtypes.AuthenticationConfig) error {
	if auth == nil {
		return fmt.Errorf("authentication config cannot be nil")
	}

	// Check that at least one authentication method is provided
	if auth.Password == "" && auth.KeyFile == "" {
		return fmt.Errorf("either password or key file must be provided")
	}

	// Validate key file if provided
	if auth.KeyFile != "" {
		if err := m.validateKeyFile(auth.KeyFile); err != nil {
			return fmt.Errorf("invalid key file: %w", err)
		}
	}

	return nil
}

// GetSupportedMethods returns supported authentication methods
func (m *Manager) GetSupportedMethods() []string {
	return []string{"password", "key", "mixed"}
}

// loadPrivateKey loads a private key from file
func (m *Manager) loadPrivateKey(keyFile string) (ssh.Signer, error) {
	privateKeyBytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	privateKey, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil
}

// validateKeyFile validates a key file
func (m *Manager) validateKeyFile(keyFile string) error {
	_, err := m.loadPrivateKey(keyFile)
	return err
}
