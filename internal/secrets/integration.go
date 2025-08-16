// Package secrets provides age encryption functionality for the spooky codebase.
package secrets

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"filippo.io/age"
	"filippo.io/age/armor"

	spookyinterfaces "spooky/internal/interfaces"
	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
)

// stringWriter implements io.WriteCloser for strings.Builder
type stringWriter struct {
	builder *strings.Builder
}

func (w *stringWriter) Write(p []byte) (n int, err error) {
	return w.builder.Write(p)
}

func (w *stringWriter) Close() error {
	return nil
}

// Integration implements the SecretsIntegration interface with age encryption
type Integration struct {
	logger       spookytypeslogging.Logger
	config       *spookytypes.AgeConfig
	hclProcessor *HCLProcessor
}

// NewIntegration creates a new age-focused secrets integration
func NewIntegration(logger spookytypeslogging.Logger, config *spookytypes.AgeConfig) spookyinterfaces.SecretsIntegration {
	return &Integration{
		logger:       logger,
		config:       config,
		hclProcessor: NewHCLProcessor(logger),
	}
}

// Encrypt encrypts data with the given key (legacy method for interface compatibility)
func (i *Integration) Encrypt(_ context.Context, _, _ []byte) ([]byte, error) {
	// This is a legacy method - we should use age encryption instead
	return nil, fmt.Errorf("legacy AES-GCM encryption not supported - use age encryption methods")
}

// Decrypt decrypts data with the given key (legacy method for interface compatibility)
func (i *Integration) Decrypt(_ context.Context, _, _ []byte) ([]byte, error) {
	// This is a legacy method - we should use age decryption instead
	return nil, fmt.Errorf("legacy AES-GCM decryption not supported - use age decryption methods")
}

// ValidateKey validates an encryption key (legacy method for interface compatibility)
func (i *Integration) ValidateKey(_ context.Context, _ []byte) error {
	// This is a legacy method - we should use age key validation instead
	return fmt.Errorf("legacy AES-GCM key validation not supported - use age key validation methods")
}

// EncryptWithAge encrypts data with age using recipient public keys
func (i *Integration) EncryptWithAge(_ context.Context, data []byte, recipients []string) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	if len(recipients) == 0 {
		return nil, fmt.Errorf("at least one recipient is required")
	}

	// Parse recipients
	var ageRecipients []age.Recipient
	for _, recipient := range recipients {
		parsed, err := age.ParseX25519Recipient(recipient)
		if err != nil {
			return nil, fmt.Errorf("invalid recipient %s: %w", recipient, err)
		}
		ageRecipients = append(ageRecipients, parsed)
	}

	// Create output buffer
	var output strings.Builder
	var writer io.WriteCloser

	// Use armored output if configured
	if i.config != nil && i.config.Encryption.Armor {
		writer = armor.NewWriter(&stringWriter{builder: &output})
	} else {
		writer = &stringWriter{builder: &output}
	}

	// Create age encryptor
	encryptor, err := age.Encrypt(writer, ageRecipients...)
	if err != nil {
		return nil, fmt.Errorf("failed to create age encryptor: %w", err)
	}

	// Write data
	if _, err := encryptor.Write(data); err != nil {
		return nil, fmt.Errorf("failed to write data to encryptor: %w", err)
	}

	// Close encryptor
	if err := encryptor.Close(); err != nil {
		return nil, fmt.Errorf("failed to close encryptor: %w", err)
	}

	// Close writer if armored
	if i.config != nil && i.config.Encryption.Armor {
		if err := writer.Close(); err != nil {
			return nil, fmt.Errorf("failed to close armored writer: %w", err)
		}
	}

	i.logger.Info("Data encrypted with age successfully", map[string]interface{}{
		"data_size":       len(data),
		"ciphertext_size": len(output.String()),
		"recipients":      len(recipients),
		"armored":         i.config != nil && i.config.Encryption.Armor,
	})

	return []byte(output.String()), nil
}

// DecryptWithAge decrypts age-encrypted data using identity file
func (i *Integration) DecryptWithAge(_ context.Context, data []byte, identityPath string) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	if identityPath == "" {
		return nil, fmt.Errorf("identity path is required")
	}

	// Read identity file
	identityFile, err := os.Open(identityPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open identity file %s: %w", identityPath, err)
	}
	defer identityFile.Close()

	identities, err := age.ParseIdentities(identityFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse identity file %s: %w", identityPath, err)
	}

	// Create input reader
	var reader io.Reader = strings.NewReader(string(data))

	// Check if data is armored
	if strings.HasPrefix(string(data), "-----BEGIN AGE ENCRYPTED FILE-----") {
		reader = armor.NewReader(reader)
	}

	// Create age decryptor
	decryptor, err := age.Decrypt(reader, identities...)
	if err != nil {
		return nil, fmt.Errorf("failed to create age decryptor: %w", err)
	}

	// Read decrypted data
	plaintext, err := io.ReadAll(decryptor)
	if err != nil {
		return nil, fmt.Errorf("failed to read decrypted data: %w", err)
	}

	i.logger.Info("Data decrypted with age successfully", map[string]interface{}{
		"ciphertext_size": len(data),
		"plaintext_size":  len(plaintext),
		"identity_path":   identityPath,
	})

	return plaintext, nil
}

