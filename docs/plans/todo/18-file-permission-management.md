# Implementation Plan: File Permission Management

## Overview
Implement proper file permission and ownership management for template deployment and file operations, replacing placeholder implementations with real syscall-based functionality.

## Task Details
- **Task ID**: 7.2
- **Priority**: Medium
- **Files**: 
  - `internal/actions/acting/actor.go`
  - `internal/ssh/acting/manager.go`
- **Functions**: File ownership, permissions, attributes

## Current State Analysis

### Existing Patterns
1. **Template Deployment**: Basic file deployment exists
2. **File Operations**: Basic file operations implemented
3. **SSH Operations**: SSH-based file operations available
4. **Error Handling**: Consistent error wrapping

### Current Placeholder Code
```go
// internal/actions/acting/actor.go lines 586, 590
// TODO: Implement actual owner setting using syscall.Chown
// For now, we'll just log the intention
```

## Implementation Requirements

### Interface Compliance
The file permission system must:
1. **Set file ownership** using syscall.Chown
2. **Set file permissions** using syscall.Chmod
3. **Handle user/group resolution** from strings to UIDs/GIDs
4. **Support remote file operations** via SSH
5. **Provide permission validation** and error handling
6. **Support recursive operations** for directories
7. **Handle different platforms** (Linux, Unix, Windows)

### Required Dependencies
- syscall package for file operations
- user/group resolution system
- SSH file operation system
- Platform abstraction layer

## Detailed Implementation Plan

### Step 1: Implement File Ownership Management

#### 1.1 User/Group Resolution
```go
// internal/actions/acting/permissions.go
package acting

import (
    "fmt"
    "os/user"
    "strconv"
    "syscall"
)

// UserGroupResolver resolves user and group names to UIDs and GIDs
type UserGroupResolver struct{}

// NewUserGroupResolver creates a new user/group resolver
func NewUserGroupResolver() *UserGroupResolver {
    return &UserGroupResolver{}
}

// ResolveUser resolves a username to UID
func (r *UserGroupResolver) ResolveUser(username string) (int, error) {
    if username == "" {
        return -1, fmt.Errorf("username cannot be empty")
    }
    
    // Try to parse as numeric UID first
    if uid, err := strconv.Atoi(username); err == nil {
        return uid, nil
    }
    
    // Look up user by name
    u, err := user.Lookup(username)
    if err != nil {
        return -1, fmt.Errorf("failed to resolve user %s: %w", username, err)
    }
    
    uid, err := strconv.Atoi(u.Uid)
    if err != nil {
        return -1, fmt.Errorf("invalid UID for user %s: %w", username, err)
    }
    
    return uid, nil
}

// ResolveGroup resolves a group name to GID
func (r *UserGroupResolver) ResolveGroup(groupname string) (int, error) {
    if groupname == "" {
        return -1, fmt.Errorf("group name cannot be empty")
    }
    
    // Try to parse as numeric GID first
    if gid, err := strconv.Atoi(groupname); err == nil {
        return gid, nil
    }
    
    // Look up group by name
    g, err := user.LookupGroup(groupname)
    if err != nil {
        return -1, fmt.Errorf("failed to resolve group %s: %w", groupname, err)
    }
    
    gid, err := strconv.Atoi(g.Gid)
    if err != nil {
        return -1, fmt.Errorf("invalid GID for group %s: %w", groupname, err)
    }
    
    return gid, nil
}

// ResolveUserGroup resolves both user and group
func (r *UserGroupResolver) ResolveUserGroup(username, groupname string) (int, int, error) {
    uid, err := r.ResolveUser(username)
    if err != nil {
        return -1, -1, fmt.Errorf("failed to resolve user: %w", err)
    }
    
    gid, err := r.ResolveGroup(groupname)
    if err != nil {
        return -1, -1, fmt.Errorf("failed to resolve group: %w", err)
    }
    
    return uid, gid, nil
}
```

