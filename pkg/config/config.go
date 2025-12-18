// Package config 提供 Steam SDK 核心配置管理能力
// 包含默认配置初始化、链式配置修改、配置校验和 HTTP Transport 构建等功能
// Package config provides core configuration management capabilities for Steam SDK
// Includes default config initialization, chain configuration modification, config validation and HTTP Transport construction

package config

import (
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/GoFurry/gf-steam-sdk/pkg/util"
	"github.com/GoFurry/gf-steam-sdk/pkg/util/errors"
)

// SteamConfig Steam 客户端核心配置结构体
// 整合 API 基础配置、代理配置、速率限制和爬虫专属配置, 支持环境变量注入
// SteamConfig is the core configuration structure for Steam client
// Integrates API basic config, proxy config, rate limit and crawler-specific config, supports environment variable injection
type SteamConfig struct {
	// 基本配置 | Basic configuration
	APIKey         string            `json:"api_key" env:"STEAM_API_KEY"`                   // Steam API 密钥
	ProxyURL       string            `json:"proxy_url" env:"STEAM_PROXY_URL"`               // 代理地址(中国区必填)
	ProxyUser      string            `json:"proxy_user" env:"STEAM_PROXY_USER"`             // 代理用户名
	ProxyPass      string            `json:"proxy_pass" env:"STEAM_PROXY_PASS"`             // 代理密码
	ProxyPool      []string          `json:"proxy_pool" env:"STEAM_PROXY_POOL"`             // 代理IP池
	ProxyStrategy  string            `json:"proxy_strategy" env:"STEAM_PROXY_STRATEGY"`     // 代理选择策略
	Timeout        time.Duration     `json:"timeout" env:"STEAM_TIMEOUT"`                   // 请求超时时间(秒)
	RetryTimes     int               `json:"retry_times" env:"STEAM_RETRY_TIMES"`           // 重试次数
	RateLimitQPS   float64           `json:"rate_limit_qps" env:"STEAM_RATE_LIMIT_QPS"`     // 限速QPS
	RateLimitBurst int               `json:"rate_limit_burst" env:"STEAM_RATE_LIMIT_BURST"` // 突发QPS上限
	Headers        map[string]string `json:"headers"`                                       // 请求头
	IsDebug        bool              `json:"is_debug"`                                      // 调试模式
	Transport      *http.Transport   `json:"-"`                                             // 构建的 Transport | Built Transport

	// 爬虫配置 | Crawler configuration
	CrawlerUserAgent   string        `json:"crawler_user_agent" env:"STEAM_CRAWLER_UA"`           // 爬虫user-agent
	CrawlerAsync       bool          `json:"crawler_async" env:"STEAM_CRAWLER_ASYNC"`             // 异步爬虫
	CrawlerMaxDepth    int           `json:"crawler_max_depth" env:"STEAM_CRAWLER_MAX_DEPTH"`     // 爬取深度
	CrawlerConcurrency int           `json:"crawler_concurrency" env:"STEAM_CRAWLER_CONCURRENCY"` // 并发数
	CrawlerDelay       time.Duration `json:"crawler_delay" env:"STEAM_CRAWLER_DELAY"`             // 每次请求延迟(毫秒)
	CrawlerQPS         float64       `json:"crawler_qps" env:"STEAM_CRAWLER_QPS"`                 // 爬虫限速QPS
	CrawlerBurst       int           `json:"crawler_burst" env:"STEAM_CRAWLER_BURST"`             // 爬虫突发QPS上限
	CrawlerCookie      string        `json:"crawler_cookie" env:"STEAM_CRAWLER_COOKIE"`           // Steam 登录 Cookie
	CrawlerStorageDir  string        `json:"crawler_storage_dir" env:"STEAM_CRAWLER_STORAGE_DIR"` // HTML存储基础目录
}

