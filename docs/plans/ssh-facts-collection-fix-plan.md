# SSH-Based Fact Collection Fix Plan

## Problem Analysis

### Current State
The current implementation has a critical gap in SSH-based fact collection:

1. **Machine Inventory Loading**: ✅ Works correctly - loads machines from `machines.hcl`
2. **Local Fact Collection**: ✅ Works but only for local machine
3. **SSH-based Fact Collection**: ❌ **BROKEN** - not actually connecting to remote machines
4. **Remote `/etc/spooky/facts.*` Reading**: ❌ **BROKEN** - not reading from remote machines

### Root Cause
The issue is in the facts integration and CLI command implementation:

1. **Facts Integration (`internal/facts/integration.go`)**:
   - `CollectFacts()` method creates a fake machine object with hardcoded values
   - Uses `source` parameter as hostname but ignores actual machine configuration
   - No SSH connection is established to remote machines

2. **CLI Command (`cmd/facts.go`)**:
   - Loads machines correctly from project
   - Passes `machine.Hostname` as source to `CollectFacts()`
   - But the integration doesn't use the actual machine configuration for SSH

3. **Missing SSH Integration**:
   - Facts collector has SSH manager but it's not being used properly
   - No actual SSH connections to remote machines
   - No reading of remote `/etc/spooky/facts.hcl` or `/etc/spooky/custom.hcl`

## Detailed Implementation Plan

### Phase 1: Fix Facts Integration (Priority: Critical)

#### 1.1 Update Facts Integration Interface
**File**: `internal/facts/integration.go`

**Current Problem**:
```go
// CollectFacts collects facts from the given source
func (i *Integration) CollectFacts(ctx context.Context, source string) (interface{}, error) {
    // Create a machine representation for local fact collection
    // The source parameter represents the local machine identifier
    machine := &spookytypes.Machine{
        Hostname: source,
        Host:     source,
        Port:     22,
        User:     "root",
    }
    // ...
}
```

**Solution**:
```go
// CollectFacts collects facts from the given machine
func (i *Integration) CollectFacts(ctx context.Context, machine *spookytypes.Machine) (interface{}, error) {
    if machine == nil {
        return nil, fmt.Errorf("machine cannot be nil")
    }

    // Collect facts using the manager with actual machine configuration
    facts, err := i.manager.CollectFacts(ctx, machine)
    if err != nil {
        return nil, fmt.Errorf("failed to collect facts from %s: %w", machine.Hostname, err)
    }

    return facts, nil
}
```

#### 1.2 Update CLI Command
**File**: `cmd/facts.go`

**Current Problem**:
```go
// Collect facts using the integration
facts, err := factsManager.CollectFacts(ctx, machine.Hostname)
```

**Solution**:
```go
// Collect facts using the integration with actual machine object
facts, err := factsManager.CollectFacts(ctx, &machine)
```

### Phase 2: Implement Proper SSH Connection (Priority: Critical)

#### 2.1 Fix SSH Command Execution
**File**: `internal/facts/collector.go`

**Current Problem**:
The `executeSSHCommand` method exists but may not be working correctly with the SSH manager.

**Solution**:
```go
// executeSSHCommand executes a command via SSH on the target machine
func (c *SystemFactCollector) executeSSHCommand(machine *spookytypes.Machine, command string) (string, error) {
    ctx := context.Background()

    // Create connection request with actual machine configuration
    connectionRequest := &spookytypes.ConnectionRequest{
        Host: machine.Hostname,
        Port: machine.Port,
        User: machine.User,
        // Add authentication from machine configuration
        Authentication: machine.Authentication,
    }

    // Establish connection
    connectionResult, err := c.sshManager.Connect(ctx, connectionRequest)
    if err != nil {
        return "", fmt.Errorf("failed to establish SSH connection to %s: %w", machine.Hostname, err)
    }

    // Create session
    session, err := c.sshManager.CreateSession(ctx, connectionResult.Connection)
    if err != nil {
        return "", fmt.Errorf("failed to create SSH session on %s: %w", machine.Hostname, err)
    }

    // Create SSH command
    sshCommand := &spookytypes.SSHCommand{
        Command: command,
        Timeout: 30 * time.Second,
    }

    // Run command
    commandResult, err := c.sshManager.RunCommand(ctx, session, sshCommand)
    if err != nil {
        return "", fmt.Errorf("failed to run SSH command on %s: %w", machine.Hostname, err)
    }

    if !commandResult.Success {
        return "", fmt.Errorf("SSH command failed on %s with exit code %d: %s", 
            machine.Hostname, commandResult.ExitCode, commandResult.Stderr)
    }

    return commandResult.Stdout, nil
}
```

