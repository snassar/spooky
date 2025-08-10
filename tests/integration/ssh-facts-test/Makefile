# SSH Facts Test Environment Makefile
# Manages QEMU VMs for testing SSH fact gathering functionality

# Configuration
VM_NAME ?= spooky-facts-test
VM_MEMORY ?= 2048
VM_CPUS ?= 2
VM_DISK_SIZE ?= 20G
SSH_PORT ?= 2222
# The SSH_USER is always spooky
SSH_USER ?= spooky

# Paths
SCRIPT_DIR := $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
ISOS_DIR := $(SCRIPT_DIR)isos
VMS_DIR := $(SCRIPT_DIR)vms
KEYS_DIR := $(SCRIPT_DIR)keys
LOGS_DIR := $(SCRIPT_DIR)logs

# VM-specific paths
VM_DISK := $(VMS_DIR)/$(VM_NAME).qcow2
VM_LOG := $(LOGS_DIR)/$(VM_NAME).log
VM_PID := $(LOGS_DIR)/$(VM_NAME).pid
SSH_KEY := $(KEYS_DIR)/$(VM_NAME)_key
CLOUD_INIT_ISO := $(VMS_DIR)/$(VM_NAME)-cloud-init.iso
CLOUD_INIT_DIR := cloud-init

# Spooky binary
SPOOKY_BIN := spooky
SPOOKY_SRC_DIR := ../../..

# ISO URLs
UBUNTU_22_04_ISO_URL := https://releases.ubuntu.com/22.04/ubuntu-22.04.3-live-server-amd64.iso
FEDORA_39_ISO_URL := https://download.fedoraproject.org/pub/fedora/linux/releases/39/Server/x86_64/iso/Fedora-Server-dvd-x86_64-39-1.5.iso
DEBIAN_12_ISO_URL := https://cdimage.debian.org/debian-cd/current/amd64/iso-cd/debian-12.11.0-amd64-netinst.iso
DEBIAN_12_CLOUD_URL := https://cdimage.debian.org/cdimage/cloud/bookworm/latest/debian-12-generic-amd64.qcow2

# Default target
.PHONY: help
help: ## Show this help message
	@echo "SSH Facts Test Environment"
	@echo "=========================="
	@echo ""
	@echo "Usage: make [target] [VM_NAME=name]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""
	@echo "Examples:"
	@echo "  make create VM_NAME=spooky-facts-test"
	@echo "  make start VM_NAME=spooky-facts-test"
	@echo "  make test-ssh VM_NAME=spooky-facts-test"
	@echo "  make test-project VM_NAME=spooky-facts-test"

# Directory setup
.PHONY: setup
setup: ## Create necessary directories
	@mkdir -p $(ISOS_DIR) $(VMS_DIR) $(KEYS_DIR) $(LOGS_DIR)
	@echo "Created directories: $(ISOS_DIR), $(ISOS_DIR), $(VMS_DIR), $(KEYS_DIR), $(LOGS_DIR)"

# Build spooky binary
.PHONY: build-spooky
build-spooky: ## Build spooky binary locally
	@echo "Building spooky binary..."
	@if [ ! -f "$(SPOOKY_BIN)" ]; then \
		echo "Building spooky from $(SPOOKY_SRC_DIR)..."; \
		cd $(SPOOKY_SRC_DIR) && go build -o $(PWD)/$(SPOOKY_BIN) .; \
		echo "Built spooky binary: $(SPOOKY_BIN)"; \
	else \
		echo "Spooky binary already exists: $(SPOOKY_BIN)"; \
	fi

.PHONY: rebuild-spooky
rebuild-spooky: ## Force rebuild spooky binary
	@echo "Force rebuilding spooky binary..."
	@rm -f $(SPOOKY_BIN)
	@cd $(SPOOKY_SRC_DIR) && go build -o $(PWD)/$(SPOOKY_BIN) .
	@echo "Rebuilt spooky binary: $(SPOOKY_BIN)"

