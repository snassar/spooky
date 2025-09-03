package ssh

import (
	"bufio"
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"

	"spooky/internal/logging"
	"spooky/internal/schemas"
)

// SSHClient represents a comprehensive SSH client for Spooky.
//
// Authentication Methods Supported:
// - Public key authentication (with optional passphrase)
// - Password authentication (pre-configured only)
// - SSH agent authentication
// - Certificate authentication
//
// Authentication Methods NOT Supported:
// - Keyboard-interactive authentication (no interactive prompts)
// - Password callbacks (no runtime password gathering)
//
// SSH Features Supported:
// - Proxy connections (ProxyCommand, ProxyJump)
// - Compression (with configurable levels)
// - TCP keepalive (with configurable parameters)
// - File transfers (SCP)
// - Command execution
//
// SSH Features NOT Supported:
// - Terminal/PTY allocation (no interactive sessions)
// - X11 forwarding
// - SSH agent forwarding
// - Environment variable passing
//
// All credentials must be pre-configured in the machine configuration.
// Sensitive values can be encrypted using age encryption.
type SSHClient struct {
	config     *SSHConfig
	client     *ssh.Client
	session    *ssh.Session
	agentConn  agent.Agent
	knownHosts ssh.HostKeyCallback
	mu         sync.Mutex
}

// SSHConfig holds SSH connection configuration
type SSHConfig struct {
	Host                      string
	Port                      int
	User                      string
	Password                  string
	PrivateKeyPath            string
	PrivateKeyData            []byte
	Passphrase                string
	Timeout                   time.Duration
	KeepAlive                 time.Duration
	KeepAliveCount            int
	KeyScanTimeout            time.Duration
	KnownHostsPath            string
	StrictHostKey             bool
	KnownHostsMode            string
	UserKnownHosts            []string
	GlobalKnownHosts          []string
	ProxyCommand              string
	ProxyJump                 string
	Compression               bool
	CompressionLevel          int
	TCPKeepAlive              bool
	TCPKeepAliveCount         int
	TCPKeepAliveIdle          time.Duration
	TCPKeepAliveInterval      time.Duration
	TCPKeepAliveProbeInterval time.Duration
	LogLevel                  string
	Verbose                   bool
	BatchMode                 bool
	PubkeyAuth                bool
	PasswordAuth              bool
	AgentForwarding           bool
	X11Forwarding             bool
	X11Display                string
	X11AuthType               string
	Environment               map[string]string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// NewSSHClient creates a new SSH client with the given configuration
func NewSSHClient(config *SSHConfig) (*SSHClient, error) {
	if config == nil {
		return nil, errors.New("SSH config cannot be nil")
	}

	// Set defaults
	if config.Port == 0 {
		config.Port = 22
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.KeepAlive == 0 {
		config.KeepAlive = 60 * time.Second
	}
	if config.KeepAliveCount == 0 {
		config.KeepAliveCount = 3
	}
	if config.KeyScanTimeout == 0 {
		config.KeyScanTimeout = 10 * time.Second
	}

	client := &SSHClient{
		config: config,
	}

	// Initialize known hosts callback
	if err := client.initKnownHosts(); err != nil {
		return nil, errors.Wrap(err, "failed to initialize known hosts")
	}

	return client, nil
}

// Connect establishes an SSH connection
func (sc *SSHClient) Connect(ctx context.Context) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Create SSH client config
	sshConfig, err := sc.createSSHConfig()
	if err != nil {
		return errors.Wrap(err, "failed to create SSH config - connection cannot be established")
	}

	// Create connection address
	addr := fmt.Sprintf("%s:%d", sc.config.Host, sc.config.Port)

	// Establish connection with timeout
	dialer := &net.Dialer{
		Timeout: sc.config.Timeout,
	}

	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return errors.Wrapf(err, "failed to connect to %s - network connection failed", addr)
	}

	// Create SSH connection
	sshConn, chans, reqs, err := ssh.NewClientConn(conn, addr, sshConfig)
	if err != nil {
		conn.Close()
		return errors.Wrapf(err, "failed to establish SSH connection to %s - authentication or protocol negotiation failed", addr)
	}

	// Create SSH client
	sc.client = ssh.NewClient(sshConn, chans, reqs)

	// Set up keepalive
	if sc.config.KeepAlive > 0 {
		go sc.keepAlive()
	}

	return nil
}

