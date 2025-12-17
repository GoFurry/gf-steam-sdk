// Package main 提供 Steam SDK 功能验证的示例测试程序
// 包含配置验证、API 调用、速率限制等多场景测试用例,
// 可作为 SDK 使用参考及功能回归测试基准
// Package main provides an example test program for verifying Steam SDK functions
// It includes multi-scenario test cases such as configuration verification, API calls, rate limiting, etc.,
// which can be used as a reference for SDK usage and a benchmark for functional regression testing

package main

import (
	"fmt"
	"sync"
	"time"

	steamConfig "github.com/GoFurry/gf-steam-sdk/pkg/config"
	"github.com/GoFurry/gf-steam-sdk/pkg/steam"
	"github.com/GoFurry/gf-steam-sdk/pkg/util"
	"github.com/GoFurry/gf-steam-sdk/pkg/util/log"
)

// 全局测试常量定义
// 集中管理测试固定参数, 提升代码可维护性
// Global test constant definitions
// Centralize management of fixed test parameters to improve code maintainability
const (
	// TestSteamID 测试用固定SteamID, 用于接口调用验证
	// TestSteamID is a fixed SteamID for testing, used for interface call verification
	TestSteamID = "76561198370695025"

	// TestAPIKey 测试用Steam API密钥（需确保有效性）
	// TestAPIKey is the Steam API key for testing (validity must be ensured)
	TestAPIKey = "********B6E87843EF78948D********"

	// TestProxyURL 适配国内访问Steam API的代理地址
	// TestProxyURL is the proxy address for accessing Steam API in China
	TestProxyURL = "http://127.0.0.1:7897"

	// TestRequestCount 限流测试并发请求数
	// TestRequestCount is the number of concurrent requests for rate limit testing
	TestRequestCount = 15

	// TestRateLimitQPS 限流测试每秒请求数限制
	// TestRateLimitQPS is the requests per second limit for rate limit testing
	TestRateLimitQPS = 2.0

	// TestRateLimitBurst 限流测试突发请求上限
	// TestRateLimitBurst is the burst request limit for rate limit testing
	TestRateLimitBurst = 5
)

// init 初始化日志组件
// SDK内置日志封装, 确保测试过程日志输出规范统一
// init initializes the log component
// The SDK has a built-in log wrapper to ensure standardized and unified log output during testing
func init() {
	if err := log.InitLogger(nil); err != nil {
		panic(fmt.Sprintf("日志初始化失败: %v", err))
	}
}

// main 测试程序入口函数
// 按需启用/注释不同测试函数, 支持单场景或多场景组合测试
// 所有测试用例执行结果会通过showResult函数统一输出
// main is the entry function of the test program
// Enable/comment different test functions as needed to support single-scenario or multi-scenario combined testing
// The execution results of all test cases are uniformly output through the showResult function
func main() {
	// 记录整体测试耗时
	// Record the total test time
	totalStart := time.Now()
	defer func() {
		fmt.Printf("\n==== 所有测试执行完成, 总耗时: %v ====\n", time.Since(totalStart))
	}()

	// 执行测试用例(按需切换)
	// Execute test cases (switch as needed)
	showResult(1, Test_1_DefaultConfig()) // 测试默认配置下的API调用
	showResult(2, Test_2_ValidConfig())   // 测试有效配置下的API调用
	showResult(3, Test_3_RateLimit())     // 测试速率限制功能
}

