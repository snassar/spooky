package ssh

import (
	"context"
	"log/slog"
	"time"

	"spooky/internal/encryption"
	"spooky/internal/logging"
	"spooky/internal/schemas"

	"github.com/pkg/errors"
)

// Manager provides SSH functionality for Spooky
type Manager struct {
	ageEncryption *encryption.AgeEncryption
	config        *schemas.SpookySSHV1
}

// NewSSHManager creates a new SSH manager
func NewSSHManager(ageEncryption *encryption.AgeEncryption, config *schemas.SpookySSHV1) *Manager {
	return &Manager{
		ageEncryption: ageEncryption,
		config:        config,
	}
}

// createSSHConfig creates a new SSH configuration for the given machine
func (sm *Manager) createSSHConfig(machine *schemas.MachinesMachineV1) *Config {
	// Use default values if config is nil
	var timeout, keepAlive, keyScanTimeout time.Duration
	var keepAliveCount int
	var knownHostsMode string

	if sm.config != nil {
		timeout = time.Duration(sm.config.Timeout) * time.Second
		keepAlive = time.Duration(sm.config.KeepaliveInterval) * time.Second
		keepAliveCount = sm.config.KeepaliveCount
		keyScanTimeout = time.Duration(sm.config.KeyScanTimeout) * time.Second
		knownHostsMode = sm.config.KnownHostsMode
	} else {
		// Default values
		timeout = 30 * time.Second
		keepAlive = 30 * time.Second
		keepAliveCount = 3
		keyScanTimeout = 5 * time.Second
		knownHostsMode = "auto"
	}

	return &Config{
		Host:           machine.Hostname,
		Port:           machine.Port,
		User:           machine.User,
		Timeout:        timeout,
		KeepAlive:      keepAlive,
		KeepAliveCount: keepAliveCount,
		KeyScanTimeout: keyScanTimeout,
		KnownHostsMode: knownHostsMode,
		PubkeyAuth:     true,
		PasswordAuth:   true,

		// Proxy configuration
		ProxyCommand: func() string {
			if sm.config != nil {
				return sm.config.ProxyCommand
			}
			return ""
		}(),
		ProxyJump: func() string {
			if sm.config != nil {
				return sm.config.ProxyJump
			}
			return ""
		}(),

		// Compression configuration
		Compression: func() bool {
			if sm.config != nil {
				return sm.config.Compression
			}
			return false
		}(),
		CompressionLevel: func() int {
			if sm.config != nil {
				return sm.config.CompressionLevel
			}
			return 6
		}(),

		// TCP keepalive configuration
		TCPKeepAlive: func() bool {
			if sm.config != nil {
				return sm.config.TCPKeepAlive
			}
			return true
		}(),
		TCPKeepAliveCount: func() int {
			if sm.config != nil {
				return sm.config.TCPKeepAliveCount
			}
			return 3
		}(),
		TCPKeepAliveIdle: func() time.Duration {
			if sm.config != nil {
				return time.Duration(sm.config.TCPKeepAliveIdle) * time.Second
			}
			return 30 * time.Second
		}(),
		TCPKeepAliveInterval: func() time.Duration {
			if sm.config != nil {
				return time.Duration(sm.config.TCPKeepAliveInterval) * time.Second
			}
			return 10 * time.Second
		}(),
		TCPKeepAliveProbeInterval: func() time.Duration {
			if sm.config != nil {
				return time.Duration(sm.config.TCPKeepAliveProbeInterval) * time.Second
			}
			return 5 * time.Second
		}(),
	}
}

