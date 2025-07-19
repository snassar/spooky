# Medium configuration for spooky SSH automation tool
# Medium hosting provider with 400 servers (100 hardware + 300 VMs)
# Data centers: FRA00 (Frankfurt) and BER0 (Berlin)
# IP range: 10.0.0.0/8
# Generated with Git-style IDs for deterministic identification
# Config ID: 20250716+3eea

# =============================================================================
# SERVERS (400 total)
# =============================================================================

server "machine-fe8d67b1073acad1" {
  host     = "10.1.1.1"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-001"
  tags = {
    os = "debian12"
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
  }
}

server "machine-d901fe5ad5a4c1da" {
  host     = "10.1.1.2"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-002"
  tags = {
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
  }
}

server "machine-09d6cdbeec529371" {
  host     = "10.1.1.3"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-003"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
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
    role = "vm-host"
    os = "debian12"
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
  }
}

server "machine-0a5d5cd4a8fce37e" {
  host     = "10.1.1.6"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-006"
  tags = {
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
  }
}

server "machine-c64a6fdcbcb92c06" {
  host     = "10.1.1.7"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-007"
  tags = {
    role = "vm-host"
    os = "debian12"
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
  }
}

server "machine-fb1728eed5b3ff78" {
  host     = "10.1.1.8"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-008"
  tags = {
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
  }
}

server "machine-29f0395903f03032" {
  host     = "10.1.1.9"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-009"
  tags = {
    os = "debian12"
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
  }
}

server "machine-2a89036cf28150c7" {
  host     = "10.1.1.10"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-010"
  tags = {
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
  }
}

server "machine-23e58e7e7ef6436a" {
  host     = "10.1.1.11"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-011"
  tags = {
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
    datacenter = "FRA00"
  }
}

