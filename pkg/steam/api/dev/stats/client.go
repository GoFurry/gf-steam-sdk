// Package stats 提供 Steam 玩家统计数据相关 API 封装
// 核心能力包括成就查询、游戏时长统计等, 支持原始字节流/结构化模型/精简模型多层级返回
// Package stats provides API encapsulation for Steam player statistics
// Core capabilities include achievement query, playtime statistics, supports multi-level return (raw bytes/structured model/simplified model)

package stats

import (
	"github.com/GoFurry/gf-steam-sdk/internal/client"
)

// StatsService 统计相关API服务核心结构体
// 依赖内部 Client 实现请求发送, 封装 Steam 玩家成就、统计数据等核心接口
// StatsService is the core structure for statistics-related API services
// Depends on internal Client for request sending, encapsulates core interfaces for Steam player achievements/statistics
type StatsService struct {
	client *client.Client
}

// NewStatsService 创建StatsService实例, 暴露初始化入口
// 参数:
//   - c: 内部Client实例 | Internal Client instance
//
// 返回值:
//   - *StatsService: 统计服务实例 | Statistics service instance
func NewStatsService(c *client.Client) *StatsService {
	return &StatsService{client: c}
}