// NewDefaultConfig 创建默认配置实例
// 优先从环境变量读取配置, 未配置则使用 util 包中定义的默认值
// 返回值:
//   - *SteamConfig: 初始化后的配置实例 | Initialized config instance
func NewDefaultConfig() *SteamConfig {
	// 解析超时时间(优先环境变量, 否则用默认值)
	// Parse timeout (env first, then default)
	timeoutSec := util.DEFAULT_TIMEOUT.Seconds()
	if envTimeout := os.Getenv("STEAM_TIMEOUT"); envTimeout != "" {
		if t, err := strconv.ParseFloat(envTimeout, 64); err == nil {
			timeoutSec = t
		}
	}

	// 解析重试次数
	// Parse retry count
	retryTimes := util.DEFAULT_RETRY_TIMES
	if envRetry := os.Getenv("STEAM_RETRY_TIMES"); envRetry != "" {
		if r, err := strconv.Atoi(envRetry); err == nil && r >= 0 {
			retryTimes = r
		}
	}

	// 解析API速率限制
	// Parse API rate limit
	rateLimitQPS := util.DEFAULT_RATE_QPS
	if envQPS := os.Getenv("STEAM_RATE_LIMIT_QPS"); envQPS != "" {
		if q, err := strconv.ParseFloat(envQPS, 64); err == nil && q > 0 {
			rateLimitQPS = q
		}
	}
	rateLimitBurst := util.DEFAULT_RATE_BURST
	if envBurst := os.Getenv("STEAM_RATE_LIMIT_BURST"); envBurst != "" {
		if b, err := strconv.Atoi(envBurst); err == nil && b > 0 {
			rateLimitBurst = b
		}
	}

	// 解析爬虫默认配置
	// Parse crawler default config
	crawlerUA := util.USER_AGENT
	if envUA := os.Getenv("STEAM_CRAWLER_UA"); envUA != "" {
		crawlerUA = envUA
	}
	crawlerAsync := false
	if envAsync := os.Getenv("STEAM_CRAWLER_ASYNC"); envAsync != "" {
		if a, err := strconv.ParseBool(envAsync); err == nil {
			crawlerAsync = a
		}
	}
	crawlerMaxDepth := util.CRAWLER_MAX_DEPTH
	if envDepth := os.Getenv("STEAM_CRAWLER_MAX_DEPTH"); envDepth != "" {
		if d, err := strconv.Atoi(envDepth); err == nil && d > 0 {
			crawlerMaxDepth = d
		}
	}
	crawlerConcurrency := util.CRAWLER_CONCURRENCY
	if envConcurrency := os.Getenv("STEAM_CRAWLER_CONCURRENCY"); envConcurrency != "" {
		if c, err := strconv.Atoi(envConcurrency); err == nil && c > 0 {
			crawlerConcurrency = c
		}
	}
	crawlerDelay := util.CRAWLER_DELAY
	if envDelay := os.Getenv("STEAM_CRAWLER_DELAY"); envDelay != "" {
		if d, err := strconv.ParseInt(envDelay, 10, 64); err == nil && d > 0 {
			crawlerDelay = time.Duration(d) * time.Millisecond
		}
	}
	crawlerQPS := util.CRAWLER_QPS
	if envQPS := os.Getenv("STEAM_CRAWLER_QPS"); envQPS != "" {
		if q, err := strconv.ParseFloat(envQPS, 64); err == nil && q > 0 {
			crawlerQPS = q
		}
	}
	crawlerBurst := util.CRAWLER_BURST
	if envBurst := os.Getenv("STEAM_CRAWLER_BURST"); envBurst != "" {
		if b, err := strconv.Atoi(envBurst); err == nil && b > 0 {
			crawlerBurst = b
		}
	}
	crawlerStorageDir := util.CRAWLER_STORAGE_DIR
	if envStorageDir := os.Getenv("STEAM_CRAWLER_STORAGE_DIR"); envStorageDir != "" {
		crawlerStorageDir = envStorageDir
	}

	// 解析代理认证配置
	// Parse proxy auth config
	proxyUser := os.Getenv("STEAM_PROXY_USER")
	proxyPass := os.Getenv("STEAM_PROXY_PASS")

	// 解析代理IP池(环境变量逗号分隔)
	// Parse proxy IP pool (env separated by commas)
	proxyPool := []string{}
	if envPool := os.Getenv("STEAM_PROXY_POOL"); envPool != "" {
		poolStr := strings.Split(envPool, ",")
		for _, p := range poolStr {
			p = strings.TrimSpace(p)
			if p != "" {
				proxyPool = append(proxyPool, p)
			}
		}
	}

	// 解析代理策略(默认轮询)
	// Parse proxy strategy (default round-robin)
	proxyStrategy := os.Getenv("STEAM_PROXY_STRATEGY")
	if proxyStrategy == "" {
		proxyStrategy = "round_robin"
	}

	// 构建配置实例
	// Build config instance
	cfg := &SteamConfig{
		// 基础配置 | Basic config
		APIKey:         os.Getenv("STEAM_API_KEY"),
		ProxyURL:       os.Getenv("STEAM_PROXY_URL"),
		ProxyUser:      proxyUser,
		ProxyPass:      proxyPass,
		ProxyPool:      proxyPool,
		ProxyStrategy:  proxyStrategy,
		Timeout:        time.Duration(timeoutSec) * time.Second,
		RetryTimes:     retryTimes,
		RateLimitQPS:   rateLimitQPS,
		RateLimitBurst: rateLimitBurst,
		IsDebug:        false,

		// 爬虫配置 | Crawler config
		CrawlerUserAgent:   crawlerUA,
		CrawlerAsync:       crawlerAsync,
		CrawlerMaxDepth:    crawlerMaxDepth,
		CrawlerConcurrency: crawlerConcurrency,
		CrawlerDelay:       crawlerDelay,
		CrawlerQPS:         crawlerQPS,
		CrawlerBurst:       crawlerBurst,
		CrawlerCookie:      os.Getenv("STEAM_CRAWLER_COOKIE"),
		CrawlerStorageDir:  crawlerStorageDir,
	}

	// 自动构建 HTTP Transport
	// Auto build HTTP Transport
	cfg.buildTransport()

	return cfg
}

