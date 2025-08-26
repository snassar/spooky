# test-syntax

Testing HCL syntax

## Overview

This is a Spooky automation project that defines configuration management, 
deployment automation, and infrastructure management tasks.

## Project Structure

- project.hcl - Project configuration and metadata
- machines.hcl - Machine inventory and connectivity settings
- actions.hcl - Automation tasks and deployment actions
- variables.hcl - Project-wide variables and configuration values
- templates/ - Template files for deployment (create as needed)
- files/ - Static files for deployment (create as needed)

## Getting Started

1. **Configure Machines**: Edit machines.hcl to define your target machines
2. **Define Actions**: Edit actions.hcl to create automation tasks
3. **Set Variables**: Edit variables.hcl to configure project variables
4. **Validate**: Run 'spooky project validate' to check configuration
5. **Execute**: Run 'spooky run <action-name>' to execute actions

## Examples

### Running Actions
```bash
# Run a specific action
spooky run deploy-application

# Run actions with specific tags
spooky run --tags deployment

# Dry run to see what would happen
spooky run --dry-run deploy-application
```

### Managing Machines
```bash
# List all machines
spooky machines list

# Test connectivity
spooky machines test-connection

# Collect facts
spooky machines collect-facts
```

## Documentation

For more information about Spooky, visit the project documentation or run:
```bash
spooky --help
spooky project --help
spooky machines --help
spooky actions --help
```

## Support

If you encounter issues or have questions, please refer to the Spooky documentation
or create an issue in the project repository.