// Package ssh provides SSH client functionality for the spooky codebase.
// This package implements SSH connections, authentication, and command execution.
package ssh

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
)

// Supported key types
const (
	KeyTypeED25519   = "ed25519"
	KeyTypeED25519SK = "ed25519-sk"
	KeyTypeRSA4096   = "rsa-4096"
	MinRSAKeySize    = 4096
)

// KeyValidationError represents key validation errors
type KeyValidationError struct {
	KeyType string
	Reason  string
}

func (e *KeyValidationError) Error() string {
	return fmt.Sprintf("key validation failed for %s: %s", e.KeyType, e.Reason)
}

// SimpleClient implements basic SSH client functionality
type SimpleClient struct {
	config      *spookytypes.ClientConfig
	logger      spookytypeslogging.Logger
	connections map[string]*ssh.Client
	mu          sync.RWMutex
	closed      bool
}

// NewSimpleClient creates a new SSH client
func NewSimpleClient(config *spookytypes.ClientConfig, logger spookytypeslogging.Logger) *SimpleClient {
	if config == nil {
		config = &spookytypes.ClientConfig{
			DefaultPort:      22,
			DefaultTimeout:   30 * time.Second,
			MaxConnections:   10,
			MaxRetryAttempts: 3,
			RetryDelay:       5 * time.Second,
			IdleTimeout:      300 * time.Second,
		}
	}

	return &SimpleClient{
		config:      config,
		logger:      logger,
		connections: make(map[string]*ssh.Client),
	}
}

// Connect establishes an SSH connection
func (c *SimpleClient) Connect(ctx context.Context, request *spookytypes.ConnectionRequest) (*spookytypes.ConnectionResult, error) {
	if c.closed {
		return nil, fmt.Errorf("client is closed")
	}

	startTime := time.Now()
	connectionKey := fmt.Sprintf("%s:%d", request.Host, request.Port)

	// Check if we already have a connection
	c.mu.RLock()
	if conn, exists := c.connections[connectionKey]; exists {
		c.mu.RUnlock()
		// Test if connection is still alive
		if _, _, err := conn.SendRequest("keepalive@openssh.com", true, nil); err == nil {
			return &spookytypes.ConnectionResult{
				Connection: &spookytypes.Connection{
					Host:        request.Host,
					Port:        request.Port,
					User:        request.User,
					Status:      spookytypes.ConnectionStatusConnected,
					ConnectedAt: &startTime,
				},
				Request:     request,
				Success:     true,
				ConnectTime: time.Since(startTime),
				CompletedAt: time.Now(),
			}, nil
		}
		// Connection is dead, remove it
		c.mu.RUnlock()
		c.mu.Lock()
		delete(c.connections, connectionKey)
		c.mu.Unlock()
	} else {
		c.mu.RUnlock()
	}

	// Create SSH client config
	sshConfig, err := c.createSSHConfig(request)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH config: %w", err)
	}

	// Establish connection with retries
	var conn *ssh.Client
	var lastErr error
	for attempt := 1; attempt <= c.config.MaxRetryAttempts; attempt++ {
		conn, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", request.Host, request.Port), sshConfig)
		if err == nil {
			break
		}

		lastErr = err
		if attempt < c.config.MaxRetryAttempts {
			time.Sleep(c.config.RetryDelay)
		}
	}

	if conn == nil {
		return &spookytypes.ConnectionResult{
			Request:       request,
			Success:       false,
			Error:         lastErr.Error(),
			ConnectTime:   time.Since(startTime),
			RetryAttempts: c.config.MaxRetryAttempts,
			CompletedAt:   time.Now(),
		}, nil
	}

	// Store connection
	c.mu.Lock()
	c.connections[connectionKey] = conn
	c.mu.Unlock()

	// Get connection info
	clientVersion := conn.ClientVersion()
	serverVersion := conn.ServerVersion()

	connection := &spookytypes.Connection{
		Host:          request.Host,
		Port:          request.Port,
		User:          request.User,
		Status:        spookytypes.ConnectionStatusConnected,
		ConnectedAt:   &startTime,
		ClientVersion: string(clientVersion),
		ServerVersion: string(serverVersion),
	}

	result := &spookytypes.ConnectionResult{
		Connection:  connection,
		Request:     request,
		Success:     true,
		ConnectTime: time.Since(startTime),
		CompletedAt: time.Now(),
	}

	return result, nil
}

