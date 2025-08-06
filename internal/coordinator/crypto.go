package coordinator

import (
	"fmt"
	"os"
	"strings"
	"time"

	spookylogging "spooky/internal/logging"
	spookysecrets "spooky/internal/secrets"
)

// CoordinatorCryptoIntegration implements crypto system integration
type CoordinatorCryptoIntegration struct {
	cryptoManager *spookysecrets.Manager
	logger        spookylogging.Logger
}

// NewCoordinatorCryptoIntegration creates a new crypto integration
func NewCoordinatorCryptoIntegration(cryptoManager *spookysecrets.Manager, logger spookylogging.Logger) *CoordinatorCryptoIntegration {
	return &CoordinatorCryptoIntegration{
		cryptoManager: cryptoManager,
		logger:        logger,
	}
}

// EncryptData encrypts data with age encryption
func (ci *CoordinatorCryptoIntegration) EncryptData(data []byte, recipients []string) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	if len(recipients) == 0 {
		return nil, fmt.Errorf("at least one recipient must be specified")
	}

	if ci.cryptoManager == nil {
		return nil, fmt.Errorf("crypto manager not available")
	}

	// Convert data to string for encryption
	dataStr := string(data)

	// Encrypt the data using the crypto manager
	encrypted, err := ci.cryptoManager.EncryptValue(dataStr, recipients)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	ci.logger.Info("Encrypted data", spookylogging.Int("data_size", len(data)), spookylogging.Int("recipients", len(recipients)))

	return []byte(encrypted), nil
}

// DecryptData decrypts data with age encryption
func (ci *CoordinatorCryptoIntegration) DecryptData(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	if ci.cryptoManager == nil {
		return nil, fmt.Errorf("crypto manager not available")
	}

	// Convert data to string for decryption
	dataStr := string(data)

	// Decrypt the data using the crypto manager
	decrypted, err := ci.cryptoManager.DecryptValue(dataStr, "")
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	ci.logger.Info("Decrypted data", spookylogging.Int("data_size", len(data)))

	return []byte(decrypted), nil
}

// ValidateEncryption validates encrypted data with enhanced checks
func (ci *CoordinatorCryptoIntegration) ValidateEncryption(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("data cannot be empty")
	}

	if ci.cryptoManager == nil {
		return fmt.Errorf("crypto manager not available")
	}

	// Check if data appears to be encrypted
	if len(data) < 100 { // Minimum size for age encryption
		return fmt.Errorf("data too small to be valid encrypted data")
	}

	// Check for age encryption header
	dataStr := string(data)
	if !strings.Contains(dataStr, "-----BEGIN AGE ENCRYPTED FILE-----") {
		return fmt.Errorf("data does not appear to be age encrypted (missing header)")
	}

	// Check for age encryption footer
	if !strings.Contains(dataStr, "-----END AGE ENCRYPTED FILE-----") {
		return fmt.Errorf("data does not appear to be age encrypted (missing footer)")
	}

	// Try to decrypt to validate
	_, err := ci.DecryptData(data)
	return err
}

// GetCryptoStatus returns crypto system status
func (ci *CoordinatorCryptoIntegration) GetCryptoStatus() map[string]interface{} {
	status := map[string]interface{}{
		"encryption_enabled": ci.cryptoManager != nil,
	}

	if ci.cryptoManager != nil {
		// Get actual crypto manager status
		cryptoStatus := ci.cryptoManager.GetStatus()

		status["algorithm"] = cryptoStatus.Algorithm
		status["key_count"] = cryptoStatus.KeyCount
		status["default_identity"] = ci.cryptoManager.GetConfig().Keys.DefaultIdentity
		status["default_recipients"] = ci.cryptoManager.GetConfig().Keys.DefaultRecipients
		status["audit_logging"] = cryptoStatus.AuditLogging
		status["default_identity_configured"] = cryptoStatus.DefaultIdentityConfigured
		status["default_recipients_configured"] = cryptoStatus.DefaultRecipientsConfigured
		status["default_recipient_count"] = cryptoStatus.DefaultRecipientCount

		if cryptoStatus.KeyError != "" {
			status["key_error"] = cryptoStatus.KeyError
		}
	}

	return status
}

