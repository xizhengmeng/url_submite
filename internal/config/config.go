package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/k12/submit-sitemap/pkg/types"
	"gopkg.in/yaml.v3"
)

// LoadAll 遍历配置目录下所有子文件夹的yaml文件，加载所有站点配置
func LoadAll(configDir string) (*types.Config, error) {
	// 检查目录是否存在
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置目录不存在: %s", configDir)
	}

	// 收集所有配置文件
	var configFiles []string
	err := filepath.WalkDir(configDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 跳过目录
		if d.IsDir() {
			return nil
		}
		// 只处理 .yaml 和 .yml 文件
		if strings.HasSuffix(strings.ToLower(path), ".yaml") || strings.HasSuffix(strings.ToLower(path), ".yml") {
			configFiles = append(configFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("遍历配置目录失败: %w", err)
	}

	if len(configFiles) == 0 {
		return nil, fmt.Errorf("配置目录下未找到任何 .yaml 或 .yml 文件: %s", configDir)
	}

	// 加载所有站点配置
	var sites []types.SiteConfig
	var globalSettings *types.GlobalSettings

	for _, configFile := range configFiles {
		siteConfigFile, err := loadSiteConfigFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("加载配置文件 %s 失败: %w", configFile, err)
		}

		// 添加站点配置
		sites = append(sites, siteConfigFile.SiteConfig)

		// 使用第一个遇到的全局设置（如果有）
		if globalSettings == nil && siteConfigFile.Settings != nil {
			globalSettings = siteConfigFile.Settings
		}
	}

	// 构建完整配置
	config := &types.Config{
		Sites: sites,
	}

	// 如果有全局设置，使用它，否则使用默认值
	if globalSettings != nil {
		config.Settings = *globalSettings
	}

	// 验证配置
	if err := validate(config); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	// 设置默认值
	setDefaults(config)

	return config, nil
}

// loadSiteConfigFile 加载单个站点配置文件
func loadSiteConfigFile(configPath string) (*types.SiteConfigFile, error) {
	// 读取文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	// 解析YAML
	var siteConfig types.SiteConfigFile
	if err := yaml.Unmarshal(data, &siteConfig); err != nil {
		return nil, fmt.Errorf("解析YAML失败: %w", err)
	}

	return &siteConfig, nil
}

// Load 加载单个配置文件（保留向后兼容性）
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
