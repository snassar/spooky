package secrets

import (
	"fmt"
	"strings"

	"spooky/internal/logging"
	"spooky/internal/templates/types"
)

// Manager implements SecretsManager interface
type Manager struct {
	config             *types.SecretsConfig
	encryptionProvider EncryptionProvider
	logger             logging.Logger
}

// NewManager creates a new secrets manager
func NewManager(config *types.SecretsConfig, logger logging.Logger) *Manager {
	return &Manager{
		config:             config,
		encryptionProvider: NewAgeEncryptionProvider(),
		logger:             logger,
	}
}

// EncryptValue encrypts a value for template use
func (m *Manager) EncryptValue(value string) (string, error) {
	if !m.config.Enabled {
		return value, nil
	}

	// Encrypt value
	encrypted, err := m.encryptionProvider.Encrypt([]byte(value))
	if err != nil {
		return "", fmt.Errorf("encryption failed: %w", err)
	}

	// Return encrypted value with prefix
	return fmt.Sprintf("encrypted:%s", string(encrypted)), nil
}

// DecryptValue decrypts an encrypted value
func (m *Manager) DecryptValue(encryptedValue string) (string, error) {
	if !m.config.Enabled {
		return encryptedValue, nil
	}

	// Check if value is encrypted
	if !m.IsEncrypted(encryptedValue) {
		return encryptedValue, nil
	}

	// Extract encrypted data
	encryptedData := strings.TrimPrefix(encryptedValue, "encrypted:")

	// Decrypt value
	decrypted, err := m.encryptionProvider.Decrypt([]byte(encryptedData))
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	return string(decrypted), nil
}

// ValidateEncryptionKey validates an encryption key
func (m *Manager) ValidateEncryptionKey(key string) error {
	if key == "" {
		return fmt.Errorf("encryption key cannot be empty")
	}

	// Basic validation - key should be at least 32 characters
	if len(key) < 32 {
		return fmt.Errorf("encryption key must be at least 32 characters long")
	}

	return nil
}

// SetEncryptionKey sets the encryption key
func (m *Manager) SetEncryptionKey(key string) error {
	if err := m.ValidateEncryptionKey(key); err != nil {
		return fmt.Errorf("invalid encryption key: %w", err)
	}

	if m.config == nil {
		m.config = &types.SecretsConfig{}
	}
	m.config.EncryptionKey = key
	return nil
}

// SetEncryptionAlgorithm sets the encryption algorithm
func (m *Manager) SetEncryptionAlgorithm(algorithm string) error {
	if m.config == nil {
		m.config = &types.SecretsConfig{}
	}
	m.config.EncryptionAlgorithm = algorithm
	return nil
}

// EnableEncryption enables or disables encryption
func (m *Manager) EnableEncryption(enabled bool) error {
	if m.config == nil {
		m.config = &types.SecretsConfig{}
	}
	m.config.Enabled = enabled
	return nil
}

// IsEncrypted checks if a value is encrypted
func (m *Manager) IsEncrypted(value string) bool {
	return strings.HasPrefix(value, "encrypted:")
}

// GetEncryptionAlgorithm returns the encryption algorithm
func (m *Manager) GetEncryptionAlgorithm() string {
	if m.config == nil {
		return "age"
	}
	return m.config.EncryptionAlgorithm
}

// Close closes the secrets manager
func (m *Manager) Close() error {
	// Cleanup resources if needed
	return nil
}
