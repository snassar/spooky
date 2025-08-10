# Invalid actions.hcl file to test enhanced validation
# This file contains various validation errors

actions {
  # Missing required fields
  action "invalid_action" {
    description = "An action with missing required fields"
    # Missing type and command/script
  }

  # Invalid action type
  action "invalid_type" {
    description = "An action with invalid type"
    type        = "invalid_type"
    command     = "echo 'test'"
  }

  # Script action without script field
  action "missing_script" {
    description = "A script action without script field"
    type        = "script"
    # Missing script field
  }

  # Template action without template block
  action "missing_template" {
    description = "A template action without template block"
    type        = "template_deploy"
    # Missing template block
  }

  # File copy action without file_copy block
  action "missing_file_copy" {
    description = "A file copy action without file_copy block"
    type        = "file_copy"
    # Missing file_copy block
  }

  # Variables on non-script action
  action "invalid_variables" {
    description = "Variables on non-script action"
    type        = "command"
    command     = "echo 'test'"
    variables = {
      "test" = "value"
    }
  }

  # Variables on non-templated script
  action "invalid_variables_script" {
    description = "Variables on non-templated script"
    type        = "script"
    script      = "files/test_script.sh"
    variables = {
      "test" = "value"
    }
  }

  # Script with invalid path format
  action "invalid_script_path" {
    description = "Script with invalid path format"
    type        = "script"
    script      = "invalid/path/script.sh"
  }

  # Template with invalid source path
  action "invalid_template_source" {
    description = "Template with invalid source path"
    type        = "template_deploy"
    template {
      source      = "invalid/path/config.tmpl"
      destination = "/etc/app/config.conf"
    }
  }

  # File copy with invalid source path
  action "invalid_file_copy_source" {
    description = "File copy with invalid source path"
    type        = "file_copy"
    file_copy {
      source      = "invalid/path/config.json"
      destination = "/etc/app/config.json"
    }
  }

  # Command with dangerous shell operators
  action "dangerous_command" {
    description = "Command with dangerous shell operators"
    type        = "command"
    command     = "echo 'test' && rm -rf /"
  }
}
