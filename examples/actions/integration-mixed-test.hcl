# Mixed test configuration for spooky SSH automation tool
# Points to both working and failure scenario SSH machines for integration testing

# Working machines
machine "spooky-test-server-1" {
  host     = "localhost"
  port     = 2221
  user     = "testuser"
  key_file = "~/.ssh/id_ed25519"
  tags = {
    environment = "testing"
    role        = "test-server"
    group       = "primary"
  }
}

machine "spooky-test-server-2" {
  host     = "localhost"
  port     = 2222
  user     = "testuser"
  key_file = "~/.ssh/id_ed25519"
  tags = {
    environment = "testing"
    role        = "test-server"
    group       = "primary"
  }
}

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

# Test actions for mixed success/failure scenarios
action "test-mixed-success-failure" {
  description = "Test mixed success and failure scenarios"
  command     = "echo 'Testing mixed scenarios on ' $(hostname) && uptime"
  machines    = ["spooky-test-server-1", "spooky-test-server-2", "spooky-test-no-ssh", "spooky-test-ssh-no-key"]
  parallel    = true
} 