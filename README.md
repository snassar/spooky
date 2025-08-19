# Spooky

A powerful SSH automation and orchestration tool built with Go, designed for declarative configuration and parallel execution across heterogeneous environments.

## Overview

Spooky is a modern automation tool that combines the power of declarative configuration with intelligent fact-driven decision making. It provides a robust CLI built with [Cobra](https://github.com/spf13/cobra) and supports HCL-based configuration for projects, machines, actions, and variables.

## Current Status

This is an active development implementation of spooky with a fully functional CLI and comprehensive automation capabilities.

### ✅ What's Working

The CLI is fully functional and supports comprehensive project management, machine inventory, variable management, facts collection, and action orchestration. For a complete reference of all available commands, see the [CLI Reference](docs/systems/CLI_REFERENCE.md).

**Available Command Categories:**
- **Project Management**: `spooky project init`, `spooky project validate`, `spooky project encrypt`
- **Machine Inventory**: `spooky machines list`, `spooky machines validate`, `spooky machines ping`, `spooky machines export`, `spooky machines encrypt`
- **Variable Management**: `spooky variables list`, `spooky variables validate`, `spooky variables resolve`, `spooky variables armor`, `spooky variables encrypt`
- **Facts Management**: `spooky facts export`
- **Actions Management**: `spooky actions list`, `spooky actions validate`, `spooky actions run`
- **Templates Management**: `spooky templates list`, `spooky templates validate`, `spooky templates render`, `spooky templates search`
- **Secrets Management**: `spooky secrets validate`
- **Schemas**: `spooky schemas list`, `spooky schemas validate`
- **Integrations**: `spooky integrations list`, `spooky integrations validate`

**Quick Examples:**
```bash
# Initialize a new project
spooky project init my-automation-project

# List machines in project
spooky machines list ./my-automation-project

# Run actions on machines
spooky actions run ./my-automation-project --dry-run

# Export facts from machines
spooky facts export ./my-automation-project --output facts.hcl

# Validate project configuration
spooky project validate ./my-automation-project
```

### 🚀 Key Features

- **Declarative Configuration**: HCL-based configuration for projects, machines, actions, and variables
- **Parallel Execution**: Run actions across multiple machines simultaneously
- **Fact-Driven Decisions**: Collect and use machine facts for intelligent automation
- **Template Rendering**: Powerful template system with variable substitution
- **SSH Integration**: Robust SSH client with connection pooling and authentication
- **Schema Validation**: Comprehensive validation using embedded HCL schemas
- **Secrets Management**: Age encryption integration for secure data handling
- **Cross-Platform**: Support for Linux, macOS, and Windows

## Building

### Using just (Recommended)

The project uses [just](https://github.com/casey/just) for build automation and [ScalVer](https://scalver.org/) for version management:

```bash
# Build with ScalVer (daily cadence: MAJOR.YYYYMMDD.PATCH)
just build

# Build with yearly ScalVer cadence (MAJOR.YYYY.0)
just build-yearly

# Build with monthly ScalVer cadence (MAJOR.YYYYMM.0)
just build-monthly

# Build for development (dev version)
just build-dev

# Build with specific ScalVer version
just build-scalver 1 20250812 5

# Run all tests and build
just dev

# Clean build artifacts
just clean

# Show ScalVer information
just version

# Show available commands
just --list
```

### Manual Build

To build spooky manually:

```bash
go build -o spooky .
```

## Running

To run spooky:

```bash
# Show help
./spooky --help

# Show version
./spooky --version

# Get help for specific commands
./spooky project --help
./spooky machines --help
./spooky variables --help
./spooky facts --help
./spooky actions --help
./spooky templates --help
./spooky secrets --help
./spooky schemas --help
./spooky integrations --help
```

For complete command documentation and examples, see the [CLI Reference](docs/systems/CLI_REFERENCE.md).

## Shell Completion

Spooky supports shell completion for better user experience:

```bash
# Generate bash completion
./spooky completion bash > ~/.local/share/bash-completion/completions/spooky

# Generate zsh completion
./spooky completion zsh > ~/.zsh/completions/_spooky

# Generate fish completion
./spooky completion fish > ~/.config/fish/completions/spooky.fish
```

## Architecture

The project follows interface-based architecture patterns with the following structure:

- `cmd/` - Cobra command implementations
  - `root.go` - Root command and CLI setup
  - `project.go` - Project management commands
  - `machines.go` - Machine inventory commands
  - `variables.go` - Variable management commands
  - `facts.go` - Facts collection commands
  - `actions.go` - Action orchestration commands
  - `templates.go` - Template management commands
  - `secrets.go` - Secrets management commands
  - `schemas.go` - Schema validation commands
  - `integrations.go` - Integration management commands
- `internal/` - Core application code
  - `types/` - Type definitions organized by domain
  - `interfaces/` - Core system interfaces
  - `cli/` - CLI implementation
  - `logging/` - Logging implementation
  - `schemas/` - Schema management and validation
  - `config/` - Configuration management
  - `project/` - Project management
  - `machines/` - Machine inventory management
  - `variables/` - Variable management
  - `facts/` - Facts collection and storage
  - `actions/` - Action orchestration
  - `templates/` - Template rendering
  - `secrets/` - Secrets management
  - `ssh/` - SSH client and connection management
  - `integration/` - Component coordination

## Development

This project is built following established patterns and rules:

- **Interface-first architecture** with comprehensive interface definitions
- **Type safety and validation** with proper argument handling
- **Comprehensive error handling** with structured error types
- **Structured logging** with secure logging practices
- **Schema-driven configuration** with embedded HCL schemas
- **Component-based design** with proper separation of concerns
- **Security-first approach** with SSH key validation and secrets management

## Project Structure

A typical spooky project follows this structure:

```
my-automation-project/
├── project.hcl          # Project configuration
├── machines.hcl         # Machine inventory
├── actions.hcl          # Action definitions
├── variables.hcl        # Variable definitions
├── templates/           # Template files
│   ├── deploy.sh.tmpl
│   └── config.yaml.tmpl
├── variables/           # Variable files (optional)
│   ├── production.hcl
│   └── staging.hcl
└── schemas/             # Custom schemas (optional)
    └── custom.hcl
```

## Requirements

- Go 1.24 or later
- HCL v2 for configuration parsing
- Cobra v1.9+ for CLI framework
- SSH access to target machines
- Age encryption tools (for secrets management)

## Versioning

This project uses [ScalVer](https://scalver.org/) (Scalable Calendar Versioning) for version management. ScalVer combines the benefits of SemVer and CalVer:

- **Calendar-aware**: You know when something was released
- **SemVer compatible**: All existing tooling works unchanged
- **Adjustable cadence**: Can adapt from yearly → monthly → daily releases
- **Date-Only-Grows (DOG)**: Prevents version confusion

### ScalVer Format

The version format is `MAJOR.DATE.PATCH` where:
- **MAJOR**: Mirrors SemVer's MAJOR component (0 = alpha/experimental, 1+ = stable)
- **DATE**: Calendar date in UTC (YYYY, YYYYMM, or YYYYMMDD)
- **PATCH**: Monotonically increasing counter for backward-compatible updates

### Examples

- `0.2025.0` - First 2025 alpha release (yearly cadence)
- `0.202508.0` - First August 2025 release (monthly cadence)
- `0.20250812.0` - First August 12, 2025 release (daily cadence)
- `0.20250812.5` - Fifth August 12, 2025 release (same day)

## Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [HCL v2](https://github.com/hashicorp/hcl) - Configuration parsing
- [just](https://github.com/casey/just) - Build automation
- [gliderlabs/ssh](https://github.com/gliderlabs/ssh) - SSH client library
- [age](https://github.com/FiloSottile/age) - Encryption library

## Documentation

- [CLI Reference](docs/systems/CLI_REFERENCE.md) - Complete command reference
- [User Guides](docs/systems/USER_GUIDES_INDEX.md) - Comprehensive user documentation
- [API Reference](docs/systems/) - System-specific API documentation
- [Development Guide](docs/DEVELOPMENT.md) - Development setup and guidelines

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