// ============================ 基础链式配置 ============================

// Debug 调试模式
// 返回值:
//   - *SteamConfig: 配置实例(支持链式调用) | Config instance (chain call supported)
func (c *SteamConfig) Debug() *SteamConfig {
	c.IsDebug = true
	return c
}

// WithAPIKey 自定义API Key
// 参数:
//   - apiKey: Steam API 密钥 | Steam API key
//
// 返回值:
//   - *SteamConfig: 配置实例(支持链式调用) | Config instance (chain call supported)
func (c *SteamConfig) WithAPIKey(apiKey string) *SteamConfig {
	c.APIKey = apiKey
	return c
}

// WithProxyURL 自定义代理地址
// 修改后自动重建 Transport
// 参数:
//   - proxyURL: 代理地址 | Proxy URL
//
// 返回值:
//   - *SteamConfig: 配置实例 | Config instance
func (c *SteamConfig) WithProxyURL(proxyURL string) *SteamConfig {
	c.ProxyURL = proxyURL
	c.buildTransport()
	return c
}

// WithTimeout 自定义超时时间
// 参数:
//   - timeout: 超时时长 | Timeout duration
//
// 返回值:
//   - *SteamConfig: 配置实例 | Config instance
func (c *SteamConfig) WithTimeout(timeout time.Duration) *SteamConfig {
	c.Timeout = timeout
	return c
}

// WithRetryTimes 自定义重试次数
// 参数:
//   - retryTimes: 重试次数(仅接受 >=0 的值) | Retry count (only >=0 accepted)
//
// 返回值:
//   - *SteamConfig: 配置实例 | Config instance
func (c *SteamConfig) WithRetryTimes(retryTimes int) *SteamConfig {
	if retryTimes >= 0 {
		c.RetryTimes = retryTimes
	}
	return c
}

