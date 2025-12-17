// Package crawler 实现 Steam 爬虫核心能力
// 包含反爬策略、代理轮换、HTML 解析和存储等功能
// Package crawler implements core capabilities of Steam crawler
// Includes anti-crawl strategy, proxy rotation, HTML parsing and storage

package crawler

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/GoFurry/gf-steam-sdk/pkg/config"
	"github.com/gocolly/colly"
	"golang.org/x/time/rate"
)

// AntiCrawl 反爬策略管理器
// 封装速率限制、随机请求头、延迟控制等反爬机制
// AntiCrawl is the anti-crawl strategy manager
// Encapsulates anti-crawl mechanisms such as rate limiting, random request headers, delay control
type AntiCrawl struct {
	cfg     *config.SteamConfig // 爬虫配置 | Crawler configuration
	limiter *rate.Limiter       // 速率限制器 | Rate limiter
	random  *rand.Rand          // 随机数生成器 | Random number generator
}

// NewAntiCrawl 创建反爬策略实例
// 参数:
//   - cfg: 爬虫配置实例 | Crawler configuration instance
//
// 返回值:
//   - *AntiCrawl: 反爬策略实例 | Anti-crawl strategy instance
func NewAntiCrawl(cfg *config.SteamConfig) *AntiCrawl {
	return &AntiCrawl{
		cfg:     cfg,
		limiter: rate.NewLimiter(rate.Limit(cfg.CrawlerQPS), cfg.CrawlerBurst),
		random:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Apply 应用反爬策略到 Colly 采集器
// 参数:
//   - c: Colly 采集器实例 | Colly collector instance
func (a *AntiCrawl) Apply(c *colly.Collector) {
	// 配置全局速率限制规则
	// Configure global rate limit rules
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: a.cfg.CrawlerConcurrency,
		Delay:       a.getRandomDelay(), // 随机请求延迟 | Random request delay
	})

	// 请求前钩子: 速率校验 + 随机请求头
	// Pre-request hook: rate check + random request headers
	c.OnRequest(func(r *colly.Request) {
		// 速率限制校验
		// Rate limit check
		if err := a.limiter.Wait(context.Background()); err != nil {
			r.Abort() // 触发限流则终止请求 | Abort request if rate limited
		}
		// 设置随机 Referer
		// Set random Referer
		r.Headers.Set("Referer", a.getRandomReferer())
		// 设置随机 Accept-Language(可取消注释启用)
		// Set random Accept-Language (uncomment to enable)
		//r.Headers.Set("Accept-Language", a.getRandomLang())
	})

	// 错误处理钩子: 5xx 错误自动重试
	// Error hook: auto retry on 5xx errors
	c.OnError(func(r *colly.Response, err error) {
		if r != nil && r.StatusCode >= 500 && r.StatusCode < 600 {
			r.Request.Retry()
		}
	})

	if a.cfg.IsDebug {
		fmt.Printf("[Info] Apply Anti-crawl Rules \n")
	}
}

// getRandomDelay 生成随机请求延迟
// 基于配置的基础延迟, 上下浮动 0-500ms, 避免固定延迟特征
// 返回值:
//   - time.Duration: 随机延迟时长 | Random delay duration
func (a *AntiCrawl) getRandomDelay() time.Duration {
	base := a.cfg.CrawlerDelay
	offset := time.Duration(a.random.Int63n(500)) * time.Millisecond
	if a.random.Intn(2) == 0 {
		return base + offset
	}
	return base - offset
}