server "machine-1126434b1bf4bac5" {
  host     = "10.1.1.12"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-012"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-fbc155b97dc1738f" {
  host     = "10.1.1.13"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-013"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-a60ef95e91267f62" {
  host     = "10.1.1.14"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-014"
  tags = {
    role = "vm-host"
    os = "debian12"
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
  }
}

server "machine-9a81ccb5a0dd9aa8" {
  host     = "10.1.1.15"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-015"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-c64f15e4c4cff512" {
  host     = "10.1.1.16"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-016"
  tags = {
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
  }
}

server "machine-29299e9377e18bc8" {
  host     = "10.1.1.17"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-017"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-ed11966b60d5984d" {
  host     = "10.1.1.18"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-018"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-37d9dc0d0c4c6a28" {
  host     = "10.1.1.19"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-019"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-a5b449a8e00bc053" {
  host     = "10.1.1.20"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-020"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-d3b6362920f429c9" {
  host     = "10.1.1.21"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-021"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-c24581523e6e3b1d" {
  host     = "10.1.1.22"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-022"
  tags = {
    os = "debian12"
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
  }
}

server "machine-e4bf3d53d5ea8ebc" {
  host     = "10.1.1.23"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-023"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-0e4f034485a59b02" {
  host     = "10.1.1.24"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-024"
  tags = {
    os = "debian12"
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
  }
}

server "machine-9bb9362712321b23" {
  host     = "10.1.1.25"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-025"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-4ad85c21c0c0d46c" {
  host     = "10.1.1.26"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-026"
  tags = {
    os = "debian12"
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
  }
}

server "machine-f57e01ff53bbaf44" {
  host     = "10.1.1.27"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-027"
  tags = {
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
  }
}

server "machine-fbea27a5a6b423dd" {
  host     = "10.1.1.28"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-028"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-cc3e2b9b3681e060" {
  host     = "10.1.1.29"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-029"
  tags = {
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
    datacenter = "FRA00"
  }
}

server "machine-3257befbe28a54f0" {
  host     = "10.1.1.30"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-030"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-3062458ed659b2ef" {
  host     = "10.1.1.31"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-031"
  tags = {
    role = "vm-host"
    os = "debian12"
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
  }
}

server "machine-74ea20890fe8b4b7" {
  host     = "10.1.1.32"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-032"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-9ac1adb2c5e950ec" {
  host     = "10.1.1.33"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-033"
  tags = {
    os = "debian12"
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
  }
}

server "machine-fd27121836a2b488" {
  host     = "10.1.1.34"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-034"
  tags = {
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
  }
}

server "machine-45ad9200435ff363" {
  host     = "10.1.1.35"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-035"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-d67b3f8ea2ceb85f" {
  host     = "10.1.1.36"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-036"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-7c71d458c3a00306" {
  host     = "10.1.1.37"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-037"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-b37cfc23ecf05852" {
  host     = "10.1.1.38"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-038"
  tags = {
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
  }
}

server "machine-76fbaf8e826f5aad" {
  host     = "10.1.1.39"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-039"
  tags = {
    os = "debian12"
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
  }
}

server "machine-2a0fd9f703780c3f" {
  host     = "10.1.1.40"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-040"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-440dafb6f5a9a770" {
  host     = "10.1.1.41"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-041"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-01352c7d11368321" {
  host     = "10.1.1.42"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-042"
  tags = {
    os = "debian12"
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
  }
}

server "machine-5dbbaef49bd8e15e" {
  host     = "10.1.1.43"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-043"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-269b5807d9804411" {
  host     = "10.1.1.44"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-044"
  tags = {
    role = "vm-host"
    os = "debian12"
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
  }
}

server "machine-19b3592237b01675" {
  host     = "10.1.1.45"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-045"
  tags = {
    os = "debian12"
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
  }
}

server "machine-cc04281b03d5c968" {
  host     = "10.1.1.46"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-046"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-bf242c886fb80cfc" {
  host     = "10.1.1.47"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-047"
  tags = {
    role = "vm-host"
    os = "debian12"
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
  }
}

server "machine-2d533f57407342ff" {
  host     = "10.1.1.48"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-048"
  tags = {
    os = "debian12"
    capacity = "high"
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
  }
}

server "machine-1365e516ea0939a9" {
  host     = "10.1.1.49"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-049"
  tags = {
    datacenter = "FRA00"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-ff49d273af23f7e4" {
  host     = "10.1.1.50"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-050"
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
  password = "hardware-secure-pass-51"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-255a10b65814aff9" {
  host     = "10.2.1.2"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-52"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-01520dfeb71cfa26" {
  host     = "10.2.1.3"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-53"
  tags = {
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
  }
}

server "machine-63ee5c040398d256" {
  host     = "10.2.1.4"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-54"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-2b6f1ee1ae7e9540" {
  host     = "10.2.1.5"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-55"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-3f9f073f8b5bfbf9" {
  host     = "10.2.1.6"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-56"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-97c23fb8365f355d" {
  host     = "10.2.1.7"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-57"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-d223e20699a8a3df" {
  host     = "10.2.1.8"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-58"
  tags = {
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
  }
}

server "machine-d023080b14244473" {
  host     = "10.2.1.9"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-59"
  tags = {
    role = "vm-host"
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
    type = "hardware"
  }
}

server "machine-32932f0046657c0c" {
  host     = "10.2.1.10"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-60"
  tags = {
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
  }
}

server "machine-8809d75289d0afce" {
  host     = "10.2.1.11"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-61"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-1f36f419f54cd235" {
  host     = "10.2.1.12"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-62"
  tags = {
    role = "vm-host"
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
    type = "hardware"
  }
}

server "machine-9f883d6f0b924ae8" {
  host     = "10.2.1.13"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-63"
  tags = {
    role = "vm-host"
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
    type = "hardware"
  }
}

server "machine-8b481da884d52a8b" {
  host     = "10.2.1.14"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-64"
  tags = {
    role = "vm-host"
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
    type = "hardware"
  }
}

server "machine-64bcffa6f0dbd6a7" {
  host     = "10.2.1.15"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-65"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-9a50157a3be5c345" {
  host     = "10.2.1.16"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-66"
  tags = {
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
  }
}

server "machine-8f9e98eb56cddaa3" {
  host     = "10.2.1.17"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-67"
  tags = {
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
  }
}

server "machine-60256ced7995be0b" {
  host     = "10.2.1.18"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-68"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-224f4189b8623579" {
  host     = "10.2.1.19"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-69"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-8a329068003c7189" {
  host     = "10.2.1.20"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-70"
  tags = {
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
  }
}

server "machine-0f72679729b4b726" {
  host     = "10.2.1.21"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-71"
  tags = {
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
  }
}

server "machine-727a1fe22f2ac1d1" {
  host     = "10.2.1.22"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-72"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-2e9848feaba3d2ba" {
  host     = "10.2.1.23"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-73"
  tags = {
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
  }
}

server "machine-bf9502132141a13a" {
  host     = "10.2.1.24"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-74"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-8791a075c92a7656" {
  host     = "10.2.1.25"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-75"
  tags = {
    capacity = "high"
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
  }
}

server "machine-58511ad23b0ff37c" {
  host     = "10.2.1.26"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-76"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-8a18c8f186ab230d" {
  host     = "10.2.1.27"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-77"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-c426aa8716901bd8" {
  host     = "10.2.1.28"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-78"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-18307d58add44de3" {
  host     = "10.2.1.29"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-79"
  tags = {
    capacity = "high"
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
  }
}

server "machine-4ef49be61c1a168f" {
  host     = "10.2.1.30"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-80"
  tags = {
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
  }
}

server "machine-b2aad324017c3633" {
  host     = "10.2.1.31"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-81"
  tags = {
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
  }
}

server "machine-1d3e125f00121665" {
  host     = "10.2.1.32"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-82"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-642c5c2c79b6a161" {
  host     = "10.2.1.33"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-83"
  tags = {
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
  }
}

server "machine-90be7c04d7c3d46d" {
  host     = "10.2.1.34"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-84"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-1b17034178573579" {
  host     = "10.2.1.35"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-85"
  tags = {
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
  }
}

server "machine-bd0f101433b2cb5d" {
  host     = "10.2.1.36"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-86"
  tags = {
    role = "vm-host"
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
    type = "hardware"
  }
}

server "machine-e0d7beb2ef6195da" {
  host     = "10.2.1.37"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-87"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-d91852a58f928d14" {
  host     = "10.2.1.38"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-88"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-88548024ad08bfe6" {
  host     = "10.2.1.39"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-89"
  tags = {
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
  }
}

server "machine-09a06d76b2931cc3" {
  host     = "10.2.1.40"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-90"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-c5c8092a7e5fb583" {
  host     = "10.2.1.41"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-91"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-c21a8a70db29152b" {
  host     = "10.2.1.42"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-92"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-92e02d056aebce3a" {
  host     = "10.2.1.43"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-93"
  tags = {
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
  }
}

server "machine-9dc2bff584e9b0e6" {
  host     = "10.2.1.44"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-94"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-f0eb02e632edc1d7" {
  host     = "10.2.1.45"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-95"
  tags = {
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
  }
}

server "machine-eafc4642db7ccbfb" {
  host     = "10.2.1.46"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-96"
  tags = {
    capacity = "high"
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
  }
}

server "machine-b3a840dca4f9a681" {
  host     = "10.2.1.47"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-97"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "machine-9cc17bb29780f53d" {
  host     = "10.2.1.48"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-98"
  tags = {
    role = "vm-host"
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
    type = "hardware"
  }
}

server "machine-91044d10a45fd7a2" {
  host     = "10.2.1.49"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-99"
  tags = {
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
    datacenter = "BER0"
  }
}

server "machine-88d9297c6509cb7d" {
  host     = "10.2.1.50"
  port     = 22
  user     = "admin"
  password = "hardware-secure-pass-100"
  tags = {
    datacenter = "BER0"
    type = "hardware"
    role = "vm-host"
    os = "debian12"
    capacity = "high"
  }
}

server "vm-e5994bec43a922bd" {
  host     = "10.1.10.1"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0001"
  tags = {
    db_type = "postgresql"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
  }
}

server "vm-673d4f53e03ed16b" {
  host     = "10.1.10.2"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0002"
  tags = {
    tier = "production"
    db_type = "mysql"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
  }
}

server "vm-1875e88a9247c59b" {
  host     = "10.1.10.3"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0003"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mongodb"
  }
}

server "vm-129ae1f1cb8670d2" {
  host     = "10.1.10.4"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0004"
  tags = {
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "postgresql"
    datacenter = "FRA00"
    type = "vm"
  }
}

server "vm-a37e39b44e2ba8bb" {
  host     = "10.1.10.5"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0005"
  tags = {
    os = "debian12"
    tier = "production"
    db_type = "mysql"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
  }
}

server "vm-a7d39c1618125ea9" {
  host     = "10.1.10.6"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0006"
  tags = {
    db_type = "mongodb"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
  }
}

server "vm-4adbe7a08ebcdde1" {
  host     = "10.1.10.7"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0007"
  tags = {
    tier = "production"
    db_type = "postgresql"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
  }
}

server "vm-b230d7206ab7e72d" {
  host     = "10.1.10.8"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0008"
  tags = {
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mysql"
    datacenter = "FRA00"
  }
}

server "vm-648d357964c92281" {
  host     = "10.1.10.9"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0009"
  tags = {
    tier = "production"
    db_type = "mongodb"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
  }
}

server "vm-c7cb0bcc82f2e6ca" {
  host     = "10.1.10.10"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0010"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "postgresql"
  }
}

server "vm-f55f2b89a2fb191c" {
  host     = "10.1.10.11"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0011"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mysql"
  }
}

server "vm-7e4dc6dc6cccee33" {
  host     = "10.1.10.12"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0012"
  tags = {
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mongodb"
    datacenter = "FRA00"
    type = "vm"
  }
}

server "vm-4e92ed7b134d533e" {
  host     = "10.1.10.13"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0013"
  tags = {
    tier = "production"
    db_type = "postgresql"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
  }
}

server "vm-790b114eb6f4c2ad" {
  host     = "10.1.10.14"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0014"
  tags = {
    os = "debian12"
    tier = "production"
    db_type = "mysql"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
  }
}

server "vm-de3390cd5ed2e2e5" {
  host     = "10.1.10.15"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0015"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mongodb"
  }
}

server "vm-3f1fa684975d8777" {
  host     = "10.1.10.16"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0016"
  tags = {
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "postgresql"
    datacenter = "FRA00"
    type = "vm"
  }
}

server "vm-2d253c6d6803e408" {
  host     = "10.1.10.17"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0017"
  tags = {
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mysql"
    datacenter = "FRA00"
  }
}

server "vm-f4ec6b9de1deeaa0" {
  host     = "10.1.10.18"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0018"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mongodb"
  }
}

server "vm-0cb722e65b1b63d8" {
  host     = "10.1.10.19"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0019"
  tags = {
    db_type = "postgresql"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
  }
}

server "vm-423eb2934fc82a15" {
  host     = "10.1.10.20"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0020"
  tags = {
    tier = "production"
    db_type = "mysql"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
  }
}

server "vm-1a4646cc65b6a7fb" {
  host     = "10.1.10.21"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0021"
  tags = {
    tier = "production"
    db_type = "mongodb"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
  }
}

server "vm-2c8ca3b2aeef8409" {
  host     = "10.1.10.22"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0022"
  tags = {
    os = "debian12"
    tier = "production"
    db_type = "postgresql"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
  }
}

server "vm-328abb5c9113d523" {
  host     = "10.1.10.23"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0023"
  tags = {
    db_type = "mysql"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
  }
}

server "vm-e43f88ec6fb6e07e" {
  host     = "10.1.10.24"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0024"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mongodb"
  }
}

server "vm-1c9230c11ab8608d" {
  host     = "10.1.10.25"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0025"
  tags = {
    tier = "production"
    db_type = "postgresql"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
  }
}

server "vm-ac63064c31e776ba" {
  host     = "10.1.10.26"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0026"
  tags = {
    os = "debian12"
    tier = "production"
    db_type = "mysql"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
  }
}

server "vm-f430d2542e0fd7d1" {
  host     = "10.1.10.27"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0027"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mongodb"
  }
}

server "vm-3d538d0965c1ae8f" {
  host     = "10.1.10.28"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0028"
  tags = {
    db_type = "postgresql"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
  }
}

server "vm-1be12064a467fa34" {
  host     = "10.1.10.29"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0029"
  tags = {
    os = "debian12"
    tier = "production"
    db_type = "mysql"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
  }
}

server "vm-0029a93e84c1e6d6" {
  host     = "10.1.10.30"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0030"
  tags = {
    os = "debian12"
    tier = "staging"
    db_type = "mongodb"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
  }
}

server "vm-62b49fe106bfc774" {
  host     = "10.1.10.31"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0031"
  tags = {
    db_type = "postgresql"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "staging"
  }
}

server "vm-fdf367a50b367329" {
  host     = "10.1.10.32"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0032"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "staging"
    db_type = "mysql"
  }
}

server "vm-15736f06185e59ec" {
  host     = "10.1.10.33"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0033"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "staging"
    db_type = "mongodb"
  }
}

server "vm-b0eb9bda06ea0bbe" {
  host     = "10.1.10.34"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0034"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "staging"
    db_type = "postgresql"
  }
}

server "vm-c8f4deb1aee12157" {
  host     = "10.1.10.35"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0035"
  tags = {
    os = "debian12"
    tier = "staging"
    db_type = "mysql"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
  }
}

server "vm-a4df34e3fadcf7cb" {
  host     = "10.1.10.36"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0036"
  tags = {
    tier = "staging"
    db_type = "mongodb"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
  }
}

server "vm-8b57a17981dc0862" {
  host     = "10.1.10.37"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0037"
  tags = {
    db_type = "postgresql"
    datacenter = "FRA00"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "staging"
  }
}

server "vm-9a5601746ea4abd3" {
  host     = "10.1.20.1"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0038"
  tags = {
    tier = "production"
    web_type = "nginx"
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
  }
}

server "vm-9309e6427ee35428" {
  host     = "10.1.20.2"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0039"
  tags = {
    tier = "production"
    web_type = "apache"
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
  }
}

server "vm-3ed9cbf44ad72ea8" {
  host     = "10.1.20.3"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0040"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "nginx"
  }
}

server "vm-9ec5635b896f585e" {
  host     = "10.1.20.4"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0041"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "apache"
  }
}

server "vm-29617d831810b034" {
  host     = "10.1.20.5"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0042"
  tags = {
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "nginx"
    datacenter = "FRA00"
    type = "vm"
  }
}

server "vm-30265c23f8747c2c" {
  host     = "10.1.20.6"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0043"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "apache"
  }
}

server "vm-94b9e67ed9dc136d" {
  host     = "10.1.20.7"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0044"
  tags = {
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "nginx"
    datacenter = "FRA00"
  }
}

server "vm-1e7d53c48ce367c3" {
  host     = "10.1.20.8"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0045"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "apache"
  }
}

server "vm-a76a80c97a5d8e64" {
  host     = "10.1.20.9"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0046"
  tags = {
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "nginx"
    datacenter = "FRA00"
    type = "vm"
  }
}

server "vm-142b780c812cdb61" {
  host     = "10.1.20.10"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0047"
  tags = {
    web_type = "apache"
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
  }
}

server "vm-469e8ba8a46c4aa8" {
  host     = "10.1.20.11"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0048"
  tags = {
    web_type = "nginx"
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
  }
}

server "vm-b7d7a45ffb28c4b4" {
  host     = "10.1.20.12"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0049"
  tags = {
    os = "debian12"
    tier = "production"
    web_type = "apache"
    datacenter = "FRA00"
    type = "vm"
    role = "web"
  }
}

server "vm-27155b43d09f048a" {
  host     = "10.1.20.13"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0050"
  tags = {
    tier = "production"
    web_type = "nginx"
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
  }
}

server "vm-230197ad59ee237d" {
  host     = "10.1.20.14"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0051"
  tags = {
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "apache"
    datacenter = "FRA00"
  }
}

server "vm-f7ad7741bdb0e6f8" {
  host     = "10.1.20.15"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0052"
  tags = {
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "nginx"
    datacenter = "FRA00"
  }
}

server "vm-bee0d4bb3038ed44" {
  host     = "10.1.20.16"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0053"
  tags = {
    os = "debian12"
    tier = "production"
    web_type = "apache"
    datacenter = "FRA00"
    type = "vm"
    role = "web"
  }
}

server "vm-8d2b63ca92372d78" {
  host     = "10.1.20.17"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0054"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "nginx"
  }
}

server "vm-a06275f21b13efb9" {
  host     = "10.1.20.18"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0055"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "apache"
  }
}

server "vm-68b60a1b3715b942" {
  host     = "10.1.20.19"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0056"
  tags = {
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "nginx"
    datacenter = "FRA00"
    type = "vm"
  }
}

server "vm-d884f7faef3fa2b7" {
  host     = "10.1.20.20"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0057"
  tags = {
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "apache"
    datacenter = "FRA00"
  }
}

server "vm-cdec72935c0bfbc9" {
  host     = "10.1.20.21"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0058"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "nginx"
  }
}

server "vm-7c30bb193c9f525b" {
  host     = "10.1.20.22"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0059"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "apache"
  }
}

server "vm-eb78742a0c45dcd5" {
  host     = "10.1.20.23"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0060"
  tags = {
    tier = "production"
    web_type = "nginx"
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
  }
}

server "vm-920d7cda9ea9449a" {
  host     = "10.1.20.24"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0061"
  tags = {
    os = "debian12"
    tier = "production"
    web_type = "apache"
    datacenter = "FRA00"
    type = "vm"
    role = "web"
  }
}

server "vm-e3c4d392d8d73664" {
  host     = "10.1.20.25"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0062"
  tags = {
    tier = "production"
    web_type = "nginx"
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
  }
}

server "vm-50befc022abb7a5b" {
  host     = "10.1.20.26"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0063"
  tags = {
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "apache"
    datacenter = "FRA00"
  }
}

server "vm-4a6c5a7013a2fe0e" {
  host     = "10.1.20.27"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0064"
  tags = {
    tier = "production"
    web_type = "nginx"
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
  }
}

server "vm-519e5afc9f7e437b" {
  host     = "10.1.20.28"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0065"
  tags = {
    os = "debian12"
    tier = "production"
    web_type = "apache"
    datacenter = "FRA00"
    type = "vm"
    role = "web"
  }
}

server "vm-4157f7b8c0c3b314" {
  host     = "10.1.20.29"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0066"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "nginx"
  }
}

server "vm-8f20dea4976fe5d0" {
  host     = "10.1.20.30"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0067"
  tags = {
    web_type = "apache"
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "staging"
  }
}

server "vm-4c196593f53226fe" {
  host     = "10.1.20.31"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0068"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "staging"
    web_type = "nginx"
  }
}

server "vm-eee5524a41d067f8" {
  host     = "10.1.20.32"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0069"
  tags = {
    tier = "staging"
    web_type = "apache"
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
  }
}

server "vm-44c61db670bc0737" {
  host     = "10.1.20.33"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0070"
  tags = {
    web_type = "nginx"
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "staging"
  }
}

server "vm-7a2f85e82d42e0f5" {
  host     = "10.1.20.34"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0071"
  tags = {
    tier = "staging"
    web_type = "apache"
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
  }
}

server "vm-42febddd4c4179cb" {
  host     = "10.1.20.35"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0072"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "staging"
    web_type = "nginx"
  }
}

server "vm-94a283c3497b0cbf" {
  host     = "10.1.20.36"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0073"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "staging"
    web_type = "apache"
  }
}

server "vm-b2cd4e3e6ab1dec6" {
  host     = "10.1.20.37"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0074"
  tags = {
    web_type = "nginx"
    datacenter = "FRA00"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "staging"
  }
}

server "vm-c5be7cbcb3bfaaa9" {
  host     = "10.1.30.1"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0075"
  tags = {
    workload_type = "compute"
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
  }
}

server "vm-91fbfa667e7d9665" {
  host     = "10.1.30.2"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0076"
  tags = {
    workload_type = "batch"
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
  }
}

server "vm-679f3c936c5ad185" {
  host     = "10.1.30.3"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0077"
  tags = {
    os = "debian12"
    tier = "production"
    workload_type = "compute"
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
  }
}

server "vm-75a4c437a3c08aef" {
  host     = "10.1.30.4"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0078"
  tags = {
    os = "debian12"
    tier = "production"
    workload_type = "batch"
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
  }
}

server "vm-85481eaea9e4eda8" {
  host     = "10.1.30.5"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0079"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "compute"
  }
}

server "vm-2c6c047123bac934" {
  host     = "10.1.30.6"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0080"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "batch"
  }
}

server "vm-16abe01ee3dae424" {
  host     = "10.1.30.7"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0081"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "compute"
  }
}

server "vm-71d17d462d3756b2" {
  host     = "10.1.30.8"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0082"
  tags = {
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "batch"
    datacenter = "FRA00"
  }
}

server "vm-89bf1db197f1444d" {
  host     = "10.1.30.9"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0083"
  tags = {
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "compute"
    datacenter = "FRA00"
  }
}

server "vm-781257f3c2ead575" {
  host     = "10.1.30.10"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0084"
  tags = {
    workload_type = "batch"
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
  }
}

server "vm-127eb2d462cd52da" {
  host     = "10.1.30.11"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0085"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "compute"
  }
}

server "vm-563c9a5a341a385b" {
  host     = "10.1.30.12"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0086"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "batch"
  }
}

server "vm-45ce763755b3fd2c" {
  host     = "10.1.30.13"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0087"
  tags = {
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "compute"
    datacenter = "FRA00"
    type = "vm"
  }
}

server "vm-09209a3aba1a63b1" {
  host     = "10.1.30.14"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0088"
  tags = {
    tier = "production"
    workload_type = "batch"
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
    os = "debian12"
  }
}

server "vm-67c3c3aed77a7804" {
  host     = "10.1.30.15"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0089"
  tags = {
    workload_type = "compute"
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
  }
}

server "vm-d0dd57b435d6d5c8" {
  host     = "10.1.30.16"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0090"
  tags = {
    os = "debian12"
    tier = "production"
    workload_type = "batch"
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
  }
}

server "vm-f968720d966d1cc3" {
  host     = "10.1.30.17"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0091"
  tags = {
    os = "debian12"
    tier = "production"
    workload_type = "compute"
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
  }
}

server "vm-b8bf84e0dc5d468a" {
  host     = "10.1.30.18"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0092"
  tags = {
    workload_type = "batch"
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
  }
}

server "vm-7ff0da72828e2f48" {
  host     = "10.1.30.19"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0093"
  tags = {
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "compute"
    datacenter = "FRA00"
    type = "vm"
  }
}

server "vm-41cdfc69011c2292" {
  host     = "10.1.30.20"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0094"
  tags = {
    os = "debian12"
    tier = "production"
    workload_type = "batch"
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
  }
}

server "vm-e2d748c6fcddc032" {
  host     = "10.1.30.21"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0095"
  tags = {
    tier = "production"
    workload_type = "compute"
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
    os = "debian12"
  }
}

server "vm-67b3d94418972570" {
  host     = "10.1.30.22"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0096"
  tags = {
    workload_type = "batch"
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
  }
}

server "vm-0a359f53b6df7123" {
  host     = "10.1.30.23"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0097"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "compute"
  }
}

server "vm-b9a7ff4dd1e951dd" {
  host     = "10.1.30.24"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0098"
  tags = {
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "batch"
    datacenter = "FRA00"
    type = "vm"
  }
}

server "vm-435381cc81958ac9" {
  host     = "10.1.30.25"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0099"
  tags = {
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "compute"
    datacenter = "FRA00"
    type = "vm"
  }
}

server "vm-7e26c9f0af406cd0" {
  host     = "10.1.30.26"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0100"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "batch"
  }
}

server "vm-ba1132bd14c6cd9c" {
  host     = "10.1.30.27"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0101"
  tags = {
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "compute"
    datacenter = "FRA00"
  }
}

server "vm-982062cdd79e86b4" {
  host     = "10.1.30.28"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0102"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "batch"
  }
}

server "vm-35bf95d31c480195" {
  host     = "10.1.30.29"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0103"
  tags = {
    workload_type = "compute"
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
  }
}

server "vm-4c7693e9322c57d0" {
  host     = "10.1.30.30"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0104"
  tags = {
    role = "workload"
    os = "debian12"
    tier = "staging"
    workload_type = "batch"
    datacenter = "FRA00"
    type = "vm"
  }
}

server "vm-b9394abdf1e15b30" {
  host     = "10.1.30.31"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0105"
  tags = {
    tier = "staging"
    workload_type = "compute"
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
    os = "debian12"
  }
}

server "vm-9bbcf3ac44fe7cd4" {
  host     = "10.1.30.32"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0106"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "staging"
    workload_type = "batch"
  }
}

server "vm-d1585b4bea3e32b4" {
  host     = "10.1.30.33"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0107"
  tags = {
    os = "debian12"
    tier = "staging"
    workload_type = "compute"
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
  }
}

server "vm-ceb5f5413ba5b7cb" {
  host     = "10.1.30.34"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0108"
  tags = {
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "staging"
    workload_type = "batch"
    datacenter = "FRA00"
  }
}

server "vm-ace0801ab4ccf654" {
  host     = "10.1.30.35"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0109"
  tags = {
    role = "workload"
    os = "debian12"
    tier = "staging"
    workload_type = "compute"
    datacenter = "FRA00"
    type = "vm"
  }
}

server "vm-fad980ccf885d3be" {
  host     = "10.1.30.36"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0110"
  tags = {
    os = "debian12"
    tier = "staging"
    workload_type = "batch"
    datacenter = "FRA00"
    type = "vm"
    role = "workload"
  }
}

server "vm-0c48e67d917ce9e8" {
  host     = "10.1.30.37"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0111"
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
  password = "vm-secure-pass-0112"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "block"
  }
}

server "vm-cf0bb901a938deeb" {
  host     = "10.1.40.2"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0113"
  tags = {
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "object"
    datacenter = "FRA00"
    type = "vm"
  }
}

server "vm-139547b30e118e44" {
  host     = "10.1.40.3"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0114"
  tags = {
    storage_type = "block"
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
  }
}

server "vm-9a1368fc7176aae0" {
  host     = "10.1.40.4"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0115"
  tags = {
    tier = "production"
    storage_type = "object"
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
  }
}

server "vm-ac27d1a2b0b8ba02" {
  host     = "10.1.40.5"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0116"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "block"
  }
}

server "vm-f8b2cd064372c204" {
  host     = "10.1.40.6"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0117"
  tags = {
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "object"
    datacenter = "FRA00"
  }
}

server "vm-becac2f4536c41f7" {
  host     = "10.1.40.7"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0118"
  tags = {
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "block"
    datacenter = "FRA00"
    type = "vm"
  }
}

server "vm-aa8170a6464227a8" {
  host     = "10.1.40.8"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0119"
  tags = {
    tier = "production"
    storage_type = "object"
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
  }
}

server "vm-019ecd600e40fc3e" {
  host     = "10.1.40.9"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0120"
  tags = {
    storage_type = "block"
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
  }
}

server "vm-61cc1c3e5727b8ae" {
  host     = "10.1.40.10"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0121"
  tags = {
    storage_type = "object"
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
  }
}

server "vm-dbd0c7cd2041444e" {
  host     = "10.1.40.11"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0122"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "block"
  }
}

server "vm-537f4a05f0e0bdf0" {
  host     = "10.1.40.12"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0123"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "object"
  }
}

server "vm-35be4e5cdcbfe861" {
  host     = "10.1.40.13"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0124"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "block"
  }
}

server "vm-a961fa4b8b5e60a2" {
  host     = "10.1.40.14"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0125"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "object"
  }
}

server "vm-69fb6e3426455ed4" {
  host     = "10.1.40.15"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0126"
  tags = {
    storage_type = "block"
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
  }
}

server "vm-93db8da1a0e89cb1" {
  host     = "10.1.40.16"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0127"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "object"
  }
}

server "vm-576144530083767f" {
  host     = "10.1.40.17"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0128"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "block"
  }
}

server "vm-a325b80603492e79" {
  host     = "10.1.40.18"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0129"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "object"
  }
}

server "vm-076b82cde8e97b66" {
  host     = "10.1.40.19"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0130"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "block"
  }
}

server "vm-d25ce9f3b7e33245" {
  host     = "10.1.40.20"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0131"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "object"
  }
}

server "vm-239e3daa8b0cc158" {
  host     = "10.1.40.21"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0132"
  tags = {
    os = "debian12"
    tier = "production"
    storage_type = "block"
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
  }
}

server "vm-c998d5c55fde56f4" {
  host     = "10.1.40.22"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0133"
  tags = {
    tier = "production"
    storage_type = "object"
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
  }
}

server "vm-c7816f02387adccf" {
  host     = "10.1.40.23"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0134"
  tags = {
    tier = "production"
    storage_type = "block"
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
  }
}

server "vm-d3f51bf8d0f6805a" {
  host     = "10.1.40.24"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0135"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "object"
  }
}

