// Package game 提供 Steam 游戏信息相关 API 封装
// 核心能力包括已拥有游戏查询、游戏详情获取等, 支持免费游戏筛选、多平台时长统计
// Package game provides API encapsulation for Steam game information
// Core capabilities include owned games query, game details retrieval, supports free game filtering and multi-platform playtime statistics

package game

import (
	"fmt"
	"net/url"

	"github.com/GoFurry/gf-steam-sdk/pkg/models"
	"github.com/GoFurry/gf-steam-sdk/pkg/util"
	"github.com/GoFurry/gf-steam-sdk/pkg/util/errors"
	"github.com/bytedance/sonic"
)

const (
	// BASE_URL 玩家已拥有游戏查询API基础地址 | Base URL for player owned games query API
	BASE_URL = "https://api.steampowered.com/IPlayerService/GetOwnedGames/v1/"
)

// ============================ 原始字节流接口 ============================

// GetOwnedGamesRawBytes 获取玩家已拥有游戏的原始字节流
// 支持筛选免费游戏，返回原始API响应字节流，适用于自定义解析场景
// 参数:
//   - steamID: 玩家SteamID | Player SteamID
//   - includeFree: 是否包含免费游戏 | Whether to include free games
//
// 返回值:
//   - []byte: 原始API响应字节流 | Raw API response bytes
//   - error: 请求/参数错误 | Request/parameter error
func (s *GameService) GetOwnedGamesRawBytes(steamID string, includeFree bool) (respBytes []byte, err error) {
	// 参数校验 | Parameter validation
	if steamID == "" {
		return respBytes, errors.ErrInvalidSteamID
	}

	// 构建API请求参数 | Build API request parameters
	params := url.Values{}
	params.Set("steamid", steamID)
	params.Set("include_appinfo", "1")                                              // 包含游戏名称/图标 | Include game name/icon
	params.Set("include_extended_appinfo", "1")                                     // 包含扩展信息 | Include extended info
	params.Set("include_played_free_games", util.Int2String(util.B2i(includeFree))) // 包含免费游戏 | Include free games

	// 调用Client发送请求 | Call Client (auto apply proxy/rate limit from chain config)
	resp, err := s.client.DoRequest("GET", BASE_URL, params)
	if err != nil {
		return respBytes, err
	}

	// 转换为字节流返回 | Convert to bytes and return
	respBytes, err = sonic.Marshal(resp)
	if err != nil {
		return respBytes, fmt.Errorf("%w: marshal resp failed: %v", errors.ErrAPIResponse, err)
	}

	return respBytes, nil
}

// ============================ 结构化原始模型接口 ============================

// GetOwnedGamesRawModel 获取玩家已拥有游戏的结构化原始模型
// 解析为Steam官方定义的原始结构体, 保留所有返回字段
// 参数:
//   - steamID: 玩家SteamID | Player SteamID
//   - includeFree: 是否包含免费游戏 | Whether to include free games
//
// 返回值:
//   - models.SteamOwnedGamesResponse: Steam原始响应结构体 | Steam raw response struct
//   - error: 请求/解析错误 | Request/parse error
func (s *GameService) GetOwnedGamesRawModel(steamID string, includeFree bool) (models.SteamOwnedGamesResponse, error) {
	// 获取原始字节流 | Get raw bytes
	bytes, err := s.GetOwnedGamesRawBytes(steamID, includeFree)
	if err != nil {
		return models.SteamOwnedGamesResponse{}, err
	}

	// 解析为原始结构体 | Unmarshal to raw struct
	var gamesResp models.SteamOwnedGamesResponse
	if err = sonic.Unmarshal(bytes, &gamesResp); err != nil {
		return models.SteamOwnedGamesResponse{}, fmt.Errorf("%w: unmarshal owned games resp failed: %v", errors.ErrAPIResponse, err)
	}

	return gamesResp, nil
}

// ============================ 精简模型接口 ============================

// GetOwnedGamesBrief 获取玩家已拥有游戏的精简模型
// 转换为业务友好的精简结构体, 补充游戏图标/封面URL、格式化时间等易用性字段
// 参数:
//   - steamID: 玩家SteamID | Player SteamID
//   - includeFree: 是否包含免费游戏 | Whether to include free games
//
// 返回值:
//   - []*models.OwnedGame: 精简游戏信息列表 | Simplified game info list
//   - error: 请求/解析错误 | Request/parse error
func (s *GameService) GetOwnedGamesBrief(steamID string, includeFree bool) ([]*models.OwnedGame, error) {
	// 获取原始结构化模型 | Get raw structured model
	rawGames, err := s.GetOwnedGamesRawModel(steamID, includeFree)
	if err != nil {
		return nil, err
	}

	// 转换为精简模型 | Convert to simplified model
	games := make([]*models.OwnedGame, 0, len(rawGames.Response.Games))
	for _, g := range rawGames.Response.Games {
		game := &models.OwnedGame{
			AppID:                  g.AppID,
			Name:                   g.Name,
			PlaytimeForever:        g.PlaytimeForever,
			Playtime2Weeks:         g.Playtime2Weeks,
			IconURL:                fmt.Sprintf(util.STEAM_ICON_URL, g.AppID, g.ImgIconURL), // 拼接图标URL | Splice icon URL
			CapsuleURL:             fmt.Sprintf(util.STEAM_CAPSULE_URL, g.AppID),            // 拼接封面URL | Splice capsule URL
			LastPlayedTime:         g.RTimeLastPlayed,
			LastPlayedTimeStr:      util.TimeUnix2String(g.RTimeLastPlayed), // 格式化最后游玩时间 | Format last played time
			HasCommunityVisible:    g.HasCommunityVisible,
			PlaytimeWindowsForever: g.PlaytimeWindowsForever,
			PlaytimeDeckForever:    g.PlaytimeDeckForever,
			HasDLC:                 g.HasDLC,
		}
		games = append(games, game)
	}

	return games, nil
}

// GetOwnedGames 精简模型接口的别名
// 简化调用方式，提供更直观的方法名
// GetOwnedGames is the alias of simplified model interface
// Simplifies calling with more intuitive method name
func (s *GameService) GetOwnedGames(steamID string, includeFree bool) ([]*models.OwnedGame, error) {
	return s.GetOwnedGamesBrief(steamID, includeFree)
}
