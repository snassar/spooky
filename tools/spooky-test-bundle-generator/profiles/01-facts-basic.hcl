profile "facts_basic" "Basic facts gathering across all operating systems" {
  description = "Tests basic facts gathering functionality with minimal configuration across Debian 13, Fedora 42, Arch Linux, Alpine 3.19, and openSUSE Leap 15.6"
  
  containers {
    debian13 {
      base_image = "debian:trixie-slim"
      static_ip = "10.0.100.11"
      ssh_port = 2222
      packages = ["openssh-server", "openssh-client", "curl", "wget"]
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
      static_ip = "10.0.100.12"
      ssh_port = 2223
      packages = ["openssh-server", "openssh-clients", "curl", "wget"]
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
      static_ip = "10.0.100.13"
      ssh_port = 2224
      packages = ["openssh", "curl", "wget"]
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
      static_ip = "10.0.100.14"
      ssh_port = 2225
      packages = ["openssh-server", "openssh-client", "curl", "wget"]
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
      static_ip = "10.0.100.15"
      ssh_port = 2226
      packages = ["openssh-server", "openssh-clients", "curl", "wget"]
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
    name = "facts-basic-test"
    description = "Cross-OS basic facts gathering test"
    
    machines {
      debian13 {
        hostname = "debian13-test"
        ip = "10.0.100.11"
        port = 2222
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
        hostname = "fedora42-test"
        ip = "10.0.100.12"
        port = 2223
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
        hostname = "arch-test"
        ip = "10.0.100.13"
        port = 2224
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
        hostname = "alpine319-test"
        ip = "10.0.100.14"
        port = 2225
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
        hostname = "opensuse156-test"
        ip = "10.0.100.15"
        port = 2226
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
      test_type = "facts-basic"
      os_list = ["debian13", "fedora42", "arch", "alpine319", "opensuse156"]
      expected_facts = ["os", "hostname", "ip", "architecture", "kernel"]
    }
    
    actions {
      gather_facts {
        name = "gather-basic-facts"
        description = "Gather basic system facts from all machines"
        command = "facts"
        tags = ["facts", "basic", "cross-os"]
        parallel = true
        timeout = "30s"
      }
      
      verify_facts {
        name = "verify-basic-facts"
        description = "Verify that basic facts are gathered correctly"
        command = "facts --verify"
        tags = ["facts", "verification", "cross-os"]
        depends_on = ["gather_facts"]
        timeout = "10s"
      }
    }
  }
}
