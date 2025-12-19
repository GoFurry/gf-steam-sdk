package models

// SteamOwnedGamesResponse Steam 原始响应结构体 IPlayerService/GetOwnedGames
type SteamOwnedGamesResponse struct {
	Response struct {
		GameCount int `json:"game_count"` // 游戏总数
		Games     []struct {
			AppID                  uint64 `json:"appid"`                       // 游戏唯一ID
			Name                   string `json:"name"`                        // 游戏名称
			Playtime2Weeks         int    `json:"playtime_2weeks"`             // 近2周游玩时长(分钟)
			PlaytimeForever        int    `json:"playtime_forever"`            // 总游玩时长(分钟)
			ImgIconURL             string `json:"img_icon_url"`                // 游戏图标哈希(小图标)
			HasCommunityVisible    bool   `json:"has_community_visible_stats"` // 是否显示社区统计
			PlaytimeWindowsForever int    `json:"playtime_windows_forever"`    // Windows端总时长
			PlaytimeMacForever     int    `json:"playtime_mac_forever"`        // Mac端总时长
			PlaytimeLinuxForever   int    `json:"playtime_linux_forever"`      // Linux端总时长
			PlaytimeDeckForever    int    `json:"playtime_deck_forever"`       // SteamDeck端总时长
			RTimeLastPlayed        int64  `json:"rtime_last_played"`           // 最后游玩时间戳
			CapsuleFilename        string `json:"capsule_filename"`            // 封面图文件名
			HasWorkshop            bool   `json:"has_workshop"`                // 是否有创意工坊
			HasMarket              bool   `json:"has_market"`                  // 是否有市场
			HasDLC                 bool   `json:"has_dlc"`                     // 是否有DLC
			ContentDescriptorIDs   []int  `json:"content_descriptorids"`
			PlaytimeDisconnected   int    `json:"playtime_disconnected"` // 离线游玩时长
		} `json:"games"`
	} `json:"response"`
}

// OwnedGame 玩家已拥有游戏信息精简模型
type OwnedGame struct {
	AppID                  uint64 `json:"app_id"`               // 游戏ID
	Name                   string `json:"name"`                 // 游戏名称
	PlaytimeForever        int    `json:"playtime_forever"`     // 总游玩时长(分钟)
	Playtime2Weeks         int    `json:"playtime_2weeks"`      // 近2周游玩时长(分钟)
	IconURL                string `json:"icon_url"`             // 游戏图标完整URL
	CapsuleURL             string `json:"capsule_url"`          // 游戏封面完整URL
	LastPlayedTime         int64  `json:"last_played_time"`     // 最后游玩时间戳
	LastPlayedTimeStr      string `json:"last_played_time_str"` // 最后游玩时间格式化字符串
	HasCommunityVisible    bool   `json:"has_community_visb"`   // 社区可见性
	PlaytimeWindowsForever int    `json:"playtime_windows"`     // Windows端时长
	PlaytimeDeckForever    int    `json:"playtime_deck"`        // SteamDeck端时长
	HasDLC                 bool   `json:"has_dlc"`              // 是否有DLC
}
