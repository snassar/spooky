# Nextcloud Verification and Service Management
# Actions for verifying installation and managing services

action "restart-services" {
  description = "Restart Apache and MariaDB services"
  command = "systemctl restart apache2 mariadb"
  tags = [
    "nextcloud",
    "services",
    "restart",
  ]
  machines = [
    "tag:role=web-server",
    "tag:role=database",
  ]
  parallel = false
  timeout = 120
}

action "verify-nextcloud-installation" {
  description = "Verify Nextcloud installation and connectivity"
  command = "curl -f http://localhost/nextcloud/status.php && echo 'Nextcloud installation verified successfully'"
  tags = [
    "nextcloud",
    "verification",
    "health-check",
  ]
  machines = [
    "tag:role=web-server",
  ]
  parallel = true
  timeout = 60
} 