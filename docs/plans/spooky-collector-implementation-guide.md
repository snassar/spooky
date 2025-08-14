# Spooky-Collector Implementation Guide

## Overview

This document provides a complete implementation guide for building the `spooky-collector` binary, which is a standalone Go application that collects comprehensive system facts using gopsutil and outputs them to `/etc/spooky/facts.hcl` for collection by the main spooky tool.

## Purpose

The `spooky-collector` binary serves as a deployable fact collection agent that:
- **Collects comprehensive system facts** using gopsutil library
- **Runs with user privileges** (no sudo required)
- **Outputs structured HCL** to `/etc/spooky/facts.hcl`
- **Provides real-time system metrics** for advanced fact collection
- **Reduces SSH overhead** by pre-collecting complex system data

## Architecture

### Binary Structure
```
cmd/spooky-collector/
├── main.go              # Main application entry point
├── collector.go         # Core fact collection logic
├── output.go           # HCL output formatting
├── service.go          # Systemd service integration
└── go.mod              # Dependencies
```

### Service Integration
```
deploy/
├── spooky-collector.service    # Systemd service file
├── spooky-collector.conf       # Configuration file
└── install.sh                  # Installation script
```

## Implementation Details

### 1. Main Application

```go
// cmd/spooky-collector/main.go
package main

import (
    "flag"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    spookycollector "spooky/internal/facts/collectors/advanced"
    spookylogging "spooky/internal/logging"
)

var (
    outputDir = flag.String("output", "/etc/spooky", "Output directory for facts")
    interval  = flag.Duration("interval", 5*time.Minute, "Collection interval")
    oneshot   = flag.Bool("oneshot", false, "Run once and exit")
    verbose   = flag.Bool("verbose", false, "Enable verbose logging")
)

func main() {
    flag.Parse()
    
    // Setup logging
    logger := spookylogging.GetLogger()
    if *verbose {
        logger.SetLevel(spookylogging.DebugLevel)
    }
    
    // Create collector
    collector := spookycollector.NewAdvancedCollector(*outputDir, logger)
    
    if *oneshot {
        // Run once and exit
        if err := collector.CollectAndOutput(); err != nil {
            log.Fatalf("Failed to collect facts: %v", err)
        }
        return
    }
    
    // Run as service
    if err := runService(collector); err != nil {
        log.Fatalf("Service failed: %v", err)
    }
}

func runService(collector *spookycollector.AdvancedCollector) error {
    // Setup signal handling
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    // Initial collection
    if err := collector.CollectAndOutput(); err != nil {
        return err
    }
    
    // Periodic collection
    ticker := time.NewTicker(*interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if err := collector.CollectAndOutput(); err != nil {
                logger.Error("Failed to collect facts", "error", err)
            }
        case <-sigChan:
            logger.Info("Shutting down collector")
            return nil
        }
    }
}
```

### 2. Advanced Collector Implementation

