# Test VM for SSH fact gathering
machines {
  machine "spooky-facts-test-vm" {
    host = "localhost"
    port = 2222
    user = "spooky"
    key_file = "../keys/spooky-facts-test_key"
    
    # Connection timeouts
    connection_timeout = 30
    command_timeout = 300
    
    # Tags for targeting and organization
    tags = {
      environment = "test"
      purpose = "ssh-facts-testing"
      vm_type = "qemu"
      os_family = "debian"
    }
    
    # Machine groups
    groups = ["test", "vm", "ssh", "facts", "debian"]
  }
} 