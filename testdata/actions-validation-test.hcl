# Test actions.hcl file to demonstrate enhanced validation
# This file contains various action types with different validation scenarios

actions {
  # Valid command action
  action "test_command" {
    description = "A simple command action"
    type        = "command"
    command     = "echo 'Hello World'"
    timeout     = 30
  }

  # Valid script action
  action "test_script" {
    description = "A script action with file reference"
    type        = "script"
    script      = "files/test_script.sh"
    timeout     = 60
  }

  # Valid templated script with variables
  action "test_templated_script" {
    description = "A templated script with variables"
    type        = "script"
    script      = "templates/test_script.tmpl"
    variables = {
      "name"    = "World"
      "version" = "1.0.0"
    }
    timeout = 90
  }

  # Valid template deploy action
  action "test_template_deploy" {
    description = "A template deployment action"
    type        = "template_deploy"
    template {
      source      = "templates/config.tmpl"
      destination = "/etc/app/config.conf"
      backup      = true
      permissions = "644"
    }
    timeout = 120
  }

  # Valid file copy action
  action "test_file_copy" {
    description = "A file copy action"
    type        = "file_copy"
    file_copy {
      source      = "files/config.json"
      destination = "/etc/app/config.json"
      backup      = true
      permissions = "644"
    }
    timeout = 60
  }

  # Action with dependencies and execution control
  action "test_dependent_action" {
    description = "An action with dependencies"
    type        = "command"
    command     = "systemctl restart service"
    dependencies = ["test_command", "test_script"]
    timeout     = 45
    critical    = true
    retries     = 3
    retry_delay = 10
  }

  # Action with environment and context
  action "test_environment_action" {
    description = "An action with environment variables"
    type        = "script"
    script      = "files/env_script.sh"
    environment = {
      "DEBUG"     = "true"
      "LOG_LEVEL" = "info"
    }
    working_directory = "/tmp"
    user             = "appuser"
    sudo             = false
    timeout          = 30
  }

  # Action with resource limits
  action "test_resource_action" {
    description = "An action with resource limits"
    type        = "script"
    script      = "files/resource_script.sh"
    timeout     = 300
    resource_limits {
      memory_mb   = 512
      cpu_percent = 50
      disk_mb     = 100
    }
  }
}
