# Machine Inventory Examples

This directory contains comprehensive examples of machine inventory configurations for the spooky machines system. These examples demonstrate various use cases, best practices, and configuration patterns.

## Examples Overview

### 1. Basic Inventory (`machines-basic-inventory.hcl`)
A simple example showing fundamental machine configuration with basic web servers, database servers, and load balancers.

**Use Case:** Getting started with spooky machines inventory
**Complexity:** Beginner
**Machines:** 3 (web server, database server, load balancer)

**Key Features:**
- Basic machine configuration
- Resource specifications
- Metadata organization
- Tag and group usage

### 2. Multi-Environment Inventory (`machines-multi-environment.hcl`)
A comprehensive example showing how to manage machines across different environments (production, staging, development).

**Use Case:** Multi-environment infrastructure management
**Complexity:** Intermediate
**Machines:** 7 (production, staging, development)

**Key Features:**
- Environment-specific configurations
- Different resource requirements per environment
- Environment-specific validation rules
- Cost center organization

### 3. Kubernetes Nodes (`machines-kubernetes-nodes.hcl`)
An example showing how to configure machines for Kubernetes cluster management, including control plane and worker nodes.

**Use Case:** Kubernetes infrastructure management
**Complexity:** Advanced
**Machines:** 7 (masters, workers, infrastructure)

**Key Features:**
- Kubernetes-specific metadata
- Node roles and labels
- GPU node configuration
- Infrastructure node separation

## Using the Examples

### Copy and Customize
```bash
# Copy an example to your project
cp docs/examples/machines/machines-basic-inventory.hcl ./my-project/machines.hcl

# Customize the configuration for your environment
# Edit hostnames, IP addresses, users, and SSH keys
```

### Test the Examples
```bash
# Test basic inventory
spooky machines validate ./my-project
spooky machines list ./my-project
spooky machines ping ./my-project

# Test with verbose output
spooky machines list ./my-project --verbose
spooky machines ping ./my-project --verbose

# Test JSON output
spooky machines ping ./my-project --format json
```

### Multi-File Organization
For larger environments, organize machines into multiple files:

```bash
# Create machines directory structure
mkdir -p ./my-project/machines

# Copy environment-specific files
cp docs/examples/machines/machines-multi-environment.hcl ./my-project/machines/production.hcl
cp docs/examples/machines/machines-kubernetes-nodes.hcl ./my-project/machines/kubernetes.hcl

# Create your own environment files
touch ./my-project/machines/staging.hcl
touch ./my-project/machines/development.hcl
```

## Configuration Patterns

### Basic Machine Configuration
```hcl
machine "server-name" {
  host = "192.168.1.10"
  user = "admin"
  port = 22
  key_file = "~/.ssh/id_rsa"
  
  tags = ["web", "production"]
  groups = ["web-servers"]
  roles = ["web-server", "nginx"]
}
```

### Resource Specifications
```hcl
resources {
  cpu_cores = 8
  memory_gb = 32
  disk_gb = 500
  network_mbps = 10000
}
```

### Comprehensive Metadata
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

### Environment-Specific Configuration
```hcl
# Production - High resources, strict validation
machine "prod-web-01" {
  host = "10.0.1.10"
  user = "admin"
  key_file = "~/.ssh/prod_key"
  
  resources {
    cpu_cores = 8
    memory_gb = 32
    disk_gb = 500
  }
  
  metadata {
    environment = "production"
    backup_schedule = "daily"
    cost_center = "IT-001"
  }
}

# Development - Lower resources, lenient validation
machine "dev-web-01" {
  host = "10.0.3.10"
  user = "developer"
  key_file = "~/.ssh/dev_key"
  
  metadata {
    environment = "development"
    owner = "developer"
    purpose = "development and testing"
  }
}
```

## Best Practices Demonstrated

### 1. Naming Conventions
- Use descriptive, consistent hostnames
- Include environment prefix: `prod-web-01`, `staging-db-01`
- Use lowercase with hyphens: `web-server-01`

### 2. Organization
- Group related machines together
- Use consistent tag names across environments
- Separate environments into different files

### 3. Security
- Use dedicated SSH keys for different environments
- Store keys securely with proper permissions (600)
- Use passphrases for additional security

### 4. Documentation
- Include comprehensive metadata
- Document maintenance windows and backup schedules
- Specify monitoring and alerting systems

### 5. Resource Planning
- Specify resource requirements for capacity planning
- Use different resource levels for different environments
- Document special requirements (GPU, high memory, etc.)

## Validation Examples

### Production Environment Validation
```bash
# Production machines should have:
# - Resource specifications
# - Backup schedule
# - Cost center information
# - Monitoring configuration

spooky machines validate ./my-project
# Should show warnings for missing production requirements
```

### Development Environment Validation
```bash
# Development machines have more lenient validation:
# - Owner information required
# - Resource specifications optional
# - Backup schedule optional

spooky machines validate ./my-project
# Should show fewer warnings for development machines
```

## Troubleshooting Examples

### Common Issues and Solutions

