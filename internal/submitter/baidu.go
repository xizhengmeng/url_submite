package submitter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/k12/submit-sitemap/pkg/types"
)

// BaiduSubmitter 百度提交器
type BaiduSubmitter struct {
	client  *http.Client
	token   string
	site    string
}

// BaiduResponse 百度API响应
type BaiduResponse struct {
	Success     int      `json:"success"`
	Remain      int      `json:"remain"`
	NotSameSite []string `json:"not_same_site"`
	NotValid    []string `json:"not_valid"`
}

// NewBaiduSubmitter 创建百度提交器
func NewBaiduSubmitter(config types.BaiduConfig, timeout int) *BaiduSubmitter {
	return &BaiduSubmitter{
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		token: config.Token,
		site:  config.Site,
	}
}

// Submit 提交URL到百度
func (b *BaiduSubmitter) Submit(urls []string) types.SubmitResult {
	result := types.SubmitResult{
		Platform:   "百度",
		TotalCount: len(urls),
	}

	if len(urls) == 0 {
		return result
	}

	// 百度API地址
	apiURL := fmt.Sprintf("http://data.zz.baidu.com/urls?site=%s&token=%s", b.site, b.token)

	// 准备请求体（每行一个URL）
	body := strings.Join(urls, "\n")

	// 创建请求
	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(body))
	if err != nil {
		result.Error = fmt.Errorf("创建请求失败: %w", err)
		result.FailedCount = len(urls)
		return result
	}

	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("User-Agent", "Submit-Sitemap-Bot/1.0")

	// 发送请求
	resp, err := b.client.Do(req)
	if err != nil {
		result.Error = fmt.Errorf("发送请求失败: %w", err)
		result.FailedCount = len(urls)
		return result
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Error = fmt.Errorf("读取响应失败: %w", err)
		result.FailedCount = len(urls)
		return result
	}

	// 解析响应
	var baiduResp BaiduResponse
	if err := json.Unmarshal(respBody, &baiduResp); err != nil {
		result.Error = fmt.Errorf("解析响应失败: %w (响应: %s)", err, string(respBody))
		result.FailedCount = len(urls)
		return result
	}

	// 设置结果
	result.SuccessCount = baiduResp.Success
	result.FailedCount = len(urls) - baiduResp.Success

	// 记录失败的URL
	result.FailedURLs = append(result.FailedURLs, baiduResp.NotSameSite...)
	result.FailedURLs = append(result.FailedURLs, baiduResp.NotValid...)

	return result
}

// Name 返回提交器名称
func (b *BaiduSubmitter) Name() string {
	return "百度"
}