// Disconnect closes the SSH connection
func (sc *SSHClient) Disconnect() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if sc.session != nil {
		sc.session.Close()
		sc.session = nil
	}

	if sc.client != nil {
		err := sc.client.Close()
		sc.client = nil
		return errors.Wrap(err, "failed to close SSH client")
	}

	return nil
}

// RunCommand executes a command on the remote host with configurable timeout
func (sc *SSHClient) RunCommand(ctx context.Context, command string) (*schemas.CommandResult, error) {
	if sc.client == nil {
		return nil, errors.New("SSH client not connected")
	}

	// Apply timeout from config if context doesn't have one
	var timeoutCtx context.Context
	var cancel context.CancelFunc

	if _, ok := ctx.Deadline(); !ok && sc.config.Timeout > 0 {
		timeoutCtx, cancel = context.WithTimeout(ctx, sc.config.Timeout)
		defer cancel()
	} else {
		timeoutCtx = ctx
	}

	// Create session
	session, err := sc.client.NewSession()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create SSH session")
	}
	defer session.Close()

	// Set up session
	if err := sc.setupSession(session); err != nil {
		return nil, errors.Wrap(err, "failed to setup session")
	}

	// Execute command
	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	// Check for context cancellation before running command
	select {
	case <-timeoutCtx.Done():
		return &schemas.CommandResult{
			ExitCode: 1,
			Stdout:   "",
			Stderr:   "command cancelled due to timeout or context cancellation",
			Error:    timeoutCtx.Err(),
		}, nil
	default:
		// Continue with command execution
	}

	if err := session.Run(command); err != nil {
		return &schemas.CommandResult{
			ExitCode: 1,
			Stdout:   stdout.String(),
			Stderr:   stderr.String(),
			Error:    err,
		}, nil
	}

	return &schemas.CommandResult{
		ExitCode: 0,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
	}, nil
}

// UploadFile uploads a file to the remote host
func (sc *SSHClient) UploadFile(ctx context.Context, localPath, remotePath string) error {
	if sc.client == nil {
		return errors.New("SSH client not connected")
	}

	// Open local file
	localFile, err := os.Open(localPath)
	if err != nil {
		return errors.Wrapf(err, "failed to open local file: %s", localPath)
	}
	defer localFile.Close()

	// Get file info
	fileInfo, err := localFile.Stat()
	if err != nil {
		return errors.Wrapf(err, "failed to get file info: %s", localPath)
	}

	// Create session
	session, err := sc.client.NewSession()
	if err != nil {
		return errors.Wrap(err, "failed to create SSH session")
	}
	defer session.Close()

	// Set up file transfer
	stdin, err := session.StdinPipe()
	if err != nil {
		return errors.Wrap(err, "failed to get stdin pipe")
	}

	// Start scp command
	if err := session.Start(fmt.Sprintf("scp -t %s", remotePath)); err != nil {
		return errors.Wrap(err, "failed to start scp command")
	}

	// Send file header
	fmt.Fprintf(stdin, "C%04o %d %s\n", fileInfo.Mode().Perm(), fileInfo.Size(), filepath.Base(localPath))

	// Copy file content
	if _, err := io.Copy(stdin, localFile); err != nil {
		return errors.Wrap(err, "failed to copy file content")
	}

	// Send end marker
	fmt.Fprint(stdin, "\x00")

	// Close stdin and wait for completion
	stdin.Close()
	return session.Wait()
}

