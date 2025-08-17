# Advanced Facts Collection Implementation Plan

## Overview

This document outlines the implementation plan for advanced fact collection using gopsutil and deployable collectors. This plan complements the SSH-based fact collection plan by providing comprehensive system metrics and privileged information collection capabilities.

## Purpose

The advanced facts collection system provides:
- **Comprehensive System Metrics**: Detailed CPU, memory, disk, and network information using gopsutil
- **Privileged Information**: Facts that require elevated access or are inefficient to gather via SSH
- **Real-time Monitoring**: Continuous fact collection and performance tracking
- **Custom Fact Collection**: Environment-specific fact collection logic
- **Reduced SSH Overhead**: Efficient collection of complex system data

## Architecture

### Advanced Collector Components

#### 1. Standalone Binary
- **Location**: `cmd/spooky-collector/`
- **Purpose**: Deployable Go binary for target machines
- **Dependencies**: gopsutil, spooky types and interfaces
- **Output**: File-based fact storage in `/etc/spooky/facts/`

#### 2. Systemd Service
- **Location**: `deploy/spooky-collector.service`
- **Purpose**: System service for continuous fact collection
- **User**: Runs as `spooky` user with appropriate permissions
- **Configuration**: Configurable collection intervals and output directories

#### 3. File-Based Output
- **Location**: `/etc/spooky/facts/`
- **Format**: JSON with timestamped files and latest symlink
- **Atomic Operations**: Temporary file writes with atomic moves
- **Permissions**: Readable by spooky user for SSH collection

## Implementation Details

### 1. Advanced Collector Binary

#### Main Application
```go
// cmd/spooky-collector/main.go
package main

import (
    "flag"
    "log"
    "os"
    "time"
    
    spookycollector "spooky/internal/facts/collectors/advanced"
    spookylogging "spooky/internal/logging"
)

var (
    outputDir string
    interval  time.Duration
)

func main() {
    flag.StringVar(&outputDir, "output-dir", "/etc/spooky/facts", "Output directory for facts")
    flag.DurationVar(&interval, "interval", 30*time.Second, "Collection interval")
    flag.Parse()
    
    // Setup logging
    logger := spookylogging.GetLogger()
    
    // Create collector
    collector := spookycollector.NewAdvancedCollector(outputDir, interval, logger)
    
    // Ensure output directory exists
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        logger.Error("Failed to create output directory", map[string]interface{}{
            "directory": outputDir,
            "error":     err.Error(),
        })
        os.Exit(1)
    }
    
    // Start collection loop
    collector.Run()
}
```

#### Advanced Collector Implementation
```go
// internal/facts/collectors/advanced/collector.go
package advanced

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "time"
    
    "github.com/shirou/gopsutil/v4/cpu"
    "github.com/shirou/gopsutil/v4/disk"
    "github.com/shirou/gopsutil/v4/host"
    "github.com/shirou/gopsutil/v4/load"
    "github.com/shirou/gopsutil/v4/mem"
    "github.com/shirou/gopsutil/v4/net"
    "github.com/shirou/gopsutil/v4/process"
    
    spookytypes "spooky/internal/types"
    spookytypesfacts "spooky/internal/types/facts"
    spookylogging "spooky/internal/logging"
)

type AdvancedCollector struct {
    outputDir string
    interval  time.Duration
    logger    spookylogging.Logger
}

func NewAdvancedCollector(outputDir string, interval time.Duration, logger spookylogging.Logger) *AdvancedCollector {
    return &AdvancedCollector{
        outputDir: outputDir,
        interval:  interval,
        logger:    logger,
    }
}

func (c *AdvancedCollector) Run() {
    ticker := time.NewTicker(c.interval)
    defer ticker.Stop()
    
    c.logger.Info("Starting advanced collector", map[string]interface{}{
        "outputDir": c.outputDir,
        "interval":  c.interval,
    })
    
    // Initial collection
    if err := c.collectAndWrite(); err != nil {
        c.logger.Error("Initial collection failed", map[string]interface{}{
            "error": err.Error(),
        })
    }
    
    // Collection loop
    for {
        select {
        case <-ticker.C:
            if err := c.collectAndWrite(); err != nil {
                c.logger.Error("Collection failed", map[string]interface{}{
                    "error": err.Error(),
                })
            }
        }
    }
}

func (c *AdvancedCollector) collectAndWrite() error {
    facts, err := c.CollectFacts(context.Background())
    if err != nil {
        return fmt.Errorf("failed to collect facts: %w", err)
    }
    
    if err := c.writeFacts(facts); err != nil {
        return fmt.Errorf("failed to write facts: %w", err)
    }
    
    return nil
}
```

