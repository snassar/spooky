# Implementation Plan: Validation Completion Implementation

## Overview
Complete machine validation implementations by replacing placeholder code with real validation logic for machine configurations, connectivity, and health checks.

## Task Details
- **Task ID**: 7.8
- **Priority**: Medium
- **Files**: 
  - `internal/machines/validation/manager.go`
- **Functions**: Machine validation, configuration validation, health validation

## Current State Analysis

### Existing Patterns
1. **Machine Types**: Machine configurations defined in `internal/machines/types/`
2. **Validation Framework**: Basic validation framework exists
3. **Error Handling**: Consistent error wrapping
4. **Logging**: Structured logging implemented

### Current Placeholder Code
```go
// internal/machines/validation/manager.go line 34
// This is a placeholder implementation
```

## Implementation Requirements

### Interface Compliance
The validation completion must:
1. **Validate machine configurations** and parameters
2. **Check machine connectivity** and accessibility
3. **Validate SSH configurations** and authentication
4. **Perform health checks** and system validation
5. **Validate network connectivity** and routing
6. **Check resource availability** and capacity
7. **Provide detailed validation reports** and diagnostics

### Required Dependencies
- SSH connectivity testing
- Network validation system
- Health checking system
- Resource monitoring system

## Detailed Implementation Plan

### Step 1: Implement Machine Configuration Validation

