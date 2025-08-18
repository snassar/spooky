# Schema Updates for Logging Configuration

## Overview

This document describes the updates made to the spooky schemas to support the new logging configuration system. The changes ensure that logging configuration is properly validated and integrated into the project structure.

## Schema Changes

### 1. Project Schema (`project.schema.hcl`)

#### Added: Logging Configuration Block

The `project.schema.hcl` file has been updated to include a comprehensive `logging` configuration block that allows projects to override global logging settings.

```hcl
# Project-specific logging configuration (overrides global logging settings)
logging = {
  type = "object"
  required = false
  description = "Project-specific logging configuration that overrides global settings"
  
  properties = {
    level = {
      type = "string"
      required = false
      enum = ["debug", "info", "warn", "error", "fatal"]
      description = "Log level for this project (overrides global setting)"
    }
    
    format = {
      type = "string"
      required = false
      enum = ["json", "text", "structured"]
      description = "Log format for this project (overrides global setting)"
    }
    
    output = {
      type = "string"
      required = false
      enum = ["stdout", "stderr", "file", "null"]
      description = "Log output destination for this project (overrides global setting)"
    }
    
    file = {
      type = "object"
      required = false
      description = "File output configuration for this project"
      
      properties = {
        path = {
          type = "string"
          required = false
          description = "Path to log file (relative to project directory or absolute)"
        }
        
        permissions = {
          type = "string"
          required = false
          pattern = "^[0-7]{3,4}$"
          default = "0644"
          description = "File permissions in octal format (e.g., 0644)"
        }
        
        append = {
          type = "boolean"
          required = false
          default = true
          description = "Append to existing log file instead of overwriting"
        }
      }
    }
    
    filtering = {
      type = "object"
      required = false
      description = "Component-specific filtering for this project"
      
      properties = {
        components = {
          type = "object"
          required = false
          additional_properties = "string"
          description = "Component-specific log levels (e.g., 'ssh' = 'debug')"
        }
        
        patterns = {
          type = "array"
          required = false
          items = "string"
          description = "Pattern-based filtering rules"
        }
      }
    }
    
    rotation = {
      type = "object"
      required = false
      description = "Log rotation configuration for this project"
      
      properties = {
        enabled = {
          type = "boolean"
          required = false
          default = false
          description = "Enable log rotation for this project"
        }
        
        max_size = {
          type = "string"
          required = false
          description = "Maximum log file size (e.g., '100MB', '1GB')"
        }
        
        max_age = {
          type = "string"
          required = false
          description = "Maximum log file age (e.g., '24h', '7d', '30d')"
        }
        
        max_backups = {
          type = "integer"
          required = false
          min = 1
          max = 100
          description = "Maximum number of backup files to keep"
        }
        
        compress = {
          type = "boolean"
          required = false
          default = true
          description = "Compress rotated log files"
        }
      }
    }
  }
}
```

#### Added: Validation Rules

New validation rules have been added to ensure logging configuration is properly formatted:

```hcl
# Logging validation
logging_level_valid = {
  rule = "enum"
  values = ["debug", "info", "warn", "error", "fatal"]
  message = "Log level must be one of: debug, info, warn, error, fatal"
}

logging_format_valid = {
  rule = "enum"
  values = ["json", "text", "structured"]
  message = "Log format must be one of: json, text, structured"
}

logging_output_valid = {
  rule = "enum"
  values = ["stdout", "stderr", "file", "null"]
  message = "Log output must be one of: stdout, stderr, file, null"
}

logging_file_permissions_valid = {
  rule = "regex"
  pattern = "^[0-7]{3,4}$"
  message = "File permissions must be in octal format (e.g., 0644)"
}

logging_rotation_backups_reasonable = {
  rule = "range"
  min = 1
  max = 100
  message = "Maximum log backups must be between 1 and 100"
}

logging_file_path_required = {
  rule = "required_if"
  condition = "output == 'file'"
  field = "logging.file.path"
  message = "File path is required when output is set to 'file'"
}
```

### 2. Project Directory Schema (`project-directory.schema.hcl`)

#### Updated: Validation Rules

The project directory schema has been updated to include validation rules for logging configuration:

```hcl
# Cross-file validation rules
validation_rules = [
  "machines_file_or_directory_exists",
  "actions_file_or_directory_exists", 
  "variables_file_or_directory_exists",
  "facts_database_initialized",
  "no_circular_references",
  "logging_file_output_requires_logs_directory",
  "logging_file_path_validation"
]
```

#### Existing: Logs Directory Support

The schema already included support for a `logs` directory:

```hcl
directory "logs" {
  type = "directory"
  required = false
  description = "Log files directory"
  validate = "directory_exists"
}
```

