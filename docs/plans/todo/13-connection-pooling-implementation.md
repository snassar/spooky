# Implementation Plan: SSH Connection Pooling Implementation

## Overview
Implement robust SSH connection pooling to replace placeholder implementations with real connection management, including connection reuse, health checking, and automatic reconnection.

## Task Details
- **Task ID**: 7.6
- **Priority**: High
- **Files**: 
  - `internal/ssh/connection_pool/manager.go`
- **Functions**: Connection pooling, health checking, reconnection

## Current State Analysis

### Existing Patterns
1. **SSH Manager**: Basic SSH connection management exists
2. **Connection Types**: SSH connection types defined
3. **Error Handling**: Consistent error wrapping
4. **Logging**: Structured logging implemented

### Current Placeholder Code
```go
// internal/ssh/connection_pool/manager.go line 67
// For now, we'll create a placeholder connection
```

## Implementation Requirements

### Interface Compliance
The connection pooling system must:
1. **Manage connection pools** for multiple machines
2. **Implement connection reuse** and lifecycle management
3. **Provide health checking** and connection validation
4. **Support automatic reconnection** on failures
5. **Handle connection limits** and resource management
6. **Implement connection timeouts** and cleanup
7. **Support connection metrics** and monitoring

### Required Dependencies
- SSH connection management
- Health checking system
- Metrics and monitoring
- Resource management system

## Detailed Implementation Plan

### Step 1: Implement Connection Pool Structure

#### 1.1 Pool Configuration
```go
// internal/ssh/connection_pool/config.go
package connection_pool

import (
    "time"
)

// PoolConfig represents connection pool configuration
type PoolConfig struct {
    MaxConnections     int           `json:"max_connections"`
    MaxIdleConnections int           `json:"max_idle_connections"`
    ConnectionTimeout  time.Duration `json:"connection_timeout"`
    IdleTimeout        time.Duration `json:"idle_timeout"`
    HealthCheckInterval time.Duration `json:"health_check_interval"`
    MaxRetries         int           `json:"max_retries"`
    RetryDelay         time.Duration `json:"retry_delay"`
    EnableMetrics      bool          `json:"enable_metrics"`
}

// DefaultPoolConfig returns default pool configuration
func DefaultPoolConfig() *PoolConfig {
    return &PoolConfig{
        MaxConnections:      100,
        MaxIdleConnections:  10,
        ConnectionTimeout:   30 * time.Second,
        IdleTimeout:         5 * time.Minute,
        HealthCheckInterval: 1 * time.Minute,
        MaxRetries:          3,
        RetryDelay:          5 * time.Second,
        EnableMetrics:       true,
    }
}

// Validate validates pool configuration
func (pc *PoolConfig) Validate() error {
    if pc.MaxConnections <= 0 {
        return fmt.Errorf("max_connections must be positive")
    }
    if pc.MaxIdleConnections > pc.MaxConnections {
        return fmt.Errorf("max_idle_connections cannot exceed max_connections")
    }
    if pc.ConnectionTimeout <= 0 {
        return fmt.Errorf("connection_timeout must be positive")
    }
    if pc.IdleTimeout <= 0 {
        return fmt.Errorf("idle_timeout must be positive")
    }
    if pc.HealthCheckInterval <= 0 {
        return fmt.Errorf("health_check_interval must be positive")
    }
    if pc.MaxRetries < 0 {
        return fmt.Errorf("max_retries cannot be negative")
    }
    if pc.RetryDelay <= 0 {
        return fmt.Errorf("retry_delay must be positive")
    }
    return nil
}
```

