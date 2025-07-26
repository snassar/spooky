# Main actions for test-unsupported-action-type project
# This file can contain core actions or reference actions in the actions/ directory

actions {
  action "check-status" {
    description = "Check server status"
    command     = "uptime && df -h"
    tags        = ["role=web"]
    parallel    = true
    timeout     = 300
  }
  
  # Unsupported action type - should cause validation error
  action "unsupported-action" {
    description = "Action with unsupported type"
    type        = "unknown"
    command     = "echo 'test'"
    tags        = ["role=web"]
    parallel    = true
    timeout     = 300
  }
}