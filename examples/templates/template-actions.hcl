# Template Actions Example
# Demonstrates server-side template deployment and evaluation

actions {
  # Deploy nginx template to servers
  action "deploy-nginx-template" {
    description = "Deploy nginx configuration template to web servers"
    type = "template_deploy"
    
    template {
      source = "examples/templates/nginx.conf.tmpl"
      destination = "/tmp/nginx.conf.tmpl"
      permissions = "644"
      owner = "root"
      group = "root"
    }
    
    machines = ["web-server-1", "web-server-2"]
    tags = ["web", "nginx"]
    parallel = true
    timeout = 300
  }

  # Evaluate nginx template on servers
  action "evaluate-nginx-config" {
    description = "Evaluate nginx configuration template on servers"
    type = "template_evaluate"
    
    template {
      source = "/tmp/nginx.conf.tmpl"
      destination = "/etc/nginx/sites-available/default"
      backup = true
      validate = true
      permissions = "644"
      owner = "root"
      group = "root"
    }
    
    machines = ["web-server-1", "web-server-2"]
    tags = ["web", "nginx"]
    parallel = true
    timeout = 300
  }

  # Validate nginx template on servers
  action "validate-nginx-template" {
    description = "Validate nginx configuration template on servers"
    type = "template_validate"
    
    template {
      source = "/tmp/nginx.conf.tmpl"
    }
    
    machines = ["web-server-1", "web-server-2"]
    tags = ["web", "nginx"]
    parallel = true
    timeout = 120
  }

  # Clean up template files
  action "cleanup-nginx-template" {
    description = "Remove nginx template files from servers"
    type = "template_cleanup"
    
    template {
      source = "/tmp/nginx.conf.tmpl"
    }
    
    machines = ["web-server-1", "web-server-2"]
    tags = ["web", "nginx"]
    parallel = true
    timeout = 60
  }

  # Deploy PHP configuration template
  action "deploy-php-config" {
    description = "Deploy PHP configuration template"
    type = "template_deploy"
    
    template {
      source = "examples/templates/php.ini.tmpl"
      destination = "/tmp/php.ini.tmpl"
      permissions = "644"
    }
    
    machines = ["web-server-1", "web-server-2"]
    tags = ["web", "php"]
    parallel = true
    timeout = 300
  }

  # Evaluate PHP configuration
  action "evaluate-php-config" {
    description = "Evaluate PHP configuration template"
    type = "template_evaluate"
    
    template {
      source = "/tmp/php.ini.tmpl"
      destination = "/etc/php/7.4/fpm/php.ini"
      backup = true
      validate = true
    }
    
    machines = ["web-server-1", "web-server-2"]
    tags = ["web", "php"]
    parallel = true
    timeout = 300
  }
} 