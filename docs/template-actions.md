# Template Actions

Template actions enable server-side template deployment and evaluation, allowing you to generate dynamic configuration files on target servers using server-specific facts.

## Overview

Template actions provide four main operations:

1. **`template_deploy`** - Deploy template files to target servers
2. **`template_evaluate`** - Evaluate templates on servers with server-specific facts
3. **`template_validate`** - Validate templates on servers
4. **`template_cleanup`** - Remove template files from servers

## Template Action Types

### template_deploy

Deploys template files to target servers without evaluation.

```hcl
action "deploy-config-template" {
  description = "Deploy configuration template to servers"
  type = "template_deploy"
  
  template {
    source = "templates/nginx.conf.tmpl"
    destination = "/tmp/nginx.conf.tmpl"
    permissions = "644"
    owner = "root"
    group = "root"
  }
  
  machines = ["web-server-1", "web-server-2"]
  parallel = true
  timeout = 300
}
```

### template_evaluate

Evaluates templates on target servers using server-specific facts and writes the result to a destination file.

```hcl
action "evaluate-config" {
  description = "Evaluate configuration template on servers"
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
  parallel = true
  timeout = 300
}
```

### template_validate

Validates templates on target servers without executing them.

```hcl
action "validate-template" {
  description = "Validate template syntax on servers"
  type = "template_validate"
  
  template {
    source = "/tmp/nginx.conf.tmpl"
  }
  
  machines = ["web-server-1", "web-server-2"]
  parallel = true
  timeout = 120
}
```

### template_cleanup

Removes template files from target servers.

```hcl
action "cleanup-templates" {
  description = "Remove template files from servers"
  type = "template_cleanup"
  
  template {
    source = "/tmp/nginx.conf.tmpl"
  }
  
  machines = ["web-server-1", "web-server-2"]
  parallel = true
  timeout = 60
}
```

## Template Configuration

The `template` block supports the following options:

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `source` | string | Yes | Source template file path (local for deploy, remote for others) |
| `destination` | string | Yes | Destination file path on target server |
| `backup` | bool | No | Create backup of existing file before overwriting |
| `validate` | bool | No | Validate the generated file after creation |
| `permissions` | string | No | File permissions (e.g., "644", "755") |
| `owner` | string | No | File owner |
| `group` | string | No | File group |

## Server-Side Template Functions

When evaluating templates on servers, the following functions are available:

### System Information Functions

- `{{machineID}}` - Get machine ID from `/etc/machine-id`
- `{{osVersion}}` - Get OS version from `uname -r`
- `{{hostname}}` - Get server hostname
- `{{ipAddress}}` - Get primary IP address
- `{{diskSpace}}` - Get available disk space
- `{{memoryInfo}}` - Get memory information

### File System Functions

- `{{fileExists 'path'}}` - Check if file exists (returns bool)
- `{{fileContent 'path'}}` - Read file content
- `{{fileSize 'path'}}` - Get file size
- `{{fileOwner 'path'}}` - Get file owner

## Example Templates

### Nginx Configuration Template

```nginx
# Nginx configuration for {{hostname}}
# Generated on {{osVersion}} with {{memoryInfo}} RAM

server {
    listen 80;
    server_name {{hostname}} {{ipAddress}};
    
    root /var/www/html;
    index index.html index.php;
    
    access_log /var/log/nginx/{{hostname}}_access.log;
    error_log /var/log/nginx/{{hostname}}_error.log;
    
    location ~ \.php$ {
        fastcgi_pass unix:/var/run/php/php7.4-fpm.sock;
        fastcgi_index index.php;
        include fastcgi_params;
    }
}

# Machine ID: {{machineID}}
# Available disk space: {{diskSpace}}
```

### PHP Configuration Template

```ini
; PHP Configuration for {{hostname}}
; Generated on {{osVersion}}

[PHP]
memory_limit = 256M
max_execution_time = 300
error_log = /var/log/php/{{hostname}}_error.log

; SSL certificate check
{{if fileExists "/etc/ssl/certs/ssl-cert-snakeoil.pem"}}
; SSL certificate available
{{else}}
; No SSL certificate found
{{end}}
```

## Usage Workflow

### 1. Deploy Template

First, deploy the template file to target servers:

```bash
spooky action run --project . deploy-nginx-template
```

### 2. Evaluate Template

Evaluate the template on each server with server-specific facts:

```bash
spooky action run --project . evaluate-nginx-config
```

### 3. Validate Template (Optional)

Validate the template syntax before evaluation:

```bash
spooky action run --project . validate-nginx-template
```

### 4. Clean Up (Optional)

Remove template files after deployment:

```bash
spooky action run --project . cleanup-nginx-template
```

## Complete Example

```hcl
actions {
  # Deploy nginx template
  action "deploy-nginx-template" {
    description = "Deploy nginx configuration template"
    type = "template_deploy"
    
    template {
      source = "templates/nginx.conf.tmpl"
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

  # Evaluate nginx template
  action "evaluate-nginx-config" {
    description = "Evaluate nginx configuration template"
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
}
```

## Benefits

- **Dynamic Configuration**: Generate server-specific configurations
- **Fact Integration**: Use server facts in template evaluation
- **Idempotent**: Safe to run multiple times
- **Validation**: Built-in template and file validation
- **Backup Support**: Automatic backup of existing files
- **Parallel Execution**: Deploy to multiple servers simultaneously

## Security Considerations

- Template functions have limited access to server information
- File operations are restricted to specified paths
- Template evaluation is sandboxed
- All operations are logged for audit purposes 