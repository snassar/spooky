# Security Audit Logging

## Overview

The spooky codebase currently lacks comprehensive security audit logging for SSH connections, authentication events, and sensitive operations. This document outlines recommendations for implementing security audit logging to improve security monitoring, compliance, and incident response capabilities.

## Current State

- Limited logging of security-relevant events
- No structured audit log format
- Missing authentication and authorization logging
- No security event correlation
- Insufficient audit trail for compliance

## Recommendations

### 1. Security Audit Interface

Define a comprehensive security audit logging interface:

```go
type SecurityAuditor interface {
    LogAuthenticationEvent(event AuthenticationEvent) error
    LogAuthorizationEvent(event AuthorizationEvent) error
    LogConnectionEvent(event ConnectionEvent) error
    LogSensitiveOperation(event SensitiveOperationEvent) error
    LogSecurityViolation(event SecurityViolationEvent) error
    GetAuditLogs(filters AuditLogFilters) ([]AuditLogEntry, error)
    ExportAuditLogs(format string, filters AuditLogFilters) ([]byte, error)
}

type AuditLogEntry struct {
    ID          string                 `json:"id"`
    Timestamp   time.Time              `json:"timestamp"`
    EventType   AuditEventType         `json:"event_type"`
    Severity    AuditSeverity          `json:"severity"`
    User        string                 `json:"user"`
    SourceIP    string                 `json:"source_ip"`
    TargetHost  string                 `json:"target_host,omitempty"`
    Operation   string                 `json:"operation"`
    Status      AuditStatus            `json:"status"`
    Details     map[string]interface{} `json:"details"`
    SessionID   string                 `json:"session_id,omitempty"`
    CorrelationID string               `json:"correlation_id,omitempty"`
}

type AuditEventType string

const (
    AuditEventTypeAuthentication AuditEventType = "authentication"
    AuditEventTypeAuthorization  AuditEventType = "authorization"
    AuditEventTypeConnection     AuditEventType = "connection"
    AuditEventTypeOperation      AuditEventType = "operation"
    AuditEventTypeViolation      AuditEventType = "violation"
)

type AuditSeverity string

const (
    AuditSeverityLow      AuditSeverity = "low"
    AuditSeverityMedium   AuditSeverity = "medium"
    AuditSeverityHigh     AuditSeverity = "high"
    AuditSeverityCritical AuditSeverity = "critical"
)

type AuditStatus string

const (
    AuditStatusSuccess AuditStatus = "success"
    AuditStatusFailure AuditStatus = "failure"
    AuditStatusDenied  AuditStatus = "denied"
)
```

### 2. Authentication Event Logging

Implement comprehensive authentication event logging:

