package secrets

import (
	"fmt"
	"time"
)

// EncryptedValue represents an encrypted value with metadata
type EncryptedValue struct {
	Data       []byte            `json:"data"`
	Recipients []string          `json:"recipients"`
	Created    time.Time         `json:"created"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// SecretsError represents a secrets operation error
type SecretsError struct {
	Operation string                 `json:"operation"`
	Cause     error                  `json:"cause"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

// Error returns the error message
func (e *SecretsError) Error() string {
	return fmt.Sprintf("secrets %s failed: %v", e.Operation, e.Cause)
}

// Unwrap returns the underlying error
func (e *SecretsError) Unwrap() error {
	return e.Cause
}

// SecretsConfig represents secrets configuration
type SecretsConfig struct {
	Enabled bool `json:"enabled"`

	Keys struct {
		DefaultIdentity   string   `json:"default_identity"`
		DefaultRecipients []string `json:"default_recipients"`
	} `json:"keys"`

	Encryption struct {
		Algorithm string `json:"algorithm"`
	} `json:"encryption"`

	Security struct {
		AuditLogging  bool `json:"audit_logging"`
		KeyValidation bool `json:"key_validation"`
		MemoryWipe    bool `json:"memory_wipe"`
	} `json:"security"`
}

// KeyMetadata represents key metadata
type KeyMetadata struct {
	Name        string    `json:"name"`
	Fingerprint string    `json:"fingerprint"`
	Type        string    `json:"type"`
	Created     time.Time `json:"created"`
	Description string    `json:"description,omitempty"`
}

// SecretsStatus represents secrets system status
type SecretsStatus struct {
	Enabled                     bool           `json:"enabled"`
	Algorithm                   string         `json:"algorithm"`
	AuditLogging                bool           `json:"audit_logging"`
	KeyCount                    int            `json:"key_count"`
	Keys                        []*KeyMetadata `json:"keys,omitempty"`
	KeyError                    string         `json:"key_error,omitempty"`
	DefaultIdentityConfigured   bool           `json:"default_identity_configured"`
	DefaultRecipientsConfigured bool           `json:"default_recipients_configured"`
	DefaultRecipientCount       int            `json:"default_recipient_count"`
}
