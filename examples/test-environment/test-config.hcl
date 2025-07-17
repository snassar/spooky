# Test configuration for spooky SSH automation tool
# Points to the test environment SSH server

# Define test server
server "spooky-test-server" {
  host     = "localhost"
  port     = 2221
  user     = "root"
  key_file = "ssh-keys/id_ed25519"
  tags = {
    environment = "testing"
    role        = "test-server"
  }
}

# Define test actions
action "check-status" {
  description = "Check system status"
  command     = "uptime && df -h"
  servers     = ["spooky-test-server"]
  parallel    = true
}

action "check-ssh-keys" {
  description = "Check SSH key configuration"
  command     = "ls -la /root/.ssh/ && echo 'SSH keys mounted successfully'"
  servers     = ["spooky-test-server"]
  parallel    = true
}

action "test-connection" {
  description = "Test basic connectivity"
  command     = "echo 'Connection test successful' && hostname && whoami"
  servers     = ["spooky-test-server"]
  parallel    = true
} 