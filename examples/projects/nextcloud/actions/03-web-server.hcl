# Nextcloud Web Server Configuration
# Actions for configuring Apache and web server settings

action "configure-apache" {
  description = "Configure Apache virtual host for Nextcloud"
  script = "scripts/configure-apache-nextcloud.sh"
  tags = [
    "nextcloud",
    "apache",
    "configuration",
  ]
  machines = [
    "tag:role=web-server",
  ]
  parallel = false
  timeout = 300
}

action "setup-ssl-certificate" {
  description = "Setup SSL certificate for Nextcloud"
  script = "scripts/setup-ssl-nextcloud.sh"
  tags = [
    "nextcloud",
    "ssl",
    "security",
  ]
  machines = [
    "tag:role=web-server",
  ]
  parallel = false
  timeout = 600
} 