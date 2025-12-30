package models

// GetEquippedProfileItemsResponse ILoyaltyRewardsService/GetEquippedProfileItems/v1
type GetEquippedProfileItemsResponse struct {
	Response struct {
		ActiveDefinitions   []ProfileItemDefinition `json:"active_definitions"`   // 已装备的道具定义列表
		InactiveDefinitions []ProfileItemDefinition `json:"inactive_definitions"` // 未装备的道具定义列表
	} `json:"response"` // 接口响应主体
}

// ProfileItemDefinition 个人资料道具定义
type ProfileItemDefinition struct {
	AppID                 int64             `json:"appid"`                   // 应用ID
	DefID                 int64             `json:"defid"`                   // 道具定义ID
	Type                  int               `json:"type"`                    // 道具类型标识
	CommunityItemClass    int               `json:"community_item_class"`    // 社区道具分类
	CommunityItemType     int               `json:"community_item_type"`     // 社区道具类型
	PointCost             string            `json:"point_cost"`              // 所需积分
	TimestampCreated      int64             `json:"timestamp_created"`       // 创建时间戳
	TimestampUpdated      int64             `json:"timestamp_updated"`       // 更新时间戳
	TimestampAvailable    int64             `json:"timestamp_available"`     // 可用开始时间戳
	TimestampAvailableEnd int64             `json:"timestamp_available_end"` // 可用结束时间戳
	Quantity              string            `json:"quantity"`                // 数量
	InternalDescription   string            `json:"internal_description"`    // 内部描述
	Active                bool              `json:"active"`                  // 是否激活
	CommunityItemData     CommunityItemData `json:"community_item_data"`     // 社区道具详情
	UsableDuration        int               `json:"usable_duration"`         // 可用时长
	BundleDiscount        int               `json:"bundle_discount"`         // 捆绑折扣
}

// CommunityItemData 社区道具详情
type CommunityItemData struct {
	ItemName        string `json:"item_name"`        // 道具名称
	ItemTitle       string `json:"item_title"`       // 道具标题
	ItemDescription string `json:"item_description"` // 道具描述
	ItemImageSmall  string `json:"item_image_small"` // 小尺寸图片哈希
	ItemImageLarge  string `json:"item_image_large"` // 大尺寸图片哈希
	ItemMovieWebm   string `json:"item_movie_webm"`  // WebM 动画哈希
	ItemMovieMp4    string `json:"item_movie_mp4"`   // MP4 动画哈希
	Animated        bool   `json:"animated"`         // 是否动画效果
	Tiled           bool   `json:"tiled"`            // 是否平铺
}

// ProfileItemBriefInfo 个人资料道具精简信息
type ProfileItemBriefInfo struct {
	ID            int64  `json:"id"`              // 道具定义ID
	AppID         int64  `json:"app_id"`          // 所属应用ID
	Name          string `json:"name"`            // 道具名称
	Title         string `json:"title"`           // 道具标题
	Description   string `json:"description"`     // 道具描述
	PointCost     string `json:"point_cost"`      // 所需积分
	IsActive      bool   `json:"is_active"`       // 是否装备中
	IsAnimated    bool   `json:"is_animated"`     // 是否动画效果
	LargeImageURL string `json:"large_image_url"` // 大尺寸图片完整地址
	SmallImageURL string `json:"small_image_url"` // 小尺寸图片完整地址
	WebmMovieURL  string `json:"webm_movie_url"`  // WebM 动画完整地址
	Mp4MovieURL   string `json:"mp4_movie_url"`   // MP4 动画完整地址
}

// GetReactionsSummaryForUserResponse ILoyaltyRewardsService/GetReactionsSummaryForUser/v1
type GetReactionsSummaryForUserResponse struct {
	Response struct {
		Total               []ReactionSummaryItem `json:"total"`                 // 所有互动类型汇总
		UserReviews         []ReactionSummaryItem `json:"user_reviews"`          // 用户评论互动汇总
		UGC                 []ReactionSummaryItem `json:"ugc"`                   // UGC内容互动汇总
		Profile             []ReactionSummaryItem `json:"profile"`               // 个人资料互动汇总
		TotalGiven          int64                 `json:"total_given"`           // 总发出互动数
		TotalReceived       int64                 `json:"total_received"`        // 总收到互动数
		TotalPointsGiven    string                `json:"total_points_given"`    // 发出互动总积分
		TotalPointsReceived string                `json:"total_points_received"` // 收到互动总积分
	} `json:"response"` // 接口响应主体
}

// ReactionSummaryItem 互动汇总项
type ReactionSummaryItem struct {
	ReactionID     int64  `json:"reactionid"`      // 互动类型ID
	Given          int64  `json:"given"`           // 发出该类型互动数
	Received       int64  `json:"received"`        // 收到该类型互动数
	PointsGiven    string `json:"points_given"`    // 发出该类型互动获得积分
	PointsReceived string `json:"points_received"` // 收到该类型互动获得积分
}

// ReactionSummaryBriefInfo 互动汇总精简信息
type ReactionSummaryBriefInfo struct {
	ReactionID     int64  `json:"reaction_id"`     // 互动类型ID
	Given          int64  `json:"given"`           // 发出该类型互动数
	Received       int64  `json:"received"`        // 收到该类型互动数
	PointsGiven    string `json:"points_given"`    // 发出该类型互动获得积分
	PointsReceived string `json:"points_received"` // 收到该类型互动获得积分
	IconURL        string `json:"icon_url"`        // 互动类型图标地址
}

// UserReactionsTotalBrief 互动汇总总览精简信息
type UserReactionsTotalBrief struct {
	TotalGiven          int64  `json:"total_given"`           // 总发出互动数
	TotalReceived       int64  `json:"total_received"`        // 总收到互动数
	TotalPointsGiven    string `json:"total_points_given"`    // 发出互动总积分
	TotalPointsReceived string `json:"total_points_received"` // 收到互动总积分
	// 各分类互动精简列表
	Total       []ReactionSummaryBriefInfo `json:"total"`        // 所有互动类型汇总
	UserReviews []ReactionSummaryBriefInfo `json:"user_reviews"` // 用户评论互动汇总
	UGC         []ReactionSummaryBriefInfo `json:"ugc"`          // UGC内容互动汇总
	Profile     []ReactionSummaryBriefInfo `json:"profile"`      // 个人资料互动汇总
}

// GetLoyaltyRewardsSummaryResponse ILoyaltyRewardsService/GetSummary/v1
type GetLoyaltyRewardsSummaryResponse struct {
	Response struct {
		Summary struct {
			Points       string `json:"points"`        // 当前剩余点数
			PointsEarned string `json:"points_earned"` // 累计获得点数
			PointsSpent  string `json:"points_spent"`  // 累计消耗点数
		} `json:"summary"` // 点数汇总信息
		TimestampUpdated int64  `json:"timestamp_updated"` // 最后更新时间戳
		AuditIDHighwater string `json:"auditid_highwater"` // 审计ID高水位值
	} `json:"response"` // 接口响应主体
}

// LoyaltyRewardsSummaryBriefInfo 忠诚度奖励点数汇总精简信息
type LoyaltyRewardsSummaryBriefInfo struct {
	CurrentPoints string `json:"current_points"` // 当前剩余点数
	EarnedPoints  string `json:"earned_points"`  // 累计获得点数
	SpentPoints   string `json:"spent_points"`   // 累计消耗点数
	UpdatedAt     int64  `json:"updated_at"`     // 最后更新时间戳（秒级）
}
