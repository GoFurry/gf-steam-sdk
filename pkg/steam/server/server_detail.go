package server

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/GoFurry/gf-steam-sdk/pkg/models"
	"github.com/GoFurry/gf-steam-sdk/pkg/util/errors"
	"github.com/rumblefrog/go-a2s"
	"golang.org/x/time/rate"
)

// QueryServerInfo 查询单个服务器的基础信息(A2S_Info接口)
// 每个请求独立创建A2S Client, 避免并发资源竞争, 保证线程安全
// 参数:
//   - addr: 服务器地址, 格式为"ip:port"
//
// 返回值:
//   - a2s.ServerInfo: 服务器基础信息结构体
//   - error: 错误信息, 包含Client创建失败/接口调用失败等场景
//
// QueryServerInfo queries the basic information of a single server (A2S_Info interface)
// Each request creates an independent A2S Client to avoid concurrent resource competition and ensure thread safety
// Parameters:
//   - addr: Server address in the format of "ip:port"
//
// Return values:
//   - a2s.ServerInfo: Server basic information structure
//   - error: Error information, including Client creation failure/interface call failure and other scenarios
func (s *ServerService) QueryServerInfo(addr string) (a2s.ServerInfo, error) {
	// 创建独立的A2S Client(每个请求独立client, 避免并发冲突)
	// Create independent A2S Client (each request has its own client to avoid concurrent conflicts)
	client, err := a2s.NewClient(addr)
	if err != nil {
		return a2s.ServerInfo{}, fmt.Errorf("%w: create a2s client failed: %v", errors.ErrRequestFailed, err)
	}
	defer client.Close() // 确保client资源释放(Ensure client resource release)

	// 调用A2S_Info接口(Call A2S_Info interface)
	info, err := client.QueryInfo()
	if err != nil {
		return a2s.ServerInfo{}, fmt.Errorf("%w: query server info failed: %v", errors.ErrRequestFailed, err)
	}

	// 类型转换(适配SDK的models层)(Type conversion (adapt to SDK's models layer))
	return *info, nil
}

// QueryServerPlayers 查询单个服务器的玩家信息(A2S_Player接口)
// 每个请求独立创建A2S Client, 避免并发资源竞争, 保证线程安全
// 参数:
//   - addr: 服务器地址, 格式为"ip:port"
//
// 返回值:
//   - a2s.PlayerInfo: 服务器玩家信息结构体
//   - error: 错误信息, 包含Client创建失败/接口调用失败等场景
//
// QueryServerPlayers queries the player information of a single server (A2S_Player interface)
// Each request creates an independent A2S Client to avoid concurrent resource competition and ensure thread safety
// Parameters:
//   - addr: Server address in the format of "ip:port"
//
// Return values:
//   - a2s.PlayerInfo: Server player information structure
//   - error: Error information, including Client creation failure/interface call failure and other scenarios
func (s *ServerService) QueryServerPlayers(addr string) (a2s.PlayerInfo, error) {
	// 创建独立的A2S Client(Create independent A2S Client)
	client, err := a2s.NewClient(addr)
	if err != nil {
		return a2s.PlayerInfo{}, fmt.Errorf("%w: create a2s client failed: %v", errors.ErrRequestFailed, err)
	}
	defer client.Close() // 确保client资源释放(Ensure client resource release)

	// 调用A2S_Player接口(Call A2S_Player interface)
	players, err := client.QueryPlayer()
	if err != nil {
		return a2s.PlayerInfo{}, fmt.Errorf("%w: query server players failed: %v", errors.ErrRequestFailed, err)
	}

	return *players, nil
}

// QueryServerRules 查询单个服务器的规则信息(A2S_Rules接口)
// 每个请求独立创建A2S Client, 避免并发资源竞争, 保证线程安全
// 参数:
//   - addr: 服务器地址, 格式为"ip:port"
//
// 返回值:
//   - a2s.RulesInfo: 服务器规则信息结构体
//   - error: 错误信息, 包含Client创建失败/接口调用失败等场景
//
// QueryServerRules queries the rule information of a single server (A2S_Rules interface)
// Each request creates an independent A2S Client to avoid concurrent resource competition and ensure thread safety
// Parameters:
//   - addr: Server address in the format of "ip:port"
//
// Return values:
//   - a2s.RulesInfo: Server rule information structure
//   - error: Error information, including Client creation failure/interface call failure and other scenarios
func (s *ServerService) QueryServerRules(addr string) (a2s.RulesInfo, error) {
	// 创建独立的A2S Client(Create independent A2S Client)
	client, err := a2s.NewClient(addr)
	if err != nil {
		return a2s.RulesInfo{}, fmt.Errorf("%w: create a2s client failed: %v", errors.ErrRequestFailed, err)
	}
	defer client.Close() // 确保client资源释放(Ensure client resource release)

	// 调用A2S_Rules接口(Call A2S_Rules interface)
	rules, err := client.QueryRules()
	if err != nil {
		return a2s.RulesInfo{}, fmt.Errorf("%w: query server rules failed: %v", errors.ErrRequestFailed, err)
	}
	return *rules, nil
}

