# Usage & CLI

[← Back to Index](index.md) | [Next: Configuration Reference](configuration.md)

---

## Basic Commands

```bash
# Execute actions from configuration file
./spooky execute config.hcl

# Validate configuration file
./spooky validate config.hcl

# List servers and actions in configuration
./spooky list config.hcl
```

## Example Configuration

```hcl
server "web-server-1" {
  host     = "192.168.1.10"
  user     = "admin"
  password = "your-password"
  tags = {
    environment = "production"
    role        = "web"
  }
}

action "check-status" {
  description = "Check system status"
  command     = "uptime && df -h"
  servers     = ["web-server-1"]
  parallel    = true
}
```

---

## Navigation

- [← Back to Index](index.md)
- [Previous: Installation](install.md)
- [Next: Configuration Reference](configuration.md) 