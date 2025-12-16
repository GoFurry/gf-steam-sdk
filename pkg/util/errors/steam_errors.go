// Package errors 提供 Steam SDK 自定义错误体系
// 包含错误类型枚举、自定义错误结构体和预设通用错误, 支持标准库 error 接口
// Package errors provides custom error system for Steam SDK
// Includes error type enumeration, custom error structure and preset common errors, supports standard library error interface

package errors

import (
	"errors"
	"fmt"
)

// ErrType 错误类型枚举
// 分类管理不同场景的错误，便于错误处理和排查
// ErrType is the error type enumeration
// Manages errors of different scenarios for easy error handling and troubleshooting
type ErrType string

const (
	ErrTypeParam   ErrType = "param"   // 参数错误(如空SteamID、无效AppID)
	ErrTypeRequest ErrType = "request" // 请求错误(如网络、超时、429)
	ErrTypeParse   ErrType = "parse"   // 解析错误(如JSON反序列化失败)
	ErrTypeAPI     ErrType = "api"     // Steam API业务错误(如success=false)
	ErrTypeCrawler ErrType = "crawler" //爬虫
)

// SteamError 自定义错误结构体
// 封装错误类型、错误码、描述和原始错误，支持标准库错误处理方法
// SteamError is the custom error structure
// Encapsulates error type, error code, description and original error, supports standard library error handling methods
type SteamError struct {
	Type    ErrType // 错误类型 | Error type
	Code    int     // 错误码 | Error code
	Message string  // 错误描述 | Error description
	Err     error   // 原始错误 | Original error
}

// Error 实现 error 接口
// 返回格式化的错误信息
// 返回值:
//   - string: 错误字符串 | Error string
func (e *SteamError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("[%s] %s", e.Type, e.Message)
	}
	return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Err)
}

// Unwrap 支持标准库 errors.Unwrap
// 返回原始错误
// 返回值:
//   - error: 原始错误 | Original error
func (e *SteamError) Unwrap() error {
	return e.Err
}

// Is 支持标准库 errors.Is
// 比较两个 SteamError 是否相同(类型+描述)
// 参数:
//   - target: 目标错误 | Target error
//
// 返回值:
//   - bool: 是否匹配 | Whether matched
func (e *SteamError) Is(target error) bool {
	t, ok := target.(*SteamError)
	if !ok {
		return false
	}
	return e.Type == t.Type && e.Message == t.Message
}

// 预设通用错误 | Preset common errors
var (
	// ErrMissingAPIKey API Key缺失
	ErrMissingAPIKey = &SteamError{
		Type:    ErrTypeParam,
		Code:    10001,
		Message: "steam api key is missing",
		Err:     errors.New("steam api key is missing"),
	}

	// ErrInvalidSteamID 无效SteamID
	ErrInvalidSteamID = &SteamError{
		Type:    ErrTypeParam,
		Code:    10002,
		Message: "invalid steam id",
		Err:     errors.New("invalid steam id"),
	}

	// ErrInvalidAppID 无效AppID
	ErrInvalidAppID = &SteamError{
		Type:    ErrTypeParam,
		Code:    10003,
		Message: "invalid app id (app id cannot be 0)",
		Err:     errors.New("invalid app id"),
	}

	// ErrAPIResponse 响应解析失败
	ErrAPIResponse = &SteamError{
		Type:    ErrTypeParse,
		Code:    20001,
		Message: "failed to parse steam api response",
		Err:     errors.New("json unmarshal failed"),
	}

	// ErrAPIQuotaExceeded API额度超限
	ErrAPIQuotaExceeded = &SteamError{
		Type:    ErrTypeRequest,
		Code:    30001,
		Message: "steam api quota exceeded (429)",
		Err:     errors.New("api rate limit exceeded"),
	}

	// ErrRequestFailed 请求失败
	ErrRequestFailed = &SteamError{
		Type:    ErrTypeRequest,
		Code:    30002,
		Message: "steam api request failed",
		Err:     errors.New("http request failed"),
	}

	// ErrAchievementFailed 成就查询失败
	ErrAchievementFailed = &SteamError{
		Type:    ErrTypeAPI,
		Code:    40001,
		Message: "steam api return success=false (achievements not found or permission denied)",
		Err:     errors.New("achievements query failed"),
	}
)

// New 快速创建自定义SteamError
// 默认类型为 ErrTypeAPI
// 参数:
//   - msg: 错误描述 | Error message
//
// 返回值:
//   - error: 自定义错误 | Custom error
func New(msg string) error {
	return &SteamError{
		Type:    ErrTypeAPI,
		Message: msg,
		Err:     errors.New(msg),
	}
}

// NewWithType 按类型创建自定义错误
// 参数:
//   - errType: 错误类型 | Error type
//   - msg: 错误描述 | Error message
//   - err: 原始错误 | Original error
//
// 返回值:
//   - error: 自定义错误 | Custom error
func NewWithType(errType ErrType, msg string, err error) error {
	return &SteamError{
		Type:    errType,
		Message: msg,
		Err:     err,
	}
}

// GetType 获取错误类型
// 参数:
//   - err: 错误实例 | Error instance
//
// 返回值:
//   - ErrType: 错误类型(非SteamError返回空) | Error type (empty for non-SteamError)
func GetType(err error) ErrType {
	var steamErr *SteamError
	if errors.As(err, &steamErr) {
		return steamErr.Type
	}
	return ""
}

// GetCode 获取错误码
// 参数:
//   - err: 错误实例 | Error instance
//
// 返回值:
//   - int: 错误码(非SteamError返回-1) | Error code (-1 for non-SteamError)
func GetCode(err error) int {
	var steamErr *SteamError
	if errors.As(err, &steamErr) {
		return steamErr.Code
	}
	return -1
}
