# Rate Limiting for SSH Connections

## Overview

The spooky codebase currently lacks rate limiting mechanisms for SSH connections, which could lead to resource exhaustion, potential security vulnerabilities, and degraded performance under high load. This document outlines recommendations for implementing rate limiting to improve security, stability, and resource management.

## Current State

- No rate limiting on SSH connection attempts
- No protection against brute force attacks
- No connection throttling under high load
- No per-user or per-host connection limits
- No adaptive rate limiting based on system resources

## Recommendations

### 1. Rate Limiting Interface

Define a comprehensive rate limiting interface:

```go
type RateLimiter interface {
    AllowConnection(user string, host string, sourceIP string) (bool, error)
    AllowAuthentication(user string, host string, sourceIP string) (bool, error)
    AllowCommand(user string, host string, sourceIP string, command string) (bool, error)
    GetRateLimitStatus(user string, host string, sourceIP string) (*RateLimitStatus, error)
    ResetRateLimit(user string, host string, sourceIP string) error
    UpdateRateLimitConfig(config RateLimitConfig) error
}

type RateLimitStatus struct {
    User           string    `json:"user"`
    Host           string    `json:"host"`
    SourceIP       string    `json:"source_ip"`
    ConnectionCount int      `json:"connection_count"`
    AuthAttempts   int      `json:"auth_attempts"`
    CommandCount   int      `json:"command_count"`
    LastAttempt    time.Time `json:"last_attempt"`
    IsBlocked      bool      `json:"is_blocked"`
    BlockUntil     time.Time `json:"block_until,omitempty"`
    RetryAfter     time.Duration `json:"retry_after,omitempty"`
}

type RateLimitConfig struct {
    MaxConnectionsPerMinute int           `hcl:"max_connections_per_minute"`
    MaxAuthAttemptsPerMinute int          `hcl:"max_auth_attempts_per_minute"`
    MaxCommandsPerMinute    int           `hcl:"max_commands_per_minute"`
    BlockDuration           time.Duration `hcl:"block_duration"`
    BurstLimit              int           `hcl:"burst_limit"`
    AdaptiveLimiting        bool          `hcl:"adaptive_limiting"`
    WhitelistIPs            []string      `hcl:"whitelist_ips,optional"`
    BlacklistIPs            []string      `hcl:"blacklist_ips,optional"`
}
```

### 2. Token Bucket Rate Limiter

Implement token bucket algorithm for rate limiting:

```go
type TokenBucketRateLimiter struct {
    config     RateLimitConfig
    buckets    map[string]*TokenBucket
    mutex      sync.RWMutex
    logger     spookylogging.Logger
    metrics    RateLimitMetrics
}

type TokenBucket struct {
    tokens     float64
    capacity   float64
    rate       float64
    lastRefill time.Time
    mutex      sync.Mutex
}

func NewTokenBucketRateLimiter(config RateLimitConfig) *TokenBucketRateLimiter {
    return &TokenBucketRateLimiter{
        config:  config,
        buckets: make(map[string]*TokenBucket),
        logger:  spookylogging.GetLogger(),
        metrics: NewRateLimitMetrics(),
    }
}

func (tbrl *TokenBucketRateLimiter) AllowConnection(user, host, sourceIP string) (bool, error) {
    key := fmt.Sprintf("connection:%s:%s:%s", user, host, sourceIP)
    
    // Check whitelist/blacklist
    if tbrl.isWhitelisted(sourceIP) {
        return true, nil
    }
    
    if tbrl.isBlacklisted(sourceIP) {
        return false, fmt.Errorf("source IP %s is blacklisted", sourceIP)
    }
    
    bucket := tbrl.getOrCreateBucket(key, tbrl.config.MaxConnectionsPerMinute)
    
    if bucket.tryConsume(1) {
        tbrl.metrics.RecordConnection(user, host, sourceIP, true)
        return true, nil
    }
    
    tbrl.metrics.RecordConnection(user, host, sourceIP, false)
    tbrl.logger.Warn("Rate limit exceeded for connection",
        "user", user, "host", host, "source_ip", sourceIP)
    
    return false, fmt.Errorf("rate limit exceeded for connections")
}

func (tbrl *TokenBucketRateLimiter) AllowAuthentication(user, host, sourceIP string) (bool, error) {
    key := fmt.Sprintf("auth:%s:%s:%s", user, host, sourceIP)
    
    bucket := tbrl.getOrCreateBucket(key, tbrl.config.MaxAuthAttemptsPerMinute)
    
    if bucket.tryConsume(1) {
        tbrl.metrics.RecordAuthentication(user, host, sourceIP, true)
        return true, nil
    }
    
    tbrl.metrics.RecordAuthentication(user, host, sourceIP, false)
    tbrl.logger.Warn("Rate limit exceeded for authentication",
        "user", user, "host", host, "source_ip", sourceIP)
    
    return false, fmt.Errorf("rate limit exceeded for authentication attempts")
}

func (tbrl *TokenBucketRateLimiter) getOrCreateBucket(key string, rate int) *TokenBucket {
    tbrl.mutex.Lock()
    defer tbrl.mutex.Unlock()
    
    if bucket, exists := tbrl.buckets[key]; exists {
        return bucket
    }
    
    bucket := &TokenBucket{
        tokens:     float64(rate),
        capacity:   float64(rate),
        rate:       float64(rate) / 60.0, // tokens per second
        lastRefill: time.Now(),
    }
    
    tbrl.buckets[key] = bucket
    return bucket
}

func (tb *TokenBucket) tryConsume(tokens float64) bool {
    tb.mutex.Lock()
    defer tb.mutex.Unlock()
    
    tb.refill()
    
    if tb.tokens >= tokens {
        tb.tokens -= tokens
        return true
    }
    
    return false
}

func (tb *TokenBucket) refill() {
    now := time.Now()
    elapsed := now.Sub(tb.lastRefill).Seconds()
    
    tokensToAdd := elapsed * tb.rate
    tb.tokens = math.Min(tb.capacity, tb.tokens+tokensToAdd)
    tb.lastRefill = now
}
```

