# SSH Fact Gathering Test Environment

This directory contains a Makefile-based system to set up QEMU virtual machines for testing SSH fact gathering functionality.

## Overview

The SSH fact collector in `internal/facts/collectors/ssh/collector.go` gathers various system facts from remote machines via SSH. This test environment provides a way to test this functionality using QEMU virtual machines managed through a Makefile.

## Prerequisites

- QEMU (qemu-system-x86_64)
- Make
- wget (for ISO downloads)
- SSH key generation tools

## Quick Start

### 1. Show Available Commands

```bash
cd tests/integration/ssh-facts-test
make help
```

### 2. Set Up Ubuntu VM

```bash
# Download Ubuntu ISO and create VM
make setup-ubuntu VM_NAME=ubuntu-server-22.04

# Install Ubuntu (interactive)
make install-os VM_NAME=ubuntu-server-22.04

# Start the VM
make start VM_NAME=ubuntu-server-22.04
```

### 3. Test SSH Fact Gathering

```bash
# Test basic SSH connectivity
make test-ssh VM_NAME=ubuntu-server-22.04

# Run SSH fact gathering test using spooky project
make test-project VM_NAME=ubuntu-server-22.04
```

## VM Management

### Creating a New VM

#### Method 1: Manual Setup (Interactive)
```bash
# Download ISO and create VM
make setup-ubuntu VM_NAME=<vm-name>
make setup-fedora VM_NAME=<vm-name>
make setup-debian VM_NAME=<vm-name>

# Or manually
make download-ubuntu
make create VM_NAME=<vm-name>
```

This will:
- Download the specified ISO if not present
- Create a QCOW2 disk image
- Set up SSH key authentication

#### Method 2: Cloud-init Setup (Automated SSH Key Configuration)
```bash
# Download ISO and create VM with cloud-init
make setup-ubuntu-cloud-init VM_NAME=<vm-name>
make setup-fedora-cloud-init VM_NAME=<vm-name>
make setup-debian-cloud-init VM_NAME=<vm-name>
```

This will:
- Download the specified ISO if not present
- Create a QCOW2 disk image
- Generate SSH key pair
- Create cloud-init configuration with SSH keys
- Automatically configure SSH access after OS installation

**Benefits of cloud-init:**
- SSH keys are automatically configured during first boot
- No manual SSH setup required
- Consistent user and SSH configuration
- Can be used for automated testing

### Starting a VM

#### Standard VM
```bash
make start VM_NAME=<vm-name>
```

#### VM with Cloud-init
```bash
make start-cloud-init VM_NAME=<vm-name>
```

The VM will start with:
- 2GB RAM
- 2 CPU cores
- SSH port forwarded to localhost (port 2222)
- Network access for package installation
- Cloud-init ISO attached (for automated setup)

### Installing OS

#### Standard Installation
```bash
make install-os VM_NAME=<vm-name>
```

#### Installation with Cloud-init
```bash
make install-os-cloud-init VM_NAME=<vm-name>
```

This will:
- Start interactive OS installation
- Attach cloud-init ISO for automated SSH configuration
- After installation, SSH keys will be automatically configured

### Stopping a VM

```bash
make stop VM_NAME=<vm-name>
```

### Connecting to VM

```bash
# SSH into the VM
make connect VM_NAME=<vm-name>

# Or manually
ssh -i keys/<vm-name>_key -p 2222 spooky@localhost
```

### Building Spooky

The test environment automatically builds the spooky binary locally:

```bash
# Build spooky binary (automatic when needed)
make build-spooky

# Force rebuild spooky binary
make rebuild-spooky
```

The binary is built from the project root and placed in the test directory.

### Cleaning Up

```bash
# Remove everything, leaving directory in git clone state
make clean
```

This removes all:
- VM disk images
- SSH keys
- ISO files
- Log files
- Cloud-init files
- PID files
- Spooky binary

## Testing SSH Fact Gathering

The SSH fact collector gathers the following information:

### System Facts
- `spooky_system`: OS type (linux, darwin, etc.)
- `spooky_architecture`: CPU architecture
- `spooky_kernel`: Kernel version
- `spooky_hostname`: System hostname
- `spooky_fqdn`: Fully qualified domain name

