# Facts System User Guide

## Overview

The spooky facts system provides fact collection and export capabilities for gathering system information from machines. This guide covers the current implementation of fact collection and export functionality.

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
4. **Custom Facts**: Read custom facts from `/etc/spooky/facts.{hcl,json}` if present
5. **Data Processing**: Convert and validate collected data
6. **Direct Export**: Write facts directly to output file in requested format
7. **Cleanup**: Memory is automatically managed during export

### Fact Storage

Facts are gathered directly and exported without intermediate storage:

- **Direct Export**: Facts are collected from machines and exported immediately
- **No Intermediate Storage**: No temporary storage step - direct gather → export
- **Minimal Memory Usage**: Only the facts being exported are kept in memory
- **Automatic Management**: Memory is automatically managed during export
- **Debugging Support**: Storage statistics available for performance monitoring

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
- `--tags`: Filter by tags (supports key=value or key-only format)
- `--groups`: Filter by groups (comma-separated list)
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

# Export facts by tags
spooky facts export ./my-project --tags "environment=production" --output prod-facts.hcl
spooky facts export ./my-project --tags "role=web,role=database" --output app-facts.hcl

# Export facts by groups
spooky facts export ./my-project --groups "webservers,database" --output app-facts.hcl

# Export facts with multiple filters
spooky facts export ./my-project --tags "role=web" --groups "production" --output web-prod-facts.hcl
```

#### Export Formats

Both JSON and HCL export formats follow the exact structure defined in `facts-structure.schema.hcl`:

**JSON Format:**
```json
[
  {
    "machine_id": "32-character-hex-string",
    "collected_at": "2024-01-01T12:00:00Z",
    "facts": {
      "system": {
        "os": { ... },
        "hardware": { ... },
        "network": { ... }
      }
    }
  }
]
```

**HCL Format:**
```hcl
facts_structure {
  machine_id = "32-character-hex-string"
  collected_at = "2024-01-01T12:00:00Z"
  facts {
    system {
      os { ... }
      hardware { ... }
      network { ... }
    }
  }
}
```

#### Filtering Options

The facts export command supports multiple filtering options to target specific subsets of machines:

**Machine Filter:**
```bash
# Filter by hostname (partial match)
spooky facts export ./my-project --machine "web" --output web-facts.hcl
```

**Tag Filtering:**
```bash
# Filter by tag key=value
spooky facts export ./my-project --tags "environment=production" --output prod-facts.hcl

# Filter by tag key only (any value)
spooky facts export ./my-project --tags "monitored" --output monitored-facts.hcl

# Multiple tag filters (OR logic)
spooky facts export ./my-project --tags "role=web,role=database" --output app-facts.hcl
```

**Group Filtering:**
```bash
# Filter by single group
spooky facts export ./my-project --groups "webservers" --output web-facts.hcl

# Filter by multiple groups (OR logic)
spooky facts export ./my-project --groups "webservers,database" --output app-facts.hcl
```

**Combined Filtering:**
```bash
# Combine machine, tag, and group filters (AND logic)
spooky facts export ./my-project \
  --machine "web" \
  --tags "environment=production" \
  --groups "webservers" \
  --output web-prod-facts.hcl
```

#### Example JSON Output

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

## Advanced Usage

### Parallel Processing

The facts export command supports parallel processing to improve performance when collecting facts from multiple machines:

```bash
# Use 4 parallel workers for fact collection
spooky facts export ./my-project --parallel 4 --output facts.hcl

