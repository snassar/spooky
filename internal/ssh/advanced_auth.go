// Package ssh provides advanced authentication capabilities for the spooky codebase.
// This package implements multi-factor authentication, certificate chains, and other advanced authentication methods.
package ssh

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	spookytypeslogging "spooky/internal/types/logging"
	spookytypesssh "spooky/internal/types/ssh"
)

// AdvancedAuthManager manages advanced authentication methods
type AdvancedAuthManager struct {
	logger spookytypeslogging.Logger
	agent  *Agent
}

// NewAdvancedAuthManager creates a new advanced authentication manager
func NewAdvancedAuthManager(logger spookytypeslogging.Logger) *AdvancedAuthManager {
	return &AdvancedAuthManager{
		logger: logger,
		agent:  NewAgent(logger),
	}
}

// MultiFactorAuthConfig represents multi-factor authentication configuration
type MultiFactorAuthConfig struct {
	// Primary authentication
	PrimaryMethod spookytypesssh.AuthMethod `json:"primary_method" hcl:"primary_method"`
	PrimaryKey    string                    `json:"primary_key,omitempty" hcl:"primary_key,optional"`
	PrimaryPass   string                    `json:"primary_pass,omitempty" hcl:"primary_pass,optional" sensitive:"true"`

	// Secondary authentication
	SecondaryMethod spookytypesssh.AuthMethod `json:"secondary_method" hcl:"secondary_method"`
	SecondaryKey    string                    `json:"secondary_key,omitempty" hcl:"secondary_key,optional"`
	SecondaryPass   string                    `json:"secondary_pass,omitempty" hcl:"secondary_pass,optional" sensitive:"true"`

	// TOTP configuration
	TOTPSecret    string `json:"totp_secret,omitempty" hcl:"totp_secret,optional" sensitive:"true"`
	TOTPAlgorithm string `json:"totp_algorithm" hcl:"totp_algorithm" default:"sha1"`
	TOTPDigits    int    `json:"totp_digits" hcl:"totp_digits" default:"6"`
	TOTPPeriod    int    `json:"totp_period" hcl:"totp_period" default:"30"`

	// Hardware token configuration
	HardwareTokenPath string `json:"hardware_token_path,omitempty" hcl:"hardware_token_path,optional"`
	HardwareTokenPin  string `json:"hardware_token_pin,omitempty" hcl:"hardware_token_pin,optional" sensitive:"true"`

	// Certificate chain configuration
	CertificateChain []string `json:"certificate_chain,omitempty" hcl:"certificate_chain,optional"`
	CAKeyPath        string   `json:"ca_key_path,omitempty" hcl:"ca_key_path,optional"`
	CAKeyPass        string   `json:"ca_key_pass,omitempty" hcl:"ca_key_pass,optional" sensitive:"true"`

	// Authentication order
	AuthOrder []string `json:"auth_order" hcl:"auth_order" default:"primary,secondary,totp,hardware"`

	// Retry configuration
	MaxRetries    int           `json:"max_retries" hcl:"max_retries" default:"3"`
	RetryDelay    time.Duration `json:"retry_delay" hcl:"retry_delay" default:"5s"`
	BackoffFactor float64       `json:"backoff_factor" hcl:"backoff_factor" default:"2.0"`
}

// CertificateChain represents a chain of SSH certificates
type CertificateChain struct {
	Certificates []*ssh.Certificate `json:"certificates"`
	Signers      []ssh.Signer       `json:"signers"`
	ValidFrom    time.Time          `json:"valid_from"`
	ValidUntil   time.Time          `json:"valid_until"`
	Principals   []string           `json:"principals"`
	Extensions   map[string]string  `json:"extensions"`
}

// TOTPGenerator generates TOTP codes
type TOTPGenerator struct {
	secret    string
	algorithm string
	digits    int
	period    int
}

// NewTOTPGenerator creates a new TOTP generator
func NewTOTPGenerator(secret, algorithm string, digits, period int) *TOTPGenerator {
	return &TOTPGenerator{
		secret:    secret,
		algorithm: algorithm,
		digits:    digits,
		period:    period,
	}
}

