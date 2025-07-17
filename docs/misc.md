# Spooky - Detailed Documentation

This document contains detailed information about Spooky that doesn't fit in the main README.md.

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Security](#security)
- [Testing](#testing)
- [Contributing](#contributing)
- [Dependencies](#dependencies)
- [Project Assumptions](#project-assumptions)

## Installation

### Prerequisites

- Go 1.24 or later
- Target systems must be accessible via SSH

### Build from Source

```bash
git clone https://github.com/snassar/spooky.git
cd spooky
go mod tidy
go build -o spooky
```

## Configuration

### Configuration File Format

Spooky uses HCL2 configuration files to define servers and actions. Here's an example:

```hcl
# Define servers
server "web-server-1" {
  host     = "192.168.1.10"
  port     = 22
  user     = "admin"
  password = "your-password"
  # key_file = "~/.ssh/id_rsa"  # Alternative to password
  tags = {
    environment = "production"
    role        = "web"
  }
}

# Define actions
action "check-status" {
  description = "Check system status"
  command     = "uptime && df -h"
  servers     = ["web-server-1"]  # Target specific servers
  # tags = ["production"]         # Or target by tags
  parallel    = true              # Execute in parallel
}
```

### Server Configuration

Each server block defines a remote server:

- `name`: Unique identifier for the server
- `host`: Server IP address or hostname
- `port`: SSH port (default: 22)
- `user`: SSH username
- `password`: SSH password (or use `key_file`)
- `key_file`: Path to SSH private key file
- `tags`: Key-value pairs for server categorization

### Action Configuration

Each action block defines an operation to perform:

- `name`: Unique identifier for the action
- `description`: Human-readable description
- `command`: Inline command to execute
- `script`: Path to script file to execute
- `servers`: List of specific server names to target
- `tags`: List of tags to match servers
- `timeout`: Custom timeout for this action
- `parallel`: Execute on servers in parallel

### Server Targeting

Actions can target servers in three ways:

1. **Specific servers**: `servers = ["server1", "server2"]`
2. **Tag-based**: `tags = ["production", "web"]`
3. **All servers**: Omit both `servers` and `tags`

## Examples

### System Maintenance

```hcl
server "prod-web-1" {
  host     = "10.0.1.10"
  user     = "admin"
  password = "secure-password"
  tags = {
    environment = "production"
    role        = "web"
  }
}

server "prod-web-2" {
  host     = "10.0.1.11"
  user     = "admin"
  password = "secure-password"
  tags = {
    environment = "production"
    role        = "web"
  }
}

action "update-system" {
  description = "Update system packages"
  command     = "sudo apt update && sudo apt upgrade -y"
  tags        = ["production"]
  parallel    = true
}

action "restart-services" {
  description = "Restart web services"
  command     = "sudo systemctl restart nginx"
  tags        = ["web"]
  parallel    = true
}
```

### Database Operations

```hcl
server "db-master" {
  host     = "10.0.2.10"
  user     = "dbadmin"
  key_file = "~/.ssh/db_key"
  tags = {
    environment = "production"
    role        = "database"
  }
}

action "backup-database" {
  description = "Create database backup"
  script      = "./scripts/backup.sh"
  servers     = ["db-master"]
}
```

## Security Considerations

- Store sensitive information (passwords, keys) securely
- Use SSH key authentication when possible
- Limit SSH user permissions on target servers
- Consider using environment variables for credentials
- Regularly rotate passwords and keys

## Error Handling

Spooky provides detailed error reporting:

- Connection failures are reported per server
- Command execution errors include stderr output
- Configuration validation errors prevent execution
- Parallel execution continues even if some servers fail

## Testing

Spooky includes comprehensive unit and integration tests to help ensure reliability.

### Running Tests

You can run tests directly using Go commands:

```bash
# Run all tests (unit + integration)
go test ./... -tags=integration

# Run unit tests only (exclude integration tests)
go test ./...

# Run integration tests only
go test -tags=integration ./tests/integration/...

# Run tests with coverage report
go test -cover ./...
```

### Test Coverage Tool

For detailed coverage analysis, install the coverage tools:

```bash
make install-coverage-tools
```

Then run:

```bash
go test ./... -coverprofile=./tests/coverage.out -covermode=atomic -coverpkg=./...
```

Run the code coverage tool:

```bash
go-test-coverage --config=./tests/testcoverage.yml
```

### Test Structure

- **Unit Tests**: Co-located with source files (e.g., `config_test.go`)
- **Integration Tests**: Located in `tests/integration/` using [gliderlabs/ssh](https://github.com/gliderlabs/ssh) for mock SSH servers
- **Test Fixtures**: Sample configurations and scripts in `tests/fixtures/`
- **Test Helpers**: Common utilities in `tests/helpers/`

### Current Test Coverage

The test suite currently covers:
- Basic configuration parsing and validation
- SSH connection and authentication
- Simple command execution
- Basic error handling

**Note**: This is a work in progress. Additional test coverage for advanced features, edge cases, and parallel execution is planned.

For detailed testing information, see [docs/coverage.md](coverage.md).

### Coverage Visualization

#### Local Development
Generate an HTML coverage report locally:
```bash
make coverage-html
```
Then open `coverage.html` in your browser to view detailed coverage information.

#### CI/CD
Coverage reports are automatically generated in CI/CD and available as workflow artifacts:

1. Go to the **Actions** tab in GitHub
2. Click on a workflow run (e.g., "Test Coverage")
3. Scroll down to **Artifacts**
4. Download `coverage-reports` to get:
   - `tests/coverage.out` - Raw coverage data
   - `coverage.html` - Interactive HTML report

#### Viewing HTML Reports
The HTML coverage report provides:
- File-by-file coverage breakdown
- Line-by-line coverage highlighting
- Overall project coverage statistics
- Coverage trends and gaps

Open `coverage.html` in any web browser to explore coverage details.

### Test Coverage Requirements

This project enforces test coverage to ensure code quality.

#### Coverage Requirements

- **Total Coverage:** Minimum 60%
- **Package Coverage:** Minimum 65%
- **File Coverage:** Minimum 50%
- **Critical Code:** Higher thresholds for SSH and configuration code

#### Running Coverage Checks Locally

```bash
# Run all tests and check coverage thresholds
make check-coverage

# Generate HTML coverage report
make coverage-html

# Run coverage tool manually
go test ./... -coverprofile=./tests/coverage.out -covermode=atomic -coverpkg=./...
go run github.com/vladopajic/go-test-coverage/v2@latest --config=./tests/testcoverage.yml
```

#### Pre-commit Hook Setup

To automatically check coverage before each commit, install the pre-commit hook:

```bash
# Build and install the pre-commit hook
make install-pre-commit-hook
```

This will:
1. Build the Go-based pre-commit hook from `tools/pre-commit/main.go`
2. Install it as a Git hook in `.git/hooks/pre-commit`
3. Automatically run coverage checks before each commit

**Manual Setup:**
```bash
# Build the hook (cross-platform)
make build-pre-commit-hook

# Or build manually:
# Unix/Linux/macOS: go build -o build/pre-commit tools/pre-commit/main.go
# Windows: go build -o build/pre-commit.exe tools/pre-commit/main.go

# Install the hook (cross-platform)
make install-pre-commit-hook

# Or install manually:
# Unix/Linux/macOS: cp build/pre-commit .git/hooks/pre-commit
# Windows: copy build\pre-commit.exe .git\hooks\pre-commit
```

**What the Hook Does:**
- Runs tests with coverage profiling
- Verifies coverage meets thresholds
- Blocks commits if coverage is insufficient
- Allows commits if coverage passes

**Note:** The pre-commit hook is a pure Go solution with no external dependencies.

#### Viewing Coverage Reports

- Open `coverage.html` in your browser for a detailed report.
- Download coverage artifacts from GitHub Actions workflow runs.

#### CI/CD

- Coverage is checked on every PR and push to main.
- PRs that drop coverage below thresholds will fail.

#### Coverage Diff Tracking

- Coverage changes are automatically tracked in pull requests
- PR comments show coverage increases/decreases
- Coverage reports are available as workflow artifacts
- Decreases in coverage trigger warnings

For detailed coverage analysis, see [docs/coverage.md](coverage.md).

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (see above for Go test commands)
6. Submit a pull request

## Dependencies

### Core Dependencies
- [Cobra](https://github.com/spf13/cobra) - CLI framework for command-line interface
- [golang.org/x/crypto/ssh](https://pkg.go.dev/golang.org/x/crypto/ssh) - SSH client implementation
- [HCL2](https://github.com/hashicorp/hcl) - HashiCorp Configuration Language for config files

### Testing Dependencies
- [gliderlabs/ssh](https://github.com/gliderlabs/ssh) - SSH server for integration testing
- [github.com/pkg/sftp](https://github.com/pkg/sftp) - SFTP client for file transfer testing

### Indirect Dependencies
- [golang.org/x/sys](https://pkg.go.dev/golang.org/x/sys) - System calls and OS-specific functionality
- [golang.org/x/text](https://pkg.go.dev/golang.org/x/text) - Text processing utilities
- [zclconf/go-cty](https://github.com/zclconf/go-cty) - Type system for HCL2
- [apparentlymart/go-textseg](https://github.com/apparentlymart/go-textseg) - Text segmentation for HCL2
- [kr/fs](https://github.com/kr/fs) - File system utilities for SFTP

## Project Assumptions

### Core Assumptions
- SSH key authentication is preferred over passwords
- We currently only support Ed25519 keys only
- We currently do not support RSA keys
- We currently do not support DSA keys
- Windows compatibility is required
- macOS compatibility is required
- linux compatibility is required
- Integration tests use https://github.com/gliderlabs/ssh
- Unit tests must be co-located with source files
- Integration tests must be located under `./tests/`
- Example configuration files and snippets must be located under `./examples/`

### Design Decisions
- Use HCL2 for configuratio files
- Do not use JSON for configuration files
- Do not use YAML for configuration files
- Parallel execution is optional per action
- Server targeting via tags or explicit names
- No persistent state or database required

### Constraints
- Make is supported
- Must work without Make
- Focus on basic functionality first 