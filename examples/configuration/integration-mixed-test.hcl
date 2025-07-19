# Mixed test configuration for spooky SSH automation tool
# Points to both working and failure scenario SSH servers for integration testing

# Working servers
server "spooky-test-server-1" {
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

server "spooky-test-server-2" {
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

# Failure scenario servers
server "spooky-test-no-ssh" {
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

server "spooky-test-ssh-no-key" {
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
  servers     = ["spooky-test-server-1", "spooky-test-server-2", "spooky-test-no-ssh", "spooky-test-ssh-no-key"]
  parallel    = true
} 