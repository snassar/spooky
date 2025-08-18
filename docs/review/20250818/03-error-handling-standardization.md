# Error Handling Standardization Issues and Solutions

**Generated:** 2025-08-18  
**Issue Type:** Code Quality - Error Handling Inconsistency  
**Priority:** High  
**Status:** Needs Standardization

## Overview

The spooky codebase has significant inconsistencies in error handling patterns across different packages and components. This document provides specific examples of problematic error handling and standardized solutions to improve code reliability and maintainability.

## Current Error Handling Issues

### 1. Inconsistent Error Types

**File:** `internal/ssh/client.go`  
**Issue:** Mixed error types without proper context

```go
// PROBLEMATIC CODE - Inconsistent error handling
func (c *Client) Connect() error {
    if c.config == nil {
        return fmt.Errorf("config is nil") // Generic error
    }
    
    if c.config.Host == "" {
        return errors.New("host is required") // Different error type
    }
    
    client, err := ssh.Dial("tcp", c.config.Host, c.config.ClientConfig)
    if err != nil {
        return err // Raw error without context
    }
    
    c.client = client
    return nil
}

// PROBLEMS:
// - Mixed error types (fmt.Errorf vs errors.New)
// - No structured error information
// - Missing context for SSH errors
// - No error categorization
// - No error recovery information
// - No error codes for programmatic handling
```

### 2. Missing Error Context

**File:** `internal/facts/manager.go`  
**Issue:** Errors lack sufficient context for debugging

```go
// PROBLEMATIC CODE - Missing error context
func (m *Manager) CollectFacts(server string) (*spookytypes.FactCollection, error) {
    if server == "" {
        return nil, fmt.Errorf("server is required")
    }
    
    facts, err := m.storage.GetFacts(server)
    if err != nil {
        return nil, err // Missing context about what failed
    }
    
    if facts == nil {
        return nil, fmt.Errorf("no facts found") // No context about server
    }
    
    return facts, nil
}

// PROBLEMS:
// - No context about which operation failed
// - No server information in error messages
// - No error categorization
// - No error recovery suggestions
// - No error codes for programmatic handling
// - No structured error information
```

### 3. Inconsistent Error Wrapping

**File:** `internal/actions/manager.go`  
**Issue:** Inconsistent error wrapping patterns

```go
// PROBLEMATIC CODE - Inconsistent error wrapping
func (m *Manager) RunAction(action spookytypes.Action, machine spookytypes.Machine) error {
    // Validate action
    if err := m.validateAction(action); err != nil {
        return fmt.Errorf("action validation failed: %v", err) // Inconsistent wrapping
    }
    
    // Validate machine
    if err := m.validateMachine(machine); err != nil {
        return err // No wrapping at all
    }
    
    // Execute action
    if err := m.executeAction(action, machine); err != nil {
        return fmt.Errorf("failed to execute action: %w", err) // Different wrapping style
    }
    
    return nil
}

// PROBLEMS:
// - Inconsistent error wrapping patterns
// - Mixed use of %v and %w
// - Some errors not wrapped at all
// - No consistent error hierarchy
// - No error categorization
// - No error recovery information
```

### 4. Missing Error Recovery

**File:** `internal/config/loader.go`  
**Issue:** No error recovery mechanisms

```go
// PROBLEMATIC CODE - No error recovery
func LoadConfig(path string) (*spookytypes.Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err // No recovery attempt
    }
    
    var config spookytypes.Config
    if err := hcl.Unmarshal(data, &config); err != nil {
        return nil, err // No recovery attempt
    }
    
    if err := config.Validate(); err != nil {
        return nil, err // No recovery attempt
    }
    
    return &config, nil
}

// PROBLEMS:
// - No fallback configuration
// - No error recovery mechanisms
// - No default value handling
// - No error categorization
// - No error recovery suggestions
// - No structured error information
```

## Current Error Handling Infrastructure

### Existing Error Types

The codebase already has a foundation of structured error types that should be leveraged:

**Common Error Types (`internal/types/common/common.go`):**
```go
// ErrorDetails provides structured error information with context and stack traces
type ErrorDetails struct {
    Code        string                 `json:"code" hcl:"code"`
    Message     string                 `json:"message" hcl:"message"`
    Context     map[string]interface{} `json:"context,omitempty" hcl:"context,optional"`
    Stack       []string               `json:"stack,omitempty" hcl:"stack,optional"`
    Recoverable bool                   `json:"recoverable" hcl:"recoverable"`
    Timestamp   time.Time              `json:"timestamp" hcl:"timestamp"`
}
```

**Domain-Specific Error Types:**
- **SSH Errors** (`internal/types/ssh/errors.go`): ConnectionError, AuthenticationError, TimeoutError
- **Action Errors** (`internal/types/actions/errors.go`): ActionError, ActingError, PlanningError
- **Machine Errors** (`internal/types/machines/errors.go`): MachineError, MachineValidationError
- **Config Errors** (`internal/types/config/config.go`): Configuration-specific errors
- **Variable Errors** (`internal/types/variables/variable.go`): VariableError, VariableWarning

## Standardized Error Handling Solutions

### Solution 1: Enhanced Error Type System

**Enhanced Error Types Building on Existing Infrastructure:**

```go
// GOOD CODE - Enhanced error types leveraging existing infrastructure
package errors

import (
    "fmt"
    "time"

    spookytypescommon "spooky/internal/types/common"
)

// Error categories (extending existing patterns)
type ErrorCategory string

const (
    ErrorCategoryValidation   ErrorCategory = "validation"
    ErrorCategoryConnection   ErrorCategory = "connection"
    ErrorCategoryAuthentication ErrorCategory = "authentication"
    ErrorCategoryPermission   ErrorCategory = "permission"
    ErrorCategoryTimeout      ErrorCategory = "timeout"
    ErrorCategoryResource     ErrorCategory = "resource"
    ErrorCategoryConfiguration ErrorCategory = "configuration"
    ErrorCategoryInternal     ErrorCategory = "internal"
)

// Error severity levels
type ErrorSeverity string

const (
    ErrorSeverityLow      ErrorSeverity = "low"
    ErrorSeverityMedium   ErrorSeverity = "medium"
    ErrorSeverityHigh     ErrorSeverity = "high"
    ErrorSeverityCritical ErrorSeverity = "critical"
)

// Enhanced error structure building on ErrorDetails
type SpookyError struct {
    spookytypescommon.ErrorDetails
    
    Category    ErrorCategory `json:"category"`
    Severity    ErrorSeverity `json:"severity"`
    Operation   string        `json:"operation"`
    Component   string        `json:"component"`
    Suggestions []string      `json:"suggestions,omitempty"`
    Cause       error         `json:"-"`
}

func (e *SpookyError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("[%s] %s: %s (caused by: %v)", e.Category, e.Message, e.Operation, e.Cause)
    }
    return fmt.Sprintf("[%s] %s: %s", e.Category, e.Message, e.Operation)
}

func (e *SpookyError) Unwrap() error {
    return e.Cause
}

// Error constructors leveraging existing patterns
func NewValidationError(code, message, operation string, context map[string]interface{}) *SpookyError {
    return &SpookyError{
        ErrorDetails: spookytypescommon.ErrorDetails{
            Code:        code,
            Message:     message,
            Context:     context,
            Timestamp:   time.Now(),
            Recoverable: true,
        },
        Category:    ErrorCategoryValidation,
        Severity:    ErrorSeverityMedium,
        Operation:   operation,
        Component:   "validation",
        Suggestions: []string{"Check input parameters", "Verify configuration values"},
    }
}

func NewConnectionError(code, message, operation string, context map[string]interface{}) *SpookyError {
    return &SpookyError{
        ErrorDetails: spookytypescommon.ErrorDetails{
            Code:        code,
            Message:     message,
            Context:     context,
            Timestamp:   time.Now(),
            Recoverable: true,
        },
        Category:    ErrorCategoryConnection,
        Severity:    ErrorSeverityHigh,
        Operation:   operation,
        Component:   "connection",
        Suggestions: []string{"Check network connectivity", "Verify host availability", "Retry operation"},
    }
}

func NewAuthenticationError(code, message, operation string, context map[string]interface{}) *SpookyError {
    return &SpookyError{
        ErrorDetails: spookytypescommon.ErrorDetails{
            Code:        code,
            Message:     message,
            Context:     context,
            Timestamp:   time.Now(),
            Recoverable: true,
        },
        Category:    ErrorCategoryAuthentication,
        Severity:    ErrorSeverityHigh,
        Operation:   operation,
        Component:   "authentication",
        Suggestions: []string{"Verify credentials", "Check SSH key permissions", "Validate authentication method"},
    }
}

func NewTimeoutError(code, message, operation string, context map[string]interface{}) *SpookyError {
    return &SpookyError{
        ErrorDetails: spookytypescommon.ErrorDetails{
            Code:        code,
            Message:     message,
            Context:     context,
            Timestamp:   time.Now(),
            Recoverable: true,
        },
        Category:    ErrorCategoryTimeout,
        Severity:    ErrorSeverityMedium,
        Operation:   operation,
        Component:   "timeout",
        Suggestions: []string{"Increase timeout value", "Check network latency", "Retry operation"},
    }
}

func NewResourceError(code, message, operation string, context map[string]interface{}) *SpookyError {
    return &SpookyError{
        ErrorDetails: spookytypescommon.ErrorDetails{
            Code:        code,
            Message:     message,
            Context:     context,
            Timestamp:   time.Now(),
            Recoverable: true,
        },
        Category:    ErrorCategoryResource,
        Severity:    ErrorSeverityHigh,
        Operation:   operation,
        Component:   "resource",
        Suggestions: []string{"Check resource availability", "Increase resource limits", "Free up resources"},
    }
}

func NewConfigurationError(code, message, operation string, context map[string]interface{}) *SpookyError {
    return &SpookyError{
        ErrorDetails: spookytypescommon.ErrorDetails{
            Code:        code,
            Message:     message,
            Context:     context,
            Timestamp:   time.Now(),
            Recoverable: true,
        },
        Category:    ErrorCategoryConfiguration,
        Severity:    ErrorSeverityMedium,
        Operation:   operation,
        Component:   "configuration",
        Suggestions: []string{"Check configuration file", "Validate configuration schema", "Use default values"},
    }
}

func NewInternalError(code, message, operation string, context map[string]interface{}) *SpookyError {
    return &SpookyError{
        ErrorDetails: spookytypescommon.ErrorDetails{
            Code:        code,
            Message:     message,
            Context:     context,
            Timestamp:   time.Now(),
            Recoverable: false,
        },
        Category:    ErrorCategoryInternal,
        Severity:    ErrorSeverityCritical,
        Operation:   operation,
        Component:   "internal",
        Suggestions: []string{"Check system logs", "Report issue to support", "Restart application"},
    }
}

// Error wrapping with context
func WrapError(err error, category ErrorCategory, code, message, operation string, context map[string]interface{}) *SpookyError {
    spookyErr := &SpookyError{
        ErrorDetails: spookytypescommon.ErrorDetails{
            Code:        code,
            Message:     message,
            Context:     context,
            Timestamp:   time.Now(),
            Recoverable: true,
        },
        Category:    category,
        Severity:    ErrorSeverityMedium,
        Operation:   operation,
        Component:   "wrapper",
        Cause:       err,
    }
    
    // Inherit suggestions from wrapped error if it's a SpookyError
    if wrappedSpookyErr, ok := err.(*SpookyError); ok {
        spookyErr.Suggestions = append(spookyErr.Suggestions, wrappedSpookyErr.Suggestions...)
    }
    
    return spookyErr
}

// Error categorization helpers
func IsValidationError(err error) bool {
    return isErrorCategory(err, ErrorCategoryValidation)
}

func IsConnectionError(err error) bool {
    return isErrorCategory(err, ErrorCategoryConnection)
}

func IsAuthenticationError(err error) bool {
    return isErrorCategory(err, ErrorCategoryAuthentication)
}

func IsTimeoutError(err error) bool {
    return isErrorCategory(err, ErrorCategoryTimeout)
}

func IsResourceError(err error) bool {
    return isErrorCategory(err, ErrorCategoryResource)
}

func IsConfigurationError(err error) bool {
    return isErrorCategory(err, ErrorCategoryConfiguration)
}

func IsInternalError(err error) bool {
    return isErrorCategory(err, ErrorCategoryInternal)
}

func isErrorCategory(err error, category ErrorCategory) bool {
    if spookyErr, ok := err.(*SpookyError); ok {
        return spookyErr.Category == category
    }
    return false
}

// Error aggregation
type AggregateError struct {
    Errors []*SpookyError `json:"errors"`
    Message string        `json:"message"`
}

func (ae *AggregateError) Error() string {
    return fmt.Sprintf("%s (%d errors)", ae.Message, len(ae.Errors))
}

func NewAggregateError(message string, errors []*SpookyError) *AggregateError {
    return &AggregateError{
        Errors: errors,
        Message: message,
    }
}

func (ae *AggregateError) AddError(err *SpookyError) {
    ae.Errors = append(ae.Errors, err)
}

func (ae *AggregateError) HasErrors() bool {
    return len(ae.Errors) > 0
}

func (ae *AggregateError) GetErrorsByCategory(category ErrorCategory) []*SpookyError {
    var filtered []*SpookyError
    for _, err := range ae.Errors {
        if err.Category == category {
            filtered = append(filtered, err)
        }
    }
    return filtered
}

func (ae *AggregateError) GetErrorsBySeverity(severity ErrorSeverity) []*SpookyError {
    var filtered []*SpookyError
    for _, err := range ae.Errors {
        if err.Severity == severity {
            filtered = append(filtered, err)
        }
    }
    return filtered
}
```

