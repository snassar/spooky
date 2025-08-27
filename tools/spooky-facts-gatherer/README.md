# Spooky Facts Gatherer

A comprehensive system information gatherer for the Spooky automation platform that collects detailed facts about the current system using [gopsutil](https://github.com/shirou/gopsutil) and outputs them in HCL format.

## Features

### 🔍 **Comprehensive System Information**
- **Host Information**: OS, platform, kernel version, virtualization details
- **CPU Details**: Vendor, model, cores, usage percentages
- **Memory Information**: RAM, swap, usage statistics
- **Disk Usage**: Partitions, mount points, filesystem types, usage stats
- **Network Interfaces**: MTU, flags, IP addresses, traffic statistics
- **System Load**: Load averages (1, 5, 15 minute)
- **Process Information**: Top processes by CPU usage
- **Environment Variables**: Common system environment variables
- **Runtime Information**: Go runtime details

### 🛠 **Cross-Platform Support**
- **Linux**: Full support for all features
- **macOS**: Full support for all features
- **Windows**: Full support for all features
- **BSD**: Full support for all features

### 📁 **Smart File Placement**
- **Unix Systems**: `/etc/spooky/facts.hcl` (or `~/.config/spooky/facts.hcl` for user mode)
- **Windows**: `%PROGRAMDATA%\spooky\facts.hcl` (or `%APPDATA%\spooky\facts.hcl` for user mode)

## Installation

### Build from Source
```bash
cd tools/spooky-facts-gatherer
go build -o spooky-facts .
```

### Install to System
```bash
# Build and install to system PATH
go install ./tools/spooky-facts-gatherer
```

## Usage

### Basic Commands

#### Gather Facts and Write to File
```bash
./spooky-facts gather
```
This command:
- Collects all available system information
- Creates the output directory if it doesn't exist
- Writes facts to the appropriate system location
- **Silent by default** - only shows output on errors

#### Gather Facts with Verbose Output
```bash
./spooky-facts gather --verbose
```
This command:
- Same as above but shows progress and success messages
- Useful for debugging or interactive use

#### Preview Facts (No File Writing)
```bash
./spooky-facts preview
```
This command:
- Collects all available system information
- Displays the generated HCL content to stdout
- Does not write any files

#### Help
```bash
./spooky-facts --help
./spooky-facts gather --help
./spooky-facts preview --help
```

### Example Output

```hcl
# System Facts gathered by spooky-facts
# Generated at: 2025-08-26T21:06:30+02:00

enhanced_facts {

  fact "hostname" {
    name = "hostname"
    value = "lathe"
    type = "string"
    description = "System hostname"
  }

  fact "os" {
    value = "linux"
    type = "string"
    description = "Operating system"
  }

  fact "platform" {
    name = "platform"
    value = "cachyos"
    type = "string"
    description = "Platform name"
  }

  fact "kernel_version" {
    name = "kernel_version"
    value = "6.16.2-2-cachyos"
    type = "string"
    description = "Kernel version"
  }

  fact "virtualization_system" {
    name = "virtualization_system"
    value = "kvm"
    type = "string"
    description = "Virtualization system"
  }

  fact "uptime" {
    name = "uptime"
    value = 127628
    type = "number"
    description = "System uptime in seconds"
  }

  fact "cpu_count" {
    name = "cpu_count"
    value = 12
    type = "number"
    description = "Number of CPU cores"
  }

  fact "cpu_model_name" {
    name = "cpu_model_name"
    value = "AMD Ryzen 5 PRO 4650G with Radeon Graphics"
    type = "string"
    description = "CPU model name"
  }

  fact "memory_total" {
    name = "memory_total"
    value = 14505799680
    type = "number"
    description = "Total memory in bytes"
  }

  fact "memory_used_percent" {
    name = "memory_used_percent"
    value = 80.15
    type = "number"
    description = "Memory usage percentage"
  }

  fact "load_1" {
    name = "load_1"
    value = 4.37
    type = "number"
    description = "1-minute load average"
  }

  fact "runtime_goos" {
    name = "runtime_goos"
    value = "linux"
    type = "string"
    description = "Go runtime OS"
  }

  fact "env_home" {
    name = "env_home"
    value = "/home/sn"
    type = "string"
    description = "Environment variable HOME"
  }

}
```

## Facts Structure

The facts gatherer generates a hierarchical facts structure that integrates with Spooky's in-memory facts representation:

```hcl
facts {
  basic_facts {}     # Gathered by normal commands run via SSH
  enhanced_facts {}  # Gathered from /etc/spooky/facts.hcl via spooky-facts
  custom_facts {}    # Gathered from /etc/spooky/custom.hcl (may contain age-encrypted facts)
}
```

Each `*_facts {}` block contains individual `fact "name" {}` blocks with the standard facts schema:
- `value` - Fact value (string, number, boolean, object, array)
- `type` - Data type
- `description` - Human-readable description
- `encrypted` - Whether the value is age-encrypted (optional)
- `tags` - Categorization tags (optional)
- `metadata` - Additional metadata (optional)

Note: The fact name is specified in the block label (`fact "name" {}`) and is not repeated as a field inside the block.

## Integration with Spooky

The facts gatherer is designed to integrate seamlessly with the Spooky automation platform:

### Template Usage
Facts can be used in Spooky templates:
```hcl
# In a template file
database_host = "{{ .Facts.enhanced_facts.hostname.value }}"
max_connections = {{ .Facts.enhanced_facts.cpu_count.value * 10 }}
memory_limit = "{{ .Facts.enhanced_facts.memory_total.value / 1024 / 1024 }}MB"
```

### Validation
The generated HCL file can be validated using Spooky's validation system:
```bash
spooky project validate /path/to/project
```

### Automation
The facts gatherer is designed for automation with silent operation:
```bash
# Gather facts before deployment (silent)
spooky-facts gather

# Gather facts with verbose output (for debugging)
spooky-facts gather --verbose

# Use facts in deployment
spooky deploy --facts-file /etc/spooky/facts.hcl
```

## Technical Details

### Dependencies
- **gopsutil/v4**: Cross-platform system information library
- **cobra**: CLI framework
- **spooky/internal/utilities**: Path configuration utilities

### Data Collection
The tool collects information from multiple sources:
- **System Calls**: Direct OS system calls via gopsutil
- **Proc Filesystem**: Linux `/proc` filesystem
- **Sysctl**: BSD/macOS system control interface
- **WMI**: Windows Management Instrumentation
- **Registry**: Windows registry (where applicable)

### Performance
- **Fast Collection**: Optimized for minimal system impact
- **Selective Gathering**: Only collects necessary information
- **Efficient Storage**: HCL format is human-readable and compact

### Security
- **Read-Only**: Only reads system information, never modifies
- **Permission Aware**: Respects system permissions and access controls
- **Safe Output**: No sensitive information is exposed in output

## Troubleshooting

### Common Issues

#### Permission Denied
```bash
Error: failed to gather disk partitions: permission denied
```
**Solution**: Run with appropriate permissions or use user mode paths.

#### Missing Information
Some information may not be available on all systems:
- Virtualization details (not available in all environments)
- Process information (may be limited by permissions)
- Network statistics (may vary by OS)

#### File Write Errors
```bash
Error: failed to write HCL file: permission denied
```
**Solution**: Ensure write permissions to the output directory.

### Debug Mode
For troubleshooting, you can examine the raw data:
```bash
# Preview mode shows all collected information
./spooky-facts preview
```

## Development

### Building
```bash
# Build for current platform
go build -o spooky-facts .

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o spooky-facts-linux .
GOOS=windows GOARCH=amd64 go build -o spooky-facts.exe .
GOOS=darwin GOARCH=amd64 go build -o spooky-facts-darwin .
```

### Testing
```bash
# Run tests
go test ./...

# Test with specific OS simulation
GOOS=linux go test ./...
```

### Contributing
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This tool is part of the Spooky automation platform and follows the same licensing terms.

## Related Tools

- **spooky-os-tool**: OS detection and path configuration
- **spooky-hcl-tool**: HCL file manipulation utilities
- **spooky-validation-rules-generator**: Schema validation rules
