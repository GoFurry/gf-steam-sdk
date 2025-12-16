package models

import "github.com/rumblefrog/go-a2s"

// SteamServerResponse 服务器数据模型
type SteamServerResponse struct {
	Server a2s.ServerInfo `json:"server"`
	Player a2s.PlayerInfo `json:"player"`
	Rules  a2s.RulesInfo  `json:"rules"`
}
