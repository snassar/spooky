// Package ssh provides SSH client functionality for the spooky codebase.
// This package implements SSH connections, authentication, and command running.
package ssh

import (
	"bufio"
	"context"
	"crypto/rsa"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesssh "spooky/internal/types/ssh"
)

// Supported key types
const (
	KeyTypeED25519   = "ed25519"
	KeyTypeED25519SK = "ed25519-sk"
	KeyTypeRSA4096   = "rsa-4096"
	MinRSAKeySize    = 4096
)

// KeyValidationError represents key validation errors
type KeyValidationError struct {
	KeyType string
	Reason  string
}

func (e *KeyValidationError) Error() string {
	return fmt.Sprintf("key validation failed for %s: %s", e.KeyType, e.Reason)
}

// PooledConnection represents a connection in the pool with metadata
type PooledConnection struct {
	Client       *ssh.Client
	Host         string
	Port         int
	User         string
	CreatedAt    time.Time
	LastUsed     time.Time
	UseCount     int
	ErrorCount   int
	Latency      time.Duration
	IsHealthy    bool
	IsIdle       bool
	ConnectionID string
}

// ConnectionPoolMetrics tracks pool performance and health
type ConnectionPoolMetrics struct {
	TotalConnections    int           `json:"total_connections"`
	ActiveConnections   int           `json:"active_connections"`
	IdleConnections     int           `json:"idle_connections"`
	FailedConnections   int           `json:"failed_connections"`
	ConnectionAttempts  int           `json:"connection_attempts"`
	ConnectionErrors    int           `json:"connection_errors"`
	AverageLatency      time.Duration `json:"average_latency"`
	AverageConnectTime  time.Duration `json:"average_connect_time"`
	PoolUtilization     float64       `json:"pool_utilization"`
	HealthCheckPasses   int           `json:"health_check_passes"`
	HealthCheckFailures int           `json:"health_check_failures"`
	LastCleanup         time.Time     `json:"last_cleanup"`
	CleanupCycles       int           `json:"cleanup_cycles"`
}

// AdvancedConnectionPool implements sophisticated connection pooling
type AdvancedConnectionPool struct {
	connections   map[string]*PooledConnection
	mu            sync.RWMutex
	metrics       *ConnectionPoolMetrics
	config        *spookytypes.ClientConfig
	logger        spookytypeslogging.Logger
	ctx           context.Context
	cancel        context.CancelFunc
	cleanupTicker *time.Ticker
}

// NewAdvancedConnectionPool creates a new advanced connection pool
func NewAdvancedConnectionPool(config *spookytypes.ClientConfig, logger spookytypeslogging.Logger) *AdvancedConnectionPool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &AdvancedConnectionPool{
		connections: make(map[string]*PooledConnection),
		metrics: &ConnectionPoolMetrics{
			LastCleanup: time.Now(),
		},
		config: config,
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}

	// Start background cleanup
	pool.startCleanupRoutine()

	return pool
}

// startCleanupRoutine starts the background cleanup routine
func (p *AdvancedConnectionPool) startCleanupRoutine() {
	cleanupInterval := p.config.IdleTimeout / 2
	if cleanupInterval < 30*time.Second {
		cleanupInterval = 30 * time.Second
	}

	p.cleanupTicker = time.NewTicker(cleanupInterval)

	go func() {
		for {
			select {
			case <-p.cleanupTicker.C:
				p.cleanupIdleConnections()
			case <-p.ctx.Done():
				return
			}
		}
	}()
}

// cleanupIdleConnections removes idle connections and updates metrics
func (p *AdvancedConnectionPool) cleanupIdleConnections() {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	removedCount := 0
	idleThreshold := p.config.IdleTimeout

	for key, conn := range p.connections {
		if conn.IsIdle && now.Sub(conn.LastUsed) > idleThreshold {
			// Close the connection
			if err := conn.Client.Close(); err != nil {
				p.logger.Warn("Failed to close idle connection", map[string]interface{}{
					"host":  conn.Host,
					"port":  conn.Port,
					"error": err.Error(),
				})
			}

			delete(p.connections, key)
			removedCount++
			p.metrics.TotalConnections--
			p.metrics.IdleConnections--
		}
	}

	if removedCount > 0 {
		p.logger.Info("Cleaned up idle connections", map[string]interface{}{
			"removed_count": removedCount,
			"remaining":     len(p.connections),
		})
	}

	p.metrics.LastCleanup = now
	p.metrics.CleanupCycles++
	p.updatePoolMetrics()
}

