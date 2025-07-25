# Nextcloud Database Configuration
# Actions for setting up and configuring the database

action "configure-mariadb" {
  description = "Configure MariaDB for Nextcloud"
  script = "scripts/configure-mariadb.sh"
  tags = [
    "nextcloud",
    "database",
    "mariadb",
  ]
  machines = [
    "tag:role=database",
  ]
  parallel = false
  timeout = 600
}

action "create-nextcloud-database" {
  description = "Create Nextcloud database and user"
  script = "scripts/create-nextcloud-db.sh"
  tags = [
    "nextcloud",
    "database",
    "setup",
  ]
  machines = [
    "tag:role=database",
  ]
  parallel = false
  timeout = 300
} 