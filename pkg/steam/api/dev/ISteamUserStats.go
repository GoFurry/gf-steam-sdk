package dev

import (
	"net/url"

	"github.com/GoFurry/gf-steam-sdk/internal/api"
	"github.com/GoFurry/gf-steam-sdk/internal/client"
	"github.com/GoFurry/gf-steam-sdk/pkg/models"
	"github.com/GoFurry/gf-steam-sdk/pkg/util"
)

const (
	ISteamUserStats = util.STEAM_API_BASE_URL + "ISteamUserStats"
)

// ============================ Structed Raw Model Raw Bytes 原始字节流接口 ============================

// GetPlayerAchievementsRawBytes get player's game achievements 获取玩家单游戏成就
//   - steamID: Player SteamID
//   - appID: Game AppID
//   - lang: Language (e.g. zh/en)
func (s *DevService) GetPlayerAchievementsRawBytes(steamID string, appID uint64, lang string) (respBytes []byte, err error) {
	return api.GetRawBytes(s.buildPlayerAchievements(steamID, appID, lang))
}

// ============================ 结构化原始模型接口 ============================

// GetPlayerAchievementsRawModel get player's game achievements 获取玩家单游戏成就
//   - steamID: Player SteamID
//   - appID: Game AppID
//   - lang: Language (e.g. zh/en)
func (s *DevService) GetPlayerAchievementsRawModel(steamID string, appID uint64, lang string) (models.SteamPlayerAchievementsResponse, error) {
	return api.GetRawModel[models.SteamPlayerAchievementsResponse](s.buildPlayerAchievements(steamID, appID, lang))
}

// ============================ Brief Model 精简模型接口 ============================

// GetPlayerAchievementsBrief get player's game achievements 获取玩家单游戏成就
//   - steamID: Player SteamID
//   - appID: Game AppID
//   - lang: Language (e.g. zh/en)
func (s *DevService) GetPlayerAchievementsBrief(steamID string, appID uint64, lang string) ([]models.PlayerAchievement, error) {
	// 获取原始结构化模型 | Get raw structured model
	rawStats, err := s.GetPlayerAchievementsRawModel(steamID, appID, lang)
	if err != nil {
		return nil, err
	}

	// 转换为精简模型 | Convert to simplified model
	achievements := make([]models.PlayerAchievement, 0, len(rawStats.PlayerStats.Achievements))
	for _, a := range rawStats.PlayerStats.Achievements {
		achievement := models.PlayerAchievement{
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

// ============================ Default Interface 默认接口 ============================

// GetPlayerAchievements get player's game achievements 获取玩家单游戏成就
//   - steamID: Player SteamID
//   - appID: Game AppID
//   - lang: Language (e.g. zh/en)
func (s *DevService) GetPlayerAchievements(steamID string, appID uint64, lang string) ([]models.PlayerAchievement, error) {
	return s.GetPlayerAchievementsBrief(steamID, appID, lang)
}

// ============================ Build 构造入参 ============================

// buildPlayerSummaries builds input params.
func (s *DevService) buildPlayerAchievements(steamID string, appID uint64, lang string) (
	c *client.Client,
	method, reqPath string,
	params url.Values,
) {
	params = url.Values{}
	params.Set("steamid", steamID)
	params.Set("appid", util.Uint642String(appID))
	params.Set("l", lang)
	return s.client, "GET", ISteamUserStats + "/GetPlayerAchievements/v1/", params
}