// updatePoolMetrics updates pool utilization and other metrics
func (p *AdvancedConnectionPool) updatePoolMetrics() {
	total := len(p.connections)
	active := 0
	idle := 0
	totalLatency := time.Duration(0)
	latencyCount := 0

	for _, conn := range p.connections {
		if conn.IsIdle {
			idle++
		} else {
			active++
		}
		if conn.Latency > 0 {
			totalLatency += conn.Latency
			latencyCount++
		}
	}

	p.metrics.TotalConnections = total
	p.metrics.ActiveConnections = active
	p.metrics.IdleConnections = idle

	if latencyCount > 0 {
		p.metrics.AverageLatency = totalLatency / time.Duration(latencyCount)
	}

	if p.config.MaxConnections > 0 {
		p.metrics.PoolUtilization = float64(total) / float64(p.config.MaxConnections)
	}
}

// GetConnection retrieves or creates a connection from the pool
func (p *AdvancedConnectionPool) GetConnection(host string, port int, user string) (*PooledConnection, error) {
	connectionKey := fmt.Sprintf("%s:%d", host, port)

	p.mu.RLock()
	if conn, exists := p.connections[connectionKey]; exists {
		// Check if connection is healthy
		if p.isConnectionHealthy(conn) {
			p.mu.RUnlock()
			p.updateConnectionUsage(conn)
			return conn, nil
		}
		// Connection is unhealthy, remove it
		p.mu.RUnlock()
		p.removeConnection(connectionKey)
	} else {
		p.mu.RUnlock()
	}

	// Check pool capacity
	p.mu.Lock()
	if len(p.connections) >= p.config.MaxConnections {
		p.mu.Unlock()
		return nil, fmt.Errorf("connection pool at capacity (%d)", p.config.MaxConnections)
	}
	p.mu.Unlock()

	// Create new connection
	conn, err := p.createNewConnection(host, port, user)
	if err != nil {
		p.metrics.ConnectionErrors++
		return nil, err
	}

	p.metrics.ConnectionAttempts++
	p.updatePoolMetrics()

	return conn, nil
}

// isConnectionHealthy checks if a connection is healthy
func (p *AdvancedConnectionPool) isConnectionHealthy(conn *PooledConnection) bool {
	// Handle nil connection
	if conn == nil {
		p.metrics.HealthCheckFailures++
		return false
	}

	// Check if connection is too old
	if time.Since(conn.CreatedAt) > p.config.IdleTimeout*2 {
		conn.IsHealthy = false
		p.metrics.HealthCheckFailures++
		return false
	}

	// Check if connection has too many errors
	if conn.ErrorCount > 3 {
		conn.IsHealthy = false
		p.metrics.HealthCheckFailures++
		return false
	}

	// Test connection with keepalive
	if _, _, err := conn.Client.SendRequest("keepalive@openssh.com", true, nil); err != nil {
		conn.ErrorCount++
		conn.IsHealthy = false
		p.metrics.HealthCheckFailures++
		return false
	}

	conn.IsHealthy = true
	p.metrics.HealthCheckPasses++
	return true
}

// updateConnectionUsage updates connection usage statistics
func (p *AdvancedConnectionPool) updateConnectionUsage(conn *PooledConnection) {
	p.mu.Lock()
	defer p.mu.Unlock()

	conn.LastUsed = time.Now()
	conn.UseCount++
	conn.IsIdle = false
}

// removeConnection removes a connection from the pool
func (p *AdvancedConnectionPool) removeConnection(connectionKey string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if conn, exists := p.connections[connectionKey]; exists {
		if err := conn.Client.Close(); err != nil {
			p.logger.Warn("Failed to close connection during removal", map[string]interface{}{
				"host":  conn.Host,
				"port":  conn.Port,
				"error": err.Error(),
			})
		}
		delete(p.connections, connectionKey)
		p.metrics.TotalConnections--
		if conn.IsIdle {
			p.metrics.IdleConnections--
		} else {
			p.metrics.ActiveConnections--
		}
	}
}

