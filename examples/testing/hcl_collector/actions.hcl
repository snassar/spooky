actions {
  action "update-system" {
    description = "Update system packages"
    command = "apt update && apt upgrade -y"
    timeout = 300
    parallel = true
    machines = ["web-server", "db-server", "app-server"]
  }

  action "install-nginx" {
    description = "Install and configure Nginx web server"
    command = "apt install -y nginx && systemctl enable nginx"
    timeout = 180
    parallel = false
    machines = ["web-server"]
  }

  action "install-mysql" {
    description = "Install and configure MySQL database"
    command = "apt install -y mysql-server && systemctl enable mysql"
    timeout = 240
    parallel = false
    machines = ["db-server"]
  }

  action "deploy-application" {
    description = "Deploy the main application"
    script = "scripts/deploy.sh"
    timeout = 600
    parallel = true
    tags = ["role = application"]
  }

  action "backup-database" {
    description = "Create database backup"
    command = "mysqldump --all-databases > /backup/db_$(date +%Y%m%d_%H%M%S).sql"
    timeout = 300
    parallel = false
    machines = ["db-server"]
  }

  action "monitor-system" {
    description = "Install monitoring tools"
    command = "apt install -y htop iotop nethogs"
    timeout = 120
    parallel = true
    tags = ["role = monitoring"]
  }

  action "security-update" {
    description = "Apply security updates"
    command = "apt update && apt upgrade --only-upgrade -y"
    timeout = 400
    parallel = true
    tags = ["environment = production"]
  }

  action "cleanup-logs" {
    description = "Clean up old log files"
    command = "find /var/log -name '*.log.*' -mtime +7 -delete"
    timeout = 60
    parallel = true
    machines = ["web-server", "db-server", "app-server", "monitoring-server"]
  }

  action "health-check" {
    description = "Perform system health check"
    command = "systemctl status --no-pager"
    timeout = 30
    parallel = true
    tags = ["team = ops"]
  }
}