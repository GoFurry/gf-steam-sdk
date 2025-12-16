// Package main 是 Steam 爬虫功能示例程序
// 演示如何使用 gf-steam-sdk 的爬虫模块获取游戏详情页原始HTML, 并支持本地保存、代理轮换等高级功能
// Package main is an example program for Steam crawler functionality
// Demonstrates how to use the crawler module of gf-steam-sdk to get raw HTML of game detail pages,
// and supports advanced features such as local saving and proxy rotation

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/GoFurry/gf-steam-sdk/pkg/config"
	"github.com/GoFurry/gf-steam-sdk/pkg/steam"
	"github.com/GoFurry/gf-steam-sdk/pkg/util/log"
)

func init() {
	// 日志配置(开发模式)
	// Log configuration (development mode)
	logCfg := log.Config{
		Level:    "debug", // 日志级别 | Log level
		Mode:     "dev",   // 输出模式(开发/生产) | Output mode (dev/prod)
		ShowLine: true,    // 显示代码行号 | Show code line numbers
	}
	if err := log.InitLogger(&logCfg); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}
}

// main 程序主入口
// 1. 解析命令行参数(游戏ID、代理、保存开关等)
// 2. 初始化爬虫配置(代理轮换、反爬策略、存储路径等)
// 3. 调用爬虫API获取游戏详情页原始HTML
// 4. 可选保存HTML到本地, 并预览内容
// main is the program entry point
// 1. Parse command line arguments (game ID, proxy, save switch, etc.)
// 2. Initialize crawler configuration (proxy rotation, anti-crawl strategy, storage path, etc.)
// 3. Call crawler API to get raw HTML of game detail page
// 4. Optionally save HTML to local and preview content
func main() {
	// 解析命令行参数(示例: go run main.go -appid 730 -proxy http://127.0.0.1:7897)
	// Parse command line arguments (Example: go run main.go -appid 730 -proxy http://127.0.0.1:7897)
	var (
		appID     uint64 // 游戏ID | Game ID
		proxyURL  string // 代理地址 | Proxy URL
		proxyUser string // 代理用户名 | Proxy username
		proxyPass string // 代理密码 | Proxy password
		saveHTML  bool   // 是否保存HTML | Whether to save HTML
	)
	flag.Uint64Var(&appID, "appid", 550, "游戏ID(示例: 550=Left 4 Dead 2, 730=CS2)")
	flag.StringVar(&proxyURL, "proxy", "", "代理地址(如 http://127.0.0.1:7897)")
	flag.StringVar(&proxyUser, "proxy-user", "", "代理用户名(可选)")
	flag.StringVar(&proxyPass, "proxy-pass", "", "代理密码(可选)")
	flag.BoolVar(&saveHTML, "save", true, "是否保存HTML到本地文件")
	flag.Parse()

	// 初始化SDK配置(适配爬虫场景)
	// Initialize SDK configuration (adapted for crawler scenarios)
	cfg := config.NewDefaultConfig().
		WithTimeout(30*time.Second).           // 爬虫超时 | Crawler timeout (longer recommended)
		WithRetryTimes(1).                     // 失败重试次数 | Retry count on failure
		WithProxyURL("http://127.0.0.1:7897"). // 默认代理地址 | Default proxy URL
		WithProxyStrategy("round_robin").      // 代理轮换策略(轮询/随机) | Proxy rotation strategy (round-robin/random)
		//WithProxyPool([]string{"http://127.0.0.1:50000", "http://127.0.0.1:50001"}). // 代理池(可选) | Proxy pool (optional)
		WithProxyAuth("", "").                                                                                                            // 代理认证(全局) | Global proxy authentication
		WithCrawlerUA("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"). // 爬虫UA | Crawler User-Agent
		WithCrawlerAsync(false).                                                                                                          // 调试关闭异步 | Disable async for debugging
		WithCrawlerMaxDepth(1).                                                                                                           // 爬取深度 | Crawl depth (current page only)
		WithCrawlerConcurrency(1).                                                                                                        // 并发数 | Concurrency (set to 1 for debugging)
		WithCrawlerDelay(1*time.Second).                                                                                                  // 爬取延迟 | Crawl delay (anti-crawl)
		WithCrawlerRateLimit(5.0, 10).                                                                                                    // 爬虫速率限制 | Crawler rate limit
		WithCrawlerStorageDir("./custom-storage/html").                                                                                   // HTML存储路径 | HTML storage path
		WithHeaders(map[string]string{                                                                                                    // 自定义请求头 | Custom request headers
			"Content-Type":    "application/json",
			"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			"Accept-Language": "zh-CN,zh",
		})

	// 配置代理(中国地址访问需要)
	if proxyURL != "" {
		cfg.WithProxyURL(proxyURL)
		// 代理认证
		if proxyUser != "" && proxyPass != "" {
			cfg.WithProxyAuth(proxyUser, proxyPass)
		}
	}

	// 创建Steam SDK实例
	// Create Steam SDK instance
	sdk, err := steam.NewSteamSDK(cfg)
	if err != nil {
		log.Fatalf("创建 Steam SDK 失败: %v", err)
	}

	// 调用爬虫API获取游戏详情页原始HTML
	// Call crawler API to get raw HTML of game detail page
	log.Infof("开始爬取游戏详情页 | appID: %d", appID)
	htmlBytes, err := sdk.Crawler.GetGameStoreRawHTML(appID)
	if err != nil {
		log.Fatalf("爬取游戏详情页失败 | appID: %d, 错误: %v", appID, err)
	}

	// 调试输出结果
	log.Infof("爬取游戏详情页成功 | appID: %d, HTML大小: %d 字节", appID, len(htmlBytes))

	// 保存HTML到本地
	// Save HTML to local
	savePath, err := sdk.Crawler.SaveGameStoreRawHTML(550, "")
	if err != nil {
		log.Error("保存HTML失败", err)
	} else {
		// 输出示例: ./custom-storage/html/20251215/store.steampowered.com_app_550.html
		fmt.Println("HTML保存路径：", savePath)
	}

	// 预览HTML前500字符
	// Preview first 500 characters of HTML (debugging)
	previewLen := 500
	if len(htmlBytes) < previewLen {
		previewLen = len(htmlBytes)
	}
	fmt.Printf("\n=== HTML 内容预览(前%d字符)===\n%s\n", previewLen, string(htmlBytes[:previewLen]))
}
