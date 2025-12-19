// pkg/models/player_models.go

package models

// SteamPlayerResponse Steam 原始响应结构体 ISteamUser/GetPlayerSummaries
type SteamPlayerResponse struct {
	Response struct {
		Players []struct {
			SteamID                  string `json:"steamid"`                  // Steam 账号唯一 ID
			CommunityVisibilityState int    `json:"communityvisibilitystate"` // 社区可见性
			ProfileState             int    `json:"profilestate"`
			PersonaName              string `json:"personaname"`       // 用户显示名称
			CommentPermission        int    `json:"commentpermission"` // 评论权限: 1=允许, 0=禁止
			ProfileURL               string `json:"profileurl"`        // 个人主页地址
			Avatar                   string `json:"avatar"`            // 小尺寸头像URL
			AvatarMedium             string `json:"avatarmedium"`      // 半尺寸头像URL
			AvatarFull               string `json:"avatarfull"`        // 全尺寸头像URL
			AvatarHash               string `json:"avatarhash"`        // 头像哈希值
			LastLogoff               int64  `json:"lastlogoff"`        // 最后离线时间戳
			PersonaState             int    `json:"personastate"`      // 在线状态: 0=离线 1=在线
			RealName                 string `json:"realname"`          // 真实姓名(用户填写)
			PrimaryClanID            string `json:"primaryclanid"`     // 主要所属群组 ID
			TimeCreated              int64  `json:"timecreated"`       // 账号创建时间戳
			PersonaStateFlags        int    `json:"personastateflags"`
			LocCountryCode           string `json:"loccountrycode"` // 国家码
			LocStateCode             string `json:"locstatecode"`   // 省份/州码
			LocCityID                int    `json:"loccityid"`      // 城市ID
		} `json:"players"`
	} `json:"response"`
}

// Player 用户基本信息精简模型
type Player struct {
	SteamID      string `json:"steam_id"`      // Steam唯一ID
	PersonaName  string `json:"persona_name"`  // 昵称
	ProfileURL   string `json:"profile_url"`   // 个人主页链接
	AvatarURL    string `json:"avatar_url"`    // 小尺寸头像
	AvatarMedium string `json:"avatar_medium"` // 半尺寸头像
	AvatarFull   string `json:"avatar_full"`   // 全尺寸头像
	TimeCreated  string `json:"time_created"`  // 账号创建时间
	LastLogoff   string `json:"last_logoff"`   // 最后离线时间
	IsOnline     bool   `json:"is_online"`     // 是否在线
	RealName     string `json:"real_name"`     // 真实姓名
	CountryCode  string `json:"country_code"`  // 国家码
}

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
