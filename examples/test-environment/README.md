# Podman-Based Test Environment

This directory contains the Podman-based test environment for spooky integration testing, replacing the custom SSH server implementations with real Debian 12 containers.

## Overview

The test environment consists of three Debian 12 containers running SSH servers with different configurations:

- **spooky-server1**: Root user with password authentication (port 2221)
- **spooky-server2**: Admin user with key-based authentication (port 2222)
- **spooky-server3**: User with SFTP support (port 2223)

## Files

### Quadlet Configuration Files

- `spooky-test.network` - Podman network configuration
- `spooky-server1.container` - Container configuration for root/password auth
- `spooky-server2.container` - Container configuration for admin/key auth
- `spooky-server3.container` - Container configuration for user/SFTP

### Test Configuration

- `test-config.hcl` - Comprehensive test configuration with all scenarios

## Usage

### Prerequisites

- Podman 4.4+ with Quadlet support
- systemd user session enabled
- Rootless Podman capability

### Quick Start

```bash
# Build the test environment tool
make build-spooky-test-env

# Run preflight check
./build/spooky-test-env preflight

# Start the test environment
./build/spooky-test-env start

# Run integration tests
go test -v -podman ./tests/podman_integration_test.go

# Stop the test environment
./build/spooky-test-env stop

# Clean up everything
./build/spooky-test-env cleanup
```

### Using the Test Configuration

```bash
# Execute all tests
go run main.go execute examples/test-environment/test-config.hcl

# Execute specific actions
go run main.go execute examples/test-environment/test-config.hcl --action test-authentication
go run main.go execute examples/test-environment/test-config.hcl --action test-sftp-operations
```

## Test Scenarios

The test configuration includes comprehensive scenarios:

### Authentication Tests
- **test-authentication**: Basic user authentication across all servers
- **test-system-info**: System information gathering

### File Operations
- **test-file-operations**: Basic file operations on standard servers
- **test-sftp-operations**: SFTP-specific operations

### Tag-Based Targeting
- **test-production-servers**: Target servers by production tier
- **test-staging-servers**: Target servers by staging tier
- **test-database-servers**: Target servers by database role
- **test-web-servers**: Target servers by web role

### Advanced Scenarios
- **test-concurrent-operations**: Parallel execution testing
- **test-error-handling**: Error handling scenarios
- **test-network-connectivity**: Network connectivity between servers

## Server Details

### spooky-server1 (Root/Password)
- **Host**: localhost:2221
- **User**: root
- **Password**: password
- **Authentication**: Password-based
- **Tags**: role=database, tier=production

### spooky-server2 (Admin/Key)
- **Host**: localhost:2222
- **User**: admin
- **Password**: adminpass (fallback)
- **Authentication**: Key-based (with password fallback)
- **Tags**: role=web, tier=staging

### spooky-server3 (User/SFTP)
- **Host**: localhost:2223
- **User**: user
- **Password**: userpass
- **Authentication**: Password-based
- **SFTP**: Enabled with chroot
- **Tags**: role=storage, tier=development

## Network Configuration

- **Network**: spooky-test
- **Subnet**: 10.1.0.0/16
- **Gateway**: 10.1.0.1
- **Server IPs**: 10.1.10.1, 10.1.10.2, 10.1.10.3

## Migration from Legacy Tests

This environment replaces the custom SSH server implementations in `tests/infrastructure/` with:

### Benefits
- **Real SSH behavior**: Actual OpenSSH servers instead of mocks
- **Multiple authentication methods**: Password, key-based, SFTP
- **Professional tooling**: Quadlet-based container management
- **Cross-platform compatibility**: Works on Linux, macOS, Windows
- **CI/CD integration**: Already configured in GitHub Actions

### Migration Path
1. **Phase 1**: Use both environments in parallel
2. **Phase 2**: Migrate critical test scenarios
3. **Phase 3**: Remove legacy infrastructure

## Troubleshooting

### Common Issues

**Port conflicts**: Ensure ports 2221-2223 are not in use
```bash
sudo netstat -tlnp | grep :222
```

**Container startup issues**: Check container logs
```bash
podman logs spooky-server1
podman logs spooky-server2
podman logs spooky-server3
```

**Network issues**: Recreate the network
```bash
podman network rm spooky-test
podman network create spooky-test
```

**Quadlet issues**: Check Quadlet support
```bash
podman quadlet --help
```

### Debug Commands

```bash
# Check environment status
./build/spooky-test-env status

# Check container IPs
podman inspect spooky-server1 --format "{{.NetworkSettings.Networks.spooky-test.IPAddress}}"

# Test SSH connectivity
ssh -p 2221 root@localhost
ssh -p 2222 admin@localhost
ssh -p 2223 user@localhost
```

## Security Note

This test environment is configured for testing purposes only:
- Root login is enabled
- Simple password authentication
- No firewall or security hardening

**Do not use this configuration in production environments.** 