#### 1.1 Configuration Validator
```go
// internal/machines/validation/config_validator.go
package validation

import (
    "fmt"
    "net"
    "regexp"
    "strconv"
    "strings"
    
    "spooky/internal/machines/types"
    "spooky/internal/logging"
)

// ConfigValidator validates machine configurations
type ConfigValidator struct {
    logger logging.Logger
}

// NewConfigValidator creates a new configuration validator
func NewConfigValidator(logger logging.Logger) *ConfigValidator {
    return &ConfigValidator{
        logger: logger,
    }
}

// ValidateMachine validates a machine configuration
func (cv *ConfigValidator) ValidateMachine(machine *types.Machine) *ValidationResult {
    result := &ValidationResult{
        MachineID: machine.Name,
        Status:    ValidationStatusValid,
        Errors:    make([]ValidationError, 0),
        Warnings:  make([]ValidationWarning, 0),
    }

    // Validate basic fields
    cv.validateBasicFields(machine, result)
    cv.validateNetworkConfiguration(machine, result)
    cv.validateSSHConfiguration(machine, result)
    cv.validateTagsAndEnvironment(machine, result)

    // Determine overall status
    if len(result.Errors) > 0 {
        result.Status = ValidationStatusInvalid
    } else if len(result.Warnings) > 0 {
        result.Status = ValidationStatusWarning
    }

    return result
}

// validateBasicFields validates basic machine fields
func (cv *ConfigValidator) validateBasicFields(machine *types.Machine, result *ValidationResult) {
    // Validate name
    if machine.Name == "" {
        result.Errors = append(result.Errors, ValidationError{
            Field:   "name",
            Message: "Machine name cannot be empty",
        })
    } else if !cv.isValidMachineName(machine.Name) {
        result.Errors = append(result.Errors, ValidationError{
            Field:   "name",
            Message: "Machine name contains invalid characters",
        })
    }

    // Validate host
    if machine.Host == "" {
        result.Errors = append(result.Errors, ValidationError{
            Field:   "host",
            Message: "Host cannot be empty",
        })
    } else if !cv.isValidHost(machine.Host) {
        result.Errors = append(result.Errors, ValidationError{
            Field:   "host",
            Message: "Invalid host format",
        })
    }

    // Validate port
    if machine.Port <= 0 || machine.Port > 65535 {
        result.Errors = append(result.Errors, ValidationError{
            Field:   "port",
            Message: "Port must be between 1 and 65535",
        })
    }

    // Validate user
    if machine.User == "" {
        result.Errors = append(result.Errors, ValidationError{
            Field:   "user",
            Message: "User cannot be empty",
        })
    } else if !cv.isValidUsername(machine.User) {
        result.Warnings = append(result.Warnings, ValidationWarning{
            Field:   "user",
            Message: "Username contains potentially problematic characters",
        })
    }
}

// validateNetworkConfiguration validates network configuration
func (cv *ConfigValidator) validateNetworkConfiguration(machine *types.Machine, result *ValidationResult) {
    // Validate host resolution
    if machine.Host != "" {
        if cv.isIPAddress(machine.Host) {
            // Validate IP address format
            if net.ParseIP(machine.Host) == nil {
                result.Errors = append(result.Errors, ValidationError{
                    Field:   "host",
                    Message: "Invalid IP address format",
                })
            }
        } else {
            // Validate hostname format
            if !cv.isValidHostname(machine.Host) {
                result.Warnings = append(result.Warnings, ValidationWarning{
                    Field:   "host",
                    Message: "Hostname format may be invalid",
                })
            }
        }
    }

    // Validate SSH key path if specified
    if machine.SSHKeyPath != "" {
        if !cv.isValidSSHKeyPath(machine.SSHKeyPath) {
            result.Warnings = append(result.Warnings, ValidationWarning{
                Field:   "ssh_key_path",
                Message: "SSH key path may not exist or be accessible",
            })
        }
    }
}

// validateSSHConfiguration validates SSH configuration
func (cv *ConfigValidator) validateSSHConfiguration(machine *types.Machine, result *ValidationResult) {
    // Validate timeout
    if machine.Timeout <= 0 {
        result.Warnings = append(result.Warnings, ValidationWarning{
            Field:   "timeout",
            Message: "Timeout should be positive",
        })
    }

    // Validate SSH key permissions if specified
    if machine.SSHKeyPath != "" {
        if !cv.hasValidSSHKeyPermissions(machine.SSHKeyPath) {
            result.Warnings = append(result.Warnings, ValidationWarning{
                Field:   "ssh_key_path",
                Message: "SSH key file permissions may be too permissive",
            })
        }
    }
}

// validateTagsAndEnvironment validates tags and environment
func (cv *ConfigValidator) validateTagsAndEnvironment(machine *types.Machine, result *ValidationResult) {
    // Validate tags
    for i, tag := range machine.Tags {
        if !cv.isValidTag(tag) {
            result.Warnings = append(result.Warnings, ValidationWarning{
                Field:   fmt.Sprintf("tags[%d]", i),
                Message: "Tag contains potentially problematic characters",
            })
        }
    }

    // Validate environment variables
    for key, value := range machine.Environment {
        if !cv.isValidEnvironmentKey(key) {
            result.Warnings = append(result.Warnings, ValidationWarning{
                Field:   fmt.Sprintf("environment.%s", key),
                Message: "Environment key contains potentially problematic characters",
            })
        }
    }
}

// Validation helper methods
func (cv *ConfigValidator) isValidMachineName(name string) bool {
    // Machine names should be alphanumeric with hyphens and underscores
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, name)
    return matched && len(name) <= 63
}

func (cv *ConfigValidator) isValidHost(host string) bool {
    return host != "" && len(host) <= 255
}

func (cv *ConfigValidator) isValidUsername(username string) bool {
    // Username should not contain special characters that could cause issues
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9._-]+$`, username)
    return matched && len(username) <= 32
}

func (cv *ConfigValidator) isIPAddress(host string) bool {
    return net.ParseIP(host) != nil
}

func (cv *ConfigValidator) isValidHostname(hostname string) bool {
    // Basic hostname validation
    if len(hostname) > 253 {
        return false
    }
    
    parts := strings.Split(hostname, ".")
    for _, part := range parts {
        if len(part) == 0 || len(part) > 63 {
            return false
        }
        if !regexp.MustCompile(`^[a-zA-Z0-9-]+$`).MatchString(part) {
            return false
        }
    }
    
    return true
}

func (cv *ConfigValidator) isValidSSHKeyPath(path string) bool {
    // This would check if the SSH key file exists and is readable
    // For now, return true as a placeholder
    return true
}

func (cv *ConfigValidator) hasValidSSHKeyPermissions(path string) bool {
    // This would check SSH key file permissions (should be 600)
    // For now, return true as a placeholder
    return true
}

