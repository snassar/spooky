# Facts System User Guide

## Overview

The spooky facts system provides fact collection and export capabilities for gathering system information from machines. This guide covers the current implementation of fact collection and export functionality.

**Status: Partially Implemented** - The facts system has basic functionality but SSH-based fact collection has known issues that need to be addressed.

## Getting Started

### Prerequisites

- spooky CLI installed and configured
- SSH access to target machines
- Basic understanding of HCL configuration syntax
- Access to create and modify project files

### Quick Start

1. **Check Available Facts Commands**
   ```bash
   spooky facts --help
   ```

2. **Export Facts from a Project**
   ```bash
   spooky facts export ./my-project --format json --output facts.json
   ```

## Facts System Concepts

### What are Facts?

Facts are system information collected from machines including:

- **Operating System**: OS name, version, architecture, kernel
- **Hardware**: CPU, memory, disk, network interfaces
- **System State**: Load average, running processes, disk usage
- **Network**: Interface configuration, network connections
- **Custom Data**: User-defined facts and metadata

### Fact Collection Process

The fact collection process works as follows:

1. **Machine Discovery**: Read machine inventory from project configuration
2. **SSH Connection**: Establish secure connection to each machine
3. **Data Collection**: Use SSH commands to gather system information
4. **Custom Facts**: Read custom facts from `/etc/spooky/facts.{hcl,json}` if present
5. **Data Processing**: Convert and validate collected data
6. **Direct Export**: Write facts directly to output file in requested format
7. **Cleanup**: Memory is automatically managed during export

### Fact Storage

Facts are gathered directly and exported without intermediate storage:

- **Direct Export**: Facts are collected from machines and exported immediately
- **No Intermediate Storage**: No temporary storage step - direct gather → export
- **Memory Management**: Memory is automatically managed during export operations
- **Cleanup**: No cleanup required - memory is freed automatically

## Current Implementation Status

### ✅ Working Features

- **Basic Fact Collection**: System fact collection using SSH commands
- **Machine Integration**: Uses project machine inventory for target identification
- **Export Functionality**: Facts export to JSON and HCL formats
- **CLI Integration**: `spooky facts export` command with filtering options
- **Basic Validation**: Fact collection validation and error handling

### ⚠️ Known Issues

- **SSH-Based Collection**: SSH-based fact collection has implementation issues
- **Remote Facts Reading**: Cannot reliably read `/etc/spooky/facts.*` files from remote machines
- **Parallel Processing**: Sequential collection only, no multi-machine parallel processing
- **SSH Integration**: Cannot fully leverage existing SSH infrastructure and machine inventory

### 🔧 Current Workarounds

- **Local Collection**: Use local fact collection for immediate needs
- **Manual Export**: Export facts manually from remote machines if needed
- **Filtered Collection**: Use machine filtering to limit collection scope
- **Monitor Updates**: Watch for improvements to SSH-based collection

## Basic Usage

### Exporting Facts

The primary way to use the facts system is through the export command:

```bash
# Export all facts from a project
spooky facts export ./my-project --format json --output facts.json

# Export with specific format
spooky facts export ./my-project --format hcl --output facts.hcl

# Export to stdout
spooky facts export ./my-project --format json
```

### Machine Filtering

Filter facts collection by specific machines:

```bash
# Export facts from specific machines
spooky facts export ./my-project --machine web-server --format json --output web-facts.json

# Export from multiple machines
spooky facts export ./my-project --machine web-server --machine db-server --format json --output server-facts.json
```

### Tag-Based Filtering

Filter facts collection by machine tags:

```bash
# Export facts from machines with specific tags
spooky facts export ./my-project --tags environment=production --format json --output prod-facts.json

# Export from machines with multiple tag criteria
spooky facts export ./my-project --tags environment=production --tags role=web --format json --output prod-web-facts.json
```

## Project Configuration

### Basic Project Setup

Create a project with machine inventory for fact collection:

```hcl
# project.hcl
project {
  name = "my-project"
  description = "Project with fact collection"
}

# machines.hcl
machines {
  machine "web-server" {
    hostname = "web.example.com"
    host = "192.168.1.10"
    user = "admin"
    port = 22
    
    tags = {
      environment = "production"
      role = "web"
    }
  }
  
  machine "db-server" {
    hostname = "db.example.com"
    host = "192.168.1.11"
    user = "admin"
    port = 22
    
    tags = {
      environment = "production"
      role = "database"
    }
  }
}
```

### Machine Authentication

Configure SSH authentication for fact collection:

```hcl
# machines.hcl
machines {
  machine "web-server" {
    hostname = "web.example.com"
    user = "admin"
    
    # SSH key authentication
    key_file = "~/.ssh/id_ed25519"
    
    # Optional passphrase for encrypted keys
    passphrase = "your-passphrase"
    
    tags = {
      environment = "production"
    }
  }
}
```

## Export Formats

### JSON Format

JSON format is ideal for machine processing and integration:

```bash
spooky facts export ./my-project --format json --output facts.json
```

**JSON Output Example**:
```json
{
  "machines": {
    "web-server": {
      "system": {
        "os": {
          "name": "Ubuntu",
          "version": "22.04.3 LTS",
          "architecture": "x86_64"
        },
        "hardware": {
          "cpu_count": 4,
          "memory_total": 8589934592,
          "disk_total": 107374182400
        }
      },
      "network": {
        "interfaces": {
          "eth0": {
            "address": "192.168.1.10",
            "netmask": "255.255.255.0"
          }
        }
      }
    }
  }
}
```

