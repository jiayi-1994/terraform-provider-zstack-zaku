# Terraform Provider for ZStack Edge

基于 [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) 构建的 ZStack Edge Terraform Provider，用于管理 Kubernetes 集群资源。

## 功能特性

- ✅ **集群管理** - 完整的 Kubernetes 集群生命周期管理（创建、读取、更新、删除）
- ✅ **高可用支持** - 支持创建高可用 Kubernetes 集群
- ✅ **GPU 集群** - 支持天数 GPU 节点配置
- ✅ **数据查询** - 查询集群详情和集群列表
- ✅ **导入功能** - 导入已有集群到 Terraform 管理
- ✅ **完整文档** - 详细的使用文档和示例

## 快速开始

### 安装

```hcl
terraform {
  required_providers {
    zstack = {
      source  = "registry.terraform.io/zstack/zstack-zaku"
      version = "~> 1.0"
    }
  }
}
```

### 配置 Provider

```hcl
provider "zstack" {
  host       = "https://your-zstack-edge-host.com"
  access_key = "your-access-key"
  secret_key = "your-secret-key"
}
```

### 创建集群

```hcl
resource "zstack_cluster" "example" {
  name             = "my-k8s-cluster"
  port             = 22
  password         = var.encrypted_password
  
  management_vip_v4 = "172.31.13.100"
  business_vip_v4   = "172.32.4.100"
  pod_cidr_v4      = "10.233.64.0/18"
  service_cidr_v4  = "10.233.0.0/18"
  dns_server       = "223.5.5.5"
  
  nodes = [
    {
      name                  = "master-1"
      roles                 = ["Master", "Worker"]
      management_ipv4_addr  = "172.31.13.101"
      business_ipv4_addr    = "172.32.4.101"
    }
  ]
  
  data_disk = {
    "master-1" = ["/dev/sdb"]
  }
}
```

## 文档

- 📖 [完整功能文档](CLUSTER_PROVIDER_README.md) - 详细的功能说明和使用指南
- 🚀 [Hermitcrab 快速上传](QUICKSTART_HERMITCRAB.md) - 上传到私有仓库快速指南
- 📦 [Hermitcrab 完整指南](UPLOAD_TO_HERMITCRAB.md) - 详细的私有仓库部署说明
- 📋 [Manifest 文件生成](MANIFEST_GENERATION.md) - Registry manifest 文件生成指南
- 🧪 [测试指南](test/README.md) - 本地测试和开发指南
- 📊 [项目总结](PROJECT_SUMMARY.md) - 项目架构和开发总结
- 📝 [示例代码](examples/) - 各种使用场景的示例

## 支持的资源

### Resources

- `zstack_cluster` - Kubernetes 集群管理

### Data Sources

- `zstack_cluster` - 查询单个集群详情
- `zstack_clusters` - 查询集群列表

## 开发

### 构建 Provider

```bash
go build -o terraform-provider-zstack.exe
```

### 本地测试

1. 配置开发环境（`.terraformrc`）：

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/zstack/zstack-zaku" = "F:/other-code/terraform-lean/terraform-provider-zstack"
  }
  direct {}
}
```

2. 运行测试：

```bash
cd test
terraform init
terraform plan
terraform apply
```

详细测试步骤请参考 [test/README.md](test/README.md)。

### 发布到私有仓库

使用 Hermitcrab 部署私有 Terraform Registry：

```bash
# Windows
.\scripts\upload-to-hermitcrab.ps1 -Host your-host -Port 5000

# Linux/Mac
./scripts/upload-to-hermitcrab.sh -H your-host -p 5000
```

详细说明请参考：
- [快速开始](QUICKSTART_HERMITCRAB.md)
- [完整指南](UPLOAD_TO_HERMITCRAB.md)

## 项目结构

```
terraform-provider-zstack/
├── internal/provider/          # Provider 实现
│   ├── provider.go            # Provider 配置
│   ├── cluster_resource.go    # 集群资源
│   ├── cluster_data_source.go # 集群数据源
│   └── clusters_data_source.go # 集群列表数据源
├── examples/                   # 示例代码
├── test/                       # 测试配置
├── vendor/                     # 依赖库
└── docs/                       # 生成的文档
```

## 系统要求

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## 技术栈

- **开发框架**: Terraform Plugin Framework
- **SDK**: ZStack Edge Go SDK
- **语言**: Go 1.21+
- **构建工具**: Go Modules

## 贡献

欢迎提交 Issue 和 Pull Request！

### 开发流程

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 许可证

本项目采用 MPL-2.0 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 相关链接

- [ZStack Edge 官方文档](https://docs.zstack.io/)
- [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
- [Terraform Registry](https://registry.terraform.io/)