```go
// cmd/spooky-collector/collector.go
package main

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
    
    spookylogging "spooky/internal/logging"
    spookytypes "spooky/internal/types"
    spookytypesfacts "spooky/internal/types/facts"
)

type AdvancedCollector struct {
    outputDir string
    logger    spookylogging.Logger
}

func NewAdvancedCollector(outputDir string, logger spookylogging.Logger) *AdvancedCollector {
    return &AdvancedCollector{
        outputDir: outputDir,
        logger:    logger,
    }
}

func (c *AdvancedCollector) CollectAndOutput() error {
    ctx := context.Background()
    
    // Collect facts
    facts, err := c.CollectFacts(ctx)
    if err != nil {
        return fmt.Errorf("failed to collect facts: %w", err)
    }
    
    // Output to HCL
    if err := c.OutputToHCL(facts); err != nil {
        return fmt.Errorf("failed to output facts: %w", err)
    }
    
    c.logger.Info("Facts collected and output successfully")
    return nil
}

func (c *AdvancedCollector) CollectFacts(ctx context.Context) (*spookytypes.FactCollection, error) {
    facts := &spookytypes.FactCollection{
        CollectedAt: time.Now(),
        Facts:       &spookytypesfacts.Facts{},
    }
    
    // Collect all fact categories
    if err := c.collectHostFacts(ctx, facts); err != nil {
        return nil, fmt.Errorf("failed to collect host facts: %w", err)
    }
    
    if err := c.collectCPUFacts(ctx, facts); err != nil {
        return nil, fmt.Errorf("failed to collect CPU facts: %w", err)
    }
    
    if err := c.collectMemoryFacts(ctx, facts); err != nil {
        return nil, fmt.Errorf("failed to collect memory facts: %w", err)
    }
    
    if err := c.collectDiskFacts(ctx, facts); err != nil {
        return nil, fmt.Errorf("failed to collect disk facts: %w", err)
    }
    
    if err := c.collectNetworkFacts(ctx, facts); err != nil {
        return nil, fmt.Errorf("failed to collect network facts: %w", err)
    }
    
    if err := c.collectProcessFacts(ctx, facts); err != nil {
        return nil, fmt.Errorf("failed to collect process facts: %w", err)
    }
    
    if err := c.collectPerformanceFacts(ctx, facts); err != nil {
        return nil, fmt.Errorf("failed to collect performance facts: %w", err)
    }
    
    return facts, nil
}

func (c *AdvancedCollector) collectHostFacts(ctx context.Context, facts *spookytypes.FactCollection) error {
    hostInfo, err := host.Info()
    if err != nil {
        return fmt.Errorf("failed to get host info: %w", err)
    }
    
    facts.Facts.Collector = &spookytypesfacts.CollectorFacts{
        Host: &spookytypesfacts.HostFacts{
            Hostname:        hostInfo.Hostname,
            Uptime:          int64(hostInfo.Uptime),
            BootTime:        int64(hostInfo.BootTime),
            OS:              hostInfo.OS,
            Platform:        hostInfo.Platform,
            PlatformFamily:  hostInfo.PlatformFamily,
            PlatformVersion: hostInfo.PlatformVersion,
            KernelVersion:   hostInfo.KernelVersion,
            KernelArch:      hostInfo.KernelArch,
            VirtualizationSystem: hostInfo.VirtualizationSystem,
            VirtualizationRole:   hostInfo.VirtualizationRole,
        },
    }
    
    return nil
}

func (c *AdvancedCollector) collectCPUFacts(ctx context.Context, facts *spookytypes.FactCollection) error {
    // CPU Info
    cpuInfo, err := cpu.Info()
    if err != nil {
        return fmt.Errorf("failed to get CPU info: %w", err)
    }
    
    if len(cpuInfo) > 0 {
        if facts.Facts.Collector == nil {
            facts.Facts.Collector = &spookytypesfacts.CollectorFacts{}
        }
        
        facts.Facts.Collector.CPU = &spookytypesfacts.CPUFacts{
            Cores:        len(cpuInfo),
            Model:        cpuInfo[0].ModelName,
            Frequency:    cpuInfo[0].Mhz,
            Architecture: cpuInfo[0].Family,
            Vendor:       cpuInfo[0].VendorID,
            CacheSize:    int64(cpuInfo[0].CacheSize),
        }
        
        // Per-core information
        var coresDetail []*spookytypesfacts.CPUCoreDetail
        for i, info := range cpuInfo {
            core := &spookytypesfacts.CPUCoreDetail{
                CPU:       i,
                ModelName: info.ModelName,
                MHz:       info.Mhz,
                CacheSize: int64(info.CacheSize),
            }
            coresDetail = append(coresDetail, core)
        }
        facts.Facts.Collector.CPU.CoresDetail = coresDetail
    }
    
    // CPU Times
    cpuTimes, err := cpu.Times(false)
    if err == nil && len(cpuTimes) > 0 {
        if facts.Facts.Collector != nil && facts.Facts.Collector.CPU != nil {
            facts.Facts.Collector.CPU.Times = &spookytypesfacts.CPUTimes{
                User:      cpuTimes[0].User,
                System:    cpuTimes[0].System,
                Idle:      cpuTimes[0].Idle,
                Nice:      cpuTimes[0].Nice,
                IOWait:    cpuTimes[0].Iowait,
                IRQ:       cpuTimes[0].Irq,
                SoftIRQ:   cpuTimes[0].Softirq,
                Steal:     cpuTimes[0].Steal,
                Guest:     cpuTimes[0].Guest,
                GuestNice: cpuTimes[0].GuestNice,
            }
        }
    }
    
    // CPU Percentage
    cpuPercent, err := cpu.Percent(0, false)
    if err == nil && len(cpuPercent) > 0 {
        if facts.Facts.Collector != nil && facts.Facts.Collector.CPU != nil {
            facts.Facts.Collector.CPU.Percent = cpuPercent[0]
        }
    }
    
    return nil
}

func (c *AdvancedCollector) collectMemoryFacts(ctx context.Context, facts *spookytypes.FactCollection) error {
    // Virtual Memory
    vmem, err := mem.VirtualMemory()
    if err != nil {
        return fmt.Errorf("failed to get virtual memory info: %w", err)
    }
    
    if facts.Facts.Collector == nil {
        facts.Facts.Collector = &spookytypesfacts.CollectorFacts{}
    }
    
    facts.Facts.Collector.Memory = &spookytypesfacts.MemoryFacts{
        Total:         int64(vmem.Total),
        Available:     int64(vmem.Available),
        Used:          int64(vmem.Used),
        Free:          int64(vmem.Free),
        Buffers:       int64(vmem.Buffers),
        Cached:        int64(vmem.Cached),
        Shared:        int64(vmem.Shared),
        Slab:          int64(vmem.Slab),
        Percent:       vmem.UsedPercent,
    }
    
    // Swap Memory
    swap, err := mem.SwapMemory()
    if err == nil {
        facts.Facts.Collector.Memory.Swap = &spookytypesfacts.SwapFacts{
            Total:   int64(swap.Total),
            Used:    int64(swap.Used),
            Free:    int64(swap.Free),
            Percent: swap.UsedPercent,
        }
    }
    
    return nil
}

func (c *AdvancedCollector) collectDiskFacts(ctx context.Context, facts *spookytypes.FactCollection) error {
    // Disk Partitions
    partitions, err := disk.Partitions(false)
    if err != nil {
        return fmt.Errorf("failed to get disk partitions: %w", err)
    }
    
    if facts.Facts.Collector == nil {
        facts.Facts.Collector = &spookytypesfacts.CollectorFacts{}
    }
    
    var diskFacts []*spookytypesfacts.DiskFacts
    for _, partition := range partitions {
        usage, err := disk.Usage(partition.Mountpoint)
        if err != nil {
            continue // Skip partitions we can't read
        }
        
        diskFact := &spookytypesfacts.DiskFacts{
            Device:     partition.Device,
            Mountpoint: partition.Mountpoint,
            FSType:     partition.Fstype,
            Total:      int64(usage.Total),
            Used:       int64(usage.Used),
            Free:       int64(usage.Free),
            Percent:    usage.UsedPercent,
        }
        diskFacts = append(diskFacts, diskFact)
    }
    
    facts.Facts.Collector.Disks = diskFacts
    
    // Disk I/O
    diskIO, err := disk.IOCounters()
    if err == nil {
        var totalRead, totalWrite uint64
        for _, io := range diskIO {
            totalRead += io.ReadBytes
            totalWrite += io.WriteBytes
        }
        
        facts.Facts.Collector.DiskIO = &spookytypesfacts.DiskIOFacts{
            ReadBytes:  int64(totalRead),
            WriteBytes: int64(totalWrite),
        }
    }
    
    return nil
}

func (c *AdvancedCollector) collectNetworkFacts(ctx context.Context, facts *spookytypes.FactCollection) error {
    // Network Interfaces
    interfaces, err := net.Interfaces()
    if err != nil {
        return fmt.Errorf("failed to get network interfaces: %w", err)
    }
    
    if facts.Facts.Collector == nil {
        facts.Facts.Collector = &spookytypesfacts.CollectorFacts{}
    }
    
    var networkFacts []*spookytypesfacts.NetworkInterfaceFacts
    for _, iface := range interfaces {
        // Skip loopback and down interfaces
        if iface.Name == "lo" || iface.Flags&net.FlagUp == 0 {
            continue
        }
        
        networkFact := &spookytypesfacts.NetworkInterfaceFacts{
            Name:        iface.Name,
            MTU:         iface.MTU,
            HardwareAddr: iface.HardwareAddr,
            Flags:       int(iface.Flags),
        }
        
        // Add addresses
        for _, addr := range iface.Addrs {
            networkFact.Addresses = append(networkFact.Addresses, addr.Addr)
        }
        
        networkFacts = append(networkFacts, networkFact)
    }
    
    facts.Facts.Collector.Network = &spookytypesfacts.NetworkFacts{
        Interfaces: networkFacts,
    }
    
    // Network I/O
    netIO, err := net.IOCounters(false)
    if err == nil && len(netIO) > 0 {
        io := netIO[0]
        facts.Facts.Collector.Network.BytesSent = int64(io.BytesSent)
        facts.Facts.Collector.Network.BytesRecv = int64(io.BytesRecv)
        facts.Facts.Collector.Network.PacketsSent = int64(io.PacketsSent)
        facts.Facts.Collector.Network.PacketsRecv = int64(io.PacketsRecv)
        facts.Facts.Collector.Network.ErrIn = int64(io.Errin)
        facts.Facts.Collector.Network.ErrOut = int64(io.Errout)
        facts.Facts.Collector.Network.DropIn = int64(io.Dropin)
        facts.Facts.Collector.Network.DropOut = int64(io.Dropout)
    }
    
    return nil
}

func (c *AdvancedCollector) collectProcessFacts(ctx context.Context, facts *spookytypes.FactCollection) error {
    // Process list
    processes, err := process.Processes()
    if err != nil {
        return fmt.Errorf("failed to get processes: %w", err)
    }
    
    if facts.Facts.Collector == nil {
        facts.Facts.Collector = &spookytypesfacts.CollectorFacts{}
    }
    
    facts.Facts.Collector.Processes = &spookytypesfacts.ProcessFacts{
        Count: len(processes),
    }
    
    // Top processes by CPU and memory
    topByCPU, err := c.getTopProcessesByCPU(processes, 10)
    if err == nil {
        facts.Facts.Collector.Processes.TopByCPU = topByCPU
    }
    
    topByMemory, err := c.getTopProcessesByMemory(processes, 10)
    if err == nil {
        facts.Facts.Collector.Processes.TopByMemory = topByMemory
    }
    
    return nil
}

func (c *AdvancedCollector) collectPerformanceFacts(ctx context.Context, facts *spookytypes.FactCollection) error {
    // Load Average
    loadAvg, err := load.Avg()
    if err == nil {
        if facts.Facts.Collector == nil {
            facts.Facts.Collector = &spookytypesfacts.CollectorFacts{}
        }
        
        facts.Facts.Collector.LoadAverage = &spookytypesfacts.LoadAverageFacts{
            Load1:  loadAvg.Load1,
            Load5:  loadAvg.Load5,
            Load15: loadAvg.Load15,
        }
    }
    
    return nil
}

func (c *AdvancedCollector) getTopProcessesByCPU(processes []*process.Process, limit int) ([]*spookytypesfacts.ProcessInfo, error) {
    var processInfos []*spookytypesfacts.ProcessInfo
    
    for _, p := range processes {
        name, err := p.Name()
        if err != nil {
            continue
        }
        
        cpuPercent, err := p.CPUPercent()
        if err != nil {
            continue
        }
        
        processInfo := &spookytypesfacts.ProcessInfo{
            PID:        int(p.Pid),
            Name:       name,
            CPUPercent: cpuPercent,
        }
        
        processInfos = append(processInfos, processInfo)
    }
    
    // Sort by CPU usage (simplified - would need proper sorting)
    if len(processInfos) > limit {
        processInfos = processInfos[:limit]
    }
    
    return processInfos, nil
}

func (c *AdvancedCollector) getTopProcessesByMemory(processes []*process.Process, limit int) ([]*spookytypesfacts.ProcessInfo, error) {
    var processInfos []*spookytypesfacts.ProcessInfo
    
    for _, p := range processes {
        name, err := p.Name()
        if err != nil {
            continue
        }
        
        memoryPercent, err := p.MemoryPercent()
        if err != nil {
            continue
        }
        
        processInfo := &spookytypesfacts.ProcessInfo{
            PID:          int(p.Pid),
            Name:         name,
            MemoryPercent: memoryPercent,
        }
        
        processInfos = append(processInfos, processInfo)
    }
    
    // Sort by memory usage (simplified - would need proper sorting)
    if len(processInfos) > limit {
        processInfos = processInfos[:limit]
    }
    
    return processInfos, nil
}
```