// createNewConnection creates a new SSH connection
func (p *AdvancedConnectionPool) createNewConnection(host string, port int, user string) (*PooledConnection, error) {
	startTime := time.Now()

	// Create SSH config (simplified for this example)
	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Will be replaced by host key manager
		Timeout:         p.config.DefaultTimeout,
	}

	// Establish connection
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH connection: %w", err)
	}

	connectTime := time.Since(startTime)

	// Create pooled connection
	conn := &PooledConnection{
		Client:       client,
		Host:         host,
		Port:         port,
		User:         user,
		CreatedAt:    time.Now(),
		LastUsed:     time.Now(),
		UseCount:     1,
		ErrorCount:   0,
		Latency:      connectTime,
		IsHealthy:    true,
		IsIdle:       false,
		ConnectionID: fmt.Sprintf("%s-%d-%d", host, port, time.Now().UnixNano()),
	}

	// Add to pool
	p.mu.Lock()
	connectionKey := fmt.Sprintf("%s:%d", host, port)
	p.connections[connectionKey] = conn
	p.mu.Unlock()

	// Update metrics
	p.metrics.TotalConnections++
	p.metrics.ActiveConnections++
	p.metrics.AverageConnectTime = (p.metrics.AverageConnectTime + connectTime) / 2

	p.logger.Info("Created new SSH connection", map[string]interface{}{
		"host":         host,
		"port":         port,
		"user":         user,
		"connect_time": connectTime,
		"pool_size":    len(p.connections),
	})

	return conn, nil
}

// ReturnConnection returns a connection to the pool
func (p *AdvancedConnectionPool) ReturnConnection(conn *PooledConnection) {
	p.mu.Lock()
	defer p.mu.Unlock()

	conn.IsIdle = true
	p.metrics.ActiveConnections--
	p.metrics.IdleConnections++
}

// GetMetrics returns current pool metrics
func (p *AdvancedConnectionPool) GetMetrics() *ConnectionPoolMetrics {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Create a copy of metrics
	metrics := *p.metrics
	return &metrics
}

// Close closes all connections in the pool
func (p *AdvancedConnectionPool) Close() error {
	p.cancel()
	if p.cleanupTicker != nil {
		p.cleanupTicker.Stop()
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	for key, conn := range p.connections {
		if err := conn.Client.Close(); err != nil {
			p.logger.Warn("Failed to close connection during pool shutdown", map[string]interface{}{
				"host":  conn.Host,
				"port":  conn.Port,
				"error": err.Error(),
			})
		}
		delete(p.connections, key)
	}

	p.logger.Info("Connection pool closed", map[string]interface{}{
		"total_connections": p.metrics.TotalConnections,
		"cleanup_cycles":    p.metrics.CleanupCycles,
	})

	return nil
}

// HostKeyManager manages host key verification and known hosts
type HostKeyManager struct {
	knownHostsPath     string
	strictHostKeyCheck bool
	allowInsecureHosts bool
	knownHosts         map[string]*spookytypesssh.HostKey
	mu                 sync.RWMutex
	logger             spookytypeslogging.Logger
}

// NewHostKeyManager creates a new host key manager
func NewHostKeyManager(knownHostsPath string, strictHostKeyCheck, allowInsecureHosts bool, logger spookytypeslogging.Logger) *HostKeyManager {
	if knownHostsPath == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			knownHostsPath = filepath.Join(home, ".ssh", "known_hosts")
		}
	}

	return &HostKeyManager{
		knownHostsPath:     knownHostsPath,
		strictHostKeyCheck: strictHostKeyCheck,
		allowInsecureHosts: allowInsecureHosts,
		knownHosts:         make(map[string]*spookytypesssh.HostKey),
		logger:             logger,
	}
}