### Solution 2: Enhanced SSH Client Error Handling

**Improved SSH Client Using Existing Error Types:**

```go
// GOOD CODE - Enhanced error handling using existing SSH error types
func (c *Client) Connect() error {
    if c.config == nil {
        return spookytypesssh.NewConnectionError(
            "SSH_CLIENT_NIL_CONFIG",
            "SSH client configuration is nil",
            "ssh_client_connect",
            map[string]interface{}{
                "component": "ssh_client",
                "operation": "connect",
            },
        )
    }
    
    if c.config.Host == "" {
        return spookytypesssh.NewValidationError(
            "SSH_CLIENT_MISSING_HOST",
            "SSH host is required",
            "ssh_client_connect",
            map[string]interface{}{
                "component": "ssh_client",
                "operation": "connect",
                "config":    "host",
            },
        )
    }
    
    if c.config.Port <= 0 {
        return spookytypesssh.NewValidationError(
            "SSH_CLIENT_INVALID_PORT",
            "SSH port must be positive",
            "ssh_client_connect",
            map[string]interface{}{
                "component": "ssh_client",
                "operation": "connect",
                "port":      c.config.Port,
            },
        )
    }
    
    // Attempt connection with timeout
    ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
    defer cancel()
    
    client, err := ssh.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", c.config.Host, c.config.Port), c.config.ClientConfig)
    if err != nil {
        // Categorize SSH connection errors using existing error types
        if isTimeoutError(err) {
            return spookytypesssh.NewTimeoutError(
                "SSH_CLIENT_CONNECTION_TIMEOUT",
                "SSH connection timed out",
                "ssh_client_connect",
                map[string]interface{}{
                    "component": "ssh_client",
                    "operation": "connect",
                    "host":      c.config.Host,
                    "port":      c.config.Port,
                    "timeout":   c.config.Timeout,
                    "cause":     err.Error(),
                },
            )
        }
        
        if isAuthenticationError(err) {
            return spookytypesssh.NewAuthenticationError(
                "SSH_CLIENT_AUTHENTICATION_FAILED",
                "SSH authentication failed",
                "ssh_client_connect",
                map[string]interface{}{
                    "component": "ssh_client",
                    "operation": "connect",
                    "host":      c.config.Host,
                    "port":      c.config.Port,
                    "user":      c.config.User,
                    "cause":     err.Error(),
                },
            )
        }
        
        if isConnectionRefusedError(err) {
            return spookytypesssh.NewConnectionError(
                "SSH_CLIENT_CONNECTION_REFUSED",
                "SSH connection refused",
                "ssh_client_connect",
                map[string]interface{}{
                    "component": "ssh_client",
                    "operation": "connect",
                    "host":      c.config.Host,
                    "port":      c.config.Port,
                    "cause":     err.Error(),
                },
            )
        }
        
        // Generic connection error
        return spookytypesssh.WrapError(
            err,
            spookytypesssh.ErrorTypeConnection,
            "SSH_CLIENT_CONNECTION_FAILED",
            "SSH connection failed",
            "ssh_client_connect",
            map[string]interface{}{
                "component": "ssh_client",
                "operation": "connect",
                "host":      c.config.Host,
                "port":      c.config.Port,
            },
        )
    }
    
    c.client = client
    return nil
}

// Helper functions for error categorization
func isTimeoutError(err error) bool {
    if netErr, ok := err.(net.Error); ok {
        return netErr.Timeout()
    }
    return false
}

func isAuthenticationError(err error) bool {
    return strings.Contains(err.Error(), "ssh: unable to authenticate")
}

func isConnectionRefusedError(err error) bool {
    return strings.Contains(err.Error(), "connection refused")
}
```