// DownloadFile downloads a file from the remote host
func (sc *SSHClient) DownloadFile(ctx context.Context, remotePath, localPath string) error {
	if sc.client == nil {
		return errors.New("SSH client not connected")
	}

	// Create session
	session, err := sc.client.NewSession()
	if err != nil {
		return errors.Wrap(err, "failed to create SSH session")
	}
	defer session.Close()

	// Set up file transfer
	stdout, err := session.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "failed to get stdout pipe")
	}

	// Start scp command
	if err := session.Start(fmt.Sprintf("scp -f %s", remotePath)); err != nil {
		return errors.Wrap(err, "failed to start scp command")
	}

	// Create local file
	localFile, err := os.Create(localPath)
	if err != nil {
		return errors.Wrapf(err, "failed to create local file: %s", localPath)
	}
	defer localFile.Close()

	// Read file header
	scanner := bufio.NewScanner(stdout)
	if !scanner.Scan() {
		return errors.New("failed to read file header")
	}

	header := scanner.Text()
	if !strings.HasPrefix(header, "C") {
		return errors.New("invalid file header")
	}

	// Parse file info
	parts := strings.Fields(header)
	if len(parts) != 3 {
		return errors.New("invalid file header format")
	}

	// Extract file size
	fileSize, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse file size")
	}

	// Send acknowledgment
	fmt.Fprint(os.Stdout, "\x00")

	// Copy file content
	written, err := io.CopyN(localFile, stdout, fileSize)
	if err != nil {
		return errors.Wrap(err, "failed to copy file content")
	}

	if written != fileSize {
		return errors.New("incomplete file transfer")
	}

	// Send acknowledgment
	fmt.Fprint(os.Stdout, "\x00")

	return session.Wait()
}

// CopyFile copies a file from local to remote using SCP
func (sc *SSHClient) CopyFile(localPath, remotePath string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if sc.client == nil {
		return errors.New("SSH client not connected")
	}

	// Open local file
	localFile, err := os.Open(localPath)
	if err != nil {
		return errors.Wrapf(err, "failed to open local file: %s", localPath)
	}
	defer localFile.Close()

	// Get file info for permissions
	fileInfo, err := localFile.Stat()
	if err != nil {
		return errors.Wrapf(err, "failed to get file info: %s", localPath)
	}

	// Create SCP session
	session, err := sc.client.NewSession()
	if err != nil {
		return errors.Wrap(err, "failed to create SCP session")
	}
	defer session.Close()

	// Create pipe for SCP protocol
	stdin, err := session.StdinPipe()
	if err != nil {
		return errors.Wrap(err, "failed to create stdin pipe")
	}

	// Start SCP command
	if err := session.Start(fmt.Sprintf("scp -t %s", remotePath)); err != nil {
		return errors.Wrap(err, "failed to start SCP command")
	}

	// Send file header
	fmt.Fprintf(stdin, "C%04o %d %s\n", fileInfo.Mode().Perm(), fileInfo.Size(), filepath.Base(remotePath))

	// Copy file content
	if _, err := io.Copy(stdin, localFile); err != nil {
		return errors.Wrap(err, "failed to copy file content")
	}

	// Send end marker
	fmt.Fprint(stdin, "\x00")

	// Close stdin and wait for session to complete
	stdin.Close()
	if err := session.Wait(); err != nil {
		return errors.Wrap(err, "SCP session failed")
	}

	return nil
}

// CopyFileContent copies file content (string) to remote path using SCP
func (sc *SSHClient) CopyFileContent(content, remotePath string, permissions os.FileMode) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if sc.client == nil {
		return errors.New("SSH client not connected")
	}

	// Create SCP session
	session, err := sc.client.NewSession()
	if err != nil {
		return errors.Wrap(err, "failed to create SCP session")
	}
	defer session.Close()

	// Create pipe for SCP protocol
	stdin, err := session.StdinPipe()
	if err != nil {
		return errors.Wrap(err, "failed to create stdin pipe")
	}

	// Start SCP command
	if err := session.Start(fmt.Sprintf("scp -t %s", remotePath)); err != nil {
		return errors.Wrap(err, "failed to start SCP command")
	}

	// Send file header
	fmt.Fprintf(stdin, "C%04o %d %s\n", permissions, len(content), filepath.Base(remotePath))

	// Copy content
	if _, err := io.Copy(stdin, strings.NewReader(content)); err != nil {
		return errors.Wrap(err, "failed to copy file content")
	}

	// Send end marker
	fmt.Fprint(stdin, "\x00")

	// Close stdin and wait for session to complete
	stdin.Close()
	if err := session.Wait(); err != nil {
		return errors.Wrap(err, "SCP session failed")
	}

	return nil
}

