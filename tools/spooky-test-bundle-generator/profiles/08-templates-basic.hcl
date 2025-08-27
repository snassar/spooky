profile "templates_basic" "Basic template rendering testing across all operating systems" {
  description = "Tests basic template rendering functionality with simple variable substitution across all operating systems"
  
  containers {
    debian13 {
      base_image = "debian:trixie-slim"
      static_ip = "10.0.100.81"
      ssh_port = 2292
      packages = ["openssh-server", "openssh-client", "curl", "wget", "nginx"]
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
      static_ip = "10.0.100.82"
      ssh_port = 2293
      packages = ["openssh-server", "openssh-clients", "curl", "wget", "nginx"]
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
      static_ip = "10.0.100.83"
      ssh_port = 2294
      packages = ["openssh", "curl", "wget", "nginx"]
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
      static_ip = "10.0.100.84"
      ssh_port = 2295
      packages = ["openssh-server", "openssh-client", "curl", "wget", "nginx"]
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
      static_ip = "10.0.100.85"
      ssh_port = 2296
      packages = ["openssh-server", "openssh-clients", "curl", "wget", "nginx"]
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
    name = "templates-basic-test"
    description = "Cross-OS basic template rendering test"
    
    machines {
      debian13 {
        hostname = "debian13-templates"
        ip = "10.0.100.81"
        port = 2292
        auth {
          password = "testpass123"
        }
        facts {
          basic = true
          enhanced = false
          custom = false
          encrypted = false
        }
      }
      
      fedora42 {
        hostname = "fedora42-templates"
        ip = "10.0.100.82"
        port = 2293
        auth {
          password = "testpass123"
        }
        facts {
          basic = true
          enhanced = false
          custom = false
          encrypted = false
        }
      }
      
      arch {
        hostname = "arch-templates"
        ip = "10.0.100.83"
        port = 2294
        auth {
          password = "testpass123"
        }
        facts {
          basic = true
          enhanced = false
          custom = false
          encrypted = false
        }
      }
      
      alpine319 {
        hostname = "alpine319-templates"
        ip = "10.0.100.84"
        port = 2295
        auth {
          password = "testpass123"
        }
        facts {
          basic = true
          enhanced = false
          custom = false
          encrypted = false
        }
      }
      
      opensuse156 {
        hostname = "opensuse156-templates"
        ip = "10.0.100.85"
        port = 2296
        auth {
          password = "testpass123"
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
      test_type = "templates-basic"
      os_list = ["debian13", "fedora42", "arch", "alpine319", "opensuse156"]
      app_name = "test-app"
      app_version = "1.0.0"
      environment = "testing"
      max_connections = 100
      log_level = "info"
    }
    
    templates {
      nginx_config {
        name = "nginx.conf.tmpl"
        path = "/etc/nginx/nginx.conf"
        content = <<-EOF
user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log {{ .log_level }};
pid /run/nginx.pid;

events {
    worker_connections {{ .max_connections }};
}

http {
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';

    access_log /var/log/nginx/access.log main;

    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;

    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    server {
        listen 80;
        server_name {{ .hostname }};
        root /var/www/html;
        index index.html index.htm;

        location / {
            try_files $uri $uri/ =404;
        }

        location /status {
            access_log off;
            return 200 "{{ .app_name }} v{{ .app_version }} on {{ .os }} is running\n";
            add_header Content-Type text/plain;
        }
    }
}
EOF
        permissions = "644"
        owner = "root"
        group = "root"
      }
      
      app_config {
        name = "app.conf.tmpl"
        path = "/etc/{{ .app_name }}/app.conf"
        content = <<-EOF
# Application configuration for {{ .app_name }}
[app]
name = {{ .app_name }}
version = {{ .app_version }}
environment = {{ .environment }}

[server]
host = {{ .ip }}
port = 8080
max_connections = {{ .max_connections }}

[logging]
level = {{ .log_level }}
file = /var/log/{{ .app_name }}/app.log

[system]
os = {{ .os }}
hostname = {{ .hostname }}
architecture = {{ .architecture }}
EOF
        permissions = "644"
        owner = "root"
        group = "root"
      }
      
      systemd_service {
        name = "app.service.tmpl"
        path = "/etc/systemd/system/{{ .app_name }}.service"
        content = <<-EOF
[Unit]
Description={{ .app_name }} Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/{{ .app_name }}
ExecStart=/usr/bin/{{ .app_name }} --config /etc/{{ .app_name }}/app.conf
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
        permissions = "644"
        owner = "root"
        group = "root"
      }
    }
    
    actions {
      render_templates {
        name = "render-templates"
        description = "Render all templates on all machines"
        command = "template --render-all"
        tags = ["templates", "rendering", "cross-os"]
        parallel = true
        timeout = "60s"
      }
      
      verify_templates {
        name = "verify-templates"
        description = "Verify that templates are rendered correctly"
        command = "template --verify"
        tags = ["templates", "verification", "cross-os"]
        depends_on = ["render_templates"]
        timeout = "30s"
      }
      
      test_template_variables {
        name = "test-template-variables"
        description = "Test template variable substitution"
        command = "template --test-variables"
        tags = ["templates", "variables", "cross-os"]
        depends_on = ["verify_templates"]
        timeout = "20s"
      }
      
      deploy_templates {
        name = "deploy-templates"
        description = "Deploy rendered templates to target locations"
        command = "template --deploy"
        tags = ["templates", "deployment", "cross-os"]
        depends_on = ["verify_templates"]
        parallel = true
        timeout = "45s"
      }
      
      restart_services {
        name = "restart-services"
        description = "Restart services after template deployment"
        command = "action --restart-services"
        tags = ["templates", "services", "cross-os"]
        depends_on = ["deploy_templates"]
        parallel = true
        timeout = "30s"
      }
    }
  }
}
