# Spooky

[![Test Coverage](coverage.svg)](https://github.com/snassar/spooky)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-AGPL%203.0-green.svg)](LICENSE)

> SSH automation tool written in Go that executes commands and scripts on multiple remote servers using HCL2 configuration files.

Spooky is a lightweight alternative to Ansible with declarative configuration and parallel execution capabilities. When spooky grows up it wants to be Ansible.

## ✨ Features

- 🚀 **Declarative Configuration** - Use HCL2 files to define servers and actions
- 🔗 **SSH Connectivity** - Secure connections with password or key-based authentication  
- ⚡ **Parallel Execution** - Run actions on multiple servers simultaneously
- 🏷️ **Tag-based Targeting** - Target servers using tags for flexible grouping
- 📝 **Script Support** - Execute both inline commands and external script files
- ✅ **Validation** - Built-in configuration validation and syntax checking

## 🚀 Quick Start

### Install

```bash
git clone https://github.com/snassar/spooky.git
cd spooky
go build -o spooky
```

### Basic Usage

```bash
# Execute actions from configuration file
./spooky execute config.hcl

# Validate configuration file
./spooky validate config.hcl

# List servers and actions in configuration
./spooky list config.hcl
```

### Example Configuration

```hcl
# Define servers
server "web-server-1" {
  host     = "192.168.1.10"
  user     = "admin"
  password = "your-password"
  tags = {
    environment = "production"
    role        = "web"
  }
}

# Define actions
action "check-status" {
  description = "Check system status"
  command     = "uptime && df -h"
  servers     = ["web-server-1"]
  parallel    = true
}
```

## 📖 Documentation

- **[Detailed Documentation](docs/misc.md)** - Installation, configuration, testing, and more
- **[Troubleshooting](docs/troubleshooting.md)** - Common issues and solutions
- **[Test Environment](docs/test-environment.md)** - Setting up test environments
- **[Tools](docs/tools.md)** - Development tools and utilities
- **[Testing](docs/coverage.md)** - Testing guidelines and coverage requirements
- **[Coverage](docs/coverage.md)** - Test coverage analysis and requirements

## 🛠️ Development

### Prerequisites

- Go 1.24 or later
- Target systems accessible via SSH

### Testing

```bash
# Run all tests
go test ./... -tags=integration

# Run with coverage
make check-coverage

# Generate coverage report
make coverage-html
```

### Test Environment

```bash
# Check requirements
go run tools/spooky-test-env/main.go preflight

# Start test environment
go run tools/spooky-test-env/main.go start
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## 📄 License

This project is licensed under the GNU Affero General Public License v3 - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- **[Issues](https://github.com/snassar/spooky/issues)** - Report bugs or request features
- **[Discussions](https://github.com/snassar/spooky/discussions)** - Ask questions and share ideas
- **[Releases](https://github.com/snassar/spooky/releases)** - Latest releases and changelog

---

**Notice**: This project is influenced heavily by agentic coding as part of the process of learning Go.
