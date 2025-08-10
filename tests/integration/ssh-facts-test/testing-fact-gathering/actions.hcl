# Test actions for SSH fact gathering

actions {
  action "test-ssh-facts" {
    description = "Test SSH fact gathering on VM"
    type = "command"
    command = "echo 'Testing SSH fact gathering'"
    
    # Target the test VM
    machines = ["spooky-facts-test-vm"]
    
    # Execution settings
    timeout = 300
    parallel = false
    dry_run = false
  }

  action "validate-facts" {
    description = "Validate collected facts"
    type = "command"
    command = "echo 'Validating facts'"
    
    machines = ["spooky-facts-test-vm"]
    
    # Execution settings
    timeout = 300
    parallel = false
    dry_run = false
  }
}
