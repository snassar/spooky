# Valid Templates Configuration
# This file should pass all schema validations

metadata {
  schema_version = "0.20250809.0"
  schema_type = "templates"
  schema_name = "Test Templates Configuration"
  last_updated = "2024-01-01"
  compatibility = ["0.20250809.0"]
  description = "A test templates configuration that should pass all validations"
}

templates {
  # Nginx configuration template
  template "nginx_config" {
    description = "Nginx server configuration template"
    source = "templates/nginx.conf.tmpl"
    destination = "/etc/nginx/nginx.conf"
    
    variables = {
      server_name = "{{ .server_name }}"
      port = "{{ .port }}"
      ssl_cert = "{{ .ssl_cert }}"
      ssl_key = "{{ .ssl_key }}"
    }
    
    validation {
      syntax_check = true
      test_nginx_config = true
      backup_existing = true
    }
    
    permissions {
      owner = "root"
      group = "root"
      mode = "0644"
    }
    
    metadata {
      author = "DevOps Team"
      version = "1.0.0"
      tags = ["nginx", "web", "configuration"]
    }
  }
  
  # Application configuration template
  template "app_config" {
    description = "Application configuration template"
    source = "templates/app.conf.tmpl"
    destination = "/opt/app/config/app.conf"
    
    variables = {
      db_host = "{{ .db_host }}"
      db_port = "{{ .db_port }}"
      db_name = "{{ .db_name }}"
      log_level = "{{ .log_level }}"
      max_connections = "{{ .max_connections }}"
    }
    
    validation {
      syntax_check = true
      test_config = true
      backup_existing = true
    }
    
    permissions {
      owner = "app"
      group = "app"
      mode = "0640"
    }
    
    metadata {
      author = "Application Team"
      version = "2.1.0"
      tags = ["application", "configuration", "database"]
    }
  }
  
  # Systemd service template
  template "systemd_service" {
    description = "Systemd service unit template"
    source = "templates/app.service.tmpl"
    destination = "/etc/systemd/system/app.service"
    
    variables = {
      app_name = "{{ .app_name }}"
      app_path = "{{ .app_path }}"
      user = "{{ .user }}"
      group = "{{ .group }}"
      environment = "{{ .environment }}"
    }
    
    validation {
      syntax_check = true
      test_systemd = true
      backup_existing = true
    }
    
    permissions {
      owner = "root"
      group = "root"
      mode = "0644"
    }
    
    metadata {
      author = "System Team"
      version = "1.0.0"
      tags = ["systemd", "service", "system"]
    }
  }
  
  # Docker Compose template
  template "docker_compose" {
    description = "Docker Compose configuration template"
    source = "templates/docker-compose.yml.tmpl"
    destination = "/opt/app/docker-compose.yml"
    
    variables = {
      app_image = "{{ .app_image }}"
      app_version = "{{ .app_version }}"
      db_image = "{{ .db_image }}"
      db_version = "{{ .db_version }}"
      network_name = "{{ .network_name }}"
    }
    
    validation {
      syntax_check = true
      test_docker_compose = true
      backup_existing = true
    }
    
    permissions {
      owner = "app"
      group = "app"
      mode = "0644"
    }
    
    metadata {
      author = "Container Team"
      version = "3.0.0"
      tags = ["docker", "compose", "container"]
    }
  }
  
  # Environment file template
  template "env_file" {
    description = "Environment variables template"
    source = "templates/.env.tmpl"
    destination = "/opt/app/.env"
    
    variables = {
      NODE_ENV = "{{ .NODE_ENV }}"
      DATABASE_URL = "{{ .DATABASE_URL }}"
      REDIS_URL = "{{ .REDIS_URL }}"
      API_KEY = "{{ .API_KEY }}"
      LOG_LEVEL = "{{ .LOG_LEVEL }}"
    }
    
    validation {
      syntax_check = true
      test_env_vars = true
      backup_existing = true
    }
    
    permissions {
      owner = "app"
      group = "app"
      mode = "0600"
    }
    
    metadata {
      author = "DevOps Team"
      version = "1.0.0"
      tags = ["environment", "variables", "configuration"]
    }
  }
}
