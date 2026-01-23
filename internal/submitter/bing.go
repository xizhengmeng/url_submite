package submitter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/k12/submit-sitemap/pkg/types"
)

// BingSubmitter Bing提交器 (使用IndexNow协议)
type BingSubmitter struct {
	client      *http.Client
	apiKey      string
	host        string
	keyLocation string
}

// IndexNowRequest IndexNow API请求结构
type IndexNowRequest struct {
	Host            string   `json:"host"`
	Key             string   `json:"key"`
	KeyLocation     string   `json:"keyLocation"`
	URLList         []string `json:"urlList"`
}

// NewBingSubmitter 创建Bing提交器
func NewBingSubmitter(config types.BingConfig, timeout int) *BingSubmitter {
	keyLocation := config.KeyLocation
	// 如果没有配置 keyLocation，使用默认格式
	if keyLocation == "" {
		keyLocation = fmt.Sprintf("https://%s/%s.txt", config.Host, config.APIKey)
	}

	return &BingSubmitter{
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		apiKey:      config.APIKey,
		host:        config.Host,
		keyLocation: keyLocation,
	}
}

// Submit 提交URL到Bing (使用IndexNow)
func (b *BingSubmitter) Submit(urls []string) types.SubmitResult {
	result := types.SubmitResult{
		Platform:   "Bing",
		TotalCount: len(urls),
	}

	if len(urls) == 0 {
		return result
	}

	// IndexNow API地址 (Bing)
	apiURL := "https://api.indexnow.org/indexnow"

	// 准备请求体
	reqBody := IndexNowRequest{
		Host:        b.host,
		Key:         b.apiKey,
		KeyLocation: b.keyLocation,
		URLList:     urls,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		result.Error = fmt.Errorf("序列化请求失败: %w", err)
		result.FailedCount = len(urls)
		return result
	}

	// 记录请求详情（用于调试）
	keyPreview := b.apiKey
	if len(keyPreview) > 8 {
		keyPreview = keyPreview[:8] + "***"
	}
	fmt.Printf("    [调试] 请求URL: %s\n", apiURL)
	fmt.Printf("    [调试] Host: %s\n", b.host)
	fmt.Printf("    [调试] Key: %s\n", keyPreview)
	fmt.Printf("    [调试] KeyLocation: %s\n", b.keyLocation)
	fmt.Printf("    [调试] URL数量: %d\n", len(urls))

	// 创建请求
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		result.Error = fmt.Errorf("创建请求失败: %w", err)
		result.FailedCount = len(urls)
		return result
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
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
	respBody, _ := io.ReadAll(resp.Body)

	// IndexNow返回200或202表示成功
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
		result.SuccessCount = len(urls)
		result.FailedCount = 0
	} else {
		result.Error = fmt.Errorf("提交失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
		result.FailedCount = len(urls)
		result.FailedURLs = urls
	}

	return result
}

// Name 返回提交器名称
func (b *BingSubmitter) Name() string {
	return "Bing"
}
