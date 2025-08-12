# Machines System: Comprehensive Implementation Plan

## Overview

This document is the authoritative source for all machines system implementation details in spooky. It covers machine inventory management, connectivity operations, enterprise-scale indexing, import/export capabilities, and integration with all other spooky systems.

**Schema Integration**: This machines system implements the schema validation patterns and machine configuration definitions defined in [Schema System](../schema-system.md) for comprehensive machine validation, inventory schema enforcement, and schema-based machine lifecycle management.

**Architecture Integration**: Machines integrate with the overall spooky architecture as described in [Spooky Design](../spooky-design.md), providing the foundational inventory management that all other systems depend on for target selection and running.

## System Integration

This machines system integrates with other core Spooky systems to provide comprehensive inventory management and connectivity operations:

### **Actions System Integration**
- **Target Selection**: Actions use machine inventory for target selection via tags, names, or filters (see [Actions System](../actions-system.md))
- **Machine Resolution**: Actions resolve machine targets through the machines system's indexing and lookup capabilities
- **Execution Context**: Actions run within machine-specific contexts with authentication and connection details
- **Parallel Running**: Machine inventory supports parallel action running across multiple targets

### **Facts System Integration**
- **Machine Identification**: Facts collection uses machine inventory for target identification and authentication (see [Facts System](../facts-system.md))
- **Machine Metadata**: Machine inventory provides metadata (hostname, tags, authentication) for facts collection
- **Facts Storage**: Machine facts are stored using machine IDs and names from inventory
- **Facts Association**: Collected facts are associated with specific machines from inventory

### **Project System Integration**
- **Project Inventory**: Machine inventory is stored in project-specific `machines.hcl` files (see [Project System](../project-system.md))
- **Project Isolation**: Each project maintains its own machine inventory with project-specific settings
- **Project Context**: Machine operations run within project context and configuration
- **Project Validation**: Machine inventory validation integrated with project validation

### **CLI System Integration**
- **Machine Commands**: Machine management through `spooky machines` commands (see [CLI System](../cli-system.md))
- **Command Patterns**: Machine commands follow the established `spooky noun verb` CLI pattern
- **Validation Commands**: `spooky machines validate` for inventory validation
- **Management Commands**: `spooky machines list`, `spooky machines ping`, `spooky machines connect`
- **Import/Export Commands**: `spooky machines import`, `spooky machines export` for external system integration

### **Configuration System Integration**
- **Global Configuration**: Machine settings inherit from global configuration (see [Configuration System](../configuration-system.md))
- **SSH Configuration**: Default SSH settings (user, port, timeout) from global config
- **Authentication Settings**: SSH key paths, connection timeouts, retry policies
- **Performance Settings**: Connection pooling, parallel connection limits

### **Variables System Integration**
- **Machine Variables**: Machine inventory can use project variables for dynamic configuration (see [Variables System](../variables-system.md))
- **Variable Interpolation**: Machine hostnames, usernames, and tags can use variable interpolation
- **Environment Variables**: Machine configuration can reference environment variables
- **Dynamic Inventory**: Machine inventory can be generated from variables and external sources

### **Template System Integration**
- **Machine Context**: Templates can access machine-specific data and facts (see [Template System](../template-system.md))
- **Machine Functions**: Template functions for accessing machine data
- **Remote Rendering**: Templates can be rendered on remote machines via SSH
- **Machine-Specific Templates**: Templates can be customized per machine
- **Template Deployment**: Machine inventory supports template deployment workflows

### **Schema System Integration**
- **Machine Schema**: Machine inventory validation against embedded schemas (see [Schema System](../schema-system.md))
- **Schema Validation**: Machine configuration validation using schema system
- **Schema Evolution**: Machine schema versioning and migration support
- **Schema Composition**: Runtime schema composition for machine validation
- **Schema Integration**: Machine system integration with all system schemas

## Current State Analysis

