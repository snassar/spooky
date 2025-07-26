# Main actions for test-invalid-actions project
# This file can contain core actions or reference actions in the actions/ directory

actions {
  action "check-status" {
    description = "Check server status"
    command     = "uptime && df -h"
    tags        = ["role=web"]
    parallel    = true
    timeout     = 300
  }
  
  # Invalid action - missing required command field
  action "invalid-action" {
    description = "Invalid action without command"
    tags        = ["role=web"]
    parallel    = true
    timeout     = 300
  }
  
  # Invalid syntax - missing closing brace
  action "broken-action" {
    description = "Action with missing closing brace"
    command     = "echo 'test'"
    tags        = ["role=web"]
    parallel    = true
    timeout     = 300
    # Missing closing brace
  }
}