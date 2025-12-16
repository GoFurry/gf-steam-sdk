// Package crawler 实现 Steam 爬虫核心能力
// 包含反爬策略、代理轮换、HTML 解析和存储等功能
// Package crawler implements core capabilities of Steam crawler
// Includes anti-crawl strategy, proxy rotation, HTML parsing and storage

package crawler

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Parser HTML 解析器
// 封装 goquery 解析逻辑, 提供通用的 HTML 解析和文本清理能力
// Parser is the HTML parser
// Encapsulates goquery parsing logic and provides universal HTML parsing and text cleaning capabilities
type Parser struct{}

// NewParser 创建 HTML 解析器实例
// 返回值:
//   - *Parser: 解析器实例 | Parser instance
func NewParser() *Parser {
	return &Parser{}
}

// ParseHTML 解析原始 HTML 字节流
// 支持自定义解析逻辑, 灵活适配不同页面结构
// 参数:
//   - html: 原始 HTML 字节流 | Raw HTML byte stream
//   - parseFn: 自定义解析函数(接收 goquery 文档实例) | Custom parse function (receives goquery document)
//
// 返回值:
//   - error: 解析失败时返回错误 | Error if parsing fails
func (p *Parser) ParseHTML(html []byte, parseFn func(doc *goquery.Document) error) error {
	if len(html) == 0 {
		return fmt.Errorf("empty html")
	}

	// 使用 goquery 解析 HTML 字节流
	// Parse HTML byte stream with goquery
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return err
	}

	// 执行自定义解析逻辑
	// Execute custom parse logic
	if err := parseFn(doc); err != nil {
		return err
	}
	return nil
}

// CleanText 清理解析后的文本
// 去除首尾空格和换行符, 提升文本可读性
// 参数:
//   - text: 原始解析文本 | Raw parsed text
//
// 返回值:
//   - string: 清理后的文本 | Cleaned text
func (p *Parser) CleanText(text string) string {
	return strings.TrimSpace(strings.ReplaceAll(text, "\n", ""))
}
