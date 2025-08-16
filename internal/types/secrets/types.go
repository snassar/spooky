package secrets

import (
	"time"
)

// AgeConfig represents age encryption configuration
type AgeConfig struct {
	Identities string `hcl:"identities,optional" json:"identities,omitempty"`
	Recipients string `hcl:"recipients,optional" json:"recipients,omitempty"`
	Passphrase string `hcl:"passphrase,optional" json:"passphrase,omitempty"`

	Validation AgeValidationConfig `hcl:"validation,block" json:"validation,omitempty"`
	Encryption AgeEncryptionConfig `hcl:"encryption,block" json:"encryption,omitempty"`
}

// AgeValidationConfig represents age validation settings
type AgeValidationConfig struct {
	StrictMode      bool `hcl:"strict_mode,optional" json:"strict_mode,omitempty"`
	CheckRecipients bool `hcl:"check_recipients,optional" json:"check_recipients,omitempty"`
	ValidateKeys    bool `hcl:"validate_keys,optional" json:"validate_keys,omitempty"`
}

// AgeEncryptionConfig represents age encryption settings
type AgeEncryptionConfig struct {
	Algorithm   string `hcl:"algorithm,optional" json:"algorithm,omitempty"`
	Compression bool   `hcl:"compression,optional" json:"compression,omitempty"`
	Armor       bool   `hcl:"armor,optional" json:"armor,omitempty"`
}

// AgeIdentity represents an age identity
type AgeIdentity struct {
	Path      string    `json:"path"`
	PublicKey string    `json:"public_key"`
	Created   time.Time `json:"created"`
	Comment   string    `json:"comment,omitempty"`
}

// AgeRecipient represents an age recipient (public key)
type AgeRecipient struct {
	PublicKey string    `json:"public_key"`
	Added     time.Time `json:"added"`
	Comment   string    `json:"comment,omitempty"`
}

// AgeEncryptedData represents age-encrypted data
type AgeEncryptedData struct {
	Data       []byte    `json:"data"`
	Recipients []string  `json:"recipients"`
	Created    time.Time `json:"created"`
	Algorithm  string    `json:"algorithm"`
	Compressed bool      `json:"compressed"`
	Armored    bool      `json:"armored"`
}

// AgeKeyInfo represents information about an age key
type AgeKeyInfo struct {
	Type      string    `json:"type"` // "identity" or "recipient"
	Path      string    `json:"path,omitempty"`
	PublicKey string    `json:"public_key"`
	Created   time.Time `json:"created"`
	Comment   string    `json:"comment,omitempty"`
	Valid     bool      `json:"valid"`
}

// AgeValidationResult represents the result of age validation
type AgeValidationResult struct {
	Valid      bool     `json:"valid"`
	Errors     []string `json:"errors,omitempty"`
	Warnings   []string `json:"warnings,omitempty"`
	Recipients []string `json:"recipients,omitempty"`
}

// AgeEncryptionResult represents the result of age encryption
type AgeEncryptionResult struct {
	EncryptedData []byte    `json:"encrypted_data"`
	Recipients    []string  `json:"recipients"`
	Algorithm     string    `json:"algorithm"`
	Compressed    bool      `json:"compressed"`
	Armored       bool      `json:"armored"`
	Created       time.Time `json:"created"`
}

// AgeDecryptionResult represents the result of age decryption
type AgeDecryptionResult struct {
	DecryptedData []byte    `json:"decrypted_data"`
	Algorithm     string    `json:"algorithm"`
	Compressed    bool      `json:"compressed"`
	Armored       bool      `json:"armored"`
	Decrypted     time.Time `json:"decrypted"`
}
