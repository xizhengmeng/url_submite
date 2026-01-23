#!/bin/bash

# Submit Sitemap 配置迁移脚本
# 将项目目录下的配置迁移到 ~/.submit

set -e

echo "🔄 Submit Sitemap 配置迁移工具"
echo "================================"
echo ""

# 检查源目录
SOURCE_CONFIG="config/sites"
SOURCE_DATA="data"
TARGET_DIR="$HOME/.submit"

if [ ! -d "$SOURCE_CONFIG" ]; then
    echo "❌ 未找到源配置目录: $SOURCE_CONFIG"
    echo "请在项目根目录下运行此脚本"
    exit 1
fi

# 显示迁移计划
echo "📋 迁移计划:"
echo "  源配置: $(pwd)/$SOURCE_CONFIG"
echo "  源数据: $(pwd)/$SOURCE_DATA"
echo "  目标:   $TARGET_DIR"
echo ""

# 确认
read -p "是否继续? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "取消迁移"
    exit 0
fi

# 创建目标目录
echo ""
echo "📂 创建目标目录..."
mkdir -p "$TARGET_DIR/config/sites"
mkdir -p "$TARGET_DIR/data/logs"
mkdir -p "$TARGET_DIR/data/submitted"

# 迁移配置文件
echo "📋 迁移配置文件..."
if [ -d "$SOURCE_CONFIG" ]; then
    # 复制所有站点配置（排除示例文件）
    find "$SOURCE_CONFIG" -type f -name "*.yaml" ! -name "*.example" | while read -r file; do
        # 获取相对路径
        rel_path="${file#$SOURCE_CONFIG/}"
        target_file="$TARGET_DIR/config/sites/$rel_path"
        target_dir=$(dirname "$target_file")

        # 创建目标目录
        mkdir -p "$target_dir"

        # 复制文件
        cp "$file" "$target_file"
        echo "  ✅ $rel_path"
    done
fi

# 迁移历史数据
echo ""
echo "📊 迁移历史数据..."
if [ -d "$SOURCE_DATA/submitted" ]; then
    cp -r "$SOURCE_DATA/submitted/"* "$TARGET_DIR/data/submitted/" 2>/dev/null || true
    echo "  ✅ 提交历史"
fi

# 迁移日志（可选）
echo ""
read -p "是否迁移日志文件? (y/N) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if [ -d "$SOURCE_DATA/logs" ]; then
        cp -r "$SOURCE_DATA/logs/"* "$TARGET_DIR/data/logs/" 2>/dev/null || true
        echo "  ✅ 日志文件"
    fi
fi

# 显示结果
echo ""
echo "✅ 迁移完成!"
echo ""
echo "📍 新配置位置:"
echo "  配置: $TARGET_DIR/config/sites/"
echo "  数据: $TARGET_DIR/data/"
echo ""
echo "📝 后续步骤:"
echo "  1. 验证配置: ./submit test"
echo "  2. 查看统计: ./submit stats"
echo "  3. 运行提交: ./submit run"
echo ""
echo "💡 提示:"
echo "  - 原配置文件已保留在项目目录"
echo "  - 可以安全删除项目目录下的 config/ 和 data/"
echo "  - 或保留作为备份"
echo ""