#### 1.2 File Ownership Setting
```go
// SetFileOwner sets file owner using syscall.Chown
func (a *actorImpl) setFileOwner(destination, owner string) error {
    a.logger.Debug("Setting file owner",
        spookylogging.String("destination", destination),
        spookylogging.String("owner", owner))

    // Parse owner string (format: "user:group" or "user")
    var username, groupname string
    if colonIndex := strings.Index(owner, ":"); colonIndex != -1 {
        username = owner[:colonIndex]
        groupname = owner[colonIndex+1:]
    } else {
        username = owner
        // Use user's primary group
        groupname = ""
    }

    // Resolve user and group
    resolver := NewUserGroupResolver()
    uid, gid, err := resolver.ResolveUserGroup(username, groupname)
    if err != nil {
        return fmt.Errorf("failed to resolve owner %s: %w", owner, err)
    }

    // Set file ownership
    err = syscall.Chown(destination, uid, gid)
    if err != nil {
        return fmt.Errorf("failed to set file owner %s:%s on %s: %w", username, groupname, destination, err)
    }

    a.logger.Debug("File owner set successfully",
        spookylogging.String("destination", destination),
        spookylogging.String("owner", owner),
        spookylogging.Int("uid", uid),
        spookylogging.Int("gid", gid))

    return nil
}
```

### Step 2: Implement File Permission Management

#### 2.1 Permission Parsing
```go
// PermissionParser parses permission strings
type PermissionParser struct{}

// NewPermissionParser creates a new permission parser
func NewPermissionParser() *PermissionParser {
    return &PermissionParser{}
}

// ParsePermissions parses permission string to mode
func (p *PermissionParser) ParsePermissions(permissions string) (os.FileMode, error) {
    if permissions == "" {
        return 0, fmt.Errorf("permissions cannot be empty")
    }

    // Handle numeric permissions (e.g., "644", "0755")
    if mode, err := strconv.ParseUint(permissions, 8, 32); err == nil {
        return os.FileMode(mode), nil
    }

    // Handle symbolic permissions (e.g., "rw-r--r--", "u+rwx")
    return p.parseSymbolicPermissions(permissions)
}

// parseSymbolicPermissions parses symbolic permission strings
func (p *PermissionParser) parseSymbolicPermissions(permissions string) (os.FileMode, error) {
    // This is a simplified implementation
    // In practice, this would handle complex symbolic permissions
    
    var mode os.FileMode
    
    // Parse basic patterns
    switch permissions {
    case "rw-r--r--", "644":
        mode = 0644
    case "rwxr-xr-x", "755":
        mode = 0755
    case "rw-------", "600":
        mode = 0600
    case "rwx------", "700":
        mode = 0700
    default:
        return 0, fmt.Errorf("unsupported permission format: %s", permissions)
    }
    
    return mode, nil
}
```

#### 2.2 File Permission Setting
```go
// setFilePermissions sets file permissions
func (a *actorImpl) setFilePermissions(destination, permissions string) error {
    a.logger.Debug("Setting file permissions",
        spookylogging.String("destination", destination),
        spookylogging.String("permissions", permissions))

    // Parse permissions
    parser := NewPermissionParser()
    mode, err := parser.ParsePermissions(permissions)
    if err != nil {
        return fmt.Errorf("failed to parse permissions %s: %w", permissions, err)
    }

    // Set file permissions
    err = os.Chmod(destination, mode)
    if err != nil {
        return fmt.Errorf("failed to set permissions %s on %s: %w", permissions, destination, err)
    }

    a.logger.Debug("File permissions set successfully",
        spookylogging.String("destination", destination),
        spookylogging.String("permissions", permissions),
        spookylogging.String("mode", mode.String()))

    return nil
}
```

### Step 3: Implement SSH File Operations