#### 1.2 Connection Pool Entry
```go
// internal/ssh/connection_pool/connection.go
package connection_pool

import (
    "sync"
    "time"
    
    "spooky/internal/ssh/types"
)

// PooledConnection represents a pooled SSH connection
type PooledConnection struct {
    connection     types.Connection
    machineID      string
    lastUsed       time.Time
    lastHealthCheck time.Time
    useCount       int64
    inUse          bool
    healthy        bool
    mutex          sync.RWMutex
    config         *PoolConfig
}

// NewPooledConnection creates a new pooled connection
func NewPooledConnection(connection types.Connection, machineID string, config *PoolConfig) *PooledConnection {
    return &PooledConnection{
        connection:      connection,
        machineID:       machineID,
        lastUsed:        time.Now(),
        lastHealthCheck: time.Now(),
        useCount:        0,
        inUse:           false,
        healthy:         true,
        config:          config,
    }
}

// GetConnection returns the underlying connection
func (pc *PooledConnection) GetConnection() types.Connection {
    return pc.connection
}

// MarkInUse marks the connection as in use
func (pc *PooledConnection) MarkInUse() {
    pc.mutex.Lock()
    defer pc.mutex.Unlock()
    
    pc.inUse = true
    pc.lastUsed = time.Now()
    pc.useCount++
}

// MarkIdle marks the connection as idle
func (pc *PooledConnection) MarkIdle() {
    pc.mutex.Lock()
    defer pc.mutex.Unlock()
    
    pc.inUse = false
    pc.lastUsed = time.Now()
}

// IsIdle checks if connection is idle
func (pc *PooledConnection) IsIdle() bool {
    pc.mutex.RLock()
    defer pc.mutex.RUnlock()
    
    return !pc.inUse && time.Since(pc.lastUsed) > pc.config.IdleTimeout
}

// IsHealthy checks if connection is healthy
func (pc *PooledConnection) IsHealthy() bool {
    pc.mutex.RLock()
    defer pc.mutex.RUnlock()
    
    return pc.healthy
}

// SetHealthy sets connection health status
func (pc *PooledConnection) SetHealthy(healthy bool) {
    pc.mutex.Lock()
    defer pc.mutex.Unlock()
    
    pc.healthy = healthy
    pc.lastHealthCheck = time.Now()
}

// NeedsHealthCheck checks if connection needs health check
func (pc *PooledConnection) NeedsHealthCheck() bool {
    pc.mutex.RLock()
    defer pc.mutex.RUnlock()
    
    return time.Since(pc.lastHealthCheck) > pc.config.HealthCheckInterval
}

// GetStats returns connection statistics
func (pc *PooledConnection) GetStats() ConnectionStats {
    pc.mutex.RLock()
    defer pc.mutex.RUnlock()
    
    return ConnectionStats{
        MachineID:        pc.machineID,
        InUse:            pc.inUse,
        Healthy:          pc.healthy,
        UseCount:         pc.useCount,
        LastUsed:         pc.lastUsed,
        LastHealthCheck:  pc.lastHealthCheck,
        IdleDuration:     time.Since(pc.lastUsed),
    }
}

// Close closes the connection
func (pc *PooledConnection) Close() error {
    pc.mutex.Lock()
    defer pc.mutex.Unlock()
    
    pc.inUse = false
    pc.healthy = false
    return pc.connection.Close()
}

// ConnectionStats represents connection statistics
type ConnectionStats struct {
    MachineID        string        `json:"machine_id"`
    InUse            bool          `json:"in_use"`
    Healthy          bool          `json:"healthy"`
    UseCount         int64         `json:"use_count"`
    LastUsed         time.Time     `json:"last_used"`
    LastHealthCheck  time.Time     `json:"last_health_check"`
    IdleDuration     time.Duration `json:"idle_duration"`
}
```

### Step 2: Implement Connection Pool Manager

