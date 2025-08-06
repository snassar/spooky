package operations

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"spooky/internal/logging"
	"spooky/internal/secrets/age"
	"spooky/internal/secrets/types"
)

// FileManager handles file encryption and decryption operations
type FileManager struct {
	ageClient age.AgeClient
	logger    logging.Logger
}

// NewFileManager creates a new file manager
func NewFileManager() *FileManager {
	return &FileManager{
		ageClient: age.NewClient(),
		logger:    logging.GetLogger(),
	}
}

// EncryptFile encrypts a file with the given recipients
func (fm *FileManager) EncryptFile(inputPath string, outputPath string, recipients []string) error {
	fm.logger.Info("Encrypting file",
		logging.String("input", inputPath),
		logging.String("output", outputPath),
		logging.Int("recipients", len(recipients)))

	// Validate input file
	if err := fm.validateInputFile(inputPath); err != nil {
		return err
	}

	// Validate recipients
	if len(recipients) == 0 {
		return &types.SecretsError{
			Operation: "validate_recipients",
			Cause:     fmt.Errorf("no recipients provided"),
		}
	}

	// Open input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return &types.SecretsError{
			Operation: "open_input_file",
			Cause:     err,
			Context: map[string]interface{}{
				"path": inputPath,
			},
		}
	}
	defer inputFile.Close()

	// Create output directory if needed
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return &types.SecretsError{
			Operation: "create_output_directory",
			Cause:     err,
			Context: map[string]interface{}{
				"path": outputDir,
			},
		}
	}

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return &types.SecretsError{
			Operation: "create_output_file",
			Cause:     err,
			Context: map[string]interface{}{
				"path": outputPath,
			},
		}
	}
	defer outputFile.Close()

	// Encrypt the file using streaming
	if err := fm.ageClient.EncryptStream(inputFile, outputFile, recipients); err != nil {
		return &types.SecretsError{
			Operation: "encrypt_file_stream",
			Cause:     err,
		}
	}

	fm.logger.Info("File encrypted successfully",
		logging.String("input", inputPath),
		logging.String("output", outputPath))

	return nil
}

// DecryptFile decrypts a file with the given identity
func (fm *FileManager) DecryptFile(inputPath string, outputPath string, identity string) error {
	fm.logger.Info("Decrypting file",
		logging.String("input", inputPath),
		logging.String("output", outputPath))

	// Validate input file
	if err := fm.validateInputFile(inputPath); err != nil {
		return err
	}

	// Validate identity
	if identity == "" {
		return &types.SecretsError{
			Operation: "validate_identity",
			Cause:     fmt.Errorf("no identity provided"),
		}
	}

	// Open input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return &types.SecretsError{
			Operation: "open_input_file",
			Cause:     err,
			Context: map[string]interface{}{
				"path": inputPath,
			},
		}
	}
	defer inputFile.Close()

	// Create output directory if needed
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return &types.SecretsError{
			Operation: "create_output_directory",
			Cause:     err,
			Context: map[string]interface{}{
				"path": outputDir,
			},
		}
	}

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return &types.SecretsError{
			Operation: "create_output_file",
			Cause:     err,
			Context: map[string]interface{}{
				"path": outputPath,
			},
		}
	}
	defer outputFile.Close()

	// Decrypt the file using streaming
	if err := fm.ageClient.DecryptStream(inputFile, outputFile, identity); err != nil {
		return &types.SecretsError{
			Operation: "decrypt_file_stream",
			Cause:     err,
		}
	}

	fm.logger.Info("File decrypted successfully",
		logging.String("input", inputPath),
		logging.String("output", outputPath))

	return nil
}

// EncryptValue encrypts a single value and returns the encrypted string
func (fm *FileManager) EncryptValue(value string, recipients []string) (string, error) {
	fm.logger.Debug("Encrypting value", logging.Int("recipients", len(recipients)))

	// Validate recipients
	if len(recipients) == 0 {
		return "", &types.SecretsError{
			Operation: "validate_recipients",
			Cause:     fmt.Errorf("no recipients provided"),
		}
	}

	// Encrypt the value
	encrypted, err := fm.ageClient.Encrypt([]byte(value), recipients)
	if err != nil {
		return "", err
	}

	// Convert to string (base64 encoded)
	encryptedStr := string(encrypted.Data)
	fm.logger.Debug("Value encrypted successfully")

	return encryptedStr, nil
}

// DecryptValue decrypts a single value and returns the decrypted string
func (fm *FileManager) DecryptValue(encryptedValue string, identity string) (string, error) {
	fm.logger.Debug("Decrypting value")

	// Validate identity
	if identity == "" {
		return "", &types.SecretsError{
			Operation: "validate_identity",
			Cause:     fmt.Errorf("no identity provided"),
		}
	}

	// Create encrypted value from string
	encrypted := &types.EncryptedValue{
		Data: []byte(encryptedValue),
	}

	// Decrypt the value
	decrypted, err := fm.ageClient.Decrypt(encrypted, identity)
	if err != nil {
		return "", err
	}

	decryptedStr := string(decrypted)
	fm.logger.Debug("Value decrypted successfully")

	return decryptedStr, nil
}

// ValidateInputFile validates that an input file exists and is readable
func (fm *FileManager) validateInputFile(path string) error {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &types.SecretsError{
			Operation: "file_not_found",
			Cause:     err,
			Context: map[string]interface{}{
				"path": path,
			},
		}
	}

	// Check if file is readable
	file, err := os.Open(path)
	if err != nil {
		return &types.SecretsError{
			Operation: "file_not_readable",
			Cause:     err,
			Context: map[string]interface{}{
				"path": path,
			},
		}
	}
	file.Close()

	return nil
}

// GetFileInfo returns information about a file
func (fm *FileManager) GetFileInfo(path string) (*FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, &types.SecretsError{
			Operation: "get_file_info",
			Cause:     err,
			Context: map[string]interface{}{
				"path": path,
			},
		}
	}

	return &FileInfo{
		Path:        path,
		Size:        info.Size(),
		Mode:        info.Mode(),
		ModTime:     info.ModTime(),
		IsDirectory: info.IsDir(),
		IsEncrypted: fm.isEncryptedFile(path),
	}, nil
}

// IsEncryptedFile checks if a file appears to be encrypted
func (fm *FileManager) isEncryptedFile(path string) bool {
	// Check if file has .age extension
	if filepath.Ext(path) == ".age" {
		return true
	}

	// Check if file starts with age armor header
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	// Read first few bytes to check for age armor header
	buffer := make([]byte, 32)
	n, err := file.Read(buffer)
	if err != nil || n < 32 {
		return false
	}

	// Check for age armor header
	header := string(buffer[:32])
	return len(header) >= 32 && header[:32] == "-----BEGIN AGE ENCRYPTED FILE-----"
}

// FileInfo contains information about a file
type FileInfo struct {
	Path        string
	Size        int64
	Mode        os.FileMode
	ModTime     time.Time
	IsDirectory bool
	IsEncrypted bool
}