server "vm-ac32766a82b56554" {
  host     = "10.1.40.25"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0136"
  tags = {
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "block"
    datacenter = "FRA00"
    type = "vm"
  }
}

server "vm-dba240742be79c28" {
  host     = "10.1.40.26"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0137"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "object"
  }
}

server "vm-423164883fcbfa21" {
  host     = "10.1.40.27"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0138"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "block"
  }
}

server "vm-96102ca7d79f964b" {
  host     = "10.1.40.28"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0139"
  tags = {
    storage_type = "object"
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
  }
}

server "vm-d9883f1d89ab84ea" {
  host     = "10.1.40.29"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0140"
  tags = {
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "block"
    datacenter = "FRA00"
  }
}

server "vm-ba4eb8a701eb4f64" {
  host     = "10.1.40.30"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0141"
  tags = {
    role = "storage"
    os = "debian12"
    tier = "staging"
    storage_type = "object"
    datacenter = "FRA00"
    type = "vm"
  }
}

server "vm-f5ea652596d7e118" {
  host     = "10.1.40.31"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0142"
  tags = {
    tier = "staging"
    storage_type = "block"
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
  }
}

server "vm-6f7cbf626beea2b7" {
  host     = "10.1.40.32"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0143"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "staging"
    storage_type = "object"
  }
}

