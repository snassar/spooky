# Actions Configuration Schema Summary

This schema defines the structure and validation rules for Spooky action definitions in `actions.hcl` files. It provides a robust framework for defining executable tasks with extensive configuration options and security controls.

## Core Action Properties

### Required Fields
- **`description`** - Human-readable action description (1-500 chars)
- **`type`** - Action execution type: `command`, `script`, `template_deploy`, `file_copy`, or `service_control`

### Execution Configuration
- **`command`** - Shell command to execute (for command type)
- **`script`** - Script file path (for script type)
- **`template`** - Template deployment configuration with source, destination, validation, backup, permissions, owner, and group settings

## Targeting and Execution Control

### Machine Targeting
- **`machines`** - Array of specific machine names
- **`tags`** - Array of tags for machine selection

### Execution Behavior
- **`timeout`** - Action timeout (1-3600 seconds, default 300)
- **`parallel`** - Execute across machines in parallel
- **`retries`** - Number of retry attempts (0-10, default 0)
- **`retry_delay`** - Delay between retries (1-300 seconds, default 5)

## Advanced Features

### Dependencies & Organization
- **`dependencies`** - Actions that must complete first
- **`category`** - Action categorization
- **`priority`** - Execution priority (1-10, default 5)
- **`critical`** - Whether failure stops execution

### Security & Environment
- **`environment`** - Environment variables
- **`working_directory`** - Execution directory
- **`user`** - Run as specific user
- **`sudo`** - Use sudo privileges
- **`dry_run`** - Preview without execution

### Performance & Resources
- **`max_concurrent`** - Maximum parallel executions (1-100)
- **`resource_limits`** - Memory, CPU, and disk limits

## Validation & Security

### Security Controls
- Command injection prevention (blocks shell operators `;&|`$`)
- File path validation patterns
- User/group name restrictions

### Validation Rules
- Action name format validation
- Type-specific requirement checks
- Circular dependency detection
- Machine targeting requirements
- Resource limit validation

### Quality Controls
- **`validate_before_run`** - Pre-execution validation (default true)
- **`allow_failure`** - Continue on failure
- **`metadata`** - Additional custom properties

## Summary

This schema provides a comprehensive, secure, and flexible system for defining infrastructure automation actions with extensive configuration options while maintaining strict validation and security controls.
