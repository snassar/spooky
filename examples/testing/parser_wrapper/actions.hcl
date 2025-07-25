# Actions for parser_wrapper project
# Add your action definitions here

actions {
  action "check-status" {
    description = "Check server status"
    command = "uptime && df -h"
    tags = ["role=web", "role=database"]
    parallel = true
    timeout = 300
  }
  
  action "update-system" {
    description = "Update system packages"
    command = "apt update && apt upgrade -y"
    tags = ["environment=testing"]
    parallel = true
    timeout = 600
  }
  
  action "deploy-nginx-config" {
    description = "Deploy nginx configuration template"
    type = "template_deploy"
    template {
      source = "templates/nginx.conf.tmpl"
      destination = "/tmp/nginx.conf.tmpl"
      backup = true
      permissions = "644"
      owner = "root"
      group = "root"
    }
    tags = ["role=web"]
  }
  
  action "evaluate-nginx-config" {
    description = "Evaluate nginx configuration on servers"
    type = "template_evaluate"
    template {
      source = "/tmp/nginx.conf.tmpl"
      destination = "/etc/nginx/nginx.conf"
      validate = true
      backup = true
    }
    tags = ["role=web"]
  }
  
  action "validate-nginx-config" {
    description = "Validate nginx configuration"
    type = "template_validate"
    template {
      source = "/tmp/nginx.conf.tmpl"
      destination = "/etc/nginx/nginx.conf"
    }
    tags = ["role=web"]
  }
  
  action "deploy-php-config" {
    description = "Deploy PHP configuration template"
    type = "template_deploy"
    template {
      source = "templates/php.ini.tmpl"
      destination = "/tmp/php.ini.tmpl"
      backup = true
      permissions = "644"
    }
    tags = ["role=web"]
  }
  
  action "cleanup-templates" {
    description = "Clean up template files"
    type = "template_cleanup"
    template {
      source = "/tmp/nginx.conf.tmpl"
      destination = "/tmp/php.ini.tmpl"
    }
    tags = ["role=web"]
  }
  
  action "database-backup" {
    description = "Create database backup"
    script = "scripts/backup-db.sh"
    tags = ["role=database"]
    timeout = 1800
  }
  
  action "restart-services" {
    description = "Restart application services"
    command = "systemctl restart nginx && systemctl restart php-fpm"
    tags = ["role=web"]
    parallel = false
  }
}