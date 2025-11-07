# 查询所有集群列表
data "zstack_clusters" "all" {
  limit  = 20
  offset = 0
}

# 输出集群列表
output "all_clusters" {
  value = data.zstack_clusters.all.clusters
}

output "total_clusters" {
  value       = data.zstack_clusters.all.total
  description = "集群总数"
}

# 输出第一个集群的信息（如果存在）
output "first_cluster" {
  value = length(data.zstack_clusters.all.clusters) > 0 ? {
    id     = data.zstack_clusters.all.clusters[0].id
    name   = data.zstack_clusters.all.clusters[0].name
    status = data.zstack_clusters.all.clusters[0].status
  } : null
  description = "第一个集群的基本信息"
}