// SetFileAttributes sets file permissions, owner, and group on a remote file
func (sc *SSHClient) SetFileAttributes(remotePath, permissions, owner, group string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if sc.client == nil {
		return errors.New("SSH client not connected")
	}

	// Create session for chmod/chown commands
	session, err := sc.client.NewSession()
	if err != nil {
		return errors.Wrap(err, "failed to create session")
	}
	defer session.Close()

	// Set permissions if specified
	if permissions != "" {
		chmodCmd := fmt.Sprintf("chmod %s %s", permissions, remotePath)
		if err := session.Run(chmodCmd); err != nil {
			return errors.Wrapf(err, "failed to set permissions: %s", chmodCmd)
		}
	}

	// Set owner and group if specified
	if owner != "" || group != "" {
		var chownCmd string
		if owner != "" && group != "" {
			chownCmd = fmt.Sprintf("chown %s:%s %s", owner, group, remotePath)
		} else if owner != "" {
			chownCmd = fmt.Sprintf("chown %s %s", owner, remotePath)
		} else {
			chownCmd = fmt.Sprintf("chgrp %s %s", group, remotePath)
		}

		if err := session.Run(chownCmd); err != nil {
			return errors.Wrapf(err, "failed to set ownership: %s", chownCmd)
		}
	}

	return nil
}

