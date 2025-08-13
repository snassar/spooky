# Facts Structure Schema
# Common fact structure definitions for all storage and export formats
# This file defines the structure of facts that can be stored and exported
# Used for validation of facts data regardless of format (memory, JSON, HCL)

# Common facts structure
facts_structure {
  # Machine ID (32-character hex string from /etc/machine-id)
  machine_id = {
    type = "string"
    required = true
    pattern = "^[a-f0-9]{32}$"
    description = "Unique machine identifier from /etc/machine-id"
  }

  # Collection timestamp
  collected_at = {
    type = "string"
    required = true
    format = "date-time"
    description = "Timestamp when facts were collected (ISO 8601 format)"
  }

  # Collection of facts
  facts = {
    type = "object"
    required = true
    description = "Collection of facts for this machine"
    
    properties = {
      # System-level facts (from gopsutil)
      system = {
        type = "object"
        required = true
        description = "System-level facts from gopsutil"
        
        properties = {
          # Operating system facts
          os = {
            type = "object"
            required = true
            description = "Operating system facts"
            
            properties = {
              name = {
                type = "string"
                required = true
                description = "OS name (e.g., Ubuntu, CentOS, Debian)"
              }
              version = {
                type = "string"
                required = true
                description = "OS version (e.g., 22.04 LTS, 8.5)"
              }
              arch = {
                type = "string"
                required = true
                description = "Architecture (e.g., x86_64, arm64)"
              }
              kernel = {
                type = "string"
                required = true
                description = "Kernel version"
              }
              platform = {
                type = "string"
                required = false
                description = "Platform identifier"
              }
              family = {
                type = "string"
                required = false
                description = "OS family (e.g., debian, rhel, suse)"
              }
            }
          }
          
          # Hardware facts
          hardware = {
            type = "object"
            required = true
            description = "Hardware facts"
            
            properties = {
              # CPU information
              cpu = {
                type = "object"
                required = true
                description = "CPU information"
                
                properties = {
                  cores = {
                    type = "integer"
                    required = true
                    description = "Number of CPU cores"
                    min = 1
                  }
                  model = {
                    type = "string"
                    required = true
                    description = "CPU model"
                  }
                  frequency = {
                    type = "number"
                    required = false
                    description = "CPU frequency in MHz"
                  }
                  architecture = {
                    type = "string"
                    required = false
                    description = "CPU architecture"
                  }
                  vendor = {
                    type = "string"
                    required = false
                    description = "CPU vendor"
                  }
                  # CPU times
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
                  # CPU percentages
                  percent = {
                    type = "number"
                    required = false
                    description = "Overall CPU usage percentage"
                  }
                  # Per-core information
                  cores_detail = {
                    type = "array"
                    required = false
                    description = "Detailed information for each CPU core"
                    items = {
                      type = "object"
                      properties = {
                        cpu = {
                          type = "integer"
                          required = true
                          description = "CPU core number"
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
                        percent = {
                          type = "number"
                          required = false
                          description = "CPU usage percentage"
                        }
                        times = {
                          type = "object"
                          required = false
                          description = "CPU time breakdown for this core"
                          
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
                      }
                    }
                  }
                }
              }
              
              # Memory information
              memory = {
                type = "object"
                required = true
                description = "Memory information"
                
                properties = {
                  total = {
                    type = "integer"
                    required = true
                    description = "Total memory in bytes"
                    min = 1
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
                  swap = {
                    type = "object"
                    required = false
                    description = "Swap memory information"
                    
                    properties = {
                      total = {
                        type = "integer"
                        required = false
                        description = "Total swap space in bytes"
                      }
                      used = {
                        type = "integer"
                        required = false
                        description = "Used swap space in bytes"
                      }
                      free = {
                        type = "integer"
                        required = false
                        description = "Free swap space in bytes"
                      }
                      percent = {
                        type = "number"
                        required = false
                        description = "Swap usage percentage"
                      }
                    }
                  }
                  virtual_memory = {
                    type = "object"
                    required = false
                    description = "Virtual memory information"
                    
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
                      percent = {
                        type = "number"
                        required = false
                        description = "Virtual memory usage percentage"
                      }
                    }
                  }
                }
              }
              
              # Disk information
              disks = {
                type = "array"
                required = true
                description = "List of disk information"
                items = {
                  type = "object"
                  properties = {
                    device = {
                      type = "string"
                      required = true
                      description = "Device name (e.g., /dev/sda1)"
                    }
                    mount_point = {
                      type = "string"
                      required = false
                      description = "Mount point (e.g., /)"
                    }
                    total = {
                      type = "integer"
                      required = true
                      description = "Total disk space in bytes"
                    }
                    used = {
                      type = "integer"
                      required = false
                      description = "Used disk space in bytes"
                    }
                    free = {
                      type = "integer"
                      required = false
                      description = "Free disk space in bytes"
                    }
                    filesystem = {
                      type = "string"
                      required = false
                      description = "Filesystem type (e.g., ext4, xfs)"
                    }
                    # Disk I/O statistics
                    io_counters = {
                      type = "object"
                      required = false
                      description = "Disk I/O statistics"
                      
                      properties = {
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
                    # Partition information
                    partition = {
                      type = "object"
                      required = false
                      description = "Partition information"
                      
                      properties = {
                        device = {
                          type = "string"
                          required = false
                          description = "Partition device name"
                        }
                        mountpoint = {
                          type = "string"
                          required = false
                          description = "Partition mount point"
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
                }
              }
              
              # Disk I/O summary
              disk_io = {
                type = "object"
                required = false
                description = "Overall disk I/O statistics"
                
                properties = {
                  read_count = {
                    type = "integer"
                    required = false
                    description = "Total read operations across all disks"
                  }
                  write_count = {
                    type = "integer"
                    required = false
                    description = "Total write operations across all disks"
                  }
                  read_bytes = {
                    type = "integer"
                    required = false
                    description = "Total bytes read across all disks"
                  }
                  write_bytes = {
                    type = "integer"
                    required = false
                    description = "Total bytes written across all disks"
                  }
                  read_time = {
                    type = "integer"
                    required = false
                    description = "Total time spent reading in milliseconds"
                  }
                  write_time = {
                    type = "integer"
                    required = false
                    description = "Total time spent writing in milliseconds"
                  }
                  io_time = {
                    type = "integer"
                    required = false
                    description = "Total time spent doing I/O in milliseconds"
                  }
                  weighted_io = {
                    type = "integer"
                    required = false
                    description = "Total weighted time spent doing I/O"
                  }
                }
              }
            }
          }
          
          # Network facts
          network = {
            type = "object"
            required = true
            description = "Network facts"
            
            properties = {
              hostname = {
                type = "string"
                required = true
                description = "System hostname"
              }
              interfaces = {
                type = "array"
                required = true
                description = "List of network interfaces"
                items = {
                  type = "object"
                  properties = {
                    name = {
                      type = "string"
                      required = true
                      description = "Interface name (e.g., eth0, ens3)"
                    }
                    mac_address = {
                      type = "string"
                      required = false
                      description = "MAC address"
                    }
                    ip_addresses = {
                      type = "array"
                      required = false
                      description = "List of IP addresses"
                      items = {
                        type = "string"
                      }
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
                  }
                }
              }
              ip_addresses = {
                type = "array"
                required = true
                description = "List of all IP addresses"
                items = {
                  type = "string"
                }
              }
              primary_ip = {
                type = "string"
                required = false
                description = "Primary IP address"
              }
              # Network statistics
              connections = {
                type = "integer"
                required = false
                description = "Number of active network connections"
              }
              listening_ports = {
                type = "array"
                required = false
                description = "List of listening ports"
                items = {
                  type = "integer"
                }
              }
              bytes_sent = {
                type = "integer"
                required = false
                description = "Total bytes sent"
              }
              bytes_recv = {
                type = "integer"
                required = false
                description = "Total bytes received"
              }
              packets_sent = {
                type = "integer"
                required = false
                description = "Total packets sent"
              }
              packets_recv = {
                type = "integer"
                required = false
                description = "Total packets received"
              }
              err_in = {
                type = "integer"
                required = false
                description = "Total input errors"
              }
              err_out = {
                type = "integer"
                required = false
                description = "Total output errors"
              }
              drop_in = {
                type = "integer"
                required = false
                description = "Total input drops"
              }
              drop_out = {
                type = "integer"
                required = false
                description = "Total output drops"
              }
              # Enhanced network information (gopsutil v4)
              protocols = {
                type = "object"
                required = false
                description = "Network protocol statistics"
                properties = {
                  tcp = {
                    type = "object"
                    required = false
                    description = "TCP protocol statistics"
                    properties = {
                      established = {
                        type = "integer"
                        required = false
                        description = "Number of established TCP connections"
                      }
                      listen = {
                        type = "integer"
                        required = false
                        description = "Number of listening TCP connections"
                      }
                      time_wait = {
                        type = "integer"
                        required = false
                        description = "Number of TCP connections in TIME_WAIT state"
                      }
                      close_wait = {
                        type = "integer"
                        required = false
                        description = "Number of TCP connections in CLOSE_WAIT state"
                      }
                      fin_wait1 = {
                        type = "integer"
                        required = false
                        description = "Number of TCP connections in FIN_WAIT1 state"
                      }
                      fin_wait2 = {
                        type = "integer"
                        required = false
                        description = "Number of TCP connections in FIN_WAIT2 state"
                      }
                      closing = {
                        type = "integer"
                        required = false
                        description = "Number of TCP connections in CLOSING state"
                      }
                      last_ack = {
                        type = "integer"
                        required = false
                        description = "Number of TCP connections in LAST_ACK state"
                      }
                      syn_sent = {
                        type = "integer"
                        required = false
                        description = "Number of TCP connections in SYN_SENT state"
                      }
                      syn_recv = {
                        type = "integer"
                        required = false
                        description = "Number of TCP connections in SYN_RECV state"
                      }
                    }
                  }
                  udp = {
                    type = "object"
                    required = false
                    description = "UDP protocol statistics"
                    properties = {
                      established = {
                        type = "integer"
                        required = false
                        description = "Number of established UDP connections"
                      }
                      listen = {
                        type = "integer"
                        required = false
                        description = "Number of listening UDP connections"
                      }
                    }
                  }
                  icmp = {
                    type = "object"
                    required = false
                    description = "ICMP protocol statistics"
                    properties = {
                      in_msgs = {
                        type = "integer"
                        required = false
                        description = "Number of ICMP messages received"
                      }
                      out_msgs = {
                        type = "integer"
                        required = false
                        description = "Number of ICMP messages sent"
                      }
                      in_errors = {
                        type = "integer"
                        required = false
                        description = "Number of ICMP errors received"
                      }
                      out_errors = {
                        type = "integer"
                        required = false
                        description = "Number of ICMP errors sent"
                      }
                    }
                  }
                }
              }
              connection_details = {
                type = "array"
                required = false
                description = "Detailed network connections"
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
                      type = "string"
                      required = false
                      description = "Local address"
                    }
                    raddr = {
                      type = "string"
                      required = false
                      description = "Remote address"
                    }
                    status = {
                      type = "string"
                      required = false
                      description = "Connection status"
                    }
                    pid = {
                      type = "integer"
                      required = false
                      description = "Process ID using this connection"
                    }
                  }
                }
              }
              interface_stats = {
                type = "array"
                required = false
                description = "Detailed interface statistics"
                items = {
                  type = "object"
                  properties = {
                    name = {
                      type = "string"
                      required = true
                      description = "Interface name"
                    }
                    bytes_sent = {
                      type = "integer"
                      required = false
                      description = "Bytes sent on this interface"
                    }
                    bytes_recv = {
                      type = "integer"
                      required = false
                      description = "Bytes received on this interface"
                    }
                    packets_sent = {
                      type = "integer"
                      required = false
                      description = "Packets sent on this interface"
                    }
                    packets_recv = {
                      type = "integer"
                      required = false
                      description = "Packets received on this interface"
                    }
                    err_in = {
                      type = "integer"
                      required = false
                      description = "Input errors on this interface"
                    }
                    err_out = {
                      type = "integer"
                      required = false
                      description = "Output errors on this interface"
                    }
                    drop_in = {
                      type = "integer"
                      required = false
                      description = "Input drops on this interface"
                    }
                    drop_out = {
                      type = "integer"
                      required = false
                      description = "Output drops on this interface"
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
              netfilter_conntrack = {
                type = "object"
                required = false
                description = "Netfilter connection tracking statistics (Linux only)"
                properties = {
                  entries = {
                    type = "integer"
                    required = false
                    description = "Number of connection tracking entries"
                  }
                  searched = {
                    type = "integer"
                    required = false
                    description = "Number of connection tracking searches"
                  }
                  found = {
                    type = "integer"
                    required = false
                    description = "Number of connection tracking matches found"
                  }
                  new = {
                    type = "integer"
                    required = false
                    description = "Number of new connection tracking entries"
                  }
                  invalid = {
                    type = "integer"
                    required = false
                    description = "Number of invalid connection tracking entries"
                  }
                  ignore = {
                    type = "integer"
                    required = false
                    description = "Number of ignored connection tracking entries"
                  }
                  delete = {
                    type = "integer"
                    required = false
                    description = "Number of deleted connection tracking entries"
                  }
                  delete_list = {
                    type = "integer"
                    required = false
                    description = "Number of connection tracking entries in delete list"
                  }
                  insert = {
                    type = "integer"
                    required = false
                    description = "Number of inserted connection tracking entries"
                  }
                  insert_failed = {
                    type = "integer"
                    required = false
                    description = "Number of failed connection tracking insertions"
                  }
                  drop = {
                    type = "integer"
                    required = false
                    description = "Number of dropped connection tracking entries"
                  }
                  early_drop = {
                    type = "integer"
                    required = false
                    description = "Number of early dropped connection tracking entries"
                  }
                  icmp_error = {
                    type = "integer"
                    required = false
                    description = "Number of ICMP errors in connection tracking"
                  }
                  expect_new = {
                    type = "integer"
                    required = false
                    description = "Number of new expected connection tracking entries"
                  }
                  expect_create = {
                    type = "integer"
                    required = false
                    description = "Number of created expected connection tracking entries"
                  }
                  expect_delete = {
                    type = "integer"
                    required = false
                    description = "Number of deleted expected connection tracking entries"
                  }
                  search_restart = {
                    type = "integer"
                    required = false
                    description = "Number of connection tracking search restarts"
                  }
                }
              }
            }
          }
          
          # Load average facts
          load_average = {
            type = "object"
            required = false
            description = "System load averages"
            
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
          
          # Process information
          processes = {
            type = "object"
            required = false
            description = "Process information"
            
            properties = {
              count = {
                type = "integer"
                required = false
                description = "Total number of processes"
              }
              running = {
                type = "integer"
                required = false
                description = "Number of running processes"
              }
              sleeping = {
                type = "integer"
                required = false
                description = "Number of sleeping processes"
              }
              stopped = {
                type = "integer"
                required = false
                description = "Number of stopped processes"
              }
              zombie = {
                type = "integer"
                required = false
                description = "Number of zombie processes"
              }
              # Top processes by resource usage
              top_by_cpu = {
                type = "array"
                required = false
                description = "Top processes by CPU usage"
                items = {
                  type = "object"
                  properties = {
                    pid = {
                      type = "integer"
                      required = true
                      description = "Process ID"
                    }
                    name = {
                      type = "string"
                      required = true
                      description = "Process name"
                    }
                    cpu_percent = {
                      type = "number"
                      required = true
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
                      required = true
                      description = "Process ID"
                    }
                    name = {
                      type = "string"
                      required = true
                      description = "Process name"
                    }
                    memory_percent = {
                      type = "number"
                      required = true
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
              # Detailed process information (gopsutil v4)
              detailed_processes = {
                type = "array"
                required = false
                description = "Detailed information for individual processes"
                items = {
                  type = "object"
                  properties = {
                    pid = {
                      type = "integer"
                      required = true
                      description = "Process ID"
                    }
                    ppid = {
                      type = "integer"
                      required = false
                      description = "Parent process ID"
                    }
                    name = {
                      type = "string"
                      required = true
                      description = "Process name"
                    }
                    cmdline = {
                      type = "string"
                      required = false
                      description = "Command line"
                    }
                    status = {
                      type = "string"
                      required = false
                      description = "Process status"
                    }
                    create_time = {
                      type = "integer"
                      required = false
                      description = "Process creation time (Unix timestamp)"
                    }
                    cwd = {
                      type = "string"
                      required = false
                      description = "Current working directory"
                    }
                    exe = {
                      type = "string"
                      required = false
                      description = "Executable path"
                    }
                    uids = {
                      type = "array"
                      required = false
                      description = "User IDs"
                      items = {
                        type = "integer"
                      }
                    }
                    gids = {
                      type = "array"
                      required = false
                      description = "Group IDs"
                      items = {
                        type = "integer"
                      }
                    }
                    terminal = {
                      type = "string"
                      required = false
                      description = "Terminal"
                    }
                    nice = {
                      type = "integer"
                      required = false
                      description = "Process nice value"
                    }
                    num_fds = {
                      type = "integer"
                      required = false
                      description = "Number of file descriptors"
                    }
                    num_ctx_switches = {
                      type = "integer"
                      required = false
                      description = "Number of context switches"
                    }
                    num_threads = {
                      type = "integer"
                      required = false
                      description = "Number of threads"
                    }
                    cpu_times = {
                      type = "object"
                      required = false
                      description = "CPU times for this process"
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
                    memory_info = {
                      type = "object"
                      required = false
                      description = "Memory information for this process"
                      properties = {
                        rss = {
                          type = "integer"
                          required = false
                          description = "RSS memory in bytes"
                        }
                        vms = {
                          type = "integer"
                          required = false
                          description = "VMS memory in bytes"
                        }
                        shared = {
                          type = "integer"
                          required = false
                          description = "Shared memory in bytes"
                        }
                        text = {
                          type = "integer"
                          required = false
                          description = "Text memory in bytes"
                        }
                        lib = {
                          type = "integer"
                          required = false
                          description = "Library memory in bytes"
                        }
                        data = {
                          type = "integer"
                          required = false
                          description = "Data memory in bytes"
                        }
                        dirty = {
                          type = "integer"
                          required = false
                          description = "Dirty memory in bytes"
                        }
                      }
                    }
                    memory_maps = {
                      type = "array"
                      required = false
                      description = "Memory maps for this process"
                      items = {
                        type = "object"
                        properties = {
                          path = {
                            type = "string"
                            required = false
                            description = "Memory map path"
                          }
                          rss = {
                            type = "integer"
                            required = false
                            description = "RSS memory in bytes"
                          }
                          size = {
                            type = "integer"
                            required = false
                            description = "Memory map size in bytes"
                          }
                          pss = {
                            type = "integer"
                            required = false
                            description = "Proportional set size in bytes"
                          }
                        }
                      }
                    }
                    open_files = {
                      type = "array"
                      required = false
                      description = "Open files for this process"
                      items = {
                        type = "object"
                        properties = {
                          fd = {
                            type = "integer"
                            required = false
                            description = "File descriptor"
                          }
                          path = {
                            type = "string"
                            required = false
                            description = "File path"
                          }
                        }
                      }
                    }
                    connections = {
                      type = "array"
                      required = false
                      description = "Network connections for this process"
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
                            type = "string"
                            required = false
                            description = "Local address"
                          }
                          raddr = {
                            type = "string"
                            required = false
                            description = "Remote address"
                          }
                          status = {
                            type = "string"
                            required = false
                            description = "Connection status"
                          }
                        }
                      }
                    }
                    cpu_affinity = {
                      type = "array"
                      required = false
                      description = "CPU affinity for this process"
                      items = {
                        type = "integer"
                      }
                    }
                    io_counters = {
                      type = "object"
                      required = false
                      description = "I/O counters for this process"
                      properties = {
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
                          description = "Bytes read"
                        }
                        write_bytes = {
                          type = "integer"
                          required = false
                          description = "Bytes written"
                        }
                      }
                    }
                    page_faults = {
                      type = "object"
                      required = false
                      description = "Page fault information"
                      properties = {
                        minor_faults = {
                          type = "integer"
                          required = false
                          description = "Minor page faults"
                        }
                        major_faults = {
                          type = "integer"
                          required = false
                          description = "Major page faults"
                        }
                        child_minor_faults = {
                          type = "integer"
                          required = false
                          description = "Child minor page faults"
                        }
                        child_major_faults = {
                          type = "integer"
                          required = false
                          description = "Child major page faults"
                        }
                      }
                    }
                    username = {
                      type = "string"
                      required = false
                      description = "Process username"
                    }
                    children = {
                      type = "array"
                      required = false
                      description = "Child process IDs"
                      items = {
                        type = "integer"
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
      
      # Enhanced system information
      enhanced = {
        type = "object"
        required = false
        description = "Enhanced system information"
        
        properties = {
          # Virtualization information
          virtualization = {
            type = "object"
            required = false
            description = "Virtualization information"
            
            properties = {
              system = {
                type = "string"
                required = false
                description = "Virtualization system (kvm, vmware, virtualbox, etc.)"
              }
              role = {
                type = "string"
                required = false
                description = "Virtualization role (host, guest)"
              }
            }
          }
          
          # Package manager detection
          package_manager = {
            type = "object"
            required = false
            description = "Package manager information"
            
            properties = {
              type = {
                type = "string"
                required = true
                description = "Package manager type (apt, yum, dnf, zypper, pacman, apk)"
                enum = ["apt", "yum", "dnf", "zypper", "pacman", "apk", "unknown"]
              }
              version = {
                type = "string"
                required = false
                description = "Package manager version"
              }
              config_path = {
                type = "string"
                required = false
                description = "Package manager configuration path"
              }
            }
          }
          
          # Service manager detection
          service_manager = {
            type = "object"
            required = false
            description = "Service manager information"
            
            properties = {
              type = {
                type = "string"
                required = true
                description = "Service manager type (systemd, upstart, init, runit, openrc)"
                enum = ["systemd", "upstart", "init", "runit", "openrc", "unknown"]
              }
              version = {
                type = "string"
                required = false
                description = "Service manager version"
              }
            }
          }
          
          # SELinux status
          selinux = {
            type = "object"
            required = false
            description = "SELinux status information"
            
            properties = {
              enabled = {
                type = "boolean"
                required = true
                description = "Whether SELinux is enabled"
              }
              mode = {
                type = "string"
                required = false
                description = "SELinux mode (enforcing, permissive, disabled)"
                enum = ["enforcing", "permissive", "disabled"]
              }
              type = {
                type = "string"
                required = false
                description = "SELinux type"
              }
            }
          }
          
          # SSH host keys
          ssh_keys = {
            type = "object"
            required = false
            description = "SSH host keys"
            
            properties = {
              rsa = {
                type = "string"
                required = false
                description = "RSA host key"
              }
              ecdsa = {
                type = "string"
                required = false
                description = "ECDSA host key"
              }
              ed25519 = {
                type = "string"
                required = false
                description = "ED25519 host key"
              }
            }
          }
          
          # BIOS information
          bios = {
            type = "object"
            required = false
            description = "BIOS information"
            
            properties = {
              vendor = {
                type = "string"
                required = false
                description = "BIOS vendor"
              }
              version = {
                type = "string"
                required = false
                description = "BIOS version"
              }
              date = {
                type = "string"
                required = false
                description = "BIOS date"
              }
              release = {
                type = "string"
                required = false
                description = "BIOS release"
              }
              board_vendor = {
                type = "string"
                required = false
                description = "Board vendor"
              }
              board_name = {
                type = "string"
                required = false
                description = "Board name"
              }
              board_version = {
                type = "string"
                required = false
                description = "Board version"
              }
            }
          }
          
          # Sensors information
          sensors = {
            type = "object"
            required = false
            description = "Hardware sensor information"
            
            properties = {
              temperatures = {
                type = "object"
                required = false
                description = "Temperature sensors"
                
                properties = {
                  cpu = {
                    type = "number"
                    required = false
                    description = "CPU temperature in Celsius"
                  }
                  gpu = {
                    type = "number"
                    required = false
                    description = "GPU temperature in Celsius"
                  }
                  motherboard = {
                    type = "number"
                    required = false
                    description = "Motherboard temperature in Celsius"
                  }
                  core_temp = {
                    type = "array"
                    required = false
                    description = "Per-core temperatures"
                    items = {
                      type = "number"
                    }
                  }
                }
              }
              fans = {
                type = "object"
                required = false
                description = "Fan speeds"
                
                properties = {
                  cpu_fan = {
                    type = "integer"
                    required = false
                    description = "CPU fan RPM"
                  }
                  system_fan = {
                    type = "integer"
                    required = false
                    description = "System fan RPM"
                  }
                  case_fan = {
                    type = "integer"
                    required = false
                    description = "Case fan RPM"
                  }
                }
              }
            }
          }
          
          # Docker information (Linux only)
          docker = {
            type = "object"
            required = false
            description = "Docker information (Linux only)"
            
            properties = {
              container_id = {
                type = "string"
                required = false
                description = "Docker container ID if running in container"
              }
              cgroups_cpu = {
                type = "object"
                required = false
                description = "Docker cgroups CPU information"
                
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
                }
              }
              cgroups_memory = {
                type = "object"
                required = false
                description = "Docker cgroups memory information"
                
                properties = {
                  usage = {
                    type = "integer"
                    required = false
                    description = "Memory usage in bytes"
                  }
                  limit = {
                    type = "integer"
                    required = false
                    description = "Memory limit in bytes"
                  }
                  max_usage = {
                    type = "integer"
                    required = false
                    description = "Maximum memory usage in bytes"
                  }
                }
              }
            }
          }
        }
      }
      
      # Application facts
      applications = {
        type = "object"
        required = false
        description = "Application-specific facts"
        
        properties = {
          # Application versions
          versions = {
            type = "object"
            required = false
            description = "Application version information"
            
            properties = {
              nginx = {
                type = "string"
                required = false
                description = "Nginx version"
              }
              apache = {
                type = "string"
                required = false
                description = "Apache version"
              }
              postgresql = {
                type = "string"
                required = false
                description = "PostgreSQL version"
              }
              mysql = {
                type = "string"
                required = false
                description = "MySQL version"
              }
              redis = {
                type = "string"
                required = false
                description = "Redis version"
              }
              docker = {
                type = "string"
                required = false
                description = "Docker version"
              }
              kubernetes = {
                type = "string"
                required = false
                description = "Kubernetes version"
              }
            }
          }
          
          # Application configuration
          config = {
            type = "object"
            required = false
            description = "Application configuration facts"
            
            properties = {
              config_paths = {
                type = "object"
                required = false
                description = "Configuration file paths"
                
                properties = {
                  nginx_conf = {
                    type = "string"
                    required = false
                    description = "Nginx configuration path"
                  }
                  apache_conf = {
                    type = "string"
                    required = false
                    description = "Apache configuration path"
                  }
                  postgres_conf = {
                    type = "string"
                    required = false
                    description = "PostgreSQL configuration path"
                  }
                  redis_conf = {
                    type = "string"
                    required = false
                    description = "Redis configuration path"
                  }
                }
              }
              
              log_paths = {
                type = "object"
                required = false
                description = "Log file paths"
                
                properties = {
                  nginx_logs = {
                    type = "string"
                    required = false
                    description = "Nginx log directory"
                  }
                  apache_logs = {
                    type = "string"
                    required = false
                    description = "Apache log directory"
                  }
                  postgres_logs = {
                    type = "string"
                    required = false
                    description = "PostgreSQL log directory"
                  }
                  redis_logs = {
                    type = "string"
                    required = false
                    description = "Redis log directory"
                  }
                }
              }
            }
          }
        }
      }
      
      # Deployment facts
      deployment = {
        type = "object"
        required = false
        description = "Deployment-related facts"
        
        properties = {
          # Deployment state
          state = {
            type = "string"
            required = false
            description = "Current deployment state"
            enum = ["not_deployed", "deploying", "deployed", "failed", "updating", "rolling_back"]
          }
          
          # Deployment information
          info = {
            type = "object"
            required = false
            description = "Deployment information"
            
            properties = {
              version = {
                type = "string"
                required = false
                description = "Deployed version"
              }
              deployed_at = {
                type = "string"
                required = false
                format = "date-time"
                description = "When deployment was completed"
              }
              deployed_by = {
                type = "string"
                required = false
                description = "Who performed the deployment"
              }
              commit_hash = {
                type = "string"
                required = false
                description = "Git commit hash of deployed version"
              }
              branch = {
                type = "string"
                required = false
                description = "Git branch of deployed version"
              }
            }
          }
          
          # Service status
          services = {
            type = "object"
            required = false
            description = "Service status information"
            
            properties = {
              nginx = {
                type = "string"
                required = false
                description = "Nginx service status"
                enum = ["active", "inactive", "failed", "unknown"]
              }
              apache = {
                type = "string"
                required = false
                description = "Apache service status"
                enum = ["active", "inactive", "failed", "unknown"]
              }
              postgresql = {
                type = "string"
                required = false
                description = "PostgreSQL service status"
                enum = ["active", "inactive", "failed", "unknown"]
              }
              redis = {
                type = "string"
                required = false
                description = "Redis service status"
                enum = ["active", "inactive", "failed", "unknown"]
              }
            }
          }
        }
      }
      
      # Environment facts
      environment = {
        type = "object"
        required = false
        description = "Environment-specific facts"
        
        properties = {
          # Environment variables
          variables = {
            type = "object"
            required = false
            description = "Environment variables"
            
            properties = {
              NODE_ENV = {
                type = "string"
                required = false
                description = "Node.js environment"
                enum = ["development", "staging", "production"]
              }
              DATABASE_URL = {
                type = "string"
                required = false
                description = "Database connection URL"
              }
              REDIS_URL = {
                type = "string"
                required = false
                description = "Redis connection URL"
              }
              LOG_LEVEL = {
                type = "string"
                required = false
                description = "Logging level"
                enum = ["debug", "info", "warn", "error"]
              }
            }
          }
          
          # Infrastructure information
          infrastructure = {
            type = "object"
            required = false
            description = "Infrastructure information"
            
            properties = {
              datacenter = {
                type = "string"
                required = false
                description = "Datacenter identifier"
              }
              rack = {
                type = "string"
                required = false
                description = "Rack identifier"
              }
              power_zone = {
                type = "string"
                required = false
                description = "Power zone identifier"
              }
              region = {
                type = "string"
                required = false
                description = "Cloud region"
              }
              availability_zone = {
                type = "string"
                required = false
                description = "Cloud availability zone"
              }
            }
          }
        }
      }
      
      # Monitoring facts
      monitoring = {
        type = "object"
        required = false
        description = "Monitoring configuration facts"
        
        properties = {
          # Monitoring endpoints
          endpoints = {
            type = "object"
            required = false
            description = "Monitoring endpoints"
            
            properties = {
              prometheus_port = {
                type = "integer"
                required = false
                description = "Prometheus metrics port"
                min = 1024
                max = 65535
              }
              grafana_port = {
                type = "integer"
                required = false
                description = "Grafana port"
                min = 1024
                max = 65535
              }
              alert_manager = {
                type = "string"
                required = false
                description = "Alert manager URL"
              }
            }
          }
          
          # Health check information
          health_checks = {
            type = "object"
            required = false
            description = "Health check configuration"
            
            properties = {
              enabled = {
                type = "boolean"
                required = false
                description = "Whether health checks are enabled"
              }
              port = {
                type = "integer"
                required = false
                description = "Health check port"
                min = 1024
                max = 65535
              }
              path = {
                type = "string"
                required = false
                description = "Health check path"
              }
              interval = {
                type = "string"
                required = false
                description = "Health check interval"
              }
            }
          }
        }
      }
      
      # Custom facts (user-defined)
      custom = {
        type = "object"
        required = false
        description = "User-defined custom facts"
      }
    }
  }

  # Common validation rules (format-agnostic)
  validation = {
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
    
    # Facts structure validation
    facts_required = {
      rule = "required"
      field = "facts"
      message = "Facts object is required"
    }
    
    # No circular references
    no_circular_refs = {
      rule = "acyclic"
      message = "No circular references allowed in facts"
    }
  }

  # Common metadata (format-agnostic)
  metadata = {
    # Schema version
    schema_version = {
      type = "string"
      value = "1.0.0"
      description = "Facts structure schema version"
    }
    
    # Schema type
    schema_type = {
      type = "string"
      value = "facts-structure"
      description = "Schema type identifier"
    }
    
    # Last updated
    last_updated = {
      type = "string"
      value = "2024-01-01"
      description = "Last schema update date"
    }
    
    # Description
    description = {
      type = "string"
      value = "Common fact structure definitions for all storage and export formats"
      description = "Schema description"
    }
  }
} 