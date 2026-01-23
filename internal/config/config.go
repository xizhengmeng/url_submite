package config

import (
	"fmt"
	"os"

	"github.com/k12/submit-sitemap/pkg/types"
	"gopkg.in/yaml.v3"
)

// Load 加载配置文件
func Load(configPath string) (*types.Config, error) {
	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %s", configPath)
	}

	// 读取文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析YAML
	var config types.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证配置
	if err := validate(&config); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	// 设置默认值
	setDefaults(&config)

	return &config, nil
}

// validate 验证配置
func validate(config *types.Config) error {
	if len(config.Sites) == 0 {
		return fmt.Errorf("至少需要配置一个网站")
	}

	for i, site := range config.Sites {
		if site.Domain == "" {
			return fmt.Errorf("网站 #%d: domain 不能为空", i+1)
		}
		if site.SitemapURL == "" {
			return fmt.Errorf("网站 #%d (%s): sitemap_url 不能为空", i+1, site.Domain)
		}

		// 检查至少配置了一个平台
		hasQuota := site.Quotas.Baidu > 0 || site.Quotas.Bing > 0 || site.Quotas.Google > 0
		if !hasQuota {
			return fmt.Errorf("网站 #%d (%s): 至少需要配置一个平台的配额", i+1, site.Domain)
		}
	}

	return nil
}

// setDefaults 设置默认值
func setDefaults(config *types.Config) {
	if config.Settings.SitemapCacheHours == 0 {
		config.Settings.SitemapCacheHours = 168 // 默认7天
	}
	if config.Settings.Timeout == 0 {
		config.Settings.Timeout = 30 // 默认30秒
	}
	if config.Settings.Concurrent == 0 {
		config.Settings.Concurrent = 3 // 默认并发3
	}
	if config.Settings.LogLevel == "" {
		config.Settings.LogLevel = "info" // 默认info级别
	}

	// 为每个网站设置默认名称
	for i := range config.Sites {
		if config.Sites[i].Name == "" {
			config.Sites[i].Name = config.Sites[i].Domain
		}
	}
}
