# Facts System User Guide

## Overview

The spooky facts system provides comprehensive fact collection and storage capabilities for gathering system information from machines. This guide covers everything from basic fact collection to advanced features like validation, export/import, and integration with other spooky components.

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
3. **Data Collection**: Use gopsutil to gather system information
4. **Data Processing**: Convert and validate collected data
5. **Storage**: Store facts in BadgerDB with machine ID as key

### Fact Storage

Facts are stored in memory during the action run with the following structure:

- **Key**: Machine ID (32-character hex string)
- **Value**: Fact collection with timestamp
- **Metadata**: Collection metadata and validation status
- **Lifetime**: Facts persist for the duration of the action run

## Basic Usage

### Exporting Facts

The `export` command automatically gathers facts from all machines in a project and exports them to a file:

```bash
# Export facts from all machines to HCL (default)
spooky facts export ./my-project --output facts.hcl

# Export facts with parallel processing
spooky facts export ./my-project --parallel 4 --output facts.hcl

# Export facts from a specific machine
spooky facts export ./my-project --machine web-server --output web-server-facts.hcl

# Export facts to JSON format
spooky facts export ./my-project --format json --output facts.json
```

#### Command Options

- `--format`: Export format (hcl, json) (default: hcl)
- `--output`: Output file path (required)
- `--machine`: Target specific machine (default: all machines)
- `--parallel`: Number of parallel workers (default: 1)
- `--verbose`: Verbose output

#### Example Output

```bash
$ spooky facts export ./my-project --output facts.hcl
INFO: Starting fact collection for export
INFO: Found 3 machines in inventory
INFO: Collecting facts from web-server (192.168.1.10)
INFO: Collecting facts from db-server (192.168.1.11)
INFO: Collecting facts from app-server (192.168.1.12)
INFO: Successfully collected facts from 3 machines
INFO: Exporting facts to facts.hcl
Successfully exported facts to: facts.hcl
```



### Fact Validation

Fact validation is performed internally during export operations to ensure data integrity and schema compliance. The validation process checks:

- **Machine ID Format**: 32-character hexadecimal string
- **Required Fields**: System, hardware, and network facts
- **Data Types**: Numeric fields are valid numbers
- **String Lengths**: String fields within limits
- **Schema Compliance**: Facts match expected structure

Validation errors and warnings are logged during export operations, but export continues to ensure data availability.

### Exporting Facts

The `export` command exports facts to various formats:

```bash
# Export all facts to JSON
spooky facts export ./my-project --format json --output facts.json

# Export facts to HCL format
spooky facts export ./my-project --format hcl --output facts.hcl

# Export facts for specific machines
spooky facts export ./my-project --machine web-server --format json --output web-server-facts.json
```

#### Export Formats

- **HCL**: HashiCorp Configuration Language format (default)
- **JSON**: Standard JSON format for data exchange

#### Example JSON Output

```json
{
  "machines": {
    "1234567890abcdef1234567890abcdef": {
      "machine_id": "1234567890abcdef1234567890abcdef",
      "collected_at": "2024-01-15T10:30:45Z",
      "facts": {
        "system": {
          "os": {
            "name": "Ubuntu",
            "version": "22.04.3 LTS",
            "arch": "x86_64",
            "kernel": "5.15.0-91-generic"
          },
          "hardware": {
            "cpu": {
              "cores": 4,
              "model": "Intel(R) Core(TM) i7-7700K",
              "frequency": 4200.0
            },
            "memory": {
              "total": 17179869184,
              "available": 8589934592,
              "used": 8589934592,
              "free": 0
            }
          }
        }
      }
    }
  }
}
```

## Advanced Usage

### Project Configuration

Facts collection can be configured in the project configuration:

```hcl
project {
  name = "my-project"
  description = "Example project with facts configuration"
  
  facts {
    # Collection settings
    parallel_workers = 4
    timeout_seconds = 60
    retry_attempts = 3
    
    # Memory settings
    memory_limit = "1GB"
    memory_efficient = true
    
    # Validation settings
    strict_validation = true
    required_facts = ["system", "hardware", "network"]
    
    # Custom facts
    custom_facts = {
      "environment" = "production"
      "team" = "platform"
      "region" = "us-west-2"
    }
  }
}
```

### Machine-Specific Configuration

