machines {
  # Kubernetes Control Plane Nodes
  machine "k8s-master-01" {
    host = "10.0.10.10"
    user = "kubernetes"
    port = 22
    key_file = "~/.ssh/k8s_master_key"
    
    tags = ["kubernetes", "control-plane", "production"]
    groups = ["k8s-masters", "production-servers"]
    roles = ["kubernetes-master", "etcd", "api-server", "scheduler", "controller-manager"]
    
    resources {
      cpu_cores = 8
      memory_gb = 32
      disk_gb = 500
      network_mbps = 10000
    }
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      rack = "K-01"
      location = "San Francisco"
      owner = "platform-team"
      department = "Engineering"
      cost_center = "IT-010"
      maintenance_window = "Sunday 2-4 AM PST"
      backup_schedule = "daily"
      monitoring = "prometheus"
      alerting = "pagerduty"
      sla = "99.9%"
      kubernetes_version = "1.28.0"
      node_role = "master"
      etcd_member = "true"
    }
  }
  
  machine "k8s-master-02" {
    host = "10.0.10.11"
    user = "kubernetes"
    port = 22
    key_file = "~/.ssh/k8s_master_key"
    
    tags = ["kubernetes", "control-plane", "production"]
    groups = ["k8s-masters", "production-servers"]
    roles = ["kubernetes-master", "etcd", "api-server", "scheduler", "controller-manager"]
    
    resources {
      cpu_cores = 8
      memory_gb = 32
      disk_gb = 500
      network_mbps = 10000
    }
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      rack = "K-02"
      location = "San Francisco"
      owner = "platform-team"
      department = "Engineering"
      cost_center = "IT-010"
      maintenance_window = "Sunday 2-4 AM PST"
      backup_schedule = "daily"
      monitoring = "prometheus"
      alerting = "pagerduty"
      sla = "99.9%"
      kubernetes_version = "1.28.0"
      node_role = "master"
      etcd_member = "true"
    }
  }
  
  machine "k8s-master-03" {
    host = "10.0.10.12"
    user = "kubernetes"
    port = 22
    key_file = "~/.ssh/k8s_master_key"
    
    tags = ["kubernetes", "control-plane", "production"]
    groups = ["k8s-masters", "production-servers"]
    roles = ["kubernetes-master", "etcd", "api-server", "scheduler", "controller-manager"]
    
    resources {
      cpu_cores = 8
      memory_gb = 32
      disk_gb = 500
      network_mbps = 10000
    }
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      rack = "K-03"
      location = "San Francisco"
      owner = "platform-team"
      department = "Engineering"
      cost_center = "IT-010"
      maintenance_window = "Sunday 2-4 AM PST"
      backup_schedule = "daily"
      monitoring = "prometheus"
      alerting = "pagerduty"
      sla = "99.9%"
      kubernetes_version = "1.28.0"
      node_role = "master"
      etcd_member = "true"
    }
  }
  
  # Kubernetes Worker Nodes
  machine "k8s-worker-01" {
    host = "10.0.11.10"
    user = "kubernetes"
    port = 22
    key_file = "~/.ssh/k8s_worker_key"
    
    tags = ["kubernetes", "worker", "production"]
    groups = ["k8s-workers", "production-servers"]
    roles = ["kubernetes-worker", "container-runtime", "kubelet", "kube-proxy"]
    
    resources {
      cpu_cores = 16
      memory_gb = 64
      disk_gb = 1000
      network_mbps = 10000
    }
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      rack = "K-11"
      location = "San Francisco"
      owner = "platform-team"
      department = "Engineering"
      cost_center = "IT-011"
      maintenance_window = "Sunday 2-4 AM PST"
      backup_schedule = "daily"
      monitoring = "prometheus"
      alerting = "pagerduty"
      sla = "99.9%"
      kubernetes_version = "1.28.0"
      node_role = "worker"
      node_labels = "node-type=general,zone=us-west-1a"
      taints = ""
    }
  }
  
  machine "k8s-worker-02" {
    host = "10.0.11.11"
    user = "kubernetes"
    port = 22
    key_file = "~/.ssh/k8s_worker_key"
    
    tags = ["kubernetes", "worker", "production"]
    groups = ["k8s-workers", "production-servers"]
    roles = ["kubernetes-worker", "container-runtime", "kubelet", "kube-proxy"]
    
    resources {
      cpu_cores = 16
      memory_gb = 64
      disk_gb = 1000
      network_mbps = 10000
    }
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      rack = "K-12"
      location = "San Francisco"
      owner = "platform-team"
      department = "Engineering"
      cost_center = "IT-011"
      maintenance_window = "Sunday 2-4 AM PST"
      backup_schedule = "daily"
      monitoring = "prometheus"
      alerting = "pagerduty"
      sla = "99.9%"
      kubernetes_version = "1.28.0"
      node_role = "worker"
      node_labels = "node-type=general,zone=us-west-1a"
      taints = ""
    }
  }
  
  machine "k8s-worker-03" {
    host = "10.0.11.12"
    user = "kubernetes"
    port = 22
    key_file = "~/.ssh/k8s_worker_key"
    
    tags = ["kubernetes", "worker", "production", "gpu"]
    groups = ["k8s-workers", "gpu-nodes", "production-servers"]
    roles = ["kubernetes-worker", "container-runtime", "kubelet", "kube-proxy", "gpu-accelerator"]
    
    resources {
      cpu_cores = 32
      memory_gb = 128
      disk_gb = 2000
      network_mbps = 10000
    }
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      rack = "K-13"
      location = "San Francisco"
      owner = "platform-team"
      department = "Engineering"
      cost_center = "IT-012"
      maintenance_window = "Sunday 2-4 AM PST"
      backup_schedule = "daily"
      monitoring = "prometheus"
      alerting = "pagerduty"
      sla = "99.9%"
      kubernetes_version = "1.28.0"
      node_role = "worker"
      node_labels = "node-type=gpu,zone=us-west-1a,gpu-type=nvidia-v100"
      taints = "nvidia.com/gpu=true:NoSchedule"
      gpu_count = "4"
      gpu_type = "nvidia-v100"
    }
  }
  
  # Kubernetes Infrastructure Nodes
  machine "k8s-infra-01" {
    host = "10.0.12.10"
    user = "kubernetes"
    port = 22
    key_file = "~/.ssh/k8s_infra_key"
    
    tags = ["kubernetes", "infrastructure", "production"]
    groups = ["k8s-infra", "production-servers"]
    roles = ["kubernetes-worker", "monitoring", "logging", "ingress"]
    
    resources {
      cpu_cores = 8
      memory_gb = 32
      disk_gb = 1000
      network_mbps = 10000
    }
    
    metadata {
      environment = "production"
      datacenter = "us-west-1"
      rack = "K-21"
      location = "San Francisco"
      owner = "platform-team"
      department = "Engineering"
      cost_center = "IT-013"
      maintenance_window = "Sunday 2-4 AM PST"
      backup_schedule = "daily"
      monitoring = "prometheus"
      alerting = "pagerduty"
      sla = "99.9%"
      kubernetes_version = "1.28.0"
      node_role = "worker"
      node_labels = "node-type=infrastructure,zone=us-west-1a"
      taints = "node-role.kubernetes.io/infrastructure=true:NoSchedule"
      dedicated_infrastructure = "true"
    }
  }
}