#### 2.2 Add Authentication Support
**File**: `internal/facts/collector.go`

**Problem**: Need to handle different authentication methods from machine configuration.

**Solution**:
```go
// getSSHConfig creates SSH configuration from machine settings
func (c *SystemFactCollector) getSSHConfig(machine *spookytypes.Machine) (*spookytypes.ClientConfig, error) {
    config := &spookytypes.ClientConfig{
        DefaultPort:        machine.Port,
        DefaultTimeout:     30 * time.Second,
        MaxConnections:     10,
        MaxRetryAttempts:   3,
        RetryDelay:         5 * time.Second,
        KeepaliveInterval:  60 * time.Second,
        KeepaliveCount:     3,
        EnableCompression:  false,
        EnableKeepalive:    true,
        StrictHostKeyCheck: true,
    }

    // Add authentication configuration
    if machine.Authentication != nil {
        switch machine.Authentication.Method {
        case "ssh_key":
            config.PrivateKeyPath = machine.Authentication.KeyPath
        case "password":
            config.Password = machine.Authentication.Password
        case "agent":
            config.UseAgent = true
        }
    }

    return config, nil
}
```

### Phase 3: Implement Remote File Reading (Priority: High)

#### 3.1 Fix Collector Facts Collection
**File**: `internal/facts/collector.go`

**Current Problem**:
```go
// Read collector facts file
factsContent, err := c.executeSSHCommand(machine, "cat /etc/spooky/facts.hcl")
```

**Solution**: Add proper error handling and file existence checks:
```go
// collectCollectorFacts collects facts from spooky-collector binary
func (c *SystemFactCollector) collectCollectorFacts(machine *spookytypes.Machine) (*spookytypesfacts.CollectorFacts, error) {
    // Check if collector facts file exists
    checkCmd := "test -f /etc/spooky/facts.hcl && echo 'exists' || echo 'not_found'"
    result, err := c.executeSSHCommand(machine, checkCmd)
    if err != nil {
        return nil, fmt.Errorf("failed to check collector facts file on %s: %w", machine.Hostname, err)
    }

    if strings.TrimSpace(result) == "not_found" {
        return nil, fmt.Errorf("collector facts file not found on %s: /etc/spooky/facts.hcl", machine.Hostname)
    }

    // Read collector facts file
    factsContent, err := c.executeSSHCommand(machine, "cat /etc/spooky/facts.hcl")
    if err != nil {
        return nil, fmt.Errorf("failed to read collector facts from %s: %w", machine.Hostname, err)
    }

    // Parse HCL content using the HCL parser
    parser := NewHCLParser()
    collectorFacts, err := parser.ParseCollectorFacts(factsContent)
    if err != nil {
        return nil, fmt.Errorf("failed to parse collector facts from %s: %w", machine.Hostname, err)
    }

    return collectorFacts, nil
}
```

#### 3.2 Fix Custom Facts Collection
**File**: `internal/facts/collector.go`

**Current Problem**:
```go
// Read custom facts file
factsContent, err := c.executeSSHCommand(machine, "cat /etc/spooky/custom.hcl")
```

