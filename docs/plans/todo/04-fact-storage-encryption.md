# Implementation Plan: Fact Storage Encryption

## Overview
Implement encryption functionality for fact storage to ensure sensitive data is protected at rest, supporting multiple encryption algorithms and key management strategies.

## Task Details
- **Task ID**: 6.2
- **Priority**: Low
- **File**: `internal/facts/storage/`
- **Functions**: Encryption for fact storage

## Current State Analysis

### Existing Patterns
1. **Fact Storage**: Multiple storage backends (BadgerDB, JSON, HCL) exist
2. **Crypto System**: Basic crypto operations exist in `internal/secrets/`
3. **Storage Interfaces**: Storage interfaces defined for different backends
4. **Configuration**: Encryption configuration patterns exist
5. **Key Management**: Basic key management system exists

### Existing Implementation Examples
- **Secrets Manager**: `internal/secrets/manager.go` provides crypto operations
- **Age Client**: `internal/secrets/age/client.go` provides age encryption
- **Storage Backends**: `internal/facts/storage/` provides storage implementations
- **Storage Interfaces**: `internal/facts/storage/interfaces.go` defines storage contracts

## Implementation Requirements

### Interface Compliance
The fact storage encryption must:
1. **Encrypt fact data** at rest in all storage backends
2. **Support multiple algorithms** (AES-256, ChaCha20-Poly1305, age)
3. **Provide key management** with rotation and backup
4. **Support transparent encryption** for existing storage
5. **Handle encryption metadata** and versioning
6. **Provide performance optimization** for encrypted storage
7. **Support encryption at rest** and in transit
8. **Integrate with existing** secrets management system

### Required Dependencies
- Secrets management system for crypto operations
- Storage backends for encrypted data
- Key management system for encryption keys
- Configuration system for encryption settings

## Detailed Implementation Plan

### Step 1: Define Encryption Interfaces

**File**: `internal/facts/storage/encryption/interfaces.go`

```go
package encryption

import (
    "context"
    "crypto/cipher"
    "io"

    "spooky/internal/facts/types"
    "spooky/internal/logging"
)

// Encryptor handles encryption and decryption of fact data
type Encryptor interface {
    // Encrypt encrypts fact data
    Encrypt(ctx context.Context, data []byte) ([]byte, error)
    
    // Decrypt decrypts fact data
    Decrypt(ctx context.Context, encryptedData []byte) ([]byte, error)
    
    // GetAlgorithm returns the encryption algorithm used
    GetAlgorithm() string
    
    // GetKeyID returns the key ID used for encryption
    GetKeyID() string
    
    // RotateKey rotates the encryption key
    RotateKey(ctx context.Context) error
    
    // ValidateKey validates the encryption key
    ValidateKey(ctx context.Context) error
}

// KeyManager handles encryption key management
type KeyManager interface {
    // GenerateKey generates a new encryption key
    GenerateKey(ctx context.Context, algorithm string) (*EncryptionKey, error)
    
    // GetKey retrieves an encryption key by ID
    GetKey(ctx context.Context, keyID string) (*EncryptionKey, error)
    
    // StoreKey stores an encryption key
    StoreKey(ctx context.Context, key *EncryptionKey) error
    
    // RotateKey rotates an encryption key
    RotateKey(ctx context.Context, keyID string) (*EncryptionKey, error)
    
    // BackupKeys backs up encryption keys
    BackupKeys(ctx context.Context, backupPath string) error
    
    // RestoreKeys restores encryption keys from backup
    RestoreKeys(ctx context.Context, backupPath string) error
    
    // ListKeys lists all available encryption keys
    ListKeys(ctx context.Context) ([]*EncryptionKey, error)
}

// EncryptionKey represents an encryption key
type EncryptionKey struct {
    ID          string                 `json:"id"`
    Algorithm   string                 `json:"algorithm"`
    Key         []byte                 `json:"key"`
    Created     time.Time              `json:"created"`
    Expires     *time.Time             `json:"expires,omitempty"`
    Metadata    map[string]interface{} `json:"metadata"`
    Version     int                    `json:"version"`
}

// EncryptionConfig represents encryption configuration
type EncryptionConfig struct {
    Enabled     bool                   `json:"enabled"`
    Algorithm   string                 `json:"algorithm"`
    KeyID       string                 `json:"key_id"`
    KeyPath     string                 `json:"key_path"`
    KeyRotation bool                   `json:"key_rotation"`
    AutoRotate  bool                   `json:"auto_rotate"`
    RotationDays int                   `json:"rotation_days"`
    BackupKeys  bool                   `json:"backup_keys"`
    BackupPath  string                 `json:"backup_path"`
    Metadata    map[string]interface{} `json:"metadata"`
}

// EncryptedData represents encrypted data with metadata
type EncryptedData struct {
    Version     int                    `json:"version"`
    Algorithm   string                 `json:"algorithm"`
    KeyID       string                 `json:"key_id"`
    IV          []byte                 `json:"iv"`
    Data        []byte                 `json:"data"`
    Metadata    map[string]interface{} `json:"metadata"`
    Timestamp   time.Time              `json:"timestamp"`
}

// EncryptionManager manages encryption operations
type EncryptionManager struct {
    config      *EncryptionConfig
    keyManager  KeyManager
    encryptor   Encryptor
    logger      logging.Logger
}

// NewEncryptionManager creates a new encryption manager
func NewEncryptionManager(config *EncryptionConfig, keyManager KeyManager, logger logging.Logger) *EncryptionManager {
    return &EncryptionManager{
        config:     config,
        keyManager: keyManager,
        logger:     logger,
    }
}
```

