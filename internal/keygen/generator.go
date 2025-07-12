package keygen

import (
	"fmt"
)

// GenerateSSHKeys generates SSH keys and writes them to files
func GenerateSSHKeys() error {
	return GenerateSSHKeysWithType(DefaultKeyType)
}

// GenerateSSHKeysWithType generates SSH keys of the specified type
func GenerateSSHKeysWithType(keyType KeyType) error {
	fmt.Println("🔐 Ed25519 SSH Key Generator")
	fmt.Println("============================")

	// Generate password
	password := GeneratePassword()
	fmt.Printf("Generated password: %s\n", password)

	// Generate Ed25519 keys
	fmt.Println("\n📝 Generating Ed25519 keys...")
	keyPair, err := GenerateKeyPair(keyType, password)
	if err != nil {
		fmt.Printf("❌ Failed to generate Ed25519 keys: %v\n", err)
		return fmt.Errorf("failed to generate keys: %w", err)
	}

	// Write Ed25519 key files
	keyFiles, err := WriteKeyPair(keyPair)
	if err != nil {
		fmt.Printf("❌ Failed to write Ed25519 key files: %v\n", err)
		return fmt.Errorf("failed to write key files: %w", err)
	}

	fmt.Println("✅ Ed25519 keys generated successfully!")

	// Display file locations
	fmt.Println("\n📁 Generated files:")
	fmt.Printf("  %s/\n", keyFiles.OutputDir)
	fmt.Println("  ├── id_ed25519          (Ed25519 private key)")
	fmt.Println("  ├── id_ed25519.pub      (Ed25519 public key)")
	fmt.Println("  ├── ed25519_keys.txt    (Ed25519 combined file)")
	fmt.Println("  └── ed25519_password.txt (Ed25519 password)")

	fmt.Printf("\n🔑 Password: %s\n", password)
	fmt.Println("\n⚠️  IMPORTANT: Keep the password secure and backup the key files!")

	return nil
}
