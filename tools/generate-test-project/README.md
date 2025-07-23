# Generate Test Project

A helper script that uses `spooky project init` to generate test projects with fake inventory and example actions for testing spooky against projects of arbitrary sizes.

## Overview

This tool replaces the old `@/generate-config` approach by leveraging the new separated project structure. It creates complete test projects with:

- **Separated inventory and actions**: Uses the new `inventory.hcl` and `actions.hcl` file structure
- **Realistic test data**: Generates machines with proper IP addressing, tags, and roles
- **Comprehensive actions**: Creates 15 different test actions covering various scenarios
- **Flexible scaling**: Supports predefined scales and custom configurations

## Usage

```bash
./build/generate-test-project <scale> [output-dir]
```

### Scales

#### Predefined Scales

- **small**: 40 machines (10 hardware + 30 VMs)
- **medium**: 400 machines (100 hardware + 300 VMs)  
- **large**: 10,000 machines (2,500 hardware + 7,500 VMs)

#### Custom Scale

Use the format `custom:hardware:vms` to specify exact counts:

```bash
./build/generate-test-project custom:50:150 ./test-projects
```

### Examples

```bash
# Generate a small test project
./build/generate-test-project small ./test-projects

# Generate a medium test project
./build/generate-test-project medium ./test-projects

# Generate a custom test project with 20 hardware and 80 VMs
./build/generate-test-project custom:20:80 ./test-projects

# Generate in current directory
./build/generate-test-project small .
```

## Generated Project Structure

Each generated project follows the new separated structure:

```
project-name/
├── project.hcl          # Project configuration
├── inventory.hcl        # Machine definitions
├── actions.hcl          # Action definitions
├── templates/           # Template directory
├── files/              # Files directory
├── logs/               # Logs directory
└── .gitignore          # Git ignore file
```

## Machine Types

The tool generates realistic machine configurations:

### Hardware Machines
- **Datacenters**: FRA00 (Frankfurt) and BER0 (Berlin)
- **IP Range**: 10.1.1.x and 10.2.1.x
- **OS**: Debian 12
- **Role**: VM host
- **Capacity**: High

### VM Machines
- **Types**: web-server, database, cache, monitoring
- **OS**: Ubuntu 22
- **IP Ranges**: 
  - Web servers: 10.1.10.x, 10.2.10.x
  - Databases: 10.1.20.x, 10.2.20.x
  - Cache: 10.1.30.x, 10.2.30.x
  - Monitoring: 10.1.40.x, 10.2.40.x

## Generated Actions

The tool creates 15 test actions covering various scenarios:

1. **system-update**: Update system packages on all machines
2. **security-scan**: Run security vulnerability scan
3. **backup-databases**: Create database backups
4. **restart-services**: Restart critical services
5. **check-disk-space**: Check available disk space
6. **update-firewall**: Update firewall rules
7. **monitor-logs**: Check system logs for errors
8. **cleanup-temp**: Clean up temporary files
9. **check-memory**: Check memory usage
10. **update-ssl-certs**: Update SSL certificates
11. **database-maintenance**: Perform database maintenance tasks
12. **network-test**: Test network connectivity
13. **update-monitoring**: Update monitoring agents
14. **backup-configs**: Backup configuration files
15. **health-check**: Perform comprehensive health check

## Testing Generated Projects

After generating a project, you can test it with spooky:

```bash
# List project contents
./build/spooky project list ./test-projects/project-name

# List with detailed output
./build/spooky project list ./test-projects/project-name --verbose

# Validate project configuration
./build/spooky project validate ./test-projects/project-name

# Validate with debug output
./build/spooky project validate ./test-projects/project-name --debug
```

## Building

```bash
go build -o build/generate-test-project tools/generate-test-project/main.go
```

## Testing

```bash
go test ./tools/generate-test-project
```

## Features

- **Deterministic IDs**: Uses Git-style SHA256 hashes for machine IDs
- **Realistic networking**: Proper IP addressing scheme across datacenters
- **Tag-based targeting**: Actions use tags for machine selection
- **Mixed execution modes**: Both parallel and sequential actions
- **Timeout configuration**: Realistic timeout values for different action types
- **Script support**: Some actions use scripts instead of commands
- **Comprehensive coverage**: Actions target different machine types and roles 