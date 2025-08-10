package secrets

import (
	"fmt"
	"os"
	"path/filepath"

	spookylogging "spooky/internal/logging"
	spookysecretsage "spooky/internal/secrets/age"
	spookysecretskeys "spooky/internal/secrets/keys"
	spookysecretsops "spooky/internal/secrets/operations"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypessecrets "spooky/internal/types/secrets"
)

// Manager implements the SecretsManager interface
type Manager struct {
	config      *spookytypessecrets.SecretsConfig
	ageClient   spookysecretsage.AgeClient
	keyManager  *spookysecretskeys.Manager
	fileManager *spookysecretsops.FileManager
	logger      spookytypeslogging.Logger
}

// NewManager creates a new secrets manager
func NewManager(config *spookytypessecrets.SecretsConfig) (*Manager, error) {
	if config == nil {
		config = &spookytypessecrets.SecretsConfig{
			Enabled: true,
		}
	}

	// Set default values
	if config.Encryption.Algorithm == "" {
		config.Encryption.Algorithm = "age"
	}

	// Create age client
	ageClient := spookysecretsage.NewClient()

	// Create key manager
	globalKeysPath := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "spooky", "keys")
	if globalKeysPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, &spookytypessecrets.SecretsError{
				Operation: "get_home_directory",
				Cause:     err,
			}
		}
		globalKeysPath = filepath.Join(homeDir, ".config", "spooky", "keys")
	}

	keyManager := spookysecretskeys.NewManager(globalKeysPath, "")

	// Create file manager
	fileManager := spookysecretsops.NewFileManager()

	manager := &Manager{
		config:      config,
		ageClient:   ageClient,
		keyManager:  keyManager,
		fileManager: fileManager,
		logger:      spookylogging.GetLogger(),
	}

	// Validate configuration
	if err := manager.Validate(); err != nil {
		return nil, err
	}

	return manager, nil
}

// EncryptVariable encrypts a variable value
func (m *Manager) EncryptVariable(name string, value string, recipients []string) (*spookytypessecrets.EncryptedValue, error) {
	m.logger.Info("Encrypting variable", spookylogging.String("name", name))

	// Use default recipients if none provided
	if len(recipients) == 0 {
		recipients = m.config.Keys.DefaultRecipients
	}

	// Validate recipients
	if len(recipients) == 0 {
		return nil, &spookytypessecrets.SecretsError{
			Operation: "validate_recipients",
			Cause:     fmt.Errorf("no recipients specified"),
		}
	}

	// Encrypt the value
	encrypted, err := m.ageClient.Encrypt([]byte(value), recipients)
	if err != nil {
		return nil, &spookytypessecrets.SecretsError{
			Operation: "encrypt_variable",
			Cause:     err,
			Context: map[string]interface{}{
				"variable": name,
			},
		}
	}

	m.logger.Info("Variable encrypted successfully")
	return encrypted, nil
}

// DecryptVariable decrypts a variable value
func (m *Manager) DecryptVariable(encrypted *spookytypessecrets.EncryptedValue, identity string) (string, error) {
	m.logger.Info("Decrypting variable")

	// Validate identity
	if identity == "" {
		return "", &spookytypessecrets.SecretsError{
			Operation: "validate_identity",
			Cause:     fmt.Errorf("no identity specified"),
		}
	}

	// Decrypt the value
	decrypted, err := m.ageClient.Decrypt(encrypted, identity)
	if err != nil {
		return "", &spookytypessecrets.SecretsError{
			Operation: "decrypt_variable",
			Cause:     err,
		}
	}

	decryptedStr := string(decrypted)
	m.logger.Info("Variable decrypted successfully")

	return decryptedStr, nil
}

// EncryptFile encrypts a file
func (m *Manager) EncryptFile(inputPath string, outputPath string, recipients []string) error {
	m.logger.Info("Encrypting file", spookylogging.String("input", inputPath), spookylogging.String("output", outputPath))

	// Use default recipients if none provided
	if len(recipients) == 0 {
		recipients = m.config.Keys.DefaultRecipients
	}

	return m.fileManager.EncryptFile(inputPath, outputPath, recipients)
}

