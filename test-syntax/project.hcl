project {
  facts_retry_delay         = 5
  name                      = "test-syntax"
  run_max_parallel          = 10
  run_validate_before_run   = true
  facts_retry_attempts      = 3
  facts_timeout             = 30
  run_backup_before_changes = false
  run_dry_run_default       = false
  age {
    default_identities_path = ""
    default_recipients_path = ""
  }
  description               = "Testing HCL syntax"
  facts_parallel_collection = 10
}
