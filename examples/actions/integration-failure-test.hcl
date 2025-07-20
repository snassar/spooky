# Failure test configuration for spooky SSH automation tool
# Points to failure scenario SSH machines for integration testing

# Failure scenario machines
machine "spooky-test-no-ssh" {
  host     = "localhost"
  port     = 2228
  user     = "testuser"
  key_file = "~/.ssh/id_ed25519"
  tags = {
    environment = "testing"
    role        = "failure-test"
    group       = "failure"
  }
}

machine "spooky-test-ssh-no-key" {
  host     = "localhost"
  port     = 2229
  user     = "testuser"
  key_file = "~/.ssh/id_ed25519"
  tags = {
    environment = "testing"
    role        = "failure-test"
    group       = "failure"
  }
}

# Test actions for connection failure scenarios
action "test-connection-failures" {
  description = "Test connection failures (no SSH, no key)"
  command     = "echo 'This should fail on failure-test machines' && hostname"
  machines    = ["spooky-test-no-ssh", "spooky-test-ssh-no-key"]
  parallel    = true
} 