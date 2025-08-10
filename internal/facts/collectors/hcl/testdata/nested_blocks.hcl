# Test HCL file with nested blocks for machine facts

# Required machine facts structure
machine_id = "1234567890abcdef1234567890abcdef"
collected_at = "2024-01-01T00:00:00Z"

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
        cores = 8
        model = "Intel Core i7"
        frequency = 2400
      }
      
      memory {
        total_gb = 16
        available_gb = 12
      }
      
      disk {
        total_gb = 500
        available_gb = 450
      }
    }
  }

  # Resource with nested blocks for machine resources
  resource "aws_instance" "web" {
    instance_type = "t3.micro"
    ami = "ami-123456"
    
    # Nested block
    root_block_device {
      volume_size = 20
      volume_type = "gp3"
    }
    
    # Another nested block
    ebs_block_device {
      device_name = "/dev/sdf"
      volume_size = 100
      volume_type = "gp3"
    }
  }

  # Fact group with nested facts for organizing machine facts
  fact_group "system" {
    description = "System information"
    
    fact "os" {
      name = "ubuntu"
      version = "22.04"
    }
    
    fact "hardware" {
      cpu_cores = 8
      memory_gb = 16
    }
  }

  # Variable with nested validation for machine configuration
  variable "app_environment" {
    type = "string"
    default = "production"
    description = "Application environment"
    
    validation {
      condition = "contains(['development', 'staging', 'production'], var.app_environment)"
      error_message = "Environment must be one of: development, staging, production"
    }
  }
}