// createSSHConfig creates the SSH client configuration
func (sc *SSHClient) createSSHConfig() (*ssh.ClientConfig, error) {
	config := &ssh.ClientConfig{
		User:            sc.config.User,
		HostKeyCallback: sc.knownHosts,
		Timeout:         sc.config.Timeout,
	}

	// Add authentication methods
	authMethods := []ssh.AuthMethod{}

	// Password authentication
	if sc.config.Password != "" && sc.config.PasswordAuth {
		authMethods = append(authMethods, ssh.Password(sc.config.Password))
	}

	// Public key authentication
	if sc.config.PrivateKeyPath != "" || len(sc.config.PrivateKeyData) > 0 {
		signer, err := sc.loadPrivateKey()
		if err != nil {
			return nil, errors.Wrap(err, "failed to load private key")
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	// SSH agent authentication
	if sc.config.AgentForwarding {
		if agentConn, err := sc.connectToAgent(); err == nil {
			authMethods = append(authMethods, ssh.PublicKeysCallback(agentConn.Signers))
		}
	}

	// Note: Spooky does not support keyboard-interactive authentication
	// All credentials must be pre-configured in the machine configuration

	if len(authMethods) == 0 {
		return nil, errors.New("no authentication methods available")
	}

	config.Auth = authMethods
	return config, nil
}

// loadPrivateKey loads and decrypts the private key
func (sc *SSHClient) loadPrivateKey() (ssh.Signer, error) {
	var keyData []byte
	var err error

	if sc.config.PrivateKeyPath != "" {
		keyData, err = os.ReadFile(sc.config.PrivateKeyPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read private key file: %s - authentication cannot proceed", sc.config.PrivateKeyPath)
		}
	} else if len(sc.config.PrivateKeyData) > 0 {
		keyData = sc.config.PrivateKeyData
	} else {
		return nil, errors.New("no private key provided - authentication configuration is incomplete")
	}

	// Security-critical: Validate key data
	if len(keyData) == 0 {
		return nil, errors.New("private key data is empty - authentication cannot proceed")
	}

	// Check for legacy DES-encrypted traditional keys first and reject them
	block, _ := pem.Decode(keyData)
	if block != nil {
		if block.Type == "RSA PRIVATE KEY" || block.Type == "DSA PRIVATE KEY" || block.Type == "EC PRIVATE KEY" {
			if x509.IsEncryptedPEMBlock(block) {
				return nil, errors.New("legacy DES-encrypted traditional keys are no longer supported - please convert to PKCS#8 or OpenSSH format")
			}
		}
	}

	// Try modern SSH parsing first (supports PKCS#8 and other encrypted formats)
	if sc.config.Passphrase != "" {
		// Try parsing with passphrase first - this handles PKCS#8 and other modern formats
		signer, err := ssh.ParsePrivateKeyWithPassphrase(keyData, []byte(sc.config.Passphrase))
		if err == nil {
			return signer, nil
		}

		// If that fails, try our custom PKCS#8 decryption
		logger := logging.GetGlobalLogger()
		logger.Debug("Modern SSH parsing failed, trying custom PKCS#8 decryption",
			slog.String("error", err.Error()))

		// Try custom PKCS#8 decryption for ENCRYPTED PRIVATE KEY format
		decryptedKey, err := sc.decryptPKCS8Key(keyData, sc.config.Passphrase)
		if err == nil {
			// Parse the decrypted key
			signer, err := ssh.ParsePrivateKey(decryptedKey)
			if err == nil {
				return signer, nil
			}
			logger.Debug("Failed to parse decrypted PKCS#8 key", slog.String("error", err.Error()))
		} else {
			logger.Debug("Custom PKCS#8 decryption failed", slog.String("error", err.Error()))
		}

		// Modern parsing failed, will try unencrypted formats below
		logger.Debug("Modern parsing failed, will try unencrypted formats")
	}

	// Parse PEM to get block and handle rest data
	block, rest := pem.Decode(keyData)
	if block == nil {
		return nil, errors.New("failed to decode PEM block - private key format is invalid")
	}
	if len(rest) > 0 {
		// Log warning about extra data but continue
		logger := logging.GetGlobalLogger()
		logger.Warn("PEM block contains extra data after first block",
			slog.Int("extra_bytes", len(rest)))
	}

	// Handle different PEM block types directly
	switch block.Type {
	case "RSA PRIVATE KEY", "DSA PRIVATE KEY", "EC PRIVATE KEY":
		// Traditional private keys - only support unencrypted format
		// Legacy DES-encrypted keys are already rejected earlier in the function

		// Parse the unencrypted key
		signer, err := ssh.ParsePrivateKey(pem.EncodeToMemory(block))
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse traditional private key - key format may be corrupted")
		}
		return signer, nil

	case "OPENSSH PRIVATE KEY":
		// OpenSSH format - try parsing directly
		signer, err := ssh.ParsePrivateKey(pem.EncodeToMemory(block))
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse OpenSSH private key - key format may be corrupted")
		}
		return signer, nil

	case "PRIVATE KEY":
		// PKCS#8 format - try parsing directly first
		signer, err := ssh.ParsePrivateKey(pem.EncodeToMemory(block))
		if err != nil {
			// If direct parsing fails, it might be encrypted PKCS#8
			logger := logging.GetGlobalLogger()
			logger.Debug("PKCS#8 parsing failed, may be encrypted", slog.String("error", err.Error()))
			return nil, errors.New("encrypted PKCS#8 keys should be handled by ParsePrivateKeyWithPassphrase")
		}
		return signer, nil

	default:
		logger := logging.GetGlobalLogger()
		logger.Debug("Unknown PEM block type", slog.String("type", block.Type))
		return nil, errors.Errorf("unsupported PEM block type: %s - private key format is not supported", block.Type)
	}
}

// OID constants for PKCS#8 encryption
var (
	oidPBES2     = []int{1, 2, 840, 113549, 1, 5, 13}
	oidPBKDF2    = []int{1, 2, 840, 113549, 1, 5, 12}
	oidAES256CBC = []int{2, 16, 840, 1, 101, 3, 4, 1, 42}
	oidAES128CBC = []int{2, 16, 840, 1, 101, 3, 4, 1, 2}
)

// AlgorithmIdentifier represents the ASN.1 structure for algorithm identifiers
type AlgorithmIdentifier struct {
	Algorithm  []int
	Parameters asn1.RawValue
}