// GetServerDetail 聚合获取单个服务器的完整信息(基础信息+玩家+规则)
// 内部调用QueryServerInfo/QueryServerPlayers/QueryServerRules三个独立方法,
// 非核心接口(玩家/规则)失败仅记录错误, 核心接口(基础信息)失败直接返回
// 参数:
//   - addr: 服务器地址, 格式为"ip:port"
//
// 返回值:
//   - models.SteamServerResponse: 聚合后的服务器完整数据
//   - error: 错误信息(任意子接口失败会返回错误, 但已获取的数据会保留)
//
// GetServerDetail aggregately gets the complete information of a single server (basic info + players + rules)
// Internally calls three independent methods: QueryServerInfo/QueryServerPlayers/QueryServerRules,
// non-core interfaces (players/rules) only record errors when failed, core interface (basic info) returns directly when failed
// Parameters:
//   - addr: Server address in the format of "ip:port"
//
// Return values:
//   - models.SteamServerResponse: Aggregated complete server data
//   - error: Error information (any sub-interface failure will return an error, but the obtained data will be retained)
func (s *ServerService) GetServerDetail(addr string) (models.SteamServerResponse, error) {
	var (
		res   models.SteamServerResponse
		errs  []string // 记录各子接口的错误信息(Record error info of each sub-interface)
		mutex sync.Mutex
	)

	// 查询服务器基础信息(核心接口, 失败则直接返回)
	// Query server basic info (core interface, return directly if failed)
	info, err := s.QueryServerInfo(addr)
	if err != nil {
		errs = append(errs, fmt.Sprintf("info: %v", err))
	} else {
		res.Server = info
	}

	// 查询玩家信息(非核心, 失败仅记录错误)
	// Query player info (non-core, only record error if failed)
	players, err := s.QueryServerPlayers(addr)
	if err != nil {
		mutex.Lock()
		errs = append(errs, fmt.Sprintf("players: %v", err))
		mutex.Unlock()
	} else {
		res.Player = players
	}

	// 查询服务器规则(非核心, 失败仅记录错误)
	// Query server rules (non-core, only record error if failed)
	rules, err := s.QueryServerRules(addr)
	if err != nil {
		mutex.Lock()
		errs = append(errs, fmt.Sprintf("rules: %v", err))
		mutex.Unlock()
	} else {
		res.Rules = rules
	}

	// 聚合错误返回
	// Return aggregated errors
	if len(errs) > 0 {
		return res, fmt.Errorf("%w: get server detail failed - %s", errors.ErrRequestFailed, strings.Join(errs, "; "))
	}
	return res, nil
}

