package ssh

import (
	"crypto/rsa"
	"fmt"
	"io"
	"testing"

	"golang.org/x/crypto/ssh"

	spookylogging "spooky/internal/logging"
)

func TestKeyValidation(t *testing.T) {
	// Create a log manager
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("ssh-test")

	// Create SSH client
	client := NewSimpleClient(nil, logger)

	// Test ed25519 key generation and validation
	t.Run("ED25519 Key", func(t *testing.T) {
		signer, err := client.generateED25519Key()
		if err != nil {
			t.Fatalf("Failed to generate ed25519 key: %v", err)
		}

		// Validate the generated key
		if err := client.validateKeyType(signer); err != nil {
			t.Fatalf("ED25519 key validation failed: %v", err)
		}

		pubKey := signer.PublicKey()
		if pubKey.Type() != ssh.KeyAlgoED25519 {
			t.Errorf("Expected ed25519 key type, got %s", pubKey.Type())
		}

		t.Logf("✓ ED25519 key validation passed, fingerprint: %s",
			ssh.FingerprintSHA256(pubKey))
	})

	// Test RSA 4096-bit key generation and validation
	t.Run("RSA 4096-bit Key", func(t *testing.T) {
		signer, err := client.generateRSA4096Key()
		if err != nil {
			t.Fatalf("Failed to generate RSA key: %v", err)
		}

		// Validate the generated key
		if err := client.validateKeyType(signer); err != nil {
			t.Fatalf("RSA key validation failed: %v", err)
		}

		pubKey := signer.PublicKey()
		if pubKey.Type() != ssh.KeyAlgoRSA {
			t.Errorf("Expected RSA key type, got %s", pubKey.Type())
		}

		// Check RSA key size
		rsaPubKey, ok := pubKey.(ssh.CryptoPublicKey)
		if !ok {
			t.Fatalf("Failed to extract RSA public key")
		}

		cryptoPubKey := rsaPubKey.CryptoPublicKey()
		rsaKey, ok := cryptoPubKey.(*rsa.PublicKey)
		if !ok {
			t.Fatalf("Failed to cast to RSA public key")
		}

		if rsaKey.Size()*8 < MinRSAKeySize {
			t.Errorf("RSA key size %d bits is less than minimum required %d bits",
				rsaKey.Size()*8, MinRSAKeySize)
		}

		t.Logf("✓ RSA 4096-bit key validation passed, fingerprint: %s",
			ssh.FingerprintSHA256(pubKey))
	})

	// Test unsupported key type validation
	t.Run("Unsupported Key Type", func(t *testing.T) {
		// Create a mock signer with unsupported key type
		mockSigner := &mockSigner{
			keyType: "unsupported-key-type",
		}

		err := client.validateKeyType(mockSigner)
		if err == nil {
			t.Error("Expected validation error for unsupported key type")
		}

		if keyErr, ok := err.(*KeyValidationError); ok {
			if keyErr.KeyType != "unsupported-key-type" {
				t.Errorf("Expected key type 'unsupported-key-type', got '%s'", keyErr.KeyType)
			}
		} else {
			t.Error("Expected KeyValidationError type")
		}

		t.Logf("✓ Unsupported key type validation correctly rejected")
	})
}

func TestSupportedKeyGeneration(t *testing.T) {
	// Create a log manager
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("ssh-test")

	// Create SSH client
	client := NewSimpleClient(nil, logger)

	// Test supported key types
	supportedTypes := []string{KeyTypeED25519, KeyTypeRSA4096}

	for _, keyType := range supportedTypes {
		t.Run(keyType, func(t *testing.T) {
			signer, err := client.generateSupportedKey(keyType)
			if err != nil {
				t.Fatalf("Failed to generate %s key: %v", keyType, err)
			}

			if signer == nil {
				t.Fatalf("Generated signer is nil for %s", keyType)
			}

			pubKey := signer.PublicKey()
			if pubKey == nil {
				t.Fatalf("Generated public key is nil for %s", keyType)
			}

			t.Logf("✓ Generated %s key successfully, fingerprint: %s",
				keyType, ssh.FingerprintSHA256(pubKey))
		})
	}

	// Test unsupported key type
	t.Run("Unsupported Type", func(t *testing.T) {
		_, err := client.generateSupportedKey("unsupported-type")
		if err == nil {
			t.Error("Expected error for unsupported key type")
		}

		t.Logf("✓ Unsupported key type generation correctly rejected: %v", err)
	})
}

// mockSigner is a mock implementation for testing unsupported key types
type mockSigner struct {
	keyType string
}

func (m *mockSigner) PublicKey() ssh.PublicKey {
	return &mockPublicKey{keyType: m.keyType}
}

func (m *mockSigner) Sign(rand io.Reader, data []byte) (*ssh.Signature, error) {
	return nil, fmt.Errorf("mock signer not implemented")
}

// mockPublicKey is a mock implementation for testing unsupported key types
type mockPublicKey struct {
	keyType string
}

func (m *mockPublicKey) Type() string {
	return m.keyType
}

func (m *mockPublicKey) Marshal() []byte {
	return []byte("mock-public-key")
}

func (m *mockPublicKey) Verify(data []byte, sig *ssh.Signature) error {
	return fmt.Errorf("mock public key not implemented")
}
