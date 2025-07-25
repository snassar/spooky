# Usage & CLI

[← Back to Index](index.md) | [Next: Configuration Reference](configuration.md)

---

## Basic Commands

```bash
# Initialize a new project
spooky project init my-project ./path/to/project

# Validate project configuration
spooky project validate ./path/to/project

# List machines and actions in project
spooky project list ./path/to/project

# Run actions on machines
spooky action run --project ./path/to/project

# List facts about machines
spooky facts list --project ./path/to/project
```

## Example Configuration

### Project File (`project.hcl`)
```hcl
project "my-project" {
  description = "Example project with wrapper blocks"
  version     = "1.0.0"
  environment = "development"
  
  inventory_file = "inventory.hcl"
  actions_file   = "actions.hcl"
}
```

### Inventory File (`inventory.hcl`)
```hcl
inventory {
  machine "web-server-1" {
    host     = "192.168.1.10"
    user     = "admin"
    password = "your-password"
    tags = {
      environment = "production"
      role        = "web"
    }
  }
}
```

### Actions File (`actions.hcl`)
```hcl
actions {
  action "check-status" {
    description = "Check system status"
    command     = "uptime && df -h"
    machines    = ["web-server-1"]
    parallel    = true
  }
}
```

---

## Navigation

- [← Back to Index](index.md)
- [Previous: Installation](install.md)
- [Next: Configuration Reference](configuration.md) 