```go
type AuthenticationEvent struct {
    User        string            `json:"user"`
    Method      AuthMethod        `json:"method"`
    SourceIP    string            `json:"source_ip"`
    TargetHost  string            `json:"target_host"`
    Status      AuthStatus        `json:"status"`
    FailureReason string          `json:"failure_reason,omitempty"`
    SessionID   string            `json:"session_id"`
    Timestamp   time.Time         `json:"timestamp"`
    Metadata    map[string]string `json:"metadata"`
}

type AuthMethod string

const (
    AuthMethodSSHKey     AuthMethod = "ssh_key"
    AuthMethodPassword   AuthMethod = "password"
    AuthMethodCertificate AuthMethod = "certificate"
    AuthMethodToken      AuthMethod = "token"
)

type AuthStatus string

const (
    AuthStatusSuccess AuthStatus = "success"
    AuthStatusFailure AuthStatus = "failure"
    AuthStatusLocked  AuthStatus = "locked"
)

type SecurityAuditorImpl struct {
    logger     spookylogging.Logger
    storage    AuditLogStorage
    config     AuditConfig
    mutex      sync.RWMutex
}

func (sa *SecurityAuditorImpl) LogAuthenticationEvent(event AuthenticationEvent) error {
    auditEntry := AuditLogEntry{
        ID:         generateAuditID(),
        Timestamp:  event.Timestamp,
        EventType:  AuditEventTypeAuthentication,
        Severity:   sa.determineAuthSeverity(event),
        User:       event.User,
        SourceIP:   event.SourceIP,
        TargetHost: event.TargetHost,
        Operation:  fmt.Sprintf("authentication_%s", event.Method),
        Status:     sa.mapAuthStatus(event.Status),
        Details: map[string]interface{}{
            "method":         event.Method,
            "session_id":     event.SessionID,
            "failure_reason": event.FailureReason,
            "metadata":       event.Metadata,
        },
        SessionID:     event.SessionID,
        CorrelationID: generateCorrelationID(),
    }
    
    return sa.storage.StoreAuditLog(auditEntry)
}

func (sa *SecurityAuditorImpl) determineAuthSeverity(event AuthenticationEvent) AuditSeverity {
    switch event.Status {
    case AuthStatusSuccess:
        return AuditSeverityLow
    case AuthStatusFailure:
        return AuditSeverityMedium
    case AuthStatusLocked:
        return AuditSeverityHigh
    default:
        return AuditSeverityMedium
    }
}
```

### 3. Connection Event Logging

Implement connection event logging for SSH operations:

```go
type ConnectionEvent struct {
    SessionID   string            `json:"session_id"`
    User        string            `json:"user"`
    SourceIP    string            `json:"source_ip"`
    TargetHost  string            `json:"target_host"`
    TargetPort  int               `json:"target_port"`
    EventType   ConnectionEventType `json:"event_type"`
    Status      ConnectionStatus  `json:"status"`
    Duration    time.Duration     `json:"duration,omitempty"`
    BytesSent   int64             `json:"bytes_sent,omitempty"`
    BytesReceived int64            `json:"bytes_received,omitempty"`
    Commands    []string          `json:"commands,omitempty"`
    Timestamp   time.Time         `json:"timestamp"`
}

type ConnectionEventType string

const (
    ConnectionEventTypeConnect    ConnectionEventType = "connect"
    ConnectionEventTypeDisconnect ConnectionEventType = "disconnect"
    ConnectionEventTypeCommand    ConnectionEventType = "command"
    ConnectionEventTypeTransfer   ConnectionEventType = "transfer"
)

type ConnectionStatus string

const (
    ConnectionStatusEstablished ConnectionStatus = "established"
    ConnectionStatusFailed      ConnectionStatus = "failed"
    ConnectionStatusTerminated  ConnectionStatus = "terminated"
    ConnectionStatusTimeout     ConnectionStatus = "timeout"
)

func (sa *SecurityAuditorImpl) LogConnectionEvent(event ConnectionEvent) error {
    auditEntry := AuditLogEntry{
        ID:         generateAuditID(),
        Timestamp:  event.Timestamp,
        EventType:  AuditEventTypeConnection,
        Severity:   sa.determineConnectionSeverity(event),
        User:       event.User,
        SourceIP:   event.SourceIP,
        TargetHost: event.TargetHost,
        Operation:  string(event.EventType),
        Status:     sa.mapConnectionStatus(event.Status),
        Details: map[string]interface{}{
            "session_id":     event.SessionID,
            "target_port":    event.TargetPort,
            "duration":       event.Duration.String(),
            "bytes_sent":     event.BytesSent,
            "bytes_received": event.BytesReceived,
            "commands":       event.Commands,
        },
        SessionID:     event.SessionID,
        CorrelationID: generateCorrelationID(),
    }
    
    return sa.storage.StoreAuditLog(auditEntry)
}
```

### 4. Sensitive Operation Logging

Implement logging for sensitive operations:

