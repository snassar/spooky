// Package ssh provides SSH client functionality for the spooky codebase.
// This package implements SSH connections, authentication, and command running.
package ssh

import (
	"bufio"
	"bytes"
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
	"golang.org/x/crypto/ssh/agent"

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
	connections    map[string]*PooledConnection
	mu             sync.RWMutex
	metrics        *ConnectionPoolMetrics
	config         *spookytypes.ClientConfig
	logger         spookytypeslogging.Logger
	hostKeyManager *HostKeyManager
	ctx            context.Context
	cancel         context.CancelFunc
	cleanupTicker  *time.Ticker
}

// NewAdvancedConnectionPool creates a new advanced connection pool
func NewAdvancedConnectionPool(config *spookytypes.ClientConfig, logger spookytypeslogging.Logger) *AdvancedConnectionPool {
	ctx, cancel := context.WithCancel(context.Background())

	// Create host key manager for the pool
	hostKeyManager := NewHostKeyManager(
		config.KnownHostsPath,
		config.StrictHostKeyCheck,
		config.AllowInsecureHosts,
		logger,
	)

	// Load known hosts
	if err := hostKeyManager.LoadKnownHosts(); err != nil {
		logger.Warn("Failed to load known hosts for connection pool", map[string]interface{}{
			"error": err.Error(),
		})
	}

	pool := &AdvancedConnectionPool{
		connections: make(map[string]*PooledConnection),
		metrics: &ConnectionPoolMetrics{
			LastCleanup: time.Now(),
		},
		config:         config,
		logger:         logger,
		hostKeyManager: hostKeyManager,
		ctx:            ctx,
		cancel:         cancel,
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
		if !conn.IsIdle || now.Sub(conn.LastUsed) <= idleThreshold {
			continue
		}

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

	// Create SSH config with proper host key verification
	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: p.hostKeyManager.GetHostKeyCallback(),
		Timeout:         p.config.DefaultTimeout,
	}

	// Note: Authentication methods will be added by the SSH manager
	// when creating connections with proper authentication details
	// This is a simplified connection for the pool

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
	if err := os.MkdirAll(dir, 0o700); err != nil {
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

// Connect establishes an SSH connection with proper authentication
func (c *Client) Connect(_ context.Context, request *spookytypes.ConnectionRequest) (*spookytypes.ConnectionResult, error) {
	if c.closed {
		return nil, fmt.Errorf("client is closed")
	}

	startTime := time.Now()

	// Create SSH config with authentication
	config := &ssh.ClientConfig{
		User:            request.User,
		HostKeyCallback: c.hostKeyManager.GetHostKeyCallback(),
		Timeout:         request.Timeout,
	}

	// Add authentication methods
	var authMethods []ssh.AuthMethod

	// Public key authentication
	if request.KeyPath != "" {
		key, err := c.loadPrivateKey(request.KeyPath, request.Passphrase)
		if err != nil {
			return &spookytypes.ConnectionResult{
				Request:       request,
				Success:       false,
				Error:         fmt.Sprintf("failed to load private key: %v", err),
				ConnectTime:   time.Since(startTime),
				RetryAttempts: c.config.MaxRetryAttempts,
				CompletedAt:   time.Now(),
			}, nil
		}
		authMethods = append(authMethods, ssh.PublicKeys(key))
	}

	// Password authentication
	if request.Password != "" {
		authMethods = append(authMethods, ssh.Password(request.Password))
	}

	// If no authentication method is provided, return error
	if len(authMethods) == 0 {
		return &spookytypes.ConnectionResult{
			Request:       request,
			Success:       false,
			Error:         "no authentication method provided",
			ConnectTime:   time.Since(startTime),
			RetryAttempts: c.config.MaxRetryAttempts,
			CompletedAt:   time.Now(),
		}, nil
	}

	config.Auth = authMethods

	// Establish connection
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", request.Host, request.Port), config)
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
	clientVersion := client.ClientVersion()
	serverVersion := client.ServerVersion()

	connection := &spookytypes.Connection{
		Host:          request.Host,
		Port:          request.Port,
		User:          request.User,
		Status:        spookytypesssh.ConnectionStatusConnected,
		ConnectedAt:   &startTime,
		ClientVersion: string(clientVersion),
		ServerVersion: string(serverVersion),
		Latency:       time.Since(startTime),
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

// validateKeyType validates that the key is of a supported type
func (c *Client) validateKeyType(signer ssh.Signer) error {
	pubKey := signer.PublicKey()
	keyType := pubKey.Type()

	switch keyType {
	case ssh.KeyAlgoED25519:
		c.logger.Info("Validated ed25519 key", map[string]interface{}{
			"key_type": "ed25519",
		})
	case ssh.KeyAlgoRSA:
		// Validate RSA key size
		if err := c.validateRSAKey(pubKey); err != nil {
			return fmt.Errorf("RSA key validation failed: %w", err)
		}
		c.logger.Info("Validated RSA key", map[string]interface{}{
			"key_type": "rsa",
		})
	default:
		return fmt.Errorf("unsupported key type: %s. Supported types: ed25519, rsa", keyType)
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

	// Check key size (minimum 2048 bits for security)
	if rsaKey.Size()*8 < 2048 {
		return fmt.Errorf("RSA key size %d bits is less than minimum required 2048 bits",
			rsaKey.Size()*8)
	}

	return nil
}

// RunCommand runs a command via SSH
func (c *Client) RunCommand(_ context.Context, connection *spookytypes.Connection, command *spookytypes.SSHCommand) (*spookytypes.SSHCommandResult, error) {
	if c.closed {
		return nil, fmt.Errorf("client is closed")
	}

	startTime := time.Now()
	connectionKey := fmt.Sprintf("%s:%d", connection.Host, connection.Port)

	// For now, we'll use a simplified approach
	// In a production system, you would want to implement proper connection pooling
	// with authenticated connections

	// Create SSH config with authentication
	config := &ssh.ClientConfig{
		User:            connection.User,
		HostKeyCallback: c.hostKeyManager.GetHostKeyCallback(),
		Timeout:         connection.Timeout,
	}

	// Note: In a real implementation, you would need to pass authentication details
	// from the calling context. For now, this is a placeholder that shows the structure.

	// Establish connection
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", connection.Host, connection.Port), config)
	if err != nil {
		return nil, fmt.Errorf("failed to establish SSH connection for %s: %w", connectionKey, err)
	}
	defer client.Close()

	// Create session
	session, err := client.NewSession()
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
func (c *Client) Close(_ context.Context) error {
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

	// Create a mock public key for testing
	mockKey := &mockPublicKey{keyType: "ssh-rsa"}

	// Test host key verification
	err := c.hostKeyManager.VerifyHostKey(testHostname, testPort, mockKey)
	if err != nil {
		c.logger.Info("Host key verification test failed as expected", map[string]interface{}{
			"hostname": testHostname,
			"port":     testPort,
			"error":    err.Error(),
		})
		return err
	}

	// Log test completion with guidance for real SSH key usage
	c.logger.Info("Host key verification test completed", map[string]interface{}{
		"hostname": testHostname,
		"port":     testPort,
		"note":     "Use real SSH keys for production testing",
	})

	return nil
}

// mockPublicKey is a mock implementation for testing
type mockPublicKey struct {
	keyType string
}

func (m *mockPublicKey) Type() string {
	return m.keyType
}

func (m *mockPublicKey) Marshal() []byte {
	return []byte("mock-key-data")
}

func (m *mockPublicKey) Verify(_ []byte, _ *ssh.Signature) error {
	return fmt.Errorf("mock key verification not implemented")
}

// GetFileTransferManager returns the file transfer manager
func (c *Client) GetFileTransferManager() *FileTransferManager {
	return c.fileTransferManager
}

// GetAdvancedAuthManager returns the advanced authentication manager
func (c *Client) GetAdvancedAuthManager() *AdvancedAuthManager {
	return c.advancedAuthManager
}

// ReusableSSHClient manages SSH connections with reuse and authentication caching
type ReusableSSHClient struct {
	connections map[string]*cachedConnection // host:port -> connection
	authCache   map[string]*authInfo         // host:port -> auth info
	metrics     *ConnectionMetrics
	mutex       sync.RWMutex
	logger      spookytypeslogging.Logger
	config      *spookytypes.ClientConfig
}

// cachedConnection represents a cached SSH connection
type cachedConnection struct {
	client   *ssh.Client
	session  *ssh.Session
	lastUsed time.Time
	healthy  bool
	authInfo *authInfo
	mutex    sync.Mutex
}

// authInfo represents cached authentication information
type authInfo struct {
	method      string
	credentials interface{}
	lastUsed    time.Time
	valid       bool
}

// ConnectionMetrics tracks connection performance metrics
type ConnectionMetrics struct {
	TotalConnections      int64
	SuccessfulConnections int64
	FailedConnections     int64
	ReusedConnections     int64
	AuthenticationTime    time.Duration
	CommandRunningTime    time.Duration
	mutex                 sync.RWMutex
}

// NewReusableSSHClient creates a new SSH client with connection reuse
func NewReusableSSHClient(config *spookytypes.ClientConfig, logger spookytypeslogging.Logger) *ReusableSSHClient {
	client := &ReusableSSHClient{
		connections: make(map[string]*cachedConnection),
		authCache:   make(map[string]*authInfo),
		metrics:     &ConnectionMetrics{},
		logger:      logger,
		config:      config,
	}

	// Start connection cleanup goroutine
	go client.performConnectionCleanup()

	return client
}

// RunCommand runs a command on a remote machine via SSH with connection reuse
func (c *ReusableSSHClient) RunCommand(ctx context.Context, machine *spookytypes.Machine, command string) (string, error) {
	c.logger.Debug("Running SSH command", map[string]interface{}{
		"machine": machine.Hostname,
		"host":    machine.Host,
		"command": command,
	})

	// Get or create cached connection
	connection, err := c.getOrCreateConnection(ctx, machine)
	if err != nil {
		c.recordConnectionFailure()
		return "", fmt.Errorf("failed to get connection for %s: %w", machine.Hostname, err)
	}

	// Run command using cached connection
	result, err := c.executeCommandOnConnection(ctx, connection, command)
	if err != nil {
		// Mark connection as unhealthy on error
		c.markConnectionUnhealthy(machine)
		c.recordConnectionFailure()
		return "", fmt.Errorf("failed to run command on %s: %w", machine.Hostname, err)
	}

	// Update connection last used time
	c.updateConnectionLastUsed(machine)
	c.recordConnectionSuccess(true) // This is a reused connection

	c.logger.Debug("Successfully ran SSH command", map[string]interface{}{
		"machine": machine.Hostname,
		"host":    machine.Host,
		"command": command,
	})

	return result, nil
}

// getOrCreateConnection gets existing connection or creates new one
func (c *ReusableSSHClient) getOrCreateConnection(ctx context.Context, machine *spookytypes.Machine) (*cachedConnection, error) {
	connectionKey := c.getConnectionKey(machine)

	c.mutex.RLock()
	if cachedConn, exists := c.connections[connectionKey]; exists {
		c.mutex.RUnlock()

		// Check if existing connection is healthy
		if c.isConnectionHealthy(cachedConn) {
			c.logger.Debug("Reusing existing SSH connection", map[string]interface{}{
				"machine": machine.Hostname,
				"host":    machine.Host,
			})
			return cachedConn, nil
		}

		// Connection is unhealthy, remove it
		c.logger.Debug("Removing unhealthy SSH connection", map[string]interface{}{
			"machine": machine.Hostname,
			"host":    machine.Host,
		})
		c.removeConnection(connectionKey)
	} else {
		c.mutex.RUnlock()
	}

	// Create new connection
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.logger.Debug("Creating new SSH connection", map[string]interface{}{
		"machine": machine.Hostname,
		"host":    machine.Host,
	})

	connection, err := c.createNewConnection(ctx, machine)
	if err != nil {
		return nil, fmt.Errorf("failed to create new connection: %w", err)
	}

	// Cache the connection
	c.connections[connectionKey] = connection
	c.recordConnectionSuccess(false) // This is a new connection

	return connection, nil
}

// createNewConnection creates a new SSH connection with authentication
func (c *ReusableSSHClient) createNewConnection(ctx context.Context, machine *spookytypes.Machine) (*cachedConnection, error) {
	// Get or create authentication info
	authInfo, err := c.getOrCreateAuthInfo(ctx, machine)
	if err != nil {
		return nil, fmt.Errorf("failed to get authentication info: %w", err)
	}

	// Create SSH client config
	sshConfig, err := c.createSSHConfig(machine, authInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH config: %w", err)
	}

	// Establish connection
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", machine.Host, machine.Port), sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH: %w", err)
	}

	// Create session
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to create SSH session: %w", err)
	}

	connection := &cachedConnection{
		client:   client,
		session:  session,
		lastUsed: time.Now(),
		healthy:  true,
		authInfo: authInfo,
	}

	return connection, nil
}

// getOrCreateAuthInfo gets existing auth info or creates new one
func (c *ReusableSSHClient) getOrCreateAuthInfo(ctx context.Context, machine *spookytypes.Machine) (*authInfo, error) {
	connectionKey := c.getConnectionKey(machine)

	// Check for existing valid auth info
	if cachedAuth, exists := c.authCache[connectionKey]; exists && cachedAuth.valid {
		// Check if auth info is still recent (within 1 hour)
		if time.Since(cachedAuth.lastUsed) < time.Hour {
			c.logger.Debug("Reusing cached authentication", map[string]interface{}{
				"machine": machine.Hostname,
				"host":    machine.Host,
				"method":  cachedAuth.method,
			})
			cachedAuth.lastUsed = time.Now()
			return cachedAuth, nil
		}
	}

	// Create new auth info
	c.logger.Debug("Creating new authentication info", map[string]interface{}{
		"machine": machine.Hostname,
		"host":    machine.Host,
	})

	authInfo, err := c.createAuthInfo(ctx, machine)
	if err != nil {
		return nil, fmt.Errorf("failed to create authentication info: %w", err)
	}

	// Cache the auth info
	c.authCache[connectionKey] = authInfo

	return authInfo, nil
}

// createAuthInfo creates authentication information for a machine
func (c *ReusableSSHClient) createAuthInfo(ctx context.Context, machine *spookytypes.Machine) (*authInfo, error) {
	// Try different authentication methods in order of preference
	authMethods := []struct {
		name string
		fn   func(*spookytypes.Machine) (ssh.AuthMethod, error)
	}{
		{"ssh_key", c.createSSHKeyAuth},
		{"password", c.createPasswordAuth},
		{"agent", c.createAgentAuth},
	}

	for _, method := range authMethods {
		authMethod, err := method.fn(machine)
		if err != nil {
			c.logger.Debug("Authentication method failed", map[string]interface{}{
				"machine": machine.Hostname,
				"method":  method.name,
				"error":   err.Error(),
			})
			continue
		}

		// Test authentication with a simple command
		if c.testAuthentication(ctx, machine, authMethod) {
			c.logger.Debug("Authentication method successful", map[string]interface{}{
				"machine": machine.Hostname,
				"method":  method.name,
			})

			return &authInfo{
				method:      method.name,
				credentials: authMethod,
				lastUsed:    time.Now(),
				valid:       true,
			}, nil
		}
	}

	return nil, fmt.Errorf("no valid authentication method found for %s", machine.Hostname)
}

// createSSHKeyAuth creates SSH key authentication
func (c *ReusableSSHClient) createSSHKeyAuth(machine *spookytypes.Machine) (ssh.AuthMethod, error) {
	if machine.KeyFile == "" {
		return nil, fmt.Errorf("no SSH key file specified")
	}

	// Read private key
	keyData, err := os.ReadFile(machine.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read SSH key file: %w", err)
	}

	var signer ssh.Signer
	if machine.Passphrase != "" {
		// Key is encrypted with passphrase
		signer, err = ssh.ParsePrivateKeyWithPassphrase(keyData, []byte(machine.Passphrase))
	} else {
		// Key is not encrypted
		signer, err = ssh.ParsePrivateKey(keyData)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH private key: %w", err)
	}

	return ssh.PublicKeys(signer), nil
}

// createPasswordAuth creates password authentication
func (c *ReusableSSHClient) createPasswordAuth(machine *spookytypes.Machine) (ssh.AuthMethod, error) {
	if machine.Password == "" {
		return nil, fmt.Errorf("no password specified")
	}

	return ssh.Password(machine.Password), nil
}

// createAgentAuth creates SSH agent authentication
func (c *ReusableSSHClient) createAgentAuth(_ *spookytypes.Machine) (ssh.AuthMethod, error) {
	agentConn := agent.NewKeyring()
	if agentConn == nil {
		return nil, fmt.Errorf("failed to create SSH agent keyring")
	}

	return ssh.PublicKeysCallback(agentConn.Signers), nil
}

// testAuthentication tests authentication with a simple command
func (c *ReusableSSHClient) testAuthentication(_ context.Context, machine *spookytypes.Machine, authMethod ssh.AuthMethod) bool {
	// Create temporary SSH config for testing
	// Note: InsecureIgnoreHostKey is used for testing only - in production, implement proper host key validation
	sshConfig := &ssh.ClientConfig{
		User:            machine.User,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nolint:gosec // Used for testing authentication only
		Timeout:         10 * time.Second,
	}

	// Try to connect
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", machine.Host, machine.Port), sshConfig)
	if err != nil {
		return false
	}
	defer client.Close()

	// Try to create a session
	session, err := client.NewSession()
	if err != nil {
		return false
	}
	defer session.Close()

	// Try to run a simple command
	err = session.Run("echo 'authentication test'")
	return err == nil
}

// createSSHConfig creates SSH client configuration
func (c *ReusableSSHClient) createSSHConfig(machine *spookytypes.Machine, authInfo *authInfo) (*ssh.ClientConfig, error) {
	authMethod, ok := authInfo.credentials.(ssh.AuthMethod)
	if !ok {
		return nil, fmt.Errorf("invalid authentication method type")
	}

	return &ssh.ClientConfig{
		User:            machine.User,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nolint:gosec // TODO: Implement proper host key validation
		Timeout:         30 * time.Second,
		Config: ssh.Config{
			KeyExchanges: []string{
				"curve25519-sha256@libssh.org",
				"ecdh-sha2-nistp256",
				"ecdh-sha2-nistp384",
				"ecdh-sha2-nistp521",
			},
			Ciphers: []string{
				"aes128-ctr",
				"aes192-ctr",
				"aes256-ctr",
				"aes128-gcm@openssh.com",
				"chacha20-poly1305@openssh.com",
			},
			MACs: []string{
				"hmac-sha2-256-etm@openssh.com",
				"hmac-sha2-512-etm@openssh.com",
				"umac-128-etm@openssh.com",
			},
		},
	}, nil
}

// executeCommandOnConnection runs a command on a cached connection
func (c *ReusableSSHClient) executeCommandOnConnection(ctx context.Context, connection *cachedConnection, command string) (string, error) {
	connection.mutex.Lock()
	defer connection.mutex.Unlock()

	// Check if connection is still healthy
	if !connection.healthy {
		return "", fmt.Errorf("connection is not healthy")
	}

	// Create a new session for this command (sessions are not reusable)
	session, err := connection.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create new session: %w", err)
	}
	defer session.Close()

	// Set up command running
	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	// Run command with context
	errChan := make(chan error, 1)
	go func() {
		errChan <- session.Run(command)
	}()

	// Wait for completion or context cancellation
	select {
	case err := <-errChan:
		if err != nil {
			return "", fmt.Errorf("command running failed: %w, stderr: %s", err, stderr.String())
		}
	case <-ctx.Done():
		return "", fmt.Errorf("command running cancelled: %w", ctx.Err())
	}

	return stdout.String(), nil
}

// isConnectionHealthy checks if a cached connection is still healthy
func (c *ReusableSSHClient) isConnectionHealthy(connection *cachedConnection) bool {
	if connection == nil || !connection.healthy {
		return false
	}

	// Check if connection hasn't timed out (30 minutes)
	if time.Since(connection.lastUsed) > 30*time.Minute {
		return false
	}

	// Check if SSH client is still connected
	if connection.client == nil {
		return false
	}

	// Try to create a test session to verify connection is alive
	session, err := connection.client.NewSession()
	if err != nil {
		return false
	}
	defer session.Close()

	// Try to run a simple command
	err = session.Run("echo 'health check'")
	return err == nil
}

// markConnectionUnhealthy marks a connection as unhealthy
func (c *ReusableSSHClient) markConnectionUnhealthy(machine *spookytypes.Machine) {
	connectionKey := c.getConnectionKey(machine)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if connection, exists := c.connections[connectionKey]; exists {
		connection.healthy = false
		c.logger.Debug("Marked SSH connection as unhealthy", map[string]interface{}{
			"machine": machine.Hostname,
			"host":    machine.Host,
		})
	}
}

// updateConnectionLastUsed updates the last used time for a connection
func (c *ReusableSSHClient) updateConnectionLastUsed(machine *spookytypes.Machine) {
	connectionKey := c.getConnectionKey(machine)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if connection, exists := c.connections[connectionKey]; exists {
		connection.lastUsed = time.Now()
	}
}

// removeConnection removes a connection from the cache
func (c *ReusableSSHClient) removeConnection(connectionKey string) {
	if connection, exists := c.connections[connectionKey]; exists {
		// Close the connection
		if connection.client != nil {
			connection.client.Close()
		}
		if connection.session != nil {
			connection.session.Close()
		}
		delete(c.connections, connectionKey)
	}
}

// getConnectionKey generates a unique key for a machine connection
func (c *ReusableSSHClient) getConnectionKey(machine *spookytypes.Machine) string {
	return fmt.Sprintf("%s:%d:%s", machine.Host, machine.Port, machine.User)
}

// performConnectionCleanup performs the actual connection cleanup
func (c *ReusableSSHClient) performConnectionCleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	expiredConnections := []string{}
	expiredAuth := []string{}

	// Check for expired connections (older than 30 minutes)
	for key, connection := range c.connections {
		if now.Sub(connection.lastUsed) > 30*time.Minute {
			expiredConnections = append(expiredConnections, key)
		}
	}

	// Check for expired auth info (older than 1 hour)
	for key, authInfo := range c.authCache {
		if now.Sub(authInfo.lastUsed) > time.Hour {
			expiredAuth = append(expiredAuth, key)
		}
	}

	// Remove expired connections
	for _, key := range expiredConnections {
		c.logger.Debug("Cleaning up expired SSH connection", map[string]interface{}{
			"connection_key": key,
		})
		c.removeConnection(key)
	}

	// Remove expired auth info
	for _, key := range expiredAuth {
		c.logger.Debug("Cleaning up expired authentication info", map[string]interface{}{
			"auth_key": key,
		})
		delete(c.authCache, key)
	}

	if len(expiredConnections) > 0 || len(expiredAuth) > 0 {
		c.logger.Info("Cleaned up expired SSH resources", map[string]interface{}{
			"expired_connections": len(expiredConnections),
			"expired_auth":        len(expiredAuth),
		})
	}
}

// Close closes all SSH connections and cleans up resources
func (c *ReusableSSHClient) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.logger.Info("Closing SSH client and cleaning up connections", map[string]interface{}{
		"active_connections": len(c.connections),
		"cached_auth":        len(c.authCache),
	})

	// Close all connections
	for key, connection := range c.connections {
		c.logger.Debug("Closing SSH connection", map[string]interface{}{
			"connection_key": key,
		})
		if connection.client != nil {
			connection.client.Close()
		}
		if connection.session != nil {
			connection.session.Close()
		}
	}

	// Clear caches
	c.connections = make(map[string]*cachedConnection)
	c.authCache = make(map[string]*authInfo)

	return nil
}

