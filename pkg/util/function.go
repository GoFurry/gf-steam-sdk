// Package util 提供 Steam SDK 通用工具函数和常量
// 包含类型转换、时间处理、字符串操作和默认配置常量等
// Package util provides common utility functions and constants for Steam SDK
// Includes type conversion, time processing, string operations and default config constants

package util

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// B2i bool 转 int
// 参数:
//   - b: 布尔值 | Boolean value
//
// 返回值:
//   - int: 1(true)/0(false)
func B2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

// TimeUnix2String 时间戳转字符串
// 使用 TIME_FORMAT_DATE 格式
// 参数:
//   - unixTime: 秒级时间戳 | Second-level timestamp
//
// 返回值:
//   - string: 格式化时间字符串（空字符串如果时间戳<=0）| Formatted time string (empty if timestamp <=0)
func TimeUnix2String(unixTime int64) string {
	if unixTime <= 0 {
		return ""
	}
	return time.Unix(unixTime, 0).Format(TIME_FORMAT_DATE)
}

// TimeString2Unix 字符串转时间戳
// 仅支持 TIME_FORMAT_DATE 格式
// 参数:
//   - timeStr: 时间字符串 | Time string
//
// 返回值:
//   - int64: 秒级时间戳 | Second-level timestamp
//   - error: 解析失败时返回错误 | Error if parse failed
func TimeString2Unix(timeStr string) (int64, error) {
	t, err := time.Parse(TIME_FORMAT_DATE, timeStr)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

// String2Int 字符串转 int
// 参数:
//   - numString: 数字字符串 | Numeric string
//
// 返回值:
//   - int: 转换后的整数 | Converted integer
//   - error: 空字符串或转换失败时返回错误 | Error if empty or conversion failed
func String2Int(numString string) (int, error) {
	if strings.TrimSpace(numString) == "" {
		return 0, errors.New("字符串不能为空") // String cannot be empty
	}
	id, err := strconv.Atoi(numString)
	return id, err
}

// Int642String int64 转字符串
// 参数:
//   - i64: int64 数值 | int64 value
//
// 返回值:
//   - string: 字符串形式 | String form
func Int642String(i64 int64) string { return strconv.FormatInt(i64, 10) }

// Uint642String uint64 转字符串
// 参数:
//   - ui64: uint64 数值 | uint64 value
//
// 返回值:
//   - string: 字符串形式 | String form
func Uint642String(ui64 uint64) string { return fmt.Sprintf("%d", ui64) }

// Int2String int 转字符串
// 参数:
//   - i: int 数值 | int value
//
// 返回值:
//   - string: 字符串形式 | String form
func Int2String(i int) string { return fmt.Sprintf("%d", i) }

// String2Int64 字符串转 int64
// 参数:
//   - numString: 数字字符串 | Numeric string
//
// 返回值:
//   - int64: 转换后的整数 | Converted integer
//   - error: 空字符串或转换失败时返回错误 | Error if empty or conversion failed
func String2Int64(numString string) (int64, error) {
	if strings.TrimSpace(numString) == "" {
		return 0, errors.New("参数不能为空") // Parameter cannot be empty
	}
	id, parseErr := strconv.ParseInt(strings.TrimSpace(numString), 10, 64)
	return id, parseErr
}

// MaskAPIKey 脱敏API Key
// 保留后8位，前面替换为********，长度<=8时不脱敏
// 参数:
//   - apiKey: Steam API Key
//
// 返回值:
//   - string: 脱敏后的API Key | Masked API Key
func MaskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return apiKey
	}

	return fmt.Sprintf("********%s", apiKey[len(apiKey)-8:])
}
