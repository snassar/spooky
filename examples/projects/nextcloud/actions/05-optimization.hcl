# Nextcloud Optimization and Maintenance
# Actions for optimizing performance and setting up maintenance tasks

action "configure-nextcloud-optimization" {
  description = "Configure Nextcloud for optimal performance"
  script = "scripts/optimize-nextcloud.sh"
  tags = [
    "nextcloud",
    "optimization",
    "performance",
  ]
  machines = [
    "tag:role=web-server",
  ]
  parallel = false
  timeout = 300
}

action "setup-backup-cron" {
  description = "Setup automated backup cron job for Nextcloud"
  script = "scripts/setup-nextcloud-backup.sh"
  tags = [
    "nextcloud",
    "backup",
    "automation",
  ]
  machines = [
    "tag:role=web-server",
  ]
  parallel = false
  timeout = 180
} 