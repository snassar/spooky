# Large actions for test-large-actions project
# This file can contain core actions or reference actions in the actions/ directory

actions {
  action "check-status" {
    description = "Check server status"
    command     = "uptime && df -h"
    tags        = ["role=web"]
    parallel    = true
    timeout     = 300
  }
  
  action "system-update" {
    description = "Update system packages"
    command     = "apt update && apt upgrade -y"
    tags        = ["role=web", "role=database"]
    parallel    = false
    timeout     = 600
  }
  
  action "check-disk-space" {
    description = "Check disk space usage"
    command     = "df -h"
    tags        = ["role=web", "role=database", "role=cache"]
    parallel    = true
    timeout     = 300
  }
  
  action "check-memory" {
    description = "Check memory usage"
    command     = "free -h"
    tags        = ["role=web", "role=database", "role=cache"]
    parallel    = true
    timeout     = 300
  }
  
  action "check-processes" {
    description = "Check running processes"
    command     = "ps aux"
    tags        = ["role=web", "role=database", "role=cache"]
    parallel    = true
    timeout     = 300
  }
  
  action "check-network" {
    description = "Check network connectivity"
    command     = "ping -c 3 8.8.8.8"
    tags        = ["role=web", "role=database", "role=cache", "role=load-balancer"]
    parallel    = true
    timeout     = 300
  }
  
  action "check-services" {
    description = "Check service status"
    command     = "systemctl status nginx mysql redis"
    tags        = ["role=web", "role=database", "role=cache"]
    parallel    = true
    timeout     = 300
  }
  
  action "backup-database" {
    description = "Backup database"
    command     = "mysqldump -u root -p database > backup.sql"
    tags        = ["role=database"]
    parallel    = false
    timeout     = 1800
  }
  
  action "restart-services" {
    description = "Restart services"
    command     = "systemctl restart nginx mysql redis"
    tags        = ["role=web", "role=database", "role=cache"]
    parallel    = false
    timeout     = 600
  }
  
  action "check-logs" {
    description = "Check recent logs"
    command     = "tail -n 100 /var/log/nginx/access.log"
    tags        = ["role=web"]
    parallel    = true
    timeout     = 300
  }
}