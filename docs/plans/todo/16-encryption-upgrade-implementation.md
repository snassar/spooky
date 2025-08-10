# Implementation Plan: Encryption Upgrade Implementation

## Overview
Upgrade encryption from simple AES to age encryption for template secrets, replacing placeholder implementations with robust age-based encryption.

## Task Details
- **Task ID**: 7.7
- **Priority**: Medium
- **Files**: 
  - `internal/templates/secrets/provider.go`
- **Functions**: Age encryption, key management, secure storage

## Current State Analysis

### Existing Patterns
1. **Template Secrets**: Basic secret management exists
2. **Encryption**: Simple AES encryption implemented
3. **Key Management**: Basic key management available
4. **Error Handling**: Consistent error wrapping

### Current Placeholder Code
```go
// internal/templates/secrets/provider.go line 45
// For now, use a simple AES encryption/decryption
```

## Implementation Requirements

### Interface Compliance
The encryption upgrade must:
1. **Replace AES with age encryption** for all secret operations
2. **Implement age key management** and key rotation
3. **Support multiple recipients** and public key encryption
4. **Provide secure key storage** and retrieval
5. **Handle key generation** and validation
6. **Support encryption/decryption** with age
7. **Implement secure random generation** for keys

### Required Dependencies
- Age encryption library
- Key management system
- Secure random generation
- File system security

## Detailed Implementation Plan

### Step 1: Implement Age Encryption Provider

#### 1.1 Age Provider Structure
```go
// internal/templates/secrets/age_provider.go
package secrets

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "io"
    "os"
    "path/filepath"
    
    "filippo.io/age"
    "filippo.io/age/armor"
    "spooky/internal/logging"
)

// AgeProvider implements age-based encryption for secrets
type AgeProvider struct {
    keyPath    string
    recipients []age.Recipient
    logger     logging.Logger
}

// NewAgeProvider creates a new age encryption provider
func NewAgeProvider(keyPath string, logger logging.Logger) (*AgeProvider, error) {
    provider := &AgeProvider{
        keyPath: keyPath,
        logger:  logger,
    }

    // Load or generate age keys
    if err := provider.loadOrGenerateKeys(); err != nil {
        return nil, fmt.Errorf("failed to initialize age provider: %w", err)
    }

    return provider, nil
}

// loadOrGenerateKeys loads existing keys or generates new ones
func (ap *AgeProvider) loadOrGenerateKeys() error {
    // Check if key file exists
    if _, err := os.Stat(ap.keyPath); os.IsNotExist(err) {
        // Generate new keys
        return ap.generateKeys()
    }

    // Load existing keys
    return ap.loadKeys()
}

// generateKeys generates new age keys
func (ap *AgeProvider) generateKeys() error {
    ap.logger.Info("Generating new age keys",
        logging.String("key_path", ap.keyPath))

    // Generate identity (private key)
    identity, err := age.GenerateX25519Identity()
    if err != nil {
        return fmt.Errorf("failed to generate age identity: %w", err)
    }

    // Ensure key directory exists
    keyDir := filepath.Dir(ap.keyPath)
    if err := os.MkdirAll(keyDir, 0700); err != nil {
        return fmt.Errorf("failed to create key directory: %w", err)
    }

    // Write private key
    privateKeyFile, err := os.OpenFile(ap.keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
        return fmt.Errorf("failed to create private key file: %w", err)
    }
    defer privateKeyFile.Close()

    if err := identity.Marshal(privateKeyFile); err != nil {
        return fmt.Errorf("failed to write private key: %w", err)
    }

    // Write public key
    publicKeyPath := ap.keyPath + ".pub"
    publicKeyFile, err := os.OpenFile(publicKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        return fmt.Errorf("failed to create public key file: %w", err)
    }
    defer publicKeyFile.Close()

    if _, err := publicKeyFile.WriteString(identity.Recipient().String()); err != nil {
        return fmt.Errorf("failed to write public key: %w", err)
    }

    // Set up recipients
    ap.recipients = []age.Recipient{identity.Recipient()}

    ap.logger.Info("Age keys generated successfully",
        logging.String("private_key", ap.keyPath),
        logging.String("public_key", publicKeyPath))

    return nil
}

// loadKeys loads existing age keys
func (ap *AgeProvider) loadKeys() error {
    ap.logger.Info("Loading existing age keys",
        logging.String("key_path", ap.keyPath))

    // Read private key file
    privateKeyFile, err := os.Open(ap.keyPath)
    if err != nil {
        return fmt.Errorf("failed to open private key file: %w", err)
    }
    defer privateKeyFile.Close()

    // Parse identity
    identities, err := age.ParseIdentities(privateKeyFile)
    if err != nil {
        return fmt.Errorf("failed to parse private key: %w", err)
    }

    if len(identities) == 0 {
        return fmt.Errorf("no identities found in private key file")
    }

    // Set up recipients from identities
    ap.recipients = make([]age.Recipient, len(identities))
    for i, identity := range identities {
        ap.recipients[i] = identity.Recipient()
    }

    ap.logger.Info("Age keys loaded successfully",
        logging.Int("recipient_count", len(ap.recipients)))

    return nil
}

// AddRecipient adds a recipient to the encryption
func (ap *AgeProvider) AddRecipient(recipientString string) error {
    recipient, err := age.ParseX25519Recipient(recipientString)
    if err != nil {
        return fmt.Errorf("failed to parse recipient: %w", err)
    }

    ap.recipients = append(ap.recipients, recipient)
    ap.logger.Debug("Recipient added",
        logging.String("recipient", recipientString))

    return nil
}

// RemoveRecipient removes a recipient from encryption
func (ap *AgeProvider) RemoveRecipient(recipientString string) error {
    for i, recipient := range ap.recipients {
        if recipient.String() == recipientString {
            ap.recipients = append(ap.recipients[:i], ap.recipients[i+1:]...)
            ap.logger.Debug("Recipient removed",
                logging.String("recipient", recipientString))
            return nil
        }
    }

    return fmt.Errorf("recipient not found: %s", recipientString)
}

// GetPublicKey returns the public key for this provider
func (ap *AgeProvider) GetPublicKey() (string, error) {
    if len(ap.recipients) == 0 {
        return "", fmt.Errorf("no recipients available")
    }

    return ap.recipients[0].String(), nil
}
```