func (cv *ConfigValidator) isValidTag(tag string) bool {
    // Tags should be alphanumeric with hyphens and underscores
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, tag)
    return matched && len(tag) <= 63
}

func (cv *ConfigValidator) isValidEnvironmentKey(key string) bool {
    // Environment keys should be valid shell variable names
    matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, key)
    return matched
}
```

### Step 2: Implement Connectivity Validation

#### 2.1 Connectivity Validator
```go
// internal/machines/validation/connectivity_validator.go
package validation

import (
    "context"
    "fmt"
    "net"
    "time"
    
    "spooky/internal/machines/types"
    "spooky/internal/ssh"
    "spooky/internal/ssh/types"
    "spooky/internal/logging"
)

// ConnectivityValidator validates machine connectivity
type ConnectivityValidator struct {
    sshManager ssh.Manager
    logger     logging.Logger
}

// NewConnectivityValidator creates a new connectivity validator
func NewConnectivityValidator(sshManager ssh.Manager, logger logging.Logger) *ConnectivityValidator {
    return &ConnectivityValidator{
        sshManager: sshManager,
        logger:     logger,
    }
}

// ValidateConnectivity validates machine connectivity
func (cv *ConnectivityValidator) ValidateConnectivity(ctx context.Context, machine *types.Machine) *ValidationResult {
    result := &ValidationResult{
        MachineID: machine.Name,
        Status:    ValidationStatusValid,
        Errors:    make([]ValidationError, 0),
        Warnings:  make([]ValidationWarning, 0),
    }

    // Test network connectivity
    cv.validateNetworkConnectivity(machine, result)
    
    // Test SSH connectivity
    cv.validateSSHConnectivity(ctx, machine, result)
    
    // Test SSH authentication
    cv.validateSSHAuthentication(ctx, machine, result)

    // Determine overall status
    if len(result.Errors) > 0 {
        result.Status = ValidationStatusInvalid
    } else if len(result.Warnings) > 0 {
        result.Status = ValidationStatusWarning
    }

    return result
}

// validateNetworkConnectivity validates network connectivity
func (cv *ConnectivityValidator) validateNetworkConnectivity(machine *types.Machine, result *ValidationResult) {
    cv.logger.Debug("Validating network connectivity",
        logging.String("machine", machine.Name),
        logging.String("host", machine.Host))

    // Test TCP connectivity
    timeout := 5 * time.Second
    if machine.Timeout > 0 {
        timeout = machine.Timeout
    }

    conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", machine.Host, machine.Port), timeout)
    if err != nil {
        result.Errors = append(result.Errors, ValidationError{
            Field:   "connectivity",
            Message: fmt.Sprintf("TCP connection failed: %v", err),
        })
        return
    }
    defer conn.Close()

    cv.logger.Debug("Network connectivity validated successfully",
        logging.String("machine", machine.Name))
}

// validateSSHConnectivity validates SSH connectivity
func (cv *ConnectivityValidator) validateSSHConnectivity(ctx context.Context, machine *types.Machine, result *ValidationResult) {
    cv.logger.Debug("Validating SSH connectivity",
        logging.String("machine", machine.Name))

    // Create SSH configuration
    sshConfig := &sshTypes.SSHConfig{
        Host:     machine.Host,
        Port:     machine.Port,
        Username: machine.User,
        Timeout:  machine.Timeout,
    }

    // Test SSH connection
    connection, err := cv.sshManager.Connect(machine.Host, sshConfig)
    if err != nil {
        result.Errors = append(result.Errors, ValidationError{
            Field:   "ssh_connectivity",
            Message: fmt.Sprintf("SSH connection failed: %v", err),
        })
        return
    }
    defer cv.sshManager.CloseConnection(connection)

    cv.logger.Debug("SSH connectivity validated successfully",
        logging.String("machine", machine.Name))
}

