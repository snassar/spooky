package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"spooky/internal/ssh"
)

func main() {
	fmt.Println("Phase 3: Comprehensive Encrypted Key Support Test")
	fmt.Println("================================================")

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "phase3-test-*")
	if err != nil {
		fmt.Printf("Failed to create temp dir: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir)

	fmt.Printf("Testing in directory: %s\n", tempDir)

	// Test Results
	var results []KeyTestResult

	// Test 1: Go-generated encrypted RSA key (DES - Legacy)
	fmt.Println("\n1. Testing Go-generated encrypted RSA key (DES)...")
	result := testGoEncryptedRSA(tempDir)
	results = append(results, result)

	// Test 2: OpenSSL-generated encrypted RSA key (AES-256)
	fmt.Println("\n2. Testing OpenSSL-generated encrypted RSA key (AES-256)...")
	result = testOpenSSLEncryptedRSA(tempDir)
	results = append(results, result)

	// Test 3: OpenSSL-generated encrypted ECDSA key (AES-256)
	fmt.Println("\n3. Testing OpenSSL-generated encrypted ECDSA key (AES-256)...")
	result = testOpenSSLEncryptedECDSA(tempDir)
	results = append(results, result)

	// Test 4: OpenSSL-generated encrypted Ed25519 key (AES-256)
	fmt.Println("\n4. Testing OpenSSL-generated encrypted Ed25519 key (AES-256)...")
	result = testOpenSSLEncryptedEd25519(tempDir)
	results = append(results, result)

	// Print comprehensive results
	printPhase3Results(results)

	fmt.Println("\nPhase 3 testing completed!")
	fmt.Printf("All test files are in: %s\n", tempDir)
}

// KeyTestResult represents the result of testing a key
type KeyTestResult struct {
	TestName    string
	KeyType     string
	Format      string
	Encryption  string
	Passphrase  string
	SSHLoadable bool
	Error       string
	Notes       string
}

// testGoEncryptedRSA tests Go-generated encrypted RSA keys (legacy DES)
func testGoEncryptedRSA(tempDir string) KeyTestResult {
	result := KeyTestResult{
		TestName:   "Go Encrypted RSA (DES)",
		KeyType:    "RSA-2048",
		Format:     "Traditional",
		Encryption: "DES",
		Passphrase: "test-passphrase-123",
	}

	// Generate RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to generate RSA key: %v", err)
		return result
	}

	// Create encrypted PEM block using deprecated function
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	encryptedBlock, err := x509.EncryptPEMBlock(rand.Reader, block.Type, block.Bytes, []byte(result.Passphrase), x509.PEMCipherDES)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to encrypt RSA key: %v", err)
		return result
	}

	// Write to file
	keyPath := filepath.Join(tempDir, "go_rsa_encrypted_des.pem")
	pemData := pem.EncodeToMemory(encryptedBlock)
	if err := os.WriteFile(keyPath, pemData, 0600); err != nil {
		result.Error = fmt.Sprintf("Failed to write encrypted RSA key: %v", err)
		return result
	}

	// Test with SSH client
	config := &ssh.SSHConfig{
		PrivateKeyPath: keyPath,
		Passphrase:     result.Passphrase,
	}

	_, err = ssh.NewSSHClient(config)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to create SSH client: %v", err)
		return result
	}

	// For this test, we'll just verify the client was created successfully
	// since we can't access the private loadPrivateKey method
	result.SSHLoadable = true
	result.Notes = "SSH client created successfully with encrypted key"

	return result
}