# ISO download targets
.PHONY: download-ubuntu download-fedora download-debian
download-ubuntu: ## Download Ubuntu 22.04 Server ISO
	@echo "Downloading Ubuntu 22.04 Server ISO..."
	@wget -O $(ISOS_DIR)/ubuntu-22.04-server-amd64.iso $(UBUNTU_22_04_ISO_URL)
	@echo "Downloaded: $(ISOS_DIR)/ubuntu-22.04-server-amd64.iso"

download-fedora: ## Download Fedora 39 Server ISO
	@echo "Downloading Fedora 39 Server ISO..."
	@wget -O $(ISOS_DIR)/fedora-39-server-amd64.iso $(FEDORA_39_ISO_URL)
	@echo "Downloaded: $(ISOS_DIR)/fedora-39-server-amd64.iso"

download-debian: ## Download Debian 12 ISO
	@echo "Downloading Debian 12 ISO..."
	@wget -O $(ISOS_DIR)/debian-12-amd64.iso $(DEBIAN_12_ISO_URL)
	@echo "Downloaded: $(ISOS_DIR)/debian-12-amd64.iso"

.PHONY: download-debian-cloud
download-debian-cloud: ## Download Debian 12 Cloud Image
	@echo "Downloading Debian 12 Cloud Image..."
	@wget -O $(VMS_DIR)/debian-12-cloud-amd64.qcow2 $(DEBIAN_12_CLOUD_URL)
	@echo "Downloaded: $(VMS_DIR)/debian-12-cloud-amd64.qcow2"

# VM creation
.PHONY: create
create: setup $(SSH_KEY) ## Create a new VM (requires ISO download first)
	@echo "Creating VM: $(VM_NAME)"
	@echo "Disk: $(VM_DISK)"
	@echo "SSH Key: $(SSH_KEY)"
	@qemu-img create -f qcow2 $(VM_DISK) $(VM_DISK_SIZE)
	@echo "Created VM disk: $(VM_DISK)"
	@echo "To install OS, run: make install-os VM_NAME=$(VM_NAME)"

.PHONY: create-from-cloud
create-from-cloud: download-debian-cloud $(SSH_KEY) ## Create VM from cloud image
	@echo "Creating VM from cloud image: $(VM_NAME)"
	@echo "Disk: $(VM_DISK)"
	@echo "SSH Key: $(SSH_KEY)"
	@cp $(VMS_DIR)/debian-12-cloud-amd64.qcow2 $(VM_DISK)
	@echo "Created VM disk from cloud image: $(VM_DISK)"
	@echo "To start VM, run: make start VM_NAME=$(VM_NAME)"

# SSH key generation
$(SSH_KEY): setup
	@echo "Generating SSH key for $(VM_NAME)..."
	@ssh-keygen -t rsa -b 4096 -f $(SSH_KEY) -N "" -C "$(VM_NAME)-test-key"
	@chmod 600 $(SSH_KEY)
	@echo "Generated SSH key: $(SSH_KEY)"

# OS installation
.PHONY: install-os
install-os: ## Install OS from ISO (interactive)
	@echo "Starting OS installation for $(VM_NAME)..."
	@echo "This will start QEMU in headless mode for installation."
	@echo "After installation, run: make start VM_NAME=$(VM_NAME)"
	@qemu-system-x86_64 \
		-hda $(VM_DISK) \
		-cdrom $(ISOS_DIR)/debian-12-amd64.iso \
		-m $(VM_MEMORY) \
		-smp $(VM_CPUS) \
		-net nic,model=virtio \
		-net user,hostfwd=tcp::$(SSH_PORT)-:22 \
		-display none \
		-boot d

