// Package ssh provides SSH connection pooling functionality for the spooky codebase.
package ssh

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	spookytypes "spooky/internal/types"
	spookytypeslogging "spooky/internal/types/logging"
)

// ConnectionPool provides SSH connection pooling functionality
type ConnectionPool struct {
	connections map[string]*PooledConnection
	mu          sync.RWMutex
	config      *spookytypes.ClientConfig
	logger      spookytypeslogging.Logger
	closed      bool
}

// PooledConnection is defined in client.go

// NewConnectionPool creates a new SSH connection pool
func NewConnectionPool(config *spookytypes.ClientConfig, logger spookytypeslogging.Logger) *ConnectionPool {
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

	pool := &ConnectionPool{
		connections: make(map[string]*PooledConnection),
		config:      config,
		logger:      logger,
	}

	// Start background cleanup routine
	go pool.cleanupRoutine()

	return pool
}

// GetConnection retrieves or creates a connection from the pool
func (p *ConnectionPool) GetConnection(host string, port int, user string) (*PooledConnection, error) {
	if p.closed {
		return nil, fmt.Errorf("connection pool is closed")
	}

	connectionKey := fmt.Sprintf("%s:%d:%s", host, port, user)

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
		return nil, err
	}

	return conn, nil
}

// ReturnConnection returns a connection to the pool
func (p *ConnectionPool) ReturnConnection(conn *PooledConnection) {
	if p.closed || conn == nil {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	conn.IsIdle = true
	conn.LastUsed = time.Now()
}

// Close closes all connections in the pool
func (p *ConnectionPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.closed = true

	var errors []error
	for key, conn := range p.connections {
		if err := conn.Client.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close connection %s: %w", key, err))
		}
	}

	// Clear the connections map
	p.connections = make(map[string]*PooledConnection)

	if len(errors) > 0 {
		return fmt.Errorf("errors closing connections: %v", errors)
	}

	return nil
}

// GetStats returns pool statistics
func (p *ConnectionPool) GetStats() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	total := len(p.connections)
	active := 0
	idle := 0
	healthy := 0

	for _, conn := range p.connections {
		if conn.IsIdle {
			idle++
		} else {
			active++
		}
		if conn.IsHealthy {
			healthy++
		}
	}

	return map[string]interface{}{
		"total_connections":   total,
		"active_connections":  active,
		"idle_connections":    idle,
		"healthy_connections": healthy,
		"max_connections":     p.config.MaxConnections,
		"pool_utilization":    float64(total) / float64(p.config.MaxConnections),
	}
}

// Helper methods

// isConnectionHealthy checks if a connection is healthy
func (p *ConnectionPool) isConnectionHealthy(conn *PooledConnection) bool {
	if conn == nil {
		return false
	}

	// Check if connection is too old
	if time.Since(conn.CreatedAt) > p.config.IdleTimeout*2 {
		conn.IsHealthy = false
		return false
	}

	// Check if connection has too many errors
	if conn.ErrorCount > 3 {
		conn.IsHealthy = false
		return false
	}

	// Test connection with keepalive
	if _, _, err := conn.Client.SendRequest("keepalive@openssh.com", true, nil); err != nil {
		conn.ErrorCount++
		conn.IsHealthy = false
		return false
	}

	conn.IsHealthy = true
	return true
}

// updateConnectionUsage updates connection usage statistics
func (p *ConnectionPool) updateConnectionUsage(conn *PooledConnection) {
	p.mu.Lock()
	defer p.mu.Unlock()

	conn.LastUsed = time.Now()
	conn.UseCount++
	conn.IsIdle = false
}

// removeConnection removes a connection from the pool
func (p *ConnectionPool) removeConnection(connectionKey string) {
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
	}
}

// createNewConnection creates a new SSH connection
func (p *ConnectionPool) createNewConnection(host string, port int, user string) (*PooledConnection, error) {
	startTime := time.Now()

	// Create SSH config
	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // In production, use proper host key verification
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
	connectionKey := fmt.Sprintf("%s:%d:%s", host, port, user)
	p.connections[connectionKey] = conn
	p.mu.Unlock()

	p.logger.Info("Created new SSH connection", map[string]interface{}{
		"host":         host,
		"port":         port,
		"user":         user,
		"connect_time": connectTime,
		"pool_size":    len(p.connections),
	})

	return conn, nil
}

// cleanupRoutine runs background cleanup of idle connections
func (p *ConnectionPool) cleanupRoutine() {
	ticker := time.NewTicker(p.config.IdleTimeout / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.cleanupIdleConnections()
		}
	}
}

// cleanupIdleConnections removes idle connections
func (p *ConnectionPool) cleanupIdleConnections() {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	removedCount := 0

	for key, conn := range p.connections {
		if conn.IsIdle && now.Sub(conn.LastUsed) > p.config.IdleTimeout {
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
		}
	}

	if removedCount > 0 {
		p.logger.Info("Cleaned up idle connections", map[string]interface{}{
			"removed_count": removedCount,
			"remaining":     len(p.connections),
		})
	}
}

// ValidateConnection validates connection parameters
func (p *ConnectionPool) ValidateConnection(host string, port int, user string) error {
	if host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	if port <= 0 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	if user == "" {
		return fmt.Errorf("user cannot be empty")
	}

	return nil
}

// TestConnection tests a connection without adding it to the pool
func (p *ConnectionPool) TestConnection(host string, port int, user string) error {
	// Validate parameters
	if err := p.ValidateConnection(host, port, user); err != nil {
		return err
	}

	// Create temporary connection for testing
	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         p.config.DefaultTimeout,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}

	// Close the test connection
	client.Close()

	return nil
}