### Solution 3: Enhanced Fact Manager Error Handling

**Improved Fact Manager Using Existing Error Types:**

```go
// GOOD CODE - Enhanced error handling with context using existing error types
func (m *Manager) CollectFacts(ctx context.Context, machine interface{}) (*spookytypes.FactCollection, error) {
    if machine == nil {
        return nil, spookytypesfacts.NewValidationError(
            "FACTS_COLLECTION_MISSING_MACHINE",
            "Machine cannot be nil for fact collection",
            "facts_collect",
            map[string]interface{}{
                "component": "facts_manager",
                "operation": "collect_facts",
                "field":     "machine",
            },
        )
    }

    // Type assert to get machine details
    machineObj, ok := machine.(*spookytypes.Machine)
    if !ok {
        return nil, spookytypesfacts.NewValidationError(
            "FACTS_COLLECTION_INVALID_MACHINE_TYPE",
            "Machine must be of type *spookytypes.Machine",
            "facts_collect",
            map[string]interface{}{
                "component": "facts_manager",
                "operation": "collect_facts",
                "field":     "machine",
                "type":      fmt.Sprintf("%T", machine),
            },
        )
    }

    if machineObj.Hostname == "" {
        return nil, spookytypesfacts.NewValidationError(
            "FACTS_COLLECTION_MISSING_HOSTNAME",
            "Machine hostname is required for fact collection",
            "facts_collect",
            map[string]interface{}{
                "component": "facts_manager",
                "operation": "collect_facts",
                "field":     "hostname",
            },
        )
    }

    m.logger.Info("Collecting facts", map[string]interface{}{
        "machine": machineObj.Hostname,
        "host":    machineObj.Host,
    })

    // Determine collection method based on machine configuration
    if machineObj.Host != "" && machineObj.Host != "localhost" && machineObj.Host != "127.0.0.1" {
        // Remote machine - use SSH-based collection
        return m.collectFactsViaSSH(ctx, machineObj)
    }

    // Local machine - use local collection
    facts, err := m.collector.Collect(ctx, machineObj)
    if err != nil {
        m.logger.Error("Failed to collect facts", err, map[string]interface{}{"machine": machineObj.Hostname})
        return nil, spookytypesfacts.WrapError(
            err,
            "FACTS_COLLECTION_FAILED",
            "Failed to collect facts from machine",
            "facts_collect",
            map[string]interface{}{
                "component": "facts_manager",
                "operation": "collect_facts",
                "machine":   machineObj.Hostname,
                "host":      machineObj.Host,
            },
        )
    }

    if facts == nil {
        return nil, spookytypesfacts.NewResourceError(
            "FACTS_COLLECTION_NIL_RESULT",
            "No facts collected from machine",
            "facts_collect",
            map[string]interface{}{
                "component": "facts_manager",
                "operation": "collect_facts",
                "machine":   machineObj.Hostname,
            },
        )
    }

    // Validate facts structure
    if err := m.validateFacts(facts); err != nil {
        return nil, spookytypesfacts.WrapError(
            err,
            "FACTS_COLLECTION_VALIDATION_FAILED",
            "Facts validation failed",
            "facts_collect",
            map[string]interface{}{
                "component": "facts_manager",
                "operation": "collect_facts",
                "machine":   machineObj.Hostname,
            },
        )
    }

    return facts, nil
}

func (m *Manager) validateFacts(facts *spookytypes.FactCollection) error {
    var errors []*spookytypesfacts.FactError
    
    if facts.Name == "" {
        errors = append(errors, spookytypesfacts.NewValidationError(
            "FACTS_VALIDATION_MISSING_NAME",
            "Fact collection name is required",
            "facts_validate",
            map[string]interface{}{
                "component": "facts_manager",
                "operation": "validate_facts",
                "field":     "name",
            },
        ))
    }
    
    if len(facts.Facts) == 0 {
        errors = append(errors, spookytypesfacts.NewValidationError(
            "FACTS_VALIDATION_EMPTY_COLLECTION",
            "Fact collection cannot be empty",
            "facts_validate",
            map[string]interface{}{
                "component": "facts_manager",
                "operation": "validate_facts",
                "field":     "facts",
            },
        ))
    }
    
    // Validate individual facts
    for i, fact := range facts.Facts {
        if err := m.validateFact(fact, i); err != nil {
            errors = append(errors, err)
        }
    }
    
    if len(errors) > 0 {
        return spookytypesfacts.NewAggregateError("Facts validation failed", errors)
    }
    
    return nil
}

func (m *Manager) validateFact(fact spookytypes.Fact, index int) error {
    if fact.Key == "" {
        return spookytypesfacts.NewValidationError(
            "FACT_VALIDATION_MISSING_KEY",
            "Fact key is required",
            "fact_validate",
            map[string]interface{}{
                "component": "facts_manager",
                "operation": "validate_fact",
                "field":     "key",
                "index":     index,
            },
        )
    }
    
    if fact.Value == "" {
        return spookytypesfacts.NewValidationError(
            "FACT_VALIDATION_MISSING_VALUE",
            "Fact value is required",
            "fact_validate",
            map[string]interface{}{
                "component": "facts_manager",
                "operation": "validate_fact",
                "field":     "value",
                "index":     index,
                "key":       fact.Key,
            },
        )
    }
    
    return nil
}
```