### Step 2: Implement Encryption Algorithms

**File**: `internal/facts/storage/encryption/algorithms.go`

```go
package encryption

import (
    "context"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "fmt"
    "io"

    "golang.org/x/crypto/chacha20poly1305"
    "spooky/internal/logging"
)

// AESEncryptor implements AES-256-GCM encryption
type AESEncryptor struct {
    key     []byte
    keyID   string
    logger  logging.Logger
}

// NewAESEncryptor creates a new AES encryptor
func NewAESEncryptor(key []byte, keyID string, logger logging.Logger) *AESEncryptor {
    return &AESEncryptor{
        key:    key,
        keyID:  keyID,
        logger: logger,
    }
}

// Encrypt encrypts data using AES-256-GCM
func (e *AESEncryptor) Encrypt(ctx context.Context, data []byte) ([]byte, error) {
    // Create AES cipher
    block, err := aes.NewCipher(e.key)
    if err != nil {
        return nil, fmt.Errorf("failed to create AES cipher: %w", err)
    }

    // Create GCM mode
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCM mode: %w", err)
    }

    // Generate nonce
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, fmt.Errorf("failed to generate nonce: %w", err)
    }

    // Encrypt data
    ciphertext := gcm.Seal(nonce, nonce, data, nil)

    return ciphertext, nil
}

// Decrypt decrypts data using AES-256-GCM
func (e *AESEncryptor) Decrypt(ctx context.Context, encryptedData []byte) ([]byte, error) {
    // Create AES cipher
    block, err := aes.NewCipher(e.key)
    if err != nil {
        return nil, fmt.Errorf("failed to create AES cipher: %w", err)
    }

    // Create GCM mode
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCM mode: %w", err)
    }

    // Extract nonce
    nonceSize := gcm.NonceSize()
    if len(encryptedData) < nonceSize {
        return nil, fmt.Errorf("encrypted data too short")
    }

    nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]

    // Decrypt data
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt data: %w", err)
    }

    return plaintext, nil
}

// GetAlgorithm returns the encryption algorithm
func (e *AESEncryptor) GetAlgorithm() string {
    return "AES-256-GCM"
}

// GetKeyID returns the key ID
func (e *AESEncryptor) GetKeyID() string {
    return e.keyID
}

// RotateKey rotates the encryption key
func (e *AESEncryptor) RotateKey(ctx context.Context) error {
    // This would generate a new key and update the encryptor
    // For now, we'll return an error indicating it's not implemented
    return fmt.Errorf("key rotation not implemented for AES encryptor")
}

// ValidateKey validates the encryption key
func (e *AESEncryptor) ValidateKey(ctx context.Context) error {
    if len(e.key) != 32 {
        return fmt.Errorf("invalid AES key length: expected 32 bytes, got %d", len(e.key))
    }
    return nil
}

// ChaCha20Encryptor implements ChaCha20-Poly1305 encryption
type ChaCha20Encryptor struct {
    key     []byte
    keyID   string
    logger  logging.Logger
}

// NewChaCha20Encryptor creates a new ChaCha20 encryptor
func NewChaCha20Encryptor(key []byte, keyID string, logger logging.Logger) *ChaCha20Encryptor {
    return &ChaCha20Encryptor{
        key:    key,
        keyID:  keyID,
        logger: logger,
    }
}

// Encrypt encrypts data using ChaCha20-Poly1305
func (e *ChaCha20Encryptor) Encrypt(ctx context.Context, data []byte) ([]byte, error) {
    // Create ChaCha20-Poly1305 cipher
    aead, err := chacha20poly1305.New(e.key)
    if err != nil {
        return nil, fmt.Errorf("failed to create ChaCha20-Poly1305 cipher: %w", err)
    }

    // Generate nonce
    nonce := make([]byte, aead.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, fmt.Errorf("failed to generate nonce: %w", err)
    }

    // Encrypt data
    ciphertext := aead.Seal(nonce, nonce, data, nil)

    return ciphertext, nil
}

// Decrypt decrypts data using ChaCha20-Poly1305
func (e *ChaCha20Encryptor) Decrypt(ctx context.Context, encryptedData []byte) ([]byte, error) {
    // Create ChaCha20-Poly1305 cipher
    aead, err := chacha20poly1305.New(e.key)
    if err != nil {
        return nil, fmt.Errorf("failed to create ChaCha20-Poly1305 cipher: %w", err)
    }

    // Extract nonce
    nonceSize := aead.NonceSize()
    if len(encryptedData) < nonceSize {
        return nil, fmt.Errorf("encrypted data too short")
    }

    nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]

    // Decrypt data
    plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt data: %w", err)
    }

    return plaintext, nil
}

// GetAlgorithm returns the encryption algorithm
func (e *ChaCha20Encryptor) GetAlgorithm() string {
    return "ChaCha20-Poly1305"
}

// GetKeyID returns the key ID
func (e *ChaCha20Encryptor) GetKeyID() string {
    return e.keyID
}

// RotateKey rotates the encryption key
func (e *ChaCha20Encryptor) RotateKey(ctx context.Context) error {
    return fmt.Errorf("key rotation not implemented for ChaCha20 encryptor")
}

// ValidateKey validates the encryption key
func (e *ChaCha20Encryptor) ValidateKey(ctx context.Context) error {
    if len(e.key) != chacha20poly1305.KeySize {
        return fmt.Errorf("invalid ChaCha20 key length: expected %d bytes, got %d", chacha20poly1305.KeySize, len(e.key))
    }
    return nil
}
```

