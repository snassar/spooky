package interfaces

// SecretsIntegration defines the interface for secrets system integration
type SecretsIntegration interface {
	// EncryptData encrypts data with age encryption
	EncryptData(data []byte, recipients []string) ([]byte, error)

	// DecryptData decrypts data with age encryption
	DecryptData(data []byte) ([]byte, error)

	// ValidateEncryption validates encrypted data
	ValidateEncryption(data []byte) error

	// EncryptFile encrypts a file
	EncryptFile(filePath string, recipients []string) error

	// DecryptFile decrypts a file
	DecryptFile(filePath string) error

	// ValidateEncryptedFile validates an encrypted file
	ValidateEncryptedFile(filePath string) error

	// GetDefaultRecipients returns the default recipients for encryption
	GetDefaultRecipients() ([]string, error)
}
