// Package crawler 实现 Steam 爬虫核心能力
// 包含反爬策略、代理轮换、HTML 解析和存储等功能
// Package crawler implements core capabilities of Steam crawler
// Includes anti-crawl strategy, proxy rotation, HTML parsing and storage

package crawler

import (
	"os"
	"path/filepath"
	"time"

	"github.com/GoFurry/gf-steam-sdk/pkg/util"
	"github.com/GoFurry/gf-steam-sdk/pkg/util/log"
)

// Storage HTML 存储管理器
// 按日期分目录存储爬取的 HTML 文件, 确保文件组织规范
// Storage is the HTML storage manager
// Stores crawled HTML files by date directory to ensure standardized file organization
type Storage struct {
	baseDir string // 基础存储目录 | Base storage directory
}

// NewStorage 创建存储管理器实例
// 参数:
//   - baseDir: 基础存储目录 | Base storage directory
//
// 返回值:
//   - *Storage: 存储管理器实例 | Storage manager instance
func NewStorage(baseDir string) *Storage {
	// 确保基础目录存在
	// Ensure base directory exists (create if not)
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		log.Error("[Storage] 创建基础目录失败", err)
	}
	return &Storage{
		baseDir: baseDir,
	}
}

// SaveHTML 保存原始 HTML 到本地文件
// 按日期分目录存储，自动创建不存在的目录
// 参数:
//   - filename: 自定义文件名(如 game_550.html) | Custom filename (e.g. game_550.html)
//   - html: 原始 HTML 字节流 | Raw HTML byte stream
//
// 返回值:
//   - string: 完整存储路径 | Full storage path
//   - error: 保存失败时返回错误 | Error if save failed
func (s *Storage) SaveHTML(filename string, html []byte) (string, error) {
	// 创建日期子目录(格式: 20060102)
	// Create date subdirectory (format: 20060102)
	dateDir := filepath.Join(s.baseDir, time.Now().Format(util.TIME_FORMAT))
	if err := os.MkdirAll(dateDir, 0755); err != nil {
		return "", err
	}

	// 拼接完整文件路径
	// Build full file path
	fullPath := filepath.Join(dateDir, filename)
	if err := os.WriteFile(fullPath, html, 0644); err != nil {
		return "", err
	}

	return fullPath, nil
}
