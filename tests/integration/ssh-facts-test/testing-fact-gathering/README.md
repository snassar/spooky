# Testing Fact Gathering Project

This is a test project for validating SSH fact gathering functionality using a QEMU VM.

## Project Structure

```
testing-fact-gathering/
├── project.hcl          # Project configuration
├── machines.hcl         # Test VM definition
├── actions.hcl          # Test actions for fact gathering
├── variables.hcl        # Test variables
├── README.md           # This file
└── facts-output.json   # Generated fact output (after running tests)
```

## Usage

### Prerequisites

1. **VM Setup**: Ensure the test VM is running
   ```bash
   cd ..
   make start VM_NAME=spooky-facts-test
   ```

2. **SSH Connectivity**: Verify SSH connection works
   ```bash
   cd ..
   make test-ssh VM_NAME=spooky-facts-test
   ```

### Running Tests

1. **Test SSH Fact Gathering**:
   ```bash
   spooky actions run test-ssh-facts testing-fact-gathering
   ```

2. **Validate Collected Facts**:
   ```bash
   spooky actions run validate-facts testing-fact-gathering
   ```

3. **List Available Actions**:
   ```bash
   spooky actions list testing-fact-gathering
   ```

4. **Validate Project Configuration**:
   ```bash
   spooky project validate testing-fact-gathering
   ```

### Expected Facts

The SSH fact collector should gather the following facts from the VM:

#### System Facts
- `spooky_system`: "linux"
- `spooky_architecture`: "x86_64" or "amd64"
- `spooky_kernel`: Kernel version
- `spooky_hostname`: VM hostname
- `spooky_fqdn`: Fully qualified domain name

#### OS Facts
- `spooky_os_name`: "Debian GNU/Linux"
- `spooky_os_version`: Debian version
- `spooky_os_family`: "Debian"
- `spooky_distribution`: "debian"
- `spooky_distribution_version`: Debian version

#### Hardware Facts
- `spooky_processor`: CPU model
- `spooky_processor_cores`: Number of CPU cores (2)
- `spooky_processor_vcpus`: Number of virtual CPUs (2)
- `spooky_memtotal_mb`: Total memory in MB (2048)

#### Network Facts
- `spooky_default_ipv4`: Default IPv4 address and interface
- `spooky_interfaces`: List of network interfaces

#### User Facts
- `spooky_user_id`: "spooky"
- `spooky_user_dir`: User home directory
- `spooky_user_shell`: User shell

#### Environment Facts
- `spooky_env`: Environment variables

## Troubleshooting

### VM Not Running
```bash
cd ..
make status VM_NAME=spooky-facts-test
make start VM_NAME=spooky-facts-test
```

### SSH Connection Issues
```bash
cd ..
make test-ssh VM_NAME=spooky-facts-test
```

### Project Validation Errors
```bash
spooky project validate testing-fact-gathering
```

### Fact Collection Issues
1. Check VM is running and accessible
2. Verify SSH key permissions: `chmod 600 ../keys/spooky-facts-test_key`
3. Check spooky binary is built and in PATH
4. Review fact collector logs for errors

## Cleanup

To clean up the test environment:
```bash
cd ..
make clean-all VM_NAME=spooky-facts-test
``` 