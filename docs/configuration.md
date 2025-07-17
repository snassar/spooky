# Configuration Reference

[← Back to Index](index.md) | [Next: Test Environment](test-environment.md)

---

## HCL2 Configuration Structure

Spooky uses HCL2 for all configuration files. Example:

```hcl
server "server-name" {
  host     = "192.168.1.100"
  port     = 22
  user     = "admin"
  password = "secret"  # or key_file = "~/.ssh/id_ed25519"
  tags = {
    environment = "production"
    datacenter  = "FRA00"
  }
}

action "update-system" {
  description = "Update system packages"
  command     = "apt update && apt upgrade -y"
  servers     = ["server1", "server2"]
  tags        = ["debian", "production"]
  timeout     = 300
  parallel    = true
}
```

## Server Block
- `host`: IP or hostname
- `port`: SSH port (default: 22)
- `user`: SSH username
- `password`: SSH password (or use `key_file`)
- `key_file`: Path to SSH private key
- `tags`: Key-value pairs for grouping

## Action Block
- `description`: Human-readable description
- `command`: Inline command to execute
- `script`: Path to script file
- `servers`: List of server names to target
- `tags`: List of tags to match servers
- `timeout`: Timeout in seconds
- `parallel`: Run in parallel (true/false)

---

## Navigation

- [← Back to Index](index.md)
- [Previous: Usage & CLI](usage.md)
- [Next: Test Environment](test-environment.md) 