### Solution 4: Enhanced Configuration Loader with Error Recovery

**Improved Configuration Loader Using Existing Error Types:**

```go
// GOOD CODE - Enhanced error handling with recovery using existing error types
func LoadConfig(path string) (*spookytypes.Config, error) {
    // Try to load configuration from specified path
    config, err := loadConfigFromPath(path)
    if err != nil {
        // Try fallback configuration
        fallbackConfig, fallbackErr := loadFallbackConfig()
        if fallbackErr != nil {
            // Both primary and fallback failed
            return nil, spookytypesconfig.NewAggregateError("Configuration loading failed", []*spookytypesconfig.Error{
                spookytypesconfig.WrapError(err, "CONFIG_LOAD_PRIMARY_FAILED", "Primary configuration loading failed", "config_load", map[string]interface{}{"path": path}),
                spookytypesconfig.WrapError(fallbackErr, "CONFIG_LOAD_FALLBACK_FAILED", "Fallback configuration loading failed", "config_load", map[string]interface{}{}),
            })
        }
        
        // Log warning about using fallback
        log.Warn("Using fallback configuration", map[string]interface{}{
            "primary_path": path,
            "fallback_reason": err.Error(),
        })
        
        return fallbackConfig, nil
    }
    
    return config, nil
}

func loadConfigFromPath(path string) (*spookytypes.Config, error) {
    // Check if file exists
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil, spookytypesconfig.NewError(
            path,
            "Configuration file not found",
        )
    }
    
    // Read file
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, spookytypesconfig.WrapError(
            err,
            "CONFIG_FILE_READ_ERROR",
            "Failed to read configuration file",
            "config_load_from_path",
            map[string]interface{}{
                "component": "config_loader",
                "operation": "load_config",
                "path":      path,
            },
        )
    }
    
    // Parse HCL
    var config spookytypes.Config
    if err := hcl.Unmarshal(data, &config); err != nil {
        return nil, spookytypesconfig.WrapError(
            err,
            "CONFIG_PARSE_ERROR",
            "Failed to parse configuration file",
            "config_load_from_path",
            map[string]interface{}{
                "component": "config_loader",
                "operation": "load_config",
                "path":      path,
                "format":    "hcl",
            },
        )
    }
    
    // Validate configuration
    if err := config.Validate(); err != nil {
        return nil, spookytypesconfig.WrapError(
            err,
            "CONFIG_VALIDATION_ERROR",
            "Configuration validation failed",
            "config_load_from_path",
            map[string]interface{}{
                "component": "config_loader",
                "operation": "load_config",
                "path":      path,
            },
        )
    }
    
    return &config, nil
}

func loadFallbackConfig() (*spookytypes.Config, error) {
    // Try to load from default locations
    defaultPaths := []string{
        "~/.config/spooky/spooky.hcl",
        "/etc/spooky/spooky.hcl",
        "./spooky.hcl",
    }
    
    for _, path := range defaultPaths {
        if expandedPath, err := expandPath(path); err == nil {
            if config, err := loadConfigFromPath(expandedPath); err == nil {
                return config, nil
            }
        }
    }
    
    // Create default configuration
    return createDefaultConfig(), nil
}

func createDefaultConfig() *spookytypes.Config {
    return &spookytypes.Config{
        Logging: spookytypes.LoggingConfig{
            Level:  "info",
            Format: "json",
        },
        SSH: spookytypes.SSHConfig{
            Timeout: 30 * time.Second,
            Port:    22,
        },
        Facts: spookytypes.FactsConfig{
            StoragePath: "~/.local/state/spooky/facts.db",
        },
    }
}

func expandPath(path string) (string, error) {
    if strings.HasPrefix(path, "~/") {
        homeDir, err := os.UserHomeDir()
        if err != nil {
            return "", err
        }
        return filepath.Join(homeDir, path[2:]), nil
    }
    return path, nil
}
```

