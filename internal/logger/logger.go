package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Logger 日志记录器
type Logger struct {
	file      *os.File
	logger    *log.Logger
	mu        sync.Mutex
	verbose   bool
	logToFile bool
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// Init 初始化全局日志记录器
func Init(dataDir string, verbose bool) error {
	var err error
	once.Do(func() {
		defaultLogger, err = New(dataDir, verbose)
	})
	return err
}

// New 创建新的日志记录器
func New(dataDir string, verbose bool) (*Logger, error) {
	// 创建日志目录
	logDir := filepath.Join(dataDir, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %w", err)
	}

	// 创建日志文件
	logFile := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("创建日志文件失败: %w", err)
	}

	// 创建多写入器（同时写入文件和标准输出）
	var writer io.Writer
	if verbose {
		writer = io.MultiWriter(file, os.Stdout)
	} else {
		writer = file
	}

	logger := &Logger{
		file:      file,
		logger:    log.New(writer, "", log.LstdFlags),
		verbose:   verbose,
		logToFile: true,
	}

	return logger, nil
}

// Info 记录信息日志
func Info(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Info(format, v...)
	}
}

// Error 记录错误日志
func Error(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Error(format, v...)
	}
}

// Debug 记录调试日志
func Debug(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debug(format, v...)
	}
}

// Success 记录成功日志
func Success(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Success(format, v...)
	}
}

// Warning 记录警告日志
func Warning(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warning(format, v...)
	}
}

// Info 记录信息日志
func (l *Logger) Info(format string, v ...interface{}) {
	l.log("INFO", format, v...)
}

// Error 记录错误日志
func (l *Logger) Error(format string, v ...interface{}) {
	l.log("ERROR", format, v...)
}

// Debug 记录调试日志
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.verbose {
		l.log("DEBUG", format, v...)
	}
}

// Success 记录成功日志
func (l *Logger) Success(format string, v ...interface{}) {
	l.log("SUCCESS", format, v...)
}

// Warning 记录警告日志
func (l *Logger) Warning(format string, v ...interface{}) {
	l.log("WARNING", format, v...)
}

// log 内部日志记录方法
func (l *Logger) log(level, format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	message := fmt.Sprintf(format, v...)
	logLine := fmt.Sprintf("[%s] %s", level, message)

	if l.logToFile && l.logger != nil {
		l.logger.Println(logLine)
	}
}

// Close 关闭日志记录器
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Close 关闭全局日志记录器
func Close() error {
	if defaultLogger != nil {
		return defaultLogger.Close()
	}
	return nil
}

// LogSubmitStart 记录提交开始
func LogSubmitStart(site, platform string, count int) {
	Info("开始提交 - 网站: %s, 平台: %s, URL数量: %d", site, platform, count)
}

// LogSubmitResult 记录提交结果
func LogSubmitResult(site, platform string, success, failed int, urls []string) {
	if failed > 0 {
		Warning("提交完成 - 网站: %s, 平台: %s, 成功: %d, 失败: %d",
			site, platform, success, failed)
		if len(urls) > 0 {
			Warning("失败的URL列表:")
			for _, url := range urls {
				Warning("  - %s", url)
			}
		}
	} else {
		Success("提交完成 - 网站: %s, 平台: %s, 成功: %d", site, platform, success)
	}
}

// LogSitemapParsed 记录sitemap解析结果
func LogSitemapParsed(url string, count int) {
	Info("Sitemap解析完成 - URL: %s, 找到 %d 个URL", url, count)
}

// LogHistoryLoaded 记录历史加载
func LogHistoryLoaded(site, platform string, count int) {
	Info("历史记录加载 - 网站: %s, 平台: %s, 已提交: %d 个URL", site, platform, count)
}

// LogFilterResult 记录过滤结果
func LogFilterResult(site, platform string, total, unsubmitted int) {
	Info("URL过滤 - 网站: %s, 平台: %s, 总数: %d, 待提交: %d",
		site, platform, total, unsubmitted)
}

// LogConfigLoaded 记录配置加载
func LogConfigLoaded(configPath string, siteCount int) {
	Info("配置加载完成 - 文件: %s, 网站数: %d", configPath, siteCount)
}

// LogError 记录通用错误
func LogError(context string, err error) {
	Error("%s: %v", context, err)
}
