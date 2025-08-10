# Machine Metadata in Templates

## Overview

Machine metadata is now available in Spooky templates, allowing you to access static configuration data about machines during template rendering. This provides a powerful way to create dynamic configurations based on machine properties.

## Machine Metadata Structure

Machine metadata is defined in your `machines.hcl` file and includes any key-value pairs you specify:

```hcl
machines {
  machine "web-server-01" {
    hostname = "web-server-01"
    host = "192.168.1.100"
    user = "admin"
    key_file = "~/.ssh/web-server-key"
    
    # Machine metadata
    metadata = {
      environment = "production"
      datacenter = "us-east-1"
      team = "platform"
      cost_center = "CC-12345"
      application = "web-app"
      version = "2.1.0"
      backup_schedule = "daily"
      monitoring_level = "high"
    }
    
    tags = {
      role = "web-server"
      tier = "frontend"
    }
  }
  
  machine "db-server-01" {
    hostname = "db-server-01"
    host = "192.168.1.101"
    user = "admin"
    key_file = "~/.ssh/db-server-key"
    
    metadata = {
      environment = "production"
      datacenter = "us-east-1"
      team = "database"
      cost_center = "CC-12346"
      application = "database"
      version = "13.4"
      backup_schedule = "hourly"
      monitoring_level = "critical"
      db_type = "postgresql"
      db_size_gb = "500"
    }
    
    tags = {
      role = "database"
      tier = "backend"
    }
  }
}
```

## Template Access

Machine metadata is available in templates through the `machines` array. Each machine object contains all its metadata:

### Basic Access

```go
{{range .machines}}
  Machine: {{.name}}
  Host: {{.host}}
  Environment: {{.metadata.environment}}
  Team: {{.metadata.team}}
{{end}}
```

### Conditional Rendering

```go
{{range .machines}}
  {{if eq .metadata.environment "production"}}
    # Production configuration for {{.name}}
    server {
      listen {{.host}}:80;
      server_name {{.name}}.{{.metadata.datacenter}}.example.com;
    }
  {{end}}
{{end}}
```

### Filtering by Metadata

```go
{{range .machines}}
  {{if eq .metadata.application "web-app"}}
    # Web application configuration
    upstream {{.name}} {
      server {{.host}}:8080;
    }
  {{end}}
{{end}}
```

### Complex Conditions

```go
{{range .machines}}
  {{if and (eq .metadata.environment "production") (eq .metadata.monitoring_level "critical")}}
    # Critical production monitoring
    monitor {{.name}} {
      interval = "30s";
      timeout = "10s";
      retries = 3;
    }
  {{end}}
{{end}}
```

## Complete Template Examples

### Nginx Configuration Template

```nginx
# nginx.conf.tmpl
events {
    worker_connections 1024;
}

http {
    # Global settings
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;
    
    # Server blocks for each web server
    {{range .machines}}
      {{if eq .metadata.application "web-app"}}
        server {
            listen 80;
            server_name {{.name}}.{{.metadata.datacenter}}.example.com;
            
            location / {
                proxy_pass http://{{.host}}:8080;
                proxy_set_header Host $host;
                proxy_set_header X-Real-IP $remote_addr;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto $scheme;
            }
            
            # Production-specific settings
            {{if eq .metadata.environment "production"}}
                access_log /var/log/nginx/{{.name}}.access.log;
                error_log /var/log/nginx/{{.name}}.error.log;
            {{end}}
        }
      {{end}}
    {{end}}
}
```

### Monitoring Configuration Template

```yaml
# monitoring.yml.tmpl
monitoring:
  global:
    scrape_interval: 15s
    evaluation_interval: 15s
  
  scrape_configs:
    {{range .machines}}
      - job_name: '{{.name}}'
        static_configs:
          - targets: ['{{.host}}:9100']
        metrics_path: /metrics
        scrape_interval: {{if eq .metadata.monitoring_level "critical"}}5s{{else}}15s{{end}}
        
        # Add labels from metadata
        relabel_configs:
          - source_labels: [__address__]
            target_label: instance
            regex: '([^:]+)(?::\d+)?'
            replacement: '${1}'
          - source_labels: [__address__]
            target_label: environment
            replacement: '{{.metadata.environment}}'
          - source_labels: [__address__]
            target_label: team
            replacement: '{{.metadata.team}}'
          - source_labels: [__address__]
            target_label: application
            replacement: '{{.metadata.application}}'
    {{end}}
```

### Backup Configuration Template

