// Package main 提供 Steam SDK A2S 功能的完整使用示例
// 包含单个服务器信息查询、批量独立接口查询、批量聚合信息查询三类场景
// 演示限流、重试、超时等核心功能的使用方式
// Package main provides a complete usage example of the Steam SDK A2S function
// Including three scenarios: single server info query, batch independent interface query, batch aggregated info query
// Demonstrates the usage of core functions such as rate limiting, retry, and timeout
package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	steamConfig "github.com/GoFurry/gf-steam-sdk/pkg/config"
	"github.com/GoFurry/gf-steam-sdk/pkg/steam"
	steamLog "github.com/GoFurry/gf-steam-sdk/pkg/util/log"
)

// 全局测试常量
// Global test constants
const (
	// TestServerAddr 测试用单个服务器地址
	// TestServerAddr is the single server address for testing
	TestServerAddr = "110.42.54.147:52023"

	// TestServerList 测试用批量服务器地址列表
	// TestServerList is the batch server address list for testing
	TestServerList = "110.42.54.147:52021,110.42.54.147:52022,110.42.54.147:52023,110.42.54.147:52024,110.42.54.147:52025"

	// 限流配置 Rate limit configuration
	TestQPS     = 2.0              // 每秒最大请求数
	TestBurst   = 5                // 突发请求上限
	TestRetry   = 3                // 单个请求重试次数
	TestTimeout = 30 * time.Second // 整体超时时间
)

// 初始化日志
// Initialize logger (using the simple log wrapper built into the SDK)
func init() {
	if err := steamLog.InitLogger(nil); err != nil {
		panic(fmt.Sprintf("日志初始化失败: %v", err))
	}
}

func main() {
	steamLog.Info("========== 开始执行 Steam A2S 功能示例 ==========")

	// 初始化 SDK 配置
	// Initialize SDK configuration
	steamLog.Info("初始化SDK配置")
	cfg := steamConfig.NewDefaultConfig().
		WithProxyURL("http://127.0.0.1:7897") // 中国地区必须配置代理
	steamLog.Infof("SDK配置完成，代理地址: %s", "http://127.0.0.1:7897")

	// 创建 Steam SDK 实例
	// Create Steam SDK instance
	steamLog.Info("创建Steam SDK实例")
	sdk, err := steam.NewSteamSDK(cfg)
	if err != nil {
		steamLog.Fatalf("[Main] 创建 Steam SDK 失败: %v", err)
		log.Fatalf("[Main] 创建 Steam SDK 失败: %v", err)
	}
	steamLog.Info("Steam SDK实例创建成功")

	// 场景1: 单个服务器完整信息查询
	// Scenario 1: Single server complete info query
	steamLog.Info("\n========== 场景1: 单个服务器完整信息查询 ==========")
	testSingleServerQuery(sdk, TestServerAddr)

	// 场景2: 批量独立接口查询（分别查询基础信息/玩家/规则）
	// Scenario 2: Batch independent interface query (query basic info/players/rules separately)
	steamLog.Info("\n========== 场景2: 批量独立接口查询 ==========")
	serverList := strings.Split(TestServerList, ",")
	testBatchIndependentQuery(sdk, serverList)

	// 场景3: 批量聚合信息查询
	// Scenario 3: Batch aggregated info query
	steamLog.Info("\n========== 场景3: 批量聚合信息查询 ==========")
	testBatchAggregatedQuery(sdk, serverList)

	steamLog.Info("\n========== Steam A2S 功能示例执行完成 ==========")
}

// testSingleServerQuery 测试单个服务器完整信息查询
// testSingleServerQuery tests single server complete info query
func testSingleServerQuery(sdk *steam.SteamSDK, addr string) {
	steamLog.Infof("开始查询单个服务器信息: %s", addr)
	start := time.Now()

	// 调用聚合接口获取完整信息
	detail, err := sdk.Server.GetServerDetail(addr)
	if err != nil {
		steamLog.Errorf("查询单个服务器信息失败: %v", err)
		fmt.Printf("查询失败: %v\n", err)
		return
	}

	// 打印查询结果
	steamLog.Infof("单个服务器信息查询完成，耗时: %v", time.Since(start))
	fmt.Println("\n=== 单个服务器查询结果 ===")
	fmt.Printf("服务器地址: %s\n", addr)
	fmt.Printf("基础信息: %+v\n", detail.Server)
	fmt.Printf("玩家信息: %+v\n", detail.Player)
	fmt.Printf("规则信息: %+v\n", detail.Rules)
}

