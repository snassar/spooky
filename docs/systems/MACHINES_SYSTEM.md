# Machines System

## Overview

The spooky machines system provides comprehensive machine inventory management including listing, validation, connectivity testing, and export capabilities. Machine inventory is defined in `machines.hcl` files within spooky projects and contains SSH connection details, authentication information, and machine metadata.

## Related Systems

This system integrates with and depends on several other spooky systems:

- **[SSH System](SSH_SYSTEM.md)** - Machine connectivity uses SSH for remote access and command running
- **[Facts System](FACTS_SYSTEM.md)** - Machines provide facts through the collection system
- **[Logging System](LOGGING_SYSTEM.md)** - Machine operations generate comprehensive logs for monitoring and debugging
- **[Integrations System](INTEGRATIONS_SYSTEM.md)** - Machine management integrates with other systems through the IntegrationManager
- **[Actions System](ACTIONS_SYSTEM.md)** - Actions run on machines defined in the inventory
- **[Projects System](PROJECTS_SYSTEM.md)** - Machine inventory is organized within projects
- **[Templates System](TEMPLATES_SYSTEM.md)** - Templates can use machine information for rendering
- **[Variables System](VARIABLES_SYSTEM.md)** - Machine metadata can be used as variables

## Core Concepts

### Machine Inventory
Machine inventory is defined in `machines.hcl` files that specify:

- **Connection Details**: Hostname, port, user, authentication method
- **Authentication**: SSH keys, passwords, or other authentication methods
- **Metadata**: Tags, groups, descriptions, and custom attributes
- **Configuration**: Timeouts, retry settings, and connection parameters

### Machine Validation
The system validates machine configurations:

- **Syntax Validation**: Ensures valid HCL syntax
- **Schema Validation**: Validates against machine schema
- **Connectivity Testing**: Tests SSH connectivity
- **Authentication Testing**: Validates authentication credentials

### Machine Export
Machines can be exported to various formats:

- **JSON Format**: Machine-readable format for integration
- **HCL Format**: Human-readable configuration format
- **Filtered Export**: Export specific machines by tags or hostname

## CLI Commands

### Machine Management Commands

#### `spooky machines list [project-path]`
List all machines defined in the project's machine inventory.

**Examples:**
```bash
# List all machines in project
spooky machines list ./my-project

# List with specific format
spooky machines list ./my-project --format json
spooky machines list ./my-project --format hcl
```

#### `spooky machines validate [project-path]`
Validate machine inventory configuration.

**Examples:**
```bash
# Validate all machines
spooky machines validate ./my-project
```

#### `spooky machines ping [project-path]`
Test connectivity to machines in the inventory.

**Flags:**
- `--machine` - Ping specific machine by hostname
- `--parallel` - Number of parallel workers (default: 1)
- `--tags` - Filter machines by tags
- `--format` - Output format (text, json)
- `--verbose` - Show detailed output for all machines
- `--auth` - Test authentication in addition to connectivity

**Examples:**
```bash
# Ping all machines
spooky machines ping ./my-project

# Ping specific machine
spooky machines ping ./my-project --machine web-server

# Ping with parallel processing
spooky machines ping ./my-project --parallel 4

# Ping with authentication test
spooky machines ping ./my-project --auth
```

#### `spooky machines export [project-path]`
Export machine inventory to various formats.

**Flags:**
- `--output` - Output file path (required)
- `--machine` - Export specific machine by hostname
- `--tags` - Filter machines by tags (key=value or key-only)

**Examples:**
```bash
# Export all machines to JSON
spooky machines export ./my-project --output machines.json

# Export specific machine
spooky machines export ./my-project --machine web-server --output web-server.json

# Export with tags filter
spooky machines export ./my-project --tags production --output prod-machines.json
```

#### `spooky machines encrypt [project-path]`
Encrypt all machines in a project that have `encrypted=true`.

**Flags:**
- `--dry-run` - Show what would be encrypted without making changes

**Examples:**
```bash
# Encrypt machines with dry-run
spooky machines encrypt ./my-project --dry-run

# Encrypt machines
spooky machines encrypt ./my-project
```

## Machine Configuration

### Basic Machine Configuration
```hcl
# machines.hcl
machines_inventory {
  machines {
    machine "web-server" {
      hostname = "web.example.com"
      port = 22
      user = "admin"
      
      authentication {
        method = "ssh_key"
        key_path = "~/.ssh/id_rsa"
      }
      
      tags = ["web", "production"]
      description = "Web application server"
    }
  }
}
```

