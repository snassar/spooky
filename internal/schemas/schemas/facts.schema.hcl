# Facts Schema
# Simple three-namespace fact collection structure
# Used for SSH-based fact collection and custom facts

metadata {
  schema_version = "0.20250809.0"
  schema_type = "facts"
  schema_name = "Facts Schema"
  last_updated = "2024-01-01"
  compatibility = ["0.20250809.0"]
  description = "Simple three-namespace fact collection structure"
  scalver_format = "0.20250809.0"
}

facts {
  # System facts (from SSH commands - user level)
  system = {
    type = "object"
    required = true
    description = "System facts collected via SSH commands"
    
    properties = {
      # OS Information
      os_name = {
        type = "string"
        required = false
        description = "Operating system name"
      }
      
      os_version = {
        type = "string"
        required = false
        description = "Operating system version"
      }
      
      kernel = {
        type = "string"
        required = false
        description = "Kernel version"
      }
      
      hostname = {
        type = "string"
        required = false
        description = "System hostname"
      }
      
      # Hardware Information
      cpu_cores = {
        type = "integer"
        required = false
        description = "Number of CPU cores"
      }
      
      memory_total = {
        type = "integer"
        required = false
        description = "Total memory in bytes"
      }
      
      disk_total = {
        type = "integer"
        required = false
        description = "Total disk space in bytes"
      }
      
      # Network Information
      primary_ip = {
        type = "string"
        required = false
        description = "Primary IP address"
      }
      
      interfaces = {
        type = "array"
        required = false
        description = "Network interface names"
        items = {
          type = "string"
        }
      }
      
      # System State
      uptime = {
        type = "integer"
        required = false
        description = "System uptime in seconds"
      }
      
      load_average = {
        type = "array"
        required = false
        description = "Load averages (1, 5, 15 minute)"
        items = {
          type = "number"
        }
      }
      
      process_count = {
        type = "integer"
        required = false
        description = "Number of running processes"
      }
    }
  }
  
  # Collector facts (from spooky-collector binary - comprehensive gopsutil coverage)
  collector = {
    type = "object"
    required = false
    description = "Comprehensive system facts collected via gopsutil"
    
    properties = {
      # Host Information (gopsutil/host)
      host = {
        type = "object"
        required = false
        description = "Host information from gopsutil"
        
        properties = {
          hostname = {
            type = "string"
            required = false
            description = "System hostname"
          }
          
          uptime = {
            type = "integer"
            required = false
            description = "System uptime in seconds"
          }
          
          boot_time = {
            type = "integer"
            required = false
            description = "System boot time (Unix timestamp)"
          }
          
          procs = {
            type = "integer"
            required = false
            description = "Number of processes"
          }
          
          os = {
            type = "object"
            required = false
            description = "Operating system information"
            
            properties = {
              family = {
                type = "string"
                required = false
                description = "OS family (e.g., debian, rhel, suse)"
              }
              
              platform = {
                type = "string"
                required = false
                description = "OS platform"
              }
              
              version = {
                type = "string"
                required = false
                description = "OS version"
              }
              
              kernel_version = {
                type = "string"
                required = false
                description = "Kernel version"
              }
              
              kernel_arch = {
                type = "string"
                required = false
                description = "Kernel architecture"
              }
            }
          }
          
          platform = {
            type = "object"
            required = false
            description = "Platform information"
            
            properties = {
              system = {
                type = "string"
                required = false
                description = "Platform system"
              }
              
              family = {
                type = "string"
                required = false
                description = "Platform family"
              }
              
              version = {
                type = "string"
                required = false
                description = "Platform version"
              }
            }
          }
        }
      }
      
      # CPU Information (gopsutil/cpu)
      cpu = {
        type = "object"
        required = false
        description = "CPU information from gopsutil"
        
        properties = {
          count = {
            type = "integer"
            required = false
            description = "Number of CPU cores"
          }
          
          count_logical = {
            type = "integer"
            required = false
            description = "Number of logical CPU cores"
          }
          
          percent = {
            type = "number"
            required = false
            description = "Overall CPU usage percentage"
          }
          
          times = {
            type = "object"
            required = false
            description = "CPU time breakdown"
            
            properties = {
              user = {
                type = "number"
                required = false
                description = "User CPU time"
              }
              
              system = {
                type = "number"
                required = false
                description = "System CPU time"
              }
              
              idle = {
                type = "number"
                required = false
                description = "Idle CPU time"
              }
              
              nice = {
                type = "number"
                required = false
                description = "Nice CPU time"
              }
              
              iowait = {
                type = "number"
                required = false
                description = "I/O wait CPU time"
              }
              
              irq = {
                type = "number"
                required = false
                description = "IRQ CPU time"
              }
              
              softirq = {
                type = "number"
                required = false
                description = "Soft IRQ CPU time"
              }
              
              steal = {
                type = "number"
                required = false
                description = "Steal CPU time"
              }
              
              guest = {
                type = "number"
                required = false
                description = "Guest CPU time"
              }
              
              guest_nice = {
                type = "number"
                required = false
                description = "Guest nice CPU time"
              }
            }
          }
          
          info = {
            type = "array"
            required = false
            description = "Detailed CPU information for each core"
            items = {
              type = "object"
              properties = {
                cpu = {
                  type = "integer"
                  required = false
                  description = "CPU core number"
                }
                
                vendor_id = {
                  type = "string"
                  required = false
                  description = "CPU vendor ID"
                }
                
                family = {
                  type = "string"
                  required = false
                  description = "CPU family"
                }
                
                model = {
                  type = "string"
                  required = false
                  description = "CPU model"
                }
                
                stepping = {
                  type = "integer"
                  required = false
                  description = "CPU stepping"
                }
                
                physical_id = {
                  type = "string"
                  required = false
                  description = "Physical CPU ID"
                }
                
                core_id = {
                  type = "string"
                  required = false
                  description = "CPU core ID"
                }
                
                cores = {
                  type = "integer"
                  required = false
                  description = "Number of cores"
                }
                
                model_name = {
                  type = "string"
                  required = false
                  description = "CPU model name"
                }
                
                mhz = {
                  type = "number"
                  required = false
                  description = "CPU frequency in MHz"
                }
                
                cache_size = {
                  type = "integer"
                  required = false
                  description = "CPU cache size in bytes"
                }
                
                flags = {
                  type = "array"
                  required = false
                  description = "CPU flags"
                  items = {
                    type = "string"
                  }
                }
              }
            }
          }
        }
      }
      
      # Memory Information (gopsutil/mem)
      memory = {
        type = "object"
        required = false
        description = "Memory information from gopsutil"
        
        properties = {
          total = {
            type = "integer"
            required = false
            description = "Total memory in bytes"
          }
          
          available = {
            type = "integer"
            required = false
            description = "Available memory in bytes"
          }
          
          used = {
            type = "integer"
            required = false
            description = "Used memory in bytes"
          }
          
          free = {
            type = "integer"
            required = false
            description = "Free memory in bytes"
          }
          
          active = {
            type = "integer"
            required = false
            description = "Active memory in bytes"
          }
          
          inactive = {
            type = "integer"
            required = false
            description = "Inactive memory in bytes"
          }
          
          wired = {
            type = "integer"
            required = false
            description = "Wired memory in bytes"
          }
          
          launder = {
            type = "integer"
            required = false
            description = "Launder memory in bytes"
          }
          
          buffers = {
            type = "integer"
            required = false
            description = "Memory used by buffers in bytes"
          }
          
          cached = {
            type = "integer"
            required = false
            description = "Memory used by cache in bytes"
          }
          
          writeback = {
            type = "integer"
            required = false
            description = "Writeback memory in bytes"
          }
          
          dirty = {
            type = "integer"
            required = false
            description = "Dirty memory in bytes"
          }
          
          writeback_tmp = {
            type = "integer"
            required = false
            description = "Writeback temporary memory in bytes"
          }
          
          shared = {
            type = "integer"
            required = false
            description = "Shared memory in bytes"
          }
          
          slab = {
            type = "integer"
            required = false
            description = "Slab memory in bytes"
          }
          
          sreclaimable = {
            type = "integer"
            required = false
            description = "SReclaimable memory in bytes"
          }
          
          sunreclaim = {
            type = "integer"
            required = false
            description = "SUnreclaim memory in bytes"
          }
          
          page_tables = {
            type = "integer"
            required = false
            description = "Page tables memory in bytes"
          }
          
          swap_cached = {
            type = "integer"
            required = false
            description = "Swap cached memory in bytes"
          }
          
          commit_limit = {
            type = "integer"
            required = false
            description = "Commit limit in bytes"
          }
          
          committed_as = {
            type = "integer"
            required = false
            description = "Committed AS memory in bytes"
          }
          
          high_total = {
            type = "integer"
            required = false
            description = "High total memory in bytes"
          }
          
          high_free = {
            type = "integer"
            required = false
            description = "High free memory in bytes"
          }
          
          low_total = {
            type = "integer"
            required = false
            description = "Low total memory in bytes"
          }
          
          low_free = {
            type = "integer"
            required = false
            description = "Low free memory in bytes"
          }
          
          swap_total = {
            type = "integer"
            required = false
            description = "Total swap space in bytes"
          }
          
          swap_free = {
            type = "integer"
            required = false
            description = "Free swap space in bytes"
          }
          
          mapped = {
            type = "integer"
            required = false
            description = "Mapped memory in bytes"
          }
          
          vmalloc_total = {
            type = "integer"
            required = false
            description = "Total vmalloc space in bytes"
          }
          
          vmalloc_used = {
            type = "integer"
            required = false
            description = "Used vmalloc space in bytes"
          }
          
          vmalloc_chunk = {
            type = "integer"
            required = false
            description = "Vmalloc chunk size in bytes"
          }
          
          huge_pages_total = {
            type = "integer"
            required = false
            description = "Total huge pages"
          }
          
          huge_pages_free = {
            type = "integer"
            required = false
            description = "Free huge pages"
          }
          
          huge_page_size = {
            type = "integer"
            required = false
            description = "Huge page size in bytes"
          }
        }
      }
      
      # Virtual Memory Information (gopsutil/mem)
      virtual_memory = {
        type = "object"
        required = false
        description = "Virtual memory information from gopsutil"
        
        properties = {
          total = {
            type = "integer"
            required = false
            description = "Total virtual memory in bytes"
          }
          
          available = {
            type = "integer"
            required = false
            description = "Available virtual memory in bytes"
          }
          
          used = {
            type = "integer"
            required = false
            description = "Used virtual memory in bytes"
          }
          
          free = {
            type = "integer"
            required = false
            description = "Free virtual memory in bytes"
          }
          
          active = {
            type = "integer"
            required = false
            description = "Active virtual memory in bytes"
          }
          
          inactive = {
            type = "integer"
            required = false
            description = "Inactive virtual memory in bytes"
          }
          
          wired = {
            type = "integer"
            required = false
            description = "Wired virtual memory in bytes"
          }
          
          launder = {
            type = "integer"
            required = false
            description = "Launder virtual memory in bytes"
          }
          
          buffers = {
            type = "integer"
            required = false
            description = "Virtual memory used by buffers in bytes"
          }
          
          cached = {
            type = "integer"
            required = false
            description = "Virtual memory used by cache in bytes"
          }
          
          writeback = {
            type = "integer"
            required = false
            description = "Writeback virtual memory in bytes"
          }
          
          dirty = {
            type = "integer"
            required = false
            description = "Dirty virtual memory in bytes"
          }
          
          writeback_tmp = {
            type = "integer"
            required = false
            description = "Writeback temporary virtual memory in bytes"
          }
          
          shared = {
            type = "integer"
            required = false
            description = "Shared virtual memory in bytes"
          }
          
          slab = {
            type = "integer"
            required = false
            description = "Slab virtual memory in bytes"
          }
          
          sreclaimable = {
            type = "integer"
            required = false
            description = "SReclaimable virtual memory in bytes"
          }
          
          sunreclaim = {
            type = "integer"
            required = false
            description = "SUnreclaim virtual memory in bytes"
          }
          
          page_tables = {
            type = "integer"
            required = false
            description = "Page tables virtual memory in bytes"
          }
          
          swap_cached = {
            type = "integer"
            required = false
            description = "Swap cached virtual memory in bytes"
          }
          
          commit_limit = {
            type = "integer"
            required = false
            description = "Commit limit in bytes"
          }
          
          committed_as = {
            type = "integer"
            required = false
            description = "Committed AS virtual memory in bytes"
          }
          
          high_total = {
            type = "integer"
            required = false
            description = "High total virtual memory in bytes"
          }
          
          high_free = {
            type = "integer"
            required = false
            description = "High free virtual memory in bytes"
          }
          
          low_total = {
            type = "integer"
            required = false
            description = "Low total virtual memory in bytes"
          }
          
          low_free = {
            type = "integer"
            required = false
            description = "Low free virtual memory in bytes"
          }
          
          swap_total = {
            type = "integer"
            required = false
            description = "Total swap space in bytes"
          }
          
          swap_free = {
            type = "integer"
            required = false
            description = "Free swap space in bytes"
          }
          
          mapped = {
            type = "integer"
            required = false
            description = "Mapped virtual memory in bytes"
          }
          
          vmalloc_total = {
            type = "integer"
            required = false
            description = "Total vmalloc space in bytes"
          }
          
          vmalloc_used = {
            type = "integer"
            required = false
            description = "Used vmalloc space in bytes"
          }
          
          vmalloc_chunk = {
            type = "integer"
            required = false
            description = "Vmalloc chunk size in bytes"
          }
          
          huge_pages_total = {
            type = "integer"
            required = false
            description = "Total huge pages"
          }
          
          huge_pages_free = {
            type = "integer"
            required = false
            description = "Free huge pages"
          }
          
          huge_page_size = {
            type = "integer"
            required = false
            description = "Huge page size in bytes"
          }
        }
      }
      
      # Disk Information (gopsutil/disk)
      disk = {
        type = "object"
        required = false
        description = "Disk information from gopsutil"
        
        properties = {
          partitions = {
            type = "array"
            required = false
            description = "Disk partition information"
            items = {
              type = "object"
              properties = {
                device = {
                  type = "string"
                  required = false
                  description = "Device name"
                }
                
                mountpoint = {
                  type = "string"
                  required = false
                  description = "Mount point"
                }
                
                fstype = {
                  type = "string"
                  required = false
                  description = "Filesystem type"
                }
                
                opts = {
                  type = "string"
                  required = false
                  description = "Mount options"
                }
              }
            }
          }
          
          usage = {
            type = "array"
            required = false
            description = "Disk usage information"
            items = {
              type = "object"
              properties = {
                path = {
                  type = "string"
                  required = false
                  description = "Path"
                }
                
                fstype = {
                  type = "string"
                  required = false
                  description = "Filesystem type"
                }
                
                total = {
                  type = "integer"
                  required = false
                  description = "Total space in bytes"
                }
                
                free = {
                  type = "integer"
                  required = false
                  description = "Free space in bytes"
                }
                
                used = {
                  type = "integer"
                  required = false
                  description = "Used space in bytes"
                }
                
                inodes_total = {
                  type = "integer"
                  required = false
                  description = "Total inodes"
                }
                
                inodes_used = {
                  type = "integer"
                  required = false
                  description = "Used inodes"
                }
                
                inodes_free = {
                  type = "integer"
                  required = false
                  description = "Free inodes"
                }
              }
            }
          }
          
          io_counters = {
            type = "array"
            required = false
            description = "Disk I/O counters"
            items = {
              type = "object"
              properties = {
                name = {
                  type = "string"
                  required = false
                  description = "Device name"
                }
                
                read_count = {
                  type = "integer"
                  required = false
                  description = "Number of read operations"
                }
                
                write_count = {
                  type = "integer"
                  required = false
                  description = "Number of write operations"
                }
                
                read_bytes = {
                  type = "integer"
                  required = false
                  description = "Total bytes read"
                }
                
                write_bytes = {
                  type = "integer"
                  required = false
                  description = "Total bytes written"
                }
                
                read_time = {
                  type = "integer"
                  required = false
                  description = "Time spent reading in milliseconds"
                }
                
                write_time = {
                  type = "integer"
                  required = false
                  description = "Time spent writing in milliseconds"
                }
                
                io_time = {
                  type = "integer"
                  required = false
                  description = "Time spent doing I/O in milliseconds"
                }
                
                weighted_io = {
                  type = "integer"
                  required = false
                  description = "Weighted time spent doing I/O"
                }
              }
            }
          }
        }
      }
      
      # Network Information (gopsutil/net)
      network = {
        type = "object"
        required = false
        description = "Network information from gopsutil"
        
        properties = {
          interfaces = {
            type = "array"
            required = false
            description = "Network interface information"
            items = {
              type = "object"
              properties = {
                name = {
                  type = "string"
                  required = false
                  description = "Interface name"
                }
                
                mtu = {
                  type = "integer"
                  required = false
                  description = "Interface MTU"
                }
                
                flags = {
                  type = "array"
                  required = false
                  description = "Interface flags"
                  items = {
                    type = "string"
                  }
                }
                
                addrs = {
                  type = "array"
                  required = false
                  description = "Interface addresses"
                  items = {
                    type = "object"
                    properties = {
                      addr = {
                        type = "string"
                        required = false
                        description = "IP address"
                      }
                      
                      netmask = {
                        type = "string"
                        required = false
                        description = "Netmask"
                      }
                    }
                  }
                }
                
                hardware_addr = {
                  type = "string"
                  required = false
                  description = "Hardware address (MAC)"
                }
              }
            }
          }
          
          io_counters = {
            type = "array"
            required = false
            description = "Network I/O counters"
            items = {
              type = "object"
              properties = {
                name = {
                  type = "string"
                  required = false
                  description = "Interface name"
                }
                
                bytes_sent = {
                  type = "integer"
                  required = false
                  description = "Bytes sent"
                }
                
                bytes_recv = {
                  type = "integer"
                  required = false
                  description = "Bytes received"
                }
                
                packets_sent = {
                  type = "integer"
                  required = false
                  description = "Packets sent"
                }
                
                packets_recv = {
                  type = "integer"
                  required = false
                  description = "Packets received"
                }
                
                err_in = {
                  type = "integer"
                  required = false
                  description = "Input errors"
                }
                
                err_out = {
                  type = "integer"
                  required = false
                  description = "Output errors"
                }
                
                drop_in = {
                  type = "integer"
                  required = false
                  description = "Input drops"
                }
                
                drop_out = {
                  type = "integer"
                  required = false
                  description = "Output drops"
                }
                
                fifo_in = {
                  type = "integer"
                  required = false
                  description = "FIFO buffer errors (input)"
                }
                
                fifo_out = {
                  type = "integer"
                  required = false
                  description = "FIFO buffer errors (output)"
                }
                
                frame_in = {
                  type = "integer"
                  required = false
                  description = "Frame errors (input)"
                }
                
                frame_out = {
                  type = "integer"
                  required = false
                  description = "Frame errors (output)"
                }
                
                compressed_in = {
                  type = "integer"
                  required = false
                  description = "Compressed packets (input)"
                }
                
                compressed_out = {
                  type = "integer"
                  required = false
                  description = "Compressed packets (output)"
                }
                
                multicast_in = {
                  type = "integer"
                  required = false
                  description = "Multicast packets (input)"
                }
                
                multicast_out = {
                  type = "integer"
                  required = false
                  description = "Multicast packets (output)"
                }
              }
            }
          }
          
          connections = {
            type = "array"
            required = false
            description = "Network connections"
            items = {
              type = "object"
              properties = {
                fd = {
                  type = "integer"
                  required = false
                  description = "File descriptor"
                }
                
                family = {
                  type = "integer"
                  required = false
                  description = "Address family"
                }
                
                type = {
                  type = "integer"
                  required = false
                  description = "Socket type"
                }
                
                laddr = {
                  type = "object"
                  required = false
                  description = "Local address"
                  properties = {
                    ip = {
                      type = "string"
                      required = false
                      description = "Local IP address"
                    }
                    
                    port = {
                      type = "integer"
                      required = false
                      description = "Local port"
                    }
                  }
                }
                
                raddr = {
                  type = "object"
                  required = false
                  description = "Remote address"
                  properties = {
                    ip = {
                      type = "string"
                      required = false
                      description = "Remote IP address"
                    }
                    
                    port = {
                      type = "integer"
                      required = false
                      description = "Remote port"
                    }
                  }
                }
                
                status = {
                  type = "string"
                  required = false
                  description = "Connection status"
                }
                
                pid = {
                  type = "integer"
                  required = false
                  description = "Process ID"
                }
              }
            }
          }
        }
      }
      
      # Load Average Information (gopsutil/load)
      load = {
        type = "object"
        required = false
        description = "Load average information from gopsutil"
        
        properties = {
          load1 = {
            type = "number"
            required = false
            description = "1-minute load average"
          }
          
          load5 = {
            type = "number"
            required = false
            description = "5-minute load average"
          }
          
          load15 = {
            type = "number"
            required = false
            description = "15-minute load average"
          }
        }
      }
      
      # Process Information (gopsutil/process)
      processes = {
        type = "object"
        required = false
        description = "Process information from gopsutil"
        
        properties = {
          count = {
            type = "integer"
            required = false
            description = "Total number of processes"
          }
          
          top_by_cpu = {
            type = "array"
            required = false
            description = "Top processes by CPU usage"
            items = {
              type = "object"
              properties = {
                pid = {
                  type = "integer"
                  required = false
                  description = "Process ID"
                }
                
                name = {
                  type = "string"
                  required = false
                  description = "Process name"
                }
                
                cpu_percent = {
                  type = "number"
                  required = false
                  description = "CPU usage percentage"
                }
                
                memory_percent = {
                  type = "number"
                  required = false
                  description = "Memory usage percentage"
                }
                
                cmdline = {
                  type = "string"
                  required = false
                  description = "Command line"
                }
              }
            }
          }
          
          top_by_memory = {
            type = "array"
            required = false
            description = "Top processes by memory usage"
            items = {
              type = "object"
              properties = {
                pid = {
                  type = "integer"
                  required = false
                  description = "Process ID"
                }
                
                name = {
                  type = "string"
                  required = false
                  description = "Process name"
                }
                
                memory_percent = {
                  type = "number"
                  required = false
                  description = "Memory usage percentage"
                }
                
                memory_rss = {
                  type = "integer"
                  required = false
                  description = "RSS memory in bytes"
                }
                
                memory_vms = {
                  type = "integer"
                  required = false
                  description = "VMS memory in bytes"
                }
              }
            }
          }
        }
      }
    }
  }
  
  # Custom facts (from /etc/spooky/custom.hcl)
  custom = {
    type = "object"
    required = false
    description = "Custom facts from /etc/spooky/custom.hcl"
    
    properties = {
      # Environment information
      environment = {
        type = "string"
        required = false
        description = "Environment name"
        enum = ["development", "staging", "production", "testing"]
      }
      
      # Application information
      application = {
        type = "object"
        required = false
        description = "Application-specific information"
      }
      
      # Infrastructure information
      infrastructure = {
        type = "object"
        required = false
        description = "Infrastructure-specific information"
      }
      
      # Business information
      business = {
        type = "object"
        required = false
        description = "Business-specific information"
      }
      
      # Monitoring information
      monitoring = {
        type = "object"
        required = false
        description = "Monitoring configuration"
      }
    }
  }
  
  # Metadata
  collected_at = {
    type = "string"
    required = true
    format = "date-time"
    description = "Timestamp when facts were collected (ISO 8601 format)"
  }
  
  machine_id = {
    type = "string"
    required = true
    pattern = "^[a-f0-9]{32}$"
    description = "32-character hexadecimal machine identifier"
  }
}

# Validation rules
validation {
  # Machine ID validation
  machine_id_format = {
    rule = "regex"
    pattern = "^[a-f0-9]{32}$"
    message = "Machine ID must be a 32-character hexadecimal string"
  }
  
  # Timestamp validation
  timestamp_format = {
    rule = "date_time"
    format = "RFC3339"
    message = "Timestamp must be in RFC3339 format"
  }
  
  # System facts required
  system_facts_required = {
    rule = "required"
    field = "system"
    message = "System facts are required"
  }
  
  # Collector facts optional
  collector_facts_optional = {
    rule = "optional"
    field = "collector"
    message = "Collector facts are optional"
  }
  
  # Custom facts optional
  custom_facts_optional = {
    rule = "optional"
    field = "custom"
    message = "Custom facts are optional"
  }
}