### OS Facts
- `spooky_os_name`: OS name (from /etc/os-release)
- `spooky_os_version`: OS version
- `spooky_os_family`: OS family
- `spooky_distribution`: Distribution name
- `spooky_distribution_version`: Distribution version

### Hardware Facts
- `spooky_processor`: CPU model
- `spooky_processor_cores`: Number of CPU cores
- `spooky_processor_vcpus`: Number of virtual CPUs
- `spooky_memtotal_mb`: Total memory in MB

### Network Facts
- `spooky_default_ipv4`: Default IPv4 address and interface
- `spooky_interfaces`: List of network interfaces

### User Facts
- `spooky_user_id`: Current user ID
- `spooky_user_dir`: User home directory
- `spooky_user_shell`: User shell

### Environment Facts
- `spooky_env`: Environment variables

## Test Commands

### test-ssh
Tests basic SSH connectivity to the VM.

```bash
make test-ssh VM_NAME=<vm-name>
```

### test-facts
Runs the SSH fact collector against the VM and displays results.

```bash
make test-facts VM_NAME=<vm-name>
```

### test-facts-json
Runs the SSH fact collector and outputs results in JSON format.

```bash
make test-facts-json VM_NAME=<vm-name>
```

## VM Configuration

VMs are configured with:
- **Memory**: 2GB
- **CPUs**: 2 cores
- **Disk**: 20GB QCOW2
- **Network**: User networking with port forwarding
- **SSH**: Port 2222 forwarded to localhost
- **User**: ubuntu (Ubuntu) or fedora (Fedora)
- **Authentication**: SSH key-based

## Troubleshooting

### VM won't start
- Check if QEMU is installed: `which qemu-system-x86_64`
- Verify ISO file exists in `isos/` directory
- Check disk space for QCOW2 images
- Check VM status: `make status VM_NAME=<vm-name>`

### SSH connection fails
- Verify VM is running: `make status VM_NAME=<vm-name>`
- Check SSH port forwarding: `ss -tlnp | grep 2222`
- Verify SSH key permissions: `chmod 600 keys/<vm-name>_key`
- Test SSH connectivity: `make test-ssh VM_NAME=<vm-name>`

### Fact gathering fails
- Check SSH connectivity first: `make test-ssh VM_NAME=<vm-name>`
- Verify user permissions on the VM
- Check if required commands are available (uname, hostname, etc.)
- Ensure spooky binary is built: `go build -o spooky main.go`

## Supported Distributions

- Ubuntu Server 22.04 LTS
- Fedora 39
- CentOS Stream 9
- Debian 12

## Directory Structure

```
tests/integration/ssh-facts-test/
├── README.md              # This file
├── Makefile              # VM management commands
├── testing-fact-gathering/ # Test project for SSH fact gathering
│   ├── project.hcl       # Project configuration
│   ├── machines.hcl      # Test VM definition
│   ├── actions.hcl       # Test actions
│   ├── variables.hcl     # Test variables
│   └── README.md         # Test project documentation
├── isos/                 # ISO files directory
├── vms/                  # VM disk images
├── keys/                 # SSH keys
└── logs/                 # VM logs and PID files
```

## Available Make Targets

Run `make help` to see all available commands, or use:

### VM Management
- `make setup-debian` - Download Debian ISO and create VM (default)
- `make setup-ubuntu` - Download Ubuntu ISO and create VM
- `make setup-fedora` - Download Fedora ISO and create VM
- `make start` - Start VM in background
- `make stop` - Stop VM
- `make status` - Show VM status
- `make connect` - SSH into VM

### Testing
- `make test-ssh` - Test SSH connectivity
- `make test-facts` - Test SSH fact gathering (direct)
- `make test-project` - Test SSH fact gathering using test project
- `make validate-project` - Validate test project configuration
- `make list-actions` - List available actions in test project

### Utilities
- `make list` - List all VMs and ISOs
- `make info` - Show VM information
- `make clean-all` - Remove VM completely 