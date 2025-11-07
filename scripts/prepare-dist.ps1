# prepare-dist.ps1
# 准备 dist 目录，确保所有必需的文件都存在

param(
    [string]$Version = "1.0.0",
    [string]$ProviderName = "terraform-provider-zstack-zaku",
    [string]$DistDir = "dist",
    [switch]$Help
)

function Show-Help {
    Write-Host @"
Usage: .\prepare-dist.ps1 [OPTIONS]

准备 dist 目录，确保所有 Terraform Registry 必需的文件都存在

OPTIONS:
    -Version <VERSION>        Provider 版本 (默认: 1.0.0)
    -ProviderName <NAME>      Provider 名称 (默认: terraform-provider-zstack-zaku)
    -DistDir <DIR>            输出目录 (默认: dist)
    -Help                     显示帮助信息

EXAMPLES:
    # 使用默认配置
    .\prepare-dist.ps1

    # 指定版本
    .\prepare-dist.ps1 -Version 1.0.1

    # 指定所有参数
    .\prepare-dist.ps1 -Version 1.0.0 -ProviderName terraform-provider-zstack-zaku -DistDir dist

"@
}

if ($Help) {
    Show-Help
    exit 0
}

Write-Host "=== Preparing Distribution Files ===" -ForegroundColor Green
Write-Host "Version:       $Version"
Write-Host "Provider Name: $ProviderName"
Write-Host "Dist Dir:      $DistDir"
Write-Host ""

# 检查 dist 目录是否存在
if (-not (Test-Path -Path $DistDir)) {
    Write-Host "Error: Dist directory '$DistDir' not found" -ForegroundColor Red
    Write-Host "Please build the provider first with: go build or goreleaser" -ForegroundColor Yellow
    exit 1
}

# 1. 生成/复制 terraform-registry-manifest.json
Write-Host "Processing terraform-registry-manifest.json..." -ForegroundColor Yellow

$manifestContent = @{
    version = 1
    metadata = @{
        protocol_versions = @("6.0")
    }
} | ConvertTo-Json -Depth 10

$manifestSourceFile = "terraform-registry-manifest.json"
$manifestDestFile = Join-Path $DistDir "${ProviderName}_${Version}_manifest.json"

# 如果根目录有 manifest 文件，使用它；否则生成新的
if (Test-Path -Path $manifestSourceFile) {
    Write-Host "  ✓ Found $manifestSourceFile in root directory" -ForegroundColor Green
    Copy-Item -Path $manifestSourceFile -Destination $manifestDestFile -Force
    Write-Host "  ✓ Copied to $manifestDestFile" -ForegroundColor Green
} else {
    Write-Host "  ! $manifestSourceFile not found, generating new one..." -ForegroundColor Yellow
    $manifestContent | Out-File -FilePath $manifestDestFile -Encoding utf8 -NoNewline
    Write-Host "  ✓ Generated $manifestDestFile" -ForegroundColor Green
}

# 2. 检查必需的文件
Write-Host "`nChecking required files..." -ForegroundColor Yellow

$requiredFiles = @{
    "ZIP packages" = "*.zip"
    "SHA256SUMS" = "*_SHA256SUMS"
}

$allFilesExist = $true

foreach ($fileType in $requiredFiles.Keys) {
    $pattern = $requiredFiles[$fileType]
    $files = Get-ChildItem -Path $DistDir -Filter $pattern
    
    if ($files.Count -gt 0) {
        Write-Host "  ✓ Found $($files.Count) $fileType file(s)" -ForegroundColor Green
        foreach ($file in $files) {
            Write-Host "    - $($file.Name)" -ForegroundColor Gray
        }
    } else {
        Write-Host "  ✗ No $fileType found (pattern: $pattern)" -ForegroundColor Red
        $allFilesExist = $false
    }
}

# 3. 检查可选的文件
Write-Host "`nChecking optional files..." -ForegroundColor Yellow

$optionalFiles = @{
    "Signature" = "*_SHA256SUMS.sig"
}

foreach ($fileType in $optionalFiles.Keys) {
    $pattern = $optionalFiles[$fileType]
    $files = Get-ChildItem -Path $DistDir -Filter $pattern
    
    if ($files.Count -gt 0) {
        Write-Host "  ✓ Found $($files.Count) $fileType file(s)" -ForegroundColor Green
        foreach ($file in $files) {
            Write-Host "    - $($file.Name)" -ForegroundColor Gray
        }
    } else {
        Write-Host "  ⊘ No $fileType found (optional)" -ForegroundColor Gray
    }
}

# 4. 验证 SHA256SUMS 内容
Write-Host "`nValidating SHA256SUMS..." -ForegroundColor Yellow

$shasumsFiles = Get-ChildItem -Path $DistDir -Filter "*_SHA256SUMS"
foreach ($shasumsFile in $shasumsFiles) {
    $content = Get-Content -Path $shasumsFile.FullName
    $lineCount = ($content | Measure-Object -Line).Lines
    
    if ($lineCount -gt 0) {
        Write-Host "  ✓ $($shasumsFile.Name) contains $lineCount hash(es)" -ForegroundColor Green
    } else {
        Write-Host "  ✗ $($shasumsFile.Name) is empty!" -ForegroundColor Red
        $allFilesExist = $false
    }
}

# 5. 生成文件清单
Write-Host "`nGenerating file manifest..." -ForegroundColor Yellow

$manifestListFile = Join-Path $DistDir "FILES_MANIFEST.txt"
$allFiles = Get-ChildItem -Path $DistDir -File | Sort-Object Name

$manifestListContent = @"
Terraform Provider Distribution Files
======================================
Generated: $(Get-Date -Format "yyyy-MM-dd HH:mm:ss")
Provider: $ProviderName
Version: $Version

Files:
"@

foreach ($file in $allFiles) {
    $size = if ($file.Length -gt 1MB) {
        "{0:N2} MB" -f ($file.Length / 1MB)
    } elseif ($file.Length -gt 1KB) {
        "{0:N2} KB" -f ($file.Length / 1KB)
    } else {
        "$($file.Length) bytes"
    }
    
    $manifestListContent += "`n  - $($file.Name) ($size)"
}

$manifestListContent | Out-File -FilePath $manifestListFile -Encoding utf8

Write-Host "  ✓ File manifest saved to $manifestListFile" -ForegroundColor Green

# 6. 总结
Write-Host "`n=== Summary ===" -ForegroundColor Green
Write-Host "Total files: $($allFiles.Count)"
Write-Host "Total size: $("{0:N2} MB" -f (($allFiles | Measure-Object -Property Length -Sum).Sum / 1MB))"

if ($allFilesExist) {
    Write-Host "`n✓ All required files are ready for upload!" -ForegroundColor Green
    exit 0
} else {
    Write-Host "`n✗ Some required files are missing. Please check the errors above." -ForegroundColor Red
    exit 1
}
