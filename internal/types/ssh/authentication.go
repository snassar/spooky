// Package ssh provides SSH authentication types for the spooky codebase.
// This package defines the data structures for SSH authentication and security.
package ssh

import (
	"time"

	spookytypescommon "spooky/internal/types/common"
)

// Authentication represents SSH authentication configuration
type Authentication struct {
	spookytypescommon.CompleteEntity

	// Authentication method
	Method AuthMethod `json:"method" hcl:"method" default:"public_key"`

	// Password authentication
	Password     string `json:"password,omitempty" hcl:"password,optional" sensitive:"true"`
	PasswordHash string `json:"password_hash,omitempty" hcl:"password_hash,optional" sensitive:"true"`

	// Public key authentication
	KeyPath    string  `json:"key_path,omitempty" hcl:"key_path,optional"`
	KeyData    []byte  `json:"key_data,omitempty" hcl:"key_data,optional" sensitive:"true"`
	Passphrase string  `json:"passphrase,omitempty" hcl:"passphrase,optional" sensitive:"true"`
	KeyType    KeyType `json:"key_type,omitempty" hcl:"key_type,optional"`

	// Key metadata
	KeyFingerprint string     `json:"key_fingerprint,omitempty" hcl:"key_fingerprint,optional"`
	KeyComment     string     `json:"key_comment,omitempty" hcl:"key_comment,optional"`
	KeyCreatedAt   *time.Time `json:"key_created_at,omitempty" hcl:"key_created_at,optional"`

	// Authentication settings
	AllowPasswordAuth  bool `json:"allow_password_auth" hcl:"allow_password_auth" default:"false"`
	AllowPublicKeyAuth bool `json:"allow_public_key_auth" hcl:"allow_public_key_auth" default:"true"`
	AllowAgentAuth     bool `json:"allow_agent_auth" hcl:"allow_agent_auth" default:"false"`

	// Security settings
	MaxAuthAttempts int           `json:"max_auth_attempts" hcl:"max_auth_attempts" default:"3"`
	AuthTimeout     time.Duration `json:"auth_timeout" hcl:"auth_timeout" default:"30s"`
	LockoutDuration time.Duration `json:"lockout_duration" hcl:"lockout_duration" default:"300s"`

	// Authentication metadata
	LastAuthAttempt *time.Time `json:"last_auth_attempt,omitempty" hcl:"last_auth_attempt,optional"`
	AuthAttempts    int        `json:"auth_attempts" hcl:"auth_attempts"`
	SuccessfulAuths int        `json:"successful_auths" hcl:"successful_auths"`
	FailedAuths     int        `json:"failed_auths" hcl:"failed_auths"`
}

// KeyType represents the type of SSH key
type KeyType string

const (
	KeyTypeRSA     KeyType = "rsa"
	KeyTypeDSA     KeyType = "dsa"
	KeyTypeECDSA   KeyType = "ecdsa"
	KeyTypeED25519 KeyType = "ed25519"
)

// Key represents an SSH key pair
type Key struct {
	spookytypescommon.CompleteEntity

	// Key details
	Type  KeyType `json:"type" hcl:"type"`
	Bits  int     `json:"bits,omitempty" hcl:"bits,optional"`
	Curve string  `json:"curve,omitempty" hcl:"curve,optional"`

	// Key data
	PrivateKey []byte `json:"private_key,omitempty" hcl:"private_key,optional" sensitive:"true"`
	PublicKey  []byte `json:"public_key,omitempty" hcl:"public_key,optional"`
	Passphrase string `json:"passphrase,omitempty" hcl:"passphrase,optional" sensitive:"true"`

	// Key metadata
	Fingerprint string     `json:"fingerprint" hcl:"fingerprint"`
	Comment     string     `json:"comment,omitempty" hcl:"comment,optional"`
	CreatedAt   time.Time  `json:"created_at" hcl:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" hcl:"expires_at,optional"`

	// Key file paths
	PrivateKeyPath string `json:"private_key_path,omitempty" hcl:"private_key_path,optional"`
	PublicKeyPath  string `json:"public_key_path,omitempty" hcl:"public_key_path,optional"`

	// Key permissions
	PrivateKeyPermissions string `json:"private_key_permissions,omitempty" hcl:"private_key_permissions,optional"`
	PublicKeyPermissions  string `json:"public_key_permissions,omitempty" hcl:"public_key_permissions,optional"`

	// Key validation
	IsValid         bool   `json:"is_valid" hcl:"is_valid"`
	ValidationError string `json:"validation_error,omitempty" hcl:"validation_error,optional"`

	// Key usage
	UsageCount int        `json:"usage_count" hcl:"usage_count"`
	LastUsed   *time.Time `json:"last_used,omitempty" hcl:"last_used,optional"`
	IsActive   bool       `json:"is_active" hcl:"is_active" default:"true"`
}

