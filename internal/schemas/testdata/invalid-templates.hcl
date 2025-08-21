# Invalid Templates Configuration
# This file should fail schema validations

metadata {
  schema_version = "invalid-version"  # Invalid ScalVer format
  schema_type = "templates"
  schema_name = "Test Templates Configuration"
  last_updated = "2024-01-01"
  compatibility = ["0.20250809.0"]
  description = "A test templates configuration that should fail validations"
}

templates {
  # Invalid template (missing required fields)
  template "invalid_template" {
    description = "Invalid template with missing required fields"
    # Missing source and destination
    # Missing variables
    
    validation {
      syntax_check = true
      test_nginx_config = true
      backup_existing = true
    }
    
    permissions {
      owner = "root"
      group = "root"
      mode = "invalid-mode"  # Invalid: not a valid file mode
    }
    
    metadata {
      author = "DevOps Team"
      version = "invalid-version"  # Invalid: not semantic version
      tags = ["nginx", "web", "configuration"]
    }
  }
  
  # Template with invalid source path
  template "invalid_source" {
    description = "Template with invalid source path"
    source = "invalid/path/with/spaces.tmpl"  # Invalid: contains spaces
    destination = "/etc/nginx/nginx.conf"
    
    variables = {
      server_name = "{{ .server_name }}"
      port = "{{ .port }}"
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
  }
  
  # Template with invalid destination path
  template "invalid_destination" {
    description = "Template with invalid destination path"
    source = "templates/nginx.conf.tmpl"
    destination = "relative/path.conf"  # Invalid: not absolute path
    
    variables = {
      server_name = "{{ .server_name }}"
      port = "{{ .port }}"
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
  }
  
  # Template with invalid permissions
  template "invalid_permissions" {
    description = "Template with invalid permissions"
    source = "templates/app.conf.tmpl"
    destination = "/opt/app/config/app.conf"
    
    variables = {
      db_host = "{{ .db_host }}"
      db_port = "{{ .db_port }}"
    }
    
    validation {
      syntax_check = true
      test_config = true
      backup_existing = true
    }
    
    permissions {
      owner = "invalid-user"  # Invalid: contains invalid characters
      group = "invalid-group"  # Invalid: contains invalid characters
      mode = "9999"  # Invalid: not a valid file mode
    }
  }
  
  # Template with invalid metadata
  template "invalid_metadata" {
    description = "Template with invalid metadata"
    source = "templates/app.service.tmpl"
    destination = "/etc/systemd/system/app.service"
    
    variables = {
      app_name = "{{ .app_name }}"
      app_path = "{{ .app_path }}"
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
      author = ""  # Invalid: empty author
      version = "not-a-version"  # Invalid: not semantic version
      tags = []  # Invalid: empty tags list
    }
  }
}