.PHONY: install-os-cloud-init
install-os-cloud-init: create-cloud-init ## Install OS with cloud-init (automated)
	@echo "Starting OS installation with cloud-init for $(VM_NAME)..."
	@echo "This will start QEMU in headless mode for installation."
	@echo "Cloud-init will automatically configure SSH keys after installation."
	@echo "After installation, run: make start-cloud-init VM_NAME=$(VM_NAME)"
	@qemu-system-x86_64 \
		-hda $(VM_DISK) \
		-cdrom $(ISOS_DIR)/debian-12-amd64.iso \
		-drive file=$(CLOUD_INIT_ISO),if=virtio,readonly=on \
		-m $(VM_MEMORY) \
		-smp $(VM_CPUS) \
		-net nic,model=virtio \
		-net user,hostfwd=tcp::$(SSH_PORT)-:22 \
		-display none \
		-boot d

# VM management
.PHONY: start
start: ## Start the VM in background
	@echo "Starting VM: $(VM_NAME)"
	@qemu-system-x86_64 \
		-hda $(VM_DISK) \
		-m $(VM_MEMORY) \
		-smp $(VM_CPUS) \
		-net nic,model=virtio \
		-net user,hostfwd=tcp::$(SSH_PORT)-:22 \
		-display none \
		-daemonize \
		-pidfile $(VM_PID) \
		-D $(VM_LOG)
	@echo "VM started. PID: $$(cat $(VM_PID))"
	@echo "SSH: ssh -i $(SSH_KEY) -p $(SSH_PORT) $(SSH_USER)@localhost"
	@echo "Log: $(VM_LOG)"

.PHONY: stop
stop: ## Stop the VM
	@if [ -f $(VM_PID) ]; then \
		echo "Stopping VM: $(VM_NAME)"; \
		kill $$(cat $(VM_PID)) 2>/dev/null || true; \
		rm -f $(VM_PID); \
		echo "VM stopped"; \
	else \
		echo "VM $(VM_NAME) is not running"; \
	fi

.PHONY: restart
restart: stop start ## Restart the VM

# Cloud-init support
.PHONY: create-cloud-init
create-cloud-init: setup ## Create cloud-init configuration and ISO
	@echo "Creating cloud-init configuration for $(VM_NAME)..."
	@./create-cloud-init.sh $(VM_NAME) $(SSH_USER) $(SSH_KEY).pub
	@echo "Creating cloud-init ISO..."
	@cloud-localds $(CLOUD_INIT_ISO) $(CLOUD_INIT_DIR)/user-data $(CLOUD_INIT_DIR)/meta-data
	@echo "Cloud-init ISO created: $(CLOUD_INIT_ISO)"

.PHONY: create-cloud-init-vm
create-cloud-init-vm: create-cloud-init ## Create VM with cloud-init support
	@echo "Creating VM with cloud-init: $(VM_NAME)"
	@qemu-img create -f qcow2 $(VM_DISK) $(VM_DISK_SIZE)
	@echo "Created VM disk: $(VM_DISK)"
	@echo "To start with cloud-init, run: make start-cloud-init VM_NAME=$(VM_NAME)"

.PHONY: create-cloud-init-vm-from-image
create-cloud-init-vm-from-image: download-debian-cloud create-cloud-init ## Create VM from cloud image with cloud-init
	@echo "Creating VM from cloud image with cloud-init: $(VM_NAME)"
	@cp $(VMS_DIR)/debian-12-cloud-amd64.qcow2 $(VM_DISK)
	@echo "Created VM disk from cloud image: $(VM_DISK)"
	@echo "To start with cloud-init, run: make start-cloud-init VM_NAME=$(VM_NAME)"

