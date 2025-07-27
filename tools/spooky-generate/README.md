# Spooky Generate

A command-line tool for generating spooky configuration files for testing purposes.

## Features

- Generate valid HCL configuration files for spooky projects
- Create diverse actions covering system administration tasks
- Generate machines with realistic roles and environments
- Intertwined tags between actions and inventory for realistic testing
- Support for different output locations and file naming

## Usage

### Basic Usage

Generate both actions and inventory files with default settings:
```bash
./spooky-generate
```

This creates:
- `./actions.hcl` with 25 diverse actions
- `./inventory.hcl` with 10 machines

### Command Line Options

```bash
./spooky-generate [flags]
```

**Flags:**
- `--generate, -g` - What to generate: `all` (default) or `inventory-only`
- `--actions, -a` - Number of actions to generate (default: 25)
- `--machines, -m` - Number of machines to generate (default: 10)
- `--output-actions` - Output path for actions.hcl (default: ./actions.hcl)
- `--output-inventory` - Output path for inventory.hcl (default: ./inventory.hcl)

### Examples

Generate only inventory with 100 machines:
```bash
./spooky-generate --generate inventory-only --machines 100
```

Generate 50 actions and 20 machines in a specific directory:
```bash
./spooky-generate --actions 50 --machines 20 --output-actions /tmp/test-project/ --output-inventory /tmp/test-project/
```

Generate files with custom names:
```bash
./spooky-generate --output-actions /tmp/my-actions.hcl --output-inventory /tmp/my-inventory.hcl
```

## Generated Content

### Actions

The tool generates diverse actions covering:
- **System monitoring**: uptime, disk space, memory, load
- **Network connectivity**: ping, DNS, port checking
- **User management**: create, delete, list users
- **Package management**: update, upgrade, install, remove packages
- **Service management**: start, stop, restart, enable, disable systemd services
- **File operations**: create directories, remove files, backups
- **Log management**: check logs, rotate logs, clear old logs
- **Security**: SSH key management, SSH configuration

### Machines

Machines are generated with:
- **IP addresses**: Using 10.0.0.0/8 private network space
- **Roles**: web, database, cache, load-balancer, monitoring, app, api, worker, storage, backup
- **Environments**: production, staging, development, testing
- **Regions**: us-east, us-west, eu-west, eu-central, ap-southeast

### Tag Integration

Actions and machines use intertwined tags:
- Actions target machines by role tags (e.g., `role=web`, `role=database`)
- Machines have matching role tags
- This creates realistic targeting scenarios for testing

## Building

```bash
cd tools/spooky-generate
go build -o spooky-generate
```

## Integration with Spooky

The generated files are valid HCL and can be used directly in spooky projects:

1. Generate configuration files:
   ```bash
   ./spooky-generate --actions 100 --machines 50
   ```

2. Copy to a spooky project:
   ```bash
   cp actions.hcl inventory.hcl /path/to/spooky-project/
   ```

3. Use with spooky commands:
   ```bash
   spooky validate
   spooky list-machines
   spooky list-actions
   ```

## Output Format

### Actions File
```hcl
# Generated actions for testing
# Contains 25 diverse actions covering system administration tasks

actions {
  action "check-uptime-1" {
    description = "Check system uptime"
    command = "uptime"
    tags = ["role=web", "role=database", "role=cache"]
    timeout = 300
    parallel = true
  }
  # ... more actions
}
```

### Inventory File
```hcl
# Generated inventory for testing
# Contains 10 machines with diverse roles and environments

inventory {
  machine "server-1" {
    host = "10.0.1.1"
    port = 22
    user = "admin"
    password = "password123"
    tags = {
      role = "web"
      environment = "production"
      region = "us-east"
      instance = "1"
    }
  }
  # ... more machines
}
``` 