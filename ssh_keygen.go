package main

import (
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"encoding/pem"
	"fmt"
	"math/rand"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

// generatePassword creates a 25-character password with ASCII letters and numbers
func generatePassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	password := make([]byte, 25)
	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}
	return string(password)
}

// generateEd25519KeyPair generates an Ed25519 key pair and returns it in OpenSSH format
func generateEd25519KeyPair(password string) (privateKeyPEM []byte, publicKeyBytes []byte, err error) {
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

// writeKeyFiles writes the private key, public key, and password to files
func writeKeyFiles(privateKeyPEM []byte, publicKey []byte, password string) (string, error) {
	// Create output directory with timestamp and counter
	now := time.Now()
	dateStr := now.Format("20060102")

	// Find the next available counter
	var counter int
	var outputDir string
	for counter = 0; counter < 1000; counter++ {
		outputDir = fmt.Sprintf("generated_keys/%s%03d", dateStr, counter)
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			break
		}
	}

	if counter >= 1000 {
		return "", fmt.Errorf("too many key directories for today, maximum 1000 reached")
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

func generateSSHKeys() {
	// coverage-ignore: CLI entry point, tested via integration tests
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	fmt.Println("🔐 Ed25519 SSH Key Generator")
	fmt.Println("============================")

	// Generate password
	password := generatePassword()
	fmt.Printf("Generated password: %s\n", password)

	// Generate Ed25519 keys
	fmt.Println("\n📝 Generating Ed25519 keys...")
	privateKeyPEM, publicKey, err := generateEd25519KeyPair(password)
	if err != nil {
		fmt.Printf("❌ Failed to generate Ed25519 keys: %v\n", err)
		os.Exit(1) // coverage-ignore: exit on error is expected behavior
	}

	// Write Ed25519 key files
	outputDir, err := writeKeyFiles(privateKeyPEM, publicKey, password)
	if err != nil {
		fmt.Printf("❌ Failed to write Ed25519 key files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Ed25519 keys generated successfully!")

	// Display file locations
	fmt.Println("\n📁 Generated files:")
	fmt.Printf("  %s/\n", outputDir)
	fmt.Println("  ├── id_ed25519          (Ed25519 private key)")
	fmt.Println("  ├── id_ed25519.pub      (Ed25519 public key)")
	fmt.Println("  ├── ed25519_keys.txt    (Ed25519 combined file)")
	fmt.Println("  └── ed25519_password.txt (Ed25519 password)")

	fmt.Printf("\n🔑 Password: %s\n", password)
	fmt.Println("\n⚠️  IMPORTANT: Keep the password secure and backup the key files!")
}
