# Main actions for test-invalid-command-syntax project
# This file can contain core actions or reference actions in the actions/ directory

actions {
  action "check-status" {
    description = "Check server status"
    command     = "uptime && df -h"
    tags        = ["role=web"]
    parallel    = true
    timeout     = 300
  }
  
  # Empty command - should cause validation error
  action "empty-command" {
    description = "Action with empty command"
    command     = ""
    tags        = ["role=web"]
    parallel    = true
    timeout     = 300
  }
  
  # Invalid shell syntax - should cause validation error
  action "invalid-syntax" {
    description = "Action with invalid shell syntax"
    command     = "echo 'unclosed quote && ls"
    tags        = ["role=web"]
    parallel    = true
    timeout     = 300
  }
}