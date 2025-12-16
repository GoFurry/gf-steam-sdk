// Package player 提供 Steam 玩家信息相关 API 封装
// 核心能力包括玩家基本信息查询、好友列表获取等, 支持批量 SteamID 处理(最大100个)
// Package player provides API encapsulation for Steam player information
// Core capabilities include player info query, friend list retrieval, supports batch SteamID processing (max 100)

package player

import (
	"github.com/GoFurry/gf-steam-sdk/internal/client"
)

// PlayerService 用户相关API服务核心结构体
// 依赖内部 Client 实现请求发送，封装 Steam 玩家信息查询类接口
// PlayerService is the core structure for player-related API services
// Depends on internal Client for request sending, encapsulates Steam player info query interfaces
type PlayerService struct {
	client *client.Client
}

// NewPlayerService 创建 PlayerService 实例, 暴露初始化入口
// 参数:
//   - c: 内部Client实例 | Internal Client instance
//
// 返回值:
//   - *PlayerService: 玩家服务实例 | Player service instance
func NewPlayerService(c *client.Client) *PlayerService {
	return &PlayerService{client: c}
}
