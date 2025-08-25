package encryption

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"filippo.io/age"
	"filippo.io/age/armor"
	"github.com/pkg/errors"

	"spooky/internal/schemas"
	"spooky/internal/utilities"
)

// AgeEncryption provides age encryption functionality
type AgeEncryption struct {
	identitiesPath string
	recipientsPath string
	identities     []age.Identity
	recipients     []age.Recipient
}

// NewAgeEncryption creates a new age encryption instance
func NewAgeEncryption(identitiesPath, recipientsPath string) (*AgeEncryption, error) {
	ae := &AgeEncryption{
		identitiesPath: identitiesPath,
		recipientsPath: recipientsPath,
	}

	// Load identities if path is provided
	if identitiesPath != "" {
		if err := ae.loadIdentities(); err != nil {
			return nil, errors.Wrap(err, "failed to load age identities")
		}
	}

	// Load recipients if path is provided
	if recipientsPath != "" {
		if err := ae.loadRecipients(); err != nil {
			return nil, errors.Wrap(err, "failed to load age recipients")
		}
	}

	return ae, nil
}

// loadIdentities loads age identities from the specified path
func (ae *AgeEncryption) loadIdentities() error {
	// Check if path is a directory
	info, err := os.Stat(ae.identitiesPath)
	if err != nil {
		return errors.Wrapf(err, "failed to stat identities path: %s", ae.identitiesPath)
	}

	if info.IsDir() {
		// Load all identity files from directory
		return ae.loadIdentitiesFromDirectory()
	} else {
		// Load single identity file
		return ae.loadIdentityFromFile(ae.identitiesPath)
	}
}

// loadIdentitiesFromDirectory loads all identity files from a directory
func (ae *AgeEncryption) loadIdentitiesFromDirectory() error {
	entries, err := os.ReadDir(ae.identitiesPath)
	if err != nil {
		return errors.Wrapf(err, "failed to read identities directory: %s", ae.identitiesPath)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Skip hidden files
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		identityPath := filepath.Join(ae.identitiesPath, entry.Name())
		if err := ae.loadIdentityFromFile(identityPath); err != nil {
			// Log error but continue loading other identities
			fmt.Printf("Warning: failed to load identity from %s: %v\n", identityPath, err)
		}
	}

	return nil
}

// loadIdentityFromFile loads a single identity from a file
func (ae *AgeEncryption) loadIdentityFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.Wrapf(err, "failed to open identity file: %s", path)
	}
	defer file.Close()

	// Try to read as armored format first
	file.Seek(0, 0)
	armoredReader := armor.NewReader(file)
	identities, err := age.ParseIdentities(armoredReader)
	if err != nil {
		// Try as raw format
		file.Seek(0, 0)
		identities, err = age.ParseIdentities(file)
		if err != nil {
			return errors.Wrapf(err, "failed to parse identity file: %s", path)
		}
	}

	ae.identities = append(ae.identities, identities...)
	return nil
}

// loadRecipients loads age recipients from the specified path
func (ae *AgeEncryption) loadRecipients() error {
	// Check if path is a directory
	info, err := os.Stat(ae.recipientsPath)
	if err != nil {
		return errors.Wrapf(err, "failed to stat recipients path: %s", ae.recipientsPath)
	}

	if info.IsDir() {
		// Load all recipient files from directory
		return ae.loadRecipientsFromDirectory()
	} else {
		// Load single recipient file
		return ae.loadRecipientsFromFile(ae.recipientsPath)
	}
}

// loadRecipientsFromDirectory loads all recipient files from a directory
func (ae *AgeEncryption) loadRecipientsFromDirectory() error {
	entries, err := os.ReadDir(ae.recipientsPath)
	if err != nil {
		return errors.Wrapf(err, "failed to read recipients directory: %s", ae.recipientsPath)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Skip hidden files
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// Skip non-text files
		if !isTextFile(entry.Name()) {
			continue
		}

		recipientPath := filepath.Join(ae.recipientsPath, entry.Name())
		if err := ae.loadRecipientsFromFile(recipientPath); err != nil {
			// Log error but continue loading other recipients
			fmt.Printf("Warning: failed to load recipients from %s: %v\n", recipientPath, err)
		}
	}

	return nil
}

// loadRecipientsFromFile loads recipients from a single file
func (ae *AgeEncryption) loadRecipientsFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.Wrapf(err, "failed to open recipients file: %s", path)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		recipient, err := age.ParseX25519Recipient(line)
		if err != nil {
			return errors.Wrapf(err, "failed to parse recipient on line %d in %s: %s", lineNum, path, line)
		}

		ae.recipients = append(ae.recipients, recipient)
	}

	if err := scanner.Err(); err != nil {
		return errors.Wrapf(err, "failed to read recipients file: %s", path)
	}

	return nil
}

