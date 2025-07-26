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

# Actions referencing an unexecutable script file

actions {
  action "run-script" {
    description = "Run a script that is not executable"
    script      = "unexecutable.sh"
    tags        = ["role=web"]
    parallel    = true
    timeout     = 300
  }
}