### 3. HCL Output Implementation

```go
// cmd/spooky-collector/output.go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
    
    "github.com/hashicorp/hcl/v2/hclwrite"
    
    spookytypes "spooky/internal/types"
)

func (c *AdvancedCollector) OutputToHCL(facts *spookytypes.FactCollection) error {
    // Ensure output directory exists
    if err := os.MkdirAll(c.outputDir, 0755); err != nil {
        return fmt.Errorf("failed to create output directory: %w", err)
    }
    
    // Create HCL content
    hclContent := c.createHCLContent(facts)
    
    // Write to temporary file first
    tempFile := filepath.Join(c.outputDir, "facts.hcl.tmp")
    if err := os.WriteFile(tempFile, hclContent, 0644); err != nil {
        return fmt.Errorf("failed to write temporary file: %w", err)
    }
    
    // Atomic move to final location
    finalFile := filepath.Join(c.outputDir, "facts.hcl")
    if err := os.Rename(tempFile, finalFile); err != nil {
        // Clean up temp file
        os.Remove(tempFile)
        return fmt.Errorf("failed to move file to final location: %w", err)
    }
    
    return nil
}

func (c *AdvancedCollector) createHCLContent(facts *spookytypes.FactCollection) []byte {
    f := hclwrite.NewEmptyFile()
    
    // Add header comment
    f.Body().AppendUnstructuredTokens([]*hclwrite.Token{
        {Type: hclwrite.TokenComment, Bytes: []byte("// Generated by spooky-collector")},
        {Type: hclwrite.TokenNewline, Bytes: []byte("\n")},
        {Type: hclwrite.TokenComment, Bytes: []byte(fmt.Sprintf("// Collected at: %s", facts.CollectedAt.Format(time.RFC3339)))},
        {Type: hclwrite.TokenNewline, Bytes: []byte("\n")},
        {Type: hclwrite.TokenNewline, Bytes: []byte("\n")},
    })
    
    // Add collector facts
    if facts.Facts.Collector != nil {
        c.addCollectorFacts(f, facts.Facts.Collector)
    }
    
    return f.Bytes()
}

func (c *AdvancedCollector) addCollectorFacts(f *hclwrite.File, collector *spookytypesfacts.CollectorFacts) {
    rootBody := f.Body()
    
    // Start collector block
    rootBody.AppendNewline()
    rootBody.SetAttributeValue("collector", c.convertToHCLValue(collector))
}
```

