# 构建指南

## 快速开始

```bash
# 构建程序
make

# 测试运行
./dist/submit test
```

## Makefile 命令

### 构建命令

| 命令 | 说明 |
|------|------|
| `make` | 构建所有 macOS 版本（默认） |
| `make build` | 同上 |
| `make build-all` | 构建所有平台（macOS/Linux/Windows） |
| `make quick` | 快速构建当前平台 |

### 管理命令

| 命令 | 说明 |
|------|------|
| `make clean` | 清理构建文件 |
| `make install` | 安装到系统 (/usr/local/bin) |
| `make uninstall` | 从系统卸载 |
| `make deps` | 安装/更新依赖 |
| `make migrate` | 运行配置迁移脚本 |

### 开发命令

| 命令 | 说明 |
|------|------|
| `make test` | 运行测试 |
| `make run` | 直接运行（开发模式） |
| `make info` | 显示项目信息 |
| `make help` | 显示帮助信息 |

## 构建示例

### 首次构建

```bash
# 1. 安装依赖
make deps

# 2. 构建程序
make

# 3. 测试
./dist/submit test
```

### 清理重建

```bash
make clean && make
```

### 快速迭代

```bash
# 开发时使用 quick 更快
make quick

# 测试
./dist/submit test
```

### 安装到系统

```bash
# 安装到 /usr/local/bin
make install

# 现在可以在任何地方运行
submit test
submit run
```

### 卸载

```bash
make uninstall
```

## 构建输出

### 默认构建 (make)

生成以下文件：
```
dist/
├── submit                  # 通用版本
└── submit-darwin-arm64     # Apple Silicon 专用版本
```

### 完整构建 (make build-all)

生成所有平台版本：
```
dist/
├── submit-darwin-arm64        # macOS Apple Silicon
├── submit-darwin-amd64        # macOS Intel
├── submit-linux-amd64         # Linux
└── submit-windows-amd64.exe   # Windows
```

## 版本信息

构建时会自动嵌入版本信息：

```bash
# 从 git tag 获取版本
make

# 输出示例：
# 版本: v1.2.0
# 构建时间: 2026-01-23 19:25:27
# Go版本: go1.24.3
```

如果在 git 仓库中，版本号会自动从 git tag 获取。

## 常见问题

### Q: make 构建失败？

A: 检查依赖：
```bash
make deps
go mod tidy
make
```

### Q: 如何指定版本号？

A: 使用 git tag：
```bash
git tag v1.2.0
make
```

或手动指定：
```bash
VERSION=1.2.0 make
```

### Q: 如何减小二进制文件大小？

A: 添加构建标志：
```bash
go build -ldflags="-s -w" -o dist/submit ./cmd/submit-sitemap
```

- `-s`: 去除符号表
- `-w`: 去除 DWARF 调试信息

### Q: 如何交叉编译？

A: 已内置在 `make build-all` 中，或手动：
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o dist/submit-linux ./cmd/submit-sitemap

# Windows
GOOS=windows GOARCH=amd64 go build -o dist/submit.exe ./cmd/submit-sitemap
```

## 开发流程

### 标准流程

```bash
# 1. 修改代码
vim cmd/submit-sitemap/main.go

# 2. 快速构建测试
make quick

# 3. 运行测试
./dist/submit test

# 4. 完整构建
make clean && make

# 5. 提交代码
git add .
git commit -m "feat: add new feature"
```

### 发布流程

```bash
# 1. 更新版本号
git tag v1.2.0

# 2. 构建所有平台
make build-all

# 3. 测试
./dist/submit test

# 4. 推送标签
git push --tags

# 5. 创建发布包
cd dist
tar czf submit-v1.2.0-darwin-arm64.tar.gz submit-darwin-arm64
tar czf submit-v1.2.0-linux-amd64.tar.gz submit-linux-amd64
zip submit-v1.2.0-windows-amd64.zip submit-windows-amd64.exe
```

## 持续集成

### GitHub Actions 示例

```yaml
name: Build

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      - run: make deps
      - run: make build-all
      - uses: actions/upload-artifact@v3
        with:
          name: binaries
          path: dist/*
```

## 性能优化

### 并行构建

Makefile 已配置支持并行构建：
```bash
make -j4  # 使用 4 个并行任务
```

### 增量构建

使用 `make quick` 进行增量构建，只编译修改的文件。

### 缓存

Go 会自动缓存编译结果，加快后续构建速度。

## 故障排除

### 构建慢

```bash
# 清理 Go 缓存
go clean -cache
go clean -modcache

# 重新构建
make deps
make
```

### 权限问题

```bash
# 添加执行权限
chmod +x dist/submit

# 或重新构建
make clean && make
```

### 依赖问题

```bash
# 重新安装依赖
rm go.sum
make deps
```

## 更多信息

- 运行 `make help` 查看所有可用命令
- 运行 `make info` 查看项目信息
- 查看 [README.md](../README.md) 了解使用说明