### HCL Format

HCL format is human-readable and follows spooky configuration patterns:

```bash
spooky facts export ./my-project --format hcl --output facts.hcl
```

**HCL Output Example**:
```hcl
machines {
  machine "web-server" {
    system {
      os {
        name = "Ubuntu"
        version = "22.04.3 LTS"
        architecture = "x86_64"
      }
      
      hardware {
        cpu_count = 4
        memory_total = 8589934592
        disk_total = 107374182400
      }
    }
    
    network {
      interfaces {
        eth0 {
          address = "192.168.1.10"
          netmask = "255.255.255.0"
        }
      }
    }
  }
}
```

## Advanced Usage

### Custom Facts

Create custom facts on remote machines:

```bash
# On remote machine, create custom facts
sudo mkdir -p /etc/spooky
sudo tee /etc/spooky/facts.hcl > /dev/null <<EOF
custom_facts {
  application {
    name = "my-app"
    version = "1.2.3"
    environment = "production"
  }
  
  monitoring {
    enabled = true
    metrics_port = 9100
  }
}
EOF
```

### Environment-Specific Configuration

Use different configurations for different environments:

```hcl
# machines.hcl
machines {
  machine "web-prod" {
    hostname = "web-prod.example.com"
    user = "admin"
    tags = {
      environment = "production"
      role = "web"
    }
  }
  
  machine "web-staging" {
    hostname = "web-staging.example.com"
    user = "admin"
    tags = {
      environment = "staging"
      role = "web"
    }
  }
}
```

```bash
# Export production facts
spooky facts export ./my-project --tags environment=production --format json --output prod-facts.json

# Export staging facts
spooky facts export ./my-project --tags environment=staging --format json --output staging-facts.json
```

## Troubleshooting

### Common Issues

#### SSH Connection Failures

**Problem**: SSH connections fail during fact collection

**Solutions**:
1. Verify SSH connectivity manually:
   ```bash
   ssh admin@web.example.com
   ```

2. Check SSH key permissions:
   ```bash
   chmod 600 ~/.ssh/id_ed25519
   ```

3. Verify machine inventory configuration:
   ```bash
   spooky machines validate ./my-project
   ```

#### Export Format Issues

**Problem**: Export fails with format errors

**Solutions**:
1. Ensure output directory exists:
   ```bash
   mkdir -p /path/to/output/directory
   ```

2. Check file permissions:
   ```bash
   chmod 755 /path/to/output/directory
   ```

3. Use absolute paths for output:
   ```bash
   spooky facts export ./my-project --format json --output /absolute/path/facts.json
   ```

#### Machine Filtering Issues

**Problem**: Machine filtering doesn't work as expected

**Solutions**:
1. Verify machine names in inventory:
   ```bash
   spooky machines list ./my-project
   ```

2. Check tag syntax:
   ```bash
   # Correct syntax
   spooky facts export ./my-project --tags environment=production
   
   # Incorrect syntax
   spooky facts export ./my-project --tags "environment=production"
   ```

### Debugging

Enable verbose output for debugging:

```bash
# Enable debug logging
export SPOOKY_LOG_LEVEL=debug

# Run export with verbose output
spooky facts export ./my-project --format json --output facts.json
```

## Best Practices

### Fact Collection Strategy

1. **Use Appropriate Filtering**: Filter by machines or tags to limit collection scope
2. **Choose Export Format**: Use JSON for machine processing, HCL for human readability
3. **Handle SSH Issues**: Be aware of current SSH-based collection limitations
4. **Validate Configuration**: Validate machine inventory before fact collection
5. **Monitor Performance**: Use appropriate filtering to avoid overwhelming systems

### Security Considerations

1. **SSH Key Management**: Use dedicated SSH keys for fact collection
2. **Key Permissions**: Ensure SSH keys have correct permissions (600)
3. **Network Security**: Use VPN or secure networks for fact collection
4. **Sensitive Data**: Be aware that facts may contain sensitive system information

### Performance Optimization

1. **Filter Appropriately**: Use machine and tag filtering to limit collection scope
2. **Batch Operations**: Export facts in batches for large environments
3. **Monitor Resources**: Watch for memory and network usage during collection
4. **Use Local Collection**: Use local fact collection when SSH-based collection has issues

## Integration with Other Systems

### Variables Integration

Use facts in variable definitions:

```hcl
# variables.hcl
variables {
  app_version = "${facts.system.os.version}"
  cpu_count = "${facts.system.hardware.cpu_count}"
  memory_gb = "${facts.system.hardware.memory_total / 1073741824}"
}
```

### Templates Integration

Use facts in template rendering:

```bash
# Render template with facts data
spooky templates render ./my-project templates/config.tmpl --data facts.json --output config.conf
```

### Actions Integration

Use facts in action configurations:

```hcl
# actions.hcl
actions {
  action "deploy-app" {
    command = "echo 'Deploying to ${facts.system.os.name} ${facts.system.os.version}'"
    machines = ["web-server"]
  }
}
```

## Remember

**Good facts system usage enables:**
- Efficient system information collection
- Machine inventory integration
- Flexible export and filtering
- Integration with other spooky systems
- Performance-optimized data gathering

**Always be aware of current SSH-based collection limitations and use appropriate workarounds until these issues are resolved.**