.PHONY: start-cloud-init
start-cloud-init: ## Start VM with cloud-init (automated setup)
	@echo "Starting VM with cloud-init: $(VM_NAME)"
	@if [ ! -f $(CLOUD_INIT_ISO) ]; then \
		echo "Cloud-init ISO not found. Creating..."; \
		make create-cloud-init VM_NAME=$(VM_NAME); \
	fi
	@qemu-system-x86_64 \
		-hda $(VM_DISK) \
		-cdrom $(CLOUD_INIT_ISO) \
		-m $(VM_MEMORY) \
		-smp $(VM_CPUS) \
		-net nic,model=virtio \
		-net user,hostfwd=tcp::$(SSH_PORT)-:22 \
		-display none \
		-daemonize \
		-pidfile $(VM_PID) \
		-D $(VM_LOG)
	@echo "VM started with cloud-init. PID: $$(cat $(VM_PID))"
	@echo "SSH: ssh -i $(SSH_KEY) -p $(SSH_PORT) $(SSH_USER)@localhost"
	@echo "Log: $(VM_LOG)"
	@echo "Waiting for cloud-init to complete..."
	@echo "Check progress with: make logs VM_NAME=$(VM_NAME)"

.PHONY: logs
logs: ## Show VM logs
	@if [ -f $(VM_LOG) ]; then \
		echo "VM Logs for $(VM_NAME):"; \
		echo "========================"; \
		tail -f $(VM_LOG); \
	else \
		echo "No log file found for $(VM_NAME)"; \
	fi

.PHONY: status
status: ## Show VM status
	@echo "VM Status: $(VM_NAME)"
	@echo "================"
	@if [ -f $(VM_PID) ]; then \
		echo "Status: Running (PID: $$(cat $(VM_PID)))"; \
		echo "SSH Port: $(SSH_PORT)"; \
		echo "SSH Command: ssh -i $(SSH_KEY) -p $(SSH_PORT) $(SSH_USER)@localhost"; \
	else \
		echo "Status: Not running"; \
	fi
	@echo "Disk: $(VM_DISK) ($$(du -h $(VM_DISK) 2>/dev/null | cut -f1 || echo 'not created'))"
	@echo "Log: $(VM_LOG)"

# SSH testing
.PHONY: test-ssh
test-ssh: ## Test SSH connectivity
	@echo "Testing SSH connectivity to $(VM_NAME)..."
	@ssh -i $(SSH_KEY) -p $(SSH_PORT) -o ConnectTimeout=10 -o StrictHostKeyChecking=no $(SSH_USER)@localhost "echo 'SSH connection successful'"

.PHONY: connect
connect: ## Connect to VM via SSH
	@echo "Connecting to $(VM_NAME) via SSH..."
	@ssh -i $(SSH_KEY) -p $(SSH_PORT) -o StrictHostKeyChecking=no $(SSH_USER)@localhost

# Fact gathering tests

.PHONY: test-project
test-project: build-spooky ## Test SSH fact gathering using the test project
	@echo "Testing SSH fact gathering using test project..."
	@echo "Running ./$(SPOOKY_BIN) actions run test-ssh-facts..."
	@cd testing-fact-gathering && ../$(SPOOKY_BIN) actions run test-ssh-facts .



.PHONY: validate-project
validate-project: build-spooky ## Validate the test project configuration
	@echo "Validating test project configuration..."
	@echo "Running ./$(SPOOKY_BIN) project validate..."
	@cd testing-fact-gathering && ../$(SPOOKY_BIN) project validate .

.PHONY: list-actions
list-actions: build-spooky ## List available actions in the test project
	@echo "Listing available actions in test project..."
	@echo "Running ./$(SPOOKY_BIN) actions list..."
	@cd testing-fact-gathering && ../$(SPOOKY_BIN) actions list .