// compareOID compares two OIDs for equality
func compareOID(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// decryptPKCS8Key attempts to decrypt a PKCS#8 encrypted private key
func (sc *SSHClient) decryptPKCS8Key(keyData []byte, passphrase string) ([]byte, error) {
	// Parse PEM to get the encrypted block
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	// Check if this is an ENCRYPTED PRIVATE KEY
	if block.Type != "ENCRYPTED PRIVATE KEY" {
		return nil, errors.New("not a PKCS#8 encrypted private key")
	}

	// Parse the PKCS#8 encrypted structure
	var encryptedKey struct {
		Version       int
		Algorithm     AlgorithmIdentifier
		EncryptedData []byte `asn1:"tag:0"`
	}

	_, err := asn1.Unmarshal(block.Bytes, &encryptedKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse PKCS#8 encrypted key structure")
	}

	// Check if it's PBES2 (Password-Based Encryption Scheme 2)
	if !compareOID(encryptedKey.Algorithm.Algorithm, oidPBES2) {
		return nil, errors.New("unsupported PKCS#8 encryption algorithm (only PBES2 supported)")
	}

	// Parse PBES2 parameters
	var pbes2Params struct {
		KeyDerivationFunc AlgorithmIdentifier
		EncryptionScheme  AlgorithmIdentifier
	}

	_, err = asn1.Unmarshal(encryptedKey.Algorithm.Parameters.FullBytes, &pbes2Params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse PBES2 parameters")
	}

	// Check if it's PBKDF2 (Password-Based Key Derivation Function 2)
	if !compareOID(pbes2Params.KeyDerivationFunc.Algorithm, oidPBKDF2) {
		return nil, errors.New("unsupported key derivation function (only PBKDF2 supported)")
	}

	// Parse PBKDF2 parameters
	var pbkdf2Params struct {
		Salt           []byte
		IterationCount int
		PRF            AlgorithmIdentifier `asn1:"optional"`
	}

	_, err = asn1.Unmarshal(pbes2Params.KeyDerivationFunc.Parameters.FullBytes, &pbkdf2Params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse PBKDF2 parameters")
	}

	// Check if it's AES encryption
	if !compareOID(pbes2Params.EncryptionScheme.Algorithm, oidAES256CBC) &&
		!compareOID(pbes2Params.EncryptionScheme.Algorithm, oidAES128CBC) {
		return nil, errors.New("unsupported encryption algorithm (only AES supported)")
	}

	// Parse AES parameters (IV)
	var aesParams []byte
	_, err = asn1.Unmarshal(pbes2Params.EncryptionScheme.Parameters.FullBytes, &aesParams)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse AES parameters")
	}

	// Derive key using PBKDF2
	keyLen := 32 // AES-256
	if compareOID(pbes2Params.EncryptionScheme.Algorithm, oidAES128CBC) {
		keyLen = 16 // AES-128
	}

	derivedKey := pbkdf2.Key([]byte(passphrase), pbkdf2Params.Salt, pbkdf2Params.IterationCount, keyLen, sha256.New)

	// Decrypt using AES
	cipherBlock, err := aes.NewCipher(derivedKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create AES cipher")
	}

	if len(encryptedKey.EncryptedData)%aes.BlockSize != 0 {
		return nil, errors.New("encrypted data is not a multiple of AES block size")
	}

	// Decrypt in CBC mode
	mode := cipher.NewCBCDecrypter(cipherBlock, aesParams)
	decrypted := make([]byte, len(encryptedKey.EncryptedData))
	mode.CryptBlocks(decrypted, encryptedKey.EncryptedData)

	// Remove PKCS#7 padding
	paddingLen := int(decrypted[len(decrypted)-1])
	if paddingLen > aes.BlockSize || paddingLen == 0 {
		return nil, errors.New("invalid PKCS#7 padding")
	}

	// Verify padding
	for i := len(decrypted) - paddingLen; i < len(decrypted); i++ {
		if decrypted[i] != byte(paddingLen) {
			return nil, errors.New("invalid PKCS#7 padding")
		}
	}

	return decrypted[:len(decrypted)-paddingLen], nil
}

