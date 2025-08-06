package interfaces

// CryptoIntegration defines the interface for crypto system integration
type CryptoIntegration interface {
	// EncryptData encrypts data with age encryption
	EncryptData(data []byte, recipients []string) ([]byte, error)

	// DecryptData decrypts data with age encryption
	DecryptData(data []byte) ([]byte, error)

	// ValidateEncryption validates encrypted data
	ValidateEncryption(data []byte) error

	// GetCryptoStatus returns crypto system status
	GetCryptoStatus() map[string]interface{}

	// EncryptFile encrypts a file
	EncryptFile(filePath string, recipients []string) error

	// DecryptFile decrypts a file
	DecryptFile(filePath string) error

	// ValidateEncryptedFile validates an encrypted file
	ValidateEncryptedFile(filePath string) error
}