# Cleanup
.PHONY: clean
clean: ## Remove everything, leaving directory in git clone state
	@echo "Cleaning up everything..."
	@echo "Stopping all VMs..."
	@for vm in $(VMS_DIR)/*.qcow2; do \
		if [ -f "$$vm" ]; then \
			vm_name=$$(basename "$$vm" .qcow2); \
			echo "Stopping $$vm_name..."; \
			make stop VM_NAME=$$vm_name || true; \
		fi; \
	done
	@echo "Removing all generated files..."
	@rm -rf $(KEYS_DIR)/*
	@rm -rf $(ISOS_DIR)/*
	@rm -rf $(LOGS_DIR)/*
	@rm -rf $(VMS_DIR)/*
	@rm -rf $(CLOUD_INIT_DIR)
	@rm -f $(SPOOKY_BIN)
	@echo "Clean complete - directory is now in git clone state"

# Information
.PHONY: info
info: ## Show VM information
	@echo "VM Information: $(VM_NAME)"
	@echo "================"
	@echo "Memory: $(VM_MEMORY)MB"
	@echo "CPUs: $(VM_CPUS)"
	@echo "Disk: $(VM_DISK_SIZE)"
	@echo "SSH Port: $(SSH_PORT)"
	@echo "SSH User: $(SSH_USER)"
	@echo "SSH Key: $(SSH_KEY)"
	@echo ""
	@echo "Available commands:"
	@echo "  make start VM_NAME=$(VM_NAME)"
	@echo "  make stop VM_NAME=$(VM_NAME)"
	@echo "  make connect VM_NAME=$(VM_NAME)"
	@echo "  make test-ssh VM_NAME=$(VM_NAME)"
	@echo "  make test-project VM_NAME=$(VM_NAME)"

# List available VMs
.PHONY: list
list: ## List all available VMs
	@echo "Available VMs:"
	@echo "=============="
	@for vm in $(VMS_DIR)/*.qcow2; do \
		if [ -f "$$vm" ]; then \
			vm_name=$$(basename "$$vm" .qcow2); \
			vm_size=$$(du -h "$$vm" | cut -f1); \
			echo "  $$vm_name ($$vm_size)"; \
		fi; \
	done
	@echo ""
	@echo "Available ISOs:"
	@echo "==============="
	@for iso in $(ISOS_DIR)/*.iso; do \
		if [ -f "$$iso" ]; then \
			iso_name=$$(basename "$$iso"); \
			iso_size=$$(du -h "$$iso" | cut -f1); \
			echo "  $$iso_name ($$iso_size)"; \
		fi; \
	done

# Quick setup for common distributions
.PHONY: setup-debian setup-ubuntu setup-fedora
setup-debian: download-debian create ## Download Debian ISO and create VM (default)
setup-debian-cloud: create-from-cloud ## Setup Debian VM from cloud image (recommended)
setup-debian-cloud-init: create-cloud-init-vm-from-image ## Setup Debian VM from cloud image with cloud-init (recommended)
	@echo "Debian VM setup complete. Run: make install-os VM_NAME=$(VM_NAME)"

setup-ubuntu: download-ubuntu create ## Download Ubuntu ISO and create VM
	@echo "Ubuntu VM setup complete. Run: make install-os VM_NAME=$(VM_NAME)"

setup-fedora: download-fedora create ## Download Fedora ISO and create VM
	@echo "Fedora VM setup complete. Run: make install-os VM_NAME=$(VM_NAME)"

# Cloud-init quick setup (automated SSH key configuration)
.PHONY: setup-debian-cloud-init setup-ubuntu-cloud-init setup-fedora-cloud-init
setup-debian-cloud-init: download-debian create-cloud-init-vm ## Download Debian ISO and create VM with cloud-init
	@echo "Debian VM with cloud-init setup complete."
	@echo "Run: make install-os-cloud-init VM_NAME=$(VM_NAME)"
	@echo "Then: make start-cloud-init VM_NAME=$(VM_NAME)"

setup-ubuntu-cloud-init: download-ubuntu create-cloud-init-vm ## Download Ubuntu ISO and create VM with cloud-init
	@echo "Ubuntu VM with cloud-init setup complete."
	@echo "Run: make install-os-cloud-init VM_NAME=$(VM_NAME)"
	@echo "Then: make start-cloud-init VM_NAME=$(VM_NAME)"

setup-fedora-cloud-init: download-fedora create-cloud-init-vm ## Download Fedora ISO and create VM with cloud-init
	@echo "Fedora VM with cloud-init setup complete."
	@echo "Run: make install-os-cloud-init VM_NAME=$(VM_NAME)"
	@echo "Then: make start-cloud-init VM_NAME=$(VM_NAME)" 