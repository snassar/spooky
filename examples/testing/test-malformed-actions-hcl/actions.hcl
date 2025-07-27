# Main actions for test-malformed-actions-hcl project
# This file has malformed HCL syntax

actions {
  action "check-status" {
    description = "Check server status"
    command     = "uptime && df -h"
    tags        = ["role=web"]
    parallel    = true
    timeout     = 300
  }
  
  # Malformed action - missing closing brace
  action "malformed-action" {
    description = "Action with malformed syntax"
    command     = "echo 'test'"
    tags        = ["role=web"]
    parallel    = true
    timeout     = 300
    # Missing closing brace