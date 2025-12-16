// Package crawler 实现 Steam 爬虫核心能力
// 包含反爬策略、代理轮换、HTML 解析和存储等功能
// Package crawler implements core capabilities of Steam crawler
// Includes anti-crawl strategy, proxy rotation, HTML parsing and storage

package crawler

import (
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	"github.com/GoFurry/gf-steam-sdk/pkg/config"
)

// ProxyRotator 代理轮换管理器
// 支持轮询/随机两种代理选择策略, 提供动态代理池管理能力
// ProxyRotator is the proxy rotation manager
// Supports round-robin/random proxy selection strategies and provides dynamic proxy pool management
type ProxyRotator struct {
	cfg        *config.SteamConfig // 全局配置 | Global configuration
	pool       []string            // 代理池 | Proxy pool
	strategy   string              // 选择策略(round_robin/random) | Selection strategy
	currentIdx uint32              // 轮询当前索引(原子操作) | Round-robin current index (atomic)
	random     *rand.Rand          // 随机数生成器 | Random number generator
	proxyUser  string              // 代理用户名 | Proxy username
	proxyPass  string              // 代理密码 | Proxy password
}

// NewProxyRotator 创建代理轮换实例
// 参数:
//   - cfg: 全局配置实例 | Global configuration instance
//
// 返回值:
//   - *ProxyRotator: 代理轮换实例 | Proxy rotation instance
func NewProxyRotator(cfg *config.SteamConfig) *ProxyRotator {
	// 构建代理池: 优先使用 ProxyPool, 无则使用 ProxyURL 作为单代理
	// Build proxy pool: use ProxyPool first, otherwise use ProxyURL as single proxy
	pool := cfg.ProxyPool
	if len(pool) == 0 && cfg.ProxyURL != "" {
		pool = []string{cfg.ProxyURL}
	}

	return &ProxyRotator{
		cfg:        cfg,
		pool:       pool,
		strategy:   cfg.ProxyStrategy,
		currentIdx: 0,
		random:     rand.New(rand.NewSource(time.Now().UnixNano())),
		proxyUser:  cfg.ProxyUser,
		proxyPass:  cfg.ProxyPass,
	}
}

// GetProxyFunc 返回 Colly 兼容的 ProxyFunc
// 每次请求调用时自动选择代理, 支持失败兜底
// 返回值:
//   - func(r *http.Request) (*url.URL, error): Colly 代理函数 | Colly proxy function
func (p *ProxyRotator) GetProxyFunc() func(r *http.Request) (*url.URL, error) {
	// 无代理池时返回空(不使用代理)
	// Return empty if no proxy pool (no proxy used)
	if len(p.pool) == 0 {
		return func(_ *http.Request) (*url.URL, error) {
			return nil, nil
		}
	}

	return func(_ *http.Request) (*url.URL, error) {
		// 选择代理地址
		// Select proxy address
		proxyAddr := p.selectProxy()
		if proxyAddr == "" {
			return nil, nil
		}

		// 解析代理 URL
		// Parse proxy URL
		proxyURL, err := url.Parse(proxyAddr)
		if err != nil {
			// 解析失败时切换下一个代理
			// Switch to next proxy if parse failed
			return p.fallbackProxy()
		}

		// 添加代理认证
		// Add proxy auth (if configured)
		if p.proxyUser != "" && p.proxyPass != "" {
			proxyURL.User = url.UserPassword(p.proxyUser, p.proxyPass)
		}

		return proxyURL, nil
	}
}

// selectProxy 根据策略选择代理
// 返回值:
//   - string: 选中的代理地址 | Selected proxy address
func (p *ProxyRotator) selectProxy() string {
	switch p.strategy {
	case "round_robin":
		return p.roundRobinSelect()
	case "random":
		return p.randomSelect()
	default:
		return p.roundRobinSelect() // 默认轮询 | Default to round-robin
	}
}

// roundRobinSelect 轮询选择代理
// 使用原子操作保证多协程下索引安全
// 返回值:
//   - string: 选中的代理地址 | Selected proxy address
func (p *ProxyRotator) roundRobinSelect() string {
	if len(p.pool) == 0 {
		return ""
	}

	// 原子操作: 获取当前索引并自增
	// Atomic operation: get current index and increment
	idx := atomic.AddUint32(&p.currentIdx, 1) - 1
	return p.pool[idx%uint32(len(p.pool))]
}

// randomSelect 随机选择代理
// 返回值:
//   - string: 选中的代理地址 | Selected proxy address
func (p *ProxyRotator) randomSelect() string {
	if len(p.pool) == 0 {
		return ""
	}
	return p.pool[p.random.Intn(len(p.pool))]
}

// fallbackProxy 代理解析失败时的兜底逻辑
// 重新选择代理并解析，避免单次失败导致请求终止
// 返回值:
//   - *url.URL: 兜底代理 URL | Fallback proxy URL
//   - error: 解析失败时返回错误 | Error if parse failed
func (p *ProxyRotator) fallbackProxy() (*url.URL, error) {
	proxyAddr := p.selectProxy()
	proxyURL, err := url.Parse(proxyAddr)
	if err != nil {
		return nil, err
	}
	if p.proxyUser != "" && p.proxyPass != "" {
		proxyURL.User = url.UserPassword(p.proxyUser, p.proxyPass)
	}
	return proxyURL, nil
}

// Pool 返回当前代理池
// 返回值:
//   - []string: 代理地址列表 | Proxy address list
func (p *ProxyRotator) Pool() []string {
	return p.pool
}

// AddProxy 动态添加代理
// 参数:
//   - proxyAddr: 代理地址 | Proxy address
func (p *ProxyRotator) AddProxy(proxyAddr string) {
	proxyAddr = strings.TrimSpace(proxyAddr)
	if proxyAddr == "" {
		return
	}
	// 避免重复添加
	// Avoid duplicate addition
	for _, p := range p.pool {
		if p == proxyAddr {
			return
		}
	}
	p.pool = append(p.pool, proxyAddr)
}

// RemoveProxy 动态移除代理
// 参数:
//   - proxyAddr: 代理地址 | Proxy address
func (p *ProxyRotator) RemoveProxy(proxyAddr string) {
	proxyAddr = strings.TrimSpace(proxyAddr)
	newPool := []string{}
	for _, p := range p.pool {
		if p != proxyAddr {
			newPool = append(newPool, p)
		}
	}
	p.pool = newPool
}
