# Testing Guide

This directory contains the test suite for the Spooky SSH automation tool.

## Test Structure

```
tests/
├── integration/           # Integration tests using gliderlabs/ssh
│   ├── config_integration_test.go
│   └── ssh_integration_test.go
├── fixtures/              # Test data and configuration files
│   ├── valid_config.hcl
│   └── invalid_config.hcl
├── helpers/               # Test utilities and helpers
│   └── test_helpers.go
└── README.md             # This file
```

## Test Types

### Unit Tests
- **Location**: Co-located with source files (e.g., `config_test.go` next to `config.go`)
- **Purpose**: Test individual functions and methods in isolation
- **Dependencies**: Minimal external dependencies
- **Speed**: Fast execution

### Integration Tests
- **Location**: `tests/integration/`
- **Purpose**: Test component interactions and end-to-end workflows
- **Dependencies**: Uses [gliderlabs/ssh](https://github.com/gliderlabs/ssh) for mock SSH servers
- **Speed**: Slower due to network operations and server setup

## Running Tests

### All Tests
```bash
make test
```

### Unit Tests Only
```bash
make test-unit
```

### Integration Tests Only
```bash
make test-integration
```

### Tests with Coverage
```bash
make test-coverage
```

### Individual Test Files
```bash
# Unit tests
go test -v ./config_test.go ./config.go

# Integration tests
go test -v -tags=integration ./tests/integration/config_integration_test.go
```

## Integration Test Features

### Mock SSH Server
The integration tests use [gliderlabs/ssh](https://github.com/gliderlabs/ssh) to create mock SSH servers that:

- Accept any password/key for testing
- Record executed commands
- Return predefined outputs
- Support custom command responses
- Run on dynamic ports to avoid conflicts

### Test Scenarios Covered

1. **Basic SSH Connection**: Test connection establishment
2. **Command Execution**: Test command execution and output capture
3. **Script Execution**: Test script file execution
4. **Sequential Execution**: Test executing actions on multiple servers sequentially
5. **Parallel Execution**: Test executing actions on multiple servers in parallel
6. **Error Handling**: Test error scenarios and edge cases
7. **Authentication Methods**: Test password and key-based authentication
8. **Timeout Handling**: Test connection and command timeouts

## Test Fixtures

### Configuration Files
- `valid_config.hcl`: Sample valid configuration for testing
- `invalid_config.hcl`: Sample invalid configuration for error testing

### Helper Functions
The `test_helpers.go` file provides utilities for:

- Generating SSH key pairs for testing
- Creating temporary configuration files
- Creating temporary script files
- Sample configuration content
- Sample script content
- Port availability checking
- Test file cleanup
- Common assertions

## Build Tags

Integration tests use the `integration` build tag to separate them from unit tests:

```go
//go:build integration
// +build integration
```

This allows running unit tests quickly during development while keeping integration tests separate for CI/CD pipelines.

## Best Practices

1. **Isolation**: Each test should be independent and not rely on other tests
2. **Cleanup**: Always clean up resources (files, connections, servers) after tests
3. **Timeouts**: Use appropriate timeouts for network operations
4. **Mocking**: Use mock SSH servers instead of real servers for consistent testing
5. **Fixtures**: Use test fixtures for consistent test data
6. **Assertions**: Use descriptive assertions with clear error messages

## Troubleshooting

### Integration Tests Failing
- Ensure no other SSH servers are running on the test ports
- Check that the gliderlabs/ssh dependency is properly installed
- Verify network connectivity for localhost connections

### Port Conflicts
- Integration tests use dynamic port allocation to avoid conflicts
- If you encounter port issues, check for other services using the same ports

### Test Timeouts
- Increase timeout values in test configuration if needed
- Check system resources and network connectivity

## Contributing

When adding new tests:

1. Follow the existing naming conventions
2. Use the provided helper functions
3. Add appropriate build tags
4. Include both positive and negative test cases
5. Update this README if adding new test categories 

# Integration Tests

This directory contains integration tests that verify the functionality of the SSH infrastructure servers.

## Test Structure

- `integration_test.go` - Main orchestration test that runs all server tests
- `simple_server_test.go` - Tests for the simple SSH server
- `public_key_server_test.go` - Tests for the public key SSH server  
- `sftp_server_test.go` - Tests for the SFTP server
- `timeout_server_test.go` - Tests for the timeout SSH server
- `helpers.go` - Helper functions for SSH and SFTP testing

## Running the Tests

### Prerequisites

Make sure all infrastructure servers have their dependencies installed:

```bash
# Install dependencies for each server
cd tests/infrastructure/simple_server && go mod tidy
cd tests/infrastructure/public_key_server && go mod tidy  
cd tests/infrastructure/sftp_server && go mod tidy
cd tests/infrastructure/timeout_server && go mod tidy
```

### Running All Tests

From the `tests` directory:

```bash
go test -v
```

### Running Individual Tests

You can run individual tests using command-line flags:

```bash
# Test simple server only
go test -v -args -simple

# Test public key server only
go test -v -args -publickey

# Test SFTP server only
go test -v -args -sftp

# Test timeout server only
go test -v -args -timeout
```

Or using the traditional Go test pattern matching:

```bash
# Test simple server only
go test -v -run TestSimpleServer

# Test public key server only
go test -v -run TestPublicKeyServer

# Test SFTP server only
go test -v -run TestSFTPServer

# Test timeout server only
go test -v -run TestTimeoutServer
```

### Running Specific Test Functions

```bash
# Test multiple connections to simple server
go test -v -run TestSimpleServerMultipleConnections

# Test timeout behavior
go test -v -run TestTimeoutServerConnectionTimeout
```

## What the Tests Do

### Simple Server Tests
- Verifies basic SSH connection and session creation
- Tests that the server responds with "Hello {username}"
- Tests multiple concurrent connections

### Public Key Server Tests  
- Generates test Ed25519 key pairs
- Tests SSH connection with public key authentication
- Verifies the server returns the public key information

### SFTP Server Tests
- Tests SFTP connection establishment
- Performs basic SFTP operations (create file, read file, list directory)
- Tests file upload/download functionality

### Timeout Server Tests
- Tests connection establishment and session creation
- Verifies the server keeps connections alive for the expected duration
- Tests timeout behavior (optional, may be flaky in CI)

## Test Infrastructure

Each test:
1. Starts the appropriate infrastructure server as a subprocess
2. Waits for the server to be ready
3. Connects to the server using SSH client
4. Performs the test operations
5. Cleans up by stopping the server process

## Troubleshooting

### Port Conflicts
If you get "address already in use" errors, make sure no other SSH servers are running on port 2222.

### Dependency Issues
If you see import errors, run `go mod tidy` in the tests directory and in each infrastructure server directory.

### Timeout Issues
Some tests may be flaky in CI environments due to timing. The timeout server tests in particular may need adjustment for different environments.