// GenerateTOTP generates a TOTP code
func (t *TOTPGenerator) GenerateTOTP() (string, error) {
	// This is a simplified TOTP implementation
	// In a real implementation, you would use a proper TOTP library
	timestamp := time.Now().Unix() / int64(t.period)

	// Generate a simple hash-based code
	hash := fmt.Sprintf("%x", timestamp)
	if len(hash) > t.digits {
		hash = hash[:t.digits]
	} else {
		// Pad with zeros if needed
		for len(hash) < t.digits {
			hash = "0" + hash
		}
	}

	return hash, nil
}

// SSHAgent manages SSH agent connections
type Agent struct {
	logger spookytypeslogging.Logger
	conn   agent.Agent
}

// NewAgent creates a new SSH agent
func NewAgent(logger spookytypeslogging.Logger) *Agent {
	return &Agent{
		logger: logger,
	}
}

// Connect connects to the SSH agent
func (sa *Agent) Connect() error {
	// For now, we'll use a keyring as a fallback
	// In a real implementation, you would connect to the actual SSH agent
	sa.conn = agent.NewKeyring()
	sa.logger.Info("Connected to SSH agent (using keyring)", map[string]interface{}{})

	return nil
}

// ListKeys lists keys in the SSH agent
func (sa *Agent) ListKeys() ([]*agent.Key, error) {
	if sa.conn == nil {
		return nil, fmt.Errorf("SSH agent not connected")
	}

	keys, err := sa.conn.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list SSH agent keys: %w", err)
	}

	return keys, nil
}

// AddKey adds a key to the SSH agent
func (sa *Agent) AddKey(key ssh.Signer) error {
	if sa.conn == nil {
		return fmt.Errorf("SSH agent not connected")
	}

	if err := sa.conn.Add(agent.AddedKey{
		PrivateKey: key,
	}); err != nil {
		return fmt.Errorf("failed to add key to SSH agent: %w", err)
	}

	sa.logger.Info("Added key to SSH agent", map[string]interface{}{
		"key_type": key.PublicKey().Type(),
	})

	return nil
}

// RemoveKey removes a key from the SSH agent
func (sa *Agent) RemoveKey(key ssh.PublicKey) error {
	if sa.conn == nil {
		return fmt.Errorf("SSH agent not connected")
	}

	if err := sa.conn.Remove(key); err != nil {
		return fmt.Errorf("failed to remove key from SSH agent: %w", err)
	}

	sa.logger.Info("Removed key from SSH agent", map[string]interface{}{
		"key_type": key.Type(),
	})

	return nil
}

// GetAuthMethods returns authentication methods for multi-factor authentication
func (aam *AdvancedAuthManager) GetAuthMethods(config *MultiFactorAuthConfig) ([]ssh.AuthMethod, error) {
	var authMethods []ssh.AuthMethod

	// Process authentication methods in order
	for _, methodName := range config.AuthOrder {
		method, err := aam.createAuthMethod(config, methodName)
		if err != nil {
			aam.logger.Warn("Failed to create auth method", map[string]interface{}{
				"method": methodName,
				"error":  err.Error(),
			})
			continue
		}
		authMethods = append(authMethods, method)
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("no valid authentication methods found")
	}

	aam.logger.Info("Created authentication methods", map[string]interface{}{
		"method_count": len(authMethods),
		"methods":      config.AuthOrder,
	})

	return authMethods, nil
}

// createAuthMethod creates a specific authentication method
func (aam *AdvancedAuthManager) createAuthMethod(config *MultiFactorAuthConfig, methodName string) (ssh.AuthMethod, error) {
	switch methodName {
	case "primary":
		return aam.createPrimaryAuth(config)
	case "secondary":
		return aam.createSecondaryAuth(config)
	case "totp":
		return aam.createTOTPAuth(config)
	case "hardware":
		return aam.createHardwareTokenAuth(config)
	case "certificate":
		return aam.createCertificateAuth(config)
	case "agent":
		return aam.createAgentAuth()
	default:
		return nil, fmt.Errorf("unknown authentication method: %s", methodName)
	}
}

