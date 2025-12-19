// Package dev provides high-level encapsulation for Steam's official core APIs
// Package dev 提供 Steam 官方核心 API 的上层封装

package dev

import (
	"fmt"

	"github.com/GoFurry/gf-steam-sdk/internal/client"
)

// DevService Steam官方API核心结构体
type DevService struct {
	client *client.Client
}

// NewDevService 创建DevService实例, 暴露初始化入口
func NewDevService(c *client.Client) *DevService {
	return &DevService{client: c}
}

// Close 释放DevService资源
func (g *DevService) Close() error {
	if g.client == nil {
		return nil
	}
	if err := g.client.Close(); err != nil {
		return fmt.Errorf("DevService Close failed: %w", err)
	}
	return nil
}
