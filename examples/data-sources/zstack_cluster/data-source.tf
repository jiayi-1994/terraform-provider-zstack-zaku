# 查询单个集群的详细信息
data "zstack_cluster" "example" {
  id = 1
}

# 输出集群详细信息
output "cluster_details" {
  value = {
    id             = data.zstack_cluster.example.id
    name           = data.zstack_cluster.example.name
    status         = data.zstack_cluster.example.status
    version        = data.zstack_cluster.example.version
    node_count     = data.zstack_cluster.example.node_count
    create_time    = data.zstack_cluster.example.create_time
    prometheus_url = data.zstack_cluster.example.prometheus_url
    cpu_usage      = data.zstack_cluster.example.cpu_usage
    memory_usage   = data.zstack_cluster.example.memory_usage
    storage_usage  = data.zstack_cluster.example.storage_usage
  }
}
