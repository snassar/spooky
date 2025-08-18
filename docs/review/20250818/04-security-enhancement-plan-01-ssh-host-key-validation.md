---
description: SSH host key validation enhancement implementation plan
globs: ["internal/ssh/**/*.go", "internal/types/ssh/**/*.go"]
alwaysApply: false
---

# Security Enhancement Plan: SSH Host Key Validation

**Generated:** 2025-08-18  
**Issue Type:** Security - SSH Host Key Validation  
**Priority:** High  
**Status:** Implementation Planning

## Overview

Enhance SSH host key validation to prevent MITM attacks and improve security monitoring. This plan focuses on implementing proper host key validation, security monitoring, and audit logging for SSH connections.

## Current State Analysis

### Existing Infrastructure
- Basic SSH client implementation in `internal/ssh/`
- SSH connection types in `internal/types/ssh/`
- Current host key validation is minimal

### Security Gaps
- No proper host key validation
- Missing security monitoring
- No audit logging for SSH connections
- No host key fingerprint verification

## Implementation Plan

### Phase 1: Enhanced Host Key Validation (Priority: High)

#### 1.1 Implement Host Key Callback
```go
// internal/ssh/host_key.go
type HostKeyValidator struct {
    knownHosts map[string]ssh.PublicKey
    logger     spookylogging.Logger
    auditLog   spookyinterfaces.AuditLogger
}

func (v *HostKeyValidator) ValidateHostKey(hostname string, remote net.Addr, key ssh.PublicKey) error {
    // Implementation with proper validation
}
```

#### 1.2 Add Host Key Storage
```go
// internal/ssh/known_hosts.go
type KnownHostsManager struct {
    filePath string
    mutex    sync.RWMutex
}

func (m *KnownHostsManager) AddHostKey(hostname string, key ssh.PublicKey) error {
    // Implementation for adding host keys
}
```

### Phase 2: Security Monitoring Integration (Priority: Medium)

#### 2.1 Security Event Logging
```go
// internal/ssh/security_monitor.go
type SSHSecurityMonitor struct {
    logger   spookylogging.Logger
    metrics  spookyinterfaces.SecurityMetrics
}

func (m *SSHSecurityMonitor) LogConnectionAttempt(hostname string, success bool, details map[string]interface{}) {
    // Implementation for security event logging
}
```

#### 2.2 Host Key Fingerprint Verification
```go
// internal/ssh/fingerprint.go
func ValidateHostKeyFingerprint(hostname, expectedFingerprint string, actualKey ssh.PublicKey) error {
    // Implementation for fingerprint validation
}
```

## Success Metrics

- **Security**: 100% of SSH connections validate host keys
- **Monitoring**: All SSH connection attempts logged with security context
- **Audit**: Complete audit trail for SSH security events
- **Performance**: Host key validation adds <10ms to connection time

## Implementation Steps

1. **Week 1**: Implement host key validation callback
2. **Week 2**: Add known hosts management
3. **Week 3**: Integrate security monitoring
4. **Week 4**: Add audit logging and testing

## Dependencies

- SSH client implementation
- Logging infrastructure
- Security metrics system
- Audit logging system

## Risk Assessment

- **Low Risk**: Leverages existing SSH infrastructure
- **Medium Impact**: Improves security posture significantly
- **High Value**: Prevents MITM attacks and provides audit trail