// Test_1_DefaultConfig 测试默认配置下的Steam API调用
// 测试目的：验证默认配置(无API Key/代理)的错误处理逻辑
// 预期结果：
//   - 国内环境：因无代理触发网络超时(context deadline exceeded)
//   - 海外环境：因无有效API Key返回403 Forbidden
//
// 返回值：bool - true表示测试符合预期(调用失败), false表示测试不符合预期(调用成功)
// Test_1_DefaultConfig tests Steam API calls under default configuration
// Test purpose: Verify the error handling logic of default configuration (no API Key/proxy)
// Expected results:
//   - Domestic environment: Network timeout triggered by no proxy (context deadline exceeded)
//   - Overseas environment: 403 Forbidden returned due to no valid API Key
//
// Return value: bool - true means the test meets expectations (call failed), false means not (call succeeded)
func Test_1_DefaultConfig() bool {
	testName := "默认配置API调用测试(Test_1_DefaultConfig)"
	fmt.Printf("\n==== 开始执行 %s ====\n", testName)
	start := time.Now()
	defer func() {
		fmt.Printf("==== %s 执行完成, 耗时: %v ====\n", testName, time.Since(start))
	}()

	// 初始化默认配置
	// Initialize default configuration
	cfg := steamConfig.NewDefaultConfig()
	log.Info("默认配置初始化完成, 未设置API Key和代理地址")

	// 创建SDK实例
	// Create SDK instance
	sdk, err := steam.NewSteamSDK(cfg)
	if err != nil {
		log.Error("SDK实例创建失败: %v", err)
		return false
	}

	// 调用玩家信息查询接口
	// Call player information query interface
	log.Infof("调用GetPlayerSummaries接口, SteamID: %s", TestSteamID)
	summaries, err := sdk.Player.GetPlayerSummaries(TestSteamID)

	// 结果处理
	// Result processing
	if err != nil {
		log.Error("API调用失败(预期结果): %v", err)
		return true
	}

	// 打印有效返回结果(非预期分支, 仅作兜底)
	// Print valid return results (unexpected branch, only for fallback)
	log.Infof("API调用成功, 玩家信息: %+v", summaries[0])
	fmt.Println("玩家信息:", summaries[0])

	return false
}

// Test_2_ValidConfig 测试有效配置下的Steam API调用
// 测试目的：验证带有效API Key和代理配置的正常API调用流程
// 测试条件：
//   - 配置有效Steam API Key
//   - 配置可访问Steam的代理地址
//
// 预期结果：成功获取玩家基本信息并打印
// 返回值：bool - true表示测试通过(调用成功), false表示测试失败(调用失败)
// Test_2_ValidConfig tests Steam API calls under valid configuration
// Test purpose: Verify the normal API call process with valid API Key and proxy configuration
// Test conditions:
//   - Configure a valid Steam API Key
//   - Configure a proxy address accessible to Steam
//
// Expected result: Successfully obtain and print basic player information
// Return value: bool - true means test passed (call succeeded), false means test failed (call failed)
func Test_2_ValidConfig() bool {
	testName := "有效配置API调用测试(Test_2_ValidConfig)"
	fmt.Printf("\n==== 开始执行 %s ====\n", testName)
	start := time.Now()
	defer func() {
		fmt.Printf("==== %s 执行完成, 耗时: %v ====\n", testName, time.Since(start))
	}()

	// 构建有效配置(链式配置)
	// Build valid configuration (chain configuration)
	cfg := steamConfig.NewDefaultConfig().
		WithAPIKey(TestAPIKey).
		WithProxyURL(TestProxyURL)
	log.Infof("有效配置初始化完成, API Key: %s, 代理地址: %s", TestAPIKey, TestProxyURL)

	// 创建SDK实例
	// Create SDK instance
	sdk, err := steam.NewSteamSDK(cfg)
	if err != nil {
		log.Error("SDK实例创建失败: %v", err)
		return false
	}

	// 调用玩家信息查询接口
	// Call player information query interface
	log.Infof("调用GetPlayerSummaries接口, SteamID: %s", TestSteamID)
	summaries, err := sdk.Player.GetPlayerSummaries(TestSteamID)

	// 结果处理
	// Result processing
	if err != nil {
		log.Error("API调用失败(非预期结果): %v", err)
		return false
	}

	// 打印玩家信息
	// Print player information
	log.Info("API调用成功, 获取玩家信息")
	fmt.Println("玩家信息:", summaries[0])
	return true
}

