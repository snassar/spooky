profile "sync_basic" "Basic file synchronization testing across all operating systems" {
  description = "Tests basic file synchronization functionality with Mutagen across all operating systems"
  
  containers {
    debian13 {
      base_image = "debian:trixie-slim"
      static_ip = "10.0.100.121"
      ssh_port = 2332
      packages = ["openssh-server", "openssh-client", "curl", "wget", "rsync", "inotify-tools"]
      ssh_config {
        port = 22
        password_auth = true
        pubkey_auth = false
        permit_root_login = "yes"
        strict_modes = false
      }
    }
    
    fedora42 {
      base_image = "fedora:42"
      static_ip = "10.0.100.122"
      ssh_port = 2333
      packages = ["openssh-server", "openssh-clients", "curl", "wget", "rsync", "inotify-tools"]
      ssh_config {
        port = 22
        password_auth = true
        pubkey_auth = false
        permit_root_login = "yes"
        strict_modes = false
      }
    }
    
    arch {
      base_image = "archlinux:base"
      static_ip = "10.0.100.123"
      ssh_port = 2334
      packages = ["openssh", "curl", "wget", "rsync", "inotify-tools"]
      ssh_config {
        port = 22
        password_auth = true
        pubkey_auth = false
        permit_root_login = "yes"
        strict_modes = false
      }
    }
    
    alpine319 {
      base_image = "alpine:3.19"
      static_ip = "10.0.100.124"
      ssh_port = 2335
      packages = ["openssh-server", "openssh-client", "curl", "wget", "rsync", "inotify-tools"]
      ssh_config {
        port = 22
        password_auth = true
        pubkey_auth = false
        permit_root_login = "yes"
        strict_modes = false
      }
    }
    
    opensuse156 {
      base_image = "opensuse/leap:15.6"
      static_ip = "10.0.100.125"
      ssh_port = 2336
      packages = ["openssh-server", "openssh-clients", "curl", "wget", "rsync", "inotify-tools"]
      ssh_config {
        port = 22
        password_auth = true
        pubkey_auth = false
        permit_root_login = "yes"
        strict_modes = false
      }
    }
  }
  
  project {
    name = "sync-basic-test"
    description = "Cross-OS basic file synchronization test"
    
    machines {
      debian13 {
        hostname = "debian13-sync"
        ip = "10.0.100.121"
        port = 2332
        auth {
          password = "testpass123"
        }
        facts {
          basic = true
          enhanced = false
          custom = false
          encrypted = false
        }
      }
      
      fedora42 {
        hostname = "fedora42-sync"
        ip = "10.0.100.122"
        port = 2333
        auth {
          password = "testpass123"
        }
        facts {
          basic = true
          enhanced = false
          custom = false
          encrypted = false
        }
      }
      
      arch {
        hostname = "arch-sync"
        ip = "10.0.100.123"
        port = 2334
        auth {
          password = "testpass123"
        }
        facts {
          basic = true
          enhanced = false
          custom = false
          encrypted = false
        }
      }
      
      alpine319 {
        hostname = "alpine319-sync"
        ip = "10.0.100.124"
        port = 2335
        auth {
          password = "testpass123"
        }
        facts {
          basic = true
          enhanced = false
          custom = false
          encrypted = false
        }
      }
      
      opensuse156 {
        hostname = "opensuse156-sync"
        ip = "10.0.100.125"
        port = 2336
        auth {
          password = "testpass123"
        }
        facts {
          basic = true
          enhanced = false
          custom = false
          encrypted = false
        }
      }
    }
    
    variables {
      test_type = "sync-basic"
      os_list = ["debian13", "fedora42", "arch", "alpine319", "opensuse156"]
      sync_mode = "two-way-resolved"
      ignore_patterns = ["*.tmp", "*.log", ".git/", "node_modules/"]
      sync_interval = "5s"
    }
    
    files {
      test_data {
        name = "test-data.txt"
        path = "/tmp/sync-test-data.txt"
        content = <<-EOF
# Test data for file synchronization
This file contains test data for cross-OS file synchronization testing.

OS: {{ .os }}
Hostname: {{ .hostname }}
IP: {{ .ip }}
Timestamp: {{ .timestamp }}

Test content:
- Line 1: Basic synchronization test
- Line 2: Cross-OS compatibility
- Line 3: Mutagen integration
- Line 4: File permissions test
- Line 5: Directory structure test

End of test data.
EOF
        permissions = "644"
        owner = "root"
        group = "root"
      }
      
      sync_config {
        name = "sync.conf"
        path = "/etc/spooky/sync.conf"
        content = <<-EOF
# File synchronization configuration
[sync]
mode = {{ .sync_mode }}
interval = {{ .sync_interval }}
max_file_size = 100MB
preserve_permissions = true
preserve_timestamps = true

[ignore]
patterns = {{ .ignore_patterns }}

[paths]
local = /tmp/spooky-sync-local
remote = /tmp/spooky-sync-remote

[logging]
level = info
file = /var/log/spooky/sync.log
EOF
        permissions = "644"
        owner = "root"
        group = "root"
      }
    }
    
    actions {
      setup_sync_directories {
        name = "setup-sync-directories"
        description = "Setup synchronization directories on all machines"
        command = "sync --setup"
        tags = ["sync", "setup", "cross-os"]
        parallel = true
        timeout = "30s"
      }
      
      create_test_files {
        name = "create-test-files"
        description = "Create test files for synchronization"
        command = "sync --create-test-files"
        tags = ["sync", "files", "cross-os"]
        depends_on = ["setup_sync_directories"]
        parallel = true
        timeout = "20s"
      }
      
      start_sync {
        name = "start-sync"
        description = "Start file synchronization on all machines"
        command = "sync --start"
        tags = ["sync", "start", "cross-os"]
        depends_on = ["create_test_files"]
        parallel = true
        timeout = "45s"
      }
      
      verify_sync {
        name = "verify-sync"
        description = "Verify that files are synchronized correctly"
        command = "sync --verify"
        tags = ["sync", "verification", "cross-os"]
        depends_on = ["start_sync"]
        timeout = "30s"
      }
      
      test_file_modifications {
        name = "test-file-modifications"
        description = "Test file modification synchronization"
        command = "sync --test-modifications"
        tags = ["sync", "testing", "cross-os"]
        depends_on = ["verify_sync"]
        parallel = true
        timeout = "40s"
      }
      
      test_conflict_resolution {
        name = "test-conflict-resolution"
        description = "Test conflict resolution in two-way sync"
        command = "sync --test-conflicts"
        tags = ["sync", "conflicts", "cross-os"]
        depends_on = ["test_file_modifications"]
        parallel = true
        timeout = "35s"
      }
      
      stop_sync {
        name = "stop-sync"
        description = "Stop file synchronization"
        command = "sync --stop"
        tags = ["sync", "stop", "cross-os"]
        depends_on = ["test_conflict_resolution"]
        parallel = true
        timeout = "20s"
      }
      
      cleanup_sync {
        name = "cleanup-sync"
        description = "Cleanup synchronization test data"
        command = "sync --cleanup"
        tags = ["sync", "cleanup", "cross-os"]
        depends_on = ["stop_sync"]
        parallel = true
        timeout = "15s"
      }
    }
  }
}
