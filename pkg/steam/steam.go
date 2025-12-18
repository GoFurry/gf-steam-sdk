// Package steam 提供 Steam SDK 核心入口及模块管理能力
// 整合玩家、游戏、统计、爬虫四大核心模块，基于模块化架构设计, 支持 Steam Web API 多场景调用
// Package steam provides core entry and module management capabilities for Steam SDK
// Integrates four core modules (Player/Game/Stats/Crawler), designed with modular architecture, supports full-scenario Steam Web API calls

package steam

import (
	"fmt"

	"github.com/GoFurry/gf-steam-sdk/internal/client"
	"github.com/GoFurry/gf-steam-sdk/pkg/config"
	"github.com/GoFurry/gf-steam-sdk/pkg/steam/api/dev/game"
	"github.com/GoFurry/gf-steam-sdk/pkg/steam/api/dev/player"
	"github.com/GoFurry/gf-steam-sdk/pkg/steam/api/dev/stats"
	"github.com/GoFurry/gf-steam-sdk/pkg/steam/crawler"
	"github.com/GoFurry/gf-steam-sdk/pkg/steam/server"
	"github.com/GoFurry/gf-steam-sdk/pkg/steam/util"
)

// 文档参考 | Documentation references:
// 	https://developer.valvesoftware.com/wiki/Steam_Web_API
// 	https://partner.steamgames.com/doc/webapi_overview
// 	https://steamapi.xpaw.me/

// SteamSDK 全局 Steam SDK 入口结构体
// 聚合所有核心业务模块, 提供统一的 SDK 调用入口, 支持链式配置扩展
// SteamSDK is the global entry structure of Steam SDK
// Aggregates all core business modules, provides a unified SDK call entry, supports chain configuration extension
// SteamSDK 全局 Steam SDK 入口
type SteamSDK struct {
	Player  *player.PlayerService   // 玩家模块 | Player module (SteamID/个人信息/好友等)
	Game    *game.GameService       // 游戏模块 | Game module (已拥有游戏/游戏信息等)
	Stats   *stats.StatsService     // 统计模块 | Stats module (成就/游戏时长等)
	Crawler *crawler.CrawlerService // 爬虫模块 | Crawler module (网页爬取/反爬策略)
	Server  *server.ServerService   // 服务器模块 | Server module (集成官方指定A2S库)
	Util    *util.UtilService       // 工具模块 | Util module (一些可能会用到的工具函数)
}

// NewSteamSDK 创建全局 Steam SDK 实例
// 初始化内部 Client 并完成所有业务模块的实例化，支持配置校验和错误兜底
// 参数:
//   - cfg: 全局配置(支持链式配置自定义代理、限流、超时等) | Global config (supports chain config for proxy/rate limit/timeout)
//
// 返回值:
//   - *SteamSDK: SDK 实例 | SDK instance
//   - error: 初始化失败错误 | Initialization error
func NewSteamSDK(cfg *config.SteamConfig) (*SteamSDK, error) {
	// 兜底API Key | Fallback API Key (avoid request failure due to empty value)
	if cfg.APIKey == "" {
		cfg.APIKey = "steam-api-key"
	}
	// 创建内部 Client | Create internal Client (integrates retry/proxy/rate limit)
	cli, err := client.NewClient(cfg)
	if err != nil {
		if cfg.IsDebug {
			fmt.Printf("[Error] NewSteamSDK NewClient error: %v\n", err)
		}
		return nil, err
	}

	// 初始化所有模块 Service | Initialize all module services
	return &SteamSDK{
		Player:  player.NewPlayerService(cli),
		Game:    game.NewGameService(cli),
		Stats:   stats.NewStatsService(cli),
		Crawler: crawler.NewCrawlerService(cfg),
		Server:  server.NewServerService(cli),
		Util:    util.NewUtilService(cli),
	}, nil
}
