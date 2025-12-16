// Package crawler 提供 Steam 网页爬虫核心能力封装
// 基于 Colly 框架构建, 整合智能反爬策略、动态代理轮换、结构化解析/存储能力, 适配 Steam 风控规则
// Package crawler provides core encapsulation for Steam web crawling capabilities
// Built on Colly framework, integrates intelligent anti-crawling strategies, dynamic proxy rotation, structured parsing/storage, adapts to Steam risk control rules

package crawler

import (
	"github.com/GoFurry/gf-steam-sdk/internal/crawler"
	"github.com/GoFurry/gf-steam-sdk/pkg/config"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

// CrawlerService 爬虫服务核心结构体
// 聚合反爬策略、代理轮换、解析器、存储管理器等核心组件, 支持链式配置自定义爬虫规则
// CrawlerService is the core structure of crawling service
// Aggregates anti-crawl strategy, proxy rotation, parser, storage manager, supports chain config for custom crawling rules
type CrawlerService struct {
	cfg          *config.SteamConfig   // 全局配置 | Global config (chain config entry)
	colly        *colly.Collector      // Colly 核心爬虫实例 | Core Colly crawler instance
	antiCrawl    *crawler.AntiCrawl    // 内部反爬策略 | Internal anti-crawl strategy (delay/QPS limit/UA random)
	parser       *crawler.Parser       // 内部解析器 | Internal parser (HTML structured parsing)
	storage      *crawler.Storage      // 内部存储管理器 | Internal storage manager (HTML/data persistence)
	proxyRotator *crawler.ProxyRotator // 代理轮换管理器 | Proxy rotation manager (dynamic proxy pool switching)
}

// NewCrawlerService 创建爬虫服务实例
// 初始化 Colly 实例并集成全套反爬策略, 支持异步爬取、深度限制、代理轮换等核心能力
// 参数:
//   - cfg: 全局配置(含爬虫UA/超时/代理池/存储路径等链式配置项) | Global config (crawler UA/timeout/proxy pool/storage path)
//
// 返回值:
//   - *CrawlerService: 爬虫服务实例 | Crawler service instance
func NewCrawlerService(cfg *config.SteamConfig) *CrawlerService {
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

	return &CrawlerService{
		cfg:          cfg,
		colly:        c,
		antiCrawl:    antiCrawl,
		parser:       parser,
		storage:      &crawler.Storage{},
		proxyRotator: proxyRotator,
	}
}

// GetProxyPool 获取当前代理池列表
// 暴露代理轮换器的代理池，支持外部监控/更新代理列表
// 返回值:
//   - []string: 代理地址列表(格式: http://ip:port) | Proxy address list (format: http://ip:port)
func (s *CrawlerService) GetProxyPool() []string {
	return s.proxyRotator.Pool()
}
