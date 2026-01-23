# 故障排除指南

## 常见问题

### 1. API 配置未设置或不完整

**问题现象：**
```
[Bing] ⚠️  跳过（API配置未设置或不完整）
[Google] ⚠️  跳过（API配置未设置或不完整）
```

**原因：**
- API 配置被注释掉或者字段为空
- 必需字段缺失

**解决方法：**

检查配置文件中的 API 配置：

#### 百度配置检查
```yaml
api:
  baidu:
    token: "your-token"     # 必需：不能为空
    site: "https://..."     # 必需：不能为空
```

#### Bing 配置检查
```yaml
api:
  bing:
    api_key: "your-key"            # 必需：不能为空
    host: "example.com"            # 必需：不能为空
    key_location: "https://..."    # 可选：如果不设置会自动生成
```

#### Google 配置检查
```yaml
api:
  google:
    api_key: "your-key"            # 必需：不能为空
    host: "example.com"            # 必需：不能为空
    key_location: "https://..."    # 可选：如果不设置会自动生成
```

**如果不想提交到某个平台：**

方法一：将配额设为 0
```yaml
quotas:
  baidu: 10
  bing: 0      # 不提交到 Bing
  google: 0    # 不提交到 Google
```

方法二：注释掉整个 API 配置段（包括配额）
```yaml
quotas:
  baidu: 10
  # bing: 200    # 注释掉
  # google: 150  # 注释掉

api:
  baidu:
    token: "..."
  # bing 配置全部注释掉
  # google 配置全部注释掉
```

---

### 2. IndexNow keyLocation 无效

**问题现象：**
```
[Bing] ❌ 提交失败: 提交失败，状态码: 400, 响应: {"errorCode":"InvalidRequestParameters","message":"Invalid request","details":[{"target":"keyLocation","message":"Keylocation url is invalid."}]}
```

**原因：**
1. key_location URL 格式不正确
2. key_location 文件不存在或无法访问
3. key_location 文件内容与 api_key 不匹配

**解决方法：**

#### 步骤 1: 检查配置
```yaml
api:
  bing:
    api_key: "9d8d89d53cd1f77b63f949ee80e64d5f"
    host: "www.kebenwang.cn"
    key_location: "https://www.kebenwang.cn/9d8d89d53cd1f77b63f949ee80e64d5f.txt"
```

注意：
- `host` 应该包含完整的主机名（如 `www.kebenwang.cn`），不要只写 `kebenwang.cn`
- `key_location` 必须是完整的 HTTPS URL
- 文件名必须与 `api_key` 一致（`{api_key}.txt`）

#### 步骤 2: 创建验证文件

```bash
# 创建验证文件
echo "9d8d89d53cd1f77b63f949ee80e64d5f" > 9d8d89d53cd1f77b63f949ee80e64d5f.txt
```

#### 步骤 3: 上传验证文件到网站根目录

将文件上传到网站根目录，确保可以通过以下 URL 访问：
```
https://www.kebenwang.cn/9d8d89d53cd1f77b63f949ee80e64d5f.txt
```

#### 步骤 4: 验证文件可访问性

在浏览器中访问 key_location URL，应该返回：
```
9d8d89d53cd1f77b63f949ee80e64d5f
```

使用 curl 测试：
```bash
curl https://www.kebenwang.cn/9d8d89d53cd1f77b63f949ee80e64d5f.txt
```

---

### 3. 百度提交失败

**问题现象：**
```
[百度] ⚠ 成功: 0, 失败: 100
```

**可能原因：**
1. Token 不正确
2. Site URL 不匹配
3. URL 格式不正确
4. 已超过每日配额

**排查步骤：**

#### 步骤 1: 检查 token 和 site
```yaml
api:
  baidu:
    token: "Z8GNNjF99Oji19zi"
    site: "https://www.kebenwang.cn"  # 必须与百度站长平台配置一致
```

