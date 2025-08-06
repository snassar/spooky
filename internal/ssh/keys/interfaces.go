package keys

import (
	"spooky/internal/ssh/types"
)

// SSHKeyManager defines the interface for SSH key management
type SSHKeyManager interface {
	LoadPrivateKey(path string) (*types.SSHKey, error)
	ValidateKeyFile(path string) error
}
