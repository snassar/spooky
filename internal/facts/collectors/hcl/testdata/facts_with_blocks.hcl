# Test HCL file with blocks for facts collection
# This file tests the block parsing functionality for single machine facts

# Required machine facts structure
machine_id = "1234567890abcdef1234567890abcdef"
collected_at = "2024-01-01T00:00:00Z"

# Machine facts collection
facts {
  # System-level facts
  system {
    os {
      name = "ubuntu"
      version = "20.04"
      arch = "x86_64"
      kernel = "5.4.0-42-generic"
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

  # Fact blocks for additional machine information
  fact "system_info" {
    os = "linux"
    version = "20.04"
    arch = "x86_64"
    kernel = "5.4.0-42-generic"
  }

  fact "hardware_info" {
    cpu_count = 8
    memory_gb = 16
    disk_gb = 500
  }

  # Resource blocks for machine resources
  resource "aws_instance" "web_server" {
    instance_type = "t3.micro"
    ami = "ami-123456"
    tags = {
      Name = "web-server"
      Environment = "production"
    }
  }

  resource "aws_security_group" "web_sg" {
    name = "web-security-group"
    description = "Security group for web servers"
    
    ingress {
      from_port = 80
      to_port = 80
      protocol = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    }
    
    ingress {
      from_port = 443
      to_port = 443
      protocol = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    }
  }

  # Variable blocks for machine configuration
  variable "environment" {
    type = "string"
    default = "production"
    description = "Environment name"
  }

  variable "instance_count" {
    type = "number"
    default = 3
    description = "Number of instances to create"
  }

  # Output blocks for machine outputs
  output "public_ip" {
    value = "192.168.1.100"
    description = "Public IP address of the web server"
  }

  output "instance_id" {
    value = "i-1234567890abcdef0"
    description = "Instance ID of the web server"
  }

  # Data blocks for machine data sources
  data "aws_ami" "ubuntu" {
    most_recent = true
    owners = ["099720109477"]
    
    filter {
      name = "name"
      values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
    }
  }

  # Local blocks for machine local values
  local "common_tags" {
    Environment = "production"
    Project = "web-app"
    Owner = "devops-team"
  }

  # Fact group blocks for organizing machine facts
  fact_group "monitoring" {
    description = "Monitoring-related facts"
    
    fact "prometheus_enabled" {
      value = true
      description = "Whether Prometheus monitoring is enabled"
    }
    
    fact "grafana_url" {
      value = "https://grafana.example.com"
      description = "Grafana dashboard URL"
    }
  }
}