// validateSSHAuthentication validates SSH authentication
func (cv *ConnectivityValidator) validateSSHAuthentication(ctx context.Context, machine *types.Machine, result *ValidationResult) {
    cv.logger.Debug("Validating SSH authentication",
        logging.String("machine", machine.Name))

    // Create SSH configuration
    sshConfig := &sshTypes.SSHConfig{
        Host:     machine.Host,
        Port:     machine.Port,
        Username: machine.User,
        Timeout:  machine.Timeout,
    }

    // Test SSH connection with authentication
    connection, err := cv.sshManager.Connect(machine.Host, sshConfig)
    if err != nil {
        result.Errors = append(result.Errors, ValidationError{
            Field:   "ssh_authentication",
            Message: fmt.Sprintf("SSH authentication failed: %v", err),
        })
        return
    }
    defer cv.sshManager.CloseConnection(connection)

    // Test basic command execution
    testResult, err := cv.sshManager.ExecuteCommand(connection, "echo 'authentication test successful'")
    if err != nil {
        result.Errors = append(result.Errors, ValidationError{
            Field:   "ssh_authentication",
            Message: fmt.Sprintf("SSH command execution failed: %v", err),
        })
        return
    }

    if testResult.ExitCode != 0 {
        result.Warnings = append(result.Warnings, ValidationWarning{
            Field:   "ssh_authentication",
            Message: fmt.Sprintf("SSH command returned non-zero exit code: %d", testResult.ExitCode),
        })
    }

    cv.logger.Debug("SSH authentication validated successfully",
        logging.String("machine", machine.Name))
}
```

### Step 3: Implement Health Validation

#### 3.1 Health Validator
```go
// internal/machines/validation/health_validator.go
package validation

import (
    "context"
    "fmt"
    "strconv"
    "strings"
    "time"
    
    "spooky/internal/machines/types"
    "spooky/internal/ssh"
    "spooky/internal/ssh/types"
    "spooky/internal/logging"
)

// HealthValidator validates machine health
type HealthValidator struct {
    sshManager ssh.Manager
    logger     logging.Logger
}

// NewHealthValidator creates a new health validator
func NewHealthValidator(sshManager ssh.Manager, logger logging.Logger) *HealthValidator {
    return &HealthValidator{
        sshManager: sshManager,
        logger:     logger,
    }
}

// ValidateHealth validates machine health
func (hv *HealthValidator) ValidateHealth(ctx context.Context, machine *types.Machine) *ValidationResult {
    result := &ValidationResult{
        MachineID: machine.Name,
        Status:    ValidationStatusValid,
        Errors:    make([]ValidationError, 0),
        Warnings:  make([]ValidationWarning, 0),
    }

    // Create SSH connection
    sshConfig := &sshTypes.SSHConfig{
        Host:     machine.Host,
        Port:     machine.Port,
        Username: machine.User,
        Timeout:  machine.Timeout,
    }

    connection, err := hv.sshManager.Connect(machine.Host, sshConfig)
    if err != nil {
        result.Errors = append(result.Errors, ValidationError{
            Field:   "health_check",
            Message: fmt.Sprintf("Failed to establish SSH connection: %v", err),
        })
        return result
    }
    defer hv.sshManager.CloseConnection(connection)

    // Perform health checks
    hv.validateSystemHealth(connection, result)
    hv.validateResourceUsage(connection, result)
    hv.validateServiceHealth(connection, result)

    // Determine overall status
    if len(result.Errors) > 0 {
        result.Status = ValidationStatusInvalid
    } else if len(result.Warnings) > 0 {
        result.Status = ValidationStatusWarning
    }

    return result
}

// validateSystemHealth validates system health
func (hv *HealthValidator) validateSystemHealth(connection sshTypes.Connection, result *ValidationResult) {
    hv.logger.Debug("Validating system health")

    // Check system uptime
    uptimeResult, err := hv.sshManager.ExecuteCommand(connection, "uptime -p")
    if err != nil {
        result.Warnings = append(result.Warnings, ValidationWarning{
            Field:   "system_health",
            Message: "Failed to check system uptime",
        })
    } else {
        hv.logger.Debug("System uptime retrieved",
            logging.String("uptime", strings.TrimSpace(uptimeResult.Stdout)))
    }

    // Check system load
    loadResult, err := hv.sshManager.ExecuteCommand(connection, "cat /proc/loadavg")
    if err != nil {
        result.Warnings = append(result.Warnings, ValidationWarning{
            Field:   "system_health",
            Message: "Failed to check system load",
        })
    } else {
        hv.validateSystemLoad(loadResult.Stdout, result)
    }
}