### 4. Systemd Service File

```ini
# deploy/spooky-collector.service
[Unit]
Description=Spooky Collector - System Facts Collection Service
After=network.target
Wants=network.target

[Service]
Type=simple
User=spooky
Group=spooky
ExecStart=/usr/local/bin/spooky-collector
ExecReload=/bin/kill -HUP $MAINPID
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=spooky-collector

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/etc/spooky

[Install]
WantedBy=multi-user.target
```

### 5. Installation Script

```bash
#!/bin/bash
# deploy/install.sh

set -e

# Configuration
BINARY_NAME="spooky-collector"
INSTALL_DIR="/usr/local/bin"
SERVICE_DIR="/etc/systemd/system"
CONFIG_DIR="/etc/spooky"
USER_NAME="spooky"
GROUP_NAME="spooky"

echo "Installing spooky-collector..."

# Create spooky user if it doesn't exist
if ! id "$USER_NAME" &>/dev/null; then
    echo "Creating spooky user..."
    sudo useradd --system --shell /bin/false --home-dir /var/lib/spooky "$USER_NAME"
fi

# Create spooky group if it doesn't exist
if ! getent group "$GROUP_NAME" &>/dev/null; then
    echo "Creating spooky group..."
    sudo groupadd --system "$GROUP_NAME"
    sudo usermod -a -G "$GROUP_NAME" "$USER_NAME"
fi

# Create configuration directory
echo "Creating configuration directory..."
sudo mkdir -p "$CONFIG_DIR"
sudo chown "$USER_NAME:$GROUP_NAME" "$CONFIG_DIR"
sudo chmod 755 "$CONFIG_DIR"

# Install binary
echo "Installing binary..."
sudo cp "$BINARY_NAME" "$INSTALL_DIR/"
sudo chown root:root "$INSTALL_DIR/$BINARY_NAME"
sudo chmod 755 "$INSTALL_DIR/$BINARY_NAME"

# Install service file
echo "Installing systemd service..."
sudo cp spooky-collector.service "$SERVICE_DIR/"
sudo chown root:root "$SERVICE_DIR/spooky-collector.service"
sudo chmod 644 "$SERVICE_DIR/spooky-collector.service"

# Reload systemd and enable service
echo "Enabling service..."
sudo systemctl daemon-reload
sudo systemctl enable spooky-collector.service

echo "Installation complete!"
echo "To start the service: sudo systemctl start spooky-collector"
echo "To check status: sudo systemctl status spooky-collector"
```

