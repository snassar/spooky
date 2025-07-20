# Example configuration for spooky SSH automation tool

# Define machines
machine "debian-tester" {
  host     = "192.168.178.73"
  port     = 22
  user     = "builder"
  password = "builder"
  # key_file = "~/.ssh/id_rsa"  # Alternative to password
  tags = {
    environment = "testing"
    role        = "undefined"
  }
}

# Define actions
action "check-status" {
  description = "Check system status"
  command     = "uptime && df -h"
  machines    = ["debian-tester"]
  parallel    = true
}