# todo

## machines

With this comprehensive documentation foundation, we're now ready to implement all the remaining features you mentioned:
SSH Connection Implementation - Complete the ping functionality
Integration with Other Systems - Facts, actions, variables, SSH
Advanced Filtering - Tag and group-based filtering
Export - Export to HCL according to machines schema
CLI Command Completion - Add missing commands
Connection Pooling - Performance optimization
Complex Query Language - Advanced filtering syntax
Import from External Sources - OpenStack, Hetzner, AWS, GCP, Azure, APIs
The documentation provides clear guidance for implementing each of these features while maintaining consistency with the established patterns and architecture.

## variables

### Minor Issues (Non-Critical)
The variables system is production-ready and fully functional. The following are cosmetic and non-functional improvements:

**Code Quality (Optional):**
- Fix linting warnings for cleaner code
  - Unused parameters - Some function parameters are marked as unused (interface compliance requirements)
  - High cyclomatic complexity - Some functions are complex but working correctly
  - Naming conventions - Some type names could be improved but don't affect functionality
- Optimize performance in complex functions
- Improve error handling patterns

### Future Enhancements (Planned)
**Template Integration:**
- Variables in templates
- Template variable substitution
- Dynamic template rendering

**Export Functionality:**
- JSON/HCL export
- Variable collection export
- Configuration export

**Advanced Features:**
- Advanced filtering and querying
- Caching for performance
- Integration with other systems (actions, facts, SSH)
- Variable dependency visualization
- Environment-specific variable management

## logging

### Future Enhancements (Planned)
The logging system is fully functional and production-ready. The following enhancements would improve configuration flexibility:

**HCL Configuration Parsing:**
- Implement `parseLoggingConfig()` in `internal/config/logging.go` - Parse HCL logging configuration files
- Implement `parseProjectLoggingConfig()` - Extract logging settings from project.hcl files
- Support for global logging configuration at `$XDG_CONFIG_HOME/spooky/logging.hcl`
- Support for project-specific logging configuration in `project.hcl`

**Enhanced Configuration Loading:**
- Load custom logging levels from configuration files
- Support for custom output formats and destinations
- Project-specific logging overrides
- Environment variable integration for logging configuration

**Advanced Logging Features:**
- Log rotation configuration from HCL files
- Performance optimization settings (buffering, async logging)
- Component-specific filtering from configuration
- Structured logging field customization