### 3. Adaptive Rate Limiting

Implement adaptive rate limiting based on system resources:

```go
type AdaptiveRateLimiter struct {
    baseLimiter *TokenBucketRateLimiter
    config      AdaptiveRateLimitConfig
    metrics     SystemMetrics
    logger      spookylogging.Logger
}

type AdaptiveRateLimitConfig struct {
    BaseConfig           RateLimitConfig `hcl:"base_config"`
    CPUThreshold         float64         `hcl:"cpu_threshold"`
    MemoryThreshold      float64         `hcl:"memory_threshold"`
    ConnectionThreshold  int             `hcl:"connection_threshold"`
    ReductionFactor      float64         `hcl:"reduction_factor"`
    RecoveryFactor       float64         `hcl:"recovery_factor"`
    CheckInterval        time.Duration   `hcl:"check_interval"`
}

func (arl *AdaptiveRateLimiter) StartAdaptiveMonitoring() {
    ticker := time.NewTicker(arl.config.CheckInterval)
    defer ticker.Stop()
    
    for range ticker.C {
        arl.adjustRateLimits()
    }
}

func (arl *AdaptiveRateLimiter) adjustRateLimits() {
    metrics := arl.metrics.GetSystemMetrics()
    
    // Check if system is under stress
    if arl.isSystemUnderStress(metrics) {
        arl.reduceRateLimits()
        arl.logger.Info("Reducing rate limits due to system stress",
            "cpu_usage", metrics.CPUUsage,
            "memory_usage", metrics.MemoryUsage,
            "active_connections", metrics.ActiveConnections)
    } else {
        arl.recoverRateLimits()
    }
}

func (arl *AdaptiveRateLimiter) isSystemUnderStress(metrics SystemMetrics) bool {
    return metrics.CPUUsage > arl.config.CPUThreshold ||
           metrics.MemoryUsage > arl.config.MemoryThreshold ||
           metrics.ActiveConnections > arl.config.ConnectionThreshold
}

func (arl *AdaptiveRateLimiter) reduceRateLimits() {
    config := arl.baseLimiter.config
    
    config.MaxConnectionsPerMinute = int(float64(config.MaxConnectionsPerMinute) * arl.config.ReductionFactor)
    config.MaxAuthAttemptsPerMinute = int(float64(config.MaxAuthAttemptsPerMinute) * arl.config.ReductionFactor)
    config.MaxCommandsPerMinute = int(float64(config.MaxCommandsPerMinute) * arl.config.ReductionFactor)
    
    arl.baseLimiter.UpdateRateLimitConfig(config)
}

func (arl *AdaptiveRateLimiter) recoverRateLimits() {
    config := arl.baseLimiter.config
    
    config.MaxConnectionsPerMinute = int(float64(config.MaxConnectionsPerMinute) * arl.config.RecoveryFactor)
    config.MaxAuthAttemptsPerMinute = int(float64(config.MaxAuthAttemptsPerMinute) * arl.config.RecoveryFactor)
    config.MaxCommandsPerMinute = int(float64(config.MaxCommandsPerMinute) * arl.config.RecoveryFactor)
    
    arl.baseLimiter.UpdateRateLimitConfig(config)
}
```

