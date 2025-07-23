# Spooky Test Environment

This directory contains the test environment configuration for spooky, which provides SSH-enabled containers for testing spooky's remote management capabilities.

## Overview

The test environment creates 9 containers that match the integration test workflow:

- **7 working SSH containers** (spooky-test-server-1 through spooky-test-server-7)
  - SSH server with key-based authentication
  - User: `testuser` with sudo access
  - Ports: 2221-2227

- **1 no-SSH container** (spooky-test-no-ssh)
  - Container without SSH server running
  - Used for testing failure scenarios
  - Port: 2228

- **1 SSH-no-key container** (spooky-test-ssh-no-key)
  - SSH server without authorized keys
  - Used for testing authentication failures
  - Port: 2229

## Prerequisites

- Podman 5.2+ with rootless support
- systemd (Linux only)
- Quadlet support
- SSH key pair (default: `~/.ssh/id_ed25519`)

## SSH Key Configuration

The tool supports multiple ways to specify the SSH key:

1. **Command line flag:**
   ```bash
   go run tools/spooky-test-env/main.go --ssh-key /path/to/your/key start
   ```

2. **Environment variable:**
   ```bash
   export SPOOKY_TEST_SSH_KEY=/path/to/your/key
   go run tools/spooky-test-env/main.go start
   ```

3. **Default location:** `~/.ssh/id_ed25519`

The tool will automatically generate a key at the specified location if it doesn't exist.

## Quick Start

1. **Check prerequisites:**
   ```bash
   go run tools/spooky-test-env/main.go preflight
   ```

2. **Build container images:**
   ```bash
   go run tools/spooky-test-env/main.go build
   ```

3. **Start the test environment:**
   ```bash
   go run tools/spooky-test-env/main.go start
   ```

4. **Check status:**
   ```bash
   go run tools/spooky-test-env/main.go status
   ```

5. **Stop the environment:**
   ```bash
   go run tools/spooky-test-env/main.go stop
   ```

6. **Clean up everything:**
   ```bash
   go run tools/spooky-test-env/main.go cleanup
   ```

## Container Images

### spooky-test-ssh
- **Base:** debian:12-slim
- **Features:** SSH server with key-based authentication
- **User:** testuser (with sudo access)
- **Authentication:** SSH public key only (no password)

### spooky-test-no-ssh
- **Base:** debian:12-slim
- **Features:** No SSH server
- **User:** testuser
- **Purpose:** Test failure scenarios

### spooky-test-ssh-no-key
- **Base:** debian:12-slim
- **Features:** SSH server without authorized keys
- **User:** testuser
- **Purpose:** Test authentication failures

## SSH Configuration

All SSH containers use:
- **User:** testuser
- **Authentication:** Public key only (no password)
- **Key:** `~/.ssh/id_ed25519` (automatically generated if missing)
- **Sudo:** testuser has passwordless sudo access
- **Security:** Root login disabled, strict modes disabled

## Quadlet Integration

The test environment uses Podman Quadlet for container management:

- **Container files:** Generated automatically as `.container` files
- **Systemd integration:** Containers run as systemd user services
- **Automatic startup:** Containers start automatically with systemd
- **Clean shutdown:** Proper systemd service management

## File Structure

```
examples/test-environment/
├── README.md                    # This file
├── Containerfile.ssh            # SSH container definition
├── Containerfile.no-ssh         # No-SSH container definition
├── Containerfile.ssh-no-key     # SSH-no-key container definition
├── test-config.hcl              # Test configuration
├── spooky-test-server.container # Legacy container file
├── spooky-test.pod              # Legacy pod file
├── spooky-test.network          # Legacy network file
└── Containerfile                # Legacy container file
```

## Testing with spooky

Once the test environment is running, you can test spooky with:

```bash
# Test SSH fact collection (using default key)
./build/spooky facts gather --ssh-user testuser --ssh-key ~/.ssh/id_ed25519 localhost:2221

# Test SSH fact collection (using custom key)
./build/spooky facts gather --ssh-user testuser --ssh-key /path/to/your/key localhost:2221

# Test with multiple servers
./build/spooky execute examples/actions/integration-working-servers.hcl

# Test failure scenarios
./build/spooky execute examples/actions/integration-failure-test.hcl
```

## SSH Access to Containers

The containers are accessible via SSH using port forwarding:

### Working SSH Containers (ports 2221-2227)
```bash
# SSH to container 1
ssh -i ~/testkeys/spooky/containers -p 2221 testuser@localhost

# SSH to container 2
ssh -i ~/testkeys/spooky/containers -p 2222 testuser@localhost

# And so on for containers 3-7
```

### Container IP Addresses
The containers also have static IP addresses on the internal `spooky-test` network:
- spooky-test-server-1: 10.1.10.1
- spooky-test-server-2: 10.1.10.2
- spooky-test-server-3: 10.1.10.3
- spooky-test-server-4: 10.1.10.4
- spooky-test-server-5: 10.1.10.5
- spooky-test-server-6: 10.1.10.6
- spooky-test-server-7: 10.1.10.7
- spooky-test-no-ssh: 10.1.10.8
- spooky-test-ssh-no-key: 10.1.10.9

**Note**: The IP addresses are only accessible from within other containers on the same network. For external access, use the port-forwarded localhost connections.

## Troubleshooting

### SSH Key Issues
If you don't have an SSH key or want to use a different one:

```bash
# Generate default key
ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -N ""

# Or generate a custom key
ssh-keygen -t ed25519 -f /path/to/your/key -N ""

# Or let the tool generate it automatically
go run tools/spooky-test-env/main.go build
```

### Quadlet Issues
If Quadlet is not working:
```bash
# Check Quadlet support
podman quadlet --help

# Check systemd user services
systemctl --user list-units --type=service
```

### Container Issues
If containers fail to start:
```bash
# Check container status
podman ps -a

# Check container logs
podman logs spooky-test-server-1

# Check systemd service status
systemctl --user status spooky-test-server-1
```

### Network Issues
If containers can't reach each other:
```bash
# Check network configuration
podman network ls
podman network inspect bridge
```

## Integration with CI

This test environment matches the containers created in the GitHub Actions integration tests workflow (`.github/workflows/integration-tests.yml`), ensuring consistent testing between local development and CI environments. 