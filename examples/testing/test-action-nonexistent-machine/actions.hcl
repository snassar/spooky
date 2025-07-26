# Main actions for test-valid-project project
# This file can contain core actions or reference actions in the actions/ directory

actions {
  action "check-status" {
    description = "Check server status"
    command     = "uptime && df -h"
    machines    = ["example-server", "ghost-machine"]
    tags        = ["role=web"]
    parallel    = true
    timeout     = 300
  }
}