### 2. Fact Collection Methods

#### System Information Collection
```go
func (c *AdvancedCollector) CollectFacts(ctx context.Context) (*spookytypes.FactCollection, error) {
    facts := &spookytypes.FactCollection{
        CollectedAt: time.Now(),
        Facts:       &spookytypesfacts.Facts{},
    }
    
    // Collect comprehensive system facts
    if err := c.collectSystemFacts(ctx, facts); err != nil {
        return nil, fmt.Errorf("failed to collect system facts: %w", err)
    }
    
    // Collect performance metrics
    if err := c.collectPerformanceFacts(ctx, facts); err != nil {
        return nil, fmt.Errorf("failed to collect performance facts: %w", err)
    }
    
    // Collect enhanced system data
    if err := c.collectEnhancedFacts(ctx, facts); err != nil {
        return nil, fmt.Errorf("failed to collect enhanced facts: %w", err)
    }
    
    return facts, nil
}

func (c *AdvancedCollector) collectSystemFacts(ctx context.Context, facts *spookytypes.FactCollection) error {
    // Host information
    hostInfo, err := host.Info()
    if err != nil {
        return fmt.Errorf("failed to get host info: %w", err)
    }
    
    facts.Facts.System = &spookytypesfacts.SystemFacts{
        OS: &spookytypesfacts.OSFacts{
            Name:     hostInfo.OS,
            Version:  hostInfo.PlatformVersion,
            Arch:     hostInfo.KernelArch,
            Kernel:   hostInfo.KernelVersion,
            Platform: hostInfo.Platform,
            Family:   hostInfo.PlatformFamily,
        },
    }
    
    // CPU information
    cpuInfo, err := cpu.Info()
    if err != nil {
        return fmt.Errorf("failed to get CPU info: %w", err)
    }
    
    if len(cpuInfo) > 0 {
        facts.Facts.System.Hardware = &spookytypesfacts.HardwareFacts{
            CPU: &spookytypesfacts.CPUFacts{
                Cores:     len(cpuInfo),
                Model:     cpuInfo[0].ModelName,
                Frequency: cpuInfo[0].Mhz,
                Vendor:    cpuInfo[0].VendorID,
            },
        }
    }
    
    // Memory information
    memInfo, err := mem.VirtualMemory()
    if err != nil {
        return fmt.Errorf("failed to get memory info: %w", err)
    }
    
    if facts.Facts.System.Hardware != nil {
        facts.Facts.System.Hardware.Memory = &spookytypesfacts.MemoryFacts{
            Total:     int64(memInfo.Total),
            Available: int64(memInfo.Available),
            Used:      int64(memInfo.Used),
            Free:      int64(memInfo.Free),
        }
    }
    
    return nil
}
```

#### Performance Metrics Collection
```go
func (c *AdvancedCollector) collectPerformanceFacts(ctx context.Context, facts *spookytypes.FactCollection) error {
    // CPU usage
    cpuPercent, err := cpu.Percent(0, false)
    if err == nil && len(cpuPercent) > 0 {
        if facts.Facts.System.Hardware != nil && facts.Facts.System.Hardware.CPU != nil {
            facts.Facts.System.Hardware.CPU.Percent = cpuPercent[0]
        }
    }
    
    // Memory usage
    memInfo, err := mem.VirtualMemory()
    if err == nil {
        if facts.Facts.System.Hardware != nil && facts.Facts.System.Hardware.Memory != nil {
            facts.Facts.System.Hardware.Memory.Percent = memInfo.UsedPercent
        }
    }
    
    // Load average
    loadAvg, err := load.Avg()
    if err == nil {
        facts.Facts.System.LoadAverage = &spookytypesfacts.LoadAverageFacts{
            Load1:  loadAvg.Load1,
            Load5:  loadAvg.Load5,
            Load15: loadAvg.Load15,
        }
    }
    
    // Disk I/O
    diskIO, err := disk.IOCounters()
    if err == nil {
        var totalRead, totalWrite uint64
        for _, io := range diskIO {
            totalRead += io.ReadBytes
            totalWrite += io.WriteBytes
        }
        
        facts.Facts.System.Hardware.DiskIO = &spookytypesfacts.DiskIOFacts{
            ReadBytes:  int64(totalRead),
            WriteBytes: int64(totalWrite),
        }
    }
    
    return nil
}
```

