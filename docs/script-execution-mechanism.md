# Script Execution Mechanism in Spooky

Based on analysis of the codebase, here's how script execution works in Spooky:

## Current Implementation (Local Only)

### Temporary Script Creation Process

**Location**: The temporary script is created on the **local machine** (where Spooky is running), not on the target machines.

**Mechanism**:
1. **`createTemporaryScript()`** function creates a temporary file using `os.CreateTemp("", "spooky-script-*.sh")`
2. **Writes script content** to the temporary file
3. **Sets executable permissions** (`chmod 0755`) on the file
4. **Returns the file path** for execution

**Code Flow**:
```go
func (a *actorImpl) createTemporaryScript(scriptContent string) (string, error) {
    // Create temporary file on local machine
    tmpFile, err := os.CreateTemp("", "spooky-script-*.sh")
    if err != nil {
        return "", fmt.Errorf("failed to create temporary file: %w", err)
    }

    // Write script content
    if _, err := tmpFile.WriteString(scriptContent); err != nil {
        tmpFile.Close()
        os.Remove(tmpFile.Name())
        return "", fmt.Errorf("failed to write script content: %w", err)
    }

    // Close file
    if err := tmpFile.Close(); err != nil {
        os.Remove(tmpFile.Name())
        return "", fmt.Errorf("failed to close temporary file: %w", err)
    }

    // Make script executable
    if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
        os.Remove(tmpFile.Name())
        return "", fmt.Errorf("failed to make script executable: %w", err)
    }

    return tmpFile.Name(), nil
}
```

### Execution Process

**Current State**: Only local execution is implemented
- Scripts are executed locally using `exec.CommandContext(ctx, scriptFile)`
- Remote machine execution shows "TODO: Implement SSH acting for remote machines"

**Local Execution Flow**:
1. Create temporary script file locally
2. Execute script directly using `exec.CommandContext(ctx, scriptFile)`
3. Capture stdout/stderr
4. Clean up temporary file with `defer os.Remove(scriptFile)`

## Planned Remote Execution (SSH)

### Future Implementation

**Documentation shows planned SSH mechanism**:

1. **Script Upload**: Upload script content to remote machine via SCP
2. **Remote Path**: Generate unique remote path like `/tmp/spooky_script_<timestamp>_<name>`
3. **Permission Setting**: Set execute permissions on remote machine
4. **Execution**: Execute script on remote machine via SSH
5. **Cleanup**: Remove remote script file after execution

**Planned Code Structure**:
```go
func (m *Manager) uploadScript(connection sshTypes.Connection, scriptPath string, context *types.ActionContext) (string, error) {
    // Read script content
    content, err := os.ReadFile(scriptPath)
    if err != nil {
        return "", fmt.Errorf("failed to read script: %w", err)
    }

    // Generate remote path
    scriptName := filepath.Base(scriptPath)
    remotePath := fmt.Sprintf("/tmp/spooky_script_%d_%s", time.Now().Unix(), scriptName)

    // Upload file via SCP
    err = m.uploadFileViaSCP(connection, content, remotePath)
    if err != nil {
        return "", fmt.Errorf("failed to upload script: %w", err)
    }

    // Set execute permissions
    chmodCmd := fmt.Sprintf("chmod +x %s", remotePath)
    chmodResult, err := m.sshManager.ExecuteCommand(connection, chmodCmd)
    if err != nil {
        return "", fmt.Errorf("failed to set script permissions: %w", err)
    }

    return remotePath, nil
}
```

## Key Points

### Current Limitations
- **Local-only execution**: Scripts only run on the machine where Spooky is executed
- **No remote support**: SSH implementation is planned but not yet implemented
- **Temporary file location**: Created in system temp directory (e.g., `/tmp/` on Linux)

### Security Considerations
- **Temporary file cleanup**: Automatic cleanup with `defer os.Remove(scriptFile)`
- **Unique naming**: Uses timestamp and random suffix to avoid conflicts
- **Permission control**: Sets appropriate execute permissions (0755)

### Future Enhancements
- **SSH integration**: Full remote execution support
- **Connection pooling**: Efficient SSH connection management
- **Parallel execution**: Execute scripts on multiple machines simultaneously
- **Error handling**: Robust error handling for network failures
