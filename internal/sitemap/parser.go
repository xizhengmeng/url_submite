package sitemap

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/k12/submit-sitemap/pkg/types"
)

// URLSet sitemap XML结构
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []URL    `xml:"url"`
}

// URL sitemap中的URL结构
type URL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod"`
	ChangeFreq string `xml:"changefreq"`
	Priority   string `xml:"priority"`
}

// SitemapIndex sitemap索引XML结构
type SitemapIndex struct {
	XMLName  xml.Name  `xml:"sitemapindex"`
	Sitemaps []Sitemap `xml:"sitemap"`
}

// Sitemap sitemap索引中的sitemap
type Sitemap struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod"`
}

// Parser sitemap解析器
type Parser struct {
	client  *http.Client
	timeout time.Duration
	verbose bool // 是否输出详细日志
}

// NewParser 创建新的解析器
func NewParser(timeout int) *Parser {
	return &Parser{
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		timeout: time.Duration(timeout) * time.Second,
		verbose: false,
	}
}

// NewParserWithVerbose 创建带详细日志的解析器
func NewParserWithVerbose(timeout int, verbose bool) *Parser {
	return &Parser{
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		timeout: time.Duration(timeout) * time.Second,
		verbose: verbose,
	}
}

// Parse 解析sitemap URL
func (p *Parser) Parse(sitemapURL string) ([]types.SitemapURL, error) {
	if p.verbose {
		fmt.Printf("📥 正在获取: %s\n", sitemapURL)
	}

	// 下载sitemap内容
	content, err := p.fetch(sitemapURL)
	if err != nil {
		return nil, fmt.Errorf("获取sitemap失败: %w", err)
	}

	// 尝试解析为sitemap索引
	var sitemapIndex SitemapIndex
	if err := xml.Unmarshal(content, &sitemapIndex); err == nil && len(sitemapIndex.Sitemaps) > 0 {
		if p.verbose {
			fmt.Printf("✓ 识别为 Sitemap Index，包含 %d 个子 sitemap\n", len(sitemapIndex.Sitemaps))
		}
		return p.parseIndex(sitemapIndex)
	}

	// 解析为普通sitemap
	var urlset URLSet
	if err := xml.Unmarshal(content, &urlset); err != nil {
		return nil, fmt.Errorf("解析sitemap失败: %w", err)
	}

	if p.verbose {
		fmt.Printf("✓ 识别为 URLSet，包含 %d 个URL\n", len(urlset.URLs))
	}

	return p.convertURLs(urlset.URLs), nil
}

// parseIndex 解析sitemap索引
func (p *Parser) parseIndex(index SitemapIndex) ([]types.SitemapURL, error) {
	var allURLs []types.SitemapURL

	if p.verbose {
		fmt.Printf("\n🔄 开始递归解析 %d 个子 sitemap...\n", len(index.Sitemaps))
	}

	for i, sitemap := range index.Sitemaps {
		if p.verbose {
			fmt.Printf("\n[%d/%d] ", i+1, len(index.Sitemaps))
		}

		urls, err := p.Parse(sitemap.Loc)
		if err != nil {
			// 记录错误但继续处理其他sitemap
			fmt.Printf("⚠️  警告: 解析子sitemap失败 (%s): %v\n", sitemap.Loc, err)
			continue
		}
		allURLs = append(allURLs, urls...)

		if p.verbose {
			fmt.Printf("   累计 URL: %d\n", len(allURLs))
		}
	}

	if p.verbose {
		fmt.Printf("\n✅ 递归解析完成！总计找到 %d 个URL\n\n", len(allURLs))
	}

	return allURLs, nil
}

// fetch 下载sitemap内容
func (p *Parser) fetch(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Submit-Sitemap-Bot/1.0")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP状态码: %d", resp.StatusCode)
	}

	var reader io.Reader = resp.Body

	// 如果是gzip压缩，解压
	if strings.HasSuffix(url, ".gz") || resp.Header.Get("Content-Encoding") == "gzip" {
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("解压gzip失败: %w", err)
		}
		defer gzReader.Close()
		reader = gzReader
	}

	return io.ReadAll(reader)
}

// convertURLs 转换URL格式
func (p *Parser) convertURLs(urls []URL) []types.SitemapURL {
	result := make([]types.SitemapURL, len(urls))
	for i, u := range urls {
		result[i] = types.SitemapURL{
			Loc:        u.Loc,
			LastMod:    u.LastMod,
			ChangeFreq: u.ChangeFreq,
			Priority:   u.Priority,
		}
	}
	return result
}