// EncryptWithPassphrase encrypts data with age using a passphrase
func (i *Integration) EncryptWithPassphrase(_ context.Context, data []byte, passphrase string) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	if passphrase == "" {
		return nil, fmt.Errorf("passphrase cannot be empty")
	}

	// Create passphrase recipient
	recipient, err := age.NewScryptRecipient(passphrase)
	if err != nil {
		return nil, fmt.Errorf("failed to create passphrase recipient: %w", err)
	}

	// Create output buffer
	var output strings.Builder
	var writer io.WriteCloser

	// Use armored output if configured
	if i.config != nil && i.config.Encryption.Armor {
		writer = armor.NewWriter(&stringWriter{builder: &output})
	} else {
		writer = &stringWriter{builder: &output}
	}

	// Create age encryptor
	encryptor, err := age.Encrypt(writer, recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to create age encryptor: %w", err)
	}

	// Write data
	if _, err := encryptor.Write(data); err != nil {
		return nil, fmt.Errorf("failed to write data to encryptor: %w", err)
	}

	// Close encryptor
	if err := encryptor.Close(); err != nil {
		return nil, fmt.Errorf("failed to close encryptor: %w", err)
	}

	// Close writer if armored
	if i.config != nil && i.config.Encryption.Armor {
		if err := writer.Close(); err != nil {
			return nil, fmt.Errorf("failed to close armored writer: %w", err)
		}
	}

	i.logger.Info("Data encrypted with passphrase successfully", map[string]interface{}{
		"data_size":       len(data),
		"ciphertext_size": len(output.String()),
		"armored":         i.config != nil && i.config.Encryption.Armor,
	})

	return []byte(output.String()), nil
}

// DecryptWithPassphrase decrypts age-encrypted data using a passphrase
func (i *Integration) DecryptWithPassphrase(_ context.Context, data []byte, passphrase string) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	if passphrase == "" {
		return nil, fmt.Errorf("passphrase cannot be empty")
	}

	// Create passphrase identity
	identity, err := age.NewScryptIdentity(passphrase)
	if err != nil {
		return nil, fmt.Errorf("failed to create passphrase identity: %w", err)
	}

	// Create input reader
	var reader io.Reader = strings.NewReader(string(data))

	// Check if data is armored
	if strings.HasPrefix(string(data), "-----BEGIN AGE ENCRYPTED FILE-----") {
		reader = armor.NewReader(reader)
	}

	// Create age decryptor
	decryptor, err := age.Decrypt(reader, identity)
	if err != nil {
		return nil, fmt.Errorf("failed to create age decryptor: %w", err)
	}

	// Read decrypted data
	plaintext, err := io.ReadAll(decryptor)
	if err != nil {
		return nil, fmt.Errorf("failed to read decrypted data: %w", err)
	}

	i.logger.Info("Data decrypted with passphrase successfully", map[string]interface{}{
		"ciphertext_size": len(data),
		"plaintext_size":  len(plaintext),
	})

	return plaintext, nil
}

// ValidateAgeKey validates an age identity file
func (i *Integration) ValidateAgeKey(_ context.Context, keyPath string) error {
	if keyPath == "" {
		return fmt.Errorf("key path cannot be empty")
	}

	// Check if file exists
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return fmt.Errorf("key file does not exist: %s", keyPath)
	}

	// Try to parse identities
	identityFile, err := os.Open(keyPath)
	if err != nil {
		return fmt.Errorf("failed to open identity file: %w", err)
	}
	defer identityFile.Close()

	identities, err := age.ParseIdentities(identityFile)
	if err != nil {
		return fmt.Errorf("failed to parse age identity file: %w", err)
	}

	if len(identities) == 0 {
		return fmt.Errorf("no valid identities found in key file")
	}

	i.logger.Info("Age key validated successfully", map[string]interface{}{
		"key_path":   keyPath,
		"identities": len(identities),
	})

	return nil
}

