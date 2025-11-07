# 创建一个 ZStack Edge 集群
resource "zstack_cluster" "example" {
  name             = "my-k8s-cluster"
  enable_ha        = true
  net_combined     = false
  port             = 22
  password         = "encrypted_password_here"
  
  # 网络配置
  management_vip_v4 = "172.31.13.100"
  business_vip_v4   = "172.32.4.100"
  
  # Kubernetes 配置
  pod_cidr_v4      = "10.233.64.0/18"
  service_cidr_v4  = "10.233.0.0/18"
  dns_server       = "223.5.5.5"
  k8s_version      = "1.24"
  max_pod_per_node = 110
  
  # 可选：启用 Istio
  istio_enabled = true
  
  # 可选：天数 GPU 配置
  iluvatar_gpu_model = "BI-V100"
  iluvatar_license   = "your_license_key_here"
  
  # 集群节点配置
  nodes = [
    {
      name                  = "master-node-1"
      roles                 = ["Master"]
      management_ipv4_addr  = "172.31.13.101"
      business_ipv4_addr    = "172.32.4.101"
    },
    {
      name                  = "worker-node-1"
      roles                 = ["Worker"]
      management_ipv4_addr  = "172.31.13.102"
      business_ipv4_addr    = "172.32.4.102"
    },
    {
      name                  = "worker-node-2"
      roles                 = ["Worker"]
      management_ipv4_addr  = "172.31.13.103"
      business_ipv4_addr    = "172.32.4.103"
    }
  ]
  
  # 数据磁盘配置
  data_disk = {
    "master-node-1" = ["/dev/sdb"]
    "worker-node-1" = ["/dev/sdb", "/dev/sdc"]
    "worker-node-2" = ["/dev/sdb", "/dev/sdc"]
  }
  
  # 可选：镜像数据磁盘配置
  image_data_disk = {
    "worker-node-1" = ["/dev/sdd"]
    "worker-node-2" = ["/dev/sdd"]
  }
}

# 输出集群信息
output "cluster_id" {
  value       = zstack_cluster.example.id
  description = "集群 ID"
}

output "cluster_status" {
  value       = zstack_cluster.example.status
  description = "集群状态"
}

output "cluster_prometheus_url" {
  value       = zstack_cluster.example.prometheus_url
  description = "Prometheus 监控地址"
}
