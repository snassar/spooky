# Integration Testing with Podman

This directory contains integration tests for spooky that use Podman to create isolated test environments with SSH-enabled containers.

## Overview

The integration tests verify that spooky can:
- Connect to SSH servers running in containers
- Execute commands and scripts remotely
- Handle parallel and sequential execution
- Process configuration files correctly
- Handle errors gracefully

## Test Structure

### Files

- `podman_ci_test.go` - CI-specific tests that run in GitHub Actions
- `podman_integration_test.go` - Local integration tests for development
- `podman_basic_test.go` - Basic Podman environment tests
- `test-config.hcl` - Sample configuration for testing
- `test-script.sh` - Sample script for testing script execution

### Test Categories

1. **Basic Environment Tests** (`podman_basic_test.go`)
   - Verify Podman is available and working
   - Test basic container operations

2. **Local Integration Tests** (`podman_integration_test.go`)
   - Full integration tests for local development
   - Requires `-podman` flag to run

3. **CI Integration Tests** (`podman_ci_test.go`)
   - Tests designed for GitHub Actions CI environment
   - Automatically skipped when not in CI

## Running Tests

### Local Development

```bash
# Run basic Podman tests
go test ./tests/integration -v -podman-basic

# Run full integration tests (requires Podman)
go test ./tests/integration -v -podman

# Run all integration tests
go test ./tests/integration -v -podman -podman-basic
```

### CI Environment

The tests run automatically in GitHub Actions when:
- Code is pushed to `main` or `develop` branches
- Pull requests are opened against `main` or `develop`

## CI Workflow

The GitHub Actions workflow (`.github/workflows/integration-tests.yml`) performs the following steps:

1. **Setup Environment**
   - Install Go 1.24+
   - Install Podman and podman-compose
   - Configure Podman for rootless operation

2. **SSH Key Generation**
   - Generate Ed25519 SSH key pair without passphrase
   - Configure SSH agent

3. **Container Setup**
   - Build custom Debian 12 image with SSH server
   - Start container with SSH service on port 2222
   - Configure SSH authentication with generated key

4. **Test Execution**
   - Run integration tests with container environment
   - Test SSH connections, command execution, and script handling
   - Verify parallel and sequential execution

5. **Cleanup**
   - Stop and remove containers
   - Clean up images and artifacts

## Test Environment

### Container Configuration

The test container is based on `debian:12-slim` and includes:
- OpenSSH server
- Sudo access for test user
- Common utilities (curl, wget, vim)
- Proper SSH configuration for key-based authentication

### SSH Configuration

- **User**: `testuser`
- **Port**: 2222 (mapped from container port 22)
- **Authentication**: Ed25519 key-based (no password)
- **Host Key**: Insecure (for testing only)

### Environment Variables

The CI environment sets these variables for tests:
- `SPOOKY_TEST_SSH_HOST`: localhost
- `SPOOKY_TEST_SSH_PORT`: 2222
- `SPOOKY_TEST_SSH_USER`: testuser
- `SPOOKY_TEST_SSH_KEY`: ~/.ssh/id_ed25519

## Test Coverage

The integration tests verify:

1. **SSH Connection**
   - Basic SSH client creation and connection
   - Command execution over SSH
   - Connection error handling

2. **Configuration Processing**
   - HCL configuration file parsing
   - Server and action validation
   - Configuration execution

3. **Script Execution**
   - Local script file execution on remote servers
   - Script output capture and validation
   - Script error handling

4. **Parallel Execution**
   - Parallel action execution across multiple servers
   - Timing validation for parallel vs sequential
   - Concurrent connection handling

5. **Error Handling**
   - Invalid command execution
   - Connection failures
   - Configuration errors

6. **CLI Commands**
   - Validate command functionality
   - List command functionality
   - Execute command functionality

## Troubleshooting

### Local Test Issues

1. **Podman not available**
   ```bash
   # Install Podman
   sudo apt-get install podman podman-compose
   ```

2. **Permission issues**
   ```bash
   # Configure user namespaces
   echo 'kernel.unprivileged_userns_clone=1' | sudo tee -a /etc/sysctl.conf
   sudo sysctl -p
   ```

3. **SSH connection failures**
   - Verify container is running: `podman ps`
   - Check SSH service: `podman exec spooky-test-server systemctl status ssh`
   - Test manual connection: `ssh -p 2222 testuser@localhost`

### CI Test Issues

1. **Container build failures**
   - Check Dockerfile syntax
   - Verify base image availability
   - Review build logs for package installation errors

2. **SSH authentication failures**
   - Verify key generation and permissions
   - Check authorized_keys file in container
   - Review SSH configuration in container

3. **Test timeouts**
   - Increase timeout values in test configuration
   - Check container resource limits
   - Review network connectivity

## Security Considerations

- Test containers use insecure host key verification
- SSH keys are generated without passphrases
- Containers run with minimal security restrictions
- These settings are for testing only and should not be used in production

## Contributing

When adding new integration tests:

1. Follow the existing test structure and naming conventions
2. Use appropriate test flags for local vs CI execution
3. Include proper cleanup in test functions
4. Add documentation for new test scenarios
5. Ensure tests are idempotent and can run multiple times 