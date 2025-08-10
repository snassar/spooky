package interfaces

import (
	"io"

	spookytypessecrets "spooky/internal/types/secrets"
)

// AgeClient interface defines the operations for age encryption/decryption
type AgeClient interface {
	// Encrypt encrypts data with the given recipients
	Encrypt(data []byte, recipients []string) (*spookytypessecrets.EncryptedValue, error)

	// Decrypt decrypts data with the given identity
	Decrypt(encrypted *spookytypessecrets.EncryptedValue, identity string) ([]byte, error)

	// EncryptStream encrypts data from a reader to a writer with streaming support
	EncryptStream(input io.Reader, output io.Writer, recipients []string) error

	// DecryptStream decrypts data from a reader to a writer with streaming support
	DecryptStream(input io.Reader, output io.Writer, identity string) error
}

// KeyManager interface defines the operations for key management
type KeyManager interface {
	// ListKeys lists all available keys
	ListKeys() ([]*spookytypessecrets.KeyMetadata, error)

	// GetKey gets a key by name
	GetKey(name string) (*spookytypessecrets.KeyMetadata, error)
}

// SecretsManager interface defines the high-level secrets operations
type SecretsManager interface {
	// EncryptVariable encrypts a variable value
	EncryptVariable(name string, value string, recipients []string) (*spookytypessecrets.EncryptedValue, error)

	// DecryptVariable decrypts a variable value
	DecryptVariable(encrypted *spookytypessecrets.EncryptedValue, identity string) (string, error)

	// EncryptValue encrypts a value
	EncryptValue(value string, recipients []string) (string, error)

	// DecryptValue decrypts a value
	DecryptValue(encrypted string, identity string) (string, error)

	// EncryptFile encrypts a file
	EncryptFile(inputPath string, outputPath string, recipients []string) error

	// DecryptFile decrypts a file
	DecryptFile(inputPath string, outputPath string, identity string) error

	// GetConfig returns the secrets configuration
	GetConfig() *spookytypessecrets.SecretsConfig

	// Validate validates the secrets configuration
	Validate() error

	// TestEncryption tests encryption and decryption
	TestEncryption() error

	// GetStatus returns the secrets system status
	GetStatus() *spookytypessecrets.SecretsStatus
}
