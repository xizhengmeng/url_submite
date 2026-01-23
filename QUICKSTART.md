# 快速开始指南

## 🚀 5分钟快速上手

### 1️⃣ 构建程序

选择以下任一方式：

```bash
# 方式1: 使用构建脚本（推荐，有彩色输出）
./build.sh

# 方式2: 使用 Makefile（简洁快速）
make

# 方式3: 快速构建当前平台
make quick
```

### 2️⃣ 配置文件

```bash
# 复制配置示例
cp config/config.example.yaml config/config.yaml

# 编辑配置文件
vim config/config.yaml  # 或使用你喜欢的编辑器
```

### 3️⃣ 填写配置

在 `config/config.yaml` 中填入你的信息：

```yaml
sites:
  - name: "我的网站"
    domain: example.com                        # 改成你的域名
    sitemap_url: "https://example.com/sitemap.xml"  # 改成你的sitemap地址

    quotas:
      baidu: 100   # 每天提交100条到百度
      bing: 200    # 每天提交200条到Bing
      google: 0    # Google暂不支持，设为0

    api:
      baidu:
        token: "your-baidu-token"        # 填入百度token
        site: "https://example.com"      # 填入你的网站地址
      bing:
        api_key: "your-bing-api-key"     # 填入Bing API密钥
```

### 4️⃣ 获取API密钥

#### 百度站长平台

1. 访问：https://ziyuan.baidu.com/
2. 添加并验证网站
3. 进入"普通收录" → "API提交"
4. 复制token（类似：`A1B2C3D4E5F6`）

#### Bing Webmaster

1. 访问：https://www.bing.com/webmasters
2. 添加并验证网站
3. 设置 → API访问 → 生成API密钥
4. 在网站根目录创建 `{api-key}.txt` 文件，内容为API密钥

### 5️⃣ 测试配置

```bash
./submit-sitemap test
```

如果看到"✨ 配置测试通过"，说明配置正确。

### 6️⃣ 运行提交

```bash
./submit-sitemap run
```

首次运行会看到类似输出：

```
🚀 URL提交工具启动...

📦 处理网站: 我的网站 (example.com)

📥 正在获取sitemap...
  ✓ 找到 1,250 个URL

📊 加载历史记录...
  • Baidu: 已提交 0 个URL
  • Bing: 已提交 0 个URL
  • Google: 已提交 0 个URL

📤 开始提交...
  [Baidu] ✓ 成功: 100, 失败: 0
  [Bing] ✓ 成功: 200, 失败: 0

✨ 完成! 本次共提交 300 个URL
```

## 📅 设置定时任务

推荐每天自动运行一次：

```bash
# 编辑crontab
crontab -e

# 添加定时任务（每天早上9点运行）
0 9 * * * cd /path/to/submit-sitemap && ./submit-sitemap run >> data/logs/cron.log 2>&1
```

## 🔍 常用命令

```bash
# 运行提交任务（带进度条）
./submit-sitemap run

# 详细输出模式（显示所有日志到控制台）
./submit-sitemap run -v

# 查看统计信息
./submit-sitemap stats

# 测试配置
./submit-sitemap test

# 清除某个站点的历史（重新提交所有URL）
./submit-sitemap reset example.com

# 查看帮助
./submit-sitemap help
```

### 📝 查看详细日志

程序会自动将详细日志保存到文件：

```bash
# 查看今天的日志
cat data/logs/$(date +%Y-%m-%d).log

# 实时查看日志
tail -f data/logs/$(date +%Y-%m-%d).log

# 查看最近的日志
ls -lt data/logs/ | head -5
```

## ⚠️ 注意事项

1. **配额设置**：不要超过平台限制
   - 百度：普通站点 500条/天，新站 3000条/天
   - Bing：建议 200-500条/天

2. **API密钥安全**：不要将 `config/config.yaml` 提交到公开仓库

3. **Sitemap更新**：工具默认缓存7天，与网站sitemap更新频率匹配

4. **失败处理**：提交失败的URL不会记录，下次运行会自动重试

## 💡 使用技巧

### 查看提交进度

```bash
./submit-sitemap stats
```

### 多站点管理

可以在一个配置文件中配置多个网站，工具会依次处理。

### 重新开始

如果想重新提交所有URL：

```bash
./submit-sitemap reset example.com
```

## 🆘 遇到问题？

1. **配置测试失败**：检查 YAML 格式是否正确
2. **Sitemap获取失败**：确认 sitemap URL 是否可访问
3. **提交失败**：验证 API 密钥是否正确

查看详细文档：`README.md`

---

**完成！** 🎉 现在你可以每天运行一次 `./submit-sitemap run` 来自动提交URL了。