// testOpenSSLEncryptedRSA tests OpenSSL-generated encrypted RSA keys (AES-256)
func testOpenSSLEncryptedRSA(tempDir string) KeyTestResult {
	result := KeyTestResult{
		TestName:   "OpenSSL Encrypted RSA (AES-256)",
		KeyType:    "RSA-2048",
		Format:     "OpenSSL",
		Encryption: "AES-256",
		Passphrase: "test-passphrase-123",
	}

	// Generate unencrypted RSA key with OpenSSL
	unencryptedPath := filepath.Join(tempDir, "openssl_rsa_2048.pem")
	cmd := exec.Command("openssl", "genrsa", "-out", unencryptedPath, "2048")
	if err := cmd.Run(); err != nil {
		result.Error = fmt.Sprintf("Failed to generate OpenSSL RSA key: %v", err)
		return result
	}

	// Encrypt with AES-256
	encryptedPath := filepath.Join(tempDir, "openssl_rsa_encrypted_aes256.pem")
	cmd = exec.Command("openssl", "rsa", "-aes256", "-in", unencryptedPath, "-out", encryptedPath, "-passout", "pass:"+result.Passphrase)
	if err := cmd.Run(); err != nil {
		result.Error = fmt.Sprintf("Failed to encrypt OpenSSL RSA key: %v", err)
		return result
	}

	// Test with SSH client
	config := &ssh.SSHConfig{
		PrivateKeyPath: encryptedPath,
		Passphrase:     result.Passphrase,
	}

	_, err := ssh.NewSSHClient(config)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to create SSH client: %v", err)
		return result
	}

	// Verify the key file exists and has correct format
	keyData, err := os.ReadFile(encryptedPath)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to read encrypted key: %v", err)
		return result
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		result.Error = "Failed to decode PEM block"
		return result
	}

	if block.Type != "ENCRYPTED PRIVATE KEY" {
		result.Error = fmt.Sprintf("Unexpected PEM type: %s (expected ENCRYPTED PRIVATE KEY)", block.Type)
		return result
	}

	result.SSHLoadable = true
	result.Notes = fmt.Sprintf("SSH client created successfully. PEM type: %s, Size: %d bytes", block.Type, len(block.Bytes))

	return result
}

// testOpenSSLEncryptedECDSA tests OpenSSL-generated encrypted ECDSA keys (AES-256)
func testOpenSSLEncryptedECDSA(tempDir string) KeyTestResult {
	result := KeyTestResult{
		TestName:   "OpenSSL Encrypted ECDSA (AES-256)",
		KeyType:    "ECDSA-P-256",
		Format:     "PKCS#8",
		Encryption: "AES-256",
		Passphrase: "test-passphrase-123",
	}

	// Generate ECDSA key
	ecdsaPath := filepath.Join(tempDir, "openssl_ecdsa_p256.pem")
	cmd := exec.Command("openssl", "ecparam", "-genkey", "-name", "prime256v1", "-out", ecdsaPath)
	if err := cmd.Run(); err != nil {
		result.Error = fmt.Sprintf("Failed to generate ECDSA key: %v", err)
		return result
	}

	// Convert to PKCS#8
	pkcs8Path := ecdsaPath + "_pkcs8.pem"
	cmd = exec.Command("openssl", "pkcs8", "-topk8", "-nocrypt", "-in", ecdsaPath, "-out", pkcs8Path)
	if err := cmd.Run(); err != nil {
		result.Error = fmt.Sprintf("Failed to convert to PKCS#8: %v", err)
		return result
	}

	// Encrypt PKCS#8 with AES-256
	encryptedPath := filepath.Join(tempDir, "openssl_ecdsa_encrypted_aes256.pem")
	cmd = exec.Command("openssl", "pkcs8", "-topk8", "-v2", "aes256", "-in", pkcs8Path, "-out", encryptedPath, "-passout", "pass:"+result.Passphrase)
	if err := cmd.Run(); err != nil {
		result.Error = fmt.Sprintf("Failed to encrypt ECDSA key: %v", err)
		return result
	}

	// Test with SSH client
	config := &ssh.SSHConfig{
		PrivateKeyPath: encryptedPath,
		Passphrase:     result.Passphrase,
	}

	_, err := ssh.NewSSHClient(config)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to create SSH client: %v", err)
		return result
	}

	// Verify the key file exists and has correct format
	keyData, err := os.ReadFile(encryptedPath)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to read encrypted key: %v", err)
		return result
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		result.Error = "Failed to decode PEM block"
		return result
	}

	if block.Type != "ENCRYPTED PRIVATE KEY" {
		result.Error = fmt.Sprintf("Unexpected PEM type: %s (expected ENCRYPTED PRIVATE KEY)", block.Type)
		return result
	}

	result.SSHLoadable = true
	result.Notes = fmt.Sprintf("SSH client created successfully. PEM type: %s, Size: %d bytes", block.Type, len(block.Bytes))

	return result
}