// ListRecipients extracts recipient information from age-encrypted data
func (i *Integration) ListRecipients(_ context.Context, encryptedData []byte) ([]string, error) {
	if len(encryptedData) == 0 {
		return nil, fmt.Errorf("encrypted data cannot be empty")
	}

	// Create input reader
	var reader io.Reader = strings.NewReader(string(encryptedData))

	// Check if data is armored
	if strings.HasPrefix(string(encryptedData), "-----BEGIN AGE ENCRYPTED FILE-----") {
		reader = armor.NewReader(reader)
	}

	// Parse recipients
	recipients, err := age.ParseRecipients(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse recipients: %w", err)
	}

	var recipientStrings []string
	for _, recipient := range recipients {
		if x25519, ok := recipient.(*age.X25519Recipient); ok {
			recipientStrings = append(recipientStrings, x25519.String())
		}
	}

	return recipientStrings, nil
}

// ValidateAgeEncryptedValue validates age-encrypted value using the age library
func (i *Integration) ValidateAgeEncryptedValue(_ context.Context, value string) error {
	if value == "" {
		return fmt.Errorf("value cannot be empty")
	}

	// Let the age library handle validation
	_, err := age.ParseRecipients(strings.NewReader(value))
	if err != nil {
		return fmt.Errorf("invalid age-encrypted value: %w", err)
	}

	return nil
}

// LoadRecipients loads recipients from a file
func (i *Integration) LoadRecipients(_ context.Context, recipientsPath string) ([]string, error) {
	if recipientsPath == "" {
		return nil, fmt.Errorf("recipients path cannot be empty")
	}

	// Read recipients file
	data, err := os.ReadFile(recipientsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read recipients file: %w", err)
	}

	// Parse recipients
	recipients, err := age.ParseRecipients(strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse recipients: %w", err)
	}

	var recipientStrings []string
	for _, recipient := range recipients {
		if x25519, ok := recipient.(*age.X25519Recipient); ok {
			recipientStrings = append(recipientStrings, x25519.String())
		}
	}

	return recipientStrings, nil
}

// LoadIdentities loads identities from a directory
func (i *Integration) LoadIdentities(_ context.Context, identitiesPath string) ([]string, error) {
	if identitiesPath == "" {
		return nil, fmt.Errorf("identities path cannot be empty")
	}

	// Check if directory exists
	if _, err := os.Stat(identitiesPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("identities directory does not exist: %s", identitiesPath)
	}

	var identityFiles []string

	// Walk directory looking for identity files
	err := filepath.Walk(identitiesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip hidden files
		if strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}

		// Try to parse as identity file
		identityFile, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer identityFile.Close()

		if _, err := age.ParseIdentities(identityFile); err == nil {
			identityFiles = append(identityFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk identities directory: %w", err)
	}

	return identityFiles, nil
}

// EncryptHCLValues encrypts values in any HCL-compatible structure
func (i *Integration) EncryptHCLValues(ctx context.Context, data interface{}, recipients []string, dryRun bool) error {
	return i.hclProcessor.EncryptHCLValues(ctx, data, i, recipients, dryRun)
}

// EncryptHCLValuesSensitive encrypts values in any HCL-compatible structure with sensitivity awareness
// This is now a simple wrapper around EncryptHCLValues since we encrypt entire objects
func (i *Integration) EncryptHCLValuesSensitive(ctx context.Context, data interface{}, recipients []string, dryRun bool, _ func(path []string, value interface{}) bool) error {
	// Since we now encrypt entire objects, we ignore the sensitivity function
	return i.hclProcessor.EncryptHCLValues(ctx, data, i, recipients, dryRun)
}

// EncryptHCLValuesWithJSONSupport encrypts values in any HCL-compatible structure with JSON serialization
// This is now a simple wrapper around EncryptHCLValues since we always use JSON for objects
func (i *Integration) EncryptHCLValuesWithJSONSupport(ctx context.Context, data interface{}, recipients []string, dryRun bool) error {
	return i.hclProcessor.EncryptHCLValues(ctx, data, i, recipients, dryRun)
}

// DecryptHCLValues decrypts age-encrypted values in any HCL-compatible structure
func (i *Integration) DecryptHCLValues(ctx context.Context, data interface{}, identityPath string) error {
	return i.hclProcessor.DecryptHCLValues(ctx, data, i, identityPath)
}

// DecryptHCLValuesWithJSONSupport decrypts age-encrypted values with JSON deserialization support
// This is now a simple wrapper around DecryptHCLValues since we always support JSON
func (i *Integration) DecryptHCLValuesWithJSONSupport(ctx context.Context, data interface{}, identityPath string) error {
	return i.hclProcessor.DecryptHCLValues(ctx, data, i, identityPath)
}