server "vm-8a87e574703460a5" {
  host     = "10.1.40.33"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0144"
  tags = {
    storage_type = "block"
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "staging"
  }
}

server "vm-14a4367fa48dbf43" {
  host     = "10.1.40.34"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0145"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "staging"
    storage_type = "object"
  }
}

server "vm-de2bdaf5afd1157b" {
  host     = "10.1.40.35"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0146"
  tags = {
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "staging"
    storage_type = "block"
    datacenter = "FRA00"
  }
}

server "vm-28ad02b488eeb156" {
  host     = "10.1.40.36"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0147"
  tags = {
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "staging"
    storage_type = "object"
  }
}

server "vm-b5871572a9d68599" {
  host     = "10.1.40.37"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0148"
  tags = {
    os = "debian12"
    tier = "staging"
    storage_type = "block"
    datacenter = "FRA00"
    type = "vm"
    role = "storage"
  }
}

server "vm-f117ad522fc8bb6c" {
  host     = "10.2.10.1"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0149"
  tags = {
    tier = "production"
    db_type = "postgresql"
    datacenter = "BER0"
    type = "vm"
    role = "database"
    os = "debian12"
  }
}

server "vm-03a35fa7b1335525" {
  host     = "10.2.10.2"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0150"
  tags = {
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mysql"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-a38bf5b746914c48" {
  host     = "10.2.10.3"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0151"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mongodb"
  }
}

server "vm-4f1d404b378ec69e" {
  host     = "10.2.10.4"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0152"
  tags = {
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "postgresql"
    datacenter = "BER0"
  }
}

server "vm-598670663e6323a4" {
  host     = "10.2.10.5"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0153"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mysql"
  }
}