#### Enhanced System Data Collection
```go
func (c *AdvancedCollector) collectEnhancedFacts(ctx context.Context, facts *spookytypes.FactCollection) error {
    // Process information
    processes, err := process.Processes()
    if err == nil {
        facts.Facts.System.Processes = &spookytypesfacts.ProcessFacts{
            Count: len(processes),
        }
    }
    
    // Network information
    interfaces, err := net.Interfaces()
    if err == nil {
        var networkInterfaces []*spookytypesfacts.NetworkInterface
        for _, iface := range interfaces {
            networkInterface := &spookytypesfacts.NetworkInterface{
                Name:       iface.Name,
                MACAddress: iface.HardwareAddr,
                MTU:        iface.MTU,
            }
            networkInterfaces = append(networkInterfaces, networkInterface)
        }
        
        facts.Facts.System.Network = &spookytypesfacts.NetworkFacts{
            Interfaces: networkInterfaces,
        }
    }
    
    // Disk partitions
    partitions, err := disk.Partitions(false)
    if err == nil {
        var diskPartitions []*spookytypesfacts.DiskPartition
        for _, partition := range partitions {
            diskPartition := &spookytypesfacts.DiskPartition{
                Device:     partition.Device,
                Mountpoint: partition.Mountpoint,
                FSType:     partition.Fstype,
                Options:    partition.Opts,
            }
            diskPartitions = append(diskPartitions, diskPartition)
        }
        
        if facts.Facts.System.Hardware != nil {
            facts.Facts.System.Hardware.DiskPartitions = diskPartitions
        }
    }
    
    return nil
}
```

### 3. File Output Management

#### Atomic File Writing
```go
func (c *AdvancedCollector) writeFacts(facts *spookytypes.FactCollection) error {
    // Create timestamped filename
    timestamp := time.Now().Format("20060102-150405")
    filename := filepath.Join(c.outputDir, fmt.Sprintf("facts-%s.json", timestamp))
    
    // Marshal facts to JSON
    data, err := json.MarshalIndent(facts, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal facts: %w", err)
    }
    
    // Write to temporary file first, then atomically move
    tempFile := filename + ".tmp"
    if err := os.WriteFile(tempFile, data, 0644); err != nil {
        return fmt.Errorf("failed to write facts file: %w", err)
    }
    
    // Atomically move to final location
    if err := os.Rename(tempFile, filename); err != nil {
        os.Remove(tempFile) // Clean up temp file
        return fmt.Errorf("failed to move facts file: %w", err)
    }
    
    // Create symlink to latest
    latestLink := filepath.Join(c.outputDir, "facts-latest.json")
    os.Remove(latestLink) // Remove existing symlink if any
    if err := os.Symlink(filename, latestLink); err != nil {
        c.logger.Warn("Failed to create latest symlink", map[string]interface{}{
            "error": err.Error(),
        })
    }
    
    c.logger.Info("Wrote facts to file", map[string]interface{}{
        "file": filename,
    })
    
    return nil
}
```

### 4. Systemd Service Configuration

#### Service File
```ini
# deploy/spooky-collector.service
[Unit]
Description=Spooky Advanced Fact Collector
After=network.target

[Service]
Type=simple
User=spooky
Group=spooky
ExecStart=/usr/local/bin/spooky-collector --output-dir /etc/spooky/facts --interval 30s
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

#### Deployment Script
```bash
#!/bin/bash
# deploy/install-collector.sh

set -e

COLLECTOR_BINARY="spooky-collector"
SERVICE_FILE="spooky-collector.service"
INSTALL_DIR="/usr/local/bin"
SERVICE_DIR="/etc/systemd/system"
FACTS_DIR="/etc/spooky/facts"

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 
   exit 1
fi

# Create spooky user if it doesn't exist
if ! id "spooky" &>/dev/null; then
    useradd -r -s /bin/false -d /var/lib/spooky spooky
fi

# Create facts directory
mkdir -p "$FACTS_DIR"
chown spooky:spooky "$FACTS_DIR"
chmod 755 "$FACTS_DIR"

