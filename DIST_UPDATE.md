# 构建目录更新说明

## 📦 变更内容

### 1. 构建目录从 `build/` 更改为 `dist/`

所有构建产物现在输出到 `dist/` 目录：

```
dist/
└── submit-sitemap-darwin-arm64    # Apple Silicon 版本
```

### 2. 只构建 Apple Silicon 版本

移除了 Intel (amd64) 版本的构建，只保留 Apple Silicon (arm64) 版本：

**之前：**
- dist/submit-sitemap-darwin-amd64 (Intel Mac)
- dist/submit-sitemap-darwin-arm64 (Apple Silicon)

**现在：**
- dist/submit-sitemap-darwin-arm64 (Apple Silicon)

### 3. 更新的文件

✅ **build.sh**
- BUILD_DIR 改为 "dist"
- 移除 Intel 版本构建
- 更新使用说明

✅ **Makefile**
- BUILD_DIR 改为 "dist"
- 移除 Intel 版本构建

✅ **.gitignore**
- build/ 改为 dist/

✅ **README.md**
- 更新构建结果说明
- 只显示 Apple Silicon 版本

✅ **PROJECT_FILES.txt**
- 更新构建产物路径

## 🚀 使用方法

### 构建

```bash
# 方式 1: 使用构建脚本
./build.sh

# 方式 2: 使用 Makefile
make build

# 方式 3: 快速构建
make quick
```

### 构建结果

```
submit-sitemap/
├── submit-sitemap                      # 当前平台通用版本
└── dist/
    └── submit-sitemap-darwin-arm64     # Apple Silicon 版本
```

### 运行

```bash
# 使用当前目录的通用版本
./submit-sitemap run

# 使用 dist 目录中的版本
./dist/submit-sitemap-darwin-arm64 run
```

## ✅ 验证

所有功能已测试通过：
- ✓ 构建脚本正常工作
- ✓ Makefile 正常工作
- ✓ 构建产物输出到 dist 目录
- ✓ 只生成 Apple Silicon 版本
- ✓ 可执行文件正常运行

## 📝 注意事项

1. **旧的 build 目录已删除**：如果你有旧的 build 目录，请手动删除
2. **gitignore 已更新**：dist 目录不会被提交到 git
3. **Intel Mac 用户**：如需 Intel 版本，可手动构建：
   ```bash
   GOOS=darwin GOARCH=amd64 go build -o submit-sitemap-intel ./cmd/submit-sitemap
   ```

## 🎯 优势

1. **更清晰的语义**：dist (distribution) 更明确表示这是发布产物
2. **符合惯例**：许多项目使用 dist 作为构建产物目录
3. **精简构建**：只构建需要的 Apple Silicon 版本，节省构建时间