### **What We Have**
- ✅ **Basic machine configuration** with HCL parsing and validation
- ✅ **Machine authentication** support for SSH keys and passwords
- ✅ **Machine tagging** for organization and filtering
- ✅ **SSH connectivity** testing and validation
- ✅ **Machine CLI commands** for list, validate, ping, connect
- ✅ **Enterprise-scale indexing** for large machine inventories
- ✅ **Import/export capabilities** for external system integration

### **What We Need**
- 🔄 **Enhanced machine schema** with comprehensive validation
- 🔄 **Machine dependency tracking** for complex deployments
- 🔄 **Machine health monitoring** and status tracking
- 🔄 **Machine grouping** and hierarchical organization
- 🔄 **Machine discovery** and auto-registration
- 🔄 **Machine backup** and recovery procedures

## Machines System Design

### **1. Machine Configuration Schema**

Based on the schema defined in `internal/schemas/schemas/machines.hcl`, machine configurations follow this structure:

```hcl
# machines.hcl
machines {
  machine "machine_name" {
    host = "192.168.1.100"           # Required: IP address or hostname
    port = 22                        # Optional: SSH port (default: 22)
    user = "admin"                   # Required: SSH username
    password = "secure-password"     # Optional: SSH password
    key_file = "~/.ssh/id_rsa"       # Optional: SSH private key file
    tags = {                         # Optional: Machine tags
      environment = "production"
      role = "web-server"
      datacenter = "primary"
    }
  }
}
```

### **2. Machine Authentication**

**SSH Key Authentication** (Recommended):
```hcl
machine "web-server" {
  host = "192.168.1.100"
  user = "admin"
  key_file = "~/.ssh/production_key"
  tags = {
    environment = "production"
    role = "web"
  }
}
```

**Password Authentication** (Less Secure):
```hcl
machine "dev-server" {
  host = "192.168.1.200"
  user = "developer"
  password = "dev-password"
  tags = {
    environment = "development"
    role = "application"
  }
}
```

**Mixed Authentication** (Different machines, different methods):
```hcl
machines {
  # Production servers with SSH keys
  machine "prod-web-01" {
    host = "192.168.1.10"
    user = "admin"
    key_file = "~/.ssh/prod_key"
    tags = { environment = "production", role = "web" }
  }
  
  # Development servers with passwords
  machine "dev-web-01" {
    host = "192.168.1.20"
    user = "dev"
    password = "dev-password"
    tags = { environment = "development", role = "web" }
  }
}
```

### **3. Machine Tagging and Organization**

**Environment-Based Tagging**:
```hcl
machines {
  machine "web-prod-01" {
    host = "192.168.1.10"
    user = "admin"
    key_file = "~/.ssh/prod_key"
    tags = {
      environment = "production"
      role = "web"
      tier = "frontend"
      region = "us-east"
    }
  }
  
  machine "db-prod-01" {
    host = "192.168.1.11"
    user = "admin"
    key_file = "~/.ssh/prod_key"
    tags = {
      environment = "production"
      role = "database"
      tier = "backend"
      region = "us-east"
    }
  }
}
```

**Role-Based Tagging**:
```hcl
machines {
  machine "web-01" {
    host = "192.168.1.10"
    user = "admin"
    tags = { role = "web", service = "nginx" }
  }
  
  machine "db-01" {
    host = "192.168.1.11"
    user = "admin"
    tags = { role = "database", service = "postgresql" }
  }
  
  machine "cache-01" {
    host = "192.168.1.12"
    user = "admin"
    tags = { role = "cache", service = "redis" }
  }
}
```

### **4. Enterprise-Scale Machine Management**

**Large Inventory Support**:
```hcl
# Support for thousands of machines with efficient indexing
machines {
  # Web servers (1000+ machines)
  machine "web-001" {
    host = "10.0.1.1"
    user = "admin"
    tags = { role = "web", region = "us-east", instance = "1" }
  }
  
  machine "web-002" {
    host = "10.0.1.2"
    user = "admin"
    tags = { role = "web", region = "us-east", instance = "2" }
  }
  
  # Database servers (500+ machines)
  machine "db-001" {
    host = "10.0.2.1"
    user = "admin"
    tags = { role = "database", region = "us-east", instance = "1" }
  }
  
  # ... thousands more machines
}
```