#### 1.2 Age Encryption Methods
```go
// Encrypt encrypts data using age encryption
func (ap *AgeProvider) Encrypt(data []byte) ([]byte, error) {
    if len(ap.recipients) == 0 {
        return nil, fmt.Errorf("no recipients configured for encryption")
    }

    ap.logger.Debug("Encrypting data with age",
        logging.Int("data_size", len(data)),
        logging.Int("recipient_count", len(ap.recipients)))

    // Create encrypted output
    encrypted := &bytes.Buffer{}

    // Create age writer with armor
    ageWriter, err := armor.NewWriter(encrypted, "AGE ENCRYPTED FILE")
    if err != nil {
        return nil, fmt.Errorf("failed to create age writer: %w", err)
    }
    defer ageWriter.Close()

    // Create encrypted writer
    encryptedWriter, err := age.Encrypt(ageWriter, ap.recipients...)
    if err != nil {
        return nil, fmt.Errorf("failed to create encrypted writer: %w", err)
    }
    defer encryptedWriter.Close()

    // Write data
    if _, err := encryptedWriter.Write(data); err != nil {
        return nil, fmt.Errorf("failed to write encrypted data: %w", err)
    }

    encryptedData := encrypted.Bytes()

    ap.logger.Debug("Data encrypted successfully",
        logging.Int("original_size", len(data)),
        logging.Int("encrypted_size", len(encryptedData)))

    return encryptedData, nil
}

// Decrypt decrypts data using age encryption
func (ap *AgeProvider) Decrypt(encryptedData []byte) ([]byte, error) {
    ap.logger.Debug("Decrypting data with age",
        logging.Int("encrypted_size", len(encryptedData)))

    // Load identities for decryption
    identities, err := ap.loadIdentities()
    if err != nil {
        return nil, fmt.Errorf("failed to load identities: %w", err)
    }

    // Create reader from encrypted data
    encryptedReader := bytes.NewReader(encryptedData)

    // Create age reader with armor
    ageReader, err := armor.NewReader(encryptedReader)
    if err != nil {
        return nil, fmt.Errorf("failed to create age reader: %w", err)
    }

    // Create decrypted reader
    decryptedReader, err := age.Decrypt(ageReader, identities...)
    if err != nil {
        return nil, fmt.Errorf("failed to create decrypted reader: %w", err)
    }

    // Read decrypted data
    decryptedData, err := io.ReadAll(decryptedReader)
    if err != nil {
        return nil, fmt.Errorf("failed to read decrypted data: %w", err)
    }

    ap.logger.Debug("Data decrypted successfully",
        logging.Int("encrypted_size", len(encryptedData)),
        logging.Int("decrypted_size", len(decryptedData)))

    return decryptedData, nil
}

// loadIdentities loads identities for decryption
func (ap *AgeProvider) loadIdentities() ([]age.Identity, error) {
    // Read private key file
    privateKeyFile, err := os.Open(ap.keyPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open private key file: %w", err)
    }
    defer privateKeyFile.Close()

    // Parse identities
    identities, err := age.ParseIdentities(privateKeyFile)
    if err != nil {
        return nil, fmt.Errorf("failed to parse identities: %w", err)
    }

    return identities, nil
}
```