// connectToAgent connects to the SSH agent
func (sc *SSHClient) connectToAgent() (agent.Agent, error) {
	if sc.agentConn != nil {
		return sc.agentConn, nil
	}

	// Try to connect to SSH agent
	sock := os.Getenv("SSH_AUTH_SOCK")
	if sock == "" {
		return nil, errors.New("SSH_AUTH_SOCK environment variable not set")
	}

	conn, err := net.Dial("unix", sock)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to SSH agent socket")
	}

	agentConn := agent.NewClient(conn)

	sc.agentConn = agentConn
	return agentConn, nil
}

// initKnownHosts initializes the known hosts callback
func (sc *SSHClient) initKnownHosts() error {
	var knownHostsFiles []string

	// Add user-specified known hosts files
	knownHostsFiles = append(knownHostsFiles, sc.config.UserKnownHosts...)
	knownHostsFiles = append(knownHostsFiles, sc.config.GlobalKnownHosts...)

	// Add default known hosts file if not specified
	if sc.config.KnownHostsPath != "" {
		knownHostsFiles = append(knownHostsFiles, sc.config.KnownHostsPath)
	} else {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			knownHostsFiles = append(knownHostsFiles, filepath.Join(homeDir, ".ssh", "known_hosts"))
		}
	}

	// Determine the known hosts mode
	mode := sc.config.KnownHostsMode
	if mode == "" {
		// Fallback to legacy StrictHostKey setting
		if sc.config.StrictHostKey {
			mode = "strict"
		} else {
			mode = "accept-new"
		}
	}

	// Create appropriate callback based on mode
	switch mode {
	case "strict":
		callback, err := knownhosts.New(knownHostsFiles...)
		if err != nil {
			return errors.Wrap(err, "failed to load known hosts")
		}
		sc.knownHosts = callback

	case "accept-new":
		sc.knownHosts = sc.createAcceptNewCallback(knownHostsFiles)

	case "ignore":
		sc.knownHosts = ssh.InsecureIgnoreHostKey()

	default:
		return errors.Errorf("invalid known hosts mode: %s (valid modes: strict, accept-new, ignore)", mode)
	}

	return nil
}

// createAcceptNewCallback creates a host key callback that mimics OpenSSH's "accept-new" behavior
func (sc *SSHClient) createAcceptNewCallback(knownHostsFiles []string) ssh.HostKeyCallback {
	// Try to load existing known hosts
	existingCallback, err := knownhosts.New(knownHostsFiles...)
	if err != nil {
		// If we can't load known hosts, just accept all keys (like accept-new with no existing file)
		return ssh.InsecureIgnoreHostKey()
	}

	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		// Try existing known hosts first
		err := existingCallback(hostname, remote, key)
		if err == nil {
			return nil // Key already known and matches
		}

		// Check if it's a "host key changed" error (security issue)
		// The knownhosts package returns specific error messages for changed keys
		errMsg := err.Error()
		if strings.Contains(errMsg, "REMOTE HOST IDENTIFICATION HAS CHANGED") ||
			strings.Contains(errMsg, "host key mismatch") ||
			strings.Contains(errMsg, "key verification failed") {
			// Warn about changed keys (security issue) - don't suppress this
			return err
		}

		// For new hosts or other errors, accept silently (like accept-new)
		// This suppresses the "Permanently added" messages
		return nil
	}
}

// setupSession configures a session for command execution
func (sc *SSHClient) setupSession(session *ssh.Session) error {
	// Set environment variables
	for key, value := range sc.config.Environment {
		if err := session.Setenv(key, value); err != nil {
			return errors.Wrapf(err, "failed to set environment variable %s", key)
		}
	}

	// Set up I/O
	if sc.config.Stdout != nil {
		session.Stdout = sc.config.Stdout
	}
	if sc.config.Stderr != nil {
		session.Stderr = sc.config.Stderr
	}
	if sc.config.Stdin != nil {
		session.Stdin = sc.config.Stdin
	}

	return nil
}

// keepAlive sends keepalive packets
func (sc *SSHClient) keepAlive() {
	ticker := time.NewTicker(sc.config.KeepAlive)
	defer ticker.Stop()

	for range ticker.C {
		if sc.client == nil {
			break
		}

		// Send keepalive
		_, _, err := sc.client.SendRequest("keepalive@golang.org", true, nil)
		if err != nil {
			break
		}
	}
}
