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

# Actions with both and with neither command/script

actions {
  action "both-command-script" {
    description = "Action with both command and script"
    command     = "uptime"
    script      = "some-script.sh"
    tags        = ["role=web"]
    parallel    = true
    timeout     = 300
  }
  action "neither-command-script" {
    description = "Action with neither command nor script"
    tags        = ["role=web"]
    parallel    = true
    timeout     = 300
  }
}