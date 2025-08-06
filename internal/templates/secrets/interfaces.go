package secrets

// SecretsManager defines the interface for template secrets integration
type SecretsManager interface {
	// Core secrets operations
	EncryptValue(value string) (string, error)
	DecryptValue(encryptedValue string) (string, error)
	ValidateEncryptionKey(key string) error

	// Configuration
	SetEncryptionKey(key string) error
	SetEncryptionAlgorithm(algorithm string) error
	EnableEncryption(enabled bool) error

	// Utility operations
	IsEncrypted(value string) bool
	GetEncryptionAlgorithm() string
	Close() error
}

// EncryptionProvider defines the interface for encryption providers
type EncryptionProvider interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
	GetAlgorithm() string
}
