package types

import "time"

// Config 主配置结构（运行时使用，包含所有站点）
type Config struct {
	Sites    []SiteConfig    `yaml:"sites"`
	Settings GlobalSettings  `yaml:"settings"`
}

// SiteConfigFile 单个站点配置文件结构（用于读取单个配置文件）
type SiteConfigFile struct {
	SiteConfig `yaml:",inline"`
	Settings   *GlobalSettings `yaml:"settings,omitempty"`
}

// SiteConfig 单个网站配置
type SiteConfig struct {
	Name       string            `yaml:"name"`
	Domain     string            `yaml:"domain"`
	SitemapURL string            `yaml:"sitemap_url"`
	Quotas     QuotaConfig       `yaml:"quotas"`
	API        APIConfig         `yaml:"api"`
}

// QuotaConfig 每日提交配额
type QuotaConfig struct {
	Baidu  int `yaml:"baidu"`
	Bing   int `yaml:"bing"`
	Google int `yaml:"google"`
}

// APIConfig API配置
type APIConfig struct {
	Baidu  BaiduConfig  `yaml:"baidu"`
	Bing   BingConfig   `yaml:"bing"`
	Google GoogleConfig `yaml:"google"`
}

// BaiduConfig 百度API配置
type BaiduConfig struct {
	Token string `yaml:"token"`
	Site  string `yaml:"site"`
}

// BingConfig Bing API配置 (IndexNow)
type BingConfig struct {
	APIKey      string `yaml:"api_key"`
	Host        string `yaml:"host"`
	KeyLocation string `yaml:"key_location"`
}

// GoogleConfig Google/IndexNow API配置
type GoogleConfig struct {
	APIKey      string `yaml:"api_key"`
	Host        string `yaml:"host"`
	KeyLocation string `yaml:"key_location"`
}

// GlobalSettings 全局设置
type GlobalSettings struct {
	SitemapCacheHours int    `yaml:"sitemap_cache_hours"`
	Timeout           int    `yaml:"timeout"`
	Concurrent        int    `yaml:"concurrent"`
	LogLevel          string `yaml:"log_level"`
}

// SubmitResult 提交结果
type SubmitResult struct {
	Platform    string
	TotalCount  int
	SuccessCount int
	FailedCount int
	FailedURLs  []string
	Error       error
}

// SubmitStats 提交统计
type SubmitStats struct {
	Site         string
	Platform     string
	SubmitCount  int
	SuccessCount int
	FailedCount  int
	TotalURLs    int
	SubmittedURLs int
	Timestamp    time.Time
}

// SitemapURL sitemap中的URL信息
type SitemapURL struct {
	Loc        string
	LastMod    string
	ChangeFreq string
	Priority   string
}