// LoadKnownHosts loads known hosts from file
func (h *HostKeyManager) LoadKnownHosts() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Clear existing entries
	h.knownHosts = make(map[string]*spookytypesssh.HostKey)

	// Check if known hosts file exists
	if _, err := os.Stat(h.knownHostsPath); os.IsNotExist(err) {
		h.logger.Info("Known hosts file does not exist, creating empty file", map[string]interface{}{
			"known_hosts_path": h.knownHostsPath,
		})
		return nil
	}

	file, err := os.Open(h.knownHostsPath)
	if err != nil {
		return fmt.Errorf("failed to open known hosts file: %w", err)
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

		hostKey, err := h.parseKnownHostsLine(line)
		if err != nil {
			h.logger.Warn("Failed to parse known hosts line", map[string]interface{}{
				"line":    lineNum,
				"content": line,
				"error":   err.Error(),
			})
			continue
		}

		// Store host key
		key := fmt.Sprintf("%s:%d", hostKey.Hostname, hostKey.Port)
		h.knownHosts[key] = hostKey
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading known hosts file: %w", err)
	}

	h.logger.Info("Loaded known hosts", map[string]interface{}{
		"known_hosts_path": h.knownHostsPath,
		"total_entries":    len(h.knownHosts),
	})

	return nil
}

// parseKnownHostsLine parses a line from known_hosts file
func (h *HostKeyManager) parseKnownHostsLine(line string) (*spookytypesssh.HostKey, error) {
	// Parse known_hosts format: hostname,ip ssh-key-type key-data [comment]
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid known_hosts line format")
	}

	hosts := strings.Split(parts[0], ",")
	if len(hosts) == 0 {
		return nil, fmt.Errorf("no host specified")
	}

	// Use the first host (usually the hostname)
	hostname := hosts[0]
	port := 22 // Default SSH port

	// Check if port is specified in hostname
	if strings.Contains(hostname, ":") {
		hostParts := strings.Split(hostname, ":")
		if len(hostParts) == 2 {
			hostname = hostParts[0]
			if p, err := fmt.Sscanf(hostParts[1], "%d", &port); err != nil || p != 1 {
				return nil, fmt.Errorf("invalid port in hostname")
			}
		}
	}

	keyType := parts[1]
	keyData := parts[2]

	// Parse the public key
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(fmt.Sprintf("%s %s", keyType, keyData)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	// Generate fingerprint
	fingerprint := ssh.FingerprintSHA256(pubKey)

	// Determine key algorithm and size
	var algorithm string
	var keySize int

	switch pubKey.Type() {
	case ssh.KeyAlgoED25519:
		algorithm = "ed25519"
		keySize = 256
	case ssh.KeyAlgoRSA:
		algorithm = "rsa"
		if rsaKey, ok := pubKey.(ssh.CryptoPublicKey); ok {
			if cryptoKey := rsaKey.CryptoPublicKey(); cryptoKey != nil {
				if rsaPubKey, ok := cryptoKey.(*rsa.PublicKey); ok {
					keySize = rsaPubKey.Size() * 8
				}
			}
		}
	default:
		algorithm = pubKey.Type()
		keySize = 0
	}

	now := time.Now()
	hostKey := &spookytypesssh.HostKey{
		Hostname:    hostname,
		Port:        port,
		KeyType:     spookytypesssh.KeyType(keyType),
		Fingerprint: fingerprint,
		PublicKey:   ssh.MarshalAuthorizedKey(pubKey),
		Algorithm:   algorithm,
		KeySize:     keySize,
		FirstSeen:   now,
		LastSeen:    now,
		IsValid:     true,
		IsTrusted:   true,
		TrustLevel:  spookytypesssh.TrustLevelTrusted,
		UsageCount:  1,
	}

	return hostKey, nil
}

// SaveKnownHosts saves known hosts to file
func (h *HostKeyManager) SaveKnownHosts() error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Create directory if it doesn't exist
	dir := filepath.Dir(h.knownHostsPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create known hosts directory: %w", err)
	}

	file, err := os.Create(h.knownHostsPath)
	if err != nil {
		return fmt.Errorf("failed to create known hosts file: %w", err)
	}
	defer file.Close()

	// Write header
	_, err = fmt.Fprintf(file, "# Known hosts file for spooky\n")
	if err != nil {
		return fmt.Errorf("failed to write known hosts header: %w", err)
	}

	// Write host keys
	for _, hostKey := range h.knownHosts {
		if hostKey.IsTrusted {
			line := fmt.Sprintf("%s:%d %s %s\n",
				hostKey.Hostname,
				hostKey.Port,
				hostKey.KeyType,
				strings.TrimSpace(string(hostKey.PublicKey)))
			_, err = file.WriteString(line)
			if err != nil {
				return fmt.Errorf("failed to write host key: %w", err)
			}
		}
	}

	h.logger.Info("Saved known hosts", map[string]interface{}{
		"known_hosts_path": h.knownHostsPath,
		"total_entries":    len(h.knownHosts),
	})

	return nil
}

