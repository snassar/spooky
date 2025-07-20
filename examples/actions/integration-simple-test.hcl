# Simple test configuration for CLI testing in integration tests
machine "test-server-1" {
  host     = "localhost"
  port     = 2221
  user     = "testuser"
  key_file = "~/.ssh/id_ed25519"
}

machine "test-server-2" {
  host     = "localhost"
  port     = 2222
  user     = "testuser"
  key_file = "~/.ssh/id_ed25519"
}

action "test-action" {
  description = "Simple test action on 2 machines"
  command     = "echo 'test successful on ' $(hostname)"
  machines    = ["test-server-1", "test-server-2"]
  parallel    = true
} 