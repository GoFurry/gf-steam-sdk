// Package player 提供 Steam 玩家信息相关 API 封装
// 核心能力包括玩家基本信息查询、好友列表获取等, 支持批量 SteamID 处理（最大100个）
// Package player provides API encapsulation for Steam player information
// Core capabilities include player info query, friend list retrieval, supports batch SteamID processing (max 100)

package player

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/GoFurry/gf-steam-sdk/pkg/models"
	"github.com/GoFurry/gf-steam-sdk/pkg/util"
	"github.com/GoFurry/gf-steam-sdk/pkg/util/errors"
	"github.com/GoFurry/gf-steam-sdk/pkg/util/log"
	"github.com/bytedance/sonic"
)

const (
	// BASE_URL 玩家信息查询API基础地址 | Base URL for player info query API
	BASE_URL = "https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v2/"
)

// ============================ 原始字节流接口 ============================

// GetPlayerSummariesRawBytes 获取玩家信息的原始字节流
// 支持批量查询(最多100个SteamID), 返回原始API响应字节流, 适用于自定义解析场景
// 参数:
//   - steamIDs: 多个 SteamID 用逗号分隔(最多100个) | Multiple SteamIDs separated by commas (max 100)
//
// 返回值:
//   - []byte: 原始API响应字节流 | Raw API response bytes
//   - error: 请求/参数错误 | Request/parameter error
func (s *PlayerService) GetPlayerSummariesRawBytes(steamIDs string) (respBytes []byte, err error) {
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
	resp, err := s.client.DoRequest("GET", BASE_URL, params)
	if err != nil {
		log.Errorf("[PlayerService] GetPlayerSummariesRawBytes: Steam API 请求失败: %v", err)
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

// GetPlayerSummariesRawModel 获取玩家信息的结构化原始模型
// 解析为Steam官方定义的原始结构体, 保留所有返回字段
// 参数:
//   - steamIDs: 多个 SteamID 用逗号分隔 | Multiple SteamIDs separated by commas
//
// 返回值:
//   - models.SteamPlayerResponse: Steam原始响应结构体 | Steam raw response struct
//   - error: 请求/解析错误 | Request/parse error
func (s *PlayerService) GetPlayerSummariesRawModel(steamIDs string) (models.SteamPlayerResponse, error) {
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

// ============================ 精简模型接口 ============================

// GetPlayerSummariesBrief 获取玩家信息的精简模型
// 转换为业务友好的精简结构体, 补充格式化时间、在线状态布尔值等易用性字段
// 参数:
//   - steamIDs: 多个 SteamID 用逗号分隔 | Multiple SteamIDs separated by commas
//
// 返回值:
//   - []*models.Player: 精简玩家信息列表 | Simplified player info list
//   - error: 请求/解析错误 | Request/parse error
func (s *PlayerService) GetPlayerSummariesBrief(steamIDs string) ([]*models.Player, error) {

	// 获取原始结构化模型 | Get raw structured model
	rawPlayers, err := s.GetPlayerSummariesRawModel(steamIDs)
	if err != nil {
		return nil, err
	}

	// 转换为精简模型 | Convert to simplified model
	players := make([]*models.Player, 0, len(rawPlayers.Response.Players))
	for _, p := range rawPlayers.Response.Players {
		player := &models.Player{
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

// GetPlayerSummaries 精简模型接口的别名
// 简化调用方式，提供更直观的方法名
// GetPlayerSummaries is the alias of simplified model interface
// Simplifies calling with more intuitive method name
func (s *PlayerService) GetPlayerSummaries(steamIDs string) ([]*models.Player, error) {
	return s.GetPlayerSummariesBrief(steamIDs)
}