// DecryptFile decrypts a file
func (m *Manager) DecryptFile(inputPath string, outputPath string, identity string) error {
	m.logger.Info("Decrypting file", spookylogging.String("input", inputPath), spookylogging.String("output", outputPath))

	// Use default identity if none provided
	if identity == "" {
		identity = m.config.Keys.DefaultIdentity
	}

	return m.fileManager.DecryptFile(inputPath, outputPath, identity)
}

// EncryptValue encrypts a single value
func (m *Manager) EncryptValue(value string, recipients []string) (string, error) {
	m.logger.Debug("Encrypting value")

	// Use default recipients if none provided
	if len(recipients) == 0 {
		recipients = m.config.Keys.DefaultRecipients
	}

	return m.fileManager.EncryptValue(value, recipients)
}

// DecryptValue decrypts a single value
func (m *Manager) DecryptValue(encryptedValue string, identity string) (string, error) {
	m.logger.Debug("Decrypting value")

	// Use default identity if none provided
	if identity == "" {
		identity = m.config.Keys.DefaultIdentity
	}

	return m.fileManager.DecryptValue(encryptedValue, identity)
}

// GetKeyManager returns the key manager
func (m *Manager) GetKeyManager() *spookysecretskeys.Manager {
	return m.keyManager
}

// GetConfig returns the secrets configuration
func (m *Manager) GetConfig() *spookytypessecrets.SecretsConfig {
	return m.config
}

// Validate validates the secrets configuration
func (m *Manager) Validate() error {
	// Validate algorithm
	if m.config.Encryption.Algorithm != "age" {
		return &spookytypessecrets.SecretsError{
			Operation: "validate_algorithm",
			Cause:     fmt.Errorf("unsupported algorithm: %s", m.config.Encryption.Algorithm),
		}
	}

	return nil
}

// TestEncryption tests encryption and decryption with the current configuration
func (m *Manager) TestEncryption() error {
	m.logger.Info("Testing encryption and decryption")

	// Generate test data
	testData := "Hello, Spooky Secrets!"

	// Get test identity and recipient
	testIdentity := m.config.Keys.DefaultIdentity
	if testIdentity == "" {
		return &spookytypessecrets.SecretsError{
			Operation: "test_encryption",
			Cause:     fmt.Errorf("no default identity configured for testing"),
		}
	}

	testRecipients := m.config.Keys.DefaultRecipients
	if len(testRecipients) == 0 {
		return &spookytypessecrets.SecretsError{
			Operation: "test_encryption",
			Cause:     fmt.Errorf("no default recipients configured for testing"),
		}
	}

	// Encrypt test data
	encrypted, err := m.ageClient.Encrypt([]byte(testData), testRecipients)
	if err != nil {
		return &spookytypessecrets.SecretsError{
			Operation: "test_encrypt",
			Cause:     err,
		}
	}

	// Decrypt the test data
	decrypted, err := m.ageClient.Decrypt(encrypted, testIdentity)
	if err != nil {
		return &spookytypessecrets.SecretsError{
			Operation: "test_decrypt",
			Cause:     err,
		}
	}

	// Verify decrypted data matches original
	if string(decrypted) != testData {
		return &spookytypessecrets.SecretsError{
			Operation: "test_verify",
			Cause:     fmt.Errorf("decrypted data does not match original"),
			Context: map[string]interface{}{
				"original":  testData,
				"decrypted": string(decrypted),
			},
		}
	}

	m.logger.Info("Encryption and decryption test passed successfully")
	return nil
}

// GetStatus returns the secrets system status
func (m *Manager) GetStatus() *spookytypessecrets.SecretsStatus {
	status := &spookytypessecrets.SecretsStatus{
		Enabled:      m.config.Enabled,
		Algorithm:    m.config.Encryption.Algorithm,
		AuditLogging: m.config.Security.AuditLogging,
	}

	// Get key information
	keys, err := m.keyManager.ListKeys()
	if err != nil {
		status.KeyError = err.Error()
	} else {
		status.KeyCount = len(keys)
		status.Keys = keys
	}

	// Check default identity
	status.DefaultIdentityConfigured = m.config.Keys.DefaultIdentity != ""

	// Check default recipients
	status.DefaultRecipientsConfigured = len(m.config.Keys.DefaultRecipients) > 0
	status.DefaultRecipientCount = len(m.config.Keys.DefaultRecipients)

	return status
}
