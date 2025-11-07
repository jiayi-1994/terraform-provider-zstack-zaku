# 使用私有 Hermitcrab 仓库示例

这个目录包含了如何配置和使用 Hermitcrab 私有仓库的示例。

## 文件说明

- `.terraformrc.example` - Terraform CLI 配置文件示例
- `main.tf` - 使用私有仓库的 Terraform 配置示例

## 快速开始

### 1. 配置 Terraform CLI

将 `.terraformrc.example` 复制并修改为实际配置：

**Windows:**
```powershell
# 复制到用户配置目录
Copy-Item .terraformrc.example $env:APPDATA\terraform.rc

# 或在项目中使用
Copy-Item .terraformrc.example .\.terraformrc
$env:TF_CLI_CONFIG_FILE = ".\.terraformrc"
```

**Linux/Mac:**
```bash
# 复制到用户配置目录
cp .terraformrc.example ~/.terraformrc

# 或在项目中使用
cp .terraformrc.example ./.terraformrc
export TF_CLI_CONFIG_FILE="./.terraformrc"
```

### 2. 修改配置

编辑 `.terraformrc` 文件，将 `your-hermitcrab-host:5000` 替换为实际的 Hermitcrab 地址：

```hcl
provider_installation {
  network_mirror {
    url = "http://your-hermitcrab-host:5000/v1/providers/"
    include = ["registry.terraform.io/zstack/*"]
  }
  
  direct {
    exclude = ["registry.terraform.io/zstack/*"]
  }
}
```

### 3. 初始化 Terraform

```bash
terraform init
```

预期输出：
```
Initializing provider plugins...
- Finding registry.terraform.io/zstack/zstack-zaku versions matching "1.0.0"...
- Installing registry.terraform.io/zstack/zstack-zaku v1.0.0...
- Installed registry.terraform.io/zstack/zstack-zaku v1.0.0 (unauthenticated)
```

### 4. 使用 Provider

编辑 `main.tf`，添加你的资源配置：

```hcl
provider "zstack-zaku" {
  # 配置你的连接参数
}

# 添加你的资源
```

## 配置选项

### 网络镜像模式

最简单的配置方式，适合大多数场景：

```hcl
provider_installation {
  network_mirror {
    url = "http://hermitcrab-host:5000/v1/providers/"
  }
}
```

### 选择性使用私有仓库

只对特定 namespace 使用私有仓库：

```hcl
provider_installation {
  network_mirror {
    url = "http://hermitcrab-host:5000/v1/providers/"
    include = ["registry.terraform.io/zstack/*"]
  }
  
  # 其他 provider 使用官方源
  direct {
    exclude = ["registry.terraform.io/zstack/*"]
  }
}
```

### 开发模式

在本地开发 provider 时使用：

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/zstack/zstack-zaku" = "/path/to/provider/build"
  }
  
  direct {}
}
```

**注意**: 使用 `dev_overrides` 时，Terraform 会跳过版本检查并使用本地二进制文件。

## 环境变量

### TF_CLI_CONFIG_FILE

指定自定义配置文件路径：

```bash
# Linux/Mac
export TF_CLI_CONFIG_FILE="/path/to/.terraformrc"

# Windows
$env:TF_CLI_CONFIG_FILE = "C:\path\to\.terraformrc"
```

### HERMITCRAB_HOST (自定义)

如果你的上传脚本使用环境变量：

```bash
# Linux/Mac
export HERMITCRAB_HOST="hermitcrab.example.com"
export HERMITCRAB_PORT="5000"

# Windows
$env:HERMITCRAB_HOST = "hermitcrab.example.com"
$env:HERMITCRAB_PORT = "5000"
```

## 验证配置

### 检查 Provider 是否可访问

```bash
# 使用 curl 测试
curl http://hermitcrab-host:5000/v1/providers/zstack/zstack-zaku/versions

# 使用 PowerShell (Windows)
Invoke-WebRequest -Uri "http://hermitcrab-host:5000/v1/providers/zstack/zstack-zaku/versions"
```

### 查看 Terraform 配置

```bash
terraform version -json
```

### 测试 Provider 安装

```bash
# 强制重新安装
rm -rf .terraform .terraform.lock.hcl
terraform init
```

## 故障排除

### Provider 无法下载

1. **检查网络连接**
   ```bash
   curl http://hermitcrab-host:5000/health
   ```

2. **检查 Provider 是否已上传**
   ```bash
   curl http://hermitcrab-host:5000/v1/providers/zstack/zstack-zaku/versions
   ```

3. **检查 .terraformrc 配置**
   ```bash
   cat ~/.terraformrc  # Linux/Mac
   type %APPDATA%\terraform.rc  # Windows
   ```

### 版本不匹配

确保 `main.tf` 中的版本号与上传的版本一致：

```hcl
required_providers {
  zstack-zaku = {
    source  = "registry.terraform.io/zstack/zstack-zaku"
    version = "1.0.0"  # 与上传的版本一致
  }
}
```

### 校验和错误

如果遇到校验和错误，重新生成并上传 SHA256SUMS：

```bash
cd dist
shasum -a 256 *.zip > terraform-provider-zstack-zaku_1.0.0_SHA256SUMS
# 重新上传
```

## 最佳实践

1. **使用版本约束**: 在 `required_providers` 中明确指定版本
2. **环境分离**: 为不同环境使用不同的 `.terraformrc` 配置
3. **CI/CD 集成**: 在 CI/CD 流程中设置 `TF_CLI_CONFIG_FILE`
4. **缓存配置**: 考虑在 `.terraformrc` 中配置 plugin cache
5. **安全性**: 如果 Hermitcrab 需要认证，使用环境变量而不是硬编码

## 参考资料

- [Terraform CLI Configuration](https://www.terraform.io/cli/config/config-file)
- [Provider Installation Methods](https://www.terraform.io/cli/config/config-file#provider-installation)
- [Network Mirror Protocol](https://www.terraform.io/internals/provider-network-mirror-protocol)
- [Hermitcrab Documentation](https://github.com/seal-io/hermitcrab)
