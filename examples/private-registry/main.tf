# 使用私有 Hermitcrab 仓库的 Terraform 配置示例

terraform {
  required_version = ">= 1.0"

  required_providers {
    # 从私有仓库使用 zstack-zaku provider
    zstack-zaku = {
      # 方法 1: 使用自定义 hostname (需要配置 .terraformrc)
      source  = "registry.terraform.io/zstack/zstack-zaku"
      version = "1.0.0"
      
      # 方法 2: 如果 Hermitcrab 配置了自定义域名
      # source  = "terraform.company.com/zstack/zstack-zaku"
      # version = "1.0.0"
    }
  }
}

# Provider 配置
provider "zstack-zaku" {
  # 根据实际情况配置 provider 参数
  # 例如：
  # api_url  = "https://zstack.example.com"
  # username = var.zstack_username
  # password = var.zstack_password
}

# 使用 provider 的示例资源
# data "zstack_cluster" "example" {
#   # 数据源配置
# }

# resource "zstack_cluster" "example" {
#   # 资源配置
# }