# Install binary
cp "$COLLECTOR_BINARY" "$INSTALL_DIR/"
chown root:root "$INSTALL_DIR/$COLLECTOR_BINARY"
chmod 755 "$INSTALL_DIR/$COLLECTOR_BINARY"

# Install service file
cp "$SERVICE_FILE" "$SERVICE_DIR/"
chown root:root "$SERVICE_DIR/$SERVICE_FILE"
chmod 644 "$SERVICE_DIR/$SERVICE_FILE"

# Reload systemd and enable service
systemctl daemon-reload
systemctl enable spooky-collector
systemctl start spooky-collector

echo "Spooky collector installed and started successfully"
```

## Integration with SSH Collection

### SSH Collector Integration
The SSH fact collector will automatically detect and read advanced collector output:

```go
// internal/facts/collectors/ssh/collector.go
func (c *SSHFactCollector) collectAdvancedFacts(ctx context.Context, machine *spookytypes.Machine) (*spookytypes.FactCollection, error) {
    // Read facts from advanced collector output files
    command := "cat /etc/spooky/facts/facts-latest.json 2>/dev/null || echo ''"
    output, err := c.executeCommand(ctx, machine, command)
    if err != nil || output == "" {
        return nil, fmt.Errorf("no advanced collector facts available")
    }
    
    var facts spookytypes.FactCollection
    if err := json.Unmarshal([]byte(output), &facts); err != nil {
        return nil, fmt.Errorf("failed to parse advanced collector facts: %w", err)
    }
    
    return &facts, nil
}
```

### Fallback Strategy
1. **Primary**: Read from advanced collector output files
2. **Fallback**: Use basic SSH command collection
3. **Graceful Degradation**: Continue with available facts

## Configuration Options

### Collector Configuration
```hcl
# /etc/spooky/collector.hcl
collector {
  output_dir = "/etc/spooky/facts"
  interval   = "30s"
  
  logging {
    level = "info"
    format = "json"
  }
  
  facts {
    system = true
    performance = true
    enhanced = true
    custom = true
  }
}
```

### Service Configuration
```bash
# Environment variables
SPOOKY_COLLECTOR_OUTPUT_DIR="/etc/spooky/facts"
SPOOKY_COLLECTOR_INTERVAL="30s"
SPOOKY_COLLECTOR_LOG_LEVEL="info"
```

## Security Considerations

### File Permissions
- **Output Directory**: `755` (spooky:spooky)
- **Fact Files**: `644` (spooky:spooky)
- **Service User**: `spooky` (non-privileged)

### Network Security
- **No Network Exposure**: File-based output only
- **No HTTP Endpoints**: Eliminates network attack surface
- **Local Access Only**: Facts accessible only via SSH

### Data Protection
- **Atomic Writes**: Prevents file corruption
- **Temporary Files**: Secure file creation process
- **Symlink Management**: Safe latest file linking

## Testing Strategy

### Unit Testing
- **Fact Collection**: Test individual collection methods
- **File Operations**: Test atomic file writing
- **Error Handling**: Test collection failures and recovery

### Integration Testing
- **Service Deployment**: Test systemd service installation
- **File Integration**: Test SSH collector reading advanced facts
- **Performance**: Test collection performance and resource usage

### Deployment Testing
- **Installation Script**: Test deployment automation
- **Service Management**: Test service start/stop/restart
- **Permission Handling**: Test file permission setup

## Success Criteria

1. **Comprehensive Fact Collection**: Successfully collect detailed system metrics using gopsutil
2. **File-Based Output**: Reliable file output with atomic operations
3. **Service Integration**: Proper systemd service deployment and management
4. **SSH Integration**: Seamless integration with SSH fact collection
5. **Performance**: Efficient collection with minimal resource usage
6. **Security**: Secure deployment and operation
7. **Reliability**: Robust error handling and recovery
8. **Usability**: Clear configuration and deployment process
9. **Testing**: Comprehensive test coverage
10. **Documentation**: Complete deployment and configuration documentation

## Next Steps

1. **Implementation**: Build advanced collector binary and service
2. **Deployment**: Create deployment scripts and documentation
3. **Integration**: Integrate with SSH fact collection
4. **Testing**: Comprehensive testing and validation
5. **Documentation**: User and deployment documentation
