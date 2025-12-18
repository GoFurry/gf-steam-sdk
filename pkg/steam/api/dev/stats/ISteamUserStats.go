// Package stats 提供 Steam 玩家统计数据相关 API 封装
// 核心能力包括成就查询、游戏时长统计等, 支持原始字节流/结构化模型/精简模型多层级返回
// Package stats provides API encapsulation for Steam player statistics
// Core capabilities include achievement query, playtime statistics, supports multi-level return (raw bytes/structured model/simplified model)

package stats

import (
	"fmt"
	"net/url"

	"github.com/GoFurry/gf-steam-sdk/pkg/models"
	"github.com/GoFurry/gf-steam-sdk/pkg/util"
	"github.com/GoFurry/gf-steam-sdk/pkg/util/errors"
	"github.com/bytedance/sonic"
)

const (
	ISteamUserStats = util.STEAM_API_BASE_URL + "ISteamUserStats"
)

// ============================ 原始字节流接口 ============================

// GetPlayerAchievementsRawBytes 获取玩家单游戏成就的原始字节流
// 适用于需要自定义解析、二次处理的场景, 保留API返回原始数据
// 参数:
//   - steamID: 玩家SteamID | Player SteamID
//   - appID: 游戏ID | Game AppID
//   - lang: 语言(比如zh/en) | Language (e.g. zh/en)
//
// 返回值:
//   - []byte: 原始API响应字节流 | Raw API response bytes
//   - error: 请求/解析错误 | Request/parse error
func (s *StatsService) GetPlayerAchievementsRawBytes(steamID string, appID uint64, lang string) (respBytes []byte, err error) {
	// 参数校验 | Parameter validation
	if steamID == "" {
		return respBytes, errors.ErrInvalidSteamID
	}
	if appID == 0 {
		return respBytes, errors.ErrInvalidAppID
	}
	// 默认语言 | Default language
	if lang == "" {
		lang = "en"
	}

	// 构建API请求参数 | Build API request parameters
	params := url.Values{}
	params.Set("steamid", steamID)
	params.Set("appid", util.Uint642String(appID))
	params.Set("l", lang) // 语言参数 | Language parameter

	// 调用Client发送请求 | Call Client to send request (auto trigger retry/proxy/rate limit)
	resp, err := s.client.DoRequest("GET", ISteamUserStats+"/GetPlayerAchievements/v1/", params)
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

// GetPlayerAchievementsRawModel 获取玩家单游戏成就的结构化原始模型
// 解析为Steam官方定义的原始结构体, 保留所有返回字段, 适用于全量数据场景
// 参数:
//   - steamID: 玩家SteamID | Player SteamID
//   - appID: 游戏ID | Game AppID
//   - lang: 语言 | Language
//
// 返回值:
//   - models.SteamPlayerAchievementsResponse: Steam原始响应结构体 | Steam raw response struct
//   - error: 请求/解析错误 | Request/parse error
func (s *StatsService) GetPlayerAchievementsRawModel(steamID string, appID uint64, lang string) (models.SteamPlayerAchievementsResponse, error) {
	// 获取原始字节流 | Get raw bytes
	bytes, err := s.GetPlayerAchievementsRawBytes(steamID, appID, lang)
	if err != nil {
		return models.SteamPlayerAchievementsResponse{}, err
	}

	// 解析为原始结构体 | Unmarshal to raw struct
	var statsResp models.SteamPlayerAchievementsResponse
	if err = sonic.Unmarshal(bytes, &statsResp); err != nil {
		return models.SteamPlayerAchievementsResponse{}, fmt.Errorf("%w: unmarshal achievements resp failed: %v", errors.ErrAPIResponse, err)
	}

	// 校验请求是否成功 | Validate request success (Steam API business status)
	if !statsResp.PlayerStats.Success {
		return models.SteamPlayerAchievementsResponse{}, errors.ErrAchievementFailed
	}

	return statsResp, nil
}

// ============================ 精简模型接口 ============================

// GetPlayerAchievementsBrief 获取玩家单游戏成就的精简模型
// 转换为业务友好的精简结构体, 剔除冗余字段, 补充格式化时间等易用性字段
// 参数:
//   - steamID: 玩家SteamID | Player SteamID
//   - appID: 游戏ID | Game AppID
//   - lang: 语言 | Language
//
// 返回值:
//   - []*models.PlayerAchievement: 精简成就信息列表 | Simplified achievement info list
//   - error: 请求/解析错误 | Request/parse error
func (s *StatsService) GetPlayerAchievementsBrief(steamID string, appID uint64, lang string) ([]*models.PlayerAchievement, error) {
	// 获取原始结构化模型 | Get raw structured model
	rawStats, err := s.GetPlayerAchievementsRawModel(steamID, appID, lang)
	if err != nil {
		return nil, err
	}

	// 转换为精简模型 | Convert to simplified model
	achievements := make([]*models.PlayerAchievement, 0, len(rawStats.PlayerStats.Achievements))
	for _, a := range rawStats.PlayerStats.Achievements {
		achievement := &models.PlayerAchievement{
			SteamID:         rawStats.PlayerStats.SteamID,
			GameName:        rawStats.PlayerStats.GameName,
			AppID:           appID,
			AchievementAPI:  a.APIName,
			AchievementName: a.Name,
			Achieved:        a.Achieved == 1, // 布尔化成就状态 | Booleanize achievement status
			UnlockTime:      a.UnlockTime,
			UnlockTimeStr:   util.TimeUnix2String(a.UnlockTime), // 格式化解锁时间 | Format unlock time
			Description:     a.Description,
		}
		achievements = append(achievements, achievement)
	}
	return achievements, nil
}

// GetPlayerAchievements 精简模型接口的别名
// 简化调用方式, 提供更直观的方法名
// GetPlayerAchievements is the alias of simplified model interface
// Simplifies calling with more intuitive method name
func (s *StatsService) GetPlayerAchievements(steamID string, appID uint64, lang string) ([]*models.PlayerAchievement, error) {
	return s.GetPlayerAchievementsBrief(steamID, appID, lang)
}
