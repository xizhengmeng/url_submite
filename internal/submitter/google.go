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

// GoogleSubmitter Google提交器 (使用IndexNow协议)
// 注意: Google官方的Indexing API仅支持特定类型内容
// 这里使用IndexNow协议，虽然Google目前不支持，但保留接口便于未来扩展
type GoogleSubmitter struct {
	client      *http.Client
	apiKey      string
	host        string
	keyLocation string
}

// NewGoogleSubmitter 创建Google提交器
func NewGoogleSubmitter(config types.GoogleConfig, timeout int) *GoogleSubmitter {
	keyLocation := config.KeyLocation
	// 如果没有配置 keyLocation，使用默认格式
	if keyLocation == "" {
		keyLocation = fmt.Sprintf("https://%s/%s.txt", config.Host, config.APIKey)
	}

	return &GoogleSubmitter{
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		apiKey:      config.APIKey,
		host:        config.Host,
		keyLocation: keyLocation,
	}
}

// Submit 提交URL到Google
func (g *GoogleSubmitter) Submit(urls []string) types.SubmitResult {
	result := types.SubmitResult{
		Platform:   "Google",
		TotalCount: len(urls),
	}

	if len(urls) == 0 {
		return result
	}

	// 使用IndexNow协议提交
	// 注意：Google目前不支持IndexNow，但保留此实现以便未来使用
	apiURL := "https://api.indexnow.org/indexnow"

	// 准备请求体
	reqBody := IndexNowRequest{
		Host:        g.host,
		Key:         g.apiKey,
		KeyLocation: g.keyLocation,
		URLList:     urls,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		result.Error = fmt.Errorf("序列化请求失败: %w", err)
		result.FailedCount = len(urls)
		return result
	}

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
	resp, err := g.client.Do(req)
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
func (g *GoogleSubmitter) Name() string {
	return "Google"
}