### Step 2: Implement Key Management

#### 2.1 Key Manager Implementation
```go
// internal/templates/secrets/key_manager.go
package secrets

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "os"
    "path/filepath"
    "time"
    
    "filippo.io/age"
    "spooky/internal/logging"
)

// KeyManager manages age keys and key rotation
type KeyManager struct {
    keyDir     string
    logger     logging.Logger
    provider   *AgeProvider
}

// NewKeyManager creates a new key manager
func NewKeyManager(keyDir string, logger logging.Logger) (*KeyManager, error) {
    // Ensure key directory exists
    if err := os.MkdirAll(keyDir, 0700); err != nil {
        return nil, fmt.Errorf("failed to create key directory: %w", err)
    }

    keyPath := filepath.Join(keyDir, "age.key")
    provider, err := NewAgeProvider(keyPath, logger)
    if err != nil {
        return nil, fmt.Errorf("failed to create age provider: %w", err)
    }

    return &KeyManager{
        keyDir:   keyDir,
        logger:   logger,
        provider: provider,
    }, nil
}

// GenerateKeyPair generates a new key pair
func (km *KeyManager) GenerateKeyPair() (*KeyPair, error) {
    km.logger.Info("Generating new age key pair")

    // Generate identity
    identity, err := age.GenerateX25519Identity()
    if err != nil {
        return nil, fmt.Errorf("failed to generate identity: %w", err)
    }

    keyPair := &KeyPair{
        PrivateKey: identity.String(),
        PublicKey:  identity.Recipient().String(),
        CreatedAt:  time.Now(),
    }

    km.logger.Info("Key pair generated successfully",
        logging.String("public_key", keyPair.PublicKey))

    return keyPair, nil
}

// SaveKeyPair saves a key pair to disk
func (km *KeyManager) SaveKeyPair(keyPair *KeyPair, name string) error {
    km.logger.Info("Saving key pair",
        logging.String("name", name))

    // Save private key
    privateKeyPath := filepath.Join(km.keyDir, name+".key")
    privateKeyFile, err := os.OpenFile(privateKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
        return fmt.Errorf("failed to create private key file: %w", err)
    }
    defer privateKeyFile.Close()

    if _, err := privateKeyFile.WriteString(keyPair.PrivateKey); err != nil {
        return fmt.Errorf("failed to write private key: %w", err)
    }

    // Save public key
    publicKeyPath := filepath.Join(km.keyDir, name+".pub")
    publicKeyFile, err := os.OpenFile(publicKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        return fmt.Errorf("failed to create public key file: %w", err)
    }
    defer publicKeyFile.Close()

    if _, err := publicKeyFile.WriteString(keyPair.PublicKey); err != nil {
        return fmt.Errorf("failed to write public key: %w", err)
    }

    km.logger.Info("Key pair saved successfully",
        logging.String("name", name),
        logging.String("private_key", privateKeyPath),
        logging.String("public_key", publicKeyPath))

    return nil
}

// LoadKeyPair loads a key pair from disk
func (km *KeyManager) LoadKeyPair(name string) (*KeyPair, error) {
    km.logger.Debug("Loading key pair",
        logging.String("name", name))

    // Load private key
    privateKeyPath := filepath.Join(km.keyDir, name+".key")
    privateKeyData, err := os.ReadFile(privateKeyPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read private key: %w", err)
    }

    // Load public key
    publicKeyPath := filepath.Join(km.keyDir, name+".pub")
    publicKeyData, err := os.ReadFile(publicKeyPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read public key: %w", err)
    }

    keyPair := &KeyPair{
        PrivateKey: string(privateKeyData),
        PublicKey:  string(publicKeyData),
        CreatedAt:  time.Now(), // This would be extracted from file metadata
    }

    km.logger.Debug("Key pair loaded successfully",
        logging.String("name", name))

    return keyPair, nil
}

// RotateKeys rotates encryption keys
func (km *KeyManager) RotateKeys() error {
    km.logger.Info("Rotating encryption keys")

    // Generate new key pair
    newKeyPair, err := km.GenerateKeyPair()
    if err != nil {
        return fmt.Errorf("failed to generate new key pair: %w", err)
    }

    // Save new key pair
    timestamp := time.Now().Format("20060102-150405")
    newKeyName := fmt.Sprintf("age-%s", timestamp)
    if err := km.SaveKeyPair(newKeyPair, newKeyName); err != nil {
        return fmt.Errorf("failed to save new key pair: %w", err)
    }

    // Update provider with new key
    if err := km.provider.updateKey(newKeyPair.PrivateKey); err != nil {
        return fmt.Errorf("failed to update provider key: %w", err)
    }

    km.logger.Info("Key rotation completed successfully",
        logging.String("new_key", newKeyName))

    return nil
}

// GenerateRandomBytes generates random bytes
func (km *KeyManager) GenerateRandomBytes(length int) ([]byte, error) {
    bytes := make([]byte, length)
    if _, err := rand.Read(bytes); err != nil {
        return nil, fmt.Errorf("failed to generate random bytes: %w", err)
    }
    return bytes, nil
}

// GenerateRandomString generates a random string
func (km *KeyManager) GenerateRandomString(length int) (string, error) {
    bytes, err := km.GenerateRandomBytes(length)
    if err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// KeyPair represents an age key pair
type KeyPair struct {
    PrivateKey string    `json:"private_key"`
    PublicKey  string    `json:"public_key"`
    CreatedAt  time.Time `json:"created_at"`
}
```

