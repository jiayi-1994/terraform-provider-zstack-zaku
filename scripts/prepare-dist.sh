#!/bin/bash
# prepare-dist.sh
# 准备 dist 目录，确保所有必需的文件都存在

set -e

# 默认配置
VERSION="${PROVIDER_VERSION:-1.0.0}"
PROVIDER_NAME="${PROVIDER_NAME:-terraform-provider-zstack-zaku}"
DIST_DIR="${DIST_DIR:-dist}"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
GRAY='\033[0;37m'
NC='\033[0m' # No Color

# 帮助信息
show_help() {
    cat << EOF
Usage: $0 [OPTIONS]

准备 dist 目录，确保所有 Terraform Registry 必需的文件都存在

OPTIONS:
    -h, --help              显示帮助信息
    -v, --version VERSION   Provider 版本 (默认: 1.0.0)
    -n, --name NAME         Provider 名称 (默认: terraform-provider-zstack-zaku)
    -d, --dist-dir DIR      输出目录 (默认: dist)

ENVIRONMENT VARIABLES:
    PROVIDER_VERSION        Provider 版本
    PROVIDER_NAME           Provider 名称
    DIST_DIR                输出目录

EXAMPLES:
    # 使用默认配置
    $0

    # 指定版本
    $0 -v 1.0.1

    # 指定所有参数
    $0 -v 1.0.0 -n terraform-provider-zstack-zaku -d dist

EOF
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -n|--name)
            PROVIDER_NAME="$2"
            shift 2
            ;;
        -d|--dist-dir)
            DIST_DIR="$2"
            shift 2
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

echo -e "${GREEN}=== Preparing Distribution Files ===${NC}"
echo "Version:       $VERSION"
echo "Provider Name: $PROVIDER_NAME"
echo "Dist Dir:      $DIST_DIR"
echo ""

# 检查 dist 目录是否存在
if [ ! -d "$DIST_DIR" ]; then
    echo -e "${RED}Error: Dist directory '$DIST_DIR' not found${NC}"
    echo -e "${YELLOW}Please build the provider first with: go build or goreleaser${NC}"
    exit 1
fi

# 1. 生成/复制 terraform-registry-manifest.json
echo -e "${YELLOW}Processing terraform-registry-manifest.json...${NC}"

MANIFEST_SOURCE="terraform-registry-manifest.json"
MANIFEST_DEST="$DIST_DIR/${PROVIDER_NAME}_${VERSION}_manifest.json"

# 如果根目录有 manifest 文件，使用它；否则生成新的
if [ -f "$MANIFEST_SOURCE" ]; then
    echo -e "${GREEN}  ✓ Found $MANIFEST_SOURCE in root directory${NC}"
    cp "$MANIFEST_SOURCE" "$MANIFEST_DEST"
    echo -e "${GREEN}  ✓ Copied to $MANIFEST_DEST${NC}"
else
    echo -e "${YELLOW}  ! $MANIFEST_SOURCE not found, generating new one...${NC}"
    cat > "$MANIFEST_DEST" << 'EOF'
{
    "version": 1,
    "metadata": {
        "protocol_versions": ["6.0"]
    }
}
EOF
    echo -e "${GREEN}  ✓ Generated $MANIFEST_DEST${NC}"
fi

# 2. 检查必需的文件
echo -e "\n${YELLOW}Checking required files...${NC}"

ALL_FILES_EXIST=true

# 检查 ZIP 文件
ZIP_COUNT=$(find "$DIST_DIR" -maxdepth 1 -name "*.zip" -type f | wc -l)
if [ "$ZIP_COUNT" -gt 0 ]; then
    echo -e "${GREEN}  ✓ Found $ZIP_COUNT ZIP package(s)${NC}"
    find "$DIST_DIR" -maxdepth 1 -name "*.zip" -type f -exec basename {} \; | while read -r file; do
        echo -e "${GRAY}    - $file${NC}"
    done
