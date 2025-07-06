package tests

import (
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"encoding/pem"
	"fmt"
	"io"
	"math/rand"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// ServerTracker tracks running server processes for cleanup
type ServerTracker struct {
	processes map[int]*exec.Cmd
	mutex     sync.RWMutex
}

var (
	serverTracker = &ServerTracker{
		processes: make(map[int]*exec.Cmd),
	}
)

// AddProcess adds a server process to the tracker
func (st *ServerTracker) AddProcess(pid int, cmd *exec.Cmd) {
	st.mutex.Lock()
	defer st.mutex.Unlock()
	st.processes[pid] = cmd
}

// RemoveProcess removes a server process from the tracker
func (st *ServerTracker) RemoveProcess(pid int) {
	st.mutex.Lock()
	defer st.mutex.Unlock()
	delete(st.processes, pid)
}

// CleanupAll kills all tracked server processes
func (st *ServerTracker) CleanupAll() {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	for pid, cmd := range st.processes {
		if cmd.Process != nil {
			fmt.Printf("Cleaning up server process PID %d\n", pid)
			cmd.Process.Kill()
			// Don't wait for the process to exit since it might hang
			// The OS will clean up the process
		}
	}
	st.processes = make(map[int]*exec.Cmd)
}

// CleanupServers is a function that can be called to clean up all server processes
func CleanupServers() {
	serverTracker.CleanupAll()
}

// connectSSH establishes an SSH connection without authentication
func connectSSH(addr, user string, privateKey []byte) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(func() ([]ssh.Signer, error) {
				if privateKey == nil {
					// For servers that don't require authentication
					return []ssh.Signer{}, nil
				}
				signer, err := ssh.ParsePrivateKey(privateKey)
				if err != nil {
					return nil, err
				}
				return []ssh.Signer{signer}, nil
			}),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	return ssh.Dial("tcp", addr, config)
}

// connectSSHWithKey establishes an SSH connection with a private key
func connectSSHWithKey(addr, user string, privateKey []byte) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(func() ([]ssh.Signer, error) {
				signer, err := ssh.ParsePrivateKey(privateKey)
				if err != nil {
					return nil, err
				}
				return []ssh.Signer{signer}, nil
			}),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	return ssh.Dial("tcp", addr, config)
}

// connectSSHWithPassword establishes an SSH connection with password authentication
func connectSSHWithPassword(addr, user, password string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	return ssh.Dial("tcp", addr, config)
}

// generateTestKeyPair generates a test Ed25519 key pair
func generateTestKeyPair() (privateKey []byte, publicKey []byte, err error) {
	// Generate Ed25519 key pair
	pub, priv, err := ed25519.GenerateKey(cryptorand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate key pair: %v", err)
	}

	// Convert private key to PEM format
	privateKeyPEM := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: append(priv.Seed(), pub...),
	}
	privateKey = pem.EncodeToMemory(privateKeyPEM)

	// Convert public key to SSH format
	sshPubKey, err := ssh.NewPublicKey(pub)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create SSH public key: %v", err)
	}
	publicKey = ssh.MarshalAuthorizedKey(sshPubKey)

	return privateKey, publicKey, nil
}

// generateTestKeyPairFromMain generates a test Ed25519 key pair using the same logic as ssh_keygen.go
func generateTestKeyPairFromMain() (privateKey []byte, publicKey []byte, err error) {
	// Generate a test password (we don't need to store it for testing)
	password := generateTestPassword()

	// Use the same key generation logic as ssh_keygen.go
	privateKeyPEM, publicKeyBytes, err := generateEd25519KeyPairFromMain(password)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate key pair using main logic: %v", err)
	}

	return privateKeyPEM, publicKeyBytes, nil
}

// generateTestPassword creates a 25-character password with ASCII letters and numbers
func generateTestPassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	password := make([]byte, 25)
	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}
	return string(password)
}

// generateEd25519KeyPairFromMain generates an Ed25519 key pair using the same logic as ssh_keygen.go
func generateEd25519KeyPairFromMain(password string) (privateKeyPEM []byte, publicKeyBytes []byte, err error) {
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

	// Create the private key in OpenSSH format
	privateKeyBlock, err := ssh.MarshalPrivateKey(privateKeyRaw, password)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal private key: %w", err)
	}

	// Encode to PEM format
	privateKeyPEM = pem.EncodeToMemory(privateKeyBlock)

	return privateKeyPEM, ssh.MarshalAuthorizedKey(publicKey), nil
}

// connectSFTP creates an SFTP client from an SSH client
func connectSFTP(sshClient *ssh.Client) (*sftp.Client, error) {
	return sftp.NewClient(sshClient)
}

// testSFTPOperations performs basic SFTP operations for testing
func testSFTPOperations(t *testing.T, sftpClient *sftp.Client) {
	// Test creating a directory - use a unique path that works on Windows
	testDir := fmt.Sprintf("./test_sftp_dir_%d", time.Now().Unix())
	err := sftpClient.Mkdir(testDir)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Test creating a file
	testFile := filepath.Join(testDir, "test.txt")
	testContent := "Hello SFTP!"

	file, err := sftpClient.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.Write([]byte(testContent))
	if err != nil {
		t.Fatalf("Failed to write to file: %v", err)
	}

	// Test reading the file
	file, err = sftpClient.Open(testFile)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("Expected content %q, got %q", testContent, string(content))
	}

	// Test listing directory
	files, err := sftpClient.ReadDir(testDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	found := false
	for _, f := range files {
		if f.Name() == "test.txt" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find test.txt in directory listing")
	}

	// Clean up
	sftpClient.Remove(testFile)
	sftpClient.RemoveDirectory(testDir)
}
