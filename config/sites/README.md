# 站点配置目录

## 目录结构

每个网站在 `config/sites/` 下有一个独立的子文件夹，文件夹名称建议使用域名。每个子文件夹中包含一个或多个 `.yaml` 配置文件。

```
config/sites/
├── kebenwang.cn/
│   └── site.yaml
├── example.com/
│   └── site.yaml
└── blog.example.com/
    └── site.yaml
```

## 配置文件格式

每个配置文件包含单个网站的配置信息：

```yaml
# 站点基本信息
name: "站点名称"
domain: example.com
sitemap_url: "https://example.com/sitemap.xml"

# 每日提交配额
quotas:
  baidu: 100      # 百度每天提交100条
  bing: 200       # Bing每天提交200条
  google: 150     # Google每天提交150条

# API配置
api:
  baidu:
    token: "your-baidu-token"
    site: "https://example.com"

  bing:
    api_key: "your-indexnow-key"
    host: "example.com"
    key_location: "https://example.com/your-key.txt"

  google:
    api_key: "your-indexnow-key"
    host: "example.com"
    key_location: "https://example.com/your-key.txt"

# 全局设置（可选）
settings:
  sitemap_cache_hours: 168
  timeout: 30
  concurrent: 3
  log_level: info
```

## 添加新网站

1. 在 `config/sites/` 下创建新的文件夹（建议以域名命名）
2. 在文件夹中创建 `site.yaml` 配置文件
3. 填写网站配置信息
4. 运行 `submit-sitemap test` 测试配置

## 注意事项

- 配置文件名可以是任意 `.yaml` 或 `.yml` 后缀
- 程序会自动遍历所有子文件夹中的 yaml 文件
- 如果多个配置文件中都设置了 `settings`，会使用第一个遇到的配置
- 建议每个网站的 `settings` 保持一致，或者只在一个配置文件中设置
