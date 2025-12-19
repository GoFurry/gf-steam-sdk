package dev

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/GoFurry/gf-steam-sdk/pkg/models"
	"github.com/GoFurry/gf-steam-sdk/pkg/util"
	"github.com/GoFurry/gf-steam-sdk/pkg/util/errors"
	"github.com/bytedance/sonic"
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

	// 构建 API 请求参数 | Build API request parameters
	params := url.Values{}
	params.Set("steamids", steamIDs)

	// 调用内部 Client 发送请求 | Call internal Client (auto apply proxy/rate limit/retry)
	resp, err := s.client.DoRequest("GET", ISteamUser+"/GetPlayerSummaries/v2/", params)
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

// GetPlayerSummariesRawModel get player's information 获取玩家信息
//   - steamIDs: Multiple SteamIDs separated by commas (max 100)
func (s *DevService) GetPlayerSummariesRawModel(steamIDs string) (models.SteamPlayerResponse, error) {
	// 获取原始字节流 | Get raw bytes
	bytes, err := s.GetPlayerSummariesRawBytes(steamIDs)
	if err != nil {
		return models.SteamPlayerResponse{}, err
	}

	// 解析为原始结构体 | Unmarshal to raw struct
	var steamResp models.SteamPlayerResponse
	if err = sonic.Unmarshal(bytes, &steamResp); err != nil {
		return models.SteamPlayerResponse{}, fmt.Errorf("%w: unmarshal player resp failed: %v", errors.ErrAPIResponse, err)
	}

	return steamResp, nil
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
