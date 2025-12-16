# GF-Steam-SDK

[![Last Version](https://img.shields.io/github/release/GoFurry/gf-steam-sdk/all.svg?logo=github&color=brightgreen)](https://github.com/GoFurry/gf-steam-sdk/releases)
[![License](https://img.shields.io/github/license/GoFurry/gf-steam-sdk)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.24-blue)](go.mod)

A lightweight, modular Go SDK for the Steam Open Platform, providing Steam WebAPI encapsulation and web crawling capabilities with built-in intelligent anti-crawling strategies.

一款轻量级、模块化的 Steam 开放平台 Go SDK, 提供Steam WebAPI 封装与网页爬虫能力, 内置智能反爬策略。

---

## 🌟 Core Features | 核心特性 🌟

### 1. 模块化架构设计 | Modular Architecture
- 拆分 **Player/Game/Stats/Crawler/Server** 五大核心模块, 职责清晰, 可按需使用
- 统一入口 `SteamSDK` 管理, 支持按需初始化, 降低资源占用

### 2. 灵活的链式配置 | Flexible Chain Configuration
- 支持 API Key、代理池、请求超时、重试次数等参数链式配置
- 示例: `config.NewSteamConfig().WithAPIKey("xxx").WithProxyPool(proxies).WithTimeout(10*time.Second)`

### 3. 多层级数据返回 | Multi-level Data Response
- 原始字节流(`RawBytes`): 保留 API 原始响应, 适用于自定义解析
- 结构化原始模型(`RawModel`): 映射 Steam 官方响应结构, 保留全量字段
- 精简业务模型(`Brief`): 剔除冗余字段, 补充格式化时间、布尔状态等易用性字段

### 4. 智能反爬策略 | Intelligent Anti-Crawling
- 动态代理轮换: 自动切换代理池, 规避 IP 封禁
- 随机 User-Agent + 合法 Referer 头: 模拟真实浏览器请求
- 请求延迟控制 + QPS 限流: 适配 Steam 风控规则
- 异步爬取 + 最大深度限制: 提升效率同时防止无限递归

### 5. 完整功能覆盖 | Comprehensive Features
| 模块      | 核心能力                              | 接口示例                                                                                                               |
|---------|-----------------------------------|--------------------------------------------------------------------------------------------------------------------|
| Player  | 玩家信息查询(批量支持)、在线状态检测               | `GetPlayerSummaries("76561198000000000")`                                                                          |
| Game    | 已拥有游戏查询、游戏详情、多平台时长统计              | `GetOwnedGames("76561198000000000", true)`                                                                         |
| Stats   | 游戏成就查询、解锁时间统计                     | `GetPlayerAchievements("7656...", 550, "zh")`                                                                      |
| Crawler | Steam 商店页爬取、HTML 存储、自定义爬取         | `GetGameStoreRawHTML(550)`<br/>`SaveGameStoreRawHTML(550, "/storage/")`                                            |
| Server  | A2S 服务器信息查询(基础/玩家/规则)、批量限流重试 | `GetServerDetail("110.42.54.147:52023")`<br/>`GetServerDetailList([]string{"ip:port"}, 2.0, 5, 30*time.Second, 3)` |

### 6. 高可用性设计 | High Availability
- 完善的错误体系: 自定义错误类型(参数错误/API 错误/爬虫错误), 便于问题定位
- 自动重试机制: 网络波动时自动重试请求, 提升稳定性
- 参数校验 + 兜底逻辑: 避免空值、非法参数导致的崩溃

---

## 🚀 Quick Start | 快速上手

### 1. Installation | 安装
```bash
go get github.com/GoFurry/gf-steam-sdk@latest
```

### 2. Basic Usage | 基础使用
#### 初始化 SDK | Initialize SDK
```go
package main

import (
	"fmt"
	"time"

	"github.com/GoFurry/gf-steam-sdk/pkg/config"
	"github.com/GoFurry/gf-steam-sdk/pkg/steam"
)

func main() {
	// 1. 配置初始化(支持链式配置)
	cfg := config.NewSteamConfig()
	.WithAPIKey("your-steam-api-key")       // 替换为你的 Steam API Key
	.WithProxyPool([]string{                // 可选: 配置代理池
		"http://127.0.0.1:7890",
		"http://127.0.0.1:7891",
	})
	.WithTimeout(10 * time.Second)          // 请求超时
	.WithRetryCount(3)                      // 重试次数
	.WithCrawlerAsync(true)                 // 爬虫异步模式
	.WithCrawlerMaxDepth(2)                 // 爬虫最大深度

	// 2. 创建 SDK 实例
	sdk, err := steam.NewSteamSDK(cfg)
	if err != nil {
		panic(fmt.Sprintf("initialize SDK failed: %v", err))
	}

	// 3. 调用接口(示例: 查询玩家信息)
	players, err := sdk.Player.GetPlayerSummaries("76561198000000000")
	if err != nil {
		panic(fmt.Sprintf("get player info failed: %v", err))
	}

	// 4. 处理结果
	for _, p := range players {
		fmt.Printf("SteamID: %s\nName: %s\nOnline: %t\nAvatar: %s\n",
			p.SteamID, p.PersonaName, p.IsOnline, p.AvatarFull)
	}
}
```

### 3. Advanced Examples | 进阶示例
#### 爬取 Steam 游戏商店页 | Crawl Steam Store Page
```go
// 爬取《Left 4 Dead 2》商店页(AppID: 550)
html, err := sdk.Crawler.GetGameStoreRawHTML(550)
if err != nil {
    panic(err)
}

// 保存 HTML 到本地(自动生成文件名: store.steampowered.com_app_550.html)
savePath, err := sdk.Crawler.SaveGameStoreRawHTML(550, "")
if err != nil {
    panic(err)
}
fmt.Printf("HTML saved to: %s\n", savePath)
```
#### 查询玩家已拥有游戏 | Get Player Owned Games
```go
// 查询玩家已拥有游戏(包含免费游戏)
games, err := sdk.Game.GetOwnedGames("76561198000000000", true)
if err != nil {
panic(err)
}

for _, game := range games {
fmt.Printf("Game: %s (AppID: %d)\nPlaytime: %d mins\nLast Played: %s\n",
game.Name, game.AppID, game.PlaytimeForever, game.LastPlayedTimeStr)
}
```
#### 查询游戏成就 | Get Game Achievements
```go
// 查询玩家在《Left 4 Dead 2》中的成就(中文)
achievements, err := sdk.Stats.GetPlayerAchievements("76561198000000000", 550, "zh")
if err != nil {
panic(err)
}

for _, a := range achievements {
fmt.Printf("Achievement: %s\nDescription: %s\nUnlocked: %t\nTime: %s\n",
a.AchievementName, a.Description, a.Achieved, a.UnlockTimeStr)
}
```
#### 查询游戏服务器信息 | Get Game Server Details
```go
// 调用聚合接口获取完整信息
detail, err := sdk.Server.GetServerDetail("110.42.54.147:52021")
if err != nil {
fmt.Printf("查询失败: %v\n", err)
return
}

// 打印查询结果
fmt.Printf("Server address: %s\n", addr)
fmt.Printf("Server info: %+v\n", detail.Server)
fmt.Printf("Player info: %+v\n", detail.Player)
fmt.Printf("Rules info: %+v\n", detail.Rules)
```
## 📋 Configuration Options | 配置项说明

| 配置项                | 类型         | 说明                                  | 默认值                  |
|-----------------------|--------------|---------------------------------------|-------------------------|
| APIKey                | string       | Steam API Key(从[Steam 开发者平台](https://steamcommunity.com/dev/apikey)获取) | 环境变量`STEAM_API_KEY`，无则为"dummy-key" |
| ProxyURL              | string       | 代理地址(中国区访问Steam必填，格式：http://ip:port) | 环境变量`STEAM_PROXY_URL`，无则为空 |
| ProxyUser             | string       | 代理认证用户名 | 环境变量`STEAM_PROXY_USER`，无则为空 |
| ProxyPass             | string       | 代理认证密码 | 环境变量`STEAM_PROXY_PASS`，无则为空 |
| ProxyPool             | []string     | 代理IP池(环境变量以逗号分隔，自动过滤空值) | 环境变量`STEAM_PROXY_POOL`，无则为空数组 |
| ProxyStrategy         | string       | 代理选择策略(仅支持 round_robin/random) | "round_robin" |
| Timeout               | time.Duration| 请求超时时间(秒) | 环境变量`STEAM_TIMEOUT`，无则为5 * time.Second |
| RetryTimes            | int          | 请求重试次数(仅接受>=0的值) | 环境变量`STEAM_RETRY_TIMES`，无则为2 |
| RateLimitQPS          | float64      | API接口限速QPS(每秒请求数) | 环境变量`STEAM_RATE_LIMIT_QPS`，无则为10.0 |
| RateLimitBurst        | int          | API接口突发QPS上限 | 环境变量`STEAM_RATE_LIMIT_BURST`，无则为20 |
| Headers               | map[string]string | 全局请求头自定义键值对 | nil |
| CrawlerUserAgent      | string       | 爬虫默认 User-Agent | 环境变量`STEAM_CRAWLER_UA`，无则为"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36" |
| CrawlerAsync          | bool         | 爬虫是否启用异步模式 | 环境变量`STEAM_CRAWLER_ASYNC`，无则为false |
| CrawlerMaxDepth       | int          | 爬虫最大爬取深度 | 环境变量`STEAM_CRAWLER_MAX_DEPTH`，无则为1 |
| CrawlerConcurrency    | int          | 爬虫并发数 | 环境变量`STEAM_CRAWLER_CONCURRENCY`，无则为5 |
| CrawlerDelay          | time.Duration| 爬虫每次请求延迟(毫秒) | 环境变量`STEAM_CRAWLER_DELAY`，无则为500 * time.Millisecond |
| CrawlerQPS            | float64      | 爬虫限速QPS | 环境变量`STEAM_CRAWLER_QPS`，无则为5.0 |
| CrawlerBurst          | int          | 爬虫突发QPS上限 | 环境变量`STEAM_CRAWLER_BURST`，无则为10 |
| CrawlerCookie         | string       | Steam登录Cookie(用于爬取需登录的内容) | 环境变量`STEAM_CRAWLER_COOKIE`，无则为空 |
| CrawlerStorageDir     | string       | 爬虫HTML存储基础目录 | 环境变量`STEAM_CRAWLER_STORAGE_DIR`，无则为"./steam-crawl-data" |

## 📚 Documentation References | 文档参考
- [Steam Web API 官方文档](https://developer.valvesoftware.com/wiki/Steam_Web_API)
- [Steam 合作伙伴 API 文档](https://partner.steamgames.com/doc/webapi_overview)
- [Steam API 非官方参考](https://steamapi.xpaw.me/)

---

## ⚠️ Notes | 注意事项
1. **API Key 申请**: 部分接口(如玩家成就、已拥有游戏)需要有效的 Steam API Key, 建议从 [Steam 开发者平台](https://steamcommunity.com/dev/apikey) 申请
2. **速率限制**: Steam API 有 QPS 限制, 建议通过 `WithQPSLimit` 配置限流, 避免账号封禁
3. **代理使用**: 爬取 Steam 商店页时建议配置代理池, 否则可能导致 IP 被封禁
4. **未完成功能**: OpenID 鉴权 API 封装正在开发中, 敬请期待

---

## 🛠️ Contributing | 贡献指南
1. Fork 本仓库
2. 创建特性分支(`git checkout -b feature/xxx`)
3. 提交代码(`git commit -m "feat: add xxx feature"`)
4. 推送分支(`git push origin feature/xxx`)
5. 提交 Pull Request

---

## 📄 License | 许可证
本项目基于 [MIT License](LICENSE) 开源, 允许商业使用、修改、分发, 无需保留原作者版权声明。

---

## 📞 Contact | 联系作者
- GitHub: [@GoFurry](https://github.com/GoFurry)
- 项目地址: [https://github.com/GoFurry/gf-steam-sdk](https://github.com/GoFurry/gf-steam-sdk)
