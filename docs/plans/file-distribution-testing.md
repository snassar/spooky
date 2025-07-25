# File Distribution Testing Plan

## Overview

This document outlines the testing plan for spooky's enhanced file distribution functionality. The goal is to verify that template files can be deployed to target servers reliably, with proper idempotency, backup support, and error handling.

## Background

File distribution is a critical component of Issue #52 (server-side template evaluation). When deploying templates to remote servers, spooky needs to:

1. **Deploy template files** to target servers via SSH
2. **Ensure idempotency** - don't overwrite files if content hasn't changed
3. **Create backups** before overwriting existing files (if requested)
4. **Set proper permissions** and ownership
5. **Validate file integrity** after deployment
6. **Handle errors gracefully** with proper logging and rollback

## Enhanced Features

### 1. Idempotency
- Check if remote file exists before deployment
- Compare local and remote file content
- Skip deployment if content is identical
- Log when files are skipped due to no changes

### 2. Backup Support
- Create backups of existing files before overwriting
- Use timestamped backup filenames
- Abort deployment if backup creation fails
- Log backup creation success/failure

### 3. File Validation
- Verify file was written correctly after deployment
- Check file size and basic integrity
- Abort deployment if validation fails
- Log validation results

### 4. Error Handling
- Graceful handling of SSH connection failures
- Proper error messages for file operations
- Rollback capabilities for failed deployments
- Comprehensive logging for debugging

## Test Setup

### 1. Test Project Structure

```
examples/projects/file-distribution-testing/
├── project.hcl
├── inventory.hcl
├── actions.hcl
├── templates/
│   ├── nginx-config.tmpl
│   ├── php-config.tmpl
│   └── test-template.tmpl
└── custom-facts.json
```

### 2. Inventory Configuration

```hcl
inventory {
  machine "test-server" {
    host     = "192.168.1.100"  # or localhost for local testing
    port     = 22
    user     = "debian"
    password = "your-password"   # or use keyfile
    tags = {
      environment = "testing"
      role = "web"
    }
  }
}
```

### 3. Template Actions Configuration

```hcl
actions {
  action "deploy-nginx-config" {
    description = "Deploy nginx configuration template"
    type = "template_deploy"
    template {
      source = "templates/nginx-config.tmpl"
      destination = "/tmp/nginx.conf.tmpl"
      backup = true
      permissions = "644"
      owner = "root"
      group = "root"
    }
    tags = ["role=web"]
  }

  action "deploy-php-config" {
    description = "Deploy PHP configuration template"
    type = "template_deploy"
    template {
      source = "templates/php-config.tmpl"
      destination = "/tmp/php.ini.tmpl"
      backup = true
      permissions = "644"
    }
    tags = ["role=web"]
  }
}
```

## Test Execution

### Step 1: Basic Template Deployment

```bash
# Deploy template to test server
./build/spooky action run --project examples/projects/file-distribution-testing deploy-nginx-config
```

**Expected Output:**
- Template file should be deployed to `/tmp/nginx.conf.tmpl`
- Backup should be created (if file existed)
- File permissions should be set to 644
- Success message should be logged

### Step 2: Idempotency Test

```bash
# Deploy the same template again
./build/spooky action run --project examples/projects/file-distribution-testing deploy-nginx-config
```

**Expected Output:**
- File content comparison should detect no changes
- Deployment should be skipped
- Log message: "File content unchanged, skipping deployment"
- No backup should be created

### Step 3: Content Change Test

1. Modify the template file locally
2. Deploy again:

```bash
./build/spooky action run --project examples/projects/file-distribution-testing deploy-nginx-config
```

**Expected Output:**
- Content comparison should detect changes
- New backup should be created
- Updated file should be deployed
- Success message should be logged

### Step 4: Backup Functionality Test

```bash
# Deploy with backup enabled
./build/spooky action run --project examples/projects/file-distribution-testing deploy-php-config
```

**Expected Output:**
- Backup file should be created with timestamp
- Original file should be preserved
- New file should be deployed
- Backup creation should be logged

### Step 5: Permission and Ownership Test

```bash
# Deploy with specific permissions and ownership
./build/spooky action run --project examples/projects/file-distribution-testing deploy-nginx-config
```

**Expected Output:**
- File permissions should be set to 644
- File ownership should be set to root:root
- Permission/ownership changes should be logged

### Step 6: Error Handling Tests

