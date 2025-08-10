# Maintenance Windows and Machine Targeting

## Overview

Maintenance windows in Spooky can affect machine targeting and action execution. This document explains how maintenance windows work and how to control their behavior.

## How Maintenance Windows Affect Targeting

### **1. Status-Based Exclusion (Default Behavior)**

Machines with `status = "maintenance"` are **automatically excluded** from actions:

```hcl
machine "web-prod-01" {
  # ... other configuration ...
  
  lifecycle = {
    status = "maintenance"  # This machine will be excluded from actions
    maintenance_window = {
      start_time = "02:00"
      end_time = "04:00"
      timezone = "America/New_York"
      days_of_week = ["sunday"]
    }
  }
}
```

**Behavior**: This machine will be skipped when running actions, even if it matches the target criteria.

### **2. Time-Based Auto-Exclusion (Optional)**

Machines can be automatically excluded during their maintenance windows:

```hcl
machine "web-prod-02" {
  # ... other configuration ...
  
  lifecycle = {
    status = "active"  # Machine is normally active
    maintenance_window = {
      start_time = "02:00"
      end_time = "04:00"
      timezone = "America/New_York"
      days_of_week = ["sunday"]
      auto_exclude = true  # Exclude during maintenance window
    }
  }
}
```

**Behavior**: This machine will be automatically excluded during its maintenance window, even with `status = "active"`.

### **3. Maintenance Window as Metadata Only**

Maintenance windows can be purely informational:

```hcl
machine "web-prod-03" {
  # ... other configuration ...
  
  lifecycle = {
    status = "active"
    maintenance_window = {
      start_time = "02:00"
      end_time = "04:00"
      timezone = "America/New_York"
      days_of_week = ["sunday"]
      auto_exclude = false  # Don't exclude during maintenance window
    }
  }
}
```

**Behavior**: This machine will always be included in actions, regardless of the current time.

## Action-Level Control

### **Respect Maintenance Windows**

Actions can be configured to respect or ignore maintenance windows:

```hcl
action "update_nginx" {
  type = "command"
  command = "apt-get update && apt-get upgrade -y nginx"
  machines = ["web-servers"]
  
  # Maintenance behavior options
  maintenance_behavior = "skip"  # Skip machines in maintenance
  # OR
  maintenance_behavior = "warn"  # Warn but continue
  # OR  
  maintenance_behavior = "force" # Ignore maintenance status
}
```

### **Emergency Override**

Force actions to run even on maintenance machines:

```bash
# Override maintenance status for emergency actions
spooky actions run --force-maintenance update_nginx /path/to/project
```

## Examples

### **Example 1: Scheduled Maintenance**

```hcl
machines {
  machine "web-prod-01" {
    hostname = "web-prod-01"
    host = "192.168.1.100"
    user = "admin"
    
    lifecycle = {
      status = "maintenance"  # Currently in maintenance
      maintenance_window = {
        start_time = "02:00"
        end_time = "04:00"
        timezone = "America/New_York"
        days_of_week = ["sunday"]
      }
    }
  }
  
  machine "web-prod-02" {
    hostname = "web-prod-02"
    host = "192.168.1.101"
    user = "admin"
    
    lifecycle = {
      status = "active"  # Available for actions
      maintenance_window = {
        start_time = "02:00"
        end_time = "04:00"
        timezone = "America/New_York"
        days_of_week = ["sunday"]
      }
    }
  }
}
```

**Result**: Only `web-prod-02` will be targeted by actions.

### **Example 2: Auto-Exclusion During Maintenance Window**

```hcl
machines {
  machine "web-prod-01" {
    hostname = "web-prod-01"
    host = "192.168.1.100"
    user = "admin"
    
    lifecycle = {
      status = "active"
      maintenance_window = {
        start_time = "02:00"
        end_time = "04:00"
        timezone = "America/New_York"
        days_of_week = ["sunday"]
        auto_exclude = true  # Auto-exclude during window
      }
    }
  }
}
```

**Result**: Machine will be excluded during Sunday 2-4 AM Eastern Time.

### **Example 3: Emergency Override**

```bash
# Force action to run on all machines, including maintenance
spooky actions run --force-maintenance emergency_patch /path/to/project

# Or in action definition
action "emergency_patch" {
  type = "command"
  command = "apply_critical_security_patch"
  machines = ["all-servers"]
  maintenance_behavior = "force"  # Ignore maintenance status
}
```

## CLI Commands and Maintenance Windows

### **List Machines with Maintenance Status**

```bash
# Show all machines with their maintenance status
spooky machines list /path/to/project --show-maintenance

# Show only machines available for actions
spooky machines list /path/to/project --available-only

# Show only machines in maintenance
spooky machines list /path/to/project --maintenance-only
```

### **Check Maintenance Windows**

```bash
# Check which machines are in maintenance windows now
spooky machines maintenance-check /path/to/project

# Check maintenance windows for specific time
spooky machines maintenance-check /path/to/project --time "2024-01-15T03:00:00Z"
```

## Template Usage

### **Filter Machines by Maintenance Status**

```go
{{/* Only target active machines */}}
{{range .machines}}
  {{if eq .lifecycle.status "active"}}
    # Configuration for {{.name}}
    server {
      listen {{.host}}:80;
    }
  {{end}}
{{end}}
```

### **Show Maintenance Information**

```go
{{/* Show maintenance schedule */}}
{{range .machines}}
  {{if .lifecycle.maintenance_window}}
    ## {{.name}} Maintenance
    **Status**: {{.lifecycle.status}}
    **Window**: {{.lifecycle.maintenance_window.start_time}}-{{.lifecycle.maintenance_window.end_time}} {{.lifecycle.maintenance_window.timezone}}
    **Days**: {{join .lifecycle.maintenance_window.days_of_week ", "}}
    {{if .lifecycle.maintenance_window.auto_exclude}}
    **Auto-exclude**: Yes
    {{end}}
  {{end}}
{{end}}
```

## Best Practices

### **1. Use Status for Manual Control**

Set `status = "maintenance"` when you want to manually exclude a machine:

```hcl
lifecycle = {
  status = "maintenance"  # Manual exclusion
}
```

### **2. Use Auto-Exclude for Scheduled Maintenance**

Use `auto_exclude = true` for regular maintenance windows:

```hcl
maintenance_window = {
  start_time = "02:00"
  end_time = "04:00"
  timezone = "America/New_York"
  days_of_week = ["sunday"]
  auto_exclude = true  # Automatic exclusion
}
```

### **3. Provide Emergency Override**

Always have a way to override maintenance status for emergencies:

```bash
spooky actions run --force-maintenance emergency_action /path/to/project
```

### **4. Document Maintenance Procedures**

Include maintenance procedures in your documentation:

```hcl
# Maintenance procedure:
# 1. Set status = "maintenance"
# 2. Perform maintenance
# 3. Set status = "active"
# 4. Verify machine is back in rotation
```

## Summary

- **Status-based exclusion**: Manual control via `status = "maintenance"`
- **Time-based exclusion**: Automatic via `auto_exclude = true`
- **Action-level control**: Via `maintenance_behavior` in actions
- **Emergency override**: Via `--force-maintenance` flag
- **Template filtering**: Filter machines by status in templates

This provides flexible control over when machines are available for actions while maintaining safety and emergency override capabilities.