server "vm-79b0c6cd90e98e4a" {
  host     = "10.2.10.6"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0154"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mongodb"
  }
}

server "vm-30e6dcf4b8aee8ce" {
  host     = "10.2.10.7"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0155"
  tags = {
    tier = "production"
    db_type = "postgresql"
    datacenter = "BER0"
    type = "vm"
    role = "database"
    os = "debian12"
  }
}

server "vm-2c6fdcdd6c756eca" {
  host     = "10.2.10.8"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0156"
  tags = {
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mysql"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-c1d1f091fea6e5ca" {
  host     = "10.2.10.9"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0157"
  tags = {
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mongodb"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-32877f1a3d3ca95e" {
  host     = "10.2.10.10"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0158"
  tags = {
    tier = "production"
    db_type = "postgresql"
    datacenter = "BER0"
    type = "vm"
    role = "database"
    os = "debian12"
  }
}

server "vm-f37bffe45c407178" {
  host     = "10.2.10.11"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0159"
  tags = {
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mysql"
    datacenter = "BER0"
  }
}

server "vm-c265eb115d0f0fd3" {
  host     = "10.2.10.12"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0160"
  tags = {
    tier = "production"
    db_type = "mongodb"
    datacenter = "BER0"
    type = "vm"
    role = "database"
    os = "debian12"
  }
}

server "vm-ea7d98e3ebaf14bf" {
  host     = "10.2.10.13"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0161"
  tags = {
    tier = "production"
    db_type = "postgresql"
    datacenter = "BER0"
    type = "vm"
    role = "database"
    os = "debian12"
  }
}

server "vm-91c0b4d2bd2fb70e" {
  host     = "10.2.10.14"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0162"
  tags = {
    db_type = "mysql"
    datacenter = "BER0"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
  }
}

server "vm-64a2c5906e53723c" {
  host     = "10.2.10.15"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0163"
  tags = {
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mongodb"
    datacenter = "BER0"
  }
}

server "vm-aa3b98bdc732cdb2" {
  host     = "10.2.10.16"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0164"
  tags = {
    os = "debian12"
    tier = "production"
    db_type = "postgresql"
    datacenter = "BER0"
    type = "vm"
    role = "database"
  }
}

server "vm-a9c7fe218dd4bc9f" {
  host     = "10.2.10.17"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0165"
  tags = {
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mysql"
    datacenter = "BER0"
  }
}

server "vm-27964bc14414f5d6" {
  host     = "10.2.10.18"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0166"
  tags = {
    tier = "production"
    db_type = "mongodb"
    datacenter = "BER0"
    type = "vm"
    role = "database"
    os = "debian12"
  }
}

server "vm-0210ee4d42435586" {
  host     = "10.2.10.19"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0167"
  tags = {
    os = "debian12"
    tier = "production"
    db_type = "postgresql"
    datacenter = "BER0"
    type = "vm"
    role = "database"
  }
}

server "vm-7639ff0e765cf2fa" {
  host     = "10.2.10.20"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0168"
  tags = {
    db_type = "mysql"
    datacenter = "BER0"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
  }
}

server "vm-c7c670c672671058" {
  host     = "10.2.10.21"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0169"
  tags = {
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mongodb"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-c96650236b685251" {
  host     = "10.2.10.22"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0170"
  tags = {
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "postgresql"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-21b32242f4e6c14c" {
  host     = "10.2.10.23"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0171"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mysql"
  }
}

server "vm-f0dd24f50b9273ac" {
  host     = "10.2.10.24"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0172"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mongodb"
  }
}

server "vm-adfd636d6a3b1783" {
  host     = "10.2.10.25"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0173"
  tags = {
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "postgresql"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-ebdb50f7a35f0aa0" {
  host     = "10.2.10.26"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0174"
  tags = {
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mysql"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-e09f80e9b6bbb2c4" {
  host     = "10.2.10.27"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0175"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mongodb"
  }
}

server "vm-46c1a625f8d0466a" {
  host     = "10.2.10.28"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0176"
  tags = {
    tier = "production"
    db_type = "postgresql"
    datacenter = "BER0"
    type = "vm"
    role = "database"
    os = "debian12"
  }
}

server "vm-208d35c6204f04e1" {
  host     = "10.2.10.29"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0177"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "production"
    db_type = "mysql"
  }
}

server "vm-3c14edf99e783465" {
  host     = "10.2.10.30"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0178"
  tags = {
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "staging"
    db_type = "mongodb"
    datacenter = "BER0"
  }
}

server "vm-05f1d4562a9a2822" {
  host     = "10.2.10.31"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0179"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "staging"
    db_type = "postgresql"
  }
}

server "vm-c48466e11a17de0b" {
  host     = "10.2.10.32"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0180"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "staging"
    db_type = "mysql"
  }
}

server "vm-52ab03c5f2da2d68" {
  host     = "10.2.10.33"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0181"
  tags = {
    tier = "staging"
    db_type = "mongodb"
    datacenter = "BER0"
    type = "vm"
    role = "database"
    os = "debian12"
  }
}

server "vm-a7a1d27230ad84d3" {
  host     = "10.2.10.34"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0182"
  tags = {
    role = "database"
    os = "debian12"
    tier = "staging"
    db_type = "postgresql"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-439160d151527835" {
  host     = "10.2.10.35"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0183"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "staging"
    db_type = "mysql"
  }
}

server "vm-f249defadda6796d" {
  host     = "10.2.10.36"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0184"
  tags = {
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "staging"
    db_type = "mongodb"
    datacenter = "BER0"
  }
}

server "vm-833b1b5795789447" {
  host     = "10.2.10.37"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0185"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "database"
    os = "debian12"
    tier = "staging"
    db_type = "postgresql"
  }
}

server "vm-6a4949e27160c48c" {
  host     = "10.2.20.1"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0186"
  tags = {
    tier = "production"
    web_type = "nginx"
    datacenter = "BER0"
    type = "vm"
    role = "web"
    os = "debian12"
  }
}

server "vm-b977828b6ce5320f" {
  host     = "10.2.20.2"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0187"
  tags = {
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "apache"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-72823f982c4b2db9" {
  host     = "10.2.20.3"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0188"
  tags = {
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "nginx"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-558c5076c0cf40db" {
  host     = "10.2.20.4"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0189"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "apache"
  }
}

server "vm-5eee875c276fba0b" {
  host     = "10.2.20.5"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0190"
  tags = {
    web_type = "nginx"
    datacenter = "BER0"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
  }
}

server "vm-60e493df1b80fd82" {
  host     = "10.2.20.6"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0191"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "apache"
  }
}

server "vm-03725f2419e2d743" {
  host     = "10.2.20.7"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0192"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "nginx"
  }
}

server "vm-f2a1616d1ae0aef9" {
  host     = "10.2.20.8"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0193"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "apache"
  }
}

server "vm-cdd3d742d567cc6c" {
  host     = "10.2.20.9"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0194"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "nginx"
  }
}

