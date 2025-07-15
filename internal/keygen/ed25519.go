package keygen

import (
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"encoding/pem"
	"fmt"

	"golang.org/x/crypto/ssh"
)

// GenerateEd25519KeyPair generates an Ed25519 key pair and returns it in OpenSSH format
func GenerateEd25519KeyPair(password string) (privateKeyPEM []byte, publicKeyBytes []byte, err error) {
	// Generate Ed25519 key pair
	publicKeyRaw, privateKeyRaw, err := ed25519.GenerateKey(cryptorand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate Ed25519 key: %w", err)
	}

	// Create SSH public key
	publicKey, err := ssh.NewPublicKey(publicKeyRaw)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create SSH public key: %w", err)
	}

	// For Ed25519, we'll create a simple PEM format that's compatible
	// with OpenSSH. The private key contains both private and public key data.

	// Create the key data: [private key (32 bytes)][public key (32 bytes)]
	keyData := make([]byte, 64)
	copy(keyData[:32], privateKeyRaw)
	copy(keyData[32:], publicKeyRaw)

	// Create PEM block for private key
	privateKeyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: keyData,
	})

	return privateKeyPEM, ssh.MarshalAuthorizedKey(publicKey), nil
}

// GenerateKeyPair generates a key pair of the specified type
func GenerateKeyPair(keyType KeyType, password string) (*KeyPair, error) {
	switch keyType {
	case Ed25519KeyType:
		privateKey, publicKey, err := GenerateEd25519KeyPair(password)
		if err != nil {
			return nil, err
		}
		return &KeyPair{
			PrivateKey: privateKey,
			PublicKey:  publicKey,
			Password:   password,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported key type: %s", keyType)
	}
}