// testBatchIndependentQuery 测试批量独立接口查询
// 分别调用QueryServerInfoList/QueryServerPlayersList/QueryServerRulesList
// testBatchIndependentQuery tests batch independent interface query
// Call QueryServerInfoList/QueryServerPlayersList/QueryServerRulesList separately
func testBatchIndependentQuery(sdk *steam.SteamSDK, addrs []string) {
	steamLog.Infof("开始批量独立接口查询，服务器数量: %d", len(addrs))
	start := time.Now()

	// 3.1 批量查询基础信息
	steamLog.Info("3.1 批量查询服务器基础信息")
	infoList, infoErrs, err := sdk.Server.QueryServerInfoList(addrs, TestQPS, TestBurst, TestTimeout, TestRetry)
	if err != nil {
		steamLog.Errorf("批量查询基础信息全局错误: %v", err)
		fmt.Printf("批量查询基础信息全局错误: %v\n", err)
	} else {
		printBatchResult("基础信息", addrs, infoList, infoErrs)
	}

	// 3.2 批量查询玩家信息
	steamLog.Info("3.2 批量查询服务器玩家信息")
	playerList, playerErrs, err := sdk.Server.QueryServerPlayersList(addrs, TestQPS, TestBurst, TestTimeout, TestRetry)
	if err != nil {
		steamLog.Errorf("批量查询玩家信息全局错误: %v", err)
		fmt.Printf("批量查询玩家信息全局错误: %v\n", err)
	} else {
		printBatchResult("玩家信息", addrs, playerList, playerErrs)
	}

	// 3.3 批量查询规则信息
	steamLog.Info("3.3 批量查询服务器规则信息")
	ruleList, ruleErrs, err := sdk.Server.QueryServerRulesList(addrs, TestQPS, TestBurst, TestTimeout, TestRetry)
	if err != nil {
		steamLog.Errorf("批量查询规则信息全局错误: %v", err)
		fmt.Printf("批量查询规则信息全局错误: %v\n", err)
	} else {
		printBatchResult("规则信息", addrs, ruleList, ruleErrs)
	}

	steamLog.Infof("批量独立接口查询完成，总耗时: %v", time.Since(start))
}

// testBatchAggregatedQuery 测试批量聚合信息查询
// testBatchAggregatedQuery tests batch aggregated info query
func testBatchAggregatedQuery(sdk *steam.SteamSDK, addrs []string) {
	steamLog.Infof("开始批量聚合信息查询，服务器数量: %d", len(addrs))
	start := time.Now()

	// 调用批量聚合接口
	detailList, detailErrs, err := sdk.Server.GetServerDetailList(addrs, TestQPS, TestBurst, TestTimeout, TestRetry)
	if err != nil {
		steamLog.Errorf("批量聚合查询全局错误: %v", err)
		fmt.Printf("批量聚合查询全局错误: %v\n", err)
		return
	}

	// 打印批量聚合结果
	steamLog.Infof("批量聚合查询完成，耗时: %v", time.Since(start))
	fmt.Println("\n=== 批量聚合查询结果 ===")
	for idx, addr := range addrs {
		if detailErrs[idx] != nil {
			fmt.Printf("服务器 %s 查询失败: %v\n", addr, detailErrs[idx])
			continue
		}
		fmt.Printf("\n服务器 %s 聚合信息:\n", addr)
		fmt.Printf("  基础信息: %+v\n", detailList[idx].Server)
		fmt.Printf("  玩家信息: %+v\n", detailList[idx].Player)
		fmt.Printf("  规则信息: %+v\n", detailList[idx].Rules)
	}
}

// printBatchResult 通用批量结果打印函数
// printBatchResult is a generic batch result printing function
func printBatchResult(resultType string, addrs []string, results interface{}, errs []error) {
	fmt.Printf("\n=== 批量%s查询结果 ===\n", resultType)
	for idx, addr := range addrs {
		if errs[idx] != nil {
			fmt.Printf("服务器 %s: 查询失败 - %v\n", addr, errs[idx])
			continue
		}
		fmt.Printf("服务器 %s: %+v\n", addr, results)
		// 修复：根据不同类型打印具体结果
		switch resultType {
		case "基础信息":
			if infoList, ok := results.([]interface{}); ok {
				fmt.Printf("服务器 %s: %+v\n", addr, infoList[idx])
			}
		case "玩家信息":
			if playerList, ok := results.([]interface{}); ok {
				fmt.Printf("服务器 %s: %+v\n", addr, playerList[idx])
			}
		case "规则信息":
			if ruleList, ok := results.([]interface{}); ok {
				fmt.Printf("服务器 %s: %+v\n", addr, ruleList[idx])
			}
		}
	}
}
