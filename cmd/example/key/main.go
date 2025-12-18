// Package main 提供 Steam 认证信息获取工具的使用示例
// 演示如何通过 SDK 的 Util 模块快速调用各类令牌/API Key 的获取指引方法
// 帮助开发者快速上手使用工具服务
// Package main provides usage examples of Steam authentication information acquisition tools
// Demonstrates how to quickly call various token/API Key acquisition guidance methods through the SDK's Util module,
// helping developers quickly get started with the utility service
package main

import (
	steamConfig "github.com/GoFurry/gf-steam-sdk/pkg/config"
	"github.com/GoFurry/gf-steam-sdk/pkg/steam"
)

// main 函数是工具使用示例的入口
// 初始化 Steam SDK 后, 依次调用 Util 模块的 API Key、Store 令牌、Community 令牌获取方法,
// 自动打开对应页面并输出操作指引
// The main function is the entry point for the tool usage example
// After initializing the Steam SDK, call the API Key, Store token, and Community token acquisition methods of the Util module in sequence,
// automatically open the corresponding pages and output operation guidelines
func main() {
	// 初始化 SDK 默认配置(Initialize SDK default configuration)
	cfg := steamConfig.NewDefaultConfig()

	// 创建 Steam SDK 实例(Create Steam SDK instance)
	sdk, _ := steam.NewSteamSDK(cfg)

	// 调用工具方法获取 API Key(Call utility method to get API Key)
	sdk.Util.GetAPIKey()
	// 调用工具方法获取 Store 令牌(Call utility method to get Store token)
	sdk.Util.GetStoreToken()
	// 调用工具方法获取 Community 令牌(Call utility method to get Community token)
	sdk.Util.GetCommunityToken()
}