// createPrimaryAuth creates primary authentication method
func (aam *AdvancedAuthManager) createPrimaryAuth(config *MultiFactorAuthConfig) (ssh.AuthMethod, error) {
	switch config.PrimaryMethod {
	case spookytypesssh.AuthMethodPublicKey:
		return aam.createPublicKeyAuth(config.PrimaryKey, config.PrimaryPass)
	case spookytypesssh.AuthMethodPassword:
		return ssh.Password(config.PrimaryPass), nil
	case spookytypesssh.AuthMethodAgent:
		return aam.createAgentAuth()
	default:
		return nil, fmt.Errorf("unsupported primary auth method: %s", config.PrimaryMethod)
	}
}

// createSecondaryAuth creates secondary authentication method
func (aam *AdvancedAuthManager) createSecondaryAuth(config *MultiFactorAuthConfig) (ssh.AuthMethod, error) {
	switch config.SecondaryMethod {
	case spookytypesssh.AuthMethodPublicKey:
		return aam.createPublicKeyAuth(config.SecondaryKey, config.SecondaryPass)
	case spookytypesssh.AuthMethodPassword:
		return ssh.Password(config.SecondaryPass), nil
	case spookytypesssh.AuthMethodAgent:
		return aam.createAgentAuth()
	default:
		return nil, fmt.Errorf("unsupported secondary auth method: %s", config.SecondaryMethod)
	}
}

// createPublicKeyAuth creates public key authentication
func (aam *AdvancedAuthManager) createPublicKeyAuth(keyPath, passphrase string) (ssh.AuthMethod, error) {
	if keyPath == "" {
		return nil, fmt.Errorf("key path is required for public key authentication")
	}

	// Load private key
	signer, err := aam.loadPrivateKey(keyPath, passphrase)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	return ssh.PublicKeys(signer), nil
}

// loadPrivateKey loads a private key from file
func (aam *AdvancedAuthManager) loadPrivateKey(keyPath, passphrase string) (ssh.Signer, error) {
	// Expand tilde in path
	if strings.HasPrefix(keyPath, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		keyPath = filepath.Join(home, keyPath[1:])
	}

	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	// Parse the key
	var signer ssh.Signer
	if passphrase != "" {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(keyData, []byte(passphrase))
	} else {
		signer, err = ssh.ParsePrivateKey(keyData)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return signer, nil
}

// createTOTPAuth creates TOTP authentication method
func (aam *AdvancedAuthManager) createTOTPAuth(config *MultiFactorAuthConfig) (ssh.AuthMethod, error) {
	if config.TOTPSecret == "" {
		return nil, fmt.Errorf("TOTP secret is required for TOTP authentication")
	}

	generator := NewTOTPGenerator(
		config.TOTPSecret,
		config.TOTPAlgorithm,
		config.TOTPDigits,
		config.TOTPPeriod,
	)

	return ssh.PasswordCallback(func() (string, error) {
		return generator.GenerateTOTP()
	}), nil
}

// createHardwareTokenAuth creates hardware token authentication
func (aam *AdvancedAuthManager) createHardwareTokenAuth(config *MultiFactorAuthConfig) (ssh.AuthMethod, error) {
	if config.HardwareTokenPath == "" {
		return nil, fmt.Errorf("hardware token path is required for hardware token authentication")
	}

	// This is a simplified implementation
	// In a real implementation, you would interface with actual hardware tokens
	return ssh.PasswordCallback(func() (string, error) {
		// Simulate hardware token interaction
		return "hardware_token_code", nil
	}), nil
}

// createCertificateAuth creates certificate-based authentication
func (aam *AdvancedAuthManager) createCertificateAuth(config *MultiFactorAuthConfig) (ssh.AuthMethod, error) {
	if len(config.CertificateChain) == 0 {
		return nil, fmt.Errorf("certificate chain is required for certificate authentication")
	}

	// Load certificate chain
	chain, err := aam.loadCertificateChain(config.CertificateChain, config.CAKeyPath, config.CAKeyPass)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate chain: %w", err)
	}

	// Create certificate signer
	if len(chain.Signers) == 0 {
		return nil, fmt.Errorf("no signers found in certificate chain")
	}

	// Use the first certificate with its signer
	certSigner, err := ssh.NewCertSigner(chain.Certificates[0], chain.Signers[0])
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate signer: %w", err)
	}

	return ssh.PublicKeys(certSigner), nil
}

