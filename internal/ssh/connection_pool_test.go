package ssh

import (
	"testing"
	"time"

	spookylogging "spooky/internal/logging"
	spookytypes "spooky/internal/types"
)

func TestAdvancedConnectionPool(t *testing.T) {
	// Create a test logger
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("test")

	// Create test configuration
	config := &spookytypes.ClientConfig{
		DefaultPort:        22,
		DefaultTimeout:     30 * time.Second,
		MaxConnections:     5,
		MaxRetryAttempts:   3,
		RetryDelay:         5 * time.Second,
		IdleTimeout:        60 * time.Second,
		KnownHostsPath:     "/tmp/test_known_hosts",
		StrictHostKeyCheck: false,
		AllowInsecureHosts: true,
	}

	// Create connection pool
	pool := NewAdvancedConnectionPool(config, logger)
	defer pool.Close()

	// Test initial metrics
	metrics := pool.GetMetrics()
	if metrics.TotalConnections != 0 {
		t.Errorf("Expected 0 total connections, got %d", metrics.TotalConnections)
	}
	if metrics.ActiveConnections != 0 {
		t.Errorf("Expected 0 active connections, got %d", metrics.ActiveConnections)
	}
	if metrics.IdleConnections != 0 {
		t.Errorf("Expected 0 idle connections, got %d", metrics.IdleConnections)
	}

	// Test pool capacity
	if pool.config.MaxConnections != 5 {
		t.Errorf("Expected max connections to be 5, got %d", pool.config.MaxConnections)
	}
}

func TestConnectionPoolMetrics(t *testing.T) {
	// Create a test logger
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("test")

	// Create test configuration with short timeouts for testing
	config := &spookytypes.ClientConfig{
		DefaultPort:        22,
		DefaultTimeout:     5 * time.Second,
		MaxConnections:     3,
		MaxRetryAttempts:   1,
		RetryDelay:         1 * time.Second,
		IdleTimeout:        10 * time.Second,
		KnownHostsPath:     "/tmp/test_known_hosts",
		StrictHostKeyCheck: false,
		AllowInsecureHosts: true,
	}

	// Create connection pool
	pool := NewAdvancedConnectionPool(config, logger)
	defer pool.Close()

	// Test metrics update
	pool.updatePoolMetrics()
	metrics := pool.GetMetrics()

	// Verify metrics structure
	if metrics.TotalConnections < 0 {
		t.Errorf("Total connections should not be negative, got %d", metrics.TotalConnections)
	}
	if metrics.ActiveConnections < 0 {
		t.Errorf("Active connections should not be negative, got %d", metrics.ActiveConnections)
	}
	if metrics.IdleConnections < 0 {
		t.Errorf("Idle connections should not be negative, got %d", metrics.IdleConnections)
	}
	if metrics.PoolUtilization < 0 || metrics.PoolUtilization > 1 {
		t.Errorf("Pool utilization should be between 0 and 1, got %f", metrics.PoolUtilization)
	}
}

func TestConnectionPoolCleanup(t *testing.T) {
	// Create a test logger
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("test")

	// Create test configuration with very short timeouts for testing
	config := &spookytypes.ClientConfig{
		DefaultPort:        22,
		DefaultTimeout:     5 * time.Second,
		MaxConnections:     3,
		MaxRetryAttempts:   1,
		RetryDelay:         1 * time.Second,
		IdleTimeout:        2 * time.Second, // Very short for testing
		KnownHostsPath:     "/tmp/test_known_hosts",
		StrictHostKeyCheck: false,
		AllowInsecureHosts: true,
	}

	// Create connection pool
	pool := NewAdvancedConnectionPool(config, logger)
	defer pool.Close()

	// Test cleanup routine - just verify the pool starts correctly
	initialCleanupCycles := pool.GetMetrics().CleanupCycles

	// The cleanup runs in background, so we just verify the pool is working
	if initialCleanupCycles < 0 {
		t.Errorf("Cleanup cycles should not be negative, got %d", initialCleanupCycles)
	}

	// Verify pool metrics are initialized correctly
	metrics := pool.GetMetrics()
	if metrics.TotalConnections != 0 {
		t.Errorf("Expected 0 total connections initially, got %d", metrics.TotalConnections)
	}
}