// WithHeaders 自定义请求头
// 参数:
//   - headers: 请求头键值对 | Request header key-value pairs
//
// 返回值:
//   - *SteamConfig: 配置实例 | Config instance
func (c *SteamConfig) WithHeaders(headers map[string]string) *SteamConfig {
	c.Headers = headers
	return c
}

// WithProxyAuth 设置代理认证信息
// 修改后自动重建 Transport
// 参数:
//   - user: 代理用户名 | Proxy username
//   - pass: 代理密码 | Proxy password
//
// 返回值:
//   - *SteamConfig: 配置实例 | Config instance
func (c *SteamConfig) WithProxyAuth(user, pass string) *SteamConfig {
	c.ProxyUser = user
	c.ProxyPass = pass
	c.buildTransport() // 重新构建 Transport | Rebuild Transport
	return c
}

// WithRateLimit 自定义API限流
// 参数:
//   - qps: 每秒请求数 | Requests per second
//   - burst: 突发请求上限 | Burst request limit
//
// 返回值:
//   - *SteamConfig: 配置实例 | Config instance
func (c *SteamConfig) WithRateLimit(qps float64, burst int) *SteamConfig {
	c.RateLimitQPS, c.RateLimitBurst = qps, burst
	return c
}

// WithProxyPool 设置代理IP池
// 自动清理空值和空格
// 参数:
//   - proxyPool: 代理地址列表 | Proxy address list
//
// 返回值:
//   - *SteamConfig: 配置实例 | Config instance
func (c *SteamConfig) WithProxyPool(proxyPool []string) *SteamConfig {
	cleanPool := []string{}
	for _, p := range proxyPool {
		p = strings.TrimSpace(p)
		if p != "" {
			cleanPool = append(cleanPool, p)
		}
	}
	c.ProxyPool = cleanPool
	return c
}

// WithProxyStrategy 设置代理选择策略
// 仅支持 round_robin/random，非法值默认使用 round_robin
// 参数:
//   - strategy: 代理策略 | Proxy strategy
//
// 返回值:
//   - *SteamConfig: 配置实例 | Config instance
func (c *SteamConfig) WithProxyStrategy(strategy string) *SteamConfig {
	if strategy == "round_robin" || strategy == "random" {
		c.ProxyStrategy = strategy
	} else {
		c.ProxyStrategy = "round_robin" // 非法策略默认轮询 | Default to round-robin for invalid strategy
	}
	return c
}

// ============================ 爬虫链式配置 ============================

// WithCrawlerUA 自定义爬虫UA
// 参数:
//   - ua: User-Agent 字符串 | User-Agent string
//
// 返回值:
//   - *SteamConfig: 配置实例 | Config instance
func (c *SteamConfig) WithCrawlerUA(ua string) *SteamConfig {
	if ua != "" {
		c.CrawlerUserAgent = ua
	}
	return c
}

// WithCrawlerAsync 自定义爬虫是否异步
// 参数:
//   - async: 是否异步 | Whether to enable async
//
// 返回值:
//   - *SteamConfig: 配置实例 | Config instance
func (c *SteamConfig) WithCrawlerAsync(async bool) *SteamConfig {
	c.CrawlerAsync = async
	return c
}

// WithCrawlerMaxDepth 自定义爬虫最大深度
// 参数:
//   - depth: 爬取深度(仅接受 >0 的值) | Crawl depth (only >0 accepted)
//
// 返回值:
//   - *SteamConfig: 配置实例 | Config instance
func (c *SteamConfig) WithCrawlerMaxDepth(depth int) *SteamConfig {
	if depth > 0 {
		c.CrawlerMaxDepth = depth
	}
	return c
}

// WithCrawlerConcurrency 自定义爬虫并发数
// 参数:
//   - concurrency: 并发数(仅接受 >0 的值) | Concurrency count (only >0 accepted)
//
// 返回值:
//   - *SteamConfig: 配置实例 | Config instance
func (c *SteamConfig) WithCrawlerConcurrency(concurrency int) *SteamConfig {
	if concurrency > 0 {
		c.CrawlerConcurrency = concurrency
	}
	return c
}