// createSSHConfig creates an SSH client configuration
func (c *SimpleClient) createSSHConfig(request *spookytypes.ConnectionRequest) (*ssh.ClientConfig, error) {
	config := &ssh.ClientConfig{
		User:            request.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: Implement proper host key verification
		Timeout:         request.Timeout,
	}

	// Add authentication methods
	var authMethods []ssh.AuthMethod

	// Public key authentication with certificate support
	if request.KeyPath != "" {
		key, err := c.loadPrivateKey(request.KeyPath, request.Passphrase)
		if err != nil {
			return nil, fmt.Errorf("failed to load private key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(key))
	}

	// SSH certificate authentication
	if request.CertificatePath != "" {
		cert, err := c.loadSSHCertificate(request.CertificatePath, request.KeyPath, request.Passphrase)
		if err != nil {
			return nil, fmt.Errorf("failed to load SSH certificate: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(cert))
	}

	// Password authentication
	if request.Password != "" {
		authMethods = append(authMethods, ssh.Password(request.Password))
	}

	// If no authentication method is provided, try to use default key
	if len(authMethods) == 0 && c.config.DefaultKeyPath != "" {
		key, err := c.loadPrivateKey(c.config.DefaultKeyPath, "")
		if err != nil {
			return nil, fmt.Errorf("failed to load default private key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(key))
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("no authentication method provided")
	}

	config.Auth = authMethods
	return config, nil
}

// loadPrivateKey loads and validates a private key from file
func (c *SimpleClient) loadPrivateKey(keyPath, passphrase string) (ssh.Signer, error) {
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

	// Validate the key type
	if err := c.validateKeyType(signer); err != nil {
		return nil, fmt.Errorf("key validation failed: %w", err)
	}

	return signer, nil
}

// loadSSHCertificate loads an SSH certificate with its private key
func (c *SimpleClient) loadSSHCertificate(certPath, keyPath, passphrase string) (ssh.Signer, error) {
	// Expand tilde in paths
	if strings.HasPrefix(certPath, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		certPath = filepath.Join(home, certPath[1:])
	}

	if strings.HasPrefix(keyPath, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		keyPath = filepath.Join(home, keyPath[1:])
	}

	// Load certificate
	certData, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}

	// Load private key
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	// Parse private key
	var signer ssh.Signer
	if passphrase != "" {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(keyData, []byte(passphrase))
	} else {
		signer, err = ssh.ParsePrivateKey(keyData)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Parse certificate
	cert, err := ssh.ParsePublicKey(certData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH certificate: %w", err)
	}

	// Create certificate signer
	certSigner, err := ssh.NewCertSigner(cert.(*ssh.Certificate), signer)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate signer: %w", err)
	}

	c.logger.Info("Loaded SSH certificate", map[string]interface{}{
		"certificate_path": certPath,
		"key_path":         keyPath,
		"certificate_type": cert.Type(),
	})

	return certSigner, nil
}

// validateKeyType validates that the key is of a supported type
func (c *SimpleClient) validateKeyType(signer ssh.Signer) error {
	pubKey := signer.PublicKey()
	keyType := pubKey.Type()

	switch keyType {
	case ssh.KeyAlgoED25519:
		// Validate ed25519 key
		if err := c.validateED25519Key(pubKey); err != nil {
			return &KeyValidationError{KeyType: KeyTypeED25519, Reason: err.Error()}
		}
		c.logger.Info("Validated ed25519 key", map[string]interface{}{
			"key_type": KeyTypeED25519,
		})

	case ssh.KeyAlgoRSA:
		// Validate RSA key size
		if err := c.validateRSAKey(pubKey); err != nil {
			return &KeyValidationError{KeyType: KeyTypeRSA4096, Reason: err.Error()}
		}
		c.logger.Info("Validated RSA key", map[string]interface{}{
			"key_type": KeyTypeRSA4096,
		})

	default:
		return &KeyValidationError{
			KeyType: keyType,
			Reason: fmt.Sprintf("unsupported key type: %s. Supported types: %s, %s, %s",
				keyType, KeyTypeED25519, KeyTypeED25519SK, KeyTypeRSA4096),
		}
	}

	return nil
}

// validateED25519Key validates an ed25519 key
func (c *SimpleClient) validateED25519Key(pubKey ssh.PublicKey) error {
	// ed25519 keys are always valid - they have a fixed size
	// We just need to ensure it's actually an ed25519 key
	if pubKey.Type() != ssh.KeyAlgoED25519 {
		return fmt.Errorf("expected ed25519 key, got %s", pubKey.Type())
	}
	return nil
}

// validateRSAKey validates an RSA key (must be 4096-bit)
func (c *SimpleClient) validateRSAKey(pubKey ssh.PublicKey) error {
	if pubKey.Type() != ssh.KeyAlgoRSA {
		return fmt.Errorf("expected RSA key, got %s", pubKey.Type())
	}

	// Extract RSA public key to check key size
	rsaPubKey, ok := pubKey.(ssh.CryptoPublicKey)
	if !ok {
		return fmt.Errorf("failed to extract RSA public key")
	}

	cryptoPubKey := rsaPubKey.CryptoPublicKey()
	rsaKey, ok := cryptoPubKey.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("failed to cast to RSA public key")
	}

	// Check key size
	if rsaKey.Size()*8 < MinRSAKeySize {
		return fmt.Errorf("RSA key size %d bits is less than minimum required %d bits",
			rsaKey.Size()*8, MinRSAKeySize)
	}

	return nil
}

// generateSupportedKey generates a supported key type for testing
func (c *SimpleClient) generateSupportedKey(keyType string) (ssh.Signer, error) {
	switch keyType {
	case KeyTypeED25519:
		return c.generateED25519Key()
	case KeyTypeRSA4096:
		return c.generateRSA4096Key()
	default:
		return nil, fmt.Errorf("unsupported key type for generation: %s", keyType)
	}
}

// generateED25519Key generates an ed25519 key pair
func (c *SimpleClient) generateED25519Key() (ssh.Signer, error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ed25519 key: %w", err)
	}

	// Convert to SSH format
	sshPubKey, err := ssh.NewPublicKey(pubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH public key: %w", err)
	}

	// Create SSH signer
	signer, err := ssh.NewSignerFromKey(privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH signer: %w", err)
	}

	c.logger.Info("Generated ed25519 key", map[string]interface{}{
		"key_type":    KeyTypeED25519,
		"fingerprint": ssh.FingerprintSHA256(sshPubKey),
	})

	return signer, nil
}

// generateRSA4096Key generates a 4096-bit RSA key pair
func (c *SimpleClient) generateRSA4096Key() (ssh.Signer, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, MinRSAKeySize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	// Create SSH signer
	signer, err := ssh.NewSignerFromKey(privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH signer: %w", err)
	}

	// Get public key for fingerprint
	pubKey := signer.PublicKey()

	c.logger.Info("Generated RSA key", map[string]interface{}{
		"key_type":    KeyTypeRSA4096,
		"key_size":    MinRSAKeySize,
		"fingerprint": ssh.FingerprintSHA256(pubKey),
	})

	return signer, nil
}

// ExecuteCommand executes a command via SSH
func (c *SimpleClient) ExecuteCommand(ctx context.Context, connection *spookytypes.Connection, command *spookytypes.SSHCommand) (*spookytypes.SSHCommandResult, error) {
	if c.closed {
		return nil, fmt.Errorf("client is closed")
	}

	startTime := time.Now()
	connectionKey := fmt.Sprintf("%s:%d", connection.Host, connection.Port)

	// Get SSH client
	c.mu.RLock()
	sshClient, exists := c.connections[connectionKey]
	c.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no SSH connection found for %s", connectionKey)
	}

	// Create session
	session, err := sshClient.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	// Set up environment
	if len(command.Environment) > 0 {
		for key, value := range command.Environment {
			session.Setenv(key, value)
		}
	}

	// Set up input/output
	var stdout, stderr strings.Builder
	if command.CaptureOutput {
		session.Stdout = &stdout
		session.Stderr = &stderr
	}

	if command.Stdin != "" {
		session.Stdin = strings.NewReader(command.Stdin)
	}

	// Execute command
	cmd := command.Command
	if len(command.Args) > 0 {
		cmd = cmd + " " + strings.Join(command.Args, " ")
	}

	err = session.Run(cmd)
	endTime := time.Now()

	result := &spookytypes.SSHCommandResult{
		Command: command,
		Session: &spookytypes.Session{
			SessionID:  fmt.Sprintf("%s-%d", connectionKey, startTime.UnixNano()),
			Connection: connection,
			Status:     spookytypes.SessionStatusCompleted,
			StartedAt:  startTime,
			EndedAt:    &endTime,
		},
		StartTime: startTime,
		EndTime:   endTime,
		Duration:  endTime.Sub(startTime),
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		if exitErr, ok := err.(*ssh.ExitError); ok {
			result.ExitCode = exitErr.ExitStatus()
		} else {
			result.ExitCode = -1
		}
	} else {
		result.Success = true
		result.ExitCode = 0
	}

	if command.CaptureOutput {
		result.Stdout = stdout.String()
		result.Stderr = stderr.String()
	}

	return result, nil
}

// Close closes all SSH connections
func (c *SimpleClient) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	for key, conn := range c.connections {
		if err := conn.Close(); err != nil {
			// Log warning but continue
		}
		delete(c.connections, key)
	}

	c.closed = true
	return nil
}
