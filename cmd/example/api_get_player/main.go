// Package main 是 Steam 玩家基本信息查询示例程序
// 演示如何使用 gf-steam-sdk 调用 Steam API 批量获取多个玩家的基础信息
// Package main is an example program for querying Steam player basic info
// Demonstrates how to use gf-steam-sdk to call Steam API to batch get basic info (nickname, avatar, online status, etc.) of multiple players

package main

import (
	"fmt"
	"time"

	steamConfig "github.com/GoFurry/gf-steam-sdk/pkg/config"
	"github.com/GoFurry/gf-steam-sdk/pkg/steam"
	"github.com/GoFurry/gf-steam-sdk/pkg/util/log"
)

// 初始化日志
func init() {
	// 提供了简易的日志封装
	// Initialize logger (using the simple log wrapper built into the SDK)
	if err := log.InitLogger(nil); err != nil {
		panic(fmt.Sprintf("日志初始化失败: %v", err))
	}
}

// main 程序主入口
// 1. 初始化SDK配置（API Key、代理、超时等）
// 2. 创建SDK实例
// 3. 调用玩家信息批量查询API
// 4. 解析并打印玩家基础信息
// main is the program entry point
// 1. Initialize SDK configuration (API Key, proxy, timeout, etc.)
// 2. Create SDK instance
// 3. Call batch player info query API
// 4. Parse and print player basic info
func main() {
	// 日志记录启动信息
	log.Info("开始执行 Steam 玩家信息查询")

	// 初始化SDK配置
	// Initialize SDK configuration
	cfg := steamConfig.NewDefaultConfig(). // 默认配置
						WithAPIKey("********B6E87843EF78948D********"). // API Key
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

	// 调用玩家信息批量查询API(支持多个SteamID, 逗号分隔)
	// Call batch player info query API (supports multiple SteamIDs separated by commas)
	players, err := sdk.Develop.GetPlayerSummaries("76561198370695025,76561198006409530")
	if err != nil {
		log.Fatal("[Main] 获取玩家信息失败: %v", err)
	}

	// 日志记录查询结果数量
	log.Infof("成功获取 %d 个玩家信息", len(players))

	// 打印 players 结果
	fmt.Printf("获取到 %d 个玩家信息：\n", len(players))
	for _, p := range players {
		fmt.Printf("SteamID: %s\n", p.SteamID)
		fmt.Printf("昵称: %s\n", p.PersonaName)
		fmt.Printf("真实姓名: %s\n", p.RealName)
		fmt.Println("大头像: ", p.AvatarFull)
		fmt.Println("中头像: ", p.AvatarMedium)
		fmt.Println("小头像: ", p.AvatarURL)
		fmt.Printf("个人主页: %s\n", p.ProfileURL)
		fmt.Printf("是否在线: %t\n", p.IsOnline)
		fmt.Printf("最后离线时间: %s\n", p.LastLogoff)
		fmt.Printf("账号创建时间: %s\n", p.TimeCreated)

		fmt.Println("------------------------")
	}

	log.Info("Steam 玩家信息查询执行完成")
}
