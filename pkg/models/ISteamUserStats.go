package models

// SteamPlayerAchievementsResponse ISteamUserStats/GetPlayerAchievements
type SteamPlayerAchievementsResponse struct {
	PlayerStats struct {
		SteamID      string `json:"steamID"`  // 玩家SteamID
		GameName     string `json:"gameName"` // 游戏名称
		Achievements []struct {
			APIName     string `json:"apiname"`     // 成就唯一标识
			Achieved    int    `json:"achieved"`    // 是否完成: 1=完成, 0=未完成
			UnlockTime  int64  `json:"unlocktime"`  // 解锁时间戳
			Name        string `json:"name"`        // 成就名称
			Description string `json:"description"` // 成就描述
		} `json:"achievements"`
		Success bool `json:"success"` // 请求是否成功
	} `json:"playerstats"`
}

// PlayerAchievement 玩家游戏成就精简模型
type PlayerAchievement struct {
	SteamID         string `json:"steam_id"`         // 玩家SteamID
	GameName        string `json:"game_name"`        // 游戏名称
	AppID           uint64 `json:"app_id"`           // 游戏ID
	AchievementAPI  string `json:"achievement_api"`  // 成就唯一标识
	AchievementName string `json:"achievement_name"` // 成就名称
	Achieved        bool   `json:"achieved"`         // 是否完成
	UnlockTime      int64  `json:"unlock_time"`      // 解锁时间戳
	UnlockTimeStr   string `json:"unlock_time_str"`  // 解锁时间格式化字符串
	Description     string `json:"description"`      // 成就描述
}