# Use 8 parallel workers for large inventories
spooky facts export ./my-project --parallel 8 --output facts.hcl
```

**Note**: The `--parallel` flag must be 2 or larger. Values 0 and 1 are invalid.

### Verbose Output

Enable verbose output to see detailed information about the fact collection process:

```bash
# Enable verbose output for detailed logging
spooky facts export ./my-project --verbose --output facts.hcl
```

Verbose output includes:
- Machine discovery information
- SSH connection details
- Fact collection progress
- Validation results
- Export statistics

### Error Handling

The facts export command handles errors gracefully:

- **Individual Machine Failures**: If fact collection fails for one machine, the command continues with other machines
- **Network Timeouts**: SSH connection timeouts are handled with retry logic
- **Authentication Failures**: Clear error messages for SSH authentication issues
- **Validation Errors**: Invalid facts are logged but don't stop the export process

## Troubleshooting

### Common Issues

#### SSH Connection Failures

```bash
# Check SSH connectivity manually
ssh user@hostname

# Verify SSH key permissions
ls -la ~/.ssh/id_rsa
# Should show: -rw------- (600 permissions)

# Test SSH connection with verbose output
ssh -v user@hostname
```

#### Machine Inventory Issues

```bash
# Validate machine inventory
spooky machines validate ./my-project

# Check machine inventory file
cat ./my-project/machines.hcl
```

#### Fact Collection Failures

```bash
# Enable verbose output for detailed error information
spooky facts export ./my-project --verbose --output facts.hcl

# Check system requirements on target machines
# - gopsutil library support
# - /etc/machine-id file exists
# - Sufficient permissions for fact collection
```

### Debug Information

When troubleshooting fact collection issues, enable verbose output to get detailed information:

```bash
# Enable verbose output
spooky facts export ./my-project --verbose --output facts.hcl
```

This will show:
- Machine discovery process
- SSH connection attempts
- Fact collection progress
- Validation results
- Error details for failed operations

## Integration with Other Systems

### Project Integration

Facts collection integrates with the project system:

- **Machine Inventory**: Uses machines.hcl for target identification
- **Project Structure**: Follows project directory structure
- **Configuration**: Uses project-specific configuration

### CLI Integration

Facts commands integrate with the CLI system:

- **Command Structure**: Follows `spooky facts export` pattern
- **Global Flags**: Supports global CLI flags
- **Error Handling**: Uses consistent error handling patterns

## Current Implementation Status

### ✅ Implemented Features

- **Basic Fact Collection**: System fact collection using gopsutil
- **Memory Storage**: Minimal memory usage during fact gathering and export
- **Export Functionality**: Facts export to JSON and HCL formats following facts-structure.schema.hcl
- **CLI Integration**: `spooky facts export` command with filtering options
- **Machine Integration**: Facts collection from project machine inventory
- **Basic Validation**: Fact collection validation and error handling
- **Parallel Processing**: Support for parallel fact collection
- **Filtering**: Support for machine, tag, and group filtering

### 📋 Future Enhancements

- **Persistent Storage**: Long-term fact storage in databases
- **Fact History**: Historical fact tracking and comparison
- **Advanced Collectors**: Custom fact collectors for specific data
- **Fact Validation**: Enhanced validation and schema checking
- **Fact Comparison**: Compare facts across machines and time periods

## Best Practices

### Performance Optimization

1. **Use Parallel Processing**: Use the `--parallel` flag for large inventories
2. **Filter Appropriately**: Use machine, tag, and group filters to target specific machines
3. **Monitor Resource Usage**: Fact collection can be resource-intensive on target machines

### Security Considerations

1. **SSH Key Management**: Ensure SSH keys have proper permissions (600)
2. **Network Security**: Use secure networks for fact collection
3. **Data Sensitivity**: Be aware of sensitive information in collected facts

### Maintenance

1. **Regular Exports**: Export facts regularly for backup and analysis
2. **Validation**: Validate fact exports for data integrity
3. **Storage Management**: Manage exported fact files appropriately

## Conclusion

The spooky facts system provides essential fact collection and export capabilities for gathering system information from machines. The current implementation focuses on the core functionality needed for basic fact management, with a clear path for future enhancements.

The system integrates well with other spooky components and provides a solid foundation for more advanced fact management features in the future.
