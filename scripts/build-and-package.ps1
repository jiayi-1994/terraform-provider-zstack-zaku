<#
.SYNOPSIS
    构建并打包 Terraform Provider (支持 linux_amd64)
    
.DESCRIPTION
    此脚本用于交叉编译 Terraform Provider 到 linux_amd64 平台，并生成 zip 和 SHA256SUMS 文件
    
.PARAMETER Version
    Provider 版本号（例如：1.0.0）
    
.PARAMETER OutputDir
    输出目录路径（默认：./dist）
    
.PARAMETER ProviderName
    Provider 名称（默认：zstack-zaku）
    
.EXAMPLE
    .\build-and-package.ps1 -Version "1.0.0"
    
.EXAMPLE
    .\build-and-package.ps1 -Version "1.0.0" -OutputDir ".\output" -ProviderName "zstack-zaku"
#>

param(
    [Parameter(Mandatory=$true)]
    [string]$Version,
    
    [Parameter(Mandatory=$false)]
    [string]$OutputDir = "dist",
    
    [Parameter(Mandatory=$false)]
    [string]$ProviderName = "zstack-zaku"
)

# 错误时停止
$ErrorActionPreference = "Stop"

# 颜色输出函数
function Write-ColorOutput {
    param(
        [string]$Message,
        [string]$Color = "White"
    )
    Write-Host $Message -ForegroundColor $Color
}

function Write-Success {
    param([string]$Message)
    Write-ColorOutput "✓ $Message" "Green"
}

function Write-Info {
    param([string]$Message)
    Write-ColorOutput "ℹ $Message" "Cyan"
}

function Write-Error {
    param([string]$Message)
    Write-ColorOutput "✗ $Message" "Red"
}

function Write-Step {
    param([string]$Message)
    Write-Host ""
    Write-ColorOutput "===> $Message" "Yellow"
}

# 检查 Go 环境
function Test-GoEnvironment {
    Write-Step "检查 Go 环境"
    
    try {
        $goVersion = go version
        Write-Success "Go 环境正常: $goVersion"
        return $true
    } catch {
        Write-Error "未找到 Go 环境，请先安装 Go"
        return $false
    }
}

# 构建 Provider
function Build-Provider {
    param(
        [string]$OS,
        [string]$Arch,
        [string]$OutputPath
    )
    
    $platform = "${OS}_${Arch}"
    Write-Info "构建平台: $platform"
    
    # 设置环境变量
    $env:GOOS = $OS
    $env:GOARCH = $Arch
    $env:CGO_ENABLED = "0"
    
    # 构建输出文件名
    $binaryName = "terraform-provider-${ProviderName}_v${Version}"
    if ($OS -eq "windows") {
        $binaryName += ".exe"
    }
    
    $binaryPath = Join-Path $OutputPath $binaryName
    
    Write-Info "输出文件: $binaryPath"
    Write-Info "开始编译..."
    
    # 执行构建
    try {
        $buildArgs = @(
            "build",
            "-o", $binaryPath,
            "-ldflags", "-s -w -X main.version=$Version",
            "."
        )
        
        & go @buildArgs
        
        if ($LASTEXITCODE -ne 0) {
            throw "构建失败，退出码: $LASTEXITCODE"
        }
        
        Write-Success "编译完成: $binaryName"
        return $binaryPath
        
    } catch {
        Write-Error "编译失败: $_"
        throw
    }
}

# 打包 Provider
function New-ProviderPackage {
    param(
        [string]$BinaryPath,
        [string]$Platform,
        [string]$OutputDirectory
    )
    
    Write-Info "打包 $Platform..."
    
    # 生成文件名
    $zipFileName = "terraform-provider-${ProviderName}_${Version}_${Platform}.zip"
    $zipFilePath = Join-Path $OutputDirectory $zipFileName
    
    # 创建临时目录
    $tempDir = Join-Path $env:TEMP "terraform-provider-package-$(Get-Random)"
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
    
    try {
        # 复制二进制文件
        $tempBinaryPath = Join-Path $tempDir ([System.IO.Path]::GetFileName($BinaryPath))
        Copy-Item -Path $BinaryPath -Destination $tempBinaryPath -Force
        
        # 删除旧的 zip 文件
        if (Test-Path $zipFilePath) {
            Remove-Item $zipFilePath -Force
        }
        
        # 创建 zip
        Add-Type -AssemblyName System.IO.Compression.FileSystem
        [System.IO.Compression.ZipFile]::CreateFromDirectory($tempDir, $zipFilePath)
        
        Write-Success "已创建: $zipFileName"
        
        # 删除二进制文件（只保留 zip）
        Remove-Item $BinaryPath -Force
        
        return $zipFilePath
        
    } finally {
        # 清理临时目录
        if (Test-Path $tempDir) {
            Remove-Item -Path $tempDir -Recurse -Force
        }
    }
}

