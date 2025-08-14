# Facts Collection System

## Overview

The spooky facts collection system provides comprehensive system information gathering through SSH-based collection, supporting three distinct namespaces for different types of facts.

## Architecture

### Three-Namespace Structure

The facts collection system organizes collected facts into three distinct namespaces:

1. **System Facts** - Basic system information collected via SSH commands
2. **Collector Facts** - Advanced system metrics from spooky-collector binary
3. **Custom Facts** - User-defined facts from HCL configuration files

### Collection Flow

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   System Facts  │    │ Collector Facts │    │   Custom Facts  │
│   (SSH Commands)│    │(spooky-collector)│   │  (HCL Files)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Facts Manager                                │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │  SSH Manager    │  │   HCL Parser    │  │   Validation    │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Fact Collection                              │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   System Facts  │  │ Collector Facts │  │   Custom Facts  │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

## System Facts (SSH Commands)

### Overview
System facts are collected using standard SSH commands that can be executed with user privileges. This provides basic system information without requiring elevated permissions.

### Collection Commands

#### Operating System Information
```bash
# OS details
uname -a
cat /etc/os-release
hostname
```

#### Hardware Information
```bash
# CPU information
cat /proc/cpuinfo
nproc
lscpu

# Memory information
cat /proc/meminfo
free -h
```

#### Disk Information
```bash
# Disk usage
df -h
lsblk
mount
```

#### Network Information
```bash
# Network interfaces
ip addr show
hostname -I
netstat -i
```

### Error Handling
- **SSH Connection Failure**: Collection fails completely
- **Command Failure**: Collection continues with partial data
- **Permission Denied**: Command is skipped, collection continues

### Example Output
```json
{
  "system": {
    "os": {
      "hostname": "web-server-01",
      "kernel": "Linux 5.15.0-generic",
      "distribution": "Ubuntu 22.04.3 LTS"
    },
    "hardware": {
      "cpu": {
        "cores": 8,
        "model": "Intel(R) Core(TM) i7-9700K",
        "frequency": 3600.0
      },
      "memory": {
        "total": 17179869184,
        "available": 8589934592,
        "used": 8589934592
      }
    },
    "network": {
      "hostname": "web-server-01",
      "interfaces": [
        {
          "name": "eth0",
          "ip": "192.168.1.100",
          "mac": "00:15:5d:01:ca:05"
        }
      ]
    }
  }
}
```

## Collector Facts (spooky-collector)

### Overview
Collector facts are provided by the `spooky-collector` binary, which is a deployable agent that collects comprehensive system metrics using gopsutil and outputs them to `/etc/spooky/facts.hcl`.

### Deployment
The spooky-collector binary should be deployed to target machines and configured to run as a service or scheduled task.

### Output Format
The collector outputs facts in HCL format to `/etc/spooky/facts.hcl`:

```hcl
collector {
  host {
    hostname = "web-server-01"
    uptime = 86400
    boot_time = 1640995200
    os = "linux"
    platform = "ubuntu"
    platform_family = "debian"
    platform_version = "22.04"
    kernel_version = "5.15.0-generic"
    kernel_arch = "x86_64"
  }

  cpu {
    cores = 8
    model = "Intel(R) Core(TM) i7-9700K"
    frequency = 3600.0
    architecture = "x86_64"
    vendor = "GenuineIntel"
    percent = 15.2
  }

  memory {
    total = 17179869184
    available = 8589934592
    used = 8589934592
    free = 4294967296
    buffers = 1073741824
    cached = 2147483648
    shared = 536870912
    slab = 268435456
  }

  disks = [
    {
      device = "/dev/sda1"
      mount_point = "/"
      total = 107374182400
      used = 53687091200
      free = 53687091200
      filesystem = "ext4"
    }
  ]

  network {
    hostname = "web-server-01"
    primary_ip = "192.168.1.100"
    bytes_sent = 1073741824
    bytes_recv = 2147483648
    packets_sent = 1000000
    packets_recv = 2000000
  }

  load_average {
    load1 = 0.5
    load5 = 0.8
    load15 = 1.2
  }
}
```

### Error Handling
- **SSH Connection Failure**: Collection fails completely
- **Collector Not Deployed**: Collection fails with clear error message
- **Insufficient Privileges**: Collection fails with privilege error
- **Invalid HCL Output**: Collection fails with parsing error

## Custom Facts (HCL Files)

### Overview
Custom facts allow users to define their own fact data in HCL format. These facts are stored in `/etc/spooky/custom.hcl` on target machines.

### File Format
Custom facts can contain any valid HCL structure:

