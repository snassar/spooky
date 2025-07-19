# Working servers only configuration for spooky SSH automation tool
# Points to 7 working SSH servers for integration testing

# Define working test servers (7 containers)
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

server "spooky-test-server-3" {
  host     = "localhost"
  port     = 2223
  user     = "testuser"
  key_file = "~/.ssh/id_ed25519"
  tags = {
    environment = "testing"
    role        = "test-server"
    group       = "secondary"
  }
}

server "spooky-test-server-4" {
  host     = "localhost"
  port     = 2224
  user     = "testuser"
  key_file = "~/.ssh/id_ed25519"
  tags = {
    environment = "testing"
    role        = "test-server"
    group       = "secondary"
  }
}

server "spooky-test-server-5" {
  host     = "localhost"
  port     = 2225
  user     = "testuser"
  key_file = "~/.ssh/id_ed25519"
  tags = {
    environment = "testing"
    role        = "test-server"
    group       = "tertiary"
  }
}

server "spooky-test-server-6" {
  host     = "localhost"
  port     = 2226
  user     = "testuser"
  key_file = "~/.ssh/id_ed25519"
  tags = {
    environment = "testing"
    role        = "test-server"
    group       = "tertiary"
  }
}

server "spooky-test-server-7" {
  host     = "localhost"
  port     = 2227
  user     = "testuser"
  key_file = "~/.ssh/id_ed25519"
  tags = {
    environment = "testing"
    role        = "test-server"
    group       = "tertiary"
  }
}

# Test actions for working servers only
action "check-status-all" {
  description = "Check system status on all working servers"
  command     = "uptime && df -h && echo 'Server: ' $(hostname)"
  servers     = ["spooky-test-server-1", "spooky-test-server-2", "spooky-test-server-3", "spooky-test-server-4", "spooky-test-server-5", "spooky-test-server-6", "spooky-test-server-7"]
  parallel    = true
}

action "check-ssh-keys-all" {
  description = "Check SSH key configuration on all working servers"
  command     = "ls -la ~/.ssh/ && echo 'SSH keys mounted successfully on ' $(hostname)"
  servers     = ["spooky-test-server-1", "spooky-test-server-2", "spooky-test-server-3", "spooky-test-server-4", "spooky-test-server-5", "spooky-test-server-6", "spooky-test-server-7"]
  parallel    = true
}

action "test-connection-all" {
  description = "Test basic connectivity on all working servers"
  command     = "echo 'Connection test successful on ' $(hostname) && whoami && pwd"
  servers     = ["spooky-test-server-1", "spooky-test-server-2", "spooky-test-server-3", "spooky-test-server-4", "spooky-test-server-5", "spooky-test-server-6", "spooky-test-server-7"]
  parallel    = true
} 