// RunCommandOnMachine executes a command on the specified machine.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - machine: Machine configuration with hostname, port, user, and authentication
//   - command: Command string to execute on the remote machine
//
// Returns:
//   - *CommandResult: Command execution results (stdout, stderr, exit code)
//   - error: Connection, authentication, or execution errors
//
// Dependencies: golang.org/x/crypto/ssh, spooky/internal/encryption for age encryption
//
// Example usage:
//
//	result, err := sshManager.RunCommandOnMachine(ctx, machine, "ls -la /etc")
//	if err != nil {
//	    return fmt.Errorf("failed to execute command: %w", err)
//	}
//	fmt.Printf("Exit code: %d\n", result.ExitCode)
//
// Performance: 100ms-30s depending on command complexity and network latency
func (sm *Manager) RunCommandOnMachine(ctx context.Context, machine *schemas.MachinesMachineV1, command string) (*schemas.CommandResult, error) {
	// Check if manager is properly initialized
	if sm == nil {
		return nil, errors.New("SSH manager is not initialized")
	}

	// Create SSH config
	sshConfig := sm.createSSHConfig(machine)

	// Set up authentication
	if err := sm.setupAuthentication(sshConfig, machine); err != nil {
		return nil, errors.Wrap(err, "failed to setup authentication")
	}

	// Create SSH client
	client, err := NewClient(sshConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create SSH client")
	}
	defer func() {
		if err := client.Disconnect(); err != nil {
			// Log the disconnect error but don't fail the operation
			logger := logging.GetGlobalLogger()
			logger.Warn("failed to disconnect SSH client", slog.String("error", err.Error()))
		}
	}()

	// Connect to machine
	if err := client.Connect(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to connect to target machine - check host configuration and network connectivity")
	}

	// Execute command
	result, err := client.RunCommand(ctx, command)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute command on target machine - check command syntax and permissions")
	}

	return result, nil
}

// setupAuthentication configures authentication methods for the SSH client
func (sm *Manager) setupAuthentication(config *Config, machine *schemas.MachinesMachineV1) error {
	auth := machine.Authentication
	logger := logging.GetGlobalLogger()

	// Try public key authentication first
	err := sm.setupPublicKeyAuth(config, &auth)
	if err == nil {
		logger.Debug("SSH public key authentication configured successfully",
			slog.String("machine", machine.Hostname))
		return nil
	}
	logger.Debug("SSH public key authentication failed",
		slog.String("machine", machine.Hostname),
		slog.String("error", err.Error()))

	// Try password authentication
	err = sm.setupPasswordAuth(config, &auth)
	if err == nil {
		logger.Debug("SSH password authentication configured successfully",
			slog.String("machine", machine.Hostname))
		return nil
	}
	logger.Debug("SSH password authentication failed",
		slog.String("machine", machine.Hostname),
		slog.String("error", err.Error()))

	// Try certificate authentication
	err = sm.setupCertificateAuth(config, &auth)
	if err == nil {
		logger.Debug("SSH certificate authentication configured successfully",
			slog.String("machine", machine.Hostname))
		return nil
	}
	logger.Debug("SSH certificate authentication failed",
		slog.String("machine", machine.Hostname),
		slog.String("error", err.Error()))

	logger.Error("All SSH authentication methods failed",
		slog.String("machine", machine.Hostname))
	return errors.New("no valid authentication method found")
}

// setupPublicKeyAuth configures public key authentication
func (sm *Manager) setupPublicKeyAuth(config *Config, auth *schemas.MachinesMachineAuthenticationV1) error {
	if auth.PublicKeyPath == "" {
		return errors.New("no public key path provided")
	}

	config.PrivateKeyPath = auth.PublicKeyPath

	// Handle passphrase if present
	if auth.Passphrase.Value != "" {
		passphrase, err := sm.decryptCredential(auth.Passphrase.Value, auth.Passphrase.Encrypted, "SSH key passphrase")
		if err != nil {
			return err
		}
		config.Passphrase = passphrase
	}

	return nil
}

// setupPasswordAuth configures password authentication
func (sm *Manager) setupPasswordAuth(config *Config, auth *schemas.MachinesMachineAuthenticationV1) error {
	if auth.Password.Value == "" {
		return errors.New("no password provided")
	}

	password, err := sm.decryptCredential(auth.Password.Value, auth.Password.Encrypted, "SSH password")
	if err != nil {
		return err
	}

	config.Password = password
	return nil
}