// QueryServerInfoList 批量查询多个服务器的基础信息(带限流、重试、超时)
// 采用并发方式查询, 保证结果顺序与输入地址列表一致, 支持指数退避重试策略
// 参数:
//   - addrs: 服务器地址列表, 格式为["ip:port", ...]
//   - qps: 每秒最大请求数(限流)
//   - burst: 突发请求上限(限流)
//   - timeout: 整体超时时间
//   - retry: 单个请求重试次数
//
// 返回值:
//   - []a2s.ServerInfo: 批量查询结果(与输入地址列表一一对应)
//   - []error: 每个地址对应的错误信息(nil表示成功)
//   - error: 全局错误(如上下文超时、限流失效)
//
// QueryServerInfoList batch queries the basic information of multiple servers (with rate limit, retry, timeout)
// Uses concurrent query method to ensure the result order is consistent with the input address list, supports exponential backoff retry strategy
// Parameters:
//   - addrs: Server address list in the format of ["ip:port", ...]
//   - qps: Maximum requests per second (rate limit)
//   - burst: Maximum burst requests (rate limit)
//   - timeout: Overall timeout time
//   - retry: Number of retries for a single request
//
// Return values:
//   - []a2s.ServerInfo: Batch query results (one-to-one correspondence with input address list)
//   - []error: Error information corresponding to each address (nil means success)
//   - error: Global error (such as context timeout, rate limit failure)
func (s *ServerService) QueryServerInfoList(
	addrs []string,
	qps float64,
	burst int,
	timeout time.Duration,
	retry int,
) ([]a2s.ServerInfo, []error, error) {
	// 初始化上下文和限流器(Initialize context and rate limiter)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	limiter := rate.NewLimiter(rate.Limit(qps), burst)

	// 预分配结果和错误切片(保证顺序与输入一致)
	// Preallocate result and error slices (ensure order is consistent with input)
	results := make([]a2s.ServerInfo, len(addrs))
	errs := make([]error, len(addrs))
	var (
		mutex sync.Mutex // 保护结果写入的线程安全(Protect thread safety of result writing)
		wg    sync.WaitGroup
	)

	// 遍历地址列表启动并发请求(Iterate address list and start concurrent requests)
	for idx, addr := range addrs {
		wg.Add(1)
		go func(index int, address string) {
			defer wg.Done()

			var (
				info a2s.ServerInfo
				err  error
			)

			// 单个请求的重试逻辑(指数退避)
			// Retry logic for single request (exponential backoff)
			for attempt := 1; attempt <= retry; attempt++ {
				// 检查上下文是否超时/取消(Check if context is timeout/canceled)
				if ctx.Err() != nil {
					err = fmt.Errorf("context canceled/timeout")
					break
				}

				// 限流等待(Rate limit waiting)
				if err = limiter.Wait(ctx); err != nil {
					err = fmt.Errorf("rate limit exceeded (attempt %d): %v", attempt, err)
					time.Sleep(time.Duration(attempt) * 100 * time.Millisecond) // 指数退避(Exponential backoff)
					continue
				}

				// 调用独立接口(Call independent interface)
				info, err = s.QueryServerInfo(address)
				if err == nil {
					break // 成功则退出重试(Exit retry if successful)
				}

				// 非最后一次重试则等待(Wait if not the last retry)
				if attempt < retry {
					time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
				}
			}

			// 线程安全写入结果(Thread-safe write result)
			mutex.Lock()
			defer mutex.Unlock()
			if err != nil {
				errs[index] = fmt.Errorf("%w: %v", errors.ErrRequestFailed, err)
			} else {
				results[index] = info
			}
		}(idx, addr)
	}

	// 等待所有请求完成(Wait for all requests to complete)
	wg.Wait()

	// 检查全局错误(Check global error)
	if ctx.Err() != nil {
		return nil, nil, fmt.Errorf("%w: batch query timeout/canceled: %v", errors.ErrRequestFailed, ctx.Err())
	}

	return results, errs, nil
}

// QueryServerPlayersList 批量查询多个服务器的玩家信息(带限流、重试、超时)
// 采用并发方式查询, 保证结果顺序与输入地址列表一致, 支持指数退避重试策略
// 参数/返回值格式与QueryServerInfoList完全一致
// QueryServerPlayersList batch queries the player information of multiple servers (with rate limit, retry, timeout)
// Uses concurrent query method to ensure the result order is consistent with the input address list, supports exponential backoff retry strategy
// Parameter/return value format is exactly the same as QueryServerInfoList
func (s *ServerService) QueryServerPlayersList(
	addrs []string,
	qps float64,
	burst int,
	timeout time.Duration,
	retry int,
) ([]a2s.PlayerInfo, []error, error) {
	// 初始化上下文和限流器(Initialize context and rate limiter)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	limiter := rate.NewLimiter(rate.Limit(qps), burst)

	// 预分配结果和错误切片(Preallocate result and error slices)
	results := make([]a2s.PlayerInfo, len(addrs))
	errs := make([]error, len(addrs))
	var (
		mutex sync.Mutex // 保护结果写入的线程安全(Protect thread safety of result writing)
		wg    sync.WaitGroup
	)

	// 遍历地址列表启动并发请求(Iterate address list and start concurrent requests)
	for idx, addr := range addrs {
		wg.Add(1)
		go func(index int, address string) {
			defer wg.Done()

			var (
				players a2s.PlayerInfo
				err     error
			)

			// 单个请求的重试逻辑(指数退避)
			// Retry logic for single request (exponential backoff)
			for attempt := 1; attempt <= retry; attempt++ {
				// 检查上下文是否超时/取消(Check if context is timeout/canceled)
				if ctx.Err() != nil {
					err = fmt.Errorf("context canceled/timeout")
					break
				}

				// 限流等待(Rate limit waiting)
				if err = limiter.Wait(ctx); err != nil {
					err = fmt.Errorf("rate limit exceeded (attempt %d): %v", attempt, err)
					time.Sleep(time.Duration(attempt) * 100 * time.Millisecond) // 指数退避(Exponential backoff)
					continue
				}

				// 调用独立接口(Call independent interface)
				players, err = s.QueryServerPlayers(address)
				if err == nil {
					break // 成功则退出重试(Exit retry if successful)
				}

				// 非最后一次重试则等待(Wait if not the last retry)
				if attempt < retry {
					time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
				}
			}

			// 线程安全写入结果(Thread-safe write result)
			mutex.Lock()
			defer mutex.Unlock()
			if err != nil {
				errs[index] = fmt.Errorf("%w: %v", errors.ErrRequestFailed, err)
			} else {
				results[index] = players
			}
		}(idx, addr)
	}

	// 等待所有请求完成(Wait for all requests to complete)
	wg.Wait()

	// 检查全局错误(Check global error)
	if ctx.Err() != nil {
		return nil, nil, fmt.Errorf("%w: batch query timeout/canceled: %v", errors.ErrRequestFailed, ctx.Err())
	}

	return results, errs, nil
}