#### 步骤 2: 查看详细日志

使用 verbose 模式运行：
```bash
./submit run -v
```

日志中会显示：
- 百度API的响应详情
- 失败的 URL 列表
- 失败原因（not_same_site 或 not_valid）

#### 步骤 3: 常见失败原因

**not_same_site**: URL 不属于配置的站点
- 检查 URL 的域名是否与 `site` 配置匹配
- 确保 sitemap 中的 URL 都属于同一个站点

**not_valid**: URL 格式无效
- URL 必须以 `http://` 或 `https://` 开头
- URL 不能包含特殊字符或空格
- URL 必须可以正常访问

#### 步骤 4: 测试单个 URL

可以使用 curl 手动测试提交：
```bash
curl -H "Content-Type: text/plain" \
  --data "https://www.kebenwang.cn/page1" \
  "http://data.zz.baidu.com/urls?site=https://www.kebenwang.cn&token=Z8GNNjF99Oji19zi"
```

---

### 4. 提交数量为 0

**问题现象：**
```
[百度] 无新URL需要提交
```

**原因：**
1. Sitemap 中的所有 URL 都已经提交过
2. Sitemap 解析失败
3. 配额为 0

**解决方法：**

#### 方法 1: 查看历史记录
```bash
./submit stats
```

#### 方法 2: 重置历史记录
```bash
# 清除指定站点的所有历史记录
./submit reset kebenwang.cn

# 确认后重新提交
./submit run
```

#### 方法 3: 检查 Sitemap
确保 sitemap.xml 可以访问且包含 URL：
```bash
curl https://www.kebenwang.cn/sitemap.xml
```

---

## 调试技巧

### 1. 启用详细日志

使用 `-v` 参数查看详细输出：
```bash
./submit run -v
```

### 2. 查看日志文件

日志文件位于：
```
data/logs/YYYY-MM-DD.log
```

查看今天的日志：
```bash
cat data/logs/$(date +%Y-%m-%d).log
```

实时监控日志：
```bash
tail -f data/logs/$(date +%Y-%m-%d).log
```

### 3. 测试配置

运行配置测试：
```bash
./submit test
```

会显示：
- 配置目录路径
- 加载的站点数量
- 每个站点的配置详情

### 4. 查看 API 配置状态

在日志中搜索 "API配置"：
```bash
grep "API配置" data/logs/$(date +%Y-%m-%d).log
```

会显示每个平台的配置状态和掩码后的 API 密钥。

---

## 获取帮助

### 查看程序版本和帮助
```bash
./submit help
```

### 常用命令
```bash
# 测试配置
./submit test

# 查看统计
./submit stats

# 运行提交
./submit run

# 详细模式运行
./submit run -v

# 重置历史
./submit reset <domain>
```

### 日志关键字

在日志中搜索以下关键字可以快速定位问题：
- `[ERROR]` - 错误信息
- `[WARNING]` - 警告信息
- `跳过平台` - 配置问题
- `提交失败` - 提交错误
- `API配置无效` - 配置验证失败

---

## 最佳实践

### 1. 渐进式测试

不要一次性提交大量 URL，建议：
```yaml
quotas:
  baidu: 10    # 先从小配额开始
  bing: 20
  google: 20
```

测试成功后再逐步增加。

### 2. 定期检查

设置定时任务前，先手动运行几次确保一切正常：
```bash
# 第一天
./submit run -v

# 检查日志
cat data/logs/$(date +%Y-%m-% d).log

# 查看统计
./submit stats

# 确认无误后设置 cron
```

### 3. 备份配置

配置文件包含敏感信息，建议：
- 不要提交到公开仓库
- 定期备份到安全位置
- 使用环境变量或密钥管理工具

### 4. 监控配额

百度站长平台有配额限制，建议：
- 登录百度站长平台查看剩余配额
- 根据实际配额调整程序配置
- 避免配置过大的配额导致浪费
