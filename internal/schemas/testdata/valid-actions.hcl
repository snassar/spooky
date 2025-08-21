# Valid Actions Configuration
# This file should pass all schema validations

metadata {
  schema_version = "0.20250809.0"
  schema_type = "actions"
  schema_name = "Test Actions Configuration"
  last_updated = "2024-01-01"
  compatibility = ["0.20250809.0"]
  description = "A test actions configuration that should pass all validations"
}

actions {
  action "deploy_webapp" {
    description = "Deploy web application to servers"
    type = "script"
    script = "files/deploy.sh"
    
    variables = {
      app_version = "1.2.3"
      environment = "production"
    }
    
    targets = ["web-server-01", "web-server-02"]
    
    execution {
      timeout = 300
      parallel = true
      max_parallel = 2
      dry_run = false
    }
    
    validation {
      pre_check = true
      post_check = true
      rollback_on_failure = true
    }
  }
  
  action "backup_database" {
    description = "Create database backup"
    type = "command"
    command = "pg_dump -h localhost -U postgres mydb > backup.sql"
    
    targets = ["db-server-01"]
    
    execution {
      timeout = 600
      parallel = false
      dry_run = false
    }
    
    validation {
      pre_check = true
      post_check = false
      rollback_on_failure = false
    }
  }
  
  action "deploy_nginx_config" {
    description = "Deploy nginx configuration using template"
    type = "template_deploy"
    
    template {
      source = "templates/nginx.conf.tmpl"
      destination = "/etc/nginx/nginx.conf"
      validate = true
      backup = true
    }
    
    targets = ["web-server-01"]
    
    execution {
      timeout = 120
      parallel = false
      dry_run = false
    }
  }
  
  action "deploy_app_config" {
    description = "Deploy application configuration using template"
    type = "template_deploy"
    
    template {
      source = "templates/app.conf.tmpl"
      destination = "/opt/app/config/app.conf"
      validate = true
      backup = true
    }
    
    targets = ["web-server-01"]
    
    execution {
      timeout = 120
      parallel = false
      dry_run = false
    }
  }
  
  action "deploy_systemd_service" {
    description = "Deploy systemd service using template"
    type = "template_deploy"
    
    template {
      source = "templates/app.service.tmpl"
      destination = "/etc/systemd/system/app.service"
      validate = true
      backup = true
    }
    
    targets = ["web-server-01"]
    
    execution {
      timeout = 120
      parallel = false
      dry_run = false
    }
  }
  
  action "deploy_docker_compose" {
    description = "Deploy Docker Compose using template"
    type = "template_deploy"
    
    template {
      source = "templates/docker-compose.yml.tmpl"
      destination = "/opt/app/docker-compose.yml"
      validate = true
      backup = true
    }
    
    targets = ["web-server-01"]
    
    execution {
      timeout = 120
      parallel = false
      dry_run = false
    }
  }
  
  action "deploy_env_file" {
    description = "Deploy environment file using template"
    type = "template_deploy"
    
    template {
      source = "templates/env.tmpl"
      destination = "/opt/app/.env"
      validate = true
      backup = true
    }
    
    targets = ["web-server-01"]
    
    execution {
      timeout = 120
      parallel = false
      dry_run = false
    }
  }
}