// setupCertificateAuth configures certificate authentication
func (sm *Manager) setupCertificateAuth(config *Config, auth *schemas.MachinesMachineAuthenticationV1) error {
	if auth.PrivateKeyPath == "" || auth.CertificatePath == "" {
		return errors.New("certificate authentication requires both private key and certificate paths")
	}

	config.PrivateKeyPath = auth.PrivateKeyPath

	// Handle certificate passphrase if present
	if auth.CertificatePassphrase.Value != "" {
		passphrase, err := sm.decryptCredential(auth.CertificatePassphrase.Value, auth.CertificatePassphrase.Encrypted, "certificate passphrase")
		if err != nil {
			return err
		}
		config.Passphrase = passphrase
	}

	return nil
}

// decryptCredential handles decryption of credentials if they are encrypted
func (sm *Manager) decryptCredential(value string, encrypted bool, credentialType string) (string, error) {
	if !encrypted {
		return value, nil
	}

	if sm.ageEncryption == nil {
		return "", errors.Errorf("encrypted %s requires age encryption, but encryption is not available", credentialType)
	}

	if !sm.ageEncryption.IsEncrypted(value) {
		return "", errors.Errorf("%s marked as encrypted but does not appear to be age-encrypted", credentialType)
	}

	decrypted, err := sm.ageEncryption.Decrypt(value)
	if err != nil {
		return "", errors.Wrapf(err, "failed to decrypt %s", credentialType)
	}

	return decrypted, nil
}

// TestConnection tests the SSH connection to a machine
func (sm *Manager) TestConnection(ctx context.Context, machine *schemas.MachinesMachineV1) error {
	result, err := sm.RunCommandOnMachine(ctx, machine, "echo 'SSH connection test successful'")
	if err != nil {
		return errors.Wrapf(err, "failed to test connection to machine %s", machine.Hostname)
	}

	if result.ExitCode != 0 {
		return errors.Errorf("test command failed on machine %s: %s", machine.Hostname, result.Stderr)
	}

	logger := logging.GetGlobalLogger()
	logger.Info("SSH connection test successful", slog.String("machine", machine.Hostname))
	return nil
}

// UploadFileToMachine uploads a file to a specific machine
func (sm *Manager) UploadFileToMachine(ctx context.Context, machine *schemas.MachinesMachineV1, localPath, remotePath string) error {
	// Create SSH config
	sshConfig := sm.createSSHConfig(machine)

	// Set up authentication
	if err := sm.setupAuthentication(sshConfig, machine); err != nil {
		return errors.Wrap(err, "failed to setup authentication")
	}

	// Create SSH client
	client, err := NewClient(sshConfig)
	if err != nil {
		return errors.Wrap(err, "failed to create SSH client")
	}
	defer func() {
		if err := client.Disconnect(); err != nil {
			// Log the disconnect error but don't fail the operation
			logger := logging.GetGlobalLogger()
			logger.Warn("failed to disconnect SSH client", slog.String("error", err.Error()))
		}
	}()

	// Connect to machine
	if err := client.Connect(ctx); err != nil {
		return errors.Wrap(err, "failed to connect to target machine - check host configuration and network connectivity")
	}

	// Upload file
	if err := client.UploadFile(ctx, localPath, remotePath); err != nil {
		return errors.Wrap(err, "failed to upload file to target machine - check file permissions and disk space")
	}

	return nil
}

// GetSSHClient gets an SSH client for a specific machine
func (sm *Manager) GetSSHClient(hostname string, machine *schemas.MachinesMachineV1) (*Client, error) {
	// Create SSH config
	sshConfig := sm.createSSHConfig(machine)

	// Set up authentication
	if err := sm.setupAuthentication(sshConfig, machine); err != nil {
		return nil, errors.Wrap(err, "failed to setup authentication")
	}

	// Create and return SSH client
	client, err := NewClient(sshConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create SSH client")
	}

	return client, nil
}
