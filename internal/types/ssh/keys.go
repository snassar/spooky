package ssh

import (
	"time"
)

// SSHKey represents an SSH key
type SSHKey struct {
	Type       string    `hcl:"type"`
	PrivateKey []byte    `hcl:"private_key"`
	PublicKey  []byte    `hcl:"public_key,optional"`
	Path       string    `hcl:"path"`
	CreatedAt  time.Time `hcl:"created_at"`
}

// KeysConfig represents SSH keys configuration
type KeysConfig struct {
	DefaultKeyPath string        `hcl:"default_key_path,optional"`
	KeyCacheTTL    time.Duration `hcl:"key_cache_ttl,optional"`
	EnableCaching  bool          `hcl:"enable_caching,optional"`
}

// KeyValidationResult represents the result of key validation
type KeyValidationResult struct {
	Valid     bool      `hcl:"valid"`
	KeyPath   string    `hcl:"key_path"`
	Error     string    `hcl:"error,optional"`
	Timestamp time.Time `hcl:"timestamp"`
}
