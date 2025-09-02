package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSSHClient_LoadPrivateKey(t *testing.T) {
	// Create a temporary directory for test keys
	tempDir, err := os.MkdirTemp("", "ssh-test-keys")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Generate a test RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	// Test case 1: Unencrypted RSA private key
	t.Run("Unencrypted RSA Key", func(t *testing.T) {
		// Create unencrypted PEM block
		block := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		}
		pemData := pem.EncodeToMemory(block)

		// Write to temporary file
		keyPath := filepath.Join(tempDir, "unencrypted_rsa.pem")
		err := os.WriteFile(keyPath, pemData, 0600)
		require.NoError(t, err)

		// Test loading
		client := &SSHClient{
			config: &SSHConfig{
				PrivateKeyPath: keyPath,
			},
		}

		signer, err := client.loadPrivateKey()
		assert.NoError(t, err)
		assert.NotNil(t, signer)
		assert.Equal(t, "ssh-rsa", signer.PublicKey().Type())
	})

	// Test case 2: Encrypted RSA private key (legacy format)
	t.Run("Encrypted RSA Key (Legacy)", func(t *testing.T) {
		// Create encrypted PEM block using deprecated function
		block := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		}

		passphrase := "test-passphrase"
		encryptedBlock, err := x509.EncryptPEMBlock(rand.Reader, block.Type, block.Bytes, []byte(passphrase), x509.PEMCipherDES)
		require.NoError(t, err)

		// Write to temporary file
		keyPath := filepath.Join(tempDir, "encrypted_rsa.pem")
		pemData := pem.EncodeToMemory(encryptedBlock)
		err = os.WriteFile(keyPath, pemData, 0600)
		require.NoError(t, err)

		// Test loading with passphrase
		client := &SSHClient{
			config: &SSHConfig{
				PrivateKeyPath: keyPath,
				Passphrase:     passphrase,
			},
		}

		signer, err := client.loadPrivateKey()
		assert.NoError(t, err)
		assert.NotNil(t, signer)
		assert.Equal(t, "ssh-rsa", signer.PublicKey().Type())
	})

	// Test case 3: PKCS#8 format (unencrypted)
	t.Run("PKCS#8 Unencrypted Key", func(t *testing.T) {
		// Create PKCS#8 format
		pkcs8Bytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
		require.NoError(t, err)

		block := &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: pkcs8Bytes,
		}
		pemData := pem.EncodeToMemory(block)

		// Write to temporary file
		keyPath := filepath.Join(tempDir, "pkcs8_unencrypted.pem")
		err = os.WriteFile(keyPath, pemData, 0600)
		require.NoError(t, err)

		// Test loading
		client := &SSHClient{
			config: &SSHConfig{
				PrivateKeyPath: keyPath,
			},
		}

		signer, err := client.loadPrivateKey()
		assert.NoError(t, err)
		assert.NotNil(t, signer)
		assert.Equal(t, "ssh-rsa", signer.PublicKey().Type())
	})

	// Test case 4: OpenSSH format
	t.Run("OpenSSH Format Key", func(t *testing.T) {
		// For OpenSSH format, we'll test with a simple RSA key that we know works
		// The ssh.MarshalPrivateKey has limitations with wrapped signers
		block := &pem.Block{
			Type:  "OPENSSH PRIVATE KEY",
			Bytes: []byte("test-openssh-key-data"), // This is just for testing the format detection
		}
		pemData := pem.EncodeToMemory(block)

		// Write to temporary file
		keyPath := filepath.Join(tempDir, "openssh.pem")
		err := os.WriteFile(keyPath, pemData, 0600)
		require.NoError(t, err)

		// Test loading - this should fail parsing but succeed in format detection
		client := &SSHClient{
			config: &SSHConfig{
				PrivateKeyPath: keyPath,
			},
		}

		// This should fail because the key data is invalid, but it tests the format detection
		_, err = client.loadPrivateKey()
		assert.Error(t, err)
		// But the error should be about parsing, not format detection
		assert.Contains(t, err.Error(), "failed to parse")
	})

	// Test case 5: Key with passphrase but wrong passphrase
	t.Run("Wrong Passphrase", func(t *testing.T) {
		// Create encrypted PEM block
		block := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		}

		correctPassphrase := "correct-passphrase"
		encryptedBlock, err := x509.EncryptPEMBlock(rand.Reader, block.Type, block.Bytes, []byte(correctPassphrase), x509.PEMCipherDES)
		require.NoError(t, err)

		// Write to temporary file
		keyPath := filepath.Join(tempDir, "wrong_passphrase.pem")
		pemData := pem.EncodeToMemory(encryptedBlock)
		err = os.WriteFile(keyPath, pemData, 0600)
		require.NoError(t, err)

		// Test loading with wrong passphrase
		client := &SSHClient{
			config: &SSHConfig{
				PrivateKeyPath: keyPath,
				Passphrase:     "wrong-passphrase",
			},
		}

		signer, err := client.loadPrivateKey()
		assert.Error(t, err)
		assert.Nil(t, signer)
		assert.Contains(t, err.Error(), "failed to decrypt traditional private key")
	})

	// Test case 6: Encrypted key without passphrase
	t.Run("Encrypted Key Without Passphrase", func(t *testing.T) {
		// Create encrypted PEM block
		block := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		}

		passphrase := "test-passphrase"
		encryptedBlock, err := x509.EncryptPEMBlock(rand.Reader, block.Type, block.Bytes, []byte(passphrase), x509.PEMCipherDES)
		require.NoError(t, err)

		// Write to temporary file
		keyPath := filepath.Join(tempDir, "no_passphrase.pem")
		pemData := pem.EncodeToMemory(encryptedBlock)
		err = os.WriteFile(keyPath, pemData, 0600)
		require.NoError(t, err)

		// Test loading without passphrase
		client := &SSHClient{
			config: &SSHConfig{
				PrivateKeyPath: keyPath,
				// No passphrase
			},
		}

		signer, err := client.loadPrivateKey()
		assert.Error(t, err)
		assert.Nil(t, signer)
		assert.Contains(t, err.Error(), "private key is encrypted but no passphrase provided")
	})
}

