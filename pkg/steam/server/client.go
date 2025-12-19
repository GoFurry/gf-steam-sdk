// Package server 提供 Steam 服务器(A2S 协议)相关查询服务
// 包含单个/批量服务器基础信息、玩家信息、规则信息的查询能力，
// 支持限流、重试、超时控制，保证线程安全, 底层基于 go-a2s 库实现
// Package server provides query services related to Steam servers (A2S protocol)
// It includes the ability to query single/batch server basic info, player info, rule info,
// supports rate limiting, retry, timeout control, ensures thread safety, and is implemented based on the go-a2s library at the bottom
package server

import (
	"fmt"

	"github.com/GoFurry/gf-steam-sdk/internal/client"
)

// ServerService Steam服务器查询服务核心结构体
// 封装所有A2S协议相关的服务器信息查询方法，通过内部client管理底层通信
// ServerService is the core structure of Steam server query service
// Encapsulates all A2S protocol-related server information query methods, and manages underlying communication through internal client
type ServerService struct {
	client *client.Client
}

// NewServerService 创建ServerService实例
// NewServerService creates a ServerService instance
// Initializes the server query service, relying on the internal client to provide underlying communication capabilities
func NewServerService(c *client.Client) *ServerService {
	return &ServerService{client: c}
}

// Close 释放ServerService资源
func (s *ServerService) Close() error {
	if s.client == nil {
		return nil
	}
	if err := s.client.Close(); err != nil {
		return fmt.Errorf("ServerService Close failed: %w", err)
	}
	return nil
}
