package keys

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"spooky/internal/logging"
	"spooky/internal/ssh/types"

	"golang.org/x/crypto/ssh"
)

// Manager implements SSHKeyManager interface
type Manager struct {
	config *types.KeysConfig
	logger logging.Logger
}

// NewManager creates a new keys manager
func NewManager(config *types.KeysConfig, logger logging.Logger) *Manager {
	return &Manager{
		config: config,
		logger: logger,
	}
}

// LoadPrivateKey loads a private key from file
func (m *Manager) LoadPrivateKey(path string) (*types.SSHKey, error) {
	if path == "" {
		return nil, fmt.Errorf("private key path cannot be empty")
	}

	// Read private key file
	privateKeyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	// Parse private key to validate it
	_, err = ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Extract public key
	signer, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer: %w", err)
	}

	publicKeyBytes := ssh.MarshalAuthorizedKey(signer.PublicKey())

	sshKey := &types.SSHKey{
		Type:       "rsa", // Default type
		PrivateKey: privateKeyBytes,
		PublicKey:  publicKeyBytes,
		Path:       path,
		CreatedAt:  time.Now(),
	}

	m.logger.Info("Private key loaded", logging.String("path", path))
	return sshKey, nil
}

// ValidateKeyFile validates a key file
func (m *Manager) ValidateKeyFile(path string) error {
	if path == "" {
		return fmt.Errorf("key file path cannot be empty")
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("key file does not exist: %s", path)
	}

	// Try to load the key to validate it
	_, err := m.LoadPrivateKey(path)
	if err != nil {
		return fmt.Errorf("invalid key file: %w", err)
	}

	m.logger.Info("Key file validated", logging.String("path", path))
	return nil
}
