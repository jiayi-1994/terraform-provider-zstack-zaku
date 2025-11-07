#!/bin/bash
# upload-to-hermitcrab.sh
# 自动上传 Terraform Provider 到 Hermitcrab 私有仓库

set -e

# 配置变量
HERMITCRAB_HOST="${HERMITCRAB_HOST:-localhost}"
HERMITCRAB_PORT="${HERMITCRAB_PORT:-5000}"
NAMESPACE="${PROVIDER_NAMESPACE:-zstack}"
TYPE="${PROVIDER_TYPE:-zstack-zaku}"
VERSION="${PROVIDER_VERSION:-1.0.0}"
DIST_DIR="${DIST_DIR:-dist}"
PROTOCOL="${HERMITCRAB_PROTOCOL:-}"
SKIP_CERT_CHECK="false"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 帮助信息
show_help() {
    cat << EOF
Usage: $0 [OPTIONS]

上传 Terraform Provider 到 Hermitcrab 私有仓库

OPTIONS:
    -h, --help              显示帮助信息
    -H, --host HOST         Hermitcrab 主机地址 (默认: localhost)
    -p, --port PORT         Hermitcrab 端口 (默认: 5000)
    -n, --namespace NS      Provider 命名空间 (默认: zstack)
    -t, --type TYPE         Provider 类型 (默认: zstack-zaku)
    -v, --version VERSION   Provider 版本 (默认: 1.0.0)
    -d, --dist-dir DIR      构建产物目录 (默认: dist)
    -m, --method METHOD     上传方法: api 或 filesystem (默认: api)
    --data-dir DIR          Hermitcrab 数据目录 (filesystem 模式使用)
    --protocol PROTOCOL     协议: http 或 https (默认: 根据端口自动判断)
    --https                 强制使用 HTTPS
    --skip-cert-check       跳过 SSL 证书验证 (用于自签名证书)

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
    export HERMITCRAB_HOST="registry.example.com"
    export HERMITCRAB_PORT="5000"
    $0

    # 使用命令行参数
    $0 -H registry.example.com -p 5000 -v 1.0.1

    # 使用 HTTPS
    $0 -H registry.example.com -p 443 --https

    # 使用 HTTPS 并跳过证书验证（自签名证书）
    $0 -H 172.26.50.232 -p 443 --skip-cert-check

    # 使用文件系统模式
    $0 -m filesystem --data-dir /var/lib/hermitcrab

EOF
}

# 解析命令行参数
UPLOAD_METHOD="api"
HERMITCRAB_DATA_DIR="${HERMITCRAB_DATA_DIR:-/var/lib/hermitcrab}"

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -H|--host)
            HERMITCRAB_HOST="$2"
            shift 2
            ;;
        -p|--port)
            HERMITCRAB_PORT="$2"
            shift 2
            ;;
        -n|--namespace)
            NAMESPACE="$2"
            shift 2
            ;;
        -t|--type)
            TYPE="$2"
            shift 2
            ;;
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -d|--dist-dir)
            DIST_DIR="$2"
            shift 2
            ;;
        -m|--method)
            UPLOAD_METHOD="$2"
            shift 2
            ;;
        --data-dir)
            HERMITCRAB_DATA_DIR="$2"
            shift 2
            ;;
        --protocol)
            PROTOCOL="$2"
            shift 2
            ;;
        --https)
            PROTOCOL="https"
            shift
            ;;
        --skip-cert-check)
            SKIP_CERT_CHECK="true"
            shift
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

# 确定协议
if [ -z "$PROTOCOL" ]; then
    if [ "$HERMITCRAB_PORT" = "443" ]; then
        PROTOCOL="https"  # 端口 443 默认使用 HTTPS
    else
        PROTOCOL="http"
    fi
else
    PROTOCOL=$(echo "$PROTOCOL" | tr '[:upper:]' '[:lower:]')
fi

# 打印配置
echo -e "${GREEN}=== Hermitcrab Provider Upload ===${NC}"
echo "Upload Method:    $UPLOAD_METHOD"
echo "Protocol:         $PROTOCOL"
echo "Hermitcrab Host:  $HERMITCRAB_HOST"
echo "Hermitcrab Port:  $HERMITCRAB_PORT"
echo "Namespace:        $NAMESPACE"
echo "Provider Type:    $TYPE"
echo "Version:          $VERSION"
echo "Dist Directory:   $DIST_DIR"

if [ "$UPLOAD_METHOD" = "filesystem" ]; then
    echo "Data Directory:   $HERMITCRAB_DATA_DIR"
fi

echo ""

# 检查 dist 目录是否存在
if [ ! -d "$DIST_DIR" ]; then
    echo -e "${RED}Error: Dist directory '$DIST_DIR' not found${NC}"
    exit 1
fi

