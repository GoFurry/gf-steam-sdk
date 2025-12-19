package steam

import (
	"fmt"

	"github.com/GoFurry/gf-steam-sdk/internal/client"
	"github.com/GoFurry/gf-steam-sdk/pkg/config"
	"github.com/GoFurry/gf-steam-sdk/pkg/steam/api/dev"
	"github.com/GoFurry/gf-steam-sdk/pkg/steam/api/store"
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
	Develop *dev.DevService         // 玩家模块 | Develop module API from api.steampowered.com
	Store   *store.StoreService     // 商店模块 | Store module (接口数据来自商店界面) API from store.steampowered.com
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
	// 兜底AccessToken | Fallback AccessToken (avoid request failure due to empty value)
	if cfg.AccessToken == "" {
		cfg.AccessToken = "steam-access-token"
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
		Develop: dev.NewDevService(cli),
		Crawler: crawler.NewCrawlerService(cfg),
		Server:  server.NewServerService(cli),
		Util:    util.NewUtilService(cli),
		Store:   store.NewStoreService(cli),
	}, nil
}

// Close 释放所有模块资源
func (s *SteamSDK) Close() error {
	defer func() {
		s.Develop, s.Store, s.Crawler, s.Server, s.Util = nil, nil, nil, nil, nil
	}()

	s.Develop.Close()
	s.Store.Close()
	s.Server.Close()
	s.Util.Close()
	return nil
}
