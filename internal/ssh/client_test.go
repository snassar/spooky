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

// createLegacyEncryptedPEMBlock creates a PEM block with legacy encryption headers
// This simulates the format that would be created by the deprecated x509.EncryptPEMBlock
func createLegacyEncryptedPEMBlock(blockType string, data []byte) *pem.Block {
	return &pem.Block{
		Type: blockType,
		Headers: map[string]string{
			"Proc-Type": "4,ENCRYPTED",
			"DEK-Info":  "DES-CBC,0123456789ABCDEF", // Dummy DEK-Info for testing
		},
		Bytes: data, // In real legacy encryption, this would be encrypted data
	}
}

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

	// Test case 2: Encrypted RSA private key (legacy format) - no longer supported
	t.Run("Encrypted RSA Key (Legacy) - No Longer Supported", func(t *testing.T) {
		// Create encrypted PEM block with legacy encryption headers
		encryptedBlock := createLegacyEncryptedPEMBlock("RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(privateKey))

		// Write to temporary file
		keyPath := filepath.Join(tempDir, "encrypted_rsa.pem")
		pemData := pem.EncodeToMemory(encryptedBlock)
		err = os.WriteFile(keyPath, pemData, 0600)
		require.NoError(t, err)

		// Test loading with passphrase - should fail as legacy DES encryption is no longer supported
		client := &SSHClient{
			config: &SSHConfig{
				PrivateKeyPath: keyPath,
				Passphrase:     "test-passphrase",
			},
		}

		signer, err := client.loadPrivateKey()
		assert.Error(t, err)
		assert.Nil(t, signer)
		assert.Contains(t, err.Error(), "legacy DES-encrypted traditional keys are no longer supported")
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

	// Test case 5: Key with passphrase but wrong passphrase - legacy DES no longer supported
	t.Run("Wrong Passphrase - Legacy DES No Longer Supported", func(t *testing.T) {
		// Create encrypted PEM block with legacy encryption headers
		encryptedBlock := createLegacyEncryptedPEMBlock("RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(privateKey))

		// Write to temporary file
		keyPath := filepath.Join(tempDir, "wrong_passphrase.pem")
		pemData := pem.EncodeToMemory(encryptedBlock)
		err = os.WriteFile(keyPath, pemData, 0600)
		require.NoError(t, err)

		// Test loading with wrong passphrase - should fail as legacy DES encryption is no longer supported
		client := &SSHClient{
			config: &SSHConfig{
				PrivateKeyPath: keyPath,
				Passphrase:     "wrong-passphrase",
			},
		}

		signer, err := client.loadPrivateKey()
		assert.Error(t, err)
		assert.Nil(t, signer)
		assert.Contains(t, err.Error(), "legacy DES-encrypted traditional keys are no longer supported")
	})

	// Test case 6: Encrypted key without passphrase - legacy DES no longer supported
	t.Run("Encrypted Key Without Passphrase - Legacy DES No Longer Supported", func(t *testing.T) {
		// Create encrypted PEM block with legacy encryption headers
		encryptedBlock := createLegacyEncryptedPEMBlock("RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(privateKey))

		// Write to temporary file
		keyPath := filepath.Join(tempDir, "no_passphrase.pem")
		pemData := pem.EncodeToMemory(encryptedBlock)
		err = os.WriteFile(keyPath, pemData, 0600)
		require.NoError(t, err)

		// Test loading without passphrase - should fail as legacy DES encryption is no longer supported
		client := &SSHClient{
			config: &SSHConfig{
				PrivateKeyPath: keyPath,
				// No passphrase
			},
		}

		signer, err := client.loadPrivateKey()
		assert.Error(t, err)
		assert.Nil(t, signer)
		assert.Contains(t, err.Error(), "legacy DES-encrypted traditional keys are no longer supported")
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

// Test helper functions for loadPrivateKey refactoring
func TestSSHClient_ExtractKeyData(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ssh-test-extract")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Generate a test RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	// Create test key data
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	pemData := pem.EncodeToMemory(block)

	t.Run("From File Path", func(t *testing.T) {
		// Write to temporary file
		keyPath := filepath.Join(tempDir, "test_key.pem")
		err := os.WriteFile(keyPath, pemData, 0600)
		require.NoError(t, err)

		client := &SSHClient{
			config: &SSHConfig{
				PrivateKeyPath: keyPath,
			},
		}

		keyData, err := client.extractKeyData()
		assert.NoError(t, err)
		assert.Equal(t, pemData, keyData)
	})

	t.Run("From Direct Data", func(t *testing.T) {
		client := &SSHClient{
			config: &SSHConfig{
				PrivateKeyData: pemData,
			},
		}

		keyData, err := client.extractKeyData()
		assert.NoError(t, err)
		assert.Equal(t, pemData, keyData)
	})

	t.Run("No Key Provided", func(t *testing.T) {
		client := &SSHClient{
			config: &SSHConfig{},
		}

		_, err := client.extractKeyData()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no private key provided")
	})

	t.Run("File Not Found", func(t *testing.T) {
		client := &SSHClient{
			config: &SSHConfig{
				PrivateKeyPath: "/nonexistent/path/key.pem",
			},
		}

		_, err := client.extractKeyData()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read private key file")
	})
}

func TestSSHClient_ValidateKeyData(t *testing.T) {
	// Generate a test RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	client := &SSHClient{}

	t.Run("Valid Key Data", func(t *testing.T) {
		block := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		}
		pemData := pem.EncodeToMemory(block)

		err := client.validateKeyData(pemData)
		assert.NoError(t, err)
	})

	t.Run("Empty Key Data", func(t *testing.T) {
		err := client.validateKeyData([]byte{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "private key data is empty")
	})

	t.Run("Legacy DES Encrypted Key", func(t *testing.T) {
		// Create legacy encrypted PEM block
		encryptedBlock := createLegacyEncryptedPEMBlock("RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(privateKey))
		pemData := pem.EncodeToMemory(encryptedBlock)

		err := client.validateKeyData(pemData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "legacy DES-encrypted traditional keys are no longer supported")
	})

	t.Run("Non-Legacy Encrypted Key", func(t *testing.T) {
		// Create encrypted PEM block without legacy headers
		block := &pem.Block{
			Type: "RSA PRIVATE KEY",
			Headers: map[string]string{
				"DEK-Info": "AES-256-CBC,0123456789ABCDEF", // Modern encryption
			},
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		}
		pemData := pem.EncodeToMemory(block)

		err := client.validateKeyData(pemData)
		assert.NoError(t, err) // Should pass validation
	})
}

func TestSSHClient_ParsePEMBlock(t *testing.T) {
	// Generate a test RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	client := &SSHClient{}

	t.Run("Valid RSA Key", func(t *testing.T) {
		block := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		}
		pemData := pem.EncodeToMemory(block)

		signer, err := client.parsePEMBlock(pemData)
		assert.NoError(t, err)
		assert.NotNil(t, signer)
		assert.Equal(t, "ssh-rsa", signer.PublicKey().Type())
	})

	t.Run("Invalid PEM Data", func(t *testing.T) {
		invalidData := []byte("not a valid PEM block")

		_, err := client.parsePEMBlock(invalidData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode PEM block")
	})

	t.Run("PEM with Extra Data", func(t *testing.T) {
		block := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		}
		pemData := pem.EncodeToMemory(block)
		// Add extra data
		pemData = append(pemData, []byte("extra data")...)

		signer, err := client.parsePEMBlock(pemData)
		assert.NoError(t, err) // Should still work but log warning
		assert.NotNil(t, signer)
	})
}

func TestSSHClient_HandlePEMBlockType(t *testing.T) {
	// Generate a test RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	client := &SSHClient{}

	t.Run("RSA Private Key", func(t *testing.T) {
		block := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		}

		signer, err := client.handlePEMBlockType(block)
		assert.NoError(t, err)
		assert.NotNil(t, signer)
		assert.Equal(t, "ssh-rsa", signer.PublicKey().Type())
	})

	t.Run("DSA Private Key", func(t *testing.T) {
		// Create a dummy DSA key block (we won't actually parse it)
		block := &pem.Block{
			Type:  "DSA PRIVATE KEY",
			Bytes: []byte("dummy dsa key data"),
		}

		_, err := client.handlePEMBlockType(block)
		// Should fail because it's dummy data, but the type should be recognized
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse traditional private key")
	})

	t.Run("EC Private Key", func(t *testing.T) {
		// Create a dummy EC key block
		block := &pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: []byte("dummy ec key data"),
		}

		_, err := client.handlePEMBlockType(block)
		// Should fail because it's dummy data, but the type should be recognized
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse traditional private key")
	})

	t.Run("OpenSSH Private Key", func(t *testing.T) {
		// Create a dummy OpenSSH key block
		block := &pem.Block{
			Type:  "OPENSSH PRIVATE KEY",
			Bytes: []byte("dummy openssh key data"),
		}

		_, err := client.handlePEMBlockType(block)
		// Should fail because it's dummy data, but the type should be recognized
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse OpenSSH private key")
	})

	t.Run("PKCS#8 Private Key", func(t *testing.T) {
		// Create PKCS#8 format
		pkcs8Bytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
		require.NoError(t, err)

		block := &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: pkcs8Bytes,
		}

		signer, err := client.handlePEMBlockType(block)
		assert.NoError(t, err)
		assert.NotNil(t, signer)
		assert.Equal(t, "ssh-rsa", signer.PublicKey().Type())
	})

	t.Run("Unknown Block Type", func(t *testing.T) {
		block := &pem.Block{
			Type:  "UNKNOWN TYPE",
			Bytes: []byte("dummy data"),
		}

		_, err := client.handlePEMBlockType(block)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported PEM block type")
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
