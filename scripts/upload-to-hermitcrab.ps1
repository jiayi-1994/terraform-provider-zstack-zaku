# upload-to-hermitcrab.ps1
# 自动上传 Terraform Provider 到 Hermitcrab 私有仓库 (PowerShell 版本)

param(
    [string]$ServerHost = $env:HERMITCRAB_HOST,
    [string]$Port = $env:HERMITCRAB_PORT,
    [string]$Namespace = $env:PROVIDER_NAMESPACE,
    [string]$Type = $env:PROVIDER_TYPE,
    [string]$Version = $env:PROVIDER_VERSION,
    [string]$DistDir = $env:DIST_DIR,
    [string]$Method = "api",
    [string]$DataDir = $env:HERMITCRAB_DATA_DIR,
    [string]$Protocol = $env:HERMITCRAB_PROTOCOL,
    [switch]$UseHttps,
    [switch]$SkipCertificateCheck,
    [switch]$Help
)

# 默认值
if ([string]::IsNullOrEmpty($ServerHost)) { $ServerHost = "172.26.50.232" }
if ([string]::IsNullOrEmpty($Port)) { $Port = "443" }
if ([string]::IsNullOrEmpty($Namespace)) { $Namespace = "zstack" }
if ([string]::IsNullOrEmpty($Type)) { $Type = "zstack-zaku" }
if ([string]::IsNullOrEmpty($Version)) { $Version = "1.0.0" }
if ([string]::IsNullOrEmpty($DistDir)) { $DistDir = "dist" }
if ([string]::IsNullOrEmpty($DataDir)) { $DataDir = "/var/lib/hermitcrab" }

# 确定协议
if ([string]::IsNullOrEmpty($Protocol)) {
    if ($UseHttps) {
        $Protocol = "https"
    } elseif ($Port -eq "443" -or $Port -eq 443) {
        $Protocol = "https"  # 端口 443 默认使用 HTTPS
    } else {
        $Protocol = "http"
    }
} else {
    $Protocol = $Protocol.ToLower()
}

# 帮助信息
function Show-Help {
    Write-Host @"
Usage: .\upload-to-hermitcrab.ps1 [OPTIONS]

上传 Terraform Provider 到 Hermitcrab 私有仓库

OPTIONS:
    -ServerHost <HOST>      Hermitcrab 主机地址 (默认: 172.26.50.232)
    -Port <PORT>            Hermitcrab 端口 (默认: 443)
    -Namespace <NS>         Provider 命名空间 (默认: zstack)
    -Type <TYPE>            Provider 类型 (默认: zstack-zaku)
    -Version <VERSION>      Provider 版本 (默认: 1.0.0)
    -DistDir <DIR>          构建产物目录 (默认: dist)
    -Method <METHOD>        上传方法: api 或 filesystem (默认: api)
    -DataDir <DIR>          Hermitcrab 数据目录 (filesystem 模式使用)
    -Protocol <PROTOCOL>    协议: http 或 https (默认: 根据端口自动判断)
    -UseHttps               强制使用 HTTPS
    -SkipCertificateCheck   跳过 SSL 证书验证 (用于自签名证书)
    -Help                   显示帮助信息

ENVIRONMENT VARIABLES:
    HERMITCRAB_HOST         Hermitcrab 主机地址
    HERMITCRAB_PORT         Hermitcrab 端口
    PROVIDER_NAMESPACE      Provider 命名空间
    PROVIDER_TYPE           Provider 类型
    PROVIDER_VERSION        Provider 版本
    DIST_DIR                构建产物目录
    HERMITCRAB_DATA_DIR     Hermitcrab 数据目录
    HERMITCRAB_PROTOCOL     协议 (http 或 https)

EXAMPLES:
    # 使用环境变量
    `$env:HERMITCRAB_HOST = "registry.example.com"
    `$env:HERMITCRAB_PORT = "5000"
    .\upload-to-hermitcrab.ps1

    # 使用命令行参数
    .\upload-to-hermitcrab.ps1 -ServerHost registry.example.com -Port 5000 -Version 1.0.1

    # 使用 HTTPS
    .\upload-to-hermitcrab.ps1 -ServerHost registry.example.com -Port 443 -UseHttps

    # 使用 HTTPS 并跳过证书验证（自签名证书）
    .\upload-to-hermitcrab.ps1 -ServerHost 172.26.50.232 -Port 443 -SkipCertificateCheck

    # 使用文件系统模式
    .\upload-to-hermitcrab.ps1 -Method filesystem -DataDir "C:\hermitcrab\data"

"@
}

