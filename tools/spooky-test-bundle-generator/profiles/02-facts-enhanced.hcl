profile "facts_enhanced" "Enhanced facts gathering with additional system information" {
  description = "Tests enhanced facts gathering functionality with additional system information across all operating systems"
  
  containers {
    debian13 {
      base_image = "debian:trixie-slim"
      static_ip = "10.0.100.21"
      ssh_port = 2232
      packages = ["openssh-server", "openssh-client", "curl", "wget", "procps", "util-linux", "dmidecode"]
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
      static_ip = "10.0.100.22"
      ssh_port = 2233
      packages = ["openssh-server", "openssh-clients", "curl", "wget", "procps-ng", "util-linux", "dmidecode"]
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
      static_ip = "10.0.100.23"
      ssh_port = 2234
      packages = ["openssh", "curl", "wget", "procps-ng", "util-linux", "dmidecode"]
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
      static_ip = "10.0.100.24"
      ssh_port = 2235
      packages = ["openssh-server", "openssh-client", "curl", "wget", "procps", "util-linux", "dmidecode"]
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
      static_ip = "10.0.100.25"
      ssh_port = 2236
      packages = ["openssh-server", "openssh-clients", "curl", "wget", "procps", "util-linux", "dmidecode"]
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
    name = "facts-enhanced-test"
    description = "Cross-OS enhanced facts gathering test"
    
    machines {
      debian13 {
        hostname = "debian13-enhanced"
        ip = "10.0.100.21"
        port = 2232
        auth {
          password = "testpass123"
        }
        facts {
          basic = true
          enhanced = true
          custom = false
          encrypted = false
        }
      }
      
      fedora42 {
        hostname = "fedora42-enhanced"
        ip = "10.0.100.22"
        port = 2233
        auth {
          password = "testpass123"
        }
        facts {
          basic = true
          enhanced = true
          custom = false
          encrypted = false
        }
      }
      
      arch {
        hostname = "arch-enhanced"
        ip = "10.0.100.23"
        port = 2234
        auth {
          password = "testpass123"
        }
        facts {
          basic = true
          enhanced = true
          custom = false
          encrypted = false
        }
      }
      
      alpine319 {
        hostname = "alpine319-enhanced"
        ip = "10.0.100.24"
        port = 2235
        auth {
          password = "testpass123"
        }
        facts {
          basic = true
          enhanced = true
          custom = false
          encrypted = false
        }
      }
      
      opensuse156 {
        hostname = "opensuse156-enhanced"
        ip = "10.0.100.25"
        port = 2236
        auth {
          password = "testpass123"
        }
        facts {
          basic = true
          enhanced = true
          custom = false
          encrypted = false
        }
      }
    }
    
    variables {
      test_type = "facts-enhanced"
      os_list = ["debian13", "fedora42", "arch", "alpine319", "opensuse156"]
      expected_facts = ["os", "hostname", "ip", "architecture", "kernel", "cpu_info", "memory_info", "disk_info", "network_info"]
      enhanced_facts_timeout = "60s"
    }
    
    actions {
      gather_enhanced_facts {
        name = "gather-enhanced-facts"
        description = "Gather enhanced system facts from all machines"
        command = "facts --enhanced"
        tags = ["facts", "enhanced", "cross-os"]
        parallel = true
        timeout = "60s"
      }
      
      verify_enhanced_facts {
        name = "verify-enhanced-facts"
        description = "Verify that enhanced facts are gathered correctly"
        command = "facts --verify --enhanced"
        tags = ["facts", "verification", "enhanced", "cross-os"]
        depends_on = ["gather_enhanced_facts"]
        timeout = "15s"
      }
      
      compare_facts {
        name = "compare-facts-across-os"
        description = "Compare facts across different operating systems"
        command = "facts --compare"
        tags = ["facts", "comparison", "cross-os"]
        depends_on = ["verify_enhanced_facts"]
        timeout = "20s"
      }
    }
  }
}