// WithCrawlerDelay 自定义爬虫请求延迟
// 参数:
//   - delay: 请求延迟时长 | Request delay duration
//
// 返回值:
//   - *SteamConfig: 配置实例 | Config instance
func (c *SteamConfig) WithCrawlerDelay(delay time.Duration) *SteamConfig {
	if delay >= 0 {
		c.CrawlerDelay = delay
	}
	return c
}

// WithCrawlerRateLimit 自定义爬虫速率限制
// 参数:
//   - qps: 爬虫每秒请求数 | Crawler requests per second
//   - burst: 爬虫突发请求上限 | Crawler burst request limit
//
// 返回值:
//   - *SteamConfig: 配置实例 | Config instance
func (c *SteamConfig) WithCrawlerRateLimit(qps float64, burst int) *SteamConfig {
	if qps > 0 {
		c.CrawlerQPS = qps
	}
	if burst > 0 {
		c.CrawlerBurst = burst
	}
	return c
}

// WithCrawlerCookie 自定义爬虫Cookie
// 参数:
//   - cookie: Steam 登录Cookie | Steam login cookie
//
// 返回值:
//   - *SteamConfig: 配置实例 | Config instance
func (c *SteamConfig) WithCrawlerCookie(cookie string) *SteamConfig {
	c.CrawlerCookie = cookie
	return c
}

// WithCrawlerStorageDir 自定义爬虫HTML存储目录
// 参数:
//   - dir: 存储目录路径 | Storage directory path
//
// 返回值:
//   - *SteamConfig: 配置实例 | Config instance
func (c *SteamConfig) WithCrawlerStorageDir(dir string) *SteamConfig {
	if dir != "" {
		c.CrawlerStorageDir = dir
	}
	return c
}

// ============================ 工具方法 ============================

// Validate 校验配置合法性
// 检查API Key、超时时间、爬虫相关配置的合法性
// 返回值:
//   - error: 配置非法时返回错误, 合法则返回nil | Error if config invalid, nil if valid
func (c *SteamConfig) Validate() error {
	if c.APIKey == "" {
		return errors.ErrMissingAPIKey
	}
	if c.Timeout <= 0 {
		return errors.New("timeout must be greater than 0")
	}
	if c.CrawlerMaxDepth < 0 {
		return errors.New("crawler max depth must be >= 0")
	}
	if c.CrawlerConcurrency < 0 {
		return errors.New("crawler concurrency must be >= 0")
	}
	if c.CrawlerDelay < 0 {
		return errors.New("crawler delay must be >= 0")
	}
	if c.CrawlerQPS < 0 {
		return errors.New("crawler qps must be >= 0")
	}
	if c.CrawlerBurst < 0 {
		return errors.New("crawler burst must be >= 0")
	}
	return nil
}

// buildTransport 构建 HTTP Transport
// 根据代理配置自动构建带/不带代理的 Transport 实例
func (c *SteamConfig) buildTransport() {
	if c.ProxyURL == "" {
		c.Transport = &http.Transport{
			MaxIdleConns:       100,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: false,
		}
		return
	}

	// 解析代理URL | Parse proxy URL
	proxyURL, err := url.Parse(c.ProxyURL)
	if err != nil {
		c.Transport = &http.Transport{}
		return
	}

	// 代理认证 | Proxy auth
	if c.ProxyUser != "" && c.ProxyPass != "" {
		proxyURL.User = url.UserPassword(c.ProxyUser, c.ProxyPass)
	}

	// 构建带代理的 Transport | Build Transport with proxy
	c.Transport = &http.Transport{
		Proxy:               http.ProxyURL(proxyURL),
		MaxIdleConns:        100,
		IdleConnTimeout:     30 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		DisableCompression:  false,
	}
}
