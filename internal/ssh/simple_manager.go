package ssh

import (
	"context"
	"fmt"
	"time"

	"spooky/internal/encryption"
	"spooky/internal/schemas"

	"github.com/pkg/errors"
)

// SimpleSSHManager provides basic SSH functionality for Spooky
type SimpleSSHManager struct {
	ageEncryption *encryption.AgeEncryption
	config        *schemas.SpookySSHV1
}

// NewSimpleSSHManager creates a new simple SSH manager
func NewSimpleSSHManager(ageEncryption *encryption.AgeEncryption, config *schemas.SpookySSHV1) *SimpleSSHManager {
	return &SimpleSSHManager{
		ageEncryption: ageEncryption,
		config:        config,
	}
}

// ExecuteCommandOnMachine executes a command on a specific machine
func (sm *SimpleSSHManager) ExecuteCommandOnMachine(ctx context.Context, machine *schemas.MachinesMachineV1, command string) (*CommandResult, error) {
	// Create SSH config
	sshConfig := &SSHConfig{
		Host:           machine.Hostname,
		Port:           machine.Port,
		User:           machine.User,
		Timeout:        time.Duration(sm.config.Timeout) * time.Second,
		KeepAlive:      time.Duration(sm.config.KeepaliveInterval) * time.Second,
		KeepAliveCount: sm.config.KeepaliveCount,
		KeyScanTimeout: time.Duration(sm.config.KeyScanTimeout) * time.Second,
		StrictHostKey:  sm.config.KnownHostsStrict,
		KnownHostsMode: sm.config.KnownHostsMode,
		PubkeyAuth:     true,
		PasswordAuth:   true,

		// Proxy configuration
		ProxyCommand: sm.config.ProxyCommand,
		ProxyJump:    sm.config.ProxyJump,

		// Compression configuration
		Compression:      sm.config.Compression,
		CompressionLevel: sm.config.CompressionLevel,

		// TCP keepalive configuration
		TCPKeepAlive:              sm.config.TCPKeepAlive,
		TCPKeepAliveCount:         sm.config.TCPKeepAliveCount,
		TCPKeepAliveIdle:          sm.config.TCPKeepAliveIdle,
		TCPKeepAliveInterval:      sm.config.TCPKeepAliveInterval,
		TCPKeepAliveProbeInterval: sm.config.TCPKeepAliveProbeInterval,
	}

	// Set up authentication
	if err := sm.setupAuthentication(sshConfig, machine); err != nil {
		return nil, errors.Wrap(err, "failed to setup authentication")
	}

	// Create SSH client
	client, err := NewSSHClient(sshConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create SSH client")
	}
	defer client.Disconnect()

	// Connect to machine
	if err := client.Connect(ctx); err != nil {
		return nil, errors.Wrapf(err, "failed to connect to machine %s", machine.Hostname)
	}

	// Execute command
	result, err := client.ExecuteCommand(ctx, command)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute command on machine %s", machine.Hostname)
	}

	return result, nil
}