#### 3.1 SSH File Permission Manager
```go
// internal/ssh/acting/permissions.go
package acting

import (
    "fmt"
    "strings"
    "time"
)

// SSHFilePermissionManager manages file permissions over SSH
type SSHFilePermissionManager struct {
    sshManager SSHManager
    logger     logging.Logger
}

// NewSSHFilePermissionManager creates a new SSH file permission manager
func NewSSHFilePermissionManager(sshManager SSHManager, logger logging.Logger) *SSHFilePermissionManager {
    return &SSHFilePermissionManager{
        sshManager: sshManager,
        logger:     logger,
    }
}

// SetFileOwner sets file owner over SSH
func (m *SSHFilePermissionManager) SetFileOwner(connection SSHConnection, destination, owner string) error {
    m.logger.Debug("Setting file owner via SSH",
        spookylogging.String("destination", destination),
        spookylogging.String("owner", owner))

    // Parse owner string
    var username, groupname string
    if colonIndex := strings.Index(owner, ":"); colonIndex != -1 {
        username = owner[:colonIndex]
        groupname = owner[colonIndex+1:]
    } else {
        username = owner
        groupname = ""
    }

    // Build chown command
    var chownCmd string
    if groupname != "" {
        chownCmd = fmt.Sprintf("chown %s:%s %s", username, groupname, destination)
    } else {
        chownCmd = fmt.Sprintf("chown %s %s", username, destination)
    }

    // Execute chown command
    result, err := m.sshManager.ExecuteCommand(connection, chownCmd)
    if err != nil {
        return fmt.Errorf("failed to execute chown command: %w", err)
    }

    if result.ExitCode != 0 {
        return fmt.Errorf("chown command failed: %s", result.Stderr)
    }

    m.logger.Debug("File owner set successfully via SSH",
        spookylogging.String("destination", destination),
        spookylogging.String("owner", owner))

    return nil
}

// SetFilePermissions sets file permissions over SSH
func (m *SSHFilePermissionManager) SetFilePermissions(connection SSHConnection, destination, permissions string) error {
    m.logger.Debug("Setting file permissions via SSH",
        spookylogging.String("destination", destination),
        spookylogging.String("permissions", permissions))

    // Build chmod command
    chmodCmd := fmt.Sprintf("chmod %s %s", permissions, destination)

    // Execute chmod command
    result, err := m.sshManager.ExecuteCommand(connection, chmodCmd)
    if err != nil {
        return fmt.Errorf("failed to execute chmod command: %w", err)
    }

    if result.ExitCode != 0 {
        return fmt.Errorf("chmod command failed: %s", result.Stderr)
    }

    m.logger.Debug("File permissions set successfully via SSH",
        spookylogging.String("destination", destination),
        spookylogging.String("permissions", permissions))

    return nil
}

// SetFileAttributes sets multiple file attributes
func (m *SSHFilePermissionManager) SetFileAttributes(connection SSHConnection, destination string, attributes *FileAttributes) error {
    m.logger.Debug("Setting file attributes via SSH",
        spookylogging.String("destination", destination))

    var errors []error

    // Set owner if specified
    if attributes.Owner != "" {
        if err := m.SetFileOwner(connection, destination, attributes.Owner); err != nil {
            errors = append(errors, fmt.Errorf("failed to set owner: %w", err))
        }
    }

    // Set permissions if specified
    if attributes.Permissions != "" {
        if err := m.SetFilePermissions(connection, destination, attributes.Permissions); err != nil {
            errors = append(errors, fmt.Errorf("failed to set permissions: %w", err))
        }
    }

    // Set timestamp if specified
    if !attributes.ModTime.IsZero() {
        touchCmd := fmt.Sprintf("touch -m -t %s %s", 
            attributes.ModTime.Format("200601021504.05"), destination)
        
        result, err := m.sshManager.ExecuteCommand(connection, touchCmd)
        if err != nil {
            errors = append(errors, fmt.Errorf("failed to set timestamp: %w", err))
        } else if result.ExitCode != 0 {
            errors = append(errors, fmt.Errorf("touch command failed: %s", result.Stderr))
        }
    }

    // Return combined errors if any
    if len(errors) > 0 {
        return fmt.Errorf("failed to set file attributes: %v", errors)
    }

    m.logger.Debug("File attributes set successfully via SSH",
        spookylogging.String("destination", destination))

    return nil
}
```

### Step 4: Implement File Attributes Structure