server "vm-cf2c74394d88325d" {
  host     = "10.2.20.10"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0195"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "apache"
  }
}

server "vm-072bbc348d9c9777" {
  host     = "10.2.20.11"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0196"
  tags = {
    tier = "production"
    web_type = "nginx"
    datacenter = "BER0"
    type = "vm"
    role = "web"
    os = "debian12"
  }
}

server "vm-1951ad7a90f065e9" {
  host     = "10.2.20.12"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0197"
  tags = {
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "apache"
    datacenter = "BER0"
  }
}

server "vm-e8fdbd92dc67619a" {
  host     = "10.2.20.13"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0198"
  tags = {
    web_type = "nginx"
    datacenter = "BER0"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
  }
}

server "vm-608f5f31c0db0d86" {
  host     = "10.2.20.14"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0199"
  tags = {
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "apache"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-571a1c25e61e23de" {
  host     = "10.2.20.15"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0200"
  tags = {
    os = "debian12"
    tier = "production"
    web_type = "nginx"
    datacenter = "BER0"
    type = "vm"
    role = "web"
  }
}

server "vm-aacd431e6c38cad7" {
  host     = "10.2.20.16"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0201"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "apache"
  }
}

server "vm-069f9d48f83a2c4b" {
  host     = "10.2.20.17"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0202"
  tags = {
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "nginx"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-c1dad26e325b5ed5" {
  host     = "10.2.20.18"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0203"
  tags = {
    os = "debian12"
    tier = "production"
    web_type = "apache"
    datacenter = "BER0"
    type = "vm"
    role = "web"
  }
}

server "vm-b18e410dd0269a5a" {
  host     = "10.2.20.19"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0204"
  tags = {
    os = "debian12"
    tier = "production"
    web_type = "nginx"
    datacenter = "BER0"
    type = "vm"
    role = "web"
  }
}

server "vm-28c99b5e15e550d2" {
  host     = "10.2.20.20"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0205"
  tags = {
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "apache"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-d315778d21cab5f6" {
  host     = "10.2.20.21"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0206"
  tags = {
    web_type = "nginx"
    datacenter = "BER0"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
  }
}

server "vm-b5170883b17a7f45" {
  host     = "10.2.20.22"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0207"
  tags = {
    tier = "production"
    web_type = "apache"
    datacenter = "BER0"
    type = "vm"
    role = "web"
    os = "debian12"
  }
}

server "vm-2f326db05472e2ec" {
  host     = "10.2.20.23"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0208"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "nginx"
  }
}

server "vm-e016d8540b288091" {
  host     = "10.2.20.24"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0209"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "apache"
  }
}

