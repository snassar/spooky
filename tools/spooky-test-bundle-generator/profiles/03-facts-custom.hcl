profile "facts_custom" "Custom facts gathering with user-defined scripts" {
  description = "Tests custom facts gathering functionality with user-defined scripts across all operating systems"
  
  containers {
    debian13 {
      base_image = "debian:trixie-slim"
      static_ip = "10.0.100.31"
      ssh_port = 2242
      packages = ["openssh-server", "openssh-client", "curl", "wget", "python3", "jq"]
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
      static_ip = "10.0.100.32"
      ssh_port = 2243
      packages = ["openssh-server", "openssh-clients", "curl", "wget", "python3", "jq"]
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
      static_ip = "10.0.100.33"
      ssh_port = 2244
      packages = ["openssh", "curl", "wget", "python", "jq"]
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
      static_ip = "10.0.100.34"
      ssh_port = 2245
      packages = ["openssh-server", "openssh-client", "curl", "wget", "python3", "jq"]
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
      static_ip = "10.0.100.35"
      ssh_port = 2246
      packages = ["openssh-server", "openssh-clients", "curl", "wget", "python3", "jq"]
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
    name = "facts-custom-test"
    description = "Cross-OS custom facts gathering test"
    
    machines {
      debian13 {
        hostname = "debian13-custom"
        ip = "10.0.100.31"
        port = 2242
        auth {
          password = "testpass123"
        }
        facts {
          basic = true
          enhanced = false
          custom = true
          encrypted = false
        }
      }
      
      fedora42 {
        hostname = "fedora42-custom"
        ip = "10.0.100.32"
        port = 2243
        auth {
          password = "testpass123"
        }
        facts {
          basic = true
          enhanced = false
          custom = true
          encrypted = false
        }
      }
      
      arch {
        hostname = "arch-custom"
        ip = "10.0.100.33"
        port = 2244
        auth {
          password = "testpass123"
        }
        facts {
          basic = true
          enhanced = false
          custom = true
          encrypted = false
        }
      }
      
      alpine319 {
        hostname = "alpine319-custom"
        ip = "10.0.100.34"
        port = 2245
        auth {
          password = "testpass123"
        }
        facts {
          basic = true
          enhanced = false
          custom = true
          encrypted = false
        }
      }
      
      opensuse156 {
        hostname = "opensuse156-custom"
        ip = "10.0.100.35"
        port = 2246
        auth {
          password = "testpass123"
        }
        facts {
          basic = true
          enhanced = false
          custom = true
          encrypted = false
        }
      }
    }
    
    variables {
      test_type = "facts-custom"
      os_list = ["debian13", "fedora42", "arch", "alpine319", "opensuse156"]
      custom_facts_dir = "/etc/spooky/facts.d"
      custom_facts_timeout = "45s"
    }
    
    files {
      custom_facts_script {
        name = "custom-facts.sh"
        path = "/etc/spooky/facts.d/custom-facts.sh"
        content = <<-EOF
#!/bin/bash
# Custom facts script for testing
echo "custom_fact_1=test_value_1"
echo "custom_fact_2=test_value_2"
echo "custom_fact_3=$(date +%s)"
echo "custom_fact_4=$(hostname)"
echo "custom_fact_5=$(uname -r)"
EOF
        permissions = "755"
        owner = "root"
        group = "root"
      }
      
      python_facts_script {
        name = "python-facts.py"
        path = "/etc/spooky/facts.d/python-facts.py"
        content = <<-EOF
#!/usr/bin/env python3
import json
import platform
import subprocess
import sys

def get_custom_facts():
    facts = {}
    
    # Python-specific facts
    facts['python_version'] = platform.python_version()
    facts['python_implementation'] = platform.python_implementation()
    
    # System-specific facts
    facts['platform'] = platform.platform()
    facts['processor'] = platform.processor()
    
    # Custom command execution
    try:
        result = subprocess.run(['uptime'], capture_output=True, text=True, timeout=5)
        facts['uptime'] = result.stdout.strip()
    except:
        facts['uptime'] = "unknown"
    
    return facts

if __name__ == "__main__":
    facts = get_custom_facts()
    for key, value in facts.items():
        print(f"{key}={value}")
EOF
        permissions = "755"
        owner = "root"
        group = "root"
      }
    }
    
    actions {
      setup_custom_facts {
        name = "setup-custom-facts"
        description = "Setup custom facts scripts on all machines"
        command = "facts --setup-custom"
        tags = ["facts", "custom", "setup", "cross-os"]
        parallel = true
        timeout = "30s"
      }
      
      gather_custom_facts {
        name = "gather-custom-facts"
        description = "Gather custom facts from all machines"
        command = "facts --custom"
        tags = ["facts", "custom", "cross-os"]
        depends_on = ["setup_custom_facts"]
        parallel = true
        timeout = "45s"
      }
      
      verify_custom_facts {
        name = "verify-custom-facts"
        description = "Verify that custom facts are gathered correctly"
        command = "facts --verify --custom"
        tags = ["facts", "verification", "custom", "cross-os"]
        depends_on = ["gather_custom_facts"]
        timeout = "15s"
      }
      
      test_custom_facts_execution {
        name = "test-custom-facts-execution"
        description = "Test custom facts script execution"
        command = "facts --test-custom"
        tags = ["facts", "testing", "custom", "cross-os"]
        depends_on = ["verify_custom_facts"]
        timeout = "20s"
      }
    }
  }
}
