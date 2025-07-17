# Test Environment

[← Back to Index](index.md) | [Next: Troubleshooting](troubleshooting.md)

---

## Overview

Spooky provides a Podman-based test environment for integration testing, using rootless containers and Quadlet (systemd user units).

The test environment consists of three Debian 12 containers running SSH servers:

- **spooky-server1**: SSH on port 2221
- **spooky-server2**: SSH on port 2222  
- **spooky-server3**: SSH on port 2223

All servers are configured with:
- Root user with password: `password`
- SSH key authentication enabled
- Root login permitted for testing

## Prerequisites

- Podman 4.4+ (with Quadlet support)
- systemd user session enabled

## Usage

### Preflight Check
```bash
go run tools/spooky-test-env/main.go preflight
```

### Start Environment
```bash
go run tools/spooky-test-env/main.go start
```

This will:
1. Create the `spooky-test` network
2. Start all three containers using Podman
3. Display the container IPs

### Stop Environment
```bash
go run tools/spooky-test-env/main.go stop
```

### Status
```bash
go run tools/spooky-test-env/main.go status
```

### Cleanup
```bash
go run tools/spooky-test-env/main.go cleanup
```

This removes all containers and networks.

## SSH Connection Details

All servers accept connections with:
- **Host**: `localhost`
- **Port**: `2221`, `2222`, or `2223` (depending on server)
- **User**: `root`
- **Password**: `password`

## Example Configuration

```hcl
server "test-server1" {
  host     = "localhost"
  port     = 2221
  user     = "root"
  password = "password"
  tags = {
    environment = "test"
    datacenter  = "local"
  }
}

server "test-server2" {
  host     = "localhost"
  port     = 2222
  user     = "root"
  password = "password"
  tags = {
    environment = "test"
    datacenter  = "local"
  }
}

server "test-server3" {
  host     = "localhost"
  port     = 2223
  user     = "root"
  password = "password"
  tags = {
    environment = "test"
    datacenter  = "local"
  }
}

action "test-command" {
  description = "Test command execution"
  command     = "echo 'Hello from spooky!' && hostname"
  servers     = ["test-server1", "test-server2", "test-server3"]
  timeout     = 30
  parallel    = true
}
```

## Troubleshooting

### Port Conflicts

If you get port binding errors, ensure ports 2221-2223 are not in use:

```bash
sudo netstat -tlnp | grep :222
```

### Container Startup Issues

Check container logs:

```bash
podman logs spooky-server1
podman logs spooky-server2
podman logs spooky-server3
```

### Network Issues

If the network doesn't exist:

```bash
podman network create spooky-test
```

## Security Note

This test environment is configured for testing purposes only:
- Root login is enabled
- Simple password authentication
- No firewall or security hardening

**Do not use this configuration in production environments.**

---

## Navigation
- [← Back to Index](index.md)
- [Previous: Configuration Reference](configuration.md)
- [Next: Troubleshooting](troubleshooting.md) 