Individual machines can have fact collection settings:

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
    
    # Fact collection settings
    facts {
      timeout_seconds = 30
      retry_attempts = 2
      include_processes = true
      include_network_connections = true
      
      # Custom facts for this machine
      custom_facts = {
        "role" = "web-server"
        "tier" = "frontend"
        "load_balancer" = "true"
      }
    }
  }
}
```

### Parallel Collection

For large environments, use parallel collection to improve performance:

```bash
# Use 8 parallel workers for faster collection
spooky facts gather ./my-project --parallel 8

# Monitor collection progress
spooky facts gather ./my-project --parallel 4 --verbose
```

### Custom Fact Collection

You can extend fact collection with custom facts:

```hcl
project {
  name = "my-project"
  
  facts {
    # Custom fact collection scripts
    custom_collectors = [
      "scripts/collect-docker-info.sh",
      "scripts/collect-kubernetes-info.sh",
      "scripts/collect-application-metrics.sh"
    ]
    
    # Custom fact validation
    custom_validators = [
      "scripts/validate-docker-facts.sh",
      "scripts/validate-kubernetes-facts.sh"
    ]
  }
}
```

### Fact Retention and Cleanup

Configure fact retention policies:

```hcl
project {
  name = "my-project"
  
  facts {
    # Memory settings
    memory_limit = "2GB"
    memory_efficient = true
    
    # Performance settings
    memory_pooling = true
    garbage_collection = true
    
    # Monitoring settings
    memory_monitoring = true
    memory_profile = true
  }
}
```

## Integration with Other Components

### Variables Integration

Facts can be used in variable resolution:

```hcl
variables {
  # Use facts in variable definitions
  web_server_cpu_cores = "{{ facts.web-server.system.hardware.cpu.cores }}"
  total_memory_gb = "{{ facts.web-server.system.hardware.memory.total / 1024 / 1024 / 1024 }}"
  
  # Conditional variables based on facts
  deployment_strategy = "{{ if facts.web-server.system.hardware.memory.total > 8589934592 }}rolling{{ else }}recreate{{ end }}"
}
```

### Templates Integration

Facts can be used in template rendering:

```bash
# Template using facts
cat > deploy.sh.tmpl << 'EOF'
#!/bin/bash
# Deploy application based on system facts

echo "Deploying to {{ .Machine.Hostname }}"
echo "CPU Cores: {{ .Facts.System.Hardware.CPU.Cores }}"
echo "Memory: {{ .Facts.System.Hardware.Memory.Total }} bytes"

# Adjust deployment based on system capabilities
if [ {{ .Facts.System.Hardware.CPU.Cores }} -gt 4 ]; then
    echo "High-performance deployment"
    export WORKER_THREADS={{ .Facts.System.Hardware.CPU.Cores }}
else
    echo "Standard deployment"
    export WORKER_THREADS=2
fi
EOF

# Render template with facts
spooky templates render deploy.sh.tmpl --machine web-server --output deploy.sh
```

### Actions Integration

Facts can be used in action definitions:

```hcl
actions {
  action "deploy-web" {
    description = "Deploy web application"
    
    machines = ["web-server"]
    parallel = true
    
    # Use facts in action logic
    condition = "facts.system.hardware.memory.available > 1073741824"  # 1GB
    
    template {
      source = "templates/deploy.sh.tmpl"
      destination = "/tmp/deploy.sh"
      permissions = "0755"
    }
    
    command = "/tmp/deploy.sh"
  }
}
```

## Monitoring and Maintenance

### Fact Collection Monitoring

Monitor fact collection health:

```bash
# Check fact collection status
spooky facts list ./my-project --verbose

# Validate fact integrity
spooky facts validate ./my-project

# Check fact freshness
spooky facts list ./my-project --check-freshness
```

### Memory Management

Monitor fact memory usage:

```bash
# Check memory usage
free -h

# Monitor memory during collection
spooky facts gather ./my-project --monitor-memory

# Profile memory usage
spooky facts gather ./my-project --profile
```

### Performance Optimization

Optimize fact collection performance:

```bash
# Use parallel collection for multiple machines
spooky facts gather ./my-project --parallel 8

# Adjust timeouts for slow networks
spooky facts gather ./my-project --timeout 120s