if ($Help) {
    Show-Help
    exit 0
}

# 打印配置
Write-Host "=== Hermitcrab Provider Upload ===" -ForegroundColor Green
Write-Host "Upload Method:    $Method"
Write-Host "Protocol:         $Protocol"
Write-Host "Hermitcrab Host:  $ServerHost"
Write-Host "Hermitcrab Port:  $Port"
Write-Host "Namespace:        $Namespace"
Write-Host "Provider Type:    $Type"
Write-Host "Version:          $Version"
Write-Host "Dist Directory:   $DistDir"

if ($Method -eq "filesystem") {
    Write-Host "Data Directory:   $DataDir"
}

Write-Host ""

# 检查 dist 目录是否存在
if (-not (Test-Path -Path $DistDir -PathType Container)) {
    Write-Host "Error: Dist directory '$DistDir' not found" -ForegroundColor Red
    exit 1
}

# 函数：通过 API 上传文件
function Upload-ViaApi {
    param(
        [string]$FilePath
    )
    
    $FileName = Split-Path -Leaf $FilePath
    $Url = "${Protocol}://${ServerHost}:${Port}/v1/providers/${Namespace}/${Type}/${Version}/upload"
    
    Write-Host "Uploading $FileName via API..." -ForegroundColor Yellow
    
    try {
        # 准备表单数据
        $boundary = [System.Guid]::NewGuid().ToString()
        $LF = "`r`n"
        
        $fileBytes = [System.IO.File]::ReadAllBytes($FilePath)
        
        $bodyLines = @(
            "--$boundary",
            "Content-Disposition: form-data; name=`"file`"; filename=`"$FileName`"",
            "Content-Type: application/octet-stream",
            "",
            [System.Text.Encoding]::GetEncoding("ISO-8859-1").GetString($fileBytes),
            "--$boundary",
            "Content-Disposition: form-data; name=`"filename`"",
            "",
            $FileName
        )
        
        # 如果是 zip 文件，添加 platform 信息
        if ($FileName -match '\.zip$') {
            if ($FileName -match '_([^_]+_[^_]+)\.zip$') {
                $platform = $matches[1]
                $bodyLines += @(
                    "--$boundary",
                    "Content-Disposition: form-data; name=`"platform`"",
                    "",
                    $platform
                )
            }
        }
        
        $bodyLines += "--$boundary--"
        
        $body = $bodyLines -join $LF
        $bodyBytes = [System.Text.Encoding]::GetEncoding("ISO-8859-1").GetBytes($body)
        
        # 发送请求
        $requestParams = @{
            Uri = $Url
            Method = 'Post'
            ContentType = "multipart/form-data; boundary=$boundary"
            Body = $bodyBytes
            UseBasicParsing = $true
        }
        
        # 如果需要跳过证书验证
        if ($SkipCertificateCheck) {
            $requestParams['SkipCertificateCheck'] = $true
        }
        
        $response = Invoke-WebRequest @requestParams
        
        if ($response.StatusCode -ge 200 -and $response.StatusCode -lt 300) {
            Write-Host "✓ Successfully uploaded $FileName" -ForegroundColor Green
            return $true
        } else {
            Write-Host "✗ Failed to upload $FileName (HTTP $($response.StatusCode))" -ForegroundColor Red
            Write-Host "Response: $($response.Content)" -ForegroundColor Red
            return $false
        }
    }
    catch {
        Write-Host "✗ Failed to upload $FileName" -ForegroundColor Red
        Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

# 函数：通过文件系统复制
function Copy-ViaFilesystem {
    param(
        [string]$FilePath
    )
    
    $FileName = Split-Path -Leaf $FilePath
    $TargetDir = Join-Path $DataDir "providers\$Namespace\$Type\$Version"
    
    Write-Host "Copying $FileName to filesystem..." -ForegroundColor Yellow
    
    try {
        # 创建目标目录
        if (-not (Test-Path -Path $TargetDir)) {
            New-Item -Path $TargetDir -ItemType Directory -Force | Out-Null
        }
        
        # 复制文件
        Copy-Item -Path $FilePath -Destination $TargetDir -Force
        Write-Host "✓ Successfully copied $FileName" -ForegroundColor Green
        return $true
    }
    catch {
        Write-Host "✗ Failed to copy $FileName" -ForegroundColor Red
        Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

# 函数：上传文件
function Upload-File {
    param(
        [string]$FilePath
    )
    
    if (-not (Test-Path -Path $FilePath -PathType Leaf)) {
        Write-Host "Warning: File not found: $FilePath" -ForegroundColor Yellow
        return $false
    }
    
    if ($Method -eq "api") {
        return Upload-ViaApi -FilePath $FilePath
    } else {
        return Copy-ViaFilesystem -FilePath $FilePath
    }
}

# 主上传逻辑
$SuccessCount = 0
$FailCount = 0

# 确保 DistDir 是绝对路径
if (-not [System.IO.Path]::IsPathRooted($DistDir)) {
    $DistDir = Join-Path $PWD $DistDir
}

# 上传所有 zip 文件
Write-Host "`nUploading provider packages..." -ForegroundColor Green
$zipFiles = Get-ChildItem -Path $DistDir -Filter "*.zip"
foreach ($zipFile in $zipFiles) {
    if (Upload-File -FilePath $zipFile.FullName) {
        $SuccessCount++
    } else {
        $FailCount++
    }
}

# 上传 SHA256SUMS 文件
Write-Host "`nUploading checksum files..." -ForegroundColor Green
$shasumsFile = Join-Path $DistDir "terraform-provider-${Type}_${Version}_SHA256SUMS"
if (Test-Path -Path $shasumsFile) {
    if (Upload-File -FilePath $shasumsFile) {
        $SuccessCount++
    } else {
        $FailCount++
    }
} else {
    Write-Host "Warning: SHA256SUMS file not found" -ForegroundColor Yellow
}

# 上传签名文件（如果存在）
$sigFile = Join-Path $DistDir "terraform-provider-${Type}_${Version}_SHA256SUMS.sig"
if (Test-Path -Path $sigFile) {
    Write-Host "`nUploading signature file..." -ForegroundColor Green
    if (Upload-File -FilePath $sigFile) {
        $SuccessCount++
    } else {
        $FailCount++
    }
} else {
    Write-Host "Note: No signature file found (optional)" -ForegroundColor Yellow
}

# 上传 manifest 文件（如果存在）
$manifestFile = Join-Path $DistDir "terraform-provider-${Type}_${Version}_manifest.json"
if (-not (Test-Path -Path $manifestFile)) {
    $manifestFile = "terraform-registry-manifest.json"
}

if (Test-Path -Path $manifestFile) {
    Write-Host "`nUploading manifest file..." -ForegroundColor Green
    if (Upload-File -FilePath $manifestFile) {
        $SuccessCount++
    } else {
        $FailCount++
    }
}

# 总结
Write-Host ""
Write-Host "=== Upload Summary ===" -ForegroundColor Green
Write-Host "Successful: $SuccessCount" -ForegroundColor Green
if ($FailCount -gt 0) {
    Write-Host "Failed: $FailCount" -ForegroundColor Red
}

# 验证上传
if ($Method -eq "api") {
    Write-Host ""
    Write-Host "=== Verification ===" -ForegroundColor Green
    $versionsUrl = "${Protocol}://${ServerHost}:${Port}/v1/providers/${Namespace}/${Type}/versions"
    Write-Host "Checking provider versions at: $versionsUrl"
    
    try {
        $verifyParams = @{
            Uri = $versionsUrl
            UseBasicParsing = $true
        }
        
        if ($SkipCertificateCheck) {
            $verifyParams['SkipCertificateCheck'] = $true
        }
        
        $response = Invoke-WebRequest @verifyParams
        Write-Host ""
        Write-Host $response.Content
    }
    catch {
        Write-Host "Could not verify (service may not be available yet)" -ForegroundColor Yellow
    }
    
    Write-Host ""
    Write-Host "Test with Terraform:" -ForegroundColor Green
    Write-Host "  terraform init"
}

# 退出码
if ($FailCount -gt 0) {
    exit 1
} else {
    Write-Host "`n✓ All uploads completed successfully!" -ForegroundColor Green
    exit 0
}