#### 2.1 Pool Manager Implementation
```go
// internal/ssh/connection_pool/manager.go
package connection_pool

import (
    "context"
    "fmt"
    "sync"
    "time"
    
    "spooky/internal/ssh"
    "spooky/internal/ssh/types"
    "spooky/internal/logging"
)

// Manager manages SSH connection pools
type Manager struct {
    pools         map[string]*ConnectionPool
    config        *PoolConfig
    sshManager    ssh.Manager
    logger        logging.Logger
    mutex         sync.RWMutex
    metrics       *PoolMetrics
    healthChecker *HealthChecker
}

// NewManager creates a new connection pool manager
func NewManager(sshManager ssh.Manager, config *PoolConfig, logger logging.Logger) (*Manager, error) {
    if err := config.Validate(); err != nil {
        return nil, fmt.Errorf("invalid pool configuration: %w", err)
    }

    manager := &Manager{
        pools:      make(map[string]*ConnectionPool),
        config:     config,
        sshManager: sshManager,
        logger:     logger,
        metrics:    NewPoolMetrics(),
    }

    // Initialize health checker
    manager.healthChecker = NewHealthChecker(manager, config, logger)

    // Start health checking
    go manager.healthChecker.Start()

    return manager, nil
}

// GetConnection gets a connection from the pool
func (m *Manager) GetConnection(ctx context.Context, machineID string, sshConfig *types.SSHConfig) (types.Connection, error) {
    m.mutex.Lock()
    defer m.mutex.Unlock()

    // Get or create pool for machine
    pool, exists := m.pools[machineID]
    if !exists {
        pool = NewConnectionPool(machineID, m.config, m.logger)
        m.pools[machineID] = pool
    }

    // Get connection from pool
    pooledConn, err := pool.GetConnection(ctx, m.sshManager, sshConfig)
    if err != nil {
        m.metrics.IncrementConnectionErrors(machineID)
        return nil, fmt.Errorf("failed to get connection from pool: %w", err)
    }

    m.metrics.IncrementConnectionRequests(machineID)
    return pooledConn.GetConnection(), nil
}

// ReturnConnection returns a connection to the pool
func (m *Manager) ReturnConnection(machineID string, connection types.Connection) error {
    m.mutex.RLock()
    pool, exists := m.pools[machineID]
    m.mutex.RUnlock()

    if !exists {
        return fmt.Errorf("pool not found for machine: %s", machineID)
    }

    return pool.ReturnConnection(connection)
}

// CloseConnection closes a connection and removes it from pool
func (m *Manager) CloseConnection(machineID string, connection types.Connection) error {
    m.mutex.RLock()
    pool, exists := m.pools[machineID]
    m.mutex.RUnlock()

    if !exists {
        return fmt.Errorf("pool not found for machine: %s", machineID)
    }

    return pool.CloseConnection(connection)
}

// ClosePool closes all connections in a pool
func (m *Manager) ClosePool(machineID string) error {
    m.mutex.Lock()
    defer m.mutex.Unlock()

    pool, exists := m.pools[machineID]
    if !exists {
        return fmt.Errorf("pool not found for machine: %s", machineID)
    }

    if err := pool.Close(); err != nil {
        return fmt.Errorf("failed to close pool: %w", err)
    }

    delete(m.pools, machineID)
    m.logger.Info("Connection pool closed",
        logging.String("machine_id", machineID))

    return nil
}

// CloseAll closes all connection pools
func (m *Manager) CloseAll() error {
    m.mutex.Lock()
    defer m.mutex.Unlock()

    var errors []error
    for machineID, pool := range m.pools {
        if err := pool.Close(); err != nil {
            errors = append(errors, fmt.Errorf("failed to close pool %s: %w", machineID, err))
        }
    }

    m.pools = make(map[string]*ConnectionPool)
    m.logger.Info("All connection pools closed")

    if len(errors) > 0 {
        return fmt.Errorf("errors closing pools: %v", errors)
    }

    return nil
}

// GetPoolStats returns pool statistics
func (m *Manager) GetPoolStats() map[string]PoolStats {
    m.mutex.RLock()
    defer m.mutex.RUnlock()

    stats := make(map[string]PoolStats)
    for machineID, pool := range m.pools {
        stats[machineID] = pool.GetStats()
    }

    return stats
}

// GetMetrics returns pool metrics
func (m *Manager) GetMetrics() *PoolMetrics {
    return m.metrics
}
```

