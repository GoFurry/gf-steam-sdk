// Package client 实现 Steam API 客户端核心逻辑
// 提供通用的 HTTP 请求封装、速率限制、重试机制和响应解析能力
// Package client implements the core logic of Steam API client
// Provides universal HTTP request encapsulation, rate limiting, retry mechanism and response parsing capabilities

package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/GoFurry/gf-steam-sdk/pkg/config"
	"github.com/GoFurry/gf-steam-sdk/pkg/util"
	"github.com/GoFurry/gf-steam-sdk/pkg/util/errors"
	"github.com/bytedance/sonic"
	"golang.org/x/time/rate"
)

// Client Steam API 客户端结构体
// 封装 HTTP 客户端、速率限制器和配置, 提供统一的 API 请求能力
// Client is the Steam API client structure
// Encapsulates HTTP client, rate limiter and configuration to provide unified API request capabilities
type Client struct {
	cfg     *config.SteamConfig // 全局配置 | Global configuration
	client  *http.Client        // 可复用 HTTP 客户端 | Reusable HTTP client
	limiter *rate.Limiter       // 速率限制器 | Rate limiter
}

// NewClient 创建 Steam API 客户端实例
// 参数:
//   - cfg: 配置实例(传 nil 则使用默认配置)| Configuration instance (use default if nil)
//
// 返回值:
//   - *Client: 客户端实例 | Client instance
//   - error: 配置校验失败时返回错误 | Error if config validation fails
func NewClient(cfg *config.SteamConfig) (*Client, error) {
	// 使用默认配置
	// Use default configuration if nil
	if cfg == nil {
		cfg = config.NewDefaultConfig()
	}

	if cfg.IsDebug {
		fmt.Printf("[Info] Start NewClient Init \n")
	}

	// 校验配置合法性
	// Validate configuration legality
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validate failed: %w", err)
	}

	// 构建可复用的 HTTP 客户端
	// Build reusable HTTP client
	httpClient := &http.Client{
		Timeout:   cfg.Timeout,
		Transport: cfg.Transport,
	}

	// 初始化速率限制器
	// Initialize rate limiter (use config value first, otherwise default)
	qps := util.DEFAULT_RATE_QPS
	if cfg.RateLimitQPS > 0 {
		qps = cfg.RateLimitQPS
	}
	burst := util.DEFAULT_RATE_BURST
	if cfg.RateLimitBurst > 0 {
		burst = cfg.RateLimitBurst
	}
	limiter := rate.NewLimiter(rate.Limit(qps), burst)

	if cfg.IsDebug {
		fmt.Printf("[Info] End NewClient Init \n")
	}

	return &Client{
		cfg:     cfg,
		client:  httpClient,
		limiter: limiter,
	}, nil
}

// DoRequest 通用 API 请求方法
// 支持速率限制、自动重试、状态码校验和 JSON 响应解析
// 参数:
//   - method: HTTP 请求方法(GET/POST 等) | HTTP request method (GET/POST, etc.)
//   - baseURL: 请求基础地址 | Request base URL
//   - params: 请求查询参数 | Request query parameters
//
// 返回值:
//   - map[string]interface{}: 解析后的 JSON 响应 | Parsed JSON response
//   - error: 请求/解析失败时返回错误 | Error if request/parsing fails
func (c *Client) DoRequest(method, baseURL string, params url.Values) (map[string]interface{}, error) {
	if c.cfg.IsDebug {
		fmt.Printf("[Info] Start DoRequest \n")
	}

	// 创建带超时的上下文
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.cfg.Timeout)
	defer cancel()

	// 速率限制: 等待获取令牌
	// Rate limit: wait for token
	if err := c.limiter.Wait(ctx); err != nil { // 等待获取令牌
		return nil, fmt.Errorf("%w: request rate limit exceeded: %v", errors.ErrRequestFailed, err)
	}

	// 追加 API Key 到请求参数
	// Append API Key to request parameters
	params.Set("key", c.cfg.APIKey)
	params.Set("access_token", c.cfg.AccessToken)

	// 构建完整请求 URL
	// Build full request URL
	requestURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse url failed: %w", err)
	}
	requestURL.RawQuery = params.Encode()

	// 带重试机制发送请求
	// Send request with retry mechanism
	var resp *http.Response
	var errRequest error
	var costTime time.Duration
	for i := 0; i <= c.cfg.RetryTimes; i++ {
		startTime := time.Now()

		// 创建 HTTP 请求
		// Create HTTP request
		req, err := http.NewRequest(method, requestURL.String(), nil)
		if err != nil {
			errRequest = fmt.Errorf("create request failed: %w", err)
			continue
		}

		// 设置请求头
		// Set request headers (default UA + custom headers)
		req.Header.Set("User-Agent", util.USER_AGENT)
		for k, v := range c.cfg.Headers {
			req.Header.Set(k, v)
		}

		if c.cfg.IsDebug {
			fmt.Printf("[Info] Start %s %s \n", method, requestURL.String())
		}
		// 发送请求
		// Send request
		resp, err = c.client.Do(req)
		costTime = time.Since(startTime)

		// 请求成功(200 状态码)则退出重试
		// Exit retry if request succeeds (200 status code)
		if err == nil && resp.StatusCode == http.StatusOK {
			if c.cfg.IsDebug {
				fmt.Printf("[Success] cost time: %v, %s %s \n", costTime, method, requestURL.String())
			}
			break
		}

		// 记录重试错误
		// Record retry error
		statusCode := "nil"
		if resp != nil {
			statusCode = util.Int2String(resp.StatusCode)
		}
		reqErr := "nil"
		if err != nil {
			reqErr = err.Error()
		}
		errRequest = fmt.Errorf("request failed (retry %d): status_code=%v, err=%w, cost_time=%v",
			i, statusCode, reqErr, costTime)

		if c.cfg.IsDebug {
			fmt.Printf("[Error] %s \n", errRequest.Error())
		}

		// 仅对 429(限流)/5xx(服务器错误) 进行重试
		// Only retry for 429 (rate limit)/5xx (server error)
		if resp != nil && (resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500) {
			// 指数退避: (i+1)*基础重试间隔
			// Exponential backoff: (i+1)*base retry interval
			time.Sleep(time.Duration((i+1)*util.RETRY_SLEEP_BASE) * time.Millisecond)
			continue
		}
		break
	}

	// 检查请求错误
	// Check request error
	if errRequest != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrRequestFailed, errRequest)
	}
	defer resp.Body.Close()

	// 校验响应状态码
	// Validate response status code
	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("request failed with status code %d", resp.StatusCode)
		switch resp.StatusCode {
		case http.StatusTooManyRequests:
			return nil, errors.ErrAPIQuotaExceeded
		case http.StatusBadRequest:
			return nil, fmt.Errorf("%w: invalid request params", errors.ErrRequestFailed)
		case http.StatusUnauthorized:
			return nil, fmt.Errorf("%w: invalid api key", errors.ErrRequestFailed)
		default:
			return nil, fmt.Errorf("%w: %s", errors.ErrRequestFailed, errMsg)
		}
	}

	// 读取并解析 JSON 响应
	// Read and parse JSON response
	var result map[string]interface{}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}
	if err := sonic.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrAPIResponse, err)
	}

	if c.cfg.IsDebug {
		fmt.Printf("[Info] End DoRequest \n")
	}
	return result, nil
}
