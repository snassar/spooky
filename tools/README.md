# Developer Tools

This directory contains all developer tools for the spooky project. Each tool is self-contained with its own justfile and documentation.

## Tool Organization

### Core Development Tools
- **dependency-graph** - Analyze package dependencies and generate visual graphs
- **terminology-checker** - Check terminology compliance across the codebase
- **pre-commit** - Pre-commit hooks and checks

### Documentation Tools
- **adr-generator** - Generate ADR (Architecture Decision Record) documentation
- **git-adr-analyzer** - Analyze ADRs from Git history
- **focused-adr-analyzer** - Perform focused ADR analysis
- **simple-adr-generator** - Generate simple ADR documentation

### Code Quality Tools
- **todo-linter** - Lint TODO comments in the codebase

## Usage

### Individual Tools
Each tool can be used independently:

```bash
# Build a specific tool
just --directory tools/dependency-graph build

# Run a specific tool
just --directory tools/dependency-graph run

# Get help for a specific tool
just --directory tools/dependency-graph help
```

### Tool Orchestration
Use the tools orchestration justfile to manage all tools:

```bash
# Build all tools
just --directory tools build

# Run all tools
just --directory tools run-all

# Clean all tools
just --directory tools clean

# Show help for all tools
just --directory tools help
```

### Main Project Integration
The main project justfile integrates with tools:

```bash
# Build all developer tools
just build-tools

# Run all developer tools
just run-tools

# Clean all developer tools
just clean-tools

# Show tools help
just tools-help
```

## Tool Structure

Each tool follows this structure:
```
tools/tool-name/
├── main.go          # Tool implementation
├── justfile         # Tool-specific tasks
├── README.md        # Tool documentation
├── go.mod           # Tool dependencies (if needed)
└── config/          # Tool configuration (if needed)
```

## Adding New Tools

To add a new tool:

1. Create a new directory in `tools/`
2. Add the tool implementation as `main.go`
3. Create a `justfile` with build, run, clean, and help targets
4. Add a `README.md` with tool documentation
5. Update `tools/justfile` to include the new tool
6. Update `tools/README.md` to document the new tool

## Tool Standards

All tools should:
- Be self-contained with their own dependencies
- Have a consistent justfile interface (build, run, clean, help)
- Include proper documentation
- Follow the project's coding standards
- Be testable and maintainable

## Dependencies

Tools may have their own dependencies separate from the main project. Each tool can have its own `go.mod` file if needed.

## Output

Tools typically generate output in the `docs/` directory or other appropriate locations. Check each tool's README for specific output locations.
