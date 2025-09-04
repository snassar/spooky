package ssh

import (
	"testing"

	"spooky/internal/schemas"

	"github.com/stretchr/testify/assert"
)

// Mock encryption for testing - we'll create a simple test without encryption for most tests
// and use real encryption for integration tests

func TestSimpleSSHManager_SetupAuthentication(t *testing.T) {
	t.Run("Public Key Authentication", func(t *testing.T) {
		manager := &SimpleSSHManager{}
		config := &Config{}
		machine := &schemas.MachinesMachineV1{
			Authentication: schemas.MachinesMachineAuthenticationV1{
				PublicKeyPath: "/path/to/private/key",
				Passphrase: schemas.MachinesMachineAuthenticationPassphraseV1{
					Value:     "test-passphrase",
					Encrypted: false,
				},
			},
		}

		err := manager.setupAuthentication(config, machine)
		assert.NoError(t, err)
		assert.Equal(t, "/path/to/private/key", config.PrivateKeyPath)
		assert.Equal(t, "test-passphrase", config.Passphrase)
	})

	t.Run("Password Authentication", func(t *testing.T) {
		manager := &SimpleSSHManager{}
		config := &Config{}
		machine := &schemas.MachinesMachineV1{
			Authentication: schemas.MachinesMachineAuthenticationV1{
				Password: schemas.MachinesMachineAuthenticationPasswordV1{
					Value:     "test-password",
					Encrypted: false,
				},
			},
		}

		err := manager.setupAuthentication(config, machine)
		assert.NoError(t, err)
		assert.Equal(t, "test-password", config.Password)
	})

	t.Run("Certificate Authentication", func(t *testing.T) {
		manager := &SimpleSSHManager{}
		config := &Config{}
		machine := &schemas.MachinesMachineV1{
			Authentication: schemas.MachinesMachineAuthenticationV1{
				PrivateKeyPath:  "/path/to/private/key",
				CertificatePath: "/path/to/certificate",
				CertificatePassphrase: schemas.MachinesMachineAuthenticationPassphraseV1{
					Value:     "cert-passphrase",
					Encrypted: false,
				},
			},
		}

		err := manager.setupAuthentication(config, machine)
		assert.NoError(t, err)
		assert.Equal(t, "/path/to/private/key", config.PrivateKeyPath)
		assert.Equal(t, "cert-passphrase", config.Passphrase)
	})

	t.Run("No Authentication Method", func(t *testing.T) {
		manager := &SimpleSSHManager{}
		config := &Config{}
		machine := &schemas.MachinesMachineV1{
			Authentication: schemas.MachinesMachineAuthenticationV1{},
		}

		err := manager.setupAuthentication(config, machine)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no valid authentication method found")
	})
}

func TestSimpleSSHManager_SetupPublicKeyAuth(t *testing.T) {
	manager := &SimpleSSHManager{}
	config := &Config{}

	t.Run("Valid Public Key", func(t *testing.T) {
		auth := &schemas.MachinesMachineAuthenticationV1{
			PublicKeyPath: "/path/to/key",
			Passphrase: schemas.MachinesMachineAuthenticationPassphraseV1{
				Value:     "passphrase",
				Encrypted: false,
			},
		}

		err := manager.setupPublicKeyAuth(config, auth)
		assert.NoError(t, err)
		assert.Equal(t, "/path/to/key", config.PrivateKeyPath)
		assert.Equal(t, "passphrase", config.Passphrase)
	})

	t.Run("No Public Key Path", func(t *testing.T) {
		auth := &schemas.MachinesMachineAuthenticationV1{}

		err := manager.setupPublicKeyAuth(config, auth)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no public key path provided")
	})
}

func TestSimpleSSHManager_SetupPasswordAuth(t *testing.T) {
	manager := &SimpleSSHManager{}
	config := &Config{}

	t.Run("Valid Password", func(t *testing.T) {
		auth := &schemas.MachinesMachineAuthenticationV1{
			Password: schemas.MachinesMachineAuthenticationPasswordV1{
				Value:     "password",
				Encrypted: false,
			},
		}

		err := manager.setupPasswordAuth(config, auth)
		assert.NoError(t, err)
		assert.Equal(t, "password", config.Password)
	})

	t.Run("No Password", func(t *testing.T) {
		auth := &schemas.MachinesMachineAuthenticationV1{}

		err := manager.setupPasswordAuth(config, auth)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no password provided")
	})
}

func TestSimpleSSHManager_SetupCertificateAuth(t *testing.T) {
	manager := &SimpleSSHManager{}
	config := &Config{}

	t.Run("Valid Certificate", func(t *testing.T) {
		auth := &schemas.MachinesMachineAuthenticationV1{
			PrivateKeyPath:  "/path/to/key",
			CertificatePath: "/path/to/cert",
			CertificatePassphrase: schemas.MachinesMachineAuthenticationPassphraseV1{
				Value:     "passphrase",
				Encrypted: false,
			},
		}

		err := manager.setupCertificateAuth(config, auth)
		assert.NoError(t, err)
		assert.Equal(t, "/path/to/key", config.PrivateKeyPath)
		assert.Equal(t, "passphrase", config.Passphrase)
	})

	t.Run("Missing Private Key Path", func(t *testing.T) {
		auth := &schemas.MachinesMachineAuthenticationV1{
			CertificatePath: "/path/to/cert",
		}

		err := manager.setupCertificateAuth(config, auth)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "certificate authentication requires both private key and certificate paths")
	})

	t.Run("Missing Certificate Path", func(t *testing.T) {
		auth := &schemas.MachinesMachineAuthenticationV1{
			PrivateKeyPath: "/path/to/key",
		}

		err := manager.setupCertificateAuth(config, auth)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "certificate authentication requires both private key and certificate paths")
	})
}

func TestSimpleSSHManager_DecryptCredential(t *testing.T) {
	t.Run("Plain Credential", func(t *testing.T) {
		manager := &SimpleSSHManager{}

		result, err := manager.decryptCredential("plain-value", false, "test credential")
		assert.NoError(t, err)
		assert.Equal(t, "plain-value", result)
	})

	t.Run("Encrypted Credential - No Encryption Available", func(t *testing.T) {
		manager := &SimpleSSHManager{}

		result, err := manager.decryptCredential("encrypted-value", true, "test credential")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "encrypted test credential requires age encryption")
		assert.Empty(t, result)
	})
}
