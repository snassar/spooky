# Dependencies setup for test-invalid-templates project

actions {
  action "install-dependencies" {
    description = "Install system dependencies"
    command     = "apt update && apt install -y curl wget git"
    tags        = ["role=web", "setup"]
    parallel    = true
    timeout     = 600
  }
}