```bash
#!/bin/bash
# backup-script.sh.tmpl

# Backup script generated from machine metadata
{{range .machines}}
  {{if .metadata.backup_schedule}}
    # Backup configuration for {{.name}}
    echo "Configuring backup for {{.name}} ({{.metadata.application}})"
    
    {{if eq .metadata.application "database"}}
      # Database backup
      pg_dump -h {{.host}} -U {{.user}} {{.metadata.db_type}} > /backups/{{.name}}-$(date +%Y%m%d).sql
    {{else}}
      # Application backup
      rsync -avz {{.user}}@{{.host}}:{{.metadata.application_path}} /backups/{{.name}}-$(date +%Y%m%d)/
    {{end}}
    
    # Schedule based on metadata
    {{if eq .metadata.backup_schedule "hourly"}}
      echo "0 * * * * /usr/local/bin/backup-{{.name}}.sh" >> /etc/crontab
    {{else if eq .metadata.backup_schedule "daily"}}
      echo "0 2 * * * /usr/local/bin/backup-{{.name}}.sh" >> /etc/crontab
    {{end}}
  {{end}}
{{end}}
```

### Cost Allocation Template

```csv
# cost-allocation.csv.tmpl
Machine,Environment,Datacenter,Team,CostCenter,Application,MonthlyCost
{{range .machines}}
{{.name}},{{.metadata.environment}},{{.metadata.datacenter}},{{.metadata.team}},{{.metadata.cost_center}},{{.metadata.application}},{{if eq .metadata.environment "production"}}$100{{else}}$50{{end}}
{{end}}
```

## Advanced Usage Patterns

### Grouping by Metadata

```go
{{/* Group machines by environment */}}
{{range $env := (groupBy .machines "metadata.environment")}}
  # Environment: {{$env.Key}}
  {{range $env.Values}}
    - {{.name}} ({{.metadata.application}})
  {{end}}
{{end}}
```

### Aggregation

```go
{{/* Count machines by team */}}
{{range $team := (groupBy .machines "metadata.team")}}
  Team {{$team.Key}}: {{len $team.Values}} machines
{{end}}
```

### Dynamic Configuration

```go
{{/* Generate configuration based on machine capabilities */}}
{{range .machines}}
  {{if .metadata.db_size_gb}}
    # Database with {{.metadata.db_size_gb}}GB storage
    database_config {
      max_connections = {{if gt (atoi .metadata.db_size_gb) 100}}200{{else}}50{{end}};
      shared_buffers = {{mul (atoi .metadata.db_size_gb) 25}}MB;
    }
  {{end}}
{{end}}
```

## Best Practices

### 1. Consistent Metadata Keys

Use consistent metadata keys across all machines:

```hcl
metadata = {
  environment = "production|staging|development"
  team = "platform|database|security|application"
  datacenter = "us-east-1|us-west-2|eu-west-1"
  cost_center = "CC-XXXXX"
  application = "web-app|api|database|cache"
  version = "semver"
}
```

### 2. Validation

Validate metadata in your schema:

```hcl
metadata = {
  type = "object"
  required = false
  description = "Machine metadata for templates"
  additional_properties = "string"
  validation = {
    required_keys = ["environment", "team", "application"]
    allowed_values = {
      environment = ["production", "staging", "development"]
      team = ["platform", "database", "security", "application"]
    }
  }
}
```

### 3. Documentation

Document your metadata schema:

```hcl
# Machine metadata schema
# environment: The deployment environment (production, staging, development)
# team: The team responsible for this machine
# datacenter: The physical datacenter location
# cost_center: The cost center for billing
# application: The primary application running on this machine
# version: The application version
# backup_schedule: Backup frequency (hourly, daily, weekly)
# monitoring_level: Monitoring intensity (low, medium, high, critical)
```

### 4. Template Organization

Organize templates by metadata:

```
templates/
├── production/
│   ├── nginx.conf.tmpl
│   └── monitoring.yml.tmpl
├── staging/
│   ├── nginx.conf.tmpl
│   └── monitoring.yml.tmpl
└── shared/
    ├── backup-script.sh.tmpl
    └── cost-allocation.csv.tmpl
```

## Migration from Tags

If you're currently using tags, you can migrate to metadata:

### Before (Tags)
```hcl
tags = {
  environment = "production"
  team = "platform"
  application = "web-app"
}
```

### After (Metadata)
```hcl
metadata = {
  environment = "production"
  team = "platform"
  application = "web-app"
}
```

### Template Access
```go
{{/* Old way with tags */}}
{{if eq .tags.environment "production"}}

{{/* New way with metadata */}}
{{if eq .metadata.environment "production"}}
```

## Troubleshooting

### Common Issues

1. **Metadata not available**: Ensure machines are loaded in template context
2. **Nil pointer errors**: Check if metadata exists before accessing
3. **Type mismatches**: Metadata values are strings, convert as needed

### Debug Template

```go
{{/* Debug template to see all available data */}}
{{range .machines}}
  Machine: {{.name}}
  Host: {{.host}}
  Metadata: {{printf "%+v" .metadata}}
  Tags: {{printf "%+v" .tags}}
{{end}}
```

## Conclusion

Machine metadata provides a powerful way to create dynamic, environment-aware configurations in Spooky templates. By leveraging metadata, you can create templates that automatically adapt to your infrastructure's characteristics and requirements.
