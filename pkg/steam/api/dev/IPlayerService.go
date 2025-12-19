package dev

import (
	"fmt"
	"net/url"

	"github.com/GoFurry/gf-steam-sdk/pkg/models"
	"github.com/GoFurry/gf-steam-sdk/pkg/util"
	"github.com/GoFurry/gf-steam-sdk/pkg/util/errors"
	"github.com/bytedance/sonic"
)

const (
	IPlayerService = util.STEAM_API_BASE_URL + "IPlayerService"
)

// ============================ Raw Bytes 原始字节流接口 ============================

// GetOwnedGamesRawBytes get player's owned games 获取玩家已拥有的游戏
//   - steamID: Player SteamID
//   - includeFree: Whether to include free games
func (s *DevService) GetOwnedGamesRawBytes(steamID string, includeFree bool) (respBytes []byte, err error) {
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
	resp, err := s.client.DoRequest("GET", IPlayerService+"/GetOwnedGames/v1/", params)
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

// ============================ Structed Raw Model 结构化原始模型接口 ============================

// GetOwnedGamesRawModel get player's owned games 获取玩家已拥有的游戏
//   - steamID: Player SteamID
//   - includeFree: Whether to include free games
func (s *DevService) GetOwnedGamesRawModel(steamID string, includeFree bool) (models.SteamOwnedGamesResponse, error) {
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

// ============================ Brief Model 精简模型接口 ============================

// GetOwnedGamesBrief get player's owned games 获取玩家已拥有的游戏
//   - steamID: Player SteamID
//   - includeFree: Whether to include free games
func (s *DevService) GetOwnedGamesBrief(steamID string, includeFree bool) ([]models.OwnedGame, error) {
	// 获取原始结构化模型 | Get raw structured model
	rawGames, err := s.GetOwnedGamesRawModel(steamID, includeFree)
	if err != nil {
		return nil, err
	}

	// 转换为精简模型 | Convert to simplified model
	games := make([]models.OwnedGame, 0, len(rawGames.Response.Games))
	for _, g := range rawGames.Response.Games {
		game := models.OwnedGame{
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

// ============================ Default Interface 默认接口 ============================

// GetOwnedGames get player's owned games 获取玩家已拥有的游戏
//   - steamID: Player SteamID
//   - includeFree: Whether to include free games
func (s *DevService) GetOwnedGames(steamID string, includeFree bool) ([]models.OwnedGame, error) {
	return s.GetOwnedGamesBrief(steamID, includeFree)
}