#### 1. Duplicate Hostnames
```bash
# Error: duplicate hostname 'web-server' found in multiple files
# Solution: Use unique hostnames across all files
machine "prod-web-01" { ... }  # Instead of "web-server"
machine "staging-web-01" { ... }  # Instead of "web-server"
```

#### 2. Missing Authentication
```bash
# Error: machine 'web-server' missing authentication method
# Solution: Add key_file attribute
machine "web-server" {
  host = "192.168.1.10"
  user = "admin"
  key_file = "~/.ssh/id_rsa"  # Add this line
}
```

#### 3. Invalid SSH Key Path
```bash
# Error: SSH key file '/path/to/key' not found
# Solution: Verify key file exists and has correct permissions
ls -la ~/.ssh/id_rsa
chmod 600 ~/.ssh/id_rsa
```

## Advanced Patterns

### Kubernetes Node Management
```hcl
# Control plane nodes
machine "k8s-master-01" {
  host = "10.0.10.10"
  user = "kubernetes"
  key_file = "~/.ssh/k8s_master_key"
  
  tags = ["kubernetes", "control-plane", "production"]
  roles = ["kubernetes-master", "etcd", "api-server"]
  
  metadata {
    kubernetes_version = "1.28.0"
    node_role = "master"
    etcd_member = "true"
  }
}

# GPU worker nodes
machine "k8s-worker-gpu-01" {
  host = "10.0.11.12"
  user = "kubernetes"
  key_file = "~/.ssh/k8s_worker_key"
  
  tags = ["kubernetes", "worker", "gpu", "production"]
  roles = ["kubernetes-worker", "gpu-accelerator"]
  
  resources {
    cpu_cores = 32
    memory_gb = 128
    disk_gb = 2000
  }
  
  metadata {
    kubernetes_version = "1.28.0"
    node_role = "worker"
    node_labels = "node-type=gpu,gpu-type=nvidia-v100"
    taints = "nvidia.com/gpu=true:NoSchedule"
    gpu_count = "4"
    gpu_type = "nvidia-v100"
  }
}
```

### Load Balancer Configuration
```hcl
machine "lb-primary" {
  host = "10.0.1.100"
  user = "admin"
  key_file = "~/.ssh/lb_key"
  
  tags = ["load-balancer", "production", "primary"]
  roles = ["load-balancer", "haproxy", "ssl-terminator"]
  
  resources {
    cpu_cores = 4
    memory_gb = 16
    disk_gb = 200
    network_mbps = 10000
  }
  
  metadata {
    environment = "production"
    sla = "99.99%"
    monitoring = "prometheus"
    alerting = "pagerduty"
  }
}
```

## Testing and Validation

### Automated Testing Script
```bash
#!/bin/bash
# test-machines.sh

PROJECT_PATH="./my-project"

echo "=== Testing Machine Inventory ==="

# Test validation
echo "1. Validating configuration..."
if spooky machines validate "$PROJECT_PATH"; then
    echo "✓ Validation passed"
else
    echo "✗ Validation failed"
    exit 1
fi

# Test listing
echo "2. Testing machine listing..."
if spooky machines list "$PROJECT_PATH"; then
    echo "✓ Listing successful"
else
    echo "✗ Listing failed"
    exit 1
fi

# Test connectivity
echo "3. Testing connectivity..."
if spooky machines ping "$PROJECT_PATH" --format json | jq -e '.status == "online"' > /dev/null; then
    echo "✓ Connectivity test passed"
else
    echo "⚠ Connectivity test failed (some machines offline)"
fi

echo "=== Testing completed ==="
```

### JSON Output Processing
```bash
# Get machine status in JSON format
spooky machines ping ./my-project --format json > machines-status.json

# Filter online machines
jq 'select(.status == "online")' machines-status.json

# Filter by tags
jq 'select(.machine.tags[]? | contains("production"))' machines-status.json

# Get machine count by environment
jq -r '.machine.metadata.environment // "unknown"' machines-status.json | sort | uniq -c
```

## Next Steps

1. **Start with Basic Inventory**: Use `machines-basic-inventory.hcl` to get familiar with the format
2. **Organize by Environment**: Create separate files for different environments
3. **Add Resource Specifications**: Include resource requirements for capacity planning
4. **Implement Validation**: Use validation to ensure configuration quality
5. **Set Up Monitoring**: Configure monitoring and alerting metadata
6. **Document Everything**: Add comprehensive metadata for better management

## Contributing

When adding new examples:

1. **Follow Naming Conventions**: Use consistent naming patterns
2. **Include Comments**: Add helpful comments explaining configuration choices
3. **Test Examples**: Ensure examples work with current spooky version
4. **Document Use Cases**: Explain when and why to use each example
5. **Follow Best Practices**: Demonstrate recommended patterns and practices

## Support

For help with machine inventory configuration:

1. **Check Documentation**: Review the main user guide and API reference
2. **Use Validation**: Run `spooky machines validate` to catch configuration issues
3. **Enable Verbose Mode**: Use `--verbose` flag for detailed output
4. **Check Examples**: Use these examples as starting points for your configuration
5. **Review Troubleshooting Guide**: Check the troubleshooting guide for common issues