### Advanced Machine Configuration
```hcl
# machines.hcl
machines_inventory {
  machines {
    machine "database-server" {
      hostname = "db.example.com"
      port = 2222
      user = "dbadmin"
      
      authentication {
        method = "ssh_key"
        key_path = "~/.ssh/db_key"
        passphrase {
          value = "encrypted_passphrase"
          encrypted = true
        }
      }
      
      connection {
        timeout = "30s"
        max_retries = 3
        retry_delay = "5s"
      }
      
      tags = ["database", "production", "critical"]
      description = "Primary database server"
      
      metadata {
        environment = "production"
        location = "us-east-1"
        role = "primary-db"
      }
    }
  }
}
```

### Authentication Methods
The system supports multiple authentication methods:

**SSH Key Authentication:**
```hcl
authentication {
  method = "ssh_key"
  key_path = "~/.ssh/id_rsa"
  passphrase = "encrypted_passphrase"  # Optional
}
```

**Password Authentication:**
```hcl
authentication {
  method = "password"
  password {
    value = "encrypted_password"
    encrypted = true
  }
}
```

**Encrypted Authentication:**
```hcl
authentication {
  method = "ssh_key"
  key_path = "~/.ssh/id_rsa"
  encrypted = true
}
```

## Machine Validation

### Validation Process
The machine validation system performs:

1. **Syntax Validation**: Validates HCL syntax
2. **Schema Validation**: Validates against machine schema
3. **Field Validation**: Validates required fields and data types
4. **Authentication Validation**: Validates authentication configuration

**Note**: For connectivity and authentication testing, use the `spooky machines ping` command instead.

### Validation Output
```bash
# Validate machines
spooky machines validate ./my-project
```

**Output:**
```
🔍 Validating machine inventory: ./my-project
✅ Machine validation passed - all machines valid
📋 Schema compliance: machines.schema.hcl ✅
📋 Authentication configuration: Valid ✅
📋 Connection settings: Valid ✅
```

### Connectivity Testing
```bash
# Test connectivity using ping command
spooky machines ping ./my-project

# Test connectivity and authentication
spooky machines ping ./my-project --auth
```

**Output:**
```
🔍 Testing connectivity for 3 machines...
✅ web-server.example.com:22 - Connected (0.045s)
✅ db-server.example.com:22 - Connected (0.078s)
❌ backup-server.example.com:22 - Connection failed: connection refused
```

## Machine Export

### Export Formats
The system supports multiple export formats:

**JSON Format:**
```json
{
  "machines": [
    {
      "hostname": "web-server",
      "host": "web.example.com",
      "port": 22,
      "user": "admin",
      "authentication": {
        "method": "ssh_key",
        "key_path": "~/.ssh/id_rsa"
      },
      "tags": ["web", "production"],
      "description": "Web application server"
    }
  ]
}
```

**HCL Format:**
```hcl
machines {
  machine "web-server" {
    hostname = "web.example.com"
    port = 22
    user = "admin"
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/id_rsa"
    }
    
    tags = ["web", "production"]
    description = "Web application server"
  }
}
```

### Filtered Export
```bash
# Export specific machine
spooky machines export ./my-project --machine web-server --output web-server.json

# Export by tags
spooky machines export ./my-project --tags production --output prod-machines.json

# Export by multiple tags
spooky machines export ./my-project --tags web,production --output web-prod-machines.json
```

## Machine Encryption

### Encryption Process
The machine encryption system:

1. **Scans Machines**: Processes `machines.hcl` files
2. **Identifies Encrypted Fields**: Finds fields with `encrypted = true`
3. **Loads Recipients**: Reads from `identities/recipients.txt`
4. **Encrypts Values**: Uses age encryption for sensitive data
5. **Updates Files**: Replaces plaintext values with encrypted armor

### Encryption Configuration
```hcl
# Encrypted machine configuration
machines {
  machine "sensitive-server" {
    hostname = "sensitive.example.com"
    user = "admin"
    password = "secret123"
    encrypted = true
    
    authentication {
      method = "password"
      password = "secret456"
      encrypted = true
    }
  }
}
```

### Encryption Commands
```bash
# Dry run to see what would be encrypted
spooky machines encrypt ./my-project --dry-run

# Actually encrypt the machines
spooky machines encrypt ./my-project
```

## Machine Connectivity Testing

### Ping Process
The ping command tests connectivity to machines:

1. **DNS Resolution**: Resolves hostname to IP address
2. **Port Connectivity**: Tests TCP connectivity to SSH port
3. **SSH Handshake**: Performs SSH protocol handshake
4. **Authentication Test**: Tests authentication (if --auth flag used)

### Ping Output
```bash
# Basic ping
spooky machines ping ./my-project
```

