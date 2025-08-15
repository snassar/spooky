# Connection Health Monitoring

## Overview

The spooky codebase currently lacks comprehensive connection health monitoring for SSH connections and other network operations. This document outlines recommendations for implementing connection health monitoring to improve reliability, performance, and operational visibility.

## Current State

- No proactive connection health checking
- Limited visibility into connection state
- No early detection of connection degradation
- Reactive connection failure handling only
- No connection metrics or monitoring

## Recommendations

### 1. Connection Health Interface

Define a comprehensive health monitoring interface:

```go
type ConnectionHealth interface {
    IsHealthy() bool
    GetHealthStatus() HealthStatus
    GetMetrics() ConnectionMetrics
    GetLastCheck() time.Time
    GetUptime() time.Duration
    GetFailureCount() int
}

type HealthStatus struct {
    Status      HealthState `json:"status"`
    LastCheck   time.Time   `json:"last_check"`
    Uptime      time.Duration `json:"uptime"`
    Failures    int         `json:"failures"`
    Latency     time.Duration `json:"latency"`
    Error       string      `json:"error,omitempty"`
}

type HealthState string

const (
    HealthStateHealthy   HealthState = "healthy"
    HealthStateDegraded  HealthState = "degraded"
    HealthStateUnhealthy HealthState = "unhealthy"
    HealthStateUnknown   HealthState = "unknown"
)

type ConnectionMetrics struct {
    TotalConnections    int64         `json:"total_connections"`
    ActiveConnections   int64         `json:"active_connections"`
    FailedConnections   int64         `json:"failed_connections"`
    AverageLatency      time.Duration `json:"average_latency"`
    LastSuccessTime     time.Time     `json:"last_success_time"`
    LastFailureTime     time.Time     `json:"last_failure_time"`
    SuccessRate         float64       `json:"success_rate"`
}
```

### 2. Health Check Implementation

Implement proactive health checking:

```go
type HealthChecker struct {
    config     HealthConfig
    logger     spookylogging.Logger
    metrics    ConnectionMetrics
    mutex      sync.RWMutex
}

type HealthConfig struct {
    CheckInterval    time.Duration `hcl:"check_interval"`
    Timeout         time.Duration `hcl:"timeout"`
    FailureThreshold int           `hcl:"failure_threshold"`
    SuccessThreshold int           `hcl:"success_threshold"`
    HealthCommands   []string      `hcl:"health_commands"`
}

func (hc *HealthChecker) StartHealthMonitoring(conn *spookyssh.Connection) {
    ticker := time.NewTicker(hc.config.CheckInterval)
    defer ticker.Stop()
    
    for range ticker.C {
        if err := hc.performHealthCheck(conn); err != nil {
            hc.logger.Warn("Health check failed", "error", err)
            hc.recordFailure()
        } else {
            hc.recordSuccess()
        }
    }
}

func (hc *HealthChecker) performHealthCheck(conn *spookyssh.Connection) error {
    ctx, cancel := context.WithTimeout(context.Background(), hc.config.Timeout)
    defer cancel()
    
    start := time.Now()
    
    // Perform simple health check command
    for _, cmd := range hc.config.HealthCommands {
        if err := hc.executeHealthCommand(ctx, conn, cmd); err != nil {
            return fmt.Errorf("health command '%s' failed: %w", cmd, err)
        }
    }
    
    latency := time.Since(start)
    hc.updateLatency(latency)
    
    return nil
}

func (hc *HealthChecker) executeHealthCommand(ctx context.Context, conn *spookyssh.Connection, cmd string) error {
    session, err := conn.NewSession()
    if err != nil {
        return fmt.Errorf("failed to create session: %w", err)
    }
    defer session.Close()
    
    // Set timeout for command execution
    session.SetDeadline(time.Now().Add(hc.config.Timeout))
    
    // Execute simple health check command
    if err := session.Run(cmd); err != nil {
        return fmt.Errorf("command execution failed: %w", err)
    }
    
    return nil
}
```

### 3. Connection Pool Health Monitoring

Integrate health monitoring into connection pool:

```go
type HealthAwareConnectionPool struct {
    connections map[string]*ManagedConnection
    healthChecker *HealthChecker
    mutex       sync.RWMutex
    logger      spookylogging.Logger
}

type ManagedConnection struct {
    connection *spookyssh.Connection
    health     *HealthChecker
    created    time.Time
    lastUsed   time.Time
    metrics    ConnectionMetrics
}

func (hcp *HealthAwareConnectionPool) GetHealthyConnection(host string, config *spookytypes.SSHConfig) (*spookyssh.Connection, error) {
    hcp.mutex.RLock()
    if managedConn, exists := hcp.connections[host]; exists {
        if managedConn.health.IsHealthy() {
            managedConn.lastUsed = time.Now()
            hcp.mutex.RUnlock()
            return managedConn.connection, nil
        }
        
        // Connection exists but is unhealthy, remove it
        hcp.logger.Warn("Removing unhealthy connection", "host", host)
        hcp.mutex.RUnlock()
        hcp.removeConnection(host)
    } else {
        hcp.mutex.RUnlock()
    }
    
    // Create new healthy connection
    return hcp.createHealthyConnection(host, config)
}

func (hcp *HealthAwareConnectionPool) createHealthyConnection(host string, config *spookytypes.SSHConfig) (*spookyssh.Connection, error) {
    // Create new connection
    conn, err := spookyssh.NewConnection(host, config)
    if err != nil {
        return nil, fmt.Errorf("failed to create connection: %w", err)
    }
    
    // Create health checker for this connection
    healthChecker := &HealthChecker{
        config: hcp.healthChecker.config,
        logger: hcp.logger,
    }
    
    // Start health monitoring in background
    go healthChecker.StartHealthMonitoring(conn)
    
    // Store managed connection
    managedConn := &ManagedConnection{
        connection: conn,
        health:     healthChecker,
        created:    time.Now(),
        lastUsed:   time.Now(),
    }
    
    hcp.mutex.Lock()
    hcp.connections[host] = managedConn
    hcp.mutex.Unlock()
    
    return conn, nil
}
```

