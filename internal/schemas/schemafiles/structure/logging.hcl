# Spooky Logging Configuration Schema
# Comprehensive schema for logging configuration with best practices for readability
# Based on log formatting best practices from Sematext

# Schema metadata
metadata {
  version = "1"
  description = "Logging configuration schema for spooky CLI"
}

# Logging configuration block
logging {
  # Log level configuration
  level = {
    type = "string"
    required = false
    default = "info"
    enum = ["debug", "info", "warn", "error", "fatal"]
    description = "Minimum log level to output (debug, info, warn, error, fatal)"
  }
  
  # Output format configuration
  format = {
    type = "string"
    required = false
    default = "json"
    enum = ["json", "text", "structured"]
    description = "Log output format (json for machine-readable, text for human-readable, structured for custom)"
  }
  
  # Output destination configuration
  output = {
    type = "string"
    required = false
    default = "stderr"
    enum = ["stdout", "stderr", "file", "null"]
    description = "Log output destination (stdout, stderr, file, null for no output)"
  }
  
  # File output configuration
  file_path = {
    type = "string"
    required = false
    description = "Path to log file (required when output is 'file')"
    pattern = "^[^<>:\"/\\|?*]+$"
  }
  
  file_permissions = {
    type = "string"
    required = false
    default = "0644"
    pattern = "^[0-7]{3,4}$"
    description = "File permissions in octal format (e.g., 0644)"
  }
  
  file_append = {
    type = "boolean"
    required = false
    default = true
    description = "Whether to append to existing file or truncate"
  }
  
  # Structured logging configuration
  structured_timestamp_enabled = {
    type = "boolean"
    required = false
    default = true
    description = "Whether to include timestamps in log entries"
  }
  
  structured_timestamp_format = {
    type = "string"
    required = false
    default = "RFC3339"
    enum = ["RFC3339", "RFC3339Nano", "Unix", "UnixNano", "ISO8601"]
    description = "Timestamp format (RFC3339, RFC3339Nano, Unix, UnixNano, ISO8601)"
  }
  
  structured_timestamp_timezone = {
    type = "string"
    required = false
    default = "UTC"
    description = "Timezone for timestamps (e.g., UTC, America/New_York)"
  }
  
  structured_level_key = {
    type = "string"
    required = false
    default = "level"
    description = "Field key for log level"
  }
  
  structured_message_key = {
    type = "string"
    required = false
    default = "message"
    description = "Field key for log message"
  }
  
  structured_error_key = {
    type = "string"
    required = false
    default = "error"
    description = "Field key for error information"
  }
  
  structured_fields_include = {
    type = "array"
    required = false
    items = {
      type = "string"
    }
    description = "Additional fields to include in structured logs"
  }
  
  structured_fields_exclude = {
    type = "array"
    required = false
    items = {
      type = "string"
    }
    description = "Fields to exclude from structured logs"
  }
  
  structured_fields_filter_sensitive = {
    type = "array"
    required = false
    items = {
      type = "string"
    }
    description = "Sensitive field names to mask or redact"
  }
  
  # Performance configuration
  performance_buffer_enabled = {
    type = "boolean"
    required = false
    default = false
    description = "Whether to use buffered logging"
  }
  
  performance_buffer_size = {
    type = "integer"
    required = false
    default = 4096
    min = 1024
    max = 1048576
    description = "Buffer size in bytes"
  }
  
  performance_buffer_flush_interval = {
    type = "string"
    required = false
    default = "1s"
    description = "Flush interval (e.g., 1s, 100ms)"
  }
  
  performance_async_enabled = {
    type = "boolean"
    required = false
    default = false
    description = "Whether to use asynchronous logging"
  }
  
  performance_async_queue_size = {
    type = "integer"
    required = false
    default = 1000
    min = 100
    max = 100000
    description = "Queue size for async logging"
  }
  
  performance_async_workers = {
    type = "integer"
    required = false
    default = 1
    min = 1
    max = 10
    description = "Number of worker goroutines for async logging"
  }
  
  performance_async_drop_when_full = {
    type = "boolean"
    required = false
    default = false
    description = "Whether to drop logs when queue is full"
  }
  
  # Filtering configuration
  filtering_components = {
    type = "object"
    required = false
    additional_properties = {
      type = "string"
      enum = ["debug", "info", "warn", "error", "fatal"]
    }
    description = "Component-specific log level configuration"
  }
  
  filtering_patterns_include = {
    type = "array"
    required = false
    items = {
      type = "string"
    }
    description = "Regex patterns to include (logs must match at least one)"
  }
  
  filtering_patterns_exclude = {
    type = "array"
    required = false
    items = {
      type = "string"
    }
    description = "Regex patterns to exclude (logs matching any are dropped)"
  }
  
  # Rotation configuration
  rotation_enabled = {
    type = "boolean"
    required = false
    default = false
    description = "Whether to enable log file rotation"
  }
  
  rotation_max_size = {
    type = "string"
    required = false
    default = "100MB"
    description = "Maximum file size before rotation (e.g., 100MB, 1GB)"
  }
  
  rotation_max_age = {
    type = "string"
    required = false
    default = "30d"
    description = "Maximum age of rotated files (e.g., 30d, 7d, 24h)"
  }
  
  rotation_max_backups = {
    type = "integer"
    required = false
    default = 5
    min = 1
    max = 100
    description = "Maximum number of backup files to keep"
  }
  
  rotation_compress = {
    type = "boolean"
    required = false
    default = true
    description = "Whether to compress rotated log files"
  }
  
  rotation_local_time = {
    type = "boolean"
    required = false
    default = false
    description = "Whether to use local time for rotation timestamps"
  }
}