### Step 3: Implement Key Management

**File**: `internal/facts/storage/encryption/keys.go`

```go
package encryption

import (
    "context"
    "crypto/rand"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "time"

    "spooky/internal/logging"
    "spooky/internal/secrets"
)

// FileKeyManager implements file-based key management
type FileKeyManager struct {
    keyPath string
    logger  logging.Logger
}

// NewFileKeyManager creates a new file-based key manager
func NewFileKeyManager(keyPath string, logger logging.Logger) *FileKeyManager {
    return &FileKeyManager{
        keyPath: keyPath,
        logger:  logger,
    }
}

// GenerateKey generates a new encryption key
func (m *FileKeyManager) GenerateKey(ctx context.Context, algorithm string) (*EncryptionKey, error) {
    var keySize int
    switch algorithm {
    case "AES-256-GCM":
        keySize = 32
    case "ChaCha20-Poly1305":
        keySize = 32
    default:
        return nil, fmt.Errorf("unsupported algorithm: %s", algorithm)
    }

    // Generate random key
    key := make([]byte, keySize)
    if _, err := rand.Read(key); err != nil {
        return nil, fmt.Errorf("failed to generate random key: %w", err)
    }

    // Create encryption key
    encryptionKey := &EncryptionKey{
        ID:        m.generateKeyID(),
        Algorithm: algorithm,
        Key:       key,
        Created:   time.Now(),
        Metadata:  make(map[string]interface{}),
        Version:   1,
    }

    return encryptionKey, nil
}

// GetKey retrieves an encryption key by ID
func (m *FileKeyManager) GetKey(ctx context.Context, keyID string) (*EncryptionKey, error) {
    keyPath := filepath.Join(m.keyPath, fmt.Sprintf("%s.key", keyID))
    
    data, err := os.ReadFile(keyPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read key file: %w", err)
    }

    var key EncryptionKey
    if err := json.Unmarshal(data, &key); err != nil {
        return nil, fmt.Errorf("failed to unmarshal key: %w", err)
    }

    return &key, nil
}

// StoreKey stores an encryption key
func (m *FileKeyManager) StoreKey(ctx context.Context, key *EncryptionKey) error {
    // Ensure key directory exists
    if err := os.MkdirAll(m.keyPath, 0700); err != nil {
        return fmt.Errorf("failed to create key directory: %w", err)
    }

    // Marshal key to JSON
    data, err := json.MarshalIndent(key, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal key: %w", err)
    }

    // Write key to file
    keyPath := filepath.Join(m.keyPath, fmt.Sprintf("%s.key", key.ID))
    if err := os.WriteFile(keyPath, data, 0600); err != nil {
        return fmt.Errorf("failed to write key file: %w", err)
    }

    m.logger.Info("Encryption key stored",
        logging.String("key_id", key.ID),
        logging.String("algorithm", key.Algorithm))

    return nil
}

// RotateKey rotates an encryption key
func (m *FileKeyManager) RotateKey(ctx context.Context, keyID string) (*EncryptionKey, error) {
    // Get existing key
    existingKey, err := m.GetKey(ctx, keyID)
    if err != nil {
        return nil, fmt.Errorf("failed to get existing key: %w", err)
    }

    // Generate new key
    newKey, err := m.GenerateKey(ctx, existingKey.Algorithm)
    if err != nil {
        return nil, fmt.Errorf("failed to generate new key: %w", err)
    }

    // Update metadata
    newKey.Metadata["rotated_from"] = keyID
    newKey.Metadata["rotation_date"] = time.Now()

    // Store new key
    if err := m.StoreKey(ctx, newKey); err != nil {
        return nil, fmt.Errorf("failed to store new key: %w", err)
    }

    m.logger.Info("Encryption key rotated",
        logging.String("old_key_id", keyID),
        logging.String("new_key_id", newKey.ID))

    return newKey, nil
}

// BackupKeys backs up encryption keys
func (m *FileKeyManager) BackupKeys(ctx context.Context, backupPath string) error {
    // Ensure backup directory exists
    if err := os.MkdirAll(filepath.Dir(backupPath), 0700); err != nil {
        return fmt.Errorf("failed to create backup directory: %w", err)
    }

    // List all keys
    keys, err := m.ListKeys(ctx)
    if err != nil {
        return fmt.Errorf("failed to list keys: %w", err)
    }

    // Create backup data
    backup := map[string]interface{}{
        "timestamp": time.Now(),
        "keys":      keys,
        "version":   1,
    }

    // Marshal backup to JSON
    data, err := json.MarshalIndent(backup, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal backup: %w", err)
    }

    // Write backup file
    if err := os.WriteFile(backupPath, data, 0600); err != nil {
        return fmt.Errorf("failed to write backup file: %w", err)
    }

    m.logger.Info("Encryption keys backed up",
        logging.String("backup_path", backupPath),
        logging.Int("key_count", len(keys)))

    return nil
}

// RestoreKeys restores encryption keys from backup
func (m *FileKeyManager) RestoreKeys(ctx context.Context, backupPath string) error {
    // Read backup file
    data, err := os.ReadFile(backupPath)
    if err != nil {
        return fmt.Errorf("failed to read backup file: %w", err)
    }

    // Unmarshal backup
    var backup map[string]interface{}
    if err := json.Unmarshal(data, &backup); err != nil {
        return fmt.Errorf("failed to unmarshal backup: %w", err)
    }

    // Extract keys
    keysData, ok := backup["keys"].([]interface{})
    if !ok {
        return fmt.Errorf("invalid backup format: keys not found")
    }

    // Restore each key
    for _, keyData := range keysData {
        keyBytes, err := json.Marshal(keyData)
        if err != nil {
            return fmt.Errorf("failed to marshal key data: %w", err)
        }

        var key EncryptionKey
        if err := json.Unmarshal(keyBytes, &key); err != nil {
            return fmt.Errorf("failed to unmarshal key: %w", err)
        }

        if err := m.StoreKey(ctx, &key); err != nil {
            return fmt.Errorf("failed to store restored key: %w", err)
        }
    }

    m.logger.Info("Encryption keys restored",
        logging.String("backup_path", backupPath),
        logging.Int("key_count", len(keysData)))

    return nil
}

// ListKeys lists all available encryption keys
func (m *FileKeyManager) ListKeys(ctx context.Context) ([]*EncryptionKey, error) {
    var keys []*EncryptionKey

    // Read key directory
    entries, err := os.ReadDir(m.keyPath)
    if err != nil {
        if os.IsNotExist(err) {
            return keys, nil // No keys directory exists
        }
        return nil, fmt.Errorf("failed to read key directory: %w", err)
    }

    // Process each key file
    for _, entry := range entries {
        if entry.IsDir() || filepath.Ext(entry.Name()) != ".key" {
            continue
        }

        keyID := strings.TrimSuffix(entry.Name(), ".key")
        key, err := m.GetKey(ctx, keyID)
        if err != nil {
            m.logger.Warn("Failed to load key",
                logging.String("key_id", keyID),
                logging.Error(err))
            continue
        }

        keys = append(keys, key)
    }

    return keys, nil
}

// generateKeyID generates a unique key ID
func (m *FileKeyManager) generateKeyID() string {
    // Generate random bytes for key ID
    idBytes := make([]byte, 16)
    rand.Read(idBytes)
    
    // Convert to hex string
    return fmt.Sprintf("%x", idBytes)
}
```

