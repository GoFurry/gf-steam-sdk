package models

// ==================== GetChangeLog ====================

// FamilyGroupChange 家庭组单条变更记录结构体
type FamilyGroupChange struct {
	Timestamp    string `json:"timestamp"`     // 变更时间戳
	ActorSteamID string `json:"actor_steamid"` // 操作人SteamID
	Type         int    `json:"type"`          // 变更类型标识
	Body         string `json:"body"`          // 变更详情JSON字符串
	BySupport    bool   `json:"by_support"`    // 是否由Steam客服执行的操作
}

// FamilyGroupChangeLogResponse IFamilyGroupsService/GetChangeLog
type FamilyGroupChangeLogResponse struct {
	Response struct {
		Changes []FamilyGroupChange `json:"changes"` // 家庭组变更记录列表
	} `json:"response"` // 接口响应主体
}

// ==================== GetFamilyGroup ====================

// FamilyGroupMember 家庭组成员信息结构体
// 对应接口返回的 members 数组中的单个成员项
type FamilyGroupMember struct {
	SteamID                  string `json:"steamid"`                    // 成员ID
	Role                     int    `json:"role"`                       // 成员角色标识
	TimeJoined               int64  `json:"time_joined"`                // 加入家庭组的时间戳
	CooldownSecondsRemaining int64  `json:"cooldown_seconds_remaining"` // 成员剩余冷却秒数
}

// FamilyGroupResponse IFamilyGroupsService/GetFamilyGroup
// 根结构匹配接口返回的JSON层级
type FamilyGroupResponse struct {
	Response FamilyGroup `json:"response"`
}
type FamilyGroup struct {
	Name                         string              `json:"name"`                            // 家庭组名称
	Members                      []FamilyGroupMember `json:"members"`                         // 家庭组成员列表
	FreeSpots                    int                 `json:"free_spots"`                      // 家庭组剩余可用名额
	Country                      string              `json:"country"`                         // 家庭组所属地区
	SlotCooldownRemainingSeconds int64               `json:"slot_cooldown_remaining_seconds"` // 名额冷却剩余秒数
	SlotCooldownOverrides        int                 `json:"slot_cooldown_overrides"`         // 名额冷却覆盖次数
}

// ==================== GetFamilyGroupForUser ====================

// FamilyGroupMembershipHistory 家庭组成员身份历史记录结构体
type FamilyGroupMembershipHistory struct {
	FamilyGroupID string `json:"family_groupid"` // 家庭组ID
	RTimeJoined   int64  `json:"rtime_joined"`   // 加入该家庭组的时间戳
	RTimeLeft     int64  `json:"rtime_left"`     // 离开该家庭组的时间戳
	Role          int    `json:"role"`           // 该历史阶段的成员角色
	Participated  bool   `json:"participated"`   // 是否实际参与该家庭组
}

// FamilyGroupForUserResponse IFamilyGroupsService/GetFamilyGroupForUser
type FamilyGroupForUserResponse struct {
	Response struct {
		FamilyGroupID             string `json:"family_groupid"`               // 当前所属家庭组ID
		IsNotMemberOfAnyGroup     bool   `json:"is_not_member_of_any_group"`   // 是否不属于任何家庭组
		LatestTimeJoined          int64  `json:"latest_time_joined"`           // 最近一次加入家庭组的时间戳
		LatestJoinedFamilyGroupID string `json:"latest_joined_family_groupid"` // 最近加入的家庭组ID
		Role                      int    `json:"role"`                         // 当前用户在家庭组中的角色
		CooldownSecondsRemaining  int64  `json:"cooldown_seconds_remaining"`   // 当前用户剩余冷却秒数
		FamilyGroup               struct { // 所属家庭组的完整信息
			Name                         string              `json:"name"`                            // 家庭组名称
			Members                      []FamilyGroupMember `json:"members"`                         // 家庭组成员列表
			FreeSpots                    int                 `json:"free_spots"`                      // 家庭组剩余可用名额
			Country                      string              `json:"country"`                         // 家庭组所属地区
			SlotCooldownRemainingSeconds int64               `json:"slot_cooldown_remaining_seconds"` // 名额冷却剩余秒数
			SlotCooldownOverrides        int                 `json:"slot_cooldown_overrides"`         // 名额冷却覆盖次数
		} `json:"family_group"`                                                                             // 家庭组详情
		CanUndeleteLastJoinedFamily bool                           `json:"can_undelete_last_joined_family"` // 是否可恢复最近退出的家庭组
		MembershipHistory           []FamilyGroupMembershipHistory `json:"membership_history"`              // 家庭组成员身份历史记录
	} `json:"response"` // 接口响应主体
}

