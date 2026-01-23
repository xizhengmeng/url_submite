# 配置文件结构迁移指南

## 变更说明

从单一配置文件模式迁移到多配置文件模式，每个网站使用独立的配置文件。

## 新架构优势

1. **独立管理**：每个网站的配置独立，便于管理和维护
2. **灵活扩展**：添加新网站只需创建新配置文件
3. **版本控制**：每个网站配置可以独立进行版本控制
4. **自动发现**：程序自动遍历配置目录，无需手动维护站点列表

## 目录结构

### 旧结构
```
config/
└── config.yaml  # 所有网站配置都在一个文件中
```

### 新结构
```
config/sites/
├── kebenwang.cn/
│   └── site.yaml
├── example.com/
│   └── site.yaml
└── _template/
    └── site.yaml.example  # 配置模板
```

## 配置文件格式变化

### 旧格式 (config/config.yaml)
```yaml
sites:
  - name: "网站1"
    domain: example.com
    # ...其他配置
  - name: "网站2"
    domain: example2.com
    # ...其他配置

settings:
  # 全局设置
```

### 新格式 (config/sites/example.com/site.yaml)
```yaml
# 直接写站点配置，不需要 sites 数组
name: "网站1"
domain: example.com
sitemap_url: "https://example.com/sitemap.xml"

quotas:
  baidu: 100
  bing: 200
  google: 150

api:
  baidu:
    token: "your-token"
    site: "https://example.com"

  # Bing 和 Google 新增配置字段
  bing:
    api_key: "your-indexnow-key"
    host: "example.com"                          # 新增
    key_location: "https://example.com/key.txt"  # 新增

  google:
    api_key: "your-indexnow-key"
    host: "example.com"                          # 新增
    key_location: "https://example.com/key.txt"  # 新增

# 可选：全局设置
settings:
  sitemap_cache_hours: 168
  timeout: 30
  concurrent: 3
  log_level: info
```

## API 配置变更

### Bing 和 Google 配置新增字段

为了更好地支持 IndexNow 协议，Bing 和 Google 的配置新增了两个字段：

- `host`: 网站主机名（不含协议）
- `key_location`: IndexNow key 文件的完整 URL

这些字段是可选的，如果不提供，程序会使用默认格式：`https://{host}/{api_key}.txt`

## 迁移步骤

### 1. 创建新的配置目录结构

```bash
mkdir -p config/sites
```

### 2. 为每个网站创建独立配置

对于旧配置文件中的每个网站：

```bash
# 创建网站配置目录
mkdir -p config/sites/example.com

# 创建配置文件
touch config/sites/example.com/site.yaml
```

### 3. 转换配置格式

将原配置文件中的每个站点配置复制到对应的新配置文件中，并：

1. 移除 `sites` 数组层级（直接从顶层开始写配置）
2. 为 Bing 和 Google 添加 `host` 和 `key_location` 字段（可选）
3. 可以保留 `settings` 部分，也可以省略使用默认值

### 4. 验证配置

```bash
# 测试配置是否正确
./submit-sitemap test

# 查看统计信息
./submit-sitemap stats
```

### 5. 备份旧配置

```bash
mv config/config.yaml config/config.yaml.backup
```

## 示例迁移

### 迁移前 (config/config.yaml)

```yaml
sites:
  - name: "课本网"
    domain: kebenwang.cn
    sitemap_url: "https://www.kebenwang.cn/sitemap.xml"
    quotas:
      baidu: 100
      bing: 200
      google: 150
    api:
      baidu:
        token: "Z8GNNjF99Oji19zi"
        site: "https://www.kebenwang.cn"
      bing:
        api_key: "9d8d89d53cd1f77b63f949ee80e64d5f"
      google:
        api_key: "9d8d89d53cd1f77b63f949ee80e64d5f"
```

### 迁移后 (config/sites/kebenwang.cn/site.yaml)

```yaml
name: "课本网"
domain: kebenwang.cn
sitemap_url: "https://www.kebenwang.cn/sitemap.xml"

quotas:
  baidu: 100
  bing: 200
  google: 150

api:
  baidu:
    token: "Z8GNNjF99Oji19zi"
    site: "https://www.kebenwang.cn"

  bing:
    api_key: "9d8d89d53cd1f77b63f949ee80e64d5f"
    host: "www.kebenwang.cn"
    key_location: "https://www.kebenwang.cn/9d8d89d53cd1f77b63f949ee80e64d5f.txt"

  google:
    api_key: "9d8d89d53cd1f77b63f949ee80e64d5f"
    host: "www.kebenwang.cn"
    key_location: "https://www.kebenwang.cn/9d8d89d53cd1f77b63f949ee80e64d5f.txt"

settings:
  sitemap_cache_hours: 168
  timeout: 30
  concurrent: 3
  log_level: info
```

## 命令行参数变更

### 旧用法
```bash
submit-sitemap run -c config/config.yaml
```

### 新用法
```bash
submit-sitemap run -c config/sites
```

现在 `-c` 参数接受配置目录而不是配置文件路径。

## 注意事项

1. **全局设置处理**：如果多个配置文件中都设置了 `settings`，程序会使用第一个遇到的配置
2. **建议实践**：所有站点的 `settings` 保持一致，或者只在一个配置文件中设置
3. **向后兼容**：保留了 `Load()` 函数用于加载单个配置文件，但主程序使用 `LoadAll()` 遍历目录
4. **IndexNow 验证**：记得将生成的 key 文件上传到网站根目录

## 获取帮助

如果遇到问题，可以：

1. 查看配置模板：`config/sites/_template/site.yaml.example`
2. 查看目录说明：`config/sites/README.md`
3. 运行测试命令：`submit-sitemap test`