**Solution**: Make custom facts optional and add proper error handling:
```go
// collectCustomFacts collects custom facts from HCL files
func (c *SystemFactCollector) collectCustomFacts(machine *spookytypes.Machine) (map[string]interface{}, error) {
    // Check if custom facts file exists (optional)
    checkCmd := "test -f /etc/spooky/custom.hcl && echo 'exists' || echo 'not_found'"
    result, err := c.executeSSHCommand(machine, checkCmd)
    if err != nil {
        return nil, fmt.Errorf("failed to check custom facts file on %s: %w", machine.Hostname, err)
    }

    if strings.TrimSpace(result) == "not_found" {
        // Custom facts are optional, return empty map
        return make(map[string]interface{}), nil
    }

    // Read custom facts file
    factsContent, err := c.executeSSHCommand(machine, "cat /etc/spooky/custom.hcl")
    if err != nil {
        return nil, fmt.Errorf("failed to read custom facts from %s: %w", machine.Hostname, err)
    }

    // Parse HCL content using the HCL parser
    parser := NewHCLParser()
    customFacts, err := parser.ParseCustomFacts(factsContent)
    if err != nil {
        return nil, fmt.Errorf("failed to parse custom facts from %s: %w", machine.Hostname, err)
    }

    return customFacts, nil
}
```

### Phase 4: Add Parallel Collection Support (Priority: Medium)

#### 4.1 Implement Parallel Fact Collection
**File**: `cmd/facts.go`

**Current Problem**: Sequential collection is slow for multiple machines.

**Solution**:
```go
// collectFactsParallel collects facts from multiple machines in parallel
func collectFactsParallel(ctx context.Context, machines []spookytypes.Machine, factsManager spookyinterfaces.FactsIntegration, parallel int) (map[string]interface{}, error) {
    results := make(map[string]interface{})
    var mu sync.Mutex
    var wg sync.WaitGroup

    // Create semaphore for parallel execution
    semaphore := make(chan struct{}, parallel)

    for _, machine := range machines {
        wg.Add(1)
        go func(m spookytypes.Machine) {
            defer wg.Done()
            semaphore <- struct{}{} // Acquire
            defer func() { <-semaphore }() // Release

            facts, err := factsManager.CollectFacts(ctx, &m)
            if err != nil {
                // Log error but continue with other machines
                log.Printf("Failed to collect facts from %s: %v", m.Hostname, err)
                return
            }

            mu.Lock()
            results[m.Hostname] = facts
            mu.Unlock()
        }(machine)
    }

    wg.Wait()
    return results, nil
}
```

### Phase 5: Add Comprehensive Error Handling (Priority: High)

#### 5.1 Implement Error Classification
**File**: `internal/facts/collector.go`

**Solution**:
```go
// FactCollectionError represents fact collection errors
type FactCollectionError struct {
    Machine     string
    ErrorType   string
    Message     string
    Recoverable bool
}

func (e *FactCollectionError) Error() string {
    return fmt.Sprintf("fact collection error on %s (%s): %s", e.Machine, e.ErrorType, e.Message)
}

// classifyError classifies fact collection errors
func classifyError(machine string, err error) *FactCollectionError {
    if err == nil {
        return nil
    }

    errorType := "unknown"
    recoverable := false

    switch {
    case strings.Contains(err.Error(), "connection refused"):
        errorType = "ssh_connection"
        recoverable = false
    case strings.Contains(err.Error(), "authentication failed"):
        errorType = "ssh_auth"
        recoverable = false
    case strings.Contains(err.Error(), "permission denied"):
        errorType = "permission"
        recoverable = true
    case strings.Contains(err.Error(), "file not found"):
        errorType = "file_not_found"
        recoverable = true
    case strings.Contains(err.Error(), "timeout"):
        errorType = "timeout"
        recoverable = true
    }

    return &FactCollectionError{
        Machine:     machine,
        ErrorType:   errorType,
        Message:     err.Error(),
        Recoverable: recoverable,
    }
}
```

### Phase 6: Add Testing and Validation (Priority: High)

