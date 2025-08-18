# Facts System

## Overview

The spooky facts collection system provides comprehensive system information gathering through SSH-based collection, supporting three distinct namespaces for different types of facts. Facts are collected from remote machines and stored for use in templates, actions, and decision-making processes.

## Related Systems

This system integrates with and depends on several other spooky systems:

- **[Machines System](MACHINES_SYSTEM.md)** - Facts are collected from machines defined in the inventory
- **[SSH System](SSH_SYSTEM.md)** - Facts collection uses SSH connections for remote machine access
- **[Logging System](LOGGING_SYSTEM.md)** - Fact collection generates comprehensive logs for monitoring and debugging
- **[Integrations System](INTEGRATIONS_SYSTEM.md)** - Facts integrate with other systems through the IntegrationManager
- **[Actions System](ACTIONS_SYSTEM.md)** - Actions use facts for context and decision-making
- **[Templates System](TEMPLATES_SYSTEM.md)** - Templates can use facts as data for rendering
- **[Projects System](PROJECTS_SYSTEM.md)** - Facts are organized within projects
- **[Variables System](VARIABLES_SYSTEM.md)** - Facts can be used as variables in templates and actions

## Core Concepts

### Fact Collection
Facts are system information collected from remote machines through SSH connections. The system supports:

- **Local Facts**: Collected from the local machine
- **SSH Facts**: Collected from remote machines via SSH
- **Custom Facts**: User-defined facts with custom collection logic

### Fact Storage
Facts are stored in BadgerDB format for efficient querying and retrieval:

- **Global Facts**: Stored in `$XDG_STATE_HOME/spooky/global-facts.db`
- **Project Facts**: Stored in `./facts.db` within each project
- **Encrypted Facts**: Sensitive facts can be encrypted using age encryption

### Fact Namespaces
The system supports three distinct fact namespaces:

1. **System Facts**: Operating system, hardware, and network information
2. **Application Facts**: Installed applications, services, and configurations
3. **Custom Facts**: User-defined facts for specific use cases

## CLI Commands

### Facts Management Commands

#### `spooky facts export [project-path]`
Export facts from a project to various formats.

**Flags:**
- `--output` - Output file path (required)
- `--format` - Export format (hcl, json) - default: hcl
- `--machine` - Filter to specific machine
- `--tags` - Filter by tags (supports key=value or key-only)
- `--groups` - Filter by groups
- `--parallel` - Number of parallel workers (default: 1)
- `--verbose` - Verbose output

**Examples:**
```bash
# Export all facts to HCL format
spooky facts export ./my-project --output facts.hcl

# Export facts for specific machine
spooky facts export ./my-project --machine web-server --output web-server-facts.hcl

# Export facts with tags filter
spooky facts export ./my-project --tags environment=production --output prod-facts.hcl

# Export with parallel processing
spooky facts export ./my-project --parallel 4 --output facts.hcl

# Export with verbose output
spooky facts export ./my-project --verbose --output facts.hcl
```

## Fact Collection Process

### Collection Methods
The system supports multiple collection methods:

1. **SSH Collection**: Primary method for remote machine facts
2. **Local Collection**: For local machine facts
3. **Custom Collection**: User-defined collection scripts

### SSH Collection Process
```bash
# SSH collection workflow
1. Establish SSH connection to target machine
2. Run fact collection commands
3. Parse command output into structured facts
4. Store facts in BadgerDB
5. Close SSH connection
```

### Collection Commands
The system runs various commands to collect facts:

**System Facts:**
- `uname -a` - Operating system information
- `cat /etc/os-release` - OS release information
- `hostname` - Machine hostname
- `ip addr show` - Network interface information
- `df -h` - Disk usage information
- `free -h` - Memory usage information
- `nproc` - CPU count information

**Application Facts:**
- `systemctl list-units --type=service` - Systemd services
- `ps aux` - Running processes
- `netstat -tlnp` - Network connections
- `docker ps` - Docker containers (if available)
- `kubectl get nodes` - Kubernetes nodes (if available)

## Fact Storage

### Storage Backends
The system supports multiple storage backends:

1. **BadgerDB**: Primary storage backend (default)
2. **JSON**: Human-readable format for debugging
3. **Memory**: In-memory storage for testing

### Storage Configuration
```hcl
# Global facts storage configuration
facts {
  storage {
    backend = "badgerdb"
    path = "$XDG_STATE_HOME/spooky/global-facts.db"
    encryption = false
  }
  
  collection {
    parallel_workers = 4
    timeout_seconds = 300
    retry_attempts = 3
  }
}
```

### Fact Structure
Facts are stored with the following structure:

```json
{
  "machine": "web-server.example.com",
  "namespace": "system",
  "fact": "os_info",
  "value": {
    "name": "Ubuntu",
    "version": "22.04.3 LTS",
    "arch": "x86_64"
  },
  "timestamp": "2024-01-01T12:00:00Z",
  "ttl": 3600
}
```

## Fact Validation

### Validation Rules
The system validates facts using schema-based validation:

1. **Type Validation**: Ensures fact values match expected types
2. **Range Validation**: Validates numeric values within acceptable ranges
3. **Format Validation**: Validates string values against expected formats
4. **Dependency Validation**: Ensures required facts are present

### Validation Configuration
```hcl
# Fact validation configuration
facts {
  validation {
    enabled = true
    strict_mode = false
    custom_rules = [
      "custom-validation.hcl"
    ]
  }
}
```

