terraform {
  required_providers {
    zstack = {
      source  = "registry.terraform.io/zstack/zstack-zaku"
      version = "~> 1.0"
    }
  }
}

# 配置 ZStack Edge Provider
provider "zstack" {
  host       = "https://your-zstack-edge-host.com"
  access_key = "your-access-key"
  secret_key = "your-secret-key"
  
  # 或者使用环境变量:
  # ZSTACK_HOST
  # ZSTACK_ACCESS_KEY
  # ZSTACK_SECRET_KEY
}
