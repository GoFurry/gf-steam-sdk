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
