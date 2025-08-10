# Global Maintenance Windows for Distributed Teams

## Overview

Spooky supports timezone-aware maintenance windows to accommodate global teams working across different time zones. This ensures that maintenance activities are scheduled at appropriate times for all team members.

## Timezone-Aware Maintenance Windows

### Basic Structure

```hcl
maintenance_window = {
  start_time = "02:00"
  end_time = "04:00"
  timezone = "America/New_York"
  days_of_week = ["sunday"]
}
```

### Global Team Examples

#### US-East Coast Team
```hcl
machine "web-prod-01" {
  # ... other configuration ...
  
  lifecycle = {
    status = "active"
    maintenance_team = "platform-us-east"
    team_timezone = "America/New_York"
    
    maintenance_window = {
      start_time = "02:00"
      end_time = "04:00"
      timezone = "America/New_York"
      days_of_week = ["sunday"]
    }
    
    team_contact = {
      primary = "platform-east@company.com"
      secondary = "+1-555-0123"
      slack_channel = "#platform-maintenance"
      pagerduty_schedule = "P123456"
    }
  }
}
```

#### Europe Team
```hcl
machine "web-prod-02" {
  # ... other configuration ...
  
  lifecycle = {
    status = "active"
    maintenance_team = "platform-europe"
    team_timezone = "Europe/London"
    
    maintenance_window = {
      start_time = "03:00"
      end_time = "05:00"
      timezone = "Europe/London"
      days_of_week = ["sunday"]
    }
    
    team_contact = {
      primary = "platform-europe@company.com"
      secondary = "+44-20-7946-0958"
      slack_channel = "#platform-maintenance-eu"
      pagerduty_schedule = "P789012"
    }
  }
}
```

#### Asia-Pacific Team
```hcl
machine "web-prod-03" {
  # ... other configuration ...
  
  lifecycle = {
    status = "active"
    maintenance_team = "platform-apac"
    team_timezone = "Asia/Tokyo"
    
    maintenance_window = {
      start_time = "01:00"
      end_time = "03:00"
      timezone = "Asia/Tokyo"
      days_of_week = ["sunday"]
    }
    
    team_contact = {
      primary = "platform-apac@company.com"
      secondary = "+81-3-1234-5678"
      slack_channel = "#platform-maintenance-apac"
      pagerduty_schedule = "P345678"
    }
  }
}
```

## Timezone Conversion Examples

### Same Time, Different Timezones

When you want maintenance to happen at the same local time across regions:

```hcl
# All teams do maintenance at 2 AM their local time on Sunday

# US East Coast
maintenance_window = {
  start_time = "02:00"
  end_time = "04:00"
  timezone = "America/New_York"
  days_of_week = ["sunday"]
}

# Europe (7 hours ahead of US East Coast)
maintenance_window = {
  start_time = "02:00"
  end_time = "04:00"
  timezone = "Europe/London"
  days_of_week = ["sunday"]
}

# Asia (14 hours ahead of US East Coast)
maintenance_window = {
  start_time = "02:00"
  end_time = "04:00"
  timezone = "Asia/Tokyo"
  days_of_week = ["sunday"]
}
```

### Coordinated Global Maintenance

When you want all teams to participate in the same maintenance window:

```hcl
# Global maintenance at 2 AM UTC on Sunday (coordinated)

# US East Coast (UTC-5, so 9 PM Saturday local time)
maintenance_window = {
  start_time = "21:00"
  end_time = "23:00"
  timezone = "America/New_York"
  days_of_week = ["saturday"]
}

# Europe (UTC+0, so 2 AM Sunday local time)
maintenance_window = {
  start_time = "02:00"
  end_time = "04:00"
  timezone = "Europe/London"
  days_of_week = ["sunday"]
}

# Asia (UTC+9, so 11 AM Sunday local time)
maintenance_window = {
  start_time = "11:00"
  end_time = "13:00"
  timezone = "Asia/Tokyo"
  days_of_week = ["sunday"]
}
```

## Template Usage

### Maintenance Schedule Template