# 生成 SHA256SUMS
function New-SHA256Sums {
    param(
        [string[]]$ZipFiles,
        [string]$OutputDirectory
    )
    
    Write-Step "生成 SHA256SUMS 文件"
    
    $sha256FileName = "terraform-provider-${ProviderName}_${Version}_SHA256SUMS"
    $sha256FilePath = Join-Path $OutputDirectory $sha256FileName
    
    $sha256Lines = @()
    
    foreach ($zipFile in $ZipFiles) {
        $fileName = [System.IO.Path]::GetFileName($zipFile)
        $hash = Get-FileHash -Path $zipFile -Algorithm SHA256
        $line = "$($hash.Hash.ToLower())  $fileName"
        $sha256Lines += $line
        Write-Info "  $fileName"
    }
    
    # 写入文件
    $sha256Lines | Set-Content -Path $sha256FilePath -Encoding ASCII
    
    Write-Success "已创建: $sha256FileName"
    
    return $sha256FilePath
}

# 显示构建摘要
function Show-BuildSummary {
    param(
        [string]$OutputDirectory,
        [string[]]$Files
    )
    
    Write-Step "构建摘要"
    
    Write-Host ""
    Write-Info "输出目录: $OutputDirectory"
    Write-Host ""
    Write-Info "生成的文件:"
    
    $totalSize = 0
    foreach ($file in $Files) {
        $item = Get-Item $file
        $size = [math]::Round($item.Length / 1MB, 2)
        $totalSize += $size
        Write-Host "  $($item.Name) ($size MB)" -ForegroundColor Yellow
    }
    
    Write-Host ""
    Write-Info "总大小: $totalSize MB"
    Write-Host ""
}

# 主函数
function Main {
    Write-Host ""
    Write-ColorOutput "=========================================" "Magenta"
    Write-ColorOutput "  Terraform Provider 构建打包工具" "Magenta"
    Write-ColorOutput "=========================================" "Magenta"
    Write-Host ""
    Write-Info "Provider: $ProviderName"
    Write-Info "Version:  $Version"
    Write-Info "Output:   $OutputDir"
    Write-Host ""
    
    # 检查 Go 环境
    if (-not (Test-GoEnvironment)) {
        exit 1
    }
    
    # 创建输出目录
    if (-not (Test-Path $OutputDir)) {
        New-Item -ItemType Directory -Path $OutputDir -Force | Out-Null
        Write-Success "已创建输出目录: $OutputDir"
    }
    
    # 获取绝对路径
    $outputDirAbsolute = (Resolve-Path $OutputDir).Path
    
    # 定义构建目标
    $targets = @(
        @{OS="linux"; Arch="amd64"}
    )
    
    $zipFiles = @()
    
    # 构建所有目标
    foreach ($target in $targets) {
        Write-Step "构建 $($target.OS)_$($target.Arch)"
        
        try {
            # 构建
            $binaryPath = Build-Provider -OS $target.OS -Arch $target.Arch -OutputPath $outputDirAbsolute
            
            # 打包
            $platform = "$($target.OS)_$($target.Arch)"
            $zipPath = New-ProviderPackage -BinaryPath $binaryPath -Platform $platform -OutputDirectory $outputDirAbsolute
            $zipFiles += $zipPath
            
        } catch {
            Write-Error "构建 $($target.OS)_$($target.Arch) 失败: $_"
            throw
        }
    }
    
    # 生成 SHA256SUMS
    $sha256File = New-SHA256Sums -ZipFiles $zipFiles -OutputDirectory $outputDirAbsolute
    
    # 显示摘要
    $allFiles = $zipFiles + $sha256File
    Show-BuildSummary -OutputDirectory $outputDirAbsolute -Files $allFiles
    
    Write-Success "✓ 所有任务完成!"
    Write-Host ""
}

# 执行
try {
    Main
} catch {
    Write-Host ""
    Write-Error "构建失败: $_"
    Write-Host $_.ScriptStackTrace -ForegroundColor Red
    exit 1
}
