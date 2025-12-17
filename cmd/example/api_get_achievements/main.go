// Package main 是 Steam 玩家成就查询示例程序
// 演示如何使用 gf-steam-sdk 调用 Steam API 获取指定玩家的游戏成就信息
// Package main is an example program for querying Steam player achievements
// Demonstrates how to use gf-steam-sdk to call Steam API to get game achievement info of a specified player

package main

import (
	"fmt"
	"log"
	"time"

	steamConfig "github.com/GoFurry/gf-steam-sdk/pkg/config"
	"github.com/GoFurry/gf-steam-sdk/pkg/steam"
	steamLog "github.com/GoFurry/gf-steam-sdk/pkg/util/log"
)

// 初始化日志
func init() {
	// 提供了简易的日志封装
	// Initialize logger (using the simple log wrapper built into the SDK)
	if err := steamLog.InitLogger(nil); err != nil {
		panic(fmt.Sprintf("日志初始化失败: %v", err))
	}
}

// main 程序主入口
// 1. 初始化SDK配置（API Key、代理、超时等）
// 2. 创建SDK实例
// 3. 调用成就查询API
// 4. 解析并打印成就统计结果
// main is the program entry point
// 1. Initialize SDK configuration (API Key, proxy, timeout, etc.)
// 2. Create SDK instance
// 3. Call achievement query API
// 4. Parse and print achievement statistics results
func main() {
	steamLog.Info("开始执行 Steam 玩家成就查询")

	// 初始化 sdk 的配置
	// Initialize SDK configuration
	cfg := steamConfig.NewDefaultConfig(). // 默认配置
						WithAPIKey("********B6E87843EF78948D********"). // Steam API Key (自行获取)
						WithProxyURL("http://127.0.0.1:7897").          // 代理IP, 中国地区需要添加
						WithProxyAuth("", "").                          // 代理认证信息, 账号, 密码
						WithTimeout(5*time.Second).                     // 请求超时时间
						WithRetryTimes(2).                              // 失败重试次数
						WithRateLimit(8.0, 15).                         // 限速
						WithHeaders(map[string]string{                  // 自定义请求头
			"Content-Type":    "application/json",
			"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			"Accept-Language": "zh-CN,zh",
		}).
		Debug() // 调试模式

	// 创建Steam SDK实例
	// Create Steam SDK instance
	sdk, err := steam.NewSteamSDK(cfg)
	if err != nil {
		log.Fatalf("[Main] 创建 Steam SDK 失败: %v", err)
	}

	// 调用成就API(示例: Left 4 Dead 2, AppID=550)
	// Call achievement query API (Example: Left 4 Dead 2, AppID=550)
	steamID := "76561198370695025" // 账户唯一ID
	appID := uint64(550)
	lang := "zh" // 中文返回成就名称/描述

	steamLog.Infof("查询 SteamID: %s 游戏 %d(%s) 的成就信息", steamID, appID, lang)
	achievements, err := sdk.Stats.GetPlayerAchievements(steamID, appID, lang)
	if err != nil {
		steamLog.Fatalf("获取成就信息失败: %v", err)
	}

	// 统计完成/未完成数量
	// Count completed/uncompleted achievements
	completed := 0
	uncompleted := 0
	for _, a := range achievements {
		if a.Achieved {
			completed++
		} else {
			uncompleted++
		}
	}

	// 打印结果
	fmt.Printf("=== SteamID: %s 游戏: %s (AppID: %d) ===\n", steamID, achievements[0].GameName, appID)
	fmt.Printf("总成就数: %d | 已完成: %d | 未完成: %d\n", len(achievements), completed, uncompleted)
	fmt.Println("----------------------------------------")

	// 打印前10项成就
	// Print details of the first 10 achievements
	limit := 10
	if len(achievements) < limit {
		limit = len(achievements)
	}
	for i := 0; i < limit; i++ {
		a := achievements[i]
		status := "已完成"
		if !a.Achieved {
			status = "未完成"
		}
		fmt.Printf("[%d] %s\n", i+1, status)
		fmt.Printf("  成就名称: %s\n", a.AchievementName)
		fmt.Printf("  成就描述: %s\n", a.Description)
		if a.Achieved {
			fmt.Printf("  解锁时间: %s\n", a.UnlockTimeStr)
		}
		fmt.Println("----------------------------------------")
	}

	if len(achievements) > limit {
		fmt.Printf("... 还有 %d 项成就未展示\n", len(achievements)-limit)
	}

	steamLog.Info("Steam 玩家成就查询执行完成")
}