#### SSH Connection Failure

```bash
# Test with unreachable server
./build/spooky action run --project examples/projects/file-distribution-testing deploy-nginx-config
```

**Expected Output:**
- SSH connection error should be logged
- Deployment should be skipped for that machine
- Other machines should continue processing

#### File Permission Denied

```bash
# Test deployment to protected directory
./build/spooky action run --project examples/projects/file-distribution-testing deploy-nginx-config
```

**Expected Output:**
- Permission denied error should be logged
- Deployment should fail gracefully
- Clear error message should be provided

#### Backup Creation Failure

```bash
# Test backup creation in read-only directory
./build/spooky action run --project examples/projects/file-distribution-testing deploy-nginx-config
```

**Expected Output:**
- Backup creation failure should be logged
- Deployment should be aborted
- Clear error message should be provided

### Step 7: File Validation Test

```bash
# Deploy template and verify integrity
./build/spooky action run --project examples/projects/file-distribution-testing deploy-nginx-config
```

**Expected Output:**
- File should be written successfully
- File validation should pass
- File size and content should be verified
- Validation success should be logged

## Validation Criteria

### ✅ Success Criteria

1. **Idempotency**: Files are not overwritten if content hasn't changed
2. **Backup Support**: Backups are created when requested and file exists
3. **File Validation**: Deployed files pass integrity checks
4. **Permission Setting**: File permissions and ownership are set correctly
5. **Error Handling**: Failures are handled gracefully with proper logging
6. **Performance**: Deployment completes within reasonable time
7. **Logging**: Comprehensive logs for debugging and monitoring

### ❌ Failure Indicators

1. **Content Loss**: Files are overwritten without backup when backup is requested
2. **Permission Issues**: Files have incorrect permissions or ownership
3. **Validation Failures**: Deployed files fail integrity checks
4. **Error Propagation**: SSH or file operation errors crash the deployment
5. **Performance Issues**: Deployment takes too long or times out
6. **Insufficient Logging**: Lack of detailed logs for troubleshooting

## Test Data Collection

### Log Analysis

Check logs for:
- SSH connection success/failure
- File existence checks
- Content comparison results
- Backup creation success/failure
- File validation results
- Permission/ownership changes
- Deployment completion status

### File System Verification

Verify on target server:
- Template files are present in correct locations
- Backup files exist with proper timestamps
- File permissions and ownership are correct
- File content matches expected template
- No temporary or partial files remain

### Performance Metrics

Monitor:
- SSH connection time
- File transfer time
- Content comparison time
- Backup creation time
- Total deployment time per machine

## Integration with Template Evaluation

### End-to-End Workflow Test

1. **Deploy Template**:
   ```bash
   ./build/spooky action run --project examples/projects/file-distribution-testing deploy-nginx-config
   ```

2. **Evaluate Template**:
   ```bash
   ./build/spooky action run --project examples/projects/file-distribution-testing evaluate-nginx-config
   ```

3. **Verify Result**:
   - Template should be evaluated with server facts
   - Final configuration should be written to destination
   - Configuration should be valid and functional

## Future Enhancements

### Checksum-based Comparison

Future versions may use checksums for more efficient content comparison:

```go
// Compare file checksums instead of full content
func (tae *TemplateActionExecutor) hasContentChanged(sshClient *SSHClient, path string, localContent []byte) (bool, error) {
    localChecksum := sha256.Sum256(localContent)
    remoteChecksum, err := sshClient.ExecuteCommand(fmt.Sprintf("sha256sum %s | awk '{print $1}'", path))
    if err != nil {
        return true, err
    }
    return fmt.Sprintf("%x", localChecksum) != strings.TrimSpace(remoteChecksum), nil
}
```

### Atomic Deployment

Consider implementing atomic deployment using temporary files and atomic moves:

```bash
# Deploy to temporary file first, then move atomically
mv /tmp/nginx.conf.tmpl.new /tmp/nginx.conf.tmpl
```

### Rollback Support

Add automatic rollback capabilities for failed deployments:

```hcl
template {
  source = "templates/nginx-config.tmpl"
  destination = "/etc/nginx/nginx.conf"
  backup = true
  rollback_on_failure = true
}
```

## Conclusion

This test plan ensures that file distribution is robust, reliable, and production-ready. Successful completion of these tests validates the enhanced functionality needed for Issue #52 implementation and provides confidence in the deployment process. 