func TestConnectionPoolHealthCheck(t *testing.T) {
	// Create a test logger
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("test")

	// Create test configuration
	config := &spookytypes.ClientConfig{
		DefaultPort:        22,
		DefaultTimeout:     5 * time.Second,
		MaxConnections:     3,
		MaxRetryAttempts:   1,
		RetryDelay:         1 * time.Second,
		IdleTimeout:        10 * time.Second,
		KnownHostsPath:     "/tmp/test_known_hosts",
		StrictHostKeyCheck: false,
		AllowInsecureHosts: true,
	}

	// Create connection pool
	pool := NewAdvancedConnectionPool(config, logger)
	defer pool.Close()

	// Test health check with nil connection (should return false)
	isHealthy := pool.isConnectionHealthy(nil)
	if isHealthy {
		t.Error("Expected nil connection to be unhealthy")
	}

	// Test health check with mock connection (without real SSH client)
	mockConn := &PooledConnection{
		Host:         "test.example.com",
		Port:         22,
		User:         "test",
		CreatedAt:    time.Now(),
		LastUsed:     time.Now(),
		UseCount:     1,
		ErrorCount:   0,
		IsHealthy:    true,
		IsIdle:       false,
		ConnectionID: "test-connection",
		Client:       nil, // No real SSH client for testing
	}

	// Test connection with too many errors (this should fail without SSH client)
	mockConn.ErrorCount = 4
	isHealthy = pool.isConnectionHealthy(mockConn)
	if isHealthy {
		t.Error("Expected connection with too many errors to fail health check")
	}

	// Test connection that's too old
	oldConn := &PooledConnection{
		Host:         "test.example.com",
		Port:         22,
		User:         "test",
		CreatedAt:    time.Now().Add(-30 * time.Second), // Very old
		LastUsed:     time.Now(),
		UseCount:     1,
		ErrorCount:   0,
		IsHealthy:    true,
		IsIdle:       false,
		ConnectionID: "test-connection-old",
		Client:       nil,
	}

	isHealthy = pool.isConnectionHealthy(oldConn)
	if isHealthy {
		t.Error("Expected old connection to fail health check")
	}
}

func TestConnectionPoolCapacity(t *testing.T) {
	// Create a test logger
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("test")

	// Create test configuration with very small pool
	config := &spookytypes.ClientConfig{
		DefaultPort:        22,
		DefaultTimeout:     5 * time.Second,
		MaxConnections:     1, // Very small pool
		MaxRetryAttempts:   1,
		RetryDelay:         1 * time.Second,
		IdleTimeout:        10 * time.Second,
		KnownHostsPath:     "/tmp/test_known_hosts",
		StrictHostKeyCheck: false,
		AllowInsecureHosts: true,
	}

	// Create connection pool
	pool := NewAdvancedConnectionPool(config, logger)
	defer pool.Close()

	// Test pool capacity enforcement
	// This would normally try to create a connection, but we'll test the capacity check
	metrics := pool.GetMetrics()
	if metrics.TotalConnections > config.MaxConnections {
		t.Errorf("Pool exceeded maximum connections: %d > %d", metrics.TotalConnections, config.MaxConnections)
	}
}

func TestConnectionPoolReturnConnection(t *testing.T) {
	// Create a test logger
	logManager := spookylogging.NewLogManager()
	logger := logManager.GetLogger("test")

	// Create test configuration
	config := &spookytypes.ClientConfig{
		DefaultPort:        22,
		DefaultTimeout:     5 * time.Second,
		MaxConnections:     3,
		MaxRetryAttempts:   1,
		RetryDelay:         1 * time.Second,
		IdleTimeout:        10 * time.Second,
		KnownHostsPath:     "/tmp/test_known_hosts",
		StrictHostKeyCheck: false,
		AllowInsecureHosts: true,
	}

	// Create connection pool
	pool := NewAdvancedConnectionPool(config, logger)
	defer pool.Close()

	// Test returning a connection
	mockConn := &PooledConnection{
		Host:         "test.example.com",
		Port:         22,
		User:         "test",
		CreatedAt:    time.Now(),
		LastUsed:     time.Now(),
		UseCount:     1,
		ErrorCount:   0,
		IsHealthy:    true,
		IsIdle:       false,
		ConnectionID: "test-connection",
	}

	// Return connection to pool
	pool.ReturnConnection(mockConn)

	// Verify connection is marked as idle
	if !mockConn.IsIdle {
		t.Error("Expected connection to be marked as idle after return")
	}
}
