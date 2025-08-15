# Connection Retry Logic with Exponential Backoff

## Overview

The spooky codebase currently lacks robust connection retry mechanisms for SSH connections and other network operations. This document outlines recommendations for implementing connection retry logic with exponential backoff to improve reliability and resilience.

## Current State

- SSH connections fail immediately on network issues
- No retry mechanism for transient network problems
- Connection timeouts are fixed and non-adaptive
- No backoff strategy for repeated connection attempts

## Recommendations

### 1. Exponential Backoff Algorithm

Implement exponential backoff with jitter to prevent thundering herd problems:

```go
type RetryConfig struct {
    MaxAttempts     int           `hcl:"max_attempts"`
    InitialDelay    time.Duration `hcl:"initial_delay"`
    MaxDelay        time.Duration `hcl:"max_delay"`
    BackoffFactor   float64       `hcl:"backoff_factor"`
    JitterFactor    float64       `hcl:"jitter_factor"`
    Timeout         time.Duration `hcl:"timeout"`
}

type RetryStrategy interface {
    NextDelay(attempt int) time.Duration
    ShouldRetry(err error, attempt int) bool
    IsRetryableError(err error) bool
}

type ExponentialBackoff struct {
    config RetryConfig
}

func (eb *ExponentialBackoff) NextDelay(attempt int) time.Duration {
    if attempt <= 0 {
        return eb.config.InitialDelay
    }
    
    // Calculate exponential delay
    delay := float64(eb.config.InitialDelay) * math.Pow(eb.config.BackoffFactor, float64(attempt-1))
    
    // Apply maximum delay cap
    if delay > float64(eb.config.MaxDelay) {
        delay = float64(eb.config.MaxDelay)
    }
    
    // Add jitter to prevent thundering herd
    jitter := delay * eb.config.JitterFactor * (rand.Float64() - 0.5)
    delay += jitter
    
    return time.Duration(delay)
}
```

### 2. Retryable Error Classification

Define which errors should trigger retries:

```go
type RetryableError interface {
    IsRetryable() bool
    RetryAfter() time.Duration
}

func IsRetryableError(err error) bool {
    if retryable, ok := err.(RetryableError); ok {
        return retryable.IsRetryable()
    }
    
    // Network-related errors that are typically transient
    retryableErrors := []string{
        "connection refused",
        "connection reset",
        "network unreachable",
        "timeout",
        "temporary failure",
        "no route to host",
        "host unreachable",
    }
    
    errStr := strings.ToLower(err.Error())
    for _, retryableErr := range retryableErrors {
        if strings.Contains(errStr, retryableErr) {
            return true
        }
    }
    
    return false
}
```

### 3. Connection Retry Wrapper

Implement a retry wrapper for SSH connections:

```go
type RetryableConnection struct {
    config     RetryConfig
    strategy   RetryStrategy
    logger     spookylogging.Logger
}

func (rc *RetryableConnection) ConnectWithRetry(host string, config *spookytypes.SSHConfig) (*spookyssh.Connection, error) {
    var lastErr error
    
    for attempt := 1; attempt <= rc.config.MaxAttempts; attempt++ {
        // Attempt connection
        conn, err := spookyssh.NewConnection(host, config)
        if err == nil {
            rc.logger.Info("Connection established", "host", host, "attempt", attempt)
            return conn, nil
        }
        
        lastErr = err
        
        // Check if error is retryable
        if !rc.strategy.ShouldRetry(err, attempt) {
            rc.logger.Error("Non-retryable error encountered", "host", host, "error", err)
            return nil, err
        }
        
        // Calculate delay for next attempt
        delay := rc.strategy.NextDelay(attempt)
        
        rc.logger.Warn("Connection attempt failed, retrying",
            "host", host,
            "attempt", attempt,
            "max_attempts", rc.config.MaxAttempts,
            "delay", delay,
            "error", err)
        
        // Wait before next attempt
        time.Sleep(delay)
    }
    
    return nil, fmt.Errorf("failed to connect after %d attempts: %w", rc.config.MaxAttempts, lastErr)
}
```

### 4. Configuration Integration

Integrate retry configuration into SSH settings:

```go
type SSHConfig struct {
    Hostname string `hcl:"hostname"`
    Port     int    `hcl:"port"`
    User     string `hcl:"user"`
    
    Authentication *SSHAuthentication `hcl:"authentication"`
    Retry          *RetryConfig       `hcl:"retry,optional"`
    
    // ... existing fields
}

func (sc *SSHConfig) GetRetryConfig() RetryConfig {
    if sc.Retry != nil {
        return *sc.Retry
    }
    
    // Default retry configuration
    return RetryConfig{
        MaxAttempts:   3,
        InitialDelay:  time.Second,
        MaxDelay:      30 * time.Second,
        BackoffFactor: 2.0,
        JitterFactor:  0.1,
        Timeout:       30 * time.Second,
    }
}
```

### 5. Connection Pool Integration

Update connection pool to use retry logic:

```go
type ConnectionPool struct {
    connections map[string]*spookyssh.Connection
    mutex       sync.RWMutex
    retryConfig RetryConfig
    strategy    RetryStrategy
    logger      spookylogging.Logger
}

func (cp *ConnectionPool) GetConnectionWithRetry(host string, config *spookytypes.SSHConfig) (*spookyssh.Connection, error) {
    // Check existing connection first
    cp.mutex.RLock()
    if conn, exists := cp.connections[host]; exists && conn.IsHealthy() {
        cp.mutex.RUnlock()
        return conn, nil
    }
    cp.mutex.RUnlock()
    
    // Create new connection with retry logic
    retryable := &RetryableConnection{
        config:   cp.retryConfig,
        strategy: cp.strategy,
        logger:   cp.logger,
    }
    
    conn, err := retryable.ConnectWithRetry(host, config)
    if err != nil {
        return nil, err
    }
    
    // Store connection in pool
    cp.mutex.Lock()
    cp.connections[host] = conn
    cp.mutex.Unlock()
    
    return conn, nil
}
```

## Implementation Plan

### Phase 1: Core Retry Logic
1. Implement exponential backoff algorithm
2. Create retryable error classification
3. Add retry configuration structures
4. Implement basic retry wrapper

### Phase 2: SSH Integration
1. Integrate retry logic into SSH connections
2. Update connection pool with retry support
3. Add retry configuration to SSH settings
4. Implement connection health monitoring

### Phase 3: Advanced Features
1. Add circuit breaker pattern
2. Implement adaptive timeouts
3. Add retry metrics and monitoring
4. Create retry policy management

## Benefits

- **Improved Reliability**: Better handling of transient network issues
- **Reduced Failures**: Automatic recovery from temporary problems
- **Better User Experience**: Fewer manual retry attempts required
- **Resource Efficiency**: Prevents overwhelming target systems
- **Operational Insights**: Better visibility into connection patterns

## Risks and Mitigation

### Risks
- Increased connection latency for retry attempts
- Potential for resource exhaustion with aggressive retries
- Complexity in retry policy configuration

### Mitigation
- Configurable retry limits and timeouts
- Circuit breaker pattern to prevent cascading failures
- Clear documentation of retry policies
- Monitoring and alerting for retry patterns

## Success Metrics

- Reduced connection failure rates
- Improved system availability
- Faster recovery from network issues
- Better user experience with fewer manual interventions

## Related Documentation

- [SSH Implementation](mdc:ssh-implementation)
- [Error Handling Standards](mdc:error-handling-standards)
- [Configuration Management](mdc:configuration-management)
