# Simple test HCL file with basic blocks for machine facts

# Required machine facts structure
machine_id = "abcdef1234567890abcdef1234567890"
collected_at = "2024-01-01T12:00:00Z"

# Machine facts collection
facts {
  # System-level facts
  system {
    os {
      name = "ubuntu"
      version = "22.04"
      arch = "x86_64"
      kernel = "5.15.0-42-generic"
      family = "debian"
    }
    
    hardware {
      cpu {
        cores = 4
        model = "Intel Core i7"
        frequency = 2400
      }
      
      memory {
        total_gb = 8
        available_gb = 6
      }
    }
  }

  # Simple fact blocks for machine information
  fact "os_info" {
    name = "ubuntu"
    version = "22.04"
  }

  fact "cpu_info" {
    cores = 4
    model = "Intel Core i7"
  }

  # Simple resource block for machine resources
  resource "local_file" "config" {
    filename = "/tmp/config.txt"
    content = "Hello, World!"
  }

  # Simple variable block for machine configuration
  variable "app_name" {
    type = "string"
    default = "my-app"
  }
}