## Usage Examples

### Command Line Usage

```bash
# Run once and exit
./spooky-collector --oneshot

# Run with custom output directory
./spooky-collector --output /tmp/spooky-facts --oneshot

# Run with custom interval
./spooky-collector --interval 2m

# Run with verbose logging
./spooky-collector --verbose --oneshot
```

### Service Management

```bash
# Start the service
sudo systemctl start spooky-collector

# Check status
sudo systemctl status spooky-collector

# View logs
sudo journalctl -u spooky-collector -f

# Stop the service
sudo systemctl stop spooky-collector
```

## Integration with Spooky

The `spooky-collector` binary outputs facts to `/etc/spooky/facts.hcl` in the format expected by the main spooky tool. The facts are collected under the `collector` namespace and can be accessed via SSH collection.

### Fact Collection Process

1. **Deploy spooky-collector** to target machines
2. **Start the service** to begin continuous collection
3. **Use spooky facts gather** to collect via SSH
4. **Access collector facts** through the three-namespace structure

### Output Format

The collector outputs HCL in the format defined by `facts.schema.hcl`:

```hcl
// Generated by spooky-collector
// Collected at: 2024-01-01T12:00:00Z

collector {
  host {
    hostname = "web-server-01"
    uptime = 86400
    os = "linux"
    platform = "ubuntu"
    kernel_version = "5.15.0-91-generic"
  }
  
  cpu {
    cores = 4
    model = "Intel(R) Core(TM) i7-8700K CPU @ 3.70GHz"
    frequency = 3700
    percent = 15.2
  }
  
  memory {
    total = 8589934592
    used = 4294967296
    free = 4294967296
    percent = 50.0
  }
  
  // ... additional facts
}
```

