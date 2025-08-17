# Facts System: Current Implementation

## Overview

This document describes the current implementation of the spooky facts system. It covers the actual implemented fact collection, storage, and export capabilities.

**Schema Integration**: The facts system uses basic validation patterns and data structures defined in the types system for fact validation and storage format consistency.

**Architecture Integration**: Facts integrate with the overall spooky architecture, providing machine data collection and export capabilities through the CLI system.

## System Integration

The facts system integrates with other core Spooky systems to provide fact management and data sharing:

### **Project System Integration**
- **Project Context**: Facts collection uses project machine inventory for target identification
- **Project Path**: Facts export uses project paths for machine inventory loading
- **Project Validation**: Facts operations validate project structure before running

### **CLI System Integration**
- **Facts Commands**: Facts management through `spooky facts export` command
- **Machine Integration**: Facts collection uses machine inventory from projects
- **Export Functionality**: Facts export to JSON and HCL formats

### **Machines System Integration**
- **Machine Inventory**: Facts collection uses machine inventory for target identification
- **Machine Metadata**: Machine inventory provides hostname and connection details
- **Machine Filtering**: Facts export supports filtering by machine, tags, and groups

## Current Implementation

### ✅ Actually Implemented

1. **Basic Fact Collection** - System fact collection using gopsutil
2. **Memory Storage** - In-memory fact storage during export operations
3. **Export Functionality** - Facts export to JSON and HCL formats
4. **CLI Integration** - `spooky facts export` command with filtering options
5. **Machine Integration** - Facts collection from project machine inventory
6. **Basic Validation** - Fact collection validation and error handling
7. **Parallel Processing** - Support for parallel fact collection
8. **Filtering** - Support for machine, tag, and group filtering

### Current Fact System Structure

```
internal/facts/
├── manager.go         # Fact collection coordination and management
├── types.go           # Fact data structures and types
├── memory_storage.go  # In-memory storage implementation
├── integration.go     # CLI integration layer
├── collector.go       # System fact collection using gopsutil
├── manager_test.go    # Manager unit tests
└── integration_test.go # Integration tests
```

## Implementation Details

### **1. Fact Collection**

The system uses a single `SystemFactCollector` that collects facts using the gopsutil library:

```go
type SystemFactCollector struct {
    name string
}
```

**Collected Facts:**
- **OS Facts**: Operating system information (name, version, architecture)
- **Hardware Facts**: CPU, memory, disk information
- **Network Facts**: Network interfaces and configuration
- **Load Average**: System load information
- **Process Facts**: Basic process information

**Collection Process:**
1. Get machine ID from `/etc/machine-id`
2. Collect system facts using gopsutil
3. Organize facts into structured format
4. Return fact collection with metadata

### **2. Memory Storage**

Facts are gathered directly and exported without intermediate storage:

```go
type MemoryFactStorage struct {
    mutex sync.RWMutex
}
```

**Features:**
- Thread-safe operations with RWMutex
- Direct export to JSON and HCL formats
- No intermediate storage (facts gathered and exported immediately)
- Storage statistics for debugging
- Minimal memory footprint

**Storage Operations:**
- `GetStats()` - Get memory usage statistics for debugging
- `ExportToJSON()` - Export facts directly to JSON format following facts-structure.schema.hcl
- `ExportToHCL()` - Export facts directly to HCL format following facts-structure.schema.hcl

### **3. CLI Integration**

The facts system integrates with the CLI through the `spooky facts export` command:

**Command Structure:**
```bash
spooky facts export [project-path] [flags]
```

**Available Flags:**
- `--format`: Export format (hcl, json) (default: hcl)
- `--output`: Output file path (required)
- `--machine`: Target specific machine (default: all machines)
- `--tags`: Filter by tags (supports key=value or key-only format)
- `--groups`: Filter by groups (comma-separated list)
- `--parallel`: Number of parallel workers (default: 1)
- `--verbose`: Verbose output

**Example Usage:**
```bash
# Export all facts to HCL
spooky facts export ./my-project --output facts.hcl

# Export facts to JSON with parallel processing
spooky facts export ./my-project --format json --parallel 4 --output facts.json

# Export facts for specific machine
spooky facts export ./my-project --machine web-server --output web-server-facts.hcl

# Export facts with filtering
spooky facts export ./my-project --tags "environment=production" --output prod-facts.hcl
```

### **4. Machine Integration**

Facts collection integrates with the machine inventory system:

**Machine Discovery:**
- Reads machine inventory from project's machines.hcl files
- Supports both single machines.hcl and machines/ directory structure
- Uses machine metadata for connection details

**Machine Filtering:**
- **Machine Filter**: Filter by hostname (partial match)
- **Tag Filtering**: Filter by tag key=value or key-only
- **Group Filtering**: Filter by machine groups
- **Combined Filtering**: Combine multiple filter types