## Fact Encryption

### Encryption Support
Sensitive facts can be encrypted using age encryption:

```hcl
# Encrypted fact configuration
facts {
  variable "database_password" {
    value = "secret123"
    encrypted = true
    description = "Database password"
  }
}
```

### Encryption Process
1. **Identify Encrypted Facts**: Facts with `encrypted = true`
2. **Load Recipients**: Read from `identities/recipients.txt`
3. **Encrypt Values**: Use age encryption
4. **Store Encrypted**: Replace plaintext with encrypted armor

## Fact Querying

### Query Language
The system supports a rich query language for fact retrieval:

```bash
# Basic queries
facts.machine == "web-server"
facts.namespace == "system"
facts.fact == "os_info"

# Complex queries
facts.machine == "web-server" AND facts.namespace == "system"
facts.value.memory_total > 8589934592  # 8GB
facts.value.cpu_count >= 4
```

### Query Examples
```bash
# Query system facts
spooky facts query "facts.namespace == 'system'"

# Query specific machine
spooky facts query "facts.machine == 'web-server'"

# Query memory usage
spooky facts query "facts.fact == 'memory_info' AND facts.value.total > 8589934592"
```

## Integration with Other Systems

### Template Integration
Facts are used in templates for dynamic configuration:

```hcl
# Template using facts
template "nginx.conf" {
  content = <<-EOF
    server {
      listen 80;
      server_name {{ facts.machine }};
      
      location / {
        root /var/www/html;
        index index.html;
      }
    }
  EOF
}
```

### Action Integration
Facts are used in actions for conditional running:

```hcl
# Action using facts
action "deploy" {
  condition = "facts.value.memory_total > 4294967296"  # 4GB
  command = "deploy.sh"
}
```

### Variable Integration
Facts can be used as variable values:

```hcl
# Variable using facts
variables {
  variable "server_hostname" {
    value = "{{ facts.machine }}"
    description = "Server hostname from facts"
  }
}
```

## Performance Optimization

### Parallel Collection
The system supports parallel fact collection:

```bash
# Collect facts in parallel
spooky facts export ./my-project --parallel 4
```

### Caching
Facts are cached to improve performance:

- **Memory Cache**: Frequently accessed facts cached in memory
- **Disk Cache**: Persistent cache on disk
- **TTL-based Expiration**: Facts expire based on TTL settings

### Optimization Strategies
1. **Batch Collection**: Collect multiple facts in single SSH session
2. **Connection Pooling**: Reuse SSH connections
3. **Incremental Updates**: Only collect changed facts
4. **Compression**: Compress fact data for storage

## Error Handling

### Collection Errors
The system handles various collection errors:

1. **SSH Connection Errors**: Network connectivity issues
2. **Authentication Errors**: SSH key or password issues
3. **Command Running Errors**: Commands failing on remote machines
4. **Parsing Errors**: Invalid command output

### Error Recovery
```bash
# Retry failed collections
spooky facts export ./my-project --retry-failed

# Validate fact collection
spooky facts validate ./my-project
```

## Troubleshooting

### Common Issues

#### SSH Connection Issues
```bash
# Test SSH connectivity
spooky machines ping ./my-project --machine web-server

# Check SSH configuration
ssh -i ~/.ssh/id_rsa user@web-server.example.com
```

#### Fact Collection Issues
```bash
# Enable verbose output
spooky facts export ./my-project --verbose --output facts.hcl

# Check fact storage
ls -la ./my-project/facts.db
```

#### Validation Issues
```bash
# Validate fact structure
spooky facts validate ./my-project

# Check fact schemas
spooky schemas validate facts.schema.hcl facts.hcl
```

### Debug Information
```bash
# Enable debug logging
export SPOOKY_LOG_LEVEL=debug

# Export with debug output
spooky facts export ./my-project --verbose --output facts.hcl
```

## Best Practices

### Collection Practices
1. **Use Appropriate Timeouts**: Set reasonable timeouts for fact collection
2. **Implement Retry Logic**: Handle temporary network issues
3. **Validate Facts**: Ensure collected facts are accurate
4. **Monitor Performance**: Track collection performance
5. **Secure Storage**: Encrypt sensitive facts

### Storage Practices
1. **Regular Backups**: Backup fact databases regularly
2. **Cleanup Old Facts**: Remove expired facts
3. **Monitor Storage**: Track storage usage
4. **Optimize Queries**: Use efficient query patterns
5. **Validate Data**: Ensure data integrity

### Security Practices
1. **Encrypt Sensitive Facts**: Use age encryption for secrets
2. **Secure SSH Keys**: Protect SSH private keys
3. **Limit Access**: Restrict access to fact databases
4. **Audit Access**: Log fact access and modifications
5. **Validate Inputs**: Sanitize fact inputs

## Future Enhancements

### Planned Features
- **Real-time Collection**: Continuous fact collection
- **Fact Streaming**: Stream facts to external systems
- **Custom Collectors**: User-defined fact collectors
- **Fact Analytics**: Advanced fact analysis and reporting
- **Fact Visualization**: Web-based fact visualization

### Extension Points
- **Custom Fact Types**: User-defined fact types
- **External Integrations**: Integration with monitoring systems
- **Fact Plugins**: Pluggable fact collection modules
- **Fact APIs**: REST API for fact access
- **Fact Webhooks**: Webhook notifications for fact changes