// getRandomReferer 生成高仿真随机 Referer
// 模拟真实用户访问链路, 按权重返回不同场景的 Referer
// 返回值:
//   - string: 随机 Referer 字符串 | Random Referer string
func (a *AntiCrawl) getRandomReferer() string {
	// 定义 Referer 场景池
	// Define Referer scenario pool (weighted by real access ratio)
	scenarios := []struct {
		weight  int           // 权重 | Weight (higher = higher probability)
		genFunc func() string // Referer 生成函数 | Referer generation function
	}{
		// Steam 商店首页
		{
			weight: 40,
			genFunc: func() string {
				return "https://store.steampowered.com/" + a.randomQueryParams()
			},
		},
		// Google 搜索
		{
			weight: 20,
			genFunc: func() string {
				keywords := []string{
					"steam 游戏 " + a.randomGameKeyword(),
					"steam " + a.randomGameKeyword() + " buy",
					"steam " + a.randomGameKeyword() + " review",
					"best steam games",
					"steam 折扣 " + a.randomGameKeyword(),
				}
				keyword := keywords[a.random.Intn(len(keywords))]
				return "https://www.google.com/search?q=" + url.QueryEscape(keyword) + "&hl=en"
			},
		},
		// Steam 社区
		{
			weight: 15,
			genFunc: func() string {
				communityPaths := []string{
					"discussions/",
					"profiles/" + a.randomString(10) + "/",
					"groups/" + a.randomString(8) + "/",
					"id/" + a.randomString(12) + "/",
				}
				return "https://steamcommunity.com/" + communityPaths[a.random.Intn(len(communityPaths))] + a.randomQueryParams()
			},
		},
		// Steam 分类页
		{
			weight: 10,
			genFunc: func() string {
				categories := []string{
					"category/action",
					"category/adventure",
					"category/indie",
					"category/rpg",
					"category/strategy",
					"tags/" + a.randomString(6),
					"sale/random",
				}
				return "https://store.steampowered.com/" + categories[a.random.Intn(len(categories))] + "/" + a.randomQueryParams()
			},
		},
		// 空 Referer
		{
			weight: 10,
			genFunc: func() string {
				return ""
			},
		},
		// 其他 Steam 子域名
		{
			weight: 5,
			genFunc: func() string {
				subDomains := []string{
					"help.steampowered.com/",
					"steamcommunity.com/market/",
					"store.steampowered.com/cart/",
				}
				return "https://" + subDomains[a.random.Intn(len(subDomains))] + a.randomQueryParams()
			},
		},
	}

	// 按权重随机选择场景
	// Randomly select scenario by weight
	totalWeight := 0
	for _, s := range scenarios {
		totalWeight += s.weight
	}
	randomWeight := a.random.Intn(totalWeight)
	accumulated := 0
	for _, s := range scenarios {
		accumulated += s.weight
		if randomWeight < accumulated {
			return s.genFunc()
		}
	}

	// 兜底: 返回 Steam 商店首页
	// Fallback: return Steam store homepage
	return "https://store.steampowered.com/"
}

// randomGameKeyword 生成随机游戏相关关键词
// 增强 Referer 真实性, 覆盖热门游戏中英文名称
// 返回值:
//   - string: 随机游戏关键词 | Random game keyword
func (a *AntiCrawl) randomGameKeyword() string {
	gameKeywords := []string{
		"cs2", "dota2", "pubg", "elden ring", "cyberpunk 2077",
		"starfield", "hogwarts legacy", "resident evil 4", "god of war",
		"minecraft", "fortnite", "valorant", "apex legends", "rust",
		"求生之路2", "赛博朋克2077", "艾尔登法环", "原神", "星露谷物语",
	}
	return gameKeywords[a.random.Intn(len(gameKeywords))]
}

// randomQueryParams 生成随机 URL 查询参数
// 模拟真实请求特征, 70% 概率返回空参数
// 返回值:
//   - string: 随机查询参数 | Random query parameters
func (a *AntiCrawl) randomQueryParams() string {
	// 70% 概率返回空参数
	// 70% probability to return empty params
	if a.random.Intn(10) < 7 {
		return ""
	}

	params := []string{
		"?l=english",
		"?l=schinese",
		"?cc=US",
		"?cc=CN",
		"&p=" + strconv.Itoa(a.random.Intn(10)+1), // 分页参数 | Pagination param
		"&sort_by=Released_DESC",
		"&filter=topsellers",
	}
	return params[a.random.Intn(len(params))]
}

// getRandomLang 生成随机 Accept-Language 头
// 模拟不同地区用户的访问特征
// 返回值:
//   - string: 随机语言头字符串 | Random language header string
func (a *AntiCrawl) getRandomLang() string {
	langs := []string{
		"en-US,en;q=0.9",
		"zh-CN,zh;q=0.9,en;q=0.8",
		"en-GB,en;q=0.9",
	}
	return langs[a.random.Intn(len(langs))]
}

// randomString 生成指定长度的随机字符串
// 参数:
//   - n: 字符串长度 | String length
//
// 返回值:
//   - string: 随机字符串 | Random string
func (a *AntiCrawl) randomString(n int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[a.random.Intn(len(chars))]
	}
	return string(b)
}
