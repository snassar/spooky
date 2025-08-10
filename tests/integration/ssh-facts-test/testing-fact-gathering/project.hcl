project {
  name = "testing-fact-gathering"
  description = "Test project for SSH fact gathering functionality"
  version = "1.0.0"
  environment = "testing"
  
  # Project tags
  tags = ["test", "ssh", "facts", "vm"]
  
  # Project structure
  structure {
    templates_dir = "templates"
    data_dir = "data"
    logs_dir = "logs"
  }
} 