else
    echo -e "${RED}  ✗ No ZIP packages found${NC}"
    ALL_FILES_EXIST=false
fi

# 检查 SHA256SUMS 文件
SHASUMS_COUNT=$(find "$DIST_DIR" -maxdepth 1 -name "*_SHA256SUMS" -type f | wc -l)
if [ "$SHASUMS_COUNT" -gt 0 ]; then
    echo -e "${GREEN}  ✓ Found $SHASUMS_COUNT SHA256SUMS file(s)${NC}"
    find "$DIST_DIR" -maxdepth 1 -name "*_SHA256SUMS" -type f -exec basename {} \; | while read -r file; do
        echo -e "${GRAY}    - $file${NC}"
    done
else
    echo -e "${RED}  ✗ No SHA256SUMS file found${NC}"
    ALL_FILES_EXIST=false
fi

# 3. 检查可选的文件
echo -e "\n${YELLOW}Checking optional files...${NC}"

SIG_COUNT=$(find "$DIST_DIR" -maxdepth 1 -name "*_SHA256SUMS.sig" -type f | wc -l)
if [ "$SIG_COUNT" -gt 0 ]; then
    echo -e "${GREEN}  ✓ Found $SIG_COUNT Signature file(s)${NC}"
    find "$DIST_DIR" -maxdepth 1 -name "*_SHA256SUMS.sig" -type f -exec basename {} \; | while read -r file; do
        echo -e "${GRAY}    - $file${NC}"
    done
else
    echo -e "${GRAY}  ⊘ No Signature file found (optional)${NC}"
fi

# 4. 验证 SHA256SUMS 内容
echo -e "\n${YELLOW}Validating SHA256SUMS...${NC}"

find "$DIST_DIR" -maxdepth 1 -name "*_SHA256SUMS" -type f | while read -r shasums_file; do
    LINE_COUNT=$(wc -l < "$shasums_file")
    FILENAME=$(basename "$shasums_file")
    
    if [ "$LINE_COUNT" -gt 0 ]; then
        echo -e "${GREEN}  ✓ $FILENAME contains $LINE_COUNT hash(es)${NC}"
    else
        echo -e "${RED}  ✗ $FILENAME is empty!${NC}"
        ALL_FILES_EXIST=false
    fi
done

# 5. 生成文件清单
echo -e "\n${YELLOW}Generating file manifest...${NC}"

MANIFEST_LIST_FILE="$DIST_DIR/FILES_MANIFEST.txt"

cat > "$MANIFEST_LIST_FILE" << EOF
Terraform Provider Distribution Files
======================================
Generated: $(date '+%Y-%m-%d %H:%M:%S')
Provider: $PROVIDER_NAME
Version: $VERSION

Files:
EOF

# 列出所有文件
find "$DIST_DIR" -maxdepth 1 -type f -exec ls -lh {} \; | awk '{print $9, "(" $5 ")"}' | sort | while read -r line; do
    FILENAME=$(basename $(echo "$line" | awk '{print $1}'))
    SIZE=$(echo "$line" | awk '{print $2}')
    echo "  - $FILENAME $SIZE" >> "$MANIFEST_LIST_FILE"
done

echo -e "${GREEN}  ✓ File manifest saved to $MANIFEST_LIST_FILE${NC}"

# 6. 总结
echo -e "\n${GREEN}=== Summary ===${NC}"

FILE_COUNT=$(find "$DIST_DIR" -maxdepth 1 -type f | wc -l)
TOTAL_SIZE=$(du -sh "$DIST_DIR" | awk '{print $1}')

echo "Total files: $FILE_COUNT"
echo "Total size: $TOTAL_SIZE"

if [ "$ALL_FILES_EXIST" = true ]; then
    echo -e "\n${GREEN}✓ All required files are ready for upload!${NC}"
    exit 0
else
    echo -e "\n${RED}✗ Some required files are missing. Please check the errors above.${NC}"
    exit 1
fi
