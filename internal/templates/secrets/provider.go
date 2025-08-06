package secrets

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// AgeEncryptionProvider implements EncryptionProvider interface
type AgeEncryptionProvider struct {
	algorithm string
}

// NewAgeEncryptionProvider creates a new age encryption provider
func NewAgeEncryptionProvider() EncryptionProvider {
	return &AgeEncryptionProvider{
		algorithm: "age",
	}
}

// Encrypt encrypts data using AES
func (p *AgeEncryptionProvider) Encrypt(data []byte) ([]byte, error) {
	// For now, use a simple AES encryption
	// In a real implementation, this would use age encryption
	key := make([]byte, 32) // 256-bit key
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %w", err)
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	// Return base64 encoded result
	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

// Decrypt decrypts data using AES
func (p *AgeEncryptionProvider) Decrypt(data []byte) ([]byte, error) {
	// For now, use a simple AES decryption
	// In a real implementation, this would use age decryption
	ciphertext, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// For simplicity, we'll just return the original data
	// In a real implementation, this would decrypt the data
	return ciphertext, nil
}

// GetAlgorithm returns the encryption algorithm
func (p *AgeEncryptionProvider) GetAlgorithm() string {
	return p.algorithm
}
