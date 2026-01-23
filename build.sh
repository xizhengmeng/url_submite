#!/bin/bash

# Submit Sitemap 构建脚本
# 用于构建 macOS 可执行文件

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目信息
APP_NAME="submit-sitemap"
VERSION="1.0.0"
BUILD_DIR="build"
CMD_PATH="./cmd/submit-sitemap"

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}  Submit Sitemap 构建工具${NC}"
echo -e "${BLUE}================================${NC}"
echo ""

# 清理旧的构建文件
echo -e "${YELLOW}🧹 清理旧的构建文件...${NC}"
rm -f ${APP_NAME}
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}
echo -e "${GREEN}✓ 清理完成${NC}"
echo ""

# 获取当前系统架构
CURRENT_ARCH=$(uname -m)
if [ "$CURRENT_ARCH" = "x86_64" ]; then
    NATIVE_ARCH="amd64"
elif [ "$CURRENT_ARCH" = "arm64" ]; then
    NATIVE_ARCH="arm64"
else
    NATIVE_ARCH="amd64"
fi

echo -e "${YELLOW}📦 开始构建 macOS 版本...${NC}"
echo ""

# 构建 macOS Intel (amd64)
echo -e "${BLUE}[1/3] 构建 macOS Intel (amd64)...${NC}"
GOOS=darwin GOARCH=amd64 go build -o ${BUILD_DIR}/${APP_NAME}-darwin-amd64 ${CMD_PATH}
if [ $? -eq 0 ]; then
    SIZE=$(du -h ${BUILD_DIR}/${APP_NAME}-darwin-amd64 | cut -f1)
    echo -e "${GREEN}✓ 构建成功: ${BUILD_DIR}/${APP_NAME}-darwin-amd64 (${SIZE})${NC}"
else
    echo -e "${RED}✗ 构建失败${NC}"
    exit 1
fi
echo ""

# 构建 macOS Apple Silicon (arm64)
echo -e "${BLUE}[2/3] 构建 macOS Apple Silicon (arm64)...${NC}"
GOOS=darwin GOARCH=arm64 go build -o ${BUILD_DIR}/${APP_NAME}-darwin-arm64 ${CMD_PATH}
if [ $? -eq 0 ]; then
    SIZE=$(du -h ${BUILD_DIR}/${APP_NAME}-darwin-arm64 | cut -f1)
    echo -e "${GREEN}✓ 构建成功: ${BUILD_DIR}/${APP_NAME}-darwin-arm64 (${SIZE})${NC}"
else
    echo -e "${RED}✗ 构建失败${NC}"
    exit 1
fi
echo ""

# 构建通用版本（当前平台）
echo -e "${BLUE}[3/3] 构建当前平台通用版本 (${NATIVE_ARCH})...${NC}"
go build -o ${APP_NAME} ${CMD_PATH}
if [ $? -eq 0 ]; then
    SIZE=$(du -h ${APP_NAME} | cut -f1)
    echo -e "${GREEN}✓ 构建成功: ${APP_NAME} (${SIZE})${NC}"

    # 添加执行权限
    chmod +x ${APP_NAME}
else
    echo -e "${RED}✗ 构建失败${NC}"
    exit 1
fi
echo ""

# 添加执行权限到所有构建文件
chmod +x ${BUILD_DIR}/*

# 显示构建结果
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}✨ 构建完成！${NC}"
echo -e "${GREEN}================================${NC}"
echo ""
echo -e "${YELLOW}📂 构建文件列表：${NC}"
echo ""
ls -lh ${BUILD_DIR}/ | tail -n +2 | awk '{printf "  • %-30s %s\n", $9, $5}'
echo ""
echo -e "  • ${YELLOW}${APP_NAME}${NC} (当前平台通用版本)"
echo ""
echo -e "${BLUE}💡 使用说明：${NC}"
echo ""
echo -e "  当前平台直接运行:"
echo -e "    ${GREEN}./${APP_NAME} help${NC}"
echo ""
echo -e "  或者使用指定架构版本:"
if [ "$NATIVE_ARCH" = "arm64" ]; then
    echo -e "    ${GREEN}./${BUILD_DIR}/${APP_NAME}-darwin-arm64 help${NC} (Apple Silicon)"
    echo -e "    ${GREEN}./${BUILD_DIR}/${APP_NAME}-darwin-amd64 help${NC} (Intel Mac)"
else
    echo -e "    ${GREEN}./${BUILD_DIR}/${APP_NAME}-darwin-amd64 help${NC} (Intel Mac)"
    echo -e "    ${GREEN}./${BUILD_DIR}/${APP_NAME}-darwin-arm64 help${NC} (Apple Silicon)"
fi
echo ""
echo -e "${BLUE}📦 版本信息：${NC}"
echo -e "  版本: ${VERSION}"
echo -e "  构建时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo -e "  Go版本: $(go version | awk '{print $3}')"
echo ""