### Step 3: Update Secret Provider

#### 3.1 Enhanced Secret Provider
```go
// internal/templates/secrets/provider.go
package secrets

import (
    "fmt"
    "sync"
    
    "spooky/internal/logging"
)

// Provider provides secret encryption and decryption
type Provider struct {
    ageProvider *AgeProvider
    keyManager  *KeyManager
    logger      logging.Logger
    mutex       sync.RWMutex
}

// NewProvider creates a new secret provider
func NewProvider(keyDir string, logger logging.Logger) (*Provider, error) {
    // Create key manager
    keyManager, err := NewKeyManager(keyDir, logger)
    if err != nil {
        return nil, fmt.Errorf("failed to create key manager: %w", err)
    }

    // Create age provider
    ageProvider, err := NewAgeProvider(filepath.Join(keyDir, "age.key"), logger)
    if err != nil {
        return nil, fmt.Errorf("failed to create age provider: %w", err)
    }

    return &Provider{
        ageProvider: ageProvider,
        keyManager:  keyManager,
        logger:      logger,
    }, nil
}

// EncryptSecret encrypts a secret value
func (p *Provider) EncryptSecret(value string) (string, error) {
    p.mutex.RLock()
    defer p.mutex.RUnlock()

    p.logger.Debug("Encrypting secret value")

    // Convert string to bytes
    data := []byte(value)

    // Encrypt with age
    encryptedData, err := p.ageProvider.Encrypt(data)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt secret: %w", err)
    }

    // Encode as base64
    encryptedString := base64.StdEncoding.EncodeToString(encryptedData)

    p.logger.Debug("Secret encrypted successfully")

    return encryptedString, nil
}

// DecryptSecret decrypts a secret value
func (p *Provider) DecryptSecret(encryptedValue string) (string, error) {
    p.mutex.RLock()
    defer p.mutex.RUnlock()

    p.logger.Debug("Decrypting secret value")

    // Decode from base64
    encryptedData, err := base64.StdEncoding.DecodeString(encryptedValue)
    if err != nil {
        return "", fmt.Errorf("failed to decode encrypted value: %w", err)
    }

    // Decrypt with age
    decryptedData, err := p.ageProvider.Decrypt(encryptedData)
    if err != nil {
        return "", fmt.Errorf("failed to decrypt secret: %w", err)
    }

    // Convert bytes to string
    decryptedString := string(decryptedData)

    p.logger.Debug("Secret decrypted successfully")

    return decryptedString, nil
}

// AddRecipient adds a recipient for encryption
func (p *Provider) AddRecipient(recipientString string) error {
    p.mutex.Lock()
    defer p.mutex.Unlock()

    return p.ageProvider.AddRecipient(recipientString)
}

// RemoveRecipient removes a recipient from encryption
func (p *Provider) RemoveRecipient(recipientString string) error {
    p.mutex.Lock()
    defer p.mutex.Unlock()

    return p.ageProvider.RemoveRecipient(recipientString)
}

// GetPublicKey returns the public key
func (p *Provider) GetPublicKey() (string, error) {
    p.mutex.RLock()
    defer p.mutex.RUnlock()

    return p.ageProvider.GetPublicKey()
}

// RotateKeys rotates encryption keys
func (p *Provider) RotateKeys() error {
    p.mutex.Lock()
    defer p.mutex.Unlock()

    return p.keyManager.RotateKeys()
}

// GenerateRandomSecret generates a random secret
func (p *Provider) GenerateRandomSecret(length int) (string, error) {
    return p.keyManager.GenerateRandomString(length)
}
```

