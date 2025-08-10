package machines

import (
	"time"

	spookytypessecrets "spooky/internal/types/secrets"
)

// MachineSecrets represents machine secrets configuration
type MachineSecrets struct {
	Secrets    map[string]spookytypessecrets.EncryptedValue `json:"secrets"`
	PrivateKey spookytypessecrets.EncryptedValue            `json:"private_key,omitempty"`
}

// MachineKeys represents encryption keys for a machine
type MachineKeys struct {
	MachineName string                            `json:"machine_name"`
	PublicKey   string                            `json:"public_key"`
	PrivateKey  spookytypessecrets.EncryptedValue `json:"private_key,omitempty"`
	Created     time.Time                         `json:"created"`
	Expires     *time.Time                        `json:"expires,omitempty"`
	Metadata    map[string]string                 `json:"metadata,omitempty"`
}

// SecretsConfig represents machine-specific secrets configuration
type SecretsConfig struct {
	Enabled bool `json:"enabled"`

	// Machine-specific encryption settings
	MachineEncryption struct {
		Algorithm string `json:"algorithm"`
	} `json:"machine_encryption"`

	// Key management
	Keys struct {
		DefaultIdentity   string              `json:"default_identity"`
		DefaultRecipients []string            `json:"default_recipients"`
		MachineRecipients map[string][]string `json:"machine_recipients,omitempty"`
	} `json:"keys"`

	// Security settings
	Security struct {
		AuditLogging     bool `json:"audit_logging"`
		KeyValidation    bool `json:"key_validation"`
		MemoryWipe       bool `json:"memory_wipe"`
		ValidateOnAccess bool `json:"validate_on_access"`
	} `json:"security"`
}

// SecretsStatus represents the status of machine secrets
type SecretsStatus struct {
	Enabled                     bool           `json:"enabled"`
	Algorithm                   string         `json:"algorithm"`
	AuditLogging                bool           `json:"audit_logging"`
	MachineCount                int            `json:"machine_count"`
	EncryptedMachineCount       int            `json:"encrypted_machine_count"`
	KeyCount                    int            `json:"key_count"`
	Keys                        []*KeyMetadata `json:"keys,omitempty"`
	KeyError                    string         `json:"key_error,omitempty"`
	DefaultIdentityConfigured   bool           `json:"default_identity_configured"`
	DefaultRecipientsConfigured bool           `json:"default_recipients_configured"`
	DefaultRecipientCount       int            `json:"default_recipient_count"`
	MachineRecipientsConfigured bool           `json:"machine_recipients_configured"`
}

// KeyMetadata represents machine key metadata
type KeyMetadata struct {
	MachineName string     `json:"machine_name"`
	Name        string     `json:"name"`
	Fingerprint string     `json:"fingerprint"`
	Type        string     `json:"type"`
	Created     time.Time  `json:"created"`
	Expires     *time.Time `json:"expires,omitempty"`
	Description string     `json:"description,omitempty"`
	IsDefault   bool       `json:"is_default"`
	IsEncrypted bool       `json:"is_encrypted"`
}

// MachineSecretsError represents a machine secrets operation error
type MachineSecretsError struct {
	MachineName string                 `json:"machine_name"`
	Operation   string                 `json:"operation"`
	Cause       error                  `json:"cause"`
	Context     map[string]interface{} `json:"context,omitempty"`
}

// Error returns the error message
func (e *MachineSecretsError) Error() string {
	return e.Cause.Error()
}

// Unwrap returns the underlying error
func (e *MachineSecretsError) Unwrap() error {
	return e.Cause
}