### Step 4: Implement Encrypted Storage Wrapper

**File**: `internal/facts/storage/encryption/storage.go`

```go
package encryption

import (
    "context"
    "encoding/json"
    "fmt"

    "spooky/internal/facts/storage"
    "spooky/internal/facts/types"
    "spooky/internal/logging"
)

// EncryptedStorage wraps a storage backend with encryption
type EncryptedStorage struct {
    storage   storage.Storage
    encryptor Encryptor
    logger    logging.Logger
}

// NewEncryptedStorage creates a new encrypted storage wrapper
func NewEncryptedStorage(storage storage.Storage, encryptor Encryptor, logger logging.Logger) *EncryptedStorage {
    return &EncryptedStorage{
        storage:   storage,
        encryptor: encryptor,
        logger:    logger,
    }
}

// Store stores encrypted fact data
func (s *EncryptedStorage) Store(ctx context.Context, collection *types.FactCollection) error {
    // Marshal collection to JSON
    data, err := json.Marshal(collection)
    if err != nil {
        return fmt.Errorf("failed to marshal collection: %w", err)
    }

    // Encrypt data
    encryptedData, err := s.encryptor.Encrypt(ctx, data)
    if err != nil {
        return fmt.Errorf("failed to encrypt collection: %w", err)
    }

    // Create encrypted collection
    encryptedCollection := &types.FactCollection{
        Server:    collection.Server,
        Timestamp: collection.Timestamp,
        Facts:     make(map[string]*types.Fact),
    }

    // Store encrypted data as a special fact
    encryptedCollection.Facts["_encrypted_data"] = &types.Fact{
        Key:       "_encrypted_data",
        Value:     encryptedData,
        Source:    "encrypted",
        Server:    collection.Server,
        Timestamp: collection.Timestamp,
    }

    // Store using underlying storage
    return s.storage.Store(ctx, encryptedCollection)
}

// Load loads and decrypts fact data
func (s *EncryptedStorage) Load(ctx context.Context, server string) (*types.FactCollection, error) {
    // Load from underlying storage
    collection, err := s.storage.Load(ctx, server)
    if err != nil {
        return nil, fmt.Errorf("failed to load collection: %w", err)
    }

    // Check if data is encrypted
    encryptedFact, exists := collection.Facts["_encrypted_data"]
    if !exists {
        // Data is not encrypted, return as-is
        return collection, nil
    }

    // Extract encrypted data
    encryptedData, ok := encryptedFact.Value.([]byte)
    if !ok {
        return nil, fmt.Errorf("invalid encrypted data format")
    }

    // Decrypt data
    decryptedData, err := s.encryptor.Decrypt(ctx, encryptedData)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt collection: %w", err)
    }

    // Unmarshal decrypted data
    var decryptedCollection types.FactCollection
    if err := json.Unmarshal(decryptedData, &decryptedCollection); err != nil {
        return nil, fmt.Errorf("failed to unmarshal decrypted collection: %w", err)
    }

    return &decryptedCollection, nil
}

// Delete deletes encrypted fact data
func (s *EncryptedStorage) Delete(ctx context.Context, server string) error {
    return s.storage.Delete(ctx, server)
}

// List lists all servers with encrypted data
func (s *EncryptedStorage) List(ctx context.Context) ([]string, error) {
    return s.storage.List(ctx)
}

// Close closes the encrypted storage
func (s *EncryptedStorage) Close() error {
    return s.storage.Close()
}

// GetEncryptionInfo returns encryption information
func (s *EncryptedStorage) GetEncryptionInfo() map[string]interface{} {
    return map[string]interface{}{
        "algorithm": s.encryptor.GetAlgorithm(),
        "key_id":    s.encryptor.GetKeyID(),
        "enabled":   true,
    }
}
```

