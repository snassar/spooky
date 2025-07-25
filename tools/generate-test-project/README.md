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
./build/generate-test-project [flags]
```

### Flags

- `-s, --scale string`: Scale configuration (small, medium, large, testing, or custom) (default "medium")
- `-w, --hardware int`: Number of hardware machines (for custom scale)
- `-v, --vms int`: Number of VM machines (for custom scale)
- `-c, --containers int`: Number of container machines (for custom scale)
- `-o, --output string`: Output directory for generated projects (default "./test-projects")
- `-h, --help`: Show help information

### Scales

#### Predefined Scales

- **small**: 25 machines (2-5 hardware, rest random mix of VMs and containers)
- **medium**: 250 machines (20-40 hardware, rest random mix of VMs and containers)
- **large**: 1500 machines (150-250 hardware, rest random mix of VMs and containers)
- **testing**: 10000 machines (completely random distribution of hardware, VMs, and containers)

#### Custom Scale

Use the `--scale custom` flag along with `--hardware`, `--vms`, and `--containers` flags to specify exact counts:

```bash
./build/generate-test-project --scale custom --hardware 50 --vms 150 --containers 100 --output ./test-projects
```

### Examples

```bash
# Generate a small test project
./build/generate-test-project --scale small --output ./test-projects

# Generate a medium test project (default)
./build/generate-test-project

# Generate a large test project
./build/generate-test-project --scale large --output ./test-projects

# Generate a custom test project with 20 hardware, 80 VMs, and 50 containers
./build/generate-test-project --scale custom --hardware 20 --vms 80 --containers 50 --output ./test-projects

# Generate in current directory
./build/generate-test-project --scale small --output .

# Show help
./build/generate-test-project --help
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

### Container Machines
- **Types**: app-server, api-gateway, worker, redis
- **OS**: Alpine
- **IP Ranges**: 
  - App servers: 10.1.50.x, 10.2.50.x
  - API gateways: 10.1.60.x, 10.2.60.x
  - Workers: 10.1.70.x, 10.2.70.x
  - Redis: 10.1.80.x, 10.2.80.x

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

## Testing the Tool

```bash
# Test help functionality
./build/generate-test-project --help

# Test different scales
./build/generate-test-project --scale small --output ./test-small
./build/generate-test-project --scale medium --output ./test-medium
./build/generate-test-project --scale large --output ./test-large

# Test custom scale
./build/generate-test-project --scale custom --hardware 10 --vms 30 --containers 20 --output ./test-custom

# Test error handling
./build/generate-test-project --scale invalid
./build/generate-test-project --scale custom  # Should fail without hardware/vms/containers
```

## Building

```bash
# Build from the project root
go build -o build/generate-test-project tools/generate-test-project/main.go

# Or build from the tool directory
cd tools/generate-test-project
go build -o generate-test-project main.go
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