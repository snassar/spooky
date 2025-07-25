# Server Facts Testing Plan

## Overview

This document outlines the testing plan for server-side fact collection in spooky's template evaluation system. The goal is to verify that all machine fields, tags, and live system facts are properly collected and made available to templates during server-side evaluation.

## Background

Server-side fact collection is a key component of Issue #52 (server-side template evaluation). When templates are evaluated on remote servers, spooky needs to:

1. SSH to the target server
2. Collect live system facts (OS version, hostname, disk space, etc.)
3. Include machine metadata (name, host, port, tags, etc.)
4. Make all facts available to templates via `serverFact` functions

## Test Template

### File: `server-facts-debug.tmpl`

Place this template in your test project's `templates/` directory:

```gotemplate
# Server Facts Debug Template

## Machine Fields
- Name: {{serverFact "name"}}
- Host: {{serverFact "host"}}
- Port: {{serverFact "port"}}
- User: {{serverFact "user"}}
- Keyfile: {{serverFact "keyfile"}}
- Password: {{serverFact "password"}}

## Tags
{{- $tags := serverFact "tags" }}
{{- if $tags }}
Tags:
{{- range $k, $v := $tags }}
  - {{ $k }}: {{ $v }}
{{- end }}
{{- else }}
No tags found.
{{- end }}

## System Facts
- Machine ID: {{serverFact "machine_id"}}
- OS Version: {{serverFact "os_version"}}
- Hostname: {{serverFact "hostname"}}
- IP Address: {{serverFact "ip_address"}}
- Disk Space: {{serverFact "disk_space"}}
- Memory Info: {{serverFact "memory_info"}}
- CPU Info: {{serverFact "cpu_info"}}
- Uptime: {{serverFact "uptime"}}
- Kernel: {{serverFact "kernel"}}
- Architecture: {{serverFact "architecture"}}

## File System Facts
- /etc/passwd exists: {{serverFact "etc_exists"}}
- /var exists: {{serverFact "var_exists"}}
- /tmp writable: {{serverFact "tmp_writable"}}
- /home exists: {{serverFact "home_exists"}}

## Service Facts
- Nginx running: {{serverFact "nginx_running"}}
- Apache running: {{serverFact "apache_running"}}
- MySQL running: {{serverFact "mysql_running"}}
- Postgres running: {{serverFact "postgres_running"}}
- Docker running: {{serverFact "docker_running"}}
```

## Test Setup

### 1. Test Project Structure

```
examples/projects/template-testing/
├── project.hcl
├── inventory.hcl
├── actions.hcl
├── templates/
│   └── server-facts-debug.tmpl
└── custom-facts.json
```

### 2. Inventory Configuration

Ensure your `inventory.hcl` contains a reachable SSH target:

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
      region = "us-west"
    }
  }
}
```

### 3. SSH Access Requirements

- Target machine must be accessible via SSH
- User must have sufficient permissions to run fact collection commands
- For local testing, consider using:
  - Localhost with SSH enabled
  - Docker container with SSH
  - Virtual machine with SSH access

## Test Execution

### Step 1: Basic Template Rendering

```bash
# Test template rendering with server facts
./build/spooky project render-template templates/server-facts-debug.tmpl examples/projects/template-testing --server test-server --dry-run
```

**Expected Output:**
- All machine fields should display correctly
- Tags should be listed with key-value pairs
- System facts should show live data from the server
- File system facts should show boolean values
- Service facts should show service status

### Step 2: Error Handling Test

```bash
# Test with non-existent server
./build/spooky project render-template templates/server-facts-debug.tmpl examples/projects/template-testing --server non-existent-server --dry-run
```

**Expected Output:**
- Clear error message indicating server not found in inventory
- No template rendering should occur

### Step 3: SSH Connection Failure Test

```bash
# Test with unreachable server (wrong IP/credentials)
./build/spooky project render-template templates/server-facts-debug.tmpl examples/projects/template-testing --server unreachable-server --dry-run
```

**Expected Output:**
- SSH connection error should be logged
- Template should render with "unknown" values for system facts
- Machine fields should still be available

### Step 4: Fact Collection Failure Test

Test scenarios where individual fact collection commands fail:

```bash
# Test with server that has limited command access
./build/spooky project render-template templates/server-facts-debug.tmpl examples/projects/template-testing --server restricted-server --dry-run
```

**Expected Output:**
- Failed facts should show as "unknown" or "false"
- Successful facts should display normally
- Warning logs should indicate which facts failed

## Validation Criteria

### ✅ Success Criteria

1. **Machine Fields**: All inventory fields (name, host, port, user, tags) are available
2. **System Facts**: Live system information is collected and displayed
3. **Error Handling**: Failed fact collection doesn't break template rendering
4. **Extensibility**: New fact types can be easily added
5. **Performance**: Fact collection completes within reasonable time (30s timeout)

### ❌ Failure Indicators

1. **Missing Facts**: Expected facts are not available in template
2. **SSH Errors**: Connection failures not handled gracefully
3. **Template Errors**: Template syntax errors due to missing functions
4. **Performance Issues**: Fact collection takes too long or times out
5. **Security Issues**: Sensitive information (passwords) exposed in templates

## Test Data Collection

### Log Analysis

Check logs for:
- SSH connection success/failure
- Individual fact collection results
- Warning messages for failed facts
- Performance metrics (collection time)

### Output Validation

Verify that:
- Machine fields match inventory configuration
- System facts reflect actual server state
- Tags are properly formatted and accessible
- Boolean facts (file system, services) show correct true/false values

## Future Enhancements

### Custom Fact Commands

Future versions may support custom fact commands defined in inventory:

```hcl
machine "test-server" {
  host = "192.168.1.100"
  user = "debian"
  
  custom_facts = {
    "app_version" = "cat /opt/app/version"
    "db_status" = "systemctl is-active postgresql"
  }
}
```

### Fact Caching

Consider implementing fact caching to improve performance for repeated evaluations.

### Fact Validation

Add validation to ensure collected facts are in expected format and range.

## Conclusion

This test plan ensures that server-side fact collection is robust, extensible, and provides all necessary data for template evaluation. Successful completion of these tests validates the core functionality needed for Issue #52 implementation. 