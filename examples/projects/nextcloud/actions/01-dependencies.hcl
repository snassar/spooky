# Nextcloud Dependencies Installation
# Actions for installing and configuring system dependencies

action "install-nextcloud-dependencies" {
  description = "Install system dependencies required for Nextcloud"
  command = "apt update && apt install -y apache2 mariadb-server mariadb-client php php-mysql php-gd php-json php-curl php-mbstring php-intl php-xml php-zip php-bz2 php-ldap php-imagick unzip wget curl"
  tags = [
    "nextcloud",
    "installation",
    "dependencies",
  ]
  machines = [
    "tag:role=web-server",
  ]
  parallel = false
  timeout = 1800
}

action "configure-php" {
  description = "Configure PHP settings for Nextcloud"
  command = "sed -i 's/upload_max_filesize = 2M/upload_max_filesize = 512M/' /etc/php/*/apache2/php.ini && sed -i 's/post_max_size = 8M/post_max_size = 512M/' /etc/php/*/apache2/php.ini && sed -i 's/memory_limit = 128M/memory_limit = 512M/' /etc/php/*/apache2/php.ini"
  tags = [
    "nextcloud",
    "php",
    "configuration",
  ]
  machines = [
    "tag:role=web-server",
  ]
  parallel = false
  timeout = 120
} 