### 4. IP-Based Rate Limiting

Implement IP-based rate limiting with blacklist/whitelist support:

```go
type IPRateLimiter struct {
    rateLimiter *TokenBucketRateLimiter
    whitelist   map[string]bool
    blacklist   map[string]bool
    mutex       sync.RWMutex
    logger      spookylogging.Logger
}

func (iprl *IPRateLimiter) isWhitelisted(sourceIP string) bool {
    iprl.mutex.RLock()
    defer iprl.mutex.RUnlock()
    
    return iprl.whitelist[sourceIP]
}

func (iprl *IPRateLimiter) isBlacklisted(sourceIP string) bool {
    iprl.mutex.RLock()
    defer iprl.mutex.RUnlock()
    
    return iprl.blacklist[sourceIP]
}

func (iprl *IPRateLimiter) AddToWhitelist(sourceIP string) error {
    if !iprl.isValidIP(sourceIP) {
        return fmt.Errorf("invalid IP address: %s", sourceIP)
    }
    
    iprl.mutex.Lock()
    defer iprl.mutex.Unlock()
    
    iprl.whitelist[sourceIP] = true
    delete(iprl.blacklist, sourceIP) // Remove from blacklist if present
    
    iprl.logger.Info("Added IP to whitelist", "source_ip", sourceIP)
    return nil
}

func (iprl *IPRateLimiter) AddToBlacklist(sourceIP string, duration time.Duration) error {
    if !iprl.isValidIP(sourceIP) {
        return fmt.Errorf("invalid IP address: %s", sourceIP)
    }
    
    iprl.mutex.Lock()
    defer iprl.mutex.Unlock()
    
    iprl.blacklist[sourceIP] = true
    delete(iprl.whitelist, sourceIP) // Remove from whitelist if present
    
    // Schedule removal from blacklist
    go func() {
        time.Sleep(duration)
        iprl.mutex.Lock()
        delete(iprl.blacklist, sourceIP)
        iprl.mutex.Unlock()
        iprl.logger.Info("Removed IP from blacklist", "source_ip", sourceIP)
    }()
    
    iprl.logger.Info("Added IP to blacklist", "source_ip", sourceIP, "duration", duration)
    return nil
}

func (iprl *IPRateLimiter) isValidIP(sourceIP string) bool {
    return net.ParseIP(sourceIP) != nil
}
```

### 5. Rate Limit Metrics and Monitoring

Implement comprehensive metrics collection:

```go
type RateLimitMetrics struct {
    connections    map[string]*ConnectionMetrics
    authentications map[string]*AuthMetrics
    commands       map[string]*CommandMetrics
    mutex          sync.RWMutex
}

type ConnectionMetrics struct {
    TotalAttempts    int64         `json:"total_attempts"`
    SuccessfulAttempts int64       `json:"successful_attempts"`
    BlockedAttempts  int64         `json:"blocked_attempts"`
    LastAttempt      time.Time     `json:"last_attempt"`
    AverageLatency   time.Duration `json:"average_latency"`
}

type AuthMetrics struct {
    TotalAttempts    int64         `json:"total_attempts"`
    SuccessfulAttempts int64       `json:"successful_attempts"`
    FailedAttempts   int64         `json:"failed_attempts"`
    BlockedAttempts  int64         `json:"blocked_attempts"`
    LastAttempt      time.Time     `json:"last_attempt"`
}

type CommandMetrics struct {
    TotalCommands    int64         `json:"total_commands"`
    SuccessfulCommands int64       `json:"successful_commands"`
    BlockedCommands  int64         `json:"blocked_commands"`
    LastCommand      time.Time     `json:"last_command"`
}

func (rlm *RateLimitMetrics) RecordConnection(user, host, sourceIP string, allowed bool) {
    key := fmt.Sprintf("%s:%s:%s", user, host, sourceIP)
    
    rlm.mutex.Lock()
    defer rlm.mutex.Unlock()
    
    if rlm.connections[key] == nil {
        rlm.connections[key] = &ConnectionMetrics{}
    }
    
    metrics := rlm.connections[key]
    metrics.TotalAttempts++
    metrics.LastAttempt = time.Now()
    
    if allowed {
        metrics.SuccessfulAttempts++
    } else {
        metrics.BlockedAttempts++
    }
}

func (rlm *RateLimitMetrics) GetMetrics() map[string]interface{} {
    rlm.mutex.RLock()
    defer rlm.mutex.RUnlock()
    
    return map[string]interface{}{
        "connections":     rlm.connections,
        "authentications": rlm.authentications,
        "commands":        rlm.commands,
    }
}
```

