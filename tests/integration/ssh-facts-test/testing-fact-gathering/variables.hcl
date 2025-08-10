# Test variables for SSH fact gathering

variable "test_vm_host" {
  description = "Test VM hostname"
  type = "string"
  default = "localhost"
}

variable "test_vm_port" {
  description = "Test VM SSH port"
  type = "integer"
  default = 2222
}

variable "test_vm_user" {
  description = "Test VM SSH user"
  type = "string"
  default = "spooky"
}

variable "test_vm_key_path" {
  description = "Path to SSH key for test VM"
  type = "string"
  default = "../keys/spooky-facts-test_key"
}

variable "test_timeout" {
  description = "SSH connection timeout"
  type = "string"
  default = "30s"
}

variable "test_retries" {
  description = "Number of SSH connection retries"
  type = "integer"
  default = 3
}

variable "facts_output_format" {
  description = "Output format for facts"
  type = "string"
  default = "json"
  validation {
    condition = "contains(['json', 'hcl'], var.facts_output_format)"
    error_message = "Output format must be 'json' or 'hcl'"
  }
} 