**Performance Optimizations**:
- **Composite Indexing**: Multi-level indexing for O(1) lookups
- **Tag-Based Caching**: Cached tag indexes for frequent queries
- **Memory Optimization**: Efficient memory usage for large inventories
- **Parallel Processing**: Concurrent machine operations

### **5. Machine Import/Export Capabilities**

**Kubernetes Import**:
```bash
# Import machines from Kubernetes
spooky machines import --from kubernetes k8s-nodes.json ./my-project
```

**Terraform State Import**:
```bash
# Import machines from Terraform state
spooky machines import --from terraform terraform.tfstate ./my-project
```

**CMDB Import**:
```bash
# Import machines from ServiceNow CMDB
spooky machines import --from servicenow cmdb-export.json ./my-project
```

**Export to External Systems**:
```bash
# Export machines to JSON
spooky machines export ./my-project --format json --output machines.json

# Export machines to HCL
spooky machines export ./my-project --format hcl --output machines.hcl
```

### **6. Machine Connectivity and Health**

**Connectivity Testing**:
```bash
# Test connectivity to all machines
spooky machines ping ./my-project

# Test specific machine
spooky machines ping ./my-project --machine web-01

# Test machines by tags
spooky machines ping ./my-project --tags "production,web"
```

## CLI Integration

### **1. Machine Management Commands**

**List Machines**:
```bash
# List all machines
spooky machines list ./my-project

# List with details
spooky machines list ./my-project --verbose

# List by tags
spooky machines list ./my-project --tags "production,web"

# List with filtering
spooky machines list ./my-project --filter "tags:production AND hostname:web*"
```

**Validate Machine Inventory**:
```bash
# Validate machine configuration
spooky machines validate ./my-project

# Validate with detailed output
spooky machines validate ./my-project --verbose
```

**Machine Connectivity**:
```bash
# Test connectivity to all machines
spooky machines ping ./my-project

# Test specific machine
spooky machines ping ./my-project --machine web-01

# Test with timeout
spooky machines ping ./my-project --timeout 30

# Test by tags
spooky machines ping ./my-project --tags "production"
```

**Machine Connection**:
```bash
# Connect to specific machine
spooky machines connect web-01 ./my-project

# Connect with custom command
spooky machines connect web-01 ./my-project --command "uptime"
```

### **2. Import/Export Commands**

**Import Machines**:
```bash
# Import from Kubernetes
spooky machines import --from kubernetes k8s-nodes.json ./my-project

# Import from Terraform state
spooky machines import --from terraform terraform.tfstate ./my-project

# Import from custom format
spooky machines import --from custom custom-inventory.json ./my-project
```

**Export Machines**:
```bash
# Export to JSON
spooky machines export ./my-project --format json --output machines.json

# Export to HCL
spooky machines export ./my-project --format hcl --output machines.hcl

# Export specific machines
spooky machines export ./my-project --machine web-01 --format json

# Export by tags
spooky machines export ./my-project --tags "production" --format json
```

## Implementation Details

### **1. Machine Configuration Types**

```go
// Machine represents a remote machine configuration
type Machine struct {
    Name     string            `hcl:"name,label" validate:"required"`
    Host     string            `hcl:"host" validate:"required"`
    Port     int               `hcl:"port,optional" validate:"omitempty,min=1,max=65535"`
    User     string            `hcl:"user" validate:"required"`
    Password string            `hcl:"password,optional"`
    KeyFile  string            `hcl:"key_file,optional"`
    Tags     map[string]string `hcl:"tags,optional" validate:"omitempty,dive,keys,required,endkeys,required"`
}

// InventoryConfig represents machine inventory configuration
type InventoryConfig struct {
    Machines []Machine `hcl:"machine,block" validate:"required,min=1,dive"`
}
```

### **2. Enterprise-Scale Indexing**

