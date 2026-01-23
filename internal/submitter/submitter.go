package submitter

import "github.com/k12/submit-sitemap/pkg/types"

// Submitter 提交器接口
type Submitter interface {
	Submit(urls []string) types.SubmitResult
	Name() string
}

// BatchSubmit 批量提交URL
// 将URLs按batchSize分批提交
func BatchSubmit(submitter Submitter, urls []string, batchSize int) []types.SubmitResult {
	if batchSize <= 0 {
		batchSize = 100 // 默认每批100条
	}

	var results []types.SubmitResult

	// 分批提交
	for i := 0; i < len(urls); i += batchSize {
		end := i + batchSize
		if end > len(urls) {
			end = len(urls)
		}

		batch := urls[i:end]
		result := submitter.Submit(batch)
		results = append(results, result)
	}

	return results
}

// MergeResults 合并多个提交结果
func MergeResults(results []types.SubmitResult) types.SubmitResult {
	if len(results) == 0 {
		return types.SubmitResult{}
	}

	merged := types.SubmitResult{
		Platform: results[0].Platform,
	}

	for _, r := range results {
		merged.TotalCount += r.TotalCount
		merged.SuccessCount += r.SuccessCount
		merged.FailedCount += r.FailedCount
		merged.FailedURLs = append(merged.FailedURLs, r.FailedURLs...)

		if r.Error != nil && merged.Error == nil {
			merged.Error = r.Error
		}
	}

	return merged
}