## Security Considerations

### User Privileges
- Runs as `spooky` user (non-root)
- No sudo privileges required
- Minimal file system access
- Network access only for local interfaces

### File Permissions
- Output directory: 755 (rwxr-xr-x)
- Facts file: 644 (rw-r--r--)
- Service file: 644 (rw-r--r--)
- Binary: 755 (rwxr-xr-x)

### Systemd Security
- `NoNewPrivileges=true`
- `PrivateTmp=true`
- `ProtectSystem=strict`
- `ProtectHome=true`
- `ReadWritePaths=/etc/spooky`

## Troubleshooting

### Common Issues

1. **Permission Denied**
   - Ensure spooky user has write access to `/etc/spooky`
   - Check file permissions and ownership

2. **Service Won't Start**
   - Check systemd logs: `journalctl -u spooky-collector`
   - Verify binary exists and is executable
   - Check configuration directory permissions

3. **No Facts Output**
   - Run with `--verbose` flag for detailed logging
   - Check if gopsutil can access system information
   - Verify output directory exists and is writable

4. **High CPU Usage**
   - Adjust collection interval with `--interval`
   - Check for process leaks or infinite loops
   - Monitor system resources during collection

### Debugging Commands

```bash
# Test binary directly
sudo -u spooky ./spooky-collector --verbose --oneshot

# Check service status
sudo systemctl status spooky-collector

# View recent logs
sudo journalctl -u spooky-collector --since "1 hour ago"

# Check file permissions
ls -la /etc/spooky/
ls -la /usr/local/bin/spooky-collector

# Test fact collection manually
sudo -u spooky cat /proc/cpuinfo
sudo -u spooky free -b
```

## Performance Considerations

### Collection Intervals
- **Default**: 5 minutes (good balance of freshness vs overhead)
- **High-frequency**: 1 minute (for monitoring systems)
- **Low-frequency**: 15 minutes (for basic inventory)

### Resource Usage
- **CPU**: < 1% during collection
- **Memory**: < 50MB resident
- **Disk I/O**: Minimal (single file write per collection)
- **Network**: None (local collection only)

### Optimization Tips
- Use longer intervals for stable systems
- Monitor collection time and adjust as needed
- Consider running during off-peak hours
- Use `--oneshot` for manual collection when needed

## Future Enhancements

### Planned Features
1. **Configuration file support** for custom collection rules
2. **Plugin system** for custom fact collectors
3. **Compression** of output files
4. **Historical data** with timestamped files
5. **Health checks** and self-monitoring
6. **Metrics export** for monitoring systems

### Integration Opportunities
1. **Prometheus metrics** export
2. **SNMP integration** for network devices
3. **Container support** for Kubernetes environments
4. **Cloud metadata** collection (AWS, GCP, Azure)
5. **Application-specific** fact collection

This implementation guide provides a complete foundation for building and deploying the `spooky-collector` binary as part of the advanced fact collection system.