```go
type SensitiveOperationEvent struct {
    User        string                 `json:"user"`
    SessionID   string                 `json:"session_id"`
    Operation   string                 `json:"operation"`
    Target      string                 `json:"target"`
    Parameters  map[string]interface{} `json:"parameters"`
    Status      OperationStatus        `json:"status"`
    Result      interface{}            `json:"result,omitempty"`
    Error       string                 `json:"error,omitempty"`
    Timestamp   time.Time              `json:"timestamp"`
}

type OperationStatus string

const (
    OperationStatusSuccess OperationStatus = "success"
    OperationStatusFailure OperationStatus = "failure"
    OperationStatusDenied  OperationStatus = "denied"
)

func (sa *SecurityAuditorImpl) LogSensitiveOperation(event SensitiveOperationEvent) error {
    auditEntry := AuditLogEntry{
        ID:         generateAuditID(),
        Timestamp:  event.Timestamp,
        EventType:  AuditEventTypeOperation,
        Severity:   sa.determineOperationSeverity(event),
        User:       event.User,
        Operation:  event.Operation,
        Status:     sa.mapOperationStatus(event.Status),
        Details: map[string]interface{}{
            "session_id": event.SessionID,
            "target":     event.Target,
            "parameters": sa.sanitizeParameters(event.Parameters),
            "result":     event.Result,
            "error":      event.Error,
        },
        SessionID:     event.SessionID,
        CorrelationID: generateCorrelationID(),
    }
    
    return sa.storage.StoreAuditLog(auditEntry)
}

func (sa *SecurityAuditorImpl) sanitizeParameters(params map[string]interface{}) map[string]interface{} {
    sanitized := make(map[string]interface{})
    
    for key, value := range params {
        // Sanitize sensitive parameters
        if sa.isSensitiveParameter(key) {
            sanitized[key] = "[REDACTED]"
        } else {
            sanitized[key] = value
        }
    }
    
    return sanitized
}

func (sa *SecurityAuditorImpl) isSensitiveParameter(key string) bool {
    sensitiveKeys := []string{
        "password", "passphrase", "secret", "key", "token",
        "credential", "auth", "private", "sensitive",
    }
    
    keyLower := strings.ToLower(key)
    for _, sensitive := range sensitiveKeys {
        if strings.Contains(keyLower, sensitive) {
            return true
        }
    }
    
    return false
}
```

### 5. Security Violation Logging

Implement security violation logging:

```go
type SecurityViolationEvent struct {
    User        string            `json:"user"`
    SourceIP    string            `json:"source_ip"`
    ViolationType ViolationType   `json:"violation_type"`
    Description string            `json:"description"`
    Severity    AuditSeverity     `json:"severity"`
    Action      ViolationAction   `json:"action"`
    Timestamp   time.Time         `json:"timestamp"`
    Context     map[string]string `json:"context"`
}

type ViolationType string

const (
    ViolationTypeBruteForce    ViolationType = "brute_force"
    ViolationTypeUnauthorized  ViolationType = "unauthorized_access"
    ViolationTypeRateLimit     ViolationType = "rate_limit_exceeded"
    ViolationTypeMalicious     ViolationType = "malicious_activity"
    ViolationTypePolicy        ViolationType = "policy_violation"
)

type ViolationAction string

const (
    ViolationActionLogged     ViolationAction = "logged"
    ViolationActionBlocked    ViolationAction = "blocked"
    ViolationActionAlerted    ViolationAction = "alerted"
    ViolationActionTerminated ViolationAction = "terminated"
)

func (sa *SecurityAuditorImpl) LogSecurityViolation(event SecurityViolationEvent) error {
    auditEntry := AuditLogEntry{
        ID:         generateAuditID(),
        Timestamp:  event.Timestamp,
        EventType:  AuditEventTypeViolation,
        Severity:   event.Severity,
        User:       event.User,
        SourceIP:   event.SourceIP,
        Operation:  string(event.ViolationType),
        Status:     AuditStatusDenied,
        Details: map[string]interface{}{
            "description": event.Description,
            "action":      event.Action,
            "context":     event.Context,
        },
        CorrelationID: generateCorrelationID(),
    }
    
    // Log to security-specific storage
    if err := sa.storage.StoreSecurityLog(auditEntry); err != nil {
        sa.logger.Error("Failed to store security log", "error", err)
    }
    
    return sa.storage.StoreAuditLog(auditEntry)
}
```