## Configuration Examples

### Basic Project with Logging Override

```hcl
project "my-project" {
  name = "my-project"
  description = "Example project with logging configuration"
  
  # Override global logging for this project
  logging {
    level = "debug"
    output = "file"
    file {
      path = "./logs/project.log"
      permissions = "0644"
      append = false
    }
  }
}
```

### Advanced Project with Component Filtering

```hcl
project "advanced-project" {
  name = "advanced-project"
  description = "Advanced project with detailed logging configuration"
  
  logging {
    level = "info"
    format = "json"
    output = "file"
    
    file {
      path = "./logs/advanced.log"
      permissions = "0640"
      append = true
    }
    
    filtering {
      components = {
        "ssh"     = "debug"
        "facts"   = "info"
        "actions" = "warn"
        "project" = "debug"
      }
    }
    
    rotation {
      enabled = true
      max_size = "50MB"
      max_age = "7d"
      max_backups = 5
      compress = true
    }
  }
}
```

## Validation Behavior

### 1. **Required Field Validation**
- When `output = "file"`, the `logging.file.path` field is required
- File permissions must be in valid octal format (e.g., `0644`, `0640`)

### 2. **Enum Validation**
- Log levels must be one of: `debug`, `info`, `warn`, `error`, `fatal`
- Log formats must be one of: `json`, `text`, `structured`
- Output destinations must be one of: `stdout`, `stderr`, `file`, `null`

### 3. **Range Validation**
- Maximum log backups must be between 1 and 100
- File permissions must be valid octal numbers

### 4. **Directory Validation**
- When using file output, the logs directory should exist (if specified)
- Project directory structure validation includes logs directory support

## Integration with Global Configuration

### Configuration Hierarchy

1. **Default Configuration** (built-in)
2. **Global Configuration** (`~/.config/spooky/logging.hcl`)
3. **Project Configuration** (`project.hcl` - if exists)

### Merging Behavior

- Project configuration **overrides** global settings
- Only specified fields are overridden
- Unspecified fields inherit from global configuration
- File paths can be relative to project directory or absolute

### Example Merging

```hcl
# Global configuration
logging {
  level = "info"
  format = "json"
  output = "stderr"
}

# Project configuration
logging {
  level = "debug"
  output = "file"
  file {
    path = "./logs/project.log"
  }
}

# Result: Merged configuration
# level = "debug" (from project)
# format = "json" (from global)
# output = "file" (from project)
# file.path = "./logs/project.log" (from project)
```

## Testing and Validation

### Test Utilities

Two test utilities have been created to validate the logging configuration:

1. **`examples/setup-logging-config-utility.go`**
   - Creates default global logging configuration
   - Sets up `~/.config/spooky/logging.hcl`

2. **`examples/test-logging-schema-validation-utility.go`**
   - Tests global configuration loading
   - Tests project configuration loading
   - Tests configuration merging
   - Tests logging setup and functionality

### Running Tests

```bash
# Create global logging configuration
go run examples/setup-logging-config-utility.go

# Test logging configuration validation
go run examples/test-logging-schema-validation-utility.go
```

## Migration Guide

### From Previous Versions

1. **No Breaking Changes**: Existing project configurations continue to work
2. **Optional Enhancement**: Add `logging` block to `project.hcl` for project-specific logging
3. **Global Configuration**: Create `~/.config/spooky/logging.hcl` for system-wide settings

### Adding Logging to Existing Projects

1. **Create global configuration**:
   ```bash
   go run examples/setup-logging-config-utility.go
   ```

2. **Add logging block to project.hcl** (optional):
   ```hcl
   project "existing-project" {
     # ... existing configuration ...
     
     logging {
       level = "debug"
       output = "file"
       file {
         path = "./logs/project.log"
       }
     }
   }
   ```

3. **Create logs directory** (if using file output):
   ```bash
   mkdir -p logs
   ```

## Benefits

### 1. **Schema-Driven Validation**
- All logging configuration is validated against schemas
- Prevents configuration errors at validation time
- Ensures consistent configuration structure

### 2. **Flexible Configuration**
- Global defaults with project-specific overrides
- Component-specific filtering
- Multiple output formats and destinations

### 3. **Integration with Existing Structure**
- Leverages existing project directory structure
- Uses established validation patterns
- Maintains backward compatibility

### 4. **Comprehensive Testing**
- Test utilities for validation
- Example configurations for reference
- Clear migration path

## Conclusion

The schema updates provide comprehensive support for the new logging configuration system while maintaining backward compatibility and following established patterns. The validation ensures that logging configuration is properly structured and integrated into the project lifecycle.