### Step 5: Implement Encryption Manager

**File**: `internal/facts/storage/encryption/manager.go`

```go
package encryption

import (
    "context"
    "fmt"
    "time"

    "spooky/internal/facts/storage"
    "spooky/internal/logging"
)

// Manager manages encryption operations for fact storage
type Manager struct {
    config      *EncryptionConfig
    keyManager  KeyManager
    encryptor   Encryptor
    logger      logging.Logger
}

// NewManager creates a new encryption manager
func NewManager(config *EncryptionConfig, keyManager KeyManager, logger logging.Logger) *Manager {
    return &Manager{
        config:     config,
        keyManager: keyManager,
        logger:     logger,
    }
}

// Initialize initializes the encryption manager
func (m *Manager) Initialize(ctx context.Context) error {
    if !m.config.Enabled {
        m.logger.Info("Encryption disabled")
        return nil
    }

    // Get or create encryption key
    key, err := m.getOrCreateKey(ctx)
    if err != nil {
        return fmt.Errorf("failed to get or create encryption key: %w", err)
    }

    // Create encryptor
    encryptor, err := m.createEncryptor(key)
    if err != nil {
        return fmt.Errorf("failed to create encryptor: %w", err)
    }

    m.encryptor = encryptor

    // Validate key
    if err := encryptor.ValidateKey(ctx); err != nil {
        return fmt.Errorf("failed to validate encryption key: %w", err)
    }

    // Setup auto-rotation if enabled
    if m.config.AutoRotate {
        go m.autoRotateKeys(ctx)
    }

    m.logger.Info("Encryption manager initialized",
        logging.String("algorithm", key.Algorithm),
        logging.String("key_id", key.ID))

    return nil
}

// WrapStorage wraps a storage backend with encryption
func (m *Manager) WrapStorage(storage storage.Storage) storage.Storage {
    if !m.config.Enabled || m.encryptor == nil {
        return storage
    }

    return NewEncryptedStorage(storage, m.encryptor, m.logger)
}

// RotateKeys rotates encryption keys
func (m *Manager) RotateKeys(ctx context.Context) error {
    if m.encryptor == nil {
        return fmt.Errorf("encryption not initialized")
    }

    // Rotate the current key
    newKey, err := m.keyManager.RotateKey(ctx, m.encryptor.GetKeyID())
    if err != nil {
        return fmt.Errorf("failed to rotate key: %w", err)
    }

    // Create new encryptor
    newEncryptor, err := m.createEncryptor(newKey)
    if err != nil {
        return fmt.Errorf("failed to create new encryptor: %w", err)
    }

    // Update encryptor
    m.encryptor = newEncryptor

    m.logger.Info("Encryption keys rotated",
        logging.String("new_key_id", newKey.ID))

    return nil
}

// BackupKeys backs up encryption keys
func (m *Manager) BackupKeys(ctx context.Context) error {
    if m.config.BackupKeys {
        return m.keyManager.BackupKeys(ctx, m.config.BackupPath)
    }
    return nil
}

// GetEncryptionInfo returns encryption information
func (m *Manager) GetEncryptionInfo() map[string]interface{} {
    if !m.config.Enabled || m.encryptor == nil {
        return map[string]interface{}{
            "enabled": false,
        }
    }

    return map[string]interface{}{
        "enabled":   true,
        "algorithm": m.encryptor.GetAlgorithm(),
        "key_id":    m.encryptor.GetKeyID(),
        "config":    m.config,
    }
}

// getOrCreateKey gets or creates an encryption key
func (m *Manager) getOrCreateKey(ctx context.Context) (*EncryptionKey, error) {
    if m.config.KeyID != "" {
        // Try to get existing key
        key, err := m.keyManager.GetKey(ctx, m.config.KeyID)
        if err == nil {
            return key, nil
        }
        m.logger.Warn("Failed to get existing key, creating new one",
            logging.String("key_id", m.config.KeyID),
            logging.Error(err))
    }

    // Generate new key
    key, err := m.keyManager.GenerateKey(ctx, m.config.Algorithm)
    if err != nil {
        return nil, fmt.Errorf("failed to generate encryption key: %w", err)
    }

    // Store key
    if err := m.keyManager.StoreKey(ctx, key); err != nil {
        return nil, fmt.Errorf("failed to store encryption key: %w", err)
    }

    return key, nil
}

// createEncryptor creates an encryptor for the given key
func (m *Manager) createEncryptor(key *EncryptionKey) (Encryptor, error) {
    switch key.Algorithm {
    case "AES-256-GCM":
        return NewAESEncryptor(key.Key, key.ID, m.logger), nil
    case "ChaCha20-Poly1305":
        return NewChaCha20Encryptor(key.Key, key.ID, m.logger), nil
    default:
        return nil, fmt.Errorf("unsupported encryption algorithm: %s", key.Algorithm)
    }
}

// autoRotateKeys automatically rotates keys
func (m *Manager) autoRotateKeys(ctx context.Context) {
    ticker := time.NewTicker(24 * time.Hour) // Check daily
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            if err := m.checkAndRotateKeys(ctx); err != nil {
                m.logger.Error("Failed to check and rotate keys",
                    logging.Error(err))
            }
        }
    }
}

// checkAndRotateKeys checks if keys need rotation and rotates them
func (m *Manager) checkAndRotateKeys(ctx context.Context) error {
    if m.encryptor == nil {
        return nil
    }

    // Get current key
    key, err := m.keyManager.GetKey(ctx, m.encryptor.GetKeyID())
    if err != nil {
        return fmt.Errorf("failed to get current key: %w", err)
    }

    // Check if key needs rotation
    if key.Expires != nil && time.Now().After(*key.Expires) {
        m.logger.Info("Key expired, rotating",
            logging.String("key_id", key.ID))
        return m.RotateKeys(ctx)
    }

    // Check rotation days
    if m.config.RotationDays > 0 {
        age := time.Since(key.Created)
        if age > time.Duration(m.config.RotationDays)*24*time.Hour {
            m.logger.Info("Key rotation period reached, rotating",
                logging.String("key_id", key.ID),
                logging.Duration("age", age))
            return m.RotateKeys(ctx)
        }
    }

    return nil
}
```







## Configuration Options

### Supported Options
- **Encryption enabled**: Enable/disable encryption
- **Algorithm**: AES-256-GCM, ChaCha20-Poly1305
- **Key rotation**: Automatic key rotation
- **Backup keys**: Automatic key backup
- **Performance**: Encryption performance settings

## Dependencies

### Internal Dependencies
- `spooky/internal/facts/storage`
- `spooky/internal/facts/types`
- `spooky/internal/secrets`
- `spooky/internal/logging`

### External Dependencies
- `golang.org/x/crypto/chacha20poly1305`
- `crypto/aes` (standard library)
- `crypto/cipher` (standard library)
- `crypto/rand` (standard library)



## Implementation Order

1. Define encryption interfaces
2. Implement encryption algorithms
3. Add key management system
4. Implement encrypted storage wrapper
5. Add encryption manager
6. Write comprehensive tests
7. Performance optimization
8. Documentation and cleanup