## Implementation Plan

### Phase 1: Error Type Standardization
1. **Leverage Existing Error Types**: Use existing domain-specific error types
2. **Enhance Error Context**: Add operation and component information to existing types
3. **Standardize Error Wrapping**: Implement consistent error wrapping patterns
4. **Add Error Recovery**: Implement error recovery suggestions

### Phase 2: Error Handling Migration
1. **Migrate SSH Package**: Update SSH client to use existing SSH error types
2. **Migrate Facts Package**: Update facts manager to use existing fact error types
3. **Migrate Actions Package**: Update actions manager to use existing action error types
4. **Migrate Config Package**: Update configuration loader to use existing config error types

### Phase 3: Error Aggregation and Recovery
1. **Implement Error Aggregation**: Add support for multiple errors using existing patterns
2. **Add Error Recovery Mechanisms**: Implement automatic error recovery
3. **Add Error Monitoring**: Implement error tracking and metrics
4. **Add Error Reporting**: Implement structured error reporting

### Phase 4: Error Testing and Validation
1. **Add Error Testing**: Implement comprehensive error testing
2. **Add Error Validation**: Validate error handling patterns
3. **Add Error Documentation**: Document error handling patterns
4. **Add Error Examples**: Provide error handling examples

## Error Handling Standards

### Error Categories (Leveraging Existing Types)
- **Validation**: Input validation and data format errors (using existing ValidationError types)
- **Connection**: Network and connection-related errors (using existing ConnectionError types)
- **Authentication**: Authentication and authorization errors (using existing AuthenticationError types)
- **Permission**: File and resource permission errors (using existing PermissionError types)
- **Timeout**: Operation timeout errors (using existing TimeoutError types)
- **Resource**: Resource availability and allocation errors (using existing ResourceError types)
- **Configuration**: Configuration and setup errors (using existing ConfigError types)
- **Internal**: Internal system and programming errors (using existing InternalError types)

