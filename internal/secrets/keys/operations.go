package keys

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"spooky/internal/logging"
	"spooky/internal/secrets/age"
	spookytypeslogging "spooky/internal/types/logging"
	"strings"
	"time"

	spookytypessecrets "spooky/internal/types/secrets"
)

// Manager implements the KeyManager interface
type Manager struct {
	globalKeysPath  string
	projectKeysPath string
	ageClient       age.AgeClient
	logger          spookytypeslogging.Logger
}

// NewManager creates a new key manager
func NewManager(globalKeysPath, projectKeysPath string) *Manager {
	return &Manager{
		globalKeysPath:  globalKeysPath,
		projectKeysPath: projectKeysPath,
		ageClient:       age.NewClient(),
		logger:          logging.GetLogger(),
	}
}

// ListKeys lists all available keys
func (m *Manager) ListKeys() ([]*spookytypessecrets.KeyMetadata, error) {
	m.logger.Debug("Listing available keys")

	var keys []*spookytypessecrets.KeyMetadata

	// List global keys
	if m.globalKeysPath != "" {
		globalKeys, err := m.listKeysFromPath(m.globalKeysPath)
		if err != nil {
			return nil, err
		}
		keys = append(keys, globalKeys...)
	}

	// List project keys (only public keys)
	if m.projectKeysPath != "" {
		projectKeys, err := m.listPublicKeysFromPath(m.projectKeysPath)
		if err != nil {
			return nil, err
		}
		keys = append(keys, projectKeys...)
	}

	return keys, nil
}

// GetKey gets a key by name
func (m *Manager) GetKey(name string) (*spookytypessecrets.KeyMetadata, error) {
	m.logger.Debug("Getting key", logging.String("name", name))

	// Try global keys first
	if m.globalKeysPath != "" {
		key, err := m.getKeyFromPath(m.globalKeysPath, name)
		if err == nil {
			return key, nil
		}
	}

	// Try project keys (only public keys)
	if m.projectKeysPath != "" {
		key, err := m.getPublicKeyFromPath(m.projectKeysPath, name)
		if err == nil {
			return key, nil
		}
	}

	return nil, &spookytypessecrets.SecretsError{
		Operation: "get_key",
		Cause:     fmt.Errorf("key not found: %s", name),
		Context: map[string]interface{}{
			"key_name": name,
		},
	}
}

// listKeysFromPath lists all keys from a specific path
func (m *Manager) listKeysFromPath(path string) ([]*spookytypessecrets.KeyMetadata, error) {
	var keys []*spookytypessecrets.KeyMetadata

	// Check if directory exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Directory doesn't exist, return empty list (no keys)
		return keys, nil
	}

	// Read directory entries
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, &spookytypessecrets.SecretsError{
			Operation: "read_keys_directory",
			Cause:     err,
			Context: map[string]interface{}{
				"path": path,
			},
		}
	}

	// Process each key file
	keyMap := make(map[string]*spookytypessecrets.KeyMetadata)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		if !strings.HasSuffix(fileName, ".identity") && !strings.HasSuffix(fileName, ".recipient") {
			continue
		}

		// Extract key name and type
		parts := strings.Split(fileName, ".")
		if len(parts) != 2 {
			continue
		}

		keyName := parts[0]
		keyType := parts[1]

		// Load key data
		keyPath := filepath.Join(path, fileName)
		keyData, err := m.loadKey(keyPath)
		if err != nil {
			m.logger.Warn("Failed to load key",
				logging.String("path", keyPath),
				logging.Error(err))
			continue
		}

		// Create or update metadata
		if keyMap[keyName] == nil {
			keyMap[keyName] = &spookytypessecrets.KeyMetadata{
				Name:    keyName,
				Created: time.Now(), // Use file modification time if available
			}
		}

		// Set key type and calculate fingerprint
		if keyType == "identity" {
			keyMap[keyName].Type = "identity"
		} else if keyType == "recipient" {
			keyMap[keyName].Type = "recipient"
		}

		// Calculate fingerprint
		hash := sha256.Sum256([]byte(keyData))
		keyMap[keyName].Fingerprint = hex.EncodeToString(hash[:])[:16]
	}

	// Convert map to slice
	for _, key := range keyMap {
		keys = append(keys, key)
	}

	return keys, nil
}