# 函数：通过 API 上传文件
upload_via_api() {
    local file=$1
    local filename=$(basename "$file")
    local url="${PROTOCOL}://${HERMITCRAB_HOST}:${HERMITCRAB_PORT}/v1/providers/${NAMESPACE}/${TYPE}/${VERSION}/upload"
    
    echo -e "${YELLOW}Uploading $filename via API...${NC}"
    
    # 准备 curl 选项
    local curl_opts="-s -w \n%{http_code}"
    if [ "$SKIP_CERT_CHECK" = "true" ]; then
        curl_opts="$curl_opts -k"
    fi
    
    # 提取平台信息（如果是 zip 文件）
    if [[ $filename == *.zip ]]; then
        platform=$(echo "$filename" | sed -E 's/.*_([^_]+_[^_]+)\.zip$/\1/')
        response=$(curl $curl_opts -X POST \
            -F "file=@${file}" \
            -F "filename=${filename}" \
            -F "platform=${platform}" \
            "$url")
    else
        response=$(curl $curl_opts -X POST \
            -F "file=@${file}" \
            -F "filename=${filename}" \
            "$url")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)
    
    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        echo -e "${GREEN}✓ Successfully uploaded $filename${NC}"
        return 0
    else
        echo -e "${RED}✗ Failed to upload $filename (HTTP $http_code)${NC}"
        echo -e "${RED}Response: $body${NC}"
        return 1
    fi
}

# 函数：通过文件系统复制
copy_via_filesystem() {
    local file=$1
    local filename=$(basename "$file")
    local target_dir="${HERMITCRAB_DATA_DIR}/providers/${NAMESPACE}/${TYPE}/${VERSION}"
    
    echo -e "${YELLOW}Copying $filename to filesystem...${NC}"
    
    # 创建目标目录
    if ! mkdir -p "$target_dir"; then
        echo -e "${RED}✗ Failed to create directory $target_dir${NC}"
        return 1
    fi
    
    # 复制文件
    if cp "$file" "$target_dir/"; then
        echo -e "${GREEN}✓ Successfully copied $filename${NC}"
        return 0
    else
        echo -e "${RED}✗ Failed to copy $filename${NC}"
        return 1
    fi
}

# 函数：上传文件
upload_file() {
    local file=$1
    
    if [ ! -f "$file" ]; then
        echo -e "${YELLOW}Warning: File not found: $file${NC}"
        return 1
    fi
    
    if [ "$UPLOAD_METHOD" = "api" ]; then
        upload_via_api "$file"
    else
        copy_via_filesystem "$file"
    fi
}

# 主上传逻辑
SUCCESS_COUNT=0
FAIL_COUNT=0

# 上传所有 zip 文件
echo -e "\n${GREEN}Uploading provider packages...${NC}"
for zipfile in "$DIST_DIR"/*.zip; do
    if [ -f "$zipfile" ]; then
        if upload_file "$zipfile"; then
            ((SUCCESS_COUNT++))
        else
            ((FAIL_COUNT++))
        fi
    fi
done

# 上传 SHA256SUMS 文件
echo -e "\n${GREEN}Uploading checksum files...${NC}"
SHASUMS_FILE="$DIST_DIR/terraform-provider-${TYPE}_${VERSION}_SHA256SUMS"
if [ -f "$SHASUMS_FILE" ]; then
    if upload_file "$SHASUMS_FILE"; then
        ((SUCCESS_COUNT++))
    else
        ((FAIL_COUNT++))
    fi
else
    echo -e "${YELLOW}Warning: SHA256SUMS file not found${NC}"
fi

# 上传签名文件（如果存在）
SIG_FILE="$DIST_DIR/terraform-provider-${TYPE}_${VERSION}_SHA256SUMS.sig"
if [ -f "$SIG_FILE" ]; then
    echo -e "\n${GREEN}Uploading signature file...${NC}"
    if upload_file "$SIG_FILE"; then
        ((SUCCESS_COUNT++))
    else
        ((FAIL_COUNT++))
    fi
else
    echo -e "${YELLOW}Note: No signature file found (optional)${NC}"
fi

# 上传 manifest 文件（如果存在）
MANIFEST_FILE="$DIST_DIR/terraform-provider-${TYPE}_${VERSION}_manifest.json"
if [ ! -f "$MANIFEST_FILE" ]; then
    MANIFEST_FILE="terraform-registry-manifest.json"
fi

if [ -f "$MANIFEST_FILE" ]; then
    echo -e "\n${GREEN}Uploading manifest file...${NC}"
    if upload_file "$MANIFEST_FILE"; then
        ((SUCCESS_COUNT++))
    else
        ((FAIL_COUNT++))
    fi
fi

# 总结
echo ""
echo -e "${GREEN}=== Upload Summary ===${NC}"
echo -e "${GREEN}Successful: $SUCCESS_COUNT${NC}"
if [ $FAIL_COUNT -gt 0 ]; then
    echo -e "${RED}Failed: $FAIL_COUNT${NC}"
fi

# 验证上传
if [ "$UPLOAD_METHOD" = "api" ]; then
    echo ""
    echo -e "${GREEN}=== Verification ===${NC}"
    VERSIONS_URL="${PROTOCOL}://${HERMITCRAB_HOST}:${HERMITCRAB_PORT}/v1/providers/${NAMESPACE}/${TYPE}/versions"
    echo "Checking provider versions at: $VERSIONS_URL"
    
    if command -v curl &> /dev/null; then
        echo ""
        if [ "$SKIP_CERT_CHECK" = "true" ]; then
            curl -s -k "$VERSIONS_URL" | head -n 20
        else
            curl -s "$VERSIONS_URL" | head -n 20
        fi
    fi
    
    echo ""
    echo -e "${GREEN}Test with Terraform:${NC}"
    echo "  terraform init"
fi

# 退出码
if [ $FAIL_COUNT -gt 0 ]; then
    exit 1
else
    echo -e "\n${GREEN}✓ All uploads completed successfully!${NC}"
    exit 0
fi