#### 4.1 File Attributes Definition
```go
// FileAttributes represents file attributes
type FileAttributes struct {
    Owner       string    `json:"owner,omitempty"`
    Group       string    `json:"group,omitempty"`
    Permissions string    `json:"permissions,omitempty"`
    ModTime     time.Time `json:"mod_time,omitempty"`
    Recursive   bool      `json:"recursive,omitempty"`
}

// NewFileAttributes creates new file attributes
func NewFileAttributes() *FileAttributes {
    return &FileAttributes{}
}

// SetOwner sets the owner
func (fa *FileAttributes) SetOwner(owner string) *FileAttributes {
    fa.Owner = owner
    return fa
}

// SetGroup sets the group
func (fa *FileAttributes) SetGroup(group string) *FileAttributes {
    fa.Group = group
    return fa
}

// SetPermissions sets the permissions
func (fa *FileAttributes) SetPermissions(permissions string) *FileAttributes {
    fa.Permissions = permissions
    return fa
}

// SetModTime sets the modification time
func (fa *FileAttributes) SetModTime(modTime time.Time) *FileAttributes {
    fa.ModTime = modTime
    return fa
}

// SetRecursive sets recursive flag
func (fa *FileAttributes) SetRecursive(recursive bool) *FileAttributes {
    fa.Recursive = recursive
    return fa
}
```

### Step 5: Implement Recursive Operations

#### 5.1 Recursive File Operations
```go
// setFileAttributesRecursive sets file attributes recursively
func (a *actorImpl) setFileAttributesRecursive(destination string, attributes *FileAttributes) error {
    a.logger.Debug("Setting file attributes recursively",
        spookylogging.String("destination", destination))

    // Walk directory tree
    err := filepath.Walk(destination, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // Skip if it's the root directory and not recursive
        if path == destination && !attributes.Recursive {
            return nil
        }

        // Set attributes for each file/directory
        if err := a.setFileAttributes(path, attributes); err != nil {
            a.logger.Warn("Failed to set attributes for file",
                spookylogging.String("path", path),
                spookylogging.Error(err))
            // Continue with other files
        }

        return nil
    })

    if err != nil {
        return fmt.Errorf("failed to set attributes recursively: %w", err)
    }

    a.logger.Debug("File attributes set recursively",
        spookylogging.String("destination", destination))

    return nil
}

// setFileAttributes sets attributes for a single file
func (a *actorImpl) setFileAttributes(destination string, attributes *FileAttributes) error {
    // Set owner if specified
    if attributes.Owner != "" {
        if err := a.setFileOwner(destination, attributes.Owner); err != nil {
            return fmt.Errorf("failed to set owner: %w", err)
        }
    }

    // Set permissions if specified
    if attributes.Permissions != "" {
        if err := a.setFilePermissions(destination, attributes.Permissions); err != nil {
            return fmt.Errorf("failed to set permissions: %w", err)
        }
    }

    // Set modification time if specified
    if !attributes.ModTime.IsZero() {
        if err := os.Chtimes(destination, attributes.ModTime, attributes.ModTime); err != nil {
            return fmt.Errorf("failed to set modification time: %w", err)
        }
    }

    return nil
}
```

### Step 6: Update Template Deployment

#### 6.1 Enhanced Template Deployment
```go
// deployTemplateWithAttributes deploys template with file attributes
func (a *actorImpl) deployTemplateWithAttributes(content, destination string, attributes *FileAttributes) error {
    // Deploy template content
    if err := a.deployTemplate(content, destination); err != nil {
        return fmt.Errorf("failed to deploy template: %w", err)
    }

    // Set file attributes
    if attributes != nil {
        if attributes.Recursive {
            if err := a.setFileAttributesRecursive(destination, attributes); err != nil {
                return fmt.Errorf("failed to set file attributes: %w", err)
            }
        } else {
            if err := a.setFileAttributes(destination, attributes); err != nil {
                return fmt.Errorf("failed to set file attributes: %w", err)
            }
        }
    }

    return nil
}
```

## Configuration Options

### Supported Options
- **DefaultOwner**: Default file owner
- **DefaultPermissions**: Default file permissions
- **RecursiveByDefault**: Enable recursive operations by default
- **PreserveTimestamps**: Preserve original timestamps

## Dependencies

### Internal Dependencies
- `spooky/internal/ssh`
- `spooky/internal/logging`

### External Dependencies
- `syscall` (standard library)
- `os/user` (standard library)
- `os` (standard library)
- `path/filepath` (standard library)
- `strconv` (standard library)
- `strings` (standard library)
- `time` (standard library)

## Implementation Order

1. Implement user/group resolution
2. Add file ownership setting
3. Implement permission parsing
4. Add file permission setting
5. Create SSH file operations
6. Implement file attributes structure
7. Add recursive operations
8. Update template deployment
9. Add comprehensive tests
10. Performance optimization
11. Documentation and cleanup