### Error Severity Levels
- **Low**: Informational errors with minimal impact
- **Medium**: Errors that may affect functionality
- **High**: Errors that significantly impact operations
- **Critical**: Errors that prevent system operation

### Error Context Requirements
- **Component**: Component or package name
- **Operation**: Specific operation being performed
- **Timestamp**: When the error occurred (from ErrorDetails)
- **Context**: Additional contextual information (from ErrorDetails)
- **Suggestions**: Error recovery suggestions
- **Recoverable**: Whether the error is recoverable (from ErrorDetails)

### Error Wrapping Standards
- **Consistent Wrapping**: Use consistent error wrapping patterns with %w
- **Context Preservation**: Preserve original error context using existing ErrorDetails
- **Categorization**: Categorize errors appropriately using existing error types
- **Recovery Information**: Include recovery suggestions
- **Error Codes**: Use consistent error codes from existing error types

## Benefits

1. **Improved Debugging**: Better error context and categorization using existing infrastructure
2. **Enhanced Recovery**: Automatic error recovery mechanisms
3. **Better Monitoring**: Comprehensive error tracking and metrics
4. **Consistent Handling**: Standardized error handling patterns leveraging existing types
5. **User Experience**: Better error messages and recovery suggestions
6. **Maintainability**: Easier error handling maintenance using established patterns
7. **Reliability**: More robust error handling and recovery
8. **Documentation**: Better error documentation and examples

