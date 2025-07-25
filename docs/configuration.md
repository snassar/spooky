# Configuration Reference

[← Back to Index](index.md) | [Next: Test Environment](test-environment.md)

---

## HCL2 Configuration Structure

Spooky uses HCL2 for all configuration files with wrapper blocks for explicit file type declaration.

### Inventory File (`inventory.hcl`)
```hcl
# Inventory configuration with wrapper block
inventory {
  machine "server-name" {
    host     = "192.168.1.100"
    port     = 22
    user     = "admin"
    password = "secret"  # or key_file = "~/.ssh/id_ed25519"
    tags = {
      environment = "production"
      datacenter  = "FRA00"
    }
  }
}
```

### Actions File (`actions.hcl`)
```hcl
# Actions configuration with wrapper block
actions {
  action "update-system" {
    description = "Update system packages"
    command     = "apt update && apt upgrade -y"
    machines    = ["server1", "server2"]
    tags        = ["debian", "production"]
    timeout     = 300
    parallel    = true
  }
}
```

## Machine Block
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
- `machines`: List of machine names to target
- `tags`: List of tags to match machines
- `timeout`: Timeout in seconds
- `parallel`: Run in parallel (true/false)

## Wrapper Block Benefits

The wrapper block format provides several advantages:

- **Explicit file type declaration**: Clear indication of what each file contains
- **Better tooling integration**: IDEs and tools can provide better autocomplete and validation
- **Future extensibility**: Foundation for advanced features like multi-file support
- **Consistent structure**: Matches the `project {}` wrapper block pattern

---

## Navigation

- [← Back to Index](index.md)
- [Previous: Usage & CLI](usage.md)
- [Next: Test Environment](test-environment.md) 