// Test_3_RateLimit 测试Steam API调用的速率限制功能
// 测试目的：验证令牌桶限流器的限流效果(QPS限制+突发请求控制)
// 测试参数：
//   - QPS限制: 2.0(每秒最多2个请求)
//   - 突发上限: 5(允许5个请求瞬间执行)
//   - 并发请求数: 15(超出突发上限, 验证排队等待逻辑)
//
// 预期结果：
//   - 前5个请求(突发上限)快速执行
//   - 后续请求按0.5秒/个的间隔执行
//   - 部分请求因超时触发限流等待失败
//
// 返回值：bool - true表示测试通过(限流逻辑生效), false表示测试失败(SDK初始化失败)
// Test_3_RateLimit tests the rate limiting function of Steam API calls
// Test purpose: Verify the rate limiting effect of the token bucket limiter (QPS limit + burst request control)
// Test parameters:
//   - QPS limit: 2.0 (maximum 2 requests per second)
//   - Burst limit: 5 (allow 5 requests to execute instantly)
//   - Number of concurrent requests: 15 (exceed burst limit, verify queue waiting logic)
//
// Expected results:
//   - The first 5 requests (burst limit) execute quickly
//   - Subsequent requests execute at intervals of 0.5 seconds per request
//   - Some requests fail due to timeout triggering rate limit waiting
//
// Return value: bool - true means test passed (rate limit logic took effect), false means test failed (SDK initialization failed)
func Test_3_RateLimit() bool {
	testName := "速率限制功能测试(Test_3_RateLimit)"
	fmt.Printf("\n==== 开始执行 %s ====\n", testName)
	start := time.Now()
	defer func() {
		fmt.Printf("==== %s 执行完成, 耗时: %v ====\n", testName, time.Since(start))
	}()

	// 构建带限流配置的实例
	// Build instance with rate limit configuration
	cfg := steamConfig.NewDefaultConfig().
		WithAPIKey(TestAPIKey).
		WithTimeout(2*time.Second).
		WithRateLimit(TestRateLimitQPS, TestRateLimitBurst).
		WithProxyURL(TestProxyURL)
	log.Infof("限流配置初始化完成, QPS: %.1f, 突发上限: %d, 超时时间: %v",
		TestRateLimitQPS, TestRateLimitBurst, 2*time.Second)

	// 创建SDK实例
	// Create SDK instance
	sdk, err := steam.NewSteamSDK(cfg)
	if err != nil {
		log.Error("SDK实例创建失败: %v", err)
		return false
	}

	// 初始化等待组, 确保所有并发请求执行完成
	// Initialize wait group to ensure all concurrent requests are completed
	var wg sync.WaitGroup
	wg.Add(TestRequestCount)

	// 启动并发请求
	// Start concurrent requests
	log.Infof("启动%d个并发请求测试限流效果", TestRequestCount)
	for i := 0; i < TestRequestCount; i++ {
		// 捕获循环变量, 避免goroutine共享同一变量导致的索引错误
		// Capture loop variable to avoid index errors caused by goroutines sharing the same variable
		reqIndex := i
		go func() {
			defer wg.Done()

			// 记录单个请求开始时间
			// Record the start time of a single request
			reqStart := time.Now()
			summaries, err := sdk.Player.GetPlayerSummaries(TestSteamID)

			// 单个请求结果处理
			// Single request result processing
			if err != nil {
				log.Warnf("第%d个请求失败: %v, 耗时: %v", reqIndex, err, time.Since(reqStart))
				fmt.Printf("第%d个请求失败: %v\n", reqIndex, err)
				return
			}

			log.Infof("第%d个请求成功, 耗时: %v", reqIndex, time.Since(reqStart))
			fmt.Printf("第%d个请求成功: %+v\n", reqIndex, summaries[0])
		}()
	}

	// 等待所有请求完成
	// Wait for all requests to complete
	wg.Wait()
	log.Info("所有并发请求执行完成")
	return true
}

// showResult 统一输出测试用例执行结果
// 入参：
//   - idx: 测试用例编号
//   - status: 测试结果(true=成功/符合预期, false=失败/不符合预期)
//
// 功能：通过日志组件输出标准化的测试结果，便于后续集成测试解析
// showResult uniformly outputs the execution results of test cases
// Parameters:
//   - idx: Test case number
//   - status: Test result (true=success/meets expectations, false=failure/does not meet expectations)
//
// Function: Output standardized test results through the log component for subsequent integration test parsing
func showResult(idx int, status bool) {
	caseID := util.Int2String(idx)
	if status {
		log.Info("测试用例 " + caseID + " Success")
		fmt.Printf("\n==== 测试用例 %s 执行成功 ====\n", caseID)
	} else {
		log.Error("测试用例 " + caseID + " Failure")
		fmt.Printf("\n==== 测试用例 %s 执行失败 ====\n", caseID)
	}
}
