.PHONY: all build clean test run help install

APP_NAME := submit-sitemap
BUILD_DIR := dist
CMD_PATH := ./cmd/submit-sitemap
VERSION := 1.0.0

# 默认目标
all: build

# 构建 Apple Silicon 版本
build:
	@echo "🚀 开始构建 Apple Silicon 版本..."
	@mkdir -p $(BUILD_DIR)
	@echo "📦 构建 macOS Apple Silicon (arm64)..."
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 $(CMD_PATH)
	@echo "✓ 完成: $(BUILD_DIR)/$(APP_NAME)-darwin-arm64"
	@echo "📦 构建通用版本..."
	@go build -o $(APP_NAME) $(CMD_PATH)
	@echo "✓ 完成: $(APP_NAME)"
	@chmod +x $(BUILD_DIR)/* $(APP_NAME)
	@echo "✨ 构建完成！"

# 快速构建（仅当前平台）
quick:
	@echo "⚡ 快速构建..."
	@go build -o $(APP_NAME) $(CMD_PATH)
	@chmod +x $(APP_NAME)
	@echo "✨ 构建完成: $(APP_NAME)"

# 清理构建文件
clean:
	@echo "🧹 清理构建文件..."
	@rm -f $(APP_NAME)
	@rm -rf $(BUILD_DIR)
	@echo "✓ 清理完成"

# 测试
test:
	@echo "🧪 运行测试..."
	@go test -v ./...

# 运行程序
run:
	@go run $(CMD_PATH) run

# 安装依赖
deps:
	@echo "📦 安装依赖..."
	@go get gopkg.in/yaml.v3
	@go mod tidy
	@echo "✓ 依赖安装完成"

# 显示帮助
help:
	@echo "Submit Sitemap - Makefile 使用说明"
	@echo ""
	@echo "可用命令:"
	@echo "  make build   - 构建所有 macOS 版本（默认）"
	@echo "  make quick   - 快速构建当前平台版本"
	@echo "  make clean   - 清理构建文件"
	@echo "  make test    - 运行测试"
	@echo "  make run     - 直接运行程序"
	@echo "  make deps    - 安装依赖"
	@echo "  make help    - 显示此帮助信息"
	@echo ""
	@echo "示例:"
	@echo "  make         - 构建所有版本"
	@echo "  make quick   - 快速构建"
	@echo "  make clean   - 清理后重新构建"
