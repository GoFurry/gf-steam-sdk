package models

// GetAppsResponse ICommunityService/GetApps
type GetAppsResponse struct {
	Response struct {
		Apps []struct {
			AppID                            int64  `json:"appid"`                             // 应用唯一标识
			Name                             string `json:"name"`                              // 应用名称
			Icon                             string `json:"icon"`                              // 应用图标哈希值
			CommunityVisibleStats            bool   `json:"community_visible_stats,omitempty"` // 是否显示社区统计
			Propagation                      string `json:"propagation"`                       // 传播范围
			AppType                          int    `json:"app_type"`                          // 应用类型标识
			ContentDescriptorIDs             []int  `json:"content_descriptorids"`
			ContentDescriptorIDsIncludingDLC []int  `json:"content_descriptorids_including_dlc"`
		} `json:"apps"` // 应用列表
	} `json:"response"` // 接口响应主体
}

// AppBriefInfo 应用信息精简模型
type AppBriefInfo struct {
	ID               int64  `json:"id"`                // 应用ID
	Name             string `json:"name"`              // 应用名称
	Icon             string `json:"icon"`              // 图标地址
	Type             int    `json:"type"`              // 应用类型(游戏/软件) game=1 software=2
	CommunityVisible bool   `json:"community_visible"` // 社区可见性
	Propagation      string `json:"propagation"`       // 传播范围
}