// KeyPair represents a complete SSH key pair
type KeyPair struct {
	spookytypescommon.CompleteEntity

	// Key pair details
	Name        string `json:"name" hcl:"name"`
	Description string `json:"description,omitempty" hcl:"description,optional"`

	// Keys
	PrivateKey *Key `json:"private_key" hcl:"private_key"`
	PublicKey  *Key `json:"public_key" hcl:"public_key"`

	// Key pair metadata
	Algorithm string     `json:"algorithm" hcl:"algorithm"`
	KeySize   int        `json:"key_size" hcl:"key_size"`
	CreatedAt time.Time  `json:"created_at" hcl:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty" hcl:"expires_at,optional"`

	// Key pair validation
	IsValid         bool   `json:"is_valid" hcl:"is_valid"`
	ValidationError string `json:"validation_error,omitempty" hcl:"validation_error,optional"`

	// Key pair usage
	UsageCount int        `json:"usage_count" hcl:"usage_count"`
	LastUsed   *time.Time `json:"last_used,omitempty" hcl:"last_used,optional"`
	IsActive   bool       `json:"is_active" hcl:"is_active" default:"true"`

	// Key pair security
	RequiresPassphrase bool `json:"requires_passphrase" hcl:"requires_passphrase" default:"false"`
	IsEncrypted        bool `json:"is_encrypted" hcl:"is_encrypted" default:"false"`
}

// HostKey represents a host key for SSH host verification
type HostKey struct {
	spookytypescommon.CompleteEntity

	// Host key details
	Hostname    string  `json:"hostname" hcl:"hostname"`
	Port        int     `json:"port" hcl:"port" default:"22"`
	KeyType     KeyType `json:"key_type" hcl:"key_type"`
	Fingerprint string  `json:"fingerprint" hcl:"fingerprint"`

	// Host key data
	PublicKey []byte `json:"public_key,omitempty" hcl:"public_key,optional"`

	// Host key metadata
	Algorithm string    `json:"algorithm" hcl:"algorithm"`
	KeySize   int       `json:"key_size" hcl:"key_size"`
	FirstSeen time.Time `json:"first_seen" hcl:"first_seen"`
	LastSeen  time.Time `json:"last_seen" hcl:"last_seen"`

	// Host key validation
	IsValid         bool   `json:"is_valid" hcl:"is_valid"`
	ValidationError string `json:"validation_error,omitempty" hcl:"validation_error,optional"`

	// Host key trust
	IsTrusted   bool       `json:"is_trusted" hcl:"is_trusted" default:"false"`
	TrustLevel  TrustLevel `json:"trust_level" hcl:"trust_level" default:"unknown"`
	ManualTrust bool       `json:"manual_trust" hcl:"manual_trust" default:"false"`

	// Host key usage
	UsageCount int        `json:"usage_count" hcl:"usage_count"`
	LastUsed   *time.Time `json:"last_used,omitempty" hcl:"last_used,optional"`
}

// TrustLevel represents the trust level of a host key
type TrustLevel string

const (
	TrustLevelUnknown   TrustLevel = "unknown"
	TrustLevelUntrusted TrustLevel = "untrusted"
	TrustLevelTrusted   TrustLevel = "trusted"
	TrustLevelVerified  TrustLevel = "verified"
)

// KnownHosts represents a known_hosts file
type KnownHosts struct {
	spookytypescommon.CompleteEntity

	// Known hosts file details
	FilePath string `json:"file_path" hcl:"file_path"`
	IsValid  bool   `json:"is_valid" hcl:"is_valid"`

	// Known hosts entries
	Entries []*HostKey `json:"entries" hcl:"entries"`

	// Known hosts metadata
	TotalEntries   int       `json:"total_entries" hcl:"total_entries"`
	TrustedEntries int       `json:"trusted_entries" hcl:"trusted_entries"`
	LastModified   time.Time `json:"last_modified" hcl:"last_modified"`

	// Known hosts settings
	StrictHostKeyChecking bool `json:"strict_host_key_checking" hcl:"strict_host_key_checking" default:"true"`
	AllowInsecureHosts    bool `json:"allow_insecure_hosts" hcl:"allow_insecure_hosts" default:"false"`

	// Known hosts validation
	ValidationError string `json:"validation_error,omitempty" hcl:"validation_error,optional"`
}

// AuthenticationResult represents the result of an authentication attempt
type AuthenticationResult struct {
	spookytypescommon.CompleteEntity

	// Authentication details
	Authentication *Authentication `json:"authentication" hcl:"authentication"`
	Hostname       string          `json:"hostname" hcl:"hostname"`
	Port           int             `json:"port" hcl:"port"`
	Username       string          `json:"username" hcl:"username"`

	// Authentication results
	Success bool       `json:"success" hcl:"success"`
	Error   string     `json:"error,omitempty" hcl:"error,optional"`
	Method  AuthMethod `json:"method" hcl:"method"`

	// Authentication metrics
	StartTime time.Time     `json:"start_time" hcl:"start_time"`
	EndTime   time.Time     `json:"end_time" hcl:"end_time"`
	Duration  time.Duration `json:"duration" hcl:"duration"`
	Attempts  int           `json:"attempts" hcl:"attempts"`

	// Authentication metadata
	ClientVersion  string `json:"client_version,omitempty" hcl:"client_version,optional"`
	ServerVersion  string `json:"server_version,omitempty" hcl:"server_version,optional"`
	KeyFingerprint string `json:"key_fingerprint,omitempty" hcl:"key_fingerprint,optional"`

	// Security information
	HostKeyFingerprint string `json:"host_key_fingerprint,omitempty" hcl:"host_key_fingerprint,optional"`
	AuditTrail         string `json:"audit_trail,omitempty" hcl:"audit_trail,optional"`
}
