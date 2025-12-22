package dev

import (
	"net/url"
	"strings"

	"github.com/GoFurry/gf-steam-sdk/internal/api"
	"github.com/GoFurry/gf-steam-sdk/internal/client"
	"github.com/GoFurry/gf-steam-sdk/pkg/models"
	"github.com/GoFurry/gf-steam-sdk/pkg/util"
	"github.com/GoFurry/gf-steam-sdk/pkg/util/errors"
)

const (
	ISteamUser = util.STEAM_API_BASE_URL + "ISteamUser"
)

// ============================ Raw Bytes 原始字节流接口 ============================

// GetPlayerSummariesRawBytes get player's information 获取玩家信息
//   - steamIDs: Multiple SteamIDs separated by commas (max 100)
func (s *DevService) GetPlayerSummariesRawBytes(steamIDs string) (respBytes []byte, err error) {
	// 参数校验 | Parameter validation
	if steamIDs == "" {
		return respBytes, errors.ErrInvalidSteamID
	}
	ids := strings.Split(steamIDs, ",")
	if len(ids) > 100 {
		return respBytes, errors.NewWithType(errors.ErrTypeParam,
			"steamids count exceeds 100 (Steam API maximum limit)", nil)
	}

	return api.GetRawBytes(s.buildPlayerSummaries(steamIDs))
}

// ============================ Structed Raw Model 结构化原始模型接口 ============================

// GetPlayerSummariesRawModel get player's information 获取玩家信息
//   - steamIDs: Multiple SteamIDs separated by commas (max 100)
func (s *DevService) GetPlayerSummariesRawModel(steamIDs string) (models.SteamPlayerResponse, error) {
	// 参数校验 | Parameter validation
	if steamIDs == "" {
		return models.SteamPlayerResponse{}, errors.ErrInvalidSteamID
	}
	ids := strings.Split(steamIDs, ",")
	if len(ids) > 100 {
		return models.SteamPlayerResponse{}, errors.NewWithType(errors.ErrTypeParam,
			"steamids count exceeds 100 (Steam API maximum limit)", nil)
	}

	return api.GetRawModel[models.SteamPlayerResponse](s.buildPlayerSummaries(steamIDs))
}

// ============================ Brief Model 精简模型接口 ============================

// GetPlayerSummariesBrief get player's information 获取玩家信息
//   - steamIDs: Multiple SteamIDs separated by commas (max 100)
func (s *DevService) GetPlayerSummariesBrief(steamIDs string) ([]models.Player, error) {

	// 获取原始结构化模型 | Get raw structured model
	rawPlayers, err := s.GetPlayerSummariesRawModel(steamIDs)
	if err != nil {
		return nil, err
	}

	// 转换为精简模型 | Convert to simplified model
	players := make([]models.Player, 0, len(rawPlayers.Response.Players))
	for _, p := range rawPlayers.Response.Players {
		player := models.Player{
			SteamID:      p.SteamID,
			PersonaName:  p.PersonaName,
			ProfileURL:   p.ProfileURL,
			AvatarURL:    p.Avatar,
			AvatarMedium: p.AvatarMedium,
			AvatarFull:   p.AvatarFull,
			LastLogoff:   util.TimeUnix2String(p.LastLogoff), // 格式化最后登录时间 | Format last logoff time
			RealName:     p.RealName,
			CountryCode:  p.LocCountryCode,
			TimeCreated:  util.TimeUnix2String(p.TimeCreated), // 格式化账号创建时间 | Format account creation time
			IsOnline:     p.PersonaState != 0,                 // 布尔化在线状态 | Booleanize online status
		}
		players = append(players, player)
	}

	return players, nil
}

// ============================ Default Interface 默认接口 ============================

// GetPlayerSummaries get player's information 获取玩家信息
//   - steamIDs: Multiple SteamIDs separated by commas (max 100)
func (s *DevService) GetPlayerSummaries(steamIDs string) ([]models.Player, error) {
	return s.GetPlayerSummariesBrief(steamIDs)
}

// ============================ Build 构造入参 ============================

// buildPlayerSummaries builds input params.
func (s *DevService) buildPlayerSummaries(steamIDs string) (
	c *client.Client,
	method, reqPath string,
	params url.Values,
) {
	params = url.Values{}
	params.Set("steamids", steamIDs)
	return s.client, "GET", ISteamUser + "/GetPlayerSummaries/v2/", params
}
