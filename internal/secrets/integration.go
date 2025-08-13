// Package secrets provides secrets management functionality for the spooky codebase.
package secrets

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypeslogging "spooky/internal/types/logging"
)

// Integration implements the SecretsIntegration interface
type Integration struct {
	logger spookytypeslogging.Logger
}

// NewIntegration creates a new secrets integration
func NewIntegration(logger spookytypeslogging.Logger) spookyinterfaces.SecretsIntegration {
	return &Integration{
		logger: logger,
	}
}

// Encrypt encrypts data with the given key
func (i *Integration) Encrypt(ctx context.Context, data []byte, key []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	if len(key) == 0 {
		return nil, fmt.Errorf("key cannot be empty")
	}

	// Create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM mode: %w", err)
	}

	// Create nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt data
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	i.logger.Info("Data encrypted successfully", map[string]interface{}{
		"data_size":       len(data),
		"ciphertext_size": len(ciphertext),
	})

	return ciphertext, nil
}

// Decrypt decrypts data with the given key
func (i *Integration) Decrypt(ctx context.Context, data []byte, key []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	if len(key) == 0 {
		return nil, fmt.Errorf("key cannot be empty")
	}

	// Create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM mode: %w", err)
	}

	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Decrypt data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	i.logger.Info("Data decrypted successfully", map[string]interface{}{
		"ciphertext_size": len(data),
		"plaintext_size":  len(plaintext),
	})

	return plaintext, nil
}

// ValidateKey validates an encryption key
func (i *Integration) ValidateKey(ctx context.Context, key []byte) error {
	if len(key) == 0 {
		return fmt.Errorf("key cannot be empty")
	}

	// Check if key is valid for AES-256
	if len(key) != 32 {
		return fmt.Errorf("key must be 32 bytes for AES-256, got %d bytes", len(key))
	}

	// Test if key can be used to create a cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("invalid key: %w", err)
	}

	// Test if GCM mode can be created
	_, err = cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("key cannot be used with GCM mode: %w", err)
	}

	i.logger.Info("Encryption key validated successfully", map[string]interface{}{
		"key_size": len(key),
	})

	return nil
}
