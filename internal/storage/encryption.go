package storage

import (
	"time"

	spookyfactstypes "spooky/internal/facts/types"
)

// EncryptionHelper provides encryption utilities for storage
type EncryptionHelper struct {
	cryptoManager interface{} // Will be properly typed when secrets package is available
}

// NewEncryptionHelper creates a new encryption helper
func NewEncryptionHelper(cryptoManager interface{}) *EncryptionHelper {
	return &EncryptionHelper{
		cryptoManager: cryptoManager,
	}
}

// AddEncryptionMetadata adds encryption metadata to a fact collection
func (eh *EncryptionHelper) AddEncryptionMetadata(collection *spookyfactstypes.FactCollection) {
	collection.EncryptionMetadata = &spookyfactstypes.EncryptionMetadata{
		EncryptedAt:       time.Now().Format(time.RFC3339),
		EncryptionVersion: "1.0",
		Recipients:        eh.getDefaultRecipients(),
	}
}

// IsEncryptedCollection checks if a fact collection is encrypted
func (eh *EncryptionHelper) IsEncryptedCollection(collection *spookyfactstypes.FactCollection) bool {
	return collection.EncryptionMetadata != nil && collection.EncryptedData != ""
}

// getDefaultRecipients returns the default encryption recipients
func (eh *EncryptionHelper) getDefaultRecipients() []string {
	// In a real implementation, this would read from configuration
	return []string{"age1example"}
}

// getDefaultIdentity returns the default encryption identity
func (eh *EncryptionHelper) getDefaultIdentity() string {
	// In a real implementation, this would read from configuration
	return "~/.config/spooky/keys/age.key"
}
