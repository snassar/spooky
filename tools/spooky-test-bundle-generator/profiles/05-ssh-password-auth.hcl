profile "ssh_password_auth" "SSH password authentication testing across all operating systems" {
  description = "Tests SSH password authentication functionality with various password configurations across all operating systems"
  
  containers {
    debian13 {
      base_image = "debian:trixie-slim"
      static_ip = "10.0.100.51"
      ssh_port = 2262
      packages = ["openssh-server", "openssh-client", "curl", "wget"]
      ssh_config {
        port = 22
        password_auth = true
        pubkey_auth = false
        permit_root_login = "yes"
        strict_modes = false
        max_auth_tries = 3
        login_grace_time = 60
      }
    }
    
    fedora42 {
      base_image = "fedora:42"
      static_ip = "10.0.100.52"
      ssh_port = 2263
      packages = ["openssh-server", "openssh-clients", "curl", "wget"]
      ssh_config {
        port = 22
        password_auth = true
        pubkey_auth = false
        permit_root_login = "yes"
        strict_modes = false
        max_auth_tries = 3
        login_grace_time = 60
      }
    }
    
    arch {
      base_image = "archlinux:base"
      static_ip = "10.0.100.53"
      ssh_port = 2264
      packages = ["openssh", "curl", "wget"]
      ssh_config {
        port = 22
        password_auth = true
        pubkey_auth = false
        permit_root_login = "yes"
        strict_modes = false
        max_auth_tries = 3
        login_grace_time = 60
      }
    }
    
    alpine319 {
      base_image = "alpine:3.19"
      static_ip = "10.0.100.54"
      ssh_port = 2265
      packages = ["openssh-server", "openssh-client", "curl", "wget"]
      ssh_config {
        port = 22
        password_auth = true
        pubkey_auth = false
        permit_root_login = "yes"
        strict_modes = false
        max_auth_tries = 3
        login_grace_time = 60
      }
    }
    
    opensuse156 {
      base_image = "opensuse/leap:15.6"
      static_ip = "10.0.100.55"
      ssh_port = 2266
      packages = ["openssh-server", "openssh-clients", "curl", "wget"]
      ssh_config {
        port = 22
        password_auth = true
        pubkey_auth = false
        permit_root_login = "yes"
        strict_modes = false
        max_auth_tries = 3
        login_grace_time = 60
      }
    }
  }
  
  project {
    name = "ssh-password-auth-test"
    description = "Cross-OS SSH password authentication test"
    
    machines {
      debian13 {
        hostname = "debian13-password"
        ip = "10.0.100.51"
        port = 2262
        auth {
          password = "SecurePass123!"
          username = "root"
        }
        facts {
          basic = true
          enhanced = false
          custom = false
          encrypted = false
        }
      }
      
      fedora42 {
        hostname = "fedora42-password"
        ip = "10.0.100.52"
        port = 2263
        auth {
          password = "SecurePass123!"
          username = "root"
        }
        facts {
          basic = true
          enhanced = false
          custom = false
          encrypted = false
        }
      }
      
      arch {
        hostname = "arch-password"
        ip = "10.0.100.53"
        port = 2264
        auth {
          password = "SecurePass123!"
          username = "root"
        }
        facts {
          basic = true
          enhanced = false
          custom = false
          encrypted = false
        }
      }
      
      alpine319 {
        hostname = "alpine319-password"
        ip = "10.0.100.54"
        port = 2265
        auth {
          password = "SecurePass123!"
          username = "root"
        }
        facts {
          basic = true
          enhanced = false
          custom = false
          encrypted = false
        }
      }
      
      opensuse156 {
        hostname = "opensuse156-password"
        ip = "10.0.100.55"
        port = 2266
        auth {
          password = "SecurePass123!"
          username = "root"
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
      test_type = "ssh-password-auth"
      os_list = ["debian13", "fedora42", "arch", "alpine319", "opensuse156"]
      auth_timeout = "30s"
      connection_retries = 3
    }
    
    actions {
      test_password_connection {
        name = "test-password-connection"
        description = "Test SSH password authentication on all machines"
        command = "ssh --test-connection"
        tags = ["ssh", "password", "authentication", "cross-os"]
        parallel = true
        timeout = "30s"
      }
      
      verify_password_auth {
        name = "verify-password-auth"
        description = "Verify password authentication works correctly"
        command = "ssh --verify-auth"
        tags = ["ssh", "verification", "password", "cross-os"]
        depends_on = ["test_password_connection"]
        timeout = "15s"
      }
      
      test_invalid_password {
        name = "test-invalid-password"
        description = "Test SSH connection with invalid password"
        command = "ssh --test-invalid-auth"
        tags = ["ssh", "negative-test", "password", "cross-os"]
        depends_on = ["verify_password_auth"]
        timeout = "20s"
      }
      
      gather_facts_via_password {
        name = "gather-facts-via-password"
        description = "Gather facts using password authentication"
        command = "facts --auth-method=password"
        tags = ["facts", "ssh", "password", "cross-os"]
        depends_on = ["verify_password_auth"]
        parallel = true
        timeout = "45s"
      }
    }
  }
}
