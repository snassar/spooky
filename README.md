# Spooky - SSH Automation Tool

[![Test Coverage](coverage.svg)](https://github.com/snassar/spooky)

Spooky is an automation tool written in Go that allows you to execute commands and scripts on multiple remote servers using HCL2 configuration files. It provides a declarative way to manage server operations with support for parallel execution and flexible server targeting. spooky uses SSH to communicate with other computers 

When spooky grows up it wants to be Ansible.

**Notice**: This project is influenced heavily by agentic coding as part of the process of learning go.

## Features

- 🚀 **Declarative Configuration**: Use HCL2 files to define servers and actions
- 🔗 **SSH Connectivity**: Secure connections with password or key-based authentication
- ⚡ **Parallel Execution**: Run actions on multiple servers simultaneously
- 🏷️ **Tag-based Targeting**: Target servers using tags for flexible grouping
- 📝 **Script Support**: Execute both inline commands and external script files
- ✅ **Validation**: Built-in configuration validation and syntax checking
- 🔍 **Listing**: View servers and actions defined in configuration files

## Installation

### Prerequisites

- Go 1.24 or later
- target systems must be accessible via SSH

### Build from Source

```bash
git clone https://github.com/snassar/spooky.git
cd spooky
go mod tidy
go build -o spooky
```

## Usage

### Basic Commands

```bash
# Execute actions from configuration file
./spooky execute config.hcl

# Validate configuration file
./spooky validate config.hcl

# List servers and actions in configuration
./spooky list config.hcl

# Execute with parallel processing
./spooky execute -p config.hcl

# Set custom timeout (in seconds)
./spooky execute -t 60 config.hcl
```

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

Spooky includes basic unit and integration tests to help ensure reliability.

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

For detailed coverage analysis, install the go-test-coverage tool:

```bash
go install github.com/vladopajic/go-test-coverage/v2@latest
```

Then run:

```bash
go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
```

Run the code coverage tool

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

For detailed testing information, see [tests/README.md](tests/README.md).

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
   - `cover.out` - Raw coverage data
   - `coverage.html` - Interactive HTML report

#### Viewing HTML Reports
The HTML coverage report provides:
- File-by-file coverage breakdown
- Line-by-line coverage highlighting
- Overall project coverage statistics
- Coverage trends and gaps

Open `coverage.html` in any web browser to explore coverage details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (see above for Go test commands)
6. Submit a pull request

## License

This project is licensed under the GNU Affero General Public License v3 - see the LICENSE file for details.

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

## Test Coverage

This project enforces test coverage to ensure code quality.

### Coverage Requirements

- **Total Coverage:** Minimum 60%
- **Package Coverage:** Minimum 65%
- **File Coverage:** Minimum 50%

### Running Coverage Checks Locally

```bash
# Run all tests and check coverage thresholds
make check-coverage

# Generate HTML coverage report
make coverage-html

# Run coverage tool manually
go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
go run github.com/vladopajic/go-test-coverage/v2@latest --config=./tests/testcoverage.yml
```

### Viewing Coverage Reports

- Open `coverage.html` in your browser for a detailed report.
- Download coverage artifacts from GitHub Actions workflow runs.

### CI/CD

- Coverage is checked on every PR and push to main.
- PRs that drop coverage below thresholds will fail.

## Troubleshooting

### Common Issues

#### Coverage Tool Problems

#### Tool Not Found
```bash
# Error: command not found: go-test-coverage
go install github.com/vladopajic/go-test-coverage/v2@latest
```

#### Permission Denied
```bash
# Error: permission denied
# On Windows, run PowerShell as Administrator
# On Linux/macOS, check file permissions
chmod +x $(go env GOPATH)/bin/go-test-coverage
```

#### Version Conflicts
```bash
# Check installed version
go-test-coverage --version

# Reinstall if needed
go install github.com/vladopajic/go-test-coverage/v2@latest
```

### Coverage Profile Issues

#### Profile Not Found
```bash
# Error: coverage profile not found
# Generate profile first
go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
```

#### Empty Profile
```bash
# Profile exists but is empty
# Check if tests are running
go test -v ./...

# Check if coverage mode is correct
go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
```

#### Profile Path Issues
```bash
# Error: cannot find coverage profile
# Use absolute path
go run github.com/vladopajic/go-test-coverage/v2@latest \
  --config=./tests/testcoverage.yml \
  --profile=$(pwd)/cover.out
```

### Threshold Failures

#### File Threshold Failed
```bash
# Error: file coverage below threshold
# Check specific file coverage
go tool cover -func=./cover.out | grep "filename.go"

# Add tests for uncovered functions
# Or add coverage-ignore comment if appropriate
```

#### Package Threshold Failed
```bash
# Error: package coverage below threshold
# Check package coverage
go tool cover -func=./cover.out

# Focus on high-priority functions
# Add integration tests if unit tests are insufficient
```

#### Total Threshold Failed
```bash
# Error: total coverage below threshold
# Review overall coverage
go tool cover -func=./cover.out

# Identify largest uncovered areas
# Prioritize testing based on business impact
```

### Configuration Issues

#### Config File Not Found
```bash
# Error: config file not found
# Check file exists
ls -la ./tests/testcoverage.yml

# Use absolute path
go run github.com/vladopajic/go-test-coverage/v2@latest \
  --config=$(pwd)/tests/testcoverage.yml
```

#### Invalid Config Syntax
```bash
# Error: invalid YAML
# Validate YAML syntax
yamllint ./tests/testcoverage.yml

# Check indentation and format
```

#### Exclusion Patterns Not Working
```bash
# Files still included despite exclusions
# Check regex patterns
# Test patterns manually
echo "filename.go" | grep -E "pattern"

# Use debug mode
go run github.com/vladopajic/go-test-coverage/v2@latest \
  --config=./tests/testcoverage.yml \
  --debug
```

### Performance Issues

#### Slow Coverage Generation
```bash
# Coverage generation is slow
# Exclude unnecessary files
# Use more specific package patterns
go test ./spooky/... -coverprofile=./cover.out -covermode=atomic -coverpkg=./spooky/...
```

#### Large Coverage Files
```bash
# Coverage file is very large
# Check what's being included
go tool cover -func=./cover.out | wc -l

# Review exclusions
# Remove unnecessary packages from coverage
```

### Integration Test Issues

#### Integration Tests Not Running
```bash
# Error: no tests found
# Check build tags
go test -tags=integration ./tests/integration/...

# Verify test files have correct tags
# //go:build integration
```

#### Mock Server Issues
```bash
# Integration tests failing
# Check mock server ports
# Verify SSH server is running
netstat -an | grep :3100

# Restart mock servers
make test-integration
```

## Getting Help

### Debug Mode
Always run with debug mode when troubleshooting:
```bash
go run github.com/vladopajic/go-test-coverage/v2@latest \
  --config=./tests/testcoverage.yml \
  --debug
```

### Logs and Output
- Check console output for error messages
- Review coverage HTML report for details
- Check GitHub Actions logs for CI issues

### Common Solutions
1. **Reinstall tools**: `go install github.com/vladopajic/go-test-coverage/v2@latest`
2. **Regenerate profiles**: `go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...`
3. **Check paths**: Use absolute paths for config and profile files
4. **Validate config**: Check YAML syntax and regex patterns
5. **Update exclusions**: Add appropriate exclusions for untestable code

## Coverage Thresholds and Requirements

### Threshold Configuration

Coverage thresholds are defined in `./tests/testcoverage.yml`:

```yaml
threshold:
  file: 50      # Individual files must have at least 50% coverage
  package: 65   # Packages must have at least 65% coverage
  total: 60     # Overall project must have at least 60% coverage
```

### Threshold Rationale

#### File-Level Threshold (50%)
- **Why 50%**: Some files contain error handling that's hard to test
- **Examples**: CLI validation, file system operations
- **Exclusions**: Files with `// coverage-ignore` comments

#### Package-Level Threshold (65%)
- **Why 65%**: Ensures core functionality is well-tested
- **Focus**: Business logic and critical paths
- **Current**: Project achieves ~61.7%, so 65% is achievable

#### Total Project Threshold (60%)
- **Why 60%**: Balances quality with practicality
- **Realistic**: Accounts for hard-to-test code
- **Maintainable**: Can be achieved and maintained

### Coverage Requirements by Code Type

#### High Priority (Must Test)
- SSH connection logic
- Configuration parsing
- Command execution
- Error handling paths

#### Medium Priority (Should Test)
- CLI argument parsing
- File operations
- Utility functions

#### Low Priority (Optional)
- Main entry points (tested via integration)
- Error exit paths (tested via integration)

### Exclusions and Ignored Code

#### Automatic Exclusions
- Test files (`_test.go`)
- Test infrastructure (`tests/`)
- Example files (`examples/`)
- Main entry point (`main.go`)

#### Manual Exclusions
- Code marked with `// coverage-ignore`
- CLI entry points
- File system error handling
- Network error handling

### Threshold Enforcement

#### CI/CD Pipeline
- Coverage checks run on every PR
- Workflow fails if thresholds not met
- Prevents merging code that reduces coverage

#### Local Development
- `make check-coverage` validates thresholds
- Immediate feedback on coverage issues
- HTML reports for detailed analysis

### Updating Thresholds

#### When to Increase
- Project has matured
- Test coverage has improved
- New testing patterns established

#### When to Decrease
- Thresholds are unrealistic
- Project scope has changed
- New code types added

#### Process
1. Update `./tests/testcoverage.yml`
2. Update documentation
3. Communicate changes to team
4. Monitor impact on development velocity

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (see above for Go test commands)
6. Submit a pull request

## License

This project is licensed under the GNU Affero General Public License v3 - see the LICENSE file for details.

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