# Spooky Facts Gathering System

A comprehensive facts gathering system that collects system information from remote machines via SSH.

## Overview

The facts gathering system collects three types of facts from remote machines:

1. **Basic Facts**: System information gathered via SSH commands
2. **Enhanced Facts**: Detailed system information from `spooky-facts` tool (if available)
3. **Custom Facts**: Age-encrypted custom facts (if available)

## Architecture

### Core Components

- **`Gatherer`**: Main facts collection engine
- **`MachineFacts`**: Container for facts from a single machine
- **SSH Integration**: Uses existing SSH infrastructure for remote connections

### Fact Types

#### Basic Facts
Collected via SSH commands:
- `hostname`: System hostname
- `os`: Operating system name
- `os_version`: OS version
- `architecture`: System architecture
- `kernel_version`: Kernel version
- `uptime`: System uptime
- `cpu_count`: Number of CPU cores
- `memory_total`: Total memory in bytes
- `memory_used`: Used memory in bytes
- `load_1/5/15`: Load averages
- `disk_usage`: Root filesystem usage
- `user`: Current user
- `home`: Home directory
- `shell`: Default shell

#### Enhanced Facts
Collected from `spooky-facts` tool:
- Comprehensive system information
- Hardware details
- Performance metrics
- Network information
- Process statistics

#### Custom Facts
Custom facts (can be age-encrypted or plain text):
- Sensitive information (typically encrypted)
- Custom metrics
- Application-specific data

## Usage

### CLI Commands

```bash
# Gather facts from all configured machines
spooky facts gather

# Gather facts and save to specific file
spooky facts gather /path/to/facts.hcl

# Export previously gathered facts
spooky facts export [output-file]
```

### Programmatic Usage

```go
// Create facts gatherer
sshManager := ssh.NewSimpleSSHManager(ageEncryption, sshConfig)
gatherer := facts.NewGatherer(sshManager, projectConfig)

// Gather facts from single machine
machineFacts, err := gatherer.GatherFactsFromMachine(ctx, machine)

// Gather facts from multiple machines in parallel
machineFacts, err := gatherer.GatherFactsFromMachines(ctx, machines)

// Export facts to HCL format
combinedFacts, err := gatherer.ExportFacts(machineFacts)
```

## Configuration

### Project Configuration

```hcl
project {
  facts_timeout = 30              # Timeout in seconds
  facts_parallel_collection = 10  # Max parallel connections
  facts_retry_attempts = 3        # Retry attempts
  facts_retry_delay = 5           # Delay between retries
}
```

### SSH Configuration

Uses existing SSH configuration from `spooky.hcl`:
- Connection timeouts
- Authentication methods
- Proxy settings
- Compression settings

## Integration Points

### Actions System
Facts can be used in action templates:
```hcl
# In action template
database_host = "{{ .Facts.basic_facts.hostname.value }}"
max_connections = {{ .Facts.basic_facts.cpu_count.value * 10 }}
```

### Export System
Facts can be exported to various formats:
- HCL files
- JSON format
- Template variables

## Error Handling

- **Graceful Degradation**: Missing enhanced/custom facts don't fail the operation
- **Parallel Processing**: Individual machine failures don't stop other machines
- **Retry Logic**: Configurable retry attempts for transient failures
- **Timeout Protection**: Configurable timeouts prevent hanging operations

## Security

- **SSH Authentication**: Uses existing SSH authentication methods
- **Age Encryption**: Custom facts can be age-encrypted
- **No Sensitive Data**: Basic facts avoid collecting sensitive information
- **Secure Transmission**: All data transmitted over encrypted SSH connections

## Performance

- **Parallel Collection**: Configurable parallel processing
- **Connection Pooling**: Reuses SSH connections where possible
- **Timeout Management**: Prevents resource exhaustion
- **Efficient Commands**: Uses lightweight system commands

## Future Enhancements

- [ ] HCL parsing for enhanced/custom facts
- [ ] Facts caching and incremental updates
- [ ] Facts validation and schema compliance
- [ ] Facts aggregation and summarization
- [ ] Facts history and trending
- [ ] Facts export to external systems
