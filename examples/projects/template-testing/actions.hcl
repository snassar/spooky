# Actions for template-testing project
# Add your action definitions here

actions {
  action "check-status" {
    description = "Check server status"
    command     = "uptime && df -h"
    tags        = ["role=web"]
    parallel    = true
    timeout     = 300
  }

  action "update-system" {
    description = "Update system packages"
    command     = "apt update && apt upgrade -y"
    tags        = ["environment=development"]
    parallel    = true
    timeout     = 600
  }
}