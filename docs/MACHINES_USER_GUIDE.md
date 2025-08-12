# Machines Inventory User Guide

## Overview

The spooky machines inventory system provides comprehensive management of remote machines for automation and orchestration. This guide covers everything from basic machine configuration to advanced features like multi-file inventories, connectivity testing, and validation.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Machine Configuration](#machine-configuration)
3. [Inventory Management](#inventory-management)
4. [Connectivity Testing](#connectivity-testing)
5. [Validation and Troubleshooting](#validation-and-troubleshooting)
6. [Advanced Features](#advanced-features)
7. [Best Practices](#best-practices)
8. [Examples](#examples)

## Getting Started

### Basic Machine Inventory

A machine inventory in spooky consists of one or more HCL files that define remote machines with their connection details, authentication methods, and metadata.

**Single File Inventory (`machines.hcl`):**
```hcl
machines {
  machine "web-server-01" {
    host = "192.168.1.10"
    user = "admin"
    port = 22
    
    key_file = "~/.ssh/id_rsa"
    passphrase = "my-secure-passphrase"
    
    tags = ["web", "production"]
    groups = ["web-servers"]
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      owner = "web-team"
    }
  }
  
  machine "db-server-01" {
    host = "db.example.com"
    user = "dbadmin"
    port = 2222
    
    key_file = "~/.ssh/db_key"
    
    tags = ["database", "production"]
    groups = ["database-servers"]
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      owner = "db-team"
    }
  }
}
```

### Multi-File Inventory

For larger environments, you can organize machines into multiple files within a `machines/` directory:

**Project Structure:**
```
my-project/
├── project.hcl
├── machines.hcl                    # Global machines
└── machines/
    ├── production.hcl              # Production machines
    ├── staging.hcl                 # Staging machines
    └── development.hcl             # Development machines
```

**Production Machines (`machines/production.hcl`):**
```hcl
machines {
  machine "prod-web-01" {
    host = "10.0.1.10"
    user = "admin"
    key_file = "~/.ssh/prod_key"
    
    tags = ["web", "production"]
    groups = ["web-servers"]
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      cost_center = "IT-001"
    }
  }
  
  machine "prod-db-01" {
    host = "10.0.1.20"
    user = "dbadmin"
    key_file = "~/.ssh/prod_db_key"
    
    tags = ["database", "production"]
    groups = ["database-servers"]
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      cost_center = "IT-002"
    }
  }
}
```

**Staging Machines (`machines/staging.hcl`):**
```hcl
machines {
  machine "staging-web-01" {
    host = "10.0.2.10"
    user = "admin"
    key_file = "~/.ssh/staging_key"
    
    tags = ["web", "staging"]
    groups = ["web-servers"]
    
    metadata {
      environment = "staging"
      datacenter = "us-west-1"
    }
  }
}
```

## Machine Configuration

### Required Fields

Every machine definition must include these required fields:

- **`hostname`** (label): Unique identifier for the machine
- **`host`**: IP address or hostname for SSH connection
- **`user`**: SSH username for authentication

### Optional Fields

- **`port`**: SSH port (default: 22)
- **`key_file`**: Path to SSH private key file
- **`passphrase`**: Passphrase for encrypted SSH keys
- **`tags`**: Key-value pairs for categorization
- **`groups`**: List of groups for organization
- **`roles`**: List of roles for automation
- **`resources`**: Machine resource specifications
- **`metadata`**: Additional machine metadata

### Authentication Methods

Spooky supports SSH key-based authentication:

```hcl
machine "secure-server" {
  host = "192.168.1.100"
  user = "admin"
  
  # SSH key authentication (recommended)
  key_file = "~/.ssh/id_rsa"
  passphrase = "my-secure-passphrase"
}
```

**Note:** Password authentication is not supported for security reasons.

### Resource Specifications

Define machine resources for capacity planning:

```hcl
machine "high-performance-server" {
  host = "192.168.1.200"
  user = "admin"
  key_file = "~/.ssh/server_key"
  
  resources {
    cpu_cores = 16
    memory_gb = 64
    disk_gb = 1000
    network_mbps = 10000
  }
}
```

### Metadata and Organization

Use metadata for better organization and management:

```hcl
machine "web-server" {
  host = "192.168.1.10"
  user = "admin"
  key_file = "~/.ssh/web_key"
  
  tags = ["web", "production", "load-balanced"]
  groups = ["web-servers", "production-servers"]
  roles = ["web-server", "nginx", "ssl-terminator"]
  
  metadata {
    environment = "production"
    datacenter = "us-west-1"
    rack = "A-01"
    location = "San Francisco"
    owner = "web-team"
    department = "Engineering"
    cost_center = "IT-001"
    maintenance_window = "Sunday 2-4 AM PST"
    backup_schedule = "daily"
  }
}
```

## Inventory Management

### Listing Machines

List all machines in a project:

```bash
# List all machines
spooky machines list ./my-project

# List with verbose output
spooky machines list ./my-project --verbose

# List machines by tags
spooky machines list ./my-project --tags "production,web"

# List machines by groups
spooky machines list ./my-project --groups "web-servers"
```

**Example Output:**
```
🔍 Loading machines from project: ./my-project
📊 Found 5 machines:

📁 Source: machines/production.hcl (3 machines)
──────────────────────────────────────────────────
1. prod-web-01 (10.0.1.10)
   User: admin
   Port: 22
   Environment: production
   Groups: [web-servers]
   Tags: [web production]

2. prod-db-01 (10.0.1.20)
   User: dbadmin
   Port: 22
   Environment: production
   Groups: [database-servers]
   Tags: [database production]

📁 Source: machines/staging.hcl (2 machines)
──────────────────────────────────────────────────
3. staging-web-01 (10.0.2.10)
   User: admin
   Port: 22
   Environment: staging
   Groups: [web-servers]
   Tags: [web staging]
```

### Validating Machine Inventory

Validate machine configuration and detect issues:

```bash
# Basic validation
spooky machines validate ./my-project

# Verbose validation with details
spooky machines validate ./my-project --verbose
```

**Example Output:**
```
🔍 Validating machines in project: ./my-project
✅ Machine validation completed successfully

📊 Validation Summary:
- Total machines: 5
- Valid machines: 5
- Warnings: 2
- Errors: 0

⚠️  Warnings:
- prod-web-01: Missing resource specifications (recommended for production)
- staging-web-01: No backup schedule specified in metadata
```

## Connectivity Testing

### Basic Connectivity Testing

Test connectivity to all machines:

```bash
# Test all machines
spooky machines ping ./my-project

# Test specific machine
spooky machines ping ./my-project --machine "prod-web-01"

# Test machines by tags
spooky machines ping ./my-project --tags "production"

# Test with verbose output
spooky machines ping ./my-project --verbose
```

**Example Output (Smart Mode):**
```
🔍 Pinging machines in ./my-project
📊 Ping Results: Total machines: 5

✅ prod-web-01: online (12ms)
✅ prod-db-01: online (8ms)
❌ staging-web-01: offline (DNS resolution failed)
✅ staging-db-01: online (15ms)
✅ dev-web-01: online (25ms)
```

**Example Output (Verbose Mode):**
```
🔍 Pinging machines in ./my-project
📊 Ping Results: Total machines: 5

✅ prod-web-01 (10.0.1.10): online (12ms)
   DNS: resolved successfully
   ICMP: reachable
   Port 22: open
   SSH: authentication successful

✅ prod-db-01 (10.0.1.20): online (8ms)
   DNS: resolved successfully
   ICMP: reachable
   Port 22: open
   SSH: authentication successful

❌ staging-web-01 (10.0.2.10): offline
   DNS: resolution failed (NXDOMAIN)
   ICMP: unreachable
   Port 22: closed
   SSH: connection failed

✅ staging-db-01 (10.0.2.20): online (15ms)
   DNS: resolved successfully
   ICMP: reachable
   Port 22: open
   SSH: authentication successful

✅ dev-web-01 (10.0.3.10): online (25ms)
   DNS: resolved successfully
   ICMP: reachable
   Port 22: open
   SSH: authentication successful
```

### JSON Output

Get machine status in JSON format for scripting:

```bash
# JSON output (streaming)
spooky machines ping ./my-project --format json

# JSON output with verbose details
spooky machines ping ./my-project --format json --verbose
```

**Example JSON Output:**
```json
{"hostname":"prod-web-01","status":"online"}
{"hostname":"prod-db-01","status":"online"}
{"hostname":"staging-web-01","status":"offline","error":"DNS resolution failed"}
{"hostname":"staging-db-01","status":"online"}
{"hostname":"dev-web-01","status":"online"}
```

## Validation and Troubleshooting

### Common Validation Issues

**1. Duplicate Hostnames**
```
❌ Error: duplicate hostname 'web-server-01' found in multiple files: 
   [machines/production.hcl machines/staging.hcl]
```
**Solution:** Ensure each hostname is unique across all inventory files.

**2. Missing Authentication**
```
❌ Error: machine 'web-server-01' missing authentication method
```
**Solution:** Add either `key_file` or ensure SSH agent has the key loaded.

**3. Invalid SSH Key Path**
```
❌ Error: SSH key file '/path/to/nonexistent/key' not found
```
**Solution:** Verify the key file path and permissions.

**4. Network Connectivity Issues**
```
❌ Error: DNS resolution failed for 'invalid-hostname.example.com'
```
**Solution:** Check hostname/IP address and DNS configuration.

### Troubleshooting Connectivity

**1. DNS Resolution Issues**
```bash
# Test DNS resolution manually
nslookup web-server-01
dig web-server-01

# Check if using IP address instead of hostname
spooky machines ping ./my-project --machine "192.168.1.10"
```

**2. SSH Authentication Issues**
```bash
# Test SSH connection manually
ssh -i ~/.ssh/id_rsa admin@192.168.1.10

# Check SSH key permissions
ls -la ~/.ssh/id_rsa
chmod 600 ~/.ssh/id_rsa

# Test with SSH agent
ssh-add ~/.ssh/id_rsa
ssh admin@192.168.1.10
```

**3. Network Connectivity Issues**
```bash
# Test basic connectivity
ping 192.168.1.10

# Test SSH port
telnet 192.168.1.10 22
nc -zv 192.168.1.10 22
```

## Advanced Features

### Environment-Specific Validation

Spooky automatically applies different validation rules based on the environment:

**Production Environment Rules:**
- Requires resource specifications
- Requires backup schedule in metadata
- Requires cost center information
- Recommends key-based authentication
- Enforces stricter timeout settings

**Development Environment Rules:**
- More lenient validation
- Allows missing resource specifications
- Allows missing metadata fields

### Cross-File Consistency

Spooky validates consistency across multiple inventory files:

**Authentication Consistency:**
- Recommends consistent authentication methods within environments
- Warns about mixed authentication methods

**Timeout Consistency:**
- Checks for consistent timeout settings within environments
- Warns about varying timeout configurations

### Duplicate Detection

Spooky detects and reports duplicates:

**Hostname Duplicates:**
- Prevents duplicate hostnames across all files
- Reports file sources for duplicates

**Host Address Duplicates:**
- Warns about multiple machines using the same IP/host
- Reports affected hostnames and file sources

## Best Practices

### 1. Organization

**Use Multi-File Inventories:**
- Separate environments into different files
- Use descriptive file names
- Group related machines together

**Example Structure:**
```
machines/
├── production/
│   ├── web-servers.hcl
│   ├── database-servers.hcl
│   └── load-balancers.hcl
├── staging/
│   ├── web-servers.hcl
│   └── database-servers.hcl
└── development/
    └── all-servers.hcl
```

### 2. Naming Conventions

**Hostname Conventions:**
- Use descriptive, consistent naming
- Include environment prefix: `prod-web-01`, `staging-db-01`
- Use lowercase with hyphens: `web-server-01`
- Avoid special characters except hyphens

**Tag Conventions:**
- Use consistent tag names across environments
- Use lowercase with hyphens: `web-server`, `production`
- Group related tags: `environment:production`, `role:web`

### 3. Security

**SSH Key Management:**
- Use dedicated SSH keys for different environments
- Store keys securely with proper permissions (600)
- Use passphrases for additional security
- Rotate keys regularly

**Access Control:**
- Use least-privilege user accounts
- Document access requirements in metadata
- Review and audit access regularly

### 4. Documentation

**Metadata Best Practices:**
- Include owner and department information
- Document maintenance windows
- Specify backup schedules
- Include cost center information

**Example Comprehensive Metadata:**
```hcl
metadata {
  environment = "production"
  datacenter = "us-west-1"
  rack = "A-01"
  location = "San Francisco"
  owner = "web-team"
  department = "Engineering"
  cost_center = "IT-001"
  maintenance_window = "Sunday 2-4 AM PST"
  backup_schedule = "daily"
  monitoring = "prometheus"
  alerting = "pagerduty"
  sla = "99.9%"
}
```

### 5. Validation

**Regular Validation:**
- Run validation before deployments
- Include validation in CI/CD pipelines
- Monitor for configuration drift
- Review warnings and address issues

**Automated Checks:**
```bash
# Pre-deployment validation
spooky machines validate ./my-project

# Connectivity testing
spooky machines ping ./my-project --tags "production"

# Configuration drift detection
spooky machines validate ./my-project --compare
```

## Examples

### Complete Production Environment

**`machines/production.hcl`:**
```hcl
machines {
  machine "prod-web-01" {
    host = "10.0.1.10"
    user = "admin"
    port = 22
    key_file = "~/.ssh/prod_web_key"
    passphrase = "secure-passphrase"
    
    tags = ["web", "production", "load-balanced"]
    groups = ["web-servers", "production-servers"]
    roles = ["web-server", "nginx", "ssl-terminator"]
    
    resources {
      cpu_cores = 8
      memory_gb = 32
      disk_gb = 500
      network_mbps = 10000
    }
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      rack = "A-01"
      location = "San Francisco"
      owner = "web-team"
      department = "Engineering"
      cost_center = "IT-001"
      maintenance_window = "Sunday 2-4 AM PST"
      backup_schedule = "daily"
      monitoring = "prometheus"
      alerting = "pagerduty"
      sla = "99.9%"
    }
  }
  
  machine "prod-db-01" {
    host = "10.0.1.20"
    user = "dbadmin"
    port = 22
    key_file = "~/.ssh/prod_db_key"
    
    tags = ["database", "production", "primary"]
    groups = ["database-servers", "production-servers"]
    roles = ["database-server", "postgresql", "primary"]
    
    resources {
      cpu_cores = 16
      memory_gb = 64
      disk_gb = 2000
      network_mbps = 10000
    }
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      rack = "A-02"
      location = "San Francisco"
      owner = "db-team"
      department = "Engineering"
      cost_center = "IT-002"
      maintenance_window = "Sunday 2-4 AM PST"
      backup_schedule = "hourly"
      monitoring = "prometheus"
      alerting = "pagerduty"
      sla = "99.99%"
    }
  }
}
```

### Development Environment

**`machines/development.hcl`:**
```hcl
machines {
  machine "dev-web-01" {
    host = "10.0.3.10"
    user = "developer"
    port = 22
    key_file = "~/.ssh/dev_key"
    
    tags = ["web", "development"]
    groups = ["web-servers", "development-servers"]
    roles = ["web-server", "nginx"]
    
    metadata {
      environment = "development"
      datacenter = "us-west-1"
      owner = "developer"
      department = "Engineering"
      purpose = "web development and testing"
    }
  }
  
  machine "dev-db-01" {
    host = "10.0.3.20"
    user = "developer"
    port = 22
    key_file = "~/.ssh/dev_key"
    
    tags = ["database", "development"]
    groups = ["database-servers", "development-servers"]
    roles = ["database-server", "postgresql"]
    
    metadata {
      environment = "development"
      datacenter = "us-west-1"
      owner = "developer"
      department = "Engineering"
      purpose = "database development and testing"
    }
  }
}
```

### Load Balancer Configuration

**`machines/load-balancers.hcl`:**
```hcl
machines {
  machine "lb-primary" {
    host = "10.0.1.100"
    user = "admin"
    port = 22
    key_file = "~/.ssh/lb_key"
    
    tags = ["load-balancer", "production", "primary"]
    groups = ["load-balancers", "production-servers"]
    roles = ["load-balancer", "haproxy", "ssl-terminator"]
    
    resources {
      cpu_cores = 4
      memory_gb = 16
      disk_gb = 200
      network_mbps = 10000
    }
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      rack = "A-00"
      location = "San Francisco"
      owner = "infrastructure-team"
      department = "Engineering"
      cost_center = "IT-000"
      maintenance_window = "Sunday 2-4 AM PST"
      backup_schedule = "daily"
      monitoring = "prometheus"
      alerting = "pagerduty"
      sla = "99.99%"
    }
  }
}
```

This comprehensive user guide provides everything needed to effectively use the spooky machines inventory system, from basic configuration to advanced features and best practices.
