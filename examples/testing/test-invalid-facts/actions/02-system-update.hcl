# System updates for test-invalid-facts project

actions {
  action "update-system" {
    description = "Update system packages"
    command     = "apt update && apt upgrade -y"
    tags        = ["environment=development", "maintenance"]
    parallel    = true
    timeout     = 600
  }
}