// listPublicKeysFromPath lists only public keys from a project path
func (m *Manager) listPublicKeysFromPath(path string) ([]*spookytypessecrets.KeyMetadata, error) {
	var keys []*spookytypessecrets.KeyMetadata

	// Check if directory exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Directory doesn't exist, return empty list (no keys)
		return keys, nil
	}

	// Read directory entries
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, &spookytypessecrets.SecretsError{
			Operation: "read_project_keys_directory",
			Cause:     err,
			Context: map[string]interface{}{
				"path": path,
			},
		}
	}

	// Process only recipient (public) key files
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		if !strings.HasSuffix(fileName, ".recipient") {
			continue
		}

		// Extract key name
		parts := strings.Split(fileName, ".")
		if len(parts) != 2 {
			continue
		}

		keyName := parts[0]

		// Load key data
		keyPath := filepath.Join(path, fileName)
		keyData, err := m.loadKey(keyPath)
		if err != nil {
			m.logger.Warn("Failed to load project key",
				logging.String("path", keyPath),
				logging.Error(err))
			continue
		}

		// Create metadata for public key
		hash := sha256.Sum256([]byte(keyData))
		key := &spookytypessecrets.KeyMetadata{
			Name:        keyName,
			Type:        "recipient",
			Created:     time.Now(),
			Fingerprint: hex.EncodeToString(hash[:])[:16],
		}

		keys = append(keys, key)
	}

	return keys, nil
}

// getKeyFromPath gets a key by name from a specific path
func (m *Manager) getKeyFromPath(path, name string) (*spookytypessecrets.KeyMetadata, error) {
	// Try identity key first
	identityPath := filepath.Join(path, fmt.Sprintf("%s.identity", name))
	if _, err := os.Stat(identityPath); err == nil {
		keyData, err := m.loadKey(identityPath)
		if err != nil {
			return nil, err
		}

		hash := sha256.Sum256([]byte(keyData))
		return &spookytypessecrets.KeyMetadata{
			Name:        name,
			Type:        "identity",
			Created:     time.Now(),
			Fingerprint: hex.EncodeToString(hash[:])[:16],
		}, nil
	}

	// Try recipient key
	recipientPath := filepath.Join(path, fmt.Sprintf("%s.recipient", name))
	if _, err := os.Stat(recipientPath); err == nil {
		keyData, err := m.loadKey(recipientPath)
		if err != nil {
			return nil, err
		}

		hash := sha256.Sum256([]byte(keyData))
		return &spookytypessecrets.KeyMetadata{
			Name:        name,
			Type:        "recipient",
			Created:     time.Now(),
			Fingerprint: hex.EncodeToString(hash[:])[:16],
		}, nil
	}

	return nil, &spookytypessecrets.SecretsError{
		Operation: "get_key_from_path",
		Cause:     fmt.Errorf("key not found: %s", name),
		Context: map[string]interface{}{
			"path":     path,
			"key_name": name,
		},
	}
}

// getPublicKeyFromPath gets only a public key by name from a project path
func (m *Manager) getPublicKeyFromPath(path, name string) (*spookytypessecrets.KeyMetadata, error) {
	// Only look for recipient (public) keys in project paths
	recipientPath := filepath.Join(path, fmt.Sprintf("%s.recipient", name))
	if _, err := os.Stat(recipientPath); err == nil {
		keyData, err := m.loadKey(recipientPath)
		if err != nil {
			return nil, err
		}

		hash := sha256.Sum256([]byte(keyData))
		return &spookytypessecrets.KeyMetadata{
			Name:        name,
			Type:        "recipient",
			Created:     time.Now(),
			Fingerprint: hex.EncodeToString(hash[:])[:16],
		}, nil
	}

	return nil, &spookytypessecrets.SecretsError{
		Operation: "get_public_key_from_path",
		Cause:     fmt.Errorf("public key not found: %s", name),
		Context: map[string]interface{}{
			"path":     path,
			"key_name": name,
		},
	}
}

// Helper functions

func (m *Manager) loadKey(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func (m *Manager) saveKey(path string, data string, permissions os.FileMode) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	// Write key data
	if err := os.WriteFile(path, []byte(data), permissions); err != nil {
		return err
	}

	return nil
}