### 6. Audit Log Storage and Retrieval

Implement secure audit log storage:

```go
type AuditLogStorage interface {
    StoreAuditLog(entry AuditLogEntry) error
    StoreSecurityLog(entry AuditLogEntry) error
    GetAuditLogs(filters AuditLogFilters) ([]AuditLogEntry, error)
    ExportAuditLogs(format string, filters AuditLogFilters) ([]byte, error)
    RotateAuditLogs() error
    ArchiveAuditLogs(before time.Time) error
}

type AuditLogFilters struct {
    StartTime     *time.Time       `json:"start_time,omitempty"`
    EndTime       *time.Time       `json:"end_time,omitempty"`
    EventTypes    []AuditEventType `json:"event_types,omitempty"`
    Severities    []AuditSeverity  `json:"severities,omitempty"`
    Users         []string         `json:"users,omitempty"`
    SourceIPs     []string         `json:"source_ips,omitempty"`
    TargetHosts   []string         `json:"target_hosts,omitempty"`
    SessionIDs    []string         `json:"session_ids,omitempty"`
    Limit         int              `json:"limit,omitempty"`
    Offset        int              `json:"offset,omitempty"`
}

type SecureAuditStorage struct {
    db           *badger.DB
    encryption   AuditLogEncryption
    logger       spookylogging.Logger
    config       AuditStorageConfig
}

func (sas *SecureAuditStorage) StoreAuditLog(entry AuditLogEntry) error {
    // Encrypt sensitive audit log data
    encryptedData, err := sas.encryption.EncryptAuditLog(entry)
    if err != nil {
        return fmt.Errorf("failed to encrypt audit log: %w", err)
    }
    
    // Store with tamper-evident structure
    key := fmt.Sprintf("audit:%s:%s", entry.Timestamp.Format("2006-01-02"), entry.ID)
    
    return sas.db.Update(func(txn *badger.Txn) error {
        return txn.Set([]byte(key), encryptedData)
    })
}
```

## Implementation Plan

### Phase 1: Core Audit Logging
1. Implement security audit interface
2. Create authentication event logging
3. Add connection event logging
4. Implement basic audit storage

### Phase 2: Advanced Logging
1. Add sensitive operation logging
2. Implement security violation logging
3. Create audit log filtering and search
4. Add audit log encryption

### Phase 3: Compliance and Monitoring
1. Implement audit log retention policies
2. Add audit log export capabilities
3. Create security monitoring dashboards
4. Implement audit log analysis tools

## Benefits

- **Security Monitoring**: Comprehensive visibility into security events
- **Compliance Support**: Detailed audit trails for regulatory requirements
- **Incident Response**: Better detection and investigation of security incidents
- **Risk Assessment**: Data-driven security risk analysis
- **Forensic Analysis**: Detailed logs for post-incident investigation

## Risks and Mitigation

### Risks
- Performance impact from extensive logging
- Storage requirements for audit logs
- Privacy concerns with detailed logging
- Potential for log tampering

### Mitigation
- Configurable logging levels and filters
- Secure, encrypted audit log storage
- Data retention and archival policies
- Tamper-evident log structures
- Regular audit log integrity checks

## Success Metrics

- Complete audit trail for all security events
- Reduced time to detect security incidents
- Improved compliance with security requirements
- Better incident response capabilities
- Enhanced security posture monitoring

## Related Documentation

- [SSH Implementation](mdc:ssh-implementation)
- [Error Handling Standards](mdc:error-handling-standards)
- [Configuration Management](mdc:configuration-management)