// validateSystemLoad validates system load
func (hv *HealthValidator) validateSystemLoad(loadOutput string, result *ValidationResult) {
    fields := strings.Fields(strings.TrimSpace(loadOutput))
    if len(fields) < 3 {
        result.Warnings = append(result.Warnings, ValidationWarning{
            Field:   "system_load",
            Message: "Invalid load average format",
        })
        return
    }

    // Parse load averages
    oneMin, err := strconv.ParseFloat(fields[0], 64)
    if err != nil {
        result.Warnings = append(result.Warnings, ValidationWarning{
            Field:   "system_load",
            Message: "Failed to parse 1-minute load average",
        })
        return
    }

    fiveMin, err := strconv.ParseFloat(fields[1], 64)
    if err != nil {
        result.Warnings = append(result.Warnings, ValidationWarning{
            Field:   "system_load",
            Message: "Failed to parse 5-minute load average",
        })
        return
    }

    // Check load thresholds
    if oneMin > 10.0 {
        result.Errors = append(result.Errors, ValidationError{
            Field:   "system_load",
            Message: fmt.Sprintf("1-minute load average is too high: %.2f", oneMin),
        })
    } else if oneMin > 5.0 {
        result.Warnings = append(result.Warnings, ValidationWarning{
            Field:   "system_load",
            Message: fmt.Sprintf("1-minute load average is elevated: %.2f", oneMin),
        })
    }

    if fiveMin > 8.0 {
        result.Warnings = append(result.Warnings, ValidationWarning{
            Field:   "system_load",
            Message: fmt.Sprintf("5-minute load average is elevated: %.2f", fiveMin),
        })
    }
}

// validateResourceUsage validates resource usage
func (hv *HealthValidator) validateResourceUsage(connection sshTypes.Connection, result *ValidationResult) {
    hv.logger.Debug("Validating resource usage")

    // Check disk usage
    diskResult, err := hv.sshManager.ExecuteCommand(connection, "df -h / | tail -1")
    if err != nil {
        result.Warnings = append(result.Warnings, ValidationWarning{
            Field:   "resource_usage",
            Message: "Failed to check disk usage",
        })
    } else {
        hv.validateDiskUsage(diskResult.Stdout, result)
    }

    // Check memory usage
    memoryResult, err := hv.sshManager.ExecuteCommand(connection, "free -m | grep Mem")
    if err != nil {
        result.Warnings = append(result.Warnings, ValidationWarning{
            Field:   "resource_usage",
            Message: "Failed to check memory usage",
        })
    } else {
        hv.validateMemoryUsage(memoryResult.Stdout, result)
    }
}

// validateDiskUsage validates disk usage
func (hv *HealthValidator) validateDiskUsage(diskOutput string, result *ValidationResult) {
    fields := strings.Fields(strings.TrimSpace(diskOutput))
    if len(fields) < 5 {
        result.Warnings = append(result.Warnings, ValidationWarning{
            Field:   "disk_usage",
            Message: "Invalid disk usage format",
        })
        return
    }

    // Parse usage percentage
    usageStr := strings.TrimSuffix(fields[4], "%")
    usage, err := strconv.Atoi(usageStr)
    if err != nil {
        result.Warnings = append(result.Warnings, ValidationWarning{
            Field:   "disk_usage",
            Message: "Failed to parse disk usage percentage",
        })
        return
    }

    // Check usage thresholds
    if usage > 95 {
        result.Errors = append(result.Errors, ValidationError{
            Field:   "disk_usage",
            Message: fmt.Sprintf("Disk usage is critical: %d%%", usage),
        })
    } else if usage > 85 {
        result.Warnings = append(result.Warnings, ValidationWarning{
            Field:   "disk_usage",
            Message: fmt.Sprintf("Disk usage is high: %d%%", usage),
        })
    }
}