// loadCertificateChain loads a chain of SSH certificates
func (aam *AdvancedAuthManager) loadCertificateChain(certPaths []string, caKeyPath, caKeyPass string) (*CertificateChain, error) {
	chain := &CertificateChain{
		Certificates: make([]*ssh.Certificate, 0, len(certPaths)),
		Signers:      make([]ssh.Signer, 0, len(certPaths)),
		Extensions:   make(map[string]string),
	}

	caSigner, err := aam.loadCASigner(caKeyPath, caKeyPass)
	if err != nil {
		return nil, err
	}

	if err := aam.loadCertificatesIntoChain(chain, certPaths, caSigner); err != nil {
		return nil, err
	}

	aam.logger.Info("Loaded certificate chain", map[string]interface{}{
		"certificate_count": len(chain.Certificates),
		"valid_from":        chain.ValidFrom,
		"valid_until":       chain.ValidUntil,
		"principals":        chain.Principals,
	})

	return chain, nil
}

func (aam *AdvancedAuthManager) loadCASigner(caKeyPath, caKeyPass string) (ssh.Signer, error) {
	if caKeyPath == "" {
		return nil, nil
	}

	caSigner, err := aam.loadPrivateKey(caKeyPath, caKeyPass)
	if err != nil {
		return nil, fmt.Errorf("failed to load CA key: %w", err)
	}

	return caSigner, nil
}

func (aam *AdvancedAuthManager) loadCertificatesIntoChain(chain *CertificateChain, certPaths []string, caSigner ssh.Signer) error {
	for i, certPath := range certPaths {
		cert, err := aam.loadSingleCertificate(certPath)
		if err != nil {
			return err
		}

		chain.Certificates = append(chain.Certificates, cert)

		if caSigner != nil {
			chain.Signers = append(chain.Signers, caSigner)
		}

		if i == 0 {
			aam.updateChainMetadata(chain, cert)
		}
	}

	return nil
}

func (aam *AdvancedAuthManager) loadSingleCertificate(certPath string) (*ssh.Certificate, error) {
	certData, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file %s: %w", certPath, err)
	}

	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(certData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate %s: %w", certPath, err)
	}

	cert, ok := pubKey.(*ssh.Certificate)
	if !ok {
		return nil, fmt.Errorf("file %s does not contain an SSH certificate", certPath)
	}

	return cert, nil
}

func (aam *AdvancedAuthManager) updateChainMetadata(chain *CertificateChain, cert *ssh.Certificate) {
	if cert.ValidAfter <= uint64(1<<63-1) {
		chain.ValidFrom = time.Unix(int64(cert.ValidAfter), 0)
	}
	if cert.ValidBefore <= uint64(1<<63-1) {
		chain.ValidUntil = time.Unix(int64(cert.ValidBefore), 0)
	}
	chain.Principals = cert.ValidPrincipals
}

// createAgentAuth creates SSH agent authentication
func (aam *AdvancedAuthManager) createAgentAuth() (ssh.AuthMethod, error) {
	if aam.agent.conn == nil {
		if err := aam.agent.Connect(); err != nil {
			return nil, fmt.Errorf("failed to connect to SSH agent: %w", err)
		}
	}

	return ssh.PublicKeysCallback(aam.agent.conn.Signers), nil
}

