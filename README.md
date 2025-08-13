# Spooky

A powerful automation and orchestration tool built with Go.

## Current Status

This is a work-in-progress implementation of spooky. The CLI is now built using [Cobra](https://github.com/spf13/cobra), the de facto standard for Go CLI applications, providing a robust and feature-rich command-line interface.

### ✅ What's Working

The CLI is fully functional and supports comprehensive project management, machine inventory, variable management, and facts collection. For a complete reference of all available commands, see the [CLI Reference](docs/CLI_REFERENCE.md).

**Available Command Categories:**
- **Project Management**: `spooky project init`, `spooky project validate`
- **Machine Inventory**: `spooky machines list`, `spooky machines validate`, `spooky machines ping`
- **Variable Management**: `spooky variables list`, `spooky variables validate`, `spooky variables resolve`
- **Facts Management**: `spooky facts export`

**Quick Examples:**
```bash
# Initialize a new project
spooky project init my-automation-project

# List machines in project
spooky machines list ./my-automation-project

# Export facts from machines
spooky facts export ./my-automation-project --output facts.hcl
```

### 🚀 Enhanced Features with Cobra

The CLI now includes advanced features provided by Cobra:

- **Subcommand-based CLI** with proper nesting (`spooky project init`, `spooky project validate`)
- **POSIX-compliant flags** with short and long versions (`-h`, `--help`, `-v`, `--version`)
- **Automatic help generation** for all commands and subcommands
- **Shell autocomplete** support (bash, zsh, fish, powershell)
- **Intelligent suggestions** for command typos
- **Command aliases** for backward compatibility
- **Proper error handling** with meaningful error messages
- **Argument validation** with clear usage instructions

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
```

For complete command documentation and examples, see the [CLI Reference](docs/CLI_REFERENCE.md).

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
- `internal/types/` - Type definitions organized by domain
- `internal/interfaces/` - Core system interfaces
- `internal/cli/` - CLI implementation (planned)
- `internal/logging/` - Logging implementation (planned)
- `internal/schemas/` - Schema management (planned)
- `internal/storage/` - Storage implementation (planned)

## Development

This project is built following established patterns and rules:

- **Interface-first architecture** with Cobra command structure
- **Type safety and validation** with proper argument handling
- **Comprehensive error handling** with meaningful error messages
- **Structured logging** (planned)
- **Schema-driven configuration** (planned)
- **Component-based design** with proper separation of concerns

## Next Steps

The following components need to be implemented:

1. **Project Management** - Actual project initialization and validation logic
2. **Schema System** - HCL schema loading and validation using the existing schemas
3. **Configuration Management** - Configuration loading and validation
4. **Logging System** - Structured logging implementation
5. **Storage System** - Data persistence layer
6. **Integration Layer** - Component coordination

## Requirements

- Go 1.24 or later
- HCL v2 for configuration parsing
- Cobra v1.9+ for CLI framework

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