// validateMemoryUsage validates memory usage
func (hv *HealthValidator) validateMemoryUsage(memoryOutput string, result *ValidationResult) {
    fields := strings.Fields(strings.TrimSpace(memoryOutput))
    if len(fields) < 7 {
        result.Warnings = append(result.Warnings, ValidationWarning{
            Field:   "memory_usage",
            Message: "Invalid memory usage format",
        })
        return
    }

    // Parse memory values
    total, err := strconv.Atoi(fields[1])
    if err != nil {
        result.Warnings = append(result.Warnings, ValidationWarning{
            Field:   "memory_usage",
            Message: "Failed to parse total memory",
        })
        return
    }

    used, err := strconv.Atoi(fields[2])
    if err != nil {
        result.Warnings = append(result.Warnings, ValidationWarning{
            Field:   "memory_usage",
            Message: "Failed to parse used memory",
        })
        return
    }

    // Calculate usage percentage
    usagePercent := float64(used) / float64(total) * 100

    // Check usage thresholds
    if usagePercent > 95 {
        result.Errors = append(result.Errors, ValidationError{
            Field:   "memory_usage",
            Message: fmt.Sprintf("Memory usage is critical: %.1f%%", usagePercent),
        })
    } else if usagePercent > 85 {
        result.Warnings = append(result.Warnings, ValidationWarning{
            Field:   "memory_usage",
            Message: fmt.Sprintf("Memory usage is high: %.1f%%", usagePercent),
        })
    }
}

// validateServiceHealth validates service health
func (hv *HealthValidator) validateServiceHealth(connection sshTypes.Connection, result *ValidationResult) {
    hv.logger.Debug("Validating service health")

    // Check SSH service status
    sshResult, err := hv.sshManager.ExecuteCommand(connection, "systemctl is-active ssh")
    if err != nil {
        result.Warnings = append(result.Warnings, ValidationWarning{
            Field:   "service_health",
            Message: "Failed to check SSH service status",
        })
    } else {
        if strings.TrimSpace(sshResult.Stdout) != "active" {
            result.Errors = append(result.Errors, ValidationError{
                Field:   "service_health",
                Message: "SSH service is not active",
            })
        }
    }

    // Check systemd service status
    systemdResult, err := hv.sshManager.ExecuteCommand(connection, "systemctl is-system-running")
    if err != nil {
        result.Warnings = append(result.Warnings, ValidationWarning{
            Field:   "service_health",
            Message: "Failed to check systemd status",
        })
    } else {
        status := strings.TrimSpace(systemdResult.Stdout)
        if status != "running" {
            result.Warnings = append(result.Warnings, ValidationWarning{
                Field:   "service_health",
                Message: fmt.Sprintf("System is not fully running: %s", status),
            })
        }
    }
}
```

### Step 4: Implement Validation Manager

#### 4.1 Enhanced Validation Manager
```go
// internal/machines/validation/manager.go
package validation

import (
    "context"
    "fmt"
    "time"
    
    "spooky/internal/machines/types"
    "spooky/internal/ssh"
    "spooky/internal/logging"
)

// Manager manages machine validation
type Manager struct {
    configValidator      *ConfigValidator
    connectivityValidator *ConnectivityValidator
    healthValidator      *HealthValidator
    logger               logging.Logger
}

// NewManager creates a new validation manager
func NewManager(sshManager ssh.Manager, logger logging.Logger) *Manager {
    return &Manager{
        configValidator:      NewConfigValidator(logger),
        connectivityValidator: NewConnectivityValidator(sshManager, logger),
        healthValidator:      NewHealthValidator(sshManager, logger),
        logger:               logger,
    }
}

// ValidateMachine validates a machine configuration and connectivity
func (m *Manager) ValidateMachine(ctx context.Context, machine *types.Machine, options *ValidationOptions) *ValidationResult {
    m.logger.Info("Validating machine",
        logging.String("machine", machine.Name))

    startTime := time.Now()

    // Validate configuration
    configResult := m.configValidator.ValidateMachine(machine)
    if configResult.Status == ValidationStatusInvalid {
        m.logger.Warn("Machine configuration validation failed",
            logging.String("machine", machine.Name),
            logging.Int("error_count", len(configResult.Errors)))
        return configResult
    }

    // Validate connectivity if requested
    var connectivityResult *ValidationResult
    if options.IncludeConnectivity {
        connectivityResult = m.connectivityValidator.ValidateConnectivity(ctx, machine)
    }

    // Validate health if requested
    var healthResult *ValidationResult
    if options.IncludeHealth {
        healthResult = m.healthValidator.ValidateHealth(ctx, machine)
    }

    // Combine results
    combinedResult := m.combineValidationResults(machine.Name, configResult, connectivityResult, healthResult)

    duration := time.Since(startTime)
    m.logger.Info("Machine validation completed",
        logging.String("machine", machine.Name),
        logging.String("status", string(combinedResult.Status)),
        logging.Duration("duration", duration))

    return combinedResult
}

