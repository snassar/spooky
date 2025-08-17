# Terminology Checker

A custom Go tool that enforces terminology standards in the spooky codebase by detecting banned terminology patterns.

## Overview

The terminology checker scans Go source files for banned terminology and reports violations. It's designed to enforce the terminology rules defined in `.cursor/rules/terminology.mdc`.

## Banned Terminology

The following terms and their variants are banned:

- `exec`, `execute`, `execution`
- `executor`, `executing`

## Allowed Patterns

The following patterns are allowed despite containing banned terms:

- `os/exec` - Standard library imports
- `exec.Command` - Standard library usage
- `.Execute()` - Framework methods (cobra, templates)
- `"exec"` - String literals for pattern matching
- Security pattern configurations (forbiddenPatterns arrays)
- Template content checking

## Usage

### Build the tool
```bash
just -f tools/terminology-checker/justfile build
```

### Run the checker
```bash
just -f tools/terminology-checker/justfile lint
```

### Run with verbose output
```bash
just -f tools/terminology-checker/justfile lint-verbose
```

### Run from main project
```bash
just lint-terminology
```

## Integration

The terminology checker is integrated into the main project's quality checks:

```bash
just check  # Includes terminology checking
```

## Output Format

The tool outputs violations in the following format:
```
file.go:line:column: message in context: text
```

Example:
```
cmd/actions.go:45:12: Use 'run' or 'orchestrate' instead of 'execute' in identifier: ExecuteAction
```

## Exit Codes

- `0` - No violations found
- `1` - Violations found
- `2` - Error occurred

## Configuration

The tool is configured through the patterns defined in `main.go`:

- `bannedTerms` - Defines banned terminology patterns
- `allowedPatterns` - Defines allowed patterns that contain banned terms
- `exactAllowed` - Defines exact string matches that should be allowed