```go
// CompositeIndex provides multi-level indexing for enterprise-scale deployments
type CompositeIndex struct {
    TagIndex        TagIndex
    MachineTagIndex MachineTagIndex
    TagCount        map[string]int
    Metrics         *IndexMetrics
}

// IndexCache provides thread-safe caching of indexes
type IndexCache struct {
    index      *CompositeIndex
    lastBuilt  time.Time
    configHash string
    mutex      sync.RWMutex
    metrics    *IndexMetrics
}
```

### **3. Machine Validation**

```go
// ValidateMachine validates machine configuration
func ValidateMachine(machine *Machine) error {
    // Validate required fields
    if machine.Name == "" {
        return fmt.Errorf("machine name is required")
    }
    if machine.Host == "" {
        return fmt.Errorf("machine host is required")
    }
    if machine.User == "" {
        return fmt.Errorf("machine user is required")
    }
    
    // Validate authentication (password or key_file)
    if machine.Password == "" && machine.KeyFile == "" {
        return fmt.Errorf("either password or key_file must be specified")
    }
    
    // Validate port range
    if machine.Port < 1 || machine.Port > 65535 {
        return fmt.Errorf("port must be between 1 and 65535")
    }
    
    return nil
}
```

### **4. Machine Connectivity Testing**

```go
// TestMachineConnectivity tests SSH connectivity to a machine
func TestMachineConnectivity(machine *Machine, timeout time.Duration) error {
    config := &ssh.ClientConfig{
        User: machine.User,
        Auth: []ssh.AuthMethod{},
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        Timeout:         timeout,
    }
    
    // Add authentication method
    if machine.Password != "" {
        config.Auth = append(config.Auth, ssh.Password(machine.Password))
    } else if machine.KeyFile != "" {
        key, err := loadPrivateKey(machine.KeyFile)
        if err != nil {
            return fmt.Errorf("failed to load private key: %w", err)
        }
        config.Auth = append(config.Auth, ssh.PublicKeys(key))
    }
    
    // Test connection
    client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", machine.Host, machine.Port), config)
    if err != nil {
        return fmt.Errorf("failed to connect: %w", err)
    }
    defer client.Close()
    
    return nil
}
```

## Performance Considerations

### **1. Large Inventory Optimization**

**Indexing Performance**:
- **Composite Indexes**: Multi-level indexing for O(1) lookups
- **Tag-Based Caching**: Cached tag indexes for frequent queries
- **Memory Optimization**: Efficient memory usage for large inventories
- **Parallel Processing**: Concurrent machine operations

**Lookup Performance**:
```go
// O(1) machine lookup by name
func (ic *IndexCache) GetMachineByName(name string) *Machine {
    // Direct map lookup
    return ic.machineMap[name]
}

// O(1) machine lookup by tags
func (ic *IndexCache) GetMachinesByTags(tags []string) []*Machine {
    // Pre-computed tag index lookup
    return ic.tagIndex[tagKey]
}
```

### **2. Connection Pooling**

**SSH Connection Pooling**:
```go
// ConnectionPool manages SSH connections
type ConnectionPool struct {
    connections map[string]*ssh.Client
    mutex       sync.RWMutex
    maxConnections int
    timeout       time.Duration
}

// GetConnection gets or creates SSH connection
func (cp *ConnectionPool) GetConnection(machine *Machine) (*ssh.Client, error) {
    cp.mutex.Lock()
    defer cp.mutex.Unlock()
    
    key := fmt.Sprintf("%s:%d", machine.Host, machine.Port)
    
    if conn, exists := cp.connections[key]; exists {
        return conn, nil
    }
    
    // Create new connection
    conn, err := createSSHConnection(machine)
    if err != nil {
        return nil, err
    }
    
    cp.connections[key] = conn
    return conn, nil
}
```

### **3. Parallel Operations**