// testOpenSSLEncryptedEd25519 tests OpenSSL-generated encrypted Ed25519 keys (AES-256)
func testOpenSSLEncryptedEd25519(tempDir string) KeyTestResult {
	result := KeyTestResult{
		TestName:   "OpenSSL Encrypted Ed25519 (AES-256)",
		KeyType:    "Ed25519",
		Format:     "PKCS#8",
		Encryption: "AES-256",
		Passphrase: "test-passphrase-123",
	}

	// Generate Ed25519 key
	ed25519Path := filepath.Join(tempDir, "openssl_ed25519.pem")
	cmd := exec.Command("openssl", "genpkey", "-algorithm", "ED25519", "-out", ed25519Path)
	if err := cmd.Run(); err != nil {
		result.Error = fmt.Sprintf("Failed to generate Ed25519 key: %v", err)
		return result
	}

	// Encrypt with AES-256
	encryptedPath := filepath.Join(tempDir, "openssl_ed25519_encrypted_aes256.pem")
	cmd = exec.Command("openssl", "pkcs8", "-topk8", "-v2", "aes256", "-in", ed25519Path, "-out", encryptedPath, "-passout", "pass:"+result.Passphrase)
	if err := cmd.Run(); err != nil {
		result.Error = fmt.Sprintf("Failed to encrypt Ed25519 key: %v", err)
		return result
	}

	// Test with SSH client
	config := &ssh.SSHConfig{
		PrivateKeyPath: encryptedPath,
		Passphrase:     result.Passphrase,
	}

	_, err := ssh.NewSSHClient(config)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to create SSH client: %v", err)
		return result
	}

	// Verify the key file exists and has correct format
	keyData, err := os.ReadFile(encryptedPath)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to read encrypted key: %v", err)
		return result
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		result.Error = "Failed to decode PEM block"
		return result
	}

	if block.Type != "ENCRYPTED PRIVATE KEY" {
		result.Error = fmt.Sprintf("Unexpected PEM type: %s (expected ENCRYPTED PRIVATE KEY)", block.Type)
		return result
	}

	result.SSHLoadable = true
	result.Notes = fmt.Sprintf("SSH client created successfully. PEM type: %s, Size: %d bytes", block.Type, len(block.Bytes))

	return result
}

// printPhase3Results prints the comprehensive test results
func printPhase3Results(results []KeyTestResult) {
	fmt.Println("\n" + strings.Repeat("=", 100))
	fmt.Println("Phase 3: Comprehensive Encrypted Key Support Test Results")
	fmt.Println(strings.Repeat("=", 100))

	fmt.Printf("%-35s %-15s %-15s %-15s %-15s %-8s\n",
		"Test Name", "Key Type", "Format", "Encryption", "Passphrase", "SSH OK")
	fmt.Println(strings.Repeat("-", 100))

	for _, result := range results {
		sshOK := "No"
		if result.SSHLoadable {
			sshOK = "Yes"
		}

		passphrase := "None"
		if result.Passphrase != "" {
			passphrase = "Set"
		}

		fmt.Printf("%-35s %-15s %-15s %-15s %-15s %-8s\n",
			result.TestName, result.KeyType, result.Format, result.Encryption, passphrase, sshOK)

		if result.Error != "" {
			fmt.Printf("  └─ Error: %s\n", result.Error)
		}

		if result.Notes != "" {
			fmt.Printf("  └─ Notes: %s\n", result.Notes)
		}
	}

	fmt.Println(strings.Repeat("-", 100))

	// Summary
	total := len(results)
	successful := 0
	for _, result := range results {
		if result.SSHLoadable {
			successful++
		}
	}

	fmt.Printf("Summary: %d/%d encrypted key types successfully supported by SSH client\n", successful, total)
	fmt.Printf("Success Rate: %.1f%%\n", float64(successful)/float64(total)*100)

	// Analysis
	fmt.Println("\nAnalysis:")
	if successful == total {
		fmt.Println("✅ All encrypted key types are supported!")
		fmt.Println("🎯 Ready to remove deprecated functions")
	} else if successful > 0 {
		fmt.Printf("⚠️  Partial support: %d/%d encrypted key types working\n", successful, total)
		fmt.Println("🔧 Some encrypted key types still need work")
	} else {
		fmt.Println("❌ No encrypted key types are supported")
		fmt.Println("🚨 Encrypted key support needs significant work")
	}
}
