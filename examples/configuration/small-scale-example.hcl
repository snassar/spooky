# Small configuration for spooky SSH automation tool
# Small hosting provider with 40 servers (10 hardware + 30 VMs)
# Data centers: FRA00 (Frankfurt) and BER0 (Berlin)
# IP range: 10.0.0.0/8
# Generated with Git-style IDs for deterministic identification
# Config ID: 20250716+e1d8

# =============================================================================
# SERVERS (40 total)
# =============================================================================

server "machine-fe8d67b1073acad1" {
  host     = "10.1.1.1"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-001"
  tags = {
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
  }
}

server "machine-d901fe5ad5a4c1da" {
  host     = "10.1.1.2"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-002"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-09d6cdbeec529371" {
  host     = "10.1.1.3"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-003"
  tags = {
    os = "debian12"
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
  }
}

server "machine-abb4ae3e23b58a38" {
  host     = "10.1.1.4"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-004"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-8fb47c05992f0abf" {
  host     = "10.1.1.5"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-005"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-9693ac14f4e38f90" {
  host     = "10.2.1.1"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-6"
  tags = {
    capacity = "high"
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
  }
}

server "machine-255a10b65814aff9" {
  host     = "10.2.1.2"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-7"
  tags = {
    capacity = "high"
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
  }
}

server "machine-01520dfeb71cfa26" {
  host     = "10.2.1.3"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-8"
  tags = {
    capacity = "high"
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
  }
}

server "machine-63ee5c040398d256" {
  host     = "10.2.1.4"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-9"
  tags = {
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
  }
}

server "machine-2b6f1ee1ae7e9540" {
  host     = "10.2.1.5"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-10"
  tags = {
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
  }
}

server "vm-e5994bec43a922bd" {
  host     = "10.1.10.1"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0001"
  tags = {
    tier = "production"
    db_type = "postgresql"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
  }
}

server "vm-673d4f53e03ed16b" {
  host     = "10.1.10.2"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0002"
  tags = {
    db_type = "mysql"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
  }
}

server "vm-1875e88a9247c59b" {
  host     = "10.1.10.3"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0003"
  tags = {
    tier = "staging"
    db_type = "mongodb"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
  }
}

server "vm-9a5601746ea4abd3" {
  host     = "10.1.20.1"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0004"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "nginx"
  }
}

server "vm-9309e6427ee35428" {
  host     = "10.1.20.2"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0005"
  tags = {
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "apache"
    datacenter = "FRA00"
  }
}

server "vm-3ed9cbf44ad72ea8" {
  host     = "10.1.20.3"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0006"
  tags = {
    role = "web"
    os = "debian12"
    tier = "staging"
    web_type = "nginx"
    datacenter = "FRA00"
    type = "vm"
  }
}

server "vm-c5be7cbcb3bfaaa9" {
  host     = "10.1.30.1"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0007"
  tags = {
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "compute"
    datacenter = "FRA00"
    type = "vm"
  }
}

server "vm-91fbfa667e7d9665" {
  host     = "10.1.30.2"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0008"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "batch"
  }
}

server "vm-679f3c936c5ad185" {
  host     = "10.1.30.3"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0009"
  tags = {
    tier = "staging"
    workload_type = "compute"
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
    os = "debian12"
  }
}

server "vm-f2fe186bf6180ee6" {
  host     = "10.1.40.1"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0010"
  tags = {
    os = "debian12"
    tier = "production"
    storage_type = "block"
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
  }
}

server "vm-cf0bb901a938deeb" {
  host     = "10.1.40.2"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0011"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "object"
  }
}

server "vm-139547b30e118e44" {
  host     = "10.1.40.3"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0012"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "staging"
    storage_type = "block"
  }
}

server "vm-f117ad522fc8bb6c" {
  host     = "10.2.10.1"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0013"
  tags = {
    os = "debian12"
    tier = "production"
    db_type = "postgresql"
    datacenter = "BER0"
    type = "vm"
    role = "database"
  }
}

