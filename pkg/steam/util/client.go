// Package util 提供 Steam 开发工具类服务
// 包含各类开发者常用的工具方法, 如 Steam 令牌/API Key 获取等
// 辅助开发者快速获取 Steam 接口调用所需的认证信息
// Package util provides Steam development tooling services
// Includes various commonly used utility methods for developers, such as Steam token/API Key acquisition, etc.,
// helping developers quickly obtain authentication information required for Steam interface calls
package util

import "github.com/GoFurry/gf-steam-sdk/internal/client"

// UtilService Steam 工具服务核心结构体
// 主要用于辅助获取 Steam 接口调用所需的各类令牌和密钥
// UtilService is the core structure of Steam utility service
// mainly used to assist in obtaining various tokens and keys required for Steam interface calls
type UtilService struct {
	client *client.Client // 内部通信客户端（Internal communication client）
}

// NewUtilService 创建UtilService实例, 暴露初始化入口
func NewUtilService(c *client.Client) *UtilService {
	return &UtilService{client: c}
}
