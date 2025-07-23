# Example template using custom facts
# This template demonstrates how to use custom facts in configuration templates

action "deploy-application" {
  description = "Deploy application with custom facts integration"
  
  # Use custom facts for application configuration
  command = "deploy.sh --app {{custom 'application.name'}} --version {{custom 'application.version'}} --port {{custom 'application.port'}}"
  
  # Use custom facts for environment configuration
  environment = {
    DATACENTER = "{{custom 'environment.datacenter'}}"
    RACK       = "{{custom 'environment.rack'}}"
    ZONE       = "{{custom 'environment.zone'}}"
  }
  
  # Use custom facts for monitoring configuration
  monitoring = {
    PROMETHEUS_PORT = "{{custom 'monitoring.prometheus_port'}}"
    ALERT_MANAGER   = "{{custom 'monitoring.alert_manager'}}"
  }
  
  # Use system facts for server information
  server_info = {
    HOSTNAME = "{{system 'hostname'}}"
    OS_NAME  = "{{system 'os.name'}}"
    OS_VERSION = "{{system 'os.version'}}"
  }
  
  # Use environment variables
  env_vars = {
    NODE_ENV = "{{env 'NODE_ENV'}}"
    LOG_LEVEL = "{{env 'LOG_LEVEL'}}"
  }
  
  # Use overrides if available
  os_override = {
    OS_NAME = "{{custom 'os.name'}}"
    OS_VERSION = "{{custom 'os.version'}}"
  }
  
  servers = ["web-001", "web-002"]
}

# Example of conditional logic based on custom facts
action "conditional-deploy" {
  description = "Conditional deployment based on custom facts"
  
  # This would be evaluated by the template engine
  command = "{{if custom 'environment.zone' == 'production'}}deploy-prod.sh{{else}}deploy-dev.sh{{end}}"
  
  servers = ["web-001"]
}

# Example of using custom facts in file templates
file "nginx.conf" {
  content = <<-EOF
    server {
        listen {{custom 'application.port'}};
        server_name {{system 'hostname'}};
        
        location / {
            root /var/www/{{custom 'application.name'}};
            index index.html;
        }
        
        # Monitoring endpoint
        location /metrics {
            proxy_pass http://localhost:{{custom 'monitoring.prometheus_port'}};
        }
    }
  EOF
  
  destination = "/etc/nginx/sites-available/{{custom 'application.name'}}"
  servers = ["web-001"]
} 