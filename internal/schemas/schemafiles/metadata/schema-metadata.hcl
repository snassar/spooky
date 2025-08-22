# Schema Metadata Schema
# Simple metadata structure for schema versioning and compatibility
# Inspired by rsync's protocol versioning approach

# Schema metadata
metadata {
  version = "1"
  description = "Simple schema metadata structure for versioning and compatibility"
}

# Schema metadata structure
schema_metadata {
  # Schema version (required)
  version = {
    type = "integer"
    required = true
    min = 1
    description = "Schema version number (incremental)"
  }
  
  # Schema description (required)
  description = {
    type = "string"
    required = true
    min_length = 1
    max_length = 200
    description = "Brief description of the schema purpose"
  }
}