**Example Machine Inventory:**
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
    
    tags = {
      "environment" = "production"
      "role" = "web-server"
    }
    
    groups = ["webservers", "production"]
  }
}
```

### **5. Export Functionality**

The facts system supports exporting facts to multiple formats:

**Supported Formats:**
- **HCL**: HashiCorp Configuration Language format (default)
- **JSON**: Standard JSON format for data exchange

**Export Process:**
1. Load machine inventory from project
2. Apply filters (machine, tags, groups)
3. Collect facts from each machine in parallel
4. Store facts in memory
5. Export facts to specified format and file

**Example JSON Output:**
```json
{
  "machines": {
    "1234567890abcdef1234567890abcdef": {
      "machine_id": "1234567890abcdef1234567890abcdef",
      "collected_at": "2024-01-15T10:30:45Z",
      "facts": {
        "os": {
          "name": "Linux",
          "version": "Ubuntu 22.04.3 LTS",
          "architecture": "x86_64",
          "kernel": "5.15.0-88-generic"
        },
        "hardware": {
          "cpu_count": 4,
          "memory_total": 8589934592,
          "disk_total": 107374182400
        },
        "network": {
          "interfaces": [
            {
              "name": "eth0",
              "address": "192.168.1.10",
              "netmask": "255.255.255.0"
            }
          ]
        }
      }
    }
  }
}
```

## Error Handling

### **Collection Errors**

The facts system handles various error conditions:

**SSH Connection Failures:**
- Connection timeouts
- Authentication failures
- Network connectivity issues
- Host key validation failures

**Fact Collection Errors:**
- Missing system files (e.g., `/etc/machine-id`)
- Insufficient permissions
- System resource limitations
- gopsutil library errors

**Error Handling Strategy:**
- Individual machine failures don't stop collection from other machines
- Errors are logged with detailed context
- Failed machines are reported in export results
- Export continues with successfully collected facts

### **Validation Errors**

Basic validation is performed during fact collection:

**Data Validation:**
- Machine ID format validation
- Required fact fields presence
- Data type validation
- String length limits

**Schema Compliance:**
- Fact structure validation
- Required nested objects
- Array format validation
- Metadata validation

## Performance Considerations

### **Parallel Processing**

The facts system supports parallel fact collection:

**Parallel Workers:**
- Configurable via `--parallel` flag
- Must be 2 or larger (0 and 1 are invalid)
- Default is 1 (sequential processing)
- Recommended: 4-8 workers for large inventories

**Performance Benefits:**
- Faster collection from multiple machines
- Better resource utilization
- Reduced total export time
- Improved user experience

### **Memory Management**

Facts are gathered directly and exported without intermediate storage:

**Memory Usage:**
- Facts gathered from machines during export process
- No intermediate storage (direct gather → export)
- Minimal memory footprint
- No cleanup needed (no persistent storage)

**Memory Benefits:**
- Fast gather and export process
- No disk I/O overhead
- No memory leaks (no persistent storage)
- Simple and efficient for ephemeral data

## Integration Patterns

### **Interface-Based Design**

The facts system uses interface-based design for flexibility:

**Core Interfaces:**
- `FactsIntegration`: Main integration interface
- `FactManager`: Fact collection and management
- `FactStorage`: Storage operations
- `FactCollector`: Fact collection from machines

**Dependency Injection:**
- Interfaces injected into CLI commands
- Testable components
- Loose coupling between components
- Extensible architecture

### **CLI Integration**

Facts commands integrate with the CLI system:

**Command Structure:**
- Follows `spooky facts export` pattern
- Consistent with other CLI commands
- Standard flag handling
- Error reporting integration

**Global Integration:**
- Uses global CLI configuration
- Integrates with logging system
- Follows error handling patterns
- Supports verbose output

## Current Limitations

### **Storage Characteristics**

1. **Ephemeral Storage**: Facts are only stored in memory during export (this is intentional)
2. **No Fact History**: No historical fact tracking or comparison (facts are gathered fresh each time)
3. **No Persistent Storage**: Facts are not saved to disk (memory-only during export)
4. **No Import Functionality**: Facts cannot be imported (only export is supported)

### **Collection Limitations**

1. **Single Collector**: Only system fact collector available
2. **No Custom Facts**: No support for custom fact collection
3. **No Fact Caching**: No caching of previously collected facts
4. **No Incremental Collection**: Always collects all facts

### **Integration Limitations**

1. **No Template Integration**: Facts not integrated with template system
2. **No Variable Integration**: Facts not used in variable resolution
3. **No Action Integration**: Facts not integrated with action system
4. **No Real-time Collection**: No real-time fact monitoring

## Future Enhancements

### **Planned Features**

1. **Persistent Storage**: Long-term fact storage in databases
2. **Fact History**: Historical fact tracking and comparison
3. **Advanced Collectors**: Custom fact collectors for specific data
4. **Fact Validation**: Enhanced validation and schema checking
5. **Fact Import**: Import facts from external sources
6. **Fact Comparison**: Compare facts across machines and time periods

### **Integration Enhancements**

1. **Template Integration**: Use facts in template rendering
2. **Variable Integration**: Use facts in variable resolution
3. **Action Integration**: Use facts in action conditions
4. **Real-time Collection**: Real-time fact monitoring

### **Performance Enhancements**

1. **Fact Caching**: Cache frequently accessed facts
2. **Incremental Collection**: Only collect changed facts
3. **Compression**: Compress fact data for storage
4. **Distributed Collection**: Distributed fact collection across multiple nodes

## Conclusion

The current facts system implementation provides essential fact collection and export capabilities for gathering system information from machines. While the implementation is focused on core functionality, it provides a solid foundation for future enhancements.

The system integrates well with other spooky components and follows established architectural patterns. The interface-based design makes it extensible and maintainable, while the CLI integration provides a user-friendly interface for fact management.

Future enhancements will build upon this foundation to provide more comprehensive fact management capabilities, including persistent storage, advanced collectors, and deeper integration with other spooky systems.