// ==================== GetPlaytimeSummary ====================

// FamilyGroupPlaytimeEntry 家庭组单条游玩时长记录结构体
type FamilyGroupPlaytimeEntry struct {
	SteamID       string `json:"steamid"`        // 游玩用户的ID
	AppID         uint64 `json:"appid"`          // 游戏/应用的Steam AppID
	FirstPlayed   int64  `json:"first_played"`   // 首次游玩该应用的时间戳
	LatestPlayed  int64  `json:"latest_played"`  // 最近一次游玩该应用的时间戳
	SecondsPlayed int64  `json:"seconds_played"` // 累计游玩时长(秒)
}

type FamilyGroupPlaytimeBrief struct {
	SteamID       string `json:"steamid"`        // 游玩用户的ID
	AppID         uint64 `json:"appid"`          // 游戏/应用的Steam AppID
	FirstPlayed   string `json:"first_played"`   // 首次游玩该应用的时间戳
	LatestPlayed  string `json:"latest_played"`  // 最近一次游玩该应用的时间戳
	SecondsPlayed int64  `json:"seconds_played"` // 累计游玩时长(秒)
}

// FamilyGroupPlaytimeSummaryResponse IFamilyGroupsService/GetPlaytimeSummary
type FamilyGroupPlaytimeSummaryResponse struct {
	Response struct {
		Entries []FamilyGroupPlaytimeEntry `json:"entries"` // 家庭组成员游玩时长记录列表
	} `json:"response"` // 接口响应主体
}

// ==================== GetSharedLibraryApps ====================

// FamilySharedLibraryApp 家庭组共享库单款应用信息结构体
type FamilySharedLibraryApp struct {
	AppID           uint64   `json:"appid"`            // 共享应用的Steam AppID
	OwnerSteamIDs   []string `json:"owner_steamids"`   // 该应用的拥有者ID
	Name            string   `json:"name"`             // 应用名称
	CapsuleFilename string   `json:"capsule_filename"` // 应用封面图文件名
	ImgIconHash     string   `json:"img_icon_hash"`    // 应用图标哈希值
	ExcludeReason   int      `json:"exclude_reason"`   // 排除共享的原因标识(0=可正常共享)
	RtTimeAcquired  int64    `json:"rt_time_acquired"` // 应用获取时间戳
	RtLastPlayed    int64    `json:"rt_last_played"`   // 最近一次游玩时间戳
	RtPlaytime      int64    `json:"rt_playtime"`      // 累计游玩时长(分钟)
	AppType         int      `json:"app_type"`         // 应用类型标识 1=游戏 2=软件
}

type FamilySharedLibraryAppBrief struct {
	OwnerSteamID string           `json:"owner_steamid"` // 共享库主拥有者的ID
	Apps         []SharedAppBrief `json:"apps"`          // 家庭组共享库应用列表
}

type SharedAppBrief struct {
	AppID          uint64   `json:"appid"`            // 共享应用的Steam AppID
	OwnerSteamIDs  []string `json:"owner_steamids"`   // 该应用的拥有者ID
	Name           string   `json:"name"`             // 应用名称
	Cover          string   `json:"cover"`            // 应用封面图文件名
	Icon           string   `json:"icon"`             // 应用图标哈希值
	ExcludeReason  int      `json:"exclude_reason"`   // 排除共享的原因标识(0=可正常共享)
	RtTimeAcquired string   `json:"rt_time_acquired"` // 应用获取时间戳
	RtLastPlayed   string   `json:"rt_last_played"`   // 最近一次游玩时间戳
	RtPlaytime     int64    `json:"rt_playtime"`      // 累计游玩时长(分钟)
	AppType        int      `json:"app_type"`         // 应用类型标识 1=游戏 2=软件
}

// FamilySharedLibraryResponse IFamilyGroupsService/GetSharedLibraryApps
type FamilySharedLibraryResponse struct {
	Response struct {
		OwnerSteamID string                   `json:"owner_steamid"` // 共享库主拥有者的ID
		Apps         []FamilySharedLibraryApp `json:"apps"`          // 家庭组共享库应用列表
	} `json:"response"`
}
