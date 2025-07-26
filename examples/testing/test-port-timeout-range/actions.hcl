# Main actions for test-valid-project project
# This file can contain core actions or reference actions in the actions/ directory

actions {
  action "check-status" {
    description = "Check server status"
    command     = "uptime && df -h"
    tags        = ["role=web"]
    parallel    = true
    timeout     = 300
  }
}

# Actions with invalid timeout values

actions {
  action "timeout-zero" {
    description = "Action with timeout 0"
    command     = "uptime"
    tags        = ["role=web"]
    parallel    = true
    timeout     = 0
  }
  action "timeout-too-high" {
    description = "Action with timeout 5000"
    command     = "uptime"
    tags        = ["role=web"]
    parallel    = true
    timeout     = 5000
  }
}