// QueryServerRulesList 批量查询多个服务器的规则信息(带限流、重试、超时)
// 采用并发方式查询, 保证结果顺序与输入地址列表一致, 支持指数退避重试策略
// 参数/返回值格式与QueryServerInfoList完全一致
// QueryServerRulesList batch queries the rule information of multiple servers (with rate limit, retry, timeout)
// Uses concurrent query method to ensure the result order is consistent with the input address list, supports exponential backoff retry strategy
// Parameter/return value format is exactly the same as QueryServerInfoList
func (s *ServerService) QueryServerRulesList(
	addrs []string,
	qps float64,
	burst int,
	timeout time.Duration,
	retry int,
) ([]a2s.RulesInfo, []error, error) {
	// 初始化上下文和限流器(Initialize context and rate limiter)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	limiter := rate.NewLimiter(rate.Limit(qps), burst)

	// 预分配结果和错误切片(Preallocate result and error slices)
	results := make([]a2s.RulesInfo, len(addrs))
	errs := make([]error, len(addrs))
	var (
		mutex sync.Mutex // 保护结果写入的线程安全(Protect thread safety of result writing)
		wg    sync.WaitGroup
	)

	// 遍历地址列表启动并发请求(Iterate address list and start concurrent requests)
	for idx, addr := range addrs {
		wg.Add(1)
		go func(index int, address string) {
			defer wg.Done()

			var (
				rules a2s.RulesInfo
				err   error
			)

			// 单个请求的重试逻辑(指数退避)
			// Retry logic for single request (exponential backoff)
			for attempt := 1; attempt <= retry; attempt++ {
				// 检查上下文是否超时/取消(Check if context is timeout/canceled)
				if ctx.Err() != nil {
					err = fmt.Errorf("context canceled/timeout")
					break
				}

				// 限流等待(Rate limit waiting)
				if err = limiter.Wait(ctx); err != nil {
					err = fmt.Errorf("rate limit exceeded (attempt %d): %v", attempt, err)
					time.Sleep(time.Duration(attempt) * 100 * time.Millisecond) // 指数退避(Exponential backoff)
					continue
				}

				// 调用独立接口(Call independent interface)
				rules, err = s.QueryServerRules(address)
				if err == nil {
					break // 成功则退出重试(Exit retry if successful)
				}

				// 非最后一次重试则等待(Wait if not the last retry)
				if attempt < retry {
					time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
				}
			}

			// 线程安全写入结果(Thread-safe write result)
			mutex.Lock()
			defer mutex.Unlock()
			if err != nil {
				errs[index] = fmt.Errorf("%w: %v", errors.ErrRequestFailed, err)
			} else {
				results[index] = rules
			}
		}(idx, addr)
	}

	// 等待所有请求完成(Wait for all requests to complete)
	wg.Wait()

	// 检查全局错误(Check global error)
	if ctx.Err() != nil {
		return nil, nil, fmt.Errorf("%w: batch query timeout/canceled: %v", errors.ErrRequestFailed, ctx.Err())
	}

	return results, errs, nil
}