#### 2.2 Connection Pool Implementation
```go
// internal/ssh/connection_pool/pool.go
package connection_pool

import (
    "context"
    "fmt"
    "sync"
    "time"
    
    "spooky/internal/ssh"
    "spooky/internal/ssh/types"
    "spooky/internal/logging"
)

// ConnectionPool manages connections for a single machine
type ConnectionPool struct {
    machineID string
    config    *PoolConfig
    logger    logging.Logger
    mutex     sync.RWMutex
    connections []*PooledConnection
    inUse      map[types.Connection]*PooledConnection
}

// NewConnectionPool creates a new connection pool
func NewConnectionPool(machineID string, config *PoolConfig, logger logging.Logger) *ConnectionPool {
    return &ConnectionPool{
        machineID:    machineID,
        config:       config,
        logger:       logger,
        connections:  make([]*PooledConnection, 0),
        inUse:        make(map[types.Connection]*PooledConnection),
    }
}

// GetConnection gets a connection from the pool
func (cp *ConnectionPool) GetConnection(ctx context.Context, sshManager ssh.Manager, sshConfig *types.SSHConfig) (*PooledConnection, error) {
    cp.mutex.Lock()
    defer cp.mutex.Unlock()

    // Try to find an available healthy connection
    for _, pooledConn := range cp.connections {
        if !pooledConn.inUse && pooledConn.healthy {
            pooledConn.MarkInUse()
            cp.inUse[pooledConn.connection] = pooledConn
            
            cp.logger.Debug("Reused connection from pool",
                logging.String("machine_id", cp.machineID))
            
            return pooledConn, nil
        }
    }

    // Create new connection if under limit
    if len(cp.connections) < cp.config.MaxConnections {
        connection, err := cp.createConnection(ctx, sshManager, sshConfig)
        if err != nil {
            return nil, fmt.Errorf("failed to create connection: %w", err)
        }

        pooledConn := NewPooledConnection(connection, cp.machineID, cp.config)
        pooledConn.MarkInUse()
        
        cp.connections = append(cp.connections, pooledConn)
        cp.inUse[connection] = pooledConn

        cp.logger.Debug("Created new connection",
            logging.String("machine_id", cp.machineID),
            logging.Int("pool_size", len(cp.connections)))

        return pooledConn, nil
    }

    // Wait for available connection
    return cp.waitForConnection(ctx)
}

// createConnection creates a new SSH connection
func (cp *ConnectionPool) createConnection(ctx context.Context, sshManager ssh.Manager, sshConfig *types.SSHConfig) (types.Connection, error) {
    var connection types.Connection
    var err error

    // Retry connection creation
    for attempt := 0; attempt <= cp.config.MaxRetries; attempt++ {
        if attempt > 0 {
            cp.logger.Debug("Retrying connection creation",
                logging.String("machine_id", cp.machineID),
                logging.Int("attempt", attempt+1))
            time.Sleep(cp.config.RetryDelay)
        }

        connection, err = sshManager.Connect(sshConfig.Host, sshConfig)
        if err == nil {
            break
        }
    }

    if err != nil {
        return nil, fmt.Errorf("failed to create connection after %d attempts: %w", cp.config.MaxRetries+1, err)
    }

    return connection, nil
}

// waitForConnection waits for an available connection
func (cp *ConnectionPool) waitForConnection(ctx context.Context) (*PooledConnection, error) {
    // This is a simplified implementation
    // In practice, this would use channels and proper waiting
    return nil, fmt.Errorf("no available connections and pool is full")
}

// ReturnConnection returns a connection to the pool
func (cp *ConnectionPool) ReturnConnection(connection types.Connection) error {
    cp.mutex.Lock()
    defer cp.mutex.Unlock()

    pooledConn, exists := cp.inUse[connection]
    if !exists {
        return fmt.Errorf("connection not found in pool")
    }

    pooledConn.MarkIdle()
    delete(cp.inUse, connection)

    cp.logger.Debug("Connection returned to pool",
        logging.String("machine_id", cp.machineID))

    return nil
}

// CloseConnection closes a connection and removes it from pool
func (cp *ConnectionPool) CloseConnection(connection types.Connection) error {
    cp.mutex.Lock()
    defer cp.mutex.Unlock()

    pooledConn, exists := cp.inUse[connection]
    if !exists {
        return fmt.Errorf("connection not found in pool")
    }

    // Remove from in-use map
    delete(cp.inUse, connection)

    // Remove from connections slice
    for i, conn := range cp.connections {
        if conn == pooledConn {
            cp.connections = append(cp.connections[:i], cp.connections[i+1:]...)
            break
        }
    }

    // Close connection
    if err := pooledConn.Close(); err != nil {
        return fmt.Errorf("failed to close connection: %w", err)
    }

    cp.logger.Debug("Connection closed and removed from pool",
        logging.String("machine_id", cp.machineID))

    return nil
}

// Close closes all connections in the pool
func (cp *ConnectionPool) Close() error {
    cp.mutex.Lock()
    defer cp.mutex.Unlock()

    var errors []error

    // Close all connections
    for _, pooledConn := range cp.connections {
        if err := pooledConn.Close(); err != nil {
            errors = append(errors, err)
        }
    }

    // Clear pools
    cp.connections = make([]*PooledConnection, 0)
    cp.inUse = make(map[types.Connection]*PooledConnection)

    if len(errors) > 0 {
        return fmt.Errorf("errors closing connections: %v", errors)
    }

    return nil
}

// GetStats returns pool statistics
func (cp *ConnectionPool) GetStats() PoolStats {
    cp.mutex.RLock()
    defer cp.mutex.RUnlock()

    totalConnections := len(cp.connections)
    inUseConnections := len(cp.inUse)
    idleConnections := totalConnections - inUseConnections

    return PoolStats{
        MachineID:         cp.machineID,
        TotalConnections:  totalConnections,
        InUseConnections:  inUseConnections,
        IdleConnections:   idleConnections,
        MaxConnections:    cp.config.MaxConnections,
        MaxIdleConnections: cp.config.MaxIdleConnections,
    }
}

// PoolStats represents pool statistics
type PoolStats struct {
    MachineID         string `json:"machine_id"`
    TotalConnections  int    `json:"total_connections"`
    InUseConnections  int    `json:"in_use_connections"`
    IdleConnections   int    `json:"idle_connections"`
    MaxConnections    int    `json:"max_connections"`
    MaxIdleConnections int   `json:"max_idle_connections"`
}
```

