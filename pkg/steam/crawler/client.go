// Package crawler 提供 Steam 网页爬虫核心能力封装
// 基于 Colly 框架构建, 整合智能反爬策略、动态代理轮换、结构化解析/存储能力, 适配 Steam 风控规则
// Package crawler provides core encapsulation for Steam web crawling capabilities
// Built on Colly framework, integrates intelligent anti-crawling strategies, dynamic proxy rotation, structured parsing/storage, adapts to Steam risk control rules

package crawler

import (
	"fmt"

	"github.com/GoFurry/gf-steam-sdk/internal/crawler"
	"github.com/GoFurry/gf-steam-sdk/pkg/config"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

// CrawlerService is the core structure of crawling service 爬虫服务核心结构体
// Aggregates anti-crawl strategy, proxy rotation, parser, storage manager, supports chain config for custom crawling rules
type CrawlerService struct {
	cfg          *config.SteamConfig   // 全局配置 | Global config (chain config entry)
	colly        *colly.Collector      // Colly 核心爬虫实例 | Core Colly crawler instance
	antiCrawl    *crawler.AntiCrawl    // 内部反爬策略 | Internal anti-crawl strategy (delay/QPS limit/UA random)
	parser       *crawler.Parser       // 内部解析器 | Internal parser (HTML structured parsing)
	storage      *crawler.Storage      // 内部存储管理器 | Internal storage manager (HTML/data persistence)
	proxyRotator *crawler.ProxyRotator // 代理轮换管理器 | Proxy rotation manager (dynamic proxy pool switching)
}

// NewCrawlerService Create crawler service instance 创建爬虫服务实例
func NewCrawlerService(cfg *config.SteamConfig) *CrawlerService {
	if cfg.IsDebug {
		fmt.Printf("[Info] Start NewCrawlerService Init \n")
	}

	// 初始化 Colly 基础实例(基础爬虫配置) | Initialize Colly base instance (basic crawler config)
	c := colly.NewCollector(
		colly.UserAgent(cfg.CrawlerUserAgent), // 自定义User-Agent | Custom User-Agent
		colly.AllowURLRevisit(),               // 允许重复访问URL | Allow URL revisit
		colly.Async(cfg.CrawlerAsync),         // 异步请求 | Async request (improve crawling efficiency)
		colly.MaxDepth(cfg.CrawlerMaxDepth),   // 爬取最大深度 | Max crawling depth (prevent infinite recursion)
	)

	// 初始化代理轮换器 | Initialize proxy rotator (dynamic proxy pool switching)
	proxyRotator := crawler.NewProxyRotator(cfg)
	// 设置代理函数 | Set proxy function (chain-configured proxy strategy)
	c.SetProxyFunc(proxyRotator.GetProxyFunc())

	// 基础反爬扩展 | Basic anti-crawl extensions
	extensions.RandomUserAgent(c)    // 随机切换User-Agent | Randomize User-Agent (avoid UA risk control)
	extensions.Referer(c)            // 设置合法Referer头 | Set valid Referer header (simulate real browser)
	c.SetRequestTimeout(cfg.Timeout) // 请求超时配置 | Request timeout config (chain-config item)

	// 集成内部反爬策略 | Integrate internal anti-crawl strategy (delay/QPS limit/retry)
	antiCrawl := crawler.NewAntiCrawl(cfg)
	antiCrawl.Apply(c)

	// 初始化内部工具 | Initialize internal tools
	parser := crawler.NewParser() // HTML结构化解析器 | HTML structured parser

	if cfg.IsDebug {
		fmt.Printf("[Info] End NewCrawlerService Init \n")
	}

	return &CrawlerService{
		cfg:          cfg,
		colly:        c,
		antiCrawl:    antiCrawl,
		parser:       parser,
		storage:      &crawler.Storage{},
		proxyRotator: proxyRotator,
	}
}

// GetProxyPool get proxy address list 获取当前代理池列表
//   - []string: Proxy address list (format: http://ip:port)
func (s *CrawlerService) GetProxyPool() []string {
	return s.proxyRotator.Pool()
}
