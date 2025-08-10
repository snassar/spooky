# HCL Fact Collector Example

This example demonstrates the comprehensive capabilities of the Spooky HCL fact collector, showcasing how to collect, process, and export facts from HCL configuration files.

## Overview

The HCL (HashiCorp Configuration Language) fact collector is designed to parse HCL files and convert them into structured fact collections. It supports:

- **Single file processing** - Parse individual HCL files
- **Directory scanning** - Process multiple HCL files in a directory
- **Selective collection** - Collect only specific facts
- **Individual fact retrieval** - Get single facts by key
- **File validation** - Validate HCL syntax
- **Export functionality** - Convert facts back to HCL format
- **Complex HCL features** - Support for blocks, expressions, and nested structures

## Features Demonstrated

### 1. Single HCL File Collection
```go
collector := hcl.NewCollector([]string{"config.hcl"})
collection, err := collector.Collect("server-name")
```

### 2. Directory Collection
```go
collector := hcl.NewCollector([]string{"/path/to/configs"})
collection, err := collector.Collect("all-servers")
```

### 3. Selective Fact Collection
```go
specificKeys := []string{"name", "version", "config.web.host"}
collection, err := collector.CollectSpecific("server", specificKeys)
```

### 4. Individual Fact Retrieval
```go
fact, err := collector.GetFact("server", "name")
```

### 5. HCL File Validation
```go
err := collector.ValidateHCLFile("config.hcl")
```

### 6. Finding HCL Files
```go
hclFiles, err := collector.FindHCLFiles("/path/to/directory")
```

### 7. Exporting Facts to HCL
```go
facts := map[string]interface{}{
    "name": "server",
    "port": 8080,
}
err := collector.ExportFactsToHCL(facts, outputFile)
```

## Sample HCL Files

The example creates three sample HCL files:

### 1. Server Configuration (`server-config.hcl`)
```hcl
# Web server configuration
name = "web-server-01"
version = "1.2.3"
enabled = true
port = 8080

# Web configuration block
config "web" {
  host = "localhost"
  ssl = true
  timeout = 30
}

# API configuration block
config "api" {
  endpoint = "/api/v1"
  rate_limit = 1000
  auth_required = true
}

# Tags
tags = ["web", "api", "production", "load-balanced"]

# Environment variables
env = {
  NODE_ENV = "production"
  LOG_LEVEL = "info"
  DEBUG = false
}
```

### 2. Database Configuration (`db-config.hcl`)
```hcl
# Database configuration
name = "db-server-01"
version = "5.7.0"
enabled = true
port = 3306

# Database configuration block
config "database" {
  host = "db.internal"
  port = 3306
  name = "app_database"
  user = "app_user"
  max_connections = 100
  ssl_mode = "require"
}

# Backup configuration
backup = {
  enabled = true
  schedule = "0 2 * * *"
  retention_days = 30
  storage_path = "/backups"
}
```

### 3. Network Configuration (`network-config.hcl`)
```hcl
# Network configuration
name = "network-config"
version = "1.0.0"

# Load balancer configuration
load_balancer "primary" {
  algorithm = "round_robin"
  health_check_path = "/health"
  health_check_interval = 10
  timeout = 5
  
  backend "web1" {
    host = "10.0.1.10"
    port = 8080
    weight = 1
  }
  
  backend "web2" {
    host = "10.0.1.11"
    port = 8080
    weight = 1
  }
}
```

## Running the Example

```bash
cd examples/hcl-collector-example
go run main.go
```

## Expected Output

The example will demonstrate:

1. **Single file collection** - Shows facts extracted from `server-config.hcl`
2. **Directory collection** - Shows facts from all HCL files in the directory
3. **Selective collection** - Shows only requested facts
4. **Individual fact retrieval** - Shows retrieving a single fact with metadata
5. **File validation** - Shows validation of valid and invalid HCL files
6. **File discovery** - Shows finding all HCL files in a directory
7. **Export functionality** - Shows converting facts back to HCL format
8. **JSON comparison** - Shows the JSON representation of collected facts

## Key Features

### HCL Block Support
The collector properly handles HCL blocks, converting them to dot-notation facts:
- `config "web" { host = "localhost" }` becomes `config.web.host = "localhost"`

### Complex Data Types
Supports all HCL data types:
- **Strings**: `name = "server"`
- **Numbers**: `port = 8080`
- **Booleans**: `enabled = true`
- **Lists**: `tags = ["web", "api"]`
- **Maps**: `env = { NODE_ENV = "production" }`

### Metadata Preservation
Each fact includes rich metadata:
- `source_machine`: Always "local" for HCL facts
- `source_file`: The HCL file the fact came from
- `source_type`: Always "hcl"
- `hcl_type`: Type of HCL element (attribute, block_attribute, etc.)
- `block_type`: For block attributes, the block type
- `block_labels`: For block attributes, the block labels

### Error Handling
Comprehensive error handling for:
- Invalid HCL syntax
- Missing files
- File size limits
- Parsing errors

## Integration with Spooky

The HCL collector integrates seamlessly with the Spooky fact collection system:

```go
// Create collector with merge policy
collector := hcl.NewCollector([]string{"config.hcl"})
collector.SetMergePolicy(types.MergePolicyReplace)

// Collect facts
collection, err := collector.Collect("server-name")

// Use facts in Spooky system
for key, fact := range collection.Facts {
    fmt.Printf("Fact: %s = %v\n", key, fact.Value)
}
```

## Advanced Usage

### Custom Parser
You can provide a custom HCL parser:

```go
type CustomParser struct{}

func (p *CustomParser) ParseFile(filePath string) (map[string]interface{}, error) {
    // Custom parsing logic
}

func (p *CustomParser) ParseContent(content []byte) (map[string]interface{}, error) {
    // Custom content parsing
}

collector := hcl.NewCollectorWithParser([]string{"config.hcl"}, &CustomParser{})
```

### File Management
```go
collector := hcl.NewCollector([]string{})

// Add files dynamically
collector.AddHCLFile("config1.hcl")
collector.AddHCLFile("config2.hcl")

// Remove files
collector.RemoveHCLFile("config1.hcl")

// Get source information
sources := collector.GetFactSources()
```

## Performance Considerations

- **File size limits**: Default 10MB limit per file
- **Directory scanning**: Recursive scanning with file filtering
- **Memory usage**: Efficient parsing with streaming where possible
- **Error recovery**: Continues processing other files if one fails

## Dependencies

- `github.com/hashicorp/hcl/v2` - Official HCL parsing library
- Standard Go libraries for file operations

## Comparison with JSON Collector

| Feature | HCL Collector | JSON Collector |
|---------|---------------|----------------|
| **Complexity** | High (supports blocks, expressions) | Low (simple JSON) |
| **Performance** | Slower (complex parsing) | Faster (native Go) |
| **Features** | Rich (blocks, expressions, validation) | Basic (simple data types) |
| **Dependencies** | External (HCL v2) | None (standard library) |
| **Use Case** | Configuration files, complex structures | Simple data exchange |

## Conclusion

The HCL fact collector provides a powerful way to extract structured facts from HCL configuration files, making it ideal for infrastructure automation, configuration management, and complex system configurations. Its support for HCL blocks, expressions, and validation makes it a robust solution for processing HashiCorp ecosystem configuration files.