**Parallel Machine Operations**:
```go
// RunParallel runs operations on machines in parallel
func RunParallel(machines []*Machine, operation func(*Machine) error, maxParallel int) error {
    semaphore := make(chan struct{}, maxParallel)
    var wg sync.WaitGroup
    errors := make(chan error, len(machines))
    
    for _, machine := range machines {
        wg.Add(1)
        go func(m *Machine) {
            defer wg.Done()
            semaphore <- struct{}{}
            defer func() { <-semaphore }()
            
            if err := operation(m); err != nil {
                errors <- fmt.Errorf("machine %s: %w", m.Name, err)
            }
        }(machine)
    }
    
    wg.Wait()
    close(errors)
    
    // Collect errors
    var errs []error
    for err := range errors {
        errs = append(errs, err)
    }
    
    if len(errs) > 0 {
        return fmt.Errorf("parallel operation failed: %v", errs)
    }
    
    return nil
}
```

## Security Considerations

### **1. Authentication Security**

**SSH Key Security**:
- **Key Validation**: Validate SSH key file permissions and format
- **Key Rotation**: Support for key rotation and management
- **Key Encryption**: Support for encrypted SSH keys with passphrases
- **Key Verification**: Verify SSH key authenticity and integrity

**Password Security**:
- **Password Validation**: Validate password strength and complexity
- **Password Encryption**: Encrypt passwords in configuration files
- **Password Rotation**: Support for password rotation policies
- **Password Policies**: Enforce password policies and requirements

### **2. Network Security**

**Connection Security**:
- **SSH Protocol**: Use SSH protocol for secure connections
- **Host Key Verification**: Verify host keys to prevent MITM attacks
- **Connection Encryption**: Encrypt all data in transit
- **Connection Timeouts**: Implement connection timeouts to prevent hanging connections

**Access Control**:
- **User Permissions**: Validate user permissions and access rights
- **Network Access**: Control network access and firewall rules
- **Audit Logging**: Log all machine access and operations
- **Access Monitoring**: Monitor and alert on suspicious access patterns

### **3. Configuration Security**

**Configuration Validation**:
- **Schema Validation**: Validate machine configuration against schema
- **Security Validation**: Validate security settings and configurations
- **Access Validation**: Validate access permissions and authentication
- **Network Validation**: Validate network connectivity and security

**Configuration Encryption**:
- **Sensitive Data**: Encrypt sensitive data in configuration files
- **Key Management**: Secure key management and storage
- **Access Control**: Control access to configuration files
- **Audit Trail**: Maintain audit trail for configuration changes

## Future Enhancements

### **1. Machine Lifecycle Management**

**Machine Provisioning**:
- **Auto-Provisioning**: Automatic machine provisioning and configuration
- **Template-Based**: Template-based machine provisioning
- **Cloud Integration**: Cloud provider integration for provisioning
- **Configuration Management**: Automated configuration management

**Machine Decommissioning**:
- **Safe Decommissioning**: Safe machine decommissioning procedures
- **Data Cleanup**: Automated data cleanup and removal
- **Access Revocation**: Automatic access revocation and cleanup
- **Audit Trail**: Complete audit trail for decommissioning

### **2. Advanced Machine Features**

**Machine Monitoring**:
- **Health Monitoring**: Real-time machine health monitoring
- **Performance Metrics**: Machine performance metrics and analytics
- **Alerting**: Automated alerting for machine issues
- **Reporting**: Comprehensive machine reporting and analytics

**Machine Automation**:
- **Auto-Recovery**: Automatic machine recovery and failover
- **Load Balancing**: Automatic load balancing and distribution
- **Scaling**: Automatic scaling based on demand
- **Optimization**: Automatic performance optimization

### **3. Integration Enhancements**

**External System Integration**:
- **CMDB Integration**: Enhanced CMDB system integration
- **Monitoring Integration**: Integration with monitoring systems
- **Ticketing Integration**: Integration with ticketing systems
- **CI/CD Integration**: Integration with CI/CD pipelines

**API Enhancements**:
- **REST API**: RESTful API for machine management
- **GraphQL API**: GraphQL API for flexible queries
- **Webhook Support**: Webhook support for events and notifications
- **Plugin System**: Plugin system for custom integrations