.PHONY: all build build-all clean test run help install deps migrate info quick

# 项目配置
APP_NAME := submit
BUILD_DIR := dist
CMD_PATH := ./cmd/submit-sitemap
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "1.2.0")
BUILD_TIME := $(shell date '+%Y-%m-%d %H:%M:%S')
GO_VERSION := $(shell go version | awk '{print $$3}')
LDFLAGS := -X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)'

# 颜色定义
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m

# 默认目标
all: build

# 构建所有版本
build: clean
	@echo "$(BLUE)================================$(NC)"
	@echo "$(BLUE)  Submit Sitemap 构建工具$(NC)"
	@echo "$(BLUE)================================$(NC)"
	@echo ""
	@echo "$(YELLOW)📦 项目信息:$(NC)"
	@echo "  版本: $(VERSION)"
	@echo "  构建时间: $(BUILD_TIME)"
	@echo "  Go版本: $(GO_VERSION)"
	@echo ""
	@mkdir -p $(BUILD_DIR)
	@echo "$(YELLOW)🔨 开始构建...$(NC)"
	@echo ""
	@echo "$(BLUE)[1/2] 构建 macOS Apple Silicon (arm64)...$(NC)"
	@GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 $(CMD_PATH)
	@SIZE=$$(du -h $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 | cut -f1); \
	echo "$(GREEN)✓ 完成: $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 ($$SIZE)$(NC)"
	@echo ""
	@echo "$(BLUE)[2/2] 构建通用版本...$(NC)"
	@go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) $(CMD_PATH)
	@SIZE=$$(du -h $(BUILD_DIR)/$(APP_NAME) | cut -f1); \
	echo "$(GREEN)✓ 完成: $(BUILD_DIR)/$(APP_NAME) ($$SIZE)$(NC)"
	@chmod +x $(BUILD_DIR)/*
	@echo ""
	@echo "$(GREEN)================================$(NC)"
	@echo "$(GREEN)✨ 构建完成！$(NC)"
	@echo "$(GREEN)================================$(NC)"
	@echo ""
	@echo "$(YELLOW)📂 构建文件:$(NC)"
	@ls -lh $(BUILD_DIR)/ | tail -n +2 | awk '{printf "  • %s (%s)\n", $$9, $$5}'
	@echo ""
	@echo "$(BLUE)💡 快速开始:$(NC)"
	@echo "  $(GREEN)./$(BUILD_DIR)/$(APP_NAME) test$(NC)     # 测试配置"
	@echo "  $(GREEN)./$(BUILD_DIR)/$(APP_NAME) run$(NC)      # 运行提交"
	@echo "  $(GREEN)./$(BUILD_DIR)/$(APP_NAME) help$(NC)     # 查看帮助"
	@echo ""

# 构建所有平台版本（包括 Linux、Windows）
build-all: clean
	@echo "$(YELLOW)🌍 构建所有平台版本...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@echo "$(BLUE)• macOS Apple Silicon (arm64)...$(NC)"
	@GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 $(CMD_PATH)
	@echo "$(GREEN)✓$(NC)"
	@echo "$(BLUE)• macOS Intel (amd64)...$(NC)"
	@GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(CMD_PATH)
	@echo "$(GREEN)✓$(NC)"
	@echo "$(BLUE)• Linux (amd64)...$(NC)"
	@GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(CMD_PATH)
	@echo "$(GREEN)✓$(NC)"
	@echo "$(BLUE)• Windows (amd64)...$(NC)"
	@GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe $(CMD_PATH)
	@echo "$(GREEN)✓$(NC)"
	@chmod +x $(BUILD_DIR)/*
	@echo ""
	@echo "$(GREEN)✨ 所有平台构建完成！$(NC)"
	@ls -lh $(BUILD_DIR)/

# 快速构建（仅当前平台）
quick:
	@echo "$(YELLOW)⚡ 快速构建...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) $(CMD_PATH)
	@chmod +x $(BUILD_DIR)/$(APP_NAME)
	@SIZE=$$(du -h $(BUILD_DIR)/$(APP_NAME) | cut -f1); \
	echo "$(GREEN)✨ 构建完成: $(BUILD_DIR)/$(APP_NAME) ($$SIZE)$(NC)"

# 清理构建文件
clean:
	@echo "$(YELLOW)🧹 清理构建文件...$(NC)"
	@rm -rf $(BUILD_DIR)
	@echo "$(GREEN)✓ 清理完成$(NC)"

# 测试
test:
	@echo "$(YELLOW)🧪 运行测试...$(NC)"
	@go test -v ./...

# 运行程序
run:
	@go run $(CMD_PATH) run

# 安装到系统路径
install: build
	@echo "$(YELLOW)📦 安装到系统...$(NC)"
	@sudo cp $(BUILD_DIR)/$(APP_NAME) /usr/local/bin/
	@sudo chmod +x /usr/local/bin/$(APP_NAME)
	@echo "$(GREEN)✓ 安装完成$(NC)"
	@echo ""
	@echo "$(BLUE)现在可以在任何位置运行:$(NC)"
	@echo "  $(GREEN)$(APP_NAME) test$(NC)"
	@echo "  $(GREEN)$(APP_NAME) run$(NC)"

# 卸载
uninstall:
	@echo "$(YELLOW)🗑️  卸载...$(NC)"
	@sudo rm -f /usr/local/bin/$(APP_NAME)
	@echo "$(GREEN)✓ 卸载完成$(NC)"

# 安装依赖
deps:
	@echo "$(YELLOW)📦 安装依赖...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✓ 依赖安装完成$(NC)"

# 迁移配置
migrate:
	@echo "$(YELLOW)🔄 运行配置迁移...$(NC)"
	@if [ -f migrate-config.sh ]; then \
		./migrate-config.sh; \
	else \
		echo "$(RED)❌ 未找到迁移脚本$(NC)"; \
		exit 1; \
	fi

# 显示项目信息
info:
	@echo "$(BLUE)================================$(NC)"
	@echo "$(BLUE)  Submit Sitemap 项目信息$(NC)"
	@echo "$(BLUE)================================$(NC)"
	@echo ""
	@echo "$(YELLOW)项目:$(NC)"
	@echo "  名称: $(APP_NAME)"
	@echo "  版本: $(VERSION)"
	@echo "  构建时间: $(BUILD_TIME)"
	@echo ""
	@echo "$(YELLOW)环境:$(NC)"
	@echo "  Go版本: $(GO_VERSION)"
	@echo "  GOOS: $(shell go env GOOS)"
	@echo "  GOARCH: $(shell go env GOARCH)"
	@echo ""
	@echo "$(YELLOW)目录:$(NC)"
	@echo "  配置: ~/.submit/config/sites/"
	@echo "  日志: ~/.submit/data/logs/"
	@echo "  历史: ~/.submit/data/submitted/"
	@echo ""
	@if [ -d $(BUILD_DIR) ]; then \
		echo "$(YELLOW)已构建文件:$(NC)"; \
		ls -lh $(BUILD_DIR)/ | tail -n +2 | awk '{printf "  • %s (%s)\n", $$9, $$5}'; \
	else \
		echo "$(YELLOW)状态:$(NC) 尚未构建"; \
	fi
	@echo ""

# 显示帮助
help:
	@echo "$(BLUE)================================$(NC)"
	@echo "$(BLUE)  Submit Sitemap - Makefile$(NC)"
	@echo "$(BLUE)================================$(NC)"
	@echo ""
	@echo "$(YELLOW)构建命令:$(NC)"
	@echo "  $(GREEN)make$(NC)              构建所有 macOS 版本（默认）"
	@echo "  $(GREEN)make build$(NC)        同上"
	@echo "  $(GREEN)make build-all$(NC)    构建所有平台版本（macOS/Linux/Windows）"
	@echo "  $(GREEN)make quick$(NC)        快速构建当前平台版本"
	@echo ""
	@echo "$(YELLOW)管理命令:$(NC)"
	@echo "  $(GREEN)make clean$(NC)        清理构建文件"
	@echo "  $(GREEN)make install$(NC)      安装到系统路径 (/usr/local/bin)"
	@echo "  $(GREEN)make uninstall$(NC)    从系统卸载"
	@echo "  $(GREEN)make deps$(NC)         安装/更新依赖"
	@echo "  $(GREEN)make migrate$(NC)      运行配置迁移脚本"
	@echo ""
	@echo "$(YELLOW)开发命令:$(NC)"
	@echo "  $(GREEN)make test$(NC)         运行测试"
	@echo "  $(GREEN)make run$(NC)          直接运行程序（开发模式）"
	@echo "  $(GREEN)make info$(NC)         显示项目信息"
	@echo "  $(GREEN)make help$(NC)         显示此帮助信息"
	@echo ""
	@echo "$(YELLOW)常用流程:$(NC)"
	@echo "  $(BLUE)# 首次构建$(NC)"
	@echo "  make deps && make build"
	@echo ""
	@echo "  $(BLUE)# 清理重建$(NC)"
	@echo "  make clean && make"
	@echo ""
	@echo "  $(BLUE)# 安装到系统$(NC)"
	@echo "  make install"
	@echo ""
	@echo "  $(BLUE)# 快速测试$(NC)"
	@echo "  make quick && ./dist/submit test"
	@echo ""