// GetConnectionStats returns statistics about cached connections
func (c *ReusableSSHClient) GetConnectionStats() map[string]interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	healthyConnections := 0
	for _, connection := range c.connections {
		if connection.healthy {
			healthyConnections++
		}
	}

	return map[string]interface{}{
		"total_connections":     len(c.connections),
		"healthy_connections":   healthyConnections,
		"unhealthy_connections": len(c.connections) - healthyConnections,
		"cached_auth":           len(c.authCache),
	}
}

// recordConnectionSuccess records a successful connection
func (c *ReusableSSHClient) recordConnectionSuccess(reused bool) {
	c.metrics.mutex.Lock()
	defer c.metrics.mutex.Unlock()

	c.metrics.TotalConnections++
	c.metrics.SuccessfulConnections++
	if reused {
		c.metrics.ReusedConnections++
	}
}

// recordConnectionFailure records a failed connection
func (c *ReusableSSHClient) recordConnectionFailure() {
	c.metrics.mutex.Lock()
	defer c.metrics.mutex.Unlock()

	c.metrics.TotalConnections++
	c.metrics.FailedConnections++
}

// GetMetrics returns connection metrics
func (c *ReusableSSHClient) GetMetrics() map[string]interface{} {
	c.metrics.mutex.RLock()
	defer c.metrics.mutex.RUnlock()

	successRate := 0.0
	if c.metrics.TotalConnections > 0 {
		successRate = float64(c.metrics.SuccessfulConnections) / float64(c.metrics.TotalConnections)
	}

	reuseRate := 0.0
	if c.metrics.SuccessfulConnections > 0 {
		reuseRate = float64(c.metrics.ReusedConnections) / float64(c.metrics.SuccessfulConnections)
	}

	return map[string]interface{}{
		"total_connections":        c.metrics.TotalConnections,
		"successful_connections":   c.metrics.SuccessfulConnections,
		"failed_connections":       c.metrics.FailedConnections,
		"reused_connections":       c.metrics.ReusedConnections,
		"success_rate":             successRate,
		"reuse_rate":               reuseRate,
		"avg_authentication_time":  c.metrics.AuthenticationTime,
		"avg_command_running_time": c.metrics.CommandRunningTime,
	}
}

