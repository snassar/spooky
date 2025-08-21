# Spooky Configuration File
# Generated automatically on first run

# Logging configuration
logging {
  level  = "info"
  format = "json"
  output = "file"
  
  file {
    path = "${log_file}"
  }
  
  structured {
    include_timestamp = true
    include_level     = true
    include_source    = false
  }
}

# Project configuration
project {
  default_timeout = "30m"
  max_parallel    = 4
  retry_attempts  = 3
}

# SSH configuration
ssh {
  default_user     = ""
  default_port     = 22
  connection_timeout = "10s"
  keepalive_interval = "30s"
}

# Template configuration
templates {
  default_engine = "go"
  cache_enabled  = true
  cache_ttl      = "1h"
}

# Variables configuration
variables {
  auto_load = true
  env_prefix = "SPOOKY_"
}

# Facts configuration
facts {
  collection_timeout = "5m"
  cache_enabled     = true
  cache_ttl         = "1h"
}

# Actions configuration
actions {
  default_timeout = "10m"
  max_retries     = 3
  parallel_limit  = 10
}

# Machines configuration
machines {
  inventory_timeout = "30s"
  connection_pool_size = 10
  max_connections_per_host = 5
}
