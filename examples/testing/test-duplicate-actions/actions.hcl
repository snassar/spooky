# Main actions for test-duplicate-actions project
# This file can contain core actions or reference actions in the actions/ directory

actions {
  action "check-status" {
    description = "Check server status"
    command     = "uptime && df -h"
    tags        = ["role=web"]
    parallel    = true
    timeout     = 300
  }
  
  # Duplicate action name - should cause validation error
  action "check-status" {
    description = "Another check status action"
    command     = "systemctl status"
    tags        = ["role=web"]
    parallel    = true
    timeout     = 300
  }
}