server "vm-08ef9cd4ebb49cdd" {
  host     = "10.2.20.25"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0210"
  tags = {
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "nginx"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-e893f4eeebe2550a" {
  host     = "10.2.20.26"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0211"
  tags = {
    tier = "production"
    web_type = "apache"
    datacenter = "BER0"
    type = "vm"
    role = "web"
    os = "debian12"
  }
}

server "vm-d46422f17d8c93d0" {
  host     = "10.2.20.27"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0212"
  tags = {
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
    web_type = "nginx"
    datacenter = "BER0"
  }
}

server "vm-d03be358b3138e80" {
  host     = "10.2.20.28"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0213"
  tags = {
    web_type = "apache"
    datacenter = "BER0"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "production"
  }
}

server "vm-4d6eb731b1ce9676" {
  host     = "10.2.20.29"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0214"
  tags = {
    os = "debian12"
    tier = "production"
    web_type = "nginx"
    datacenter = "BER0"
    type = "vm"
    role = "web"
  }
}

server "vm-32a5cdaf510d9e7a" {
  host     = "10.2.20.30"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0215"
  tags = {
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "staging"
    web_type = "apache"
    datacenter = "BER0"
  }
}

server "vm-2fbe3ca3cbda9cf8" {
  host     = "10.2.20.31"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0216"
  tags = {
    web_type = "nginx"
    datacenter = "BER0"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "staging"
  }
}

server "vm-a7d1e4ed8d7711d0" {
  host     = "10.2.20.32"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0217"
  tags = {
    role = "web"
    os = "debian12"
    tier = "staging"
    web_type = "apache"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-fe8df73b7d7bdd1f" {
  host     = "10.2.20.33"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0218"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "staging"
    web_type = "nginx"
  }
}

server "vm-28e802b8b304733a" {
  host     = "10.2.20.34"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0219"
  tags = {
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "staging"
    web_type = "apache"
    datacenter = "BER0"
  }
}

server "vm-bc700bbe2845a96d" {
  host     = "10.2.20.35"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0220"
  tags = {
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "staging"
    web_type = "nginx"
    datacenter = "BER0"
  }
}

server "vm-d213a895e6423bf7" {
  host     = "10.2.20.36"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0221"
  tags = {
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "staging"
    web_type = "apache"
    datacenter = "BER0"
  }
}

server "vm-257bcb1130549fcd" {
  host     = "10.2.20.37"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0222"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "web"
    os = "debian12"
    tier = "staging"
    web_type = "nginx"
  }
}

server "vm-3f17439f62e204cd" {
  host     = "10.2.30.1"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0223"
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
  password = "vm-secure-pass-0224"
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
  password = "vm-secure-pass-0225"
  tags = {
    tier = "production"
    workload_type = "compute"
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
  }
}

server "vm-d30875a45d508d42" {
  host     = "10.2.30.4"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0226"
  tags = {
    tier = "production"
    workload_type = "batch"
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
  }
}

server "vm-fc082d2ba4f306da" {
  host     = "10.2.30.5"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0227"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "compute"
  }
}

server "vm-e2acfdcb12580438" {
  host     = "10.2.30.6"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0228"
  tags = {
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "batch"
    datacenter = "BER0"
  }
}

server "vm-494954534aadf512" {
  host     = "10.2.30.7"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0229"
  tags = {
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "compute"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-8c82839dcff71483" {
  host     = "10.2.30.8"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0230"
  tags = {
    tier = "production"
    workload_type = "batch"
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
  }
}

server "vm-1f0f650209b4bb23" {
  host     = "10.2.30.9"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0231"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "compute"
  }
}

server "vm-83e5487ed3a847a9" {
  host     = "10.2.30.10"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0232"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "batch"
  }
}

server "vm-0194e154fabde521" {
  host     = "10.2.30.11"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0233"
  tags = {
    tier = "production"
    workload_type = "compute"
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
  }
}

server "vm-bfa6146bd49f1ad9" {
  host     = "10.2.30.12"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0234"
  tags = {
    workload_type = "batch"
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
  }
}

server "vm-79801914f931fff9" {
  host     = "10.2.30.13"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0235"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "compute"
  }
}

server "vm-5149b2a5edb6cb62" {
  host     = "10.2.30.14"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0236"
  tags = {
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "batch"
    datacenter = "BER0"
  }
}

server "vm-3fa443d3a81cf751" {
  host     = "10.2.30.15"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0237"
  tags = {
    os = "debian12"
    tier = "production"
    workload_type = "compute"
    datacenter = "BER0"
    type = "vm"
    role = "workload"
  }
}

server "vm-cbffd09363e9aa70" {
  host     = "10.2.30.16"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0238"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "batch"
  }
}

server "vm-dff0cbe1f69c072a" {
  host     = "10.2.30.17"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0239"
  tags = {
    tier = "production"
    workload_type = "compute"
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
  }
}

server "vm-18415aef908b8589" {
  host     = "10.2.30.18"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0240"
  tags = {
    tier = "production"
    workload_type = "batch"
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
  }
}

server "vm-4ea6ffd36720e89e" {
  host     = "10.2.30.19"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0241"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "compute"
  }
}

server "vm-dcf24143477590aa" {
  host     = "10.2.30.20"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0242"
  tags = {
    os = "debian12"
    tier = "production"
    workload_type = "batch"
    datacenter = "BER0"
    type = "vm"
    role = "workload"
  }
}

server "vm-6bb0a1ec418a98f1" {
  host     = "10.2.30.21"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0243"
  tags = {
    workload_type = "compute"
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
  }
}

server "vm-10e3747bb454c037" {
  host     = "10.2.30.22"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0244"
  tags = {
    os = "debian12"
    tier = "production"
    workload_type = "batch"
    datacenter = "BER0"
    type = "vm"
    role = "workload"
  }
}

server "vm-c0dd18b1732a176a" {
  host     = "10.2.30.23"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0245"
  tags = {
    os = "debian12"
    tier = "production"
    workload_type = "compute"
    datacenter = "BER0"
    type = "vm"
    role = "workload"
  }
}

server "vm-cd9eea06e78ed6cd" {
  host     = "10.2.30.24"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0246"
  tags = {
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "batch"
    datacenter = "BER0"
  }
}

server "vm-d9ebe48efedad86c" {
  host     = "10.2.30.25"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0247"
  tags = {
    tier = "production"
    workload_type = "compute"
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
  }
}

server "vm-0007a70c320ddbbf" {
  host     = "10.2.30.26"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0248"
  tags = {
    tier = "production"
    workload_type = "batch"
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
  }
}

server "vm-1ecf20087fd1a8cb" {
  host     = "10.2.30.27"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0249"
  tags = {
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "compute"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-80251f41cd4a5134" {
  host     = "10.2.30.28"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0250"
  tags = {
    tier = "production"
    workload_type = "batch"
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
  }
}

server "vm-0fd00b345a7d251a" {
  host     = "10.2.30.29"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0251"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "production"
    workload_type = "compute"
  }
}