**Output:**
```
🔍 Testing connectivity for 3 machines...
✅ web-server.example.com:22 - Connected (0.045s)
✅ db-server.example.com:22 - Connected (0.078s)
❌ backup-server.example.com:22 - Connection failed: connection refused
```

### Parallel Ping
```bash
# Ping with parallel processing
spooky machines ping ./my-project --parallel 4
```

### Authentication Testing
```bash
# Test authentication
spooky machines ping ./my-project --auth
```

**Output:**
```
🔍 Testing connectivity and authentication for 3 machines...
✅ web-server.example.com:22 - Connected, Auth OK (0.045s)
✅ db-server.example.com:22 - Connected, Auth OK (0.078s)
❌ backup-server.example.com:22 - Auth failed: invalid key
```

## Machine Filtering

### Tag-based Filtering
```bash
# Filter by single tag
spooky machines ping ./my-project --tags production

# Filter by multiple tags
spooky machines ping ./my-project --tags web,production

# Filter by key-value tags
spooky machines ping ./my-project --tags environment=production
```

### Machine-specific Operations
```bash
# Validate specific machine
spooky machines validate ./my-project --machine web-server

# Ping specific machine
spooky machines ping ./my-project --machine web-server

# Export specific machine
spooky machines export ./my-project --machine web-server --output web-server.json
```

## Integration with Other Systems

### SSH Integration
Machines integrate with the SSH system for connectivity:

```bash
# Test SSH connectivity
spooky machines ping ./my-project

# Use machines in actions
spooky actions run ./my-project --machine web-server
```

### Facts Integration
Machines provide targets for fact collection:

```bash
# Collect facts from machines
spooky facts export ./my-project --machine web-server
```

### Actions Integration
Machines serve as targets for action running:

```bash
# Run actions on specific machines
spooky actions run ./my-project --machine web-server --action deploy
```

## Error Handling

### Common Validation Errors
- **Missing Required Fields**: Hostname, user, or authentication missing
- **Invalid Authentication**: Invalid SSH key or password
- **Connection Failures**: Network connectivity issues
- **Schema Violations**: Configuration doesn't match schema

### Error Recovery
```bash
# Fix validation errors
spooky machines validate ./my-project

# Address specific issues
# ... fix machine configuration

# Re-validate
spooky machines validate ./my-project
```

## Troubleshooting

### Common Issues

#### Connection Issues
```bash
# Test basic connectivity
ping web-server.example.com

# Test SSH connectivity
ssh -i ~/.ssh/id_rsa admin@web-server.example.com

# Test with verbose output
spooky machines ping ./my-project --verbose
```

#### Authentication Issues
```bash
# Check SSH key permissions
ls -la ~/.ssh/id_rsa

# Test SSH key
ssh-keygen -l -f ~/.ssh/id_rsa

# Test authentication
spooky machines ping ./my-project --auth
```

#### Configuration Issues
```bash
# Validate machine configuration
spooky machines validate ./my-project

# Check HCL syntax
hclfmt -w machines.hcl
```

### Debug Information
```bash
# Enable verbose logging
export SPOOKY_LOG_LEVEL=debug

# Run with debug output
spooky machines ping ./my-project --verbose
```

## Best Practices

### Machine Organization
1. **Use Descriptive Names**: Choose clear, descriptive machine names
2. **Organize by Purpose**: Group machines by function or environment
3. **Use Tags**: Apply meaningful tags for filtering and organization
4. **Documentation**: Include descriptions for complex configurations
5. **Version Control**: Use git for machine inventory version control

### Security Practices
1. **Use SSH Keys**: Prefer SSH keys over passwords
2. **Encrypt Sensitive Data**: Use `encrypted = true` for secrets
3. **Secure File Permissions**: Set appropriate file permissions
4. **Regular Key Rotation**: Rotate SSH keys regularly
5. **Access Auditing**: Monitor machine access and changes

### Performance Practices
1. **Optimize Parallel Operations**: Use appropriate parallel worker counts
2. **Connection Pooling**: Reuse SSH connections when possible
3. **Timeout Configuration**: Set appropriate timeouts for operations
4. **Retry Logic**: Implement retry logic for transient failures
5. **Resource Monitoring**: Monitor resource usage during operations

## Future Enhancements

### Planned Features
- **Machine Discovery**: Automatic machine discovery
- **Machine Templates**: Reusable machine templates
- **Machine Backup**: Automated machine configuration backup
- **Machine Analytics**: Usage analytics and performance metrics
- **Machine Collaboration**: Multi-user machine management

### Extension Points
- **Custom Validators**: Machine-specific validation rules
- **Custom Exporters**: Extended export formats
- **Custom Filters**: Advanced filtering capabilities
- **External Integrations**: Integration with cloud providers
- **Custom Workflows**: Automated machine workflows