server "vm-03a35fa7b1335525" {
  host     = "10.2.10.2"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0014"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mysql"
  }
}

server "vm-a38bf5b746914c48" {
  host     = "10.2.10.3"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0015"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "staging"
    db_type = "mongodb"
  }
}

server "vm-6a4949e27160c48c" {
  host     = "10.2.20.1"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0016"
  tags = {
    web_type = "nginx"
    datacenter = "BER0"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
  }
}

server "vm-b977828b6ce5320f" {
  host     = "10.2.20.2"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0017"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "apache"
  }
}

server "vm-72823f982c4b2db9" {
  host     = "10.2.20.3"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0018"
  tags = {
    web_type = "nginx"
    datacenter = "BER0"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "staging"
  }
}

server "vm-3f17439f62e204cd" {
  host     = "10.2.30.1"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0019"
  tags = {
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "compute"
    datacenter = "BER0"
  }
}

server "vm-7db875d2dd81636c" {
  host     = "10.2.30.2"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0020"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "batch"
  }
}

server "vm-8ae1811ed5a97cae" {
  host     = "10.2.30.3"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0021"
  tags = {
    role = "workload"
    os = "debian12"
    tier = "staging"
    workload_type = "compute"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-44e2c3366a307954" {
  host     = "10.2.40.1"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0022"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "block"
  }
}

server "vm-0a7a0f62e44c67d0" {
  host     = "10.2.40.2"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0023"
  tags = {
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "object"
    datacenter = "BER0"
  }
}

server "vm-1c08dc3b8226a3f2" {
  host     = "10.2.40.3"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0024"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "staging"
    storage_type = "block"
  }
}

# =============================================================================
# ACTIONS FOR TESTING
# =============================================================================

action "check-production-status" {
  description = "Check status of all production servers"
  command     = "uptime && df -h && systemctl status --no-pager"
  tags        = ["tier=production"]
  parallel    = true
  timeout     = 300
}

action "update-databases" {
  description = "Update all database servers"
  command     = "apt update && apt upgrade -y"
  tags        = ["role=database"]
  parallel    = true
  timeout     = 600
}

action "check-fra00-web" {
  description = "Check FRA00 web servers specifically"
  command     = "systemctl status nginx apache2 --no-pager"
  tags        = ["datacenter=FRA00", "role=web"]
  parallel    = true
  timeout     = 120
}

action "backup-storage" {
  description = "Create backups on all storage servers"
  script      = "/usr/local/bin/backup-storage.sh"
  tags        = ["role=storage"]
  parallel    = false
  timeout     = 1800
}

action "check-hardware" {
  description = "Check hardware server status"
  command     = "lscpu && free -h && df -h"
  tags        = ["type=hardware"]
  parallel    = true
  timeout     = 180
}

action "update-staging" {
  description = "Update all staging servers"
  command     = "apt update && apt upgrade -y"
  tags        = ["tier=staging"]
  parallel    = true
  timeout     = 600
}

action "check-ber0-db" {
  description = "Check BER0 database servers"
  command     = "systemctl status postgresql mysql mongod --no-pager"
  tags        = ["datacenter=BER0", "role=database"]
  parallel    = true
  timeout     = 120
}

action "monitor-compute" {
  description = "Monitor compute workload servers"
  command     = "htop --batch --iterations=1 && nvidia-smi"
  tags        = ["workload_type=compute"]
  parallel    = true
  timeout     = 60
}

action "check-nginx" {
  description = "Check all nginx web servers"
  command     = "nginx -t && systemctl status nginx --no-pager"
  tags        = ["web_type=nginx"]
  parallel    = true
  timeout     = 90
}

action "full-system-check" {
  description = "Comprehensive system check"
  command     = "uptime && df -h && free -h && systemctl --failed --no-pager"
  servers     = ["machine-550e8400e29b41d4", "vm-550e8400e29b41da", "vm-550e8400e29b41e6"]
  parallel    = true
  timeout     = 300
}