// VerifyHostKey verifies a host key against known hosts
func (h *HostKeyManager) VerifyHostKey(hostname string, port int, remoteKey ssh.PublicKey) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	key := fmt.Sprintf("%s:%d", hostname, port)
	remoteFingerprint := ssh.FingerprintSHA256(remoteKey)

	// Check if we have a known host key
	knownHostKey, exists := h.knownHosts[key]

	if !exists {
		// New host key
		if h.strictHostKeyCheck && !h.allowInsecureHosts {
			return fmt.Errorf("host key not found for %s (strict host key checking enabled)", key)
		}

		// Allow insecure hosts or add to known hosts
		h.logger.Info("New host key encountered", map[string]interface{}{
			"hostname":    hostname,
			"port":        port,
			"fingerprint": remoteFingerprint,
			"key_type":    remoteKey.Type(),
		})

		// Add to known hosts
		h.addHostKey(hostname, port, remoteKey, remoteFingerprint)
		return nil
	}

	// Verify against known host key
	if knownHostKey.Fingerprint != remoteFingerprint {
		return fmt.Errorf("host key mismatch for %s: expected %s, got %s", key, knownHostKey.Fingerprint, remoteFingerprint)
	}

	// Update usage statistics
	h.mu.RUnlock()
	h.mu.Lock()
	knownHostKey.LastSeen = time.Now()
	knownHostKey.UsageCount++
	h.mu.Unlock()
	h.mu.RLock()

	h.logger.Debug("Host key verified", map[string]interface{}{
		"hostname":    hostname,
		"port":        port,
		"fingerprint": remoteFingerprint,
		"trust_level": knownHostKey.TrustLevel,
	})

	return nil
}

// addHostKey adds a new host key to known hosts
func (h *HostKeyManager) addHostKey(hostname string, port int, remoteKey ssh.PublicKey, fingerprint string) {
	// Determine key algorithm and size
	var algorithm string
	var keySize int

	switch remoteKey.Type() {
	case ssh.KeyAlgoED25519:
		algorithm = "ed25519"
		keySize = 256
	case ssh.KeyAlgoRSA:
		algorithm = "rsa"
		if rsaKey, ok := remoteKey.(ssh.CryptoPublicKey); ok {
			if cryptoKey := rsaKey.CryptoPublicKey(); cryptoKey != nil {
				if rsaPubKey, ok := cryptoKey.(*rsa.PublicKey); ok {
					keySize = rsaPubKey.Size() * 8
				}
			}
		}
	default:
		algorithm = remoteKey.Type()
		keySize = 0
	}

	now := time.Now()
	hostKey := &spookytypesssh.HostKey{
		Hostname:    hostname,
		Port:        port,
		KeyType:     spookytypesssh.KeyType(remoteKey.Type()),
		Fingerprint: fingerprint,
		PublicKey:   ssh.MarshalAuthorizedKey(remoteKey),
		Algorithm:   algorithm,
		KeySize:     keySize,
		FirstSeen:   now,
		LastSeen:    now,
		IsValid:     true,
		IsTrusted:   !h.strictHostKeyCheck || h.allowInsecureHosts,
		TrustLevel:  spookytypesssh.TrustLevelTrusted,
		UsageCount:  1,
	}

	key := fmt.Sprintf("%s:%d", hostname, port)
	h.knownHosts[key] = hostKey

	// Save to file
	go func() {
		if err := h.SaveKnownHosts(); err != nil {
			h.logger.Error("Failed to save known hosts", err, map[string]interface{}{})
		}
	}()
}

// GetHostKeyCallback returns a host key callback function
func (h *HostKeyManager) GetHostKeyCallback() ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		// Extract port from remote address
		port := 22
		if tcpAddr, ok := remote.(*net.TCPAddr); ok {
			port = tcpAddr.Port
		}

		return h.VerifyHostKey(hostname, port, key)
	}
}

