package keygen

import (
	"fmt"
	"os"
	"time"
)

// WriteKeyFiles writes the private key, public key, and password to files
func WriteKeyFiles(privateKeyPEM []byte, publicKey []byte, password string) (string, error) {
	// Create output directory with timestamp and counter
	now := time.Now()
	dateStr := now.Format("20060102")

	// Find the next available counter
	var counter int
	var outputDir string
	for counter = 0; counter < MaxKeyDirectories; counter++ {
		outputDir = fmt.Sprintf("generated_keys/%s%03d", dateStr, counter)
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			break
		}
	}

	if counter >= MaxKeyDirectories {
		return "", fmt.Errorf("too many key directories for today, maximum %d reached", MaxKeyDirectories)
	}

	// Create the directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write private key file
	privateKeyPath := fmt.Sprintf("%s/id_ed25519", outputDir)
	if err := os.WriteFile(privateKeyPath, privateKeyPEM, 0600); err != nil {
		return "", fmt.Errorf("failed to write private key: %w", err)
	}

	// Write public key file
	publicKeyPath := fmt.Sprintf("%s/id_ed25519.pub", outputDir)
	if err := os.WriteFile(publicKeyPath, publicKey, 0644); err != nil {
		return "", fmt.Errorf("failed to write public key: %w", err)
	}

	// Write password file
	passwordPath := fmt.Sprintf("%s/ed25519_password.txt", outputDir)
	passwordContent := fmt.Sprintf("Generated password for Ed25519 key: %s\n", password)
	if err := os.WriteFile(passwordPath, []byte(passwordContent), 0600); err != nil {
		return "", fmt.Errorf("failed to write password file: %w", err)
	}

	// Write combined file with all information
	combinedPath := fmt.Sprintf("%s/ed25519_keys.txt", outputDir)
	combinedContent := fmt.Sprintf("=== ED25519 SSH Key Pair ===\n\n")
	combinedContent += "Private Key:\n"
	combinedContent += string(privateKeyPEM)
	combinedContent += "\n\nPublic Key:\n"
	combinedContent += string(publicKey)
	combinedContent += "\n\nPassword: " + password + "\n"

	if err := os.WriteFile(combinedPath, []byte(combinedContent), 0600); err != nil {
		return "", fmt.Errorf("failed to write combined file: %w", err)
	}

	return outputDir, nil
}

// WriteKeyPair writes a complete key pair to files
func WriteKeyPair(keyPair *KeyPair) (*KeyFiles, error) {
	outputDir, err := WriteKeyFiles(keyPair.PrivateKey, keyPair.PublicKey, keyPair.Password)
	if err != nil {
		return nil, err
	}

	return &KeyFiles{
		PrivateKeyPath: fmt.Sprintf("%s/id_ed25519", outputDir),
		PublicKeyPath:  fmt.Sprintf("%s/id_ed25519.pub", outputDir),
		PasswordPath:   fmt.Sprintf("%s/ed25519_password.txt", outputDir),
		CombinedPath:   fmt.Sprintf("%s/ed25519_keys.txt", outputDir),
		OutputDir:      outputDir,
	}, nil
}
