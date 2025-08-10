# Machine Metadata Template Test

This example demonstrates how machine metadata can be used in Spooky templates.

## Project Structure

```
machine-metadata-template-test/
├── machines.hcl
├── templates/
│   ├── nginx.conf.tmpl
│   ├── monitoring.yml.tmpl
│   └── debug.tmpl
└── README.md
```

## Usage

1. Create the project directory structure
2. Add your machine definitions with metadata
3. Create templates that use the metadata
4. Render templates with: `spooky templates render <project> <template> --output <output_file>`

## Example Commands

```bash
# Render nginx configuration
spooky templates render . templates/nginx.conf.tmpl --output nginx.conf

# Render monitoring configuration  
spooky templates render . templates/monitoring.yml.tmpl --output monitoring.yml

# Debug template to see all available data
spooky templates render . templates/debug.tmpl --output debug.txt
```

## Machine Metadata Examples

The `machines.hcl` file should contain machine definitions with metadata:

```hcl
machines {
  machine "web-server-01" {
    hostname = "web-server-01"
    host = "192.168.1.100"
    user = "admin"
    
    metadata = {
      environment = "production"
      datacenter = "us-east-1"
      team = "platform"
      application = "web-app"
      version = "2.1.0"
      monitoring_level = "high"
    }
  }
  
  machine "db-server-01" {
    hostname = "db-server-01"
    host = "192.168.1.101"
    user = "admin"
    
    metadata = {
      environment = "production"
      datacenter = "us-east-1"
      team = "database"
      application = "database"
      version = "13.4"
      monitoring_level = "critical"
      db_type = "postgresql"
    }
  }
}
```

## Template Examples

### nginx.conf.tmpl
```nginx
events {
    worker_connections 1024;
}

http {
    {{range .machines}}
      {{if eq .metadata.application "web-app"}}
        server {
            listen 80;
            server_name {{.name}}.{{.metadata.datacenter}}.example.com;
            
            location / {
                proxy_pass http://{{.host}}:8080;
            }
            
            {{if eq .metadata.environment "production"}}
                access_log /var/log/nginx/{{.name}}.access.log;
                error_log /var/log/nginx/{{.name}}.error.log;
            {{end}}
        }
      {{end}}
    {{end}}
}
```

### monitoring.yml.tmpl
```yaml
monitoring:
  scrape_configs:
    {{range .machines}}
      - job_name: '{{.name}}'
        static_configs:
          - targets: ['{{.host}}:9100']
        scrape_interval: {{if eq .metadata.monitoring_level "critical"}}5s{{else}}15s{{end}}
        
        relabel_configs:
          - source_labels: [__address__]
            target_label: environment
            replacement: '{{.metadata.environment}}'
          - source_labels: [__address__]
            target_label: team
            replacement: '{{.metadata.team}}'
    {{end}}
```

### debug.tmpl
```go
# Debug template to see all available machine data
{{range .machines}}
Machine: {{.name}}
Host: {{.host}}
User: {{.user}}
Metadata: {{printf "%+v" .metadata}}
Tags: {{printf "%+v" .tags}}
---
{{end}}
```
