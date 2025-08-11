# Spooky

A powerful automation and orchestration tool built with Go.

## Current Status

This is a work-in-progress implementation of spooky. The CLI is now built using [Cobra](https://github.com/spf13/cobra), the de facto standard for Go CLI applications, providing a robust and feature-rich command-line interface.

### ✅ What's Working

The basic CLI structure is in place and supports the following commands:

- `spooky --version` - Show version information
- `spooky --help` - Show comprehensive help information
- `spooky project init <project>` - Initialize a new spooky project (placeholder implementation)
- `spooky project validate <project>` - Validate a spooky project (placeholder implementation)

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

To build spooky:

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

# Project commands
./spooky project init my-project
./spooky project validate my-project

# Get help for specific commands
./spooky project --help
./spooky project init --help
```

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

## Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [HCL v2](https://github.com/hashicorp/hcl) - Configuration parsing