### Step 4: Implement Key Rotation

#### 4.1 Key Rotation Manager
```go
// internal/templates/secrets/rotation.go
package secrets

import (
    "fmt"
    "time"
    
    "spooky/internal/logging"
)

// RotationManager manages key rotation
type RotationManager struct {
    provider    *Provider
    logger      logging.Logger
    rotationInterval time.Duration
    lastRotation    time.Time
}

// NewRotationManager creates a new rotation manager
func NewRotationManager(provider *Provider, rotationInterval time.Duration, logger logging.Logger) *RotationManager {
    return &RotationManager{
        provider:         provider,
        logger:           logger,
        rotationInterval: rotationInterval,
        lastRotation:     time.Now(),
    }
}

// ShouldRotate checks if keys should be rotated
func (rm *RotationManager) ShouldRotate() bool {
    return time.Since(rm.lastRotation) > rm.rotationInterval
}

// RotateKeysIfNeeded rotates keys if needed
func (rm *RotationManager) RotateKeysIfNeeded() error {
    if !rm.ShouldRotate() {
        return nil
    }

    rm.logger.Info("Performing scheduled key rotation")

    if err := rm.provider.RotateKeys(); err != nil {
        return fmt.Errorf("failed to rotate keys: %w", err)
    }

    rm.lastRotation = time.Now()

    rm.logger.Info("Key rotation completed successfully")

    return nil
}

// ForceRotation forces immediate key rotation
func (rm *RotationManager) ForceRotation() error {
    rm.logger.Info("Forcing key rotation")

    if err := rm.provider.RotateKeys(); err != nil {
        return fmt.Errorf("failed to force key rotation: %w", err)
    }

    rm.lastRotation = time.Now()

    rm.logger.Info("Forced key rotation completed successfully")

    return nil
}

// GetLastRotation returns the last rotation time
func (rm *RotationManager) GetLastRotation() time.Time {
    return rm.lastRotation
}

// GetNextRotation returns the next scheduled rotation time
func (rm *RotationManager) GetNextRotation() time.Time {
    return rm.lastRotation.Add(rm.rotationInterval)
}
```

## Configuration Options

### Supported Options
- **KeyDirectory**: Directory for storing keys
- **RotationInterval**: Key rotation interval
- **EnableArmor**: Enable/disable armor encoding
- **MultipleRecipients**: Enable/disable multiple recipients
- **AutoRotation**: Enable/disable automatic rotation

## Dependencies

### Internal Dependencies
- `spooky/internal/logging`

### External Dependencies
- `filippo.io/age`
- `filippo.io/age/armor`
- `crypto/rand` (standard library)
- `encoding/base64` (standard library)
- `fmt` (standard library)
- `io` (standard library)
- `os` (standard library)
- `path/filepath` (standard library)
- `sync` (standard library)
- `time` (standard library)

## Implementation Order

1. Implement age provider structure
2. Add age encryption methods
3. Create key manager implementation
4. Update secret provider
5. Add key rotation manager
6. Implement secure key storage
7. Add comprehensive tests
8. Performance optimization
9. Documentation and cleanup
