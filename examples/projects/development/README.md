# Development Environment

This directory contains development-specific files and databases for the Spooky project.

## Facts Database

The facts database (`facts.db`) is stored in this directory to avoid cluttering the project root. This follows best practices for development environments.

### Configuration

The facts database path can be configured using the `SPOOKY_FACTS_PATH` environment variable:

```bash
# For development/testing in the git repository
export SPOOKY_FACTS_PATH="examples/projects/development/facts.db"

# For production use (default if not set)
# ~/.local/state/spooky/facts.db
```

### Usage

The facts database is automatically used by all Spooky facts commands:

- `spooky facts gather` - Collect facts from machines
- `spooky facts query` - Query stored facts
- `spooky facts export` - Export facts to JSON
- `spooky facts import` - Import facts from JSON
- `spooky facts validate` - Validate stored facts
- `spooky facts list` - List all facts

### File Location

- **Default**: `examples/projects/development/facts.db`
- **Configurable**: Set `SPOOKY_FACTS_PATH` environment variable
- **Type**: BadgerDB (embedded key-value store)

### Backup and Version Control

The facts database file should not be committed to version control as it contains runtime data. It's already included in `.gitignore`. 