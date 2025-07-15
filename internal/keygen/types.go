package keygen

// KeyPair represents a generated SSH key pair
type KeyPair struct {
	PrivateKey []byte
	PublicKey  []byte
	Password   string
}

// KeyFiles represents the files generated for a key pair
type KeyFiles struct {
	PrivateKeyPath string
	PublicKeyPath  string
	PasswordPath   string
	CombinedPath   string
	OutputDir      string
}

// KeyType represents the type of SSH key to generate
type KeyType string

const (
	// Ed25519KeyType represents Ed25519 key generation
	Ed25519KeyType KeyType = "ed25519"

	// DefaultPasswordLength is the default length for generated passwords
	DefaultPasswordLength = 25

	// MaxKeyDirectories is the maximum number of key directories per day
	MaxKeyDirectories = 1000

	// DefaultKeyType is the default key type to generate
	DefaultKeyType = Ed25519KeyType
)
