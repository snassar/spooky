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

// createSSHConfig creates a new SSH configuration for the given machine
func (sm *SimpleSSHManager) createSSHConfig(machine *schemas.MachinesMachineV1) *SSHConfig {
	return &SSHConfig{
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
}

// ExecuteCommandOnMachine executes a command on a remote machine via SSH.
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
//	result, err := sshManager.ExecuteCommandOnMachine(ctx, machine, "ls -la /etc")
//	if err != nil {
//	    return fmt.Errorf("failed to execute command: %w", err)
//	}
//	fmt.Printf("Exit code: %d\n", result.ExitCode)
//
// Performance: 100ms-30s depending on command complexity and network latency
func (sm *SimpleSSHManager) RunCommandOnMachine(ctx context.Context, machine *schemas.MachinesMachineV1, command string) (*schemas.CommandResult, error) {
	// Create SSH config
	sshConfig := sm.createSSHConfig(machine)

	// Set up authentication
	if err := sm.setupAuthentication(sshConfig, machine); err != nil {
		return nil, errors.Wrap(err, "failed to setup authentication")
	}

	// Create SSH client
	client, err := NewSSHClient(sshConfig)
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
		return nil, errors.Wrapf(err, "failed to connect to machine %s", machine.Hostname)
	}

	// Execute command
	result, err := client.RunCommand(ctx, command)
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
			if auth.Passphrase.Encrypted {
				if sm.ageEncryption == nil {
					return errors.New("encrypted passphrase requires age encryption, but encryption is not available")
				}
				if sm.ageEncryption.IsEncrypted(passphrase) {
					decrypted, err := sm.ageEncryption.Decrypt(passphrase)
					if err != nil {
						return errors.Wrap(err, "failed to decrypt SSH key passphrase")
					}
					passphrase = decrypted
				} else {
					return errors.New("passphrase marked as encrypted but does not appear to be age-encrypted")
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
		if auth.Password.Encrypted {
			if sm.ageEncryption == nil {
				return errors.New("encrypted password requires age encryption, but encryption is not available")
			}
			if sm.ageEncryption.IsEncrypted(passphrase) {
				decrypted, err := sm.ageEncryption.Decrypt(passphrase)
				if err != nil {
					return errors.Wrap(err, "failed to decrypt SSH password")
				}
				passphrase = decrypted
			} else {
				return errors.New("password marked as encrypted but does not appear to be age-encrypted")
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
			if auth.CertificatePassphrase.Encrypted {
				if sm.ageEncryption == nil {
					return errors.New("encrypted certificate passphrase requires age encryption, but encryption is not available")
				}
				if sm.ageEncryption.IsEncrypted(passphrase) {
					decrypted, err := sm.ageEncryption.Decrypt(passphrase)
					if err != nil {
						return errors.Wrap(err, "failed to decrypt certificate passphrase")
					}
					passphrase = decrypted
				} else {
					return errors.New("certificate passphrase marked as encrypted but does not appear to be age-encrypted")
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
func (sm *SimpleSSHManager) UploadFileToMachine(ctx context.Context, machine *schemas.MachinesMachineV1, localPath, remotePath string) error {
	// Create SSH config
	sshConfig := sm.createSSHConfig(machine)

	// Set up authentication
	if err := sm.setupAuthentication(sshConfig, machine); err != nil {
		return errors.Wrap(err, "failed to setup authentication")
	}

	// Create SSH client
	client, err := NewSSHClient(sshConfig)
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
		return errors.Wrapf(err, "failed to connect to machine %s", machine.Hostname)
	}

	// Upload file
	if err := client.UploadFile(ctx, localPath, remotePath); err != nil {
		return errors.Wrapf(err, "failed to upload file to machine %s", machine.Hostname)
	}

	return nil
}
