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

	baiduResp, statusCode, respBody, err := b.submitRaw(urls)
	if err != nil {
		result.Error = err
		result.FailedCount = len(urls)
		result.FailedURLs = append(result.FailedURLs, urls...)
		return result
	}

	if statusCode != http.StatusOK {
		if statusCode == http.StatusBadRequest && isOverQuota(respBody) && len(urls) > 1 {
			return b.submitOneByOne(urls)
		}
		result.Error = fmt.Errorf("HTTP错误 - 状态码: %d, 响应: %s", statusCode, string(respBody))
		result.FailedCount = len(urls)
		result.FailedURLs = append(result.FailedURLs, urls...)
		return result
	}

	// 设置结果
	result.SuccessCount = baiduResp.Success
	result.FailedCount = len(urls) - baiduResp.Success

	// 记录失败的URL
	result.FailedURLs = append(result.FailedURLs, baiduResp.NotSameSite...)
	result.FailedURLs = append(result.FailedURLs, baiduResp.NotValid...)

	// 如果有失败的URL，添加错误信息
	if len(result.FailedURLs) > 0 {
		failedInfo := fmt.Sprintf("not_same_site: %d, not_valid: %d, remain: %d",
			len(baiduResp.NotSameSite), len(baiduResp.NotValid), baiduResp.Remain)
		result.Error = fmt.Errorf("部分URL提交失败: %s", failedInfo)
	}

	return result
}

func (b *BaiduSubmitter) submitRaw(urls []string) (BaiduResponse, int, []byte, error) {
	var empty BaiduResponse

	apiURL := fmt.Sprintf("http://data.zz.baidu.com/urls?site=%s&token=%s", b.site, b.token)
	body := strings.Join(urls, "\n")

	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(body))
	if err != nil {
		return empty, 0, nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "text/plain")

	resp, err := b.client.Do(req)
	if err != nil {
		return empty, 0, nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return empty, resp.StatusCode, nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return empty, resp.StatusCode, respBody, nil
	}

	var baiduResp BaiduResponse
	if err := json.Unmarshal(respBody, &baiduResp); err != nil {
		return empty, resp.StatusCode, respBody, fmt.Errorf("解析响应失败: %w (响应: %s)", err, string(respBody))
	}

	return baiduResp, resp.StatusCode, respBody, nil
}

func (b *BaiduSubmitter) submitOneByOne(urls []string) types.SubmitResult {
	result := types.SubmitResult{
		Platform:   "百度",
		TotalCount: len(urls),
	}

	for i, u := range urls {
		baiduResp, statusCode, respBody, err := b.submitRaw([]string{u})
		if err != nil {
			result.Error = err
			result.FailedURLs = append(result.FailedURLs, urls[i:]...)
			result.FailedCount = len(result.FailedURLs)
			return result
		}

		if statusCode != http.StatusOK {
			if statusCode == http.StatusBadRequest && isOverQuota(respBody) {
				result.Error = fmt.Errorf("over quota")
				result.FailedURLs = append(result.FailedURLs, urls[i:]...)
				result.FailedCount = len(result.FailedURLs)
				return result
			}
			result.Error = fmt.Errorf("HTTP错误 - 状态码: %d, 响应: %s", statusCode, string(respBody))
			result.FailedURLs = append(result.FailedURLs, u)
			result.FailedCount++
			continue
		}

		if baiduResp.Success > 0 {
			result.SuccessCount += baiduResp.Success
			continue
		}

		if len(baiduResp.NotSameSite) > 0 || len(baiduResp.NotValid) > 0 {
			result.FailedURLs = append(result.FailedURLs, baiduResp.NotSameSite...)
			result.FailedURLs = append(result.FailedURLs, baiduResp.NotValid...)
			result.FailedCount += len(baiduResp.NotSameSite) + len(baiduResp.NotValid)
			continue
		}

		result.FailedURLs = append(result.FailedURLs, u)
		result.FailedCount++
	}

	if result.FailedCount > 0 {
		result.Error = fmt.Errorf("部分URL提交失败: %d", result.FailedCount)
	}

	return result
}

func isOverQuota(respBody []byte) bool {
	msg := strings.ToLower(string(respBody))
	return strings.Contains(msg, "over quota")
}

// Name 返回提交器名称
func (b *BaiduSubmitter) Name() string {
	return "百度"
}
