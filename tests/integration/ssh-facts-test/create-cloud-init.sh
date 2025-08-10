#!/bin/bash

# Create cloud-init configuration for VM setup
# This script generates the necessary cloud-init files to set up SSH keys and user

set -e

# Configuration
VM_NAME=${1:-"spooky-facts-test"}
SSH_USER=${2:-"spooky"}
SSH_KEY_PATH=${3:-"keys/${VM_NAME}_key.pub"}
CLOUD_INIT_DIR="cloud-init"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Creating cloud-init configuration for VM: ${VM_NAME}${NC}"

# Create cloud-init directory
mkdir -p "${CLOUD_INIT_DIR}"

# Check if SSH public key exists
if [ ! -f "${SSH_KEY_PATH}" ]; then
    echo -e "${YELLOW}Warning: SSH public key not found at ${SSH_KEY_PATH}${NC}"
    echo -e "${YELLOW}Generating new SSH key pair...${NC}"
    
    # Generate SSH key pair
    mkdir -p keys
    ssh-keygen -t rsa -b 4096 -f "keys/${VM_NAME}_key" -N "" -C "spooky-test@${VM_NAME}"
    SSH_KEY_PATH="keys/${VM_NAME}_key.pub"
fi

# Read the public key
SSH_PUBLIC_KEY=$(cat "${SSH_KEY_PATH}")

# Create user-data file (cloud-init configuration)
cat > "${CLOUD_INIT_DIR}/user-data" << EOF
#cloud-config
hostname: ${VM_NAME}
manage_etc_hosts: true

# Create user
users:
  - name: ${SSH_USER}
    sudo: ALL=(ALL) NOPASSWD:ALL
    shell: /bin/bash
    ssh_authorized_keys:
      - ${SSH_PUBLIC_KEY}
    # Set password for console access (optional)
    lock_passwd: false
    # Default password: spooky123
    hashed_passwd: \$6\$rounds=656000\$spooky\$spooky123

# Update system
package_update: true
package_upgrade: true

# Install packages
packages:
  - openssh-server
  - curl
  - wget
  - vim
  - htop
  - net-tools

# Configure SSH
ssh_pwauth: false
disable_root: true

# Write files
write_files:
  - path: /etc/ssh/sshd_config.d/99-spooky.conf
    content: |
      PasswordAuthentication no
      PubkeyAuthentication yes
      PermitRootLogin no
      AllowUsers ${SSH_USER}
    owner: root:root
    permissions: '0644'

# Run commands
runcmd:
  - systemctl enable ssh
  - systemctl restart ssh
  - echo "SSH key setup complete for user ${SSH_USER}" > /home/${SSH_USER}/ssh-setup.log
  - chown ${SSH_USER}:${SSH_USER} /home/${SSH_USER}/ssh-setup.log

# Final message
final_message: "Cloud-init setup complete for ${VM_NAME}"
EOF

# Create meta-data file
cat > "${CLOUD_INIT_DIR}/meta-data" << EOF
instance-id: ${VM_NAME}
local-hostname: ${VM_NAME}
EOF

# Create network-config file (optional, for static IP if needed)
cat > "${CLOUD_INIT_DIR}/network-config" << EOF
version: 2
ethernets:
  eth0:
    dhcp4: true
    dhcp6: true
EOF

echo -e "${GREEN}Cloud-init configuration created:${NC}"
echo -e "  - ${CLOUD_INIT_DIR}/user-data"
echo -e "  - ${CLOUD_INIT_DIR}/meta-data"
echo -e "  - ${CLOUD_INIT_DIR}/network-config"
echo -e ""
echo -e "${GREEN}SSH Configuration:${NC}"
echo -e "  - User: ${SSH_USER}"
echo -e "  - Key: ${SSH_KEY_PATH}"
echo -e "  - Password: spooky123 (for console access)"
echo -e ""
echo -e "${YELLOW}Next steps:${NC}"
echo -e "  1. Use 'make create-cloud-init' to create VM with cloud-init"
echo -e "  2. Use 'make start-cloud-init' to start VM with cloud-init"
echo -e "  3. Wait for cloud-init to complete (check logs with 'make logs')" 