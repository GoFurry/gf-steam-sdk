package crawler

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/GoFurry/gf-steam-sdk/internal/crawler"
	"github.com/GoFurry/gf-steam-sdk/pkg/util/errors"
	"github.com/gocolly/colly"
)

// ============================ Raw HTML 通用获取原始 HTML ============================

// GetRawHTML crawl any address to get HTML 爬取任意地址的原始 HTML (通用爬取, 无任何跳过验证的策略)
func (s *CrawlerService) GetRawHTML(targetURL string) ([]byte, error) {
	// 参数校验 | Parameter validation
	if targetURL == "" {
		return nil, errors.NewWithType(errors.ErrTypeParam, "target URL is empty", nil)
	}

	// 初始化变量 | Initialize variables
	var html []byte
	var reqErr error

	// 注册响应回调 | Register response callback (capture raw HTML)
	s.colly.OnResponse(func(r *colly.Response) {
		html = r.Body
	})

	// 注册错误回调 | Register error callback (capture request/response error)
	s.colly.OnError(func(r *colly.Response, err error) {
		if r != nil {
			reqErr = fmt.Errorf("response error (status: %d): %w", r.StatusCode, err)
		} else {
			reqErr = fmt.Errorf("request failed (no response): %w", err)
		}
	})

	if s.cfg.IsDebug {
		fmt.Printf("[Info] Start colly.Visit: %s \n", targetURL)
	}
	// 执行请求 | Execute request (auto trigger anti-crawl strategy)
	if err := s.colly.Visit(targetURL); err != nil {
		return nil, fmt.Errorf("%w: crawl URL failed: %v", errors.ErrTypeCrawler, err)
	}
	s.colly.Wait() // 等待异步请求完成 | Wait for async requests to complete

	// 错误检查 | Error check
	if reqErr != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrTypeCrawler, reqErr)
	}
	if len(html) == 0 {
		return nil, fmt.Errorf("%w: empty html response for URL: %s", errors.ErrTypeCrawler, targetURL)
	}

	return html, nil
}

// ============================ Save HTML 通用保存原始 HTML ============================

// SaveRawHTML crawl any address to save HTML
//   - targetURL: Target URL
//   - filename: Custom filename (auto-generate if empty)
func (s *CrawlerService) SaveRawHTML(targetURL string, filename string) (string, error) {
	if s.cfg.IsDebug {
		fmt.Printf("[Info] Start GetRawHTML \n")
	}
	// 获取原始 HTML | Get raw HTML
	html, err := s.GetRawHTML(targetURL)
	if err != nil {
		return "", err
	}

	// 自动生成文件名 | Auto-generate filename (avoid special chars/long filename)
	if filename == "" {
		filename = generateFilenameFromURL(targetURL)
	}
	if s.cfg.IsDebug {
		fmt.Printf("[Info] generateFilenameFromURL %s\n", filename)
	}

	if s.cfg.IsDebug {
		fmt.Printf("[Info] NewStorage Init \n")
	}
	// 初始化存储管理器并保存 | Initialize storage manager and save
	s.storage = crawler.NewStorage(s.cfg.CrawlerStorageDir)

	if s.cfg.IsDebug {
		fmt.Printf("[Info] Start SaveHTML \n")
	}
	fullPath, err := s.storage.SaveHTML(filename, html)
	if err != nil {
		return "", fmt.Errorf("%w: save generic HTML failed: %v", errors.ErrTypeCrawler, err)
	}

	return fullPath, nil
}

// generateFilenameFromURL 基于URL自动生成合法文件名
// 替换特殊字符、截断超长名称，兜底使用时间戳命名，适配跨平台存储规则
// eg: https://store.steampowered.com/app/550/ → store.steampowered.com_app_550.html
func generateFilenameFromURL(targetURL string) string {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		// 兜底: 纳秒级时间戳命名 | Fallback: nanosecond timestamp (avoid duplication)
		return fmt.Sprintf("crawl_%d.html", time.Now().UnixNano())
	}

	// 处理路径 | Process path (replace special chars, keep core identifier)
	path := strings.Trim(parsedURL.Path, "/")
	path = strings.ReplaceAll(path, "/", "_")
	path = strings.ReplaceAll(path, "?", "_")
	path = strings.ReplaceAll(path, "&", "_")
	path = strings.ReplaceAll(path, "=", "_")

	// 生成文件名 | Generate filename
	filename := fmt.Sprintf("%s_%s.html", parsedURL.Host, path)
	// 截断超长文件名 | Truncate long filename (avoid system limits)
	if len(filename) > 50 {
		filename = filename[:50] + ".html"
	}

	// 路径为空时补充标识 | Add identifier if path is empty
	if filename == fmt.Sprintf("%s_.html", parsedURL.Host) {
		filename = fmt.Sprintf("%s_home_%d.html", parsedURL.Host, time.Now().Unix())
	}

	return filename
}