### Step 3: Implement Health Checking

#### 3.1 Health Checker Implementation
```go
// internal/ssh/connection_pool/health.go
package connection_pool

import (
    "sync"
    "time"
    
    "spooky/internal/logging"
)

// HealthChecker performs health checks on connections
type HealthChecker struct {
    manager *Manager
    config  *PoolConfig
    logger  logging.Logger
    stop    chan struct{}
    wg      sync.WaitGroup
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(manager *Manager, config *PoolConfig, logger logging.Logger) *HealthChecker {
    return &HealthChecker{
        manager: manager,
        config:  config,
        logger:  logger,
        stop:    make(chan struct{}),
    }
}

// Start starts the health checker
func (hc *HealthChecker) Start() {
    hc.wg.Add(1)
    go hc.run()
}

// Stop stops the health checker
func (hc *HealthChecker) Stop() {
    close(hc.stop)
    hc.wg.Wait()
}

// run runs the health checker loop
func (hc *HealthChecker) run() {
    defer hc.wg.Done()

    ticker := time.NewTicker(hc.config.HealthCheckInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            hc.performHealthChecks()
        case <-hc.stop:
            return
        }
    }
}

// performHealthChecks performs health checks on all pools
func (hc *HealthChecker) performHealthChecks() {
    hc.manager.mutex.RLock()
    pools := make(map[string]*ConnectionPool)
    for machineID, pool := range hc.manager.pools {
        pools[machineID] = pool
    }
    hc.manager.mutex.RUnlock()

    for machineID, pool := range pools {
        hc.checkPoolHealth(machineID, pool)
    }
}

// checkPoolHealth checks health of connections in a pool
func (hc *HealthChecker) checkPoolHealth(machineID string, pool *ConnectionPool) {
    pool.mutex.RLock()
    connections := make([]*PooledConnection, len(pool.connections))
    copy(connections, pool.connections)
    pool.mutex.RUnlock()

    for _, pooledConn := range connections {
        if pooledConn.NeedsHealthCheck() {
            hc.checkConnectionHealth(pooledConn)
        }
    }
}

// checkConnectionHealth checks health of a single connection
func (hc *HealthChecker) checkConnectionHealth(pooledConn *PooledConnection) {
    // Perform health check
    healthy := hc.performHealthCheck(pooledConn.GetConnection())
    pooledConn.SetHealthy(healthy)

    if !healthy {
        hc.logger.Warn("Connection health check failed",
            logging.String("machine_id", pooledConn.machineID))
    }
}

// performHealthCheck performs a health check on a connection
func (hc *HealthChecker) performHealthCheck(connection types.Connection) bool {
    // This would implement actual health check logic
    // For now, return true (healthy)
    return true
}
```

