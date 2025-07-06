# Spooky - SSH Automation Tool

Spooky is a powerful SSH automation tool written in Go that allows you to execute commands and scripts on multiple remote servers using HCL2 configuration files. It provides a declarative way to manage server operations with support for parallel execution and flexible server targeting.

## Features

- 🚀 **Declarative Configuration**: Use HCL2 files to define servers and actions
- 🔗 **SSH Connectivity**: Secure connections with password or key-based authentication
- ⚡ **Parallel Execution**: Run actions on multiple servers simultaneously
- 🏷️ **Tag-based Targeting**: Target servers using tags for flexible grouping
- 📝 **Script Support**: Execute both inline commands and external script files
- ✅ **Validation**: Built-in configuration validation and syntax checking
- 🔍 **Listing**: View servers and actions defined in configuration files

## Installation

### Prerequisites

- Go 1.21 or later
- SSH access to target servers

### Build from Source

```bash
git clone <repository-url>
cd spooky
go mod tidy
go build -o spooky
```

## Usage

### Basic Commands

```bash
# Execute actions from configuration file
./spooky execute config.hcl

# Validate configuration file
./spooky validate config.hcl

# List servers and actions in configuration
./spooky list config.hcl

# Execute with parallel processing
./spooky execute -p config.hcl

# Set custom timeout (in seconds)
./spooky execute -t 60 config.hcl
```

### Configuration File Format

Spooky uses HCL2 configuration files to define servers and actions. Here's an example:

```hcl
# Define servers
server "web-server-1" {
  host     = "192.168.1.10"
  port     = 22
  user     = "admin"
  password = "your-password"
  # key_file = "~/.ssh/id_rsa"  # Alternative to password
  tags = {
    environment = "production"
    role        = "web"
  }
}

# Define actions
action "check-status" {
  description = "Check system status"
  command     = "uptime && df -h"
  servers     = ["web-server-1"]  # Target specific servers
  # tags = ["production"]         # Or target by tags
  parallel    = true              # Execute in parallel
}
```

### Server Configuration

Each server block defines a remote server:

- `name`: Unique identifier for the server
- `host`: Server IP address or hostname
- `port`: SSH port (default: 22)
- `user`: SSH username
- `password`: SSH password (or use `key_file`)
- `key_file`: Path to SSH private key file
- `tags`: Key-value pairs for server categorization

### Action Configuration

Each action block defines an operation to perform:

- `name`: Unique identifier for the action
- `description`: Human-readable description
- `command`: Inline command to execute
- `script`: Path to script file to execute
- `servers`: List of specific server names to target
- `tags`: List of tags to match servers
- `timeout`: Custom timeout for this action
- `parallel`: Execute on servers in parallel

### Server Targeting

Actions can target servers in three ways:

1. **Specific servers**: `servers = ["server1", "server2"]`
2. **Tag-based**: `tags = ["production", "web"]`
3. **All servers**: Omit both `servers` and `tags`

## Examples

### System Maintenance

```hcl
server "prod-web-1" {
  host     = "10.0.1.10"
  user     = "admin"
  password = "secure-password"
  tags = {
    environment = "production"
    role        = "web"
  }
}

server "prod-web-2" {
  host     = "10.0.1.11"
  user     = "admin"
  password = "secure-password"
  tags = {
    environment = "production"
    role        = "web"
  }
}

action "update-system" {
  description = "Update system packages"
  command     = "sudo apt update && sudo apt upgrade -y"
  tags        = ["production"]
  parallel    = true
}

action "restart-services" {
  description = "Restart web services"
  command     = "sudo systemctl restart nginx"
  tags        = ["web"]
  parallel    = true
}
```

### Database Operations

```hcl
server "db-master" {
  host     = "10.0.2.10"
  user     = "dbadmin"
  key_file = "~/.ssh/db_key"
  tags = {
    environment = "production"
    role        = "database"
  }
}

action "backup-database" {
  description = "Create database backup"
  script      = "./scripts/backup.sh"
  servers     = ["db-master"]
}
```

## Security Considerations

- Store sensitive information (passwords, keys) securely
- Use SSH key authentication when possible
- Limit SSH user permissions on target servers
- Consider using environment variables for credentials
- Regularly rotate passwords and keys

## Error Handling

Spooky provides detailed error reporting:

- Connection failures are reported per server
- Command execution errors include stderr output
- Configuration validation errors prevent execution
- Parallel execution continues even if some servers fail

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [golang.org/x/crypto/ssh](https://pkg.go.dev/golang.org/x/crypto/ssh) - SSH client
- [HCL2](https://github.com/hashicorp/hcl) - Configuration language 