package history

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Manager 历史记录管理器
type Manager struct {
	dataDir string
	cache   map[string]map[string]map[string]bool // domain -> platform -> url -> bool
	mu      sync.RWMutex
}

// NewManager 创建历史记录管理器
func NewManager(dataDir string) *Manager {
	return &Manager{
		dataDir: dataDir,
		cache:   make(map[string]map[string]map[string]bool),
	}
}

// Load 加载指定站点和平台的历史记录
func (m *Manager) Load(domain, platform string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	filePath := m.getFilePath(domain, platform)

	// 初始化缓存
	if m.cache[domain] == nil {
		m.cache[domain] = make(map[string]map[string]bool)
	}
	if m.cache[domain][platform] == nil {
		m.cache[domain][platform] = make(map[string]bool)
	}

	// 如果文件不存在，返回空记录
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	// 读取文件
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("打开历史文件失败: %w", err)
	}
	defer file.Close()

	// 逐行读取URL
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()
		if url != "" {
			m.cache[domain][platform][url] = true
		}
	}

	return scanner.Err()
}

// IsSubmitted 检查URL是否已提交
func (m *Manager) IsSubmitted(domain, platform, url string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.cache[domain] == nil || m.cache[domain][platform] == nil {
		return false
	}

	return m.cache[domain][platform][url]
}

// Save 保存成功提交的URL
func (m *Manager) Save(domain, platform string, urls []string) error {
	if len(urls) == 0 {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	filePath := m.getFilePath(domain, platform)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 以追加模式打开文件
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("打开历史文件失败: %w", err)
	}
	defer file.Close()

	// 写入URL
	writer := bufio.NewWriter(file)
	for _, url := range urls {
		if _, err := writer.WriteString(url + "\n"); err != nil {
			return fmt.Errorf("写入URL失败: %w", err)
		}
		// 更新缓存
		if m.cache[domain] == nil {
			m.cache[domain] = make(map[string]map[string]bool)
		}
		if m.cache[domain][platform] == nil {
			m.cache[domain][platform] = make(map[string]bool)
		}
		m.cache[domain][platform][url] = true
	}

	return writer.Flush()
}

// GetCount 获取已提交的URL数量
func (m *Manager) GetCount(domain, platform string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.cache[domain] == nil || m.cache[domain][platform] == nil {
		return 0
	}

	return len(m.cache[domain][platform])
}

// Reset 清除指定站点的历史记录
func (m *Manager) Reset(domain string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	domainDir := filepath.Join(m.dataDir, "submitted", domain)
	if err := os.RemoveAll(domainDir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除历史记录失败: %w", err)
	}

	// 清除缓存
	delete(m.cache, domain)

	return nil
}

// getFilePath 获取历史文件路径
func (m *Manager) getFilePath(domain, platform string) string {
	return filepath.Join(m.dataDir, "submitted", domain, platform+".txt")
}

// FilterUnsubmitted 过滤出未提交的URL
func (m *Manager) FilterUnsubmitted(domain, platform string, urls []string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var unsubmitted []string
	for _, url := range urls {
		if !m.IsSubmitted(domain, platform, url) {
			unsubmitted = append(unsubmitted, url)
		}
	}

	return unsubmitted
}