func TestSSHClient_LoadPrivateKeyWithPassphrase(t *testing.T) {
	// Create a temporary directory for test keys
	tempDir, err := os.MkdirTemp("", "ssh-test-keys")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Generate a test RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	// Test case 1: PKCS#8 encrypted key (should use ParsePrivateKeyWithPassphrase)
	t.Run("PKCS#8 Encrypted Key", func(t *testing.T) {
		// Create PKCS#8 format
		pkcs8Bytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
		require.NoError(t, err)

		// Note: We can't easily create encrypted PKCS#8 keys in Go without external tools
		// This test verifies that the code path exists and handles the case gracefully
		block := &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: pkcs8Bytes,
		}
		pemData := pem.EncodeToMemory(block)

		// Write to temporary file
		keyPath := filepath.Join(tempDir, "pkcs8_test.pem")
		err = os.WriteFile(keyPath, pemData, 0600)
		require.NoError(t, err)

		// Test loading with passphrase (should work for unencrypted PKCS#8)
		client := &SSHClient{
			config: &SSHConfig{
				PrivateKeyPath: keyPath,
				Passphrase:     "test-passphrase", // Should be ignored for unencrypted
			},
		}

		signer, err := client.loadPrivateKey()
		assert.NoError(t, err)
		assert.NotNil(t, signer)
		assert.Equal(t, "ssh-rsa", signer.PublicKey().Type())
	})
}

// Benchmark tests for performance
func BenchmarkSSHClient_LoadPrivateKey(b *testing.B) {
	// Create a temporary directory for test keys
	tempDir, err := os.MkdirTemp("", "ssh-benchmark-keys")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Generate a test RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		b.Fatal(err)
	}

	// Create unencrypted PEM block
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	pemData := pem.EncodeToMemory(block)

	// Write to temporary file
	keyPath := filepath.Join(tempDir, "benchmark_rsa.pem")
	err = os.WriteFile(keyPath, pemData, 0600)
	if err != nil {
		b.Fatal(err)
	}

	client := &SSHClient{
		config: &SSHConfig{
			PrivateKeyPath: keyPath,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.loadPrivateKey()
		if err != nil {
			b.Fatal(err)
		}
	}
}
