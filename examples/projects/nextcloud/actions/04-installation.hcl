# Nextcloud Installation
# Actions for downloading and installing Nextcloud

action "download-nextcloud" {
  description = "Download and extract Nextcloud"
  command = "cd /var/www && wget https://download.nextcloud.com/server/releases/latest.zip && unzip latest.zip && chown -R www-data:www-data nextcloud && chmod -R 755 nextcloud"
  tags = [
    "nextcloud",
    "download",
    "extract",
  ]
  machines = [
    "tag:role=web-server",
  ]
  parallel = false
  timeout = 900
}

action "install-nextcloud" {
  description = "Run Nextcloud installation wizard"
  script = "scripts/install-nextcloud.sh"
  tags = [
    "nextcloud",
    "installation",
    "setup",
  ]
  machines = [
    "tag:role=web-server",
  ]
  parallel = false
  timeout = 1200
} 