// GetServerDetailList 批量查询多个服务器的完整聚合信息(带限流、重试、超时)
// 内部调用GetServerDetail聚合接口, 采用并发方式查询, 保证结果顺序与输入地址列表一致
// 支持指数退避重试策略, 线程安全的结果写入机制
// 参数:
//   - addrs: 服务器地址列表, 格式为["ip:port", ...]
//   - qps: 每秒最大请求数(限流)
//   - burst: 突发请求上限(限流)
//   - timeout: 整体超时时间
//   - retry: 单个请求重试次数
//
// 返回值:
//   - []models.SteamServerResponse: 批量聚合结果(与输入地址一一对应)
//   - []error: 每个地址对应的错误信息
//   - error: 全局错误(如上下文超时、限流失效)
//
// GetServerDetailList batch queries the complete aggregated information of multiple servers (with rate limit, retry, timeout)
// Internally calls the GetServerDetail aggregation interface, uses concurrent query method to ensure the result order is consistent with the input address list
// Supports exponential backoff retry strategy and thread-safe result writing mechanism
// Parameters:
//   - addrs: Server address list in the format of ["ip:port", ...]
//   - qps: Maximum requests per second (rate limit)
//   - burst: Maximum burst requests (rate limit)
//   - timeout: Overall timeout time
//   - retry: Number of retries for a single request
//
// Return values:
//   - []models.SteamServerResponse: Batch aggregation results (one-to-one correspondence with input addresses)
//   - []error: Error information corresponding to each address
//   - error: Global error (such as context timeout, rate limit failure)
func (s *ServerService) GetServerDetailList(
	addrs []string,
	qps float64,
	burst int,
	timeout time.Duration,
	retry int,
) ([]models.SteamServerResponse, []error, error) {
	// 初始化上下文和限流器(Initialize context and rate limiter)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	limiter := rate.NewLimiter(rate.Limit(qps), burst)

	// 预分配结果和错误切片(Preallocate result and error slices)
	results := make([]models.SteamServerResponse, len(addrs))
	errs := make([]error, len(addrs))
	var (
		mutex sync.Mutex // 保护结果写入的线程安全(Protect thread safety of result writing)
		wg    sync.WaitGroup
	)

	// 遍历地址启动并发请求(Iterate addresses and start concurrent requests)
	for idx, addr := range addrs {
		wg.Add(1)
		// 捕获循环变量, 避免goroutine共享变量问题
		// Capture loop variables to avoid goroutine shared variable issues
		go func(index int, address string) {
			defer wg.Done()

			var (
				detail models.SteamServerResponse
				err    error
			)

			// 单个请求的重试逻辑(指数退避)
			// Retry logic for single request (exponential backoff)
			for attempt := 1; attempt <= retry; attempt++ {
				// 检查上下文是否已超时/取消(Check if context is timeout/canceled)
				if ctx.Err() != nil {
					err = fmt.Errorf("context canceled/timeout")
					break
				}

				// 限流等待(遵守QPS限制)(Rate limit waiting (comply with QPS limit))
				if err = limiter.Wait(ctx); err != nil {
					err = fmt.Errorf("rate limit exceeded (attempt %d): %v", attempt, err)
					// 指数退避: 重试次数越多, 等待时间越长
					// Exponential backoff: the more retries, the longer the waiting time
					time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
					continue
				}

				// 调用聚合接口(内部已调用三个独立方法)
				// Call aggregation interface (internally calls three independent methods)
				detail, err = s.GetServerDetail(address)
				if err == nil {
					break // 成功则退出重试(Exit retry if successful)
				}

				// 非最后一次重试, 等待后继续(Wait and continue if not the last retry)
				if attempt < retry {
					time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
				}
			}

			// 线程安全写入结果(避免并发读写冲突)
			// Thread-safe write result (avoid concurrent read/write conflicts)
			mutex.Lock()
			defer mutex.Unlock()
			if err != nil {
				errs[index] = fmt.Errorf("%w: %v", errors.ErrRequestFailed, err)
			} else {
				results[index] = detail
			}
		}(idx, addr)
	}

	// 等待所有并发请求完成(Wait for all concurrent requests to complete)
	wg.Wait()

	// 检查全局错误(如整体超时)(Check global error (such as overall timeout))
	if ctx.Err() != nil {
		return nil, nil, fmt.Errorf("%w: batch detail query timeout/canceled: %v", errors.ErrRequestFailed, ctx.Err())
	}

	return results, errs, nil
}