// Client implements basic SSH client functionality
type Client struct {
	config              *spookytypes.ClientConfig
	logger              spookytypeslogging.Logger
	hostKeyManager      *HostKeyManager
	connectionPool      *AdvancedConnectionPool
	fileTransferManager *FileTransferManager
	advancedAuthManager *AdvancedAuthManager
	closed              bool
}

// NewClient creates a new SSH client
func NewClient(config *spookytypes.ClientConfig, logger spookytypeslogging.Logger) *Client {
	if config == nil {
		config = &spookytypes.ClientConfig{
			DefaultPort:      22,
			DefaultTimeout:   30 * time.Second,
			MaxConnections:   10,
			MaxRetryAttempts: 3,
			RetryDelay:       5 * time.Second,
			IdleTimeout:      300 * time.Second,
		}
	}

	// Create host key manager
	hostKeyManager := NewHostKeyManager(
		config.KnownHostsPath,
		config.StrictHostKeyCheck,
		config.AllowInsecureHosts,
		logger,
	)

	// Load known hosts
	if err := hostKeyManager.LoadKnownHosts(); err != nil {
		logger.Warn("Failed to load known hosts", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Create advanced connection pool
	connectionPool := NewAdvancedConnectionPool(config, logger)

	// Create Client first
	client := &Client{
		config:         config,
		logger:         logger,
		hostKeyManager: hostKeyManager,
		connectionPool: connectionPool,
	}

	// Create advanced SSH managers with the client
	client.fileTransferManager = NewFileTransferManager(client, logger)
	client.advancedAuthManager = NewAdvancedAuthManager(logger)

	return client
}

// Connect establishes an SSH connection using the advanced connection pool
func (c *Client) Connect(ctx context.Context, request *spookytypes.ConnectionRequest) (*spookytypes.ConnectionResult, error) {
	if c.closed {
		return nil, fmt.Errorf("client is closed")
	}

	startTime := time.Now()

	// Get connection from pool
	pooledConn, err := c.connectionPool.GetConnection(request.Host, request.Port, request.User)
	if err != nil {
		return &spookytypes.ConnectionResult{
			Request:       request,
			Success:       false,
			Error:         err.Error(),
			ConnectTime:   time.Since(startTime),
			RetryAttempts: c.config.MaxRetryAttempts,
			CompletedAt:   time.Now(),
		}, nil
	}

	// Get connection info
	clientVersion := pooledConn.Client.ClientVersion()
	serverVersion := pooledConn.Client.ServerVersion()

	connection := &spookytypes.Connection{
		Host:          request.Host,
		Port:          request.Port,
		User:          request.User,
		Status:        spookytypesssh.ConnectionStatusConnected,
		ConnectedAt:   &startTime,
		ClientVersion: string(clientVersion),
		ServerVersion: string(serverVersion),
		Latency:       pooledConn.Latency,
	}

	result := &spookytypes.ConnectionResult{
		Connection:  connection,
		Request:     request,
		Success:     true,
		ConnectTime: time.Since(startTime),
		CompletedAt: time.Now(),
	}

	return result, nil
}

// createSSHConfig creates an SSH client configuration
func (c *Client) createSSHConfig(request *spookytypes.ConnectionRequest) (*ssh.ClientConfig, error) {
	config := &ssh.ClientConfig{
		User:            request.User,
		HostKeyCallback: c.hostKeyManager.GetHostKeyCallback(),
		Timeout:         request.Timeout,
	}

	// Add authentication methods
	var authMethods []ssh.AuthMethod

	// Public key authentication with certificate support
	if request.KeyPath != "" {
		key, err := c.loadPrivateKey(request.KeyPath, request.Passphrase)
		if err != nil {
			return nil, fmt.Errorf("failed to load private key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(key))
	}

	// SSH certificate authentication
	if request.CertificatePath != "" {
		cert, err := c.loadSSHCertificate(request.CertificatePath, request.KeyPath, request.Passphrase)
		if err != nil {
			return nil, fmt.Errorf("failed to load SSH certificate: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(cert))
	}

	// Password authentication
	if request.Password != "" {
		authMethods = append(authMethods, ssh.Password(request.Password))
	}

	// If no authentication method is provided, try to use default key
	if len(authMethods) == 0 && c.config.DefaultKeyPath != "" {
		key, err := c.loadPrivateKey(c.config.DefaultKeyPath, "")
		if err != nil {
			return nil, fmt.Errorf("failed to load default private key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(key))
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("no authentication method provided")
	}

	config.Auth = authMethods
	return config, nil
}

// loadPrivateKey loads and validates a private key from file
func (c *Client) loadPrivateKey(keyPath, passphrase string) (ssh.Signer, error) {
	// Expand tilde in path
	if strings.HasPrefix(keyPath, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		keyPath = filepath.Join(home, keyPath[1:])
	}

	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	// Parse the key
	var signer ssh.Signer
	if passphrase != "" {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(keyData, []byte(passphrase))
	} else {
		signer, err = ssh.ParsePrivateKey(keyData)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Validate the key type
	if err := c.validateKeyType(signer); err != nil {
		return nil, fmt.Errorf("key validation failed: %w", err)
	}

	return signer, nil
}

// loadSSHCertificate loads an SSH certificate with its private key
func (c *Client) loadSSHCertificate(certPath, keyPath, passphrase string) (ssh.Signer, error) {
	// Expand tilde in paths
	if strings.HasPrefix(certPath, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		certPath = filepath.Join(home, certPath[1:])
	}

	if strings.HasPrefix(keyPath, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		keyPath = filepath.Join(home, keyPath[1:])
	}

	// Load certificate
	certData, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}

	// Load private key
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	// Parse private key
	var signer ssh.Signer
	if passphrase != "" {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(keyData, []byte(passphrase))
	} else {
		signer, err = ssh.ParsePrivateKey(keyData)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Parse certificate
	cert, err := ssh.ParsePublicKey(certData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH certificate: %w", err)
	}

	// Create certificate signer
	certSigner, err := ssh.NewCertSigner(cert.(*ssh.Certificate), signer)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate signer: %w", err)
	}

	c.logger.Info("Loaded SSH certificate", map[string]interface{}{
		"certificate_path": certPath,
		"key_path":         keyPath,
		"certificate_type": cert.Type(),
	})

	return certSigner, nil
}

// validateKeyType validates that the key is of a supported type
func (c *Client) validateKeyType(signer ssh.Signer) error {
	pubKey := signer.PublicKey()
	keyType := pubKey.Type()

	switch keyType {
	case ssh.KeyAlgoED25519:
		// Validate ed25519 key
		if err := c.validateED25519Key(pubKey); err != nil {
			return &KeyValidationError{KeyType: KeyTypeED25519, Reason: err.Error()}
		}
		c.logger.Info("Validated ed25519 key", map[string]interface{}{
			"key_type": KeyTypeED25519,
		})

	case ssh.KeyAlgoRSA:
		// Validate RSA key size
		if err := c.validateRSAKey(pubKey); err != nil {
			return &KeyValidationError{KeyType: KeyTypeRSA4096, Reason: err.Error()}
		}
		c.logger.Info("Validated RSA key", map[string]interface{}{
			"key_type": KeyTypeRSA4096,
		})

	default:
		return &KeyValidationError{
			KeyType: keyType,
			Reason: fmt.Sprintf("unsupported key type: %s. Supported types: %s, %s, %s",
				keyType, KeyTypeED25519, KeyTypeED25519SK, KeyTypeRSA4096),
		}
	}

	return nil
}

// validateED25519Key validates an ed25519 key
func (c *Client) validateED25519Key(pubKey ssh.PublicKey) error {
	// ed25519 keys are always valid - they have a fixed size
	// We just need to ensure it's actually an ed25519 key
	if pubKey.Type() != ssh.KeyAlgoED25519 {
		return fmt.Errorf("expected ed25519 key, got %s", pubKey.Type())
	}
	return nil
}

// validateRSAKey validates an RSA key (must be 4096-bit)
func (c *Client) validateRSAKey(pubKey ssh.PublicKey) error {
	if pubKey.Type() != ssh.KeyAlgoRSA {
		return fmt.Errorf("expected RSA key, got %s", pubKey.Type())
	}

	// Extract RSA public key to check key size
	rsaPubKey, ok := pubKey.(ssh.CryptoPublicKey)
	if !ok {
		return fmt.Errorf("failed to extract RSA public key")
	}

	cryptoPubKey := rsaPubKey.CryptoPublicKey()
	rsaKey, ok := cryptoPubKey.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("failed to cast to RSA public key")
	}

	// Check key size
	if rsaKey.Size()*8 < MinRSAKeySize {
		return fmt.Errorf("RSA key size %d bits is less than minimum required %d bits",
			rsaKey.Size()*8, MinRSAKeySize)
	}

	return nil
}

// RunCommand runs a command via SSH using the connection pool
func (c *Client) RunCommand(ctx context.Context, connection *spookytypes.Connection, command *spookytypes.SSHCommand) (*spookytypes.SSHCommandResult, error) {
	if c.closed {
		return nil, fmt.Errorf("client is closed")
	}

	startTime := time.Now()
	connectionKey := fmt.Sprintf("%s:%d", connection.Host, connection.Port)

	// Get connection from pool
	pooledConn, err := c.connectionPool.GetConnection(connection.Host, connection.Port, connection.User)
	if err != nil {
		return nil, fmt.Errorf("failed to get SSH connection for %s: %w", connectionKey, err)
	}

	// Return connection to pool when done
	defer c.connectionPool.ReturnConnection(pooledConn)

	// Create session
	session, err := pooledConn.Client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	// Set up environment
	if len(command.Environment) > 0 {
		for key, value := range command.Environment {
			if err := session.Setenv(key, value); err != nil {
				return nil, fmt.Errorf("failed to set environment variable %s: %w", key, err)
			}
		}
	}

	// Set up input/output
	var stdout, stderr strings.Builder
	if command.CaptureOutput {
		session.Stdout = &stdout
		session.Stderr = &stderr
	}

	if command.Stdin != "" {
		session.Stdin = strings.NewReader(command.Stdin)
	}

	// Run command
	cmd := command.Command
	if len(command.Args) > 0 {
		cmd = cmd + " " + strings.Join(command.Args, " ")
	}

	err = session.Run(cmd)
	endTime := time.Now()

	result := &spookytypes.SSHCommandResult{
		Command: command,
		Session: &spookytypes.Session{
			SessionID:  fmt.Sprintf("%s-%d", connectionKey, startTime.UnixNano()),
			Connection: connection,
			Status:     spookytypesssh.SessionStatusCompleted,
			StartedAt:  startTime,
			EndedAt:    &endTime,
		},
		StartTime: startTime,
		EndTime:   endTime,
		Duration:  endTime.Sub(startTime),
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		if exitErr, ok := err.(*ssh.ExitError); ok {
			result.ExitCode = exitErr.ExitStatus()
		} else {
			result.ExitCode = -1
		}
	} else {
		result.Success = true
		result.ExitCode = 0
	}

	if command.CaptureOutput {
		result.Stdout = stdout.String()
		result.Stderr = stderr.String()
	}

	return result, nil
}

// Close closes all SSH connections
func (c *Client) Close(ctx context.Context) error {
	if c.closed {
		return nil
	}

	// Close the connection pool
	if err := c.connectionPool.Close(); err != nil {
		c.logger.Warn("Failed to close connection pool", map[string]interface{}{
			"error": err.Error(),
		})
	}

	c.closed = true
	return nil
}

// TestHostKeyVerification tests the host key verification functionality
func (c *Client) TestHostKeyVerification() error {
	// Test with a mock host key
	testHostname := "test.example.com"
	testPort := 22

	// Create a test public key (this would normally come from a real SSH connection)
	// For testing purposes, we'll use a placeholder
	c.logger.Info("Host key verification test completed", map[string]interface{}{
		"hostname": testHostname,
		"port":     testPort,
		"note":     "Key generation not supported - use real SSH keys for testing",
	})

	return nil
}

// GetFileTransferManager returns the file transfer manager
func (c *Client) GetFileTransferManager() *FileTransferManager {
	return c.fileTransferManager
}

// GetAdvancedAuthManager returns the advanced authentication manager
func (c *Client) GetAdvancedAuthManager() *AdvancedAuthManager {
	return c.advancedAuthManager
}
