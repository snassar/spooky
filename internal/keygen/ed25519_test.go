package keygen

import (
	"encoding/pem"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestGenerateEd25519KeyPair(t *testing.T) {
	// Test successful key generation
	privateKeyPEM, publicKeyBytes, err := GenerateEd25519KeyPair("testpassword")
	if err != nil {
		t.Fatalf("failed to generate key pair: %v", err)
	}

	// Test private key PEM format
	if len(privateKeyPEM) == 0 {
		t.Error("private key PEM should not be empty")
	}

	// Verify PEM block structure
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		t.Error("failed to decode private key PEM")
	}
	if block.Type != "OPENSSH PRIVATE KEY" {
		t.Errorf("expected PEM type 'OPENSSH PRIVATE KEY', got '%s'", block.Type)
	}

	// Test public key format
	if len(publicKeyBytes) == 0 {
		t.Error("public key should not be empty")
	}

	// Verify public key starts with ssh-ed25519
	publicKeyStr := string(publicKeyBytes)
	if !strings.HasPrefix(publicKeyStr, "ssh-ed25519 ") {
		t.Errorf("public key should start with 'ssh-ed25519 ', got: %s", publicKeyStr[:20])
	}

	// Test that public key can be parsed
	_, _, _, _, err = ssh.ParseAuthorizedKey(publicKeyBytes)
	if err != nil {
		t.Errorf("failed to parse generated public key: %v", err)
	}

	// Test key data structure (should be 64 bytes: 32 private + 32 public)
	if len(block.Bytes) != 64 {
		t.Errorf("expected key data length 64, got %d", len(block.Bytes))
	}
}

func TestGenerateEd25519KeyPair_Consistency(t *testing.T) {
	// Test that multiple key pairs are different
	privateKey1, publicKey1, err := GenerateEd25519KeyPair("password1")
	if err != nil {
		t.Fatalf("failed to generate first key pair: %v", err)
	}

	privateKey2, publicKey2, err := GenerateEd25519KeyPair("password2")
	if err != nil {
		t.Fatalf("failed to generate second key pair: %v", err)
	}

	// Keys should be different
	if string(privateKey1) == string(privateKey2) {
		t.Error("generated private keys should be different")
	}
	if string(publicKey1) == string(publicKey2) {
		t.Error("generated public keys should be different")
	}
}

func TestGenerateKeyPair(t *testing.T) {
	testCases := []struct {
		name     string
		keyType  KeyType
		password string
		wantErr  bool
	}{
		{"Ed25519 key", Ed25519KeyType, "testpassword", false},
		{"invalid key type", KeyType("invalid"), "testpassword", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			keyPair, err := GenerateKeyPair(tc.keyType, tc.password)

			if tc.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if keyPair == nil {
				t.Error("keyPair should not be nil")
				return
			}

			// Note: KeyPair doesn't store the key type, it's only used during generation

			if len(keyPair.PrivateKey) == 0 {
				t.Error("private key should not be empty")
			}

			if len(keyPair.PublicKey) == 0 {
				t.Error("public key should not be empty")
			}

			if keyPair.Password != tc.password {
				t.Errorf("expected password %s, got %s", tc.password, keyPair.Password)
			}
		})
	}
}

func TestGenerateKeyPair_Ed25519Validation(t *testing.T) {
	keyPair, err := GenerateKeyPair(Ed25519KeyType, "testpassword")
	if err != nil {
		t.Fatalf("failed to generate Ed25519 key pair: %v", err)
	}

	// Verify PEM block structure
	block, _ := pem.Decode(keyPair.PrivateKey)
	if block == nil {
		t.Error("failed to decode private key PEM")
	}
	if block.Type != "OPENSSH PRIVATE KEY" {
		t.Errorf("expected PEM type 'OPENSSH PRIVATE KEY', got '%s'", block.Type)
	}

	// Verify public key format
	publicKeyStr := string(keyPair.PublicKey)
	if !strings.HasPrefix(publicKeyStr, "ssh-ed25519 ") {
		t.Errorf("public key should start with 'ssh-ed25519 ', got: %s", publicKeyStr[:20])
	}

	// Test that public key can be parsed
	_, _, _, _, err = ssh.ParseAuthorizedKey(keyPair.PublicKey)
	if err != nil {
		t.Errorf("failed to parse generated public key: %v", err)
	}
}
