server "test-server" {
  host = "localhost"
  port = 2222
  user = "testuser"
  key_file = "~/.ssh/id_ed25519"
  host_key_type = "insecure"
}

action "test-command" {
  description = "Test command execution"
  command = "echo 'Hello from spooky!' && whoami && pwd"
  servers = ["test-server"]
  parallel = false
  timeout = 30
}

action "test-script" {
  description = "Test script execution"
  script = "tests/integration/test-script.sh"
  servers = ["test-server"]
  parallel = false
  timeout = 30
}

action "test-parallel" {
  description = "Test parallel execution"
  command = "sleep 2 && echo 'Parallel execution test'"
  servers = ["test-server"]
  parallel = true
  timeout = 30
} 