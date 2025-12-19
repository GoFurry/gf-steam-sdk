// Package main 是 Steam 玩家已拥有游戏查询示例程序
// 演示如何使用 gf-steam-sdk 调用 Steam API 获取指定玩家的游戏库信息
// Package main is an example program for querying Steam owned games
// Demonstrates how to use gf-steam-sdk to call Steam API to get game library info (playtime, last played time, etc.) of a specified player

package main

import (
	"fmt"
	"time"

	steamConfig "github.com/GoFurry/gf-steam-sdk/pkg/config"
	"github.com/GoFurry/gf-steam-sdk/pkg/steam"
	"github.com/GoFurry/gf-steam-sdk/pkg/util/log"
)

func init() {
	// 提供了简易的日志封装
	// Initialize logger (using the simple log wrapper built into the SDK)
	if err := log.InitLogger(nil); err != nil {
		panic(fmt.Sprintf("logger init failed: %v", err))
	}
}

// main 程序主入口
// 1. 初始化SDK配置（API Key、代理、超时等）
// 2. 创建SDK实例
// 3. 调用游戏库查询API
// 4. 解析并打印游戏库统计结果
// main is the program entry point
// 1. Initialize SDK configuration (API Key, proxy, timeout, etc.)
// 2. Create SDK instance
// 3. Call game library query API
// 4. Parse and print game library statistics results
func main() {
	log.Info("开始执行 Steam 游戏信息查询")

	// 初始化 sdk 的配置
	// Initialize SDK configuration
	cfg := steamConfig.NewDefaultConfig(). // 默认配置
						WithAPIKey("********B6E87843EF78948D********"). // Steam API Key
						WithProxyURL("http://127.0.0.1:7897").          // 代理IP, 中国地区需要添加
						WithProxyAuth("", "").                          // 代理认证信息, 账号, 密码
						WithTimeout(5*time.Second).                     // 请求超时时间
						WithRetryTimes(2).                              // 失败重试次数
						WithRateLimit(8.0, 15).                         // 限速
						WithHeaders(map[string]string{                  // 自定义请求头
			"Content-Type":    "application/json",
			"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			"Accept-Language": "zh-CN,zh",
		})

	// 创建Steam SDK实例
	// Create Steam SDK instance
	sdk, err := steam.NewSteamSDK(cfg)
	if err != nil {
		log.Fatal("[Main] 创建 Steam SDK 失败: %v", err)
	}

	// 调用API获取用户已拥有游戏信息
	steamID := "76561198370695025" // 用户ID
	log.Infof("开始查询 SteamID: %s 的游戏信息", steamID)

	// 调用API获取用户已拥有游戏信息
	// Call API to get user's owned games info
	games, err := sdk.Develop.GetOwnedGames(steamID, true)
	if err != nil {
		log.Fatalf("获取已拥有游戏失败: %v", err)
	}

	// 打印结果
	log.Infof("SteamID %s 共拥有 %d 款游戏", steamID, len(games))
	fmt.Printf("SteamID %s 共拥有 %d 款游戏: \n", steamID, len(games))

	// 打印前50款游戏详情
	// Print details of the first 50 games (avoid excessive output)
	limit := 50
	if len(games) < limit {
		limit = len(games)
	}
	for i := 0; i < limit; i++ {
		g := games[i]
		fmt.Printf("[%d]%s (ID: %d)\n", i+1, g.Name, g.AppID)
		fmt.Printf("  总游玩时长: %d分钟（%.1f小时）\n", g.PlaytimeForever, float64(g.PlaytimeForever)/60)
		fmt.Printf("  近2周游玩时长: %d分钟\n", g.Playtime2Weeks)
		fmt.Printf("  最后游玩时间: %s\n", g.LastPlayedTimeStr)
		fmt.Printf("  是否有DLC: %t\n", g.HasDLC)
		fmt.Printf("  SteamDeck游玩时长: %d分钟\n", g.PlaytimeDeckForever)
		fmt.Printf("  游戏图标: %s\n", g.IconURL)
		fmt.Printf("  游戏封面: %s\n", g.CapsuleURL)
		fmt.Println("------------------------")
	}

	if len(games) > limit {
		remaining := len(games) - limit
		log.Infof("还有 %d 款游戏未展示", remaining)
		fmt.Printf("... 还有 %d 款游戏未展示\n", remaining)
	}

	log.Info("Steam 游戏信息查询完成")
}