### Step 4: Implement Metrics

#### 4.1 Pool Metrics Implementation
```go
// internal/ssh/connection_pool/metrics.go
package connection_pool

import (
    "sync"
    "sync/atomic"
    "time"
)

// PoolMetrics tracks pool metrics
type PoolMetrics struct {
    connectionRequests int64
    connectionErrors   int64
    connectionReuses   int64
    connectionCreates  int64
    healthCheckPasses  int64
    healthCheckFails   int64
    mutex              sync.RWMutex
    machineMetrics     map[string]*MachineMetrics
}

// NewPoolMetrics creates new pool metrics
func NewPoolMetrics() *PoolMetrics {
    return &PoolMetrics{
        machineMetrics: make(map[string]*MachineMetrics),
    }
}

// IncrementConnectionRequests increments connection requests
func (pm *PoolMetrics) IncrementConnectionRequests(machineID string) {
    atomic.AddInt64(&pm.connectionRequests, 1)
    pm.getMachineMetrics(machineID).IncrementConnectionRequests()
}

// IncrementConnectionErrors increments connection errors
func (pm *PoolMetrics) IncrementConnectionErrors(machineID string) {
    atomic.AddInt64(&pm.connectionErrors, 1)
    pm.getMachineMetrics(machineID).IncrementConnectionErrors()
}

// IncrementConnectionReuses increments connection reuses
func (pm *PoolMetrics) IncrementConnectionReuses(machineID string) {
    atomic.AddInt64(&pm.connectionReuses, 1)
    pm.getMachineMetrics(machineID).IncrementConnectionReuses()
}

// IncrementConnectionCreates increments connection creates
func (pm *PoolMetrics) IncrementConnectionCreates(machineID string) {
    atomic.AddInt64(&pm.connectionCreates, 1)
    pm.getMachineMetrics(machineID).IncrementConnectionCreates()
}

// getMachineMetrics gets or creates machine metrics
func (pm *PoolMetrics) getMachineMetrics(machineID string) *MachineMetrics {
    pm.mutex.Lock()
    defer pm.mutex.Unlock()

    if metrics, exists := pm.machineMetrics[machineID]; exists {
        return metrics
    }

    metrics := NewMachineMetrics(machineID)
    pm.machineMetrics[machineID] = metrics
    return metrics
}

// GetMetrics returns current metrics
func (pm *PoolMetrics) GetMetrics() Metrics {
    pm.mutex.RLock()
    defer pm.mutex.RUnlock()

    machineMetrics := make(map[string]MachineMetrics)
    for machineID, metrics := range pm.machineMetrics {
        machineMetrics[machineID] = metrics.GetMetrics()
    }

    return Metrics{
        ConnectionRequests: atomic.LoadInt64(&pm.connectionRequests),
        ConnectionErrors:   atomic.LoadInt64(&pm.connectionErrors),
        ConnectionReuses:   atomic.LoadInt64(&pm.connectionReuses),
        ConnectionCreates:  atomic.LoadInt64(&pm.connectionCreates),
        HealthCheckPasses:  atomic.LoadInt64(&pm.healthCheckPasses),
        HealthCheckFails:   atomic.LoadInt64(&pm.healthCheckFails),
        MachineMetrics:     machineMetrics,
        Timestamp:          time.Now(),
    }
}

// Metrics represents pool metrics
type Metrics struct {
    ConnectionRequests int64                    `json:"connection_requests"`
    ConnectionErrors   int64                    `json:"connection_errors"`
    ConnectionReuses   int64                    `json:"connection_reuses"`
    ConnectionCreates  int64                    `json:"connection_creates"`
    HealthCheckPasses  int64                    `json:"health_check_passes"`
    HealthCheckFails   int64                    `json:"health_check_fails"`
    MachineMetrics     map[string]MachineMetrics `json:"machine_metrics"`
    Timestamp          time.Time                `json:"timestamp"`
}

// MachineMetrics tracks metrics for a specific machine
type MachineMetrics struct {
    machineID          string
    connectionRequests int64
    connectionErrors   int64
    connectionReuses   int64
    connectionCreates  int64
}

// NewMachineMetrics creates new machine metrics
func NewMachineMetrics(machineID string) *MachineMetrics {
    return &MachineMetrics{
        machineID: machineID,
    }
}

// IncrementConnectionRequests increments connection requests
func (mm *MachineMetrics) IncrementConnectionRequests() {
    atomic.AddInt64(&mm.connectionRequests, 1)
}

// IncrementConnectionErrors increments connection errors
func (mm *MachineMetrics) IncrementConnectionErrors() {
    atomic.AddInt64(&mm.connectionErrors, 1)
}

// IncrementConnectionReuses increments connection reuses
func (mm *MachineMetrics) IncrementConnectionReuses() {
    atomic.AddInt64(&mm.connectionReuses, 1)
}

// IncrementConnectionCreates increments connection creates
func (mm *MachineMetrics) IncrementConnectionCreates() {
    atomic.AddInt64(&mm.connectionCreates, 1)
}

// GetMetrics returns current machine metrics
func (mm *MachineMetrics) GetMetrics() MachineMetrics {
    return MachineMetrics{
        machineID:          mm.machineID,
        connectionRequests: atomic.LoadInt64(&mm.connectionRequests),
        connectionErrors:   atomic.LoadInt64(&mm.connectionErrors),
        connectionReuses:   atomic.LoadInt64(&mm.connectionReuses),
        connectionCreates:  atomic.LoadInt64(&mm.connectionCreates),
    }
}
```

## Configuration Options

### Supported Options
- **MaxConnections**: Maximum connections per pool
- **MaxIdleConnections**: Maximum idle connections
- **ConnectionTimeout**: Connection creation timeout
- **IdleTimeout**: Connection idle timeout
- **HealthCheckInterval**: Health check interval
- **MaxRetries**: Maximum connection retries
- **RetryDelay**: Delay between retries
- **EnableMetrics**: Enable/disable metrics

## Dependencies

### Internal Dependencies
- `spooky/internal/ssh`
- `spooky/internal/ssh/types`
- `spooky/internal/logging`

### External Dependencies
- `context` (standard library)
- `fmt` (standard library)
- `sync` (standard library)
- `time` (standard library)

## Implementation Order

1. Implement pool configuration
2. Add connection pool entry structure
3. Create connection pool manager
4. Implement connection pool
5. Add health checking system
6. Implement metrics tracking
7. Add comprehensive tests
8. Performance optimization
9. Documentation and cleanup