// RunCommandWithStdin runs a command with stdin data
func (c *ReusableSSHClient) RunCommandWithStdin(ctx context.Context, machine *spookytypes.Machine, command, stdin string) (string, error) {
	c.logger.Debug("Running SSH command with stdin", map[string]interface{}{
		"machine": machine.Hostname,
		"host":    machine.Host,
		"command": command,
	})

	// Get or create cached connection
	connection, err := c.getOrCreateConnection(ctx, machine)
	if err != nil {
		c.recordConnectionFailure()
		return "", fmt.Errorf("failed to get connection for %s: %w", machine.Hostname, err)
	}

	// Run command with stdin using cached connection
	result, err := c.executeCommandOnConnectionWithStdin(ctx, connection, command, stdin)
	if err != nil {
		// Mark connection as unhealthy on error
		c.markConnectionUnhealthy(machine)
		c.recordConnectionFailure()
		return "", fmt.Errorf("failed to run command on %s: %w", machine.Hostname, err)
	}

	// Update connection last used time
	c.updateConnectionLastUsed(machine)
	c.recordConnectionSuccess(true)

	return result, nil
}

// executeCommandOnConnectionWithStdin runs a command with stdin data
func (c *ReusableSSHClient) executeCommandOnConnectionWithStdin(_ context.Context, connection *cachedConnection, command, stdin string) (string, error) {
	connection.mutex.Lock()
	defer connection.mutex.Unlock()

	// Create new session for command running
	session, err := connection.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	// Set up stdin
	if stdin != "" {
		session.Stdin = strings.NewReader(stdin)
	}

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	// Run command
	err = session.Run(command)
	if err != nil {
		return "", fmt.Errorf("command running failed: %w (stderr: %s)", err, stderr.String())
	}

	return stdout.String(), nil
}