// isTextFile checks if a file appears to be a text file based on extension
func isTextFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	textExtensions := []string{".txt", ".pub", ".key", ".age", ""} // empty extension for files without extension

	for _, textExt := range textExtensions {
		if ext == textExt {
			return true
		}
	}
	return false
}

// Encrypt encrypts a plaintext value using age encryption
func (ae *AgeEncryption) Encrypt(plaintext string) (string, error) {
	if len(ae.recipients) == 0 {
		return "", errors.New("no recipients available for encryption")
	}

	// Create encrypted output
	var encrypted strings.Builder
	output := armor.NewWriter(&encrypted)
	defer output.Close()

	// Create recipient list
	recipients := make([]age.Recipient, len(ae.recipients))
	copy(recipients, ae.recipients)

	// Create encrypted writer
	encryptedWriter, err := age.Encrypt(output, recipients...)
	if err != nil {
		return "", errors.Wrap(err, "failed to create encrypted writer")
	}
	defer encryptedWriter.Close()

	// Write the plaintext
	if _, err := encryptedWriter.Write([]byte(plaintext)); err != nil {
		return "", errors.Wrap(err, "failed to write plaintext to encrypted writer")
	}

	// Close the encrypted writer to finalize encryption
	if err := encryptedWriter.Close(); err != nil {
		return "", errors.Wrap(err, "failed to finalize encryption")
	}

	// Close the armor writer
	if err := output.Close(); err != nil {
		return "", errors.Wrap(err, "failed to close armor writer")
	}

	return encrypted.String(), nil
}

// Decrypt decrypts an age-encrypted value
func (ae *AgeEncryption) Decrypt(encrypted string) (string, error) {
	if len(ae.identities) == 0 {
		return "", errors.New("no identities available for decryption")
	}

	// Create armored reader
	armoredReader := armor.NewReader(strings.NewReader(encrypted))

	// Create decrypted reader
	decryptedReader, err := age.Decrypt(armoredReader, ae.identities...)
	if err != nil {
		return "", errors.Wrap(err, "failed to create decrypted reader")
	}

	// Read the decrypted content
	decryptedBytes, err := io.ReadAll(decryptedReader)
	if err != nil {
		return "", errors.Wrap(err, "failed to read decrypted content")
	}

	return string(decryptedBytes), nil
}

// IsEncrypted checks if a value appears to be age-encrypted
func (ae *AgeEncryption) IsEncrypted(value string) bool {
	return strings.HasPrefix(strings.TrimSpace(value), "-----BEGIN AGE ENCRYPTED FILE-----")
}

// GetIdentitiesCount returns the number of loaded identities
func (ae *AgeEncryption) GetIdentitiesCount() int {
	return len(ae.identities)
}

// GetRecipientsCount returns the number of loaded recipients
func (ae *AgeEncryption) GetRecipientsCount() int {
	return len(ae.recipients)
}

// ValidateConfiguration validates that the age encryption is properly configured
func (ae *AgeEncryption) ValidateConfiguration() error {
	if len(ae.recipients) == 0 {
		return errors.New("no recipients available for encryption")
	}

	if len(ae.identities) == 0 {
		return errors.New("no identities available for decryption")
	}

	return nil
}

// GetDefaultAgePaths returns the default paths for age configuration
func GetDefaultAgePaths() (identitiesPath, recipientsPath string) {
	// Use OS detection utility to get proper config directory
	pathConfig, err := utilities.GetPathConfig("spooky")
	if err != nil {
		// Fallback to current directory if OS detection fails
		identitiesPath = ".spooky/age/identities"
		recipientsPath = ".spooky/age/recipients"
		return
	}

	// Use the config directory from OS detection
	identitiesPath = filepath.Join(pathConfig.ConfigDir, "age", "identities")
	recipientsPath = filepath.Join(pathConfig.ConfigDir, "age", "recipients")
	return
}

// GetProjectAgePaths returns age paths for a project, using project-specific overrides if available
func GetProjectAgePaths(projectAge *schemas.ProjectAgeV1) (identitiesPath, recipientsPath string) {
	defaultIdentities, defaultRecipients := GetDefaultAgePaths()

	// Use project-specific paths if provided, otherwise use defaults
	if projectAge != nil && projectAge.DefaultIdentitiesPath != "" {
		identitiesPath = projectAge.DefaultIdentitiesPath
	} else {
		identitiesPath = defaultIdentities
	}

	if projectAge != nil && projectAge.DefaultRecipientsPath != "" {
		recipientsPath = projectAge.DefaultRecipientsPath
	} else {
		recipientsPath = defaultRecipients
	}

	return
}