server "vm-6eb57bcd39cc5cb7" {
  host     = "10.2.30.30"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0252"
  tags = {
    tier = "staging"
    workload_type = "batch"
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
  }
}

server "vm-e657651ae23f5bae" {
  host     = "10.2.30.31"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0253"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "staging"
    workload_type = "compute"
  }
}

server "vm-70a4a2d9b7984d90" {
  host     = "10.2.30.32"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0254"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "staging"
    workload_type = "batch"
  }
}

server "vm-74a11aa4f9977722" {
  host     = "10.2.30.33"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0255"
  tags = {
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "staging"
    workload_type = "compute"
    datacenter = "BER0"
  }
}

server "vm-d941bed40e97532c" {
  host     = "10.2.30.34"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0256"
  tags = {
    workload_type = "batch"
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "staging"
  }
}

server "vm-a5ee17a099db66a5" {
  host     = "10.2.30.35"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0257"
  tags = {
    tier = "staging"
    workload_type = "compute"
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
  }
}

server "vm-82e0e3fa2ef0d20d" {
  host     = "10.2.30.36"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0258"
  tags = {
    role = "workload"
    os = "debian12"
    tier = "staging"
    workload_type = "batch"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-3c311507596c857b" {
  host     = "10.2.30.37"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0259"
  tags = {
    workload_type = "compute"
    datacenter = "BER0"
    type = "vm"
    role = "workload"
    os = "debian12"
    tier = "staging"
  }
}

server "vm-44e2c3366a307954" {
  host     = "10.2.40.1"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0260"
  tags = {
    tier = "production"
    storage_type = "block"
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
  }
}

server "vm-0a7a0f62e44c67d0" {
  host     = "10.2.40.2"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0261"
  tags = {
    storage_type = "object"
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
  }
}

server "vm-1c08dc3b8226a3f2" {
  host     = "10.2.40.3"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0262"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "block"
  }
}

server "vm-ff812ab07f53dc37" {
  host     = "10.2.40.4"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0263"
  tags = {
    storage_type = "object"
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
  }
}

server "vm-5df3c4c88581f39b" {
  host     = "10.2.40.5"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0264"
  tags = {
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "block"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-5a271fd1de8ba29a" {
  host     = "10.2.40.6"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0265"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "object"
  }
}

server "vm-ca6baa580d787977" {
  host     = "10.2.40.7"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0266"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "block"
  }
}

server "vm-8b42cd724a7099db" {
  host     = "10.2.40.8"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0267"
  tags = {
    storage_type = "object"
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
  }
}

server "vm-536b7877b2cb96c5" {
  host     = "10.2.40.9"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0268"
  tags = {
    os = "debian12"
    tier = "production"
    storage_type = "block"
    datacenter = "BER0"
    type = "vm"
    role = "storage"
  }
}

server "vm-697676ac4d6752c3" {
  host     = "10.2.40.10"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0269"
  tags = {
    tier = "production"
    storage_type = "object"
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
  }
}

server "vm-20f4671e2387da4d" {
  host     = "10.2.40.11"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0270"
  tags = {
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "block"
    datacenter = "BER0"
  }
}

server "vm-19ddcf6903afacf4" {
  host     = "10.2.40.12"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0271"
  tags = {
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "object"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-3b366ee0ae48e223" {
  host     = "10.2.40.13"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0272"
  tags = {
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "block"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-3c199a6a97691b63" {
  host     = "10.2.40.14"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0273"
  tags = {
    tier = "production"
    storage_type = "object"
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
  }
}

server "vm-d86f1951146b72f0" {
  host     = "10.2.40.15"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0274"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "block"
  }
}

server "vm-a08e41352ba0f265" {
  host     = "10.2.40.16"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0275"
  tags = {
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "object"
    datacenter = "BER0"
  }
}

server "vm-f33b172a3877a612" {
  host     = "10.2.40.17"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0276"
  tags = {
    storage_type = "block"
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
  }
}

server "vm-da40e02b81a87c5d" {
  host     = "10.2.40.18"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0277"
  tags = {
    tier = "production"
    storage_type = "object"
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
  }
}

server "vm-0ffbd0d326600c51" {
  host     = "10.2.40.19"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0278"
  tags = {
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "block"
    datacenter = "BER0"
  }
}

server "vm-01fbeb6c3be8a939" {
  host     = "10.2.40.20"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0279"
  tags = {
    storage_type = "object"
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
  }
}

server "vm-9c765d067f1a348b" {
  host     = "10.2.40.21"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0280"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "block"
  }
}

server "vm-bbd8375c79d5358f" {
  host     = "10.2.40.22"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0281"
  tags = {
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "object"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-abb7569bbbc39c21" {
  host     = "10.2.40.23"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0282"
  tags = {
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "block"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-9ed07c6766bac427" {
  host     = "10.2.40.24"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0283"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "object"
  }
}

server "vm-d187ec14142d92ef" {
  host     = "10.2.40.25"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0284"
  tags = {
    tier = "production"
    storage_type = "block"
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
  }
}

server "vm-82df5fbb91bf1411" {
  host     = "10.2.40.26"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0285"
  tags = {
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "object"
    datacenter = "BER0"
  }
}

server "vm-f3caa98864e173d6" {
  host     = "10.2.40.27"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0286"
  tags = {
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "block"
    datacenter = "BER0"
    type = "vm"
  }
}

server "vm-31a129fb61b827be" {
  host     = "10.2.40.28"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0287"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "production"
    storage_type = "object"
  }
}

server "vm-95ca7466e63440c5" {
  host     = "10.2.40.29"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0288"
  tags = {
    tier = "production"
    storage_type = "block"
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
  }
}

server "vm-a947e2210198b463" {
  host     = "10.2.40.30"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0289"
  tags = {
    os = "debian12"
    tier = "staging"
    storage_type = "object"
    datacenter = "BER0"
    type = "vm"
    role = "storage"
  }
}

server "vm-18f606efdb9d6495" {
  host     = "10.2.40.31"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0290"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "staging"
    storage_type = "block"
  }
}

server "vm-569bd2141dd72bbb" {
  host     = "10.2.40.32"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0291"
  tags = {
    os = "debian12"
    tier = "staging"
    storage_type = "object"
    datacenter = "BER0"
    type = "vm"
    role = "storage"
  }
}

server "vm-86eed7c17beeadb6" {
  host     = "10.2.40.33"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0292"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "staging"
    storage_type = "block"
  }
}

server "vm-d58e4573b39d276e" {
  host     = "10.2.40.34"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0293"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "staging"
    storage_type = "object"
  }
}

server "vm-951f7860620a3e80" {
  host     = "10.2.40.35"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0294"
  tags = {
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "staging"
    storage_type = "block"
  }
}

server "vm-a3ea4fa92664c899" {
  host     = "10.2.40.36"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0295"
  tags = {
    type = "vm"
    role = "storage"
    os = "debian12"
    tier = "staging"
    storage_type = "object"
    datacenter = "BER0"
  }
}

server "vm-f8851a924915703c" {
  host     = "10.2.40.37"
  port     = 22
  user     = "debian"
  password = "vm-secure-pass-0296"
  tags = {
    tier = "staging"
    storage_type = "block"
    datacenter = "BER0"
    type = "vm"
    role = "storage"
    os = "debian12"
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
  servers     = ["machine-fe8d67b1073acad1", "machine-d901fe5ad5a4c1da", "machine-09d6cdbeec529371"]
  parallel    = true
  timeout     = 300
}