// setupAuthentication configures authentication methods for the SSH client
func (sm *SimpleSSHManager) setupAuthentication(config *SSHConfig, machine *schemas.MachinesMachineV1) error {
	auth := machine.Authentication

	// Check if we have public key authentication
	if auth.PublicKeyPath != "" {
		config.PrivateKeyPath = auth.PublicKeyPath

		// Handle passphrase if present
		if auth.Passphrase.Value != "" {
			passphrase := auth.Passphrase.Value

			// Decrypt passphrase if it's encrypted
			if auth.Passphrase.Encrypted && sm.ageEncryption != nil {
				if sm.ageEncryption.IsEncrypted(passphrase) {
					decrypted, err := sm.ageEncryption.Decrypt(passphrase)
					if err != nil {
						return errors.Wrap(err, "failed to decrypt SSH key passphrase")
					}
					passphrase = decrypted
				}
			}

			config.Passphrase = passphrase
		}
		return nil
	}

	// Check if we have password authentication
	if auth.Password.Value != "" {
		passphrase := auth.Password.Value

		// Decrypt password if it's encrypted
		if auth.Password.Encrypted && sm.ageEncryption != nil {
			if sm.ageEncryption.IsEncrypted(passphrase) {
				decrypted, err := sm.ageEncryption.Decrypt(passphrase)
				if err != nil {
					return errors.Wrap(err, "failed to decrypt SSH password")
				}
				passphrase = decrypted
			}
		}

		config.Password = passphrase
		return nil
	}

	// Check if we have certificate authentication
	if auth.PrivateKeyPath != "" && auth.CertificatePath != "" {
		config.PrivateKeyPath = auth.PrivateKeyPath

		// Handle certificate passphrase if present
		if auth.CertificatePassphrase.Value != "" {
			passphrase := auth.CertificatePassphrase.Value

			// Decrypt passphrase if it's encrypted
			if auth.CertificatePassphrase.Encrypted && sm.ageEncryption != nil {
				if sm.ageEncryption.IsEncrypted(passphrase) {
					decrypted, err := sm.ageEncryption.Decrypt(passphrase)
					if err != nil {
						return errors.Wrap(err, "failed to decrypt certificate passphrase")
					}
					passphrase = decrypted
				}
			}

			config.Passphrase = passphrase
		}
		return nil
	}

	return errors.New("no valid authentication method found")
}

// TestConnection tests the SSH connection to a machine
func (sm *SimpleSSHManager) TestConnection(ctx context.Context, machine *schemas.MachinesMachineV1) error {
	result, err := sm.ExecuteCommandOnMachine(ctx, machine, "echo 'SSH connection test successful'")
	if err != nil {
		return errors.Wrapf(err, "failed to test connection to machine %s", machine.Hostname)
	}

	if result.ExitCode != 0 {
		return errors.Errorf("test command failed on machine %s: %s", machine.Hostname, result.Stderr)
	}

	fmt.Printf("SSH connection test successful for machine %s\n", machine.Hostname)
	return nil
}

// UploadFileToMachine uploads a file to a specific machine
func (sm *SimpleSSHManager) UploadFileToMachine(ctx context.Context, machine *schemas.MachinesMachineV1, localPath, remotePath string) error {
	// Create SSH config
	sshConfig := &SSHConfig{
		Host:           machine.Hostname,
		Port:           machine.Port,
		User:           machine.User,
		Timeout:        time.Duration(sm.config.Timeout) * time.Second,
		KeepAlive:      time.Duration(sm.config.KeepaliveInterval) * time.Second,
		KeepAliveCount: sm.config.KeepaliveCount,
		KeyScanTimeout: time.Duration(sm.config.KeyScanTimeout) * time.Second,
		StrictHostKey:  sm.config.KnownHostsStrict,
		KnownHostsMode: sm.config.KnownHostsMode,
		PubkeyAuth:     true,
		PasswordAuth:   true,

		// Proxy configuration
		ProxyCommand: sm.config.ProxyCommand,
		ProxyJump:    sm.config.ProxyJump,

		// Compression configuration
		Compression:      sm.config.Compression,
		CompressionLevel: sm.config.CompressionLevel,

		// TCP keepalive configuration
		TCPKeepAlive:              sm.config.TCPKeepAlive,
		TCPKeepAliveCount:         sm.config.TCPKeepAliveCount,
		TCPKeepAliveIdle:          sm.config.TCPKeepAliveIdle,
		TCPKeepAliveInterval:      sm.config.TCPKeepAliveInterval,
		TCPKeepAliveProbeInterval: sm.config.TCPKeepAliveProbeInterval,
	}

	// Set up authentication
	if err := sm.setupAuthentication(sshConfig, machine); err != nil {
		return errors.Wrap(err, "failed to setup authentication")
	}

	// Create SSH client
	client, err := NewSSHClient(sshConfig)
	if err != nil {
		return errors.Wrap(err, "failed to create SSH client")
	}
	defer client.Disconnect()

	// Connect to machine
	if err := client.Connect(ctx); err != nil {
		return errors.Wrapf(err, "failed to connect to machine %s", machine.Hostname)
	}

	// Upload file
	if err := client.UploadFile(ctx, localPath, remotePath); err != nil {
		return errors.Wrapf(err, "failed to upload file to machine %s", machine.Hostname)
	}

	return nil
}
