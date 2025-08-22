# Facts Schema
# Schema for fact collection and storage
# Defines the structure and validation rules for facts

# Schema metadata
metadata {
  version = "1"
  description = "Facts configuration schema for spooky fact collection"
}

# Facts structure
facts {
  # Fact name
  name = {
    type = "string"
    required = true
    pattern = "^[a-zA-Z0-9_.-]+$"
    min_length = 1
    max_length = 128
    description = "Name of the fact"
  }
  
  # Fact value
  value = {
    type = "any"
    required = true
    description = "Value of the fact - can be string, number, boolean, object, or age-encrypted string"
  }
  
  # Fact type
  type = {
    type = "string"
    required = true
    enum = ["string", "number", "boolean", "object", "array", "encrypted"]
    description = "Type of the fact value"
  }
  
  # Encryption flag
  encrypted = {
    type = "boolean"
    required = false
    default = false
    description = "Whether the fact value is age-encrypted"
  }
  
  # Age encryption metadata (optional, only when encrypted = true)
  encryption_metadata = {
    type = "object"
    required = false
    description = "Age encryption metadata - only present when encrypted = true"
  }
  
  # Fact description
  description = {
    type = "string"
    required = false
    max_length = 256
    description = "Description of the fact"
  }
  
  # Fact tags
  tags = {
    type = "array"
    required = false
    description = "Tags for categorizing facts"
    max_items = 10
    items = {
      type = "string"
      pattern = "^[a-zA-Z0-9_-]+$"
      min_length = 1
      max_length = 32
    }
  }
  
  # Fact metadata
  metadata = {
    type = "object"
    required = false
    description = "Additional metadata for the fact"
  }
}