## Conclusion

The error handling standardization effort will significantly improve the reliability and maintainability of the spooky codebase by leveraging existing structured error types and providing consistent, comprehensive error handling patterns with proper context and recovery mechanisms.

**Priority Actions:**
1. **Immediate**: Leverage existing error types and enhance them with additional context
2. **High**: Migrate existing error handling to use established error type patterns
3. **Medium**: Add error aggregation and recovery mechanisms using existing infrastructure
4. **Long-term**: Implement comprehensive error monitoring and reporting

**Expected Outcomes:**
- Consistent error handling across all packages using existing error types
- Better error context and debugging information
- Automatic error recovery mechanisms
- Comprehensive error monitoring and metrics
- Improved user experience with better error messages
- Enhanced system reliability and maintainability

The standardization effort will make the spooky codebase more robust and easier to maintain while building on the existing error handling infrastructure rather than replacing it.

**Key Corrections from Original Report:**
1. **Leverage Existing Infrastructure**: The codebase already has comprehensive error types that should be used
2. **Use Domain-Specific Errors**: Each domain has its own error types that should be leveraged
3. **Build on ErrorDetails**: The common ErrorDetails type provides the foundation for structured errors
4. **Maintain Interface Patterns**: Error handling should follow established interface patterns
5. **Use Existing Error Categories**: The existing error categorization should be extended rather than replaced