### 6. SSH Integration

Integrate rate limiting into SSH connection management:

```go
type RateLimitedSSHManager struct {
    sshManager  spookyinterfaces.SSHManager
    rateLimiter RateLimiter
    logger      spookylogging.Logger
}

func (rlssh *RateLimitedSSHManager) Connect(host string, config *spookytypes.SSHConfig) (*spookyssh.Connection, error) {
    user := config.User
    sourceIP := getSourceIP()
    
    // Check rate limits before attempting connection
    allowed, err := rlssh.rateLimiter.AllowConnection(user, host, sourceIP)
    if err != nil {
        return nil, fmt.Errorf("rate limit check failed: %w", err)
    }
    
    if !allowed {
        return nil, fmt.Errorf("connection rate limit exceeded for %s@%s from %s", user, host, sourceIP)
    }
    
    // Proceed with connection
    return rlssh.sshManager.Connect(host, config)
}

func (rlssh *RateLimitedSSHManager) Authenticate(conn *spookyssh.Connection, config *spookytypes.SSHConfig) error {
    user := config.User
    host := conn.RemoteAddr().String()
    sourceIP := getSourceIP()
    
    // Check authentication rate limits
    allowed, err := rlssh.rateLimiter.AllowAuthentication(user, host, sourceIP)
    if err != nil {
        return fmt.Errorf("authentication rate limit check failed: %w", err)
    }
    
    if !allowed {
        return fmt.Errorf("authentication rate limit exceeded for %s@%s from %s", user, host, sourceIP)
    }
    
    // Proceed with authentication
    return rlssh.sshManager.Authenticate(conn, config)
}

func (rlssh *RateLimitedSSHManager) ExecuteCommand(conn *spookyssh.Connection, command string) (*spookyssh.CommandResult, error) {
    user := conn.User()
    host := conn.RemoteAddr().String()
    sourceIP := getSourceIP()
    
    // Check command rate limits
    allowed, err := rlssh.rateLimiter.AllowCommand(user, host, sourceIP, command)
    if err != nil {
        return nil, fmt.Errorf("command rate limit check failed: %w", err)
    }
    
    if !allowed {
        return nil, fmt.Errorf("command rate limit exceeded for %s@%s from %s", user, host, sourceIP)
    }
    
    // Proceed with command execution
    return rlssh.sshManager.ExecuteCommand(conn, command)
}
```

## Implementation Plan

### Phase 1: Core Rate Limiting
1. Implement token bucket rate limiter
2. Create rate limiting interfaces
3. Add basic rate limit configuration
4. Implement IP-based filtering

### Phase 2: SSH Integration
1. Integrate rate limiting into SSH connections
2. Add authentication rate limiting
3. Implement command rate limiting
4. Add rate limit metrics collection

### Phase 3: Advanced Features
1. Implement adaptive rate limiting
2. Add blacklist/whitelist management
3. Create rate limit monitoring dashboards
4. Implement rate limit analytics

## Benefits

- **Security Enhancement**: Protection against brute force attacks
- **Resource Protection**: Prevention of resource exhaustion
- **Performance Stability**: Consistent performance under load
- **Operational Control**: Granular control over connection rates
- **Monitoring Capabilities**: Better visibility into connection patterns

## Risks and Mitigation

### Risks
- Potential for legitimate users to be blocked
- Complexity in rate limit configuration
- Performance overhead from rate limiting checks

### Mitigation
- Configurable whitelist for trusted sources
- Adaptive rate limiting based on system load
- Efficient rate limiting algorithms
- Clear documentation and monitoring

## Success Metrics

- Reduced brute force attack success rates
- Improved system stability under load
- Better resource utilization
- Reduced security incidents
- Improved user experience for legitimate users

## Related Documentation

- [SSH Implementation](mdc:ssh-implementation)
- [Connection Health Monitoring](mdc:connection-health-monitoring)
- [Security Audit Logging](mdc:security-audit-logging)