### 4. Health Metrics Collection

Implement comprehensive metrics collection:

```go
type HealthMetricsCollector struct {
    metrics map[string]*ConnectionMetrics
    mutex   sync.RWMutex
    logger  spookylogging.Logger
}

func (hmc *HealthMetricsCollector) RecordConnection(host string, success bool, latency time.Duration) {
    hmc.mutex.Lock()
    defer hmc.mutex.Unlock()
    
    if hmc.metrics[host] == nil {
        hmc.metrics[host] = &ConnectionMetrics{}
    }
    
    metrics := hmc.metrics[host]
    metrics.TotalConnections++
    
    if success {
        metrics.ActiveConnections++
        metrics.LastSuccessTime = time.Now()
        metrics.AverageLatency = hmc.calculateAverageLatency(metrics.AverageLatency, latency, metrics.TotalConnections)
    } else {
        metrics.FailedConnections++
        metrics.LastFailureTime = time.Now()
    }
    
    metrics.SuccessRate = float64(metrics.ActiveConnections) / float64(metrics.TotalConnections)
}

func (hmc *HealthMetricsCollector) GetMetrics(host string) (*ConnectionMetrics, error) {
    hmc.mutex.RLock()
    defer hmc.mutex.RUnlock()
    
    if metrics, exists := hmc.metrics[host]; exists {
        return metrics, nil
    }
    
    return nil, fmt.Errorf("no metrics found for host: %s", host)
}

func (hmc *HealthMetricsCollector) GetAllMetrics() map[string]*ConnectionMetrics {
    hmc.mutex.RLock()
    defer hmc.mutex.RUnlock()
    
    result := make(map[string]*ConnectionMetrics)
    for host, metrics := range hmc.metrics {
        result[host] = metrics
    }
    
    return result
}
```

### 5. Health Status Reporting

Implement health status reporting and alerting:

```go
type HealthReporter struct {
    collector *HealthMetricsCollector
    logger    spookylogging.Logger
    alerts    []HealthAlert
}

type HealthAlert struct {
    Host      string      `json:"host"`
    Status    HealthState `json:"status"`
    Message   string      `json:"message"`
    Timestamp time.Time   `json:"timestamp"`
    Severity  string      `json:"severity"`
}

func (hr *HealthReporter) GenerateHealthReport() *HealthReport {
    metrics := hr.collector.GetAllMetrics()
    
    report := &HealthReport{
        Timestamp: time.Now(),
        Summary:   hr.generateSummary(metrics),
        Details:   hr.generateDetails(metrics),
        Alerts:    hr.alerts,
    }
    
    return report
}

type HealthReport struct {
    Timestamp time.Time              `json:"timestamp"`
    Summary   HealthSummary          `json:"summary"`
    Details   map[string]HealthStatus `json:"details"`
    Alerts    []HealthAlert          `json:"alerts"`
}

type HealthSummary struct {
    TotalConnections    int     `json:"total_connections"`
    HealthyConnections  int     `json:"healthy_connections"`
    DegradedConnections int     `json:"degraded_connections"`
    UnhealthyConnections int    `json:"unhealthy_connections"`
    OverallHealth       float64 `json:"overall_health"`
}
```

## Implementation Plan

### Phase 1: Core Health Monitoring
1. Implement health checking interface
2. Create basic health check commands
3. Add health status tracking
4. Implement metrics collection

### Phase 2: Connection Pool Integration
1. Integrate health monitoring into connection pool
2. Add automatic unhealthy connection removal
3. Implement health-aware connection selection
4. Add connection lifecycle management

### Phase 3: Advanced Monitoring
1. Add health reporting and alerting
2. Implement health dashboards
3. Add predictive health analysis
4. Create health-based load balancing

## Benefits

- **Proactive Problem Detection**: Identify issues before they cause failures
- **Improved Reliability**: Better connection management and failover
- **Operational Visibility**: Clear view of connection health and performance
- **Reduced Downtime**: Faster detection and resolution of connection issues
- **Performance Optimization**: Better resource utilization based on health data

## Risks and Mitigation

### Risks
- Increased overhead from health checking
- Potential for false positive health failures
- Complexity in health check configuration

### Mitigation
- Configurable health check intervals and timeouts
- Sophisticated health check algorithms
- Clear health check documentation and examples
- Monitoring of health check overhead

## Success Metrics

- Reduced connection failure rates
- Faster detection of connection issues
- Improved system availability
- Better resource utilization
- Reduced manual intervention requirements

## Related Documentation

- [SSH Implementation](mdc:ssh-implementation)
- [Connection Retry Logic](mdc:connection-retry-logic)
- [Error Handling Standards](mdc:error-handling-standards)
