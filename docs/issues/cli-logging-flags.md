# Add CLI Logging Flags for Debugging and Verbose Output

## Problem Statement

Currently, spooky only supports logging configuration through HCL files (`$XDG_CONFIG_HOME/spooky/spooky.hcl` and `$XDG_CONFIG_HOME/spooky/logging.hcl`). This makes it difficult to quickly increase log verbosity for debugging purposes without editing configuration files.

**Current limitations:**
- No way to increase log level from CLI for debugging
- No verbose output option for troubleshooting
- No debug mode for development
- Users must edit config files to change logging behavior

## Use Cases

1. **Debugging Issues**: Users need to quickly enable debug logging to troubleshoot problems
2. **Development**: Developers need verbose output during development and testing
3. **Troubleshooting**: Users need to see detailed output without editing config files
4. **CI/CD**: Automated systems need to control logging verbosity via CLI flags

## Proposed Solution

Add standard Go CLI logging flags following idiomatic patterns:

### Global Flags (on root command)

```bash
# Verbose output (info level)
spooky --verbose <command>

# Debug output (debug level)  
spooky --debug <command>

# Quiet mode (error level only)
spooky --quiet <command>

# Explicit log level
spooky --log-level=debug <command>
spooky --log-level=info <command>
spooky --log-level=warn <command>
spooky --log-level=error <command>

# Log format override
spooky --log-format=text <command>
spooky --log-format=json <command>
```

### Flag Precedence

CLI flags should override HCL configuration with this precedence order:
1. CLI flags (highest priority)
2. Environment variables
3. HCL configuration files (lowest priority)

## Implementation Details

### 1. Add Global Flags to Root Command

```go
// In cmd/root.go
var (
    verbose   bool
    debug     bool
    quiet     bool
    logLevel  string
    logFormat string
)

func init() {
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output (info level)")
    rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug output (debug level)")
    rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "quiet mode (error level only)")
    rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "", "log level (debug, info, warn, error)")
    rootCmd.PersistentFlags().StringVar(&logFormat, "log-format", "", "log format (text, json)")
}
```

### 2. Update Logging Integration

```go
// In internal/config/integration.go
func (i *Integration) LoadLoggingConfig() (*spookytypeslogging.Config, error) {
    config, err := i.loadLoggingConfigFromFiles()
    if err != nil {
        return nil, err
    }
    
    // Override with CLI flags
    if i.cliFlags.verbose {
        config.Level = "info"
    }
    if i.cliFlags.debug {
        config.Level = "debug"
    }
    if i.cliFlags.quiet {
        config.Level = "error"
    }
    if i.cliFlags.logLevel != "" {
        config.Level = i.cliFlags.logLevel
    }
    if i.cliFlags.logFormat != "" {
        config.Format = i.cliFlags.logFormat
    }
    
    return config, nil
}
```

### 3. Update CLI Help

Update the main help output to include logging flags:

```bash
$ spooky --help
Usage:
  spooky [flags]
  spooky [command]

Available Commands:
  actions     Manage and run actions
  facts       Manage fact collection
  machines    Manage machine inventory
  project     Manage projects
  schemas     Manage schemas
  secrets     Manage secrets
  templates   Manage templates
  variables   Manage variables

Global Flags:
  -h, --help              help for spooky
  -v, --version           version for spooky
  --verbose, -v           verbose output (info level)
  --debug                 debug output (debug level)
  --quiet, -q             quiet mode (error level only)
  --log-level string      log level (debug, info, warn, error)
  --log-format string     log format (text, json)
```

## Examples

### Debug a Command
```bash
# Debug machine ping
spooky --debug machines ping myproject

# Verbose output for action run
spooky --verbose actions run myproject

# Custom log level
spooky --log-level=debug --log-format=text actions run myproject
```

### CI/CD Usage
```bash
# Quiet mode for automated scripts
spooky --quiet actions run myproject

# Debug mode for troubleshooting
spooky --debug machines ping myproject
```

## Testing Requirements

1. **Flag Precedence**: Verify CLI flags override HCL config
2. **Flag Combinations**: Test various flag combinations
3. **Help Output**: Verify help shows all logging flags
4. **Integration**: Test with existing logging system
5. **Backward Compatibility**: Ensure existing HCL config still works

## Acceptance Criteria

- [ ] `--verbose` flag sets log level to info
- [ ] `--debug` flag sets log level to debug  
- [ ] `--quiet` flag sets log level to error
- [ ] `--log-level` accepts debug, info, warn, error
- [ ] `--log-format` accepts text, json
- [ ] CLI flags override HCL configuration
- [ ] Help output includes all logging flags
- [ ] Backward compatibility maintained
- [ ] Tests cover all flag combinations

## Related Issues

- Current logging system uses HCL configuration only
- No CLI-based debugging capabilities
- Users must edit config files for debugging

## Labels

- `enhancement`
- `cli`
- `logging`
- `debugging`
- `user-experience`

## Priority

**Medium** - This improves developer experience and debugging capabilities but doesn't block core functionality.
