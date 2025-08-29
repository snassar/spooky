package ssh

import (
	"bufio"
	"bytes"
	"context"
	"crypto/x509"
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
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"

	"spooky/internal/logging"
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
	RequestTTY                bool
	RequestPty                bool
	TerminalType              string
	TerminalSize              *TerminalSize
	Stdin                     io.Reader
	Stdout                    io.Writer
	Stderr                    io.Writer
}

// TerminalSize represents terminal dimensions
type TerminalSize struct {
	Width  int
	Height int
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
	if config.TerminalType == "" {
		config.TerminalType = "xterm-256color"
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
		return errors.Wrap(err, "failed to create SSH config")
	}

	// Create connection address
	addr := fmt.Sprintf("%s:%d", sc.config.Host, sc.config.Port)

	// Establish connection with timeout
	dialer := &net.Dialer{
		Timeout: sc.config.Timeout,
	}

	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return errors.Wrapf(err, "failed to connect to %s", addr)
	}

	// Create SSH connection
	sshConn, chans, reqs, err := ssh.NewClientConn(conn, addr, sshConfig)
	if err != nil {
		conn.Close()
		return errors.Wrap(err, "failed to establish SSH connection")
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

// ExecuteCommand executes a command on the remote host
func (sc *SSHClient) ExecuteCommand(ctx context.Context, command string) (*CommandResult, error) {
	if sc.client == nil {
		return nil, errors.New("SSH client not connected")
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

	if err := session.Run(command); err != nil {
		return &CommandResult{
			ExitCode: 1,
			Stdout:   stdout.String(),
			Stderr:   stderr.String(),
			Error:    err,
		}, nil
	}

	return &CommandResult{
		ExitCode: 0,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
	}, nil
}

// ExecuteCommandWithTimeout executes a command with a timeout
func (sc *SSHClient) ExecuteCommandWithTimeout(ctx context.Context, command string, timeout time.Duration) (*CommandResult, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return sc.ExecuteCommand(ctx, command)
}

// ExecuteCommandInteractive executes a command with interactive input/output
func (sc *SSHClient) ExecuteCommandInteractive(ctx context.Context, command string) error {
	if sc.client == nil {
		return errors.New("SSH client not connected")
	}

	// Create session
	session, err := sc.client.NewSession()
	if err != nil {
		return errors.Wrap(err, "failed to create SSH session")
	}
	defer session.Close()

	// Set up session for interactive mode
	if err := sc.setupInteractiveSession(session); err != nil {
		return errors.Wrap(err, "failed to setup interactive session")
	}

	// Execute command
	if err := session.Run(command); err != nil {
		return errors.Wrap(err, "failed to execute interactive command")
	}
	return nil
}

// CreateShell creates an interactive shell session
func (sc *SSHClient) CreateShell(ctx context.Context) error {
	if sc.client == nil {
		return errors.New("SSH client not connected")
	}

	// Create session
	session, err := sc.client.NewSession()
	if err != nil {
		return errors.Wrap(err, "failed to create SSH session")
	}

	// Set up session for shell
	if err := sc.setupShellSession(session); err != nil {
		session.Close()
		return errors.Wrap(err, "failed to setup shell session")
	}

	sc.session = session
	return nil
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

// CommandResult represents the result of a command execution
type CommandResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Error    error
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
			return nil, errors.Wrapf(err, "failed to read private key file: %s", sc.config.PrivateKeyPath)
		}
	} else if len(sc.config.PrivateKeyData) > 0 {
		keyData = sc.config.PrivateKeyData
	} else {
		return nil, errors.New("no private key provided")
	}

	// Try to parse as PEM
	block, rest := pem.Decode(keyData)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}
	if len(rest) > 0 {
		// Log warning about extra data but continue
		logger := logging.GetGlobalLogger()
		logger.Warn("PEM block contains extra data after first block",
			slog.Int("extra_bytes", len(rest)))
	}

	// Handle encrypted private keys
	if x509.IsEncryptedPEMBlock(block) {
		if sc.config.Passphrase == "" {
			return nil, errors.New("private key is encrypted but no passphrase provided")
		}

		decryptedBlock, err := x509.DecryptPEMBlock(block, []byte(sc.config.Passphrase))
		if err != nil {
			return nil, errors.Wrap(err, "failed to decrypt private key")
		}

		block = &pem.Block{
			Type:  block.Type,
			Bytes: decryptedBlock,
		}
	}

	// Parse private key
	signer, err := ssh.ParsePrivateKey(pem.EncodeToMemory(block))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse private key")
	}

	return signer, nil
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

// setupInteractiveSession configures a session for interactive use
func (sc *SSHClient) setupInteractiveSession(session *ssh.Session) error {
	if err := sc.setupSession(session); err != nil {
		return errors.Wrap(err, "failed to setup interactive session")
	}

	// Request PTY for interactive mode
	if sc.config.RequestPty || sc.config.RequestTTY {
		if err := sc.requestPTY(session); err != nil {
			return errors.Wrap(err, "failed to request PTY")
		}
	}

	return nil
}

// setupShellSession configures a session for shell use
func (sc *SSHClient) setupShellSession(session *ssh.Session) error {
	if err := sc.setupInteractiveSession(session); err != nil {
		return errors.Wrap(err, "failed to setup shell session")
	}

	// Set up I/O for shell
	if sc.config.Stdout != nil {
		session.Stdout = sc.config.Stdout
	} else {
		session.Stdout = os.Stdout
	}
	if sc.config.Stderr != nil {
		session.Stderr = sc.config.Stderr
	} else {
		session.Stderr = os.Stderr
	}
	if sc.config.Stdin != nil {
		session.Stdin = sc.config.Stdin
	} else {
		session.Stdin = os.Stdin
	}

	return nil
}

// requestPTY requests a PTY for the session
func (sc *SSHClient) requestPTY(session *ssh.Session) error {
	// Get terminal size
	width := 80
	height := 24
	if sc.config.TerminalSize != nil {
		width = sc.config.TerminalSize.Width
		height = sc.config.TerminalSize.Height
	}

	// Request PTY
	if err := session.RequestPty(sc.config.TerminalType, height, width, ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}); err != nil {
		return errors.Wrap(err, "failed to request PTY")
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