#### 6.1 Create Integration Tests
**File**: `internal/facts/collector_test.go`

**Solution**:
```go
func TestSSHFactCollection(t *testing.T) {
    // Test with mock SSH server
    // Test with real SSH server
    // Test error conditions
    // Test parallel collection
}
```

#### 6.2 Add Validation Tests
**File**: `cmd/facts_test.go`

**Solution**:
```go
func TestFactsExportWithSSH(t *testing.T) {
    // Test CLI command with SSH
    // Test machine filtering
    // Test parallel execution
    // Test error handling
}
```

## Implementation Steps

### Step 1: Fix Facts Integration (1-2 hours)
1. Update `internal/facts/integration.go` to accept `*spookytypes.Machine`
2. Update `cmd/facts.go` to pass machine objects
3. Test basic integration

### Step 2: Fix SSH Connection (2-3 hours)
1. Update `executeSSHCommand` to use proper machine configuration
2. Add authentication support
3. Test SSH connections

### Step 3: Implement Remote File Reading (2-3 hours)
1. Fix collector facts collection with proper error handling
2. Fix custom facts collection (make optional)
3. Test remote file reading

### Step 4: Add Parallel Support (1-2 hours)
1. Implement parallel fact collection
2. Add proper synchronization
3. Test parallel execution

### Step 5: Add Error Handling (1-2 hours)
1. Implement error classification
2. Add recoverable vs non-recoverable error handling
3. Test error scenarios

### Step 6: Add Testing (2-3 hours)
1. Create integration tests
2. Create validation tests
3. Test with real SSH servers

## Success Criteria

### Functional Requirements
- ✅ SSH connections to remote machines work correctly
- ✅ System facts collected via SSH commands
- ✅ Collector facts read from remote `/etc/spooky/facts.hcl`
- ✅ Custom facts read from remote `/etc/spooky/custom.hcl` (optional)
- ✅ Parallel collection works for multiple machines
- ✅ Proper error handling and classification
- ✅ Authentication support for different methods

### Performance Requirements
- ✅ Collection time < 30 seconds per machine
- ✅ Parallel collection supports 4-8 concurrent connections
- ✅ Memory usage remains reasonable for large inventories

### Quality Requirements
- ✅ All tests pass
- ✅ No compilation errors
- ✅ Proper error messages for troubleshooting
- ✅ Logging for debugging and monitoring

## Risk Assessment

### High Risk
- **SSH Authentication**: Different authentication methods may not work
- **Network Issues**: SSH connections may fail in various network conditions
- **Performance**: Parallel collection may overwhelm target machines

### Medium Risk
- **File Permissions**: Remote files may have different permissions
- **HCL Parsing**: Remote HCL files may have different formats
- **Error Handling**: Complex error scenarios may not be handled properly

### Low Risk
- **CLI Integration**: Changes to CLI are straightforward
- **Type System**: Type changes are minimal and backward compatible

## Dependencies

### Internal Dependencies
- SSH manager must be working correctly
- Machine configuration must be properly loaded
- HCL parser must handle remote file formats

### External Dependencies
- Target machines must have SSH access
- Target machines must have proper permissions
- Network connectivity must be stable

## Timeline

**Total Estimated Time**: 9-15 hours

- **Phase 1**: 1-2 hours (Critical)
- **Phase 2**: 2-3 hours (Critical)
- **Phase 3**: 2-3 hours (High)
- **Phase 4**: 1-2 hours (Medium)
- **Phase 5**: 1-2 hours (High)
- **Phase 6**: 2-3 hours (High)

**Priority Order**:
1. Fix Facts Integration (Critical)
2. Fix SSH Connection (Critical)
3. Implement Remote File Reading (High)
4. Add Error Handling (High)
5. Add Testing (High)
6. Add Parallel Support (Medium)

This plan addresses the core issue where the facts export command is essentially broken for multi-machine environments and provides a comprehensive solution to implement proper SSH-based fact collection.
