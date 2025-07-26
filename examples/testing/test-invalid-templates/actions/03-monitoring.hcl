# Monitoring actions for test-invalid-templates project

actions {
  action "check-disk-space" {
    description = "Check disk space usage"
    command     = "df -h"
    tags        = ["monitoring"]
    parallel    = true
    timeout     = 60
  }

  action "check-memory" {
    description = "Check memory usage"
    command     = "free -h"
    tags        = ["monitoring"]
    parallel    = true
    timeout     = 60
  }
}