```go
{{/* maintenance-schedule.tmpl */}}
# Maintenance Schedule for {{.name}}

## Machine Information
- **Name**: {{.name}}
- **Environment**: {{.metadata.environment}}
- **Team**: {{.lifecycle.maintenance_team}}

## Maintenance Window
{{if .lifecycle.maintenance_window.simple}}
**Simple Window**: {{.lifecycle.maintenance_window.simple}} (local timezone)
{{else}}
**Timezone-Aware Window**:
- **Start**: {{.lifecycle.maintenance_window.start_time}} {{.lifecycle.maintenance_window.timezone}}
- **End**: {{.lifecycle.maintenance_window.end_time}} {{.lifecycle.maintenance_window.timezone}}
- **Days**: {{join .lifecycle.maintenance_window.days_of_week ", "}}
{{end}}

## Team Contact
{{if .lifecycle.team_contact}}
- **Primary**: {{.lifecycle.team_contact.primary}}
- **Secondary**: {{.lifecycle.team_contact.secondary}}
- **Slack**: {{.lifecycle.team_contact.slack_channel}}
- **PagerDuty**: {{.lifecycle.team_contact.pagerduty_schedule}}
{{end}}

## Timezone Conversions
{{if .lifecycle.maintenance_window.timezone}}
**Maintenance window in different timezones**:
- **UTC**: {{convertTime .lifecycle.maintenance_window.start_time .lifecycle.maintenance_window.timezone "UTC"}}
- **US East Coast**: {{convertTime .lifecycle.maintenance_window.start_time .lifecycle.maintenance_window.timezone "America/New_York"}}
- **Europe**: {{convertTime .lifecycle.maintenance_window.start_time .lifecycle.maintenance_window.timezone "Europe/London"}}
- **Asia**: {{convertTime .lifecycle.maintenance_window.start_time .lifecycle.maintenance_window.timezone "Asia/Tokyo"}}
{{end}}
```

### Global Maintenance Calendar

```go
{{/* global-maintenance-calendar.tmpl */}}
# Global Maintenance Calendar

{{range .machines}}
{{if .lifecycle.maintenance_window}}
## {{.name}} ({{.lifecycle.maintenance_team}})

**Window**: {{.lifecycle.maintenance_window.start_time}}-{{.lifecycle.maintenance_window.end_time}} {{.lifecycle.maintenance_window.timezone}}
**Days**: {{join .lifecycle.maintenance_window.days_of_week ", "}}
**Contact**: {{.lifecycle.team_contact.primary}}

{{if .lifecycle.maintenance_window.timezone}}
**Local Times**:
- **Team Timezone**: {{.lifecycle.maintenance_window.start_time}}-{{.lifecycle.maintenance_window.end_time}}
- **UTC**: {{convertTime .lifecycle.maintenance_window.start_time .lifecycle.maintenance_window.timezone "UTC"}}-{{convertTime .lifecycle.maintenance_window.end_time .lifecycle.maintenance_window.timezone "UTC"}}
{{end}}

---
{{end}}
{{end}}
```

## Best Practices

### 1. Use IANA Timezone Identifiers

Always use IANA timezone identifiers instead of abbreviations:

```hcl
# ✅ Good
timezone = "America/New_York"
timezone = "Europe/London"
timezone = "Asia/Tokyo"

# ❌ Avoid
timezone = "EST"
timezone = "GMT"
timezone = "JST"
```

### 2. Consider Team Availability

Schedule maintenance windows when the responsible team is available:

```hcl
# For US-based team
maintenance_window = {
  start_time = "02:00"
  end_time = "04:00"
  timezone = "America/New_York"
  days_of_week = ["sunday"]  # Weekend when traffic is lower
}

# For Europe-based team
maintenance_window = {
  start_time = "03:00"
  end_time = "05:00"
  timezone = "Europe/London"
  days_of_week = ["sunday"]
}
```

### 3. Account for Daylight Saving Time

IANA timezone identifiers automatically handle DST transitions:

```hcl
# This automatically adjusts for DST
timezone = "America/New_York"  # EST/EDT
timezone = "Europe/London"     # GMT/BST
```

### 4. Provide Multiple Contact Methods

Include various contact methods for different scenarios:

```hcl
team_contact = {
  primary = "team@company.com"      # General contact
  secondary = "+1-555-0123"         # Emergency phone
  slack_channel = "#team-maintenance" # Real-time coordination
  pagerduty_schedule = "P123456"    # Escalation
}
```

### 5. Use Consistent Days

Standardize maintenance days across teams:

```hcl
# All teams use Sunday for maintenance
days_of_week = ["sunday"]

# Or use multiple days for flexibility
days_of_week = ["saturday", "sunday"]
```

## Migration from Simple Format

If you're currently using the simple format, you can migrate gradually:

### Before
```hcl
maintenance_window = "02:00-04:00"
```

### After
```hcl
maintenance_window = {
  simple = "02:00-04:00"  # Backward compatibility
  # OR
  start_time = "02:00"
  end_time = "04:00"
  timezone = "America/New_York"
  days_of_week = ["sunday"]
}
```

## Validation

The schema validates that:

1. Timezone is a valid IANA identifier
2. Times are in HH:MM format
3. Days of week are valid
4. Either simple format OR timezone-aware format is used

This ensures that maintenance windows are properly configured for global team coordination.