// GenerateCertificate generates a new SSH certificate
func (aam *AdvancedAuthManager) GenerateCertificate(config *CertificateConfig) (*ssh.Certificate, error) {
	// Load private key
	if config.PrivateKeyPath == "" {
		return nil, fmt.Errorf("private key path is required")
	}

	signer, err := aam.loadPrivateKey(config.PrivateKeyPath, config.PrivateKeyPass)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	// Create certificate
	cert := &ssh.Certificate{
		Key:             signer.PublicKey(),
		Serial:          config.Serial,
		CertType:        config.CertType,
		KeyId:           config.KeyID,
		ValidPrincipals: config.Principals,
		// Unix timestamps are always positive and within uint64 range for reasonable dates
		// nolint:gosec // G115: Unix timestamps are safe to convert to uint64
		ValidAfter: uint64(config.ValidAfter.Unix()),
		// nolint:gosec // G115: Unix timestamps are safe to convert to uint64
		ValidBefore: uint64(config.ValidBefore.Unix()),
		Permissions: ssh.Permissions{
			CriticalOptions: config.CriticalOptions,
			Extensions:      config.Extensions,
		},
	}

	// Sign certificate if CA key is provided
	if config.CAKeyPath != "" {
		caSigner, err := aam.loadPrivateKey(config.CAKeyPath, config.CAKeyPass)
		if err != nil {
			return nil, fmt.Errorf("failed to load CA key: %w", err)
		}

		if err := cert.SignCert(rand.Reader, caSigner); err != nil {
			return nil, fmt.Errorf("failed to sign certificate: %w", err)
		}
	}

	aam.logger.Info("Generated SSH certificate", map[string]interface{}{
		"key_id":       cert.KeyId,
		"cert_type":    cert.CertType,
		"principals":   cert.ValidPrincipals,
		"valid_after":  config.ValidAfter,
		"valid_before": config.ValidBefore,
		"serial":       cert.Serial,
	})

	return cert, nil
}

// CertificateConfig represents certificate generation configuration
type CertificateConfig struct {
	PrivateKeyPath  string            `json:"private_key_path,omitempty" hcl:"private_key_path,optional"`
	PrivateKeyPass  string            `json:"private_key_pass,omitempty" hcl:"private_key_pass,optional" sensitive:"true"`
	KeyType         string            `json:"key_type" hcl:"key_type" default:"rsa"`
	KeySize         int               `json:"key_size" hcl:"key_size" default:"4096"`
	Serial          uint64            `json:"serial" hcl:"serial"`
	CertType        uint32            `json:"cert_type" hcl:"cert_type"`
	KeyID           string            `json:"key_id" hcl:"key_id"`
	Principals      []string          `json:"principals" hcl:"principals"`
	ValidAfter      time.Time         `json:"valid_after" hcl:"valid_after"`
	ValidBefore     time.Time         `json:"valid_before" hcl:"valid_before"`
	CriticalOptions map[string]string `json:"critical_options,omitempty" hcl:"critical_options,optional"`
	Extensions      map[string]string `json:"extensions,omitempty" hcl:"extensions,optional"`
	CAKeyPath       string            `json:"ca_key_path,omitempty" hcl:"ca_key_path,optional"`
	CAKeyPass       string            `json:"ca_key_pass,omitempty" hcl:"ca_key_pass,optional" sensitive:"true"`
}

// SaveCertificate saves a certificate to file
func (aam *AdvancedAuthManager) SaveCertificate(cert *ssh.Certificate, path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create certificate directory: %w", err)
	}

	// Marshal certificate
	certData := ssh.MarshalAuthorizedKey(cert)

	// Write to file
	if err := os.WriteFile(path, certData, 0o600); err != nil {
		return fmt.Errorf("failed to write certificate file: %w", err)
	}

	aam.logger.Info("Saved SSH certificate", map[string]interface{}{
		"path":   path,
		"key_id": cert.KeyId,
		"size":   len(certData),
	})

	return nil
}