```hcl
# Application-specific facts
application {
  name = "web-app"
  version = "1.2.3"
  environment = "production"
  port = 8080
}

# Environment facts
environment {
  region = "us-west-2"
  datacenter = "dc-01"
  rack = "rack-03"
  slot = "slot-05"
}

# Business facts
business {
  owner = "engineering-team"
  cost_center = "CC-12345"
  project = "web-infrastructure"
  sla_tier = "gold"
}

# Custom metrics
metrics {
  response_time_avg = 150.5
  error_rate = 0.02
  active_connections = 1250
}
```

### Error Handling
- **SSH Connection Failure**: Collection fails completely
- **File Not Found**: Collection continues without custom facts
- **Invalid HCL**: Collection fails with parsing error

## CLI Usage

### Basic Fact Collection
```bash
# Collect facts from all machines in a project
spooky facts export ./my-project --output facts.hcl

# Collect facts from specific machine
spooky facts export ./my-project --machine web-server --output web-server-facts.hcl

# Collect facts with filtering
spooky facts export ./my-project --tags environment=production --output prod-facts.hcl
```

### Export Formats
```bash
# Export to HCL format (default)
spooky facts export ./my-project --output facts.hcl

# Export to JSON format
spooky facts export ./my-project --format json --output facts.json
```

### Parallel Collection
```bash
# Collect facts with parallel processing
spooky facts export ./my-project --parallel 4 --output facts.hcl
```

## Error Handling Strategy

### System Facts
- **SSH Connection Failure**: Fail completely
- **Basic Command Failure**: Continue with partial data
- **Permission Issues**: Skip failed commands, continue collection

### Collector Facts
- **SSH Connection Failure**: Fail completely
- **Collector Not Deployed**: Fail with deployment error
- **Insufficient Privileges**: Fail with privilege error
- **Invalid HCL Output**: Fail with parsing error

### Custom Facts
- **SSH Connection Failure**: Fail completely
- **File Not Found**: Continue without custom facts (optional)
- **Invalid HCL**: Fail with parsing error

## Configuration

### SSH Configuration
SSH connections are configured through the standard spooky SSH management system:

```hcl
# machines.hcl
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
```

### Fact Collection Configuration
Fact collection can be configured through project settings:

```hcl
# project.hcl
project {
  name = "my-project"
  
  settings {
    fact_collection {
      timeout = 300
      parallel_workers = 4
      retry_attempts = 3
    }
  }
}
```

## Security Considerations

### User Privileges
- System facts collection uses only user-level privileges
- No sudo or elevated permissions required
- Reduces security risk and simplifies deployment

### SSH Security
- All SSH connections use standard SSH security features
- Key-based authentication recommended
- Connection pooling for efficiency

### Data Privacy
- Facts are collected in memory and exported directly
- No persistent storage of sensitive data
- Custom facts allow users to control what data is collected

## Troubleshooting

### Common Issues

#### SSH Connection Failures
```bash
# Test SSH connectivity
spooky machines ping ./my-project --machine web-server

# Check SSH configuration
ssh -i ~/.ssh/id_rsa admin@web.example.com
```

#### Collector Deployment Issues
```bash
# Check if collector is installed
ssh admin@web.example.com "ls -la /etc/spooky/facts.hcl"

# Check collector service status
ssh admin@web.example.com "systemctl status spooky-collector"
```

#### HCL Parsing Errors
```bash
# Validate HCL syntax
ssh admin@web.example.com "hclfmt -check /etc/spooky/custom.hcl"

# Test HCL parsing
ssh admin@web.example.com "cat /etc/spooky/custom.hcl | hclfmt"
```

### Debug Mode
Enable verbose logging for troubleshooting:

```bash
spooky facts export ./my-project --verbose --output facts.hcl
```

## Performance Considerations

### Parallel Collection
- Use `--parallel` flag to collect facts from multiple machines simultaneously
- Default is 1, recommended maximum is 8-16 depending on network capacity

### Timeout Settings
- Default timeout is 30 seconds per machine
- Adjust based on network latency and machine performance

### Caching
- Facts are collected fresh each time
- Consider implementing caching for large environments

## Future Enhancements

### Planned Features
- **Fact Caching**: Cache collected facts to improve performance
- **Incremental Collection**: Only collect changed facts
- **Fact Validation**: Validate facts against schemas
- **Fact Dependencies**: Support for fact dependencies and relationships

### Integration Opportunities
- **Monitoring Systems**: Export facts to Prometheus, Grafana, etc.
- **Configuration Management**: Use facts for dynamic configuration
- **Compliance**: Use facts for compliance reporting and auditing
