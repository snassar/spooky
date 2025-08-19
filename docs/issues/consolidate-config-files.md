# 4
# Consolidate configuration files and expand spooky.hcl schema

## Problem Statement
Currently, the system generates multiple small configuration files (`spooky.hcl` and `logging.hcl`) which creates complexity and confusion. The `spooky.hcl` schema already includes a basic logging section, but we're generating a separate `logging.hcl` file with more detailed logging configuration.

### Current Issues
1. **Multiple configuration files** - Users have to manage both `spooky.hcl` and `logging.hcl`
2. **Configuration overlap** - Both files define logging settings
3. **User confusion** - Which file should be edited for logging configuration?
4. **Maintenance complexity** - More files to keep in sync and validate
5. **Schema inconsistency** - `spooky.hcl` schema has basic logging, but we generate separate detailed logging file

### Current Implementation
```go
// In auto_setup.go
func ensureConfigFiles(configDir string) error {
    // Creates spooky.hcl
    if err := createDefaultSpookyConfig(configDir); err != nil {
        return fmt.Errorf("failed to create default spooky.hcl: %w", err)
    }

    // Creates separate logging.hcl
    if err := createDefaultLoggingConfig(configDir); err != nil {
        return fmt.Errorf("failed to create default logging.hcl: %w", err)
    }
    return nil
}
```

## Proposed Solution

### 1. Consolidate to Single Configuration File
- **Use only `spooky.hcl`** for all configuration
- **Remove `logging.hcl` generation** from auto-setup
- **Simplify user experience** - one file to manage

### 2. Expand spooky.hcl Schema
The current `spooky.hcl` schema has a basic logging section:
```hcl
logging {
  level = {
    type = "string"
    required = false
    enum = ["debug", "info", "warn", "error"]
    default = "info"
    description = "Default logging level"
  }
  
  format = {
    type = "string"
    required = false
    enum = ["json", "text"]
    default = "text"
    description = "Log format"
  }
  
  output = {
    type = "string"
    required = false
    default = "stderr"
    description = "Log output destination"
  }
}
```

**Expand this to include:**
- File output configuration
- Structured logging options
- Performance settings
- Component filtering
- Log rotation settings
- All the detailed logging options currently in `logging.hcl`

### 3. Update Configuration Generation
- **Remove `createDefaultLoggingConfig`** function
- **Expand `createDefaultSpookyConfig`** to include comprehensive logging configuration
- **Use schema-driven generation** for the expanded spooky.hcl content

## Benefits
- **Simplified user experience** - One configuration file
- **Reduced complexity** - No file overlap or confusion
- **Better schema utilization** - Use the existing spooky.hcl schema properly
- **Easier maintenance** - Single file to validate and manage
- **Clearer architecture** - All configuration in one place

## Priority
**High Priority** - Configuration simplification and user experience improvement

## Related Files
- `internal/config/auto_setup.go` - Current implementation
- `internal/schemas/schemas/structure/spooky.hcl` - Schema to expand
- `internal/schemas/schemas/structure/logging.hcl` - Schema to consolidate
- `internal/config/logging.go` - May need updates for single file approach

## Acceptance Criteria
- [x] Remove `createDefaultLoggingConfig` function
- [x] Expand `spooky.hcl` schema to include comprehensive logging configuration
- [x] Update `createDefaultSpookyConfig` to generate complete configuration
- [x] Remove `logging.hcl` file generation from auto-setup
- [x] Update configuration loading to use single file approach
- [x] Ensure backward compatibility with existing `logging.hcl` files
- [x] Update documentation to reflect single configuration file approach
- [x] Add tests for consolidated configuration generation

## Implementation Notes
- Follow schema-driven generation pattern for expanded spooky.hcl
- Ensure the expanded logging section in spooky.hcl covers all current logging.hcl functionality
- Consider migration path for users with existing separate logging.hcl files
- Update configuration loading logic to handle both old and new approaches during transition
