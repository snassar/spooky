# Development Tools

[← Back to Index](index.md) | [Next: FAQ](faq.md)

---

## Pre-commit Hook

The `pre-commit/` directory contains a Go tool that runs test coverage checks before each commit.

### Building the Pre-commit Hook

```bash
# From the project root (cross-platform)
make build-pre-commit-hook

# Manual build (Unix/Linux/macOS)
go build -o build/pre-commit tools/pre-commit/main.go

# Manual build (Windows)
go build -o build/pre-commit.exe tools/pre-commit/main.go

# Or run directly
go run tools/pre-commit/main.go
```

### Installing the Pre-commit Hook

After building, install the hook:

```bash
# Automated install (cross-platform)
make install-pre-commit-hook

# Manual install (Unix/Linux/macOS)
cp build/pre-commit .git/hooks/pre-commit

# Manual install (Windows)
copy build\pre-commit.exe .git\hooks\pre-commit
```

### What the Hook Does

The pre-commit hook:
1. Checks if you're in a git repository
2. Identifies staged Go files
3. Runs tests with coverage profiling
4. Verifies coverage meets thresholds defined in `tests/testcoverage.yml`
5. Blocks commits if coverage is insufficient
6. Allows commits if coverage passes

### Manual Usage

You can also run the coverage check manually:

```bash
go test ./... -coverprofile=./tests/coverage.out -covermode=atomic -coverpkg=./... && go run github.com/vladopajic/go-test-coverage/v2@latest --config=./tests/testcoverage.yml
```

## Test Environment Tool

The `spooky-test-env/` directory contains a unified Go tool for managing the Podman-based test environment.

### Building Test Environment Tool

```bash
# Build the test environment tool
make build-test-env

# Use the built binary
./build/spooky-test-env preflight
./build/spooky-test-env start
./build/spooky-test-env status
./build/spooky-test-env stop
./build/spooky-test-env cleanup
./build/spooky-test-env --help
```

### Available Commands

- **`preflight`**: Check that podman, systemd user session, and Quadlet are available
- **`start`**: Create network, start containers, and display IP addresses
- **`stop`**: Stop containers using Podman
- **`status`**: Show current status of containers and networks
- **`cleanup`**: Remove containers, networks, and all resources

The tool provides:
- Unified interface for all test environment operations
- Prerequisite checking before starting
- Status monitoring and reporting
- Complete cleanup capabilities
- Error handling with informative messages

## Configuration Generator

The `generate-config/` directory contains a tool for generating large-scale configuration files for testing.

### Building Configuration Generator

```bash
# Build the configuration generator
make build-generate-config

# Or run directly
go run tools/generate-config/main.go
```

### Using Configuration Generator

```bash
# Generate configuration files
go run tools/generate-config/main.go

# This creates:
# - examples/configuration/small-scale-example.hcl
# - examples/configuration/medium-scale-example.hcl
# - examples/configuration/large-scale-example.hcl
```

## Directory Structure

```
tools/
├── spooky-test-env/     # Test environment management
│   └── main.go
├── pre-commit/          # Pre-commit hook for coverage checks
│   └── main.go
├── generate-config/     # Configuration file generator
│   └── main.go
└── README.md           # This file
```

---

## Navigation
- [← Back to Index](index.md)
- [Previous: Contributing](contributing.md)
- [Next: FAQ](faq.md) 