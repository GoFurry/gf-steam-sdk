// Package game 提供 Steam 游戏信息相关 API 封装
// 核心能力包括已拥有游戏查询、游戏详情获取等, 支持免费游戏筛选、多平台时长统计
// Package game provides API encapsulation for Steam game information
// Core capabilities include owned games query, game details retrieval, supports free game filtering and multi-platform playtime statistics

package game

import (
	"github.com/GoFurry/gf-steam-sdk/internal/client"
)

// GameService 游戏相关API服务核心结构体
// 依赖内部 Client 实现请求发送, 封装 Steam 游戏信息查询类接口
// GameService is the core structure for game-related API services
// Depends on internal Client for request sending, encapsulates Steam game info query interfaces
type GameService struct {
	client *client.Client
}

// NewGameService 创建GameService实例, 暴露初始化入口
// 参数:
//   - c: 内部Client实例 | Internal Client instance
//
// 返回值:
//   - *GameService: 游戏服务实例 | Game service instance
func NewGameService(c *client.Client) *GameService {
	return &GameService{client: c}
}
