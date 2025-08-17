# ADR Generator Tool

## Purpose
Generates ADR (Architecture Decision Record) documentation from the codebase.

## Usage
```bash
# Build the tool
just build

# Generate ADR documentation
just run
```

## Output
Generates ADR documentation in the `docs/adr/` directory.

## Configuration
This tool analyzes the codebase to generate ADR recommendations and documentation.

## Dependencies
- Go 1.24+
- Access to the spooky codebase
