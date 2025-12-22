// Package util 提供 Steam SDK 通用工具函数和常量
// 包含类型转换、时间处理、字符串操作和默认配置常量等
// Package util provides common utility functions and constants for Steam SDK
// Includes type conversion, time processing, string operations and default config constants
package util

import "time"

// 时间常量 | Time constants
const (
	TIME_FORMAT_DIGIT = "20060102150405"      // 数字格式时间戳 | Numeric format timestamp
	TIME_FORMAT_DATE  = "2006-01-02 15:04:05" // 标准日期时间格式 | Standard datetime format
	TIME_FORMAT_DAY   = "2006-01-02"          // 日期格式 | Date format
	TIME_FORMAT       = "20060102"            // 简化日期格式(用于目录/文件名) | Simplified date format (for dir/filename)
)

// 请求头常量 | Request header constants
const (
	USER_AGENT      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36" // 默认User-Agent
	ACCEPT_LANGUAGE = "zh-CN,zh"                                                                                                        // 默认接受语言 | Default accept language
	APPLICATION     = "application/json"                                                                                                // JSON内容类型 | JSON content type
)

// Steam API 常量 | Steam API constants
const (
	STEAM_STORE_BASE_URL = "https://store.steampowered.com/"
	STEAM_API_BASE_URL   = "https://api.steampowered.com/"                                              // Steam API基础地址 | Steam API base URL
	STEAM_ICON_URL       = "https://media.steampowered.com/steamcommunity/public/images/apps/%d/%s.jpg" // 游戏图标URL模板 | Game icon URL template
	STEAM_CAPSULE_URL    = "https://cdn.akamai.steamstatic.com/steam/apps/%d/header.jpg"                // 游戏封面URL模板 | Game capsule URL template
)

// 基础默认配置 | Basic default config
const (
	DEFAULT_TIMEOUT     = 5 * time.Second // 默认请求超时 | Default request timeout
	DEFAULT_RATE_QPS    = 10.0            // 默认API限速QPS | Default API rate limit QPS
	DEFAULT_RATE_BURST  = 20              // 默认API突发请求上限 | Default API burst limit
	DEFAULT_RETRY_TIMES = 2               // 默认重试次数 | Default retry count
	RETRY_SLEEP_BASE    = 300             // 重试基础延迟(毫秒) | Retry base delay (milliseconds)
)

// 爬虫默认配置 | Crawler default config
const (
	CRAWLER_MAX_DEPTH   = 1                        // 默认爬虫深度 | Default crawler depth
	CRAWLER_CONCURRENCY = 1                        // 默认爬虫并发数 | Default crawler concurrency
	CRAWLER_DELAY       = 1000 * time.Millisecond  // 默认爬虫请求延迟 | Default crawler request delay
	CRAWLER_QPS         = 5.0                      // 默认爬虫限速QPS | Default crawler rate limit QPS
	CRAWLER_BURST       = 10                       // 默认爬虫突发请求上限 | Default crawler burst limit
	CRAWLER_STORAGE_DIR = "./storage/crawler/html" // 默认爬虫HTML存储目录 | Default crawler HTML storage dir
)