# Use connection pooling
spooky facts gather ./my-project --connection-pool 10
```

## Best Practices

### Fact Collection

- **Regular Collection**: Collect facts regularly (daily or weekly)
- **Parallel Processing**: Use parallel collection for multiple machines
- **Error Handling**: Monitor collection failures and retry
- **Validation**: Always validate collected facts

### Memory Management

- **Memory Monitoring**: Monitor memory usage during fact collection
- **Garbage Collection**: Proper cleanup of temporary objects
- **Memory Pooling**: Reuse memory structures for better performance
- **Memory Limits**: Set appropriate memory limits for large collections

### Security

- **Memory Protection**: Facts stored in memory are not persisted to disk
- **Access Control**: Facts are only accessible during the action run
- **Audit Logging**: Log fact collection and usage
- **Data Sanitization**: Remove sensitive information before collection

### Performance

- **Parallel Collection**: Use parallel workers for faster collection
- **Timeout Tuning**: Adjust timeouts based on network conditions
- **Memory Tuning**: Optimize memory allocation and usage
- **Memory Management**: Monitor memory usage during collection

## Troubleshooting

### Common Issues

#### Collection Failures

```bash
# Check SSH connectivity
ssh -i ~/.ssh/id_rsa user@machine.example.com

# Test fact collection manually
spooky facts gather ./my-project --machine machine.example.com --verbose

# Check machine configuration
spooky machines list ./my-project
```

#### Memory Issues

```bash
# Check memory usage
free -h

# Monitor memory during collection
spooky facts gather ./my-project --monitor-memory

# Validate facts integrity
spooky facts validate ./my-project
```

#### Validation Errors

```bash
# Check validation errors
spooky facts validate ./my-project --verbose

# Export and inspect facts
spooky facts export ./my-project --format json --output debug.json

# Re-collect facts if needed
spooky facts gather ./my-project --machine problematic-machine
```

### Debugging

#### Enable Debug Logging

```bash
# Enable debug logging for facts
export SPOOKY_LOG_LEVEL=debug
spooky facts gather ./my-project --verbose
```

#### Check Fact Details

```bash
# Export facts for inspection
spooky facts export ./my-project --format json --output facts-debug.json

# Validate specific machine
spooky facts validate ./my-project --machine specific-machine --verbose
```

#### Performance Analysis

```bash
# Profile fact collection
spooky facts gather ./my-project --profile

# Check collection timing
spooky facts gather ./my-project --timing
```

## Examples

### Basic Project Setup

```bash
# Create a new project
mkdir my-project
cd my-project

# Create project configuration
cat > project.hcl << 'EOF'
project {
  name = "my-project"
  description = "Example project with facts"
  
  facts {
    parallel_workers = 4
    timeout_seconds = 60
  }
}

machines {
  machine "web-server" {
    hostname = "web.example.com"
    port = 22
    user = "admin"
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/id_rsa"
    }
  }
}
EOF

# Gather facts
spooky facts gather .

# List facts
spooky facts list .

# Validate facts
spooky facts validate .
```

### Advanced Configuration

```hcl
project {
  name = "production-cluster"
  description = "Production cluster with comprehensive fact collection"
  
  facts {
    # Performance settings
    parallel_workers = 8
    timeout_seconds = 120
    retry_attempts = 3
    
    # Storage settings
    storage_path = "cluster-facts.db"
    compression_enabled = true
    
    # Validation settings
    strict_validation = true
    required_facts = ["system", "hardware", "network", "processes"]
    
    # Custom facts
    custom_facts = {
      "environment" = "production"
      "cluster" = "primary"
      "datacenter" = "us-west-2"
    }
    
    # Retention settings
    retention_days = 90
    auto_cleanup = true
    cleanup_schedule = "0 3 * * *"  # Daily at 3 AM
  }
}

machines {
  machine "web-1" {
    hostname = "web-1.example.com"
    port = 22
    user = "admin"
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/production-key"
    }
    
    facts {
      timeout_seconds = 60
      include_processes = true
      include_network_connections = true
      
      custom_facts = {
        "role" = "web-server"
        "tier" = "frontend"
        "load_balancer" = "true"
      }
    }
  }
  
  machine "db-1" {
    hostname = "db-1.example.com"
    port = 22
    user = "admin"
    
    authentication {
      method = "ssh_key"
      key_path = "~/.ssh/production-key"
    }
    
    facts {
      timeout_seconds = 90
      include_processes = true
      
      custom_facts = {
        "role" = "database"
        "tier" = "backend"
        "database_type" = "postgresql"
      }
    }
  }
}
```

This comprehensive user guide provides everything needed to effectively use the spooky facts system for collecting, managing, and utilizing machine facts in your infrastructure automation workflows.
