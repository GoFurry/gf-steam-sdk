// Package dev provides high-level encapsulation for Steam's store page APIs
// Package dev 提供 Steam 商店页面 API 的上层封装

package store

import (
	"fmt"

	"github.com/GoFurry/gf-steam-sdk/internal/client"
)

// StoreService Steam官方商店页面API核心结构体
type StoreService struct {
	client *client.Client
}

// NewStoreService 创建StoreService实例, 暴露初始化入口
func NewStoreService(c *client.Client) *StoreService {
	return &StoreService{client: c}
}

// Close 释放StoreService资源
func (s *StoreService) Close() error {
	if s.client == nil {
		return nil
	}
	if err := s.client.Close(); err != nil {
		return fmt.Errorf("StoreService Close failed: %w", err)
	}
	return nil
}

// Client 对外暴露 client
func (s *StoreService) Client() *client.Client {
	return s.client
}