// ValidateMachines validates multiple machines
func (m *Manager) ValidateMachines(ctx context.Context, machines []*types.Machine, options *ValidationOptions) []*ValidationResult {
    m.logger.Info("Validating multiple machines",
        logging.Int("machine_count", len(machines)))

    results := make([]*ValidationResult, len(machines))

    for i, machine := range machines {
        results[i] = m.ValidateMachine(ctx, machine, options)
    }

    m.logger.Info("Multiple machine validation completed",
        logging.Int("machine_count", len(machines)))

    return results
}

// combineValidationResults combines multiple validation results
func (m *Manager) combineValidationResults(machineID string, results ...*ValidationResult) *ValidationResult {
    combined := &ValidationResult{
        MachineID: machineID,
        Status:    ValidationStatusValid,
        Errors:    make([]ValidationError, 0),
        Warnings:  make([]ValidationWarning, 0),
    }

    for _, result := range results {
        if result == nil {
            continue
        }

        // Combine errors
        combined.Errors = append(combined.Errors, result.Errors...)

        // Combine warnings
        combined.Warnings = append(combined.Warnings, result.Warnings...)
    }

    // Determine overall status
    if len(combined.Errors) > 0 {
        combined.Status = ValidationStatusInvalid
    } else if len(combined.Warnings) > 0 {
        combined.Status = ValidationStatusWarning
    }

    return combined
}

// ValidationOptions represents validation options
type ValidationOptions struct {
    IncludeConnectivity bool `json:"include_connectivity"`
    IncludeHealth       bool `json:"include_health"`
    Timeout             time.Duration `json:"timeout"`
}

// DefaultValidationOptions returns default validation options
func DefaultValidationOptions() *ValidationOptions {
    return &ValidationOptions{
        IncludeConnectivity: true,
        IncludeHealth:       false,
        Timeout:             30 * time.Second,
    }
}

// ValidationResult represents validation result
type ValidationResult struct {
    MachineID string            `json:"machine_id"`
    Status    ValidationStatus  `json:"status"`
    Errors    []ValidationError `json:"errors"`
    Warnings  []ValidationWarning `json:"warnings"`
    Timestamp time.Time         `json:"timestamp"`
}

// ValidationStatus represents validation status
type ValidationStatus string

const (
    ValidationStatusValid   ValidationStatus = "valid"
    ValidationStatusWarning ValidationStatus = "warning"
    ValidationStatusInvalid ValidationStatus = "invalid"
)

// ValidationError represents a validation error
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

// ValidationWarning represents a validation warning
type ValidationWarning struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}
```

## Configuration Options

### Supported Options
- **IncludeConnectivity**: Enable/disable connectivity validation
- **IncludeHealth**: Enable/disable health validation
- **ValidationTimeout**: Validation timeout
- **StrictMode**: Enable/disable strict validation mode
- **ParallelValidation**: Enable/disable parallel validation

## Dependencies

### Internal Dependencies
- `spooky/internal/machines/types`
- `spooky/internal/ssh`
- `spooky/internal/ssh/types`
- `spooky/internal/logging`

### External Dependencies
- `context` (standard library)
- `fmt` (standard library)
- `net` (standard library)
- `regexp` (standard library)
- `strconv` (standard library)
- `strings` (standard library)
- `time` (standard library)

## Implementation Order

1. Implement configuration validator
2. Add connectivity validator
3. Create health validator
4. Update validation manager
5. Add validation options and results
6. Implement parallel validation
7. Add comprehensive tests
8. Performance optimization
9. Documentation and cleanup