// EncryptFile encrypts a file with atomic operations
func (ci *CoordinatorCryptoIntegration) EncryptFile(filePath string, recipients []string) error {
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	if len(recipients) == 0 {
		return fmt.Errorf("at least one recipient must be specified")
	}

	if ci.cryptoManager == nil {
		return fmt.Errorf("crypto manager not available")
	}

	// Create output path for encrypted file
	outputPath := filePath + ".age"

	// Encrypt the file using the crypto manager
	err := ci.cryptoManager.EncryptFile(filePath, outputPath, recipients)
	if err != nil {
		return fmt.Errorf("failed to encrypt file: %w", err)
	}

	// Verify the encrypted file
	if err := ci.ValidateEncryptedFile(outputPath); err != nil {
		return fmt.Errorf("encrypted file validation failed: %w", err)
	}

	ci.logger.Info("Encrypted file", spookylogging.String("file", filePath), spookylogging.String("output", outputPath), spookylogging.Int("recipients", len(recipients)))

	return nil
}

// DecryptFile decrypts a file with atomic operations
func (ci *CoordinatorCryptoIntegration) DecryptFile(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	if ci.cryptoManager == nil {
		return fmt.Errorf("crypto manager not available")
	}

	// Create output path for decrypted file (remove .age extension if present)
	outputPath := filePath
	if len(filePath) > 4 && filePath[len(filePath)-4:] == ".age" {
		outputPath = filePath[:len(filePath)-4]
	}

	// Decrypt the file using the crypto manager
	err := ci.cryptoManager.DecryptFile(filePath, outputPath, "")
	if err != nil {
		return fmt.Errorf("failed to decrypt file: %w", err)
	}

	// Verify the decrypted file
	if err := ci.validateDecryptedFile(outputPath); err != nil {
		return fmt.Errorf("decrypted file validation failed: %w", err)
	}

	ci.logger.Info("Decrypted file", spookylogging.String("file", filePath), spookylogging.String("output", outputPath))

	return nil
}

// ValidateEncryptedFile validates an encrypted file with comprehensive checks
func (ci *CoordinatorCryptoIntegration) ValidateEncryptedFile(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	if ci.cryptoManager == nil {
		return fmt.Errorf("crypto manager not available")
	}

	// Check if file exists
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("encrypted file does not exist: %w", err)
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read encrypted file: %w", err)
	}

	// Validate encryption format
	if err := ci.ValidateEncryption(content); err != nil {
		return fmt.Errorf("encryption validation failed: %w", err)
	}

	ci.logger.Info("Validated encrypted file", spookylogging.String("file", filePath))

	return nil
}

// validateDecryptedFile validates a decrypted file
func (ci *CoordinatorCryptoIntegration) validateDecryptedFile(filePath string) error {
	// Check if file exists
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("decrypted file does not exist: %w", err)
	}

	// Check if file is readable
	if _, err := os.ReadFile(filePath); err != nil {
		return fmt.Errorf("decrypted file is not readable: %w", err)
	}

	ci.logger.Debug("Validated decrypted file", spookylogging.String("file", filePath))

	return nil
}

// AuditLog logs crypto operations for audit purposes
func (ci *CoordinatorCryptoIntegration) AuditLog(operation string, details map[string]interface{}) {
	if ci.cryptoManager == nil {
		return
	}

	// Add timestamp and operation type
	auditEntry := map[string]interface{}{
		"timestamp": time.Now(),
		"operation": operation,
	}

	// Add details
	for key, value := range details {
		auditEntry[key] = value
	}

	ci.logger.Info("Crypto audit log",
		spookylogging.String("operation", operation),
		spookylogging.String("details", fmt.Sprintf("%v", details)))

	// In a real implementation, this would write to an audit log file
	// or send to an audit system
}
